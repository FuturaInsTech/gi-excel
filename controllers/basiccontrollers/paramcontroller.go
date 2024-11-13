package basiccontrollers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	// "github.com/FuturaInsTech/gi/excelGenerator"
	"github.com/FuturaInsTech/GoExcel/excelGenerator"
	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/paramTypes"
	"github.com/FuturaInsTech/GoExcel/pdfGenerator"

	// "github.com/FuturaInsTech/gi/pdfGenerator"
	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/FuturaInsTech/GoExcel/utilities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateParam(c *gin.Context) {

	var param models.Param
	var paramdesc models.ParamDesc
	var paramVals struct {
		CompanyId      uint16
		Name           string
		Type           string
		Longdesc       string
		LanguageId     uint8
		ExtraDataExist bool
	}
	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if c.Bind(&paramVals) != nil {
		shortCode := "GL075"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return
	}

	param.CreatedAt = time.Now()
	paramdesc.CreatedAt = time.Now()

	param.CompanyId = paramVals.CompanyId
	param.Name = paramVals.Name
	param.EndDate = "0"
	param.StartDate = "0"
	param.Is_valid = true
	param.RecType = "HE"
	param.LastModUser = 1

	if paramVals.Type != "dated" {

		paramVals.Type = ""
	}

	param.Data = map[string]interface{}{
		"paramType":      paramVals.Type,
		"extraDataExist": paramVals.ExtraDataExist,
	}

	paramdesc.CompanyId = paramVals.CompanyId
	paramdesc.Name = paramVals.Name
	paramdesc.LanguageId = paramVals.LanguageId
	paramdesc.RecType = "HE"
	paramdesc.Longdesc = paramVals.Longdesc
	paramdesc.LastModUser = 1
	result := initializers.DB.Create(&param)

	if result.Error != nil {
		shortCode := "GL076"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	result = initializers.DB.Create(&paramdesc)

	if result.Error != nil {
		shortCode := "GL077"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "param " + paramVals.Name + " is created"})

}

