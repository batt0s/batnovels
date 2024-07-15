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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd))
	if err != nil {
		return user, ErrorIncorrectPassword
	}
	return user, err
}
