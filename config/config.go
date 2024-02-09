package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DbHost  string `mapstructure:"DB_HOST"`
	AppPort string `mapstructure:"APP_PORT"`
	DB      *pgxpool.Pool
}

func (e *Env) LoadEnv() {
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
	e.DbHost = env.DbHost
	e.AppPort = env.AppPort
}

func (e *Env) ConnectDB() {
	conn, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:5432/rinha", e.DbHost))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database successfully...")
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	e.DB = conn
}

func (e *Env) CloseDB() {
	log.Println("Closing the database connection...")
	e.DB.Close()
	log.Println("Database connection closed successfully...")
}
