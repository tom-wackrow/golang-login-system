package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

var users = map[string]string{
	"admin": "admin",
	"test": "test",
}

var sessions = map[string]session{}

type session struct {
	username string
	expiry time.Time
	PreviousPage string
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func GetCredentials(w http.ResponseWriter, r *http.Request) Credentials {
	// var creds Credentials
	// if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return Credentials{}, err
	// }
	// fmt.Printf("%v", creds)
	// return creds, nil

	creds := Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	return creds
}

func GetSession(w http.ResponseWriter, r *http.Request) (session, error) {
	cookie, err := r.Cookie("session")

	if err != nil {
		return session{}, http.ErrNoCookie
	}

	sessionToken := cookie.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		return session{}, http.ErrNoCookie
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		return session{}, http.ErrNoCookie
	}

	return userSession, nil
}

func Auth(w http.ResponseWriter, r *http.Request) bool{
	creds := GetCredentials(w, r)

	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		return false
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[newSessionToken] = session{
		username: creds.Username,
		expiry: expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: newSessionToken,
		Expires: expiresAt,
	})

	return true
}


func RefreshAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}

	sessionToken := cookie.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		http.Redirect(w, r, "/login", 303)
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[newSessionToken] = session{
		username: userSession.username,
		expiry: expiresAt,
	}
	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: newSessionToken,
		Expires: expiresAt,
	})
	
}

func Logout(w http.ResponseWriter, r  *http.Request) {
	cookie, _ := r.Cookie("session")

	sessionToken := cookie.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: "",
		Expires: time.Now(),
	})
}

func RequireAuth(f http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")

		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", 303)
				return
			}
			http.Redirect(w, r, "/login", 303)
			return
		}

		sessionToken := cookie.Value

		userSession, exists := sessions[sessionToken]

		if !exists {
			http.Redirect(w, r, "/login", 303)
			return
		}

		if userSession.isExpired() {
			delete(sessions, sessionToken)
			http.Redirect(w, r, "/refresh", 303)
			return
		}

		// w.Write([]byte(fmt.Sprintf("Token: %s", sessionToken)))



		f(w, r)
		return
	}
}