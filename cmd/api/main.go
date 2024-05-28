package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/kaffeed/bingoscape/db"
	"github.com/kaffeed/bingoscape/handlers"
	"github.com/kaffeed/bingoscape/services"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setupImageDirectories(path string) {
	subPaths := []string{"tiles", "bingo", "submissions"}

	for _, p := range subPaths {
		fp := filepath.Join(path, p)
		if _, err := os.Stat(fp); !os.IsNotExist(err) { // image path exists
			log.Printf("Path %s exists, skipping creation\n", fp)
			continue
		}
		err := os.MkdirAll(filepath.Join(path, p), os.ModePerm)
		if err != nil {
			panic("Could not create directories!")
		}
	}
}

func main() {
	godotenv.Load()
	e := echo.New()

	e.Static("/", "assets")
	p := filepath.Join(os.Getenv("IMAGE_PATH"))
	setupImageDirectories(p)

	imageGroup := e.Group("/img")
	imageGroup.Use(middleware.Static(p))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))))

	store, err := db.NewStore(os.Getenv("DB"))

	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	us := services.NewUserServices(store)
	bs := services.NewBingoService(store)
	ah := handlers.NewAuthHandler(us)
	bh := handlers.NewBingoHandler(bs, us)

	handlers.SetupRoutes(e, ah, bh)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
