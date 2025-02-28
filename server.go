package main

import (
	// Go Internal Packages
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// External Packages
	"github.com/gorilla/mux"
)

var POSTGRES_URI = ""
var REDIS_URI = ""

func main() {
	db, err := NewPostgresSQL(POSTGRES_URI)
	if err != nil {
		log.Fatalf("Could not initialize Postgres due to: %s", err)
	}
	defer db.Close()

	redis, err := NewRedis(REDIS_URI)
	if err != nil {
		log.Fatalf("Could not initialize Redis due to: %s", err)
	}

	renderJSON := func(w http.ResponseWriter, val interface{}, statusCode int) {
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(val)
	}

	router := mux.NewRouter()
	router.HandleFunc("/emp-salary/{empID}", func(w http.ResponseWriter, r *http.Request) {
		empID := mux.Vars(r)["empID"]

		val, err := redis.Get(r.Context(), empID)
		if err == nil {
			renderJSON(w, &val, http.StatusOK)
			return
		}

		empSalary, err := db.GetEmployeeSalary(empID)
		if err != nil {
			renderJSON(w, &Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		_ = redis.Set(r.Context(), empSalary)
		renderJSON(w, &empSalary, http.StatusOK)
	})

	fmt.Println("Starting server :8080")
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8888",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
