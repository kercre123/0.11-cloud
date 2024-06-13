# 0.11-cloud

This is almost-fully-working cloud software specifically designed for 0.11.19. All support for newer bots is removed.

This is designed for 0.11.19, but older versions should work too.

This is essentially just a bunch of stuff copied from wire-pod. This is more designed to be run as server software rather than software meant to be run on the same network as Vector. So, a bunch of stuff is taken out.

I have gotten it working with 0.11.19 on my public server. I'll release firmware soon which uses this.

I will probably implement a public web interface. Anyone can register their bot's ESN with a password (once it's implemented). Weather API will be handled by me, but you will be able to enter your knowledge graph credentials.

## Env vars

```
	CertFileEnv      = "TLS_CERT_PATH"
	KeyFileEnv       = "TLS_KEY_PATH"
	WeatherAPIKeyEnv = "WEATHER_API_KEY"
	VoskModelPathEnv = "VOSK_MODEL_PATH"
	ChipperPortEnv   = "CHIPPER_PORT"
	WebPortEnv       = "WEB_PORT"
```