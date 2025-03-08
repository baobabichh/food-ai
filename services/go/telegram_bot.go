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

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		if update.Message != nil {
			if update.Message.Photo != nil {
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

				// Here, modify to a GET request with parameters
				endpoint := "http://localhost:5555/createRequest"
				params := url.Values{}
				params.Add("uuid", uuid.(string))
				params.Add("img_base64", base64Str)
				params.Add("img_type", mimeType)

				urlWithParams := fmt.Sprintf("%s?%s", endpoint, params.Encode())
				getReq, err := http.NewRequest("GET", urlWithParams, nil)
				if err != nil {
					log.Println("Error creating GET request:", err)
					continue
				}

				client := &http.Client{}
				resp, err = client.Do(getReq)
				if err != nil {
					log.Println("Error making GET request:", err)
					continue
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("Error reading response body:", err)
					continue
				}

				log.Printf("Response from server: %s", body)

				responseMsg := fmt.Sprintf("Received photo: \nMIME Type: %s\nBase64 length: %d", mimeType, len(base64Str))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseMsg)
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "That is not an image!")
				bot.Send(msg)
			}
		}
	}
}