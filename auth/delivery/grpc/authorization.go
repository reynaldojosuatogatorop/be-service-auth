package grpc

import (
	auth "be-service-auth/auth/delivery/grpc/authorization"
	"be-service-auth/domain"
	"context"
)

type Auth struct {
	auth.UnimplementedAuthorizationServiceServer
	Auth domain.AuthUseCase
}

func NewGRPCAuth(usecaseAuth domain.AuthUseCase) *Auth {
	return &Auth{
		Auth: usecaseAuth,
	}
}

func (s *Auth) GetSessionServiceAuth(ctx context.Context, req *auth.AuthorizationAuthServiceRequest) (response *auth.AuthorizationAuthServiceResponse, err error) {
	res, err := s.Auth.AuthorizeAuth(ctx, req.Token)
	if err != nil {
		return response, err
	}

	response = &auth.AuthorizationAuthServiceResponse{
		Id:              int64(res.ID),
		Email:           res.Email,
		Role:            res.Role,
		Token:           res.Token,
		ExpiredDatetime: res.ExpiredDatetime.String(),
	}
	return
}
