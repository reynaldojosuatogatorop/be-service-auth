package handler

import (
	"be-service-auth/domain"
	"be-service-auth/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/valyala/fasthttp"
)

type AuthHandler struct {
	AuthUseCase domain.AuthUseCase
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
		return helper.HttpSimpleResponse(c, fasthttp.StatusInternalServerError)

	}

	res := domain.ResponseLoginDTO{
		Token:           session.Token,
		ExpiredDatetime: session.ExpiredDatetime,
	}
	return c.Status(200).JSON(res)
}
