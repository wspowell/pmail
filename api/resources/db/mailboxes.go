package db

import (
	"math/rand"

	"github.com/wspowell/errors"
	"github.com/wspowell/pmail/resources"
)

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

func (self *Mailboxes) CreateMailbox(attributes resources.MailboxAttributes) (uint32, error) {
	mailboxId := rand.Uint32()

	self.mailboxIdToMailbox[mailboxId] = resources.Mailbox{
		Attributes: attributes,
	}

	return mailboxId, nil
}
func (self *Mailboxes) SetHomeMailbox(mailboxId uint32, userId uint32) error {
	if _, exists := self.mailboxIdToUserId[mailboxId]; exists {
		return errors.Wrap(icMailboxesMailboxOwned, resources.ErrMailboxOwned)
	}
	if _, exists := self.userIdToMailboxId[userId]; exists {
		return errors.Wrap(icMailboxesUserHomeMailboxExists, resources.ErrHomeMailboxExists)
	}

	self.mailboxIdToUserId[mailboxId] = userId
	self.userIdToMailboxId[userId] = mailboxId

	return nil
}
func (self *Mailboxes) RemoveHomeMailbox(mailboxId uint32, userId uint32) error {
	panic("not implemented") // TODO: Implement
}
func (self *Mailboxes) FindNearbyMailboxes(location resources.GeoCoordinate, radius float32) error {
	panic("not implemented") // TODO: Implement
}
