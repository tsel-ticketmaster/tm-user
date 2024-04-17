package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	c              *Config
	configSyncOnce sync.Once
)

type Config struct {
	Application struct {
		Name        string
		Port        int
		Environment string
		Debug       bool
		Timeout     time.Duration
	}
	OpenTelemetry struct {
		Collector struct {
			Endpoint string
		}
	}
	CORS struct {
		AllowedOrigins   []string
		AllowedMethods   []string
		AllowedHeaders   []string
		ExposedHeaders   []string
		MaxAge           int
		AllowCredentials bool
	}
	JWT struct {
		PrivateKey []byte
		PublicKey  []byte
	}
	Postgresql struct {
		Host         string
		Port         int
		User         string
		Password     string
		DBName       string
		SSLMode      string
		MaxOpenConns int
		MaxIdleConns int
	}
	Redis struct {
		Addrs    []string
		Username string
		Password string
		DB       int
	}
}

func (cfg *Config) application() {
	cfg.Application.Name = os.Getenv("APP_NAME")
	cfg.Application.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	cfg.Application.Environment = os.Getenv("APP_ENVIRONMENT")
	cfg.Application.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))

	timeoutInSec, _ := strconv.Atoi(os.Getenv("APP_TIMEOUT"))
	cfg.Application.Timeout = time.Duration(timeoutInSec) * time.Second
}

func (cfg *Config) openTelemetry() {
	collectorEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
	cfg.OpenTelemetry.Collector.Endpoint = collectorEndpoint
}

func (cfg *Config) cors() {
	cfg.CORS.AllowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	cfg.CORS.AllowedMethods = strings.Split(os.Getenv("CORS_ALLOWED_METHODS"), ",")
	cfg.CORS.AllowedHeaders = strings.Split(os.Getenv("CORS_ALLOWED_HEADERS"), ",")
	cfg.CORS.ExposedHeaders = strings.Split(os.Getenv("CORS_EXPOSED_HEADERS"), ",")
	cfg.CORS.AllowCredentials, _ = strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))
	cfg.CORS.MaxAge, _ = strconv.Atoi(os.Getenv("CORS_MAX_AGE"))
}

func (c *Config) jwt() {
	jwtRsaPlain := os.Getenv("JWT_RSA")
	var jwtRsa = struct {
		PrivateKey string `json:"private"`
		PublicKey  string `json:"public"`
	}{}

	json.Unmarshal([]byte(jwtRsaPlain), &jwtRsa)

	c.JWT.PrivateKey = []byte(jwtRsa.PrivateKey)
	c.JWT.PublicKey = []byte(jwtRsa.PublicKey)

	fmt.Println(string(jwtRsa.PrivateKey))
}

func (c *Config) postgresql() {
	c.Postgresql.Host = os.Getenv("POSTGRESQL_HOST")
	c.Postgresql.Port, _ = strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))
	c.Postgresql.User = os.Getenv("POSTGRESQL_USER")
	c.Postgresql.Password = os.Getenv("POSTGRESQL_PASSWORD")
	c.Postgresql.DBName = os.Getenv("POSTGRESQL_DBNAME")
	c.Postgresql.SSLMode = os.Getenv("POSTGRESQL_SSLMODE")
	c.Postgresql.MaxOpenConns, _ = strconv.Atoi(os.Getenv("POSTGRESQL_MAX_OPEN_CONNS"))
	c.Postgresql.MaxIdleConns, _ = strconv.Atoi(os.Getenv("POSTGRESQL_MAX_IDLE_CONNS"))
}

func (cfg *Config) redis() {
	cfg.Redis.Addrs = strings.Split(os.Getenv("REDIS_HOSTS"), ",")
	cfg.Redis.Username = os.Getenv("REDIS_USERNAME")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.DB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
}

func load() *Config {
	cfg := new(Config)
	cfg.application()
	cfg.openTelemetry()
	cfg.jwt()
	cfg.postgresql()
	cfg.cors()
	cfg.redis()
	return cfg
}

func Get() *Config {
	configSyncOnce.Do(func() {
		c = load()
	})

	return c
}
