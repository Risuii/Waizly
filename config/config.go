package config

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Port string
	}
	Database struct {
		DSN string
	}
	Bcrypt struct {
		HashCost int
	}
	Jwt struct {
		PrivateKey *rsa.PrivateKey
		PublicKey  *rsa.PublicKey
	}
	BasicAuth struct {
		Username string
		Password string
	}
}

func New() *Config {
	c := new(Config)
	c.loadApp()
	c.loadDatabase()
	c.loadBcrypt()

	return c
}

func (c *Config) loadApp() *Config {
	port := os.Getenv("PORT")

	c.App.Port = port

	return c
}

func (c *Config) loadDatabase() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get Env Value

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE_NAME")

	connVal := url.Values{}
	// connVal.Add("parseTime", "1")
	connVal.Add("loc", "Asia/Jakarta")

	dbConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)

	dsn := fmt.Sprintf("%s?%s", dbConnection, connVal.Encode())

	c.Database.DSN = dsn

	c.Database.DSN = dbConnection

	return c
}

func (c *Config) loadBcrypt() *Config {
	// env value
	hashCost := os.Getenv("BCRYPT_HASH_COST")

	c.Bcrypt.HashCost, _ = strconv.Atoi(hashCost)

	return c
}
