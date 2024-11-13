package routes

import (
	"github.com/FuturaInsTech/GoExcel/excel"
	"github.com/FuturaInsTech/GoExcel/middleware"
	"github.com/gin-gonic/gin"
)

func ExcelRoutes(route *gin.RouterGroup) {
	services := route.Group("/excelservices", middleware.RequiredAuth)
	{
		services.POST("/:serviceName", excel.ProcessExcel1)

	}

}
