package resources

import (
	"github.com/wspowell/context"
	"github.com/wspowell/log"
	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
)

type Resources struct {
	LogConfig log.LoggerConfig
	Datastore db.Datastore
	UserAuth  auth.User
}

func Load() *Resources {
	ctx := context.Background()

	// TODO: Pull log level from AppConfig, env, or something.
	logConfig := log.NewConfig().WithLevel(log.LevelDebug)
	log.WithContext(ctx, logConfig)

	var err error

	datastore := db.NewMySql()
	if err = datastore.Connect(ctx); err != nil {
		log.Fatal(ctx, "failed to connect to datastore: %v", err)
	}
	if err = datastore.Migrate(); err != nil {
		log.Fatal(ctx, "failed to migrate datastore: %v", err)
	}

	signingKey, err := auth.GetSigningKey(ctx)
	if err != nil {
		// FIXME: Need to add Context() to restful.Server
		log.Fatal(ctx, "failed to get jwt signing key: %s", err)
	}

	jwtAuth := auth.NewJwt(signingKey)
	userAuth := auth.User{
		JwtAuth: jwtAuth,
	}

	return &Resources{
		LogConfig: logConfig,
		Datastore: datastore,
		UserAuth:  userAuth,
	}
}
