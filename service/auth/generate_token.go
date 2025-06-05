package authservice

import "context"

func (s Service) GenerateToken(ctx context.Context, phoneNumber string) (string, error) {
	const op = "authservice.GenerateToken"

}
