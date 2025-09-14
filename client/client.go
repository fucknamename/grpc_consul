package main

import (
	"context"
	"fmt"
	"grpcdemo/common/registry"
	"grpcdemo/proto/payment"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 测试 grpc 客户端

// 保存每个渠道的轮询游标
var counters = make(map[string]*uint64)

func getNextIndex(channel string, n int) int {
	if _, ok := counters[channel]; !ok {
		var c uint64 = 0
		counters[channel] = &c
	}
	// 原子自增，避免并发问题
	idx := atomic.AddUint64(counters[channel], 1)
	return int(idx % uint64(n))
}

func main() {
	// // 建立 gRPC 连接
	// host := "127.0.0.1"
	// port := 53001
	// target := fmt.Sprintf("%s:%d", host, port)
	// conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	reg := registry.NewConsulRegistry()
	// 从 Consul 发现服务
	channel := "apay" // 比如 "wechat" / "alipay"
	serviceName := "payment-" + channel
	services, err := reg.Discover(serviceName)
	if err != nil || len(services) == 0 {
		fmt.Println("未找到服务：" + err.Error())
		return
	}

	// --- 负载均衡策略 ---
	// 随机方式
	// idx := rand.Intn(len(services))

	// 轮询方式
	idx := getNextIndex(channel, len(services))
	srv := services[idx].Service
	target := fmt.Sprintf("%s:%d", srv.Address, srv.Port)
	// 建立 gRPC 连接
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return
	}
	defer conn.Close()

	client := payment.NewPaymentClient(conn)

	// 调用 gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 超时5秒
	defer cancel()

	resp, err := client.Pay(ctx, &payment.PayRequest{
		Channel:  "apay",
		OrderId:  "U50001",
		Amount:   528.8,
		Currency: "IN",
		// ExtraData: req.ExtraData,
	})
	if err != nil {
		return
	}

	fmt.Println(resp)
}
