package mysql

import (
	"be-service-auth/domain"
	"context"
	"database/sql"
	"errors"
)

type mysqlAuthRepository struct {
	Conn *sql.DB
}

func NewMySQLAuthRepository(Conn *sql.DB) domain.AuthMySQLRepo {
	return &mysqlAuthRepository{Conn}
}

func (db *mysqlAuthRepository) GetUserByEmail(ctx context.Context, email string) (response domain.ResponseLoginDTO, err error) {
	query := `SELECT id, email, password, role
	FROM auth
	WHERE email = ?;`

	row := db.Conn.QueryRowContext(ctx, query, email)

	var user domain.ResponseLoginDTO
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("Not found")
		}
		return
	}

	return user, nil
}
