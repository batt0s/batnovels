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
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	Addr      string
	AppMode   string
	Secret    string
	AuthToken *jwtauth.JWTAuth
	Router    *chi.Mux
	Server    http.Server
	Database  *database.Database
}

func (app *App) Init() error {
	database, err := database.New("sqlite", "dev.db", &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	app.Database = database

	var host, port string
	host = strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = "127.0.0.1"
	}
	port = strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8090"
	}
	addr := host + ":" + port

	var secret string
	secret = strings.TrimSpace(os.Getenv("SECRET"))
	if secret == "" {
		secret = "Qgh]9.sNsTY<GA]H3>2k@6[Na7oIB,$IjNXN:g^51a]bhrl9u<jATQW2I2HLV0"
	}

	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(120 * time.Second))

	r.Route("/api", func(api chi.Router) {
		api.Route("/user", func(user chi.Router) {
			user.Post("/login", app.LoginHandler)
			user.Post("/register", app.RegisterHandler)
		})
		api.Route("/project", func(project chi.Router) {
			project.Get("/", app.ProjectList)
			project.Get("/{slug}", app.ProjectDetail)
			project.Get("/featured", app.FeaturedProjectList)
			project.Get("/latest", app.LatestProjectList)
			project.Get("/{slug}/chapters", app.ChapterList)

			project.Group(func(projectAuth chi.Router) {
				projectAuth.Use(jwtauth.Verifier(tokenAuth))
				projectAuth.Use(jwtauth.Authenticator(tokenAuth))

				projectAuth.Post("/", app.ProjectAdd)
				projectAuth.Post("/{slug}/chapters", app.ChapterAdd)
			})
		})
		api.Route("/chapter", func(chapter chi.Router) {
			chapter.Get("/{slug}", app.Chapter)
		})
	})

	app.Router = r
	app.Addr = addr
	app.Server = http.Server{
		Addr:    addr,
		Handler: r,
	}
	app.Secret = secret
	app.AuthToken = tokenAuth

	log.Printf("App Inited\n Addr: %s\n App Mode: %s", app.Addr, app.AppMode)

	return nil
}

func (app *App) Run() {
	log.Printf("[info] App starting on %s", app.Addr)
	app.Server.ListenAndServe()
}
