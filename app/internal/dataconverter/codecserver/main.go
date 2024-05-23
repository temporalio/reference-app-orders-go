package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/temporalio/orders-reference-app-go/app/internal/dataconverter"

	"go.temporal.io/sdk/converter"
)

var portFlag int
var urlFlag string

func init() {
	flag.IntVar(&portFlag, "port", 8089, "Port number on which this codec server will listen for requests")
	flag.StringVar(&urlFlag, "url", "", "Temporal Web UI base URL. This option enables CORS, as required by that UI.")
}

func newCORSHTTPHandler(web string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", web)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Namespace")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	encKeyID := os.Getenv("CLIENT_ENCRYPTION_KEY_ID")
	if encKeyID == "" {
		log.Fatalf("identifier for client encryption key is undefined\n")
	}

	handler := converter.NewPayloadCodecHTTPHandler(&dataconverter.Codec{
		EncryptionKeyID: encKeyID,
	})

	if urlFlag != "" {
		handler = newCORSHTTPHandler(urlFlag, handler)
	} else {
		log.Println("Temporal Web UI support is not enabled")
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(portFlag),
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		_ = srv.Close()
	case err := <-errCh:
		log.Fatal(err)
	}
}
