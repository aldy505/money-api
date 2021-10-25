package auth

import (
	"context"
	"strconv"
	"strings"

	"github.com/allegro/bigcache/v3"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/jmoiron/sqlx"
)

func CheckIfUserExists(user User, db *sqlx.DB, ctx context.Context) (bool, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return false, err
	}
	defer c.Close()

	var res string
	r, err := c.QueryxContext(ctx, "SELECT name FROM users WHERE email = ?", user.Email)
	if err != nil {
		return false, err
	}
	defer r.Close()

	err = r.Scan(&res)
	if err != nil {
		return false, err
	}

	return res == user.Email, nil
}

func GetUserByEmail(email string, db *sqlx.DB, ctx context.Context) (User, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return User{}, err
	}
	defer c.Close()

	var u User
	r, err := c.QueryContext(ctx, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return User{}, err
	}
	defer r.Close()

	err = sqlscan.ScanOne(&u, r)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetAllUsers(db *sqlx.DB, ctx context.Context) ([]User, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return []User{}, err
	}
	defer c.Close()

	var u []User
	r, err := c.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return []User{}, err
	}
	defer r.Close()

	err = sqlscan.ScanAll(&u, r)
	if err != nil {
		return []User{}, err
	}

	return u, nil
}

func RefreshMemory(u []User, mem *bigcache.BigCache) error {
	var users []string
	for _, v := range u {
		users = append(users, strconv.Itoa(v.ID))
	}

	err := mem.Delete("UsersID")
	if err != nil {
		return err
	}

	err = mem.Set("UsersID", []byte(strings.Join(users, ",")))
	if err != nil {
		return err
	}

	return nil
}