func ModifyParam(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	var sourceMap map[string]interface{}
	if c.Bind(&sourceMap) != nil {
		shortCode := "GL075"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return
	}

	/*Param Update */
	var param models.Param

	result := initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", sourceMap["companyId"], sourceMap["name"], "HE", "", 0)

	if result.Error != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	var targetMapParam map[string]interface{}

	paramdata, _ := json.Marshal(param)

	json.Unmarshal(paramdata, &targetMapParam)

	for key, _ := range targetMapParam {
		key1 := strings.ToLower(string(key[0])) + key[1:]
		if val1, ok := sourceMap[key1]; ok {
			targetMapParam[key] = val1
		}

	}

	paramdata, _ = json.Marshal(targetMapParam)
	json.Unmarshal(paramdata, &param)

	param.UpdatedAt = time.Now()
	param.LastModUser = 1
	param.Seqno = 0

	if sourceMap["type"] != "dated" {

		sourceMap["type"] = ""
	}

	extraDataExist := param.Data["extraDataExist"]

	val, ok := sourceMap["extraDataExist"]
	// If the key exists
	if ok {
		extraDataExist = val.(bool)
	}

	param.Data = map[string]interface{}{
		"paramType":      sourceMap["type"],
		"extraDataExist": extraDataExist,
	}

	result = initializers.DB.Model(&param).Updates(param)

	if result.Error != nil {
		shortCode := "GL079"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	/*Param Update ends*/

	/*ParamDesc Update */

	var paramdesc models.ParamDesc

	result = initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", sourceMap["companyId"], sourceMap["name"], "HE", "", sourceMap["languageId"])

	if result.Error != nil {
		shortCode := "GL080"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	var targetMapParamDesc map[string]interface{}

	paramdata, _ = json.Marshal(paramdesc)

	json.Unmarshal(paramdata, &targetMapParamDesc)

	for key, _ := range targetMapParamDesc {
		key1 := strings.ToLower(string(key[0])) + key[1:]
		if val1, ok := sourceMap[key1]; ok {
			targetMapParamDesc[key] = val1
		}

	}

	paramdata, _ = json.Marshal(targetMapParamDesc)
	json.Unmarshal(paramdata, &paramdesc)

	paramdesc.UpdatedAt = time.Now()
	paramdesc.LastModUser = 1

	result = initializers.DB.Model(&paramdesc).Updates(paramdesc)

	if result.Error != nil {
		shortCode := "GL081"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	/*ParamDesc Update ends*/

	c.JSON(http.StatusOK, gin.H{"message": "param " + sourceMap["name"].(string) + " is modified"})

}

func EnquireParam(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	queryParams := c.Request.URL.Query()

	var param models.Param

	result := initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", queryParams["companyId"], queryParams["name"], "HE", "", 0)

	if result.Error != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	var paramdesc models.ParamDesc

	result = initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", queryParams["companyId"], queryParams["name"], "HE", "", queryParams["languageId"])

	if result.Error != nil {
		shortCode := "GL080"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "" + result.Error.Error(),
		})

		return

	}
	extraDataExist := false

	val, ok := param.Data["extraDataExist"]
	// If the key exists
	if ok {
		extraDataExist = val.(bool)
	}

	paramOut := map[string]interface{}{
		"companyId":      param.CompanyId,
		"name":           param.Name,
		"type":           param.Data["paramType"],
		"extraDataExist": extraDataExist,
		"languageId":     paramdesc.LanguageId,
		"longdesc":       paramdesc.Longdesc,
	}

	c.JSON(http.StatusOK, gin.H{"param": paramOut})

}

func DeleteParam(c *gin.Context) {
	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	queryParams := c.Request.URL.Query()

	var param models.Param

	result := initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", queryParams.Get("companyId"), queryParams.Get("name"), "HE", "", 0)

	if result.Error != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	//check any  items exists then don't allow delete

	var totalRecords int64 = 0

	result = initializers.DB.Model(&models.Param{}).Where("rec_type = ? AND company_id = ? AND name = ? ", "IT", queryParams.Get("companyId"), queryParams.Get("name")).Count(&totalRecords)

	if result.Error != nil {
		shortCode := "GL082"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	if totalRecords > 0 {
		shortCode := "GL083"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	result = initializers.DB.Delete(&param)

	if result.Error != nil {
		shortCode := "GL084"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	var paramdesc models.ParamDesc

	result = initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", queryParams.Get("companyId"), queryParams.Get("name"), "HE", "", queryParams.Get("languageId"))

	if result.Error != nil {
		shortCode := "GL080"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "" + result.Error.Error(),
		})

		return

	}

	result = initializers.DB.Delete(&paramdesc)

	if result.Error != nil {
		shortCode := "GL084"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{"message": gin.H{"message": "param " + queryParams.Get("companyId") + "-" + queryParams.Get("name") + " is deleted"}})

}

func EnquireParams(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	queryParams := c.Request.URL.Query()

	pageNum := 1
	pageSize := 10
	searchString := ""
	searchCriteria := ""
	sortColumn := "name"
	sortDirection := "asc"
	firstTime := false
	company_id := 1
	language_id := 1
	reportOnly := false
	reportType := "pdf"

	if queryParams.Has("pageNum") {
		pageNum, _ = strconv.Atoi(queryParams.Get("pageNum"))

	}
	if queryParams.Has("companyId") {
		company_id, _ = strconv.Atoi(queryParams.Get("companyId"))

	}

	if queryParams.Has("languageId") {
		language_id, _ = strconv.Atoi(queryParams.Get("languageId"))

	}

	if queryParams.Has("pageSize") {
		pageSize, _ = strconv.Atoi(queryParams.Get("pageSize"))
	}

	if queryParams.Has("searchString") {
		searchString = queryParams.Get("searchString")
	}

	if queryParams.Has("searchCriteria") {
		searchCriteria = queryParams.Get("searchCriteria")
	}

	if queryParams.Has("sortColumn") {

		sortColumn = queryParams.Get("sortColumn")

	}

	if queryParams.Has("sortDirection") {

		sortDirection = queryParams.Get("sortDirection")

	}

	if queryParams.Has("firstTime") {

		firstTime, _ = strconv.ParseBool(queryParams.Get("firstTime"))

	}

	if queryParams.Has("reportType") {

		reportType = queryParams.Get("reportType")
		reportOnly = true

	}

	offset := (pageNum - 1) * pageSize

	var params []models.Param
	resp := make([]interface{}, 0)

	var result *gorm.DB
	var totalRecords int64 = 0

	if searchString == "" || searchCriteria == "" {
		result = initializers.DB.Model(&models.Param{}).Where("rec_type = ? AND company_id = ?", "HE", company_id).Count(&totalRecords)
		if result.Error == nil {
			//if for reporting purpose , pagination not required as the whole data will be downloaded in excel/pdf
			if reportOnly {
				result = initializers.DB.Order(sortColumn+" "+sortDirection).Find(&params, "rec_type = ? AND company_id = ?", "HE", company_id)
			} else {

				result = initializers.DB.Order(sortColumn+" "+sortDirection).Limit(pageSize).Offset(offset).Find(&params, "rec_type = ? AND company_id = ?", "HE", company_id)
			}

		}

	} else {

		result = initializers.DB.Model(&models.Param{}).Where(searchCriteria+" LIKE ? AND rec_type = ? AND company_id = ?", "%"+searchString+"%", "HE", company_id).Count(&totalRecords)
		if result.Error == nil {

			//if for reporting purpose , pagination not required as the whole data will be downloaded in excel/pdf
			if reportOnly {
				result = initializers.DB.Order(sortColumn+" "+sortDirection).Find(&params, searchCriteria+" LIKE ? AND rec_type = ? AND company_id = ?", "%"+searchString+"%", "HE", company_id)
			} else {

				result = initializers.DB.Order(sortColumn+" "+sortDirection).Limit(pageSize).Offset(offset).Find(&params, searchCriteria+" LIKE ? AND rec_type = ? AND company_id = ?", "%"+searchString+"%", "HE", company_id)
			}
		}
	}

	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	for i := 0; i < len(params); i++ {

		var paramdesc models.ParamDesc

		result = initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", params[i].CompanyId, params[i].Name, "HE", "", language_id)
		var longdesc = ""
		var languageId = 0
		if result.Error == nil {
			longdesc = paramdesc.Longdesc
			languageId = int(paramdesc.LanguageId)
		}

		paramOut := map[string]interface{}{
			"companyId":  params[i].CompanyId,
			"name":       params[i].Name,
			"type":       params[i].Data["paramType"],
			"languageId": languageId,
			"longdesc":   longdesc,
		}
		resp = append(resp, paramOut)
	}

	if reportOnly {

		var reportData types.ReportData
		reportData.Title = "Business Rules List"
		reportData.ColumnHeadings = []string{"Company", "Param Name", "Param Type", "Short Description", "Long Description"}
		reportData.FirstPageRecCount = 29
		reportData.SubseqPageRecCount = 33
		reportData.ReportType = reportType
		reportData.DataRows = make([][]interface{}, 0)
		for i := 0; i < len(resp); i++ {
			valuesMap := resp[i].(map[string]interface{})
			row := []interface{}{valuesMap["companyId"], valuesMap["name"], valuesMap["type"], valuesMap["shortdesc"], valuesMap["longdesc"]}
			reportData.DataRows = append(reportData.DataRows, row)
		}
		reportData.FileName = "BusinessRulesList.xlsx"
		if reportType == "pdf" {
			reportData.FileName = "BusinessRulesList.pdf"
			reportData.TemplateName = "reportTemplates/ListScreenReport.gohtml"
		}

		sendReportResponse(c, reportData)
		return

	}

	paginationData := map[string]interface{}{
		"totalRecords": totalRecords,
	}

	if firstTime {

		fieldMappings := [1]map[string]string{{
			"displayName": "Param Name",
			"fieldName":   "name",
		}}

		c.JSON(http.StatusOK, gin.H{
			"data":           resp,
			"fieldMapping":   fieldMappings,
			"paginationData": paginationData,
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"data":           resp,
			"paginationData": paginationData,
		})

	}

}
func EnquireParamItem(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	queryParams := c.Request.URL.Query()
	qp := make(map[string]string)
	qp["companyId"] = queryParams.Get("companyId")
	qp["name"] = queryParams.Get("name")

	paramHeader, err := getParamHeader(qp)

	if err != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	reportType := "pdf"
	reportOnly := false
	if queryParams.Has("reportType") {

		reportType = queryParams.Get("reportType")
		reportOnly = true

	}
	//paramType := paramHeader.Data["paramType"]

	paramType := "0"
	seqno := 0
	extraDataExist, _ := paramHeader.Data["extraDataExist"].(bool)

	if paramHeader.Data["paramType"] == "dated" {
		paramType = "D"
		if queryParams.Has("seqno") && queryParams.Get("seqno") != "" {
			seqno, _ = strconv.Atoi(queryParams.Get("seqno"))
		}
	} else {

		if extraDataExist {
			paramType = "1"
		}
	}

	var paramdesc models.ParamDesc

	result := initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", queryParams.Get("companyId"), queryParams.Get("name"), "IT", queryParams.Get("item"), queryParams.Get("languageId"))

	if result.Error != nil {
		shortCode := "GL080"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	/*if paramType == "dated" {

		var params []models.Param
		var resp []interface{}
		result = initializers.DB.Find(&params, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ? ", queryParams.Get("companyId"), queryParams.Get("name"), "IT", queryParams.Get("item"))

		if result.Error != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get param items :" + result.Error.Error(),
			})

			return

		}

		for i := 0; i < len(params); i++ {

			paramOut := map[string]interface{}{
				"companyId":  params[i].CompanyId,
				"name":       params[i].Name,
				"item":       params[i].Item,
				"data":       params[i].Data,
				"languageId": paramdesc.LanguageId,
				"longdesc":   paramdesc.Longdesc,
				"shortdesc":  paramdesc.Shortdesc,
			}
			resp = append(resp, paramOut)

		}

		c.JSON(http.StatusOK, gin.H{"param": resp})

	} else {
	*/
	//get total dated item count
	var totalRecords int64 = 0
	if paramType == "D" {
		result = initializers.DB.Model(&models.Param{}).Where("rec_type = ? AND company_id = ? AND name = ? AND item = ? ", "IT", queryParams.Get("companyId"), queryParams.Get("name"), queryParams.Get("item")).Count(&totalRecords)

		if result.Error != nil {
			shortCode := "GL085"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
			})

			return

		}

	}

	//end get total dated item count

	var param models.Param
	result = initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ? AND seqno=?", queryParams.Get("companyId"), queryParams.Get("name"), "IT", queryParams.Get("item"), seqno)

	if result.Error != nil {
		shortCode := "GL082"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}
	if reportOnly {

		var reportData types.ReportData
		reportData.Title = "Param Data"
		if paramType != "D" {
			reportData.CommonDataHeads = []string{"Company", "Param Name", "Item Name", "Short Description", "Long Description"}
			reportData.CommonDataDetails = []string{strconv.Itoa(int(param.CompanyId)), param.Name, param.Item, paramdesc.Shortdesc, paramdesc.Longdesc}
		} else {

			reportData.CommonDataHeads = []string{"Company", "Param Name", "Item Name", "Language Id", "Short Description", "Long Description", "Start Date", "End Date", "Seq No"}
			date, err := time.Parse("20060102", param.StartDate)
			date1 := ""
			date2 := ""
			if err == nil {
				date1 = date.Format("02-01-2006")
			}

			date, err = time.Parse("20060102", param.EndDate)
			if err == nil {
				date2 = date.Format("02-01-2006")
			}

			reportData.CommonDataDetails = []string{strconv.Itoa(int(param.CompanyId)), param.Name, param.Item, strconv.Itoa(int(paramdesc.LanguageId)), paramdesc.Shortdesc, paramdesc.Longdesc, date1, date2, strconv.Itoa(int(param.Seqno))}
		}

		reportData.ReportType = reportType
		//check if the top level field is an array ..if yes then set is_slice field to true . this needs different treatment
		is_slice := false
		key := ""

		if len(param.Data) == 1 {

			for k, v := range param.Data {

				if reflect.TypeOf(v).Kind() == reflect.Slice {
					is_slice = true
					key = k
					break
				}

			}

		}

		if !is_slice {

			for k, v := range param.Data {
				value := ""
				switch v := v.(type) {
				case int:
					value = strconv.Itoa(v)
				case float64:
					value = fmt.Sprintf("%f", v)
				case string:
					value = v
				case []interface{}:
					arr := make([]string, 0)

					for _, v1 := range v {

						switch i := v1.(type) {
						case int:
							arr = append(arr, strconv.Itoa(i))
						case string:
							arr = append(arr, i)
						case float64:
							arr = append(arr, fmt.Sprintf("%f", i))
						default:
							arr = append(arr, "unknownType")
						}

					}
					value = strings.Join(arr, ",")

				default:

					value = "unknownType"
				}

				reportData.CommonDataHeads = append(reportData.CommonDataHeads, strings.ToUpper(string(k[0]))+k[1:])
				reportData.CommonDataDetails = append(reportData.CommonDataDetails, value)

			}

		} else {

			i := 0
			reportData.DataRows = make([][]interface{}, 0)
			dataKeys := make([]string, 0)
			for _, v := range param.Data[key].([]interface{}) {
				valueMap := v.(map[string]interface{})
				row := make([]interface{}, 0)
				if i == 0 {
					for k := range valueMap {
						dataKeys = append(dataKeys, k)
					}
					sort.Strings(dataKeys)
				}

				for _, key := range dataKeys {

					switch l := valueMap[key].(type) {
					case float64:
						row = append(row, fmt.Sprintf("%f", l))
					case []interface{}:
						arr := make([]string, 0)
						for _, v1 := range l {

							switch k := v1.(type) {
							case int:
								arr = append(arr, strconv.Itoa(k))
							case string:
								arr = append(arr, k)
							case float64:
								arr = append(arr, fmt.Sprintf("%f", k))
							default:
								arr = append(arr, "unknown_type")
							}

						}

						row = append(row, strings.Join(arr, ","))

					default:
						row = append(row, valueMap[key])
					}

				}
				reportData.DataRows = append(reportData.DataRows, row)
				i++
			}

			for _, v := range dataKeys {
				reportData.ColumnHeadings = append(reportData.ColumnHeadings, strings.ToUpper(string(v[0]))+v[1:])
			}

			reportData.FirstPageRecCount = 29
			reportData.SubseqPageRecCount = 33

		}
		reportData.FileName = "ParamData" + "_" + param.Name + "_" + param.Item + ".xlsx"
		if reportType == "pdf" {
			reportData.FileName = "ParamData" + "_" + param.Name + "_" + param.Item + ".pdf"
			reportData.TemplateName = "reportTemplates/ListScreenReport.gohtml"
		}

		sendReportResponse(c, reportData)
		return

	}

	if paramType != "D" {

		paramOut := map[string]interface{}{
			"companyId":  param.CompanyId,
			"name":       param.Name,
			"item":       param.Item,
			"data":       param.Data,
			"type":       paramType,
			"languageId": paramdesc.LanguageId,
			"longdesc":   paramdesc.Longdesc,
			"shortdesc":  paramdesc.Shortdesc,
		}

		c.JSON(http.StatusOK, gin.H{"param": paramOut})

	} else {

		paginationData := map[string]interface{}{
			"totalRecords": totalRecords,
		}

		paramOut := map[string]interface{}{
			"companyId":  param.CompanyId,
			"name":       param.Name,
			"item":       param.Item,
			"data":       param.Data,
			"type":       paramType,
			"startDate":  param.StartDate,
			"endDate":    param.EndDate,
			"languageId": paramdesc.LanguageId,
			"longdesc":   paramdesc.Longdesc,
			"shortdesc":  paramdesc.Shortdesc,
		}

		c.JSON(http.StatusOK, gin.H{"param": paramOut,
			"paginationData": paginationData,
		})

	}

}

func getParamHeader(queryParams map[string]string) (models.Param, error) {

	var param models.Param

	result := initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", queryParams["companyId"], queryParams["name"], "HE", "", 0)

	if result.Error != nil {

		return param, errors.New(result.Error.Error())

	}

	return param, nil

}

func CreateParamItem(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	var param models.Param
	var paramVals struct {
		CompanyId  uint16
		Name       string
		Item       string
		Data       map[string]interface{}
		StartDate  string
		EndDate    string
		Seqno      uint16
		Longdesc   string
		Shortdesc  string
		LanguageId uint8
	}

	if c.Bind(&paramVals) != nil {
		shortCode := "GL086"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	qp := make(map[string]string)
	qp["companyId"] = fmt.Sprintf("%v", paramVals.CompanyId)
	qp["name"] = paramVals.Name
	paramHeader, err := getParamHeader(qp)

	if err != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + err.Error(),
		})

		return

	}

	paramType := paramHeader.Data["paramType"]

	if paramType == "dated" {

		if paramVals.EndDate == "" || paramVals.StartDate == "" {
			shortCode := "GL087"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})
			return
		}

	} else {
		if paramVals.Seqno > 0 {
			shortCode := "GL088"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})
			return
		}

	}

	if paramVals.Seqno > 0 {
		var params_existing []models.Param
		result := initializers.DB.Order("seqno asc").Find(&params_existing, "company_id  = ? AND  name = ?  AND item = ? AND  rec_type = ?", paramVals.CompanyId, paramVals.Name, paramVals.Item, "IT")

		if result.Error != nil {
			shortCode := "GL089"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})
			return
		}
		if result.RowsAffected != int64(paramVals.Seqno) {
			shortCode := "GL088"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})
			return
		}

	}

	param.CreatedAt = time.Now()
	param.CompanyId = paramVals.CompanyId
	param.Name = paramVals.Name
	param.Item = paramVals.Item
	if paramType != "dated" {
		param.EndDate = "0"
		param.StartDate = "0"
	} else {
		param.EndDate = paramVals.EndDate
		param.StartDate = paramVals.StartDate
	}
	param.Is_valid = true
	param.RecType = "IT"
	param.LastModUser = 1
	param.Data = paramVals.Data
	param.Seqno = paramVals.Seqno

	result := initializers.DB.Create(&param)

	if result.Error != nil {
		shortCode := "GL076"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	//create param description only for the first dated sequence item
	if paramVals.Seqno == 0 {
		var paramdesc models.ParamDesc
		paramdesc.CreatedAt = time.Now()
		paramdesc.CompanyId = paramVals.CompanyId
		paramdesc.Name = paramVals.Name
		paramdesc.Item = paramVals.Item
		paramdesc.LanguageId = paramVals.LanguageId
		paramdesc.RecType = "IT"
		paramdesc.Longdesc = paramVals.Longdesc
		paramdesc.Shortdesc = paramVals.Shortdesc
		paramdesc.LastModUser = 1

		result = initializers.DB.Create(&paramdesc)

		if result.Error != nil {
			shortCode := "GL076"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})

			return

		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "param item " + paramVals.Name + "-" + paramVals.Item + "-" + strconv.Itoa(int(paramVals.Seqno)) + " is created"})

}

