package account_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	newJWT "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"waizly/config/jwt"
	"waizly/helpers/response"
	"waizly/internal/account"
	"waizly/internal/account/mocks"
	"waizly/models"
)

func TestHandler_Register(t *testing.T) {
	t.Run("Register Success", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "test-1",
			Password: "test-1",
			Email:    "test@gmail.com",
		}

		resp := response.Success(response.StatusCreated, models.Account{})

		validate := validator.New()
		accountUseCase := new(mocks.AccountUseCase)
		accountUseCase.On("Register", mock.Anything, mock.AnythingOfType("models.RegisterRequest")).Return(resp)

		newReq, err := json.Marshal(req)
		if err != nil {
			t.Error(err)
			return
		}

		accountHandler := account.AccountHandler{
			Validate: validate,
			UseCase:  accountUseCase,
		}

		r := httptest.NewRequest(http.MethodPost, "/just/for/testing", bytes.NewReader(newReq))
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.Register)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusCreated, rb.Status, fmt.Sprintf("Should be status '%s'", response.StatusCreated))
		assert.NotNil(t, rb.Data, "Should be not nil")

		accountUseCase.AssertExpectations(t)
	})

	t.Run("Register Failed", func(t *testing.T) {
		type invalidReq struct {
			Data string
		}

		req := invalidReq{
			Data: "error",
		}

		validate := validator.New()
		accountUseCase := new(mocks.AccountUseCase)

		newReq, err := json.Marshal(req)
		if err != nil {
			t.Error(err)
			return
		}

		accountHandler := account.AccountHandler{
			Validate: validate,
			UseCase:  accountUseCase,
		}

		r := httptest.NewRequest(http.MethodPost, "/just/for/testing", bytes.NewReader(newReq))
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.Register)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusBadRequest, rb.Status, fmt.Sprintf("Should be '%s'", response.StatusBadRequest))
		assert.Nil(t, rb.Data, "Should be nil")

		accountUseCase.AssertExpectations(t)
	})
}

func TestHandler_login(t *testing.T) {
	t.Run("Login Success", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "test@gmail.com",
			Password: "test-1",
		}

		resp := response.Success(response.StatusOK, models.LoginRequest{})

		validate := validator.New()
		accountUseCase := new(mocks.AccountUseCase)
		accountUseCase.On("Login", mock.Anything, mock.AnythingOfType("models.LoginRequest")).Return(resp, models.Token{})

		newReq, err := json.Marshal(req)
		if err != nil {
			t.Error(err)
			return
		}

		accountHandler := account.AccountHandler{
			Validate: validate,
			UseCase:  accountUseCase,
		}

		r := httptest.NewRequest(http.MethodPost, "/just/for/testing", bytes.NewReader(newReq))
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.Login)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusOK, rb.Status, fmt.Sprintf("Should be status '%s'", response.StatusOK))
		assert.NotNil(t, rb.Data, "Should be not nil")

		accountUseCase.AssertExpectations(t)
	})

	t.Run("Login Failed", func(t *testing.T) {
		type invalidReq struct {
			Data string
		}

		req := invalidReq{
			Data: "error",
		}

		validate := validator.New()
		accountUseCase := new(mocks.AccountUseCase)

		accountHandler := account.AccountHandler{
			Validate: validate,
			UseCase:  accountUseCase,
		}

		newReq, err := json.Marshal(req)
		if err != nil {
			t.Error(err)
			return
		}

		r := httptest.NewRequest(http.MethodPost, "/just/for/testing", bytes.NewReader(newReq))
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.Login)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusBadRequest, rb.Status, fmt.Sprintf("Should be status '%s'", response.StatusBadRequest))
		assert.Nil(t, rb.Data, "Should be nil")

		accountUseCase.AssertExpectations(t)
	})
}

func TestHandler_DetailAccount(t *testing.T) {
	t.Run("Get Detail Success", func(t *testing.T) {

		resp := response.Success(response.StatusOK, models.Account{})

		accountUseCase := new(mocks.AccountUseCase)
		accountUseCase.On("DetailAccount", mock.Anything, mock.AnythingOfType("int64")).Return(resp)

		accountHandler := account.AccountHandler{
			UseCase: accountUseCase,
		}

		mockToken := &jwt.JWTclaim{
			ID:    1,
			Email: "test@test.com",
			StandardClaims: newJWT.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
			},
		}

		tokenAlgo := newJWT.NewWithClaims(newJWT.SigningMethodHS256, mockToken)

		token, err := tokenAlgo.SignedString(jwt.JWT_KEY)
		if err != nil {
			t.Error(err)
			return
		}

		r := httptest.NewRequest(http.MethodGet, "/just/for/testing", nil)
		r.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})

		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.DetailAccount)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusOK, rb.Status, fmt.Sprintf("Should be status '%s'", response.StatusOK))
		assert.NotNil(t, rb.Data, "Should be not nil")

		accountUseCase.AssertExpectations(t)
	})

	t.Run("Get Detail Unauthorized", func(t *testing.T) {

		accountUseCase := new(mocks.AccountUseCase)

		accountHandler := account.AccountHandler{
			UseCase: accountUseCase,
		}

		r := httptest.NewRequest(http.MethodGet, "/just/for/testing", nil)
		r.AddCookie(&http.Cookie{
			Name:  "",
			Value: "",
		})
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.DetailAccount)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusUnauthorized, rb.Status, fmt.Sprintf("Should be status '%s'", response.StatusUnauthorized))
		assert.Nil(t, rb.Data, "Should be not nil")
	})
}

