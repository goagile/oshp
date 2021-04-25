
.PHONY: grpc

grpc:
	$(shell ./build/order_grpc.sh)
	$(shell ./build/delivery_grpc.sh)
