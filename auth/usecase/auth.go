package usecase

import (
	"be-service-auth/domain"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	authMySQLRepo domain.AuthMySQLRepo
	authRedisRepo domain.AuthRedisRepo
}

func NewAuthUsecase(AuthMySQL domain.AuthMySQLRepo, AuthRedisRepo domain.AuthRedisRepo) domain.AuthUseCase {
	return &authUsecase{
		authMySQLRepo: AuthMySQL,
		authRedisRepo: AuthRedisRepo,
	}
}

func (au *authUsecase) clearSessionIfExist(ctx context.Context, username string) {
	oldSession, err := au.authRedisRepo.GetSession(ctx, username)
	if err == nil {
		log.Debug("User " + oldSession.Email + " already have session, clean up old session")
		log.Debug(oldSession.Token)
		err = au.authRedisRepo.DeleteSession(ctx, oldSession.Token)

		if err != nil {
			log.Error("Error delete session token")
			return
		}
	}
	// return
}
func (au *authUsecase) Login(ctx context.Context, req domain.LoginRequestDTO) (session domain.ResponseLoginDTO, err error) {
	res, err := au.authMySQLRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return
	}

	if res.Role != "admin" {
		err = errors.New("Forbidden")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = errors.New("Not match")
		}
		return
	}

	au.clearSessionIfExist(ctx, res.Email)

	date := time.Now().Format(time.RFC3339)
	tokenByte := sha256.Sum256([]byte(strconv.FormatInt(int64(res.ID), 10) + "_" + res.Email + date))
	token := base64.URLEncoding.EncodeToString(tokenByte[:])
	if err != nil {
		fmt.Println("Gagal mengkonversi string ke int:", err)
		return
	}
	session = domain.ResponseLoginDTO{
		ID:    res.ID,
		Email: res.Email,
		Role:  res.Role,
	}

	session, err = au.authRedisRepo.CreateSession(ctx, session, token)
	if err != nil {
		return
	}
	return
}

func (au *authUsecase) AuthorizeAuth(ctx context.Context, token string) (authorization domain.ResponseLoginDTO, err error) {
	authorization, err = au.authRedisRepo.GetAuth(ctx, token)
	if err != nil {
		err = errors.New("user only")
		return
	}

	return
}

func (au *authUsecase) Logout(ctx context.Context, token string) (err error) {
	err = au.authRedisRepo.DeleteSession(ctx, token)
	if err != nil {
		log.Info(err)
	}
	return
}
