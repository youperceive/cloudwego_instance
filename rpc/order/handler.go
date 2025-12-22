package main

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/youperceive/cloudwego_instance/rpc/order/kitex_gen/base"
	order "github.com/youperceive/cloudwego_instance/rpc/order/kitex_gen/order"
	"github.com/youperceive/cloudwego_instance/rpc/order/pkg/mongo"
	"github.com/youperceive/cloudwego_instance/rpc/order/pkg/snowflake"
	"github.com/youperceive/cloudwego_instance/rpc/order/pkg/trans"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoOfficial "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Coll = mongo.Cli.Database("order_db").Collection("total")

var (
	internalErrMsg = "Internal Error."
	successMsg     = "Success."
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

func validateCreateReq(req *order.CreateRequest) error {
	// 没有什么要校验的
	_ = req
	return nil
}

// Create implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) Create(ctx context.Context, req *order.CreateRequest) (resp *order.CreateResponse, err error) {
	klogErr := func(msg string) {
		target := ""
		if req != nil {
			target = req.String()
		}
		klog.Error(
			"method: ", "Create",
			"target: ", target,
			"message: ", msg,
		)
	}

	if err = validateCreateReq(req); err != nil {
		klogErr("invalid params." + err.Error())
		resp = &order.CreateResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  err.Error(), // 这里是参数校验错误，可以返回给客户端
			},
		}
		return
	}

	objectId := primitive.NewObjectID()
	orderIdStr := objectId.Hex()

	var items []*order.OrderItem
	for _, item := range req.Items {
		itemId := snowflake.SfNode.Generate().Int64()
		items = append(items, trans.OrderItem(*item, itemId, orderIdStr))
	}

	doc := trans.OrderDoc{
		ID: objectId,
		Order: order.Order{
			Id:         orderIdStr,
			Type:       req.Type,
			Status:     req.Status,
			ReqUserId:  req.ReqUserId,
			RespUserId: req.RespUserId,
			Items:      items,
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
			Ext:        req.Ext,
		},
	}

	_, err = Coll.InsertOne(ctx, doc)
	if err != nil {
		klogErr("mongo insert operator failed. " + err.Error())
		resp = &order.CreateResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  internalErrMsg,
			},
		}
		return
	}

	resp = &order.CreateResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  successMsg,
		},
		OrderId: orderIdStr,
	}

	return
}

func validateUpdateReq(req *order.UpdateRequest) error {
	if req.Id == "" {
		return errors.New("订单ID不能为空")
	}
	return nil
}

// Update implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) Update(ctx context.Context, req *order.UpdateRequest) (resp *order.UpdateResponse, err error) {

	klogErr := func(msg string) {
		target := ""
		if req != nil {
			target = req.String()
		}
		klog.Error(
			"method: ", "Update",
			"req: ", target,
			"message: ", msg,
		)
	}

	if err = validateUpdateReq(req); err != nil {
		klogErr("invalid params: " + err.Error())
		resp = &order.UpdateResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  err.Error(),
			},
		}
		return resp, nil
	}

	val := bson.D{
		bson.E{
			Key:   "updated_at",
			Value: time.Now().Unix(),
		},
	}
	if req.Status != nil {
		val = append(
			val,
			bson.E{
				Key:   "status",
				Value: *req.Status,
			},
		)
	}
	blackList := map[string]bool{"_id": true, "updated_at": true, "status": true, "req_user_id": true}
	for k, v := range req.Ext {
		if !blackList[k] {
			val = append(
				val,
				bson.E{
					Key:   k,
					Value: v,
				},
			)
		}
	}
	data := bson.D{
		bson.E{
			Key:   "$set",
			Value: val,
		},
	}

	objectId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		klogErr("fail to convert req.Id to objectId: " + err.Error())
		resp = &order.UpdateResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  "订单ID格式错误",
			},
		}
		return
	}

	_, err = Coll.UpdateByID(ctx, objectId, data)
	if err != nil {
		if errors.Is(err, mongoOfficial.ErrNoDocuments) {
			klogErr("fail to update. " + err.Error())
			resp = &order.UpdateResponse{
				BaseResp: &base.BaseResponse{
					Code: base.Code_NOT_FOUND,
					Msg:  "订单不存在",
				},
			}
		} else {
			klogErr("fail to update. " + err.Error())
			resp = &order.UpdateResponse{
				BaseResp: &base.BaseResponse{
					Code: base.Code_DB_ERR,
					Msg:  internalErrMsg,
				},
			}
		}
		return
	}

	resp = &order.UpdateResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  successMsg,
		},
	}

	return
}

func validateQueryOrderInfoReq(req *order.QueryOrderInfoRequest) error {
	if req.Id == "" {
		return errors.New("订单ID不能为空")
	}
	return nil
}

