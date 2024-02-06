package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arthurqueiroz4/rinha-go/config"
)

func main() {
	log.Println("Starting the application...")
	env := &config.Env{}
	env.LoadEnv()
	env.ConnectDB()
	defer env.CloseDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	log.Printf("Server started on: http://localhost:%s", env.AppPort)
	http.ListenAndServe(fmt.Sprintf(":%s", env.AppPort), nil)
}
