package streamserver

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MiddleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := MiddleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func (m MiddleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		SendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/viders/:id-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)
	// router.GET("/testpage", testPageHandler)
	return router
}

func main() {
	r := RegisterHandlers()
	ml := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", ml)
}
