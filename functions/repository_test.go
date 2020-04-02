package functions_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/codebender/firechat/functions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}

var TestFirestoreClient *firestore.Client

var _ = BeforeSuite(createFSTestClient)
var _ = AfterSuite(closeFSTestClient)

const serviceAccountKey = "./serviceAccountKey.json"
const projectID = "codebender-12e17"

func createFSTestClient() {
	opt := option.WithCredentialsFile(serviceAccountKey)
	ctx := context.Background()
	TestFirestoreClient, _ = firestore.NewClient(ctx, projectID, opt)
}

func closeFSTestClient() {
	TestFirestoreClient.Close()
}

var _ = Describe("Repository Client", func() {
	Describe("NewRepo", func() {
		It("initializes an instance of a FirestoreRepository", func() {
			functions.NewRepo(projectID, serviceAccountKey) // fails if not configured correctly
		})
	})

	Describe("AddMessage", func() {
		var (
			repo              functions.FirestoreRepository
			chatMessage       functions.ChatMessage
			conversationID    = "test123"
			author            = "TestAuthor"
			message           = "Hello World?"
			timestamp         = time.Now().Unix()
			expectedConvoPath = fmt.Sprintf("conversations/%s/messages", conversationID)
		)

		BeforeEach(func() {
			repo = functions.NewRepo(projectID, serviceAccountKey)
			chatMessage = functions.ChatMessage{Author: author, Message: message, Timestamp: timestamp}
		})

		AfterEach(func() {
			colRef := TestFirestoreClient.Collection(expectedConvoPath)
			ctx := context.Background()
			docIterator := colRef.DocumentRefs(ctx)
			for {
				docRef, err := docIterator.Next()

				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Println("delete collection error", err)
				}
				docRef.Delete(ctx)
			}
		})

		It("adds the chat message to the given conversation", func() {
			err := repo.AddMessage(conversationID, chatMessage)
			Expect(err).ToNot(HaveOccurred())

			colRef := TestFirestoreClient.Collection(expectedConvoPath)
			ctx := context.Background()
			docSnaps, err := colRef.Documents(ctx).GetAll()
			Expect(err).ToNot(HaveOccurred())
			Expect(docSnaps).To(HaveLen(1))

			var actual functions.ChatMessage
			err = docSnaps[0].DataTo(&actual)
			Expect(err).To(BeNil())

			expectedChatMessage := functions.ChatMessage{
				Author:    author,
				Message:   message,
				Timestamp: timestamp,
			}

			Expect(actual).To(Equal(expectedChatMessage))
		})

		It("returns an error if the conversation ID is invalid", func() {
			err := repo.AddMessage("__INVALID__CONVO_ID__", chatMessage)
			Expect(err).To(HaveOccurred())
		})
	})
})
