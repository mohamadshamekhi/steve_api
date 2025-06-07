package authservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func (s Service) GenerateToken(ctx context.Context, phoneNumber string, subject string, kind string, expireDuration time.Duration) (string, error) {
	const op = "authservice.GenerateToken"
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		Kind:        kind,
		PhoneNumber: phoneNumber,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.SignKey))
	if err != nil {
		return "", richerror.New(op).WithWarpError(err).WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return tokenString, nil
}
