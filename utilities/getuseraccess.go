package utilities

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
)

func GetUserAccess(iuser interface{}, imethod string) (map[string]interface{}, error) {

	//user, _ := c.Get("user")
	var sourceMap map[string]interface{}
	// converting an entity(Model) to Json
	data1, _ := json.Marshal(iuser)
	//converting Json to Source Map
	json.Unmarshal(data1, &sourceMap)

	iid := sourceMap["Id"]
	iusergroupid := sourceMap["UserGroupId"]
	// fmt.Println("Ranga - Input User")
	// fmt.Println(iid)
	// fmt.Println(imethod)
	// fmt.Println("Input UserGroup")
	// fmt.Println(iusergroupid)
	ico := sourceMap["CompanyId"]
	fmt.Println("ICO Value ", ico)
	//?? Shijith - how to get ID
	var iCompany models.Company

	var getpermission models.Permission

	initializers.DB.Where("id", ico).Find(&iCompany)
	sourceMap["companyName"] = iCompany.CompanyName
	//db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
	result := initializers.DB.Where("method  = ? AND user_group_ID = ?", imethod, iusergroupid).Find(&getpermission)
	//fmt.Println(result.RowsAffected)
	fmt.Println("RAnga1")
	fmt.Println(result)
	fmt.Println(getpermission)
	if result.Error != nil {
		//result := initializers.DB.Where("method  = ? AND user_id = ?", imethod, iid).Find(&getpermission)
		//if result.Error != nil {
		fmt.Println("RAnga2")
		fmt.Println(result)
		/*fmt.Println(getpermission)*/
		return sourceMap, errors.New(result.Error.Error())
		//	}

	}
	if result.RowsAffected == 0 {

		result := initializers.DB.Where("method  = ? AND user_ID = ?", imethod, iid).Find(&getpermission)

		if result.Error != nil {
			fmt.Println("RAnga2")
			fmt.Println(result)
			/*fmt.Println(getpermission)*/
			return sourceMap, errors.New(result.Error.Error())
		}
		if result.RowsAffected == 0 {
			fmt.Println("RAnga3")
			fmt.Println(result)
			fmt.Println(getpermission)
			return sourceMap, errors.New("User Not Authorized to this Function - " + imethod)
		}

	}
	//var err error
	return sourceMap, nil
}

func GetUserAccessNew(iuser interface{}, imethod string) (map[string]interface{}, error) {

	//user, _ := c.Get("user")
	var sourceMap map[string]interface{}
	// converting an entity(Model) to Json
	data1, _ := json.Marshal(iuser)
	//converting Json to Source Map
	json.Unmarshal(data1, &sourceMap)

	iusergroupid := sourceMap["UserGroupId"]
	// fmt.Println("Ranga - Input User")
	// fmt.Println(iid)
	// fmt.Println(imethod)
	// fmt.Println("Input UserGroup")
	fmt.Println(iusergroupid)
	ico := sourceMap["CompanyId"]
	fmt.Println("ICO Value ", ico)
	//?? Shijith - how to get ID
	var iCompany models.Company
	initializers.DB.Where("id", ico).Find(&iCompany)
	sourceMap["companyName"] = iCompany.CompanyName

	//var err error
	return sourceMap, nil
}
