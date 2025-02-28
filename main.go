package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"live-stream/backend/handler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	r := mux.NewRouter()

	r.HandleFunc("/api/livestream/create", handler.CreateLiveInput).Methods("POST")
	r.HandleFunc("/api/livestream/uploadStream", handler.UploadStream).Methods("POST")
	r.HandleFunc("/api/livestream/join", handler.GetVideos).Methods("GET")

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./frontend/assets"))))

	r.HandleFunc("/broadcaster", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/view/broadcaster.html")
	})

	r.HandleFunc("/viewer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/view/viewer.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
