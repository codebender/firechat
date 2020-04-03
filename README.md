# Firechat
A simple chat for new friends https://codebender-12e17.web.app/

## Dependencies
- Firebase Project: https://console.firebase.google.com/
- Go 1.13: `brew install go`

## Cloud Function
- Dev
  - change to the function directory `cd functions`
  - Install Go dependencies: `go mod vendor`
  - Run tests: `ENV=test go test`

- Deployment using gCloud SDK
  - Download & install: https://cloud.google.com/sdk/docs/quickstart-macos
  - Authenticate `gcloud auth login`
  - Set Project ID: `gcloud config set project <YOUR PROJECT ID HERE>`
  - Deploy: `gcloud functions deploy SendMessage --runtime go113 --trigger-http --allow-unauthenticated`

## Frontend
- Dev
  - Open `frontend/index.html`

- Deployment using Firebase hosting
  - Download NodeJS `brew install node`
  - Install firebase-tools library: `npm install firebase-tool`
  - Login to firebase: `firebase login`
  - Deploy to Firebase Hosting: `firebase deploy`
