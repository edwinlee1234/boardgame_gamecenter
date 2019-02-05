package main

import (
	"context"
	"log"
	"net"

	pb "boardgame_gamecenter/proto"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// Redis
var goRedis *redis.Client

// 遊戲中心
var gameCenter *Center

// server is used to implement helloworld.GreeterServer.
type server struct{}

func init() {
	gameCenter = newCenter()

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
