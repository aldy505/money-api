package transaction

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/jmoiron/sqlx"
)

var ErrInvalidSwitcher = errors.New("switcher must be one of: user, id, friend")

// Checks whether a transaction by a certain ID is exists or not.
func IsTransactionExists(id int, db *sqlx.DB, ctx context.Context, mem *bigcache.BigCache) (bool, error) {
	cache, err := mem.Get("TransactionIDs")
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
		c, err := db.Connx(ctx)
		if err != nil {
			return false, err
		}
		defer c.Close()

		r, err := c.QueryContext(ctx, "SELECT id FROM transactions")
		if err != nil {
			return false, err
		}
		defer r.Close()

		var t []Transaction
		err = sqlscan.ScanAll(&t, r)
		if err != nil {
			return false, err
		}

		err = RefreshMemory(t, mem)
		if err != nil {
			return false, err
		}

		for _, v := range t {
			if v.ID == id {
				return true, nil
			}
		}
		return false, nil
	}

	return false, err
}

// Switcher (1st param argument) accepts one of: "user", "id" or "friend".
// Which will be used by the SQL query from the `t Transaction` readout.
func GetTransactionBy(switcher string, t Transaction, db *sqlx.DB, ctx context.Context) ([]Transaction, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return []Transaction{}, err
	}
	defer c.Close()

	switch switcher {
	case "user":
		r, err := c.QueryContext(ctx, "SELECT * FROM transactions WHERE sender = ?", t.Sender.ID)
		if err != nil {
			return []Transaction{}, err
		}
		defer r.Close()

		var t []Transaction
		err = sqlscan.ScanOne(&t, r)
		if err != nil {
			return []Transaction{}, err
		}

		return t, err
	case "id":
		r, err := c.QueryContext(ctx, "SELECT * FROM transactions WHERE id = ?", t.ID)
		if err != nil {
			return []Transaction{}, err
		}
		defer r.Close()

		var t []Transaction
		err = sqlscan.ScanAll(&t, r)
		if err != nil {
			return []Transaction{}, err
		}

		return t, err

	case "friend":
		r, err := c.QueryContext(ctx, "SELECT * FROM transactions WHERE sender = ? AND recipient = ?", t.Sender.ID, t.Recipient.ID)
		if err != nil {
			return []Transaction{}, err
		}
		defer r.Close()

		var t []Transaction
		err = sqlscan.ScanAll(&t, r)
		if err != nil {
			return []Transaction{}, err
		}

		return t, err

	default:
		return []Transaction{}, ErrInvalidSwitcher
	}
}

// Returns the transaction ID
func CreateTransaction(i Intermmediate, db *sqlx.DB, ctx context.Context) (int, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer c.Close()

	r, err := c.QueryContext(
		ctx,
		`INSERT INTO transactions
		  (sender, recipient, amount, message, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			RETURNING id;`,
		i.Sender,
		i.Recipient,
		i.Message,
		i.Status,
		time.Now().Unix(),
		time.Now().Unix(),
	)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	var t int
	err = r.Scan(&t)
	if err != nil {
		return 0, err
	}

	return t, nil
}

// Only updates the status for now. Returns the transaction ID.
func UpdateTransaction(t Transaction, db *sqlx.DB, ctx context.Context) (int, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer c.Close()

	_, err = c.ExecContext(
		ctx,
		`UPDATE transactions
			SET status = ?, updated_at = ?
			WHERE id = ?;`,
		t.Status,
		time.Now().Unix(),
		t.ID,
	)
	if err != nil {
		return 0, err
	}
	return t.ID, nil
}

// Returns the transaction ID.
func CreateRequest(i Intermmediate, db *sqlx.DB, ctx context.Context) (int, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer c.Close()

	r, err := c.QueryContext(
		ctx,
		`INSERT INTO transactions
			(sender, recipient, amount, message, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			RETURNING id;`,
		i.Sender,
		i.Recipient,
		i.Amount,
		i.Message,
		StatusRequested,
		time.Now().Unix(),
		time.Now().Unix(),
	)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	var t int
	err = r.Scan(&t)
	if err != nil {
		return 0, err
	}

	return t, nil
}

func RefreshMemory(t []Transaction, mem *bigcache.BigCache) error {
	var ids []string
	for _, v := range t {
		ids = append(ids, strconv.Itoa(v.ID))
	}

	err := mem.Delete("TransactionIDs")
	if err != nil {
		return err
	}

	err = mem.Set("TransactionIDs", []byte(strings.Join(ids, ",")))
	if err != nil {
		return err
	}

	return nil
}
