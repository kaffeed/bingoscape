package main

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/kaffeed/bingoscape/db"
	"github.com/kaffeed/bingoscape/handlers"
	"github.com/kaffeed/bingoscape/services"
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

	us := services.NewUserServices(store)
	// userservice.CreateUser(services.User{})
	authhandler := handlers.NewAuthHandler(us)
	bs := services.NewBingoService(store)
	bingohandler := handlers.NewBingoHandler(bs, us)

	handlers.SetupRoutes(e, authhandler, bingohandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
