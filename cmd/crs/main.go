// General purpose data capture API.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jgrprior/crs"
)

func main() {

	port := flag.String("port", getEnv("CRSPORT", "8080"), "Port to listen on")
	dburl := flag.String("db-url", getEnv("CRSDBURL", "localhost/capture"), "Database connection URL")
	dbtbl := flag.String("db-table", getEnv("CRSDBTBL", "entry"), "Database table or collection")
	authUser := flag.String("username", getEnv("CRSUSER", "user"), "Basic auth user name")
	authPass := flag.String("password", getEnv("CRSPASS", "pass"), "Basic auth password")
	capturePath := flag.String("path", getEnv("CRSPATH", "campaign"), "Path to capture handler")
	flag.Parse()

	db, err := crs.Connect(*dburl, *dbtbl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	handler := crs.NewHandler(db, *authUser, *authPass)
	mux.Handle(fmt.Sprintf("/%s", *capturePath), handler)

	log.Printf("Listening on http://%s:%s", *capturePath, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), mux))
}

// getEnv returns environment variable at key, or def if key is not set.
func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}
