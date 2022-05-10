package db

import (
	gocontext "context"
	"math"
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
	gormlogger "gorm.io/gorm/logger"

	"github.com/wspowell/snailmail/resources/aws"
	"github.com/wspowell/snailmail/resources/models"
)

// SELECT
// 	address,
// 	# (LONG, LAT)
// 	ST_Distance_Sphere(location, ST_SRID(POINT(-96.9073063043077, 33.09389393115629), 4326)) AS distance_m
// FROM mailboxes
// WHERE
//  # FLOOR (long * 10) +/- 1
// 	long_1000_floor IN (-969, -968, -970) AND
//  # FLOOR (lat * 10) +/- 1
// 	lat_1000_floor IN (330, 331, 329) AND
// 	ST_Distance_Sphere(location, ST_SRID(POINT(-96.9073063043077, 33.09389393115629), 4326)) <= 1000
// ;

// https://aaronfrancis.com/2021/efficient-distance-querying-in-my-sql
// https://en.wikipedia.org/wiki/Decimal_degrees#Precision
// https://dba.stackexchange.com/questions/242001/mysql-8-st-geomfromtext-giving-error-latitude-out-of-range-in-function-st-geomfr
// https://stackoverflow.com/questions/7477003/calculating-new-longitude-latitude-from-old-n-meters

// rdsConnectionInfo is a secrets model in AWS SecretsManager.
type rdsConnectionInfo struct {
	Username string `env:"MYSQL_USERNAME" json:"username" envDefault:"root"`
	Password string `env:"MYSQL_PASSWORD" json:"password" envDefault:"password"`
	Host     string `env:"MYSQL_HOST"     json:"host"     envDefault:"mysql"`
	Port     int    `env:"MYSQL_PORT"     json:"port"     envDefault:"3306"`
}

type dbLogger struct {
	logger log.Logger
}

func (self *dbLogger) LogMode(_ gormlogger.LogLevel) gormlogger.Interface {
	return self
}

func (self *dbLogger) Info(ctx gocontext.Context, format string, values ...interface{}) {
	self.logger.Info(format, values...)
}

func (self *dbLogger) Warn(ctx gocontext.Context, format string, values ...interface{}) {
	self.logger.Warn(format, values...)
}

func (self *dbLogger) Error(ctx gocontext.Context, format string, values ...interface{}) {
	self.logger.Error(format, values...)
}

func (self *dbLogger) Trace(ctx gocontext.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// Ignore
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
		return err
	}

	// See: https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dsn := connectionInfo.Username + ":" + connectionInfo.Password + "@tcp(" + connectionInfo.Host + ":" + strconv.Itoa(connectionInfo.Port) + ")/snailmail?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"

	log.Error(ctx, "dsn: %s", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &dbLogger{
			logger: log.NewLog(log.NewConfig().WithLevel(log.LevelInfo)),
		},
	})
	if err != nil {
		return err
	}

	self.db = db

	return nil
}

