package excel

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/paramTypes"
	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/FuturaInsTech/GoExcel/utilities"
	"github.com/gin-gonic/gin"
)

func ProcessExcel(c *gin.Context) {

	var req models.RequestData

	// Bind the incoming JSON to the RequestData struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputMap := make(map[string]interface{})

	inputMap["p_product"] = []interface{}{[]interface{}{req.Product}}
	inputMap["p_sa"] = []interface{}{[]interface{}{req.SumAssured}}
	inputMap["p_term"] = []interface{}{[]interface{}{req.Term}}
	inputMap["p_dob"] = []interface{}{[]interface{}{req.Dob}}
	inputMap["p_gender"] = []interface{}{[]interface{}{req.Gender}}
	inputMap["p_startdate"] = []interface{}{[]interface{}{req.StartDate}}
	inputMap["p_basiccoverage"] = []interface{}{[]interface{}{req.Coverage}}
	inputMap["p_ppt"] = []interface{}{[]interface{}{req.Ppt}}
	inputMap["p_freq"] = []interface{}{[]interface{}{req.Freq}}
	inputMap["p_smoker"] = []interface{}{[]interface{}{req.Smoker}}

	outputNames := []interface{}{"p_premium", "p_age", "c_age", "c_prem", "c_term"}

	excelmanager, err := utilities.NewExcelManager("D:\\Go\\END.xlsx")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get excel manager" + err.Error(),
		})

		return
	}

	defer excelmanager.Close()

	outputMap, err := excelmanager.NamedRangeSetAndGet(inputMap, outputNames)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get results:" + err.Error(),
		})

		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"Result": outputMap,
	})
}

