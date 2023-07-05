package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/shijith.chand/go-jwt/initializers"
	"github.com/shijith.chand/go-jwt/models"
	"github.com/shijith.chand/go-jwt/models/quotation"
	"github.com/shijith.chand/go-jwt/types"
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

func ParamValidator(CompanyId uint, LanguageId uint, ParamRule string, ParamItem string) bool {
	var iparam models.Param

	result := initializers.DB.Select("data").Where("company_id = ? AND name =? ", CompanyId, ParamRule).Find(&iparam)
	if result.Error != nil {
		return false

	} else {
		return true
	}

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

// ITDMIO  itemcoy = 2, itemtabl = 'T5687'  itemitem = '091R'  ITMFRM = 20220101
func GetParamD(CompanyID uint, ParamCode string, ParamItem string, EffectiveDate string) (map[string]interface{}, error) {
	var iparamd models.Param
	// WITEMREC    FROM ITEMREC
	var paramdata map[string]interface{}
	// WT5687REC   FROM T5687REC

	result := initializers.DB.Select("data").Where("company_id = ? AND name = ? AND ITEM = ? AND rec_type = ? and start_date >= ? and end_date <= ? and is_valid = ?", CompanyID, ParamCode, ParamItem, "IT", EffectiveDate, EffectiveDate, "1").Find(&iparamd)

	if result.RowsAffected == 0 {
		var stringV string = strconv.FormatUint(uint64(CompanyID), 10)
		error := errors.New(stringV + ParamCode + ParamItem + EffectiveDate + "Parameter Not Found ")
		return nil, error
	}
	result.Scan(&paramdata)
	return paramdata, nil
}

//ITEMIO

func GetParam(CompanyID uint, ParamCode string, ParamItem string) (map[string]interface{}, error) {
	var iparamd models.Param
	var paramdata map[string]interface{}

	result := initializers.DB.Select("data").Where("company_id = ? AND name = ? AND ITEM = ? AND rec_type = ? and is_valid = ?", CompanyID, ParamCode, ParamItem, "IT", "1").Find(&iparamd)
	// select data from param where company id = '1' and name = 'P0004' AND ITEM = ? AND  REC_TYPE  = 'IT' AND IS_VALID = '1'

	if result.RowsAffected == 0 {
		var stringV string = strconv.FormatUint(uint64(CompanyID), 10)
		error := errors.New(stringV + ParamCode + ParamItem + "Parameter Not Found ")
		return nil, error
	}
	result.Scan(&paramdata)
	return paramdata, nil
}

func CheckParam(CompanyID uint, ParamCode string, ParamItem string) error {
	var iparamd models.Param
	fmt.Println(CompanyID, ParamCode, ParamItem)

	result := initializers.DB.Select("data").Where("company_id = ? AND name = ? AND ITEM = ? AND rec_type = ? and is_valid = ?", CompanyID, ParamCode, ParamItem, "IT", "1").Find(&iparamd)
	// select data from param where company id = '1' and name = 'P0004' AND ITEM = ? AND  REC_TYPE  = 'IT' AND IS_VALID = '1'

	if result.RowsAffected == 0 {
		var stringV string = strconv.FormatUint(uint64(CompanyID), 10)
		error := errors.New(stringV + ParamCode + ParamItem + "Parameter Not Found ")
		return error
	}

	return nil
}

type DbValError struct {
	DbErrors []interface{}
}

func ValidateData(dbdata map[string]interface{}, modelname string) (error, DbValError) {

	if modelname == "QDetail" {
		err, dbvalerror := validateQDetail(dbdata)
		return err, dbvalerror
	}
	if modelname == "DeathH" {
		err, dbvalerror := validateQDetail(dbdata)
		return err, dbvalerror
	}
	if modelname == "Error" {
		err, dbvalerror := validateError1(dbdata)
		return err, dbvalerror
	}
	if modelname == "Permission" {
		err, dbvalerror := validatePermission(dbdata)
		return err, dbvalerror
	}
	if modelname == "QHeader" {
		err, dbvalerror := validateQHeader(dbdata)
		return err, dbvalerror
	}
	if modelname == "Agency" {
		err, dbvalerror := validateAgency(dbdata)
		return err, dbvalerror
	}
	if modelname == "Proposal" {
		err, dbvalerror := validateProposal(dbdata)
		return err, dbvalerror
	}
	if modelname == "Benefits" {
		err, dbvalerror := ValidateBenefit(dbdata)
		return err, dbvalerror
	}

	if modelname == "Policy" {
		err, dbvalerror := ValidatePolicy(dbdata)
		return err, dbvalerror
	}
	// Add New Models to Validate
	return errors.New("Model Not Found"), DbValError{}

}

// QDetail Validatation
func validateQDetail(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var qdetail quotation.QDetail

	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &qdetail); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	iCompany := qdetail.CompanyID
	iCoverage := qdetail.QCoverage
	iDate := qdetail.QDate

	var q0006data types.Q0006Data
	var extradataq0006 types.Extradata = &q0006data
	err = GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}
	fmt.Println("yukesh SA ", qdetail.QSumAssured, q0006data.MinSA)
	if uint(qdetail.QSumAssured) < q0006data.MinSA {
		errorDescription := GetError(iCompany, 1, "E0021")
		fielderror := map[string]interface{}{
			"QSumAssured": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if qdetail.QAge > q0006data.MaxRiskCessAge {
		errorDescription := GetError(iCompany, 1, "E0012")
		fielderror := map[string]interface{}{
			"QAge": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)
	}
	if qdetail.QAge < q0006data.MinRiskCessAge {
		errorDescription := GetError(iCompany, 1, "E0013")
		fielderror := map[string]interface{}{
			"QAge": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)
	}

	// if qdetail.QRiskCessTerm > q0006data.MaxRiskCessT {
	// 	errorDescription := GetError(iCompany, 1, "E0014")
	// 	fielderror := map[string]interface{}{
	// 		"QRiskCessTerm": errorDescription,
	// 	}
	// 	fieldvalerrors = append(fieldvalerrors, fielderror)
	// }
	// if qdetail.QRiskCessTerm < q0006data.MinTerm {
	// 	errorDescription := GetError(iCompany, 1, "E0015")
	// 	fielderror := map[string]interface{}{
	// 		"QRiskCessTerm": errorDescription,
	// 	}

	// 	fieldvalerrors = append(fieldvalerrors, fielderror)
	// }
	// if qdetail.QPremCessTerm > q0006data.MaxP {
	// 	errorDescription := GetError(iCompany, 1, "E0016")
	// 	fielderror := map[string]interface{}{
	// 		"QRiskCessTerm": errorDescription,
	// 	}
	// 	fieldvalerrors = append(fieldvalerrors, fielderror)
	// }
	// if qdetail.QPremCessTerm < q0006data.MinPpt {
	// 	errorDescription := GetError(iCompany, 1, "E0017")
	// 	fmt.Println("Barath", qdetail.QPremCessTerm, q0006data.MinPpt)
	// 	fielderror := map[string]interface{}{
	// 		"QPremCessTerm": errorDescription,
	// 	}
	// 	fieldvalerrors = append(fieldvalerrors, fielderror)
	// }

	if qdetail.QDate == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"QDate": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)
	}

	if qdetail.QCoverage == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"QCoverage": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)
	}
	if qdetail.QAge == 0 {
		errorDescription := GetError(iCompany, 1, "E0019")
		fielderror := map[string]interface{}{
			"QAge": errorDescription,
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

func validateDeathH(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var deathh models.DeathH
	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &deathh); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	if deathh.EffectiveDate == "" {
		errorDescription := GetError(1, 1, "E0018")
		fielderror := map[string]interface{}{
			"EffectiveDate": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)
	}

	if deathh.DeathDate == "" {
		errorDescription := GetError(1, 1, "E0018")
		fielderror := map[string]interface{}{
			"DeathDate": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)
	}
	if deathh.Cause == "" {
		errorDescription := GetError(1, 1, "E0019")
		fielderror := map[string]interface{}{
			"Cause of Death": errorDescription,
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

func validateQHeader(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var qheader quotation.QHeader
	iCompany := qheader.CompanyID

	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &qheader); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	if qheader.QProduct == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"QProduct": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if qheader.AddressID == 0 {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"AddressID": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if qheader.ClientID == 0 {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"ClientID": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if qheader.QAnnualIncome == 0 {
		errorDescription := GetError(iCompany, 1, "E0019")
		fielderror := map[string]interface{}{
			"QAnnualIncome": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if qheader.QOccupation == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"QOccupation": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if qheader.QNri == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"QNri": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if qheader.QuoteDate == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"QuoteDate": errorDescription,
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

func validateAgency(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var agency models.Agency
	iCompany := agency.CompanyID

	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &agency); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	if agency.LicenseNo == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"LicenseNo": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if agency.Startdate == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"Startdate": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	fmt.Println("Comapny Code is ..............", iCompany)
	if agency.EndDate == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"EndDate": errorDescription,
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

func validateProposal(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var policy models.Policy
	iCompany := uint(1)
	fmt.Println("Company Code", iCompany)

	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &policy); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}
	if policy.PRCD > policy.PReceivedDate {
		errorDescription := GetError(iCompany, 1, "E0020")
		fielderror := map[string]interface{}{
			"PDate": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if policy.PRCD == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"PDate": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if policy.PReceivedDate == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"PReceivedDate": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if policy.PProduct == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"PProduct": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if policy.PFreq == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"PFreq": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}
	if policy.PContractCurr == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"PContractCurr": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if policy.PBillCurr == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"PBillCurr": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if policy.POffice == "" {
		errorDescription := GetError(iCompany, 1, "E0018")
		fielderror := map[string]interface{}{
			"POffice": errorDescription,
		}
		fieldvalerrors = append(fieldvalerrors, fielderror)

	}

	if policy.POffice != "" {
		err := CheckParam(iCompany, "P0018", policy.POffice)
		if err != nil {
			fieldvalerrors = append(fieldvalerrors, err)
		}
	}

	if policy.PProduct != "" {
		err := CheckParam(iCompany, "Q0005", policy.PProduct)
		if err != nil {
			fmt.Println("Found Error")
			fieldvalerrors = append(fieldvalerrors, err)
		}
		fmt.Println(err)
	}
	if len(fieldvalerrors) > 0 {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Field Errors"), dbvalerror
	} else {
		return nil, dbvalerror
	}

}

func ValidateBenefit(dbdata map[string]interface{}) (error, DbValError) {

	var dbvalerror DbValError
	fieldvalerrors := make([]interface{}, 0)
	var benefit models.Benefit
	iCompany := uint(1)
	fmt.Println("Company Code", iCompany)

	jsonStr, err := json.Marshal(dbdata)
	if err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &benefit); err != nil {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}
	if len(fieldvalerrors) > 0 {
		dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Field Errors"), dbvalerror
	} else {
		return nil, dbvalerror
	}

}

func ValidatePolicy(dbdata map[string]interface{}) (error, DbValError) {
	var dbvalerror DbValError

	var policy models.Policy
	iCompany := uint(1)
	fmt.Println("Company Code", iCompany)

	jsonStr, err := json.Marshal(dbdata)
	if err != nil {

		// dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &policy); err != nil {
		// dbvalerror.DbErrors = fieldvalerrors
		return errors.New("Marshalign Error"), dbvalerror
	}
	// Actual Validations of Fields Here
	fieldvalerrors := make([]interface{}, 0)
	if policy.PRCD == "" || policy.PProduct == "" || policy.PFreq == "" ||
		policy.AgencyID == 0 || policy.ClientID == 0 || policy.PBillCurr == "" ||
		policy.PContractCurr == "" || policy.POffice == "" || policy.PReceivedDate == "" ||
		policy.PUWDate == "" {
		err := CheckParam(iCompany, "P0018", policy.PRCD)
		if err != nil {
			fieldvalerrors = append(fieldvalerrors, err)
		}
	}

	if policy.PProduct != "" {
		err := CheckParam(iCompany, "E0018", policy.PProduct)
		if err != nil {
			fieldvalerrors = append(fieldvalerrors, err)
		}
	}

	if len(fieldvalerrors) > 0 {
		dbvalerror.DbErrors = fieldvalerrors
		fmt.Println("Shijith Field Val Error", fieldvalerrors)
		dbvalerror.DbErrors = fieldvalerrors
		fmt.Println("Shijith Field Val Error", dbvalerror.DbErrors)
		return errors.New("Field Errors"), dbvalerror

	} else {
		return nil, dbvalerror
	}

}
