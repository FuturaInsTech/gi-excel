package middleware

import (
	"strconv"

	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/gin-gonic/gin"
)

func SearchPagination(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	var message types.SearchPagination

	message.SearchString = ""
	message.SearchCriteria = ""
	message.PageNum = 1
	message.PageSize = 5
	message.SortColumn = ""
	message.SortDirection = ""
	message.FirstTime = true
	if queryParams.Has("searchString") {
		message.SearchString = queryParams.Get("searchString")
	}

	if queryParams.Has("searchCriteria") {
		message.SearchCriteria = queryParams.Get("searchCriteria")
	}

	if queryParams.Has("firstTime") {
		message.FirstTime, _ = strconv.ParseBool(queryParams.Get("firstTime"))

	}

	if queryParams.Has("pageNum") {
		message.PageNum, _ = strconv.Atoi(queryParams.Get("pageNum"))

	}

	if queryParams.Has("pageSize") {
		message.PageSize, _ = strconv.Atoi(queryParams.Get("pageSize"))
	}

	if queryParams.Has("sortColumn") {
		message.SortColumn = queryParams.Get("sortColumn")
	}

	if queryParams.Has("sortDirection") {
		message.SortDirection = queryParams.Get("sortDirection")
	}

	message.Offset = (message.PageNum - 1) * message.PageSize

	c.Set("searchpagination", message)

	c.Next()

}
