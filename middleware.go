package jsm

import "net/http"

func (j *Jsm) SessionLoad(next http.Handler) http.Handler {
	return j.Session.LoadAndSave(next)
}
