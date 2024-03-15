package logic

import (
	"context"

	"game_assistantor/debug/go_zero/api_mode/internal/svc"
	"game_assistantor/debug/go_zero/api_mode/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderInfoLogic {
	return &OrderInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderInfoLogic) OrderInfo(req *types.OrderInfoReq) (resp *types.OrderInfoResp, err error) {
	// todo: add your logic here and delete this line
	order_id := req.OrderId
	resp = new(types.OrderInfoResp)
	resp.GoodsName = "雪茄"
	resp.OrderId = order_id
	return
}
