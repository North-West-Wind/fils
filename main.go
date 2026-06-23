package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goware/urlx"
	"github.com/joho/godotenv"
)

var (
	ErrLog = log.New(os.Stderr, "", 0)
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Found .env file. Loading environment...")
		err = godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	passphrase := os.Getenv("PASSPHRASE")
	if passphrase == "" {
		log.Println("No passphrase configured. Everyone can use this service!")
	}

	store := os.Getenv("STORE")
	if store == "" {
		log.Println("No store specified. Defaulting to store.txt...")
		store = "store.txt"
	}

	// Create directory of store if not exist
	os.MkdirAll(path.Dir(store), 0644)

	LoadStore(store)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/s/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to parse ID"))
			return
		}
		url, err := LookupID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to lookup ID"))
			ErrLog.Printf("Failed to lookup ID: %v", err)
			return
		}
		if url == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("ID not found"))
			return
		}
		w.Header().Add("location", url)
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte("Redirect to " + url))
	})

	r.Post("/new", func(w http.ResponseWriter, r *http.Request) {
		if passphrase != "" {
			// Not very secure huh
			if r.Header.Get("authorization") != passphrase {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Passphrase doesn't match"))
				return
			}
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read body"))
			ErrLog.Printf("Failed to read body: %v", err)
			return
		}
		url, err := urlx.NormalizeString(string(data))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to normalize URL"))
			ErrLog.Printf("Failed to normalize URL: %v", err)
			return
		}
		id, err := AddURL(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to store URL"))
			ErrLog.Printf("Failed to store URL: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatUint(id, 10)))
		log.Printf("Added URL %s as ID %d", url, id)
	})

	r.Handle("/*", http.FileServer(http.Dir("./client")))

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}
	http.ListenAndServe(port, r)
}
