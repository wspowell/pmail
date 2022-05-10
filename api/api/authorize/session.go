package authorize

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
)

type authSessionRequest struct {
	mime.Json

	OneTimePassword   string `json:"otp"`
	SessionTtlMinutes int    `json:"ttlMinutes"`
}

type sessionAuthMessage struct {
	mime.Json

	JwtToken  string    `json:"jwtToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type authSession struct {
	auth.User
	Datastore db.Datastore
	body.Request[authSessionRequest]
}

func (self *authSession) Handle(ctx context.Context) (int, error) {
	expiresAt := time.Now().UTC().Add(time.Duration(self.RequestBody.SessionTtlMinutes) * time.Minute)
	jwtToken, err := self.JwtAuth.UserToken(&self.AuthorizedUser.User, expiresAt)
	if err != nil {
		return httpstatus.InternalServerError, errors.Wrap(err, errJwtError)
	}

	// TODO: Send jwtToken to websocket that holds the OTP.
	log.Debug(ctx, "otp: %s, jwtToken: %s", self.RequestBody.OneTimePassword, jwtToken)

	connectionId, err := self.Datastore.GetOneTimePassword(ctx, self.RequestBody.OneTimePassword)
	if err != nil {
		if errors.Is(err, db.ErrOneTimePasswordNotFound) {
			return httpstatus.NotFound, errors.Wrap(err, errNotFound)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, err
		} else {
			return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	awsSession, err := session.NewSession(&aws.Config{})
	if err != nil {
		return httpstatus.InternalServerError, errors.Wrap(err, errAwsError)
	}

	apiGatewayManagementApi := apigatewaymanagementapi.New(awsSession, aws.NewConfig().WithEndpoint("50elu9kia5.execute-api.us-east-1.amazonaws.com/test"))

	message := &sessionAuthMessage{
		JwtToken:  jwtToken,
		ExpiresAt: expiresAt,
	}
	messageBytes := []byte{}
	if err := message.MarshalMimeTypeJson(&messageBytes, message); err != nil {
		return httpstatus.InternalServerError, err
	}

	_, err = apiGatewayManagementApi.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionId),
		Data:         messageBytes,
	})

	apiGatewayManagementApi.DeleteConnection(&apigatewaymanagementapi.DeleteConnectionInput{
		ConnectionId: aws.String(connectionId),
	})

	return httpstatus.NoContent, nil
}
