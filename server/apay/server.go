package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"grpcdemo/common/registry"
	"grpcdemo/proto/payment"
)

/*
	https://github.com/mbobakov/grpc-consul-resolver
	https://github.com/helpleness/IMChatRpc/tree/main
	https://github.com/baiyuze/go-development-template
*/

type paymentServer struct {
	payment.UnimplementedPaymentServer
	channel string
}

func (s *paymentServer) Pay(ctx context.Context, req *payment.PayRequest) (*payment.PayReply, error) {
	log.Printf("[%s] received order=%s amount=%.2f currency=%s", s.channel, req.OrderId, req.Amount, req.Currency)
	return &payment.PayReply{
		Success: true,
		Message: fmt.Sprintf("[%s] processed order %s", s.channel, req.OrderId),
	}, nil
}

func main() {
	ip := "127.0.0.1"
	port := 53001     // 不同服务要改端口
	channel := "apay" // 比如 "wechat" / "alipay"
	serviceName := "payment-" + channel
	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, ip, port)

	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	paySrv := grpc.NewServer()
	payment.RegisterPaymentServer(paySrv, &paymentServer{channel: channel})

	// 注册健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(paySrv, healthServer)
	// 设置某个服务为健康状态
	healthServer.SetServingStatus("helloworld.Greeter", grpc_health_v1.HealthCheckResponse_SERVING)

	// 注册到 Consul
	reg := registry.NewConsulRegistry()
	if e := reg.Register(serviceName, serviceID, ip, port); e == nil {
		fmt.Printf("[%s] service registered to consul", channel)
	}
	// go func() {
	// 	utils.RetryForever(func() error {
	// 		fmt.Println("register consule ...")
	// 		e := reg.Register(serviceName, serviceID, ip, port)
	// 		if e == nil {
	// 			fmt.Printf("[%s] service registered to consul", channel)
	// 		} else {
	// 			fmt.Println(e.Error())
	// 		}
	// 		return e
	// 	}, utils.DefaultRetryConfig())
	// }()

	// 优雅退出
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Printf("[%s] shutting down", channel)
		reg.Deregister(serviceID)
		log.Println("consul deregister ")
		paySrv.GracefulStop()
		log.Println("grpc server Stoped")
	}()

	log.Printf("[%s] listening at %v", channel, lis.Addr())
	if err := paySrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
