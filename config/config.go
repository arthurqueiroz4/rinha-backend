package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Env struct {
	DbHost  string `mapstructure:"DB_HOST"`
	AppPort string `mapstructure:"APP_PORT"`
	DB *gorm.DB
}

func (e *Env) LoadEnv() {
    log.Println("Loading environment variables...")
    viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	env := &Env{}
	if err := viper.Unmarshal(env); err != nil {
		panic(err)
	}
    log.Println("Environment variables loaded successfully...")
	log.Println(env)
	e.DbHost = env.DbHost
	e.AppPort = env.AppPort
}

func (e *Env) ConnectDB() {
	log.Println("Connecting to the database...")
	dsn := fmt.Sprintf(
		"host=%s user=postgres password=postgres dbname=rinha port=5432 sslmode=disable",
		e.DbHost,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database successfully...")
	e.DB = db
}

func (e *Env) CloseDB() {
	log.Println("Closing the database connection...")
	db, err := e.DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	log.Println("Database connection closed successfully...")
}