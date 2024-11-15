package excelmanagement

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type ExcelManager struct {
	excelApp  *ole.IDispatch
	workbook  *ole.IDispatch
	workbooks *ole.IDispatch
	names     *ole.IDispatch
}

func NewExcelManager(workbookPath string) (*ExcelManager, error) {
	// Lock the OS thread to ensure all COM operations are on the same thread
	runtime.LockOSThread()

	// Initialize COM
	ole.CoInitialize(0)

	// Create Excel Application
	excelApp, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		ole.CoUninitialize()
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("failed to start Excel: %v", err)
	}

	// Get IDispatch object from Excel Application
	excelAppDispatch, err := excelApp.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		excelApp.Release()
		ole.CoUninitialize()
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("failed to get IDispatch: %v", err)
	}

	// Make Excel invisible (optional)
	oleutil.PutProperty(excelAppDispatch, "Visible", false)

	// Open the workbook
	workbooks := oleutil.MustGetProperty(excelAppDispatch, "Workbooks").ToIDispatch()
	workbook := oleutil.MustCallMethod(workbooks, "Open", workbookPath).ToIDispatch()
	// Access the "Names" collection from the workbook
	names := oleutil.MustGetProperty(workbook, "Names").ToIDispatch()

	return &ExcelManager{
		excelApp:  excelAppDispatch,
		workbooks: workbooks,
		workbook:  workbook,
		names:     names,
	}, nil
}

// 	// Initialize

// GetCellValue retrieves a value from a specific cell in the workbook
func (m *ExcelManager) GetCellValue(sheetName string, cell string) (interface{}, error) {
	sheet := oleutil.MustGetProperty(m.workbook, "Sheets", sheetName).ToIDispatch()
	cellValue := oleutil.MustGetProperty(sheet, "Range", cell).ToIDispatch()
	value := oleutil.MustGetProperty(cellValue, "Value")
	return value.Value(), nil
}

// SetCellValue sets a value to a specific cell in the workbook
func (m *ExcelManager) SetCellValue(sheetName string, cell string, value interface{}) error {
	sheet := oleutil.MustGetProperty(m.workbook, "Sheets", sheetName).ToIDispatch()
	cellRef := oleutil.MustGetProperty(sheet, "Range", cell).ToIDispatch()
	_, err := oleutil.PutProperty(cellRef, "Value", value)
	return err
}

func (m *ExcelManager) Close() {
	// Ensure cleanup is done on the same thread
	defer runtime.UnlockOSThread()

	oleutil.CallMethod(m.workbook, "Close", false) // Close workbook without saving
	oleutil.CallMethod(m.excelApp, "Quit")         // Quit Excel application

	// Release the COM objects
	m.names.Release()
	m.workbook.Release()
	m.workbooks.Release()
	m.excelApp.Release()

	// Uninitialize COM
	ole.CoUninitialize()
}

func (m *ExcelManager) NamedRangeSetAndGet(input map[string]interface{}, output []interface{}) (map[string]interface{}, error) {
	//process input
	for key, value := range input {
		namedRange := oleutil.MustCallMethod(m.names, "Item", key).ToIDispatch()
		rangeObj := oleutil.MustGetProperty(namedRange, "RefersToRange").ToIDispatch()
		// Get number of rows and columns in the named range
		rows := oleutil.MustGetProperty(rangeObj, "Rows").ToIDispatch()
		columns := oleutil.MustGetProperty(rangeObj, "Columns").ToIDispatch()
		rowCount := int(oleutil.MustGetProperty(rows, "Count").Val)       // Get number of rows
		columnCount := int(oleutil.MustGetProperty(columns, "Count").Val) // Get number of columns
		input2dArray := make([][]interface{}, 0)
		for _, val := range value.([]interface{}) {

			input2dArray = append(input2dArray, val.([]interface{}))

		}
		inputRowCount := len(input2dArray)
		inputColCount := len(input2dArray[0])
		if inputRowCount > rowCount || inputColCount > columnCount {
			return nil, errors.New("input size mismatch with named range for " + key)
		}

		for i := 1; i <= inputRowCount; i++ {
			//skip some elements if no value provided
			if input2dArray[i-1] == nil {
				continue
			}
			for j := 1; j <= inputColCount; j++ {
				//skip some elements if no value provided
				if input2dArray[i-1][j-1] == nil {
					continue
				}
				cell := oleutil.MustGetProperty(rangeObj, "Cells", i, j).ToIDispatch()
				oleutil.MustPutProperty(cell, "Value", input2dArray[i-1][j-1])
				cell.Release()
			}
		}
		rows.Release()
		columns.Release()
		namedRange.Release()
		rangeObj.Release()
	}
	outputMap := make(map[string]interface{})
	for _, key := range output {
		namedRange := oleutil.MustCallMethod(m.names, "Item", key).ToIDispatch()
		rangeObj := oleutil.MustGetProperty(namedRange, "RefersToRange").ToIDispatch()

		// Get number of rows and columns in the named range
		rows := oleutil.MustGetProperty(rangeObj, "Rows").ToIDispatch()
		columns := oleutil.MustGetProperty(rangeObj, "Columns").ToIDispatch()

		rowCount := int(oleutil.MustGetProperty(rows, "Count").Val)       // Get number of rows
		columnCount := int(oleutil.MustGetProperty(columns, "Count").Val) // Get number of columns
		valueArray := make([][]interface{}, rowCount)
		for i := 1; i <= rowCount; i++ {
			valueArray[i-1] = make([]interface{}, columnCount)
			for j := 1; j <= columnCount; j++ {
				cell := oleutil.MustGetProperty(rangeObj, "Cells", i, j).ToIDispatch()
				valueArray[i-1][j-1] = oleutil.MustGetProperty(cell, "Value").Value()
				cell.Release()
			}
		}

		outputMap[key.(string)] = valueArray

		rows.Release()
		columns.Release()
		namedRange.Release()
		rangeObj.Release()
	}

	return outputMap, nil
}

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
