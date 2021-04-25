package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	delivery "github.com/goagile/oshp/pkg/api/grpc/delivery"
	order "github.com/goagile/oshp/pkg/api/grpc/order"
	"google.golang.org/grpc"
)

var (
	deliveryPort = flag.Int("delivery_port", 8086, "delivery_port")
)

var (
	orderServerAddress = flag.String("order_server", "localhost:8085", "Address of gRPC server 'localhost:8085'")
	orderClient        order.OrderClient
)

func main() {
	flag.Parse()
	p := fmt.Sprintf("localhost:%d", *deliveryPort)
	fmt.Println("Delivery gRPC server listen on", p)
	tcp, err := net.Listen("tcp", p)
	if err != nil {
		log.Fatalln("Error With TCP Listen", err)
	}

	//
	// Order Client
	//
	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithInsecure())
	conn, err := grpc.Dial(*orderServerAddress, dialOpts...)
	if err != nil {
		log.Fatalln("grpc.Dial", err)
	}
	defer conn.Close()
	orderClient = order.NewOrderClient(conn)

	//
	// Delivery Server
	//
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	delivery.RegisterDeliveryServer(s, &deliveryServer{})
	s.Serve(tcp)
}

type deliveryServer struct {
	delivery.UnimplementedDeliveryServer
}

func (*deliveryServer) ScheduleDelivery(ctx context.Context, req *delivery.ScheduleDeliveryRequest) (*delivery.ScheduleDeliveryResponse, error) {
	log.Println("\n\nRecieve Schedule Request from Order", req)
	resp := new(delivery.ScheduleDeliveryResponse)
	resp.OrderId = req.OrderId
	log.Println("calculating ...")
	go calculate(req.OrderId)
	return resp, nil
}

func calculate(orderID int32) {
	ctx := context.Background()
	time.Sleep(3 * time.Second)
	log.Println("calculated")
	req := &order.UpdateOrderRequest{
		OrderId:      orderID,
		DeliveryDate: fmt.Sprintf("%v.04.2021", orderID),
	}
	log.Println("orderClient.UpdateOrder", req)
	resp, err := orderClient.UpdateOrder(ctx, req)
	if err != nil {
		log.Fatalln("client.SayHello", err)
	}
	log.Println("orderClient.UpdateOrder", resp)
}
