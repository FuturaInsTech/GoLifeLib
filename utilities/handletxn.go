package utilities

import (
	"_/C_/Go/GoLifeLib/utilities"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func HandleTxn(c *gin.Context, txn *gorm.DB, method string, txnErr *models.TxnError, userco, userlan *uint) {
	if r := recover(); r != nil {
		txn.Rollback()

		var errMsg string
		switch x := r.(type) {
		case error:
			errMsg = x.Error()
		case string:
			errMsg = x
		default:
			errMsg = fmt.Sprintf("%v", x)
		}

		fmt.Println("PANIC TRACE:", string(debug.Stack()))

		// Map to APERR
		longDesc, _ := utilities.GetErrorDesc(*userco, *userlan, "APERR")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("APERR : %s - %s (in %s)", longDesc, errMsg, method),
		})
		return
	}

	if txnErr.ErrorCode != "" || txn.Error != nil {
		txn.Rollback()
		code := txnErr.ErrorCode
		if code == "" {
			code = "GL999"
		}

		var errMsg string
		if code == "PERME" {
			errMsg = fmt.Sprintf("%s : Access denied. You are not authorized to perform '%s'.", code, method)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMsg,
			})
			return
		}
		if code == "PARME" {
			longDesc, _ := utilities.GetErrorDesc(*userco, *userlan, code)

			errMsg = fmt.Sprintf(
				"%s : %s - %s | ParamName: %s | ParamItem: %s",
				code,
				longDesc,
				method,
				txnErr.ParamName,
				txnErr.ParamItem,
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMsg,
			})
			return
		}
		if code == "DBERR" {

			longDesc, _ := utilities.GetErrorDesc(*userco, *userlan, code)

			dbErrMsg := txnErr.DbError.Error()
			cleanErr := dbErrMsg
			if strings.HasPrefix(dbErrMsg, "Error") {
				cleanErr = strings.TrimPrefix(dbErrMsg, "Error ")
			}
			// Final top-level error
			errMsg = fmt.Sprintf("%s : %s (in %s) - %s",
				code, longDesc, method, cleanErr)

			errMsg = fmt.Sprintf("%s : %s (in %s) - %s",
				code, longDesc, method, cleanErr)

			c.JSON(http.StatusBadRequest, gin.H{
				"error": errMsg,
			})
			return

		}

		if txnErr.DbError != nil {
			longDesc, _ := utilities.GetErrorDesc(*userco, *userlan, code)
			dbErrMsg := txnErr.DbError.Error()
			cleanErr := dbErrMsg
			if strings.HasPrefix(dbErrMsg, "Error") {
				cleanErr = strings.TrimPrefix(dbErrMsg, "Error ")
			}
			// Final top-level error
			errMsg = fmt.Sprintf("%s : %s (in %s) - %s",
				code, longDesc, method, cleanErr)

			c.JSON(http.StatusBadRequest, gin.H{
				"error": errMsg,
			})
			return
		} else {
			longDesc, _ := utilities.GetErrorDesc(*userco, *userlan, code)
			errMsg = fmt.Sprintf("%s : %s - %s", code, longDesc, method)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errMsg,
			})
			return
		}

	}

	txn.Commit()
}
