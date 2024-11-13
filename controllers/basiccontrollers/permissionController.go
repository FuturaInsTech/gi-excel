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
func GetAllPermission(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetAllPermission" //B0034
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

	var getallpermission []models.Permission
	//userco := userdatamap["CompanyId"]

	var result *gorm.DB

	if searchpagination.SearchString != "" && searchpagination.SearchCriteria != "" {
		result = initializers.DB.Model(&models.Permission{}).Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.Permission{}).
			Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallpermission)

	} else {
		fmt.Println("No Selection ")
		fmt.Println(searchpagination.SearchCriteria)
		fmt.Println(searchpagination.SearchString)
		result = initializers.DB.Model(&models.Permission{}).Where("company_id = ?", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.Permission{}).
			Where("company_id = ?", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallpermission)
	}

	// if result is null, then give an language ..
	if result.Error != nil {
		shortCode := "GL098"
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
			"displayName": "Model Name",
			"fieldName":   "model_name",
			"dataType":    "string"},
			{"displayName": "Method Name",
				"fieldName": "method",
				"dataType":  "string"},
		}

		c.JSON(200, gin.H{

			"All Permissions": getallpermission,
			"Field Map":       fieldMappings,
			"paginationData":  paginationData,
		})

	} else {
		c.JSON(200, gin.H{

			"All Permissions": getallpermission,
			"paginationData":  paginationData,
		})
	}

}

func GetAllPermissionByUserID(c *gin.Context) {
	// Here User Id, User Group and Method of the function passed to GetUserAccess function
	/*user, _ := c.Get("userid")
	method := "GetAllPermissionByUserID"  //B0035
	fmt.Println(user)
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)

	if err != nil {
		fmt.Println("Permissions")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Access Failed " + err.Permission(),
		})

		return
	}
	fmt.Println(userdatamap)
	//Get All in an array
	*/
	igetuserid := c.Param("userid")
	fmt.Println("User ID Passed" + igetuserid)
	var getallpermission []models.Permission
	user, _ := c.Get("user")
	method := "GetAllPermissionByUserID" //B0034
	//var userdatamap map[string]interface{}
	//_, err := utilities.GetUserAccess(user, method)
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))
	// run the query and save values into results.
	// similar to select * from <Table Name> and save it to results

	result := initializers.DB.Where("user_id LIKE ?", "%"+igetuserid+"%").Find(&getallpermission)
	// if result is null, then give an permission ..
	if result.Error != nil {
		shortCode := "GL098"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})
		// skip the funciton
		return

	}

	// return the values to Postman in JSON format
	c.JSON(200, gin.H{
		"All Permission": getallpermission,
	})

}

func GetAllPermissionByUserGroupID(c *gin.Context) {
	// Here User Id, User Group and Method of the function passed to GetUserAccess function
	/*user, _ := c.Get("userid")
	method := "GetAllPermissionByUserGroupID" //B0036
	fmt.Println(user)
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)

	if err != nil {
		fmt.Println("Permissions")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Access Failed " + err.Permission(),
		})

		return
	}
	fmt.Println(userdatamap)
	//Get All in an array
	*/
	igetusergroupid := c.Param("userid")
	fmt.Println("User ID Passed" + igetusergroupid)
	var getallpermission []models.Permission
	user, _ := c.Get("user")
	method := "GetAllPermissionByUserGroupID" //B0034
	//var userdatamap map[string]interface{}
	//_, err := utilities.GetUserAccess(user, method)
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))
	// run the query and save values into results.
	// similar to select * from <Table Name> and save it to results

	result := initializers.DB.Where("user_group_id LIKE ?", "%"+igetusergroupid+"%").Find(&getallpermission)
	// if result is null, then give an permission ..
	if result.Error != nil {
		shortCode := "GL098"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})
		// skip the funciton
		return

	}

	// return the values to Postman in JSON format
	c.JSON(200, gin.H{
		"All Permission": getallpermission,
	})

}

// Create Function

func CreatePermission(c *gin.Context) {
	// store it in working storage variable
	// field description in cobol
	user, _ := c.Get("user")
	method := "CreatePermission" //B0037
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
	var createpermission models.Permission
	fmt.Println(createpermission)
	err = c.Bind(&createpermission)

	//if c.Bind(&createpermission) != nil {
	if err != nil {
		shortCode := "GL098"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return
	}

	// Validations
	// Validate Inputs Fields Before Updating into DB
	jsonStr, err := json.Marshal(createpermission)
	if err != nil {
		fmt.Println(err)
	}
	var datamap map[string]interface{}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &datamap); err != nil {
		fmt.Println(err)
	}
	err, dbvalerror := utilities.ValidateData(datamap, "Permission")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error ": dbvalerror,
		})
		return
	}

	// Validations Ended
	result := initializers.DB.Create(&createpermission)
	fmt.Println(createpermission)
	if result.Error != nil {
		shortCode := "GL100"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	//response
	shortCode := "GL358"
	longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"error": shortCode + " : " + longDesc,
	})

}

//Delete Function

func DeletePermission(c *gin.Context) {

	user, _ := c.Get("user")
	method := "DeletePermission" //B0038
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

	var delpermission models.Permission
	result := initializers.DB.First(&delpermission, "id  = ?", delid)

	fmt.Println(delpermission)
	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	result = initializers.DB.Delete(&delpermission)

	if result.Error != nil {
		shortCode := "GL038"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, "Permission ID "+delid+" is deleted")

}

// Get Function
func GetPermission(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetPermission" //B0039
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

	var getpermission models.Permission
	result := initializers.DB.First(&getpermission, "id  = ?", getid)
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
		"Permission": getpermission,
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
	if result.Permission != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get :" + result.Permission.Permission(),
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
	if result.Permission != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Clone :" + result.Permission.Permission(),
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
func ClonePermission(c *gin.Context) {
	user, _ := c.Get("user")
	method := "ClonePermission" //B0040
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

	var spermission models.Permission
	result := initializers.DB.First(&spermission, "id  = ?", sourceid)
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
	data, _ := json.Marshal(spermission)
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
	var tpermission models.Permission
	// converting json to a model
	json.Unmarshal(data, &tpermission)
	// edecuting query persisting the model
	result = initializers.DB.Create(&tpermission)
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
		"Cloned": tpermission,
	})

}

// Modify Function
func ModifyPermission(c *gin.Context) {
	// mapping json to sourceMap
	user, _ := c.Get("user")
	method := "ModifyPermission" //B0041
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
		shortCode := "GL099"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	var opermission models.Permission
	result := initializers.DB.First(&opermission, "id  = ?", sourceMap["ID"])
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
	data, _ := json.Marshal(opermission)

	json.Unmarshal(data, &targetMap)

	for key, _ := range targetMap {

		if val1, ok := sourceMap[key]; ok {
			targetMap[key] = val1
		}

	}

	data, _ = json.Marshal(targetMap)
	json.Unmarshal(data, &opermission)
	// Validations
	// Validate Inputs Fields Before Updating into DB
	jsonStr, err := json.Marshal(opermission)
	if err != nil {
		fmt.Println(err)
	}
	var datamap map[string]interface{}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &datamap); err != nil {
		fmt.Println(err)
	}
	err, dbvalerror := utilities.ValidateData(datamap, "Permission")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error ": dbvalerror,
		})
		return
	}

	// Validations Ended

	// update modified time
	//opermission.UpdatedAt = time.Now()
	result = initializers.DB.Save(&opermission)

	if result.Error != nil {
		shortCode := "GL041"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(200, gin.H{
		"Permission": opermission,
	})

}
