package utilities

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/paramTypes"
)

type ErrorInfo struct {
	SNo         int    `json:"s_no"`
	KeyField    string `json:"key_field"`
	SourceField string `json:"source_field"`
	Error       string `json:"error"`
}

// Yukesh's Code
func Date2String(iDate time.Time) (odate string) {

	var temp string
	temp = iDate.String()
	temp1 := temp[0:4] + temp[5:7] + temp[8:10]
	// fmt.Println("Rangarajan Ramaujam ***********")
	// fmt.Println(iDate)
	// fmt.Println(temp1)
	odate = temp1
	return odate

}

func GetErrorDesc(iCompany uint, iLanguage uint, iShortCode string) (string, error) {
	var errorenq models.Error

	result := initializers.DB.Find(&errorenq, "company_id = ? and language_id = ? and short_code = ?", iCompany, iLanguage, iShortCode)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", errors.New(" -" + strconv.FormatUint(uint64(iCompany), 10) + "-" + "-" + strconv.FormatUint(uint64(iLanguage), 10) + "-" + " is missing")
	}

	return errorenq.LongCode, nil
}

func GetBusinessDate(iCompany uint, iUser uint, iDepartment uint) (oDate string) {
	var businessdate models.BusinessDate
	// Get with User
	result := initializers.DB.Find(&businessdate, "company_id = ? and user_id = ? and department = ? and user_id IS NOT NULL and department IS NOT NULL", iCompany, iUser, iDepartment)
	if result.RowsAffected == 0 {
		// If User Not Found, get with Department
		result = initializers.DB.Find(&businessdate, "company_id = ? and department = ? and user_id IS NULL ", iCompany, iDepartment)
		if result.RowsAffected == 0 {
			// If Department Not Found, get with comapny
			result = initializers.DB.Find(&businessdate, "company_id = ? and department IS NULL and user_id IS NULL", iCompany)
			if result.RowsAffected == 0 {
				return Date2String(time.Now())

			} else {
				oDate := businessdate.Date
				return oDate
			}
		} else {
			oDate := businessdate.Date
			return oDate
		}

	} else {
		oDate := businessdate.Date
		return oDate
	}

}

func GetItemD(iCompany int, iTable string, iItem string, iFrom string, data *paramTypes.Extradata) error {

	//var sourceMap map[string]interface{}
	var itemparam models.Param
	//	fmt.Println(iCompany, iItem, iFrom)
	results := initializers.DB.Find(&itemparam, "company_id =? and name= ? and item = ? and rec_type = ? and ? between start_date  and  end_date", iCompany, iTable, iItem, "IT", iFrom)

	if results.Error == nil && results.RowsAffected != 0 {
		(*data).ParseData(itemparam.Data)
		return nil
	} else {
		if results.Error != nil {
			return errors.New(results.Error.Error())
		} else {
			return errors.New("No Item Found " + iTable + iItem)
		}

	}
}

// func InputToExcel(c *gin.Context) {

// 	var req RequestData

// 	// Bind the incoming JSON to the RequestData struct
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Define the names referring to the ranges of cells
// 	definedName1 := "p_term"
// 	definedName2 := "p_premium"
// 	definedName3 := "p_sa"

// 	// Extract values from the request
// 	oTerm := req.Term
// 	oPremium := req.Premium

// 	fmt.Println("****************Term******************", oTerm)
// 	fmt.Println("****************Premium******************", oPremium)

// 	// Now generate a new Excel file with the extracted values
// 	newFile := excelize.NewFile()

// 	// Write the values to the new Excel file
// 	newFile.SetCellValue("Sheet1", "D5", oTerm)    // p_term value in cell D5
// 	newFile.SetCellValue("Sheet1", "D6", oPremium) // p_premium value in cell D6

// 	// Use Excel formula to calculate p_sa (D7 = D5 * D6)
// 	newFile.SetCellFormula("Sheet1", "D7", "=D5*D6")

// 	// Set defined names for the new file, similar to the original file
// 	newFile.SetDefinedName(&excelize.DefinedName{
// 		Name:     definedName1,
// 		RefersTo: "Sheet1!$D$5",
// 	})
// 	newFile.SetDefinedName(&excelize.DefinedName{
// 		Name:     definedName2,
// 		RefersTo: "Sheet1!$D$6",
// 	})
// 	newFile.SetDefinedName(&excelize.DefinedName{
// 		Name:     definedName3,
// 		RefersTo: "Sheet1!$D$7",
// 	})

// 	// Save the new file
// 	err := newFile.SaveAs("D:/Go/Output.xlsx")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to save new Excel file",
// 		})
// 		return
// 	}

// 	// Return success response
// 	c.JSON(http.StatusOK, gin.H{
// 		"Result": "New Excel generated successfully",
// 	})
// }

// func ExcelCalculation(c *gin.Context) {

// 	// Assuming RequestData struct has the fields `Term`, `Premium`
// 	var req RequestData

