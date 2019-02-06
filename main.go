package main

import (
	"context"
	"log"
	"net"
	"os"

	lib "boardgame_gamecenter/lib"
	pb "boardgame_gamecenter/proto"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// 遊戲中心
var gameCenter *Center

// server is used to implement helloworld.GreeterServer.
type server struct{}

func init() {
	initCenter()
	connectDB()
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
}

func initCenter() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	wsURL := os.Getenv("WS_URL")
	WS := lib.NewWS(wsURL)

	gameCenter = newCenter(WS)
}
