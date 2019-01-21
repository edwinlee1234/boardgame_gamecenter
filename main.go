package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	jaipurClass "boardgame_gamecenter/games/jaipur"
	pb "boardgame_gamecenter/proto"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

// var (
// 	key   = []byte("super-secret-key")
// 	store = sessions.NewCookieStore(key)
// )

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

func allowOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8989")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header, x-xsrf-token")
}

var game *jaipurClass.Jaipur

func index(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)
	// fmt.Println("OK!")
	game = jaipurClass.NewJaipur(1, nil)
	game.Init(map[int32]string{1: "111", 2: "222"})
}

func test(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)
}

func show(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)
	info := game.Info()
	fmt.Printf("deskCard: %v \n", info.DeskCard)
	fmt.Printf("CardsPoint: %v \n", info.CardsPoint)
	fmt.Printf("PlayerCard: %v \n", info.Players)
	fmt.Printf("ActionPlayer: %d \n", info.ActionPlayer)
	fmt.Printf("FoldCardNum: %d \n", info.FoldCardNum)
	fmt.Printf("*********\n")
}

func take(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)
	cardIDtmp, _ := r.URL.Query()["cardID"]
	userIDtmp, _ := r.URL.Query()["userID"]
	cardInt, _ := strconv.Atoi(cardIDtmp[0])
	userInt, _ := strconv.Atoi(userIDtmp[0])
	userID := int32(userInt)
	cardID := int32(cardInt)

	var act jaipurClass.Action
	act.Type = "take"
	act.Take = cardID
	err := game.Action(userID, act)
	if err != nil {
		log.Println(err)
	}
}

func sell(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)

	userIDtmp, _ := r.URL.Query()["userID"]
	userInt, _ := strconv.Atoi(userIDtmp[0])
	userID := int32(userInt)

	var cardsID []int32
	cardIDtmpString, _ := r.URL.Query()["cardID"]
	cardIDtmp := strings.Split(cardIDtmpString[0], ",")
	for _, v := range cardIDtmp {
		cardInt, _ := strconv.Atoi(v)
		cardsID = append(cardsID, int32(cardInt))
	}

	var act jaipurClass.Action
	act.Type = "sell"
	act.Sell = cardsID
	err := game.Action(userID, act)
	if err != nil {
		log.Println(err)
	}
}

func switchCard(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w, r)
	userIDtmp, _ := r.URL.Query()["userID"]
	userInt, _ := strconv.Atoi(userIDtmp[0])
	userID := int32(userInt)

	var selfCardsID []int32
	selftmpString, _ := r.URL.Query()["self"]
	selftmp := strings.Split(selftmpString[0], ",")
	for _, v := range selftmp {
		selfInt, _ := strconv.Atoi(v)
		selfCardsID = append(selfCardsID, int32(selfInt))
	}

	var switchCardsID []int32
	sswitchtmpString, _ := r.URL.Query()["switch"]
	switchtmp := strings.Split(sswitchtmpString[0], ",")
	for _, v := range switchtmp {
		switchInt, _ := strconv.Atoi(v)
		switchCardsID = append(switchCardsID, int32(switchInt))
	}

	var act jaipurClass.Action
	act.Type = "switch"
	act.SwitchSelfCard = selfCardsID
	act.SwitchTargetCard = switchCardsID
	err := game.Action(userID, act)
	if err != nil {
		log.Println(err)
	}
}
