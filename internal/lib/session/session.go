package session

import (
	"fmt"
	"jurry_dev/internal/lib/utils"
)

type sessionData struct {
	Login  string
	UserID int
}
type Session struct {
	data map[string]*sessionData
}

func NewSession() *Session {
	s := new(Session)
	s.data = make(map[string]*sessionData)
	return s
}

var GlobalSession *Session

func InitGlobalSession() {
	GlobalSession = NewSession()
}
func (s *Session) SetLogin(login string, userID int) string {
	sessionId := utils.GenerateId()
	data := &sessionData{Login: login, UserID: userID}
	s.data[sessionId] = data
	return sessionId

}

func (s *Session) Get(sessionId string) (string, int) {
	data := s.data[sessionId]
	if data == nil {
		fmt.Println("data is nil: ", data)
		fmt.Println("sessionID: ", sessionId)
		fmt.Println("s.data[sessionID]: ", s.data[sessionId])
		fmt.Println("data.UserID: ", data.UserID)
		fmt.Println("data.Login: ", data.Login)
		return "", 0
	}
	return data.Login, data.UserID
}
