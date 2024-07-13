package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/batt0s/batnovels/authentication"
	"github.com/batt0s/batnovels/database"
	"github.com/go-chi/jwtauth/v5"
)

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequestBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (app App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body, err := getRequestBody[RegisterRequestBody](w, r)
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
	user := database.User{
		Username: body.Username,
		Email:    body.Email,
		Name:     body.Name,
		Password: body.Password,
	}
	err = app.Database.Users.Add(context.Background(), user)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, nil)
}

func (app App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := getRequestBody[LoginRequestBody](w, r)
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
	user, err := authentication.Authenticate(body.Username, body.Password, app.Database.Users)
	if err != nil {
		log.Println(err)
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	tokenAuth := jwtauth.New("HS256", []byte(app.Secret), nil)
	claims := map[string]interface{}{"username": user.Username}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "Token generation error")
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"token": tokenString})
}
