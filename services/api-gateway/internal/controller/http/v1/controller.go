package v1

import (
	v1 "api-gateway/docs/proto/v1"

	"go.uber.org/zap"
)

type V1 struct {
	authClient v1.AuthClient
	//userClient v1.UserClient

	l *zap.SugaredLogger
}
