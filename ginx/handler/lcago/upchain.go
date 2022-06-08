package lcago

import (
	"lcago/ginx/model/lcago"
	"lcago/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Mint(ctx *gin.Context) {
	// 1.req
	req := new(lcago.MintRequest)
	resp := new(lcago.MintResponse)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		log.Error(err)
		ctx.Abort()
		resp.Msg = "fialed"
		resp.Code = -1
		resp.Data = lcago.Data{
			Key:   "key",
			Value: "hello",
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.Msg = "success"
	resp.Code = 0
	resp.Data = lcago.Data{
		Key:   "ok",
		Value: "world",
	}
	ctx.JSON(http.StatusOK, resp)
}
