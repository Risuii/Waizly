package account_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"waizly/internal/account"
	"waizly/internal/constant"
	"waizly/internal/mock"
	"waizly/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var currentTime = time.Date(2021, 12, 12, 0, 0, 0, 0, &time.Location{})
var accountStruct = models.Account{
	ID:        1,
	Username:  "username-test",
	Password:  "password-test",
	Email:     "email@test.com",
	CreatedAt: currentTime,
	UpdateAt:  currentTime,
}

func TestCreat(t *testing.T) {
	t.Run("Test Create Success", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`INSERT INTO %s`, constant.TableAccount)
		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(accountStruct.Username, accountStruct.Password, accountStruct.Email, accountStruct.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

		ID, err := repo.Create(ctx, accountStruct)

		assert.Equal(t, int64(1), ID)
		assert.NoError(t, err)
	})

	t.Run("Test Create Error", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`INSERT INTO %s`, constant.TableAccount)
		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(accountStruct.Username, accountStruct.Password, accountStruct.Email, accountStruct.CreatedAt).WillReturnResult(sqlmock.NewResult(0, 0))

		ID, err := repo.Create(ctx, accountStruct)

		assert.Equal(t, int64(0), ID)
		assert.NoError(t, err)
	})
}

func TestFindByID(t *testing.T) {
	t.Run("Test FindByID Success", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`SELECT id, username, password, email, created_at, update_at FROM %s WHERE id = ?`, constant.TableAccount)
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at", "update_at"}).AddRow(accountStruct.ID, accountStruct.Username, accountStruct.Password, accountStruct.Email, accountStruct.CreatedAt, accountStruct.UpdateAt)

		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(accountStruct.ID).WillReturnRows(rows)

		accountStruct, err := repo.FindByID(ctx, accountStruct.ID)

		assert.NotNil(t, accountStruct)
		assert.NoError(t, err)
	})

	t.Run("Test FindByID Error", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`SELECT id, username, password, email, created_at, update_at FROM %s WHERE id = ?`, constant.TableAccount)
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at", "update_at"})

		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(accountStruct.ID).WillReturnRows(rows)

		accountStruct, err := repo.FindByID(ctx, accountStruct.ID)

		assert.Empty(t, accountStruct)
		assert.Error(t, err)
	})
}

func TestFindByEmail(t *testing.T) {
	t.Run("Test FindByEmail Success", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`SELECT id, username, password, email, created_at, update_at FROM %s WHERE email = ?`, constant.TableAccount)
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at", "update_at"}).AddRow(accountStruct.ID, accountStruct.Username, accountStruct.Password, accountStruct.Email, accountStruct.CreatedAt, accountStruct.UpdateAt)

		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(accountStruct.Email).WillReturnRows(rows)

		accountStruct, err := repo.FindByEmail(ctx, accountStruct.Email)

		assert.NotNil(t, accountStruct)
		assert.NoError(t, err)
	})

	t.Run("Test FindByEmail Error", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`SELECT id, username, password, email, created_at, update_at FROM %s WHERE email = ?`, constant.TableAccount)
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at", "update_at"})

		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(accountStruct.Email).WillReturnRows(rows)

		accountStruct, err := repo.FindByEmail(ctx, accountStruct.Email)

		assert.Empty(t, accountStruct)
		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update Success", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`UPDATE %s SET`, constant.TableAccount)

		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(accountStruct.Username, accountStruct.Password, accountStruct.Email, accountStruct.UpdateAt).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(ctx, accountStruct.ID, accountStruct)

		assert.NoError(t, err)
	})

	t.Run("Test Update Error", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`UPDATE %s SET`, constant.TableAccount)

		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectExec().WithArgs(accountStruct.Username, accountStruct.Password, accountStruct.Email, accountStruct.UpdateAt).WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(ctx, accountStruct.ID, accountStruct)

		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Delete Success", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, constant.TableAccount)
		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, accountStruct.ID)

		assert.NoError(t, err)
	})

	t.Run("Test Delete Error", func(t *testing.T) {
		db, mock := mock.NewMock()
		repo := account.NewAccountRepository(db, constant.TableAccount)

		defer db.Close()

		query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, constant.TableAccount)
		ctx := context.TODO()

		mock.ExpectPrepare(query).ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(ctx, accountStruct.ID)

		assert.Error(t, err)
	})
}
