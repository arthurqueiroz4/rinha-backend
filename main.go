package main

import (
	"log"
	"net/http"
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
)

func main(){
	db := connectDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello World"))
	})

	log.Println("Server started on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func connectDB() *gorm.DB {
	dsn := "
		host=localhost
		user=postgres
		password=postgres
		dbname=postgres
		port=5432
		sslmode=disable
	"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}