func (self *MySql) Migrate() error {
	db, err := self.db.DB()
	if err != nil {
		return err
	}

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://"+self.migrationsFolder,
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err = migrations.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (self *MySql) CreateOneTimePassword(ctx context.Context, connectionId string, oneTimePassword string) error {
	query := `
		INSERT INTO one_time_passwords
		(connection_id, otp, created_on)
		VALUES
		(?, ?, ?)
	`

	result := self.db.Exec(query, connectionId, oneTimePassword, sqlDate(time.Now().UTC()))
	if result.Error != nil {
		// FIXME: Check for duplicate errors and return ErrOneTimePasswordExists
		return ErrInternalFailure
	}

	return nil
}

func (self *MySql) GetOneTimePassword(ctx context.Context, oneTimePassword string) (string, error) {
	type otpRow struct {
		ConnectionId string
	}

	query := `
		SELECT
			connection_id
		FROM one_time_passwords
		WHERE
			otp = ?
	`

	connectionResult := otpRow{}
	result := self.db.Raw(query, oneTimePassword).Scan(&connectionResult)
	if result.Error != nil {
		return "", ErrInternalFailure
	}

	if result.RowsAffected == 0 {
		return "", ErrOneTimePasswordNotFound
	}

	return connectionResult.ConnectionId, nil
}

func (self *MySql) DeleteOneTimePassword(ctx context.Context, connectionId string) error {
	query := `
		DELETE
		FROM one_time_passwords
		WHERE
			connection_id = ?
	`

	result := self.db.Exec(query, connectionId)
	if result.Error != nil {
		return ErrInternalFailure
	}

	return nil
}

func (self *MySql) CreateUser(ctx context.Context, newUser models.User) error {
	return self.db.Transaction(func(tx *gorm.DB) error {
		{
			query := `
				INSERT INTO users
				(guid, public_key, signature, created_on)
				VALUES
				(?, ?, ?, ?)
				`

			result := tx.Exec(query, newUser.Guid, newUser.PublicKey, newUser.Signature, sqlDate(newUser.CreatedOn))
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}

		{
			query := `
				INSERT INTO locations
				(address, geo_coordinate, long_1000_floor, lat_1000_floor)
				VALUES
				# (LONG, LAT)
				(?, ST_SRID(POINT(?, ?), 4326), FLOOR(?), FLOOR(?))
				`

			result := tx.Exec(
				query,
				newUser.Mailbox.Address,
				float32(newUser.Mailbox.Location.Longitude),
				float32(newUser.Mailbox.Location.Latitude),
				float32(math.Floor(float64(newUser.Mailbox.Location.Longitude)*1000)),
				float32(math.Floor(float64(newUser.Mailbox.Location.Latitude)*1000)),
			)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}

		{
			query := `
				INSERT INTO mailboxes
				(location_id, user_id)
				VALUES
				(LAST_INSERT_ID(), (SELECT id FROM users WHERE guid = ?))
				`

			result := tx.Exec(query, newUser.Guid)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}

		return nil
	})
}

func (self *MySql) GetUser(ctx context.Context, userGuid string) (*models.User, error) {
	type userRow struct {
		Guid      string
		PublicKey string
		Signature string
		CreatedOn time.Time
	}

	query := `
		SELECT
			u.guid,
			u.public_key,
			u.signature,
			u.created_on
		FROM users u
		WHERE
			u.guid = ?
	`

	userResult := userRow{}
	result := self.db.Raw(query, userGuid).Scan(&userResult)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, ErrInternalFailure)
	}

	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return &models.User{
		Guid:      userResult.Guid,
		PublicKey: userResult.PublicKey,
		Signature: userResult.Signature,
		CreatedOn: userResult.CreatedOn,
	}, nil
}

