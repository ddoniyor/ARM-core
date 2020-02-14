package core

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestLogin_QueryError(t *testing.T) {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = Login("", "", db)
	// errors.Is vs errors.As
	var typedErr *QueryError
	if ok := errors.As(err, &typedErr); !ok {
		t.Errorf("error not maptch QueryError: %v", err)
	}
}

func TestLogin_NoSuchLoginForEmptyDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	// Crash Early
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query: %v", err)
	}

	result, err := Login("", "", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != false {
		t.Error("Login result not false for empty table")
	}
}

func TestLogin_LoginOk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO managers(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	result, err := Login("vasya", "secret", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != true {
		t.Error("Login result not true for existing account")
	}
}

func TestLogin_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO managers(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = Login("vasya", "password", db)
	if !errors.Is(err, ErrInvalidPass) {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}

func TestLoginClient(t *testing.T) {
	type args struct {
		login    string
		password string
		db       *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoginClient(tt.args.login, tt.args.password, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoginClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoginClient_QueryError(t *testing.T) {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = LoginClient("", "", db)
	// errors.Is vs errors.As
	var typedErr *QueryError
	if ok := errors.As(err, &typedErr); !ok {
		t.Errorf("error not maptch QueryError: %v", err)
	}
}

func TestLoginClient_NoSuchLoginForEmptyDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query: %v", err)
	}

	result, err := LoginClient("", "", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != false {
		t.Error("Login result not false for empty table")
	}
}

func TestLoginClient_LoginOk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO clients(id, login, password) VALUES (1, 'don', 'don')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	result, err := LoginClient("don", "don", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != true {
		t.Error("Login result not true for existing account")
	}
}

func TestLoginClient_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO clients(id, login, password) VALUES (1, 'don', 'don')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = LoginClient("don", "secret", db)
	if !errors.Is(err, ErrInvalidPass) {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}

