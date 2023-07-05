package utilities

import (
	"github.com/shijith.chand/go-jwt/initializers"
	"github.com/shijith.chand/go-jwt/models"
)

func AppendData(data map[string]interface{}, icoid uint) map[string]interface{} {

	returndata := make(map[string]interface{})
	var getcompany models.Company

	for k, v := range data {
		returndata[k] = v
	}
	initializers.DB.First(&getcompany, "id  = ?", icoid)
	returndata["company"] =
		map[string]interface{}{
			"companyId":   icoid,
			"Companyname": getcompany.CompanyName,
		}

	return returndata
}
