package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
}

type CustomValidator struct {
	validator *validator.Validate
}