// 	// Bind the incoming JSON to the RequestData struct
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Open the existing Excel file
// 	filePath := "D:/Go/Vengy.xlsx"
// 	f, err := excelize.OpenFile(filePath)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to open Excel file",
// 		})
// 		return
// 	}
// 	defer f.Close()

// 	// Define the names referring to the ranges of cells
// 	definedName1 := "p_product"
// 	definedName2 := "p_sa"
// 	definedName3 := "p_term"
// 	definedName4 := "p_dob"
// 	definedName5 := "p_gender"
// 	definedName6 := "p_startdate"

// 	// Extract values from the request
// 	oProduct := req.Product
// 	oSA := req.SumAssured
// 	oTerm := req.Term
// 	oDob := String2Date(req.Dob)
// 	oGender := req.Gender
// 	oStart := String2Date(req.StartDate)

// 	// Get all the defined names in the workbook
// 	definedNames := f.GetDefinedName()

// 	for _, dn := range definedNames {

// 		// Find and update the cell that corresponds to the defined name `p_term`
// 		if dn.Name == definedName1 {
// 			ref := strings.TrimPrefix(dn.RefersTo, "=")
// 			parts := strings.Split(ref, "!")
// 			if len(parts) != 2 {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Invalid defined name reference format",
// 				})
// 				return
// 			}
// 			sheetName := parts[0]
// 			cellLocation := strings.ReplaceAll(parts[1], "$", "")

// 			// Set the value of the cell referenced by `p_term`
// 			if err := f.SetCellValue(sheetName, cellLocation, oProduct); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Failed to set cell value for " + definedName1,
// 				})
// 				return
// 			}
// 		}

// 		// Find and update the cell that corresponds to the defined name `p_premium`
// 		if dn.Name == definedName2 {
// 			ref := strings.TrimPrefix(dn.RefersTo, "=")
// 			parts := strings.Split(ref, "!")
// 			if len(parts) != 2 {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Invalid defined name reference format",
// 				})
// 				return
// 			}
// 			sheetName := parts[0]
// 			cellLocation := strings.ReplaceAll(parts[1], "$", "")

// 			// Set the value of the cell referenced by `p_premium`
// 			if err := f.SetCellValue(sheetName, cellLocation, oSA); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Failed to set cell value for  " + definedName2,
// 				})
// 				return
// 			}
// 		}

// 		// Set the formula for `p_sa` (p_term * p_premium)
// 		if dn.Name == definedName3 {
// 			ref := strings.TrimPrefix(dn.RefersTo, "=")
// 			parts := strings.Split(ref, "!")
// 			if len(parts) != 2 {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Invalid defined name reference format",
// 				})
// 				return
// 			}
// 			sheetName := parts[0]
// 			cellLocation := strings.ReplaceAll(parts[1], "$", "")

// 			// Set the value of the cell referenced by `p_premium`
// 			if err := f.SetCellValue(sheetName, cellLocation, oTerm); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Failed to set cell value for " + definedName3,
// 				})
// 				return
// 			}

// 			// Set the formula for `p_sa` (p_term * p_premium)
// 			if dn.Name == definedName4 {
// 				ref := strings.TrimPrefix(dn.RefersTo, "=")
// 				parts := strings.Split(ref, "!")
// 				if len(parts) != 2 {
// 					c.JSON(http.StatusBadRequest, gin.H{
// 						"error": "Invalid defined name reference format",
// 					})
// 					return
// 				}
// 				sheetName := parts[0]
// 				cellLocation := strings.ReplaceAll(parts[1], "$", "")

// 				// Set the value of the cell referenced by `p_premium`
// 				if err := f.SetCellValue(sheetName, cellLocation, oDob); err != nil {
// 					c.JSON(http.StatusBadRequest, gin.H{
// 						"error": "Failed to set cell value for " + definedName4,
// 					})
// 					return
// 				}
// 			}
// 		}

// 		// Set the formula for `p_sa` (p_term * p_premium)
// 		if dn.Name == definedName4 {
// 			ref := strings.TrimPrefix(dn.RefersTo, "=")
// 			parts := strings.Split(ref, "!")
// 			if len(parts) != 2 {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Invalid defined name reference format",
// 				})
// 				return
// 			}
// 			sheetName := parts[0]
// 			cellLocation := strings.ReplaceAll(parts[1], "$", "")

// 			// Set the value of the cell referenced by `p_premium`
// 			if err := f.SetCellValue(sheetName, cellLocation, oDob); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Failed to set cell value for " + definedName4,
// 				})
// 				return
// 			}
// 		}

// 		// Set the formula for `p_sa` (p_term * p_premium)
// 		if dn.Name == definedName5 {
// 			ref := strings.TrimPrefix(dn.RefersTo, "=")
// 			parts := strings.Split(ref, "!")
// 			if len(parts) != 2 {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Invalid defined name reference format",
// 				})
// 				return
// 			}
// 			sheetName := parts[0]
// 			cellLocation := strings.ReplaceAll(parts[1], "$", "")