func (self *MySql) DeleteUser(ctx context.Context, userGuid string) error {
	query := `
		DELETE 
		FROM users
		WHERE
			guid = ?
	`

	result := self.db.Exec(query, userGuid)
	if result.Error != nil {
		return errors.Wrap(result.Error, ErrInternalFailure)
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (self *MySql) UpdateUser(ctx context.Context, updatedUser models.User) error {
	return nil
}

func (self *MySql) CreateMail(ctx context.Context, newMail models.Mail) error {
	query := `
		INSERT INTO mail
			(
				guid,
				sender, 
				` + "`from`" + `, 
				from_address,
				recipient, 
				to_address, 
				` + "`to`" + `, 
				body,
				destination,
				sent_on
			)
		VALUES
			(
				# guid
				?,
				# sender 
				(
					SELECT 
						id 
					FROM users 
					WHERE 
						guid = ?
				), 
				# from
				?, 
				# from_address
				?,
				# recipient
				(
					SELECT
						u.id
					FROM locations l
					JOIN mailboxes m ON m.location_id = l.id
					JOIN users u ON u.id = m.user_id
					WHERE
						l.address = ?
				), 
				# to_address
				?, 
				# to
				?, 
				# body
				?,
				# destination
				(
					SELECT
						l.id
					FROM locations l
					JOIN mailboxes m ON m.location_id = l.id
					WHERE
						l.address = ?
				), 
				# sent_on
				?
		)
	`

	result := self.db.Exec(query,
		newMail.Guid,
		newMail.FromGuid,
		newMail.Contents.From,
		newMail.FromMailboxAddress,
		newMail.ToMailboxAddress,
		newMail.ToMailboxAddress,
		newMail.Contents.To,
		newMail.Contents.Body,
		newMail.ToMailboxAddress,
		sqlDate(newMail.SentOn),
	)
	if result.Error != nil {
		return errors.Wrap(result.Error, ErrInternalFailure)
	}

	return nil
}

func (self *MySql) DeliverMail(ctx context.Context, userGuid string, location models.Coordinate) error {
	return self.db.Transaction(func(tx *gorm.DB) error {
		// The location is checked for any carried mail destinations.
		{
			type mailTrackingRow struct {
				MailId int
			}

			query := `
				SELECT 
					mt.mail_id,
					ST_Distance_Sphere(l.geo_coordinate, ST_SRID(POINT(?, ?), 4326)) AS distance_m
				FROM mail_tracking mt
				JOIN mail m ON m.id = mt.mail_id
				JOIN mailboxes mb ON mb.location_id = m.destination
				JOIN locations l ON l.id = mb.location_id
				WHERE
					# Mail that the user is carrying.
					mt.carrier = (SELECT id FROM users WHERE guid = ?) AND 

					# Mailboxes that are at the location.
					l.long_1000_floor IN (?,?,?) AND
					l.lat_1000_floor IN (?,?,?)
				HAVING
					# Less than 100 meters
					distance_m <= 100
			`

			mailTrackingResult := []mailTrackingRow{}
			result := tx.Raw(
				query,

				location.Longitude,
				location.Latitude,

				userGuid,

				// Check this block and the surrounding blocks in case it is on the edge.
				float32(math.Floor(float64(location.Longitude)*1000)),
				float32(math.Floor(float64(location.Longitude)*1000))+1,
				float32(math.Floor(float64(location.Longitude)*1000))-1,

				float32(math.Floor(float64(location.Latitude)*1000)),
				float32(math.Floor(float64(location.Latitude)*1000))+1,
				float32(math.Floor(float64(location.Latitude)*1000))-1,
			).Scan(&mailTrackingResult)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}

			// Stop if there are no mailboxes.
			if result.RowsAffected == 0 {
				return nil
			}

			for mailIndex := range mailTrackingResult {
				// Remove from mail tracking.
				{
					query := `
						DELETE 
						FROM mail_tracking
						WHERE
							mail_id = ?
					`

					result := tx.Exec(
						query,
						mailTrackingResult[mailIndex].MailId,
					)
					if result.Error != nil {
						return errors.Wrap(result.Error, ErrInternalFailure)
					}
				}

				// Update mail delivery date.
				{
					query := `
						UPDATE mail
						SET 
							delivered_on = CURRENT_TIMESTAMP()
						WHERE
							id = ?
					`

					result := tx.Exec(
						query,
						mailTrackingResult[mailIndex].MailId,
					)
					if result.Error != nil {
						return errors.Wrap(result.Error, ErrInternalFailure)
					}
				}
			}

		}

		return nil
	})
}

func (self *MySql) ExchangeMail(ctx context.Context, userGuid string, location models.Coordinate) error {
	if err := self.DeliverMail(ctx, userGuid, location); err != nil {
		return err
	}

	return self.db.Transaction(func(tx *gorm.DB) error {
		// The location is checked for an exchange.
		var exchangeId int
		var exchangeLocationId int

		{
			type exchangeRow struct {
				ExchangeId int
				LocationId int
				DistanceM  float32
			}

			query := `
				SELECT 
					e.id as exchange_id,
					l.id as location_id,
					# (LONG, LAT)
					ST_Distance_Sphere(l.geo_coordinate, ST_SRID(POINT(?, ?), 4326)) AS distance_m
				FROM locations l
				JOIN exchanges e ON e.location_id = l.id
				WHERE 
					l.long_1000_floor IN (?,?,?) AND
					l.lat_1000_floor IN (?,?,?)
				HAVING
					# Less than 100 meters
					distance_m <= 100
				ORDER BY
					# Select the closest exchange, if multiple.
					distance_m ASC
				LIMIT 1
			`

			exchangeResult := exchangeRow{}
			result := tx.Raw(
				query,

				location.Longitude,
				location.Latitude,

				// Check this block and the surrounding blocks in case it is on the edge.
				float32(math.Floor(float64(location.Longitude)*1000)),
				float32(math.Floor(float64(location.Longitude)*1000))+1,
				float32(math.Floor(float64(location.Longitude)*1000))-1,

				float32(math.Floor(float64(location.Latitude)*1000)),
				float32(math.Floor(float64(location.Latitude)*1000))+1,
				float32(math.Floor(float64(location.Latitude)*1000))-1,
			).Scan(&exchangeResult)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}

			// Stop if there are no exchanges.
			if result.RowsAffected == 0 {
				return nil
			}

			log.Debug(ctx, "found exchange %d %s meters away", exchangeResult.ExchangeId, exchangeResult.DistanceM)
			exchangeId = exchangeResult.ExchangeId
			exchangeLocationId = exchangeResult.LocationId
		}

		// Any new mail the user has created is dropped off
		{
			type mailRow struct {
				Id int
			}

			query := `
				SELECT 
					id
				FROM mail
				WHERE 
					starting_location IS NULL AND
					sender = (SELECT id FROM users WHERE guid = ?)
			`

			newMailResult := []mailRow{}
			result := tx.Raw(
				query,
				userGuid,
			).Scan(&newMailResult)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}

			for mailIndex := range newMailResult {
				mailId := newMailResult[mailIndex]

				// Add to mail tracking.
				{
					query := `
						INSERT INTO mail_tracking
						(mail_id, exchange, last_carrier)
						VALUES
						(?, ?, (SELECT id FROM users WHERE guid = ?))
					`

					result := tx.Exec(
						query,
						mailId,
						exchangeId,
						userGuid,
					)
					if result.Error != nil {
						return errors.Wrap(result.Error, ErrInternalFailure)
					}
				}

				// Update mail starting location.
				{
					query := `
						UPDATE mail
						SET 
							starting_location = ?
						WHERE
							id = ?
					`

					result := tx.Exec(
						query,
						exchangeLocationId,
						mailId,
					)
					if result.Error != nil {
						return errors.Wrap(result.Error, ErrInternalFailure)
					}
				}

				log.Debug(ctx, "dropped off new mail: %d", mailId)
			}

		}

		// Any mail the user is carrying is dropped off
		{
			query := `
				UPDATE mail_tracking 
				SET 
					carrier = NULL,
					exchange = ?,
					last_carrier = (SELECT id FROM users WHERE guid = ?)
				WHERE
					carrier = (SELECT id FROM users WHERE guid = ?)
			`

			result := tx.Exec(
				query,
				exchangeId,
				userGuid,
				userGuid,
			)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}

		// Any mail in the exchange has a chance to be picked up
		// The same user may not pick up mail they dropped off within the last 24 hours
		{
			query := `
				UPDATE mail_tracking 
				SET 
					carrier = (SELECT id FROM users WHERE guid = ?),
					exchange = NULL
				WHERE
					exchange = ? AND
					(
						SELECT
							# Take the ratio of closest user owned mailbox and the total distance the mail must travel. This creates relative probabilities to ensure mail tends to travel towards its destination.
							# Distance calculated is one dimensional and does not take into account both x,y coordinates. This could lead to backtracking or mail traveling tangentially to its destination.
							#
							# Keep at least a small probability (1%) that the user will pick up the mail. 
							# This should prevent mail getting stuck for any edge case where the mail distance is shorter than the user mailbox distance.
							RAND() > GREATEST(0.01, ST_Distance_Sphere(ul.geo_coordinate, m_dest.geo_coordinate) / ST_Distance_Sphere(m_start.geo_coordinate, m_dest.geo_coordinate)) AS should_pickup
						FROM mail_tracking mt
						JOIN exchanges e ON e.id = mt.exchange
						JOIN locations l ON l.id = e.location_id
						JOIN mail m ON m.id = mt.mail_id
						JOIN locations m_start ON m_start.id = m.starting_location
						JOIN locations m_dest ON m_dest.id = m.destination
						JOIN mailboxes mb
						JOIN locations ul ON ul.id = mb.location_id
						JOIN users u ON u.id = mb.user_id
						WHERE 
							# Use the current user's mailbox locations.
							u.guid = ? AND

							mt.exchange = ? AND

							# Mail that no user is carrying.
							mt.carrier IS NULL AND
							
							# Mail that the user did not last touch.
							# In an edge case where only one carrier ever crosses an exchange, that mail would become stuck. 
							# To prevent this, allow a carrier to pick up mail that has been at the exchange for at least one day. 
							(
								mt.last_carrier IS NULL OR 
								mt.last_carrier != (SELECT id FROM users WHERE guid = ?) OR 
								DATEDIFF(CURRENT_TIMESTAMP(), updated_at) >= 1
							)
					)
			`

			result := tx.Exec(
				query,
				userGuid,
				exchangeId,
				userGuid,
				exchangeId,
				userGuid,
			)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}

		return nil
	})
}