func CloneParamItem(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CloneParamItem" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	var paramVals struct {
		CompanyId  uint16
		Name       string
		FromItem   string
		ToItem     string
		Longdesc   string
		Shortdesc  string
		LanguageId uint8
	}

	if c.Bind(&paramVals) != nil {
		shortCode := "GL086"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}
	if paramVals.ToItem == "" {

		shortCode := "GL086"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "to item cannot be blank",
		})

		return

	}

	qp := make(map[string]string)
	qp["companyId"] = fmt.Sprintf("%v", paramVals.CompanyId)
	qp["name"] = paramVals.Name
	_, err := getParamHeader(qp)

	if err != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + err.Error(),
		})

		return

	}

	seqno := 0
	var params []models.Param
	//get all the instances of clone from item..means all seqnumbers
	result := initializers.DB.Find(&params, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ? ", paramVals.CompanyId, paramVals.Name, "IT", paramVals.FromItem)

	if result.Error != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}
	//if no items exists return error
	if len(params) == 0 {

		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + "no from item found",
		})

		return

	}

	//check if the clone to item already exists
	var exists bool
	result = initializers.DB.Model(models.Param{}).Select("count(*) > 0").Where("company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", paramVals.CompanyId, paramVals.Name, "IT", paramVals.ToItem, seqno).Find(&exists)
	// return error even if already one exists
	if result.Error != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	if exists {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + "to item already exists",
		})

		return

	}
	params1 := make([]models.Param, 0)
	for _, param := range params {

		param.Item = paramVals.ToItem
		param.CreatedAt = time.Now()
		params1 = append(params1, param)
	}

	result = initializers.DB.Create(&params1)

	if result.Error != nil {
		shortCode := "GL076"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + result.Error.Error(),
		})

		return

	}

	//create descrition
	var paramdesc models.ParamDesc
	paramdesc.CreatedAt = time.Now()
	paramdesc.CompanyId = paramVals.CompanyId
	paramdesc.Name = paramVals.Name
	paramdesc.Item = paramVals.ToItem
	paramdesc.LanguageId = paramVals.LanguageId
	paramdesc.RecType = "IT"
	paramdesc.Longdesc = paramVals.Longdesc
	paramdesc.Shortdesc = paramVals.Shortdesc
	paramdesc.LastModUser = 1

	result = initializers.DB.Create(&paramdesc)

	if result.Error != nil {
		shortCode := "GL076"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "param item " + paramVals.Name + "-" + paramVals.ToItem + " is cloned from " + paramVals.FromItem})

}

