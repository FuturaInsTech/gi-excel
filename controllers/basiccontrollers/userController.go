package basiccontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/FuturaInsTech/GoExcel/utilities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context) {
	var body struct {
		Phone           string
		Password        string
		Name            string
		Email           string
		Language        uint
		Profile         string `gorm:"type:longtext"`
		VerficationCode string `gorm:"type:varchar(10)"`
		UserStatusID    uint
		LanguageId      uint
		CompanyId       uint
		UserGroupId     uint
		Gender          string `gorm:"type:varchar(10)"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	now_time := time.Now()

	fmt.Println(body.Name)
	// IsGenderValid := false
	// if body.Gender != "" {
	// 	IsGenderValid = utilities.ParamValidator(body.CompanyId, body.LanguageId, "P0002", body.Gender)

	// }
	// if !IsGenderValid {
	// 	errorDescription := utilities.GetError(body.CompanyId, body.LanguageId, "E0005")

	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": errorDescription,
	// 	})
	// 	return
	// }

	user := models.User{Name: body.Name, Email: body.Email,
		Is_valid: true, Password: string(hash), Phone: body.Phone,
		Last_logged_in_datime: now_time, LanguageID: body.LanguageId, CompanyID: body.CompanyId,
		UserStatusID: body.UserStatusID, UserGroupID: body.UserGroupId, Gender: body.Gender}
	result := initializers.DB.Create(&user)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to insert",
		})

		return

	}

	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	var body struct {
		Phone    string
		Password string
		Channel  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	webLogin := true

	if body.Channel == "app" {
		webLogin = false
	}

	var user models.User
	initializers.DB.First(&user, "phone = ?", body.Phone)

	if user.Id == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user don't exist",
		})

		return

	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": body.Phone,
		"aud": c.ClientIP(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed tokenize",
		})

		return

	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.Id,
		"aud": c.ClientIP(),
		"exp": time.Now().Add(time.Hour * 30 * 24).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_REFRESH_TOKEN")))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed tokenize",
		})

		return

	}

	if webLogin {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)
		c.SetCookie("RefreshToken", refreshTokenString, 30*24*3600, "/api/v1/auth/refresh", "", false, true)
	}

	var message struct {
		Id              uint64 `json:"id"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		Phone           string `json:"phone"`
		UserGroupId     uint   `json:"userGroupId"`
		CompanyId       uint   `json:"companyId"`
		LanguageId      uint   `json:"languageId"`
		AuthToken       string `json:"authToken"`
		RefreshToken    string `json:"refreshToken"`
		CompanyLongName string `json:"companyName"`
		UserGroupName   string `json:"userGroupName"`
	}

	var companyenq models.Company
	initializers.DB.Find(&companyenq, "id = ?", user.CompanyID)
	var usergroupenq models.UserGroup
	initializers.DB.Find(&usergroupenq, "id = ?", user.UserGroupID)

	message.Id = user.Id
	message.Email = user.Email
	message.Phone = user.Phone
	message.UserGroupId = user.UserGroupID
	message.Name = user.Name
	message.CompanyId = user.CompanyID
	message.LanguageId = user.LanguageID
	message.UserGroupName = usergroupenq.GroupName
	message.CompanyLongName = companyenq.CompanyName
	if !webLogin {
		message.AuthToken = tokenString
		message.RefreshToken = refreshTokenString
	}
	//c.JSON(http.StatusOK, gin.H{})
	fmt.Println(tokenString + "authToken=======================================")
	c.JSON(http.StatusOK, gin.H{"message": message})

}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": user})

}

