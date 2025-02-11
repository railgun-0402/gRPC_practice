package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "mygrpc/pkg/grpc"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func NewMyServer() *myServer {
	return &myServer{}
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func main() {
	// 8080ポートでListener作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバー作成
	server := grpc.NewServer()

	hellopb.RegisterGreetingServiceServer(server, NewMyServer())

	reflection.Register(server)

	// 無名関数でサーバー起動する
	go func () {
		log.Printf("start gRPC server port: %v", port)
		server.Serve(listener)
	} ()

	// どうやらこれで Ctrl + C でシャットダウンされるらしい
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit // ここでシグナルを受け取るまで以降の処理はされない
	log.Println("stopping gRPC server...")

	// GracefulStop() は進行中のリクエストを完了させた後にサーバーを停止する
	// Stop()は即終了なので、こっちのが安全そう
	server.GracefulStop()
}