func ModifyParamItem(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	var sourceMap map[string]interface{}
	if c.Bind(&sourceMap) != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

		return

	}

	qp := make(map[string]string)
	qp["companyId"] = fmt.Sprintf("%v", sourceMap["companyId"])
	qp["name"] = sourceMap["name"].(string)
	_, err := getParamHeader(qp)

	if err != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + err.Error(),
		})

		return
	}
	seqno := 0
	val, ok := sourceMap["seqno"]

	if ok {
		seqno = int(val.(float64))
	}

	//paramType := paramHeader.Data["paramType"]
	/*
		if paramType == "dated" {
			_, ok := sourceMap["startDate"]

			if !ok {

				c.JSON(http.StatusBadRequest, gin.H{
					"error": "start date is mandatory for dated ",
				})
				return
			} else {
				if sourceMap["startDate"] == "" {

					c.JSON(http.StatusBadRequest, gin.H{
						"error": "start date is mandatory for dated ",
					})
					return

				}
			}

		}
	*/
	/*Param Update */
	var param models.Param

	/*if paramType == "dated" {
		result := initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and start_date =?", sourceMap["companyId"], sourceMap["name"], "IT", sourceMap["item"], sourceMap["startDate"])

		if result.Error != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get param :" + result.Error.Error(),
			})

			return

		}
	} else { */

	result := initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", sourceMap["companyId"], sourceMap["name"], "IT", sourceMap["item"], seqno)

	if result.Error != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	//}

	var targetMapParam map[string]interface{}

	paramdata, _ := json.Marshal(param)

	json.Unmarshal(paramdata, &targetMapParam)

	for key, _ := range targetMapParam {
		key1 := strings.ToLower(string(key[0])) + key[1:]
		if val1, ok := sourceMap[key1]; ok {
			//	if key == "Data" {
			//		jsonstr, _ := json.Marshal(val1)
			//		targetMapParam[key] = string(jsonstr)

			//	} else {

			targetMapParam[key] = val1
			//	}

		}

	}
	param.Data = nil
	paramdata, _ = json.Marshal(targetMapParam)
	json.Unmarshal(paramdata, &param)

	param.UpdatedAt = time.Now()
	param.LastModUser = 1

	result = initializers.DB.Model(&param).Where("company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", sourceMap["companyId"], sourceMap["name"], "IT", sourceMap["item"], seqno).Updates(param)

	if result.Error != nil {
		shortCode := "GL079"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	/*Param Update ends*/

	/*ParamDesc Update */

	var paramdesc models.ParamDesc

	result = initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", sourceMap["companyId"], sourceMap["name"], "IT", sourceMap["item"], sourceMap["languageId"])

	if result.Error != nil {
		shortCode := "GL080"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	var targetMapParamDesc map[string]interface{}

	paramdata, _ = json.Marshal(paramdesc)

	json.Unmarshal(paramdata, &targetMapParamDesc)

	for key, _ := range targetMapParamDesc {
		key1 := strings.ToLower(string(key[0])) + key[1:]
		if val1, ok := sourceMap[key1]; ok {
			targetMapParamDesc[key] = val1
		}

	}

	paramdata, _ = json.Marshal(targetMapParamDesc)
	json.Unmarshal(paramdata, &paramdesc)

	paramdesc.UpdatedAt = time.Now()
	paramdesc.LastModUser = 1

	result = initializers.DB.Model(&paramdesc).Updates(paramdesc)

	if result.Error != nil {
		shortCode := "GL081"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}
	/*ParamDesc Update ends*/

	c.JSON(http.StatusOK, gin.H{"message": "param " + sourceMap["name"].(string) + "-" + sourceMap["item"].(string) + " is modified"})

}

func DeleteParamItem(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	queryParams := c.Request.URL.Query()

	qp := make(map[string]string)
	qp["companyId"] = fmt.Sprintf("%v", queryParams.Get("companyId"))
	qp["name"] = queryParams.Get("name")
	paramHeader, err := getParamHeader(qp)

	if err != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + err.Error(),
		})

		return

	}

	paramType := paramHeader.Data["paramType"]

	seqno := 0

	if paramType == "dated" && queryParams.Has("seqno") && queryParams.Get("seqno") != "" {
		seqno, _ = strconv.Atoi(queryParams.Get("seqno"))
	}

	/*	if paramType == "dated" {
			ok := queryParams.Has("startDate")

			if !ok {

				c.JSON(http.StatusBadRequest, gin.H{
					"error": "start date is mandatory for dated ",
				})
				return
			} else {
				if queryParams.Get("startDate") == "" {

					c.JSON(http.StatusBadRequest, gin.H{
						"error": "start date is mandatory for dated ",
					})
					return

				}
			}

		}
	*/
	var param models.Param
	if paramType == "dated" {

		//first check additional dated items exists with seqno> seqno passed, if yes not allowed to delete if those are not deleted

		var totalRecords int64 = 0

		result := initializers.DB.Model(&models.Param{}).Where("rec_type = ? AND company_id = ? AND name = ?  AND item = ? and seqno > ? ", "IT", queryParams.Get("companyId"), queryParams.Get("name"), queryParams.Get("item"), seqno).Count(&totalRecords)

		if result.Error != nil {
			shortCode := "GL082"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
			})

			return

		}

		if totalRecords > 0 {
			shortCode := "GL090"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
			})

			return

		}

		result = initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", queryParams.Get("companyId"), queryParams.Get("name"), "IT", queryParams.Get("item"), seqno)

		if result.Error != nil {
			shortCode := "GL078"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
			})

			return

		}

		result = initializers.DB.Delete(&param)

		if result.Error != nil {
			shortCode := "GL084"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})

			return

		}

	} else {

		result := initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", queryParams.Get("companyId"), queryParams.Get("name"), "IT", queryParams.Get("item"), 0)

		if result.Error != nil {

			shortCode := "GL078"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
			})

			return

		}

		result = initializers.DB.Delete(&param)

		if result.Error != nil {
			shortCode := "GL084"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})

			return

		}

	}
	//delete description only if the seqno 0 item is deleted
	if seqno == 0 {

		var paramdesc models.ParamDesc

		result := initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", queryParams.Get("companyId"), queryParams.Get("name"), "IT", queryParams.Get("item"), queryParams.Get("languageId"))

		if result.Error != nil {
			shortCode := "GL080"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
			})

			return

		}

		result = initializers.DB.Delete(&paramdesc)

		if result.Error != nil {
			shortCode := "GL084"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})

			return

		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "param " + queryParams.Get("companyId") + "-" + queryParams.Get("name") + "-" + queryParams.Get("item") + " is deleted"})

}

