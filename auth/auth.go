package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	singingKey := []byte(os.Getenv("SIGNINGKEY"))

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["authorized"] = true
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	tokenString, err := token.SignedString(singingKey)
	if err != nil {
		fmt.Println("Error signing token: ", err)
	}

	return tokenString, nil

}

func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
			return
		}

		singingKey := []byte(os.Getenv("SIGNINGKEY"))

		fullstring := r.Header["Authorization"][0]
		if fullstring == "" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "User not authorized please login!")
			return
		}

		tokenString := strings.Split(fullstring, " ")

		token, err := jwt.Parse(tokenString[0], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%v", "There was an error in parsing token.")
			}
			return singingKey, nil
		})

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Token is either invalid or expired, please login!")
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		exp, _ := claims.GetExpirationTime()

		if time.Until(exp.Time) < 1*time.Minute {
			claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			newTkn, err := token.SignedString(singingKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				fmt.Fprintln(w, "Error creating signature")
			}

			w.Header().Set("Authorization", newTkn)
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
		}
	})

}

func GetUser(authToken string) (string, error) {

	tokenString := strings.Split(authToken, " ")

	signingkey := []byte(os.Getenv("SIGNINGKEY"))

	token, err := jwt.Parse(tokenString[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", "There was an error in parsing token.")
		}
		return signingkey, nil
	})
	if err != nil {
		return "", err
	}
	// ... error handling

	// do something with decoded claims
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error decoding user info from token")
	}

	username := claim["username"].(string)

	return username, nil
}
