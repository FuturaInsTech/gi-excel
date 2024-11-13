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
// This function Name we need to add it in main.go
func GetAllError(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetAllErrors" //B0028
	//var userdatamap map[string]interface{}
	//_, err := utilities.GetUserAccess(user, method)
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0116" // Access failed
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

	var getallerror []models.Error
	//userco := userdatamap["CompanyId"]

	var result *gorm.DB

	if searchpagination.SearchString != "" && searchpagination.SearchCriteria != "" {
		result = initializers.DB.Model(&models.Error{}).Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.Error{}).
			Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallerror)

	} else {
		fmt.Println("No Selection ")
		fmt.Println(searchpagination.SearchCriteria)
		fmt.Println(searchpagination.SearchString)
		result = initializers.DB.Model(&models.Error{}).Where("company_id = ?", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.Error{}).
			Where("company_id = ?", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallerror)
	}

	// if result is null, then give an language ..
	if result.Error != nil {
		shortCode := "K0117" // Failed to fetch error
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
			"displayName": "Short Code",
			"fieldName":   "short_code",
			"dataType":    "string"},
			{"displayName": "Long Code",
				"fieldName": "long_code",
				"dataType":  "string"},
		}

		c.JSON(200, gin.H{

			"All Errors":     getallerror,
			"Field Map":      fieldMappings,
			"paginationData": paginationData,
		})

	} else {
		c.JSON(200, gin.H{

			"All Errors":     getallerror,
			"paginationData": paginationData,
		})
	}

}

// Create Function

func CreateError(c *gin.Context) {
	// store it in working storage variable
	// field description in cobol
	user, _ := c.Get("user")
	method := "CreateError" //B0029
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0116"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var createerror models.Error

	err1 := c.Bind(&createerror)
	//if c.Bind(&createerror) != nil {
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err1.Error(),
		})

		return
	}
	// Validations
	// Validate Inputs Fields Before Updating into DB
	jsonStr, err := json.Marshal(createerror)
	if err != nil {
		fmt.Println(err)
	}
	var datamap map[string]interface{}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &datamap); err != nil {
		fmt.Println(err)
	}
	err, dbvalerror := utilities.ValidateData(datamap, "Error")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error ": dbvalerror,
		})
		return
	}

	// Validations Ended
	result := initializers.DB.Create(&createerror)

	if result.Error != nil {
		shortCode := "K0118" // failed to create error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	//response
	shortCode := "K0119" // successfully created error
	longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"Result": shortCode + " : " + longDesc,
	})

}

//Delete Function

func DeleteError(c *gin.Context) {

	user, _ := c.Get("user")
	method := "DeleteError" //B0030
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0116"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	delid := c.Param("id")

	var delerror models.Error
	result := initializers.DB.First(&delerror, "id  = ?", delid)
	fmt.Println(delerror)
	if result.Error != nil {
		shortCode := "K0120" // failed to get error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	result = initializers.DB.Delete(&delerror)

	if result.Error != nil {
		shortCode := "K0121" // failed to delete error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, "Error ID "+delid+" is deleted")

}

// Get Function
func GetError(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetError" //B0031

	//var userdatamap map[string]interface{}
	iSkipPermission := "Y"
	userdatamap, err := utilities.GetUserAccessNew(user, method, iSkipPermission)
	//userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0116"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println("User Data map")
	fmt.Println(userdatamap)
	getid := c.Param("id")

	//var out1 = utilities.sumofvariables(1, 2)

	fmt.Println("rangarajan !!!!!")

	var sourceMap map[string]interface{}
	// converting an entity(Model) to Json
	data1, _ := json.Marshal(user)
	//converting Json to Source Map
	json.Unmarshal(data1, &sourceMap)

	icompany := userdatamap["CompanyId"]
	fmt.Println(sourceMap)
	fmt.Println(icompany)

	var geterror models.Error
	result := initializers.DB.First(&geterror, "id  = ?", getid)
	if result.Error != nil {
		shortCode := "K0120" // failed to get error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"Error": geterror,
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
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get :" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	var tco models.Company
	tco.CompanyAddress1 = sco.CompanyAddress1
	tco.CompanyAddress2 = sco.CompanyAddress2
	tco.CompanyAddress3 = sco.CompanyAddress3
	tco.CompanyAddress4 = sco.CompanyAddress4
	tco.CompanyAddress5 = sco.CompanyAddress5
	tco.CompanyGst = sco.CompanyGst
	tco.CompanyIncorporationDate = sco.CompanyIncorporationDate
	tco.CompanyName = sco.CompanyName
	tco.CompanyUid = sco.CompanyUid
	tco.CreatedAt = time.Now()
	tco.IsActive = sco.IsActive
	tco.OwnerId = sco.OwnerId

	result = initializers.DB.Create(&tco)
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Clone :" + result.Error.Error(),
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
func CloneError(c *gin.Context) {
	user, _ := c.Get("user")
	method := "CloneError" //B0032
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0116"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	sourceid := c.Param("id")
	//fmt.Println("ID " + sid)

	var serror models.Error
	result := initializers.DB.First(&serror, "id  = ?", sourceid)
	//fmt.Println((sco.CompanyName))
	if result.Error != nil {
		shortCode := "K0120" // failed to get error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// Declaring a Map so that i could move all or selective values into my Map
	var sourceMap map[string]interface{}
	// converting an entity(Model) to Json
	data, _ := json.Marshal(serror)
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
	var terror models.Error
	// converting json to a model
	json.Unmarshal(data, &terror)
	// edecuting query persisting the model
	result = initializers.DB.Create(&terror)
	if result.Error != nil {
		shortCode := "K0122" // failed to create error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"Cloned": terror,
	})

}

// Modify Function
func ModifyError(c *gin.Context) {
	// mapping json to sourceMap
	user, _ := c.Get("user")
	method := "ModifyError" //B0033
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0115"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var sourceMap map[string]interface{}

	if c.Bind(&sourceMap) != nil {
		shortCode := "K0123" // failed to fetch error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	var oerror models.Error

	result := initializers.DB.First(&oerror, "id  = ?", sourceMap["ID"])
	if result.Error != nil {
		shortCode := "K0120" // failed to get error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	var targetMap map[string]interface{}
	fmt.Println((targetMap))
	data, _ := json.Marshal(oerror)

	json.Unmarshal(data, &targetMap)

	for key, _ := range targetMap {

		if val1, ok := sourceMap[key]; ok {
			targetMap[key] = val1
		}

	}

	data, _ = json.Marshal(targetMap)
	json.Unmarshal(data, &oerror)
	// update modified time
	//oerror.UpdatedAt = time.Now()
	//oerror.UpdatedID := iid
	fmt.Println("MOdified User")
	updateduserid := userdatamap["Id"]
	fmt.Println(updateduserid)
	//oerror.UpdatedID = updateduserid

	// Validations
	// Validate Inputs Fields Before Updating into DB
	jsonStr, err := json.Marshal(oerror)
	if err != nil {
		fmt.Println(err)
	}
	var datamap map[string]interface{}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &datamap); err != nil {
		fmt.Println(err)
	}
	err, dbvalerror := utilities.ValidateData(datamap, "Error")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error ": dbvalerror,
		})
		return
	}

	result = initializers.DB.Save(&oerror)

	if result.Error != nil {
		shortCode := "K0124" // failed to save error
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusOK, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{"outputs": oerror})

}
