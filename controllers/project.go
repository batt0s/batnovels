package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/batt0s/batnovels/database"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ProjectRequestBody struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Synopsis  string    `json:"synopsis"`
	Author    string    `json:"author"`
	Status    string    `json:"status"`
	Tags      string    `json:"tags"` // , ile ayırarak
	Views     int32     `json:"views"`
	Image     string    `json:"image"`
}

func (app *App) ProjectList(w http.ResponseWriter, r *http.Request) {
	var projects []database.Project
	var err error
	projects, err = app.Database.Projects.List(context.Background(), 100, 0, "created_at")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sendResponse(w, http.StatusNotFound, nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, projects)
}

func (app *App) FeaturedProjectList(w http.ResponseWriter, r *http.Request) {
	var projects []database.Project
	var err error
	projects, err = app.Database.Projects.List(context.Background(), 100, 0, "views desc, created_at desc")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sendResponse(w, http.StatusNotFound, nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, projects)
}

func (app *App) LatestProjectList(w http.ResponseWriter, r *http.Request) {
	var projects []database.Project
	results := app.Database.DB.Table("projects").
		Select("projects.*, MAX(chapters.created_at) as last_chapter_created_at").
		Joins("JOIN chapters ON chapters.project_id = projects.id").
		Where("chapters.deleted_at IS NULL").
		Group("projects.id").
		Order("last_chapter_created_at DESC").
		Find(&projects)
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			sendResponse(w, http.StatusNotFound, nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, map[string]string{"error": results.Error.Error()})
		}
		log.Println(results.Error)
		return
	}
	sendResponse(w, http.StatusOK, projects)
}

func (app *App) ProjectDetail(w http.ResponseWriter, r *http.Request) {
	project_slug := chi.URLParam(r, "slug")
	if project_slug == "" {
		sendResponse(w, http.StatusBadRequest, nil)
		return
	}
	var project database.Project
	var err error
	project, err = app.Database.Projects.FindBySlug(context.Background(), project_slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sendResponse(w, http.StatusNotFound, nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, project)
}

func (app *App) ProjectAdd(w http.ResponseWriter, r *http.Request) {
	user, err := userContextBody(app.Database.Users, r.Context())
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	if !user.IsStaff {
		sendResponse(w, http.StatusUnauthorized, nil)
		return
	}
	body, err := getRequestBody[ProjectRequestBody](w, r)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println(err)
		return
	}
	project := database.Project{
		Title:    body.Title,
		Synopsis: body.Synopsis,
		Author:   body.Author,
		Status:   body.Status,
		Tags:     body.Tags,
		Image:    body.Image,
	}
	project, err = app.Database.Projects.Add(context.Background(), project)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, project)
}