func EnquireParamItems(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	queryParams := c.Request.URL.Query()
	//
	pageNum := 1
	pageSize := 100
	searchString := ""
	searchCriteria := ""
	sortColumn := "name"
	sortDirection := "asc"
	firstTime := false
	getAllInstances := false
	company_id := 1
	language_id := 1
	reportOnly := false
	reportType := "pdf"

	if queryParams.Has("pageNum") {
		pageNum, _ = strconv.Atoi(queryParams.Get("pageNum"))

	}

	if queryParams.Has("pageSize") {
		pageSize, _ = strconv.Atoi(queryParams.Get("pageSize"))
	}

	if queryParams.Has("searchString") {
		searchString = queryParams.Get("searchString")
	}

	if queryParams.Has("searchCriteria") {
		searchCriteria = queryParams.Get("searchCriteria")
	}

	if queryParams.Has("sortColumn") {

		sortColumn = queryParams.Get("sortColumn")

	}

	if queryParams.Has("sortDirection") {

		sortDirection = queryParams.Get("sortDirection")

	}

	if queryParams.Has("firstTime") {

		firstTime, _ = strconv.ParseBool(queryParams.Get("firstTime"))

	}

	if queryParams.Has("getAllInstances") {

		getAllInstances, _ = strconv.ParseBool(queryParams.Get("getAllInstances"))

	}

	if queryParams.Has("languageId") {

		language_id, _ = strconv.Atoi(queryParams.Get("languageId"))

	}

	if queryParams.Has("reportType") {

		reportType = queryParams.Get("reportType")
		reportOnly = true

	}

	offset := (pageNum - 1) * pageSize

	var result *gorm.DB
	var headerParam models.Param

	// get two fields from param header
	result = initializers.DB.Model(&models.Param{}).Find(&headerParam, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =? and is_valid =?", queryParams["companyId"], queryParams["name"], "HE", "", 0, true)

	if result.Error != nil {
		shortCode := "GL078"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}

	if result.RowsAffected == 0 {
		shortCode := "GL092"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + " - " + result.Error.Error(),
		})

		return

	}
	paramType := "0"
	extraDataExist, _ := headerParam.Data["extraDataExist"].(bool)

	if headerParam.Data["paramType"] == "dated" {
		paramType = "D"
	} else {

		if extraDataExist {
			paramType = "1"
		}
	}
	// get param header desc

	var paramHeaderdesc models.ParamDesc

	result = initializers.DB.First(&paramHeaderdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", queryParams["companyId"], queryParams["name"], "HE", "", language_id)

	//
	var params []models.Param
	resp := make([]interface{}, 0)

	var totalRecords int64 = 0

	if searchString == "" || searchCriteria == "" {

		query := make(map[string]interface{})
		query["rec_type"] = "IT"
		query["company_id"] = company_id
		query["name"] = queryParams.Get("name")
		if !getAllInstances {
			query["seqno"] = 0
		}

		result = initializers.DB.Model(&models.Param{}).Where(query).Count(&totalRecords)
		if result.Error == nil {
			if reportOnly {
				result = initializers.DB.Order(sortColumn+" "+sortDirection).Find(&params, query)

			} else {
				result = initializers.DB.Order(sortColumn+" "+sortDirection).Limit(pageSize).Offset(offset).Find(&params, query)
			}

		}

	} else {

		query := searchCriteria + " LIKE ? AND rec_type = ? AND company_id = ? AND name = ?"
		values := []interface{}{"%" + searchString + "%", "IT", company_id, queryParams.Get("name")}
		if !getAllInstances {
			query = query + " AND seqno = ?"
			values = append(values, 0)
		}
		result = initializers.DB.Model(&models.Param{}).Where(query, values...).Count(&totalRecords)

		if result.Error == nil {

			if reportOnly {
				result = initializers.DB.Order(sortColumn+" "+sortDirection).Where(query, values...).Find(&params)

			} else {
				result = initializers.DB.Order(sortColumn+" "+sortDirection).Limit(pageSize).Offset(offset).Where(query, values...).Find(&params)
			}

		}

	}

	if result.Error != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + result.Error.Error(),
		})

		return

	}

	for i := 0; i < len(params); i++ {

		var paramdesc models.ParamDesc

		result = initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", params[i].CompanyId, params[i].Name, "IT", params[i].Item, language_id)
		var longdesc = ""
		var shortdesc = ""
		var languageId = 0
		if result.Error == nil {
			longdesc = paramdesc.Longdesc
			shortdesc = paramdesc.Shortdesc
			languageId = int(paramdesc.LanguageId)
		}

		paramOut := map[string]interface{}{
			"companyId": params[i].CompanyId,
			"name":      params[i].Name,
			"item":      params[i].Item,
			"seqno":     params[i].Seqno,
			//"data":       params[i].Data,
			"languageId": languageId,
			"longdesc":   longdesc,
			"shortdesc":  shortdesc,
		}

		if paramType == "D" && getAllInstances {

			date, err := time.Parse("20060102", params[i].StartDate)
			if err == nil {
				paramOut["startDate"] = date.Format("02/01/2006")
			}

			date, err = time.Parse("20060102", params[i].EndDate)
			if err == nil {
				paramOut["endDate"] = date.Format("02/01/2006")
			}

		}
		resp = append(resp, paramOut)
	}

	if reportOnly {

		var reportData types.ReportData
		reportData.Title = "Param Table Items List"
		reportData.ColumnHeadings = []string{"Company", "Param Name", "Item Name", "Short Description", "Long Description"}
		reportData.FirstPageRecCount = 29
		reportData.SubseqPageRecCount = 33
		reportData.ReportType = reportType
		reportData.DataRows = make([][]interface{}, 0)
		for i := 0; i < len(resp); i++ {
			valuesMap := resp[i].(map[string]interface{})
			row := []interface{}{valuesMap["companyId"], valuesMap["name"], valuesMap["item"], valuesMap["shortdesc"], valuesMap["longdesc"]}
			reportData.DataRows = append(reportData.DataRows, row)
		}
		reportData.FileName = "ParamItemsList.xlsx"
		if reportType == "pdf" {
			reportData.FileName = "ParamItemsList.pdf"
			reportData.TemplateName = "reportTemplates/ListScreenReport.gohtml"
		}

		sendReportResponse(c, reportData)
		return

	}

	paginationData := map[string]interface{}{
		"totalRecords": totalRecords,
	}

	if firstTime {

		fieldMappings := [1]map[string]string{{
			"displayName": "Item Name",
			"fieldName":   "item",
		}}

		c.JSON(http.StatusOK, gin.H{
			"data":           resp,
			"fieldMapping":   fieldMappings,
			"paginationData": paginationData,
			"paramType":      paramType,
			"paramLongDesc":  paramHeaderdesc.Longdesc,
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"data":           resp,
			"paginationData": paginationData,
			"paramType":      paramType,
			"paramLongDesc":  paramHeaderdesc.Longdesc,
		})

	}
}

