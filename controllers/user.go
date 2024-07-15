package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

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

// JWT tokendeki user objesini çekmek için
func userContextBody(users database.UserRepo, ctx context.Context) (database.User, error) {
	var user database.User
	var err error
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return user, err
	}
	username := claims["user"].(string)
	if username == "" {
		return user, errors.New("no username in claims")
	}

	user, err = users.FindByUsername(context.Background(), username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
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
	claims := map[string]interface{}{
		"authorized": true,
		"user":       user.Username,
		"exp":        time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	_, tokenString, err := app.AuthToken.Encode(claims)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "token generation error"})
		log.Println(err)
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"token": tokenString})
}
