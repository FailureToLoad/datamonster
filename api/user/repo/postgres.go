package repo

import (
	"context"
	"datamonster/store"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id       int
	Username string
	Hash     []byte
}

type PostGres struct {
	conn store.Connection
}

func New(d store.Connection) *PostGres {
	return &PostGres{conn: d}
}

func (pg PostGres) Insert(ctx context.Context, user User) (userId int, err error) {
	tx, err := pg.conn.Begin(ctx)
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
	err := pg.conn.QueryRow(ctx, query, username).Scan(&u.Id, &u.Username, &u.Hash)
	return u, err
}
