package session

import (
	"jurry_dev/internal/lib/utils"
)

type sessionData struct {
	Login string
}
type Session struct {
	data map[string]*sessionData
}

func NewSession() *Session {
	s := new(Session)
	s.data = make(map[string]*sessionData)
	return s
}

func (s *Session) SetLogin(login string) string {
	sessionId := utils.GenerateId()
	data := &sessionData{Login: login}
	s.data[sessionId] = data
	return sessionId

}

func (s *Session) Get(sessionId string) string {
	data := s.data[sessionId]
	if data == nil {
		return ""
	}
	return data.Login
}
