package encryption

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"go.temporal.io/sdk/converter"
)

func newCORSHTTPHandler(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Namespace")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RunCodecServer launches the Codec Server on the specified port, enabling
// CORS for the Temporal Web UI at the specified URL
func RunCodecServer(port int, url string) {
	// The EncryptionKeyID attribute is omitted when creating the Codec
	// instance below because the Codec Server only decrypts. It locates
	// the encryption key ID from the payload's metadata.
	handler := converter.NewPayloadCodecHTTPHandler(&Codec{})

	if url != "" {
		log.Printf("Codec Server will support Temporal Web UI at %s\n", url)
		handler = newCORSHTTPHandler(url, handler)
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(port),
		Handler: handler,
	}

	log.Printf("Starting Codec Server on port %d\n", port)

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
