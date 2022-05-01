package server

import (
	"fmt"
	"log"
	"main/cache"
	"net/http"
	"strings"
)

const defaultBasePath = "/_bincache/"

type HTTPPool struct {
	self     string
	basePath string
}

func New(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	urls := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(urls) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	groupName := urls[0]
	key := urls[1]

	group := cache.GetGroup(groupName)
	if group == nil {
		http.Error(w, http.StatusText(http.StatusNotFound) + " " + groupName, http.StatusNotFound)
		return
	}

	byteView, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(byteView.ByteSlice())
}