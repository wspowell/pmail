package db

import "github.com/wspowell/pmail/resources"

type Mail struct{}

func (self *Mail) SendMail(attributes resources.MailAttributes) (uint, error) {
	panic("not implemented") // TODO: Implement
}

func (self *Mail) TrackMail(mailId uint) error {
	panic("not implemented") // TODO: Implement
}

func (self *Mail) CollectMail(mailId uint) error {
	panic("not implemented") // TODO: Implement
}
