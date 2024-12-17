chicago.wthr.cloud
---

Every 5 hours, post the current conditions to @chicago.wthr.cloud

## Basic Flow

Periodic nomad job set to a cron for 7am, 11pm, 3pm, 7pm, 11pm, 3am

#### Cron wakes up, and gets the current conditions from the wthr.cloud/api

<!-- target: wthr -->
```bash
curl -s 'https://wthr.cloud/api/forecast?latitude=41.975844&longitude=-87.6633969' \
  | jq -r .current.properties.periods[0].detailedForecast
```

<!-- name: wthr -->
```
Cloudy, with a high near 39. West wind 5 to 10 mph, with gusts as high as 20 mph.
```

#### Gets summary from LLM

<!-- target: llm -->
```bash
curl -s -H "Authorization: Bearer ${OLLAMA_KEY}" https://ollama.home.benoist.dev/v1/chat/completions -d '{"model": "llama3.2", "messages": [{"role": "user", "content": "Reformat the following Chicago weather report from 12/17/2024 @8am as a tweet less than 240 characters. Please use emoji. Do not make up anything. Do not editorialize. Do not add any hashtags. Cloudy, with a high near 39. West wind 5 to 10 mph, with gusts as high as 20 mph."}] }' | jq -r .choices[0].message.content
```

<!-- name: llm -->
```
"8am update: cloudy skies @ 39°F, west winds 5-10mph & gusts up to 20mph 🌫️💨👍 #WeatherUpdate"
```

#### Post result to bluesky account @chicago.wthr.cloud

```python
from atproto import Client
import os

client = Client()
client.login('chicago.wthr.cloud', os.environ['CHICAGO_WTHR_BSKY_PASS'])
post = client.send_post('8am update: cloudy skies @ 39°F, west winds 5-10mph & gusts up to 20mph 🌫️💨👍')
```
