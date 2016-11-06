// General purpose data capture API.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jgrprior/crs/capture"
	"gopkg.in/mgo.v2"
)

var (
	port            int
	mongoURL        string
	defaultMongoURL = os.Getenv("GOCAPTURE_MONGOURL")
)

func init() {
	const (
		defaultPort          = 8888
		usageDefaultPort     = "port to listen on"
		usagedefaultMongoURL = "MongoDb URL. Falls back to $GOCAPTURE_MONGOURL"
	)
	flag.IntVar(&port, "port", defaultPort, usageDefaultPort)
	flag.StringVar(&mongoURL, "mongourl", defaultMongoURL, usagedefaultMongoURL)
}

func main() {
	flag.Parse()

	if mongoURL == "" {
		mongoURL = "localhost"
	}

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	if err := session.Ping(); err != nil {
		panic(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Handling request")

		var e capture.Entry = req.Context().Value(capture.PayloadContextKey).(capture.Entry)
		if saveErr := e.Save(session.Copy()); saveErr != nil {
			// Let the panic handler deal
			panic(saveErr)
		}

		resp := &capture.JSONResponse{http.StatusOK, e.PublicID}
		w.Header().Set("Content-Type", "aplication/json")
		w.WriteHeader(http.StatusOK)
		s, _ := resp.JSONEncode()
		io.WriteString(w, s)
	})

	mux := http.NewServeMux()
	mux.Handle("/campaign",
		capture.PanicHandler(
			capture.AuthHandler(
				capture.PostHandler(
					capture.ValidJSONHandler(handler)))))
	log.Printf("Listening on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
