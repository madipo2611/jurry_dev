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
		DB:       0,            // Используемая база данных
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
	sessionID := utils.GenerateId()
	data := &sessionData{Login: login, UserID: userID}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("Ошибка формирования sessionData:", sl.Err(err))
		return ""
	}

	// Сохранение данных в Redis с временем жизни (например, 30 минут)
	s.redisClient.Set(ctx, sessionID, dataBytes, 30*time.Minute)
	if err != nil {
		slog.Error("Ошибка сохранения сессии в Redis:", sl.Err(err))
		return ""
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

	var data sessionData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		slog.Error("Ошибка json.Unmarshal при получении данных:", sl.Err(err))
		return nil, err
	}

	return &data, nil
}

// DeleteSession удаляет данные сессии из Redis
func (s *Session) DeleteSession(ctx context.Context, sessionID string) error {
	return s.redisClient.Del(ctx, sessionID).Err()
}
