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
func GetAllUserGroup(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetAllUserGroups" //B0049
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

	var getallusergroup []models.UserGroup
	//userco := userdatamap["CompanyId"]

	var result *gorm.DB

	if searchpagination.SearchString != "" && searchpagination.SearchCriteria != "" {
		result = initializers.DB.Model(&models.UserGroup{}).Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.UserGroup{}).
			Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallusergroup)

	} else {
		fmt.Println("No Selection ")
		fmt.Println(searchpagination.SearchCriteria)
		fmt.Println(searchpagination.SearchString)
		result = initializers.DB.Model(&models.UserGroup{}).Where("company_id = ?", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.UserGroup{}).
			Where("company_id = ?", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallusergroup)
	}

	// if result is null, then give an language ..
	if result.Error != nil {
		shortCode := "GL124"
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
		fieldMappings := [1]map[string]string{{
			"displayName": "Group Name",
			"fieldName":   "group_name",
			"dataType":    "string"},
		}

		c.JSON(200, gin.H{

			"All UserGroups": getallusergroup,
			"Field Map":      fieldMappings,
			"paginationData": paginationData,
		})

	} else {
		c.JSON(200, gin.H{

			"All UserGroups": getallusergroup,
			"paginationData": paginationData,
		})
	}

}

// Create Function

func CreateUserGroup(c *gin.Context) {
	// store it in working storage variable
	// field description in cobol
	user, _ := c.Get("user")
	method := "CreateUserGroup" //B0050
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
	var createusergroup models.UserGroup

	err1 := c.Bind(&createusergroup)
	//if c.Bind(&createusergroup) != nil {
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err1.Error(),
		})

		return
	}
	// Date Compare
	errorlongcode := utilities.CompareDate(createusergroup.ValidFrom, createusergroup.ValidTo, uint(userdatamap["LanguageId"].(float64)))

	if errorlongcode != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorlongcode.Error(),
		})

		return
	}
	// Date Should Not Be Zero

	errorlongcode = utilities.DateZero(createusergroup.ValidFrom, uint(userdatamap["LanguageId"].(float64)))

	if errorlongcode != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorlongcode.Error(),
		})

		return
	}
	result := initializers.DB.Create(&createusergroup)

	if result.Error != nil {
		shortCode := "GL125"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	//response
	shortCode := "GL362"
	longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"Result": shortCode + " : " + longDesc,
	})

}

//Delete Function

func DeleteUserGroup(c *gin.Context) {

	user, _ := c.Get("user")
	method := "DeleteUserGroup" //B0051
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

	var delusergroup models.UserGroup
	result := initializers.DB.First(&delusergroup, "id  = ?", delid)
	fmt.Println(delusergroup)
	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	result = initializers.DB.Delete(&delusergroup)

	if result.Error != nil {
		shortCode := "GL038"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, "UserGroup ID "+delid+" is deleted")

}

// Get Function
func GetUserGroup(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetUserGroup" //B0052

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

	var getusergroup models.UserGroup
	result := initializers.DB.First(&getusergroup, "id  = ?", getid)
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
		"UserGroup": getusergroup,
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
	tco.CompanyUserGroup1 = sco.CompanyUserGroup1
	tco.CompanyUserGroup2 = sco.CompanyUserGroup2
	tco.CompanyUserGroup3 = sco.CompanyUserGroup3
	tco.CompanyUserGroup4 = sco.CompanyUserGroup4
	tco.CompanyUserGroup5 = sco.CompanyUserGroup5
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
func CloneUserGroup(c *gin.Context) {
	user, _ := c.Get("user")
	method := "CloneUserGroup" //B0053
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

	var susergroup models.UserGroup
	result := initializers.DB.First(&susergroup, "id  = ?", sourceid)
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
	data, _ := json.Marshal(susergroup)
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
	var tusergroup models.UserGroup
	// converting json to a model
	json.Unmarshal(data, &tusergroup)
	// edecuting query persisting the model
	result = initializers.DB.Create(&tusergroup)
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
		"Cloned": tusergroup,
	})

}

// Modify Function
func ModifyUserGroup(c *gin.Context) {
	// mapping json to sourceMap
	user, _ := c.Get("user")
	method := "ModifyUserGroup" //B0054
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
		shortCode := "GL126"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	var ousergroup models.UserGroup

	result := initializers.DB.First(&ousergroup, "id  = ?", sourceMap["ID"])
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
	data, _ := json.Marshal(ousergroup)

	json.Unmarshal(data, &targetMap)

	for key, _ := range targetMap {

		if val1, ok := sourceMap[key]; ok {
			targetMap[key] = val1
		}

	}

	data, _ = json.Marshal(targetMap)
	json.Unmarshal(data, &ousergroup)
	// update modified time
	//ousergroup.UpdatedAt = time.Now()
	//ousergroup.UpdatedID := iid
	fmt.Println("MOdified User")
	updateduserid := userdatamap["Id"]
	fmt.Println(updateduserid)
	//ousergroup.UpdatedID = updateduserid

	result = initializers.DB.Save(&ousergroup)

	if result.Error != nil {
		shortCode := "GL041"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusOK, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{"outputs": ousergroup})

}
