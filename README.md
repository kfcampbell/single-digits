# single-digits

The hunt for single digits is on.

Usage:

- Install [tesseract-ocr](https://github.com/tesseract-ocr/tessdoc/blob/main/Installation.md#linux)
- create a `.env` file that looks like the following:

```
DISCORD_BOT_TOKEN=...
DISCORD_SERVER_ID=...
DISCORD_CHANNEL_ID=...
DISCORD_OTHER_BOT_ID=...
ENVIRONMENT=DEV
```

- `source .env`
- `go run main.go` to run main function
- `go run cmd/generate_test_string.go -- testdata/<FILE>` to write the detected text from a file to standard output
- `go test ./...` to run unit tests
- Use the VSCode launch.json task for debugging
