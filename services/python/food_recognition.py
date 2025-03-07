from openai import OpenAI
client = OpenAI()

response = client.chat.completions.create(
  model="gpt-4o",
  messages=[
    {
      "role": "user",
      "content": [
        {
          "type": "image_url",
          "image_url": {
            "url": "data:image/jpeg;base64,..."
          }
        },
        {
          "type": "text",
          "text": "write me every product on this photo and amount of calories of every product on photo"
        }
      ]
    },
    {
      "role": "assistant",
      "content": [
        {
          "type": "text",
          "text": "The dish in the photo includes:\n\n1. **Buckwheat Groats:**\n   - Amount: It's hard to determine the exact quantity, but let's assume there's about 1 cup.\n   - Calories: Approximately 155 calories per cup cooked.\n\n2. **Chicken Breast:**\n   - Amount: Approximately 4-5 slices, which could be around 150 grams.\n   - Calories: Approximately 165 calories per 100 grams, so around 247 calories for 150 grams.\n\n3. **Shredded Cheese:**\n   - Amount: Judging visually, there might be about 1/4 cup.\n   - Calories: Approximately 110 calories per 1/4 cup, depending on the type of cheese.\n\nThese are rough estimates. The exact calorie count can vary based on specific quantities and preparation methods."
        }
      ]
    }
  ],
  response_format={
    "type": "text"
  },
  temperature=1,
  max_completion_tokens=2048,
  top_p=1,
  frequency_penalty=0,
  presence_penalty=0
)