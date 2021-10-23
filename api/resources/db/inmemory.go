package db

import (
	"sync"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"

	"github.com/wspowell/snailmail/resources/models/geo"
	"github.com/wspowell/snailmail/resources/models/mail"
	"github.com/wspowell/snailmail/resources/models/mailbox"
	"github.com/wspowell/snailmail/resources/models/user"
)

var _ Datastore = (*InMemory)(nil)

type InMemory struct {
	// Users
	userMutex          *sync.RWMutex
	userGuidToUser     map[user.Guid]user.User
	usernameToUserGuid map[string]user.Guid

	// Mail
	mailMutex      *sync.RWMutex
	mailGuidToMail map[mail.Guid]mail.Mail

	// Mailboxes
	mailboxMutex              *sync.RWMutex
	mailboxGuidToMailbox      map[mailbox.Guid]mailbox.Mailbox
	mailboxLabelToMailboxGuid map[string]mailbox.Guid
	userGuidToMailboxGuid     map[user.Guid]mailbox.Guid
	mailboxGuidToMailGuids    map[mailbox.Guid][]mail.Guid
}

func NewInMemory() *InMemory {
	return &InMemory{
		// Users
		userMutex:          &sync.RWMutex{},
		userGuidToUser:     map[user.Guid]user.User{},
		usernameToUserGuid: map[string]user.Guid{},

		// Mail
		mailMutex:      &sync.RWMutex{},
		mailGuidToMail: map[mail.Guid]mail.Mail{},

		// Mailboxes
		mailboxMutex:              &sync.RWMutex{},
		mailboxGuidToMailbox:      map[mailbox.Guid]mailbox.Mailbox{},
		mailboxLabelToMailboxGuid: map[string]mailbox.Guid{},
		userGuidToMailboxGuid:     map[user.Guid]mailbox.Guid{},
		mailboxGuidToMailGuids:    map[mailbox.Guid][]mail.Guid{},
	}
}

func (self *InMemory) CreateUser(ctx context.Context, newUser user.User) error {
	if self.userGuidExists(ctx, newUser.UserGuid) {
		return errors.Propagate(icCreateUserGuidConflict, ErrUserGuidExists)
	}

	if self.usernameExists(ctx, newUser.Attributes.Username) {
		return errors.Propagate(icCreateUserUsernameConflict, ErrUsernameExists)
	}

	self.userMutex.Lock()
	self.userGuidToUser[newUser.UserGuid] = newUser
	self.usernameToUserGuid[newUser.Attributes.Username] = newUser.UserGuid
	self.userMutex.Unlock()

	return nil
}

func (self *InMemory) GetUser(ctx context.Context, userGuid user.Guid) (*user.User, error) {
	self.userMutex.RLock()
	defer self.userMutex.RUnlock()

	if user, exists := self.userGuidToUser[userGuid]; exists {
		return &user, nil
	}

	return nil, errors.Propagate(icGetUserUserNotFound, ErrUserNotFound)
}

func (self *InMemory) DeleteUser(ctx context.Context, userGuid user.Guid) error {
	self.userMutex.Lock()
	defer self.userMutex.Unlock()

	if userToDelete, exists := self.userGuidToUser[userGuid]; exists {
		delete(self.userGuidToUser, userGuid)
		delete(self.usernameToUserGuid, userToDelete.Username)
	}

	return nil
}

func (self *InMemory) UpdateUser(ctx context.Context, updatedUser user.User) error {
	self.userMutex.Lock()
	defer self.userMutex.Unlock()

	if userToUpdate, exists := self.userGuidToUser[updatedUser.UserGuid]; exists {
		self.userGuidToUser[userToUpdate.UserGuid] = updatedUser
		return nil
	}

	return errors.Propagate(icUpdateUserUserNotFound, ErrUserNotFound)
}

func (self *InMemory) userGuidExists(ctx context.Context, userGuid user.Guid) bool {
	self.userMutex.RLock()
	defer self.userMutex.RUnlock()

	_, exists := self.userGuidToUser[userGuid]
	return exists
}

func (self *InMemory) usernameExists(ctx context.Context, username string) bool {
	self.userMutex.RLock()
	defer self.userMutex.RUnlock()

	_, exists := self.usernameToUserGuid[username]
	return exists
}

func (self *InMemory) CreateMail(ctx context.Context, newMail mail.Mail) error {
	if self.mailGuidExists(ctx, newMail.MailGuid) {
		return errors.Propagate(icCreateMailGuidConflict, ErrMailGuidExists)
	}

	self.mailMutex.Lock()
	self.mailGuidToMail[mail.Guid(newMail.MailGuid)] = newMail
	self.mailMutex.Unlock()

	return nil
}

func (self *InMemory) GetMail(ctx context.Context, mailGuid mail.Guid) (*mail.Mail, error) {
	self.mailMutex.RLock()
	defer self.mailMutex.RUnlock()

	if mail, exists := self.mailGuidToMail[mailGuid]; exists {
		return &mail, nil
	}

	return nil, errors.Propagate(icGetMailGuidNotFound, ErrMailNotFound)
}

func (self *InMemory) DeleteMail(ctx context.Context, mailGuid mail.Guid) error {
	self.mailMutex.Lock()
	defer self.mailMutex.Unlock()

	delete(self.mailGuidToMail, mailGuid)

	return nil
}

func (self *InMemory) mailGuidExists(ctx context.Context, mailGuid mail.Guid) bool {
	self.mailMutex.RLock()
	defer self.mailMutex.RUnlock()

	_, exists := self.mailGuidToMail[mailGuid]
	return exists
}

