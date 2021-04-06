package resources

import "time"

type Mail struct {
	From     uint32
	To       uint32
	Contents string

	// Metadata
	SentOn      time.Time
	DeliveredOn time.Time
}

type MailStore interface {
	CreateMail(mail Mail) (uint32, error)
	ReadMail(mailId uint32) (*Mail, error)
	TrackMail(mailId uint32) error
	CollectMail(mailboxId uint32) ([]uint32, error)
	DepositMail(mailId uint32, mailboxId uint32) error
}
