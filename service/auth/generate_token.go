package authservice

import (
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func (s Service) generateToken(id string, subject string, kind string, expireDuration time.Duration) (string, error) {
	const op = "authservice.GenerateToken"
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		Kind:     kind,
		EntityID: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.SignKey))
	if err != nil {
		return "", richerror.New(op).WithWarpError(err).WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return tokenString, nil
}

func (s Service) GenerateAccessToken(id string, kind string) (string, error) {
	return s.generateToken(id ,s.cfg.AccessSubject ,kind , s.cfg.AccessExpirationTime)
}
