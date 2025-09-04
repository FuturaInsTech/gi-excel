package excelparamTypes

import (
	"encoding/json"
	"fmt"

	"github.com/FuturaInsTech/gi-excel/exceltypes"
)

type Extradata interface {
	// Methods
	ParseData(map[string]interface{})
	GetFormattedData(datamap map[string]string) map[string]interface{}
}

type E0001Data struct {
	ExcelPath      string
	TemplatePath   string
	InputExcelPath string
	InputSheet     string
	OutputSheet    string
	ImagePath      string
	PdfPath        string
}

func (m *E0001Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *E0001Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type E0002Data struct {
	FieldArray []E0002
}
type E0002 struct {
	JsonName    string
	ExcelName   string
	FieldMode   exceltypes.ServiceFieldMode
	FieldType   exceltypes.ServiceFieldType
	OuterKeys   string
	InnerKeys   string
	Mandatory   bool
	Orientation exceltypes.ServiceFieldOrient
}

func (m *E0002Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *E0002Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}
