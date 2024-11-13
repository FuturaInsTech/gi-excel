package basiccontrollers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/FuturaInsTech/GoExcel/utilities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllBusinessDate(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetAllBusinessDate" //B0069
	//var userdatamap map[string]interface{}
	//_, err := utilities.GetUserAccess(user, method)
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0086"
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

	var getallbusinessdate []models.BusinessDate
	//userco := userdatamap["CompanyId"]

	var result *gorm.DB

	if searchpagination.SearchString != "" && searchpagination.SearchCriteria != "" {
		result = initializers.DB.Model(&models.BusinessDate{}).Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.BusinessDate{}).
			Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallbusinessdate)

	} else {
		fmt.Println("No Selection ")
		fmt.Println(searchpagination.SearchCriteria)
		fmt.Println(searchpagination.SearchString)
		result = initializers.DB.Model(&models.BusinessDate{}).Where("company_id = ?", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.BusinessDate{}).
			Where("company_id = ?", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getallbusinessdate)
	}

	// if result is null, then give an language ..
	if result.Error != nil {
		shortCode := "K0087" // failed to fetch businessdate
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
			"displayName": "Business Date",
			"fieldName":   "date",
			"dataType":    "string"},
		}

		c.JSON(200, gin.H{

			"AllBusinessDate": getallbusinessdate,
			"Field Map":       fieldMappings,
			"paginationData":  paginationData,
		})

	} else {
		c.JSON(200, gin.H{

			"AllBusinessDate": getallbusinessdate,
			"paginationData":  paginationData,
		})
	}

}

func CreateBusinessDate(c *gin.Context) {
	// store it in working storage variable
	// field description in cobol
	user, _ := c.Get("user") //B0067
	method := "CreateBusinessDate"
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0086"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var createbusinessdate models.BusinessDate

	if c.Bind(&createbusinessdate) != nil {
		shortCode := "K0088" // failed to bind businessdate
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return
	}

	uid := createbusinessdate.UserID
	deptid := createbusinessdate.Department

	var getbusinessdate models.BusinessDate
	result := initializers.DB.First(&getbusinessdate, "user_id  = ? and department = ?", uid, deptid)

	if result.RowsAffected != 0 {
		shortCode := "K0087" // no record found
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	result = initializers.DB.Create(&createbusinessdate)

	if result.Error != nil {
		shortCode := "K0089" // Failed to create businessdate
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	//response
	shortCode := "K0090" // created businessdate
	longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"Result": shortCode + " : " + longDesc,
	})

}

func DeleteBusinessDate(c *gin.Context) {

	user, _ := c.Get("user")
	method := "DeleteBusinessDate" //B0059
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0086"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	delid := c.Param("id")

	var delbusinessdate models.BusinessDate
	result := initializers.DB.First(&delbusinessdate, "id  = ?", delid)

	fmt.Println(delbusinessdate)
	if result.Error != nil {
		shortCode := "K0087" // failed to get
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	result = initializers.DB.Delete(&delbusinessdate)

	if result.Error != nil {
		shortCode := "K0091" // failed to delete
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, "BusinessDate ID "+delid+" is deleted")

}

func GetBusinessDate(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetBusinessDate" //B0058
	// //var userdatamap map[string]interface{}
	iSkipPermission := "Y"
	userdatamap, err := utilities.GetUserAccessNew(user, method, iSkipPermission)

	//userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		fmt.Println("BusinessDate")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Access Failed " + err.Error(),
		})

		return
	}
	// fmt.Println(userdatamap)
	getid := c.Param("id")

	var getbusinessdate models.BusinessDate
	result := initializers.DB.First(&getbusinessdate, "id  = ?", getid)
	if result.Error != nil {
		shortCode := "K0087" // failed to get
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"BusinessDate": getbusinessdate,
	})

}

