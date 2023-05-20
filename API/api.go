package API

import (
	"awesomeProject/database"
	"awesomeProject/jwt_auth"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

type API struct {
	router *mux.Router
	db     *sql.DB
}

func NewAPI() (*API, error) {
	res := new(API)
	res.router = mux.NewRouter()

	return res, nil
}

func (a *API) Start() error {
	db, err := database.Connect(login + ":" + pass + "@tcp(localhost:3306)/universities")
	if err != nil {
		fmt.Println(err)
		return err
	}
	a.db = db

	a.router.HandleFunc("/signin", a.signin())
	a.router.HandleFunc("/platform", jwt_auth.VerifyJWT(a.handleGamePlatform()))
	a.router.HandleFunc("/game", jwt_auth.VerifyJWT(a.handleAddGame()))
	a.router.HandleFunc("/update/release_year", jwt_auth.VerifyJWT(a.handleUpdateReleaseYear()))

	return http.ListenAndServe(port, a.router)
}

func (a *API) Stop() {
	fmt.Println("Stopping API...")
	a.db.Close()
	fmt.Println("API stopped...")
}
func (a *API) signin() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var creds jwt_auth.Credentials
		// Get the JSON body and decode into credentials
		err := json.NewDecoder(request.Body).Decode(&creds)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get the expected password from our in memory map
		expectedPassword, ok := jwt_auth.Users[creds.Username]

		if !ok || expectedPassword != creds.Password {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		adminClaim := false
		if creds.Username == "admin" {
			adminClaim = true
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &jwt_auth.Claims{
			Admin:    adminClaim,
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(jwt_auth.JwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(writer, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
}
func (a *API) handleGamePlatform() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}

		var msg gamePlatform
		err = json.Unmarshal(body, &msg)
		if err != nil {
			http.Error(writer, "error during unmarshal", http.StatusBadRequest)
			return
		}

		switch request.Method {
		case "POST":
			_, err = a.db.Exec(database.AddPlatformQuery, msg.PlatformName)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		case "DELETE":
			ins, err := a.db.Prepare(database.DeletePlatformQuery1)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			ins.Exec(msg.PlatformName)

			ins, err = a.db.Prepare(database.DeletePlatformQuery2)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			ins.Exec(msg.PlatformName)
			ins, err = a.db.Prepare(database.DeletePlatformQuery3)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			ins.Exec(msg.PlatformName)
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleAddGame() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}

		var msg AddGame
		err = json.Unmarshal(body, &msg)
		if err != nil {
			http.Error(writer, "error during unmarshal", http.StatusBadRequest)
			return
		}

		_, err = a.db.Exec(database.AddGameQuery, msg.GenreName, msg.GameName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleUpdateReleaseYear() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}

		var msg UpdateReleaseYear
		err = json.Unmarshal(body, &msg)
		if err != nil {
			http.Error(writer, "error during unmarshal", http.StatusBadRequest)
			return
		}

		_, err = a.db.Exec(database.UpdateGameReleaseYear, msg.Year, msg.GameName /*msg.PublisherName, msg.PublisherName*/)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
