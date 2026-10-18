package backofficeservice

import (
	"context"

	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// IssueServiceToken mints a long-lived access token for an existing
// backoffice user, meant for service/automation accounts (e.g. the external
// SEO/content agent) that cannot go through the OTP login flow. The token
// carries exactly the target user's current role + direct permissions, same
// as a normal login token, just with a much longer TTL.
func (s Service) IssueServiceToken(ctx context.Context, req backofficeparam.IssueServiceTokenRequest) (backofficeparam.IssueServiceTokenResponse, error) {
	const op = "backofficeservice.IssueServiceToken"

	backoffice, err := s.repo.GetBackOfficeUserByID(ctx, req.ID)
	if err != nil {
		return backofficeparam.IssueServiceTokenResponse{}, richerror.New(op).WithWarpError(err)
	}

	rolePermissionKeys, rErr := s.roleRepo.GetPermissionKeysOfRole(ctx, backoffice.Roles)
	if rErr != nil {
		return backofficeparam.IssueServiceTokenResponse{}, richerror.New(op).WithWarpError(rErr)
	}
	permissionKeys := mergePermissionKeys(backoffice.Permissions, rolePermissionKeys)

	token, gErr := s.authSvc.GenerateServiceToken(backoffice.ID, s.config.Prefix, backoffice.Roles, permissionKeys)
	if gErr != nil {
		return backofficeparam.IssueServiceTokenResponse{}, richerror.New(op).WithWarpError(gErr)
	}

	return backofficeparam.IssueServiceTokenResponse{Token: token}, nil
}
