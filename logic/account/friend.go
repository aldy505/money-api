package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/jmoiron/sqlx"
)

func GetFriends(whom int, db *sqlx.DB, ctx context.Context) ([]Friend, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return []Friend{}, err
	}
	defer c.Close()

	r, err := c.QueryContext(ctx, "SELECT with FROM friends WHERE id = ?", whom)
	if err != nil {
		return []Friend{}, err
	}
	defer r.Close()

	var f []Friend
	err = sqlscan.ScanAll(&f, r)
	if err != nil {
		return []Friend{}, err
	}

	return f, nil
}

func AddFriend(whom int, with int, db *sqlx.DB, ctx context.Context) (Account, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return Account{}, err
	}
	defer c.Close()

	_, err = c.ExecContext(ctx, "INSERT INTO friends (id, friend) VALUES (?, ?)", whom, with)
	if err != nil {
		return Account{}, err
	}

	a, err := GetAccountBy("id", Account{ID: whom}, db, ctx)
	if err != nil {
		return Account{}, err
	}

	return a, nil
}

func RemoveFriend(whom int, with int, db *sqlx.DB, ctx context.Context) (Account, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return Account{}, err
	}
	defer c.Close()

	_, err = c.ExecContext(ctx, "DELETE FROM friends WHERE id = ? AND friend = ?", whom, with)
	if err != nil {
		return Account{}, err
	}

	a, err := GetAccountBy("id", Account{ID: whom}, db, ctx)
	if err != nil {
		return Account{}, err
	}

	return a, nil
}

func IsFriend(whom int, with int, db *sqlx.DB, ctx context.Context) (bool, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return false, err
	}
	defer c.Close()

	r, err := c.QueryContext(ctx, "SELECT friend FROM friends WHERE id = ? AND friend = ?", whom, with)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	defer r.Close()

	var f int
	err = r.Scan(&f)
	if err != nil {
		return false, err
	}

	return f == with, nil
}
