# single-digits

The hunt for single digits is on.

Usage:

- Install [tesseract-ocr](https://github.com/tesseract-ocr/tessdoc/blob/main/Installation.md#linux)
- export `DISCORD_CHANNEL_ID` and `DISCORD_BOT_TOKEN` environment variables
- `go run main.go` to run main function
- `go run utils/generate_test_string.go -- testdata/<FILE>` to write the detected text from a file to standard output
- `go test ./...` to run unit tests