func (self *MySql) GetMail(ctx context.Context, mailGuid string) (*models.Mail, error) {
	return nil, nil
}

func (self *MySql) GetUserMail(ctx context.Context, userGuid string) ([]models.Mail, error) {
	// type mailRow struct {
	// 	MailGuid    string
	// 	From        string
	// 	To          string
	// 	Contents    string
	// 	SentOn      time.Time
	// 	DeliveredOn time.Time
	// 	OpenedOn    time.Time
	// }

	// query := `
	// 	SELECT
	// 		m.mail_guid,
	// 		u2.user_guid AS "from",
	// 		u.user_guid AS "to",
	// 		m.contents,
	// 		m.sent_on,
	// 		m.delivered_on,
	// 		m.opened_on
	// 	FROM users u
	// 	JOIN mail m ON m.` + "`to`" + ` = u.id
	// 	JOIN users u2 ON u2.id = m.` + "`from`" + `
	// 	JOIN user_inbox ui ON ui.mail_id = m.id
	// 	WHERE
	// 		u.user_guid = ?
	// `

	// mailRows := []mailRow{}
	// result := self.db.Raw(query, string(userGuid)).Scan(&mailRows)
	// if result.Error != nil {
	// 	return nil, ErrInternalFailure
	// }

	// log.Debug(ctx, "mail: %+v", mailRows)

	// userMail := make([]mail.Mail, len(mailRows))
	// for rowIndex := range mailRows {
	// 	userMail[rowIndex] = mail.Mail{
	// 		Guid: mail.Guid(mailRows[rowIndex].MailGuid),
	// 		Attributes: mail.Attributes{
	// 			From:     string(mailRows[rowIndex].From),
	// 			To:       string(mailRows[rowIndex].To),
	// 			Contents: mailRows[rowIndex].Contents,
	// 		},
	// 		SentOn:      mailRows[rowIndex].SentOn,
	// 		DeliveredOn: mailRows[rowIndex].DeliveredOn,
	// 		OpenedOn:    mailRows[rowIndex].OpenedOn,
	// 	}
	// }

	// return userMail, nil

	return nil, nil
}

