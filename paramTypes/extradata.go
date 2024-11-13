package paramTypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/types"
)

type Extradata interface {
	// Methods
	ParseData(map[string]interface{})
	GetFormattedData(datamap map[string]string) map[string]interface{}
}

func GetParamDesc(iCompany uint, iParam string, iItem string, iLanguage uint) (string, string, error) {
	var paramdesc models.ParamDesc

	results := initializers.DB.Where("company_id = ? AND name = ? and item = ? and language_id = ?", iCompany, iParam, iItem, iLanguage).Find(&paramdesc)
	if results.Error != nil || results.RowsAffected == 0 {

		return "", "", errors.New(" -" + strconv.FormatUint(uint64(iCompany), 10) + "-" + iParam + "-" + "-" + iItem + "-" + strconv.FormatUint(uint64(iLanguage), 10) + "-" + " is missing")
		//return errors.New(results.Error.Error())
	}
	return paramdesc.Shortdesc, paramdesc.Longdesc, nil
}

type E0001Data struct {
	ExcelPath string
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
	FieldMode   types.ServiceFieldMode
	FieldType   types.ServiceFieldType
	OuterKeys   string
	InnerKeys   string
	Mandatory   bool
	Orientation types.ServiceFieldOrient
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
