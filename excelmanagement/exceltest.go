package excelmanagement

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FuturaInsTech/gi-excel/exceltypes"
	"github.com/FuturaInsTech/gi-excel/paramTypes"
)

func ExcelProcessor(serviceName string, requestMap map[string]interface{}, e0001data paramTypes.E0001Data, e0002data paramTypes.E0002Data) (map[string]interface{}, error) {
	outputfieldDataMap := make(map[string]paramTypes.E0002)
	outputFields := make([]interface{}, 0)
	inputMap := make(map[string]interface{})

	for _, field := range e0002data.FieldArray {
		fmt.Println(field.Mandatory, "**********************")
		if field.FieldMode == exceltypes.Input {
			// val, ok := requestMap[field.JsonName]
			val, err := GetNestedValue(requestMap, field.JsonName)
			// If the key does not exist error
			if err != nil {
				if field.Mandatory {
					return nil, err
				}
				// If the key exist proceed
			} else {

				switch field.FieldType {
				case exceltypes.Single:
					inputMap[field.ExcelName] = []interface{}{[]interface{}{val}}
				case exceltypes.OneDArray:
					if field.Orientation == exceltypes.Horizontal {
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

				case exceltypes.TwoDArray:
					if field.Orientation == exceltypes.Horizontal {
						inputMap[field.ExcelName] = val

					} else {

						inputMap[field.ExcelName] = Transpose(val.([]interface{}))
					}
				case exceltypes.OneDMap:
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
					if field.Orientation == exceltypes.Horizontal {

						inputMap[field.ExcelName] = []interface{}{valArray}
					} else {
						interfaceSlice := make([]interface{}, len(valArray))

						for i, v := range valArray {
							interfaceSlice[i] = []interface{}{v}
						}
						inputMap[field.ExcelName] = interfaceSlice

					}

				case exceltypes.TwoDMap:
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
					if field.Orientation == exceltypes.Horizontal {

						inputMap[field.ExcelName] = valArray
					} else {

						inputMap[field.ExcelName] = Transpose(valArray)

					}

				case exceltypes.TwoDArrayMap:
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
					if field.Orientation == exceltypes.Horizontal {

						inputMap[field.ExcelName] = opArray
					} else {

						inputMap[field.ExcelName] = Transpose(opArray)

					}

				default:
					iErrDesc := "unknown data type for input json field " + field.JsonName
					return nil, errors.New(iErrDesc)
				}

			}
		} else {

			outputfieldDataMap[field.ExcelName] = field
			outputFields = append(outputFields, field.ExcelName)

		}
	}

	excelmanager, err := NewExcelManager(e0001data.ExcelPath)
	if err != nil {
		return nil, err
	}

	defer excelmanager.Close()

	outputMap, err := excelmanager.NamedRangeSetAndGet(inputMap, outputFields)
	if err != nil {
		return nil, err
	}
	formatted_outputmap := make(map[string]interface{})
	for key, field := range outputfieldDataMap {

		switch field.FieldType {
		case exceltypes.Single:
			// formatted_outputmap[field.JsonName] = outputMap[key].([][]interface{})[0][0]
			AddNestedValue(formatted_outputmap, field.JsonName, outputMap[key].([][]interface{})[0][0])
		case exceltypes.OneDArray:
			if field.Orientation == exceltypes.Horizontal {
				// formatted_outputmap[field.JsonName] = outputMap[key].([][]interface{})[0]
				AddNestedValue(formatted_outputmap, field.JsonName, outputMap[key].([][]interface{})[0][0])
			} else {
				array := outputMap[key].([][]interface{})
				interfaceSlice := make([]interface{}, len(array))

				// Fill the interface slice
				for i, v := range array {
					interfaceSlice[i] = v[0]
				}
				// formatted_outputmap[field.JsonName] = interfaceSlice
				AddNestedValue(formatted_outputmap, field.JsonName, interfaceSlice)
			}

		case exceltypes.TwoDArray:
			if field.Orientation == exceltypes.Horizontal {
				// formatted_outputmap[field.JsonName] = outputMap[key]
				AddNestedValue(formatted_outputmap, field.JsonName, outputMap[key])

			} else {

				// formatted_outputmap[field.JsonName] = Transpose1(outputMap[key].([][]interface{}))
				AddNestedValue(formatted_outputmap, field.JsonName, Transpose1(outputMap[key].([][]interface{})))
			}
		case exceltypes.OneDMap:
			var outerkeys []string

			// Unmarshal JSON string to a slice of strings
			err := json.Unmarshal([]byte(field.OuterKeys), &outerkeys)
			if err != nil {
				fmt.Println("Error:Unable parse the map keys: ", err)
			}

			outputvalMap := make(map[string]interface{})

			if field.Orientation == exceltypes.Horizontal {
				for i, val := range outputMap[key].([][]interface{})[0] {
					outputvalMap[outerkeys[i]] = val
				}

			} else {
				for i, val := range outputMap[key].([][]interface{}) {
					outputvalMap[outerkeys[i]] = val[0]
				}

			}
			// formatted_outputmap[field.JsonName] = outputvalMap
			AddNestedValue(formatted_outputmap, field.JsonName, outputvalMap)

		case exceltypes.TwoDMap:
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

			if field.Orientation == exceltypes.Horizontal {

				valArray = outputMap[key].([][]interface{})
			} else {

				valArray = Transpose1(outputMap[key].([][]interface{}))

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
			AddNestedValue(formatted_outputmap, field.JsonName, outputvalMap)
		case exceltypes.TwoDArrayMap:
			var innerkeys []string

			// Unmarshal JSON string to a slice of strings

			err = json.Unmarshal([]byte(field.InnerKeys), &innerkeys)
			if err != nil {
				fmt.Println("Error:Unable parse the map inner keys: ", err)
			}

			var valArray [][]interface{}

			if field.Orientation == exceltypes.Horizontal {

				valArray = outputMap[key].([][]interface{})
			} else {

				valArray = Transpose1(outputMap[key].([][]interface{}))

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
			AddNestedValue(formatted_outputmap, field.JsonName, outputvalArray)

		default:
			iErrDesc := "unknown data type for input json field " + field.JsonName
			return nil, errors.New(iErrDesc)

		}

	}

	return formatted_outputmap, nil
}
