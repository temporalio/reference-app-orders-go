package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/encryption"
)

const (
	defaultPort = 8089
	defaultURL  = "http://localhost:8233"
)

var (
	port int
	url  string
)

var rootCmd = &cobra.Command{
	Use:   "codec-server",
	Short: "Codec Server decrypts payloads displayed by Temporal CLI and Web UI",
	RunE: func(*cobra.Command, []string) error {

		if url != "" {
			log.Printf("Codec Server will support Temporal Web UI at %s\n", url)

			if strings.HasSuffix(url, "/") {
				// In my experience, a slash character at the end of the URL will
				// result in a "Codec server could not connect" in the Web UI and
				// the cause will not be obvious. I don't want to strip it off, in
				// case there really is a valid reason to have one, but warning the
				// user could help them to more quickly spot the problem otherwise.
				log.Println("Warning: Temporal Web UI base URL ends with '/'")
			}
		}

		log.Printf("Starting Codec Server on port %d\n", port)
		encryption.StartCodecServer(port, url)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", defaultPort,
		"Port number on which this Codec Server will listen for requests")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", defaultURL,
		"Temporal Web UI base URL (option enables CORS for that origin)")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
