package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	bcryptmocks "waizly/config/bcrypt/mocks"
	"waizly/helpers/exception"
	"waizly/internal/account"
	"waizly/internal/account/mocks"
	"waizly/models"
)

func TestRegister(t *testing.T) {
	t.Run("Success Register", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		registerRepository := new(mocks.AccountRepository)

		registerRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, exception.ErrNotFound)
		registerRepository.On("Create", mock.Anything, mock.AnythingOfType("models.Account")).Return(int64(1), nil)
		bcrypt.On("HashPassword", mock.AnythingOfType("string")).Return("hashed password", nil)

		registerUseCase := account.NewAccountUseCase(
			registerRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.RegisterRequest{
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
		}

		resp := registerUseCase.Register(ctx, params)

		assert.NoError(t, resp.Err())

		registerRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Error Hash Password", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		registerRepository := new(mocks.AccountRepository)

		registerRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, exception.ErrNotFound)
		// registerRepository.On("Create", mock.Anything, mock.AnythingOfType("models.Account")).Return(int64(1), nil)
		bcrypt.On("HashPassword", mock.AnythingOfType("string")).Return("", exception.ErrInternalServer)

		registerUseCase := account.NewAccountUseCase(
			registerRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.RegisterRequest{
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
		}

		resp := registerUseCase.Register(ctx, params)

		assert.Error(t, resp.Err())

		registerRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Error Create", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		registerRepository := new(mocks.AccountRepository)

		registerRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, exception.ErrNotFound)
		registerRepository.On("Create", mock.Anything, mock.AnythingOfType("models.Account")).Return(int64(0), exception.ErrInternalServer)
		bcrypt.On("HashPassword", mock.AnythingOfType("string")).Return("", nil)

		accountUseCase := account.NewAccountUseCase(
			registerRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.RegisterRequest{
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
		}

		resp := accountUseCase.Register(ctx, params)

		assert.Error(t, resp.Err())
		registerRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Conflict Error", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		registerRepository := new(mocks.AccountRepository)

		registerRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, nil)

		accountUseCase := account.NewAccountUseCase(
			registerRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.RegisterRequest{
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
		}

		resp := accountUseCase.Register(ctx, params)

		assert.Error(t, resp.Err())
		registerRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)

	})

	t.Run("Error Query To DB", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		registerRepository := new(mocks.AccountRepository)

		registerRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, exception.ErrInternalServer)

		accountUseCase := account.NewAccountUseCase(
			registerRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.RegisterRequest{
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
		}

		resp := accountUseCase.Register(ctx, params)

		assert.Error(t, resp.Err())
		registerRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Account Not Found", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		loginRepository := new(mocks.AccountRepository)

		loginRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, exception.ErrNotFound)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.LoginRequest{
			Email: "email@test.com",
		}

		resp, _ := accountUseCase.Login(ctx, params)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Error query to DB", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		loginRepository := new(mocks.AccountRepository)

		loginRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, exception.ErrInternalServer)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.LoginRequest{
			Email: "email@test.com",
		}

		resp, _ := accountUseCase.Login(ctx, params)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Password not valid", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		loginRepository := new(mocks.AccountRepository)

		password := "hashed"

		mockAccount := models.Account{
			Password: password,
		}

		bcrypt.On("ComparePasswordHash", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false)

		loginRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockAccount, nil)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.LoginRequest{
			Email: "email@test.com",
		}

		resp, _ := accountUseCase.Login(ctx, params)

		assert.NoError(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Token Empty", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		loginRepository := new(mocks.AccountRepository)

		loginRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(models.Account{}, nil)
		bcrypt.On("ComparePasswordHash", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(true)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()
		params := models.LoginRequest{
			Email: "email@test.com",
		}

		resp, token := accountUseCase.Login(ctx, params)

		token = models.Token{
			Token: "",
		}

		assert.NoError(t, resp.Err())
		assert.Empty(t, token)

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Login Success", func(t *testing.T) {
		bcrypt := new(bcryptmocks.Bcrypt)
		loginRepository := new(mocks.AccountRepository)

		password := "hashed"

		mockAccount := models.Account{
			Password: password,
		}

		loginRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockAccount, nil)
		bcrypt.On("ComparePasswordHash", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(true)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()
		params := models.LoginRequest{
			Email: "email@test.com",
		}

		resp, token := accountUseCase.Login(ctx, params)

		token = models.Token{
			Token: "jwt-token-test",
		}

		assert.NoError(t, resp.Err())
		assert.NotEmpty(t, token)

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})
}

func TestDetailAccount(t *testing.T) {
	t.Run("Account Not Found", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Account{}, exception.ErrParams)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.Account{
			ID: 1,
		}

		resp := accountUseCase.DetailAccount(ctx, params.ID)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Query error to DB", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Account{}, exception.ErrInternalServer)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.Account{
			ID: 1,
		}

		resp := accountUseCase.DetailAccount(ctx, params.ID)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Get detail account success", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Account{}, nil)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.Account{
			ID: 1,
		}

		resp := accountUseCase.DetailAccount(ctx, params.ID)

		assert.NoError(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})
}

func TestUpdateAccount(t *testing.T) {

	t.Run("Account Not Found", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Account{}, exception.ErrNotFound)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()
		params := models.Account{
			ID:       1,
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
			UpdateAt: time.Time{},
		}

		resp := accountUseCase.UpdateAccount(ctx, params.ID, params)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Query error to DB", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Account{}, exception.ErrInternalServer)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()
		params := models.Account{
			ID:       1,
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
			UpdateAt: time.Time{},
		}

		resp := accountUseCase.UpdateAccount(ctx, params.ID, params)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Update Success", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Account{}, nil)
		loginRepository.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("models.Account")).Return(nil)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()
		params := models.Account{
			ID:       1,
			Username: "username-test",
			Password: "password-test",
			Email:    "email@test.com",
			UpdateAt: time.Time{},
		}

		resp := accountUseCase.UpdateAccount(ctx, params.ID, params)

		assert.NoError(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})
}

func TestDeleteAcco(t *testing.T) {
	t.Run("Account Not Found", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(exception.ErrNotFound)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.Account{
			ID: 1,
		}

		resp := accountUseCase.DeleteAccount(ctx, params.ID)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Query error to DB", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(exception.ErrInternalServer)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.Account{
			ID: 1,
		}

		resp := accountUseCase.DeleteAccount(ctx, params.ID)

		assert.Error(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})

	t.Run("Delete account success", func(t *testing.T) {
		loginRepository := new(mocks.AccountRepository)
		bcrypt := new(bcryptmocks.Bcrypt)

		loginRepository.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		accountUseCase := account.NewAccountUseCase(
			loginRepository,
			bcrypt,
		)

		ctx := context.TODO()

		params := models.Account{
			ID: 1,
		}

		resp := accountUseCase.DeleteAccount(ctx, params.ID)

		assert.NoError(t, resp.Err())

		loginRepository.AssertExpectations(t)
		bcrypt.AssertExpectations(t)
	})
}
