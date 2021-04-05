package resources

type MailAttributes struct{}

type MailStore interface {
	SendMail(attributes MailAttributes) (uint, error)
	TrackMail(mailId uint) error
	CollectMail(mailId uint) error
}
