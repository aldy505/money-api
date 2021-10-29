package account

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/jmoiron/sqlx"
)

func AddBalance(a Account, db *sqlx.DB, ctx context.Context) (Account, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return Account{}, err
	}
	defer c.Close()

	t, err := c.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return Account{}, err
	}

	b, err := t.QueryContext(ctx, "SELECT balance FROM accounts WHERE id = ?", a.ID)
	if err != nil {
		t.Rollback()
		return Account{}, err
	}
	defer b.Close()

	var balance int64
	for b.Next() {
		err = b.Scan(&balance)
		if err != nil {
			t.Rollback()
			return Account{}, err
		}
	}

	_, err = t.ExecContext(
		ctx,
		"UPDATE accounts SET balance = ? WHERE id = ?",
		balance+a.Balance,
		a.ID,
	)
	if err != nil {
		t.Rollback()
		return Account{}, err
	}

	err = t.Commit()
	if err != nil {
		t.Rollback()
		return Account{}, err
	}

	updatedAccount, err := c.QueryContext(ctx, "SELECT * FROM accounts WHERE id = ?", a.ID)
	if err != nil {
		return Account{}, err
	}
	defer updatedAccount.Close()

	var o Account
	err = sqlscan.ScanOne(&o, updatedAccount)
	if err != nil {
		return Account{}, err
	}

	return o, nil
}

func SubtractBalance(a Account, db *sqlx.DB, ctx context.Context) (Account, error) {
	c, err := db.Connx(ctx)
	if err != nil {
		return Account{}, err
	}
	defer c.Close()

	t, err := c.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return Account{}, err
	}

	b, err := t.QueryContext(ctx, "SELECT balance FROM accounts WHERE id = ?", a.ID)
	if err != nil {
		t.Rollback()
		return Account{}, err
	}
	defer b.Close()

	var balance int64
	for b.Next() {
		err = b.Scan(&balance)
		if err != nil {
			t.Rollback()
			return Account{}, err
		}
	}

	_, err = t.ExecContext(
		ctx,
		"UPDATE accounts SET balance = ? WHERE id = ?",
		balance-a.Balance,
		a.ID,
	)
	if err != nil {
		t.Rollback()
		return Account{}, err
	}

	err = t.Commit()
	if err != nil {
		t.Rollback()
		return Account{}, err
	}

	updatedAccount, err := c.QueryContext(ctx, "SELECT * FROM accounts WHERE id = ?", a.ID)
	if err != nil {
		return Account{}, err
	}
	defer updatedAccount.Close()

	var o Account
	err = sqlscan.ScanOne(&o, updatedAccount)
	if err != nil {
		return Account{}, err
	}

	return o, nil
}