func (self *MySql) DeleteMail(ctx context.Context, mailGuid string) error {
	return nil
}

func (self *MySql) CreateMailbox(ctx context.Context, userGuid string, newMailbox models.Mailbox) error {
	return self.db.Transaction(func(tx *gorm.DB) error {
		{
			query := `
				INSERT INTO locations
				(address, geo_coordinate, long_1000_floor, lat_1000_floor)
				VALUES
				# (LONG, LAT)
				(?, ST_SRID(POINT(?, ?), 4326), FLOOR(?), FLOOR(?))
			`

			result := tx.Exec(
				query,
				newMailbox.Address,
				float32(newMailbox.Location.Longitude),
				float32(newMailbox.Location.Latitude),
				float32(math.Floor(float64(newMailbox.Location.Longitude)*1000)),
				float32(math.Floor(float64(newMailbox.Location.Latitude)*1000)),
			)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}

		{
			query := `
				INSERT INTO mailboxes
				(location_id, user_id)
				VALUES
				(LAST_INSERT_ID(), (SELECT id FROM users WHERE guid = ?))
			`

			result := tx.Exec(query, userGuid)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}
		return nil
	})
}

func (self *MySql) GetMailbox(ctx context.Context, mailboxAddress string) (*models.Mailbox, error) {
	type mailboxRow struct {
		Address   string
		UserGuid  string
		Capacity  uint32
		Longitude float32
		Latitude  float32
	}

	query := `
		SELECT
			l.address,
			ST_Longitude(l.geo_coordinate) AS longitude,
			ST_Latitude(l.geo_coordinate) AS latitude,
			u.guid AS user_guid
		FROM locations l
		JOIN mailboxes m ON m.location_id = l.id
		JOIN users u ON u.id = m.user_id
		WHERE
			l.address = ?
	`

	mailboxResult := mailboxRow{}
	result := self.db.Raw(query, mailboxAddress).Scan(&mailboxResult)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, ErrInternalFailure)
	}

	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return &models.Mailbox{
		UserGuid: mailboxResult.UserGuid,
		Address:  mailboxResult.Address,
		Location: models.Coordinate{
			Latitude:  mailboxResult.Latitude,
			Longitude: mailboxResult.Longitude,
		},
	}, nil
}

