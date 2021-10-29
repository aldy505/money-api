package account

import (
	"context"
	"errors"
	"money-api/platform/decorator"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/jmoiron/sqlx"
)

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
			return false, decorator.Err(err)
		}

		err = RefreshMemory(accounts, mem)
		if err != nil {
			return false, decorator.Err(err)
		}

		for _, v := range accounts {
			if v.Tag == tag {
				return true, nil
			}
		}

		return false, nil
	}

	return false, decorator.Err(err)
}

// Returns the ID of the account
func UpdateTag(a Account, db *sqlx.DB, ctx context.Context) (int, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return 0, decorator.Err(err)
	}
	defer c.Close()

	_, err = c.ExecContext(
		ctx,
		`UPDATE accounts SET tag=?, updated_at=? WHERE id=?`,
		a.Tag,
		time.Now().Unix(),
		a.ID,
	)
	if err != nil {
		return 0, decorator.Err(err)
	}

	return a.ID, nil
}
