package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	delivery "github.com/goagile/oshp/pkg/api/grpc/delivery"
	order "github.com/goagile/oshp/pkg/api/grpc/order"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

var (
	ordergRESTPort = flag.Int("order_rest_port", 8084, "Order REST API port")
	ordergRPCPort  = flag.Int("order_grpc_port", 8085, "Order gRPC API port")
)

var (
	deliveryServerAddress = flag.String("delivery_server", "localhost:8086", "Address of gRPC server 'localhost:8086'")
	deliveryClient        delivery.DeliveryClient
)

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// Delivery gRPC Client
	//
	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithInsecure())
	conn, err := grpc.Dial(*deliveryServerAddress, dialOpts...)
	if err != nil {
		log.Fatalln("grpc.Dial", err)
	}
	defer conn.Close()
	deliveryClient = delivery.NewDeliveryClient(conn)

	// ctx := context.Background()
	// go scheduleDelivery(ctx)

	//
	// Order gRPC Server
	//
	p := fmt.Sprintf("localhost:%d", *ordergRPCPort)
	fmt.Println("Order gRPC server listen on", p)
	tcp, err := net.Listen("tcp", p)
	if err != nil {
		log.Fatalln("Error With TCP Listen", err)
	}
	var serverOpts []grpc.ServerOption
	s := grpc.NewServer(serverOpts...)
	order.RegisterOrderServer(s, &orderServer{})
	go s.Serve(tcp)

	//
	// Order REST API Server
	//
	ordergRESTPortStr := fmt.Sprintf(":%v", *ordergRESTPort)
	log.Println("REST API listen at", ordergRESTPortStr)
	if err := setupRESTServer().Run(ordergRESTPortStr); err != nil {
		log.Fatalln("Run REST server error", err)
	}
}

func setupRESTServer() *gin.Engine {
	r := gin.Default()

	//
	// curl -i \\
	// 		-X POST http://127.0.0.1:8084/orders
	//      -H 'Content-Type: application/json'
	//      -d '{"user_id":"777"}'
	//
	r.POST("/orders", CreateOrder)

	return r
}

// CreateOrder - POST callback
func CreateOrder(c *gin.Context) {
	var r CreateOrderRequest
	if err := c.BindJSON(&r); err != nil {
		log.Println("CreateOrder BindJSON", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "fail to create order"},
		)
		return
	}
	fmt.Println("Create Order Request", r)

	//
	// Call Delivery Service
	//
	// time.Sleep(time.Duration(i) * 3 * time.Second)
	req := &delivery.ScheduleDeliveryRequest{
		OrderId: 333,
	}
	log.Println("\n\nSend Request to Delivery Service", req)
	resp, err := deliveryClient.ScheduleDelivery(c.Request.Context(), req)
	if err != nil {
		log.Println("deliveryClient.ScheduleDelivery", err)
	}
	log.Println("deliveryClient.ScheduleDelivery Resp -> ", resp)

	c.JSON(http.StatusCreated, gin.H{"data": "delivery scheduled"})
}

type CreateOrderRequest struct {
	UserID string `json:"user_id"`
}

type orderServer struct {
	order.UnimplementedOrderServer
}

func (*orderServer) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.UpdateOrderResponse, error) {
	log.Println("req.OrderId:", req)
	resp := new(order.UpdateOrderResponse)
	resp.OrderId = req.OrderId
	resp.OrderStatus = "OK"
	return resp, nil
}

func scheduleDelivery(ctx context.Context) {
	for i := 1; i < 4; i++ {
		time.Sleep(time.Duration(i) * 3 * time.Second)
		req := &delivery.ScheduleDeliveryRequest{
			OrderId: int32(i),
		}
		log.Println("\n\nSend Request to Delivery Service", req)
		resp, err := deliveryClient.ScheduleDelivery(ctx, req)
		if err != nil {
			log.Fatalln("deliveryClient.ScheduleDelivery", err)
		}
		log.Println("deliveryClient.ScheduleDelivery Resp -> ", resp)
	}
}
