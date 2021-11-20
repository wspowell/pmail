package db

import (
	"sync"
	"time"

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
	usernameToPassword map[string]string
	userToMailGuids    map[user.Guid][]mail.Guid

	// Mail
	mailMutex          *sync.RWMutex
	mailGuidToMail     map[mail.Guid]mail.Mail
	carrierToMailGuids map[user.Guid][]mail.Guid

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
		usernameToPassword: map[string]string{},
		userToMailGuids:    map[user.Guid][]mail.Guid{},

		// Mail
		mailMutex:          &sync.RWMutex{},
		mailGuidToMail:     map[mail.Guid]mail.Mail{},
		carrierToMailGuids: map[user.Guid][]mail.Guid{},

		// Mailboxes
		mailboxMutex:              &sync.RWMutex{},
		mailboxGuidToMailbox:      map[mailbox.Guid]mailbox.Mailbox{},
		mailboxLabelToMailboxGuid: map[string]mailbox.Guid{},
		userGuidToMailboxGuid:     map[user.Guid]mailbox.Guid{},
		mailboxGuidToMailGuids:    map[mailbox.Guid][]mail.Guid{},
	}
}

func (self *InMemory) CreateUser(ctx context.Context, newUser user.User, password string) error {
	if self.userGuidExists(ctx, newUser.UserGuid) {
		return errors.Propagate(icCreateUserGuidConflict, ErrUserGuidExists)
	}

	if self.usernameExists(ctx, newUser.Attributes.Username) {
		return errors.Propagate(icCreateUserUsernameConflict, ErrUsernameExists)
	}

	self.userMutex.Lock()
	self.userGuidToUser[newUser.UserGuid] = newUser
	self.usernameToUserGuid[newUser.Attributes.Username] = newUser.UserGuid
	self.usernameToPassword[newUser.Username] = password
	self.userMutex.Unlock()

	return nil
}

func (self *InMemory) GetUser(ctx context.Context, userGuid user.Guid) (*user.User, error) {
	self.userMutex.RLock()
	defer self.userMutex.RUnlock()

	if foundUser, exists := self.userGuidToUser[userGuid]; exists {
		return &foundUser, nil
	}

	return nil, errors.Propagate(icGetUserUserNotFound, ErrUserNotFound)
}

func (self *InMemory) GetUserMail(ctx context.Context, userGuid user.Guid) ([]mail.Mail, error) {
	self.userMutex.RLock()
	defer self.userMutex.RUnlock()

	if _, exists := self.userGuidToUser[userGuid]; exists {
		if userMail, exists := self.userToMailGuids[userGuid]; exists {
			mailList := make([]mail.Mail, 0, len(userMail))
			self.mailMutex.RLock()
			defer self.mailMutex.RUnlock()

			for _, mailGuid := range userMail {
				mailList = append(mailList, self.mailGuidToMail[mailGuid])
			}
			return mailList, nil
		}
		return nil, nil
	}

	return nil, errors.Propagate(icGetUserUserNotFound, ErrUserNotFound)
}

