package sessions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var CookieStore *sessions.CookieStore

func init(){
	CookieStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
}

type Session struct {
	w http.ResponseWriter
	r *http.Request
	session *sessions.Session
	first_run bool
}

type Sessions interface {
	RetrieveSession(w http.ResponseWriter, r *http.Request) (*Session, error)
}

func (s *Session) FirstRun() Sessions {
	s.first_run = true
	return s
}

func (s *Session) RetrieveSession(w http.ResponseWriter, r *http.Request) (*Session, error){
	// If it is the first run, we need to delete existing session
	// Because that session is encrypted and we already have new key
	// So, we cannot decrypt that session and we need a new one
	if s.first_run {
		cookie := http.Cookie {
			Name: "filemanager",
			Value: "",
			Expires: time.Now().Add(time.Hour * -1),
		}
		http.SetCookie(w, &cookie)
	}
	session, err := CookieStore.Get(r, "filemanager")
	if err != nil {
		if err.Error() == securecookie.ErrMacInvalid.Error() {
			// If there's an error about that the existing session is encrypted and unreadable
			// We are calling refresh because we already made a new session
			http.Redirect(w, r, "/", http.StatusFound)
		}
		return nil, err

	}
	session.Save(r, w)
	return &Session{session: session, w: w, r: r,}, nil
}

func (s *Session) Update(name string, value string) error {
	if s.session == nil {
		return fmt.Errorf("sessions: Cannot request the value of %s because session is nil", name)
	}
	s.session.Values[name] = value
	s.session.Save(s.r, s.w)
	return nil
}

func (s *Session) Get(name string) (string, error){
	if s.session == nil {
		return "", fmt.Errorf("sessions: Cannot request the value of %s because session is nil", name)
	}
	return s.session.Values[name].(string), nil
}

func (s *Session) Compare(name string, compareWith interface{}) (bool, error){
	if s.session == nil {
		return false, fmt.Errorf("sessions: Cannot request the value of %s because session is nil", name)
	}
	return s.session.Values[name] == compareWith, nil
}