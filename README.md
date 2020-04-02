# Firechat
A simple chat for friends

## Dependencies
- Firebase Project: https://console.firebase.google.com/
- Go 1.13: `brew install go`

## Cloud Functions
- Dev
  - Install Go dependencies: `go mod vendor`
  - Run tests: `ENV=test go test`

- Deployment using  gCloud SDK
  - Download & install: https://cloud.google.com/sdk/docs/quickstart-macos
  - Authenticate `gcloud auth login`
  - Set Project ID: `gcloud config set project <YOUR PROJECT ID HERE>`
  - Deploy: `gcloud functions deploy SendMessage --runtime go113 --trigger-http --allow-unauthenticated`

