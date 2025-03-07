import openai
import base64
import json
import requests
import base64
import os

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

from openai import OpenAI
client = OpenAI(api_key=get_api_key("../../credentials.json"))

img_ext, img_base64 = get_image_ext_base64("/home/goshan9/food-ai/test_data/1.jpg")

response = client.chat.completions.create(
  model="gpt-4o",
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
          "text": "how much calories is here? Output in json format as an array of products (name, frams and calories)"
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

print(response)