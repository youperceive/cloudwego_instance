package trans

import (
	"github.com/youperceive/cloudwego_instance/rpc/order/kitex_gen/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderDoc struct {
	order.Order
	ID primitive.ObjectID `bson:"_id"`
}

func OrderItem(orderItemForCreate order.OrderItemForCreate, itemId int64, orderId string) *order.OrderItem {
	orderItem := order.OrderItem{
		Id:        itemId,
		OrderId:   orderId,
		ProductId: orderItemForCreate.ProductId,
		SkuId:     orderItemForCreate.SkuId,
		Count:     orderItemForCreate.Count,
		Price:     orderItemForCreate.Price,
		Ext:       orderItemForCreate.Ext,
	}
	return &orderItem
}
