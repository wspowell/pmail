package resources

type MailboxAttributes struct{}

type MailboxStore interface {
	CreateMailbox(attributes MailboxAttributes) (uint, error)
	CheckMailbox(mailboxId uint) error
	UpdateUser(mailboxId uint, newAttributes MailboxAttributes) error
	DeleteMailbox(mailboxId uint) error
}
