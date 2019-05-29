package api

import (
	"encoding/json"
	"mime"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thanatkrit289/wongnok/internal/auth"
)

// API handler
type API struct {
	Auth *auth.Auth
}

type AuthService interface {

}

// Handler returns api's handler
func (api API) Handler() http.Handler {
	router := httprouter.New()

	// auth
	router.POST("/auth/signup", api.authSignUp)

	router.POST("/auth/signin", api.authSignIn)

	router.POST("/auth/signout", api.authSignOut)

	return router
}

func decodeJSON(r *http.Request, v interface{}) error {
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if mt != "application/json" {
		return json.NewDecoder(r.Body).Decode(v)
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}
