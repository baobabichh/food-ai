import openai
import base64
import json
import requests
import base64
import os
from openai import OpenAI

def get_api_key(file_path):
    with open(file_path, 'r') as file:
        data = json.load(file)
    api_key = data.get('openai_api_key')
    if api_key is None:
        raise KeyError("API key not found in the JSON file.")
    
    return api_key

def get_image_ext_base64(image_path):
    _, extension = os.path.splitext(image_path)
    
    with open(image_path, "rb") as image_file:
        encoded_image = base64.b64encode(image_file.read()).decode('utf-8')
    
    return extension, encoded_image

def gpt_request_text_img(key, model, promt, base64_img, img_type):
    client = OpenAI(api_key=key)
    response = client.chat.completions.create(
        model=f"{model}",
        messages=[
            {
            "role": "user",
            "content": [
                {
                "type": "image_url",
                "image_url": {
                    "url": f"data:image/{img_ext};base64,{img_base64}"
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

img_ext, img_base64 = get_image_ext_base64("/home/goshan9/food-ai/test_data/4.jpg")
api_key = get_api_key("/home/goshan9/food-ai/credentials.json")

promt = ("""
Write me a report in json format about the food on the photo. Provide name, grams, calories per each product.

Example: 
{"products": [
    {"name": "Buckwheat", "grams": 150, "calories": 164},
    {"name": "Chicken Breast", "grams": 200, "calories": 330},
    {"name": "Cheese", "grams": 50, "calories": 200}
]}
""")
try:
    res = gpt_request_text_img(api_key, "gpt-4o", promt, img_base64, img_ext)
except:
    res = "Error"

print(res)