package controllers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/batt0s/batnovels/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	Addr     string
	AppMode  string
	Router   *chi.Mux
	Server   http.Server
	Database *database.Database
}

func (app *App) Init() error {
	database, err := database.New("sqlite", "dev.db", &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	app.Database = database

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(120 * time.Second))

	r.Route("/api", func(api chi.Router) {
		api.Get("/project", app.ProjectList)
		api.Get("/project/{id}", app.ProjectDetail)
		api.Get("/project/{id}/chapters", app.ChapterList)
		api.Post("/project", app.ProjectAdd)
		api.Post("/project/{id}/chapters", app.ChapterAdd)
	})

	var host, port string
	host = strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = "127.0.0.1"
	}
	port = strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port

	app.Router = r
	app.Addr = addr
	app.Server = http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("App Inited\n Addr: %s\n App Mode: %s", app.Addr, app.AppMode)

	return nil
}

func (app *App) Run() {
	log.Printf("[info] App starting on %s", app.Addr)
	app.Server.ListenAndServe()
}
