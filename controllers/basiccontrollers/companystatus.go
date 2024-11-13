package basiccontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/utilities"
	"github.com/gin-gonic/gin"
)

// Get All Function
// This function Name we need to add it in main.go
func GetAllCompanyStatus(c *gin.Context) {
	user, _ := c.Get("user")
	method := "GetAllCompanyStatus" //B0055
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0109" // Access Failed
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	//Get All in an array
	var getallcompanystatus []models.CompanyStatus
	// run the query and save values into results.
	// similar to select * from <Table Name> and save it to results
	result := initializers.DB.Find(&getallcompanystatus)
	// if result is null, then give an error ..
	if result.Error != nil {
		shortCode := "K0110" //Failed to get company status
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})
		// skip the funciton
		return

	}

	// return the values to Postman in JSON format
	c.JSON(200, gin.H{
		"All Status": getallcompanystatus,
	})

}

// Create Function

func CreateCompanyStatus(c *gin.Context) {
	// store it in working storage variable
	// field description in cobol
	user, _ := c.Get("user")
	method := "CreateCompanyStatus"
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0109"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var createcompanystatus models.CompanyStatus

	if c.Bind(&createcompanystatus) != nil {
		shortCode := "K0111" // failed to fetch company status
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return
	}

	result := initializers.DB.Create(&createcompanystatus)

	if result.Error != nil {
		shortCode := "K0112" // Failed to create company status
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	//response
	shortCode := "K0113" // Successfully created company
	longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"Result": shortCode + " : " + longDesc,
	})

}

//Delete Function

func DeleteCompanyStatus(c *gin.Context) {

	user, _ := c.Get("user")
	method := "DeleteCompanyStatus"
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0109"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	delid := c.Param("id")

	var delcompanystatus models.CompanyStatus
	//	delcompanystatus.UpdatedID = uint64(userdatamap["Id"].(float64))
	result := initializers.DB.First(&delcompanystatus, "id  = ?", delid)
	fmt.Println(delcompanystatus)
	if result.Error != nil {
		shortCode := "K0110"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	result = initializers.DB.Delete(&delcompanystatus)

	if result.Error != nil {
		shortCode := "K0114" // failed to delete
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, "CompanyStatus ID "+delid+" is deleted")

}

// Get Function
func GetCompanyStatus(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetCompanyStatus"

	//var userdatamap map[string]interface{}
	iSkipPermission := "Y"
	userdatamap, err := utilities.GetUserAccessNew(user, method, iSkipPermission)
	//userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0109"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	getid := c.Param("id")

	var getcompanystatus models.CompanyStatus

	result := initializers.DB.First(&getcompanystatus, "id  = ?", getid)
	if result.Error != nil {
		shortCode := "K0110"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"CompanyStatus": getcompanystatus,
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
func CloneCompanyStatus(c *gin.Context) {
	user, _ := c.Get("user")
	method := "CloneCompanyStatus"
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0109"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	sourceid := c.Param("id")
	//fmt.Println("ID " + sid)

	var scompanystatus models.CompanyStatus
	result := initializers.DB.First(&scompanystatus, "id  = ?", sourceid)
	//fmt.Println((sco.CompanyName))
	if result.Error != nil {
		shortCode := "K0110"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// Declaring a Map so that i could move all or selective values into my Map
	var sourceMap map[string]interface{}
	// converting an entity(Model) to Json
	data, _ := json.Marshal(scompanystatus)
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
	var tcompanystatus models.CompanyStatus
	// converting json to a model
	json.Unmarshal(data, &tcompanystatus)
	// edecuting query persisting the model
	//tcompanystatus.UpdatedID = uint64(userdatamap["Id"].(float64))
	result = initializers.DB.Create(&tcompanystatus)
	if result.Error != nil {
		shortCode := "K0112"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"Cloned": tcompanystatus,
	})

}

// Modify Function
func ModifyCompanyStatus(c *gin.Context) {
	// mapping json to sourceMap
	user, _ := c.Get("user")
	method := "ModifyCompanyStatus"
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0109"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var sourceMap map[string]interface{}

	if c.Bind(&sourceMap) != nil {
		shortCode := "K0111" // failed to fetch company status
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	var ocompanystatus models.CompanyStatus
	//	ocompanystatus.UpdatedID = uint64(userdatamap["Id"].(float64))
	result := initializers.DB.First(&ocompanystatus, "id  = ?", sourceMap["ID"])
	if result.Error != nil {
		shortCode := "K0110"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	var targetMap map[string]interface{}
	fmt.Println((targetMap))
	data, _ := json.Marshal(ocompanystatus)

	json.Unmarshal(data, &targetMap)

	for key, _ := range targetMap {

		if val1, ok := sourceMap[key]; ok {
			targetMap[key] = val1
		}

	}

	data, _ = json.Marshal(targetMap)
	json.Unmarshal(data, &ocompanystatus)
	// update modified time
	//ocompanystatus.UpdatedAt = time.Now()
	// update modified user into this field
	fmt.Println("Modified")
	fmt.Println(userdatamap)
	fmt.Println(userdatamap["Id"])
	//	ocompanystatus.UpdatedID = uint64(userdatamap["Id"].(float64))
	//fmt.Println("MOdified User")

	//ocompanystatus.UpdatedID = updateduserid
	result = initializers.DB.Save(&ocompanystatus)

	if result.Error != nil {
		shortCode := "K0115" // Failed to save company
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{"outputs": ocompanystatus})

}
