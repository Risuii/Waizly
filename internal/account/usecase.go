package account

import (
	"context"
	"log"
	"time"

	newJWT "github.com/dgrijalva/jwt-go"

	"waizly/config/bcrypt"
	"waizly/config/jwt"
	"waizly/helpers/exception"
	"waizly/helpers/response"
	"waizly/models"
)

type (
	AccountUseCase interface {
		Register(ctx context.Context, params models.RegisterRequest) response.Response
		Login(ctx context.Context, params models.LoginRequest) (response.Response, models.Token)
		DetailAccount(ctx context.Context, id int64) response.Response
		UpdateAccount(ctx context.Context, id int64, params models.Account) response.Response
		DeleteAccount(ctx context.Context, id int64) response.Response
	}

	accountUseCaseImpl struct {
		repository AccountRepository
		bcrypt     bcrypt.Bcrypt
	}
)

func NewAccountUseCase(repo AccountRepository, bcrypt bcrypt.Bcrypt) AccountUseCase {
	return &accountUseCaseImpl{
		repository: repo,
		bcrypt:     bcrypt,
	}
}

func (au *accountUseCaseImpl) Register(ctx context.Context, params models.RegisterRequest) response.Response {
	_, err := au.repository.FindByEmail(ctx, params.Email)

	if err == nil {
		return response.Error(response.StatusConflicted, exception.ErrConflicted)
	}

	if err != exception.ErrNotFound {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	hashedPassword, err := au.bcrypt.HashPassword(params.Password)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	account := models.Account{
		Username:  params.Username,
		Password:  hashedPassword,
		Email:     params.Email,
		CreatedAt: time.Now(),
	}

	ID, err := au.repository.Create(ctx, account)
	if err != nil {
		return response.Error(response.StatusInternalServerError, err)
	}

	account.ID = ID
	account.Password = ""

	return response.Success(response.StatusCreated, account)
}

func (au *accountUseCaseImpl) Login(ctx context.Context, params models.LoginRequest) (response.Response, models.Token) {
	account, err := au.repository.FindByEmail(ctx, params.Email)

	if err == exception.ErrNotFound {
		log.Println(err)
		return response.Error(response.StatusNotFound, exception.ErrNotFound), models.Token{}
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer), models.Token{}
	}

	isPasswordValid := au.bcrypt.ComparePasswordHash(params.Password, account.Password)

	if !isPasswordValid {
		return response.Error(response.StatusUnauthorized, err), models.Token{}
	}

	account.Password = ""

	claims := &jwt.JWTclaim{
		ID:    account.ID,
		Email: account.Email,
		StandardClaims: newJWT.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
		},
	}

	tokenAlgo := newJWT.NewWithClaims(newJWT.SigningMethodHS256, claims)

	token, err := tokenAlgo.SignedString(jwt.JWT_KEY)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer), models.Token{}
	}

	newToken := models.Token{
		Token: token,
	}

	return response.Success(response.StatusOK, account), newToken
}

func (au *accountUseCaseImpl) DetailAccount(ctx context.Context, id int64) response.Response {
	account, err := au.repository.FindByID(ctx, id)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrParams)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	account.Password = ""

	return response.Success(response.StatusOK, account)
}

func (au *accountUseCaseImpl) UpdateAccount(ctx context.Context, id int64, params models.Account) response.Response {

	account, err := au.repository.FindByID(ctx, id)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	account.ID = params.ID
	account.Username = params.Username
	account.Password = params.Password
	account.Email = params.Password
	account.UpdateAt = time.Now()

	err = au.repository.Update(ctx, account.ID, account)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, account)
}

func (au *accountUseCaseImpl) DeleteAccount(ctx context.Context, id int64) response.Response {

	err := au.repository.Delete(ctx, id)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	msg := "Success Delete Data"

	return response.Success(response.StatusOK, msg)
}
