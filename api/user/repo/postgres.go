package repo

import (
	"context"
	"datamonster/dao"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id       int
	Username string
	Hash     []byte
}

type PostGres struct {
	dao dao.Connection
}

func New(d dao.Connection) *PostGres {
	return &PostGres{dao: d}
}

func (pg PostGres) Insert(ctx context.Context, user User) (userId int, err error) {
	tx, err := pg.dao.Begin(ctx)
	if err != nil {
		return 0, err
	}

	insert := "INSERT INTO private.user (username, hash) VALUES(@username, @hash) RETURNING id"
	args := pgx.NamedArgs{
		"username": user.Username,
		"hash":     user.Hash,
	}
	result := tx.QueryRow(ctx, insert, args)
	err = result.Scan(&userId)
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
	return userId, err
}

func (pg PostGres) Get(ctx context.Context, username string) (User, error) {
	query := `SELECT * FROM private.user WHERE username = $1`
	var u User
	err := pg.dao.QueryRow(ctx, query, username).Scan(&u.Id, &u.Username, &u.Hash)
	return u, err
}
