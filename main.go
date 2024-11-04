package main

import (
	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/log"

	"github.com/zhaohaibing666/wallet-chain-account/chaindispatcher"
	"github.com/zhaohaibing666/wallet-chain-account/config"
	wallet2 "github.com/zhaohaibing666/wallet-chain-account/rpc/account"
)

func main() {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	config, err := config.New(*f)
	if err != nil {
		panic(err)
	}

	dispatcher, err := chaindispatcher.New(config)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer grpcServer.GracefulStop()

	wallet2.RegisterWalletAccoutServiceServer(grpcServer, dispatcher)

	listen, err := net.Listen("tcp", ":"+config.Server.Port)
	if err != nil {
		log.Error("net listen failed", "err", err)
		panic(err)
	}
	reflection.Register(grpcServer)

	log.Info("dapplink wallet rpc services start success", "port", config.Server.Port)

	if err := grpcServer.Serve(listen); err != nil {
		log.Error("grpc server serve failed", "err", err)
		panic(err)
	}

}
