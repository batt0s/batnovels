package database

import (
	"context"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             string         `gorm:"type:uuid;primary_key;"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	LastLogin      time.Time      `json:"last_login"`
	IsAdmin        bool           `json:"is_admin"`
	IsStaff        bool           `json:"is_staff"`
	Username       string         `gorm:"not null;size:256;unique;;" json:"username"`
	Email          string         `gorm:"not null;size:256;unique;;" json:"email"`
	Name           string         `gorm:"not null;size:128;;" json:"name"`
	Password       string         `gorm:"not null;size:128;;" json:"-"`
	ProfilePicture string         `gorm:"size:128;" json:"profile_picture"`
}

type UserRepo interface {
	Find(ctx context.Context, id string) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Add(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, user User) error
}

type SqlUserRepo struct {
	db *gorm.DB
}

func NewSqlUserRepo(db *gorm.DB) *SqlUserRepo {
	return &SqlUserRepo{
		db: db,
	}
}

func (repo SqlUserRepo) Find(ctx context.Context, id string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ErrorOperationCanceled
	default:
		var user User
		result := repo.db.First(&user, "id = ?", id)
		return user, result.Error
	}
}

func (repo SqlUserRepo) FindByUsername(ctx context.Context, username string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ErrorOperationCanceled
	default:
		var user User
		result := repo.db.First(&user, "username = ?", username)
		return user, result.Error
	}
}

func (repo SqlUserRepo) FindByEmail(ctx context.Context, email string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ErrorOperationCanceled
	default:
		var user User
		result := repo.db.First(&user, "email = ?", email)
		return user, result.Error
	}
}

func (repo SqlUserRepo) Add(ctx context.Context, user User) error {
	select {
	case <-ctx.Done():
		return ErrorOperationCanceled
	default:
		if !user.IsValid() {
			return ErrorInvalidUser
		}
		user.ID = uuid.New().String()
		user.SetPassword(user.Password)
		result := repo.db.Create(&user)
		return result.Error
	}
}

func (repo SqlUserRepo) Update(ctx context.Context, user User) error {
	select {
	case <-ctx.Done():
		return ErrorOperationCanceled
	default:
		result := repo.db.Save(&user)
		return result.Error
	}
}

func (repo SqlUserRepo) Delete(ctx context.Context, user User) error {
	select {
	case <-ctx.Done():
		return ErrorOperationCanceled
	default:
		result := repo.db.Delete(&user)
		return result.Error
	}
}

func (u User) IsValid() bool {
	if len(u.Username) > 256 || len(u.Username) < 4 {
		return false
	}
	if len(u.Email) > 256 {
		return false
	}
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return false
	}
	if len(u.Name) > 128 || len(u.Name) < 3 {
		return false
	}
	return true
}

// does not save with new password
func (u *User) SetPassword(passwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}