func (self *MySql) DeleteMailbox(ctx context.Context, mailboxAddress string) error {
	return nil
}

func (self *MySql) GetUserMailbox(ctx context.Context, userGuid string) (*models.Mailbox, error) {
	type mailboxRow struct {
		Address   string
		UserGuid  string
		Longitude float32
		Latitude  float32
	}

	query := `
		SELECT
			l.address,
			ST_Longitude(l.geo_coordinate) AS longitude,
			ST_Latitude(l.geo_coordinate) AS latitude,
			u.guid AS user_guid
		FROM users u 
		JOIN mailboxes m ON m.user_id = u.id
		JOIN locations l ON l.id = m.location_id
		WHERE
			u.guid = ?
	`

	mailboxResult := mailboxRow{}
	result := self.db.Raw(query, userGuid).Scan(&mailboxResult)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, ErrInternalFailure)
	}

	if result.RowsAffected == 0 {
		return nil, ErrMailboxNotFound
	}

	return &models.Mailbox{
		UserGuid: mailboxResult.UserGuid,
		Address:  mailboxResult.Address,
		Location: models.Coordinate{
			Longitude: mailboxResult.Longitude,
			Latitude:  mailboxResult.Latitude,
		},
	}, nil
}

func (self *MySql) GetNearbyMailboxes(ctx context.Context, location models.Coordinate, radiusMeters float32) ([]models.Mailbox, error) {
	// type mailboxRow struct {
	// 	Address   string
	// 	Owner     string
	// 	Capacity  uint32
	// 	Latitude  float32
	// 	Longitude float32
	// }

	// // See: https://martech.zone/calculate-great-circle-distance/
	// // TODO: It would be good to add a bounding box to the lat/lng values so that this query is not searching the entire table.
	// query := `
	// SELECT
	// 	m.address,
	// 	u.user_guid AS "owner",
	// 	m.capacity,
	// 	m.latitude,
	// 	m.longitude
	// FROM mailboxes m
	// JOIN users u ON u.id = m.owner
	// WHERE (((ACOS(SIN((? * PI()/180)) * SIN((m.latitude* PI()/180)) + COS((?* PI()/180)) * COS((m.latitude* PI()/180)) * COS(((? - m.longitude)* PI()/180)))) * 180/ PI()) * 60 * 1.1515* 1.609344) <= ?
	// `

	// mailboxRows := []mailboxRow{}
	// result := self.db.Raw(query, location.Lat, location.Lat, location.Lng, metersToKilometers(radiusMeters)).Scan(&mailboxRows)
	// if result.Error != nil {
	// 	return nil, ErrInternalFailure
	// }

	// mailboxes := make([]mailbox.Mailbox, len(mailboxRows))
	// for rowIndex := range mailboxRows {
	// 	mailboxes[rowIndex] = mailbox.Mailbox{
	// 		Address: mailboxRows[rowIndex].Address,
	// 		Attributes: mailbox.Attributes{
	// 			Owner:    string(mailboxRows[rowIndex].Owner),
	// 			Capacity: mailboxRows[rowIndex].Capacity,
	// 			Location: geo.Coordinate{
	// 				Lat: geo.Latitude(mailboxRows[rowIndex].Latitude),
	// 				Lng: geo.Longitude(mailboxRows[rowIndex].Longitude),
	// 			},
	// 		},
	// 	}
	// }

	// return mailboxes, nil

	return nil, nil
}

