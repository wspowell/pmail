package db

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/wspowell/snailmail/resources/aws"
	"github.com/wspowell/snailmail/resources/models/geo"
	"github.com/wspowell/snailmail/resources/models/mail"
	"github.com/wspowell/snailmail/resources/models/mailbox"
	"github.com/wspowell/snailmail/resources/models/user"
)

// rdsConnectionInfo is a secrets model in AWS SecretsManager.
type rdsConnectionInfo struct {
	Username string `env:"MYSQL_USERNAME" json:"username" envDefault:"root"`
	Password string `env:"MYSQL_PASSWORD" json:"password" envDefault:"password"`
	Host     string `env:"MYSQL_HOST"     json:"host"     envDefault:"mysql"`
	Port     int    `env:"MYSQL_PORT"     json:"port"     envDefault:"3306"`
}

type MySql struct {
	db               *gorm.DB
	migrationsFolder string
}

func NewMySql() *MySql {
	return &MySql{
		migrationsFolder: "/app/api/resources/db/migrations",
	}
}

func (self *MySql) Connect(ctx context.Context) error {
	var connectionInfo rdsConnectionInfo
	if err := aws.GetSecret(ctx, &connectionInfo); err != nil {
		return errors.Propagate("TODO:", err)
	}

	// See: https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dsn := connectionInfo.Username + ":" + connectionInfo.Password + "@tcp(" + connectionInfo.Host + ":" + strconv.Itoa(connectionInfo.Port) + ")/snailmail?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.Propagate(icConnectFailure, err)
	}

	self.db = db

	return nil
}

func (self *MySql) Migrate() error {
	db, err := self.db.DB()
	if err != nil {
		return errors.Propagate("TODO:", err)
	}

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{})
	if err != nil {
		return errors.Propagate(icMigrateInitError, err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://"+self.migrationsFolder,
		"mysql",
		driver,
	)
	if err != nil {
		return errors.Propagate(icMigrateNewInstanceError, err)
	}

	if err = migrations.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Propagate(icMigrateUpError, err)
	}

	return nil
}

func (self *MySql) CreateUser(ctx context.Context, newUser user.User, password string) error {
	query := `
		INSERT INTO users
		(user_guid, username, pineapple_on_pizza, mailbag_capacity, created_on, secret, salt)
		VALUES
		(?, ?, ?, ?, ?, SHA2(?+SHA2(RAND(?),512),512), SHA2(RAND(?),512))
	`

	randValue, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return errors.Propagate(icCreateUserRandSeedError, err)
	}

	seed := randValue.Uint64()

	result := self.db.Exec(query, string(newUser.UserGuid), newUser.Username, newUser.PineappleOnPizza, newUser.MailCarryCapacity, sqlDate(newUser.CreatedOn), password, seed, seed)
	if result.Error != nil {
		return errors.Propagate("TODO:", ErrInternalFailure)
	}

	return nil
}

func (self *MySql) GetUser(ctx context.Context, userGuid user.Guid) (*user.User, error) {
	type userRow struct {
		UserGuid         string
		Username         string
		PineappleOnPizza bool
		MailbagCapacity  uint32
		CreatedOn        time.Time
	}

	query := `
		SELECT
			user_guid, 
			username,
			pineapple_on_pizza, 
			mailbag_capacity, 
			created_on
		FROM users
		WHERE
			user_guid = ?
	`

	userResult := userRow{}
	result := self.db.Raw(query, string(userGuid)).Scan(&userResult)
	if result.Error != nil {
		return nil, errors.Propagate("TODO:", ErrInternalFailure)
	}

	if result.RowsAffected == 0 {
		return nil, errors.Propagate("TODO:", ErrUserNotFound)
	}

	return &user.User{
		UserGuid: user.Guid(userResult.UserGuid),
		Attributes: user.Attributes{
			Username:          userResult.Username,
			PineappleOnPizza:  userResult.PineappleOnPizza,
			MailCarryCapacity: userResult.MailbagCapacity,
			CreatedOn:         userResult.CreatedOn,
		},
	}, nil
}

