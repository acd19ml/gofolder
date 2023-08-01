package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//'h': is a pointer to the handler
	// 'h.l': is a pointer to the logger
	h.l.Println("Hello World")
	//'r.Body' is a stream of data
	//'ioutil': is a package that contains a lot of helper functions for io
	//'ReadAll': returns a byte slice and an error
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Oops"))
		return
	}

	//curl -v -d 'Oli' localhost:9090
	//'-v': 用于显示请求的细节 verbose
	//'-d': 用于发送数据 data
	//'%s': 用于格式化字符串
	//'\n': 用于换行
	fmt.Fprintf(rw, "Hello %s\n", d)
}
