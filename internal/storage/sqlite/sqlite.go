package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"jurry_dev/internal/storage"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(ps string) (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем соединение с базой данных
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Login(login string) (string, string, error) {
	const op = "storage.sqlite.Login"

	stmt, err := s.db.Prepare("SELECT password, id FROM users WHERE login = $1;")
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	var dbPass, userID string

	err = stmt.QueryRow(login).Scan(&dbPass, &userID)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Errorf("Ошибка: %s: %w", op, err)
		return "", "", storage.ErrLoginNotFound
	}
	if err != nil {
		fmt.Errorf("login error: %s, %s", dbPass)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return dbPass, userID, nil

}

func (s *Storage) Register(login string, password string, gender string) (bool, error) {
	const op = "storage.sqlite.Register"

	stmt, err := s.db.Prepare("INSERT INTO users(login, password, gender) VALUES($1, $2, $3)")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(login, password, gender)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil

}

func (s *Storage) AddPost(text string, image string, user string) (bool, error) {
	const op = "storage.sqlite.Register"

	stmt, err := s.db.Prepare("INSERT INTO posts(text, image, userID) VALUES($1, $2, $3)")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(text, image, user)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}
