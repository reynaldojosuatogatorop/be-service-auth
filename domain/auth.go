package domain

import (
	"context"
	"time"
)

type LoginRequestDTO struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type ResponseLoginDTO struct {
	ID              int       `redis:"id" json:"id,omitempty"`
	Email           string    `redis:"email" json:"email,omitempty"`
	Password        string    `redis:"password" json:"password,omitempty"`
	Role            string    `redis:"role" json:"role,omitempty"`
	Token           string    `redis:"token" json:"token"`
	ExpiredDatetime time.Time `redis:"expired_datetime" json:"expired_datetime"`
}

type AuthUseCase interface {
	Login(ctx context.Context, request LoginRequestDTO) (response ResponseLoginDTO, err error)
	AuthorizeAuth(ctx context.Context, request string) (response ResponseLoginDTO, err error)
	Logout(ctx context.Context, token string) (err error)
}

type AuthMySQLRepo interface {
	GetUserByEmail(ctx context.Context, email string) (res ResponseLoginDTO, err error)
}

type AuthRedisRepo interface {
	CreateSession(ctx context.Context, user ResponseLoginDTO, token string) (session ResponseLoginDTO, err error)
	GetSession(ctx context.Context, username string) (session ResponseLoginDTO, err error)
	DeleteSession(ctx context.Context, token string) (err error)
	GetAuth(ctx context.Context, token string) (session ResponseLoginDTO, err error)
}

type AuthGRPCRepo interface {
	GetFromAuth(ctx context.Context, token string) (response ResponseLoginDTO, err error)
}
