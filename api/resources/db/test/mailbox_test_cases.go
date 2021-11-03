package test

/*
import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	MailboxStoreTestCases = [...]func(*testing.T, MailboxStore){
		TestCase_CreateMailbox_new,
		// TestCase_CreateUser_username_conflict,
		// TestCase_GetUser_exists,
		// TestCase_GetUser_not_exists,
		// TestCase_DeleteUser_exists,
		// TestCase_DeleteUser_not_exists,
		// TestCase_UpdateUser_exists,
		// TestCase_UpdateUser_not_exists,
	}
)

// RunMailboxStoreTestCases for any given user store implementation.
// Runs all unit tests that must pass regardless of implementation.
func RunMailboxStoreTestCases(t *testing.T, mailboxStore MailboxStore) {
	for _, testCase := range MailboxStoreTestCases {
		t.Run(functionName(testCase), func(t *testing.T) {
			testCase(t, mailboxStore)
		})
	}
}

func TestCase_CreateMailbox_new(t *testing.T, mailboxes MailboxStore) {
	ctx := context.Background()

	mailboxAttributes := MailboxAttributes{
		Label: uuid.New().String(),
	}
	mailbox1, err := mailboxes.CreateMailbox(ctx, mailboxAttributes)
	assert.Nil(t, err)
	// User Attributes should be the same.
	assert.Equal(t, mailbox1.MailboxAttributes, mailboxAttributes)

	mailboxAttributes = MailboxAttributes{
		Label: uuid.New().String(),
	}
	mailbox2, err := mailboxes.CreateMailbox(ctx, mailboxAttributes)
	assert.Nil(t, err)
	// User ID should be some random value other than an auto increment.
	assert.NotEqual(t, 1, mailbox2.Id-mailbox1.Id)
	assert.NotEqual(t, mailbox1.Id, mailbox2.Id)
}
*/
