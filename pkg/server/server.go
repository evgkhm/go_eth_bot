package server

import (
	"errors"
	"net/http"
)

// New создание сервера заглушки
func New(port string) error {
	var err error
	http.HandleFunc("/", MainHandler)
	go func() {
		err = http.ListenAndServe(":"+port, nil)
		if err != nil {
			errors.New("can't create server")
		}
	}()
	return err
}

// MainHandler функция приветствия для правильной работы с heroku
func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	_, err := resp.Write([]byte("Hi there! I'm Bot!"))
	if err != nil {
		errors.New("can't write response")
	}
}