func GetParamExtraData(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	qp := make(map[string]string)
	for k := range queryParams {
		//fmt.Println(k, " => ", v)
		qp[k] = queryParams.Get(k)
	}

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	iSkipPermission := "Y"
	userdatamap, err := utilities.GetUserAccessNew(user, method, iSkipPermission)
	//userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))
	var userMap map[string]interface{}
	data1, _ := json.Marshal(user)
	json.Unmarshal(data1, &userMap)

	qp["LanguageId"] = fmt.Sprintf("%v", userMap["LanguageId"].(float64))
	var extradataparam paramTypes.Extradata
	// if qp["name"] == "Q0011" {
	// 	var q0011data paramTypes.Q0011Data
	// 	extradataparam = &q0011data
	// }
	// if qp["name"] == "Q0017" {
	// 	var q0017data paramTypes.Q0017Data
	// 	extradataparam = &q0017data
	// }
	// if qp["name"] == "Q0016" {
	// 	var q0016data paramTypes.Q0016Data
	// 	extradataparam = &q0016data
	// }
	// if qp["name"] == "Q0015" {
	// 	var q0015data paramTypes.Q0015Data
	// 	extradataparam = &q0015data
	// }
	// if qp["name"] == "Q0014" {
	// 	var q0014data paramTypes.Q0014Data
	// 	extradataparam = &q0014data
	// }
	// if qp["name"] == "Q0005" {
	// 	var q0005data paramTypes.Q0005Data
	// 	extradataparam = &q0005data
	// }
	// if qp["name"] == "Q0006" {
	// 	var q0006data paramTypes.Q0006Data
	// 	extradataparam = &q0006data
	// }
	// if qp["name"] == "P0044" {
	// 	var p0044data paramTypes.P0044Data
	// 	extradataparam = &p0044data
	// }
	// if qp["name"] == "P0061" {
	// 	var p0061data paramTypes.P0061Data
	// 	extradataparam = &p0061data
	// }
	companyID, _ := strconv.Atoi(qp["company_id"])
	err = utilities.GetItemD(companyID, qp["name"], qp["item"], qp["date"], &extradataparam)
	if err != nil {
		shortCode := "GL058"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc + "-" + err.Error(),
		})

		return

	}
	resp := extradataparam.GetFormattedData(qp)

	c.JSON(http.StatusOK, resp)
}

