package usecase

import (
	"be-service-auth/domain"
	"context"

	nethttp "net/http"

	"github.com/labstack/gommon/log"
	serveroauth2 "gopkg.in/oauth2.v3/server"
)

type oAuthUsecase struct {
	oAuthMySQLRepo domain.OAuthMySQLRepo
	authRedisRepo  domain.AuthRedisRepo
	oautHttp       *serveroauth2.Server
}

func NewOAuthUsecase(OAuthMySQL domain.OAuthMySQLRepo, AuthRedisRepo domain.AuthRedisRepo, oautHttp *serveroauth2.Server) domain.OAuthUseCase {
	return &oAuthUsecase{
		oAuthMySQLRepo: OAuthMySQL,
		authRedisRepo:  AuthRedisRepo,
		oautHttp:       oautHttp,
	}
}

func (oau *oAuthUsecase) TokenOAuth(ctx context.Context, w nethttp.ResponseWriter, r *nethttp.Request) (err error) {
	err = oau.oautHttp.HandleTokenRequest(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (oau *oAuthUsecase) ValidateBarrerToken(ctx context.Context, r *nethttp.Request) (err error) {
	tokenInfo, err := oau.oautHttp.ValidationBearerToken(r)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(tokenInfo.GetClientID())
	return
}

func (oau *oAuthUsecase) GetAllB2BData(ctx context.Context) (response []domain.ResponseB2BDTO, err error) {
	response, err = oau.oAuthMySQLRepo.GetAllB2BData(ctx)
	if err != nil {
		return
	}
	return
}
