package handler

import (
	"net/http"

	"game_assistantor/debug/go_zero/api_mode/internal/logic"
	"game_assistantor/debug/go_zero/api_mode/internal/svc"
	"game_assistantor/debug/go_zero/api_mode/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func orderInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOrderInfoLogic(r.Context(), svcCtx)
		resp, err := l.OrderInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