func (self *MySql) GetMailboxMail(ctx context.Context, mailboxAddress string) ([]models.MailPreview, error) {
	type mailRow struct {
		Guid string
		From string
		To   string
	}

	query := `
		SELECT 
			m.guid,
			` + "`from`" + `, 
			` + "`to`" + `
		FROM mailboxes mb
		JOIN locations l ON l.id = mb.location_id
		JOIN mail m ON m.destination = l.id
		WHERE
			l.address = ? AND 
			m.delivered_on IS NOT NULL
	`

	mailResult := []mailRow{}
	result := self.db.Raw(
		query,
		mailboxAddress,
	).Scan(&mailResult)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, ErrInternalFailure)
	}

	mailboxMail := make([]models.MailPreview, len(mailResult))
	for mailIndex := range mailResult {
		mailboxMail[mailIndex] = models.MailPreview{
			Guid: mailResult[mailIndex].Guid,
			From: mailResult[mailIndex].From,
			To:   mailResult[mailIndex].To,
		}
	}

	return mailboxMail, nil
}

func (self *MySql) DropOffMail(ctx context.Context, carrierGuid string, mailboxAddress string) ([]string, error) {
	// err := self.db.Transaction(func(tx *gorm.DB) error {
	// 	type mailboxRow struct {
	// 		Id       uint64
	// 		Owner    string
	// 		Capacity uint32
	// 	}
	// 	mailboxLookupQuery := `
	// 		SELECT
	// 			id,
	// 			owner,
	// 			capacity
	// 		FROM mailboxes
	// 		WHERE
	// 			address = ?
	// 	`
	// 	dropOffMailbox := mailboxRow{}
	// 	// TODO: Could potentially cache this query instead of looking it up every time.
	// 	if err := tx.Raw(mailboxLookupQuery, mailboxAddress).Scan(&dropOffMailbox).Error; err != nil {
	// 		return errors.Wrap(err, ErrInternalFailure)
	// 	}

	// 	dropOffLimitQuery := `
	// 		SELECT
	// 			capacity - (SELECT COUNT(*) FROM mailbox_mail WHERE id = ? FOR UPDATE) AS drop_off_limit
	// 		FROM mailboxes
	// 		WHERE
	// 			id = ?
	// 	`
	// 	var dropOffLimit int
	// 	if err := tx.Raw(dropOffLimitQuery, dropOffMailbox.Id, dropOffMailbox.Id).Scan(&dropOffLimit).Error; err != nil {
	// 		return errors.Wrap(err, ErrInternalFailure)
	// 	}

	// 	if dropOffLimit <= 0 {
	// 		return errors.New("mailbox full")
	// 	}

	// 	carrierMailQuery := `
	// 		SELECT
	// 			mail_id
	// 		FROM mail_carriers
	// 		WHERE
	// 			user_id = (SELECT id FROM users WHERE user_guid = ?)
	// 		LIMIT ?
	// 		FOR UPDATE
	// 	`
	// 	dropOffMail := make([]interface{}, 0, dropOffLimit)
	// 	if err := tx.Raw(carrierMailQuery, carrierGuid, dropOffLimit).Scan(&dropOffMail).Error; err != nil {
	// 		return errors.Wrap(err, ErrInternalFailure)
	// 	}

	// 	if len(dropOffMail) == 0 {
	// 		return errors.New("carrier has no mail")
	// 	}

	// 	dropOffMailValues := make([]interface{}, dropOffLimit*2)
	// 	//deleteIn := ""
	// 	var dropOffMailIndex int
	// 	dropOffMailQuery := `
	// 		INSERT INTO mailbox_mail
	// 		(mailbox_id, mail_id)
	// 		VALUES
	// 		`
	// 	for index, mailDrop := range dropOffMail {
	// 		if index != 0 {
	// 			dropOffMailQuery += `,`
	// 			//deleteIn += ","
	// 		}
	// 		dropOffMailQuery += `(?,?)`
	// 		dropOffMailValues[dropOffMailIndex] = dropOffMailbox.Id
	// 		dropOffMailIndex++
	// 		dropOffMailValues[dropOffMailIndex] = mailDrop
	// 		dropOffMailIndex++

	// 		//deleteIn += "?"
	// 	}
	// 	if err := tx.Exec(dropOffMailQuery, dropOffMailValues...).Error; err != nil {
	// 		return errors.Wrap(err, ErrInternalFailure)
	// 	}

	// 	// deleteMailCarrierQuery := `
	// 	// 	DELETE
	// 	// 	FROM mail_carriers
	// 	// 	WHERE
	// 	// 		mail_id IN (` + deleteIn + `)
	// 	// `
	// 	// if err := tx.Exec(deleteMailCarrierQuery, dropOffMail...).Error; err != nil {
	// 	// 	return err
	// 	// }

	// 	return nil
	// })

	// if err != nil {
	// 	return nil, errors.Wrap(err, ErrInternalFailure)
	// }

	return []string{}, nil
}