func ModifyBusinessDate(c *gin.Context) {
	// mapping json to sourceMap
	user, _ := c.Get("user")
	method := "ModifyBusinessDate" //B0059
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0086"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	var sourceMap map[string]interface{}

	if c.Bind(&sourceMap) != nil {
		shortCode := "K0088" // failed to bind businessdate
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	var modifybusinessdate models.BusinessDate
	result := initializers.DB.First(&modifybusinessdate, "id  = ?", sourceMap["ID"])
	if result.Error != nil {
		shortCode := "K0087" // failed to get
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	var targetMap map[string]interface{}
	fmt.Println((targetMap))
	data, _ := json.Marshal(modifybusinessdate)

	json.Unmarshal(data, &targetMap)

	for key, _ := range targetMap {

		if val1, ok := sourceMap[key]; ok {
			targetMap[key] = val1
		}

	}

	data, _ = json.Marshal(targetMap)
	json.Unmarshal(data, &modifybusinessdate)

	result = initializers.DB.Save(&modifybusinessdate)

	if result.Error != nil {
		shortCode := "K0092" // Failed to save
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	c.JSON(200, gin.H{
		"BusinessDate": result,
	})

}

func CloneBusinessDate(c *gin.Context) {
	user, _ := c.Get("user")
	method := "CloneBusinessDate" //B0060
	//var userdatamap map[string]interface{}
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		shortCode := "K0086"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + method,
		})

		return
	}
	fmt.Println(userdatamap)
	sourceid := c.Param("id")
	//fmt.Println("ID " + sid)

	var clonebusinessdate models.BusinessDate
	result := initializers.DB.First(&clonebusinessdate, "id  = ?", sourceid)
	if result.Error != nil {
		shortCode := "K0087" // failed to get
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// Declaring a Map so that i could move all or selective values into my Map
	var sourceMap map[string]interface{}
	// converting an entity(Model) to Json
	data, _ := json.Marshal(clonebusinessdate)
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
	var clnbusinessdate models.BusinessDate
	// converting json to a model
	json.Unmarshal(data, &clnbusinessdate)
	// edecuting query persisting the model
	result = initializers.DB.Create(&clnbusinessdate)
	if result.Error != nil {
		shortCode := "K0092" // Failed to save
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}
	// return the values to postman
	c.JSON(200, gin.H{
		"Cloned": clnbusinessdate,
	})

}

func GetCompanyBusinessDate(c *gin.Context) {
	getid := c.Param("coid")
	getdept := c.Param("deptcode")
	getuserid := c.Param("usercode")

	userco, err := strconv.Atoi(getid)
	userid, err := strconv.Atoi(getuserid)
	deptid, err := strconv.Atoi(getdept)
	if err != nil {
		fmt.Println("Conversion failed:", err)
	}

	businessDate := utilities.GetBusinessDate(uint(userco), uint(userid), uint(deptid))

	c.JSON(200, gin.H{
		"BusinessDate": businessDate,
	})
}

func AddPermissionAuto(c *gin.Context) {
	user, _ := c.Get("user")
	// method := "GetComponentAddEnq"

	userdatamap, _ := utilities.GetUserAccess(user, "")

	// if err != nil {

	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Access Failed " + err.Error(),
	// 	})
	// 	return
	// }

	// userco := userdatamap["CompanyId"]
	// iuserID := int64(userdatamap["Id"].(float64))
	//userco := 1
	iuserID := 1
	paramID := c.Param("id")
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))
	var transactionenq []models.Transaction
	var result *gorm.DB

	result = initializers.DB.Find(&transactionenq, "company_id = ?", userco)

	if result.Error != nil {
		shortCode := "GL037"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return
	}

	var userenq models.User

	result = initializers.DB.First(&userenq, "company_id = ? and id = ?", userco, paramID)

	if result.Error != nil {
		shortCode := "GL254"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + ":" + longDesc + "-" + result.Error.Error(),
		})

		return
	}

	iUserGrp := userenq.UserGroupID

	var permission models.Permission

	for i := 0; i < len(transactionenq); i++ {

		iMethod := transactionenq[i].Description
		a, _ := strconv.Atoi(paramID)
		// permission.UserGroupID = sql.NullInt64{int64(iUserGrp), true}
		// permission.TransactionID = transactionenq[i].ID
		// permission.CompanyID = transactionenq[i].CompanyID

		result = initializers.DB.First(&permission, "company_id = ? and user_group_id = ? and method = ?", userco, iUserGrp, iMethod)
		if result.RowsAffected == 0 {
			var permission models.Permission
			result = initializers.DB.First(&permission, "company_id = ? and user_id = ? and method = ?", userco, paramID, iMethod)
			if result.RowsAffected == 0 {

				var permissionupd models.Permission
				permissionupd.CompanyID = userenq.CompanyID
				permissionupd.Method = iMethod
				permissionupd.ModelName = iMethod
				permissionupd.TransactionID = transactionenq[i].ID
				permissionupd.UserID = sql.NullInt64{int64(a), true}
				permissionupd.UserGroupID = sql.NullInt64{int64(iUserGrp), true}
				permissionupd.UpdatedID = uint64(iuserID)
				result = initializers.DB.Create(&permissionupd)

				if result.Error != nil {
					longDesc, _ := utilities.GetErrorDesc(userco, userlan, "GL255")
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "GL255" + "-" + longDesc,
					})

					return
				}

			}
		}
	}

	c.JSON(http.StatusOK, gin.H{

		"Created Permissions for user": paramID,
	})

}
