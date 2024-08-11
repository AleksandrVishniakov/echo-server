package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

func main() {
	port := flag.Int("port", 8080, "listening http port")
	host := flag.String("host", "localhost", "listening http host")
	flag.Parse()

	if port == nil && host == nil {
		log.Fatal("--port or --host is not defined")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandler)

	log.Println("Listening on", fmt.Sprintf("%s:%d", *host, *port))

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), mux); err != nil {
		log.Fatal("server error: ", color.RedString(err.Error()))
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Print()
	cyanBold := color.New(color.FgCyan, color.Bold)
	whiteBold := color.New(color.FgWhite, color.Bold)
	_, err := cyanBold.Print(r.Method + " ")
	if err != nil {
		log.Fatal("failed to write output: ", err.Error())
	}

	fmt.Print(r.Host + r.URL.Path)

	if len(r.URL.Query()) > 0 {
		fmt.Println(color.BlackString("?"))

		for key, values := range r.URL.Query() {
			fmt.Printf("  %s=%s%s\n", key, color.CyanString(strings.Join(values, ", ")), color.BlackString("&"))
		}
	} else {
		fmt.Println()
	}

	fmt.Println()

	if len(r.Header) > 0 {
		_, err = whiteBold.Printf("Headers:\n")
		if err != nil {
			log.Fatal("failed to write output: ", err.Error())
		}

		for key, values := range r.Header {
			fmt.Printf("%s: %s\n", key, color.CyanString(strings.Join(values, ", ")))
		}

		fmt.Println()
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("failed to read body: ", err.Error())
	}

	defer r.Body.Close()

	if len(body) > 0 {
		_, err = whiteBold.Printf("Body:\n")
		if err != nil {
			log.Fatal("failed to write output: ", err.Error())
		}

		fmt.Printf("%s\n", string(body))
		fmt.Println()
	}
}