func TestHandler_UpdateAccount(t *testing.T) {
	t.Run("Update Account Success", func(t *testing.T) {
		mockData := models.Account{
			ID:        1,
			Username:  "test-1",
			Password:  "test-password",
			Email:     "test@test.com",
			CreatedAt: time.Time{},
			UpdateAt:  time.Time{},
		}

		mockToken := &jwt.JWTclaim{
			ID:    1,
			Email: "test@test.com",
			StandardClaims: newJWT.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
			},
		}

		tokenAlgo := newJWT.NewWithClaims(newJWT.SigningMethodHS256, mockToken)

		token, err := tokenAlgo.SignedString(jwt.JWT_KEY)
		if err != nil {
			t.Error(err)
			return
		}

		validate := validator.New()
		resp := response.Success(response.StatusOK, models.Account{})
		accountUseCase := new(mocks.AccountUseCase)
		accountUseCase.On("UpdateAccount", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("models.Account")).Return(resp)

		reqData, err := json.Marshal(mockData)
		if err != nil {
			t.Error(err)
			return
		}

		accountHandler := account.AccountHandler{
			Validate: validate,
			UseCase:  accountUseCase,
		}

		r := httptest.NewRequest(http.MethodPatch, "/just/for/testing", bytes.NewReader(reqData))
		r.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.UpdateAccount)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusOK, rb.Status, fmt.Sprintf("Should be status %s", response.StatusOK))
		assert.NotNil(t, rb.Data, "Should be not nil")

		accountUseCase.AssertExpectations(t)
	})

	t.Run("Update Account unauthorized", func(t *testing.T) {
		mockData := models.Account{
			ID:        1,
			Username:  "test-1",
			Password:  "test-password",
			Email:     "test@test.com",
			CreatedAt: time.Time{},
			UpdateAt:  time.Time{},
		}

		validate := validator.New()

		accountUseCase := new(mocks.AccountUseCase)

		reqData, err := json.Marshal(mockData)
		if err != nil {
			t.Error(err)
			return
		}

		accountHandler := account.AccountHandler{
			Validate: validate,
			UseCase:  accountUseCase,
		}

		r := httptest.NewRequest(http.MethodPatch, "/just/for/testing", bytes.NewReader(reqData))
		r.AddCookie(&http.Cookie{
			Name:     "",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
		})
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.UpdateAccount)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusUnauthorized, rb.Status, fmt.Sprintf("Should be status %s", response.StatusUnauthorized))
		assert.Nil(t, rb.Data, "Should be nil")

		accountUseCase.AssertExpectations(t)
	})

	t.Run("Update Account Error", func(t *testing.T) {
		type invalidReq struct {
			Data string
		}
		mockData := invalidReq{
			Data: "error",
		}

		mockToken := &jwt.JWTclaim{
			ID:    1,
			Email: "test@test.com",
			StandardClaims: newJWT.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
			},
		}

		tokenAlgo := newJWT.NewWithClaims(newJWT.SigningMethodHS256, mockToken)

		token, err := tokenAlgo.SignedString(jwt.JWT_KEY)
		if err != nil {
			t.Error(err)
			return
		}

		validate := validator.New()
		accountUseCase := new(mocks.AccountUseCase)

		reqData, err := json.Marshal(mockData)
		if err != nil {
			t.Error(err)
			return
		}

		accountHandler := account.AccountHandler{
			Validate: validate,
			UseCase:  accountUseCase,
		}

		r := httptest.NewRequest(http.MethodPatch, "/just/for/testing", bytes.NewReader(reqData))
		r.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.UpdateAccount)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusBadRequest, rb.Status, fmt.Sprintf("Should be status %s", response.StatusBadRequest))
		assert.Nil(t, rb.Data, "Should be not nil")

		accountUseCase.AssertExpectations(t)
	})
}

func TestHandler_DeleteAccount(t *testing.T) {
	t.Run("Delete Account Success", func(t *testing.T) {
		type testReq struct {
			msg string
		}

		data := testReq{
			msg: "Success Delete Data",
		}

		resp := response.Success(response.StatusOK, data)

		accountUseCase := new(mocks.AccountUseCase)
		accountUseCase.On("DeleteAccount", mock.Anything, mock.AnythingOfType("int64")).Return(resp)

		accountHandler := account.AccountHandler{
			UseCase: accountUseCase,
		}

		mockToken := &jwt.JWTclaim{
			ID:    1,
			Email: "test@test.com",
			StandardClaims: newJWT.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
			},
		}

		tokenAlgo := newJWT.NewWithClaims(newJWT.SigningMethodHS256, mockToken)

		token, err := tokenAlgo.SignedString(jwt.JWT_KEY)
		if err != nil {
			t.Error(err)
			return
		}

		r := httptest.NewRequest(http.MethodDelete, "/just/for/testing", nil)
		r.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.DeleteAccount)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusOK, rb.Status, fmt.Sprintf("Should be status %s", response.StatusOK))
		assert.NotNil(t, rb.Data, "Should be not nil")

		accountUseCase.AssertExpectations(t)
	})

	t.Run("Delete Account Unauthorized", func(t *testing.T) {
		accountUseCase := new(mocks.AccountUseCase)

		accountHandler := account.AccountHandler{
			UseCase: accountUseCase,
		}

		r := httptest.NewRequest(http.MethodDelete, "/just/for/testing", nil)
		r.AddCookie(&http.Cookie{
			Name:     "",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
		})
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(accountHandler.DeleteAccount)
		handler.ServeHTTP(recorder, r)

		rb := response.ResponseImpl{}
		if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusUnauthorized, rb.Status, fmt.Sprintf("Should be status %s", response.StatusUnauthorized))
		assert.Nil(t, rb.Data, "Should be nil")
	})
}
