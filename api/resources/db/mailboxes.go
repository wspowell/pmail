package db

import (
	"context"
	"math/rand"

	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources"
)

var _ resources.MailboxStore = (*Mailboxes)(nil)

type Mailboxes struct {
	mailboxIdToMailbox map[uint32]resources.Mailbox
	userIdToMailboxId  map[uint32]uint32
	mailboxIdToUserId  map[uint32]uint32
}

func NewMailboxes() *Mailboxes {
	return &Mailboxes{
		mailboxIdToMailbox: map[uint32]resources.Mailbox{},
		userIdToMailboxId:  map[uint32]uint32{},
		mailboxIdToUserId:  map[uint32]uint32{},
	}
}

func (self *Mailboxes) CreateMailbox(ctx context.Context, userId uint32, attributes resources.MailboxAttributes) (uint32, error) {
	if _, exists := self.userIdToMailboxId[userId]; exists {
		return 0, errors.Wrap(icMailboxesUserHomeMailboxExists, resources.ErrHomeMailboxExists)
	}

	mailboxId := rand.Uint32()

	self.mailboxIdToMailbox[mailboxId] = resources.Mailbox{
		Id:         mailboxId,
		Attributes: attributes,
	}
	self.mailboxIdToUserId[mailboxId] = userId
	self.userIdToMailboxId[userId] = mailboxId

	return mailboxId, nil
}

func (self *Mailboxes) GetMailboxById(ctx context.Context, mailboxId uint32) (*resources.Mailbox, error) {
	if mailbox, exists := self.mailboxIdToMailbox[mailboxId]; exists {
		return &mailbox, nil
	}

	return nil, errors.Wrap(icMailboxesMailboxNotFound, resources.ErrorMailboxNotFound)
}

func (self *Mailboxes) GetMailboxByUserId(ctx context.Context, userId uint32) (*resources.Mailbox, error) {
	if mailboxId, exists := self.userIdToMailboxId[userId]; exists {
		if mailbox, exists := self.mailboxIdToMailbox[mailboxId]; exists {
			return &mailbox, nil
		}
	}

	return nil, errors.Wrap(icMailboxesMailboxNotFound, resources.ErrorMailboxNotFound)
}

func (self *Mailboxes) FindNearbyMailboxes(ctx context.Context, location resources.GeoCoordinate, radius float32) error {
	panic("not implemented") // TODO: Implement
}
