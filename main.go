package main

import (
	"context"
	"crypto/rand"
	"log"
	"money-api/handlers"
	"money-api/platform/migration"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Setup sqlite database
	db, err := sqlx.Open("sqlite3", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Migrate first
	err = migration.Migrate(db, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Setup in memory database
	mem, err := bigcache.NewBigCache(bigcache.DefaultConfig(3 * time.Hour))
	if err != nil {
		log.Fatal(err)
	}
	defer mem.Close()

	// Generate RSA Public/Private key pair for JWT
	jwtSecret := make([]byte, 64)
	_, err = rand.Read(jwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	app := echo.New()

	h := handlers.Dependency{
		DB:        db,
		Memory:    mem,
		JWTSecret: jwtSecret,
	}

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("ALLOW_URL"), ","),
		AllowHeaders: []string{echo.HeaderAccept, echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderCookie},
		MaxAge:       60 * 60 * 3,
	}))

	auth := app.Group("/auth")
	auth.GET("/", h.Login)
	auth.POST("/", h.Signup)

	account := app.Group("/account")
	account.Use(h.MustLogin)
	account.GET("/my", h.GetMyAccount)
	account.PATCH("/my", h.UpdateAccount)
	account.GET("/friends", h.GetFriends)
	account.PUT("/friends", h.AddFriend)
	account.DELETE("/friends/:tag", h.RemoveFriend)

	transaction := app.Group("/transaction")
	transaction.Use(h.MustLogin)
	transaction.GET("/my", h.GetAllTransaction)
	transaction.GET("/friend/:tag", h.GetTransactionByFriend)
	transaction.GET("/id/:id", h.GetTransactionByID)
	transaction.POST("/send", h.SendMoney)
	transaction.POST("/request", h.RequestMoney)
	transaction.PATCH("/cancel/:id", h.UpdateStatus)
	transaction.PATCH("/reject/:id", h.UpdateStatus)

	// Start server
	go func() {
		if err := app.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}
}
