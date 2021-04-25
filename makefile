
.PHONY: build

build:
	$(shell ./build/order_grpc.sh)
	$(shell ./build/delivery_grpc.sh)
	$(shell ./build/order.sh)
	$(shell ./build/delivery.sh)

.PHONY: clear

clear:
	rm delivery_service
	rm order_service
