package handler

import (
	"be-service-auth/domain"
	"be-service-auth/helper"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/valyala/fasthttp"
)

type AuthHandler struct {
	AuthUseCase domain.AuthUseCase
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
		return c.SendString(err.Error())
	}

	return c.Status(200).JSON(session)
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
