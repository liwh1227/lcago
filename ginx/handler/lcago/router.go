package lcago

import "github.com/gin-gonic/gin"

// Router 路由分组
func Router(r *gin.Engine) {
	router := r.Group("/lcago")
	{
		router.POST("upChain/mint", Mint)
	}
}
