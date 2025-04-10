import time
import openai
import base64
import json
import requests
import base64
import os
from openai import OpenAI
import mysql.connector
from mysql.connector import Error

def get_value_from_cfg_file(file_path, key):
    with open(file_path, 'r') as file:
        data = json.load(file)
    api_key = data.get(key)
    if api_key is None:
        raise KeyError("API key not found in the JSON file.")
    
    return api_key

def get_image_ext_base64(image_path):
    _, extension = os.path.splitext(image_path)

    extension = extension.replace('.', '')
    
    with open(image_path, "rb") as image_file:
        encoded_image = base64.b64encode(image_file.read()).decode('utf-8')
    
    return extension, encoded_image

def gpt_request_text_img(key, model, promt, base64_img, img_type):
    client = OpenAI(api_key=key, timeout=10)
    response = client.chat.completions.create(
        model=f"{model}",
        messages=[
            {
            "role": "user",
            "content": [
                {
                "type": "image_url",
                "image_url": {
                    "url": f"data:image/{img_type};base64,{base64_img}"
                }
                },
                {
                "type": "text",
                "text": f"{promt}"
                }
            ]
            },
        ],
        response_format={
            "type": "json_object"
        },
        temperature=1,
        max_completion_tokens=2048,
        top_p=1,
        frequency_penalty=0,
        presence_penalty=0
    )

    return response.choices[0].message.content

api_key = get_value_from_cfg_file("../../credentials.json", "openai_api_key")
db_user = get_value_from_cfg_file("../../credentials.json", "db_user")
db_pass = get_value_from_cfg_file("../../credentials.json", "db_pass")

promt = ("""
Write me a report in json format about the food on the photo. Provide name, grams, calories per each product.

Example: 
{"products": [
    {"name": "Buckwheat", "grams": 150, "calories": 164},
    {"name": "Chicken Breast", "grams": 200, "calories": 330},
    {"name": "Cheese", "grams": 50, "calories": 200}
]}
""")

def process_requests():
    try:
        connection = mysql.connector.connect(
            host='127.0.0.1',
            user=db_user,
            password=db_pass,
            database='food_ai'
        )

        if connection.is_connected():
            cursor = connection.cursor(dictionary=True)
            
            print("Successfully connected to the database")

            query = "select ID, ImgBase64, ImgType FROM FoodRecognitionRequests where Status = 1 limit 100"
            cursor.execute(query)
            rows = cursor.fetchall()

            for row in rows:
                print(f"Updating: {row['ID']}")
                query = "update FoodRecognitionRequests set Status = 2 where id = %s"
                values = (row['ID'],)
                cursor.execute(query, values)
                connection.commit()

                res_json = "{}"
                is_error = False
                try:
                    if row["ImgBase64"] == "" or row["ImgType"] == "":
                        res_json = "{}"
                        is_error = True
                    else:
                        res_json = gpt_request_text_img(api_key, "gpt-4o", promt, row["ImgBase64"], row["ImgType"])
                        is_error = False
                        print(f"Sucess: {row['ID']}")
                except Exception as e:
                    res_json = "{}"
                    is_error = True
                    print(f"Error: {row['ID']}: {e}")

                query = "update FoodRecognitionRequests set Response = %s, Status = %s where id = %s"
                if is_error:
                    values = (res_json, 3, row['ID'],)
                else:
                    values = (res_json, 4, row['ID'],)
                cursor.execute(query, values)
                connection.commit()

    except Exception as e:
        print(f"Error: {e}")

    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()
            print("MySQL connection is closed")


while True:
    time.sleep(1)
    process_requests()