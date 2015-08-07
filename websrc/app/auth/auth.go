package auth

import (
	"encoding/json"
	"errors"
	//	"fmt"
	"io/ioutil"
	"net/http"

	ctx "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore(
	[]byte(securecookie.GenerateRandomKey(64)), //Signing key
	[]byte(securecookie.GenerateRandomKey(32)))

var credentials auth

type auth struct {
	User     string
	Password []byte
}

func Login(r *http.Request) (bool, error) {
	err := errors.New("Error")
	username := r.FormValue("username")
	password := r.FormValue("password")
	session, _ := Store.Get(r, "admin")
	err = bcrypt.CompareHashAndPassword([]byte(password), credentials.Password)
	if err != nil {
		return false, err
	}
	if username != credentials.User {
		return false, errors.New("Wrong Username")
	}
	//This is the serialized information from username if we were using a real model
	ctx.Set(r, "user", "holden")
	session.Values["id"] = 123
	return true, nil
}

func Init() {
	creds, err := ioutil.ReadFile("./websrc/cred.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(creds, &credentials)
	if err != nil {
		panic(err)
	}
}
