package handler

import (
	"github.com/labstack/echo/v4"
)

type ServerInterface interface {
	UserRegistration() echo.HandlerFunc
	UserLogin() echo.HandlerFunc
	GetUser() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
}
