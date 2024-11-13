package basiccontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/FuturaInsTech/GoExcel/utilities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get All Function

func GetAllTransaction(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetAllTransaction" //B0042
	//var userdatamap map[string]interface{}
	//_, err := utilities.GetUserAccess(user, method)
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "GL001"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	// Filter Variables
	// search and pagination
	var searchpagination types.SearchPagination

	temp, _ := c.Get("searchpagination")
	searchpagination, ok := temp.(types.SearchPagination)
	fmt.Println("OK Value")
	fmt.Println(ok)

	if searchpagination.SortColumn == "" {
		searchpagination.SortColumn = "id"
		searchpagination.SortDirection = "asc"
	}

	fmt.Println(ok)
	var totalRecords int64 = 0

	var getallTransaction []models.Transaction
	//userco := userdatamap["CompanyId"]

	var result *gorm.DB

	if searchpagination.SearchString != "" && searchpagination.SearchCriteria != "" {
		result = initializers.DB.Model(&models.Transaction{}).Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.Transaction{}).
			Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallTransaction)

	} else {
		fmt.Println("No Selection ")
		fmt.Println(searchpagination.SearchCriteria)
		fmt.Println(searchpagination.SearchString)
		result = initializers.DB.Model(&models.Transaction{}).Where("company_id = ?", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.Transaction{}).
			Where("company_id = ?", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallTransaction)
	}

	// if result is null, then give an language ..
	if result.Error != nil {
		shortCode := "GL042"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})
		// skip the funciton
		return

	}
	paginationData := map[string]interface{}{
		"totalRecords": totalRecords,
	}
	// return the values to Postman in JSON format
	// Provide Search Fields... currently 2 fields are used.

	if searchpagination.FirstTime {
		fieldMappings := [2]map[string]string{{
			"displayName": "Transaction",
			"fieldName":   "method",
			"dataType":    "string"},
			{"displayName": "Description",
				"fieldName": "description",
				"dataType":  "string"},
		}

		c.JSON(200, gin.H{

			"All Transactions": getallTransaction,
			"Field Map":        fieldMappings,
			"paginationData":   paginationData,
		})

	} else {
		c.JSON(200, gin.H{

			"All Transactions": getallTransaction,
			"paginationData":   paginationData,
		})
	}

}

func GetAllTransactionByClientID(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetAllTransactionByClientID" //B0043
	//var userdatamap map[string]interface{}
	//_, err := utilities.GetUserAccess(user, method)
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "GL001"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	// Filter Variables
	// search and pagination
	var searchpagination types.SearchPagination
	getclientid := c.Param("id")

	temp, _ := c.Get("searchpagination")
	searchpagination, ok := temp.(types.SearchPagination)
	fmt.Println("OK Value")
	fmt.Println(ok)

	if searchpagination.SortColumn == "" {
		searchpagination.SortColumn = "id"
		searchpagination.SortDirection = "asc"
	}

	fmt.Println(ok)
	var totalRecords int64 = 0

	var getallTransaction []models.Transaction
	//userco := userdatamap["CompanyId"]

	var result *gorm.DB

	if searchpagination.SearchString != "" && searchpagination.SearchCriteria != "" {
		result = initializers.DB.Model(&models.Transaction{}).Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ? and client_id = ?", "%"+searchpagination.SearchString+"%", userco, getclientid).Count(&totalRecords)
		result = initializers.DB.Model(&models.Transaction{}).
			Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ? and client_id = ?", "%"+searchpagination.SearchString+"%", userco, getclientid).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallTransaction)

	} else {
		fmt.Println("No Selection ")
		fmt.Println(searchpagination.SearchCriteria)
		fmt.Println(searchpagination.SearchString)
		result = initializers.DB.Model(&models.Transaction{}).Where("company_id = ? and client_id = ?", userco, getclientid).Count(&totalRecords)
		result = initializers.DB.Model(&models.Transaction{}).
			Where("company_id = ? and client_id = ?", userco, getclientid).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallTransaction)
	}

	// if result is null, then give an language ..
	if result.Error != nil {
		shortCode := "GL042"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})
		// skip the funciton
		return

	}
	paginationData := map[string]interface{}{
		"totalRecords": totalRecords,
	}
	// return the values to Postman in JSON format
	// Provide Search Fields... currently 2 fields are used.

	if searchpagination.FirstTime {
		fieldMappings := [6]map[string]string{{
			"displayName": "Transaction Line 1",
			"fieldName":   "Transaction_line1",
			"dataType":    "string"},
			{"displayName": "Transaction Line 2",
				"fieldName": "Transaction_line2",
				"dataType":  "string"},
			{"displayName": "Transaction Line 3",
				"fieldName": "Transaction_line3",
				"dataType":  "string"},
			{"displayName": "Transaction Line 4",
				"fieldName": "Transaction_line4",
				"dataType":  "string"},
			{"displayName": "Transaction Line 5",
				"fieldName": "Transaction_line5",
				"dataType":  "string"},
			{"displayName": "Postal Code",
				"fieldName": "Transaction_postal_code",
				"dataType":  "string"}}

		c.JSON(200, gin.H{

			"TransactionByClientID": getallTransaction,
			"Field Map":             fieldMappings,
			"paginationData":        paginationData,
		})

	} else {
		c.JSON(200, gin.H{

			"TransactionByClientID": getallTransaction,
			"paginationData":        paginationData,
		})
	}

}

// Create Function