func sendReportResponse(c *gin.Context, reportData types.ReportData) {
	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if reportData.ReportType == "excel" {
		buf := &bytes.Buffer{}
		err := excelGenerator.GenerateExcelReport(reportData, buf)
		if err != nil {
			shortCode := "GL093"
			longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": shortCode + " : " + longDesc,
			})

		} else {

			c.Header("Content-Disposition", "inline; filename="+reportData.FileName)
			c.Header("Access-Control-Expose-Headers", "Content-Disposition")
			c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
		}

	} else if reportData.ReportType == "pdf" {

		r := pdfGenerator.NewRequestPdf("")

		if err := r.ParseTemplate(reportData.TemplateName, reportData); err == nil {
			buf := new(bytes.Buffer)
			r.GeneratePDF(buf)
			c.Header("Content-Disposition", "inline; filename="+reportData.FileName)
			c.Header("Access-Control-Expose-Headers", "Content-Disposition")
			c.Data(http.StatusOK, "application/pdf", buf.Bytes())
		} else {
			// shortCode := "GL094"
			// longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		}

	} else {
		shortCode := "GL095"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})

	}

}

func UploadParamData(c *gin.Context) {

	user, _ := c.Get("user")
	method := "CreateParam" //B0054
	//var userdatamap map[string]interface{}
	userdatamap, _ := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		shortCode := "GL096"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})
		return
	}

	defer file.Close()

	err = utilities.UploadParamDataItems(file)

	if err != nil {
		shortCode := "GL097"
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shortCode + " : " + longDesc,
		})
		return

	}
	shortCode := "GL357"
	longDesc, _ := utilities.GetErrorDesc(userco, userlan, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"data": shortCode + " : " + longDesc,
	})

}
