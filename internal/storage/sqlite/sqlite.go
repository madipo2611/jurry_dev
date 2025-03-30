package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"jurry_dev/internal/storage"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	st, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    password CHAR(120) NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    balans INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    role VARCHAR(10) NOT NULL DEFAULT 'normal', -- Изменено с enum на varchar
    last_seen TIMESTAMP NULL DEFAULT NULL,
    gender VARCHAR(50) NOT NULL,
    language VARCHAR(10) NOT NULL,
    active_status_online VARCHAR(3) NOT NULL DEFAULT 'yes', -- Изменено с enum на varchar
    posts_privacy TINYINT(1) NOT NULL DEFAULT 1,
    allow_dm TINYINT(1) NOT NULL DEFAULT 1,
    allow_comments TINYINT(1) NOT NULL DEFAULT 1
)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = st.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Register(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO users(login, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Login(login string, password string) (bool, error) {
	const op = "storage.sqlite.Login"

	stmt, err := s.db.Prepare("SELECT login, password FROM users WHERE login = ? AND password = ?;")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	var dbLogin, dbPass string

	err = stmt.QueryRow(login, password).Scan(&dbLogin, &dbPass)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Errorf("Ошибка: %s: %w", op, err)
		return false, storage.ErrLoginNotFound
	}
	if err != nil {
		fmt.Errorf("login error: %s, %s", dbLogin, dbPass)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil

}
