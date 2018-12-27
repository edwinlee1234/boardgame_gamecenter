package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	pb "./proto"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
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
var gameCenter *Center

// server is used to implement helloworld.GreeterServer.
type server struct{}

func init() {
	gameCenter = newCenter()

	connectDb()
	connectRedis()
	createGrpcServer()
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
	// r := mux.NewRouter()

	// r.HandleFunc("/", index).Methods("GET")
	// r.HandleFunc("/test", test).Methods("GET")
	// r.HandleFunc("/show", show).Methods("GET")
	// r.HandleFunc("/take", take).Methods("GET")
	// r.HandleFunc("/sell", sell).Methods("GET")
	// r.HandleFunc("/switchCard", switchCard).Methods("GET")

	// err := http.ListenAndServe(":8888", r)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}
