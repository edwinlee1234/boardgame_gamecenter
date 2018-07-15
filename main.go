package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	pb "./proto"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"google.golang.org/grpc"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

const (
	port = ":50051"
)

// MySQL
var db *sql.DB

// Redis
var goRedis *redis.Client

// 遊戲中心
// var gameCenter *Center

// server is used to implement helloworld.GreeterServer.
type server struct{}

func init() {
	connectDb()
	connectRedis()
	createGrpcServer()
	// gameCenter = newCenter()
	// gameCenter.CreateGame(1, "jaipur")
	// gameClass := gameCenter.FindGames(1)
	// gameClass.Init()
}

// gRPC的func
func (s *server) Ping(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	return &pb.TestReply{
		State: "Pong",
	}, nil
}

// 開啓gRPC服務
func createGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGameCenterServer(s, &server{})
	s.Serve(lis)
}

func main() {
	// 這個服務不會開放http直接打進來
	// 測試用而己，之後拿掉
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
