package domain

import (
	"context"
	"net/http"
)

type OAuth2ClientCredential struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ResponseB2BDTO struct {
	ID           int    `json:"ïd"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Role         string `json:"role"`
	Domain       string `json:"domain"`
	DtmCrt       string `json:"dtm_crt"`
	DtmUpd       string `json:"dtm_upd"`
}

type OAuthUseCase interface {
	TokenOAuth(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error)
	ValidateBarrerToken(ctx context.Context, r *http.Request) (err error)
	GetAllB2BData(ctx context.Context) (response []ResponseB2BDTO, err error)
}

type OAuthMySQLRepo interface {
	GetAllB2BData(ctx context.Context) (response []ResponseB2BDTO, err error)
}
