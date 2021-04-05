package resources

type Latitude float32
type Longitude float32

type GeoCoordinate struct {
	Lat Latitude
	Lng Longitude
}

type MailboxType uint8

const (
	// Home of a user. Only the owner can dropoff mail.
	Home = MailboxType(1)
	// Dropoff mail only. Public dropoff location.
	Dropoff = MailboxType(2)
	// Exchange mail by dropping off and picking up. Public location.
	Exchange = MailboxType(3)
)

type Mailbox struct {
	Attributes MailboxAttributes
}

type MailboxAttributes struct {
	Type     MailboxType
	Location GeoCoordinate
}

type MailboxStore interface {
	CreateMailbox(attributes MailboxAttributes) (uint32, error)
	SetHomeMailbox(mailboxId uint32, userId uint32) error
	RemoveHomeMailbox(mailboxId uint32, userId uint32) error
	FindNearbyMailboxes(location GeoCoordinate, radius float32) error
}
