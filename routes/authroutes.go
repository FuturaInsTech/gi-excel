package routes

import (
	"github.com/FuturaInsTech/GoExcel/controllers/basiccontrollers"
	"github.com/FuturaInsTech/GoExcel/middleware"
	"github.com/gin-gonic/gin"
)

func Authroutes(route *gin.RouterGroup) {
	auth := route.Group("/auth")
	{
		auth.POST("/signup", basiccontrollers.Signup)
		auth.GET("/getallusers", middleware.RequiredAuth, middleware.SearchPagination, basiccontrollers.GetAllUsers)
		auth.POST("/login", basiccontrollers.Login)
		auth.GET("/validateUser", middleware.RequiredAuth, basiccontrollers.Validate)
		auth.DELETE("/deleteuser/:rangaid", basiccontrollers.DeleteUser)
		auth.GET("/GetUser/:id", basiccontrollers.GetUser)
		auth.PUT("/user", basiccontrollers.ModifyUser)
		auth.POST("/logout", basiccontrollers.Logout)
		auth.POST("/refresh", middleware.CheckRefreshToken, basiccontrollers.Refresh)
		//	auth..GET("/user",  basiccontrollers.Validate)

	}

}
