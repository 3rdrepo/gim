package main

import (
	"gim/config"
	"gim/internal/business/api"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/urlwhitelist"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("business_interceptor", urlwhitelist.Business)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop", zap.Any("signal", s))
		server.GracefulStop()
	}()

	pb.RegisterBusinessIntServer(server, &api.BusinessIntServer{})
	pb.RegisterBusinessExtServer(server, &api.BusinessExtServer{})
	listen, err := net.Listen("tcp", config.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}
