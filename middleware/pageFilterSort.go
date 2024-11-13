package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/gin-gonic/gin"
)

func PageFilterSortHandler(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	pageFilterSort := types.PageFilterSort{Offset: 0, PageSize: 5, Filters: nil, SortingList: ""}

	filters := "[]"
	sorting := "[]"

	if queryParams.Has("start") {
		pageFilterSort.Offset, _ = strconv.Atoi(queryParams.Get("start"))

	}
	if queryParams.Has("size") {
		pageFilterSort.PageSize, _ = strconv.Atoi(queryParams.Get("size"))
	}

	if queryParams.Has("filters") {
		filters = queryParams.Get("filters")
	}

	if queryParams.Has("sorting") {
		sorting = queryParams.Get("sorting")
	}

	var sortings []types.Sorting

	err := json.Unmarshal([]byte(sorting), &sortings)

	//Checks whether the error is nil or not
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while decoding soring the data" + err.Error(),
		})

		return
	}

	err = json.Unmarshal([]byte(filters), &pageFilterSort.Filters)

	//Checks whether the error is nil or not
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while decoding filter data" + err.Error(),
		})

		return
	}

	sortingArray := make([]string, 0)
	for _, val := range sortings {
		// perform an operation

		sortDir := ""
		if val.Desc {
			sortDir = "desc"
		}

		sortingArray = append(sortingArray, val.Id+" "+sortDir)
	}
	pageFilterSort.SortingList = strings.Join(sortingArray, ",")

	c.Set("pageFilterSort", pageFilterSort)

	c.Next()

}

func PageFilterSortHandlerNew(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	pageFilterSort := types.PageFilterSort{Offset: 0, PageSize: 9999, Filters: nil, SortingList: ""}

	filters := "[]"
	sorting := "[]"

	if queryParams.Has("start") {
		pageFilterSort.Offset, _ = strconv.Atoi(queryParams.Get("start"))

	}
	if queryParams.Has("size") {
		pageFilterSort.PageSize, _ = strconv.Atoi(queryParams.Get("size"))
	}

	if queryParams.Has("filters") {
		filters = queryParams.Get("filters")
	}

	if queryParams.Has("sorting") {
		sorting = queryParams.Get("sorting")
	}

	var sortings []types.Sorting

	err := json.Unmarshal([]byte(sorting), &sortings)

	//Checks whether the error is nil or not
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while decoding soring the data" + err.Error(),
		})

		return
	}

	err = json.Unmarshal([]byte(filters), &pageFilterSort.Filters)

	//Checks whether the error is nil or not
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while decoding filter data" + err.Error(),
		})

		return
	}

	sortingArray := make([]string, 0)
	for _, val := range sortings {
		// perform an operation

		sortDir := ""
		if val.Desc {
			sortDir = "desc"
		}

		sortingArray = append(sortingArray, val.Id+" "+sortDir)
	}
	pageFilterSort.SortingList = strings.Join(sortingArray, ",")

	c.Set("pageFilterSort", pageFilterSort)

	c.Next()

}
