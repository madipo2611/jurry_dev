package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"jurry_dev/internal/storage"
	_ "modernc.org/sqlite"
	"time"
)

type Storage struct {
	db *sql.DB
}

type Posts struct {
	Id        int
	UserID    int
	Image     string
	Text      string
	Likes     int
	CreatedAt time.Time
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

func (s *Storage) Login(login string) (string, int, error) {
	const op = "storage.sqlite.Login"

	stmt, err := s.db.Prepare("SELECT password, id FROM users WHERE login = $1;")
	if err != nil {
		return "", 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	var dbPass string
	var userID int

	err = stmt.QueryRow(login).Scan(&dbPass, &userID)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Errorf("Ошибка: %s: %w", op, err)
		return "", 0, storage.ErrLoginNotFound
	}
	if err != nil {
		fmt.Errorf("login error: %s, %s", dbPass)
		return "", 0, fmt.Errorf("%s: %w", op, err)
	}
	return dbPass, userID, nil

}

func (s *Storage) Register(login string, password string, gender string) (bool, error) {
	const op = "storage.sqlite.Register"

	stmt, err := s.db.Prepare("INSERT INTO users(login, password, gender) VALUES($1, $2, $3)")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(login, password, gender)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil

}

func (s *Storage) AddPost(text string, image string, user int) (bool, error) {
	const op = "storage.sqlite.Register"

	stmt, err := s.db.Prepare("INSERT INTO posts(text, image, userID) VALUES($1, $2, $3)")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(text, image, user)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func (s *Storage) GetPost(offset, limit int) ([]Posts, error) {
	const op = "storage.sqlite.GetPost"

	stmt, err := s.db.Query("SELECT * FROM posts ORDER BY created_at DESC OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	var data []Posts
	for stmt.Next() {
		var post Posts
		if err := stmt.Scan(&post.Id, &post.UserID, &post.Image, &post.Text, &post.Likes, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		data = append(data, post)
		return data, nil
	}
	return data, nil
}
