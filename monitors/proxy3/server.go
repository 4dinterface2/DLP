package proxy3

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
)

func CheckRequest(r *http.Request) bool {
	//host := r.Host;
	// Отправка ответа с содержимым тела запроса
	// w.Header().Set("Content-Type", "text/plain") ;

	// fmt.Println(str);

	m := r.Header.Get("method")
	fmt.Println("hello " + r.RequestURI + " " + r.Method + "/ " + m)
	if r.Method != "CONNECT" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return true
		}

		fmt.Println(r.Method)
		fmt.Printf("Body: %s", body)

	}

	return true
}

type MyHandler struct {
	Message string
	proxy   *goproxy.ProxyHttpServer
}

// Добавляем метод к структуре
func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if CheckRequest(r) {
		h.proxy.ServeHTTP(w, r)
	}
}

func NewProxy() *MyHandler {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	return &MyHandler{
		proxy:   proxy,
		Message: "Привет, мир!",
	}
}
