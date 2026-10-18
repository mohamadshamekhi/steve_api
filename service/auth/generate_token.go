package authservice

import (
	"time"

	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
)

func (s Service) generateToken(id string, subject string, kind string, expireDuration time.Duration, roleKeys []string, permissionKeys []string) (string, error) {
	const op = "authservice.GenerateToken"
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			ID:        ulid.Make().String(),
		},
		Kind:           kind,
		EntityID:       id,
		RoleKeys:       roleKeys,
		PermissionKeys: permissionKeys,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.SignKey))
	if err != nil {
		return "", richerror.New(op).WithWarpError(err).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return tokenString, nil
}

func (s Service) GenerateAccessToken(id string, kind string, roleKeys []string, permissionKeys []string) (string, error) {
	return s.generateToken(id, s.cfg.AccessSubject, kind, s.cfg.AccessExpirationTime, roleKeys, permissionKeys)
}

// GenerateServiceToken issues a long-lived access token (TTL from
// cfg.ServiceTokenExpirationTime) for a service/automation backoffice
// account, e.g. the external SEO/content agent. It carries the same claim
// shape as a normal access token, so it works with the existing
// BackofficeAuth + ACL middleware unchanged.
func (s Service) GenerateServiceToken(id string, kind string, roleKeys []string, permissionKeys []string) (string, error) {
	return s.generateToken(id, s.cfg.AccessSubject, kind, s.cfg.ServiceTokenExpirationTime, roleKeys, permissionKeys)
}
