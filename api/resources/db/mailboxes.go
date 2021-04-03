package db

import "github.com/wspowell/pmail/resources"

type Mailboxes struct{}

func (self *Mailboxes) SendMail(attributes resources.MailAttributes) (uint, error) {
	panic("not implemented") // TODO: Implement
}

func (self *Mailboxes) TrackMail(mailId uint) error {
	panic("not implemented") // TODO: Implement
}

func (self *Mailboxes) CollectMail(mailId uint) error {
	panic("not implemented") // TODO: Implement
}