// QueryOrderInfo implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) QueryOrderInfo(ctx context.Context, req *order.QueryOrderInfoRequest) (resp *order.QueryOrderInfoResponse, err error) {

	klogErr := func(msg string) {
		target := ""
		if req != nil {
			target = req.String()
		}
		klog.Error(
			"method: ", "QueryOrderInfo",
			"req: ", target,
			"message: ", msg,
		)
	}

	if err = validateQueryOrderInfoReq(req); err != nil {
		klogErr("invalid params: " + err.Error())
		resp = &order.QueryOrderInfoResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  err.Error(),
			},
		}
		return
	}

	objectId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		klogErr("fail to convert req.Id to objectId: " + err.Error())
		resp = &order.QueryOrderInfoResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  "订单ID格式错误",
			},
		}
		return
	}

	filter := bson.M{"_id": objectId}

	var orderDoc trans.OrderDoc
	err = Coll.FindOne(ctx, filter).Decode(&orderDoc)

	if err != nil {
		if errors.Is(err, mongoOfficial.ErrNoDocuments) {
			klogErr("order not found.")
			resp = &order.QueryOrderInfoResponse{
				BaseResp: &base.BaseResponse{
					Code: base.Code_NOT_FOUND,
					Msg:  "订单不存在",
				},
			}
			return
		}
		klogErr("fail to query order: " + err.Error())
		resp = &order.QueryOrderInfoResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  internalErrMsg,
			},
		}
		return
	}

	// 还真是，意料之外
	resp = &order.QueryOrderInfoResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  successMsg,
		},
		Order: &orderDoc.Order,
	}

	return
}

func validateQueryOrderIdReq(req *order.QueryOrderIdRequest) error {
	if (req.Type == order.QueryOrderIdType_REQ_USER || req.Type == order.QueryOrderIdType_RESP_USER) && req.UserId == nil {
		return errors.New("想要通过用户 id 查询订单，但 user_id 为空.")
	}
	if req.Type == order.QueryOrderIdType_EXT_KEY && (req.ExtKey == nil || req.ExtVal == nil) {
		return errors.New("想要通过扩展字段查询订单，但扩展字段或字段值为空.")
	}
	if req.Type != order.QueryOrderIdType_REQ_USER && req.Type != order.QueryOrderIdType_RESP_USER && req.Type != order.QueryOrderIdType_EXT_KEY {
		return errors.New("查询类型不存在.")
	}
	if req.Page <= 0 {
		return errors.New("页数为负.")
	}
	if req.PageSize <= 0 {
		return errors.New("每页条目为负.")
	}
	if req.PageSize > 100 {
		return errors.New("每页条目不能超过 100.")
	}
	return nil
}

// QueryOrderId implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) QueryOrderId(ctx context.Context, req *order.QueryOrderIdRequest) (resp *order.QueryOrderIdResponse, err error) {

	klogErr := func(msg string) {
		target := ""
		if req != nil {
			target = req.String()
		}
		klog.Error(
			"method: ", "QueryOrderId",
			"target: ", target,
			"message: ", msg,
		)
	}

	if err = validateQueryOrderIdReq(req); err != nil {
		klogErr("invalid params: " + err.Error())
		resp = &order.QueryOrderIdResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  err.Error(),
			},
		}
		return
	}

	var filter bson.M
	switch req.Type {
	case order.QueryOrderIdType_REQ_USER:
		filter = bson.M{"req_user_id": *req.UserId}
	case order.QueryOrderIdType_RESP_USER:
		filter = bson.M{"resp_user_id": *req.UserId}
	case order.QueryOrderIdType_EXT_KEY:
		filter = bson.M{*req.ExtKey: *req.ExtVal}
	}

	skip := (req.Page - 1) * req.PageSize

	findOpts := options.Find().
		SetProjection(bson.M{"_id": 1}).
		SetLimit(int64(req.PageSize)).
		SetSkip(int64(skip))

	cur, err := Coll.Find(ctx, filter, findOpts)
	if err != nil {
		klogErr("fail to query order: " + err.Error())
		resp = &order.QueryOrderIdResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  internalErrMsg,
			},
		}
		return
	}
	defer func() {
		if closeErr := cur.Close(ctx); closeErr != nil {
			klogErr("fail to close cursor: " + closeErr.Error())
		}
	}()

	var orderIds []string
	for cur.Next(ctx) {
		var tmpDoc struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err = cur.Decode(&tmpDoc); err != nil {
			klogErr("fail to decode order doc: " + err.Error())
			resp = &order.QueryOrderIdResponse{
				BaseResp: &base.BaseResponse{
					Code: base.Code_DB_ERR,
					Msg:  internalErrMsg,
				},
			}
			return
		}
		orderIds = append(orderIds, tmpDoc.ID.Hex())
	}

	if cur.Err() != nil {
		klogErr("cursor iteration error: " + cur.Err().Error())
		resp = &order.QueryOrderIdResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  internalErrMsg,
			},
		}
		return
	}

	total, err := Coll.CountDocuments(ctx, filter)
	if err != nil {
		klogErr("fail to count order total: " + err.Error())
		resp = &order.QueryOrderIdResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  internalErrMsg,
			},
		}
		return resp, err
	}

	resp = &order.QueryOrderIdResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  successMsg,
		},
		OrderId:  orderIds,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if len(orderIds) == 0 {
		resp.BaseResp.Msg = "当前页无数据."
	}

	return
}
