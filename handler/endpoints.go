package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) UserRegistration() echo.HandlerFunc {
	return s.userRegistration
}

func (s *Server) userRegistration(ctx echo.Context) error {
	spanCtx := context.Background()
	writeResponse := func(id int, err string, status int) error {
		return ctx.JSON(status, UserRegistrationResponse{
			ID:    id,
			Error: err,
		})
	}

	request := UserRegistrationRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		return writeResponse(0, err.Error(), http.StatusBadRequest)
	}

	if err := request.Validate(); err != "" {
		return writeResponse(0, err, http.StatusBadRequest)
	}

	generatedPassword, err := request.HashAndSaltPassword()
	if err != nil {
		return writeResponse(0, err.Error(), http.StatusInternalServerError)
	}

	id, err := s.Repository.InsertUser(spanCtx, repository.InsertUser{
		PhoneNumber: request.PhoneNumber,
		Name:        request.Name,
		Password:    generatedPassword,
	})
	if err != nil {
		return writeResponse(0, err.Error(), http.StatusInternalServerError)
	}

	return writeResponse(id, "", http.StatusOK)
}

func (s *Server) UserLogin() echo.HandlerFunc {
	return s.userLogin
}

func (s *Server) userLogin(ctx echo.Context) error {
	spanCtx := context.Background()
	writeResponse := func(user repository.User, errMsg string, status int) error {
		var (
			bearerToken string
			err         error
		)
		if errMsg == "" {
			bearerToken, err = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
				Issuer:   "sawitpro",
				Subject:  fmt.Sprintf("%d", user.ID),
				IssuedAt: jwt.NewNumericDate(time.Now()),
				ID:       uuid.NewString(),
			}).SignedString(s.PrivateKey)
			if err != nil {
				return err
			}
		}

		return ctx.JSON(status, UserLoginResponse{
			ID:    user.ID,
			Token: bearerToken,
			Error: errMsg,
		})
	}

	request := UserLoginRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		return writeResponse(repository.User{}, err.Error(), http.StatusBadRequest)
	}

	user, err := s.Repository.GetUserByPhoneNumber(spanCtx, request.PhoneNumber)
	switch errors.Cause(err) {
	case nil:
		break
	case sql.ErrNoRows:
		return writeResponse(repository.User{}, err.Error(), http.StatusNotFound)
	default:
		return writeResponse(repository.User{}, err.Error(), http.StatusInternalServerError)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return writeResponse(repository.User{}, err.Error(), http.StatusInternalServerError)
	}

	if err := s.Repository.InsertUserLoginHistory(spanCtx, repository.InsertUserLoginHistory{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		Password:    user.Password,
	}); err != nil {
		return writeResponse(repository.User{}, err.Error(), http.StatusInternalServerError)
	}

	return writeResponse(user, "", http.StatusOK)
}

func (s *Server) GetUser() echo.HandlerFunc {
	return s.getUser
}

func (s *Server) getUser(ctx echo.Context) error {
	spanCtx := context.Background()
	writeResponse := func(user repository.User, errMsg string, status int) error {
		return ctx.JSON(status, GetUserResponse{
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
			Error:       errMsg,
		})
	}

	authorization := ctx.Request().Header["Authorization"]
	if len(authorization) < 1 {
		return writeResponse(repository.User{}, "Forbidden", http.StatusForbidden)
	}

	tokenString := strings.TrimSpace(strings.Replace(authorization[0], "Bearer", "", 1))
	var tokenClaims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(t *jwt.Token) (interface{}, error) {
		return s.PublicKey, nil
	})
	if err != nil {
		return writeResponse(repository.User{}, "Forbidden "+err.Error(), http.StatusForbidden)
	}

	id, err := strconv.Atoi(tokenClaims.Subject)
	if err != nil {
		return writeResponse(repository.User{}, "Forbidden "+err.Error(), http.StatusForbidden)
	}

	switch user, err := s.Repository.GetUserByID(spanCtx, id); errors.Cause(err) {
	case nil:
		return writeResponse(user, "", http.StatusOK)
	case sql.ErrNoRows:
		return writeResponse(repository.User{}, err.Error(), http.StatusNotFound)
	default:
		return writeResponse(repository.User{}, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) UpdateUser() echo.HandlerFunc {
	return s.updateUser
}

func (s *Server) updateUser(ctx echo.Context) error {
	spanCtx := context.Background()
	writeResponse := func(user repository.User, errMsg string, status int) error {
		return ctx.JSON(status, GetUserResponse{
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
			Error:       errMsg,
		})
	}

	authorization := ctx.Request().Header["Authorization"]
	if len(authorization) < 1 {
		return writeResponse(repository.User{}, "Forbidden", http.StatusForbidden)
	}

	tokenString := strings.TrimSpace(strings.Replace(authorization[0], "Bearer", "", 1))
	var tokenClaims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(t *jwt.Token) (interface{}, error) {
		return s.PublicKey, nil
	})
	if err != nil {
		return writeResponse(repository.User{}, "Forbidden "+err.Error(), http.StatusForbidden)
	}

	id, err := strconv.Atoi(tokenClaims.Subject)
	if err != nil {
		return writeResponse(repository.User{}, "Forbidden "+err.Error(), http.StatusForbidden)
	}

	request := UpdateUserRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		return writeResponse(repository.User{}, err.Error(), http.StatusBadRequest)
	}

	if err := request.Validate(); err != "" {
		return writeResponse(repository.User{}, err, http.StatusBadRequest)
	}

	user, err := s.Repository.UpdateUser(spanCtx, repository.UpdateUser{
		ID:          id,
		PhoneNumber: request.PhoneNumber,
		Name:        request.Name,
	})
	if err != nil {
		if strings.Contains(err.Error(), "user_phone_number_key") {
			return writeResponse(repository.User{}, err.Error(), http.StatusConflict)
		}

		return writeResponse(repository.User{}, err.Error(), http.StatusInternalServerError)
	}

	return writeResponse(user, "", http.StatusOK)
}
