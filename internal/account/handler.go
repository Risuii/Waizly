package account

import (
	"encoding/json"
	"net/http"

	newJWT "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"waizly/config/jwt"
	"waizly/helpers/exception"
	"waizly/helpers/response"
	"waizly/models"
)

type AccountHandler struct {
	Validate *validator.Validate
	UseCase  AccountUseCase
}

func NewAccountHandler(router *mux.Router, validate *validator.Validate, usecase AccountUseCase) {
	handler := &AccountHandler{
		Validate: validate,
		UseCase:  usecase,
	}

	router.HandleFunc("/account/register", handler.Register).Methods(http.MethodPost)
	router.HandleFunc("/account/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/account/detail", handler.DetailAccount).Methods(http.MethodGet)
	router.HandleFunc("/account/update", handler.UpdateAccount).Methods(http.MethodPatch)
	router.HandleFunc("/account/delete", handler.DeleteAccount).Methods(http.MethodDelete)
}

func (handler *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var params models.RegisterRequest

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		res = response.Error(response.StatusUnprocessableEntity, err)
		res.JSON(w)
		return
	}

	err = handler.Validate.StructCtx(ctx, params)
	if err != nil {
		res = response.Error(response.StatusBadRequest, err)
		res.JSON(w)
		return
	}

	res = handler.UseCase.Register(ctx, params)

	res.JSON(w)
}

func (handler *AccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var params models.LoginRequest

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		res = response.Error(response.StatusUnprocessableEntity, err)
		res.JSON(w)
		return
	}

	err = handler.Validate.StructCtx(ctx, params)
	if err != nil {
		res = response.Error(response.StatusBadRequest, err)
		res.JSON(w)
		return
	}

	res, token := handler.UseCase.Login(ctx, params)

	if token.Token == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "",
			Path:     "",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token.Token,
		HttpOnly: true,
	})

	res.JSON(w)
}

func (handler *AccountHandler) DetailAccount(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	ctx := r.Context()

	c, err := r.Cookie("token")
	if err != nil {
		res = response.Error(response.StatusUnauthorized, exception.ErrUnauthorized)
		res.JSON(w)
		return
	}

	tokenString := c.Value
	claims := &jwt.JWTclaim{}

	newJWT.ParseWithClaims(tokenString, claims, func(t *newJWT.Token) (interface{}, error) {
		return jwt.JWT_KEY, nil
	})

	if claims.ID == 0 {
		res = response.Error(response.StatusUnauthorized, exception.ErrUnauthorized)
		res.JSON(w)
		return
	}

	res = handler.UseCase.DetailAccount(ctx, claims.ID)

	res.JSON(w)
}

func (handler *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var account models.Account
	ctx := r.Context()

	c, err := r.Cookie("token")
	if err != nil {
		res = response.Error(response.StatusUnauthorized, exception.ErrUnauthorized)
		res.JSON(w)
		return
	}

	tokenString := c.Value
	claims := &jwt.JWTclaim{}

	newJWT.ParseWithClaims(tokenString, claims, func(t *newJWT.Token) (interface{}, error) {
		return jwt.JWT_KEY, nil
	})

	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		res = response.Error(response.StatusUnprocessableEntity, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	err = handler.Validate.StructCtx(ctx, account)
	if err != nil {
		res = response.Error(response.StatusBadRequest, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	res = handler.UseCase.UpdateAccount(ctx, claims.ID, account)

	res.JSON(w)
}

func (handler *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	ctx := r.Context()

	c, err := r.Cookie("token")
	if err != nil {
		res = response.Error(response.StatusUnauthorized, exception.ErrUnauthorized)
		res.JSON(w)
		return
	}

	tokenString := c.Value
	claims := &jwt.JWTclaim{}

	newJWT.ParseWithClaims(tokenString, claims, func(t *newJWT.Token) (interface{}, error) {
		return jwt.JWT_KEY, nil
	})

	res = handler.UseCase.DeleteAccount(ctx, claims.ID)

	res.JSON(w)
}
