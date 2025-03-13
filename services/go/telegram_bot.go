package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Product struct {
	Name     string `json:"name"`
	Grams    int    `json:"grams"`
	Calories int    `json:"calories"`
}

type ProductsList struct {
	Products []Product `json:"products"`
}

func getValueFromJSONFile(filePath string, key string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	value, exists := data[key]
	if !exists {
		return nil, fmt.Errorf("key %q not found", key)
	}

	return value, nil
}

func parseAndSummarizeProducts(jsonData string) string {
	var productsList ProductsList
	err := json.Unmarshal([]byte(jsonData), &productsList)
	if err != nil {
		log.Fatal(err)
	}

	// Collect results in a string
	var result string

	// Add each product to the result string
	for _, product := range productsList.Products {
		result += fmt.Sprintf("Product: %s, Grams: %d, Calories: %d\n", product.Name, product.Grams, product.Calories)
	}

	// Calculate totals
	totalGrams := 0
	totalCalories := 0
	for _, product := range productsList.Products {
		totalGrams += product.Grams
		totalCalories += product.Calories
	}

	// Add summary to the result string
	result += fmt.Sprintf("Summary:\nTotal Grams: %d, Total Calories: %d\n", totalGrams, totalCalories)

	return result
}

func main() {
	botKey, err := getValueFromJSONFile("../../credentials.json", "tg_bot_key")
	if err != nil {
		log.Panic(err)
	}

	uuid, err := getValueFromJSONFile("../../credentials.json", "uuid_for_tg_bot")
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(botKey.(string))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Photo != nil {
			photo := update.Message.Photo[len(update.Message.Photo)-1] // Get highest resolution photo
			file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
			if err != nil {
				log.Println("Error getting file:", err)
				continue
			}

			fileURL := file.Link(bot.Token)
			resp, err := http.Get(fileURL)
			if err != nil {
				log.Println("Error downloading file:", err)
				continue
			}
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading file data:", err)
				continue
			}

			base64Str := base64.StdEncoding.EncodeToString(data)
			mimeType := http.DetectContentType(data)
			parts := strings.Split(mimeType, "/")
			if len(parts) == 2 {
				mimeType = parts[1]
			}

			/* Step 1: Send to createRequest */
			endpoint := "http://localhost:5555/createRequest"
			form := url.Values{}
			form.Add("uuid", uuid.(string))
			form.Add("img_base64", base64Str)
			form.Add("img_type", mimeType)

			// Create a POST request
			postReq, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
			if err != nil {
				log.Println("Error creating POST request:", err)
				continue
			}
			postReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			client := &http.Client{}
			resp_p, err := client.Do(postReq)
			if err != nil {
				log.Println("Error making POST request:", err)
				continue
			}
			defer resp_p.Body.Close()

			body, err := ioutil.ReadAll(resp_p.Body)
			if err != nil {
				log.Println("Error reading response body:", err)
				continue
			}

			// Assume the createRequest returns a JSON with a request_id
			var createResponse map[string]interface{}
			if err = json.Unmarshal(body, &createResponse); err != nil {
				log.Println("Error decoding createRequest response:", err)
				continue
			}

			requestID, ok := createResponse["RequestID"].(string)
			if !ok {
				log.Println("Invalid request_id in response")
				continue
			}

			is_error := false

			/* Step 2: Check status until complete */
			checkStatusEndpoint := "http://localhost:5555/checkRequestStatus"
			for {
				statusParams := url.Values{}
				statusParams.Add("uuid", uuid.(string))
				statusParams.Add("request_id", requestID)

				statusURL := fmt.Sprintf("%s?%s", checkStatusEndpoint, statusParams.Encode())
				checkStatusReq, err := http.NewRequest("GET", statusURL, nil)
				if err != nil {
					log.Println("Error creating check status request:", err)
					continue
				}

				resp, err = client.Do(checkStatusReq)
				if err != nil {
					log.Println("Error making status check request:", err)
					continue
				}
				defer resp.Body.Close()

				statusBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("Error reading status response body:", err)
					continue
				}

				var statusResponse map[string]interface{}
				if err = json.Unmarshal(statusBody, &statusResponse); err != nil {
					log.Println("Error decoding status response:", err)
					continue
				}

				status, ok := statusResponse["RequestStatus"].(string)
				if !ok {
					log.Println("Invalid status in response")
					continue
				}

				if status == "Success" {
					is_error = false
					break
				} else if status == "Error" {
					is_error = true
					break
				}

				time.Sleep(2 * time.Second) // Wait and then retry
			}

			if is_error {
				log.Println("Error after check status")
				continue
			}

			/* Step 3: Get Request Result */
			getResultEndpoint := "http://localhost:5555/getRequest"
			resultParams := url.Values{}
			resultParams.Add("uuid", uuid.(string))
			resultParams.Add("request_id", requestID)

			resultURL := fmt.Sprintf("%s?%s", getResultEndpoint, resultParams.Encode())
			getResultReq, err := http.NewRequest("GET", resultURL, nil)
			if err != nil {
				log.Println("Error creating get result request:", err)
				continue
			}

			resp, err = client.Do(getResultReq)
			if err != nil {
				log.Println("Error making get result request:", err)
				continue
			}
			defer resp.Body.Close()

			resultBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading result response body:", err)
				continue
			}

			var resultResponse map[string]interface{}
			if err = json.Unmarshal(resultBody, &resultResponse); err != nil {
				log.Println("Error decoding resultBody response:", err)
				continue
			}

			response_t, ok := resultResponse["Response"].(string)
			if !ok {
				log.Println("Invalid Response in response")
				continue
			}

			resultString := parseAndSummarizeProducts(response_t)

			// Ideally, we could do some processing or verification of the result here.
			log.Printf("Received final result: %s", string(resultString))


			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(resultString))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "That is not an image!")
			bot.Send(msg)
		}
	}
}