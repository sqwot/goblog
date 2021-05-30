package session

import (
	"GoBlog/utils"
	"net/http"
	"time"

	"github.com/codegangsta/martini"
)

const (
	COOKIE_NAME = "sessionId"
)

type Session struct {
	Id       string
	Username string
}

type SessionStore struct {
	data map[string]*Session
}

func NewSessionStore() *SessionStore {
	s := new(SessionStore)
	s.data = make(map[string]*Session)
	return s
}

func ensureCookie(w http.ResponseWriter, r *http.Request) string {
	cookie, _ := r.Cookie(COOKIE_NAME)
	if cookie != nil {
		return cookie.Value
	}
	sessonId := utils.GenerateId()
	cookie = &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   sessonId,
		Expires: time.Now().Add(5 * time.Minute),
	}
	http.SetCookie(w, cookie)

	return sessonId
}
func (store *SessionStore) Get(sessionId string) *Session {
	session := store.data[sessionId]
	if session == nil {
		return &Session{
			Id: sessionId,
			//Name: COOKIE_NAME,
		}
	}
	return session
}

func (store *SessionStore) Set(session *Session) {
	store.data[session.Id] = session
}

var sessionStore = NewSessionStore()

func Middleware(w http.ResponseWriter, r *http.Request, ctx martini.Context) {
	sessionId := ensureCookie(w, r)
	session := sessionStore.Get(sessionId)

	ctx.Map(session)

	ctx.Next()

	sessionStore.Set(session)
}
