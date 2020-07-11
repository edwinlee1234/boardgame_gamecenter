package main

import (
	"context"

	pb "boardgame_gamecenter/proto"
)

// 新增遊戲
func (s *server) CreateGame(ctx context.Context, in *pb.CreateGameRequest) (*pb.CreateGameReply, error) {
	err := gameCenter.CreateGame(in.GameID, in.GameType, in.Players, map[string]interface{}{})

	if err != nil {
		return &pb.CreateGameReply{
			State: "error",
		}, err
	}

	return &pb.CreateGameReply{
		State: "success",
	}, nil
}

// 遊戲資訊
func (s *server) GameInfo(ctx context.Context, in *pb.GameInfoRequest) (*pb.GameInfoReply, error) {
	err := gameCenter.GameInfo(in.UserID, in.GameID, in.GameType)

	if err != nil {
		return &pb.GameInfoReply{
			State: "error",
		}, err
	}

	return &pb.GameInfoReply{
		State: "success",
	}, nil
}

// Action
func (s *server) JaipurAction(ctx context.Context, in *pb.JaipurActionRequest) (*pb.JaipurActionReply, error) {
	err := gameCenter.ActionProcess(in.UserID, in.GameID, in.GameType, in.Action)

	if err != nil {
		return &pb.JaipurActionReply{
			State: "error",
		}, err
	}

	return &pb.JaipurActionReply{
		State: "success",
	}, nil
}
