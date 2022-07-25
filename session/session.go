package session

import "C"
import (
	"github.com/alexedwards/scs/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	CookieName     string
	CookieLifeTime string
	CookiePersist  string
	CookieSecure   string
	CookieDomain   string
	SessionType    string
}

func (c *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// How long should sessions last?
	minutes, err := strconv.Atoi(c.CookieLifeTime)
	if err != nil {
		minutes = 60
	}

	// Should cookies persist?
	if strings.ToLower(c.CookiePersist) == "true" {
		persist = true
	}

	// Must cookies be secure?
	if strings.ToLower(c.CookieSecure) == "true" {
		secure = true
	}

	// Create session
	session := scs.New()
	session.Lifetime = time.Minute * time.Duration(minutes)
	session.Cookie.Name = c.CookieName
	session.Cookie.Persist = persist
	session.Cookie.Secure = secure
	session.Cookie.Domain = c.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	// Which session store
	switch strings.ToLower(c.SessionType) {
	case "redis":

	case "mysql, mariadb":

	case "postgres, postgresql":

	default:
		//cookies
	}

	return session
}