func ProcessExcel1(c *gin.Context) {
	var requestMap map[string]interface{}
	if c.Bind(&requestMap) != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unable to bind values",
		})

		return
	}

	serviceName := c.Param("serviceName")

	iDate := "20240101"

	var e0001data paramTypes.E0001Data
	var extradataE0001 paramTypes.Extradata = &e0001data
	err := utilities.GetItemD(1, "E0001", serviceName, iDate, &extradataE0001)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Service Not Found" + err.Error(),
		})
		return

	}

	var e0002data paramTypes.E0002Data
	var extradataE0002 paramTypes.Extradata = &e0002data

	err = utilities.GetItemD(1, "E0002", serviceName, iDate, &extradataE0002)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Service Not Found" + err.Error(),
		})
	}
	outputfieldDataMap := make(map[string]paramTypes.E0002)
	outputFields := make([]interface{}, 0)
	inputMap := make(map[string]interface{})

	for _, field := range e0002data.FieldArray {
		fmt.Println(field.Mandatory, "**********************")
		if field.FieldMode == types.Input {
			// val, ok := requestMap[field.JsonName]
			val, err := utilities.GetNestedValue(requestMap, field.JsonName)
			// If the key does not exist error
			if err != nil {
				if field.Mandatory {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "mandatory field " + field.JsonName + " must be present in input",
					})

					return
				}
				// If the key exist proceed
			} else {

				switch field.FieldType {
				case types.Single:
					inputMap[field.ExcelName] = []interface{}{[]interface{}{val}}
				case types.OneDArray:
					if field.Orientation == types.Horizontal {
						inputMap[field.ExcelName] = []interface{}{val.([]interface{})}
					} else {
						array := val.([]interface{})
						interfaceSlice := make([]interface{}, len(array))

						// Fill the interface slice
						for i, v := range array {
							interfaceSlice[i] = []interface{}{v}
						}
						inputMap[field.ExcelName] = interfaceSlice
					}

				case types.TwoDArray:
					if field.Orientation == types.Horizontal {
						inputMap[field.ExcelName] = val

					} else {

						inputMap[field.ExcelName] = utilities.Transpose(val.([]interface{}))
					}
				case types.OneDMap:
					var outerkeys []string

					// Unmarshal JSON string to a slice of strings
					err := json.Unmarshal([]byte(field.OuterKeys), &outerkeys)
					if err != nil {
						fmt.Println("Error:Unable parse the map keys: ", err)
					}

					valMap := val.(map[string]interface{})
					valArray := make([]interface{}, 0)
					for _, mapkey := range outerkeys {

						if mapvalue, exists := valMap[mapkey]; exists {
							valArray = append(valArray, mapvalue)
						} else {
							valArray = append(valArray, nil)
						}
					}
					if field.Orientation == types.Horizontal {

						inputMap[field.ExcelName] = []interface{}{valArray}
					} else {
						interfaceSlice := make([]interface{}, len(valArray))

						for i, v := range valArray {
							interfaceSlice[i] = []interface{}{v}
						}
						inputMap[field.ExcelName] = interfaceSlice

					}

				case types.TwoDMap:
					var outerkeys []string
					var innerkeys []string

					// Unmarshal JSON string to a slice of strings
					err := json.Unmarshal([]byte(field.OuterKeys), &outerkeys)
					if err != nil {
						fmt.Println("Error:Unable parse the map outer keys: ", err)
					}

					err = json.Unmarshal([]byte(field.InnerKeys), &innerkeys)
					if err != nil {
						fmt.Println("Error:Unable parse the map inner keys: ", err)
					}

					valMap := val.(map[string]interface{})
					valArray := make([]interface{}, 0)

					for _, mapkey1 := range outerkeys {

						if mapvalue1, exists := valMap[mapkey1]; exists {
							valArray1 := make([]interface{}, 0)
							for _, mapkey2 := range innerkeys {

								mapvalue2 := mapvalue1.(map[string]interface{})

								if mapvalue, exists := mapvalue2[mapkey2]; exists {

									valArray1 = append(valArray1, mapvalue)

								} else {

									valArray1 = append(valArray1, nil)

								}
							}

							valArray = append(valArray, valArray1)
						} else {
							valArray = append(valArray, nil)
						}
					}
					if field.Orientation == types.Horizontal {

						inputMap[field.ExcelName] = valArray
					} else {

						inputMap[field.ExcelName] = utilities.Transpose(valArray)

					}

				case types.TwoDArrayMap:
					var innerkeys []string

					// Unmarshal JSON string to a slice of strings

					err := json.Unmarshal([]byte(field.InnerKeys), &innerkeys)
					if err != nil {
						fmt.Println("Error:Unable parse the map inner keys: ", err)
					}

					valArray := val.([]interface{})
					opArray := make([]interface{}, 0)

					for _, value := range valArray {

						opArray1 := make([]interface{}, 0)
						for _, mapkey := range innerkeys {

							value2 := value.(map[string]interface{})

							if value3, exists := value2[mapkey]; exists {

								opArray1 = append(opArray1, value3)

							} else {

								opArray1 = append(opArray1, nil)

							}
						}

						opArray = append(opArray, opArray1)

					}
					if field.Orientation == types.Horizontal {

						inputMap[field.ExcelName] = opArray
					} else {

						inputMap[field.ExcelName] = utilities.Transpose(opArray)

					}

				default:
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "unknown data type for input json field " + field.JsonName,
					})
					return
				}

			}
		} else {

			outputfieldDataMap[field.ExcelName] = field
			outputFields = append(outputFields, field.ExcelName)

		}
	}

	excelmanager, err := utilities.NewExcelManager(e0001data.ExcelPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get excel manager" + err.Error(),
		})

		return
	}

	defer excelmanager.Close()

	outputMap, err := excelmanager.NamedRangeSetAndGet(inputMap, outputFields)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get results:" + err.Error(),
		})

		return
	}
	formatted_outputmap := make(map[string]interface{})
	for key, field := range outputfieldDataMap {

		switch field.FieldType {
		case types.Single:
			// formatted_outputmap[field.JsonName] = outputMap[key].([][]interface{})[0][0]
			utilities.AddNestedValue(formatted_outputmap, field.JsonName, outputMap[key].([][]interface{})[0][0])
		case types.OneDArray:
			if field.Orientation == types.Horizontal {
				// formatted_outputmap[field.JsonName] = outputMap[key].([][]interface{})[0]
				utilities.AddNestedValue(formatted_outputmap, field.JsonName, outputMap[key].([][]interface{})[0][0])
			} else {
				array := outputMap[key].([][]interface{})
				interfaceSlice := make([]interface{}, len(array))

				// Fill the interface slice
				for i, v := range array {
					interfaceSlice[i] = v[0]
				}
				// formatted_outputmap[field.JsonName] = interfaceSlice
				utilities.AddNestedValue(formatted_outputmap, field.JsonName, interfaceSlice)
			}

		case types.TwoDArray:
			if field.Orientation == types.Horizontal {
				// formatted_outputmap[field.JsonName] = outputMap[key]
				utilities.AddNestedValue(formatted_outputmap, field.JsonName, outputMap[key])

			} else {

				// formatted_outputmap[field.JsonName] = utilities.Transpose1(outputMap[key].([][]interface{}))
				utilities.AddNestedValue(formatted_outputmap, field.JsonName, utilities.Transpose1(outputMap[key].([][]interface{})))
			}
		case types.OneDMap:
			var outerkeys []string

			// Unmarshal JSON string to a slice of strings
			err := json.Unmarshal([]byte(field.OuterKeys), &outerkeys)
			if err != nil {
				fmt.Println("Error:Unable parse the map keys: ", err)
			}

			outputvalMap := make(map[string]interface{})

			if field.Orientation == types.Horizontal {
				for i, val := range outputMap[key].([][]interface{})[0] {
					outputvalMap[outerkeys[i]] = val
				}

			} else {
				for i, val := range outputMap[key].([][]interface{}) {
					outputvalMap[outerkeys[i]] = val[0]
				}

			}
			// formatted_outputmap[field.JsonName] = outputvalMap
			utilities.AddNestedValue(formatted_outputmap, field.JsonName, outputvalMap)

		case types.TwoDMap:
			var outerkeys []string
			var innerkeys []string

			// Unmarshal JSON string to a slice of strings
			err := json.Unmarshal([]byte(field.OuterKeys), &outerkeys)
			if err != nil {
				fmt.Println("Error:Unable parse the map outer keys: ", err)
			}

			err = json.Unmarshal([]byte(field.InnerKeys), &innerkeys)
			if err != nil {
				fmt.Println("Error:Unable parse the map inner keys: ", err)
			}

			var valArray [][]interface{}

			if field.Orientation == types.Horizontal {

				valArray = outputMap[key].([][]interface{})
			} else {

				valArray = utilities.Transpose1(outputMap[key].([][]interface{}))

			}

			outputvalMap := make(map[string]interface{})
			for i, mapkey1 := range outerkeys {
				outputvalMap1 := make(map[string]interface{})
				for j, mapkey2 := range innerkeys {
					outputvalMap1[mapkey2] = valArray[i][j]
				}
				outputvalMap[mapkey1] = outputvalMap1
			}

			// formatted_outputmap[field.JsonName] = outputvalMap
			utilities.AddNestedValue(formatted_outputmap, field.JsonName, outputvalMap)
		case types.TwoDArrayMap:
			var innerkeys []string

			// Unmarshal JSON string to a slice of strings

			err = json.Unmarshal([]byte(field.InnerKeys), &innerkeys)
			if err != nil {
				fmt.Println("Error:Unable parse the map inner keys: ", err)
			}

			var valArray [][]interface{}

			if field.Orientation == types.Horizontal {

				valArray = outputMap[key].([][]interface{})
			} else {

				valArray = utilities.Transpose1(outputMap[key].([][]interface{}))

			}

			outputvalArray := make([]interface{}, 0)
			for _, value := range valArray {
				outputvalMap1 := make(map[string]interface{})
				for j, mapkey2 := range innerkeys {
					outputvalMap1[mapkey2] = value[j]
				}
				outputvalArray = append(outputvalArray, outputvalMap1)

			}

			// formatted_outputmap[field.JsonName] = outputvalArray
			utilities.AddNestedValue(formatted_outputmap, field.JsonName, outputvalArray)

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "unknown data type for json field " + field.JsonName,
			})
			return

		}

	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"Output": formatted_outputmap,
	})
}
