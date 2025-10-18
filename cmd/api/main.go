package main

import (
	"database/sql"
	users "dijam-ecommerce/internal/user"
	userHTTPHandler "dijam-ecommerce/internal/user/http"
	"dijam-ecommerce/internal/user/repository"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	host         string
	username     string
	password     string
	db           string
	sslmode      string
	port         int
	maxIdleConns int
	maxOpenConns int
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	var dbCfg dbConfig

	// Read values from the environment
	dbCfg.host = os.Getenv("DB_HOST")
	dbCfg.db = os.Getenv("DB_NAME")
	dbCfg.username = os.Getenv("DB_USER")
	dbCfg.password = os.Getenv("DB_PASSWORD")
	dbCfg.sslmode = os.Getenv("DB_SSLMODE")

	// Convert string values from .env to integers
	dbCfg.port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error parsing DB_PORT: %v", err)
	}

	dbCfg.maxIdleConns, err = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		log.Fatalf("Error parsing DB_MAX_IDLE_CONNS: %v", err)
	}

	dbCfg.maxOpenConns, err = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		log.Fatalf("Error parsing DB_MAX_OPEN_CONNS: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", dbCfg)

	db, err := openDB(dbCfg)
	if err != nil {
		log.Fatalf("error connecting to the db: %e", err)
	}
	postgresRepo := repository.NewPostgresRepository(db)
	userService := users.NewUserService(postgresRepo)
	validator := validator.New()
	handler := userHTTPHandler.NewHandler(userService, validator)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/vi/users/register", handler.Register)
	mux.HandleFunc("POST /api/vi/users/login", handler.Login)
	mux.HandleFunc("GET /api/vi/health", healthCheck)

	slog.Info("Server starting on port 8080")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func openDB(cfg dbConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.username, cfg.password, cfg.host, cfg.db, cfg.sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to open DB connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetMaxOpenConns(cfg.maxOpenConns)

	return db, nil
}