func CreateTransaction(c *gin.Context) {
	// store it in working storage variable
	// field description in cobol
	user, _ := c.Get("user")
	method := "CreateTransaction" //B0044
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "GL001"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var createTransaction models.Transaction

	err1 := c.Bind(&createTransaction)
	//if c.Bind(&createTransaction) != nil {
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err1.Error(),
		})

		return
	}

	result := initializers.DB.Create(&createTransaction)

	if result.Error != nil {
		shortCode := "GL110"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	//response
	shortCode := "GL361"
	longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"Result": shortCode + " : " + longDesc,
	})

}

//Delete Function

func DeleteTransaction(c *gin.Context) {

	user, _ := c.Get("user")
	method := "DeleteTransaction" //B0045
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "GL001"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	delid := c.Param("id")

	var delTransaction models.Transaction
	result := initializers.DB.First(&delTransaction, "id  = ?", delid)
	fmt.Println(delTransaction)
	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	result = initializers.DB.Delete(&delTransaction)

	if result.Error != nil {
		shortCode := "GL038"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, "Transaction ID "+delid+" is deleted")

}

// Get Function
func GetTransaction(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetTransaction" //B0046

	//var userdatamap map[string]interface{}
	iSkipPermission := "Y"
	userdatamap, err := utilities.GetUserAccessNew(user, method, iSkipPermission)
	//userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "GL001"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	getid := c.Param("id")

	//var out1 = utilities.sumofvariables(1, 2)

	fmt.Println("rangarajan ")

	var getTransaction models.Transaction
	result := initializers.DB.First(&getTransaction, "id  = ?", getid)
	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"Transaction": getTransaction,
	})

}

/*
// clone only selective fields  Method 1  Field by Field
func CloneCompany(c *gin.Context) {
	sid := c.Param("id")
	//fmt.Println("ID " + sid)

	var sco models.Company
	result := initializers.DB.First(&sco, "id  = ?", sid)
	//fmt.Println((sco.CompanyName))
	if result.Transaction != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get :" + result.Transaction.Transaction(),
		})

		return

	}
	// return the values to postman
	var tco models.Company
	tco.CompanyTransaction1 = sco.CompanyTransaction1
	tco.CompanyTransaction2 = sco.CompanyTransaction2
	tco.CompanyTransaction3 = sco.CompanyTransaction3
	tco.CompanyTransaction4 = sco.CompanyTransaction4
	tco.CompanyTransaction5 = sco.CompanyTransaction5
	tco.CompanyGst = sco.CompanyGst
	tco.CompanyIncorporationDate = sco.CompanyIncorporationDate
	tco.CompanyName = sco.CompanyName
	tco.CompanyUid = sco.CompanyUid
	tco.CreatedAt = time.Now()
	tco.IsActive = sco.IsActive
	tco.OwnerId = sco.OwnerId

	result = initializers.DB.Create(&tco)
	if result.Transaction != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Clone :" + result.Transaction.Transaction(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"Company": tco,
	})

}
*/
// Clone all fields
func CloneTransaction(c *gin.Context) {
	user, _ := c.Get("user")
	method := "CloneTransaction" //B0047
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "GL001"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	sourceid := c.Param("id")
	//fmt.Println("ID " + sid)

	var sTransaction models.Transaction
	result := initializers.DB.First(&sTransaction, "id  = ?", sourceid)
	//fmt.Println((sco.CompanyName))
	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// Declaring a Map so that i could move all or selective values into my Map
	var sourceMap map[string]interface{}
	// converting an entity(Model) to Json
	data, _ := json.Marshal(sTransaction)
	//converting Json to Source Map
	json.Unmarshal(data, &sourceMap)

	var targetMap = make(map[string]interface{})

	// moving all values except ID
	for key, val := range sourceMap {

		if key != "ID" {
			targetMap[key] = val
		}

	}
	// converting target map to a json
	data, _ = json.Marshal(targetMap)
	// creating a local model
	var tTransaction models.Transaction
	// converting json to a model
	json.Unmarshal(data, &tTransaction)
	// edecuting query persisting the model
	result = initializers.DB.Create(&tTransaction)
	if result.Error != nil {
		shortCode := "GL039"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"Cloned": tTransaction,
	})

}

// Modify Function
func ModifyTransaction(c *gin.Context) {
	// mapping json to sourceMap
	user, _ := c.Get("user")
	method := "ModifyTransaction" //B0048
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "GL001"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var sourceMap map[string]interface{}

	if c.Bind(&sourceMap) != nil {
		shortCode := "GL111"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	var oTransaction models.Transaction

	result := initializers.DB.First(&oTransaction, "id  = ?", sourceMap["ID"])
	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	var targetMap map[string]interface{}
	fmt.Println((targetMap))
	data, _ := json.Marshal(oTransaction)

	json.Unmarshal(data, &targetMap)

	for key, _ := range targetMap {

		if val1, ok := sourceMap[key]; ok {
			targetMap[key] = val1
		}

	}

	data, _ = json.Marshal(targetMap)
	json.Unmarshal(data, &oTransaction)
	// update modified time
	//oTransaction.UpdatedAt = time.Now()
	//oTransaction.UpdatedID := iid
	fmt.Println("MOdified User")
	updateduserid := userdatamap["Id"]
	fmt.Println(updateduserid)
	//oTransaction.UpdatedID = updateduserid

	result = initializers.DB.Save(&oTransaction)

	if result.Error != nil {
		shortCode := "GL041"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusOK, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{"outputs": oTransaction})

}
