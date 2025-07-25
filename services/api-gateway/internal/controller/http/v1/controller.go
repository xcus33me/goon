package v1

import (
	v1 "github.com/xcus33me/protos/gen/go/auth"

	"go.uber.org/zap"
)

type V1 struct {
	authClient v1.AuthClient
	//userClient v1.UserClient

	l *zap.SugaredLogger
}
