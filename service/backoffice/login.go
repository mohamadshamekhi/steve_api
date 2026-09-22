package backofficeservice

import (
	"context"

	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) SendOTPCode(ctx context.Context, req backofficeauthparam.SendLoginCodeRequest) (backofficeauthparam.SendLoginCodeResponse, error) {
	const op = "backofficeservice.Login"

	err := s.authSvc.GenerateOTPCode(ctx, req.PhoneNumber, s.config.Prefix)
	if err != nil {
		return backofficeauthparam.SendLoginCodeResponse{}, richerror.New(op)
	}

	// return resp.
	return backofficeauthparam.SendLoginCodeResponse{}, nil
}

func (s Service) VerifyOTPCode(ctx context.Context, req backofficeauthparam.VerifyLoginCodeRequest) (backofficeauthparam.VerifyLoginCodeResponse, error) {
	const op = "backofficeservice.VerifyOTPCode"

	resp, err := s.authSvc.VerifyOTPCode(ctx, req.PhoneNumber, req.Code, s.config.Prefix)
	if err != nil {
		return backofficeauthparam.VerifyLoginCodeResponse{}, err
	}
	if resp {
		backoffice, iErr := s.repo.GetBackOfficeUserByPhoneNumber(ctx, req.PhoneNumber)
		if iErr != nil {
			return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(iErr)
		}
		rolePermissionKeys, rErr := s.roleRepo.GetPermissionKeysOfRole(ctx, backoffice.Roles)
		if rErr != nil {
			return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(rErr)
		}
		permissionKeys := mergePermissionKeys(backoffice.Permissions, rolePermissionKeys)

		token, gErr := s.authSvc.GenerateAccessToken(backoffice.ID, s.config.Prefix, backoffice.Roles, permissionKeys)
		if gErr != nil {
			return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(gErr)
		}
		return backofficeauthparam.VerifyLoginCodeResponse{
			Token: token,
		}, nil
	}
	return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).
		WithMessageKey("validation.incorrect_opt_code")
}

// mergePermissionKeys returns the de-duplicated union of the user's direct
// permissions and the permissions granted by their roles (same rule as
// aclservice.CheckAccess).
func mergePermissionKeys(direct []string, fromRoles []string) []string {
	seen := make(map[string]struct{}, len(direct)+len(fromRoles))
	merged := make([]string, 0, len(direct)+len(fromRoles))
	for _, list := range [][]string{direct, fromRoles} {
		for _, k := range list {
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			merged = append(merged, k)
		}
	}
	return merged
}