func (self *InMemory) CreateMailbox(ctx context.Context, newMailbox mailbox.Mailbox) error {
	if self.mailboxGuidExists(ctx, newMailbox.MailboxGuid) {
		return errors.Propagate(icCreateMailboxGuidConflict, ErrMailboxGuidExists)
	}

	if self.mailboxLabelExists(ctx, newMailbox.Label) {
		return errors.Propagate(icCreateMailboxLabelConflict, ErrMailboxLabelExists)
	}

	self.mailboxMutex.Lock()
	self.mailboxGuidToMailbox[newMailbox.MailboxGuid] = newMailbox
	self.mailboxLabelToMailboxGuid[newMailbox.Label] = newMailbox.MailboxGuid
	if newMailbox.Owner != "" {
		self.userGuidToMailboxGuid[newMailbox.Owner] = newMailbox.MailboxGuid
	}
	self.mailboxGuidToMailGuids[newMailbox.MailboxGuid] = []mail.Guid{}
	self.mailboxMutex.Unlock()

	return nil
}

func (self *InMemory) GetMailbox(ctx context.Context, mailboxGuid mailbox.Guid) (*mailbox.Mailbox, error) {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	if mail, exists := self.mailboxGuidToMailbox[mailboxGuid]; exists {
		return &mail, nil
	}

	return nil, errors.Propagate(icGetMailboxGuidNotFound, ErrMailboxNotFound)
}

func (self *InMemory) DeleteMailbox(ctx context.Context, mailboxGuid mailbox.Guid) error {
	self.mailboxMutex.Lock()
	defer self.mailboxMutex.Unlock()

	if mailboxToDelete, exists := self.mailboxGuidToMailbox[mailboxGuid]; exists {
		delete(self.mailboxGuidToMailbox, mailboxGuid)
		delete(self.mailboxLabelToMailboxGuid, mailboxToDelete.Label)
		delete(self.userGuidToMailboxGuid, mailboxToDelete.Owner)
		delete(self.mailboxGuidToMailGuids, mailboxGuid)
	}

	return nil
}

func (self *InMemory) GetUserMailbox(ctx context.Context, userGuid user.Guid) (*mailbox.Mailbox, error) {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	if mailboxGuid, exists := self.userGuidToMailboxGuid[userGuid]; exists {
		if mail, exists := self.mailboxGuidToMailbox[mailboxGuid]; exists {
			return &mail, nil
		}
		return nil, errors.Propagate(icGetUserMailboxGuidNotFound, ErrMailboxNotFound)
	}

	return nil, errors.Propagate(icGetUserMailboxUserMailboxNotFound, ErrMailboxNotFound)
}

func (self *InMemory) GetNearbyMailboxes(ctx context.Context, location geo.Coordinate, radiusMeters uint32) ([]mailbox.Mailbox, error) {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	nearbyMailboxes := []mailbox.Mailbox{}

	for _, mailbox := range self.mailboxGuidToMailbox {
		if mailbox.IsNearby(location, radiusMeters) {
			nearbyMailboxes = append(nearbyMailboxes, mailbox)
		}
	}

	return nearbyMailboxes, nil
}

func (self *InMemory) GetMailboxMail(ctx context.Context, mailboxGuid mailbox.Guid) ([]mail.Mail, error) {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	if mailbox, exists := self.mailboxGuidToMailbox[mailboxGuid]; exists {
		mailboxMail := make([]mail.Mail, 0, mailbox.Capacity)

		for _, mailGuid := range self.mailboxGuidToMailGuids[mailboxGuid] {
			mailboxMail = append(mailboxMail, self.mailGuidToMail[mailGuid])
		}

		return mailboxMail, nil
	}

	return nil, errors.Propagate(icGetMailboxMailMailboxNotFound, ErrMailboxNotFound)
}

func (self *InMemory) DropOffMail(ctx context.Context, mailboxGuid mailbox.Guid, mailGuids []mail.Guid) error {
	mailbox, err := self.GetMailbox(ctx, mailboxGuid)
	if err != nil {
		return errors.Propagate(icDropOffMailMailboxNotFound, err)
	}

	self.mailboxMutex.Lock()
	defer self.mailboxMutex.Unlock()

	numberOfMail := len(self.mailboxGuidToMailGuids[mailboxGuid])
	if numberOfMail >= int(mailbox.Capacity) {
		return errors.Propagate(icDropOffMailMailboxFull, ErrMailboxFull)
	}

	// FIXME: Needs to check for the mail carrier.

	self.mailboxGuidToMailGuids[mailboxGuid] = append(self.mailboxGuidToMailGuids[mailboxGuid], mailGuids...)

	return nil
}

func (self *InMemory) PickUpMail(ctx context.Context, carrier user.Guid, mailboxGuid mailbox.Guid, mailGuids []mail.Guid) error {
	if !self.mailboxGuidExists(ctx, mailboxGuid) {
		return errors.Propagate(icDropOffMailMailboxNotFound, ErrMailboxNotFound)
	}

	self.mailboxMutex.Lock()
	defer self.mailboxMutex.Unlock()

	for _, mailGuid := range self.mailboxGuidToMailGuids[mailboxGuid] {
		self.mailMutex.Lock()
		if mail, exists := self.mailGuidToMail[mailGuid]; exists {
			mail.Carrier = carrier
			self.mailGuidToMail[mailGuid] = mail
		}
		self.mailMutex.Unlock()
	}

	return nil
}

func (self *InMemory) mailboxGuidExists(ctx context.Context, mailboxGuid mailbox.Guid) bool {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	_, exists := self.mailboxGuidToMailbox[mailboxGuid]
	return exists
}

func (self *InMemory) mailboxLabelExists(ctx context.Context, mailboxLabel string) bool {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	_, exists := self.mailboxLabelToMailboxGuid[mailboxLabel]
	return exists
}
