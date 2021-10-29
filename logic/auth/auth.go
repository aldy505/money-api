package auth

import (
	"context"
	"money-api/platform/decorator"
	"strconv"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/jmoiron/sqlx"
)

func CheckIfUserExists(user User, db *sqlx.DB, ctx context.Context) (bool, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return false, decorator.Err(err)
	}
	defer c.Close()

	var res string
	r, err := c.QueryContext(ctx, "SELECT email FROM users WHERE email = ?", user.Email)
	if err != nil {
		return false, decorator.Err(err)
	}
	defer r.Close()

	for r.Next() {
		err = r.Scan(&res)
		if err != nil {
			return false, decorator.Err(err)
		}
	}

	return res == user.Email, nil
}

func GetUserByEmail(email string, db *sqlx.DB, ctx context.Context) (User, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return User{}, decorator.Err(err)
	}
	defer c.Close()

	var u User
	r, err := c.QueryContext(ctx, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return User{}, decorator.Err(err)
	}
	defer r.Close()

	err = sqlscan.ScanOne(&u, r)
	if err != nil {
		return User{}, decorator.Err(err)
	}

	return u, nil
}

func GetAllUsers(db *sqlx.DB, ctx context.Context) ([]User, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return []User{}, decorator.Err(err)
	}
	defer c.Close()

	var u []User
	r, err := c.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return []User{}, decorator.Err(err)
	}
	defer r.Close()

	err = sqlscan.ScanAll(&u, r)
	if err != nil {
		return []User{}, decorator.Err(err)
	}

	return u, nil
}

func RegisterUser(u User, db *sqlx.DB, ctx context.Context) (User, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return User{}, decorator.Err(err)
	}
	defer c.Close()

	_, err = c.ExecContext(
		ctx,
		`INSERT INTO users
			(name, password, email, address, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?);`,
		u.Name,
		u.Password,
		u.Email,
		u.Address,
		time.Now().Unix(),
		time.Now().Unix(),
	)
	if err != nil {
		return User{}, decorator.Err(err)
	}

	r, err := c.QueryContext(ctx, "SELECT * FROM users WHERE email = ?", u.Email)
	if err != nil {
		return User{}, decorator.Err(err)
	}
	defer r.Close()

	var user User
	err = sqlscan.ScanOne(&user, r)
	if err != nil {
		return User{}, decorator.Err(err)
	}

	return user, nil
}

func RefreshMemory(u []User, mem *bigcache.BigCache) error {
	var users []string
	for _, v := range u {
		users = append(users, strconv.Itoa(v.ID))
	}

	err := mem.Delete("UsersID")
	if err != nil {
		return decorator.Err(err)
	}

	err = mem.Set("UsersID", []byte(strings.Join(users, ",")))
	if err != nil {
		return decorator.Err(err)
	}

	return nil
}
