package session

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"reflect"
	"testing"
)

func TestSession_InitSession(t *testing.T) {
	c := &Session{
		CookieLifeTime: "100",
		CookiePersist:  "true",
		CookieName:     "JSM",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}

	var sm *scs.SessionManager

	ses := c.InitSession()

	var sesKind reflect.Kind
	var sesType reflect.Type

	rv := reflect.ValueOf(ses)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println("For Loop: ", rv.Kind(), rv.Type(), rv)
		sesKind = rv.Kind()
		sesType = rv.Type()

		rv = rv.Elem()
	}

	if !rv.IsValid() {
		t.Error("Invalid type or kind; kind:", rv.Kind(), "type:", rv.Type())
	}

	if sesKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind returned testing cookie session. Expected",
			reflect.ValueOf(sm).Kind(), "and got", sesKind)
	}

	if sesType != reflect.ValueOf(sm).Type() {
		t.Error("wrong kind returned testing cookie session. Expected",
			reflect.ValueOf(sm).Type(), "and got", sesType)
	}
}
