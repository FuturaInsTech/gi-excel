package routes

import (
	"github.com/FuturaInsTech/GoExcel/controllers/basiccontrollers"
	"github.com/FuturaInsTech/GoExcel/excel"
	"github.com/FuturaInsTech/GoExcel/middleware"
	"github.com/gin-gonic/gin"
)

func Basicroutes(route *gin.RouterGroup) {
	services := route.Group("/basicservices", middleware.RequiredAuth)
	{
		services.GET("/getjson", excel.ProcessExcel)
		//Param
		services.POST("/param", basiccontrollers.CreateParam)
		services.PUT("/param", basiccontrollers.ModifyParam)
		services.GET("/param", basiccontrollers.EnquireParam)
		services.DELETE("/param", basiccontrollers.DeleteParam)
		services.GET("/paramItem", basiccontrollers.EnquireParamItem)
		services.GET("/paramItems", basiccontrollers.EnquireParamItems)
		services.POST("/paramItem", basiccontrollers.CreateParamItem)
		services.POST("/cloneParamItem", basiccontrollers.CloneParamItem)
		services.PUT("/paramItem", basiccontrollers.ModifyParamItem)
		services.DELETE("/paramItem", basiccontrollers.DeleteParamItem)
		services.GET("/params", basiccontrollers.EnquireParams)
		services.GET("/paramextradata", basiccontrollers.GetParamExtraData)
		services.POST("/paramDataUpload", basiccontrollers.UploadParamData)
		services.POST("/autoperm/:id", basiccontrollers.AddPermissionAuto)

	}

}
