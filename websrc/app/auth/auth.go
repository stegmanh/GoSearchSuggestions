package auth

import (
	"errors"
	"fmt"
	"net/http"

	ctx "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore(
	[]byte(securecookie.GenerateRandomKey(64)), //Signing key
	[]byte(securecookie.GenerateRandomKey(32)))

func Login(r *http.Request) (bool, error) {
	err := errors.New("Error")
	password := r.FormValue("password")
	session, _ := Store.Get(r, "admin")
	fmt.Println(bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost))
	if password != "wow" {
		return false, err
	}
	ctx.Set(r, "user", "holden")
	session.Values["id"] = 123
	return true, nil
}