// 			// Set the value of the cell referenced by `p_premium`
// 			if err := f.SetCellValue(sheetName, cellLocation, oGender); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Failed to set cell value for " + definedName5,
// 				})
// 				return
// 			}
// 		}

// 		// Set the formula for `p_sa` (p_term * p_premium)
// 		if dn.Name == definedName6 {
// 			ref := strings.TrimPrefix(dn.RefersTo, "=")
// 			parts := strings.Split(ref, "!")
// 			if len(parts) != 2 {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Invalid defined name reference format",
// 				})
// 				return
// 			}
// 			sheetName := parts[0]
// 			cellLocation := strings.ReplaceAll(parts[1], "$", "")

// 			// Set the value of the cell referenced by `p_premium`
// 			if err := f.SetCellValue(sheetName, cellLocation, oStart); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "Failed to set cell value for " + definedName6,
// 				})
// 				return
// 			}
// 		}
// 	}

// 	// Save the updated file with a new name or overwrite the existing one
// 	if err := f.SaveAs(filePath); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to save the updated Excel file",
// 		})
// 		return
// 	}

// 	// Return success response
// 	c.JSON(http.StatusOK, gin.H{
// 		"Result": "Excel updated successfully",
// 	})
// }

// func ExcelTest() {
// 	unknown, _ := oleutil.CreateObject("Excel.Application")
// 	defer unknown.Release()
// 	excel, _ := unknown.QueryInterface(ole.IID_IDispatch)
// 	defer excel.Release()
// 	oleutil.PutProperty(excel, "Visible", false)
// 	workbooks := oleutil.MustGetProperty(excel, "Workbooks").ToIDispatch()
// 	defer workbooks.Release()
// 	workbookPath := `D:\Go\data.xlsx`
// 	workbook := oleutil.MustCallMethod(workbooks, "Open", workbookPath).ToIDispatch()
// 	defer workbook.Release()
// 	worksheet := oleutil.MustGetProperty(workbook, "Worksheets", "Sheet1").ToIDispatch()
// 	defer worksheet.Release()
// 	oleutil.MustGetProperty(worksheet, "Range", "NUMBER1").ToIDispatch().PutProperty("Value", 50)
// 	oleutil.MustGetProperty(worksheet, "Range", "NUMBER2").ToIDispatch().PutProperty("Value", 20)
// 	a3 := oleutil.MustGetProperty(worksheet, "Range", "SUM").ToIDispatch()
// 	val, _ := a3.GetProperty("Value")
// 	fmt.Println(val.Value())
// 	oleutil.CallMethod(workbook, "Close", false)
// 	oleutil.CallMethod(excel, "Quit")
// }

// type TestSum struct {

func Transpose(matrix []interface{}) []interface{} {
	if len(matrix) == 0 {
		return nil
	}

	// Determine the size of the transposed matrix
	rows, cols := len(matrix), len(matrix[0].([]interface{}))
	transposed := make([]interface{}, cols)
	for i := range transposed {
		transposed[i] = make([]interface{}, rows)
	}

	// Fill the transposed matrix
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j].([]interface{})[i] = matrix[i].([]interface{})[j]
		}
	}

	return transposed
}

func Transpose1(matrix [][]interface{}) [][]interface{} {
	if len(matrix) == 0 {
		return nil
	}

	// Determine the size of the transposed matrix
	rows, cols := len(matrix), len(matrix[0])
	transposed := make([][]interface{}, cols)
	for i := range transposed {
		transposed[i] = make([]interface{}, rows)
	}

	// Fill the transposed matrix
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}

	return transposed
}

func AddNestedValue(jsonMap map[string]interface{}, key string, value interface{}) {
	// Split the key by "#"
	parts := strings.Split(key, "#")

	// Iterate through the parts, creating nested maps as necessary
	currentMap := jsonMap
	for i, part := range parts {
		// If it's the last part, assign the value
		if i == len(parts)-1 {
			currentMap[part] = value
		} else {
			// Otherwise, check if the next level exists as a map
			if _, ok := currentMap[part]; !ok {
				// If it doesn't exist, create a new map for the current part
				currentMap[part] = make(map[string]interface{})
			}
			// Move deeper into the nested map
			currentMap = currentMap[part].(map[string]interface{})
		}
	}
}

func GetNestedValue(jsonMap map[string]interface{}, key string) (interface{}, error) {
	// Split the key by "#"
	parts := strings.Split(key, "#")

	// Navigate through the nested map
	currentMap := jsonMap
	for i, part := range parts {
		// Check if the current part exists in the map
		if value, ok := currentMap[part]; ok {
			// If it's the last part, return the value
			if i == len(parts)-1 {
				return value, nil
			}

			// Otherwise, ensure the next part is also a map[string]interface{}
			if nestedMap, isMap := value.(map[string]interface{}); isMap {
				currentMap = nestedMap
			} else {
				return nil, errors.New("intermediate value is not a map")
			}
		} else {
			// Part not found
			return nil, errors.New("key not found in the map")
		}
	}

	return nil, errors.New("unexpected error")
}

// 	NUMBER1 float64
// 	NUMBER2 float64
// }
