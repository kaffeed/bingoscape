package main

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/kaffeed/topez-bingomania/db"
	"github.com/kaffeed/topez-bingomania/services"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	e := echo.New()

	e.Static("/", "assets")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))))

	store, err := db.NewStore(os.Getenv("DB"))
	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	userservice := services.NewUserServices(services.User{}, store)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hallo, welt!")
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
