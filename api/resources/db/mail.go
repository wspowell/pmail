package db

import (
	"context"

	"github.com/wspowell/snailmail/resources"
)

var _ resources.MailStore = (*Mails)(nil)

type Mails struct{}

func NewMails() *Mails {
	return &Mails{}
}

func (self *Mails) CreateMail(ctx context.Context, mail resources.Mail) (uint32, error) {
	panic("not implemented") // TODO: Implement
}

func (self *Mails) ReadMail(ctx context.Context, mailId uint32) (*resources.Mail, error) {
	return nil, nil
}

func (self *Mails) TrackMail(ctx context.Context, mailId uint32) error {
	panic("not implemented") // TODO: Implement
}

func (self *Mails) CollectMail(ctx context.Context, mailboxId uint32) ([]uint32, error) {
	return []uint32{}, nil
}

func (self *Mails) DepositMail(ctx context.Context, mailId uint32, mailboxId uint32) error {
	panic("not implemented") // TODO: Implement
}
