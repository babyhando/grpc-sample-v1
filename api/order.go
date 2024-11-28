package api

import (
	"context"
	"proto/quera/pb"
)

type orderGRPCApi struct {
	pb.UnimplementedOrderServiceServer
}

func NewOrderGRPCServer() pb.OrderServiceServer {
	return new(orderGRPCApi)
}

func (s *orderGRPCApi) GetOrder(ctx context.Context, req *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	return &pb.GetOrderResponse{
		Orders: []*pb.Order{{
			Id:        req.ID,
			Quantity:  1,
			TotalBill: 550000,
			IpgMethod: "saman",
			Items: []*pb.OrderItem{
				{
					Id:              1,
					OrderId:         1,
					ItemName:        "discrete structures book",
					Quantity:        1,
					UnitPrice:       550000,
					ItemDescription: "Discrete mathematics book by donal knuth",
				},
			},
			Status: pb.OrderStatus_PAID,
			User: &pb.User{
				Id:        1,
				FirstName: "dara",
				LastName:  "nasibi",
			},
			ExternalPayment: &pb.Order_DigiPay{
				DigiPay: 50000,
			},
			Headers: map[string]string{
				"SAMAN-X-UUID": "xxx-xxx-xxx",
			},
		}},
	}, nil
}
