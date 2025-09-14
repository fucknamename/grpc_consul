package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"grpcdemo/common/registry"
	"grpcdemo/proto/payment"
)

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
	reg := registry.NewConsulRegistry()

	rand.New(rand.NewSource(time.Now().UnixNano()))

	r := gin.Default()
	// 通用入口：/pay/:channel
	r.POST("/pay/:channel", func(c *gin.Context) {
		channel := c.Param("channel")
		serviceName := "payment-" + channel

		// 从 Consul 发现服务
		services, err := reg.Discover(serviceName)
		if err != nil || len(services) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no service found"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "grpc dial failed"})
			return
		}
		defer conn.Close()

		client := payment.NewPaymentClient(conn)

		var req struct {
			OrderID   string            `json:"order_id"`
			Amount    float64           `json:"amount"`
			Currency  string            `json:"currency"`
			ExtraData map[string]string `json:"extra_data"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 调用 gRPC
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 超时5秒
		defer cancel()

		resp, err := client.Pay(ctx, &payment.PayRequest{
			Channel:  channel,
			OrderId:  req.OrderID,
			Amount:   req.Amount,
			Currency: req.Currency,
			// ExtraData: req.ExtraData,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})
	})

	log.Println("Gateway listening on :8080")
	r.Run(":8080")
}
