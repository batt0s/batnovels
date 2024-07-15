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

type ChapterRequestBody struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ProjectID string    `json:"project_id"`
	Slug      string    `json:"slug"`
}

func (app *App) ChapterList(w http.ResponseWriter, r *http.Request) {
	project_slug := chi.URLParam(r, "slug")
	if project_slug == "" {
		sendResponse(w, http.StatusBadRequest, nil)
		return
	}
	var chapters []database.Chapter
	var err error
	chapters, err = app.Database.Chapters.ListBySlug(context.Background(), project_slug, 100, 0, "created_at asc")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sendResponse(w, http.StatusNotFound, nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		log.Println(err)
		return
	}
	var requestBodies []ChapterRequestBody
	for _, chapter := range chapters {
		requestBody := ChapterRequestBody{
			ID:        chapter.ID,
			CreatedAt: chapter.CreatedAt,
			UpdatedAt: chapter.UpdatedAt,
			Title:     chapter.Title,
			Content:   chapter.Content[:61] + "...",
			ProjectID: chapter.ProjectID,
			Slug:      chapter.Slug,
		}
		requestBodies = append(requestBodies, requestBody)
	}
	sendResponse(w, http.StatusOK, requestBodies)
}

func (app *App) Chapter(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		sendResponse(w, http.StatusBadRequest, nil)
		return
	}
	var chapter database.Chapter
	var err error
	chapter, err = app.Database.Chapters.FindBySlug(context.Background(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sendResponse(w, http.StatusNotFound, nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, chapter)
}

func (app *App) ChapterAdd(w http.ResponseWriter, r *http.Request) {
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
	project_slug := chi.URLParam(r, "slug")
	if project_slug == "" {
		sendResponse(w, http.StatusBadRequest, nil)
		return
	}
	body, err := getRequestBody[ChapterRequestBody](w, r)
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
	project, err := app.Database.Projects.FindBySlug(context.Background(), project_slug)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	chapter := database.Chapter{
		Title:     body.Title,
		Content:   body.Content,
		ProjectID: project.ID,
	}
	chapter, err = app.Database.Chapters.Add(context.Background(), chapter)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, chapter)
}
