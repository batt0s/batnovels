package authentication

import (
	"context"
	"time"

	"github.com/batt0s/batnovels/database"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(username, passwd string, users database.UserRepo) (database.User, error) {
	var user database.User
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	user, err = users.FindByUsername(ctx, username)
	if err != nil {
		return user, err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	if user.Password != string(bytes) {
		return user, ErrorIncorrectPassword
	}
	return user, err
}