func (self *MySql) AuthUser(ctx context.Context, username string, password string) (*user.User, error) {
	type userRow struct {
		UserGuid         string
		Username         string
		PineappleOnPizza bool
		MailbagCapacity  uint32
		CreatedOn        time.Time
	}

	query := `
		SELECT
			u.user_guid, 
			u.username,
			u.pineapple_on_pizza, 
			u.mailbag_capacity, 
			u.created_on
		FROM users u
		WHERE
			u.username = ? AND
			u.secret = SHA2(?+u.salt,512)
	`

	userResult := userRow{}
	result := self.db.Raw(query, username, password).Scan(&userResult)
	if result.Error != nil {
		return nil, errors.Propagate("TODO:", ErrInternalFailure)
	}

	if result.RowsAffected == 0 {
		return nil, errors.Propagate("TODO:", ErrUserNotFound)
	}

	log.Debug(ctx, "user: %+v", userResult)

	return &user.User{
		UserGuid: user.Guid(userResult.UserGuid),
		Attributes: user.Attributes{
			Username:          userResult.Username,
			PineappleOnPizza:  userResult.PineappleOnPizza,
			MailCarryCapacity: userResult.MailbagCapacity,
			CreatedOn:         userResult.CreatedOn,
		},
	}, nil
}

func (self *MySql) DeleteUser(ctx context.Context, userGuid user.Guid) error {
	return nil
}

func (self *MySql) UpdateUser(ctx context.Context, updatedUser user.User) error {
	return nil
}

func (self *MySql) CreateMail(ctx context.Context, newMail mail.Mail) error {
	query := `
		INSERT INTO mail
		(mail_guid, ` + "`from`" + `, ` + "`to`" + `, contents, sent_on, delivered_on, opened_on)
		VALUES
		(?, (SELECT id FROM users WHERE user_guid = ?), (SELECT id FROM users WHERE user_guid = ?), ?, ?, ?, ?)
	`

	result := self.db.Exec(query, string(newMail.MailGuid), string(newMail.From), string(newMail.To), newMail.Contents, sqlDate(newMail.SentOn), sqlDate(newMail.DeliveredOn), sqlDate(newMail.OpenedOn))
	if result.Error != nil {
		return errors.Propagate("TODO:", ErrInternalFailure)
	}

	return nil
}

func (self *MySql) GetMail(ctx context.Context, mailGuid mail.Guid) (*mail.Mail, error) {
	return nil, nil
}

func (self *MySql) GetUserMail(ctx context.Context, userGuid user.Guid) ([]mail.Mail, error) {
	type mailRow struct {
		MailGuid    string
		From        string
		To          string
		Contents    string
		SentOn      time.Time
		DeliveredOn time.Time
		OpenedOn    time.Time
	}

	query := `
		SELECT
			m.mail_guid,
			u2.user_guid AS "from",
			u.user_guid AS "to",
			m.contents,
			m.sent_on,
			m.delivered_on,
			m.opened_on
		FROM users u
		JOIN mail m ON m.` + "`to`" + ` = u.id
		JOIN users u2 ON u2.id = m.` + "`from`" + `
		JOIN user_inbox ui ON ui.mail_id = m.id
		WHERE
			u.user_guid = ?
	`

	mailRows := []mailRow{}
	result := self.db.Raw(query, string(userGuid)).Scan(&mailRows)
	if result.Error != nil {
		return nil, errors.Propagate("TODO:", ErrInternalFailure)
	}

	log.Debug(ctx, "mail: %+v", mailRows)

	userMail := make([]mail.Mail, len(mailRows))
	for rowIndex := range mailRows {
		userMail[rowIndex] = mail.Mail{
			MailGuid: mail.Guid(mailRows[rowIndex].MailGuid),
			Attributes: mail.Attributes{
				From:     user.Guid(mailRows[rowIndex].From),
				To:       user.Guid(mailRows[rowIndex].To),
				Contents: mailRows[rowIndex].Contents,
			},
			SentOn:      mailRows[rowIndex].SentOn,
			DeliveredOn: mailRows[rowIndex].DeliveredOn,
			OpenedOn:    mailRows[rowIndex].OpenedOn,
		}
	}

	return userMail, nil
}

func (self *MySql) DeleteMail(ctx context.Context, mailGuid mail.Guid) error {
	return nil
}

func (self *MySql) CreateMailbox(ctx context.Context, newMailbox mailbox.Mailbox) error {
	query := `
		INSERT INTO mailboxes
		(address, ` + "`owner`" + `, capacity, latitude, longitude)
		VALUES
		(?, (SELECT id FROM users WHERE user_guid = ?), ?, ?, ?)
	`

	result := self.db.Exec(query, newMailbox.Address, string(newMailbox.Owner), newMailbox.Capacity, float32(newMailbox.Location.Lat), float32(newMailbox.Location.Lng))
	if result.Error != nil {
		return errors.Propagate("TODO:", ErrInternalFailure)
	}

	return nil
}

