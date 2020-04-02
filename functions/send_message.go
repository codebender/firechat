package functions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// PROJECT_ID is set by the Cloud Function runtime
var projectID = os.Getenv("PROJECT_ID")

// repo is the datastore that is reused between function calls
var repo FirestoreRepository

// SendMessageRequest represents the JSON object sent as the request
type SendMessageRequest struct {
	ConversationID string `json:"conversation_id"`
	Author         string `json:"author"`
	Message        string `json:"message"`
}

// init initializes the Firstore client
// init is ONLY run during the function cold start
func init() {
	if strings.ToLower(os.Getenv("ENV")) == "test" || strings.ToLower(os.Getenv("CI")) == "true" {
		return
	}

	repo = NewRepo(projectID, "")

	log.Println("Project ID", projectID, "Function Initialized Successfully")
}

// SendMessage is the entry point for an http trigger Google Clound Function.
// It parses the request and adds the message to the given conversation in Firestore.
func SendMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions: // CORS - the fun part
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodPost:
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var request SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("Unable to parse JSON, err=%v", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		handleMessageRequest(w, request)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func handleMessageRequest(w http.ResponseWriter, request SendMessageRequest) {
	log.Println("method", "handleMessageRequest", "request.ConversationID", request.ConversationID,
		"request.Author", request.Author,
		"request.Message", request.Message)

	err := validateRequest(request)
	if err != nil {
		log.Println("method", "handleMessageRequest", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	chatMessage := ChatMessage{Author: request.Author, Message: request.Message, Timestamp: time.Now().Unix()}
	err = repo.AddMessage(request.ConversationID, chatMessage)
	if err != nil {
		log.Println("Unable to save Chat Message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func validateRequest(request SendMessageRequest) error {
	switch {
	case request.ConversationID == "":
		return fmt.Errorf("ConversationID is Required")
	case request.Author == "":
		return fmt.Errorf("Author is Required")
	case request.Message == "":
		return fmt.Errorf("Message is Required")
	}

	return nil
}
