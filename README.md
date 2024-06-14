# 0.11-cloud

This is fully functional cloud software for Anki Vector robots running 0.11.19 firmware (and other firmware from that era).

This also serves as experimental grounds. If I learn how to better implement something, I will implement it that way here.

It is meant to be deployed on a server. To make a server deployment tar, run `sudo ./deploy-create.sh`

I will be running this on a server soon.

## Webroot

The web interface lets you set the bot's location and lets you set the API credentials for Houndify, OpenAI, or Together.

The web interface only lets you use it once you have said a voice command to the bot. It links your network's public IP with your bot's ESN. Then, you are allowed to go to it and enter your bot's ESN to modify settings.

## Env vars

These must be set for this to run (in source.sh).

```
CertFileEnv      = "TLS_CERT_PATH"
KeyFileEnv       = "TLS_KEY_PATH"
WeatherAPIKeyEnv = "WEATHER_API_KEY"
VoskModelPathEnv = "VOSK_MODEL_PATH"
ChipperPortEnv   = "CHIPPER_PORT"
WebPortEnv       = "WEB_PORT"
```