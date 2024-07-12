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

func (app App) ProjectList(w http.ResponseWriter, r *http.Request) {
	var projects []database.Project
	var err error
	projects, err = app.Database.Projects.List(context.Background())
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

func (app App) ProjectDetail(w http.ResponseWriter, r *http.Request) {
	project_id := chi.URLParam(r, "id")
	if project_id == "" {
		sendResponse(w, http.StatusBadRequest, nil)
		return
	}
	var project database.Project
	var err error
	project, err = app.Database.Projects.Find(context.Background(), project_id)
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

func (app App) ProjectAdd(w http.ResponseWriter, r *http.Request) {
	var err error
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
	err = app.Database.Projects.Add(context.Background(), project)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, project)
}
