package transaction

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func MoveFund(from int, to int, amount int64, db *sqlx.DB, ctx context.Context) error {
	c, err := db.Connx(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	t, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Get current balance for the from user
	bfq, err := t.QueryContext(
		ctx,
		"SELECT balance FROM accounts WHERE id = ?",
		from,
	)
	if err != nil {
		t.Rollback()
		return err
	}
	defer bfq.Close()

	var bf int64
	err = bfq.Scan(&bf)
	if err != nil {
		t.Rollback()
		return err
	}

	// Update the balance for the from user
	_, err = t.ExecContext(
		ctx,
		"UPDATE accounts SET balance = ? WHERE id = ?",
		bf-amount,
		from,
	)
	if err != nil {
		t.Rollback()
		return err
	}

	// Get current balance for the to user
	btq, err := t.QueryContext(
		ctx,
		"SELECT balance FROM accounts WHERE id = ?",
		from,
	)
	if err != nil {
		t.Rollback()
		return err
	}
	defer btq.Close()

	var bt int64
	err = btq.Scan(&bt)
	if err != nil {
		t.Rollback()
		return err
	}

	// Update the balance for the to user
	_, err = t.ExecContext(
		ctx,
		"UPDATE accounts SET balance = ? WHERE id = ?",
		bt+amount,
		from,
	)
	if err != nil {
		t.Rollback()
		return err
	}

	err = t.Commit()
	if err != nil {
		t.Rollback()
		return err
	}

	return nil
}
