package account

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/allegro/bigcache/v3"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/jmoiron/sqlx"
)

var ErrInvalidSwitcher = errors.New("switcher must be one of: id, tag")

func GetAccountBy(switcher string, a Account, db *sqlx.DB, ctx context.Context) (Account, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return Account{}, err
	}
	defer c.Close()

	switch switcher {
	case "id":
		r, err := c.QueryContext(ctx, "SELECT * FROM accounts WHERE id=? LIMIT 1", a.ID)
		if err != nil {
			return Account{}, err
		}
		defer r.Close()

		var a Account
		err = sqlscan.ScanOne(&a, r)
		if err != nil {
			return Account{}, err
		}

		f, err := GetFriends(a.ID, db, ctx)
		if err != nil {
			return Account{}, err
		}
		a.Friends = f

		return a, nil

	case "tag":
		r, err := c.QueryContext(ctx, "SELECT * FROM accounts WHERE tag=? LIMIT 1", a.Tag)
		if err != nil {
			return Account{}, err
		}
		defer r.Close()

		var a Account
		err = sqlscan.ScanOne(&a, r)
		if err != nil {
			return Account{}, err
		}
		f, err := GetFriends(a.ID, db, ctx)
		if err != nil {
			return Account{}, err
		}
		a.Friends = f

		return a, nil

	default:
		return Account{}, ErrInvalidSwitcher
	}
}

func GetAllAccounts(db *sqlx.DB, ctx context.Context) ([]Account, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return []Account{}, err
	}
	defer c.Close()

	r, err := c.QueryContext(ctx, "SELECT * FROM accounts")
	if err != nil {
		return []Account{}, err
	}
	defer r.Close()

	var a []Account
	err = sqlscan.ScanAll(&a, r)
	if err != nil {
		return []Account{}, err
	}

	return a, nil
}

func IsIDExists(id int, db *sqlx.DB, ctx context.Context, mem *bigcache.BigCache) (bool, error) {
	cache, err := mem.Get("AccountIDs")
	if err == nil {
		ids := strings.Split(string(cache), ",")
		for _, v := range ids {
			s, err := strconv.Atoi(v)
			if err != nil {
				return false, err
			}

			if s == id {
				return true, nil
			}
		}
		return false, nil
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		accounts, err := GetAllAccounts(db, ctx)
		if err != nil {
			return false, err
		}

		err = RefreshMemory(accounts, mem)
		if err != nil {
			return false, err
		}

		for _, v := range accounts {
			if v.ID == id {
				return true, nil
			}
		}

		return false, nil
	}

	return false, err
}

func IsTagExists(tag string, db *sqlx.DB, ctx context.Context, mem *bigcache.BigCache) (bool, error) {
	cache, err := mem.Get("AccountTags")
	if err == nil {
		ids := strings.Split(string(cache), ",")
		for _, v := range ids {
			if v == tag {
				return true, nil
			}
		}
		return false, nil
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		accounts, err := GetAllAccounts(db, ctx)
		if err != nil {
			return false, err
		}

		err = RefreshMemory(accounts, mem)
		if err != nil {
			return false, err
		}

		for _, v := range accounts {
			if v.Tag == tag {
				return true, nil
			}
		}

		return false, nil
	}

	return false, err
}

func RefreshMemory(a []Account, mem *bigcache.BigCache) error {
	var ids []string
	var tags []string

	for _, v := range a {
		ids = append(ids, strconv.Itoa(v.ID))
		tags = append(tags, v.Tag)
	}

	err := mem.Delete("AccountsID")
	if err != nil {
		return err
	}

	err = mem.Delete("AccountTags")
	if err != nil {
		return err
	}

	err = mem.Set("AccountIDs", []byte(strings.Join(ids, ",")))
	if err != nil {
		return err
	}

	err = mem.Set("AccountTags", []byte(strings.Join(tags, ",")))
	if err != nil {
		return err
	}
	return nil
}