func (self *InMemory) AuthUser(ctx context.Context, username string, password string) (*user.User, error) {
	self.userMutex.RLock()
	defer self.userMutex.RUnlock()

	if expectedPassword, exists := self.usernameToPassword[username]; exists && expectedPassword == password {
		if foundUserGuid, exists := self.usernameToUserGuid[username]; exists {
			if foundUser, exists := self.userGuidToUser[foundUserGuid]; exists {
				return &foundUser, nil
			}
		}
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

	self.mailMutex.Lock()
	defer self.mailMutex.Unlock()
	if _, exists := self.carrierToMailGuids[newMail.From]; !exists {
		self.carrierToMailGuids[newMail.From] = []mail.Guid{}
	}

	// FIXME: How to handle carry capacity when creating new mail? Maybe there should be a carry mail space and new mail space in a user mailbag.
	self.carrierToMailGuids[newMail.From] = append(self.carrierToMailGuids[newMail.From], newMail.MailGuid)
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

	if newMailbox.Owner != "" && self.userMailboxExists(ctx, newMailbox.Owner) {
		return errors.Propagate(icCreateMailboxUserMailboxConflict, ErrUserMailboxExists)
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

	if foundMailbox, exists := self.mailboxGuidToMailbox[mailboxGuid]; exists {
		return &foundMailbox, nil
	}

	return nil, errors.Propagate(icGetMailboxGuidNotFound, ErrMailboxNotFound)
}

func (self *InMemory) GetMailboxByLabel(ctx context.Context, mailboxLabel string) (*mailbox.Mailbox, error) {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	if mailboxGuid, exists := self.mailboxLabelToMailboxGuid[mailboxLabel]; exists {
		if foundMailbox, exists := self.mailboxGuidToMailbox[mailboxGuid]; exists {
			return &foundMailbox, nil
		}
	}

	return nil, errors.Propagate(icGetMailboxByLabelLabelNotFound, ErrMailboxNotFound)
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
		if userMailbox, exists := self.mailboxGuidToMailbox[mailboxGuid]; exists {
			return &userMailbox, nil
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

func (self *InMemory) DropOffMail(ctx context.Context, carrierGuid user.Guid, mailboxGuid mailbox.Guid) ([]mail.Guid, error) {
	if !self.userGuidExists(ctx, carrierGuid) {
		return nil, errors.Propagate(icDropOffMailUserNotFound, ErrUserNotFound)
	}

	mailbox, err := self.GetMailbox(ctx, mailboxGuid)
	if err != nil {
		return nil, errors.Propagate(icDropOffMailMailboxNotFound, err)
	}

	self.mailboxMutex.Lock()
	defer self.mailboxMutex.Unlock()

	self.mailMutex.Lock()
	defer self.mailMutex.Unlock()

	updatedCarrierMail := make([]mail.Guid, 0, len(self.carrierToMailGuids[carrierGuid]))
	droppedOffMail := make([]mail.Guid, 0, len(self.carrierToMailGuids[carrierGuid]))
	for _, mailGuid := range self.carrierToMailGuids[carrierGuid] {
		numberOfMailboxMail := len(self.mailboxGuidToMailGuids[mailboxGuid])
		if numberOfMailboxMail >= int(mailbox.Capacity) {
			updatedCarrierMail = append(updatedCarrierMail, mailGuid)
			continue
		}

		if mail, exists := self.mailGuidToMail[mailGuid]; exists {
			if mailbox.Owner != "" && mail.To != mailbox.Owner {
				updatedCarrierMail = append(updatedCarrierMail, mailGuid)
				continue
			}

			mail.Carrier = ""
			if mail.To == mailbox.Owner {
				mail.DeliveredOn = time.Now().UTC()
			}

			self.mailGuidToMail[mailGuid] = mail
			self.mailboxGuidToMailGuids[mailboxGuid] = append(self.mailboxGuidToMailGuids[mailboxGuid], mailGuid)
			droppedOffMail = append(droppedOffMail, mailGuid)
		}
	}

	self.carrierToMailGuids[carrierGuid] = updatedCarrierMail

	return droppedOffMail, nil
}

func (self *InMemory) PickUpMail(ctx context.Context, carrierGuid user.Guid, mailboxGuid mailbox.Guid) ([]mail.Guid, error) {
	if !self.mailboxGuidExists(ctx, mailboxGuid) {
		return nil, errors.Propagate(icDropOffMailMailboxNotFound, ErrMailboxNotFound)
	}

	mailbox, err := self.GetMailbox(ctx, mailboxGuid)
	if err != nil {
		return nil, errors.Propagate(icPickUpMailMailboxNotFound, err)
	}

	self.userMutex.RLock()
	carrierUser, exists := self.userGuidToUser[carrierGuid]
	if !exists {
		return nil, errors.Propagate(icPickUpMailUserNotFound, ErrUserNotFound)
	}
	self.userMutex.RUnlock()

	self.mailboxMutex.Lock()
	defer self.mailboxMutex.Unlock()

	self.mailMutex.Lock()
	defer self.mailMutex.Unlock()

	updatedMailboxMail := make([]mail.Guid, 0, len(self.mailboxGuidToMailGuids[mailboxGuid]))
	pickedUpMail := make([]mail.Guid, 0, len(self.mailboxGuidToMailGuids[mailboxGuid]))
	for _, mailGuid := range self.mailboxGuidToMailGuids[mailboxGuid] {
		if mailbox.Owner == carrierGuid {
			// Picking up mail from owned mailbox.

			_, exists := self.userToMailGuids[mailbox.Owner]
			if !exists {
				self.userToMailGuids[mailbox.Owner] = []mail.Guid{}
			}

			self.userToMailGuids[mailbox.Owner] = append(self.userToMailGuids[mailbox.Owner], mailGuid)

			// Mail picked up from own mailbox does not go into mail bag.
			pickedUpMail = append(pickedUpMail, mailGuid)
		} else {
			// Picking up mail from public mail exchange.

			carrierMail, exists := self.carrierToMailGuids[carrierGuid]
			if !exists {
				self.carrierToMailGuids[carrierGuid] = []mail.Guid{}
			}

			if len(carrierMail) >= int(carrierUser.MailCarryCapacity) {
				updatedMailboxMail = append(updatedMailboxMail, mailGuid)
				continue
			}

			if mail, exists := self.mailGuidToMail[mailGuid]; exists {
				mail.Carrier = carrierGuid
				self.mailGuidToMail[mailGuid] = mail
				self.carrierToMailGuids[carrierGuid] = append(self.carrierToMailGuids[carrierGuid], mailGuid)
				pickedUpMail = append(pickedUpMail, mailGuid)
			}
		}
	}

	self.mailboxGuidToMailGuids[mailboxGuid] = updatedMailboxMail

	return pickedUpMail, nil
}

func (self *InMemory) OpenMail(ctx context.Context, mailGuid mail.Guid, openedAt time.Time) error {
	self.mailMutex.Lock()
	defer self.mailMutex.Unlock()

	if foundMail, exists := self.mailGuidToMail[mailGuid]; exists {
		foundMail.OpenedOn = openedAt
		self.mailGuidToMail[mailGuid] = foundMail
		return nil
	}

	return errors.Propagate(icOpenMailGuidNotFound, ErrMailNotFound)
}

func (self *InMemory) mailboxGuidExists(ctx context.Context, mailboxGuid mailbox.Guid) bool {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	_, exists := self.mailboxGuidToMailbox[mailboxGuid]
	return exists
}

func (self *InMemory) userMailboxExists(ctx context.Context, userGuid user.Guid) bool {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	_, exists := self.userGuidToMailboxGuid[userGuid]
	return exists
}

func (self *InMemory) mailboxLabelExists(ctx context.Context, mailboxLabel string) bool {
	self.mailboxMutex.RLock()
	defer self.mailboxMutex.RUnlock()

	_, exists := self.mailboxLabelToMailboxGuid[mailboxLabel]
	return exists
}
