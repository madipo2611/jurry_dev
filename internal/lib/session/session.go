package session

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jurry_dev/internal/config"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/lib/utils"
	"log/slog"
	"time"
)

var ctx = context.Background()
var redAd = config.MustLoad()
var Redis = NewSession(redAd.RedisAddr)

type sessionData struct {
	Login      string
	UserID     int
	LastAccess time.Time
}

// Session управляет сессиями
type Session struct {
	redisClient *redis.Client
}

// NewSession создает новый экземпляр Session
func NewSession(redisAddr string) *Session {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,    // Адрес сервера Redis
		Password: "Sanata2426", // Пароль (если установлен)
		DB:       1,            // Используемая база данных
	})

	// Проверяем подключение
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Ошибка подключения к Redis:", err)
		return nil
	}

	return &Session{redisClient: client}
}

// SaveSession сохраняет данные сессии в Redis
func (s *Session) SaveSession(ctx context.Context, login string, userID int) string {
	const op = "handler.session.SaveSession"

	sessionID := utils.GenerateId()
	data := sessionData{Login: login, UserID: userID} // Убрал указатель (&)

	dataBytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("Ошибка сериализации sessionData:", sl.Err(err))
		return ""
	}

	// Явно проверяем ошибку Set
	if err = s.redisClient.Set(ctx, sessionID, dataBytes, 3000*time.Minute).Err(); err != nil {
		slog.Error("Ошибка записи в Redis:", sl.Err(err))
		return ""
	}

	// Логируем успешное сохранение
	slog.Info("Данные сохранены в Redis",
		slog.String("sessionID", sessionID),
		slog.String("data", string(dataBytes)),
	)

	val, err := s.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		slog.Error("Ключ не найден после сохранения!", sl.Err(err))
	} else {
		slog.Info("Ключ успешно прочитан после сохранения", slog.String("value", val))
	}

	exists, err := s.redisClient.Exists(ctx, sessionID).Result()
	if err != nil {
		slog.Error("Ошибка проверки ключа в Redis:", sl.Err(err))
	} else if exists == 0 {
		slog.Error("Ключ исчез сразу после сохранения!")
	}

	return sessionID
}

// LoadSession загружает данные сессии из Redis
func (s *Session) LoadSession(ctx context.Context, sessionID string) (*sessionData, error) {
	dataBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, err // Сессия не найдена
		}
		slog.Error("Ошибка доступа к redis:", sl.Err(err))
		return nil, err // Ошибка при доступе к Redis
	}
	slog.Debug("Получаем dataBytes:", dataBytes)

	var data sessionData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		slog.Error("Ошибка json.Unmarshal при получении данных:", sl.Err(err))
		return nil, err
	}
	slog.Debug("Получаем data из редиски:", data, &data, data.UserID)

	return &data, nil
}

// DeleteSession удаляет данные сессии из Redis
func (s *Session) DeleteSession(ctx context.Context, sessionID string) error {
	return s.redisClient.Del(ctx, sessionID).Err()
}
