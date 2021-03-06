package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	//	Flags
	port           = flag.Int("port", 3000, "The port to listen on")
	authEmail      = flag.String("authEmail", "ReplaceWithSvcAcctEmail", "Service account email address")
	authSubject    = flag.String("authSubject", "user@domain.com", "Impersonated user email address")
	allowedOrigins = flag.String("allowedOrigins", "*", "A comma-separated list of valid CORS origins")
	keyFilePath    = flag.String("keyFile", "key.pem", "The location of the PEM encoded private key")
)

func parseEnvironment() {
	//	Check for the listen port
	if env_port := os.Getenv("CALENDAR_PORT"); env_port != "" {
		*port, _ = strconv.Atoi(env_port)
	}

	//	Check for allowed origins
	if env_origins := os.Getenv("CALENDAR_ALLOWED_ORIGINS"); env_origins != "" {
		*allowedOrigins = env_origins
	}

	//	Auth email
	if env_auth_email := os.Getenv("CALENDAR_AUTHEMAIL"); env_auth_email != "" {
		*authEmail = env_auth_email
	}

	//	Auth subject
	if env_auth_subject := os.Getenv("CALENDAR_AUTHSUBJECT"); env_auth_subject != "" {
		*authSubject = env_auth_subject
	}

	//	Key file path
	if env_keyfilepath := os.Getenv("CALENDAR_KEYFILEPATH"); env_keyfilepath != "" {
		*keyFilePath = env_keyfilepath
	}

}

func main() {
	//	Parse environment variables:
	parseEnvironment()

	//	Parse the command line for flags:
	flag.Parse()

	//	Read the key file in:
	keydata, err := ioutil.ReadFile(*keyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf := &jwt.Config{
		Email: *authEmail,
		// The contents of your RSA private key or your PEM file
		// that contains a private key.
		// If you have a p12 file instead, you
		// can use `openssl` to export the private key into a pem file.
		//
		//    $ openssl pkcs12 -in key.p12 -passin pass:notasecret -out key.pem -nodes
		//
		// The field only supports PEM containers with no passphrase.
		// The openssl command will convert p12 keys to passphrase-less PEM containers.
		PrivateKey: keydata,
		Scopes: []string{
			calendar.CalendarScope,
			calendar.CalendarReadonlyScope,
		},
		TokenURL: google.JWTTokenURL,
		// If you would like to impersonate a user, you can
		// create a transport with a subject. The following GET
		// request will be made on the behalf of user@example.com.
		// Optional.
		Subject: *authSubject,
	}

	// Initiate an http.Client, the following GET request will be
	// authorized and authenticated on the behalf of user@example.com.
	client := conf.Client(oauth2.NoContext)

	r := mux.NewRouter()
	r.HandleFunc("/calendar/{calendarid}", func(w http.ResponseWriter, r *http.Request) {

		//	Parse the calendarid from the url
		id := mux.Vars(r)["calendarid"]

		//	Get a connection to the calendar service
		//	If we have errors, return them using standard HTTP service method
		svc, err := calendar.New(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//	Get the list of events from now until the end of today
		now := time.Now()
		end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Add(24 * time.Hour)

		events, err := svc.Events.List(id).
			TimeMin(now.Format(time.RFC3339)).
			TimeMax(end.Format(time.RFC3339)).
			SingleEvents(true).
			OrderBy("startTime").Do()

		//	If we have errors, return them using standard HTTP service method
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//	Set the content type header and return the JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(events)
	})

	//	CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(*allowedOrigins, ","),
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	//	Indicate what port we're starting the service on
	portString := strconv.Itoa(*port)
	fmt.Println("Allowed origins: ", *allowedOrigins)
	fmt.Println("Starting server on :", portString)
	http.ListenAndServe(":"+portString, handler)
}
