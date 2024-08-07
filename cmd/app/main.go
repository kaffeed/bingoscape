package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/kaffeed/bingoscape/app/db"
	"github.com/kaffeed/bingoscape/app/handlers"
	"github.com/kaffeed/bingoscape/app/services"
	"github.com/kaffeed/bingoscape/public"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pressly/goose/v3"
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

	ctx := context.Background()
	connpool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error connecting to db: %#v", err)
	}
	defer connpool.Close()
	p := filepath.Join(os.Getenv("IMAGE_PATH"))
	setupImageDirectories(p)

	imageGroup := e.Group("/img")
	imageGroup.Use(middleware.Static(p))

	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "assets", // because files are located in `assets` directory in `webAssets` fs
		Filesystem: http.FS(public.AssetsFS),
	}))

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))))

	goose.SetBaseFS(db.EmbedMigrations)

	if err := goose.SetDialect(os.Getenv("DB_DRIVER")); err != nil {
		log.Fatalf(err.Error())
	}

	if err := goose.Up(stdlib.OpenDBFromPool(connpool), "migrations"); err != nil {
		log.Fatalf(err.Error())
	}

	store := db.New(connpool)

	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	us := services.NewUserServices(store)
	bs := services.NewBingoService(store, connpool)
	ts := services.NewTileService(store, connpool)

	authHandler := handlers.NewAuthHandler(us)
	bingoHandler := handlers.NewBingoHandler(bs, us, ts)
	tileHandler := handlers.NewTileHandler(ts, us)
	apiHandler := handlers.NewApiHandler(ts, us, bs)

	handlers.SetupRoutes(e, authHandler, bingoHandler, tileHandler, apiHandler)

	e.Logger.Fatal(e.Start(os.Getenv("HTTP_LISTEN_ADDR")))
}