func (self *MySql) GetMailbox(ctx context.Context, mailboxAddress string) (*mailbox.Mailbox, error) {
	type mailboxRow struct {
		Address   string
		Owner     string
		Capacity  uint32
		Latitude  float32
		Longitude float32
	}

	query := `
		SELECT
			m.address,
			u.user_guid AS "owner",
			m.capacity,
			m.latitude,
			m.longitude
		FROM mailboxes m
		JOIN users u ON u.id = m.owner
		WHERE
			address = ?
	`

	mailboxResult := mailboxRow{}
	result := self.db.Raw(query, mailboxAddress).Scan(&mailboxResult)
	if result.Error != nil {
		return nil, errors.Propagate("TODO:", ErrInternalFailure)
	}

	if result.RowsAffected == 0 {
		return nil, errors.Propagate("TODO:", ErrUserNotFound)
	}

	return &mailbox.Mailbox{
		Address: mailboxResult.Address,
		Attributes: mailbox.Attributes{
			Owner:    user.Guid(mailboxResult.Owner),
			Capacity: mailboxResult.Capacity,
			Location: geo.Coordinate{
				Lat: geo.Latitude(mailboxResult.Latitude),
				Lng: geo.Longitude(mailboxResult.Longitude),
			},
		},
	}, nil
}

func (self *MySql) DeleteMailbox(ctx context.Context, mailboxAddress string) error {
	return nil
}

func (self *MySql) GetUserMailbox(ctx context.Context, userGuid user.Guid) (*mailbox.Mailbox, error) {
	return nil, nil
}

func (self *MySql) GetNearbyMailboxes(ctx context.Context, location geo.Coordinate, radiusMeters float32) ([]mailbox.Mailbox, error) {
	type mailboxRow struct {
		Address   string
		Owner     string
		Capacity  uint32
		Latitude  float32
		Longitude float32
	}

	// See: https://martech.zone/calculate-great-circle-distance/
	// TODO: It would be good to add a bounding box to the lat/lng values so that this query is not searching the entire table.
	query := `
	SELECT
		m.address,
		u.user_guid AS "owner",
		m.capacity,
		m.latitude,
		m.longitude
	FROM mailboxes m
	JOIN users u ON u.id = m.owner
	WHERE (((ACOS(SIN((? * PI()/180)) * SIN((m.latitude* PI()/180)) + COS((?* PI()/180)) * COS((m.latitude* PI()/180)) * COS(((? - m.longitude)* PI()/180)))) * 180/ PI()) * 60 * 1.1515* 1.609344) <= ?
	`

	mailboxRows := []mailboxRow{}
	result := self.db.Raw(query, location.Lat, location.Lat, location.Lng, metersToKilometers(radiusMeters)).Scan(&mailboxRows)
	if result.Error != nil {
		return nil, errors.Propagate("TODO:", ErrInternalFailure)
	}

	mailboxes := make([]mailbox.Mailbox, len(mailboxRows))
	for rowIndex := range mailboxRows {
		mailboxes[rowIndex] = mailbox.Mailbox{
			Address: mailboxRows[rowIndex].Address,
			Attributes: mailbox.Attributes{
				Owner:    user.Guid(mailboxRows[rowIndex].Owner),
				Capacity: mailboxRows[rowIndex].Capacity,
				Location: geo.Coordinate{
					Lat: geo.Latitude(mailboxRows[rowIndex].Latitude),
					Lng: geo.Longitude(mailboxRows[rowIndex].Longitude),
				},
			},
		}
	}

	return mailboxes, nil
}

func (self *MySql) GetMailboxMail(ctx context.Context, mailboxAddress string) ([]mail.Mail, error) {
	return nil, nil
}

func (self *MySql) DropOffMail(ctx context.Context, carrierGuid user.Guid, mailboxAddress string) ([]mail.Guid, error) {
	return nil, nil
}

func (self *MySql) PickUpMail(ctx context.Context, carrierGuid user.Guid, mailboxAddress string) ([]mail.Guid, error) {
	return nil, nil
}

func (self *MySql) OpenMail(ctx context.Context, mailGuid mail.Guid, openedAt time.Time) error {
	return nil
}

func metersToKilometers(meters float32) float32 {
	return meters / 1000.0
}

// sqlDate converts a time.Time to a value that can be consumed by SQL.
// Returns a *string so that gorm can either send a string time or NULL.
func sqlDate(timeValue time.Time) *string {
	const sqlDateTimeFormat = "2006-01-02 15:04:05"

	if timeValue.IsZero() {
		return nil
	}

	sqlString := timeValue.Format(sqlDateTimeFormat)

	return &sqlString
}
