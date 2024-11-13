package initializers

import "github.com/FuturaInsTech/GoExcel/models"

// Table Name should start with Capital Letter
func SyncDatabase() {

	DB.AutoMigrate(&models.Currency{})
	DB.AutoMigrate(&models.CompanyStatus{})
	DB.AutoMigrate(&models.Company{})
	DB.AutoMigrate(&models.Transaction{})
	DB.AutoMigrate(&models.Permission{})
	DB.AutoMigrate(&models.Language{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Error{})
	DB.AutoMigrate(&models.UserGroup{})
	DB.AutoMigrate(&models.BusinessDate{})
	DB.AutoMigrate(&models.UserStatus{})
	DB.AutoMigrate(&models.Param{})
	DB.AutoMigrate(&models.ParamDesc{})

}
