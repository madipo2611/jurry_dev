package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"jurry_dev/internal/storage"
	"log/slog"
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

type User struct {
	Id                   int
	Name                 string
	Login                string
	Balans               int
	Status               string
	Role                 string
	Last_seen            time.Time
	Gender               string
	Language             string
	Active_status_online bool
	Posts_privacy        int
	Allow_dm             int
	Allow_comments       int
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

func (s *Storage) GetPost(offset, limit int) ([]Posts, int, error) {
	const op = "storage.sqlite.GetPost"

	// Получаем общее количество постов
	var totalCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем посты с пагинацией
	stmt, err := s.db.Query("SELECT * FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var data []Posts
	for stmt.Next() {
		var post Posts
		if err := stmt.Scan(&post.Id, &post.UserID, &post.Image, &post.Text, &post.Likes, &post.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("%s: %w", op, err)
		}
		data = append(data, post)
	}

	if err := stmt.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return data, totalCount, nil
}

func (s *Storage) DelPost(id int) error {
	const op = "storage.sqlite.Register"

	stmt, err := s.db.Prepare("DELETE FROM posts WHERE id=$1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetUser(userID int) (User, error) {
	const op = "storage.sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT id, name, login, balans, status, role, last_seen, gender, language, active_status_online, posts_privacy, allow_dm, allow_comments FROM users WHERE id = $1;")
	if err != nil {
		return User{}, nil
	}
	defer stmt.Close()

	var data User
	slog.Info("Передает userID в БД: ", userID)
	err = stmt.QueryRow(userID).Scan(&data.Id, &data.Name, &data.Login, &data.Balans, &data.Status, &data.Role, &data.Last_seen, &data.Gender, &data.Language, &data.Active_status_online, &data.Posts_privacy, &data.Allow_dm, &data.Allow_comments)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error("Ошибка запроса в БД:", err)
		return User{}, nil
	}
	if err != nil {
		slog.Error("Ошибка получения данных из БД:", err)
		return User{}, nil
	}
	slog.Info("Получаем id из БД: ", data.Id)
	return data, nil
}
