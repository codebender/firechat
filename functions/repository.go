package functions

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// ChatMessage represents the message stored in a conversation in the datastore
type ChatMessage struct {
	Author    string `json:"author"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// FirestoreRepository is the datastore interface
type FirestoreRepository interface {
	AddMessage(conversationID string, chatMessage ChatMessage) error
}

const messagesCollectionPathTemplate string = "conversations/%s/messages"

// NewRepo initializes a FirestoreRepository
func NewRepo(projectID, credentialFilePath string) FirestoreRepository {
	ctx := context.Background()

	var err error
	var fsClient *firestore.Client
	if credentialFilePath != "" {
		opt := option.WithCredentialsFile(credentialFilePath)
		fsClient, err = firestore.NewClient(ctx, projectID, opt)
	} else {
		fsClient, err = firestore.NewClient(ctx, projectID)
	}

	if err != nil {
		log.Fatalf("Error while configuring Firestore: err=%s", err)
	}

	return firestoreRepository{firestoreClient: fsClient}
}

type firestoreRepository struct {
	firestoreClient *firestore.Client
}

// AddMessage add a new ChatMessage to a conversation.
func (ar firestoreRepository) AddMessage(conversationID string, chatMessage ChatMessage) error {
	messagesCollectionPath := fmt.Sprintf(messagesCollectionPathTemplate, conversationID)
	colRef := ar.firestoreClient.Collection(messagesCollectionPath)
	ctx := context.Background()
	_, _, err := colRef.Add(ctx, chatMessage)

	return err
}
