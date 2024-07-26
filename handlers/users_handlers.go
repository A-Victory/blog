package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/A-Victory/blog/auth"
	"github.com/A-Victory/blog/database/conn"
	"github.com/A-Victory/blog/models"
	"golang.org/x/crypto/bcrypt"
)

type HttpHandler struct {
	db *conn.DB
	va *auth.Validation
}

type Config struct {
	Database  *conn.DB
	Validator *auth.Validation
}

type customResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func NewHttpHandler(opt *Config) *HttpHandler {
	return &HttpHandler{
		db: opt.Database,
		va: opt.Validator,
	}
}

func (httpConfig *HttpHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	newUser := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"message": err.Error()}}
		json.NewEncoder(w).Encode(response)
	}

	hashedpass, err := hashpassword(newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := customResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"message": err.Error()}}
		json.NewEncoder(w).Encode(response)
	}

	if err = httpConfig.va.ValidateUserInfo(newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := customResponse{Status: http.StatusBadRequest, Message: "invalid request", Data: map[string]interface{}{"message": err.Error()}}
		json.NewEncoder(w).Encode(response)
	}

	newUser.Password = hashedpass

	field, err := httpConfig.searchUser(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	if field == "" {
		id, err := httpConfig.db.SaveUser(newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		type userResponse struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			UserID   int    `json:"userID"`
		}

		resp := userResponse{
			Email:    newUser.Email,
			Username: newUser.Username,
			UserID:   id,
		}

		w.WriteHeader(http.StatusOK)
		response := customResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"user": resp}}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := customResponse{Status: http.StatusBadRequest, Message: fmt.Sprintf("%s in use", field), Data: map[string]interface{}{"msg": fmt.Sprintf("%s already exists, try again...", field)}}
		json.NewEncoder(w).Encode(response)
		return
	}

}

func (httpConfig *HttpHandler) Login(w http.ResponseWriter, r *http.Request) {

	login := models.LoginDetails{}

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"message": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	// write another for when the email is not found in the database

	user, err := httpConfig.db.GetUser("email", login.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	if user == (models.User{}) {
		w.WriteHeader(http.StatusNotFound)
		response := customResponse{Status: http.StatusNotFound, Message: "email not registered", Data: map[string]interface{}{"msg": "email not associated to a user, proceed to register page to signup..."}}
		json.NewEncoder(w).Encode(response)
		return
	}

	valid := comparePassword(login.Password, user.Password)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		response := customResponse{Status: http.StatusUnauthorized, Message: "error", Data: map[string]interface{}{"message": "incorrect password"}}
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := auth.GenerateJWT(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "failed to generate token: " + err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
	response := customResponse{Status: http.StatusOK, Message: "login successful", Data: map[string]interface{}{"authorization": fmt.Sprintf("your generated token is %s attach to subsequest request with the header %s", token, "Authorization")}}
	json.NewEncoder(w).Encode(response)

	// if successful, add the authorization header to response and return token as json response as well

}

func (httpConfig *HttpHandler) Profile(w http.ResponseWriter, r *http.Request) {

	user, err := httpConfig.getUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "failed to retrieve user's details: " + err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := customResponse{Status: http.StatusOK, Message: "user's profile", Data: map[string]interface{}{"user": user}}
	json.NewEncoder(w).Encode(response)

	// get the user's details from the database
}

func hashpassword(password string) (string, error) {
	encrytedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(encrytedPassword), nil
}

func comparePassword(inputPassword, dbPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(inputPassword)); err != nil {
		return false
	}

	return true
}

func (httpConfig *HttpHandler) getUser(r *http.Request) (models.User, error) {
	token := r.Header.Get("Authorization")
	username, err := auth.GetUser(token)
	if err != nil {
		return models.User{}, err
	}

	user, err := httpConfig.db.GetUser("username", username)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (httpConfig *HttpHandler) searchUser(newUser models.User) (field string, err error) {
	checks := []struct {
		field string
		value string
	}{
		{"email", newUser.Email},
		{"username", newUser.Username},
	}

	for _, check := range checks {
		user, err := httpConfig.db.GetUser(check.field, check.value)
		if err != nil {
			if err != sql.ErrNoRows {
				return "", nil
			}
		}
		if err == nil {
			if user.Username == newUser.Username {
				return "username", nil
			}
			if user.Email == newUser.Email {
				return "email", nil
			}
		}
	}

	return "", nil
}
