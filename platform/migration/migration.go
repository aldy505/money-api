package migration

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB, ctx context.Context) error {
	err := tableUsers(db, ctx)
	if err != nil {
		return err
	}

	err = tableAccounts(db, ctx)
	if err != nil {
		return err
	}

	err = tableTransactions(db, ctx)
	if err != nil {
		return err
	}

	err = tableFriends(db, ctx)
	if err != nil {
		return err
	}

	return nil
}

func tableUsers(db *sqlx.DB, ctx context.Context) error {
	c, err := db.Connx(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (
		id          INTEGER                 PRIMARY KEY AUTOINCREMENT,
		name        VARCHAR(255)  NOT NULL,
		password    TEXT          NOT NULL,
		email       VARCHAR(255)  NOT NULL 	UNIQUE,
		address     TEXT,
		updated_at  DATETIME      NOT NULL,
		created_at  DATETIME      NOT NULL
	)`)
	if err != nil {
		return err
	}

	return nil
}

func tableAccounts(db *sqlx.DB, ctx context.Context) error {
	c, err := db.Connx(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS accounts (
		id         	INTEGER            	  PRIMARY KEY,
		balance    	MEDIUMINT   NOT NULL  DEFAULT 0,
		tag        	TEXT               	  UNIQUE,
		updated_at  DATETIME 	  NOT NULL,
		created_at  DATETIME 	  NOT NULL,
		FOREIGN KEY (id)        REFERENCES users (id)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func tableTransactions(db *sqlx.DB, ctx context.Context) error {
	c, err := db.Connx(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	t, err := c.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = t.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS transactions (
		id         INTEGER              PRIMARY KEY AUTOINCREMENT,
		sender     INTEGER    NOT NULL,
	  recipient  INTEGER    NOT NULL,
		amount     MEDIUMINT  NOT NULL,
		message    TEXT,
		status     SMALLINT   NOT NULL DEFAULT 0,
		updated_at DATETIME   NOT NULL,
		created_at DATETIME   NOT NULL,

		FOREIGN KEY (sender)      REFERENCES users (id),
		FOREIGN KEY (recipient)   REFERENCES users (id)
	)`)
	if err != nil {
		t.Rollback()
		return err
	}

	_, err = t.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_sender ON transactions (sender)`)
	if err != nil {
		t.Rollback()
		return err
	}

	_, err = t.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_recipient ON transactions (recipient)`)
	if err != nil {
		t.Rollback()
		return err
	}

	err = t.Commit()
	if err != nil {
		return err
	}

	return nil
}

func tableFriends(db *sqlx.DB, ctx context.Context) error {
	c, err := db.Connx(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	t, err := c.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = t.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS friends (
		id      INTEGER   NOT NULL,
		friend  INTEGER   NOT NULL,
		FOREIGN KEY (id)     REFERENCES users (id)
		FOREIGN KEY (friend) REFERENCES users (id)
	)`)
	if err != nil {
		t.Rollback()
		return err
	}
	_, err = t.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_id ON friends (id)`)
	if err != nil {
		t.Rollback()
		return err
	}
	_, err = t.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_friend ON friends (friend)`)
	if err != nil {
		t.Rollback()
		return err
	}

	err = t.Commit()
	if err != nil {
		return err
	}

	return nil
}
