package logic

import (
	"context"

	"game_assistantor/debug/go_zero/protoc_mode/internal/svc"
	"game_assistantor/debug/go_zero/protoc_mode/types/goods"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsLogic {
	return &GetGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// rpc方法
func (l *GetGoodsLogic) GetGoods(in *goods.GoodsRequest) (*goods.GoodsResponse, error) {
	// todo: add your logic here and delete this line

	goodsId := in.GoodsId
	goods := new(goods.GoodsResponse)
	goods.GoodsId = goodsId
	goods.Name = "茅台"

	return goods, nil
}
