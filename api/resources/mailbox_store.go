package resources

type Latitude float32
type Longitude float32

type GeoCoordinate struct {
	Lat Latitude
	Lng Longitude
}

type Mailbox struct {
	MailboxId  uint32
	Attributes MailboxAttributes
}

type MailboxAttributes struct {
	Location GeoCoordinate
}

type MailboxStore interface {
	CreateMailbox(userId uint32, attributes MailboxAttributes) (uint32, error)
	GetMailboxById(mailboxId uint32) (*Mailbox, error)
	GetMailboxByUserId(userId uint32) (*Mailbox, error)
	FindNearbyMailboxes(location GeoCoordinate, radius float32) error
}
