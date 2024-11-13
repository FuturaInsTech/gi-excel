package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
)

func CompareDate(fromdate, todate string, language uint) error {

	var getError models.Error

	fromdateint, err := strconv.ParseUint(fromdate, 10, 64)
	if err != nil {
		panic(err)
	}
	todateint, err := strconv.ParseUint(todate, 10, 64)
	if err != nil {
		result := initializers.DB.Select("long_code").Where("short_code = ? AND language_id = ?", "E0001", language).Find(&getError)

		if result.RowsAffected == 0 {
			err1 := errors.New("error code not found1 ")
			return err1
		}
	}

	if fromdateint > todateint {
		fmt.Println(fromdateint)
		fmt.Println(todateint)
		var longcode string

		//result = initializers.DB.Where("bank_code LIKE ?", "%"+isearch+"%").Find(&getallbank)
		result := initializers.DB.Select("long_code").Where("short_code = ? AND language_id = ?", "E0001", language).Find(&getError)

		if result.RowsAffected == 0 {
			err1 := errors.New("error code not found1 ")
			return err1
		}
		result.Scan(&longcode)

		var err1 = errors.New(longcode)
		return err1
	}
	return nil
	//return output1, output2

}

func DateBlank(fromdate string, language uint) error {

	var getError models.Error

	if fromdate == "" {
		fmt.Println("i am inside zero")
		var longcode string

		result := initializers.DB.Select("long_code").Where("short_code = ? AND language_id = ?", "E0002", language).Find(&getError)

		if result.RowsAffected == 0 {
			err1 := errors.New("error code not found1 ")
			return err1
		}
		result.Scan(&longcode)

		var err1 = errors.New(longcode)
		return err1
	}
	return nil
	//return output1, output2

}
func DateZero(fromdate string, language uint) error {

	var getError models.Error
	fromdateint, err := strconv.ParseUint(fromdate, 10, 64)
	if err != nil {
		panic(err)
	}

	if fromdateint == 0 {
		fmt.Println("i am inside zero")
		var longcode string

		result := initializers.DB.Select("long_code").Where("short_code = ? AND language_id = ?", "E0002", language).Find(&getError)

		if result.RowsAffected == 0 {
			err1 := errors.New("error code not found1 ")
			return err1
		}
		result.Scan(&longcode)

		var err1 = errors.New(longcode)
		return err1
	}
	return nil
	//return output1, output2

}

func GetError(CompanyId uint, LanguageId uint, ErrorCode string) string {
	var ierror models.Error
	errordesc := "Error Code Not Found"
	fmt.Println("Error Code inside Validators  2", CompanyId, LanguageId, ErrorCode)

	result := initializers.DB.Select("long_code").Where("company_id = ? AND language_id = ? and short_code = ? ",
		CompanyId, LanguageId, ErrorCode).Find(&ierror)
	if result.Error == nil {
		result.Scan(&errordesc)

	}
	return errordesc
}

type DbValError struct {
	DbErrors []interface{}
}

func validateError1(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var errors1 models.Error
	iCompany := errors1.CompanyID
	iLanguage := errors1.LanguageID
	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &errors1); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	if errors1.LongCode == "" {
		errorDescription := GetError(iCompany, iLanguage, "E0018")
		fielderror := map[string]interface{}{
			"LongCode": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if errors1.ShortCode == "" {
		errorDescription := GetError(iCompany, iLanguage, "E0018")
		fielderror := map[string]interface{}{
			"ShortCode": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if len(fieldvalerrors) > 0 {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Field Errors"), dbvalerror
	} else {
		return nil, dbvalerror
	}

}

func validatePermission(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var permission models.Permission
	iCompany := permission.CompanyID
	fmt.Println("Company Code ", iCompany)
	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &permission); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	if permission.Method == "" {
		errorDescription := GetError(1, 1, "E0018")
		fielderror := map[string]interface{}{
			"Method": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if permission.ModelName == "" {
		errorDescription := GetError(1, 1, "E0018")
		fielderror := map[string]interface{}{
			"ModelName": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if len(fieldvalerrors) > 0 {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Field Errors"), dbvalerror
	} else {
		return nil, dbvalerror
	}

}

func validateUserGroup(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var usergroup models.UserGroup
	iCompany := usergroup.CompanyID

	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &usergroup); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	if usergroup.GroupName == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"GroupName": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if usergroup.ValidFrom > usergroup.ValidTo {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"ValidFrom": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if len(fieldvalerrors) > 0 {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Field Errors"), dbvalerror
	} else {
		return nil, dbvalerror
	}

}

func ValidateData(dbdata map[string]interface{}, modelname string) (error, DbValError) {

	if modelname == "Error" {
		err, dbvalerror := validateError1(dbdata)
		return err, dbvalerror
	}
	if modelname == "Permission" {
		err, dbvalerror := validatePermission(dbdata)
		return err, dbvalerror
	}
	// Add New Models to Validate
	return errors.New("Model Not Found"), DbValError{}

}
