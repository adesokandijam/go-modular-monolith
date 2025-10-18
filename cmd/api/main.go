package main

import (
	"database/sql"
	"dijam-ecommerce/internal/product"
	productHTTPHandler "dijam-ecommerce/internal/product/http"
	productRepository "dijam-ecommerce/internal/product/repository"
	users "dijam-ecommerce/internal/user"
	userHTTPHandler "dijam-ecommerce/internal/user/http"
	"dijam-ecommerce/internal/user/repository"
	"dijam-ecommerce/pkg/middleware"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

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

type config struct {
	dbconfig       dbConfig
	accessTokenTTL time.Duration
	jwtSecret      []byte
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
	var cfg config
	cfg.dbconfig = dbCfg
	cfg.jwtSecret = []byte(os.Getenv("JWT_SECRET"))
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
	intAccessTokenTTL, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		log.Fatalf("Error parsing ACCESS_TOKEN_TTL: %v", err)
	}

	cfg.accessTokenTTL = time.Duration(intAccessTokenTTL)

	fmt.Printf("Loaded config: %+v\n", dbCfg)

	db, err := openDB(dbCfg)
	if err != nil {
		log.Fatalf("error connecting to the db: %e", err)
	}
	validator := validator.New()

	userPostgresRepo := repository.NewUserPostgresRepository(db)
	userService := users.NewUserService(userPostgresRepo, cfg.jwtSecret, cfg.accessTokenTTL)
	handler := userHTTPHandler.NewHandler(userService, validator)

	productPostgresRepo := productRepository.NewProductRepository(db)
	productService := product.NewProductService(productPostgresRepo)
	productHandler := productHTTPHandler.NewProductHTTPHandler(productService, validator)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/vi/users/register", handler.Register)
	mux.HandleFunc("POST /api/vi/users/login", handler.Login)
	mux.HandleFunc("POST /api/vi/products/list", productHandler.ListProduct)
	mux.HandleFunc("GET /api/vi/health", healthCheck)

	slog.Info("Server starting on port 8080")
	err = http.ListenAndServe(":8080", middleware.LogRequestMetrics(mux))
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