func GetAllUsers(c *gin.Context) {

	user, _ := c.Get("user")
	method := "GetAllUsers"
	//var userdatamap map[string]interface{}
	//_, err := utilities.GetUserAccess(user, method)
	userdatamap, err := utilities.GetUserAccess(user, method)
	userco := uint(userdatamap["CompanyId"].(float64))
	userlan := uint(userdatamap["LanguageId"].(float64))

	if err != nil {
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, "GL001")
		fmt.Println("Errors")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "GL001" + " : " + longDesc + "-" + method,
		})

		return
	}
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

	var getalluser []models.User
	//userco := userdatamap["CompanyId"]

	var result *gorm.DB

	if searchpagination.SearchString != "" && searchpagination.SearchCriteria != "" {
		result = initializers.DB.Model(&models.User{}).Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.User{}).
			Where(searchpagination.SearchCriteria+" LIKE ? AND company_id = ?", "%"+searchpagination.SearchString+"%", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getalluser)

	} else {
		fmt.Println("No Selection ")
		fmt.Println(searchpagination.SearchCriteria)
		fmt.Println(searchpagination.SearchString)
		result = initializers.DB.Model(&models.User{}).Where("company_id = ?", userco).Count(&totalRecords)
		result = initializers.DB.Model(&models.User{}).
			Where("company_id = ?", userco).
			Order(searchpagination.SortColumn + " " + searchpagination.SortDirection).
			Limit(searchpagination.PageSize).Offset(searchpagination.Offset).
			Find(&getalluser)
	}

	// if result is null, then give an language ..
	if result.Error != nil {
		longDesc, _ := utilities.GetErrorDesc(userco, userlan, "GL120")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "GL120" + " : " + longDesc,
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
		fieldMappings := [3]map[string]string{{
			"displayName": "User Name",
			"fieldName":   "name",
			"dataType":    "string"},
			{"displayName": "Email",
				"fieldName": "email",
				"dataType":  "string"},
			{"displayName": "Gender Code",
				"fieldName": "gender",
				"dataType":  "string"},
		}

		c.JSON(200, gin.H{

			"All Users":      getalluser,
			"Field Map":      fieldMappings,
			"paginationData": paginationData,
		})

	} else {
		c.JSON(200, gin.H{

			"All Users":      getalluser,
			"paginationData": paginationData,
		})
	}

}

//Delete Function

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var cont models.User
	result := initializers.DB.First(&cont, "id  = ?", id)
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get :" + result.Error.Error(),
		})

		return

	}
}

func DeleteUser(c *gin.Context) {
	// http parameter get value from postman and storing it in local as shijithid
	shijithid := c.Param("rangaid")
	// copy the structure of users and store it in wroking storage as wuser
	var wuser models.User
	// get result to check whether it is successful or null. if null give error
	// this is nothing but select statement
	result := initializers.DB.First(&wuser, "id  = ?", shijithid)
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get :" + result.Error.Error(),
		})

		return

	}

	result = initializers.DB.Delete(&wuser)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete :" + result.Error.Error(),
		})

		return

	}

	c.JSON(http.StatusOK, "userid "+shijithid+" is deleted")

}

func ModifyUser(c *gin.Context) {
	// mapping json to sourceMap

	var sourceMap map[string]interface{}
	fmt.Println(sourceMap)
	if c.Bind(&sourceMap) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user",
		})

		return

	}

	var uuser models.User
	fmt.Println((uuser))
	result := initializers.DB.First(&uuser, "id  = ?", sourceMap["organgeid"])
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get :" + result.Error.Error(),
		})

		return

	}

	var targetMap map[string]interface{}
	fmt.Println((targetMap))
	data, _ := json.Marshal(uuser)

	json.Unmarshal(data, &targetMap)

	for key, _ := range targetMap {

		if val1, ok := sourceMap[key]; ok {
			targetMap[key] = val1
		}

	}

	data, _ = json.Marshal(targetMap)
	json.Unmarshal(data, &uuser)

	result = initializers.DB.Save(&uuser)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update :" + result.Error.Error(),
		})

		return

	}
	fmt.Println((uuser.Name))
	c.JSON(http.StatusOK, gin.H{"outputs": uuser})

}

func Refresh(c *gin.Context) {
	//var token_expry = 10

	var body struct {
		Channel string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	webLogin := true

	if body.Channel == "app" {
		webLogin = false
	}

	data, _ := c.Get("user")
	var phone string
	if val, ok := data.(string); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse data",
		})

		return

	} else {

		phone = val
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": phone,
		"aud": c.ClientIP(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed tokenize",
		})

		return

	}

	if webLogin {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(200, gin.H{

			"authToken": tokenString,
		})
	}

}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.SetCookie("RefreshToken", "", -1, "/api/v1/auth/refresh", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "user is signed out"})
}
