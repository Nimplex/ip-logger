package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	realip "github.com/tomasen/realip"
)

var (
	port = 8080
	file = "logs.txt"
)


func logRequest(handler http.Handler, f *os.File) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := realip.FromRequest(r)
		
		_, err := f.WriteString(fmt.Sprintf("IP: %s,\n", address))

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("New request: %s\n", address)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.IntVar(&port, "port", 8080, "port that logger listens to.")
	flag.StringVar(&file, "file", "logs.txt", "file that logger will log to.")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "<h1>Error 503</h1><p>Service is currently unavailable, please try again later!</p>")
	})

	f, err := os.Create(file)

	if err != nil {
		log.Fatal(fmt.Sprintf("FATAL %s", err))
	}

	log.Println(fmt.Sprintf("Logger is listening to port %d", port))
	log.Println(fmt.Sprintf("Logger will log to %s", file))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), logRequest(http.DefaultServeMux, f))
	if err != nil {
		log.Fatal(fmt.Sprintf("FATAL %s", err))
	}
}