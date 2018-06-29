package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// MySQL
var db *sql.DB

// Redis
var goRedis *redis.Client

func init() {
	connectDb()
	connectRedis()
	move()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/test", test).Methods("GET")

	err := http.ListenAndServe(":8888", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func allowOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8989")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header, x-xsrf-token")
}

func index(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)
}

func test(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)
}
