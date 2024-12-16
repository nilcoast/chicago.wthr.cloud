chicago.wthr.cloud
---

Every 5 hours, post the current conditions to @chicago.wthr.cloud

## Basic Flow

```
Periodic nomad job set to a cron for 7am, 11pm, 3pm, 7pm, 11pm, 3am

Cron wakes up, and gets the current conditions from the wthr.cloud/api
```

<!-- target: wthr -->
```bash
curl -s 'https://wthr.cloud/api/forecast?latitude=41.975844&longitude=-87.6633969' \
  | jq .current.properties.periods[0].detailedForecast
```

<!-- name: wthr -->
```
"Patchy fog before noon. Mostly cloudy. High near 50, with temperatures falling to around 43 in the afternoon. West southwest wind 10 to 15 mph, with gusts as high as 30 mph. New rainfall amounts less than a tenth of an inch possible."
```

Gets summary from LLM


<!-- target: llm -->
```bash
curl -s -H "Authorization: Bearer ${OLLAMA_KEY}" https://ollama.home.benoist.dev/v1/chat/completions -d '{"model": "llama3.2", "messages": [{"role": "user", "content": "Reformat the following Chicago weather report from 12/16/2024 @ 9am as a tweet less than 240 characters. Please use emoji. Do not make up anything. Do not editorialize. Do not add any hashtags. Always end the message with a link to https://chicago.wthr.cloud: Patchy fog before noon. Mostly cloudy. High near 50, with temperatures falling to around 43 in the afternoon. West southwest wind 10 to 15 mph, with gusts as high as 30 mph. New rainfall amounts less than a tenth of an inch possible."}] }' | jq .choices[0].message.content
```

<!-- name: llm -->
```
"\"Morning weather alert 🛋️🌫️: Patchy fog before noon, mostly cloudy, high near 50, falling to 43 in the afternoon. Winds WSW 10-15mph with gusts up to 30mph. <https://chicago.wthr.cloud>"
```

Post result to bluesky account @chicago.wthr.cloud

```python
from atproto import Client
import os

client = Client()
client.login('chicago.wthr.cloud', os.environ['CHICAGO_WTHR_BSKY_PASS'])
post = client.send_post('\"Morning weather alert 🛋️🌫️: Patchy fog before noon, mostly cloudy, high near 50, falling to 43 in the afternoon. Winds WSW 10-15mph with gusts up to 30mph. <https://chicago.wthr.cloud>')
```
