package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/batt0s/batnovels/database"
	"gorm.io/gorm"
)

var (
	db   *database.Database
	user = database.User{
		Username: "test",
		Email:    "test@gmail.com",
		Name:     "tester",
		Password: "test",
	}
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
)

func TestMain(m *testing.M) {
	log.Println("Starting testing...")
	exitVal := m.Run()
	log.Println("Done testing.")
	err := os.Remove("test.db")
	if err != nil {
		log.Println("Could not remove test.db")
	}
	os.Exit(exitVal)
}

func TestNew(t *testing.T) {
	var err error
	db, err = database.New("sqlite", "test.db", &gorm.Config{})
	if err != nil {
		t.Errorf("[ERROR] -> %v", err)
	}
}

func TestAddUser(t *testing.T) {
	err := db.Users.Add(ctx, user)
	if err != nil {
		t.Errorf("[ERROR] -> %v", err)
	}
}

func TestFindUserByUsername(t *testing.T) {
	usr, err := db.Users.FindByUsername(ctx, user.Username)
	if err != nil {
		t.Errorf("[ERROR] -> %v", err)
	}
	if usr.Username != user.Username {
		t.Errorf("Want %s, got %s", user.Username, usr.Username)
	}
	if usr.Email != user.Email {
		t.Errorf("Want %s, got %s", user.Email, usr.Email)
	}
	user.ID = usr.ID
}