func (self *MySql) PickUpMail(ctx context.Context, carrierGuid string, mailboxAddress string) ([]string, error) {
	return nil, nil
}

func (self *MySql) OpenMail(ctx context.Context, mailGuid string) (*models.Mail, error) {
	var openedMail *models.Mail
	return openedMail, self.db.Transaction(func(tx *gorm.DB) error {
		// Select the mail record.
		{
			type mailRow struct {
				Guid string

				RecipientGuid string

				FromMailboxAddress string
				ToMailboxAddress   string

				From string
				To   string
				Body string

				// Metadata
				SentOn      time.Time
				DeliveredOn time.Time
				OpenedOn    time.Time
			}

			query := `
				SELECT 
					m.guid,
					` + "m.`from`" + `, 
					` + "m.`to`" + `,
					from_address AS from_mailbox_address,
					to_address AS to_mailbox_address,
					m.body,
					m.sent_on,
					m.delivered_on,
					m.opened_on,
					u.guid AS recipient_guid
				FROM mail m
				LEFT JOIN locations from_loc ON from_loc.id = m.id
				LEFT JOIN locations to_loc ON to_loc.id = m.id
				LEFT JOIN users u ON u.id = m.recipient
				WHERE
					m.guid = ?
			`

			mailResult := mailRow{}
			result := self.db.Raw(
				query,
				mailGuid,
			).Scan(&mailResult)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}

			if result.RowsAffected == 0 {
				return ErrMailNotFound
			}

			openedMail = &models.Mail{
				Guid: mailResult.Guid,

				ToGuid: mailResult.RecipientGuid,

				FromMailboxAddress: mailResult.FromMailboxAddress,
				ToMailboxAddress:   mailResult.ToMailboxAddress,

				Contents: models.MailContents{
					From: mailResult.From,
					To:   mailResult.To,
					Body: mailResult.Body,
				},

				// Metadata
				SentOn:      mailResult.SentOn,
				DeliveredOn: mailResult.DeliveredOn,
				OpenedOn:    mailResult.OpenedOn,
			}
		}

		// If the mail has not been opened before, update the mail record.
		if !openedMail.IsOpened() {
			openedMail.OpenedOn = time.Now().UTC()

			query := `
				UPDATE mail
				SET 
					opened_on = ?
				WHERE
					guid = ?
			`

			result := tx.Exec(
				query,
				sqlDate(openedMail.OpenedOn),
				mailGuid,
			)
			if result.Error != nil {
				return errors.Wrap(result.Error, ErrInternalFailure)
			}
		}

		return nil
	})

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
