package account

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"waizly/helpers/exception"
	"waizly/models"
)

type (
	AccountRepository interface {
		Create(ctx context.Context, params models.Account) (int64, error)
		FindByID(ctx context.Context, id int64) (models.Account, error)
		FindByEmail(ctx context.Context, email string) (models.Account, error)
		Update(ctx context.Context, id int64, params models.Account) error
		Delete(ctx context.Context, id int64) error
	}

	accountRepositoryImpl struct {
		db        *sql.DB
		tableName string
	}
)

func NewAccountRepository(db *sql.DB, tableName string) AccountRepository {
	return &accountRepositoryImpl{
		db:        db,
		tableName: tableName,
	}
}

func (ar *accountRepositoryImpl) Create(ctx context.Context, params models.Account) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, password, email, created_at) VALUES (?, ?, ?, ?)", ar.tableName)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		params.Username,
		params.Password,
		params.Email,
		params.CreatedAt,
	)

	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	ID, err := result.LastInsertId()

	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	return ID, nil
}

func (ar *accountRepositoryImpl) FindByID(ctx context.Context, id int64) (models.Account, error) {
	account := models.Account{}

	query := fmt.Sprintf(`SELECT id, username, password, email, created_at, update_at FROM %s WHERE id = ?`, ar.tableName)

	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return account, exception.ErrInternalServer
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	var password sql.NullString
	var updateAt sql.NullTime

	err = row.Scan(
		&account.ID,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.CreatedAt,
		&account.UpdateAt,
	)

	if err != nil {
		log.Println(err)
		return account, exception.ErrInternalServer
	}

	if password.Valid {
		account.Password = password.String
	}

	if updateAt.Valid {
		account.UpdateAt = updateAt.Time
	}

	return account, nil
}

func (ar *accountRepositoryImpl) FindByEmail(ctx context.Context, email string) (models.Account, error) {
	account := models.Account{}

	query := fmt.Sprintf(`SELECT id, username, password, email, created_at, update_at FROM %s WHERE email = ?`, ar.tableName)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return account, exception.ErrInternalServer
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	var password sql.NullString
	var updateAt sql.NullTime

	err = row.Scan(
		&account.ID,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.CreatedAt,
		&account.UpdateAt,
	)

	if err != nil {
		log.Println(err)

		return account, exception.ErrNotFound
	}

	if password.Valid {
		account.Password = password.String
	}

	if updateAt.Valid {
		account.UpdateAt = updateAt.Time
	}

	return account, nil
}

func (ar *accountRepositoryImpl) Update(ctx context.Context, id int64, params models.Account) error {
	query := fmt.Sprintf(`UPDATE %s SET username = ?, password = ?, email = ?, update_at = ? WHERE id = %d`, ar.tableName, id)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		params.Username,
		params.Password,
		params.Email,
		params.UpdateAt,
	)

	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected < 1 {
		return exception.ErrNotFound
	}

	return nil
}

func (ar *accountRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = %d`, ar.tableName, id)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
	)

	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected < 1 {
		return exception.ErrNotFound
	}

	return nil
}
