package handler

import (
	"be-service-auth/domain"
	"be-service-auth/helper"
	"context"
	"errors"
	"io"
	"strings"

	nethttp "net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/labstack/gommon/log"
	"github.com/valyala/fasthttp"
)

type AuthHandler struct {
	AuthUseCase  domain.AuthUseCase
	OAuthUseCase domain.OAuthUseCase
}

func (ah *AuthHandler) authorizationAuth(c *fiber.Ctx) (resAuth domain.ResponseLoginDTO, err error) {
	var token string
	auth := c.GetReqHeaders()
	authorization := auth["Authorization"]
	if len(authorization) != 0 {
		token = authorization[0][7:]
	}
	if token == "" {
		err = errors.New("token not found")
		return
	}
	c.Locals("token", token)

	resAuth, err = ah.AuthUseCase.AuthorizeAuth(c.Context(), token)
	if err != nil {
		log.Error(err)
	}

	return
}
func (az *AuthHandler) AuthorizationAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth, err := az.authorizationAuth(c)

		if err != nil {
			if err.Error() == "Expired" {
				return c.Status(fasthttp.StatusForbidden).SendString("Invalid token")
			}
			return helper.HttpSimpleResponse(c, fasthttp.StatusUnauthorized)
		}

		c.Locals("session", domain.ResponseLoginDTO{
			ID:              auth.ID,
			Email:           auth.Email,
			Role:            auth.Role,
			Token:           auth.Token,
			ExpiredDatetime: auth.ExpiredDatetime,
		})

		return c.Next()
	}
}

func (ah *AuthHandler) Login(c *fiber.Ctx) (err error) {
	var req domain.LoginRequestDTO
	err = c.BodyParser(&req)
	if err != nil {
		return helper.HttpSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		log.Error(err)
		return helper.HttpSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	session, err := ah.AuthUseCase.Login(c.Context(), req)
	if err != nil {
		log.Error(err)
		if err.Error() == "Not found" {
			log.Error("Data not found")
			return c.SendStatus(fasthttp.StatusUnauthorized)
		}

		if err.Error() == "User Data Not Found" {
			log.Error("User data not found")
			return c.SendStatus(fasthttp.StatusUnauthorized)
		}

		if err.Error() == "Not match" {
			log.Error("Password not match")
			return c.SendStatus(fasthttp.StatusUnauthorized)
		}

		if err.Error() == "Forbidden" {
			log.Error("Forbidden")
			return c.SendStatus(fasthttp.StatusForbidden)
		}
		return helper.HttpSimpleResponse(c, fasthttp.StatusInternalServerError)

	}

	res := domain.ResponseLoginDTO{
		Token:           session.Token,
		ExpiredDatetime: session.ExpiredDatetime,
	}
	return c.Status(200).JSON(res)
}

func (ah *AuthHandler) Auth(c *fiber.Ctx) (err error) {
	session := c.Locals("session")

	if err != nil {
		if err.Error() == "Expired" {
			return c.Status(fasthttp.StatusForbidden).SendString("Token invalid")
		}
		return c.SendString(err.Error())
	}

	return c.Status(200).JSON(session)
}

func (ah *AuthHandler) OAuth2(c *fiber.Ctx) (err error) {
	return
}

func (ah *AuthHandler) PostTokenOAuth2(w nethttp.ResponseWriter, r *nethttp.Request) {
	var ctx context.Context
	// Generate Credential OAUTH2
	CredentialValues, _ := helper.GenerateOAuthCredential(ctx, r)

	log.Info(CredentialValues)

	// Data baru yang ingin ditambahkan ke tubuh permintaan
	newPayloadData := "&client_id=" + CredentialValues.ClientID + "&client_secret=" + CredentialValues.ClientSecret

	existingBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Error reading existing request body:", err)
		return
	}
	// r.Body.Close()

	// Menggabungkan data lama dan data baru
	newRequestBody := append(existingBody, []byte(newPayloadData)...)

	// Membuat ulang tubuh permintaan dengan data yang diperbarui
	r.Body = io.NopCloser(strings.NewReader(string(newRequestBody)))

	responseToken := ah.OAuthUseCase.TokenOAuth(ctx, w, r)
	if responseToken != nil {
		log.Error(responseToken.Error())
	}
}
func (ah *AuthHandler) Logout(c *fiber.Ctx) (err error) {
	// var sessionData domain.ResponseSessionLogin
	session := c.Locals("session")

	sessionData, ok := session.(domain.ResponseLoginDTO)
	if !ok {
		return c.Status(500).SendString("Invalid session type") // Internal Server Error
	}

	err = ah.AuthUseCase.Logout(c.Context(), sessionData.Token)

	if err != nil {
		return
	}

	return helper.HttpSimpleResponse(c, fasthttp.StatusOK)
}

func (ah *AuthHandler) TokenOauth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		httpReq, err := adaptor.ConvertRequest(c, false)
		if err != nil {
			return err
		}

		err = ah.OAuthUseCase.ValidateBarrerToken(c.Context(), httpReq)

		log.Print(err)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		c.Locals("isPartner", true)

		return c.Next()
	}
}
