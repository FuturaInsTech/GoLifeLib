package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

// *********************************************************************************************
// # 1
// Find out Difference between two Dates  Date should be in 2029-11-05 00:00:00 +0000 UTC
//
// # Outputs are Years, Months, Days, Hours, Minutes, Seconds
//
// ©  FuturaInsTech
// *********************************************************************************************
func DateDiff(a, b time.Time, m string) (year, month, day, hour, min, sec int) {
	// method = N means age nearer birthday
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// *********************************************************************************************
// # 2
// From Date Generally DOB
// To Date Generally RCD
// In case From Date is Lesser than To Date, it will be swapped internally
// Change YYYYMMDD format into Golang Date Format
// imethod is method used .  3 Methods are used .
// Age Last Birth Day - "L"
// Age Next Birth Day  - "X"
// Age Nearest Birth Day - "N"
// It returns values in year, month, day, hour, min and sec.  But we generally use only Year
// This Calcualte Age will call DateDiff
//
// ©  FuturaInsTech
// *********************************************************************************************
func CalculateAge(fromDate, toDate, imethod string) (year, month, day, hour, min, sec int) {
	//
	fromDate1 := fromDate
	toDate1 := toDate
	if fromDate > toDate {
		temp2 := fromDate
		fromDate1 = toDate
		toDate1 = temp2
	}
	dob := String2Date(fromDate1)
	rcd := String2Date(toDate1)

	method := imethod
	year1, month1, day1, hour1, min1, sec1 := DateDiff(dob, rcd, method)
	//Age Nearer Birthday
	if method == "N" {
		if month1 > 5 {
			year1 = year1 + 1
			month1 = 0
		}
	}
	// Age Last Birthday
	if method == "L" {
		month1 = 0
	}
	// Age neXt BirthDay
	if method == "X" {
		year1 = year1 + 1
		month1 = 0
	}
	// fmt.Println(year, month, day, hour, min, sec)
	return year1, month1, day1, hour1, min1, sec1
}

// # 3
// *********************************************************************************************
// This Function Take input as YYYYMMDD and Give Result in
// this format 2023-02-01 00:00:00 +0000 UTC
//
// ©  FuturaInsTech
// *********************************************************************************************
func String2Date(iDate string) (oDate time.Time) {
	//
	format := "20060102"
	//dob, err := time.Parse("20060102", temp)
	oDate, err := time.Parse(format, iDate)
	if err == nil {
		fmt.Println("New Values")
		fmt.Println(oDate)
	}
	return oDate

}

// # 4
// *********************************************************************************************
// Author Ranga
//
//	Convert Date format into String as YYYYMMDD
//
// FuturaInsTech
// *********************************************************************************************
func Date2String(iDate time.Time) (odate string) {

	var temp string
	temp = iDate.String()
	temp1 := temp[0:4] + temp[5:7] + temp[8:10]
	// fmt.Println("Rangarajan Ramaujam ***********")
	// fmt.Println(iDate)
	// fmt.Println(temp1)
	odate = temp1
	return odate

}

// # 5
// *********************************************************************************************
// This GetNextDue Function Give NextDue Date based on the frequency provided
//
// Inputs Date YYYYMMDD
// Frequency (Y/H/Q/M/S)
// Reversal Indicator (R)
// It Convert string into Date Fromat using DateFormat Function
// Then depends upon the frequency next due date is arrived
//
// If Reversal is "R" then previous Due Date will be calcualted.
//
// ©  FuturaInsTech
// *********************************************************************************************
func GetNextDue(iDate string, iFrequency string, iReversal string) (a time.Time) {
	//yyymmdd format to 2023-02-01 00:00:00 +0000 UTC
	x := String2Date(iDate)

	if iReversal == "R" {
		switch {
		case iFrequency == "Y":
			a := AddMonth(x, -12)
			return a
		case iFrequency == "H":
			a := AddMonth(x, -6)
			return a
		case iFrequency == "Q":
			a := AddMonth(x, -3)
			return a
		case iFrequency == "M":
			a := AddMonth(x, -1)
			return a
		case iFrequency == "S":
			a := x.AddDate(0, 0, 0)
			return a

		}
	} else {

		switch {
		case iFrequency == "Y":
			a := AddMonth(x, 12)
			return a
		case iFrequency == "H":
			a := AddMonth(x, 6)
			return a
		case iFrequency == "Q":
			a := AddMonth(x, 3)
			return a
		case iFrequency == "M":
			a := AddMonth(x, 1)
			return a
		case iFrequency == "S":
			a := x.AddDate(0, 0, 0)
			return a
		}
	}
	return
}

// *********************************************************************************************
// # 6
// Add Month should ended with month
// ©  FuturaInsTech
// *********************************************************************************************

func AddMonth(t time.Time, m int) time.Time {
	x := t.AddDate(0, m, 0)
	if d := x.Day(); d != t.Day() {
		return x.AddDate(0, 0, -d)
	}
	return x
}

// # 7
// *********************************************************************************************
// To Calculate No of instalments Paid
//
// Inputs are From Date in YYYYMMDD, To Date in YYYYMMDD, Frequency (Y/H/Q/M/S)
//
// # Output Value is No of instalements
//
// ©  FuturaInsTech
// *********************************************************************************************
func GetNoIstalments(iFromDate, iToDate, iFrequency string) (oInstalments int) {

	fromDate := String2Date(iFromDate)
	toDate := String2Date(iToDate)
	method := "M"
	var noinstalments float64

	year1, month1, _, _, _, _ := DateDiff(fromDate, toDate, method)
	fmt.Println("Shubham", iFromDate, iToDate, iFrequency, year1, month1)
	switch {

	case iFrequency == "Y":
		// 10 and 0  10
		noinstalments = (float64(year1) / 1)
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments
	case iFrequency == "H":
		// 10 and 6  10*2 + 6/6 =  21
		noinstalments := float32((year1 * 2) + (month1 / 6))
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments

	case iFrequency == "Q":
		// 5 9   = 5 * 4  + 9/3    20 + 3
		noinstalments := float32((year1 * 4) + (month1 / 3))
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments
	case iFrequency == "M":
		noinstalments := float32((year1 * 12) + (month1))
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments
	case iFrequency == "S":
		oInstalments := 1
		return oInstalments

	}
	return
}

// #8
// *********************************************************************************************
// Function : GetPremium Paid
//
// Purpose : To Get Sum of Preium Paid between Two Duration for the given Premium based on Frequency
//
// Inputs FromDate, Todate in YYYYMMDD Format
// Frequency (Y/H/Q/M/S)
// Model Premium (eg., 1234.22)
//
// Output Premium is Mo of instalments Paid * Premium
//
// ©  FuturaInsTech
// *********************************************************************************************
func GetPremiumPaid(iFromDate, iToDate, iFrequency string, iModelPrem float64) (oPremium float64) {

	var oNoOfDues int
	oNoOfDues = GetNoIstalments(iFromDate, iToDate, iFrequency)
	oPrem := float64(float64(oNoOfDues) * iModelPrem)
	oPrem = RoundFloat(oPrem, 2)
	return float64(oPrem)

}

// # 9
// *********************************************************************************************
// Function Name : GetPaidUp
//
// # Calculate Paidup Up Value
//
// Input Values RCD Date, Paid To Date, Premium Cessation Date in YYYYMMDD Format
// Frequency (Y/H/Q/M/S)
//
// Output is (Paid/Payable) * SA as float
//
// ©  FuturaInsTech
// *********************************************************************************************
func GetPaidUp(iFromDate, iToDate, iPremCessDate, iFrequency string, iSumAssured float32) (oPaidUpValue float32) {

	a := GetNoIstalments(iFromDate, iToDate, iFrequency)
	b := GetNoIstalments(iFromDate, iPremCessDate, iFrequency)

	oPaidUpValue = float32((a / b) * int(iSumAssured))
	// fmt.Println(" a and b values are", a, b)
	// fmt.Println(int(oPaidUpValue))
	return oPaidUpValue

}

// # 10
// func RoundAmt(iAmount float64, iMethod string) (oAmount float64) {
// Function Name : Round Float To Round to 2 or 3 or 4 Decimals
//
// # Inputs Amount and Precision which is 2 or 3 or 4
//
// # Output Return Value after Round
//
// ©  FuturaInsTech
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// # 11
// Function Name : AddLeadDays
//
// Inputs Date String YYYYMMDD and Days integer. (days can be -negative as well)
//
// # Output Date String YYYYMMDD
//
// ©  FuturaInsTech
func AddLeadDays(iDate string, iDays int) (oDate string) {
	a := String2Date(iDate)
	b := a.AddDate(0, 0, iDays)
	c := Date2String(b)
	// fmt.Println("Lead Days ", a, b, c)
	return c

}

// # 12
// Function Name : Simple Interest Calculation
//
// # Inputs Principal, Interest and No of Days
//
// # Output Interest only
//
// ©  FuturaInsTech
func SimpleInterest(iPrincipal, iInterest, iDays float64) (oInterest float64) {
	oInterest = iPrincipal * (iInterest / 100) * (iDays / 365)
	return oInterest
}

// # 13
// Function : CompoundInterest
//
// # Purpose to Calculate Compounding Interest for a Given Days
//
// # Inputs Amount, Interest as float , Days
//
// # Outputs Actual Interest
//
// ©  FuturaInsTech
func CompoundInterest(iPrincipal, iInterest, iDays float64) (oInterest float64) {
	oDays := iDays / 365
	oInt := 1 + (iInterest / 100)
	oInterest1 := iPrincipal * (math.Pow(oInt, oDays))
	oInterest = oInterest1 - iPrincipal
	//	fmt.Println("Compounding Interest", iPrincipal, iInterest, iDays, oInterest1, oInterest, oInt)
	return oInterest
}

// # 14
// Function : ValidateFields
//
// # Function, FieldName, FieldValue, User ID and User Type
//
// # Outputs Error
//
// ©  FuturaInsTech
func ValidateFields(iFunction string, iFieldName string, iFieldVal string, iUserId uint64, iFieldType string) error {
	var fieldvalidators models.FieldValidator
	var getUser models.User
	results := initializers.DB.First(&getUser, "id = ?", iUserId)

	if results.Error != nil {
		return errors.New(results.Error.Error())
	} else {
		oLanguageId := getUser.LanguageID
		oCompanyId := getUser.CompanyID
		results := initializers.DB.First(&fieldvalidators, "function_name = ? and company_id = ? and language_id =? and field_name = ?", iFunction, oCompanyId, oLanguageId, iFieldName)

		if results.Error != nil {
			return nil
		} else {
			if fieldvalidators.ParamName != "" {
				fmt.Println("I am here Ranga")
				//	fmt.Println(iFieldName, iFieldVal, fieldvalidators.ParamName, fieldvalidators.ErrorDescription)
				// fmt.Println("I am inside ", fieldvalidators.ParamName, iFieldVal, iFieldName)
				err := ValidateItem(iUserId, fieldvalidators.ParamName, iFieldVal, iFieldName, fieldvalidators.ErrorDescription)
				if err != nil {
					fmt.Println("Error inside Validator ", err)
					return err
				}
			} else if iFieldType == "string" {
				if fieldvalidators.BlankAllowed == "N" && iFieldVal == "" {
					return errors.New(fieldvalidators.ErrorDescription + "-" + iFieldName)
				}

			} else if iFieldType != "string" {
				fmt.Println("Field Type is ", iFieldType, iFieldVal)

				if fieldvalidators.ZeroAllowed == "N" && iFieldVal == "0" {
					return errors.New(fieldvalidators.ErrorDescription + "-" + iFieldName)
				}
			}

		}
		return nil
	}

}

// # 15
// Function : ValidateItem
//
// # Inputs UserId, Name, Item, FieldName, Errors
//
// # Outputs Error
//
// ©  FuturaInsTech
func ValidateItem(iUserId uint64, iName string, iItem string, iFieldName string, iErros string) error {
	var getUser models.User
	results := initializers.DB.First(&getUser, "id = ?", iUserId)
	if results.Error != nil {
		fmt.Println(results.Error)
		return errors.New(results.Error.Error())
	}
	var valdiateparam models.ParamDesc
	oLanguageId := getUser.LanguageID
	oCompanyId := getUser.CompanyID
	results = initializers.DB.Where("company_id = ? AND name = ? and item = ? and language_id = ?", oCompanyId, iName, iItem, oLanguageId).Find(&valdiateparam)
	if results.Error != nil || results.RowsAffected == 0 {

		return errors.New(" -" + strconv.FormatUint(uint64(oCompanyId), 10) + "-" + iName + "-" + strconv.FormatUint(uint64(oLanguageId), 10) + "-" + "-" + iFieldName + iErros + " is missing")
		//return errors.New(results.Error.Error())
	}
	return nil
}

// # 16
// Function : ValidateItem
//
// # Inputs Company, Business Rule, Business Rule ID, From Date, Data
//
// # Outputs Error
//
// ©  FuturaInsTech
func GetItemD(iCompany int, iTable string, iItem string, iFrom string, data *paramTypes.Extradata) error {

	//var sourceMap map[string]interface{}
	var itemparam models.Param
	//	fmt.Println(iCompany, iItem, iFrom)
	results := initializers.DB.Find(&itemparam, "company_id =? and name= ? and item = ? and rec_type = ? and ? between start_date  and  end_date", iCompany, iTable, iItem, "IT", iFrom)

	if results.Error == nil && results.RowsAffected != 0 {
		(*data).ParseData(itemparam.Data)
		return nil
	} else {
		if results.Error != nil {
			return errors.New(results.Error.Error())
		} else {
			return errors.New("No Item Found " + iTable + iItem)
		}

	}
}

// #17
// GetAnnualPrem
//
// Inputs Premium float64, Frequency (Y/H/Q/M/S)
//
// # Outputs Annualized Premium float64
//
// ©  FuturaInsTech
func GetAnnualPrem(iPrem float64, iFrequency string) (oPremium float64) {

	switch {
	case iFrequency == "Y":
		a := iPrem
		return a
	case iFrequency == "H":
		a := iPrem * 2
		return a
	case iFrequency == "Q":
		a := iPrem * 4
		return a
	case iFrequency == "M":
		a := iPrem * 12
		return a
	case iFrequency == "S":
		a := 0 * iPrem
		return a
	}
	return
}

// #18
// GetNextYr
//
// Inputs FromDate String (YYYYMMDD)
//
// Outputs ToDate String (YYYYMMDD)
//
// ©  FuturaInsTech
func GetNextYr(iFrom string) (iTo string) {
	a := String2Date(iFrom)
	b := a.AddDate(1, 0, 0)
	c := Date2String(b)
	return c

}

// #19
// ModeChange
//
// Inputs Old Frequency and New Frequency (Y/H/Q/M).  Old Premium eg., 1234.22
//
// # Outputs New Premium as 14808.02
//
// ©  FuturaInsTech
func ModeChage(iOldFreq, iNewFreq string, iOldPrem float64) (oNewPrem float64) {
	var ofreq int
	var nfreq int

	switch {
	case iOldFreq == "Y":
		ofreq = 01
	case iOldFreq == "H":
		ofreq = 02
	case iOldFreq == "Q":
		ofreq = 04
	case iOldFreq == "M":
		ofreq = 12
	}
	switch {
	case iNewFreq == "Y":
		nfreq = 01
	case iNewFreq == "H":
		nfreq = 02
	case iNewFreq == "Q":
		nfreq = 04
	case iNewFreq == "M":
		nfreq = 12
	}
	oNewPrem = float64(ofreq/nfreq) * iOldPrem
	return oNewPrem

}

// # 20
// AddYears2Date
//
// Inputs: Date String in YYYYMMDD and Years months and days to be added (eg., 22 Yrs  11 Months 2 Days)
//
// # Outputs Date in String YYYYMMDD   20220220
//
// ©  FuturaInsTech
func AddYears2Date(iDate string, iYrs int, iMonth int, iDays int) (oDate string) {
	a := String2Date(iDate)
	b := a.AddDate(iYrs, iMonth, iDays)
	c := Date2String(b)
	return c
}

// # 21
// Convert Map to Strucutre
//
// Inputs: Map and Interface
//
// # Outputs Error
//
// ©  FuturaInsTech
func ConvertMapToStruct(m map[string]interface{}, s interface{}) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		if value, ok := m[field.Name]; ok {
			stValue.Field(i).Set(reflect.ValueOf(value))
		}
	}
	return nil
}

// # 22
// GetAnnualRate - Get Annual Rate of the Coverage - No Model Discount/Staff Discount/SA/Prem Discount
//
// Inputs: Company, Coverage, Age (Attained Age), Gender(F/N/U), Term (2 Characters), Premium Method as
// PM001 - Term Based , PM002 Age Based, Mortality Clause "S" Smoker, "N" Non Smoker
//
// Outputs Annualized Premium as float (124.22)
//
// ©  FuturaInsTech

func GetAnnualRate(iCompany uint, iCoverage string, iAge uint, iGender string, iTerm uint, iPremTerm uint, iPremMethod string, iDate string, iMortality string) (float64, error) {

	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)

	var q0010data paramTypes.Q0010Data
	var extradataq0010 paramTypes.Extradata = &q0010data
	var q0010key string
	var prem float64
	//term := strconv.FormatUint(uint64(iTerm), 10)
	//premTerm := strconv.FormatUint(uint64(iPremTerm), 10)

	term := fmt.Sprintf("%02d", iTerm)
	premTerm := fmt.Sprintf("%02d", iPremTerm)

	//fmt.Println("****************", iCompany, iCoverage, iAge, iGender, iTerm, iPremMethod, iDate, iMortality)
	if q0006data.PremCalcType == "A" || q0006data.PremCalcType == "U" {
		if q0006data.PremiumMethod == "PM002" {
			// END1 + Male
			q0010key = iCoverage + iGender
		}
	} else if q0006data.PremCalcType == "P" {
		// END1 + Male + Term + Premium Term
		if q0006data.PremiumMethod == "PM001" || q0006data.PremiumMethod == "PM003" {
			q0010key = iCoverage + iGender + term + premTerm

		}

	} else if q0006data.PremCalcType == "H" {
		// HIP1 + Male
		if q0006data.PremiumMethod == "PM005" {
			q0010key = iCoverage + iGender
		}
	}
	fmt.Println("Premium Key ******", iCoverage, iGender, term, premTerm, q0006data.PremCalcType, q0010key)
	err := GetItemD(int(iCompany), "Q0010", q0010key, iDate, &extradataq0010)
	if err != nil {
		return 0, err

	}
	fmt.Println("************", iCompany, iAge, q0010key, iDate)

	for i := 0; i < len(q0010data.Rates); i++ {
		if q0010data.Rates[i].Age == uint(iAge) {
			prem = q0010data.Rates[i].Rate
			break
		}
	}
	fmt.Println("************", iCompany, iAge, q0010key, iDate, prem)
	return prem, nil
}

// # 23
// ValidateCoverageQ0011 - Rider is Allowed for Product or Not Validation
//
// Inputs: Company, Product, Coverage and Date String in YYYYMMDD
//
// Outputs Product Found or Not  "Y" Means Found "N" Means Not Found
//
// ©  FuturaInsTech
func ValidateCoverageQ0011(iCompany uint, iProduct, iCoverage, iDate string) string {

	fmt.Println("Coverages Q0011", iCompany, iProduct, iCoverage, iDate)
	var q0011data paramTypes.Q0011Data
	var extradataq0011 paramTypes.Extradata = &q0011data
	productFound := "N"
	err := GetItemD(int(iCompany), "Q0011", iProduct, iDate, &extradataq0011)
	if err != nil {
		return productFound
	}

	for i := 0; i < len(q0011data.Coverages); i++ {
		if q0011data.Coverages[i].CoverageName == iCoverage {
			productFound = "Y"
			return productFound
		}
	}
	return productFound
}

// # 24
// Levels
// AddYears2Date
//
// Inputs:
//
// # Outputs
//
// ©  FuturaInsTech
func CustomizedPreload(d *gorm.DB) *gorm.DB {
	return d.Preload("Levels", CustomizedPreload)
}

// # 25
// ValidateQ0012 - Survival Benefit (Term Based)
//
// Inputs: Company, Coverage and Date String in YYYYMMDD
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func ValidateQ0012(iCompany uint, iCoverage string, iDate string) error {
	var q0012data paramTypes.Q0012Data
	var extradataq0012 paramTypes.Extradata = &q0012data

	err := GetItemD(int(iCompany), "Q0012", "AED1", iDate, &extradataq0012)

	if err != nil {
		return err

	}

	for i := 0; i < len(q0012data.SbRates); i++ {
		fmt.Println("Survival Benefits .......")
		fmt.Println(q0012data.SbRates[i].Term)
		fmt.Println(q0012data.SbRates[i].Percentage)
	}
	return nil
}

// # 26
// ValidateQ0013 - Survival Benefit (Age Based)
//
// Inputs: Company, Coverage and Date String in YYYYMMDD
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func ValidateQ0013(iCompany uint, iCoverage string, iDate string) error {
	var q0013data paramTypes.Q0013Data
	var extradataq0013 paramTypes.Extradata = &q0013data

	err := GetItemD(int(iCompany), "Q0013", "AEDR", iDate, &extradataq0013)

	if err != nil {
		return err

	}
	fmt.Println(q0013data.SbRates[0].Percentage)
	for i := 0; i < len(q0013data.SbRates); i++ {
		fmt.Println("Survival Benefits .......")
		fmt.Println(q0013data.SbRates[i].Age)
		fmt.Println(q0013data.SbRates[i].Percentage)
	}
	return nil
}

// # 27
// ValidateQ0013 - Survival Benefit (Age Based)
//
// Inputs: Company, Coverage and Date String in YYYYMMDD
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func GetSBByYear(iCompany uint, iCoverage string, iDate string, iSA float64, iType string, iMethod string, iYear int, iAge int) float64 {

	if iType == "T" {
		var q0012data paramTypes.Q0012Data
		var extradataq0012 paramTypes.Extradata = &q0012data
		// fmt.Println("SB Parameters", iCompany, iType, iMethod, iYear, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0012", iMethod, iDate, &extradataq0012)

		if err != nil {
			return 0

		}
		// fmt.Println(q0012data.SBRates[0].Percentage)
		for i := 0; i < len(q0012data.SbRates); i++ {
			if iYear == int(q0012data.SbRates[i].Term) {
				oSB := q0012data.SbRates[i].Percentage * iSA
				return oSB
			}
		}
	}
	if iType == "A" {
		var q0013data paramTypes.Q0013Data
		var extradataq0013 paramTypes.Extradata = &q0013data
		fmt.Println("SB Parameters", iCompany, iType, iMethod, iAge, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0013", iMethod, iDate, &extradataq0013)
		fmt.Println("SB Parameters", iCompany, iCoverage, iDate)

		if err != nil {
			return 0

		}
		fmt.Println(q0013data.SbRates[0].Percentage)
		for i := 0; i < len(q0013data.SbRates); i++ {
			if iAge == int(q0013data.SbRates[i].Age) {
				oSB := q0013data.SbRates[i].Percentage * iSA
				return oSB
			}
		}
	}
	return 0
}

// # 28
// GetBonus - Get Bonus for a Given Duration
//
// Inputs: Company, Bonus Method, Status, Coverage Start Date, Year of Policy, Policy Status, SA
//
// # Date in YYYYMMDD as a string
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func GetBonus(iCompany uint, iCoverage string, iStartDate string, iEndDate string, iStatus string, iTerm uint, iSA uint) uint {

	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iStartDate, &extradata)

	iRBMethod := q0006data.RevBonus
	// iIBMethod := q0006data.IBonus
	// iTBMethod := q0006data.TBonus
	// iLBMethod := q0006data.LoyaltyBonus
	// iSSVMethod := q0006data.SSVMethod
	// iGSVMethod := q0006data.GSVMethod
	// iBSVMethod := q0006data.BSVMethod
	var q0014data paramTypes.Q0014Data
	var extradata1 paramTypes.Extradata = &q0014data

	iKey := iRBMethod + iStatus

	fmt.Println("Sreeeram XXXXXXXXXXXXXX", iCompany, iRBMethod, iStatus, iStartDate, iEndDate, iSA, iKey)

	var term uint

	term = 1

	for term <= iTerm {

		GetItemD(int(iCompany), "Q0014", iKey, iStartDate, &extradata1)

		for i := 0; i < len(q0014data.BRates); i++ {
			if term <= uint(q0014data.BRates[i].Term) {
				fmt.Println("one by one ", iSA, q0014data.BRates[i].Percentage, q0014data.BRates[i].Term, term)
				oBonus := float64(iSA) * ((q0014data.BRates[i].Percentage) / 100)
				fmt.Println(q0014data.BRates[i].Term, q0014data.BRates[i].Percentage, i)
				fmt.Println("bonus Values ********", int(oBonus))
				break
			}
		}
		term = term + 1
		temp := AddYears2Date(iStartDate, 1, 0, 0)
		iStartDate = temp

	}
	return 0

}

// # 29
// GetBonusByYear - Get Bonus for a Given Year Array of 10 Allowed
//
// Inputs: Company, Coverage , Bonus Method, Status, Coverage Start Date, Year of Policy, Policy Status, SA
//
// Date in YYYYMMDD as a string. Bonus Method such as RB/IB/TB/LB/SSV/GSV/BSV
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func GetBonusByYear(iCompany uint, iCoverage string, iBonusMethod string, iDate string, iYear uint, iStatus string, iSA uint) uint64 {
	//	fmt.Println("inside Bonus ", iCoverage, iCompany, iBonusMethod, iDate, iYear, iStatus, iSA)
	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)

	var key1 string
	if iBonusMethod == "RB" {
		key1 = q0006data.RevBonus
	} else if iBonusMethod == "IB" {
		key1 = q0006data.IBonus
	} else if iBonusMethod == "TB" {
		key1 = q0006data.TBonus
	} else if iBonusMethod == "LB" {
		key1 = q0006data.LoyaltyBonus
	} else if iBonusMethod == "SSV" {
		key1 = q0006data.SsvMethod
	} else if iBonusMethod == "GSV" {
		key1 = q0006data.GsvMethod
	} else if iBonusMethod == "BSV" {
		key1 = q0006data.BsvMethod
	} else if iBonusMethod == "GA" {
		key1 = q0006data.GBonus
	} else if iBonusMethod == "DV" {
		key1 = q0006data.DivMethod
	} else if iBonusMethod == "DI" {
		key1 = q0006data.DivIMethod
	}

	// Method is Not Found Exit
	if key1 == "" {
		oBonus := 0
		return uint64(oBonus)
	}
	var q0014data paramTypes.Q0014Data
	var extradata1 paramTypes.Extradata = &q0014data
	iKey := key1 + iStatus

	GetItemD(int(iCompany), "Q0014", iKey, iDate, &extradata1)
	for i := 0; i < len(q0014data.BRates); i++ {
		if iYear <= q0014data.BRates[i].Term {
			oBonus := iSA * uint(q0014data.BRates[i].Percentage) / 100
			return uint64(oBonus)
		}
	}
	return 0
}

// # 30 (Redundant)
// GetTerm - Term
//
// Inputs: Company, Coverage , Bonus Method, Status, Coverage Start Date, Year of Policy, Policy Status, SA
//
// Date in YYYYMMDD as a string. Bonus Method such as RB/IB/TB/LB/SSV/GSV/BSV
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech

func GetTerm(iCompany uint, iCoverage string, iDate string) {
	var q0015data paramTypes.Q0015Data
	var extradata paramTypes.Extradata = &q0015data
	iKey := iCoverage
	fmt.Println(iKey)

	GetItemD(int(iCompany), "Q0015", iKey, iDate, &extradata)
	for i := 0; i < len(q0015data.Terms); i++ {
		term := q0015data.Terms[i].Term
		fmt.Println(term)

	}
	return
}

// # 31 Redundant
func GetPTerm(iCompany uint, iCoverage string, iDate string) {
	var q0016data paramTypes.Q0016Data
	var extradata paramTypes.Extradata = &q0016data
	iKey := iCoverage
	fmt.Println(iKey)

	GetItemD(int(iCompany), "Q0016", iKey, iDate, &extradata)
	// for i := 0; i <= 98; i++ {
	// 	term := q0016data.PTerms[i].PTerm

	// 	if term == 0 {
	// 		break
	// 	}
	// 	fmt.Println("PPT ", term)

	// }
	// return

	for i := 0; i < len(q0016data.PTerms); i++ {
		term := q0016data.PTerms[i].PTerm
		if term == 0 {
			break
		}
		fmt.Println("Shijith PPT ", term)

	}

}

// # 32
// CalcSaPremDiscount - Calculate Discounted Amount based on SA or Annualised Prem
//
// Inputs: Company, Discount Type (S/P) , Discount Method (As per Product), Annualised Prem
// SA Amount
//
// # Outputs Discounted Amount as float
//
// ©  FuturaInsTech
func CalcSaPremDiscount(iCompany uint, iDiscType string, iDiscMethod string, iAnnPrem float64, iSA uint, iDate string) float64 {
	// SA Discount

	if iDiscType == "S" {
		var q0017data paramTypes.Q0017Data
		var extradataq0017 paramTypes.Extradata = &q0017data
		err := GetItemD(int(iCompany), "Q0017", iDiscMethod, iDate, &extradataq0017)

		if err != nil {
			return 0

		}

		for i := 0; i < len(q0017data.SaBand); i++ {
			if int(iSA) <= int(q0017data.SaBand[i].Sa) {
				oDiscount := uint(q0017data.SaBand[i].Discount) * uint(iAnnPrem) / 100
				return float64(oDiscount)
			}
		}
	}
	// Premium Discount
	if iDiscType == "P" {
		var q0018data paramTypes.Q0018Data
		var extradataq0018 paramTypes.Extradata = &q0018data

		err := GetItemD(int(iCompany), "Q0018", iDiscMethod, iDate, &extradataq0018)

		if err != nil {
			return 0

		}

		for i := 0; i < len(q0018data.PremBand); i++ {
			if int(iAnnPrem) <= int(q0018data.PremBand[i].AnnPrem) {
				oDiscount := uint(q0018data.PremBand[i].Discount) * uint(iAnnPrem) / 100
				return float64(oDiscount)
			}
		}
	}
	return 0
}

// # 33
// CalcFrequencyPrem - Calculate Frequency Premium as per Model Factor Provided
//
// Inputs: Company, Frequency Factor Method as mentioned in Q0006, Current Frequency, Annualized Premium of the Coverage
//
// Output Model Premium =  Model Factor * Annualized Premium.
//
// ©  FuturaInsTech
func CalcFrequencyPrem(iCompany uint, iDate, iFreqMethod string, iFreq string, iAnnPrem float64) float64 {
	var q0019data paramTypes.Q0019Data
	var extradataq0019 paramTypes.Extradata = &q0019data
	err := GetItemD(int(iCompany), "Q0019", iFreqMethod, iDate, &extradataq0019)

	if err != nil {
		return iAnnPrem

	}

	for i := 0; i < len(q0019data.FreqFactor); i++ {
		if iFreq == q0019data.FreqFactor[i].Frequency {
			var a, oPrem float64
			a = float64(q0019data.FreqFactor[i].Factor)

			oPrem = float64(iAnnPrem) * a
			oPrem = RoundFloat(oPrem, 2)
			return oPrem
		}
	}
	return 0

}

// # 34
// GetWaiverSA - Calculate Waiver SA of a Policy
//
// Inputs: Company, All Coverages under the policy , Waiver Method as per Q0006, Waiver Coverage Start Date, Premium of the Current Coverage
//
// First Check whether given coverage is available in Q0020 for the Waiver Method.
// Foud, Add it in output and return SA
//
// ©  FuturaInsTech
func GetWaiverSA(iCompany uint, iCoverage string, iMethod string, iDate string, iPrem float64) float64 {

	var q0020data paramTypes.Q0020Data
	var extradataq0020 paramTypes.Extradata = &q0020data
	err := GetItemD(int(iCompany), "Q0020", iMethod, iDate, &extradataq0020)

	if err != nil {
		return 0

	}
	for i := 0; i < len(q0020data.WaiverCoverages); i++ {
		if iCoverage == q0020data.WaiverCoverages[i].Coverage {
			oSA := iPrem
			return oSA
		}
	}
	return 0
}

// # 35
// GetULAllocRates - Get Unit Linked Allocation Percentage
//
// Inputs: Company,  Date String in YYYYMMDD, Allocation Method as defined in Q0006, Frequency of the Policy
// RCD Date and Current Paid Todate in YYYYMMDD Format.
//
// First it converted no of dues in months and then look at table Q0021 and get applicable rates
// Return Percentage.  It has to be posted to different GL Accounts
//
// ©  FuturaInsTech
func GetULAllocRates(iCompany uint, iDate string, iAllMethod string, iFrequency string, iFromDate string, iToDate string) float64 {

	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data

	err := GetItemD(int(iCompany), "P0060", iAllMethod, iDate, &extradatap0060)
	if err != nil {
		return 0

	}

	noofdues := GetNoIstalments(iFromDate, iToDate, "M")
	fmt.Println("Inside Allocation", iCompany, iDate, iAllMethod, iFrequency, iFromDate, iToDate)

	for i := 0; i < len(p0060data.AlBand); i++ {
		if uint(noofdues) <= uint(p0060data.AlBand[i].Months) {
			iRate := p0060data.AlBand[i].Percentage
			return iRate
		}
	}
	return 0
}

// # 36
// GetULMortPrem - Get Unit Linked Mortality Prem for a given duration
//
// Inputs: Company,  Coverage and Date String in YYYYMMDD, SA, Fund Value, Attained Age, Gender
//
// # Outputs Premium
//
// ©  FuturaInsTech
func GetULMortPrem(iCompany uint, iCoverage string, iDate string, iSA uint64, iFund uint64, iAge uint, iGender string) float64 {

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	// Get Coverage Rules
	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
	if err != nil {
		return 0

	}
	// Check Basis  1 = SAR  2 = SA  3 = SA + Fund
	var oSA uint64
	if q0006data.UlMorttMethod == "UM001" {
		oSA = iSA - iFund
	} else if q0006data.UlMorttMethod == "UM002" {
		oSA = iSA
	} else if q0006data.UlMorttMethod == "UM003" {
		oSA = iSA + iFund
	}

	var q0022data paramTypes.Q0022Data
	var extradataq0022 paramTypes.Extradata = &q0022data
	key := q0006data.UlMorttMethod + iGender
	// Get Premium Rate
	err = GetItemD(int(iCompany), "Q0022", key, iDate, &extradataq0022)
	if err != nil {
		return 0

	}
	var aPrem float64
	for i := 0; i < len(q0022data.Rates); i++ {
		if q0022data.Rates[i].Age == uint(iAge) {
			aPrem = q0022data.Rates[i].Rate * float64(oSA)
		}
	}
	// Apply Model Factor
	if q0006data.UlMortFreq == "M" {
		aPrem = aPrem * 0.0833
	} else if q0006data.UlMortFreq == "Q" {
		aPrem = aPrem * 0.25
	} else if q0006data.UlMortFreq == "H" {
		aPrem = aPrem * 0.5
	}

	return aPrem

}

// # 37
// GetGSTPercentage - Get GST Percemtage for a given months
//
// Inputs: Company,  Coverage and Date String in YYYYMMDD (Current Date), Key is Coverage Code, No of Months, Amount to be charged
//
// # Outputs GST Amount
//
// ©  FuturaInsTech
func GetGSTAmount(iCompany uint, iDate string, iKey string, iMonths uint64, iAmount float64) float64 {

	var q0023data paramTypes.Q0023Data
	var extradataq0023 paramTypes.Extradata = &q0023data

	// Get Premium Rate
	err := GetItemD(int(iCompany), "Q0023", iKey, iDate, &extradataq0023)
	if err != nil {
		return 0
	}

	for i := 0; i < len(q0023data.Gst); i++ {
		if uint(iMonths) <= q0023data.Gst[i].Month {
			oAmount := iAmount * q0023data.Gst[i].Rate
			oAmount = RoundFloat(oAmount, 2)
			return oAmount
		}
	}
	return 0
}

// # 38
// GetMaxTranno - Get Transaction No and History Code
//
// Inputs: Company,  Policy No, Method, Effective Date, User
//
// # Outputs History Code and New Tranno
//
// # It update PHISTORY Table
//
// ©  FuturaInsTech
func GetMaxTranno(iCompany uint, iPolicy uint, iMethod string, iEffDate string, iuser uint64, historyMap map[string]interface{}) (string, uint) {
	var permission models.Permission
	var result *gorm.DB

	result = initializers.DB.First(&permission, "company_id = ? and method = ?", iCompany, iMethod)
	if result.Error != nil {
		return iMethod, 0
	}
	iHistoryCode := permission.TransactionID
	var transaction models.Transaction
	result = initializers.DB.Find(&transaction, "ID = ?", iHistoryCode)
	if result.Error != nil {
		return iMethod, 0
	}
	iHistoryCD := transaction.TranCode
	var phistory models.PHistory
	var maxtranno float64 = 0

	fmt.Println(iCompany, iPolicy, iHistoryCD, iEffDate)

	result1 := initializers.DB.Table("p_histories").Where("company_id = ? and policy_id= ?", iCompany, iPolicy).Select("max(tranno)")

	if result1.Error != nil {
		fmt.Println(result1.Error)

	}
	err := result1.Row().Scan(&maxtranno)
	fmt.Println("Error ", err)
	phistory.CompanyID = iCompany
	phistory.Tranno = uint(maxtranno) + 1
	phistory.PolicyID = iPolicy
	phistory.HistoryCode = iHistoryCD
	phistory.EffectiveDate = iEffDate
	phistory.Is_reversed = false
	phistory.IsValid = "1"
	if historyMap != nil {
		phistory.PrevData = historyMap
	}
	a := time.Now()
	b := Date2String(a)
	phistory.CurrentDate = b
	phistory.UpdatedID = iuser
	result1 = initializers.DB.Create(&phistory)
	if result1.Error != nil {
		fmt.Println(result1.Error)

	}

	return phistory.HistoryCode, phistory.Tranno

}

// # 39
// Post GL - Get Transaction No and History Code
//
// ©  FuturaInsTech
func PostGlMove(iCompany uint, iContractCurry string, iEffectiveDate string,
	iTranno int, iGlAmount float64, iAccAmount float64, iAccountCodeID uint, iGlRdocno uint,
	iGlRldgAcct string, iSeqnno uint64, iGlSign string, iAccountCode string, iHistoryCode string, iRevInd string, iCoverage string) error {

	iAccAmount = RoundFloat(iAccAmount, 2)

	var glmove models.GlMove
	var company models.Company
	glmove.ContractCurry = iContractCurry
	glmove.ContractAmount = iAccAmount
	initializers.DB.Find(&company, "id = ?", iCompany)
	var currency models.Currency
	// fmt.Println("Currency Code is .... ", company.CurrencyID)
	initializers.DB.Find(&currency, "id = ?", company.CurrencyID)
	iGlCurry := currency.CurrencyShortName
	glmove.CurrencyRate = 1
	if glmove.GlCurry != glmove.ContractCurry {
		var p0031data paramTypes.P0031Data
		var extradata paramTypes.Extradata = &p0031data
		iKey := iContractCurry + "2" + iGlCurry
		// fmt.Println("i key ", iKey)
		err := GetItemD(int(iCompany), "P0031", iKey, iEffectiveDate, &extradata)
		if err != nil {
			fmt.Println("I am inside Error in Exchange ")
			glmove.CurrencyRate = 1
		} else {
			for i := 0; i < len(p0031data.CurrencyRates); i++ {
				// fmt.Println("Exchange Rates", p0031data.CurrencyRates[i].Action, p0031data.CurrencyRates[i].Action)
			}
			glmove.CurrencyRate = p0031data.CurrencyRates[0].Rate
			// fmt.Println("I am outside Error in Exchange ")

		}
	}

	// fmt.Println("Exchange Rate &&&&&&&&&&&&&&&&&&&&", glmove.CurrencyRate)

	glmove.AccountCode = iAccountCode
	glmove.AccountCodeID = iAccountCodeID
	glmove.CompanyID = iCompany
	glmove.EffectiveDate = iEffectiveDate
	glmove.GlAmount = iAccAmount * glmove.CurrencyRate
	glmove.GlCurry = iGlCurry
	glmove.GlRdocno = strconv.Itoa(int(iGlRdocno))
	glmove.GlRldgAcct = iGlRldgAcct
	glmove.GlSign = iGlSign
	glmove.HistoryCode = iHistoryCode
	glmove.Tranno = uint(iTranno)
	glmove.SequenceNo = uint64(iSeqnno)
	curr_date := time.Now()
	glmove.CurrentDate = Date2String(curr_date)

	GlRdocno := glmove.GlRdocno
	glmove.ID = 0
	glmove.ReversalIndicator = iRevInd
	glmove.BCoverage = iCoverage
	tx := initializers.DB.Save(&glmove)
	tx.Commit()

	UpdateGlBal(iCompany, iGlRldgAcct, iAccountCode, iContractCurry, iAccAmount, iGlSign, GlRdocno)
	return nil
}

// # 40
// GetCommissionRates - Get Commission Rates
//
// Inputs: Company,  Coverage, Nof Instalments Collected (so far) and Date String in YYYYMMDD
//
// # Outputs Commission Rate
//
// ©  FuturaInsTech
func GetCommissionRates(iCompany uint, iCoverage string, iNofInstalemnts uint, iDate string) float64 {

	var p0028data paramTypes.P0028Data
	var extradatap0028 paramTypes.Extradata = &p0028data
	iKey := iCoverage
	fmt.Println("commission Rates **********", iCompany, iCoverage, iDate, iNofInstalemnts, iKey)
	// Get Premium Rate
	err := GetItemD(int(iCompany), "P0028", iKey, iDate, &extradatap0028)
	if err != nil {
		return 0
	}

	for i := 0; i < len(p0028data.Commissions); i++ {
		if uint(iNofInstalemnts) <= p0028data.Commissions[i].Ppt {
			fmt.Println("Iam inside the array", p0028data.Commissions[i].Ppt)
			oRate := p0028data.Commissions[i].Rate
			fmt.Println("i am getting in ", p0028data.Commissions[i].Rate)
			return oRate
		}
	}
	return 0
}

// # 41
// UpdateGlBal
//
// ©  FuturaInsTech
func UpdateGlBal(iCompany uint, iGlRldgAcct string, iGlAccountCode string, iContCurry string, iAmount float64, iGLSign string, iGlRdocno string) float64 {
	var glbal models.GlBal
	var temp float64
	if iGLSign == "-" {
		temp = iAmount * -1

	} else {
		temp = iAmount * 1
	}
	var company []models.Company
	initializers.DB.First(&company, "id = ?", iCompany)

	results := initializers.DB.First(&glbal, "company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno)
	if results.RowsAffected == 0 {
		glbal.ContractAmount = temp
		glbal.CompanyID = iCompany
		glbal.GlAccountno = iGlAccountCode
		glbal.GlRldgAcct = iGlRldgAcct
		glbal.ContractCurry = iContCurry
		glbal.GlRdocno = iGlRdocno
		initializers.DB.Save(&glbal)
	} else {
		iAmount := glbal.ContractAmount + temp
		// fmt.Println("I am inside update.....2", iAmount, glbal.ContractAmount)
		initializers.DB.Model(&glbal).Where("company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno).Update("contract_amount", iAmount)
	}
	results.Commit()
	return glbal.ContractAmount
}

// # 42
// Post GL - Get Transaction No and History Code
//
// ©  FuturaInsTech
func ValidateStatus(iCompany uint, iMethod string, iDate string, iStatus string) (string, string) {
	var permission models.Permission
	var result *gorm.DB
	oStatus := ""
	oHistory := ""

	result = initializers.DB.First(&permission, "company_id = ? and method = ?", iCompany, iMethod)
	if result.Error != nil {
		return oStatus, oHistory
	}

	iHistoryCode := permission.TransactionID
	var transaction models.Transaction
	result = initializers.DB.First(&transaction, "ID = ?", iHistoryCode)
	if result.Error != nil {
		return oStatus, oHistory
	}

	iHistoryCD := transaction.TranCode
	oHistory = iHistoryCD
	var p0029data paramTypes.P0029Data
	var extradata paramTypes.Extradata = &p0029data
	fmt.Println("Transaction Foound !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", iHistoryCode)
	err := GetItemD(int(iCompany), "P0029", iHistoryCD, iDate, &extradata)

	// fmt.Println("Newstatus", iStatus, p0029data.Statuses[0].CurrentStatus, p0029data.Statuses[0].ToBeStatus)
	if err != nil {
		return oStatus, oHistory
	} else {
		for i := 0; i < len(p0029data.Statuses); i++ {
			fmt.Println("Newstatus", iStatus, p0029data.Statuses[i].CurrentStatus,
				p0029data.Statuses[i].ToBeStatus)
			if iStatus == p0029data.Statuses[i].CurrentStatus {
				oStatus = p0029data.Statuses[i].ToBeStatus
				return oStatus, oHistory
			}
		}
	}
	return oStatus, oHistory
}

// # 43
// GetParamDesc - Get Long and Short Description of an item
//
// Inputs: Company, Param , Param ITem and Language
//
// # Outputs  Short Description, Long Description and Error
//
// ©  FuturaInsTech
func GetParamDesc(iCompany uint, iParam string, iItem string, iLanguage uint) (string, string, error) {
	var paramdesc models.ParamDesc

	results := initializers.DB.Where("company_id = ? AND name = ? and item = ? and language_id = ?", iCompany, iParam, iItem, iLanguage).Find(&paramdesc)
	if results.Error != nil || results.RowsAffected == 0 {

		return "", "", errors.New(" -" + strconv.FormatUint(uint64(iCompany), 10) + "-" + iParam + "-" + "-" + iItem + "-" + strconv.FormatUint(uint64(iLanguage), 10) + "-" + " is missing")
		//return errors.New(results.Error.Error())
	}
	return paramdesc.Shortdesc, paramdesc.Longdesc, nil
}

// # 44
// TDFBillD - Time Driven Function - Update Next Bill Date
//
// Inputs: Company, Policy, Function BILLD, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFBillD(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"
	initializers.DB.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	for i := 0; i < len(benefitenq); i++ {
		if benefitenq[i].BPremCessDate > odate {
			odate = benefitenq[i].BPremCessDate
		}
	}

	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if iRevFlag == "R" {
		var q0005data paramTypes.Q0005Data
		var extradataq0005 paramTypes.Extradata = &q0005data
		err := GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)
		if err != nil {
			return "", err
		}

		nxtBtdate := AddLeadDays(policy.PaidToDate, (-1 * q0005data.BillingLeadDays))
		policy.NxtBTDate = nxtBtdate
	}

	if result.Error != nil {
		return "", result.Error
	}

	if policy.PaidToDate >= odate {
		return "Date Exceeded", errors.New("Premium Cessation Date is Exceeded")
	}

	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = policy.NxtBTDate
		tdfpolicy.Tranno = iTranno
		tdfpolicy.Seqno = tdfrule.Seqno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	} else {
		initializers.DB.Delete(&tdfpolicy)
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = policy.NxtBTDate
		tdfpolicy.Tranno = iTranno

		initializers.DB.Create(&tdfpolicy)
		return "", nil
	}
}

func TDFBillDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string, txn *gorm.DB) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	for i := 0; i < len(benefitenq); i++ {
		if benefitenq[i].BPremCessDate > odate {
			odate = benefitenq[i].BPremCessDate
		}
	}

	result = txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	if iRevFlag == "R" {
		var q0005data paramTypes.Q0005Data
		var extradataq0005 paramTypes.Extradata = &q0005data
		err := GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)
		if err != nil {
			txn.Rollback()
			return "", err
		}

		nxtBtdate := AddLeadDays(policy.PaidToDate, (-1 * q0005data.BillingLeadDays))
		policy.NxtBTDate = nxtBtdate
	}

	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	if policy.PaidToDate >= odate {
		// return "Date Exceeded", errors.New("Premium Cessation Date is Exceeded")
		var tdfpolicyupd models.TDFPolicy
		result = txn.Find(&tdfpolicyupd, "company_id = ? AND policy_id = ? and tdf_type= ?", iCompany, iPolicy, "BILLD")
		if result.Error != nil {
			txn.Rollback()
			return "", nil
		}
		result = txn.Delete(&tdfpolicyupd)
		return "", nil
	}

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = policy.NxtBTDate
		tdfpolicy.Tranno = iTranno
		tdfpolicy.Seqno = tdfrule.Seqno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = policy.NxtBTDate
		tdfpolicy.Tranno = iTranno

		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		return "", nil
	}
}

// # 45
// TDFAnniD - Time Driven Function - Update Anniversary Date
//
// Inputs: Company, Policy, Function ANNID, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFAnniD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	result := initializers.DB.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", result.Error
	}
	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)

	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = policy.AnnivDate
		tdfpolicy.Tranno = iTranno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	} else {
		initializers.DB.Delete(&tdfpolicy)
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = policy.AnnivDate
		tdfpolicy.Tranno = iTranno

		initializers.DB.Create(&tdfpolicy)
		return "", nil
	}
}

func TDFAnniDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", result.Error
	}
	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = policy.AnnivDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = policy.AnnivDate
		tdfpolicy.Tranno = iTranno

		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	}
}

// # 46
// TDFReraD - Time Driven Function - ReRate Date Updation
//
// Inputs: Company, Policy, Function RERAD, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
// It is commented now.. It will be uncommented after Premium Increase Logic is added
func TDFReraD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	// var benefits []models.Benefit
	// var tdfpolicy models.TDFPolicy
	// var tdfrule models.TDFRule
	// var extraenq []models.Extra

	// oDate := ""

	// results := initializers.DB.Find(&extraenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	// if results.Error == nil {
	// 	if results.RowsAffected > 1 {
	// 		for i := 0; i < len(extraenq); i++ {
	// 			if oDate == "" {
	// 				oDate = extraenq[i].ToDate
	// 			}
	// 		}
	// 	}

	// }
	// initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	// result := initializers.DB.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	// if result.Error != nil {
	// 	return "", result.Error
	// }

	// for i := 0; i < len(benefits); i++ {
	// 	if benefits[i].BPremCessDate > benefits[i].BRerate {
	// 		if oDate == "" {
	// 			oDate = benefits[i].BRerate
	// 		}

	// 		if benefits[i].BRerate < oDate {
	// 			oDate = benefits[i].BRerate
	// 		}
	// 	}

	// }
	// if oDate != "" {
	// 	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	// 	if results.Error != nil {
	// 		tdfpolicy.CompanyID = iCompany
	// 		tdfpolicy.PolicyID = iPolicy
	// 		tdfpolicy.Seqno = tdfrule.Seqno
	// 		tdfpolicy.TDFType = iFunction
	// 		tdfpolicy.EffectiveDate = oDate
	// 		tdfpolicy.Tranno = iTranno
	// 		initializers.DB.Create(&tdfpolicy)
	// 		return "", nil
	// 	} else {
	// 		initializers.DB.Delete(&tdfpolicy)
	// 		var tdfpolicy models.TDFPolicy
	// 		tdfpolicy.CompanyID = iCompany
	// 		tdfpolicy.PolicyID = iPolicy
	// 		tdfpolicy.Seqno = tdfrule.Seqno
	// 		tdfpolicy.TDFType = iFunction
	// 		tdfpolicy.ID = 0
	// 		tdfpolicy.EffectiveDate = oDate
	// 		tdfpolicy.Tranno = iTranno

	// 		initializers.DB.Create(&tdfpolicy)
	// 		return "", nil
	// 	}
	// }
	return "", nil
}

func TDFReraDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	// var benefits []models.Benefit
	// var tdfpolicy models.TDFPolicy
	// var tdfrule models.TDFRule
	// var extraenq []models.Extra

	// oDate := ""

	// results := initializers.DB.Find(&extraenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	// if results.Error == nil {
	// 	if results.RowsAffected > 1 {
	// 		for i := 0; i < len(extraenq); i++ {
	// 			if oDate == "" {
	// 				oDate = extraenq[i].ToDate
	// 			}
	// 		}
	// 	}

	// }
	// initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	// result := initializers.DB.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	// if result.Error != nil {
	// 	return "", result.Error
	// }

	// for i := 0; i < len(benefits); i++ {
	// 	if benefits[i].BPremCessDate > benefits[i].BRerate {
	// 		if oDate == "" {
	// 			oDate = benefits[i].BRerate
	// 		}

	// 		if benefits[i].BRerate < oDate {
	// 			oDate = benefits[i].BRerate
	// 		}
	// 	}

	// }
	// if oDate != "" {
	// 	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	// 	if results.Error != nil {
	// 		tdfpolicy.CompanyID = iCompany
	// 		tdfpolicy.PolicyID = iPolicy
	// 		tdfpolicy.Seqno = tdfrule.Seqno
	// 		tdfpolicy.TDFType = iFunction
	// 		tdfpolicy.EffectiveDate = oDate
	// 		tdfpolicy.Tranno = iTranno
	// 		initializers.DB.Create(&tdfpolicy)
	// 		return "", nil
	// 	} else {
	// 		initializers.DB.Delete(&tdfpolicy)
	// 		var tdfpolicy models.TDFPolicy
	// 		tdfpolicy.CompanyID = iCompany
	// 		tdfpolicy.PolicyID = iPolicy
	// 		tdfpolicy.Seqno = tdfrule.Seqno
	// 		tdfpolicy.TDFType = iFunction
	// 		tdfpolicy.ID = 0
	// 		tdfpolicy.EffectiveDate = oDate
	// 		tdfpolicy.Tranno = iTranno

	// 		initializers.DB.Create(&tdfpolicy)
	// 		return "", nil
	// 	}
	// }
	return "", nil
}

// # 47
// TDFExpidD - Time Driven Function - Expiry Date Updation
//
// Inputs: Company, Policy, Function EXPID, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFExpiD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	if result.Error != nil {
		return "", result.Error
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		if benefits[i].BStatus != "EX" {
			iCoverage := benefits[i].BCoverage
			iDate := benefits[i].BStartDate
			var q0006data paramTypes.Q0006Data
			var extradataq0006 paramTypes.Extradata = &q0006data
			GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
			if q0006data.MatMethod == "" {
				if oDate == "" {
					oDate = benefits[i].BRiskCessDate
				}
				if benefits[i].BRiskCessDate < oDate {
					oDate = benefits[i].BRiskCessDate
				}
			}
		}
	}
	if oDate != "" {
		results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			initializers.DB.Create(&tdfpolicy)
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			initializers.DB.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

func TDFExpiDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	result = txn.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		if benefits[i].BStatus != "EX" {
			iCoverage := benefits[i].BCoverage
			iDate := benefits[i].BStartDate
			var q0006data paramTypes.Q0006Data
			var extradataq0006 paramTypes.Extradata = &q0006data
			err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
			if err != nil {
				txn.Rollback()
				return "", err
			}
			if q0006data.MatMethod == "" {
				if oDate == "" {
					oDate = benefits[i].BRiskCessDate
				}
				if benefits[i].BRiskCessDate < oDate {
					oDate = benefits[i].BRiskCessDate
				}
			}
		}
	}
	if oDate != "" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}

			return "", nil
		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}

			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}

			return "", nil
		}
	}
	return "", nil
}

// # 48
// TDFExpidS - Time Driven Function - Expiry Date Updation
//
// Inputs: Company, Policy, Function EXPID, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
// For Single Premium
func TDFExpiDS(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "SP")
	if result.Error != nil {
		return "", result.Error
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		if benefits[i].BStatus != "EX" {
			iCoverage := benefits[i].BCoverage
			iDate := benefits[i].BStartDate
			var q0006data paramTypes.Q0006Data
			var extradataq0006 paramTypes.Extradata = &q0006data
			GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
			if q0006data.MatMethod == "" {
				if oDate == "" {
					oDate = benefits[i].BRiskCessDate
				}
				if benefits[i].BRiskCessDate < oDate {
					oDate = benefits[i].BRiskCessDate
				}
			}
		}
	}
	if oDate != "" {
		results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			initializers.DB.Create(&tdfpolicy)
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			initializers.DB.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

func TDFExpiDSN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	result = txn.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "SP")
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		if benefits[i].BStatus != "EX" {
			iCoverage := benefits[i].BCoverage
			iDate := benefits[i].BStartDate
			var q0006data paramTypes.Q0006Data
			var extradataq0006 paramTypes.Extradata = &q0006data
			GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
			if q0006data.MatMethod == "" {
				if oDate == "" {
					oDate = benefits[i].BRiskCessDate
				}
				if benefits[i].BRiskCessDate < oDate {
					oDate = benefits[i].BRiskCessDate
				}
			}
		}
	}
	if oDate != "" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)

		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			result = txn.Create(&tdfpolicy)

			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			return "", nil
		} else {
			result = txn.Delete(&tdfpolicy)

			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			result = txn.Create(&tdfpolicy)

			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			return "", nil
		}
	}
	return "", nil
}

// # 49
// TDFReraD - Time Driven Function - Expiry Date Updation
//
// Inputs: Company, Policy, Function EXPID, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFMatD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	if result.Error != nil {
		return "", result.Error
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		iCoverage := benefits[i].BCoverage
		iDate := benefits[i].BStartDate
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data
		GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
		if q0006data.MatMethod != "" {
			if oDate == "" {
				oDate = benefits[i].BRiskCessDate
			}
			if benefits[i].BRiskCessDate < oDate {
				oDate = benefits[i].BRiskCessDate
			}
		}
	}
	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if oDate != "" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			initializers.DB.Create(&tdfpolicy)
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			initializers.DB.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

func TDFMatDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	result = txn.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		iCoverage := benefits[i].BCoverage
		iDate := benefits[i].BStartDate
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data
		GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
		if q0006data.MatMethod != "" {
			if oDate == "" {
				oDate = benefits[i].BRiskCessDate
			}
			if benefits[i].BRiskCessDate < oDate {
				oDate = benefits[i].BRiskCessDate
			}
		}
	}
	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if oDate != "" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			result = txn.Create(&tdfpolicy)

			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			return "", nil
		} else {
			result = txn.Delete(&tdfpolicy)

			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			result = txn.Create(&tdfpolicy)

			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			return "", nil
		}
	}
	return "", nil
}

// # 50
// TDFSurbD - Time Driven Function - Survival Benefit Date Updation
//
// Inputs: Company, Policy, Function SURVB, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFSurvbD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var survb models.SurvB
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	results := initializers.DB.First(&survb, "company_id = ? and policy_id = ? and paid_date = ?", iCompany, iPolicy, "")

	if results.Error != nil {
		return "", results.Error
	}
	result := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ? ", iCompany, iPolicy, iFunction)

	if result.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = survb.EffectiveDate
		tdfpolicy.Tranno = iTranno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	} else {
		initializers.DB.Delete(&tdfpolicy)
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = survb.EffectiveDate
		tdfpolicy.Tranno = iTranno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	}
}

func TDFSurvbDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var survb models.SurvB
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error

	}
	result = txn.First(&survb, "company_id = ? and policy_id = ? and paid_date = ?", iCompany, iPolicy, "")

	if result.Error != nil {
		//	txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ? ", iCompany, iPolicy, iFunction)

	if result.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = survb.EffectiveDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = survb.EffectiveDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	}
}

// # 51
// Survival Benefit Creation
//
// Inputs: Company, Coverage, Date (Inception) in YYYYMMDD, SA, Type as A/T (Age or Term), Method of SB, Term of the Policy, Age at Inception, Trannsaction No
//
// # Outputs SB Rates are creatd in SURVB Table
//
// ©  FuturaInsTech
func SBCreate(iCompany uint, iPolicy uint, iBenefitID uint, iCoverage string, iDate string, iSA float64, iType string, iMethod string, iYear int, iAge int, iTranno uint) error {

	var survb models.SurvB
	fmt.Println("Values", iCompany, iPolicy, iBenefitID, iCoverage, iDate, iSA, iType, iMethod, iYear, iAge, iTranno)
	if iType == "T" {
		var q0012data paramTypes.Q0012Data
		var extradataq0012 paramTypes.Extradata = &q0012data
		// fmt.Println("SB Parameters", iCompany, iType, iMethod, iYear, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0012", iMethod, iDate, &extradataq0012)
		fmt.Println("I am inside Term Based ")
		if err != nil {
			return err

		}
		// fmt.Println(q0012data.SBRates[0].Percentage)
		for x1 := 0; x1 <= iYear; x1++ {
			fmt.Println("X1Values are ", x1)
			for i := 0; i < len(q0012data.SbRates); i++ {
				fmt.Println("i Values are ", x1, i)
				if x1 == int(q0012data.SbRates[i].Term) {
					oSB := q0012data.SbRates[i].Percentage * iSA / 100
					// Write it in SB Table
					fmt.Println("Values of X and I", x1, i, iYear)
					survb.CompanyID = iCompany
					survb.PolicyID = iPolicy
					survb.PaidDate = ""
					survb.EffectiveDate = AddYears2Date(iDate, x1, 0, 0)
					survb.SBPercentage = q0012data.SbRates[i].Percentage
					survb.Amount = oSB
					survb.Tranno = uint(iTranno)
					survb.Sequence++
					survb.BenefitID = iBenefitID
					survb.ID = 0
					err1 := initializers.DB.Create(&survb)
					if err1.Error != nil {
						fmt.Println("I am inside Error")
						return err1.Error
					}

				}

			}

		}
	}
	if iType == "A" {
		var q0013data paramTypes.Q0013Data
		var extradataq0013 paramTypes.Extradata = &q0013data
		fmt.Println("SB Parameters", iCompany, iType, iMethod, iAge, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0013", iMethod, iDate, &extradataq0013)
		fmt.Println("SB Parameters", iCompany, iCoverage, iDate)

		if err != nil {
			return err

		}
		fmt.Println(q0013data.SbRates[0].Percentage)
		for x := 0; x <= iAge; x++ {
			for i := 0; i < len(q0013data.SbRates); i++ {
				if x == int(q0013data.SbRates[i].Age) {
					oSB := q0013data.SbRates[i].Percentage * iSA / 100
					// Write it in SB Table
					survb.CompanyID = iCompany
					survb.PolicyID = iPolicy
					survb.PaidDate = ""
					survb.EffectiveDate = AddYears2Date(iDate, x, 0, 0)
					survb.SBPercentage = q0013data.SbRates[i].Percentage
					survb.Amount = oSB
					survb.Tranno = uint(iTranno)
					survb.Sequence++
					survb.BenefitID = iBenefitID
					err1 := initializers.DB.Create(&survb)
					if err1 != nil {
						return err1.Error
					}
					continue
				}

			}

		}

	}
	return nil
}

func SBCreateN(iCompany uint, iPolicy uint, iBenefitID uint, iCoverage string, iDate string, iSA float64, iType string, iMethod string, iYear int, iAge int, iTranno uint, txn *gorm.DB) error {

	var survb models.SurvB
	fmt.Println("Values", iCompany, iPolicy, iBenefitID, iCoverage, iDate, iSA, iType, iMethod, iYear, iAge, iTranno)
	if iType == "T" {
		var q0012data paramTypes.Q0012Data
		var extradataq0012 paramTypes.Extradata = &q0012data
		// fmt.Println("SB Parameters", iCompany, iType, iMethod, iYear, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0012", iMethod, iDate, &extradataq0012)
		fmt.Println("I am inside Term Based ")
		if err != nil {
			txn.Rollback()
			return err

		}
		// fmt.Println(q0012data.SBRates[0].Percentage)
		for x1 := 0; x1 <= iYear; x1++ {
			fmt.Println("X1Values are ", x1)
			for i := 0; i < len(q0012data.SbRates); i++ {
				fmt.Println("i Values are ", x1, i)
				if x1 == int(q0012data.SbRates[i].Term) {
					oSB := q0012data.SbRates[i].Percentage * iSA / 100
					// Write it in SB Table
					fmt.Println("Values of X and I", x1, i, iYear)
					survb.CompanyID = iCompany
					survb.PolicyID = iPolicy
					survb.PaidDate = ""
					survb.EffectiveDate = AddYears2Date(iDate, x1, 0, 0)
					survb.SBPercentage = q0012data.SbRates[i].Percentage
					survb.Amount = oSB
					survb.Tranno = uint(iTranno)
					survb.Sequence++
					survb.BenefitID = iBenefitID
					survb.ID = 0
					err1 := txn.Create(&survb)
					if err1.Error != nil {
						txn.Rollback()
						fmt.Println("I am inside Error")
						return err1.Error
					}

				}

			}

		}
	}
	if iType == "A" {
		var q0013data paramTypes.Q0013Data
		var extradataq0013 paramTypes.Extradata = &q0013data
		fmt.Println("SB Parameters", iCompany, iType, iMethod, iAge, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0013", iMethod, iDate, &extradataq0013)
		fmt.Println("SB Parameters", iCompany, iCoverage, iDate)

		if err != nil {
			return err

		}
		fmt.Println(q0013data.SbRates[0].Percentage)
		for x := 0; x <= iAge; x++ {
			for i := 0; i < len(q0013data.SbRates); i++ {
				if x == int(q0013data.SbRates[i].Age) {
					oSB := q0013data.SbRates[i].Percentage * iSA / 100
					// Write it in SB Table
					survb.CompanyID = iCompany
					survb.PolicyID = iPolicy
					survb.PaidDate = ""
					survb.EffectiveDate = AddYears2Date(iDate, x, 0, 0)
					survb.SBPercentage = q0013data.SbRates[i].Percentage
					survb.Amount = oSB
					survb.Tranno = uint(iTranno)
					survb.Sequence++
					survb.BenefitID = iBenefitID
					err1 := txn.Create(&survb)
					if err1 != nil {
						txn.Rollback()
						return err1.Error
					}
					continue
				}

			}

		}

	}
	return nil
}

// # 52 (Redundant)  It is replaced by Create Communication
// LetterCreate - To Create Letters
//
// Inputs: Company,  Policy, Product, History Code and iDate
//
// # Stamp Duty Value
//
// # Parameter Used P0033 and P0034
//
// ©  FuturaInsTech
func LetterCreate(iCompany int, iPolicy uint, iTransaction string, iDate string, idata map[string]interface{}) {
	var policy models.Policy
	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data

	results := initializers.DB.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if results.Error != nil {
		return
	}
	iKey := iTransaction + policy.PProduct

	err1 := GetItemD(iCompany, "P0034", iKey, iDate, &extradatap0034)
	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(iCompany, "P0034", iKey, iDate, &extradatap0034)
		if err1 != nil {
			return
		}

	}
	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			err := GetItemD(iCompany, "P0033", iKey, iDate, &extradatap0033)
			if err != nil {
				return
			}
			var communication models.Communication
			communication.AgencyID = policy.AgencyID
			communication.AgentEmailAllowed = p0033data.AgentEmailAllowed
			communication.AgentSMSAllowed = p0033data.AgentSMSAllowed
			communication.AgentWhatsAppAllowed = p0033data.AgentWhatsAppAllowed

			communication.ClientID = policy.ClientID
			communication.PolicyID = policy.ID
			communication.CompanyID = uint(iCompany)

			communication.EmailAllowed = p0033data.EmailAllowed
			communication.SMSAllowed = p0033data.SMSAllowed
			communication.WhatsAppAllowed = p0033data.WhatsAppAllowed
			communication.DepartmentHead = p0033data.DepartmentHead
			communication.DepartmentName = p0033data.DepartmentName
			communication.CompanyPhone = p0033data.CompanyPhone
			communication.CompanyEmail = p0033data.CompanyEmail
			communication.Tranno = policy.Tranno
			communication.TemplateName = iKey
			communication.EffectiveDate = policy.PRCD
			communication.ExtractedData = idata
			communication.PDFPath = p0034data.Letters[i].PdfLocation
			communication.TemplatePath = p0034data.Letters[i].ReportTemplateLocation

			results := initializers.DB.Create(&communication)
			if results.Error != nil {
				return
			}

		}

	}

}

// # 53
// CalculateStampDuty - To Calculate Stamp Duty
//
// Inputs: Company,  Coverage, Date in YYYYMMDD and SA YYYYMMDD generally Inception Date
//
// # Stamp Duty Value
//
// # Parameter USed P0036  This Function has to be used during NB as well as when SA increased
//
// ©  FuturaInsTech
func CalculateStampDuty(iCompany uint, iCoverage string, iInstalment int, iDate string, iSA float64) float64 {
	var p0036data paramTypes.P0036Data
	var extradata paramTypes.Extradata = &p0036data
	iKey := iCoverage
	// fmt.Println("i key ", iKey)
	err := GetItemD(int(iCompany), "P0036", iKey, iDate, &extradata)
	if err != nil {
		return 0

	}
	for i := 0; i < len(p0036data.StampDuties); i++ {
		if iInstalment <= p0036data.StampDuties[i].Noofinstalments {
			if iSA < p0036data.StampDuties[i].Sa {
				oStampDuty := p0036data.StampDuties[i].Rate * iSA
				oStampDuty = RoundFloat(oStampDuty, 2)
				return oStampDuty
			}
		}
	}
	return 0

}

// # 54
// GetGlBal - To Get GL Balance for a given account code
//
// Inputs: Company,  Policy, GL Account Code
//
// # Output GL Amount
//
// ©  FuturaInsTech
func GetGlBal(iCompany uint, iPolicy uint, iGlaccount string) float64 {
	var glbal models.GlBal
	result := initializers.DB.Find(&glbal, "company_id = ? and gl_rdocno = ? and gl_accountno = ?", iCompany, iPolicy, iGlaccount)
	if result.Error != nil {
		return 0
	}
	return glbal.ContractAmount

}

// # 55
// GetTolerance - To Get Tolerance for a Given Freqquency
//
// Inputs: Company,  Transaciton Code, Currency, Product, Date
//
// # Output Tolerance Amount
//
// ©  FuturaInsTech
func GetTolerance(iCompany uint, iTransaction string, iCurrency string, iProduct string, iFrequency string, iDate string) float64 {
	var p0043data paramTypes.P0043Data
	var extradata paramTypes.Extradata = &p0043data
	iKey := iTransaction + iCurrency + iProduct
	// fmt.Println("i key ", iKey)
	err := GetItemD(int(iCompany), "P0043", iKey, iDate, &extradata)
	if err != nil {
		iKey := iTransaction + iCurrency
		err := GetItemD(int(iCompany), "P0043", iKey, iDate, &extradata)
		if err != nil {
			return 0
		}
	}

	for i := 0; i < len(p0043data.Frequencies); i++ {
		if p0043data.Frequencies[i].Frequency == iFrequency {
			return p0043data.Frequencies[i].Amount
		}

	}
	return 0
}

// # 56
// GetDeathAmount - Give Death Amount based on coverage and reason of death
//
// Inputs: Company Code, Policy, Coverage, Effective Date and cause of Death
//
// # Death Amount
//
// ©  FuturaInsTech
func GetDeathAmount(iCompany uint, iPolicy uint, iProduct string, iCoverage string, iEffectiveDate string, iCause string, iHistoryCD string) (oAmount float64) {
	var benefit models.Benefit
	result := initializers.DB.Find(&benefit, "company_id = ? and policy_id = ? and b_coverage = ?", iCompany, iPolicy, iCoverage)

	if result.Error != nil {
		oAmount = 0
		return
	}

	iSA := float64(benefit.BSumAssured)
	iStartDate := benefit.BStartDate
	iDate := benefit.BStartDate
	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data

	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)
	if err != nil {
		oAmount = 0
		return
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data

	err = GetItemD(int(iCompany), "Q0005", iProduct, iDate, &extradataq0005)
	if err != nil {
		oAmount = 0
		return
	}

	ideathMethod := q0006data.DeathMethod //DC001
	oAmount = 0
	var p0049data paramTypes.P0049Data
	var extradata1 paramTypes.Extradata = &p0049data
	iKey := iCause + iCoverage

	err = GetItemD(int(iCompany), "P0049", iKey, iDate, &extradata1)
	if err != nil {
		iKey = iCoverage
		err = GetItemD(int(iCompany), "P0049", iKey, iDate, &extradata1)
		if err != nil {
			oAmount = 0
		}

	}
	iNoofMonths := NewNoOfInstalments(iStartDate, iEffectiveDate)
	oPercentage := 100.00
	var oDeathMethod string

	for i := 0; i < len(p0049data.Months); i++ {
		if iNoofMonths <= int(p0049data.Months[i].Month) {
			oPercentage = p0049data.Months[i].Percentage
			oDeathMethod = p0049data.Months[i].DeathMethod
			break
		}
	}
	if oDeathMethod != "" { //DC006
		ideathMethod = oDeathMethod
	}
	iFund := 0.0

	switch {
	case ideathMethod == "DC001": // Return of SA
		oAmount = iSA
		break
	case ideathMethod == "DC002": // Return of FV
		if q0005data.NoLapseGuarantee == "Y" {
			if iNoofMonths <= q0005data.NoLapseGuaranteeMonths {
				iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
				if iFund <= 0 {
					oAmount = iSA
					break
				}
			}
		} else {
			iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
			oAmount = iFund
			break
		}
	case ideathMethod == "DC003": // Return of SA or Fund Value whichever is Highter
		iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
		if iSA >= iFund {
			oAmount = iSA
		} else {
			oAmount = iFund
		}
		break
	case ideathMethod == "DC004": // Return of SA + Fund Value

		iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)

		oAmount = iSA + iFund
		break
	case ideathMethod == "DC005": // Return of Premium Paid (All Coverages)
		var policy models.Policy
		initializers.DB.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		inoofinstalments := NewNoOfInstalments(policy.PRCD, policy.PaidToDate)
		switch {
		case policy.PFreq == "M":
			break
		case policy.PFreq == "Q":
			inoofinstalments = (inoofinstalments) / 3
			break
		case policy.PFreq == "H":
			inoofinstalments = (inoofinstalments) / 6
			break
		case policy.PFreq == "Y":
			inoofinstalments = (inoofinstalments) / 12
			break
		}
		oAmount = float64(inoofinstalments) * policy.InstalmentPrem
		//oAmount = GetPremiumPaid(policy.PRCD, policy.PaidToDate, policy.PFreq, policy.InstalmentPrem)
		break
	case ideathMethod == "DC006": // Return of Premium Paid (Given Coverages)
		oAmount = float64(benefit.BSumAssured)
		break
	case ideathMethod == "DC007": // Return of Premium Paid (All Coverages excluding Extra)
		oAmount = float64(benefit.BSumAssured)
		break
	case ideathMethod == "DC008": // Return of Premium Paid (Given  Coverage excluding extra)
		oAmount = float64(benefit.BSumAssured)
		break
	case ideathMethod == "MRTA1":
		a := benefit.BStartDate
		var noofyears int
		for {

			b := GetNextDue(a, "Y", "")
			c := Date2String(b)
			a = c
			noofyears++
			if a > iEffectiveDate {
				break
			}
		}
		var mrtaenq models.Mrta
		result := initializers.DB.First(&mrtaenq, "company_id = ? and policy_id = ? and prem_paying_term= ? ", iCompany, iPolicy, noofyears)
		if result.Error != nil {
			oAmount = 0
			return
		}
		oAmount = mrtaenq.BSumAssured
		return

	default:
		oAmount = 0
		return
	}
	oAmount = oAmount * oPercentage / 100
	return
}

// # 56
// GetDeathAmountN - Give Death Amount based on coverage and reason of death using txn
//
// Inputs: Company Code, Policy, Coverage, Effective Date and cause of Death
//
// # Death Amount
//
// ©  FuturaInsTech
func GetDeathAmountN(iCompany uint, iPolicy uint, iProduct string, iCoverage string, iEffectiveDate string, iCause string, iHistoryCD string, txn *gorm.DB) (oAmount float64) {
	var benefit models.Benefit
	result := txn.Find(&benefit, "company_id = ? and policy_id = ? and b_coverage = ?", iCompany, iPolicy, iCoverage)

	if result.Error != nil {
		oAmount = 0
		return
	}

	iSA := float64(benefit.BSumAssured)
	iStartDate := benefit.BStartDate
	iDate := benefit.BStartDate
	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data

	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)
	if err != nil {
		oAmount = 0
		return
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data

	err = GetItemD(int(iCompany), "Q0005", iProduct, iDate, &extradataq0005)
	if err != nil {
		oAmount = 0
		return
	}

	ideathMethod := q0006data.DeathMethod //DC001
	oAmount = 0
	var p0049data paramTypes.P0049Data
	var extradata1 paramTypes.Extradata = &p0049data
	iKey := iCause + iCoverage

	err = GetItemD(int(iCompany), "P0049", iKey, iDate, &extradata1)
	if err != nil {
		iKey = iCoverage
		err = GetItemD(int(iCompany), "P0049", iKey, iDate, &extradata1)
		if err != nil {
			oAmount = 0
		}

	}
	iNoofMonths := NewNoOfInstalments(iStartDate, iEffectiveDate)
	oPercentage := 100.00
	var oDeathMethod string

	for i := 0; i < len(p0049data.Months); i++ {
		if iNoofMonths <= int(p0049data.Months[i].Month) {
			oPercentage = p0049data.Months[i].Percentage
			oDeathMethod = p0049data.Months[i].DeathMethod
			break
		}
	}
	if oDeathMethod != "" { //DC006
		ideathMethod = oDeathMethod
	}
	iFund := 0.0

	switch {
	case ideathMethod == "DC001": // Return of SA
		oAmount = iSA
		break
	case ideathMethod == "DC002": // Return of FV
		if q0005data.NoLapseGuarantee == "Y" {
			if iNoofMonths <= q0005data.NoLapseGuaranteeMonths {
				oIlpMortality, oIlpFee := GetIlpMortalityFee(iCompany, benefit.ID)

				if oIlpMortality != 0 {
					err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, oIlpMortality, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
					if err != nil {
						txn.Rollback()
						return
					}

				}
				if oIlpFee != 0 {
					err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, oIlpFee, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
					if err != nil {
						txn.Rollback()
						return
					}
				}
				iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
				err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, iFund, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
				if err != nil {
					txn.Rollback()
					return
				}
				if iFund <= 0 {
					oAmount = iSA
					break
				}
			}
		} else {
			iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
			oAmount = iFund
			break
		}
	case ideathMethod == "DC003": // Return of SA or Fund Value whichever is Highter

		oIlpMortality, oIlpFee := GetIlpMortalityFee(iCompany, benefit.ID)

		if oIlpMortality != 0 {
			err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, oIlpMortality, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
			if err != nil {
				txn.Rollback()
				return
			}

		}
		if oIlpFee != 0 {
			err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, oIlpFee, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
			if err != nil {
				txn.Rollback()
				return
			}
		}
		iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
		err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, iFund, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
		if err != nil {
			txn.Rollback()
			return
		}
		if iSA >= iFund {
			oAmount = iSA
		} else {
			oAmount = iFund
		}
		break
	case ideathMethod == "DC004": // Return of SA + Fund Value

		//Mortality amount should come hear
		oIlpMortality, oIlpFee := GetIlpMortalityFee(iCompany, benefit.ID)

		if oIlpMortality != 0 {
			err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, oIlpMortality, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
			if err != nil {
				txn.Rollback()
				return
			}

		}
		if oIlpFee != 0 {
			err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, oIlpFee, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
			if err != nil {
				txn.Rollback()
				return
			}
		}
		iFund, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
		err = PostUlpDeductionByAmountN(iCompany, iPolicy, benefit.ID, iFund, iHistoryCD, iCoverage, iDate, iEffectiveDate, 0, "", txn)
		if err != nil {
			txn.Rollback()
			return
		}
		oAmount = iSA

		break
	case ideathMethod == "DC005": // Return of Premium Paid (All Coverages)
		var policy models.Policy
		txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		inoofinstalments := NewNoOfInstalments(policy.PRCD, policy.PaidToDate)
		switch {
		case policy.PFreq == "M":
			break
		case policy.PFreq == "Q":
			inoofinstalments = (inoofinstalments) / 3
			break
		case policy.PFreq == "H":
			inoofinstalments = (inoofinstalments) / 6
			break
		case policy.PFreq == "Y":
			inoofinstalments = (inoofinstalments) / 12
			break
		}
		oAmount = float64(inoofinstalments) * policy.InstalmentPrem
		//oAmount = GetPremiumPaid(policy.PRCD, policy.PaidToDate, policy.PFreq, policy.InstalmentPrem)
		break
	case ideathMethod == "DC006": // Return of Premium Paid (Given Coverages)
		oAmount = float64(benefit.BSumAssured)
		break
	case ideathMethod == "DC007": // Return of Premium Paid (All Coverages excluding Extra)
		oAmount = float64(benefit.BSumAssured)
		break
	case ideathMethod == "DC008": // Return of Premium Paid (Given  Coverage excluding extra)
		oAmount = float64(benefit.BSumAssured)
		break
	case ideathMethod == "MRTA1":
		a := benefit.BStartDate
		var noofyears int
		for {

			b := GetNextDue(a, "Y", "")
			c := Date2String(b)
			a = c
			noofyears++
			if a > iEffectiveDate {
				break
			}
		}
		var mrtaenq models.Mrta
		result := txn.First(&mrtaenq, "company_id = ? and policy_id = ? and prem_paying_term= ? ", iCompany, iPolicy, noofyears)
		if result.Error != nil {
			oAmount = 0
			return
		}
		oAmount = mrtaenq.BSumAssured
		return
	case ideathMethod == "DC009":
		// Important. Note. Annuity has multiple records, we need to pick up the latest record in the table
		var annuity models.Annuity
		result = txn.Last(&annuity, "policy_id = ?", iPolicy)
		if result.Error != nil {
			return
		}
		oIntrestRate := 6.00
		_, _, _, days, _, _, _, _ := NoOfDays(iEffectiveDate, annuity.AnnStartDate)
		inoofinstalments := NewNoOfInstalments(annuity.AnnStartDate, annuity.AnnCurrDate) + 1
		oCompoundint := CompoundInterest(iSA, oIntrestRate, float64(days))
		oPaidValue := inoofinstalments * int(annuity.AnnAmount)
		oAmount = iSA + oCompoundint - float64(oPaidValue)
		return
	default:
		oAmount = 0
		return
	}
	oAmount = oAmount * oPercentage / 100
	return
}

// # 57
// NewNofInstalments - Get No of instalments in Months
//
// Inputs: From Date, To Date
//
// # No of Instalments in Months
//
// ©  FuturaInsTech
func NewNoOfInstalments(iFromDate string, iToDate string) (oinstalment int) {
	tempFromDate := String2Date(iFromDate)
	tempToDate := String2Date(iToDate)
	a := tempFromDate
	for i := 0; i < 1000; i++ {
		a = AddMonth(a, 1)
		if a == tempToDate || a.After(tempToDate) {
			oinstalment = i + 1
			return oinstalment

		}
	}
	return
}

// # 58
// DateConvert - Convert Date into DD / MM / YYYY Format
//
// Inputs: Date in YYYYMMDD
//
// Outputs : DD/MM/YYYY
//
// ©  FuturaInsTech
func DateConvert(iDate string) (oDate string) {
	if iDate == "" {
		return iDate
	}
	dd := iDate[6:8]
	mm := iDate[4:6]
	yy := iDate[0:4]
	oDate = dd + "/" + mm + "/" + yy
	return

}

// # 59
// GetTranCode - To Get TranCode
//
// Inputs: Company,  Program/Function Name
//
// Outputs: Transaction Code
//
// ©  FuturaInsTech
func GetTranCode(iCompany uint, iDescription string) (otrancode string) {

	var result *gorm.DB
	var transaction models.Transaction

	result = initializers.DB.Find(&transaction, "company_id = ? and description = ?", iCompany, iDescription)
	if result.Error != nil {
		return "Not Found"
	}

	otrancode = transaction.TranCode
	return
}

// # 60
// WrapInArray
// Inputs: Interfaces
//
// Outputs: Interfaces
//
// ©  FuturaInsTech
func WrapInArray(obj interface{}) interface{} {
	sliceType := reflect.SliceOf(reflect.TypeOf(obj))
	slice := reflect.MakeSlice(sliceType, 1, 1)
	slice.Index(0).Set(reflect.ValueOf(obj))
	return slice.Interface()
}

// # 61
// Number Func  Convert Float to string
// Inputs: Float
//
// Outputs: String
//
// ©  FuturaInsTech
func NumberFunc(iAmount float64) (oAmount string) {

	p := message.NewPrinter(language.English)
	oAmount = p.Sprintf("%15.2f", iAmount)
	return
}

// # 62
func NumbertoPrint(iAmount float64) (oAmount string) {

	p := message.NewPrinter(language.English)
	oAmount = p.Sprintf("%15.2f", iAmount)
	return
}

// # 63
func IDtoPrint(iID uint) (oID string) {
	oID = strconv.FormatUint(uint64(iID), 10)
	fmt.Println(oID, reflect.TypeOf(oID))
	return
}

// # 64
// GetTotalGSTPercentage - Get Unit Linked Mortality Prem for a given duration
//
// Inputs: Company,  Coverage and Date String in YYYYMMDD (Current Date), Key is Coverage Code, No of Months, Amount to be charged
//
// # Outputs GST Amount
//
// ©  FuturaInsTech
func GetTotalGSTAmount(iCompany uint, iPolicy uint, iFromDate string, iToDate string) float64 {

	var policyenq models.Policy
	var oAmount float64
	oAmount = 0
	result := initializers.DB.First(&policyenq, "company_id =? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return 0
	}
	iFrequency := policyenq.PFreq
	//iRCD := policyenq.PRCD
	var benefitenq1 []models.Benefit

	results := initializers.DB.Find(&benefitenq1, "company_id =? and policy_id = ? ", iCompany, iPolicy)
	if results.Error != nil {
		return 0
	}

	for a := 0; a < len(benefitenq1); a++ {
		FromDate := iFromDate
		ToDate := iToDate
		iCovRcd := benefitenq1[a].BStartDate

		iKey := benefitenq1[a].BCoverage
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data

		err := GetItemD(int(iCompany), "Q0006", iKey, FromDate, &extradataq0006)
		if err != nil {
			return 0
		}
		if q0006data.PremCalcType != "U" {
			var q0023data paramTypes.Q0023Data
			var extradataq0023 paramTypes.Extradata = &q0023data

			iAmount := benefitenq1[a].BPrem
			for b := FromDate; b < ToDate; {
				// Get Premium Rate
				err := GetItemD(int(iCompany), "Q0023", iKey, FromDate, &extradataq0023)
				if err != nil {
					return 0
				}

				date := GetNextDue(FromDate, iFrequency, "")
				FromDate = Date2String(date)
				b = FromDate
				iMonths := NewNoOfInstalments(iCovRcd, FromDate)
				for i := 0; i < len(q0023data.Gst); i++ {
					if uint(iMonths) <= q0023data.Gst[i].Month {
						oAmount = float64(iAmount)*q0023data.Gst[i].Rate + oAmount
						oAmount = RoundFloat(oAmount, 2)
						break

					}
				}
			}
		}

	}
	return oAmount
}

// # 65
// GetMRTABenefit - SA, Interest Rate, Policy Year, Interim Period, Term
//
// Inputs: Company,  Coverage and Date String in YYYYMMDD (Current Date), Key is Coverage Code, No of Months, Amount to be charged
//
// # Outputs Benefit
//
// ©  FuturaInsTech
func GetMRTABen(iSA float64, iInterest float64, iPolYear float64, iInterimPeriod float64, iTerm float64) float64 {
	a := math.Pow((1 + ((iInterest / 100) / 12)), ((iPolYear - iInterimPeriod) * 12))
	b := math.Pow((1 + ((iInterest / 100) / 12)), (iTerm * 12))
	c := (1 - (a-1)/(b-1))
	oSA := RoundFloat(iSA*c, 2)
	return oSA

}

// #66
// GetMrtaPremO - MRTA Premium
//
// Inputs: Company,  Coverage , Age Gender , Term, Premium Paying Term, Prem Method, Date String in YYYYMMDD and Mortality
//
// # Outputs MRTA Premium and Error
//
// ©  FuturaInsTech
func GetMrtaPremO(iCompany uint, iPolicy uint, iCoverage string, iAge uint, iGender string, iTerm uint, iPremTerm uint, iPremMethod string, iDate string, iMortality string) (float64, error) {

	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)

	var q0010data paramTypes.Q0010Data
	var extradataq0010 paramTypes.Extradata = &q0010data
	var q0010key string
	var prem float64
	prem = 0
	var prem1 float64
	prem1 = 0
	term := strconv.FormatUint(uint64(iTerm), 10)
	var mrtaenq []models.Mrta

	result := initializers.DB.Find(&mrtaenq, "company_id = ? and policy_id =? and b_coverage = ?", iCompany, iPolicy, iCoverage)

	if result.Error != nil {
		return 0, result.Error
	}

	for x := 0; x < len(mrtaenq); x++ {
		// Single Premium it is always 1
		//premTerm := strconv.FormatUint(uint64(iTerm-mrtaenq[x].PremPayingTerm), 10)
		premTerm := "1"
		//fmt.Println("****************", iCompany, iCoverage, iAge, iGender, iTerm, iPremMethod, iDate, iMortality)
		if q0006data.PremCalcType == "A" {
			q0010key = iCoverage + iGender
		} else if q0006data.PremCalcType == "P" {
			q0010key = iCoverage + iGender + term + premTerm
			// END1 + Male + Term + Premium Term
		}
		fmt.Println("Premium Key ******", iCoverage, iGender, term, premTerm, q0006data.PremCalcType, q0010key)
		err := GetItemD(int(iCompany), "Q0010", q0010key, iDate, &extradataq0010)
		if err != nil {
			return 0, err

		}
		fmt.Println("************", iCompany, iAge, q0010key, iDate)

		for i := 0; i < len(q0010data.Rates); i++ {
			if q0010data.Rates[i].Age == uint(iAge) {
				prem = q0010data.Rates[i].Rate / 10000
				prem1 = prem*mrtaenq[x].BSumAssured + prem1
				iAge = iAge + 1
				break
			}
		}
	}
	prem = prem1
	fmt.Println("************", iCompany, iAge, q0010key, iDate, prem)
	return prem, nil

}

// #67
// RevGL Move - Reverse GL Move
//
// Inputs: Transaction No, User Company and Policy
//
// # Outputs Reversal Records
//
// ©  FuturaInsTech
func RevGlMove(tranno, userco, ipolicy float64) error {
	var glmoveenq []models.GlMove
	opol := strconv.Itoa(int(ipolicy))
	results := initializers.DB.Where("gl_rldg_acct LIKE ?", "%"+opol+"%").Find(&glmoveenq, "tranno = ? and company_id = ? ", tranno, userco)

	if results.Error != nil {
		return nil
	}

	for i := 0; i < len(glmoveenq); i++ {
		oglamount := glmoveenq[i].GlAmount * -1
		ocontractamt := glmoveenq[i].ContractAmount * -1
		iCompany := userco
		iContractCurry := glmoveenq[i].ContractCurry
		iEffectiveDate := glmoveenq[i].EffectiveDate
		iGlAmount := oglamount
		iAccAmount := ocontractamt
		iAccountCodeID := glmoveenq[i].AccountCodeID

		iGlRdocno, _ := strconv.Atoi(glmoveenq[i].GlRdocno)
		iGlRldgAcct := glmoveenq[i].GlRldgAcct
		iSeqnno := glmoveenq[i].SequenceNo
		iGlSign := glmoveenq[i].GlSign
		iAccountCode := glmoveenq[i].AccountCode
		iHistoryCode := glmoveenq[i].HistoryCode
		iTranno := tranno
		iRevInd := "R"
		iCoverage := glmoveenq[i].BCoverage
		//glmoveupd.UpdatedID = userid
		err := PostGlMove(uint(iCompany), iContractCurry, iEffectiveDate, int(iTranno), iGlAmount, iAccAmount, iAccountCodeID, uint(iGlRdocno), iGlRldgAcct, iSeqnno, iGlSign, iAccountCode, iHistoryCode, iRevInd, iCoverage)
		if err != nil {
			return err
		}
	}

	return nil
}

// #68
// Surrender Amount
//
// Inputs: Company,  Policy, Coverage , Effective Date, Term, Premium Term, Status, SA, Start Date, Surrender Method, NO of Instalments
//
// # Outputs Surrender Amount
//
// ©  FuturaInsTech
func GetSurrenderAmount(iCompany uint, iPolicy uint, iCoverage string, iEffectiveDate string, iTerm uint, iPremTerm uint, iStatus string, iSumAssured float64, iPaidTerm int, iStartDate string, iSurrMethod string, iInstallments int) (oAmount float64) {

	oAmount = 0
	var p0053data paramTypes.P0053Data
	var extradatap0053 paramTypes.Extradata = &p0053data

	iKey := iCoverage + iStatus + strconv.Itoa(int(iTerm)) + strconv.Itoa(int(iPremTerm))

	err := GetItemD(int(iCompany), "P0053", iKey, iStartDate, &extradatap0053)
	if err != nil {
		oAmount = 0
		return oAmount

	}
	a := uint64(iPaidTerm)
	b := uint64(iPremTerm)
	var c float64
	c = float64(a) / float64(b)

	switch {
	case iSurrMethod == "SM001": // Return of SA
		for i := 0; i < len(p0053data.Rates); i++ {
			if iInstallments <= int(p0053data.Rates[i].Month) {

				oAmount = (p0053data.Rates[i].Percentage / 100) * iSumAssured * c
				return oAmount
			}
		}

	case iSurrMethod == "SM002": // Return of FV
		oAmount = 0
		break
	case iSurrMethod == "SM003": // Return of SA or Fund Value whichever is Highter
		oAmount = 0
		break
	case iSurrMethod == "SM004":
		var annuity models.Annuity
		txn := initializers.DB.Begin()
		result := txn.Last(&annuity, "policy_id = ?", iPolicy)
		if result.Error != nil {
			return 0
		}
		oIntrestRate := 6.00
		_, _, _, days, _, _, _, _ := NoOfDays(iEffectiveDate, annuity.AnnStartDate)
		inoofinstalments := NewNoOfInstalments(annuity.AnnStartDate, annuity.AnnCurrDate) + 1
		oCompoundint := CompoundInterest(iSumAssured, oIntrestRate, float64(days))
		oPaidValue := inoofinstalments * int(annuity.AnnAmount)
		oAmount = iSumAssured + oCompoundint - float64(oPaidValue)
	default:
		oAmount = 0
		return
	}
	return
}

// #69
// CalculateStampDutyByPolicy - Stamp Duty for a Policy
//
// Inputs: Company,  Policy
//
// # Outputs Stamp Duty
//
// ©  FuturaInsTech
func CalculateStampDutyByPolicy(iCompanyId uint, iPolicyId uint) float64 {

	tStampDuty := 0.0
	var policyenq models.Policy
	result := initializers.DB.First(&policyenq, "company_id =? and id = ?", iCompanyId, iPolicyId)
	iDate := policyenq.PRCD

	if result.Error != nil {
		return 0.0

	}
	if policyenq.PolStatus == "PV" || policyenq.PolStatus == "PC" || policyenq.PolStatus == "UW" {
		policyenq.PaidToDate = policyenq.PRCD
	}

	var benefitsenq []models.Benefit
	results := initializers.DB.Find(&benefitsenq, "company_id = ? and policy_id = ? ", iCompanyId, iPolicyId)

	if results.Error != nil {
		return 0.0
	}

	for i := 0; i < len(benefitsenq); i++ {
		iCoverage := benefitsenq[i].BCoverage
		FromDate := benefitsenq[i].BStartDate
		iCompany := benefitsenq[i].CompanyID

		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data

		err := GetItemD(int(iCompany), "Q0006", iCoverage, FromDate, &extradataq0006)
		if err != nil {
			return 0
		}
		if q0006data.PremCalcType != "U" {

			iSA := benefitsenq[i].BSumAssured
			iInstalmentPaid := GetNoIstalments(benefitsenq[i].BStartDate, policyenq.PaidToDate, policyenq.PFreq)

			iStampDuty := CalculateStampDuty(iCompanyId, iCoverage, iInstalmentPaid, iDate, float64(iSA))

			tStampDuty = tStampDuty + iStampDuty

		}
	}

	return tStampDuty

}

// #70
// NoOfDays - Get No of Days between two dates
//
// Inputs: From and To Dates
//
// # Outputs Year, Month, Week, Days , hrs , minutes, seconds, millie seconds and nano seconds
//
// ©  FuturaInsTech
func NoOfDays(startDate string, endDate string) (year int64, month int64, week int64, days int64, hrs float64, mm float64, ss float64, nss int64) {
	a := String2Date(startDate)
	b := String2Date(endDate)
	difference := a.Sub(b)

	year = int64(difference.Hours() / 24 / 365)
	month = int64(difference.Hours() / 24 / 30)
	week = int64(difference.Hours() / 24 / 7)
	days = int64(difference.Hours() / 24)
	hrs = difference.Hours()
	mm = difference.Minutes()
	ss = difference.Seconds()
	nss = difference.Nanoseconds()
	return year, month, week, days, hrs, mm, ss, nss

}

// #80
// Get Business Date
//
// Inputs: Company,  User and Departemet
//
// # Outputs Business Date
//
// ©  FuturaInsTech
// 01 - NB
// 02 - Cash and Payment
// 03 - Maturity
// 04 - Death Claim
// 05 - Customer Service
// 06
func GetBusinessDate(iCompany uint, iUser uint, iDepartment uint) (oDate string) {
	var businessdate models.BusinessDate
	// Get with User
	result := initializers.DB.Find(&businessdate, "company_id = ? and user_id = ? and department = ? and user_id IS NOT NULL and department IS NOT NULL", iCompany, iUser, iDepartment)
	if result.RowsAffected == 0 {
		// If User Not Found, get with Department
		result = initializers.DB.Find(&businessdate, "company_id = ? and department = ? and user_id IS NULL ", iCompany, iDepartment)
		if result.RowsAffected == 0 {
			// If Department Not Found, get with comapny
			result = initializers.DB.Find(&businessdate, "company_id = ? and department IS NULL and user_id IS NULL", iCompany)
			if result.RowsAffected == 0 {
				return Date2String(time.Now())

			} else {
				oDate := businessdate.Date
				return oDate
			}
		} else {
			oDate := businessdate.Date
			return oDate
		}

	} else {
		oDate := businessdate.Date
		return oDate
	}

}

// #81
// TDFLapsD - Time Driven Function - Update Lapse Date as per Q0005 Parameter
//
// Inputs: Company, Policy, Function LAPSD, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFLapsD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.Error != nil {
		return "", result.Error
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	err := GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)

	if err != nil {
		return "", err
	}
	iLapsedDate := AddLeadDays(policy.PaidToDate, q0005data.LapsedDays)

	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = iLapsedDate
		tdfpolicy.Tranno = iTranno
		tdfpolicy.Seqno = tdfrule.Seqno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	} else {
		initializers.DB.Delete(&tdfpolicy)
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = iLapsedDate
		tdfpolicy.Tranno = iTranno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	}
}

func TDFLapsDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	err := GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)

	if err != nil {
		txn.Rollback()
		return "", err
	}
	iLapsedDate := AddLeadDays(policy.PaidToDate, q0005data.LapsedDays)

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = iLapsedDate
		tdfpolicy.Tranno = iTranno
		tdfpolicy.Seqno = tdfrule.Seqno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = iLapsedDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		return "", nil
	}

}

// #82
// TdfhUpdate - Time Driven Function - Update TDF Header File
//
// Inputs: Company, Policy
//
// # It has to loop through TDFPOLICIES and update earliest due in Tdfh
//
// # Outputs  Old Record is Soft Deleted and New Record is Created in TDFH
//
// ©  FuturaInsTech
func TdfhUpdate(iCompany uint, iPolicy uint) error {
	var tdfhupd models.Tdfh
	var tdfpolicyenq []models.TDFPolicy
	iDate := "29991231"

	results := initializers.DB.Find(&tdfpolicyenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if results.Error != nil {

		return results.Error
	}
	for i := 0; i < len(tdfpolicyenq); i++ {
		if tdfpolicyenq[i].EffectiveDate <= iDate {
			iDate = tdfpolicyenq[i].EffectiveDate
		}
	}
	result := initializers.DB.Find(&tdfhupd, "company_id =? and policy_id = ?", iCompany, iPolicy)

	if result.Error == nil {
		if result.RowsAffected == 0 {
			tdfhupd.CompanyID = iCompany
			tdfhupd.PolicyID = iPolicy
			tdfhupd.EffectiveDate = iDate
			result = initializers.DB.Create(&tdfhupd)
		} else {
			result = initializers.DB.Delete(&tdfhupd)
			var tdfhupd models.Tdfh
			tdfhupd.CompanyID = iCompany
			tdfhupd.PolicyID = iPolicy
			tdfhupd.EffectiveDate = iDate
			tdfhupd.ID = 0
			result = initializers.DB.Create(&tdfhupd)
		}

	}
	return nil
}

// #82
// TdfhUpdateN - Time Driven Function - Update TDF Header File
//
// Inputs: Company, Policy  (New Version with Rollback)
//
// # It has to loop through TDFPOLICIES and update earliest due in Tdfh
//
// # Outputs  Old Record is Soft Deleted and New Record is Created in TDFH
//
// ©  FuturaInsTech
func TdfhUpdateN(iCompany uint, iPolicy uint, txn *gorm.DB) error {
	var tdfhupd models.Tdfh
	var tdfpolicyenq []models.TDFPolicy

	iDate := "29991231"

	results := txn.Find(&tdfpolicyenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if results.Error != nil {
		txn.Rollback()
		return results.Error
	}
	for i := 0; i < len(tdfpolicyenq); i++ {
		if tdfpolicyenq[i].EffectiveDate <= iDate {
			iDate = tdfpolicyenq[i].EffectiveDate
		}
	}
	result := txn.Find(&tdfhupd, "company_id =? and policy_id = ?", iCompany, iPolicy)

	if result.Error == nil {
		if result.RowsAffected == 0 {
			tdfhupd.CompanyID = iCompany
			tdfhupd.PolicyID = iPolicy
			tdfhupd.EffectiveDate = iDate
			result = txn.Create(&tdfhupd)
			if result.Error != nil {
				txn.Rollback()
				return results.Error
			}
		} else {
			result = txn.Delete(&tdfhupd)
			var tdfhupd models.Tdfh
			tdfhupd.CompanyID = iCompany
			tdfhupd.PolicyID = iPolicy
			tdfhupd.EffectiveDate = iDate
			tdfhupd.ID = 0
			result = txn.Create(&tdfhupd)
			if result.Error != nil {
				txn.Rollback()
				return results.Error
			}
		}

	}
	return nil
}

// #83
// TDFColl - Time Driven Function - Create Collection Record in TDF
//
// Inputs: Company, Policy, Function CollD and iDate (Which is passed).
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFCollD(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iDate string) (string, error) {
	//iBusinssdate := GetBusinessDate(iCompany, 1, "02")
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.Error != nil {
		return "", result.Error
	}
	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = iDate
		tdfpolicy.Tranno = policy.Tranno
		tdfpolicy.Seqno = tdfrule.Seqno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	} else {
		initializers.DB.Delete(&tdfpolicy)
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = iDate
		tdfpolicy.Tranno = policy.Tranno

		initializers.DB.Create(&tdfpolicy)
		return "", nil
	}
}
func TDFCollDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iDate string, txn *gorm.DB) (string, error) {
	//iBusinssdate := GetBusinessDate(iCompany, 1, "02")
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {

		return "", result.Error
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.Error != nil {
		return "", result.Error
	}
	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = iDate
		tdfpolicy.Tranno = policy.Tranno
		tdfpolicy.Seqno = tdfrule.Seqno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			return "", result.Error
		}

		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {

			return "", result.Error
		}
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = iDate
		tdfpolicy.Tranno = policy.Tranno

		result = txn.Create(&tdfpolicy)
		if result.Error != nil {

			return "", result.Error
		}
		return "", nil
	}
}

// #84 Get Future Date (Redundant)  GetNextDue Could be used
//
// Inputs: From Date, Todate and Frequency
//
// # Outputs  New Date
//
// ©  FuturaInsTech
func GetFutureDue(iFromDate string, iToDate string, iFreq string) (oDate string) {
	a := iFromDate
	for {
		nxtDate := GetNextDue(a, iFreq, "")
		a = Date2String(nxtDate)
		if a > iToDate {
			break
		}
		return a
	}

	return a
}

// #85
// GetMaturityAmount - Get Maturity Amount
// Inputs: Company, Policy, Coverage, Effective Date
//
// # Outputs  Maturity Amount
//
// ©  FuturaInsTech
func GetMaturityAmount(iCompany uint, iPolicy uint, iCoverage string, iEffectiveDate string) (oAmount float64) {
	var benefit models.Benefit
	result := initializers.DB.Find(&benefit, "company_id = ? and policy_id = ? and b_coverage = ?", iCompany, iPolicy, iCoverage)

	if result.Error != nil {
		oAmount = 0
		return
	}

	// iFund := float64(70000.00)
	iSA := float64(benefit.BSumAssured)
	// iStartDate := benefit.BStartDate
	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	iDate := benefit.BStartDate

	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)
	if err != nil {
		oAmount = 0
		return
	}

	imatMethod := q0006data.MatMethod //MAT001
	oAmount = 0

	switch {
	case imatMethod == "MAT001": // Return of SA
		oAmount = iSA
		break
	case imatMethod == "MAT002": // Double SA
		oAmount = iSA * 2
		break
	case imatMethod == "MAT003": // Return of Premium (To be Developed)
		oAmount = iSA * 2
		break
	case imatMethod == "MAT004": // No Maturity Value
		oAmount = 0
		break
	case imatMethod == "MAT005": // No Maturity Value
		oAmount = 0
		break
	case imatMethod == "MAT006": // No Maturity Value
		oAmount = 0
		break
	case imatMethod == "MAT007": // Return of Final Survival Benefit Amount
		// Survival Benefit Amount is already paid through TDF
		// So Maturity Amount is set to Zero
		var survb models.SurvB
		result := initializers.DB.Find(&survb, "company_id = ? and policy_id = ? and b_coverage = ? and paid_date = ?", iCompany, iPolicy, iCoverage, "")
		if result.Error != nil {
			oAmount = 0
		} else {
			oAmount = survb.Amount
		}
		return oAmount
	// case imatMethod == "MAT008": // ilp fund value
	// 	iFund, _, _ := GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
	// 	oAmount = iFund
	// 	return oAmount
	// case imatMethod == "MAT009": // ilp fund value + Sumassured
	// 	iFund, _, _ := GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
	// 	oAmount = iFund + iSA
	// 	return oAmount
	// case imatMethod == "MAT010": // ilp fund value or Sumassured which one is greater
	// 	iFund, _, _ := GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)
	// 	if iFund < iSA {
	// 		oAmount = iSA
	// 	} else {
	// 		oAmount = iFund
	// 	}
	// 	return oAmount
	case imatMethod == "MAT099": // No Maturity Value
		oAmount = 0
		break
	default:
		oAmount = 0
		return
	}

	return oAmount
}

// #97  (Redundant )
// func GetSurrDData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {

// 	var surrdenq []models.SurrD

// 	initializers.DB.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
// 	surrdarray := make([]interface{}, 0)

// 	for k := 0; k < len(surrdenq); k++ {
// 		resultOut := map[string]interface{}{
// 			"ID":              IDtoPrint(surrdenq[k].ID),
// 			"PolicyID":        IDtoPrint(surrdenq[k].PolicyID),
// 			"ClientID":        IDtoPrint(surrdenq[k].ClientID),
// 			"BenefitID":       IDtoPrint(surrdenq[k].ID),
// 			"BCoverage":       surrdenq[k].BCoverage,
// 			"BSumAssured":     surrdenq[k].BSumAssured,
// 			"SurrAmount":      float64(surrdenq[k].SurrAmount),
// 			"RevBonus":        float64(surrdenq[k].RevBonus),
// 			"AddlBonus":       float64(surrdenq[k].AddlBonus),
// 			"TerminalBonus":   float64(surrdenq[k].TerminalBonus),
// 			"InterimBonus":    float64(surrdenq[k].InterimBonus),
// 			"LoyaltyBonus":    float64(surrdenq[k].LoyaltyBonus),
// 			"OtherAmount":     float64(surrdenq[k].OtherAmount),
// 			"AccumDividend":   float64(surrdenq[k].AccumDividend),
// 			"AccumDivInt":     float64(surrdenq[k].AccumDivInt),
// 			"TotalFundValue":  float64(surrdenq[k].TotalFundValue),
// 			"TotalSurrAmount": float64(surrdenq[k].TotalSurrAmount),
// 		}
// 		surrdarray = append(surrdarray, resultOut)
// 	}

// 	return surrdarray

// }

// #103
// Check Status
//
// # This function, take company code, history code, date and status as inputs
//
// # It returns status which is boolean and also output status which is string
//
// ©  FuturaInsTech
func CheckStatus(iCompany uint, iHistoryCD string, iDate string, iStatus string) (status bool, oStatus string) {
	var p0029data paramTypes.P0029Data
	var extradata paramTypes.Extradata = &p0029data

	err := GetItemD(int(iCompany), "P0029", iHistoryCD, iDate, &extradata)
	if err != nil {
		return true, ""
	}
	for i := 0; i < len(p0029data.Statuses); i++ {
		if iStatus == p0029data.Statuses[i].CurrentStatus {
			return false, p0029data.Statuses[i].ToBeStatus
		}
	}
	return true, ""
}

// #105
// Get Name
//
// # This function, Return Name of the Client in Long Name + Short Name + Sur Name Format
//
// #  Input Variables Company Code and Client Code
// #  Return is Name
//
// ©  FuturaInsTech
func GetName(iCompany uint, iClient uint) string {
	var clientenq models.Client
	oName := ""
	result := initializers.DB.Find(&clientenq, "company_id = ? and id = ?", iCompany, iClient)

	if result.Error != nil {
		return oName
	}
	oName = clientenq.ClientLongName + " " + clientenq.ClientShortName + " " + clientenq.ClientSurName
	return oName
}

// #106
// New Function to do Amount in Words in Receipts
//
// Inputs: Company, Amount, Currency
//
// # Outputs  Values in Words
// ©  FuturaInsTech
func AmountinWords(iCompany uint, amount float64, curr string) (aiw string, csym string) {
	iBusinssdate := GetBusinessDate(iCompany, 1, 2)

	var p0023data paramTypes.P0023Data
	var extradatap0023 paramTypes.Extradata = &p0023data

	err := GetItemD(int(iCompany), "P0023", curr, iBusinssdate, &extradatap0023)
	if err != nil {
		return err.Error(), ""

	}
	// Alternate Words type between Millions and Lakhs
	switch {
	case p0023data.CurWordType == "L":
		aiw = WordsinLakhs(amount, p0023data.CurSymbol, p0023data.CurBill, p0023data.CurCoin)
		return aiw, p0023data.CurSymbol

	case p0023data.CurWordType == "M":
		aiw = WordsinMillions(amount, p0023data.CurSymbol, p0023data.CurBill, p0023data.CurCoin)
		return aiw, p0023data.CurSymbol

	default:
		return "error", ""
	}
}

// #107
// Function to convert INR amounts in words
//
// Inputs: For INR (Lakhs)
//
// # Outputs  Values in Words
// ©  FuturaInsTech
func WordsinLakhs(camt float64, csym string, cname string, ccoin string) (aiw string) {
	// Function to convert a number to its corresponding words

	// Define word representations for numbers from 0 to 19
	ones := []string{"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}

	// Define word representations for multiples of 10 from 20 to 90
	tens := []string{"", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}

	if camt == 0 {
		return "zero"
	}

	if camt < 0 {
		return "minus " + WordsinLakhs(-camt, csym, cname, ccoin)
	}

	num := int(camt)
	dec := int((camt - float64(num)) * 100)

	// Process the number in crore, lakh, thousand, and units

	crores := num / 10000000
	num %= 10000000

	lakhs := num / 100000
	num %= 100000

	thousands := num / 1000
	num %= 1000

	units := num

	words := ""
	words += cname + " "

	// Convert crores to words
	if croreWords := HundredsInWords(crores, ones, tens); len(croreWords) > 0 {
		words += croreWords + " Crore "
	}

	// Convert lakhs to words
	if lakhWords := HundredsInWords(lakhs, ones, tens); len(lakhWords) > 0 {
		words += lakhWords + " Lakh "
	}

	// Convert thousands to words
	if thousandWords := HundredsInWords(thousands, ones, tens); len(thousandWords) > 0 {
		words += thousandWords + " Thousand "
	}

	// Convert units to words
	if unitWords := HundredsInWords(units, ones, tens); len(unitWords) > 0 {
		words += unitWords
	}

	// Convert dec to decwords

	if decWords := HundredsInWords(dec, ones, tens); len(decWords) > 0 {
		words += " and "
		words += decWords
		words += " " + ccoin
	}

	return strings.TrimSpace(words)
}

// #108
// Function to convert USD,SGD,EURO,GBP... amounts in words
//
// Inputs: For USD (Millions)
//
// # Outputs  Values in Words
// ©  FuturaInsTech
func WordsinMillions(camt float64, csym string, cname string, ccoin string) (aiw string) {
	// Define word representations for numbers from 0 to 19
	ones := []string{"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}

	// Define word representations for multiples of 10 from 20 to 90
	tens := []string{"", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}

	if camt == 0 {
		return "zero"
	}

	if camt < 0 {
		return "minus " + WordsinMillions(-camt, csym, cname, ccoin)
	}

	num := int64(camt)
	dec := int64((camt - float64(num)) * 100)

	// Process the number in trillion, billion, million, thousand, and units

	trillions := num / 1000000000000
	num %= 1000000000000

	billions := num / 1000000000
	num %= 1000000000

	millions := num / 1000000
	num %= 1000000

	thousands := num / 1000
	num %= 1000

	units := num

	words := ""

	// Convert trillions to words
	trillionsInt := int(trillions) // Convert int64 to int
	if trillionWords := HundredsInWords(trillionsInt, ones, tens); len(trillionWords) > 0 {
		words += trillionWords + " Trillion "
	}

	// Convert billions to words
	// Convert billions to words
	if billionWords := HundredsInWords(int(billions), ones, tens); len(billionWords) > 0 {
		words += billionWords + " Billion "
	}

	// Convert millions to words
	if millionWords := HundredsInWords(int(millions), ones, tens); len(millionWords) > 0 {
		words += millionWords + " Million "
	}

	// Convert thousands to words
	if thousandWords := HundredsInWords(int(thousands), ones, tens); len(thousandWords) > 0 {
		words += thousandWords + " Thousand "
	}

	// Convert units to words
	if unitWords := HundredsInWords(int(units), ones, tens); len(unitWords) > 0 {
		words += unitWords
	}

	words += " " + cname

	// Convert dec to decwords
	if decWords := HundredsInWords(int(dec), ones, tens); len(decWords) > 0 {
		words += " and "
		words += decWords
		words += " " + ccoin
	}

	return strings.TrimSpace(words)
}

// #109
// Function to convert all digits of given number in words
// Inputs: For INR (Hundreds)
// # Outputs  Values in Words
// ©  FuturaInsTech
func HundredsInWords(num int, ones, tens []string) (aiw string) {
	if num == 0 {
		return ""
	}

	words := ""

	// Convert hundreds to words
	if num >= 100 {
		words += ones[num/100] + " Hundred "
		num %= 100
	}

	// Convert tens and ones to words
	if num > 0 {
		if num < 20 {
			words += ones[num]
		} else {
			words += tens[num/10]
			if onesPlace := num % 10; onesPlace > 0 {
				words += " " + ones[onesPlace]
			}
		}
	}

	return strings.TrimSpace(words)
}

// #110
// *********************************************************************************************
// This Function Take a Date in YYYYMMDD format and Freq as Input
// and return the Premium Due Dates as a String value
//
// E.g:  " " for "S", "1 of Mar Jun Sep Dec" for "Q", "1 of Every Month" for "M"
//
// ©  FuturaInsTech
// *********************************************************************************************
func GetPremDueDates(iStartDate string, freq string) string {
	var months []string
	var dueDay string
	var dueDates string
	x := 0
	dueDay = iStartDate[6:8]

	freqs := []string{"S", "Y", "H", "Q", "M"}
	intervalMonths := []int{0, 12, 6, 3, 1}

	for i := 0; i <= len(freqs); i++ {
		if freq == freqs[i] {
			x = i
			break
		}
		continue
	}

	if intervalMonths[x] <= 0 || intervalMonths[x] > 12 {
		dueDates = " "
		return dueDates
	}

	if intervalMonths[x] == 1 {
		dueDates = dueDay + " of Every Month"
		return dueDates
	}

	dueDates = dueDay + " of"
	months = append(months, dueDates)

	// Get the start year and the last day of that year
	startDate := String2Date(iStartDate)
	endDate := AddMonth(startDate, 12)

	// Loop through all the months in the given year
	intervalDate := startDate
	for intervalDate.Before(endDate) {
		months = append(months, intervalDate.Month().String()[:3])
		intervalDate = intervalDate.AddDate(0, intervalMonths[x], 0)
	}

	dueDates = strings.Join(months, " ")

	return dueDates
}

// #111
// *********************************************************************************************
// This Function is to CreateReceipt
// Input Values are Company Code, Policy , Address, Amount, Collection Date, Currency, Collection Type, Reference
// Output Values are Receipt No and Error
//
// ©  FuturaInsTech
// ***
func CreateReceiptB(iCompany uint, iPolicy uint, iAmount float64, iCollDate string, iCollCurr string, iCollType string, iRef string, iMethod string, iIFSC string, iBankAc string) (oreceipt uint, oerror error) {
	iBusinssdate := GetBusinessDate(iCompany, 1, 2)

	var policyenq models.Policy
	var receiptupd models.Receipt
	var result *gorm.DB
	var clientenq models.Client

	result = initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.Error != nil {
		return 0, errors.New(result.Error.Error())
	}

	iClient := policyenq.ClientID

	result = initializers.DB.Find(&clientenq, "company_id = ? and Id = ?", iCompany, iClient)

	if result.Error != nil {
		return 0, errors.New(result.Error.Error())
	}

	if clientenq.ClientStatus != "AC" {
		return 0, errors.New(result.Error.Error())
	}

	var p0055data paramTypes.P0055Data
	var extradatap0055 paramTypes.Extradata = &p0055data
	iKey := iCollType
	err := GetItemD(int(iCompany), "P0055", iKey, iBusinssdate, &extradatap0055)

	if err != nil {
		return 0, errors.New(err.Error())
	}

	var p0027data paramTypes.P0027Data
	var extradata paramTypes.Extradata = &p0027data

	err = GetItemD(int(iCompany), "P0027", iMethod, iBusinssdate, &extradata)

	if err != nil {
		return 0, errors.New(err.Error())
	}

	receiptupd.AccAmount = iAmount
	receiptupd.AccCurry = iCollCurr
	receiptupd.AddressID = policyenq.AddressID
	receiptupd.BankAccountNo = iBankAc
	receiptupd.BankIFSC = iIFSC
	receiptupd.InsurerBankAccNo = p0055data.BankAccount
	receiptupd.InsurerBankIFSC = p0055data.BankCode
	receiptupd.CurrentDate = iBusinssdate
	receiptupd.DateOfCollection = iCollDate
	receiptupd.BankReferenceNo = iRef
	receiptupd.Branch = "HO"
	receiptupd.ClientID = policyenq.ClientID
	receiptupd.ReceiptRefNo = iPolicy
	receiptupd.ReceiptAmount = policyenq.InstalmentPrem
	receiptupd.ReceiptDueDate = policyenq.PaidToDate
	receiptupd.Tranno = policyenq.Tranno
	receiptupd.TypeOfReceipt = iCollType
	receiptupd.CompanyID = iCompany

	// Save Receipt
	result = initializers.DB.Create(&receiptupd)

	// Debit Entry
	glcode := p0027data.GlMovements[0].AccountCode
	var acccode models.AccountCode
	result = initializers.DB.First(&acccode, "company_id = ? and account_code = ? ", iCompany, glcode)
	if result.RowsAffected == 0 {
		return 0, errors.New(err.Error())
	}
	var iSequenceno uint64
	iSequenceno++
	iAccountCodeID := acccode.ID
	iAccAmount := receiptupd.AccAmount
	iAccountCode := glcode + receiptupd.Branch + p0055data.GlAccount
	iEffectiveDate := receiptupd.DateOfCollection
	iGlAmount := receiptupd.AccAmount

	iGlRdocno := receiptupd.ID
	var iGlRldgAcct string
	//iGlRldgAcct := strconv.Itoa(int(iClient))
	// As per our discussion on 22/06/2023, it is decided to use policy no in RLDGACCT
	iGlRldgAcct = strconv.Itoa(int(iPolicy))
	iGlSign := p0027data.GlMovements[0].GlSign
	iTranno := 0

	err = PostGlMove(iCompany, iCollCurr, iEffectiveDate, int(iTranno), iGlAmount,
		iAccAmount, iAccountCodeID, uint(iGlRdocno), string(iGlRldgAcct), iSequenceno, iGlSign, iAccountCode, iMethod, "", "")

	if err != nil {
		return 0, errors.New(err.Error())
	}

	// Credit Entry
	glcode = p0027data.GlMovements[1].AccountCode
	var acccode1 models.AccountCode
	result = initializers.DB.First(&acccode1, "company_id = ? and account_code = ? ", iCompany, glcode)
	if result.RowsAffected == 0 {
		return 0, errors.New(err.Error())
	}

	iSequenceno++
	iAccountCodeID = acccode.ID
	iAccAmount = receiptupd.AccAmount
	iAccountCode = glcode
	iEffectiveDate = receiptupd.DateOfCollection
	iGlAmount = receiptupd.AccAmount
	iGlRdocno = iPolicy
	iGlRldgAcct = strconv.Itoa(int(iPolicy))
	iGlSign = p0027data.GlMovements[1].GlSign
	iTranno = 0

	err = PostGlMove(iCompany, iCollCurr, iEffectiveDate, int(iTranno), iGlAmount,
		iAccAmount, iAccountCodeID, uint(iGlRdocno), iGlRldgAcct, iSequenceno, iGlSign, iAccountCode, iMethod, "", "")

	if err != nil {
		return 0, errors.New(err.Error())

	}
	if policyenq.PolStatus == "IF" {
		iNextDueDate := Date2String(GetNextDue(policyenq.PaidToDate, policyenq.PFreq, ""))
		gstamountneeded := GetTotalGSTAmount(iCompany, iPolicy, policyenq.PaidToDate, iNextDueDate)
		iPolicyDeposit := GetGlBal(iCompany, iPolicy, "PolicyDeposit")
		iStampDuty := CalculateStampDutyByPolicy(iCompany, iPolicy)
		iPayable := policyenq.InstalmentPrem + gstamountneeded + iStampDuty + iPolicyDeposit
		if iPayable <= 0 {
			TDFCollD(iCompany, iPolicy, "COLLD", 0, policyenq.PaidToDate)
			TdfhUpdate(iCompany, iPolicy)
		}
	}

	iAgency := policyenq.AgencyID

	err = CreateCommunications(iCompany, iMethod, uint(iTranno), iBusinssdate, iPolicy, receiptupd.ClientID, receiptupd.AddressID, receiptupd.ID, 0, iAgency, "", "", "", "", "", 0, 0, 0)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	return receiptupd.ID, nil
}

// # 112
// CalBonus - Calculate Bonus due on Anniversary Date
//
// NOTE: THIS IS CLONED FROM GetBonusByYear not to impact to other functions. //
//
// # Input:  Company, Coverage , Bonus Method, Coverage Start Date, Anniversary Date, Policy Status, SA, Term, Premium Term
// # Output: Calculted Bonus Amount as float64
//
// # Date in YYYYMMDD as a string
//
// ©  FuturaInsTech
func CalcBonus(iCompany uint, iCoverage string, iBonusMethod string, iDate string, iAnnivDate string, iStatus string, iSA uint, iTerm uint, iPTerm uint) float64 {
	//	fmt.Println("inside Bonus ", iCoverage, iCompany, iBonusMethod, iDate, iYear, iStatus, iSA, iTerm)

	var iKey string
	var oBonus float64

	if iBonusMethod == "" {
		// No Bonus Method exists hence return bonus as zero and exit
		oBonus = 0
		return oBonus
	}

	var q0014data paramTypes.Q0014Data
	var extradata1 paramTypes.Extradata = &q0014data

	iKey = iBonusMethod + iStatus + strconv.Itoa(int(iTerm)) + strconv.Itoa(int(iPTerm))
	err := GetItemD(int(iCompany), "Q0014", iKey, iAnnivDate, &extradata1)
	if err != nil {
		oBonus = 0
		return oBonus
	}

	iYear, _, _, _, _, _, _, _ := NoOfDays(iAnnivDate, iDate)

	for i := 0; i < len(q0014data.BRates); i++ {
		if iYear <= int64(q0014data.BRates[i].Term) {
			oBonus = float64(iSA) * (q0014data.BRates[i].Percentage) / 100
			return oBonus
		}
	}
	return 0
}

// # 113
// CalcInterimBonus - Calculate Interim Bonus using Interim Bonus Rates
//
// # Input:  Company, Coverage, BonusMethod, StartDate, EffectiveDate, Bonus Date, Coverage Status, Sum Assured, Term, Premium Term
// # Output: Calculted Interim Bonus Amount as float64
//
// # Date in YYYYMMDD as a string
//
// ©  FuturaInsTech
func CalcIBonus(iCompany uint, iCoverage string, iBonusMethod string, iDate string, iEffectiveDate string, iBonusDate string, iStatus string, iSA uint, iTerm uint, iPTerm uint) float64 {
	//	fmt.Println("inside Bonus ", iCompany, iCoverage, iBonusMethod, iDate, iEffectiveDate, iBonusDate, iStatus, iSA, iTerm, iPterm)

	if iBonusDate == "" {
		iBonusDate = iDate
	}

	var iKey string
	var oBonus float64 = 0
	var rateYear int64 = 0
	var prorateDays int64 = 0

	if iBonusMethod == "" {
		// No Interim Bonus Method exists hence return bonus as zero and exit
		oBonus = 0
		return oBonus
	}

	var q0014data paramTypes.Q0014Data
	var extradata1 paramTypes.Extradata = &q0014data

	iKey = iBonusMethod + iStatus + strconv.Itoa(int(iTerm)) + strconv.Itoa(int(iPTerm))
	err := GetItemD(int(iCompany), "Q0014", iKey, iEffectiveDate, &extradata1)
	if err != nil {
		oBonus = 0
		return oBonus
	}

	iYear, _, iDays, _, _, _, _, _ := NoOfDays(iEffectiveDate, iDate)
	iyearsindays := iYear * 365
	if iDays > iyearsindays {
		rateYear = iYear + 1
	}

	iYear, _, iDays, _, _, _, _, _ = NoOfDays(iEffectiveDate, iBonusDate)
	iyearsindays = iYear * 365
	if iDays > 0 {
		prorateDays = iDays
	}

	for i := 0; i < len(q0014data.BRates); i++ {
		if rateYear <= int64(q0014data.BRates[i].Term) {
			oBonus = float64(iSA) * (q0014data.BRates[i].Percentage) / 100 * float64(prorateDays/365)
			return oBonus
		}
	}
	return 0
}

// # 114
// TDFIBD - Time Driven Function - Income Benefit Date Updation
//
// Inputs: Company, Policy, Function IBEN, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFIBD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var incomeb models.IBenefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	results := initializers.DB.First(&incomeb, "company_id = ? and policy_id = ? and paid_date = ?", iCompany, iPolicy, "")

	if results.Error != nil {
		return "", results.Error
	}
	result := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ? ", iCompany, iPolicy, iFunction)

	if result.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = incomeb.NextPayDate
		tdfpolicy.Tranno = iTranno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	} else {
		initializers.DB.Delete(&tdfpolicy)
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = incomeb.NextPayDate
		tdfpolicy.Tranno = iTranno
		initializers.DB.Create(&tdfpolicy)
		return "", nil
	}
}
func TDFIBDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var incomeb models.IBenefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		// txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&incomeb, "company_id = ? and policy_id = ? and paid_date = ?", iCompany, iPolicy, "")
	if result.Error != nil {
		// txn.Rollback()
		return "", result.Error
	}

	result = initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ? ", iCompany, iPolicy, iFunction)

	if result.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = incomeb.NextPayDate
		tdfpolicy.Tranno = iTranno
		result = initializers.DB.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

	} else {
		result = initializers.DB.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = incomeb.NextPayDate
		tdfpolicy.Tranno = iTranno
		result = initializers.DB.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}
		return "", nil
	}
	return "", nil
}

// # 115
// CalRBonus - Calculate Bonus due on Annniversary Date OR Bonus Date
//
// # Input:  Company, Coverage , Bonus Method, Coverage Start Date, Anniversary Date OR Bonus Date, Policy Status, Coverage SA, Coverage Term, Coverage Premium Term
// # Output: Calculted Bonus Amount as float64
//
// # Date in YYYYMMDD as a string
//
// ©  FuturaInsTech
func CalcRBonus(iCompany uint, iCoverage string, iBonusMethod string, iDate string, iBonusDate string, iEffectiveDate string, iStatus string, iSA uint, iTerm uint, iPTerm uint) (oBonus float64) {

	var iKey string
	var rateYear int64 = 0

	var q0014data paramTypes.Q0014Data
	var extradataq0014 paramTypes.Extradata = &q0014data
	iKey = iBonusMethod + iStatus + strconv.Itoa(int(iTerm)) + strconv.Itoa(int(iPTerm))
	err := GetItemD(int(iCompany), "Q0014", iKey, iEffectiveDate, &extradataq0014)
	if err != nil {
		oBonus := 0.0
		return oBonus
	}

	iYear, iMonth, iDay, _, _, _ := StringDateDiff(iEffectiveDate, iDate, "")
	rateYear = int64(iYear)

	if iMonth > 0 || iDay > 0 {
		rateYear++
	}

	if rateYear == 0 {
		oBonus := 0.0
		return oBonus
	}

	for i := 0; i < len(q0014data.BRates); i++ {
		if rateYear <= int64(q0014data.BRates[i].Term) {
			oBonus := float64(iSA) * (q0014data.BRates[i].Percentage) / 100
			return oBonus
		}
	}
	return 0
}

// # 116
// StringDateDiff - Calculate Bonus due on Annniversary Date OR Bonus Date
//
// # Input:
// # Output:
//
// #
//
// ©  FuturaInsTech
func StringDateDiff(as, bs string, m string) (year, month, day, hour, min, sec int) {
	// method = N means age nearer birthday
	var a time.Time
	var b time.Time
	a = String2Date(as)
	b = String2Date(bs)

	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// # 118
// TDFExtrD - Time Driven Function - Expiry Date Updation
//
// Inputs: Company, Policy, Functio EXTRD, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFExtrD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var extraenq []models.Extra
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var policyenq models.Policy
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.Find(&extraenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return "", result.Error
		}
	}
	oDate := ""
	for i := 0; i < len(extraenq); i++ {
		if oDate == "" {
			oDate = extraenq[i].ToDate
		}
		if extraenq[i].ToDate < oDate {
			oDate = extraenq[i].ToDate
		}
	}
	// Subtract Billing Lead Days as well
	result = initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	err := GetItemD(int(iCompany), "Q0005", policyenq.PProduct, policyenq.PRCD, &extradataq0005)
	if err != nil {
		return "", err
	}
	if oDate != "" {
		oDate = AddLeadDays(oDate, (-1 * q0005data.BillingLeadDays))
	}

	if oDate != "" {
		results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			initializers.DB.Create(&tdfpolicy)
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			initializers.DB.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

func TDFExtrDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var extraenq []models.Extra
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var policyenq models.Policy
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.Find(&extraenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return "", result.Error
		}
	}
	oDate := ""
	for i := 0; i < len(extraenq); i++ {
		if oDate == "" {
			oDate = extraenq[i].ToDate
		}
		if extraenq[i].ToDate < oDate {
			oDate = extraenq[i].ToDate
		}
	}
	// Subtract Billing Lead Days as well
	result = initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	err := GetItemD(int(iCompany), "Q0005", policyenq.PProduct, policyenq.PRCD, &extradataq0005)
	if err != nil {
		return "", err
	}
	if oDate != "" {
		oDate = AddLeadDays(oDate, (-1 * q0005data.BillingLeadDays))
	}

	if oDate != "" {
		results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				return "", result.Error
			}
			return "", nil
		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				return "", result.Error
			}
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			result = txn.Create(&tdfpolicy)

			if result.Error != nil {
				return "", result.Error
			}
			return "", nil
		}
	}
	return "", nil
}

// # 119
// GetErrorDesc - Get Error Description
//
// Inputs: Company, Language, Short Code
//
// # Long Description, Error
//
// ©  FuturaInsTech
func GetErrorDesc(iCompany uint, iLanguage uint, iShortCode string) (string, error) {
	var errorenq models.Error

	result := initializers.DB.Find(&errorenq, "company_id = ? and language_id = ? and short_code = ?", iCompany, iLanguage, iShortCode)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", errors.New(" -" + strconv.FormatUint(uint64(iCompany), 10) + "-" + "-" + strconv.FormatUint(uint64(iLanguage), 10) + "-" + " is missing")
	}

	return errorenq.LongCode, nil
}

// # 120
// PostAllocation - This function apportion amount into different funds and investible and non investible
//
// Inputs:
//
// # Success/Failure
//
// ©  FuturaInsTech
func PostAllocation(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iFrequency string, iStartDate string, iEffDate string, iGender string, iAllocMethod string, iTranno uint) error {

	var policyenq models.Policy

	result := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data
	iDate := iStartDate
	iKey := iAllocMethod + iGender
	err := GetItemD(int(iCompany), "P0060", iKey, iDate, &extradatap0060)
	if err != nil {
		return errors.New(err.Error())
	}
	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey = iHistoryCode + iBenefitCode
	err = GetItemD(int(iCompany), "P0059", iKey, iDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}
	iAllocationDate := ""
	if iEffDate == iStartDate {
		a := GetNextDue(iStartDate, iFrequency, "")
		iAllocationDate = Date2String(a)
	}
	iNoofMonths := NewNoOfInstalments(iStartDate, iAllocationDate)
	iAllocPercentage := 0.00
	for i := 0; i < len(p0060data.AlBand); i++ {
		if uint(iNoofMonths) <= p0060data.AlBand[i].Months {
			iAllocPercentage = p0060data.AlBand[i].Percentage
			break
		}
	}
	iInvested := RoundFloat(iAmount*(iAllocPercentage/100), 2)
	iNonInvested := RoundFloat(iAmount*((100-iAllocPercentage)/100), 2)

	var ilpfundenq []models.IlpFund

	result = initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	for j := 0; j < len(ilpfundenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}

		var ilptrancrt models.IlpTransaction
		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpfundenq[j].FundCode
		ilptrancrt.FundType = ilpfundenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundEffDate = iBusinessDate
		ilptrancrt.FundAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.FundCurr = ilpfundenq[j].FundCurr
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = ilpfundenq[j].FundPercentage
		ilptrancrt.AccountCode = "Invested" // ranga

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "P"
		result = initializers.DB.Create(&ilptrancrt)
	}
	// Non Invested Amount Updation

	var ilptrancrt models.IlpTransaction
	// Move Variables
	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = "NONIN"
	ilptrancrt.FundType = "NI"
	ilptrancrt.TransactionDate = iEffDate
	ilptrancrt.FundEffDate = iEffDate
	ilptrancrt.FundAmount = iNonInvested
	ilptrancrt.FundCurr = ""
	ilptrancrt.FundUnits = 0
	ilptrancrt.FundPrice = 0
	ilptrancrt.CurrentOrFuture = "C"
	ilptrancrt.OriginalAmount = iNonInvested
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "NI"
	ilptrancrt.AllocationCategory = "NI"
	ilptrancrt.InvNonInvPercentage = 0
	ilptrancrt.Tranno = iTranno

	ilptrancrt.AccountCode = "NonInvested"

	ilptrancrt.CurrencyRate = 1.00 // ranga
	ilptrancrt.MortalityIndicator = ""
	ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = initializers.DB.Create(&ilptrancrt)
	return nil
}

func PostAllocationN(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iFrequency string, iStartDate string, iEffDate string, iGender string, iAllocMethod string, iTranno uint, txn *gorm.DB) error {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}

	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data
	iDate := iStartDate
	iKey := iAllocMethod + iGender
	err := GetItemD(int(iCompany), "P0060", iKey, iDate, &extradatap0060)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}
	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey = iHistoryCode + iBenefitCode
	err = GetItemD(int(iCompany), "P0059", iKey, iDate, &extradatap0059)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}
	iAllocationDate := ""
	if iEffDate == iStartDate {
		a := GetNextDue(iStartDate, iFrequency, "")
		iAllocationDate = Date2String(a)
	}
	iNoofMonths := NewNoOfInstalments(iStartDate, iAllocationDate)
	iAllocPercentage := 0.00
	for i := 0; i < len(p0060data.AlBand); i++ {
		if uint(iNoofMonths) <= p0060data.AlBand[i].Months {
			iAllocPercentage = p0060data.AlBand[i].Percentage
			break
		}
	}
	iInvested := RoundFloat(iAmount*(iAllocPercentage/100), 2)
	iNonInvested := RoundFloat(iAmount*((100-iAllocPercentage)/100), 2)

	var ilpfundenq []models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	for j := 0; j < len(ilpfundenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}

		var ilptrancrt models.IlpTransaction
		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpfundenq[j].FundCode
		ilptrancrt.FundType = ilpfundenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundEffDate = iBusinessDate
		ilptrancrt.FundAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.FundCurr = ilpfundenq[j].FundCurr
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = ilpfundenq[j].FundPercentage
		ilptrancrt.AccountCode = "Invested" // ranga

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "P"
		result = txn.Create(&ilptrancrt)
		if result.Error != nil {
			txn.Rollback()
			return result.Error
		}
	}
	// Non Invested Amount Updation

	var ilptrancrt models.IlpTransaction
	// Move Variables
	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = "NONIN"
	ilptrancrt.FundType = "NI"
	ilptrancrt.TransactionDate = iEffDate
	ilptrancrt.FundEffDate = iEffDate
	ilptrancrt.FundAmount = iNonInvested
	ilptrancrt.FundCurr = ""
	ilptrancrt.FundUnits = 0
	ilptrancrt.FundPrice = 0
	ilptrancrt.CurrentOrFuture = "C"
	ilptrancrt.OriginalAmount = iNonInvested
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "NI"
	ilptrancrt.AllocationCategory = "NI"
	ilptrancrt.InvNonInvPercentage = 0
	ilptrancrt.Tranno = iTranno

	ilptrancrt.AccountCode = "NonInvested"

	ilptrancrt.CurrencyRate = 1.00
	ilptrancrt.MortalityIndicator = ""
	ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = txn.Create(&ilptrancrt)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}
	return nil
}

// # 121
// TDFFUNDP - Time Driven Function - Update Fund Price
//
// Inputs: Company, Policy, Function FUNDP, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFFundP(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var ilptransenq []models.IlpTransaction
	odate := "00000000"

	result := initializers.DB.Where("company_id = ? and policy_id = ? and ul_process_flag = ?", iCompany, iPolicy, "P").Order("fund_eff_date").Find(&ilptransenq)
	for i := 0; i < len(ilptransenq); i++ {
		if ilptransenq[i].FundEffDate > odate {
			odate = ilptransenq[i].FundEffDate
		}
	}

	result = initializers.DB.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", result.Error
	}

	result = initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		return "", result.Error
	}

	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if odate != "00000000" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno
			tdfpolicy.Seqno = tdfrule.Seqno
			initializers.DB.Create(&tdfpolicy)
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno

			initializers.DB.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

func TDFFundPN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string, txn *gorm.DB) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var ilptransenq []models.IlpTransaction
	odate := "00000000"

	result := txn.Where("company_id = ? and policy_id = ? and ul_process_flag = ?", iCompany, iPolicy, "P").Order("fund_eff_date").Find(&ilptransenq)
	if result.Error != nil {
		return "", result.Error
	}

	for i := 0; i < len(ilptransenq); i++ {
		if ilptransenq[i].FundEffDate > odate {
			odate = ilptransenq[i].FundEffDate
		}
	}

	result = txn.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", result.Error
	}

	result = txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		return "", result.Error
	}

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if odate != "00000000" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno
			tdfpolicy.Seqno = tdfrule.Seqno
			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				return "", result.Error
			}
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno

			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				return "", result.Error
			}
			return "", nil
		}
	}
	return "", nil
}

// # 122
func GetAllFundValueByPol(iCompany uint, iPolicy uint, iDate string) (float64, float64, string) {
	if iDate == "" {
		iDate = GetBusinessDate(iCompany, 0, 0)
	}

	var ilpsummaryenq []models.IlpSummary
	result := initializers.DB.Find(&ilpsummaryenq, "company_id = ?  and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		return 0, 0, iDate
	}

	bpfv := 0.0
	opfv := 0.0
	for i := 0; i < len(ilpsummaryenq); i++ {
		bpv, opv, _ := GetaFundValue(iCompany, iPolicy, ilpsummaryenq[i].FundCode, iDate)
		bpfv = RoundFloat(bpfv+bpv, 2)
		opfv = RoundFloat(opfv+opv, 2)
	}
	return bpfv, opfv, iDate
}

func GetAllFundValueByPolN(iCompany uint, iPolicy uint, iDate string, txn *gorm.DB) (float64, float64, string) {
	if iDate == "" {
		iDate = GetBusinessDate(iCompany, 0, 0)
	}

	var ilpsummaryenq []models.IlpSummary
	result := txn.Find(&ilpsummaryenq, "company_id = ?  and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		return 0, 0, iDate
	}

	bpfv := 0.0
	opfv := 0.0
	for i := 0; i < len(ilpsummaryenq); i++ {
		bpv, opv, _ := GetaFundValueN(iCompany, iPolicy, ilpsummaryenq[i].FundCode, iDate, txn)
		bpfv = RoundFloat(bpfv+bpv, 5)
		opfv = RoundFloat(opfv+opv, 5)
	}
	return bpfv, opfv, iDate
}

// # 123
// GetaFundValue - Get Fund Value per Policy (Across All Benefits)
//
// Inputs: Company, Policy, Fund Code, Fund Price Date
//
// # Outputs  Bid Amount, Offer Amount, Effective Date
//
// ©  FuturaInsTech
func GetaFundValue(iCompany uint, iPolicy uint, iFundCode string, iDate string) (float64, float64, string) {
	if iCompany == 0 || iPolicy == 0 || iFundCode == "" || iDate == "" {
		return 0, 0, iDate
	}

	bpfundvalue := 0.0
	opfundvalue := 0.0
	var ilpsummaryenq []models.IlpSummary
	result := initializers.DB.Order("fund_code").
		Find(&ilpsummaryenq, "company_id = ?  and policy_id = ? and fund_code = ?", iCompany, iPolicy, iFundCode)
	if result.Error != nil {
		return 0, 0, iDate
	}

	var ilppriceenq models.IlpPrice
	var iPriceDateUsed = "00000000"
	result = initializers.DB.Where("company_id = ? and fund_code = ? and approval_flag = ? and fund_eff_date <= ?", iCompany, iFundCode, "AP", iDate).Order("fund_eff_date DESC").First(&ilppriceenq)
	if result.Error != nil {
		return 0, 0, iDate
	}

	iPriceDateUsed = ilppriceenq.FundEffDate
	fmt.Println("******* Price Date Used|BidPrice|OfferPrice ********", iPriceDateUsed, ilppriceenq.FundBidPrice, ilppriceenq.FundOfferPrice)
	for i := 0; i < len(ilpsummaryenq); i++ {
		bpfundvalue = RoundFloat(ilpsummaryenq[i].FundUnits*ilppriceenq.FundBidPrice, 2)
		opfundvalue = RoundFloat(ilpsummaryenq[i].FundUnits*ilppriceenq.FundOfferPrice, 2)
	}
	return bpfundvalue, opfundvalue, iPriceDateUsed

}

func GetaFundValueN(iCompany uint, iPolicy uint, iFundCode string, iDate string, txn *gorm.DB) (float64, float64, string) {
	if iCompany == 0 || iPolicy == 0 || iFundCode == "" || iDate == "" {
		return 0, 0, iDate
	}

	bpfundvalue := 0.0
	opfundvalue := 0.0
	var ilpsummaryenq []models.IlpSummary
	result := txn.Order("fund_code").
		Find(&ilpsummaryenq, "company_id = ?  and policy_id = ? and fund_code = ?", iCompany, iPolicy, iFundCode)
	if result.Error != nil {
		return 0, 0, iDate
	}

	var ilppriceenq models.IlpPrice
	var iPriceDateUsed = "00000000"
	result = txn.Where("company_id = ? and fund_code = ? and approval_flag = ? and fund_eff_date <= ?", iCompany, iFundCode, "AP", iDate).Order("fund_eff_date DESC").First(&ilppriceenq)
	if result.Error != nil {
		return 0, 0, iDate
	}

	iPriceDateUsed = ilppriceenq.FundEffDate
	fmt.Println("******* Price Date Used|BidPrice|OfferPrice ********", iPriceDateUsed, ilppriceenq.FundBidPrice, ilppriceenq.FundOfferPrice)
	for i := 0; i < len(ilpsummaryenq); i++ {
		bpfundvalue = RoundFloat(ilpsummaryenq[i].FundUnits*ilppriceenq.FundBidPrice, 5)
		opfundvalue = RoundFloat(ilpsummaryenq[i].FundUnits*ilppriceenq.FundOfferPrice, 5)
	}
	return bpfundvalue, opfundvalue, iPriceDateUsed

}

// # 124
// GetAllFundValueByBenefit - Get All Fund Values for a Given Benefit ID
//
// Inputs: Company, Policy, Benefit ID, Fund Code, Fund Price Date
//
// # Outputs  Bid Amount, Offer Amount, Effective Date
//
// ©  FuturaInsTech
func GetAllFundValueByBenefit(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string, iDate string) (float64, float64, string) {
	if iDate == "" {
		iDate = GetBusinessDate(iCompany, 0, 0)
	}

	var ilpsummaryenq []models.IlpSummary
	if iFundCode != "" {
		result := initializers.DB.Find(&ilpsummaryenq, "company_id = ?  and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
		if result.Error != nil {
			return 0, 0, iDate
		}
	} else {
		result := initializers.DB.Find(&ilpsummaryenq, "company_id = ?  and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
		if result.Error != nil {
			return 0, 0, iDate
		}
	}

	bpfv := 0.0
	opfv := 0.0
	for i := 0; i < len(ilpsummaryenq); i++ {
		iFundCode := ilpsummaryenq[i].FundCode
		bpv, opv, _ := GetaFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iDate)
		bpfv = RoundFloat(bpfv+bpv, 2)
		opfv = RoundFloat(opfv+opv, 2)
	}
	return bpfv, opfv, iDate
}

// # 125
// GetFundValueByBeneift - Get Fund Summary based on a Benefit ID
//
// Inputs: Company, Policy, Benefit ID, Fund Code, Fund Price Date
//
// # Outputs  Bid Amount, Offer Amount, Effective Date
//
// ©  FuturaInsTech
func GetaFundValueByBenefit(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string, iDate string) (float64, float64, string) {
	if iCompany == 0 || iPolicy == 0 || iFundCode == "" || iDate == "" {
		return 0, 0, iDate
	}

	bpfundvalue := 0.0
	opfundvalue := 0.0
	var ilpsummaryenq models.IlpSummary
	result := initializers.DB.Find(&ilpsummaryenq, "company_id = ?  and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		return 0, 0, iDate
	}

	var ilppriceenq models.IlpPrice
	var iPriceDateUsed = "00000000"
	result = initializers.DB.Where("company_id = ? and fund_code = ? and  approval_flag = ? and fund_eff_date <= ?", iCompany, iFundCode, "AP", iDate).Order("fund_eff_date DESC").First(&ilppriceenq)
	if result.Error != nil {
		return 0, 0, iDate
	}

	iPriceDateUsed = ilppriceenq.FundEffDate
	fmt.Println("******* Price Date Used|BidPrice|OfferPrice ********", iPriceDateUsed, ilppriceenq.FundBidPrice, ilppriceenq.FundOfferPrice)

	bpfundvalue = RoundFloat(ilpsummaryenq.FundUnits*ilppriceenq.FundBidPrice, 2)
	opfundvalue = RoundFloat(ilpsummaryenq.FundUnits*ilppriceenq.FundOfferPrice, 2)

	return bpfundvalue, opfundvalue, iPriceDateUsed

}

// #126
//
// # GetILPAlloc - Get ILP Allocation Based on From and To Date
//
// Inputs: Company, Policy, Frequency, Start Date of Policy, Premium Adjusted Date, Coverage Code, Allocation Method, Gender, Amount Collected
//
// # Outputs  Allocated , Non Allocated and Error
//
// ©  FuturaInsTech
func GetILPAlloc(iCompany uint, iPolicy uint, iFreq string, iStartDate string, iEndDate string, iCoverage string, iAllocMethod string, iGender string, iAmount float64) (float64, float64, error) {
	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data
	iDate := iStartDate
	iKey := iAllocMethod + iGender
	err := GetItemD(int(iCompany), "P0060", iKey, iDate, &extradatap0060)
	if err != nil {
		return 0, 0, errors.New(err.Error())
	}

	if iStartDate == iEndDate {
		a := GetNextDue(iStartDate, iFreq, "")
		iEndDate = Date2String(a)
	}
	iNoofMonths := NewNoOfInstalments(iStartDate, iEndDate)
	iAllocPercentage := 0.00
	for i := 0; i < len(p0060data.AlBand); i++ {
		if uint(iNoofMonths) <= p0060data.AlBand[i].Months {
			iAllocPercentage = p0060data.AlBand[i].Percentage
			break
		}
	}
	iInvested := RoundFloat(iAmount*(iAllocPercentage/100), 2)
	iNonInvested := RoundFloat(iAmount*((100-iAllocPercentage)/100), 2)
	return iInvested, iNonInvested, nil
}

// # 127
// PostBuySell - Buy or Sell or Non Invested Posting
//
// Inputs: Function (Buy/Sell/NonInvested), Company, Policy, Contract Currency, Transaciton Code, Coverage Code, Coverage Start Date, Adjusted Date, Coverage ID, Amount (Non Invested/Invested/Mortality/Fee/Surrender), Tranno
//
// # Outputs  Allocated , Non Allocated and Error
//
// ©  FuturaInsTech
func PostBuySell(iFunction string, iCompany uint, iPolicy uint, iContractCurr string, iHistoryCode string, iCoverage string, iStartDate string, iEffDate string, iBenefit uint, iAmount float64, iTranno uint) error {
	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iCoverage
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}
	iCurrentOrFuture := p0059data.CurrentOrFuture

	var ilpfundenq []models.IlpFund

	result := initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	if iFunction == "NonInvested" {
		var ilptrancrt models.IlpTransaction
		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = "NonInvested"
		ilptrancrt.FundType = "NI"
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundEffDate = iEffDate
		ilptrancrt.FundAmount = RoundFloat(iAmount, 2)
		ilptrancrt.FundCurr = ""
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = "C"
		ilptrancrt.OriginalAmount = RoundFloat(iAmount, 2)
		ilptrancrt.ContractCurry = iContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "NI"
		ilptrancrt.InvNonInvPercentage = 0.00
		ilptrancrt.AccountCode = "NonInvested"

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "C"
		result = initializers.DB.Create(&ilptrancrt)

	} else {
		// Invested Posting
		for j := 0; j < len(ilpfundenq); j++ {
			iBusinessDate := GetBusinessDate(iCompany, 0, 0)
			if iCurrentOrFuture == "F" {
				iBusinessDate = AddLeadDays(iBusinessDate, 1)
			} else if iCurrentOrFuture == "E" {
				iBusinessDate = iEffDate
			}

			var ilptrancrt models.IlpTransaction
			ilptrancrt.CompanyID = iCompany
			ilptrancrt.PolicyID = iPolicy
			ilptrancrt.BenefitID = iBenefit
			ilptrancrt.FundCode = ilpfundenq[j].FundCode
			ilptrancrt.FundType = ilpfundenq[j].FundType
			ilptrancrt.TransactionDate = iEffDate
			ilptrancrt.FundEffDate = iBusinessDate
			ilptrancrt.FundAmount = RoundFloat(((iAmount * ilpfundenq[j].FundPercentage) / 100), 2)
			ilptrancrt.FundCurr = ilpfundenq[j].FundCurr
			ilptrancrt.FundUnits = 0
			ilptrancrt.FundPrice = 0
			ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
			ilptrancrt.OriginalAmount = RoundFloat(((iAmount * ilpfundenq[j].FundPercentage) / 100), 2)
			ilptrancrt.ContractCurry = iContractCurr
			ilptrancrt.HistoryCode = iHistoryCode
			ilptrancrt.InvNonInvFlag = ilpfundenq[j].FundType
			ilptrancrt.AllocationCategory = p0059data.AllocationCategory

			ilptrancrt.InvNonInvPercentage = ilpfundenq[j].FundPercentage
			ilptrancrt.AccountCode = iFunction + "Invested"

			ilptrancrt.CurrencyRate = 1.00 // ranga
			ilptrancrt.MortalityIndicator = ""
			ilptrancrt.SurrenderPercentage = 0
			ilptrancrt.Tranno = iTranno
			ilptrancrt.Seqno = uint(p0059data.SeqNo)
			ilptrancrt.UlProcessFlag = "P"
			result = initializers.DB.Create(&ilptrancrt)
		}
	}

	return nil
}

// # 128
// TDFFUNDM - Time Driven Function - Mortality Premium
//
// Inputs: Company, Policy, Function FUNDM, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFFundM(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {

	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"

	result := initializers.DB.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", result.Error
	}

	result = initializers.DB.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	for i := 0; i < len(benefitenq); i++ {
		if benefitenq[i].IlpMortalityDate > odate {
			odate = benefitenq[i].IlpMortalityDate
		}
	}

	result = initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		return "", result.Error
	}

	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if odate != "00000000" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno
			tdfpolicy.Seqno = tdfrule.Seqno
			initializers.DB.Create(&tdfpolicy)
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno

			initializers.DB.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}
func TDFFundMN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {

	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"

	result := txn.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	for i := 0; i < len(benefitenq); i++ {

		iCoverage := benefitenq[i].BCoverage
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data
		err := GetItemD(int(iCompany), "Q0006", iCoverage, benefitenq[i].BStartDate, &extradataq0006)
		if err != nil {
			err := errors.New("Q0006 Not Found")
			return "", err
		}
		if q0006data.PremCalcType == "U" {

			if benefitenq[i].IlpMortalityDate > odate {
				odate = benefitenq[i].IlpMortalityDate
			}
		}
	}

	result = initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	if odate != "00000000" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno
			tdfpolicy.Seqno = tdfrule.Seqno
			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}

		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno

			result = initializers.DB.Create(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			return "", nil
		}
	}
	return "", nil
}

// # 128
// TDFFUNDF - Time Driven Function - ILP Fee
//
// Inputs: Company, Policy, Function FUNDF, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFFundF(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {

	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"

	result := initializers.DB.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", result.Error
	}

	result = initializers.DB.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	for i := 0; i < len(benefitenq); i++ {
		if benefitenq[i].IlpFeeDate > odate {
			odate = benefitenq[i].IlpFeeDate
		}
	}

	result = initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		return "", result.Error
	}

	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if odate != "00000000" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno
			tdfpolicy.Seqno = tdfrule.Seqno
			initializers.DB.Create(&tdfpolicy)
			return "", nil
		} else {
			initializers.DB.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno

			initializers.DB.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

func TDFFundFN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {

	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"

	result := txn.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	for i := 0; i < len(benefitenq); i++ {

		iCoverage := benefitenq[i].BCoverage
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data
		err := GetItemD(int(iCompany), "Q0006", iCoverage, benefitenq[i].BStartDate, &extradataq0006)
		if err != nil {
			err := errors.New("Q0006 Not Found")
			return "", err
		}
		if q0006data.PremCalcType == "U" {
			if benefitenq[i].IlpFeeDate > odate {
				odate = benefitenq[i].IlpFeeDate
			}
		}
	}

	result = txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)

	if odate != "00000000" {
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno
			tdfpolicy.Seqno = tdfrule.Seqno
			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			return "", nil
		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = odate
			tdfpolicy.Tranno = iTranno

			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				txn.Rollback()
				return "", result.Error
			}
			return "", nil
		}
	}
	return "", nil
}

// # 129
// CalcMortPrem - Calculate Mortality Premium for ILP
//
// Inputs: Company, Policy, Benefit ID, History Code, Effective Date
//
// # Outputs  Mortality Premium for the Frequency, Next Due Date and Error
//
// ©  FuturaInsTech
func CalcMortPrem(iCompany uint, iPolicy uint, iBenefit uint, iHistoryCode string, iEffDate string) (float64, string, error) {

	var policyenq models.Policy

	results := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	if results.Error != nil {
		return 0, "", results.Error
	}

	iDate := policyenq.PRCD
	boolstat, _ := CheckStatus(iCompany, iHistoryCode, iDate, policyenq.PolStatus)

	if boolstat {
		err := errors.New("Invalid Policy Status")
		return 0, "", err
	}

	var benefitupd models.Benefit

	results = initializers.DB.Find(&benefitupd, "company_id = ? and policy_id = ? and id = ?", iCompany, iPolicy, iBenefit)

	if results.Error != nil {
		return 0, "", results.Error
	}

	iDate = benefitupd.BStartDate
	boolstat, _ = CheckStatus(iCompany, iHistoryCode, iDate, policyenq.PolStatus)

	if boolstat {
		err := errors.New("Invalid Benefit Status")
		return 0, "", err
	}

	iCoverage := benefitupd.BCoverage
	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
	if err != nil {
		err := errors.New("Q0006 Not Found")
		return 0, "", err
	}

	if q0006data.PremCalcType != "U" {
		err := errors.New("Not Unit Linked")
		return 0, "", err
	}

	if q0006data.UlMorttMethod == "" {
		err := errors.New("Mortality Method Not Found")
		return 0, "", err

	}

	if q0006data.UlMortFreq == "" {
		err := errors.New("Mortality Frequency Not Found")
		return 0, "", err
	}
	iNextDue := Date2String(GetNextDue(iEffDate, q0006data.UlMortFreq, ""))

	iUlMortalityMethod := q0006data.UlMorttMethod
	oSA := 0.0
	oAmount := 0.0
	switch iUlMortalityMethod {
	case "ULM001":
		oSA = float64(benefitupd.BSumAssured)
	case "ULM002":
		oFund, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)
		oSA = float64(benefitupd.BSumAssured) - oFund

	case "ULM003":
		oFund, _, _ := GetAllFundValueByPol(iCompany, iPolicy, iEffDate)
		oSA = float64(benefitupd.BSumAssured) - oFund
	case "ULM004":
		oFund, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)
		oSA = float64(benefitupd.BSumAssured) + oFund
	case "ULM005":
		oFund, _, _ := GetAllFundValueByPol(iCompany, iPolicy, iEffDate)
		oSA = float64(benefitupd.BSumAssured) + oFund
	default:
		oSA = 0
	}
	fmt.Println("New SA", oSA)
	if oSA < 0 {
		oSA = 0
	}
	iPrem := 0.00
	iAge, _, _, _, _, _ := CalculateAge(benefitupd.BDOB, iEffDate, q0006data.AgeCalcMethod)
	iGender := benefitupd.BGender
	iTerm := benefitupd.BTerm
	iPremTerm := benefitupd.BPTerm
	iPremMethod := q0006data.PremiumMethod
	iMortalityClass := benefitupd.BMortality
	iPrem, err = GetAnnualRate(iCompany, iCoverage, uint(iAge), iGender, iTerm, iPremTerm, iPremMethod, iEffDate, iMortalityClass)
	iPrem = iPrem * oSA / 10000
	iMortFreq := q0006data.UlMortFreq
	mPrem := 0.00
	switch iMortFreq {
	case "M":
		mPrem = iPrem / 12
	case "Q":
		mPrem = iPrem / 4
	case "H":
		mPrem = iPrem / 2
	case "Y":
		mPrem = iPrem / 1
	}

	fmt.Println(iPrem, mPrem)
	oAmount = RoundFloat(mPrem, 2)

	return oAmount, iNextDue, nil
}

// # 130
// CalcPolicyFee - Calculate PolicyFee  for ILP
//
// Inputs: Company, Policy, Benefit Code, Benefit Code, Start Date of Benefit, Effective Date, Fee Method and Fee Frequency
//
// # Outputs  Policy Fee for the Frequency, Next Due Date and Error
//
// ©  FuturaInsTech
func CalcUlPolicyFee(iCompany uint, iPolicy uint, iBenefitID uint, iCoverage string, iStartDate string, iEffDate string, iFeeMethod string, iFeeFreq string) (float64, string, error) {

	iKey := iFeeMethod
	iNextDue := Date2String(GetNextDue(iEffDate, iFeeFreq, ""))
	oAmount := 0.00
	var p0063data paramTypes.P0063Data
	var extradatap0063 paramTypes.Extradata = &p0063data
	err := GetItemD(int(iCompany), "P0063", iKey, iStartDate, &extradatap0063)
	if err != nil {
		err := errors.New("P0063 Not Found")
		return 0, "", err
	}
	switch iKey {
	case "ULFEE01":
		oAmount = p0063data.FlatAmount
	case "ULFEE02":
		tempdate := iStartDate
		i := 0
		for tempdate < iEffDate {
			a := GetNextDue(tempdate, "Y", "")
			tempdate = Date2String(a)
			i++
		}
		iDays := i * 365
		oAmount = SimpleInterest(p0063data.FlatAmount, p0063data.Percentage, float64(iDays))
	case "ULFEE03":
		tempdate := iStartDate
		i := 0
		for tempdate < iEffDate {
			a := GetNextDue(tempdate, "Y", "")
			tempdate = Date2String(a)
			i++
		}
		iDays := i * 365
		oAmount = CompoundInterest(p0063data.FlatAmount, p0063data.Percentage, float64(iDays))

	case "ULFEE04":
		oAmount, _, _ = GetAllFundValueByBenefit(iCompany, iPolicy, iBenefitID, "", iEffDate)
		oAmount = oAmount * p0063data.FundValPercentage / 100

	case "ULFEE05":
		oAmount, _, _ = GetAllFundValueByPol(iCompany, iPolicy, iEffDate)
		oAmount = oAmount * p0063data.FundValPercentage / 100

	}

	if p0063data.CapAmount != 0 {
		if oAmount > p0063data.CapAmount {
			oAmount = p0063data.CapAmount
		}
	}

	mPrem := 0.00
	switch iFeeFreq {
	case "M":
		mPrem = oAmount / 12
	case "Q":
		mPrem = oAmount / 4
	case "H":
		mPrem = oAmount / 2
	case "Y":
		mPrem = oAmount / 1
	}

	oAmount = RoundFloat(mPrem, 2)

	return oAmount, iNextDue, err
}

// # 131
// PostUlpDeduction - Post ILP Deductions (Mortality/Policy Fee or Any other deductions)
//
// Inputs: Company, Policy, Benefit Code, Benefit ID, Amount to be deducted, History Code, Benefit Code, Start Date of Benefit, Effective Date and Tranno
//
// # Outputs  Record is written in ILP Transaction Table
//
// ©  FuturaInsTech
func PostUlpDeduction(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, NegativeUnitsOrAmt string) error {

	var policyenq models.Policy

	result := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}

	var ilpfundenq []models.IlpFund

	result = initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	var ilpsumenq []models.IlpSummary

	result = initializers.DB.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	for j := 0; j < len(ilpsumenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}
		iFundCode := ilpsumenq[j].FundCode
		iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
		if iFundValue == 0 {
			if p0059data.NegativeAccum == "N" {
				// triger TdfLaps
				return nil
			}

		}
		var ilptrancrt models.IlpTransaction
		iKey := ilpsumenq[j].FundCode
		err := GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
		if err != nil {
			return errors.New(err.Error())
		}
		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpsumenq[j].FundCode
		ilptrancrt.FundType = ilpsumenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundEffDate = iBusinessDate
		//ilptrancrt.FundAmount = RoundFloat(((iAmount * ilpfundenq[j].FundPercentage) / 100), 2)

		if NegativeUnitsOrAmt == "U" {
			if iTotalFundValue <= 0 {
				ilptrancrt.FundAmount = RoundFloat(iAmount, 2)
				ilptrancrt.OriginalAmount = RoundFloat(iAmount, 2)
				ilptrancrt.InvNonInvPercentage = 100
			} else {
				ilptrancrt.FundAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
				ilptrancrt.OriginalAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
				ilptrancrt.InvNonInvPercentage = RoundFloat((iFundValue / iTotalFundValue), 5)
			}
		} else {
			ilptrancrt.FundAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
			ilptrancrt.OriginalAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
			ilptrancrt.InvNonInvPercentage = RoundFloat((iFundValue / iTotalFundValue), 5)

		}
		ilptrancrt.FundCurr = p0061data.FundCurr
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture

		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.AccountCode = p0059data.AccountCode
		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "P"
		result = initializers.DB.Create(&ilptrancrt)
	}
	var tdfpolicyupd models.TDFPolicy
	iType := "FUNDM"
	if iHistoryCode == "H0132" {
		iType = "FUNDM"
	} else if iHistoryCode == "H0133" {
		iType = "FUNDF"
	}
	var tdfrule models.TDFRule
	result = initializers.DB.Find(&tdfrule, "company_id = ? and tdf_type  = ?", iCompany, iType)
	result = initializers.DB.Find(&tdfpolicyupd, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iType)
	if result.RowsAffected == 0 {
		tdfpolicyupd.CompanyID = iCompany
		tdfpolicyupd.PolicyID = iPolicy
		tdfpolicyupd.EffectiveDate = iStartDate
		tdfpolicyupd.TDFType = iType
		tdfpolicyupd.Tranno = iTranno
		tdfpolicyupd.Seqno = tdfrule.Seqno
		initializers.DB.Create(&tdfpolicyupd)
	}
	return nil
}

func PostUlpDeductionN(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, NegativeUnitsOrAmt string, txn *gorm.DB) error {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpfundenq []models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpsumenq []models.IlpSummary

	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	for j := 0; j < len(ilpsumenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}
		iFundCode := ilpsumenq[j].FundCode
		iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
		if iFundValue == 0 {
			if p0059data.NegativeAccum == "N" {
				// triger TdfLaps
				return nil
			}

		}
		var ilptrancrt models.IlpTransaction
		iKey := ilpsumenq[j].FundCode
		err := GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
		if err != nil {
			return errors.New(err.Error())
		}
		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpsumenq[j].FundCode
		ilptrancrt.FundType = ilpsumenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundEffDate = iBusinessDate
		//ilptrancrt.FundAmount = RoundFloat(((iAmount * ilpfundenq[j].FundPercentage) / 100), 2)

		if NegativeUnitsOrAmt == "U" {
			if iTotalFundValue <= 0 {
				ilptrancrt.FundAmount = RoundFloat(iAmount, 2)
				ilptrancrt.OriginalAmount = RoundFloat(iAmount, 2)
				ilptrancrt.InvNonInvPercentage = 100
			} else {
				ilptrancrt.FundAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
				ilptrancrt.OriginalAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
				ilptrancrt.InvNonInvPercentage = RoundFloat((iFundValue / iTotalFundValue), 5)
			}
		} else {
			ilptrancrt.FundAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
			ilptrancrt.OriginalAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
			ilptrancrt.InvNonInvPercentage = RoundFloat((iFundValue / iTotalFundValue), 5)

		}
		ilptrancrt.FundCurr = p0061data.FundCurr
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture

		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.AccountCode = p0059data.AccountCode
		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "P"
		result = txn.Create(&ilptrancrt)
		if result.Error != nil {
			txn.Rollback()
			return result.Error
		}

	}
	var tdfpolicyupd models.TDFPolicy
	iType := "FUNDM"
	if iHistoryCode == "H0132" {
		iType = "FUNDM"
	} else if iHistoryCode == "H0133" {
		iType = "FUNDF"
	}
	var tdfrule models.TDFRule
	result = txn.Find(&tdfrule, "company_id = ? and tdf_type  = ?", iCompany, iType)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}
	result = txn.Find(&tdfpolicyupd, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iType)

	if result.RowsAffected == 0 {
		tdfpolicyupd.CompanyID = iCompany
		tdfpolicyupd.PolicyID = iPolicy
		tdfpolicyupd.EffectiveDate = iStartDate
		tdfpolicyupd.TDFType = iType
		tdfpolicyupd.Tranno = iTranno
		tdfpolicyupd.Seqno = tdfrule.Seqno
		result = txn.Create(&tdfpolicyupd)
		if result.Error != nil {
			txn.Rollback()
			return result.Error
		}
	}
	return nil
}

// # 132
// CheckPendingILP - Check Pending ILP Transaction on a Policy
//
// Inputs: Company, Policy, Benefit Code
//
// # Outputs  Error Description
//
// ©  FuturaInsTech
func CheckPendingILP(iCompany uint, iPolicy uint, iLanguage uint) string {

	var ilptransenq models.IlpTransaction

	result := initializers.DB.Find(&ilptransenq, "company_id = ? and policy_id = ? and ul_process_flag = ?", iCompany, iPolicy, "P")
	if result.RowsAffected != 0 {
		shortCode := "E0005"
		longdesc, _ := GetErrorDesc(uint(iCompany), iLanguage, shortCode)
		return shortCode + ": -" + longdesc
	}
	return ""
}

// # 133
// GetMrtaPrem - calculate MRTA Premium (New Version)
//
// Inputs: Company, Benefit Code, Initial SA, Initial Age, Gender, Term , Premium Paying Term, Interest, Interim Period, Start Date
//
// # Outputs  Premium and Error Description
//
// ©  FuturaInsTech
func GetMrtaPrem(iCompany uint, iCoverage string, iSA float64, iAge uint, iGender string, iTerm uint, iPremTerm uint, iInterest float64, iInterimPeriod uint, iDate string) (float64, error) {

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
	if err != nil {
		return 0, err
	}

	var q0010data paramTypes.Q0010Data
	var extradataq0010 paramTypes.Extradata = &q0010data
	var q0010key string
	var prem float64
	prem = 0
	var prem1 float64
	prem1 = 0
	oSA := iSA
	term := strconv.FormatUint(uint64(iTerm), 10)
	premTerm := strconv.FormatUint(uint64(iTerm), 10)

	if q0006data.PremCalcType == "A" {
		q0010key = iCoverage + iGender
	} else if q0006data.PremCalcType == "P" {
		q0010key = iCoverage + iGender + term + premTerm
		// END1 + Male + Term + Premium Term
	}
	err = GetItemD(int(iCompany), "Q0010", q0010key, iDate, &extradataq0010)
	if err != nil {
		return 0, err
	}

	for x := 0; x < int(iTerm); x++ {
		rSA := GetMRTABen(float64(oSA), float64(iInterest), float64(x+1), float64(iInterimPeriod), float64(iTerm))
		for i := 0; i < len(q0010data.Rates); i++ {
			if q0010data.Rates[i].Age == uint(iAge) {
				prem = q0010data.Rates[i].Rate / 10000
				prem1 = (prem * rSA) + prem1
				iAge = iAge + 1
				break
			}
		}
		oSA = rSA
	}
	prem = prem1

	return prem, nil

}

// # 134
// PostTopAllocation - This function apportion amount into different funds and investible and non investible (Top up Only)
//
// Inputs:
//
// # Success/Failure
//
// ©  FuturaInsTech
func PostTopAllocation(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iFrequency string, iStartDate string, iEffDate string, iGender string, iAllocMethod string, iTranno uint) error {

	var policyenq models.Policy

	result := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data
	iDate := iStartDate
	iKey := iAllocMethod + iGender
	err := GetItemD(int(iCompany), "P0060", iKey, iDate, &extradatap0060)
	if err != nil {
		return errors.New(err.Error())
	}
	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey = iHistoryCode + iBenefitCode
	err = GetItemD(int(iCompany), "P0059", iKey, iDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}

	if iEffDate == iStartDate {
		a := GetNextDue(iStartDate, iFrequency, "")
		iEffDate = Date2String(a)
	}
	iNoofMonths := NewNoOfInstalments(iStartDate, iEffDate)
	iAllocPercentage := 0.00
	for i := 0; i < len(p0060data.AlBand); i++ {
		if uint(iNoofMonths) <= p0060data.AlBand[i].Months {
			iAllocPercentage = p0060data.AlBand[i].Percentage
			break
		}
	}
	iInvested := RoundFloat(iAmount*(iAllocPercentage/100), 2)
	iNonInvested := RoundFloat(iAmount*((100-iAllocPercentage)/100), 2)

	var ilpfundenq []models.IlpFund
	// Select  Top-up Funds ONly
	result = initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ? and history_code= ?", iCompany, iPolicy, iBenefit, iHistoryCode)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	for j := 0; j < len(ilpfundenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}

		var ilptrancrt models.IlpTransaction
		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpfundenq[j].FundCode
		ilptrancrt.FundType = ilpfundenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundEffDate = iBusinessDate
		ilptrancrt.FundAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.FundCurr = ilpfundenq[j].FundCurr
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = ilpfundenq[j].FundPercentage
		ilptrancrt.AccountCode = "Invested" // ranga

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "P"
		result = initializers.DB.Create(&ilptrancrt)
	}
	// Non Invested Amount Updation

	var ilptrancrt models.IlpTransaction
	// Move Variables
	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = "NONIN"
	ilptrancrt.FundType = "NI"
	ilptrancrt.TransactionDate = iEffDate
	ilptrancrt.FundEffDate = iEffDate
	ilptrancrt.FundAmount = iNonInvested
	ilptrancrt.FundCurr = ""
	ilptrancrt.FundUnits = 0
	ilptrancrt.FundPrice = 0
	ilptrancrt.CurrentOrFuture = "C"
	ilptrancrt.OriginalAmount = iNonInvested
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "NI"
	ilptrancrt.AllocationCategory = "NI"
	ilptrancrt.InvNonInvPercentage = 0
	ilptrancrt.Tranno = iTranno

	ilptrancrt.AccountCode = "NonInvested"

	ilptrancrt.CurrencyRate = 1.00 // ranga
	ilptrancrt.MortalityIndicator = ""
	ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = initializers.DB.Create(&ilptrancrt)

	// Delete Newly Cleared Fund Rules which is created for Top-up
	initializers.DB.Delete(ilpfundenq)

	return nil

}

func PostTopAllocationN(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iFrequency string, iStartDate string, iEffDate string, iGender string, iAllocMethod string, iTranno uint, txn *gorm.DB) error {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data
	iDate := iStartDate
	iKey := iAllocMethod + iGender
	err := GetItemD(int(iCompany), "P0060", iKey, iDate, &extradatap0060)
	if err != nil {
		return errors.New(err.Error())
	}
	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey = iHistoryCode + iBenefitCode
	err = GetItemD(int(iCompany), "P0059", iKey, iDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}

	if iEffDate == iStartDate {
		a := GetNextDue(iStartDate, iFrequency, "")
		iEffDate = Date2String(a)
	}
	iNoofMonths := NewNoOfInstalments(iStartDate, iEffDate)
	iAllocPercentage := 0.00
	for i := 0; i < len(p0060data.AlBand); i++ {
		if uint(iNoofMonths) <= p0060data.AlBand[i].Months {
			iAllocPercentage = p0060data.AlBand[i].Percentage
			break
		}
	}
	iInvested := RoundFloat(iAmount*(iAllocPercentage/100), 2)
	iNonInvested := RoundFloat(iAmount*((100-iAllocPercentage)/100), 2)

	var ilpfundenq []models.IlpFund
	// Select  Top-up Funds ONly
	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ? and history_code= ?", iCompany, iPolicy, iBenefit, iHistoryCode)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	for j := 0; j < len(ilpfundenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}

		var ilptrancrt models.IlpTransaction
		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpfundenq[j].FundCode
		ilptrancrt.FundType = ilpfundenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundEffDate = iBusinessDate
		ilptrancrt.FundAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.FundCurr = ilpfundenq[j].FundCurr
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat(((iInvested * ilpfundenq[j].FundPercentage) / 100), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = ilpfundenq[j].FundPercentage
		ilptrancrt.AccountCode = "Invested" // ranga

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "P"
		result = txn.Create(&ilptrancrt)
		if result.Error != nil {
			txn.Rollback()
			return result.Error
		}
	}
	// Non Invested Amount Updation

	var ilptrancrt models.IlpTransaction
	// Move Variables
	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = "NONIN"
	ilptrancrt.FundType = "NI"
	ilptrancrt.TransactionDate = iEffDate
	ilptrancrt.FundEffDate = iEffDate
	ilptrancrt.FundAmount = iNonInvested
	ilptrancrt.FundCurr = ""
	ilptrancrt.FundUnits = 0
	ilptrancrt.FundPrice = 0
	ilptrancrt.CurrentOrFuture = "C"
	ilptrancrt.OriginalAmount = iNonInvested
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "NI"
	ilptrancrt.AllocationCategory = "NI"
	ilptrancrt.InvNonInvPercentage = 0
	ilptrancrt.Tranno = iTranno

	ilptrancrt.AccountCode = "NonInvested"

	ilptrancrt.CurrencyRate = 1.00 // ranga
	ilptrancrt.MortalityIndicator = ""
	ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = txn.Create(&ilptrancrt)
	if result.Error != nil {
		return result.Error
	}

	// Delete Newly Cleared Fund Rules which is created for Top-up
	result = txn.Delete(ilpfundenq)
	if result.Error != nil {
		return result.Error
	}

	return nil

}

// # 134
// GetAllowedFunds - This function return all allowed funds for a particular benefit
//
// Inputs: Company, Coverage Code and Coverage Start Date
//
// Outputs : Fund Array consist of fund code, fund category, fund type and fund currency and error code
//
// ©  FuturaInsTech
func GetAllowedFunds(iCompany uint, iCoverage string, iDate string) ([]interface{}, error) {

	fundlist := make([]interface{}, 0)

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
	if err != nil {
		return nil, err
	}
	if q0006data.FUNDCODE == nil {
		return nil, err
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	for i := 0; i < len(q0006data.FUNDCODE); i++ {
		err = GetItemD(int(iCompany), "P0061", q0006data.FUNDCODE[i], iDate, &extradatap0061)
		if err != nil {
			return nil, err
		}
		resultOut := map[string]interface{}{
			"FundCode":     p0061data.FundCode,
			"FundCategory": p0061data.FundCategory,
			"FundCurr":     p0061data.FundCurr,
			"FundType":     p0061data.FundType,
		}
		fmt.Print(fundlist)
		fundlist = append(fundlist, resultOut)
	}
	return fundlist, nil
}

// # 135
// Validate the Address Table Fields mandatory as required by P0065 Rules
//
// Inputs: Address Model, Company, Language, Key (Program name)
//
// Outputs : Error
//
// ©  FuturaInsTech
func ValidateAddress(addressval models.Address, userco uint, userlan uint, iKey string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return errors.New(err.Error())
	}

	for i := 0; i < len(p0065data.FieldList); i++ {
		var fv interface{}
		r := reflect.ValueOf(addressval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) == true {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	return
}

// # 136
// Validate the Client Table Fields mandatory as required by P0065 Rules
//
// Inputs: Cleint Mode, Company, Language, Key (Program name)
//
// Outputs : Error
//
// ©  FuturaInsTech
func ValidateClient(clientval models.Client, userco uint, userlan uint, iKey string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	iClientType := clientval.ClientType
	if iClientType == "I" || iClientType == "C" {
		iKey = iKey + iClientType
	}

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return errors.New(err.Error())
	}

	for i := 0; i < len(p0065data.FieldList); i++ {

		var fv interface{}
		r := reflect.ValueOf(clientval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	validemail := isValidEmail(clientval.ClientEmail)
	if !validemail {
		shortCode := "GL477" // Email Format is invalid
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)
	}

	_, err = strconv.Atoi(clientval.ClientMobile)
	if err != nil {
		shortCode := "GL478" // MobileNumber is not Numeric
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)
	}

	ibusinessdate := GetBusinessDate(userco, 0, 0)
	if clientval.ClientDob > ibusinessdate {
		shortCode := ""
		if clientval.ClientType == "C" {
			shortCode = "GL586" // Incorrect Date of Incorporation
		} else {
			shortCode = "GL566" // Incorrect Date of Birth
		}
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)
	}

	if clientval.ClientDod != "" {
		if clientval.ClientDod <= clientval.ClientDob {
			shortCode := ""
			if clientval.ClientType == "C" {
				shortCode = "GL587" // Incorrect Date of Termination
			} else {
				shortCode = "GL567" // Date of Birth/Death Incorrect
			}
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}
	}

	return
}

// # 137
// Validate the Bank Table Fields mandatory as required by P0065 Rules
//
// Inputs: Bank Model, Company, Language, Key (Program name)
//
// Outputs : Error
//
// ©  FuturaInsTech
func ValidateBank(bankval models.Bank, userco uint, userlan uint, iKey string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return errors.New(err.Error())
	}

	for i := 0; i < len(p0065data.FieldList); i++ {

		var fv interface{}
		r := reflect.ValueOf(bankval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) == true {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	if bankval.StartDate > bankval.EndDate {
		shortCode := "GL563"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)
	}

	return
}

// # 138
//
// # GetIlpFundUnits - To return the current available Units in a Fund
//
// Inputs: Company, Policy, Benefit, Fund Code
//
// # Outputs:  Return the available Units against the given Fund
// # Error:  If given fund does not exist in the policy, then return zeroes.
//
// ©  FuturaInsTech
func GetIlpFundUnits(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string) (float64, error) {
	var ilpsummaryenq models.IlpSummary
	result := initializers.DB.First(&ilpsummaryenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		return 0.0, nil
	}
	oUnits := ilpsummaryenq.FundUnits
	return oUnits, nil
}

// # 142
// *********************************************************************************************
// Method : GetNewPremium
//
// Calculate Annual Premium  and Model Premium  (With Discount and Frequency Loading)
// # Inputs Company, Coverage Code, Start Date of Coverage, Age , Gender, Risk Term, Premium Term, Mortality Class, Sum Assured, Premium Method (Q0006), Discount Type, Discount Method, Frequency Method, Frequency
// # Outputs Annual Premium and Model Premium
// ©  FuturaInsTech
func GetNewPremium(iCompany uint, iCoverage string, iDate string, iAge uint, iGender string, iRiskTerm uint, iPremTerm uint, iMortality string, iSA float64, iPremMethod string, iDiscType string, iDiscMethod string, iFrqMethod string, iFrequency string) (oAnnualPrem float64, oBasePrem float64) {

	prem, _ := GetAnnualRate(iCompany, iCoverage, iAge, iGender, iRiskTerm, iPremTerm, iPremMethod, iDate, iMortality)
	oAnnualPrem = RoundFloat((prem * float64(iSA) / 10000), 2)
	discount := RoundFloat(CalcSaPremDiscount(iCompany, iDiscType, iDiscMethod, oAnnualPrem, uint(iSA), iDate), 2)
	prem1 := oAnnualPrem - discount
	oBasePrem = CalcFrequencyPrem(iCompany, iDate, iFrqMethod, iFrequency, prem1)

	return oAnnualPrem, oBasePrem
}

// #  143
// GetFundCPrice - Fetch the Available Fund Price approved to use as Current Price in calculations
//
// Inputs: Company, Fund Code, Processing Date
//
// # Outputs:  Bid Price, Offer Price and Date of Price
//
// ©  FuturaInsTech
func GetFundCPrice(iCompany uint, iFundCode string, iDate string) (float64, float64, string) {
	var ilppriceenq models.IlpPrice
	var iPriceDateUsed = "00000000"
	result := initializers.DB.Where("company_id = ? and fund_code = ? and approval_flag = ? and fund_eff_date <= ?", iCompany, iFundCode, "AP", iDate).Order("fund_eff_date DESC").First(&ilppriceenq)
	if result.Error != nil {
		return 0.0, 0.0, iPriceDateUsed
	}
	iPriceDateUsed = ilppriceenq.FundEffDate
	iBidPrice := ilppriceenq.FundBidPrice
	iOfferPrice := ilppriceenq.FundOfferPrice
	return iBidPrice, iOfferPrice, iPriceDateUsed
}

// # 144
//
// # PostUlpDeductionByFundAmount - Post ILP Deductions by alloctype for a Specific Fund
//
// Inputs: Company, Policy, Benefit, Fund Code, Amount to be deducted, History Code, Benefit Code, Start Date of Benefit, Effective Date, Tranno and Allocation Type
//
// # Outputs  Record is written in ILP Transaction Table
//
// ©  FuturaInsTech
func PostUlpDeductionByFundAmount(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string) error {

	var policyenq models.Policy

	result := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}

	var ilpfundenq models.IlpFund

	result = initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	var ilpsumenq models.IlpSummary

	result = initializers.DB.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	iBusinessDate := GetBusinessDate(iCompany, 0, 0)
	if p0059data.CurrentOrFuture == "F" {
		iBusinessDate = AddLeadDays(iBusinessDate, 1)
	} else if p0059data.CurrentOrFuture == "E" {
		iBusinessDate = iEffDate
	}

	iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
	var ilptrancrt models.IlpTransaction
	iKey = iFundCode
	err = GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
	if err != nil {
		return errors.New(err.Error())
	}

	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = iFundCode
	ilptrancrt.FundType = ilpsumenq.FundType
	ilptrancrt.TransactionDate = iEffDate

	ilptrancrt.FundAmount = RoundFloat(iAmount, 2)
	ilptrancrt.FundCurr = p0061data.FundCurr

	ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq.FundCode, iBusinessDate)
	ilptrancrt.FundPrice = ibidprice
	ilptrancrt.FundEffDate = ipriceuseddate
	ilptrancrt.FundUnits = RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)

	ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
	ilptrancrt.OriginalAmount = RoundFloat(iAmount, 2)
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.SurrenderPercentage = RoundFloat(((ilptrancrt.FundAmount / iFundValue) * 100), 2)
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "AC"
	ilptrancrt.AllocationCategory = p0059data.AllocationCategory
	ilptrancrt.InvNonInvPercentage = RoundFloat(((iFundValue / iTotalFundValue) * 100), 2)
	ilptrancrt.AccountCode = p0059data.AccountCode

	ilptrancrt.CurrencyRate = 1.00 // ranga
	ilptrancrt.MortalityIndicator = ""
	//ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Tranno = iTranno
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = initializers.DB.Create(&ilptrancrt)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	//update ilpsummary
	var ilpsummupd models.IlpSummary
	result = initializers.DB.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

	if result.RowsAffected != 0 {
		ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
		initializers.DB.Save(&ilpsummupd)
	} else {
		return errors.New(err.Error())
	}

	return nil
}

func PostUlpDeductionByFundAmountN(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string, txn *gorm.DB) error {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpfundenq models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpsumenq models.IlpSummary

	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	iBusinessDate := GetBusinessDate(iCompany, 0, 0)
	if p0059data.CurrentOrFuture == "F" {
		iBusinessDate = AddLeadDays(iBusinessDate, 1)
	} else if p0059data.CurrentOrFuture == "E" {
		iBusinessDate = iEffDate
	}

	iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
	var ilptrancrt models.IlpTransaction
	iKey = iFundCode
	err = GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = iFundCode
	ilptrancrt.FundType = ilpsumenq.FundType
	ilptrancrt.TransactionDate = iEffDate

	ilptrancrt.FundAmount = RoundFloat(iAmount, 2)
	ilptrancrt.FundCurr = p0061data.FundCurr

	ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq.FundCode, iBusinessDate)
	ilptrancrt.FundPrice = ibidprice
	ilptrancrt.FundEffDate = ipriceuseddate
	ilptrancrt.FundUnits = RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)

	ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
	ilptrancrt.OriginalAmount = RoundFloat(iAmount, 2)
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.SurrenderPercentage = RoundFloat(((ilptrancrt.FundAmount / iFundValue) * 100), 2)
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "AC"
	ilptrancrt.AllocationCategory = p0059data.AllocationCategory
	ilptrancrt.InvNonInvPercentage = RoundFloat(((iFundValue / iTotalFundValue) * 100), 2)
	ilptrancrt.AccountCode = p0059data.AccountCode

	ilptrancrt.CurrencyRate = 1.00 // ranga
	ilptrancrt.MortalityIndicator = ""
	//ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Tranno = iTranno
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = txn.Create(&ilptrancrt)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	//update ilpsummary
	var ilpsummupd models.IlpSummary
	result = txn.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

	if result.RowsAffected != 0 {
		ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
		txn.Save(&ilpsummupd)
	} else {
		txn.Rollback()
		return errors.New(err.Error())
	}

	return nil
}

// # 145
//
// PostUlpDeductionByAmount - Post ILP Deductions by alloctype (used in Surrender & PartSurrender Penalty & GST Postings )
//
// Inputs: Company, Policy, Benefit Code, Benefit ID, Amount to be deducted, History Code, Benefit Code, Start Date of Benefit, Effective Date, Tranno and Allocation Type
//
// # Outputs  Record is written in ILP Transaction Table
//
// ©  FuturaInsTech
func PostUlpDeductionByAmount(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string) error {

	var policyenq models.Policy

	result := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}

	var ilpfundenq []models.IlpFund

	result = initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	var ilpsumenq []models.IlpSummary

	result = initializers.DB.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	for j := 0; j < len(ilpsumenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}
		iFundCode := ilpsumenq[j].FundCode
		iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
		var ilptrancrt models.IlpTransaction
		iKey := ilpsumenq[j].FundCode
		err := GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
		if err != nil {
			return errors.New(err.Error())
		}

		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpsumenq[j].FundCode
		ilptrancrt.FundType = ilpsumenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundAmount = RoundFloat((iAmount * iFundValue / iTotalFundValue), 2)
		ilptrancrt.FundCurr = p0061data.FundCurr

		ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq[j].FundCode, iBusinessDate)
		ilptrancrt.FundPrice = ibidprice
		ilptrancrt.FundEffDate = ipriceuseddate
		ilptrancrt.FundUnits = RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)

		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat((iAmount * iFundValue / iTotalFundValue), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr

		ilptrancrt.SurrenderPercentage = RoundFloat(((iFundValue / iTotalFundValue) * 100), 2)
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = 0
		ilptrancrt.AccountCode = p0059data.AccountCode

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		//ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "C"
		result = initializers.DB.Create(&ilptrancrt)
		if result.Error != nil {
			return errors.New(err.Error())
		}

		//update ilpsummary
		var ilpsummupd models.IlpSummary
		result = initializers.DB.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

		if result.RowsAffected != 0 {
			ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
			initializers.DB.Save(&ilpsummupd)
		} else {
			return errors.New(err.Error())
		}
	}
	return nil
}

// # 145
//
// PostUlpDeductionByAmountN - Post ILP Deductions by alloctype (used in Surrender & PartSurrender Penalty & GST Postings ) Using txn
//
// Inputs: Company, Policy, Benefit Code, Benefit ID, Amount to be deducted, History Code, Benefit Code, Start Date of Benefit, Effective Date, Tranno and Allocation Type
//
// # Outputs  Record is written in ILP Transaction Table
//
// ©  FuturaInsTech
func PostUlpDeductionByAmountN(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string, txn *gorm.DB) error {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpfundenq []models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpsumenq []models.IlpSummary

	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	for j := 0; j < len(ilpsumenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}
		iFundCode := ilpsumenq[j].FundCode
		iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
		var ilptrancrt models.IlpTransaction
		iKey := ilpsumenq[j].FundCode
		err := GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
		if err != nil {
			txn.Rollback()
			return errors.New(err.Error())
		}

		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpsumenq[j].FundCode
		ilptrancrt.FundType = ilpsumenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ilptrancrt.FundAmount = RoundFloat((iAmount * iFundValue / iTotalFundValue), 2)
		ilptrancrt.FundCurr = p0061data.FundCurr

		ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq[j].FundCode, iBusinessDate)
		ilptrancrt.FundPrice = ibidprice
		ilptrancrt.FundEffDate = ipriceuseddate
		ilptrancrt.FundUnits = RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)

		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat((iAmount * iFundValue / iTotalFundValue), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr

		ilptrancrt.SurrenderPercentage = RoundFloat(((iFundValue / iTotalFundValue) * 100), 2)
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = 0
		ilptrancrt.AccountCode = p0059data.AccountCode

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		//ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "C"
		result = txn.Create(&ilptrancrt)
		if result.Error != nil {
			txn.Rollback()
			return errors.New(err.Error())
		}

		//update ilpsummary
		var ilpsummupd models.IlpSummary
		result = txn.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

		if result.RowsAffected != 0 {
			fundunit := 0.0

			fundunit = RoundFloat(ilpsummupd.FundUnits*-1, 5)

			ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+fundunit, 5)
			txn.Save(&ilpsummupd)
		} else {
			txn.Rollback()
			return errors.New(err.Error())
		}
	}
	return nil
}

// # 148
//
// # PostUlpDeductionByFundUnits - Post ILP Deductions by alloctype for a Specific Fund
//
// Inputs: Company, Policy, Benefit, Fund Code, % of Units to be deducted, History Code, Benefit Code, Start Date of Benefit, Effective Date, Tranno and Allocation Type
//
// # Outputs  Record is written in ILP Transaction Table
//
// ©  FuturaInsTech
func PostUlpDeductionByFundUnits(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string, iSurrPercentage float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string) error {

	var policyenq models.Policy

	result := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}

	var ilpfundenq models.IlpFund

	result = initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	var ilpsumenq models.IlpSummary

	result = initializers.DB.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	iBusinessDate := GetBusinessDate(iCompany, 0, 0)
	if p0059data.CurrentOrFuture == "F" {
		iBusinessDate = AddLeadDays(iBusinessDate, 1)
	} else if p0059data.CurrentOrFuture == "E" {
		iBusinessDate = iEffDate
	}

	iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
	var ilptrancrt models.IlpTransaction
	iKey = iFundCode
	err = GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
	if err != nil {
		return errors.New(err.Error())
	}

	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = ilpsumenq.FundCode
	ilptrancrt.FundType = ilpsumenq.FundType
	ilptrancrt.TransactionDate = iEffDate
	ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq.FundCode, iBusinessDate)
	ilptrancrt.FundPrice = ibidprice
	ilptrancrt.FundEffDate = ipriceuseddate
	iUnits, _ := GetIlpFundUnits(iCompany, iPolicy, iBenefit, iFundCode)
	// Full Withdrawl is -100% and Part Withdrawl is -20% or -30% etc
	iSurrUnits := iUnits * iSurrPercentage / 100
	ilptrancrt.FundUnits = RoundFloat(iSurrUnits, 5)
	//utilities.RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)
	ilptrancrt.FundAmount = RoundFloat((iSurrUnits * ibidprice), 2)
	ilptrancrt.FundCurr = p0061data.FundCurr
	ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
	ilptrancrt.OriginalAmount = RoundFloat((iSurrUnits * ibidprice), 2)
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.SurrenderPercentage = RoundFloat(((ilptrancrt.FundAmount / iFundValue) * 100), 2)
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "AC"
	ilptrancrt.AllocationCategory = p0059data.AllocationCategory
	ilptrancrt.InvNonInvPercentage = RoundFloat(((ilptrancrt.FundAmount / iTotalFundValue) * 100), 2)
	ilptrancrt.AccountCode = p0059data.AccountCode

	ilptrancrt.CurrencyRate = 1.00 // ranga
	ilptrancrt.MortalityIndicator = ""
	//ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Tranno = iTranno
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = initializers.DB.Create(&ilptrancrt)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	//update ilpsummary
	var ilpsummupd models.IlpSummary
	result = initializers.DB.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

	if result.RowsAffected != 0 {
		ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
		initializers.DB.Save(&ilpsummupd)
	} else {
		return errors.New(err.Error())
	}

	return nil
}

func PostUlpDeductionByFundUnitsN(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string, iSurrPercentage float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string, txn *gorm.DB) error {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpfundenq models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpsumenq models.IlpSummary

	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	iBusinessDate := GetBusinessDate(iCompany, 0, 0)
	if p0059data.CurrentOrFuture == "F" {
		iBusinessDate = AddLeadDays(iBusinessDate, 1)
	} else if p0059data.CurrentOrFuture == "E" {
		iBusinessDate = iEffDate
	}

	iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
	var ilptrancrt models.IlpTransaction
	iKey = iFundCode
	err = GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	ilptrancrt.CompanyID = iCompany
	ilptrancrt.PolicyID = iPolicy
	ilptrancrt.BenefitID = iBenefit
	ilptrancrt.FundCode = ilpsumenq.FundCode
	ilptrancrt.FundType = ilpsumenq.FundType
	ilptrancrt.TransactionDate = iEffDate
	ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq.FundCode, iBusinessDate)
	ilptrancrt.FundPrice = ibidprice
	ilptrancrt.FundEffDate = ipriceuseddate
	iUnits, _ := GetIlpFundUnits(iCompany, iPolicy, iBenefit, iFundCode)
	// Full Withdrawl is -100% and Part Withdrawl is -20% or -30% etc
	iSurrUnits := iUnits * iSurrPercentage / 100
	ilptrancrt.FundUnits = RoundFloat(iSurrUnits, 5)
	//utilities.RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)
	ilptrancrt.FundAmount = RoundFloat((iSurrUnits * ibidprice), 2)
	ilptrancrt.FundCurr = p0061data.FundCurr
	ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
	ilptrancrt.OriginalAmount = RoundFloat((iSurrUnits * ibidprice), 2)
	ilptrancrt.ContractCurry = policyenq.PContractCurr
	ilptrancrt.SurrenderPercentage = RoundFloat(((ilptrancrt.FundAmount / iFundValue) * 100), 2)
	ilptrancrt.HistoryCode = iHistoryCode
	ilptrancrt.InvNonInvFlag = "AC"
	ilptrancrt.AllocationCategory = p0059data.AllocationCategory
	ilptrancrt.InvNonInvPercentage = RoundFloat(((ilptrancrt.FundAmount / iTotalFundValue) * 100), 2)
	ilptrancrt.AccountCode = p0059data.AccountCode

	ilptrancrt.CurrencyRate = 1.00 // ranga
	ilptrancrt.MortalityIndicator = ""
	//ilptrancrt.SurrenderPercentage = 0
	ilptrancrt.Tranno = iTranno
	ilptrancrt.Seqno = uint(p0059data.SeqNo)
	ilptrancrt.UlProcessFlag = "C"
	result = txn.Create(&ilptrancrt)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	//update ilpsummary
	var ilpsummupd models.IlpSummary
	result = txn.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

	if result.RowsAffected != 0 {
		ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
		txn.Save(&ilpsummupd)
	} else {
		txn.Rollback()
		return errors.New(err.Error())
	}

	return nil
}

// # 149
//
// PostUlpDeductionByUnits - Post ILP Deductions by alloctype (used in PartSurrender Penalty & GST Postings )
//
// Inputs: Company, Policy, Benefit Code, Benefit ID, % of Units to be deducted, History Code, Benefit Code, Start Date of Benefit, Effective Date, Tranno and Allocation Type
//
// # Outputs  Record is written in ILP Transaction Table
//
// ©  FuturaInsTech
func PostUlpDeductionByUnits(iCompany uint, iPolicy uint, iBenefit uint, iSurrPercentage float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string) error {

	var policyenq models.Policy

	result := initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		return errors.New(err.Error())
	}

	var ilpfundenq []models.IlpFund

	result = initializers.DB.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	var ilpsumenq []models.IlpSummary

	result = initializers.DB.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	for j := 0; j < len(ilpsumenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}
		iFundCode := ilpsumenq[j].FundCode
		iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
		var ilptrancrt models.IlpTransaction
		iKey := ilpsumenq[j].FundCode
		err := GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
		if err != nil {
			return errors.New(err.Error())
		}

		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpsumenq[j].FundCode
		ilptrancrt.FundType = ilpsumenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq[j].FundCode, iBusinessDate)
		ilptrancrt.FundPrice = ibidprice
		ilptrancrt.FundEffDate = ipriceuseddate
		iUnits, _ := GetIlpFundUnits(iCompany, iPolicy, iBenefit, iFundCode)
		// Full Withdrawl is -100% and Part Withdrawl is -20% or -30% etc
		iSurrUnits := iUnits * iSurrPercentage / 100
		ilptrancrt.FundUnits = RoundFloat(iSurrUnits, 5)
		//utilities.RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)
		ilptrancrt.FundAmount = RoundFloat((iSurrUnits * ibidprice), 2)
		ilptrancrt.FundCurr = p0061data.FundCurr
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat((iSurrUnits * ibidprice), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.SurrenderPercentage = RoundFloat(((ilptrancrt.FundAmount / iFundValue) * 100), 2)
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = RoundFloat(((ilptrancrt.FundAmount / iTotalFundValue) * 100), 2)
		ilptrancrt.AccountCode = p0059data.AccountCode

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		//ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "C"
		result = initializers.DB.Create(&ilptrancrt)
		if result.Error != nil {
			return errors.New(err.Error())
		}

		//update ilpsummary
		var ilpsummupd models.IlpSummary
		result = initializers.DB.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

		if result.RowsAffected != 0 {
			ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
			initializers.DB.Save(&ilpsummupd)
		} else {
			return errors.New(err.Error())
		}
	}
	return nil
}

func PostUlpDeductionByUnitsN(iCompany uint, iPolicy uint, iBenefit uint, iSurrPercentage float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string, txn *gorm.DB) error {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	err := GetItemD(int(iCompany), "P0059", iKey, iStartDate, &extradatap0059)
	if err != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpfundenq []models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	var ilpsumenq []models.IlpSummary

	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		txn.Rollback()
		return errors.New(err.Error())
	}

	// Get Total Fund Value
	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

	for j := 0; j < len(ilpsumenq); j++ {
		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
		if p0059data.CurrentOrFuture == "F" {
			iBusinessDate = AddLeadDays(iBusinessDate, 1)
		} else if p0059data.CurrentOrFuture == "E" {
			iBusinessDate = iEffDate
		}
		iFundCode := ilpsumenq[j].FundCode
		iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
		var ilptrancrt models.IlpTransaction
		iKey := ilpsumenq[j].FundCode
		err := GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
		if err != nil {
			txn.Rollback()
			return errors.New(err.Error())
		}

		ilptrancrt.CompanyID = iCompany
		ilptrancrt.PolicyID = iPolicy
		ilptrancrt.BenefitID = iBenefit
		ilptrancrt.FundCode = ilpsumenq[j].FundCode
		ilptrancrt.FundType = ilpsumenq[j].FundType
		ilptrancrt.TransactionDate = iEffDate
		ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq[j].FundCode, iBusinessDate)
		ilptrancrt.FundPrice = ibidprice
		ilptrancrt.FundEffDate = ipriceuseddate
		iUnits, _ := GetIlpFundUnits(iCompany, iPolicy, iBenefit, iFundCode)
		// Full Withdrawl is -100% and Part Withdrawl is -20% or -30% etc
		iSurrUnits := iUnits * iSurrPercentage / 100
		ilptrancrt.FundUnits = RoundFloat(iSurrUnits, 5)
		//utilities.RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)
		ilptrancrt.FundAmount = RoundFloat((iSurrUnits * ibidprice), 2)
		ilptrancrt.FundCurr = p0061data.FundCurr
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat((iSurrUnits * ibidprice), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.SurrenderPercentage = RoundFloat(((ilptrancrt.FundAmount / iFundValue) * 100), 2)
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = RoundFloat(((ilptrancrt.FundAmount / iTotalFundValue) * 100), 2)
		ilptrancrt.AccountCode = p0059data.AccountCode

		ilptrancrt.CurrencyRate = 1.00 // ranga
		ilptrancrt.MortalityIndicator = ""
		//ilptrancrt.SurrenderPercentage = 0
		ilptrancrt.Tranno = iTranno
		ilptrancrt.Seqno = uint(p0059data.SeqNo)
		ilptrancrt.UlProcessFlag = "C"
		result = txn.Create(&ilptrancrt)
		if result.Error != nil {
			txn.Rollback()
			return errors.New(err.Error())
		}

		//update ilpsummary
		var ilpsummupd models.IlpSummary
		result = txn.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

		if result.RowsAffected != 0 {
			ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
			txn.Save(&ilpsummupd)
		} else {
			txn.Rollback()
			return errors.New(err.Error())
		}
	}
	return nil
}

// # 150
//
// # GetIlpFundUnits - To return the current available Units in a Fund
//
// Inputs: Company, Policy, Benefit, Fund Code
//
// # Outputs:  Return the available Units against the given Fund
// # Error:  If given fund does not exist in the policy, then return zeroes.
//
// ©  FuturaInsTech
func GetIlpFundByUnits(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string) (float64, error) {
	var ilpsummaryenq models.IlpSummary
	result := initializers.DB.First(&ilpsummaryenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.Error != nil {
		return 0.0, nil
	}
	oUnits := ilpsummaryenq.FundUnits
	return oUnits, nil
}

// # 151
// Validate the ValidatePolicyFields mandatory as required by P0065 Rules
// Input: Policy Model, Company, User Language, Transaction Code and Coverage Code
// Output: Error
//
// ©  FuturaInsTech
func ValidatePolicyFields(policyval models.Policy, userco uint, userlan uint, iKey string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return errors.New(err.Error())
	}

	for i := 0; i < len(p0065data.FieldList); i++ {

		var fv interface{}
		r := reflect.ValueOf(policyval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	return
}

// # 152
// Validate the BenefitFields mandatory as required by P0065 Rules
// Input: Benefit Model, Company, User Language, Transaction Code and Coverage Code
// Output: Error
//
// ©  FuturaInsTech
func ValidateBenefitFields(benefitval models.Benefit, userco uint, userlan uint, iHistoryCode string, iCoverage string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	iKey := iHistoryCode + iCoverage

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		iKey = iHistoryCode
		err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	for i := 0; i < len(p0065data.FieldList); i++ {
		var fv interface{}
		r := reflect.ValueOf(benefitval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	return
}

// # 153
// Validate the MrtaFields mandatory as required by P0065 Rules
// ???
//
// Input: MRTA Model, Company, User Language, Transaction Code and Coverage Code
// Output: Error
//
// ©  FuturaInsTech
func ValidateMrtaFields(mrtaval models.Mrta, userco uint, userlan uint, iHistoryCode string, iCoverage string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	iKey := iHistoryCode + iCoverage

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		iKey = iHistoryCode
		err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	for i := 0; i < len(p0065data.FieldList); i++ {
		var fv interface{}
		r := reflect.ValueOf(mrtaval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	return
}

// # 154
// Validate given field if  zero or blank
// Input: Field name
// Output: True when the field is zero or blank and False otherwise...
//
// ©  FuturaInsTech
func isFieldZero(field interface{}) bool {
	v := reflect.ValueOf(field)

	// Check if the field is a valid type
	if v.IsValid() {
		zero := reflect.Zero(v.Type()).Interface()
		return reflect.DeepEqual(field, zero)
	}

	return false // Field is not a valid type
}

// # 155
// ILP Products Only.  Return SA based on Premium and Frequency
// Input: Company, Policy No, Coverage Code, History Code, Effective Date, Age (Inception Age), Instalment Premium (Not Annualized), Frequency, SA
// Output: Error and New Calcualted SA
//
// ©  FuturaInsTech
func CalcILPSA(iCompany uint, iPolicy uint, iCoverage string, iHistoryCD string, iDate string, iAge uint, iPrem float64, iFreq string, iSA float64) (oErr error, oSA float64) {
	var p0068data paramTypes.P0068Data
	var extradatap0068 paramTypes.Extradata = &p0068data
	iKey := iHistoryCD + iCoverage
	prem := 0.0
	switch iFreq {
	case "M":
		prem = iPrem * 12
	case "Q":
		prem = iPrem * 4
	case "H":
		prem = iPrem * 2
	case "Y":
		prem = iPrem * 1
	case "S":
		prem = iPrem * 1
	}
	err := GetItemD(int(iCompany), "P0068", iKey, iDate, &extradatap0068)
	if err != nil {
		return err, 0

	}
	err = errors.New("In Valid SA")
	// Multiplier Logic
	for i := 0; i < len(p0068data.RangeArray); i++ {
		if p0068data.RangeArray[i].P0068Basis == "M" {
			if iAge <= p0068data.RangeArray[i].Age {
				oSA = p0068data.RangeArray[i].Factor * prem
				return nil, oSA
			}
		}
		if p0068data.RangeArray[i].P0068Basis == "R" {
			if iAge <= p0068data.RangeArray[i].Age {
				if iSA < p0068data.RangeArray[i].FromSA {
					oSA = p0068data.RangeArray[i].FromSA
					return nil, oSA
				}
				if iSA > p0068data.RangeArray[i].ToSA {
					oSA = p0068data.RangeArray[i].ToSA
					return nil, oSA
				}
				return nil, iSA
			}
		}
	}

	return err, oSA
}

// # 156
// ILP Products Only.  Switch Fund to In Active Fund During lapsation
// Input: Company, Policy No, Benefit ID, Tranno, TargetFund, Effective Date
// Output: Error // Write Records in ILPSWITCHHEADER and ILPSWITCHFUNDS
//
// ©  FuturaInsTech
func FundSwitch(iCompany uint, iPolicy uint, iBenefit uint, iTranno uint, iTargetFund string, iEffectiveDate string) (oerror error) {
	var ilpswitchheader models.IlpSwitchHeader
	var ilpswitchfunds models.IlpSwitchFund
	var ilpsummary []models.IlpSummary

	result := initializers.DB.Find(&ilpsummary, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.Error != nil {
		return errors.New("Funds Not Found")
	}
	// Switch 100 % from All funds and transfer it to Target Fund

	ilpswitchheader.BenefitID = iBenefit
	ilpswitchheader.CompanyID = iCompany
	ilpswitchheader.EffectiveDate = iEffectiveDate
	ilpswitchheader.FundSwitchBasis = "P"
	ilpswitchheader.PolicyID = iPolicy
	ilpswitchheader.Tranno = iTranno
	iTotalAmount := 0.0
	initializers.DB.Create(&ilpswitchheader)
	// We do not need to delete ilpsummary
	for i := 0; i < len(ilpsummary); i++ {
		if ilpsummary[i].FundUnits > 0 {
			ilpswitchfunds.ID = 0
			ilpswitchfunds.BenefitID = ilpsummary[i].BenefitID
			ilpswitchfunds.CompanyID = ilpsummary[i].CompanyID
			ilpswitchfunds.FundCode = ilpsummary[i].FundCode
			ilpswitchfunds.FundCurr = ilpsummary[i].FundCurr
			ilpswitchfunds.FundType = ilpsummary[i].FundType
			ilpswitchfunds.PolicyID = ilpsummary[i].PolicyID
			ilpswitchfunds.Tranno = iTranno
			ilpswitchfunds.EffectiveDate = iEffectiveDate
			ilpswitchfunds.FundPercentage = -100
			ilpswitchfunds.SwitchDirection = "S"
			_, ilpswitchfunds.FundPrice, _ = GetFundCPrice(iCompany, ilpsummary[i].FundCode, iEffectiveDate)
			_, iFundAmount, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, ilpsummary[i].FundCode, iEffectiveDate)
			iFundAmount = RoundFloat(iFundAmount, 2) * -1
			ilpswitchfunds.FundAmount = iFundAmount
			ilpswitchfunds.FundUnits = ilpsummary[i].FundUnits * -1
			ilpswitchfunds.IlpSwitchHeaderID = ilpswitchheader.ID
			iTotalAmount = iFundAmount + iTotalAmount
			initializers.DB.Create(&ilpswitchfunds)
			// Set Summary Units to Zero.
			ilpsummary[i].FundUnits = 0
			initializers.DB.Save(&ilpsummary[i])
		}
	}
	// Write Target
	iKey := iTargetFund

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	err := GetItemD(int(iCompany), "P0061", iKey, iEffectiveDate, &extradatap0061)

	if err != nil {
		shortCode := "GL442"
		longDesc, _ := GetErrorDesc(iCompany, 1, shortCode)
		return errors.New(shortCode + " : " + longDesc)

	}

	ilpswitchfunds.ID = 0
	ilpswitchfunds.BenefitID = iBenefit
	ilpswitchfunds.CompanyID = iCompany
	ilpswitchfunds.FundCode = iTargetFund
	ilpswitchfunds.FundCurr = p0061data.FundCurr
	ilpswitchfunds.FundType = p0061data.FundType
	ilpswitchfunds.PolicyID = iPolicy
	ilpswitchfunds.Tranno = iTranno
	ilpswitchfunds.EffectiveDate = iEffectiveDate
	ilpswitchfunds.FundPercentage = 100
	ilpswitchfunds.SwitchDirection = "T"
	ilpswitchfunds.FundAmount = RoundFloat(iTotalAmount, 2) * -1
	ilpswitchfunds.FundPrice, _, _ = GetFundCPrice(iCompany, iTargetFund, iEffectiveDate)
	ilpswitchfunds.FundUnits = RoundFloat(iTotalAmount/ilpswitchfunds.FundPrice*-1, 5)
	ilpswitchfunds.IlpSwitchHeaderID = ilpswitchheader.ID
	initializers.DB.Create(&ilpswitchfunds)

	// Create ILP Summary
	var ilpsum models.IlpSummary

	ilpsum.PolicyID = ilpswitchfunds.PolicyID
	ilpsum.BenefitID = ilpswitchfunds.BenefitID
	ilpsum.FundCode = ilpswitchfunds.FundCode
	ilpsum.FundType = ilpswitchfunds.FundType
	ilpsum.FundUnits = ilpswitchfunds.FundUnits
	ilpsum.FundCurr = ilpswitchfunds.FundCurr
	ilpsum.CompanyID = ilpswitchfunds.CompanyID

	initializers.DB.Create(&ilpsum)

	return nil
}

// #159
// ILP Products Only.  Calculate Switch Fee
// Input: Company, Policy No, Fee Method, Benefit Start Date
// Output: Error , Fee Amount, Fee Percentage
//
// ©  FuturaInsTech
func CalcSwitchFee(iCompany uint, iPolicy uint, iFeeMethod string, iEffectiveDate string) (oError error, oAmount float64, oPercentage float64) {

	var p0070data paramTypes.P0070Data
	var extradatap0070 paramTypes.Extradata = &p0070data
	iKey := iFeeMethod
	err := GetItemD(int(iCompany), "P0070", iKey, iEffectiveDate, &extradatap0070)

	if err != nil {
		shortCode := "GL535"
		longDesc, _ := GetErrorDesc(iCompany, 1, shortCode)
		return errors.New(shortCode + " : " + longDesc), 0, 0

	}
	// Get Next Policy Anniversary
	iBusinessDate := GetBusinessDate(iCompany, 0, 0)
	iPolicyAnniversary := iEffectiveDate
	// 20200328
	// a = 2024:03:28:00:00:00:00:00:00
	// 20230328
	//

	for {

		a := GetNextDue(iPolicyAnniversary, "Y", "")
		iPolicyAnniversary = Date2String(a)

		if iPolicyAnniversary >= iBusinessDate {
			break
		}

	}
	b := GetNextDue(iPolicyAnniversary, "Y", "R")
	iPolicyAnniversary = Date2String(b)
	// Get No of Free Done in Policy Anniversary
	var policyhistory []models.PHistory
	results := initializers.DB.Find(&policyhistory, "company_id = ? and policy_id = ? and effective_date >=? and effective_date <=? and is_reversed = ? and history_code = ?", iCompany, iPolicy, iPolicyAnniversary, iBusinessDate, 0, "H0139")
	switchcount := 0
	if results.Error != nil {
		switchcount = 0
	} else {
		switchcount = len(policyhistory)

	}
	// Percentage
	if p0070data.SwitchFeeBasis == "P" {
		if uint(switchcount) <= p0070data.FreeSwitches {
			return nil, 0, 0
		} else {
			return nil, 0, p0070data.FeePercentage
		}
	}
	// Fixed Amount
	if p0070data.SwitchFeeBasis == "F" {
		if uint(switchcount) <= p0070data.FreeSwitches {
			return nil, 0, 0
		} else {
			return nil, p0070data.FeeAmount, 0
		}
	}

	return oError, 0, 0
}

// # 160
// # To Validate Email.
func isValidEmail(email string) bool {
	// Define a regular expression pattern for a standard email format
	// This pattern is a simplified version and may not catch all edge cases
	// For a more comprehensive pattern, consider using a library like "github.com/badoux/checkmail"
	// or a more complex regex.

	//pattern := `/.+@.+\\..+/i`
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Use the MatchString function to check if the email matches the pattern
	return regex.MatchString(email)

}

// # 39
// Post GL - Get Transaction No and History Code (New Version with Rollback)
//
// ©  FuturaInsTech
func PostGlMoveN(iCompany uint, iContractCurry string, iEffectiveDate string,
	iTranno int, iGlAmount float64, iAccAmount float64, iAccountCodeID uint, iGlRdocno uint,
	iGlRldgAcct string, iSeqnno uint64, iGlSign string, iAccountCode string, iHistoryCode string, iRevInd string, iCoverage string, txn *gorm.DB) error {

	iAccAmount = RoundFloat(iAccAmount, 2)

	var glmove models.GlMove
	var company models.Company
	glmove.ContractCurry = iContractCurry
	glmove.ContractAmount = iAccAmount
	result := txn.Find(&company, "id = ?", iCompany)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}
	var currency models.Currency
	// fmt.Println("Currency Code is .... ", company.CurrencyID)
	result = txn.Find(&currency, "id = ?", company.CurrencyID)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}

	iGlCurry := currency.CurrencyShortName
	glmove.CurrencyRate = 1
	if glmove.GlCurry != glmove.ContractCurry {
		var p0031data paramTypes.P0031Data
		var extradata paramTypes.Extradata = &p0031data
		iKey := iContractCurry + "2" + iGlCurry
		// fmt.Println("i key ", iKey)
		err := GetItemD(int(iCompany), "P0031", iKey, iEffectiveDate, &extradata)
		if err != nil {
			fmt.Println("I am inside Error in Exchange ")
			glmove.CurrencyRate = 1
		} else {
			for i := 0; i < len(p0031data.CurrencyRates); i++ {
				// fmt.Println("Exchange Rates", p0031data.CurrencyRates[i].Action, p0031data.CurrencyRates[i].Action)
			}
			glmove.CurrencyRate = p0031data.CurrencyRates[0].Rate
			// fmt.Println("I am outside Error in Exchange ")

		}
	}

	// fmt.Println("Exchange Rate &&&&&&&&&&&&&&&&&&&&", glmove.CurrencyRate)

	glmove.AccountCode = iAccountCode
	glmove.AccountCodeID = iAccountCodeID
	glmove.CompanyID = iCompany
	glmove.EffectiveDate = iEffectiveDate
	glmove.GlAmount = iAccAmount * glmove.CurrencyRate
	glmove.GlCurry = iGlCurry
	glmove.GlRdocno = strconv.Itoa(int(iGlRdocno))
	glmove.GlRldgAcct = iGlRldgAcct
	glmove.GlSign = iGlSign
	glmove.HistoryCode = iHistoryCode
	glmove.Tranno = uint(iTranno)
	glmove.SequenceNo = uint64(iSeqnno)
	curr_date := time.Now()
	glmove.CurrentDate = Date2String(curr_date)

	GlRdocno := glmove.GlRdocno
	glmove.ID = 0
	glmove.ReversalIndicator = iRevInd
	glmove.BCoverage = iCoverage
	result = txn.Save(&glmove)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}
	//tx := initializers.DB.Save(&glmove)
	//tx.Commit()

	UpdateGlBalN(iCompany, iGlRldgAcct, iAccountCode, iContractCurry, iAccAmount, iGlSign, GlRdocno, txn)
	return nil
}

// # 41
// UpdateGlBal (New Version with Rollback)
//
// ©  FuturaInsTech
func UpdateGlBalN(iCompany uint, iGlRldgAcct string, iGlAccountCode string, iContCurry string, iAmount float64, iGLSign string, iGlRdocno string, txn *gorm.DB) (error, float64) {
	var glbal models.GlBal
	var temp float64
	if iGLSign == "-" {
		temp = iAmount * -1

	} else {
		temp = iAmount * 1
	}
	var company []models.Company
	result := txn.First(&company, "id = ?", iCompany)
	if result.Error != nil {
		txn.Rollback()
		return result.Error, 0
	}
	results := txn.First(&glbal, "company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno)
	// if results.Error != nil {
	// 	return errors.New("Account Code Not Found"), glbal.ContractAmount
	// }
	if results.RowsAffected == 0 {
		glbal.ContractAmount = temp
		glbal.CompanyID = iCompany
		glbal.GlAccountno = iGlAccountCode
		glbal.GlRldgAcct = iGlRldgAcct
		glbal.ContractCurry = iContCurry
		glbal.GlRdocno = iGlRdocno
		//initializers.DB.Save(&glbal)
		result = txn.Create(&glbal)
		if result.Error != nil {
			txn.Rollback()
			return result.Error, glbal.ContractAmount
		}
		return nil, glbal.ContractAmount
	} else {
		iAmount := glbal.ContractAmount + temp
		// fmt.Println("I am inside update.....2", iAmount, glbal.ContractAmount)
		//initializers.DB.Model(&glbal).Where("company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno).Update("contract_amount", iAmount)
		result = txn.Model(&glbal).Where("company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno).Update("contract_amount", iAmount)
		if result.Error != nil {
			txn.Rollback()
			return result.Error, glbal.ContractAmount
		}

		return nil, glbal.ContractAmount
	}
	//results.Commit()

}

// # 38
// GetMaxTranno - Get Transaction No and History Code
//
// Inputs: Company,  Policy No, Method, Effective Date, User
//
// # Outputs History Code and New Tranno
//
// # It update PHISTORY Table
//
// ©  FuturaInsTech
func GetMaxTrannoN(iCompany uint, iPolicy uint, iMethod string, iEffDate string, iuser uint64, historyMap map[string]interface{}, txn *gorm.DB) (error, string, uint) {
	var permission models.Permission
	var result *gorm.DB

	result = initializers.DB.First(&permission, "company_id = ? and method = ?", iCompany, iMethod)
	if result.Error != nil {
		return result.Error, iMethod, 0
	}
	iHistoryCode := permission.TransactionID
	var transaction models.Transaction
	result = initializers.DB.Find(&transaction, "ID = ?", iHistoryCode)
	if result.Error != nil {
		return result.Error, iMethod, 0
	}
	iHistoryCD := transaction.TranCode
	var phistory models.PHistory
	var maxtranno float64 = 0

	fmt.Println(iCompany, iPolicy, iHistoryCD, iEffDate)

	result1 := initializers.DB.Table("p_histories").Where("company_id = ? and policy_id= ?", iCompany, iPolicy).Select("max(tranno)")

	if result1.Error != nil {
		fmt.Println(result1.Error)

	}
	err := result1.Row().Scan(&maxtranno)
	fmt.Println("Error ", err)
	phistory.CompanyID = iCompany
	phistory.Tranno = uint(maxtranno) + 1
	phistory.PolicyID = iPolicy
	phistory.HistoryCode = iHistoryCD
	phistory.EffectiveDate = iEffDate
	phistory.Is_reversed = false
	phistory.IsValid = "1"
	if historyMap != nil {
		phistory.PrevData = historyMap
	}
	a := time.Now()
	b := Date2String(a)
	phistory.CurrentDate = b
	phistory.UpdatedID = iuser
	result1 = txn.Create(&phistory)
	if result1.Error != nil {
		return result1.Error, phistory.HistoryCode, phistory.Tranno
	}

	return nil, phistory.HistoryCode, phistory.Tranno

}

// # 161
// GetP0050ItemCodeDesc - Get the Description of an item's Code
//
// Inputs: Company, ParamItem and Language
//
// # Outputs:  Description
//
// ©  FuturaInsTech
func GetP0050ItemCodeDesc(iCompany uint, iItem string, iLanguage uint, iCode string) string {
	var paramdata models.Param
	var paramdatamap map[string]interface{}

	var idescription string = ""
	var iParam = "P0050"
	results := initializers.DB.Where("company_id = ? AND name = ? and item = ? and is_valid = ?", iCompany, iParam, iItem, 1).Find(&paramdata)
	if results.Error != nil || results.RowsAffected == 0 {
		return ""
	}

	datamap, _ := json.Marshal(paramdata.Data)
	json.Unmarshal(datamap, &paramdatamap)

	// use the iCode and return its corresponding Description
	type jsondata struct {
		Code        string
		Description string
	}
	var jd []jsondata

	datamap, _ = json.Marshal(paramdatamap["dataPairs"])
	json.Unmarshal(datamap, &jd)

	for i := 0; i < len(jd); i++ {
		if jd[i].Code == iCode {
			idescription = jd[i].Description
			break
		}
	}

	return idescription
}

// # 163
// Validate the PolicyFields mandatory as required by P0065 Rules
// Inputs: Company,
//
// # Outputs:
//
// ©  FuturaInsTech
func ValidatePolicyData(policyenq models.Policy, langid uint, iHistoryCode string) (string error) {
	businessdate := GetBusinessDate(policyenq.CompanyID, 0, 0)
	var clientenq models.Client
	result := initializers.DB.First(&clientenq, "company_id  = ? and id = ?", policyenq.CompanyID, policyenq.ClientID)
	if result.Error != nil {
		shortCode := "GL212" // Client Not Found
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)

	}

	var agencyenq models.Agency
	result = initializers.DB.First(&agencyenq, "company_id  = ? and id = ?", policyenq.CompanyID, policyenq.AgencyID)
	if result.Error != nil {
		shortCode := "GL275" // Agent Not Found
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	err := GetItemD(int(policyenq.CompanyID), "Q0005", policyenq.PProduct, policyenq.PRCD, &extradataq0005)
	if err != nil {
		shortCode := "GL385" // Q0005 not configured
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#001 RCD is less than PropsalDate
	if q0005data.BackDateAllowed == "N" {
		if policyenq.PRCD < policyenq.ProposalDate {
			shortCode := "GL539" // RCD is less than PropsalDate
			longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
			return errors.New(shortCode + ":" + longDesc)
		}
	}

	//#002 UW Date is less than PropsalDate
	if policyenq.PUWDate != "" {
		if policyenq.PUWDate < policyenq.ProposalDate {
			shortCode := "GL540" // UW Date is less than PropsalDate
			longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
			return errors.New(shortCode + ":" + longDesc)
		}
	}

	//#003 Frequency is Inalid
	var iFreqFound bool = false
	for i := 0; i < len(q0005data.Frequencies); i++ {
		if policyenq.PFreq == q0005data.Frequencies[i] {
			iFreqFound = true
			break
		}
	}
	if !iFreqFound {
		shortCode := "GL541" // Frequency is Inalid
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#004 Contract Curr is Inalid
	var iCCurrFound bool = false
	for i := 0; i < len(q0005data.ContractCurr); i++ {
		if policyenq.PContractCurr == q0005data.ContractCurr[i] {
			iCCurrFound = true
			break
		}
	}
	if !iCCurrFound {
		shortCode := "GL542" // Contract Curr is Inalid
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#005 Billing Curr is Inalid
	var iBCurrFound bool = false
	for i := 0; i < len(q0005data.ContractCurr); i++ {
		if policyenq.PBillCurr == q0005data.BillingCurr[i] {
			iBCurrFound = true
			break
		}
	}
	if !iBCurrFound {
		shortCode := "GL543" // Billing Curr is Inalid
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#006 Backdataing not Allowed
	if policyenq.PRCD < businessdate &&
		q0005data.BackDateAllowed == "N" {
		shortCode := "GL544" // Backdataing not Allowed
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#007 Agency Channel Not Allowed
	var iAgencyChannelFound bool = false
	for i := 0; i < len(q0005data.AgencyChannel); i++ {
		if agencyenq.AgencyChannel == q0005data.AgencyChannel[i] {
			iAgencyChannelFound = true
			break
		}
	}
	if !iAgencyChannelFound {
		shortCode := "GL545" // Agency Channel Not Allowed
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#008 Client is Invalid
	if clientenq.ClientStatus != "AC" {
		shortCode := "GL546" // Invalid Client
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#009 Deceased Client
	if !isFieldZero(clientenq.ClientDod) {
		shortCode := "GL547" // Deceased Client
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if policyenq.PRCD > businessdate {
		shortCode := "GL568" // RCD is greter than businessdate
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	return nil
}

// # 164
// Validate the PolicyFields mandatory as required by P0065 Rules
// Inputs: Company,
//
// # Outputs:
//
// ©  FuturaInsTech
func ValidateBenefitData(benefitenq models.Benefit, langid uint, iHistoryCode string) (string error) {
	//businessdate := GetBusinessDate(benefitenq.CompanyID, 0, 0)
	var clientenq models.Client
	result := initializers.DB.First(&clientenq, "company_id  = ? and id = ?", benefitenq.CompanyID, benefitenq.ClientID)
	if result.Error != nil {
		shortCode := "GL212" // Client Not Found
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	err := GetItemD(int(benefitenq.CompanyID), "Q0006", benefitenq.BCoverage, benefitenq.BStartDate, &extradataq0006)
	if err != nil {
		shortCode := "GL172" // Failed to Get Q0006
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#001 Age Not Allowed
	var iAllowedAge bool = false
	for i := 0; i < len(q0006data.AgeRange); i++ {
		if benefitenq.BAge == q0006data.AgeRange[i] {
			iAllowedAge = true
			break
		}
	}
	if !iAllowedAge {
		shortCode := "GL548" // Age Not Allowed
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#002 Policy Term not Allowed
	var iAllowedPolTerm bool = false
	for i := 0; i < len(q0006data.TermRange); i++ {
		if benefitenq.BTerm == q0006data.TermRange[i] {
			iAllowedPolTerm = true
			break
		}
	}
	if !iAllowedPolTerm {
		shortCode := "GL549" // Policy Term not Allowed
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#003 Premium Paying Term not Allowed
	var iAllowedPPT bool = false
	for i := 0; i < len(q0006data.PptRange); i++ {
		if benefitenq.BPTerm == q0006data.PptRange[i] {
			iAllowedPPT = true
			break
		}
	}
	if !iAllowedPPT {
		shortCode := "GL550" // Premium Paying Term not Allowed
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#004 Risk cess Age not Allowed
	benriskcessage := benefitenq.BAge + benefitenq.BTerm
	if benriskcessage < q0006data.MinRiskCessAge ||
		benriskcessage > q0006data.MaxRiskCessAge {
		shortCode := "GL551" // Risk cess Age not Allowed
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#005 Premium cess Age not Allowed
	benpremcessage := benefitenq.BAge + benefitenq.BPTerm
	if benpremcessage < q0006data.MinPremCessAge ||
		benpremcessage > q0006data.MaxPremCessAge {
		shortCode := "GL552" // Premium cess Age not Allowed
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	//#006 Min Sum Assured not met
	if uint(benefitenq.BSumAssured) < q0006data.MinSA {
		shortCode := "GL553" // Min Sum Assured not met
		longDesc, _ := GetErrorDesc(benefitenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	return nil
}

// # 165
// GetMaxTranno (New Version)
// Inputs: Company, Policy, Method, Effective Date, user
//
// # Outputs: uint
//
// ©  FuturaInsTech
func GetMaxTranno2(iCompany uint, iPolicy uint, iMethod string, iEffDate string, iuser uint64, txn *gorm.DB) (error, uint) {

	var maxtranno = 0
	var phistories models.PHistory

	result1 := txn.Order("tranno DESC").Find(&phistories, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	if result1.Error != nil {
		txn.Rollback()
		return result1.Error, 0
	}
	maxtranno = int(phistories.Tranno)
	return nil, uint(maxtranno)
}

// # 166
// Create Hisotry  Records (New Version)
// Inputs: Company, Policy, Method, Effective Date, Max Tranno, user,Hisotry Map
//
// # Outputs: error,
//
// ©  FuturaInsTech
func CreatePHistory(iCompany uint, iPolicy uint, iMethod string, iEffDate string, maxTranno uint, iuser uint64, historyMap map[string]interface{}, txn *gorm.DB) error {

	iHistoryCD := iMethod
	var phistory models.PHistory
	phistory.CompanyID = iCompany
	phistory.Tranno = maxTranno
	phistory.PolicyID = iPolicy
	phistory.HistoryCode = iHistoryCD
	phistory.EffectiveDate = iEffDate
	phistory.Is_reversed = false
	phistory.IsValid = "1"
	if historyMap != nil {
		phistory.PrevData = historyMap
	}
	a := time.Now()
	b := Date2String(a)
	phistory.CurrentDate = b
	phistory.UpdatedID = iuser
	result := txn.Create(&phistory)
	if result.Error != nil {
		txn.Rollback()
		return result.Error
	}
	return nil
}

// # 167
// Validate Nominee (New Version)
// Inputs: Nominee Model, Company id, User Language, History Code
//
// # Outputs: error,
//
// ©  FuturaInsTech
func ValidateNominee(nomineeval models.Nominee, userco uint, userlan uint, iKey string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return errors.New(err.Error())
	}
	for i := 0; i < len(p0065data.FieldList); i++ {

		var fv interface{}
		r := reflect.ValueOf(nomineeval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}
	}
	var clientenq models.Client
	result := initializers.DB.First(&clientenq, "company_id  = ? and id = ?", nomineeval.CompanyID, nomineeval.ClientID)
	if result.Error != nil {
		shortCode := "GL212" // Client Not Found
		longDesc, _ := GetErrorDesc(nomineeval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if clientenq.ClientStatus != "AC" ||
		clientenq.ClientDod != "" {
		shortCode := "GL546" // Invalid Client
		longDesc, _ := GetErrorDesc(nomineeval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	var p0045data paramTypes.P0045Data
	var extradatap0045 paramTypes.Extradata = &p0045data
	err = GetItemD(int(nomineeval.CompanyID), "P0045", nomineeval.NomineeRelationship, "0", &extradatap0045)
	if err != nil {
		shortCode := "GL573" // P0045 not configured
		longDesc, _ := GetErrorDesc(nomineeval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	var iGender bool = false
	for i := 0; i < len(p0045data.Gender); i++ {
		if clientenq.Gender == p0045data.Gender {
			iGender = true
			break
		}
	}
	if !iGender {
		shortCode := "GL572" // gender is not same in relationship
		longDesc, _ := GetErrorDesc(nomineeval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	// Owner cannot be Nominee
	var policyenq models.Policy
	result = initializers.DB.First(&policyenq, "company_id  = ? and id = ?", nomineeval.CompanyID, nomineeval.PolicyID)
	if result.Error != nil {
		shortCode := "GL210" // Policy Not Found
		longDesc, _ := GetErrorDesc(nomineeval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if nomineeval.ClientID == policyenq.ClientID {
		shortCode := "GL589" // Owner cannot be Nominee
		longDesc, _ := GetErrorDesc(nomineeval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	return
}

// # 168
// Validate Payer (New Version)
// Inputs: Payer Model, Company id, User Language, History Code
//
// # Outputs: error
//
// ©  FuturaInsTech

func ValidatePayer(payerval models.Payer, userco uint, userlan uint, iKey string) (string error) {
	businessdate := GetBusinessDate(payerval.CompanyID, 0, 0)
	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return errors.New(err.Error())
	}
	for i := 0; i < len(p0065data.FieldList); i++ {

		var fv interface{}
		r := reflect.ValueOf(payerval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}
	}

	iPolicy := payerval.PolicyID
	var policy models.Policy
	result := initializers.DB.Find(&policy, "company_id = ? and id = ?", userco, iPolicy)
	if result.Error != nil {
		shortCode := "GL175"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)

	}
	if payerval.FromDate > businessdate {
		shortCode := "GL616" // From date is greater than business date
		longDesc, _ := GetErrorDesc(payerval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}
	if payerval.FromDate < policy.PRCD {
		shortCode := "GL617" // From Date is lesser than RCD Date
		longDesc, _ := GetErrorDesc(payerval.CompanyID, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if payerval.FromDate > payerval.ToDate {
		shortCode := "GL901" // FromDate greater than ToDate
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)
	}

	return
}

// # 171
// Validate Frequency (New Version)
// Return False when current frequency premium dues are pending, else true
// Inputs: RCD, PTD, Curr Freq, New Freq
//
// # Outputs: True or False
//
// ©  FuturaInsTech
func ValidateFreq(PrcdDate string, Ptdate string, Currfreq string, Newfreq string) bool {
	tdate := PrcdDate
	for {
		DueDate := GetNextDue(tdate, Newfreq, "")
		NextDueDate := Date2String(DueDate)
		if NextDueDate == Ptdate {
			return true
		} else if NextDueDate > Ptdate {
			return false
		}
		tdate = NextDueDate
	}
}

// #172
// Validate Agency (New Version)
// Return Error when Agency is not valid
// Inputs: agency model, user company, user language, validdate
//
// # Outputs: nil or error
//
// ©  FuturaInsTech
func ValidateAgency(agencyenq models.Agency, userco uint, userlan uint, iDate string) (string error) {

	if agencyenq.AgencySt != "AC" {
		shortCode := "GL221" // InValid Status
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if agencyenq.LicenseStartDate > iDate {
		shortCode := "GL577"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if agencyenq.LicenseEndDate < iDate {
		shortCode := "GL578"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}
	return nil
}

// # 173
// Validate the PolicyData and Benefits Data as required by Q0011 Rules
// Inputs: Policy Data, Benefit(s) Data, Lang Id,
//
// # Outputs: error
//
// ©  FuturaInsTech
func ValidatePolicyBenefitsData(policyenq models.Policy, benefitenq []models.Benefit, langid uint) (string error) {

	var q0011data paramTypes.Q0011Data
	var extradataq0011 paramTypes.Extradata = &q0011data
	err := GetItemD(int(policyenq.CompanyID), "Q0011", policyenq.PProduct, policyenq.PRCD, &extradataq0011)
	if err != nil {
		shortCode := "GL387" // Q0011 not configured
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	// Duplicate Benefits Check
	duplicatebenefits := false
	var benefitenq1 []models.Benefit = benefitenq
	for i := 0; i < len(benefitenq); i++ {
		duplicatebenefits = false
		for j := 0; j < len(benefitenq1); j++ {
			if benefitenq[i].ID != benefitenq1[j].ID &&
				benefitenq[i].BCoverage == benefitenq1[j].BCoverage &&
				benefitenq[i].ClientID == benefitenq1[j].ClientID {
				duplicatebenefits = true
				break
			}
		}
	}
	if duplicatebenefits {
		shortCode := "GL619" // Duplicate Benefits Exist
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	// basicbenefit selection
	var basicbenefit models.Benefit
	basicbenefit.ID = 0
	for i := 0; i < len(q0011data.Coverages); i++ {
		for j := 0; j < len(benefitenq); j++ {
			if q0011data.Coverages[i].CoverageName == benefitenq[j].BCoverage &&
				q0011data.Coverages[i].BasicorRider == "B" {
				basicbenefit = benefitenq[j]
				break
			}
		}
	}

	// Mandatory Benefits check
	mandatorybenefits := true
	for i := 0; i < len(q0011data.Coverages); i++ {
		if q0011data.Coverages[i].Mandatory == "Y" {
			mandatorybenefits = false
			for j := 0; j < len(benefitenq); j++ {
				if q0011data.Coverages[i].CoverageName == benefitenq[j].BCoverage {
					mandatorybenefits = true
					break
				}
			}
		}
	}

	if !mandatorybenefits {
		shortCode := "GL621" // Mandatory Coverage(s) not Found
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	// Benefits Validation
	for i := 0; i < len(benefitenq); i++ {
		//#001 Policy RCD > Benefit Start Date
		if policyenq.PRCD > benefitenq[i].BStartDate {
			shortCode := "GL622" // Policy RCD > Benefit Start Date
			longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
			return errors.New(shortCode + ":" + longDesc)
		}
		if benefitenq[i].BCoverage != basicbenefit.BCoverage {
			for j := 0; j < len(q0011data.Coverages); j++ {
				if benefitenq[i].BCoverage == q0011data.Coverages[j].CoverageName {
					//#002 Benefit Term Exceed Basic Benefit Term
					if q0011data.Coverages[i].TermCanExceed == "N" &&
						benefitenq[i].BTerm > basicbenefit.BTerm {
						shortCode := "GL623" // Benefit Term Exceed Basic Benefit Term
						longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
						return errors.New(shortCode + ":" + longDesc)
					}
					//#003 Benefit Prem Term Exceed Basic Benefit Prem Term
					if q0011data.Coverages[i].PptCanExceed == "N" &&
						benefitenq[i].BPTerm > basicbenefit.BPTerm {
						shortCode := "GL624" // Benefit Prem Term Exceed Basic Benefit Prem Term
						longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
						return errors.New(shortCode + ":" + longDesc)
					}
					//#004 Benefit SA Exceed Basic Benefit SA
					if q0011data.Coverages[i].SaCanExceed == "N" &&
						benefitenq[i].BSumAssured > basicbenefit.BSumAssured {
						shortCode := "GL625" // Benefit SA Exceed Basic Benefit SA
						longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
						return errors.New(shortCode + ":" + longDesc)
					}
				}
			}
		}
	}

	return nil
}

// #175
// GetCurrencyNamebyId
// Inputs: curr_id
//
// # Outputs  Currency ShortName and Currency LongName
//
// ©  FuturaInsTech
func GetCurrencyName(iCurr uint) (string, string) {
	var curry models.Currency
	result := initializers.DB.Find(&curry, "id = ?", iCurr)
	if result.Error != nil {
		return "", ""
	}
	return curry.CurrencyShortName, curry.CurrencyLongName
}

// #176
// GetCoCurrIdName
// Inputs: Company Id
//
// # Outputs  Currency ShortName and Currency LongName
//
// ©  FuturaInsTech
func GetCoCurrIdName(iCompany uint) (string, string) {
	var company models.Company
	result := initializers.DB.Find(&company, "id = ?", iCompany)
	if result.Error != nil {
		return "", ""
	}
	var curry models.Currency
	result = initializers.DB.Find(&curry, "id = ?", company.CurrencyID)
	if result.Error != nil {
		return "", ""
	}
	return curry.CurrencyShortName, curry.CurrencyLongName
}

// #177
// Intf2Uint
// Inputs: Interface Variable Value
//
// # Outputs  Uint variable Value
//
// ©  FuturaInsTech
func Intf2Uint(i interface{}) uint {

	value, ok := i.(float64)
	if !ok {
		return 0
	}
	return uint(value)
}

// #178
// Intf2Int
// Inputs: Interface Variable Value
//
// # Outputs  Int variable Value
//
// ©  FuturaInsTech
func Intf2Int(i interface{}) int {

	value, ok := i.(float64)
	if !ok {
		return 0
	}
	return int(value)
}

// #179
// Intf2String
// Inputs: Interface Variable Value
//
// # Outputs  String variable Value
//
// ©  FuturaInsTech
func Intf2String(i interface{}) string {

	value, ok := i.(string)
	if !ok {
		j := Intf2Int(i)
		return strconv.Itoa(j)
	}
	return value
}

// #180
// SplitDateString
// Inputs: String Date Value in YYYYMMDD Format
//
// # Outputs  Year, Month and Date
//
// ©  FuturaInsTech
func SplitDateString(ds string) (string, string, string) {
	// Check if the input string has at least 8 characters
	if len(ds) < 8 {
		return "", "", ""
	}

	// Extract year, month, and day using string slicing
	year := ds[:4]
	month := ds[4:6]
	day := ds[6:8]

	return year, month, day
}

// #181
// Validate BillType(New Version)
// Return Error when BillType is not valid
// Inputs: billtype string,payingAuthority uint, user company, user language,
//
// # Outputs: nil or error
//
// # Only SSI validation implemented
//
// ©  FuturaInsTech
func ValidateBillType(policyenq models.Policy, userco uint, userlan uint, iDate string, iBillType string, iPayingAuthority uint) (string error) {

	var p0055data paramTypes.P0055Data
	var extradatap0055 paramTypes.Extradata = &p0055data

	err := GetItemD(int(userco), "P0055", iBillType, iDate, &extradatap0055)

	if err != nil {
		shortCode := "GL279"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}
	// Validate SSI Bill Type

	if p0055data.PayingAuthority == "N" &&
		iBillType == policyenq.BillingType {
		shortCode := "GL637" // Existing and new bill type shuld not be the same
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if p0055data.PayingAuthority == "Y" &&
		iBillType == policyenq.BillingType &&
		iPayingAuthority == policyenq.PayingAuthority {

		shortCode := "GL638" // existing  and new Pa should not be same
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if p0055data.PayingAuthority == "N" {
		if iPayingAuthority != 0 {
			shortCode := "GL700" // Paying Authority Is Not Provided
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + ":" + longDesc)

		}
	}

	if p0055data.PayingAuthority == "Y" {
		if iPayingAuthority == 0 {
			shortCode := "GL701" // Paying Authority Not Required
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + ":" + longDesc)

		}
	}

	// validate Paying authority
	err = ValidatePayingAuthority(userco, userlan, iDate, iPayingAuthority)
	if err != nil {
		shortCode := "GL639" // No item found in validate Paying Authority
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	// P0055 Bank Extration Types like cBank,DBank,NEFT,UPI validation are to be added

	return nil
}

// #182
// Validate Paying Authority(New Version)
// Return Error when Paying Authority is not valid
// Inputs: payingAuthority uint, user company, user language,
//
// # Outputs: nil or error
//
// ©  FuturaInsTech
func ValidatePayingAuthority(userco uint, userlan uint, iDate string, iPayingAuthority uint) (string error) {

	var payingauth models.PayingAuthority
	result := initializers.DB.First(&payingauth, "company_id = ? and id = ?", userco, iPayingAuthority)
	if result.Error != nil {
		shortCode := "GL671" //Failed to get Paying Authority
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if payingauth.PaStatus != "AC" {
		shortCode := "GL640" // InValid PA Status
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if payingauth.StartDate > iDate {
		shortCode := "GL641" // PA Start Date Should be Greater than Current Date
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if payingauth.EndDate < iDate {
		shortCode := "GL642" // PA End Date Should be Greater than Current Date
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	return nil
}

// #184
// Get Validate Client Work
// Inputs: ClientWork data, Company, Language, Date and Key Detail
//
// # Outputs  string error
//
// ©  FuturaInsTech
func ValidateClientWork(clientwork models.ClientWork, userco uint, userlan uint, iDate string, iKey string) (string error) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return errors.New(err.Error())
	}

	for i := 0; i < len(p0065data.FieldList); i++ {

		var fv interface{}
		r := reflect.ValueOf(clientwork)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)
		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	var client models.Client
	clientid := clientwork.ClientID
	initializers.DB.Find(&client, "company_id = ? and id = ?", userco, clientid)

	if client.ClientStatus != "AC" {
		shortCode := "GL221" // InValid Status
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}
	var employer models.Client
	employerid := clientwork.EmployerID
	initializers.DB.Find(&employer, "company_id = ? and id = ?", userco, employerid)

	if employer.ClientStatus != "AC" {
		shortCode := "GL221" // InValid Status
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if clientwork.StartDate > iDate {
		shortCode := "GL656"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}

	if clientwork.EndDate < iDate {
		shortCode := "GL657"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + ":" + longDesc)
	}
	return nil
}

// #185
// Stamp Duty for Loan Calculation
// Inputs: CompanyID, SD Rate, Loan Amt, PolicyID
//
// # Outputs  Stamp Duty Amount
//
// ©  FuturaInsTech
func CalculateStampDutyforLoan(iCompany uint, iRate float64, iDate string, iLoanAmount float64, iPolicy uint) float64 {

	oStampDuty := iRate * iLoanAmount
	oStampDuty = RoundFloat(oStampDuty, 2)

	return oStampDuty
}

// #187
// All OS Loan and OS Loan Interest by Loan Type
// Inputs: CompanyID, PolicyID, BenefitID, Effective Date, Loan Type
//
// # Outputs  OS Loan, OS Loan Int, Total Loan Int, Total Loan, Loan Currency
//
// ©  FuturaInsTech

func GetAllLoanOSByType(iCompany uint, iPolicy uint, iBenefit uint, iEffectiveDate string) (oLoanOS float64, oLoanIntOS float64, oLoanInt float64, oTotalLoan float64, oLoanCurr string) {

	var loanenq []models.Loan
	result := initializers.DB.Find(&loanenq, "company_id = ? and policy_id = ? and benefit_id = ? and loan_status = ? ", iCompany, iPolicy, iBenefit, "AC")
	if result.Error != nil {
		// txn.Rollback()
		return 0, 0, 0, 0, ""
	}
	var p0072data paramTypes.P0072Data
	var extradata3 paramTypes.Extradata = &p0072data
	GetItemD(int(iCompany), "P0072", "LN001", iEffectiveDate, &extradata3)

	var itemp float64
	var LoanIntOS float64
	var loanIntPaid float64
	var iAmount float64
	var brokenperiodint float64

	for i := 0; i < len(loanenq); i++ {

		_, _, _, iNoOfDays, _, _, _, _ := NoOfDays(iEffectiveDate, loanenq[i].LastIntBillDate)
		oLoanInt = p0072data.RateOfInterest
		LoanIntamt := loanenq[i].LastCapAmount

		if p0072data.LoanInterestType == "C" {
			iAmount = CompoundInterest(LoanIntamt, oLoanInt, float64(iNoOfDays))
		} else if p0072data.LoanInterestType == "S" {
			iAmount = SimpleInterest(oLoanOS, oLoanInt, float64(iNoOfDays))
		}
		brokenperiodint = RoundFloat(iAmount, 2)

		if loanenq[i].NextCapDate > iEffectiveDate {
			oTotalLoan = oLoanOS
			oLoanOS = oLoanOS + loanenq[i].LastCapAmount
			oLoanInt = loanenq[i].LoanIntRate
			_, _, _, iNoOfDays, _, _, _, _ := NoOfDays(iEffectiveDate, loanenq[i].LastCapDate)
			itemp := CompoundInterest(oLoanOS, oLoanInt, float64(iNoOfDays))
			oLoanIntOS = itemp
			oLoanOS = oLoanOS + brokenperiodint
			// oLoanIntOS = itemp + loanenq[i].LastCapAmount - oLoanIntOS
			// oTotalLoan = oLoanOS
			oLoanCurr = loanenq[i].LoanCurrency
		}
		oTotalLoan = oLoanOS
	}
	var loanbillupd []models.LoanBill
	initializers.DB.Find(&loanbillupd, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	for i := 0; i < len(loanbillupd); i++ {

		if loanbillupd[i].ReceiptNo == 0 {

			itemp = loanbillupd[i].LoanIntAmount

			LoanIntOS += RoundFloat(itemp, 2)

		}

		if loanbillupd[i].ReceiptNo != 0 {

			itemp = loanbillupd[i].LoanIntAmount

			loanIntPaid += RoundFloat(itemp, 2)

		}

	}

	oLoanOS = oLoanOS + LoanIntOS - loanIntPaid
	return oLoanOS, oLoanIntOS, oLoanInt, oLoanOS, oLoanCurr

}

// #189
// Get Next Loan Number (Increment and return to use for new Loan)
// Inputs: CompanyID, PolicyID
//
// # Outputs  Next Loan Sequence Number (Maximum Loan Number + 1)
//
// ©  FuturaInsTech
func GetMaxLoanSeqNo(iCompany uint, iPolicy uint) (error, uint) {

	var result *gorm.DB

	result1 := initializers.DB.Table("loans").Where("company_id = ? and policy_id= ?", iCompany, iPolicy).Select("max(loan_seq_number)")

	if result1.Error != nil {
		return result.Error, 0
	}

	var loanseqno uint = 0
	var newloanseqno uint = 0

	err := result1.Row().Scan(&loanseqno)

	fmt.Println("Error ", err)

	newloanseqno = uint(loanseqno) + 1

	return nil, newloanseqno

}

// #190
// TDF Loan Interest Billing Process
// Inputs: CompanyID, PolicyID, Function, Tran No,
//
// # Outputs  Error
//
// ©  FuturaInsTech
func TDFLoanInt(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var loanenq []models.Loan
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_status = ? ", iCompany, iPolicy, "AC")
	loandelete := "N"
	if result.Error != nil || result.RowsAffected == 0 {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			return "", result.Error
		} else {
			loandelete = "Y"
		}
	}

	if loandelete == "Y" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error == nil {
			txn.Delete(&tdfpolicy)
			return "", nil
		}

	}
	oDate := ""
	for i := 0; i < len(loanenq); i++ {

		if oDate == "" {
			oDate = loanenq[i].NextIntBillDate
		}
		if loanenq[i].NextIntBillDate < oDate {
			oDate = loanenq[i].NextIntBillDate
		}
	}

	if oDate != "" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			txn.Create(&tdfpolicy)
			return "", nil
		} else {
			txn.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			txn.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

// #191
// TDF Loan Interest Capitalization Process
// Inputs: CompanyID, PolicyID, Function, Tran No,
//
// # Outputs  Error
//
// ©  FuturaInsTech
func TDFLoanCap(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var loanenq []models.Loan
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_status = ? ", iCompany, iPolicy, "AC")
	loandelete := "N"
	if result.Error != nil || result.RowsAffected == 0 {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			return "", result.Error
		} else {
			loandelete = "Y"
		}
	}

	if loandelete == "Y" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error == nil {
			txn.Delete(&tdfpolicy)
			return "", nil
		}

	}
	oDate := ""
	for i := 0; i < len(loanenq); i++ {

		if oDate == "" {
			oDate = loanenq[i].NextCapDate
		}
		if loanenq[i].NextCapDate < oDate {
			oDate = loanenq[i].NextCapDate
		}
	}

	if oDate != "" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno
			txn.Create(&tdfpolicy)
			return "", nil
		} else {
			txn.Delete(&tdfpolicy)
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			txn.Create(&tdfpolicy)
			return "", nil
		}
	}
	return "", nil
}

// #195
// TDF Loan Deposit Adjustment process
// Inputs: CompanyID, PolicyID, Function, Tran No,
//
// # Outputs  Error
//
// ©  FuturaInsTech
func TDFLoanDN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iDate string, txn *gorm.DB) (string, error) {
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = iDate
		tdfpolicy.Tranno = iTranno
		tdfpolicy.Seqno = tdfrule.Seqno
		txn.Create(&tdfpolicy)
		return "", nil
	} else {
		txn.Delete(&tdfpolicy)
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = iDate
		tdfpolicy.Tranno = iTranno

		txn.Create(&tdfpolicy)
		return "", nil
	}
}

// #197
// Get IlpMortality and IlpFee values from the benefit
// Inputs: CompanyID, iBenefit
//
// # Outputs: IlpMortality,IlpFee,totUnpaidInterest
//
// ©  FuturaInsTech
func GetIlpMortalityFee(iCompany uint, iBenefit uint) (oIlpMortality float64, oIlpFee float64) {

	var benefitenq []models.Benefit

	result := initializers.DB.Find(&benefitenq, "id= ? and company_id = ? ", iBenefit, iCompany)
	if result.Error != nil {
		return
	}

	for i := 0; i < len(benefitenq); i++ {

		oIlpMortality += benefitenq[i].IlpMortality
		oIlpFee += benefitenq[i].IlpFee
	}

	return oIlpMortality, oIlpFee
}

// #198
// CheckNegativeFund from the fund value
// Inputs: CompanyID, iPolicy, iEffectiveDate, iHistoryCode
//
// # Outputs: oAmount,oStatus, ONegativeUnitsOrAmt
//
// ©  FuturaInsTech
func CheckNegativeFund(iCompany uint, iPolicy uint, iEffectiveDate string, iHistoryCode string, txn *gorm.DB) (oAmount float64, oStatus string, ONegativeUnitsOrAmt string) {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return 0, "N", ""
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	err := GetItemD(int(iCompany), "Q0005", policyenq.PProduct, policyenq.PRCD, &extradataq0005)

	if err != nil {
		return 0, "N", ""
	}

	if q0005data.ProductFamily != "RUL" {
		return 0, "N", ""
	}

	iFund, _, _ := GetAllFundValueByPol(iCompany, iPolicy, iEffectiveDate)

	var q0011data paramTypes.Q0011Data
	var extradata1 paramTypes.Extradata = &q0011data
	iProduct := policyenq.PProduct
	iDate := policyenq.PRCD

	err = GetItemD(int(iCompany), "Q0011", iProduct, iDate, &extradata1)
	if err != nil {
		return 0, "N", ""
	}
	var iBasicCover string
	for i := 0; i < len(q0011data.Coverages); i++ {
		if q0011data.Coverages[i].BasicorRider == "B" {
			iBasicCover = q0011data.Coverages[i].CoverageName
			break
		}
	}

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBasicCover
	err = GetItemD(int(iCompany), "P0059", iKey, iDate, &extradatap0059)
	if err != nil {
		return 0, "N", ""
	}

	if iFund == 0 {
		if p0059data.NegativeAccum == "Y" {
			return iFund, "Y", p0059data.NegativeUnitsOrAmt
		}
	}

	return iFund, "N", p0059data.NegativeUnitsOrAmt

}

// #199
// GetP0059  from the param
// Inputs: CompanyID, iPolicy, iEffectiveDate, iHistoryCode
//
// # Outputs: oNegativeAmount, oNegativeMonths
//
// ©  FuturaInsTech
func GetP0059(iCompany uint, iPolicy uint, iEffectiveDate string, iHistoryCode string, txn *gorm.DB) (oNegativeAmount string, oNegativeMonths float64) {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", 0
	}

	var q0011data paramTypes.Q0011Data
	var extradata1 paramTypes.Extradata = &q0011data
	iProduct := policyenq.PProduct
	iDate := policyenq.PRCD

	err := GetItemD(int(iCompany), "Q0011", iProduct, iDate, &extradata1)
	if err != nil {
		return "", 0
	}
	var iBasicCover string
	for i := 0; i < len(q0011data.Coverages); i++ {
		if q0011data.Coverages[i].BasicorRider == "B" {
			iBasicCover = q0011data.Coverages[i].CoverageName
			break
		}
	}

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBasicCover
	err = GetItemD(int(iCompany), "P0059", iKey, iDate, &extradatap0059)
	if err != nil {
		return "", 0
	}

	return p0059data.NegativeUnitsOrAmt, p0059data.NegativeAccumMonths

}

// #200
// GetP0059  from the param
// Inputs: CompanyID, iPolicy, iEffectiveDate
//
// # Outputs: p0069Month, NoOfMonths
// ©  FuturaInsTech
func GetP0069Data(iCompany uint, iPolicy uint, iDate string) (p0069Month int, NoOfMonths int) {

	txn := initializers.DB.Begin()

	var benefitenq []models.Benefit
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		txn.Rollback()
		return 0, 0
	}

	var p0069data1 paramTypes.P0069Data
	var extradatap00691 paramTypes.Extradata = &p0069data1
	// var p0069Month int
	// var NoOfMonths int
	for i := 0; i < len(benefitenq); i++ {
		iKey := benefitenq[i].BCoverage
		iBenefitDate := benefitenq[i].BStartDate

		err := GetItemD(int(iCompany), "P0069", iKey, iBenefitDate, &extradatap00691)
		if err != nil {
			txn.Rollback()
			return 0, 0
		}

		NoOfMonths = NewNoOfInstalments(benefitenq[i].BStartDate, iDate)
		for e := 0; e < len(p0069data1.P0069); e++ {
			if p0069data1.P0069[e].LiquidatedIlpFund == "Y" {
				p0069Month = int(p0069data1.P0069[e].Months)

			}
		}
	}
	return p0069Month, NoOfMonths
}

// #201
// GetLoanIntrest  from the loan table
// Inputs: CompanyID, iPolicy, iEffectiveDate,loantype
//
// # Outputs: LoanAmount, oLoanInt
// ©  FuturaInsTech
func GetLoanAndIntrestD(iCompany uint, iPolicy uint, iDate string, itype string, txn *gorm.DB) (oLoanAmount float64, oLoanInt float64) {

	var policy models.Policy
	result := initializers.DB.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return 0, 0
	}

	var loanenq []models.Loan
	var loanbillenq []models.LoanBill

	result = txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_status =? and loan_type=?", iCompany, iPolicy, "AC", itype)
	if result.Error != nil {
		return 0, 0
	}
	result = txn.Find(&loanbillenq, "company_id = ? and policy_id = ? and loan_type=? and billing_status=?", iCompany, iPolicy, itype, "OP")
	if result.Error != nil {
		return 0, 0
	}

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data

	var p0072data paramTypes.P0072Data
	var extradata3 paramTypes.Extradata = &p0072data
	iAmount := 0.0
	brokenperiodint := 0.0
	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)
	oLoanAmount = 0.0
	for i := 0; i < len(loanenq); i++ {
		GetItemD(int(iCompany), "Q0006", loanenq[i].BCoverage, iDate, &extradataq0006)
		if itype == "A" {
			GetItemD(int(iCompany), "P0072", q0005data.AplLoanMethod, iDate, &extradata3)
		} else if itype == "P" {
			GetItemD(int(iCompany), "P0072", q0006data.LoanMethod, iDate, &extradata3)
		}
		oLoanAmount += loanenq[i].LastCapAmount

		if iDate > loanenq[i].LastIntBillDate {
			_, _, _, iNoOfDays, _, _, _, _ := NoOfDays(iDate, loanenq[i].LastIntBillDate)
			oLoanIntRat := p0072data.RateOfInterest
			LoanIntamt := loanenq[i].LastCapAmount

			if p0072data.LoanInterestType == "C" {
				iAmount = CompoundInterest(LoanIntamt, oLoanIntRat, float64(iNoOfDays))
			} else if p0072data.LoanInterestType == "S" {
				iAmount = SimpleInterest(LoanIntamt, oLoanIntRat, float64(iNoOfDays))
			}
			brokenperiodint += iAmount
		}

	}

	for i := 0; i < len(loanbillenq); i++ {
		oLoanInt += loanbillenq[i].LoanIntAmount
	}

	oLoanInt = oLoanInt + brokenperiodint
	return RoundFloat(oLoanAmount, 2), RoundFloat(oLoanInt, 2)
}

// #202
// xlsx2json - xlsxfile data is transformed into jsonfile
// Inputs: xlsxFile
//
// # Outputs: jsonFile
// ©  FuturaInsTech
func xlsx2json(excelFileName string, jsonFileName string) {

	// Open the Excel file
	f, err := excelize.OpenFile(excelFileName)
	if err != nil {
		log.Fatalf("Error opening Excel file: %v", err)
	}

	// Get all sheet names
	sheetNames := f.GetSheetList()

	// Create a map to hold the data for all sheets
	allSheetsData := make(map[string][]map[string]interface{})

	// Process each sheet
	for _, sheetName := range sheetNames {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("Error reading rows for sheet %s: %v", sheetName, err)
			continue
		}

		if len(rows) == 0 {
			log.Fatal("Excel file is empty")
		}

		// Extract headers from the first row
		headers := rows[0]

		// Create a slice to hold the JSON objects for this sheet
		var sheetData []map[string]interface{}

		// Process rows to create JSON objects
		if len(rows) > 1 {
			// Process rows starting from the second row if they exist
			for _, row := range rows[1:] {
				obj := make(map[string]interface{})
				for colIndex, value := range row {
					if colIndex < len(headers) {
						obj[headers[colIndex]] = value // Use headers to map keys
					}
				}
				sheetData = append(sheetData, obj)
			}
		} else {
			// No data rows; write an empty object with only headers
			obj := make(map[string]interface{})
			for _, header := range headers {
				obj[header] = nil // Set values as `nil` for missing data
			}
			sheetData = append(sheetData, obj)
		}

		// Add the sheet data to the map
		allSheetsData[sheetName] = sheetData
	}

	// Write the JSON data to a file
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		log.Fatalf("Error creating JSON file: %v", err)
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ") // Format the JSON output with indentation
	if err := encoder.Encode(allSheetsData); err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}

	fmt.Println("JSON file successfully created:", jsonFileName)
}

// #203
// json2xlsx - jsonfile data is transformed into xlsxfile
// Inputs: jsonFile
//
// # Outputs: xlsxFile
// ©  FuturaInsTech
func json2xlsx(jsonFile string, excelFile string) {
	// Read JSON file
	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var jsondata1 map[string][]map[string]string // nested data structure
	var jsondata2 []map[string]interface{}       // slice of map string interface
	var jsondata3 []map[string]string            // slice of map string of strings
	var jsondata4 map[string]interface{}         // map string interface
	var jsondata5 map[string]string              // map string of string
	var jsondata6 []string                       // array of string

	jsonType := ""

	err = json.Unmarshal(jsonData, &jsondata1)
	if err == nil {
		if jsonType == "" {
			jsonType = "1"
		}
	}

	err = json.Unmarshal(jsonData, &jsondata2)
	if err == nil {
		if jsonType == "" {
			jsonType = "2"
		}
	}

	err = json.Unmarshal(jsonData, &jsondata3)
	if err == nil {
		if jsonType == "" {
			jsonType = "3"
		}
	}

	err = json.Unmarshal(jsonData, &jsondata4)
	if err == nil {
		if jsonType == "" {
			jsonType = "4"
		}
	}

	err = json.Unmarshal(jsonData, &jsondata5)
	if err == nil {
		if jsonType == "" {
			jsonType = "5"
		}
	}

	err = json.Unmarshal(jsonData, &jsondata6)
	if err == nil {
		if jsonType == "" {
			jsonType = "6"
		}
	}

	switch jsonType {
	case "1": // Parse JSON data into a nested data structure
		var sheetsData map[string][]map[string]string
		err = json.Unmarshal(jsonData, &sheetsData)
		if err != nil {
			log.Fatalf("Failed to parse JSON data: %v", err)
		}

		// Create a new Excel file
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("Failed to close Excel file: %v", err)
			}
		}()

		// Populate the Excel file with data for each sheet
		for sheetName, records := range sheetsData {
			// Add a new sheet
			index, _ := f.NewSheet(sheetName)
			if index == -1 {
				log.Fatalf("Failed to create sheet: %s", sheetName)
			}

			if len(records) == 0 {
				log.Printf("Skipping empty sheet: %s", sheetName)
				continue
			}

			// Create a map to store unique keys
			uniqueKeys := make(map[string]struct{})

			// Collect all keys from each JSON object
			for _, obj := range records {
				for key := range obj {
					uniqueKeys[key] = struct{}{} // Using empty struct for memory efficiency
				}
			}

			// Convert map keys to a slice
			headers := make([]string, 0, len(uniqueKeys))
			for key := range uniqueKeys {
				headers = append(headers, key)
			}

			//Write header rows
			for i, header := range headers {
				cell, _ := excelize.CoordinatesToCellName(i+1, 1)
				if err := f.SetCellValue(sheetName, cell, header); err != nil {
					log.Fatalf("Failed to write header to sheet %s: %v", sheetName, err)
				}
			}

			// Write data rows
			for rowIdx, record := range records {
				for colIdx, header := range headers {
					cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
					if err := f.SetCellValue(sheetName, cell, record[header]); err != nil {
						log.Fatalf("Failed to write data to sheet %s: %v", sheetName, err)
					}
				}
			}
		}
		// Delete the default "Sheet1"
		defaultSheet := "Sheet1"
		if err := f.DeleteSheet(defaultSheet); err != nil {
			log.Fatalf("Failed to delete default sheet: %v", err)
		}
		// Save the Excel file
		if err := f.SaveAs(excelFile); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		} else {
			fmt.Printf("Excel file has been created: %s\n", excelFile)
		}

	case "2": // Parse JSON data into slice of map string interface
		var data []map[string]interface{}
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			log.Fatalf("Failed to parse JSON data: %v", err)
		}

		// Create a new Excel file
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("Failed to close Excel file: %v", err)
			}
		}()

		sheetName := "Sheet1"
		index, _ := f.NewSheet(sheetName)

		// Find all unique keys from the JSON data to create the headers
		headers := make(map[string]bool) // To track unique keys
		for _, item := range data {
			for key := range item {
				headers[key] = true
			}
		}

		// Write the column headers in the first row
		col := 1
		for key := range headers {
			cellHeader, _ := excelize.CoordinatesToCellName(col, 1) // Column (col), Row 1
			f.SetCellValue(sheetName, cellHeader, key)
			col++
		}

		// Write the values under each header
		row := 2 // Start from the second row
		for _, item := range data {
			col = 1
			for key := range headers {
				// Write the value in the correct column
				cellValue, _ := excelize.CoordinatesToCellName(col, row) // Column (col), Row (row)
				if value, exists := item[key]; exists {
					f.SetCellValue(sheetName, cellValue, value)
				}
				col++
			}
			row++
		}

		// Set the active sheet
		f.SetActiveSheet(index)

		// Save the Excel file
		if err := f.SaveAs(excelFile); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		} else {
			fmt.Printf("Excel file has been created: %s\n", excelFile)
		}

	case "3": // Parse JSON data into slice of map string of strings
		var data []map[string]string
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			log.Fatalf("Failed to parse JSON data: %v", err)
		}

		// Collect all unique headers
		headerSet := make(map[string]bool)
		for _, record := range data {
			for key := range record {
				headerSet[key] = true
			}
		}

		// Convert headers to a sorted slice
		headers := make([]string, 0, len(headerSet))
		for header := range headerSet {
			headers = append(headers, header)
		}
		sort.Strings(headers)

		// Create a new Excel file
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("Failed to close Excel file: %v", err)
			}
		}()

		sheetName := "Sheet1"
		f.SetSheetName(f.GetSheetName(0), sheetName)

		// Write headers to the first row
		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// Write data to subsequent rows
		for rowNum, record := range data {
			for colNum, header := range headers {
				value, exists := record[header]
				if exists {
					cell, _ := excelize.CoordinatesToCellName(colNum+1, rowNum+2)
					f.SetCellValue(sheetName, cell, value)
				}
			}
		}

		// Save the Excel file
		if err := f.SaveAs(excelFile); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		} else {
			fmt.Printf("Excel file has been created: %s\n", excelFile)
		}

	case "4": // Parse JSON data into a map string interface
		// Parse JSON data into a map
		var data map[string]interface{}
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			log.Fatalf("Failed to parse JSON data: %v", err)
		}

		// Create a new Excel file
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("Failed to close Excel file: %v", err)
			}
		}()

		sheetName := "Sheet1"
		index, _ := f.NewSheet(sheetName)

		// Write the column header in the first row
		col := 1 // Start with column 1 (A)
		for key, value := range data {
			// Write the column header (key) in the first row
			cellHeader, _ := excelize.CoordinatesToCellName(col, 1) // Column (col), Row 1
			f.SetCellValue(sheetName, cellHeader, key)

			// Write the value under the column header in the second row
			cellValue, _ := excelize.CoordinatesToCellName(col, 2) // Column (col), Row 2
			f.SetCellValue(sheetName, cellValue, value)

			// Move to the next column
			col++
		}

		// Set active sheet
		f.SetActiveSheet(index)

		// Save the Excel file
		if err := f.SaveAs(excelFile); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		} else {
			fmt.Printf("Excel file has been created: %s\n", excelFile)
		}

	case "5": // Parse JSON data into a map string of string
		var data map[string]string
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			log.Fatalf("Failed to parse JSON data: %v", err)
		}

		// Create a new Excel file
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("Failed to close Excel file: %v", err)
			}
		}()

		sheetName := "Sheet1"
		index, _ := f.NewSheet(sheetName)

		// Maintain order by using a slice of keys
		keys := []string{}
		for key := range data {
			keys = append(keys, key)
		}

		// Write headers (keys from JSON)
		for col, key := range keys {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1) // Header in row 1
			if err := f.SetCellValue(sheetName, cell, key); err != nil {
				log.Fatalf("Error writing header: %v", err)
			}
		}

		// Write values (values from JSON)
		for col, key := range keys {
			cell, _ := excelize.CoordinatesToCellName(col+1, 2) // Values in row 2
			if err := f.SetCellValue(sheetName, cell, data[key]); err != nil {
				log.Fatalf("Error writing value: %v", err)
			}
		}

		// Set active sheet
		f.SetActiveSheet(index)

		// Save the Excel file
		if err := f.SaveAs(excelFile); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		} else {
			fmt.Printf("Excel file has been created: %s\n", excelFile)
		}

	case "6": // Parse JSON data into a JSON array
		var data []string
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			log.Fatalf("Failed to parse JSON data: %v", err)
		}

		// Create a new Excel file
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("Failed to close Excel file: %v", err)
			}
		}()

		sheetName := "Sheet1"
		index, _ := f.NewSheet(sheetName)

		// Write the data horizontally (each string in a new column)
		for i, value := range data {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1) // Column i+1, Row 1
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				log.Fatalf("Error writing value: %v", err)
			}
		}

		// Set active sheet
		f.SetActiveSheet(index)

		// Save the Excel file
		if err := f.SaveAs(excelFile); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		} else {
			fmt.Printf("Excel file has been created: %s\n", excelFile)
		}

	default:
		log.Fatalf("JSON Data Format Invalid: %v", err)
	}
}

// #204
// GetCalcDate - Utility to return a calculated date
// Returns Calculated Date
//   - adjusted by mnth
//   - optionally to return [B]eginning or [E]nding of the month or [N] day in iDate
//   - additionally verifying to return the given day if available in the month
//
// Inputs: YYYYMMDD Date, month between -12 and 12, option B/E/N, original day in iDate
//
// # Outputs: YYYYMMDD calc Date and Time format calc Date.
// ©  FuturaInsTech
func GetCalcDate(iDate string, mnth int, opt string, day int) (osDate string, otDate time.Time) {
	// Sanity edits
	if opt == "" {
		opt = "N"
	}
	if opt != "B" && opt != "E" && opt != "N" {
		fmt.Println("Invalid Option.. correction needed")
		return
	}
	if mnth > 12 || mnth < -12 {
		fmt.Println("Invalid Months.. correction needed")
		return
	}

	var iyear int
	var imonth time.Month
	var iday int
	var oyear int
	var omonth time.Month
	day31mnths := []time.Month{
		time.January,
		time.March,
		time.May,
		time.July,
		time.August,
		time.October,
		time.December,
	}
	itDate, _ := time.Parse("20060102", iDate)
	iyear, imonth, iday = itDate.Date()

	var cmnth time.Month
	cmnth = imonth + time.Month(mnth)
	if opt == "B" {
		iday = 1
	}
	if opt == "E" {
		iday = 0
	}
	if opt == "N" {
		iday = day
	}

	if cmnth <= 0 {
		cmnth = cmnth + 12
		iyear = iyear - 1
	} else if cmnth > 12 {
		cmnth = cmnth - 12
		iyear = iyear + 1
	}

	if opt == "B" {
		otDate = time.Date(iyear, cmnth, iday, 0, 0, 0, 0, itDate.Location())
	} else if opt == "E" {
		otDate = time.Date(iyear, cmnth+1, iday, 0, 0, 0, 0, itDate.Location())
	} else if opt == "N" {
		otDate = time.Date(iyear, cmnth, iday, 0, 0, 0, 0, itDate.Location())
		//Special Handling for February Dates...
		cday := iday
		oyear, omonth, _ = otDate.Date()

		for omonth != cmnth {
			cday = cday - 1
			otDate = time.Date(oyear, cmnth, cday, 0, 0, 0, 0, itDate.Location())
			oyear, omonth, _ = otDate.Date()
		}
		for _, month := range day31mnths {
			if month == cmnth && day != 0 && day > cday {
				otDate = time.Date(iyear, cmnth, day, 0, 0, 0, 0, itDate.Location())
				break
			}
		}

	}
	osDate = otDate.Format("20060102")
	return osDate, otDate
}

// #205
// GetCalcDates - Utility to return a series of calculated dates
// Returns a Series of Calculated Dates
//   - adjusted by mnth
//   - optionally to return [B]eginning or [E]nding of the month or [N] day in iDate
//
// Inputs: DateFrom & DateTo in YYYYMMDD, month between -12 and 12, option B/E/N
// Dependant on #204 GetCalcDate
// # Outputs: An array of Calc Dates in YYYYMMDD between DateFrom & DateTo and Time format calc Date.
// ©  FuturaInsTech
func GetAllCalcDates(iDatefrom string, iDateto string, mnth int, opt string) (osDate []string) {
	// Sanity edits
	if opt == "" {
		opt = "N"
	}
	if opt != "B" && opt != "E" && opt != "N" {
		fmt.Println("Invalid Option.. correction needed")
		return
	}
	if mnth > 12 || mnth < -12 {
		fmt.Println("Invalid Months.. correction needed")
		return
	}
	// Convert the string dates into time format dates
	tDatefrom, _ := time.Parse("20060102", iDatefrom)
	tDateto, _ := time.Parse("20060102", iDateto)
	_, _, iday := tDatefrom.Date()
	// Loop to find end of month dates
	if mnth > 0 && iDatefrom < iDateto {
		for current := tDatefrom; !current.After(tDateto); {
			// Get the current date in string format
			scurrDate := current.Format("20060102")
			sDate, tDate := GetCalcDate(scurrDate, mnth, opt, iday)
			if sDate < iDateto {
				fmt.Println(sDate)
				osDate = append(osDate, sDate)
			}
			year, month, day := tDate.Date()
			current = time.Date(year, month, day, 0, 0, 0, 0, current.Location())
		}
	} else if mnth < 0 && iDatefrom > iDateto {
		for current := tDatefrom; !current.Before(tDateto); {
			// Get the current date in string format
			scurrDate := current.Format("20060102")
			sDate, tDate := GetCalcDate(scurrDate, mnth, opt, iday)
			if sDate > iDateto {
				fmt.Println(sDate)
				osDate = append(osDate, sDate)
			}
			year, month, day := tDate.Date()
			current = time.Date(year, month, day, 0, 0, 0, 0, current.Location())
		}
	} else {
		fmt.Println("Dates are jumbled.. correction needed")
		return
	}
	return osDate
}

// #206
// GetDateOpt - Utility to find if the date is [B]eginning or [E]nding or [N]ormal Date
// Returns the option value as B/E/N of given date
//
// Input: iDate
//
// # Output: B or E or N.
// ©  FuturaInsTech
func GetDateOpt(iDate string) (opt string) {
	tdate, _ := time.Parse("20060102", iDate)
	// Subtract one day to the date
	prevDay := tdate.AddDate(0, 0, -1)
	if tdate.Month() != prevDay.Month() {
		return "B"
	}
	// Add one day to the date
	nextDay := tdate.AddDate(0, 0, 1)
	if tdate.Month() != nextDay.Month() {
		return "E"
	}
	return "N"
}

// #207
// CheckDateOpt - Utility to confirm if Date's opt value given is correct or incorrect
// Returns true if the given date satifies the opt value given
//
// Input: iDate
//
// # Output: true or false
// ©  FuturaInsTech
func CheckDateOpt(iDate string, opt string) bool {
	tdate, _ := time.Parse("20060102", iDate)
	// Subtract one day to the date
	prevDay := tdate.AddDate(0, 0, -1)
	if tdate.Month() != prevDay.Month() && (opt == "B") {
		return true
	}
	// Add one day to the date
	nextDay := tdate.AddDate(0, 0, 1)
	if tdate.Month() != nextDay.Month() && (opt == "E") {
		return true
	}
	return false
}

// #208
func GetDepDes(iDepCoad string, iCompany uint, iLanguage uint) (oDepCoad string) {

	shortdes, _, _ := GetParamDesc(iCompany, "W0006", iDepCoad, iLanguage)

	return shortdes
}

// #209
func GetTeamDes(iTeamCoad string, iCompany uint, iLanguage uint) (oDepCoad string) {

	shortdes, _, _ := GetParamDesc(iCompany, "W0008", iTeamCoad, iLanguage)

	return shortdes

}

// #210
func GetMedInfo(iCompany uint, iMedProv uint, txn *gorm.DB) (oName string, oAddress string, oPin string, oState string, oPhone string, oEmail string, oBank string, oErr string) {
	var medprov models.MedProvider
	result := txn.Find(&medprov, "company_id = ? and id = ?", iCompany, iMedProv)
	if result.Error != nil {
		oErr = "Medical Provider Not Found"
		return "", "", "", "", "", "", "", oErr
	}
	oName = medprov.MedProviderName
	var clnt models.Client
	result = txn.Find(&clnt, "company_id = ? and id = ?", iCompany, medprov.ClientID)
	if result.Error != nil {
		oErr = "Medical Provider Client Not Found"
		return "", "", "", "", "", "", "", oErr
	}
	oPhone = clnt.ClientAltMobCode + " " + clnt.ClientMobile
	oEmail = clnt.ClientEmail + "&" + clnt.ClientAltEmail

	var address models.Address
	result = txn.Find(&address, "company_id = ? and id = ?", iCompany, medprov.AddressID)
	if result.Error != nil {
		oErr = "Medical Provider Address Not Found"
		return "", "", "", "", "", "", "", oErr
	}

	oAddress = address.AddressLine1 + "," + address.AddressLine2 + "," + address.AddressLine3 + "," + address.AddressLine4 + "," + address.AddressLine5
	oPin = address.AddressPostCode
	oState = address.AddressState
	var bank models.Bank
	result = txn.Find(&address, "company_id = ? and id = ?", iCompany, medprov.BankID)
	if result.Error != nil {
		oErr = "Medical Provider Bank Not Found"
		return "", "", "", "", "", "", "", oErr
	}
	oBank = bank.BankCode + "-" + bank.BankAccountNo

	return
}

// #211
func TDFAnnPN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var annuity models.Annuity
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error

	}
	result = txn.First(&annuity, "company_id = ? and policy_id = ? and paystatus = ?", iCompany, iPolicy, "PN")

	if result.Error != nil {
		//	txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ? ", iCompany, iPolicy, iFunction)

	if result.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = annuity.AnnNxtDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = annuity.AnnNxtDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	}
}

// #212
// This function to get User Name by providing Company No and User id
func GetUserName(iCompany uint, iUserId uint) (oName string, oErr error) {
	var usrenq models.User
	result := initializers.DB.Find(&usrenq, "company_id = ? and id = ?", iCompany, iUserId)
	if result.Error != nil {
		return "", result.Error
	}
	return usrenq.Name, nil
}

// #213
// TDF Function For Annuities
func TDFAnniPN(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, error) {
	var annuity models.Annuity
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error

	}
	result = txn.First(&annuity, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	if result.Error != nil {
		//	txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ? ", iCompany, iPolicy, iFunction)

	if result.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = annuity.AnnNxtDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = annuity.AnnNxtDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	}
}

// #214
func GetReqComm(iCompany uint, iPolicy uint, iClient uint, txn *gorm.DB) (map[string]interface{}, error) {
	var reqcall []models.ReqCall
	var client models.Client
	var address models.Address

	medDetailsArray := make([]string, 0) // Array for medDetails
	reqCodeArray := make([]string, 0)    // Array for ReqCode
	reqIDArray := make([]uint, 0)        // Array for Req.ID
	remiderdateArray := make([]string, 0)

	// txn := initializers.DB.Begin()

	txn.Find(&reqcall, "company_id = ? and policy_id = ? and req_status = ?", iCompany, iPolicy, "P")
	txn.Find(&client, "company_id = ? and id = ?", iCompany, iClient)
	txn.Find(&address, "company_id = ? and client_id = ?", iCompany, iClient)

	// Populate data from reqcall
	for _, req := range reqcall {
		oMedName, oMedAddress, oMedPin, oMedState, oMedPhone, _, _, _ := GetMedInfo(iCompany, req.MedId, txn)
		oDesc := GetP0050ItemCodeDesc(iCompany, "REQCODE", 1, req.ReqCode)
		effDate, _ := ConvertYYYYMMDD(req.RemindDate)
		// When Medical Provider is Empty Do not Print Blank Address
		medDetails := ""
		if req.MedId != 0 {
			medDetails = fmt.Sprintf("%s, %s, %s, %s, %s", oMedName, oMedAddress, oMedPin, oMedState, oMedPhone)
		}
		medDetailsArray = append(medDetailsArray, medDetails) // Store medDetails in the array

		reqCodeArray = append(reqCodeArray, oDesc) // Store ReqCode in the array

		reqIDArray = append(reqIDArray, req.ID) // Store Req.ID in the array
		remiderdateArray = append(remiderdateArray, effDate)
	}

	// Create resultMap and include all data from clientInfo
	resultMap := make(map[string]interface{})

	// Include all data directly (without the "clientInfo" label)
	clientInfo := map[string]interface{}{
		"ClientFullName":   client.ClientLongName,
		"ClientSalutation": client.Salutation,
		"AddressLine1":     address.AddressLine1,
		"AddressLine2":     address.AddressLine2,
		"AddressLine3":     address.AddressLine3,
		"AddressLine4":     address.AddressLine4,
		"AddressState":     address.AddressState,
		"AddressCountry":   address.AddressCountry,
		"PolicyId":         IDtoPrint(iPolicy),
		"MedDetails":       medDetailsArray,
		"ReqCodes":         reqCodeArray,
		"ReqIDs":           reqIDArray,
		"Reminderdates":    remiderdateArray,
	}

	for key, value := range clientInfo {
		resultMap[key] = value
	}

	return resultMap, nil
}

// #215
// This Method to create payments for the payable entry.  It can be used wherever we need
// Automatic Approval and Payment Creation
func AutoPayCreate(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iBank uint, iAccCurr string, iAmount float64, iDate string, iDrAcc string, iCrAcc string, iTypeofPayment string, iUserID uint, iReason string, iHistoryCode string, iTranno uint, iPayStatus string, iCoverage string, txn *gorm.DB) (oPayno uint, oErr error) {
	if iPayStatus == "PN" {
		var payosbal models.PayOsBal
		result := txn.Find(&payosbal, "company_id = ? and gl_accountno = ? and gl_rldg_acct =? and contract_curry = ?", iCompany, iDrAcc, iPolicy, iAccCurr)
		iErr := "Payment Already Processed"
		if result.RowsAffected > 0 {
			return 0, errors.New(iErr)
		}
	}

	oPayno = 0
	var bankenq models.Bank
	result := txn.Find(&bankenq, "id = ?", iBank)
	if result.Error != nil {
		return oPayno, result.Error
	}
	iDrSign := "+"
	iCrSign := "-"

	// Following change has been commented. It is not required.  This entry is handled through
	// P0027 parameter set up
	// if iHistoryCode == "H0211" {
	// 	iDrSign = "-"
	// 	iCrSign = "+"
	// }
	// Get Payment Type Accounting Code for Creation
	var p0055data paramTypes.P0055Data
	var extradatap0055 paramTypes.Extradata = &p0055data

	err := GetItemD(int(iCompany), "P0055", iTypeofPayment, iDate, &extradatap0055)
	if err != nil {
		return oPayno, err
	}
	iCrBank := p0055data.GlAccount
	iFSC := p0055data.BankCode
	iCrAccount := iCrAcc + "-" + iCrBank // BankAccount-KVB
	iEffectiveDate := iDate
	// Create Payment First.  Then when it is Auto Approved Payment write accounting entries
	// Write Payment

	var paycrt models.Payment
	paycrt.AccAmount = iAmount
	paycrt.AccCurry = iAccCurr
	paycrt.AddressID = iAddress
	paycrt.BankAccountNo = bankenq.BankAccountNo
	paycrt.BankIFSC = bankenq.BankCode
	paycrt.Branch = "HO"
	paycrt.CheckerUserID = 1
	paycrt.ClientID = iClient
	paycrt.CompanyID = iCompany
	paycrt.CurrentDate = iEffectiveDate
	paycrt.DateOfPayment = iEffectiveDate
	paycrt.InsurerBankAccNo = iCrAccount
	paycrt.InsurerBankIFSC = iFSC
	paycrt.PaymentAccount = iDrAcc + iCoverage
	paycrt.PolicyID = iPolicy
	paycrt.TypeOfPayment = iTypeofPayment
	paycrt.UpdatedID = 1
	paycrt.MakerUserID = 2
	paycrt.Reason = iReason
	paycrt.Status = iPayStatus
	result = txn.Save(&paycrt)

	if result.Error != nil {
		return oPayno, result.Error
	}

	oPayno = paycrt.ID
	oPolicy := strconv.Itoa(int(iPolicy))
	if iPayStatus == "PN" {
		var payosbalcrt models.PayOsBal
		payosbalcrt.CompanyID = iCompany
		payosbalcrt.GlRldgAcct = oPolicy
		payosbalcrt.GlRdocno = oPolicy
		payosbalcrt.GlAccountno = iDrAcc
		payosbalcrt.ContractCurry = iAccCurr
		payosbalcrt.PaymentNo = oPayno
		payosbalcrt.ContractAmount = iAmount
		result = txn.Create(&payosbalcrt)
		if result.Error != nil {
			return 0, result.Error
		}

	}
	if iPayStatus == "AP" {
		// Debit
		glcode := iDrAcc
		var acccode models.AccountCode
		result = txn.First(&acccode, "company_id = ? and account_code = ? ", iCompany, glcode)
		if result.RowsAffected == 0 {
			return oPayno, result.Error
		}
		var iSequenceno uint64
		iSequenceno++
		iAccountCodeID := acccode.ID
		iAccAmount := iAmount
		iAccCurry := iAccCurr
		iAccountCode := glcode + iCoverage

		iGlAmount := iAmount

		iGlRdocno := int(iPolicy)
		var iGlRldgAcct string
		//iGlRldgAcct := strconv.Itoa(int(iClient))
		// As per our discussion on 22/06/2023, it is decided to use policy no in RLDGACCT
		iGlRldgAcct = strconv.Itoa(int(iPolicy))
		iGlSign := iDrSign

		err = PostGlMoveN(iCompany, iAccCurry, iEffectiveDate, int(iTranno), iGlAmount,
			iAccAmount, iAccountCodeID, uint(iGlRdocno), string(iGlRldgAcct), iSequenceno, iGlSign, iAccountCode, iHistoryCode, "", "", txn)

		if err != nil {
			return oPayno, err
		}
		// Credit

		glcode = iCrAcc
		var acccode1 models.AccountCode
		result = txn.First(&acccode1, "company_id = ? and account_code = ? ", iCompany, glcode)
		if result.RowsAffected == 0 {
			return oPayno, result.Error
		}

		iSequenceno++
		iAccountCodeID = acccode1.ID
		iAccAmount = iAmount
		iAccCurry = iAccCurr
		iAccountCode = iCrAccount
		iEffectiveDate = iDate
		iGlAmount = iAmount

		//iGlRdocno = int(iPolicy)
		iGlRdocno = int(oPayno)

		//iGlRldgAcct := strconv.Itoa(int(iClient))
		// As per our discussion on 22/06/2023, it is decided to use policy no in RLDGACCT
		iGlRldgAcct = strconv.Itoa(int(iPolicy))
		iGlSign = iCrSign

		err = PostGlMoveN(iCompany, iAccCurry, iEffectiveDate, int(iTranno), iGlAmount,
			iAccAmount, iAccountCodeID, uint(iGlRdocno), string(iGlRldgAcct), iSequenceno, iGlSign, iAccountCode, iHistoryCode, "", "", txn)

		if err != nil {
			return oPayno, err
		}
	}

	return oPayno, nil
}

// #216
func PolicyDep(iCompany uint, iPolicy uint) map[string]interface{} {
	var polenq models.Policy
	result := initializers.DB.Where("company_id = ? AND id = ?", iCompany, iPolicy).Find(&polenq)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Policy not found"}
	}

	var clnt models.Client
	result = initializers.DB.Where("company_id = ? AND id = ?", iCompany, polenq.ClientID).Find(&clnt)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Client not found"}
	}

	var address models.Address
	result = initializers.DB.Where("company_id = ? AND id = ?", iCompany, polenq.ClientID).Find(&address)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Address not found"}
	}

	var pymt models.Payment
	result = initializers.DB.Where("company_id = ? AND policy_id = ?", iCompany, iPolicy).Find(&pymt)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Payment not found"}
	}

	var glbal models.GlBal
	result = initializers.DB.Where("company_id = ? AND gl_rdocno = ? ", iCompany, iPolicy).Find(&glbal)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "GlBal not found"}
	}

	// Create a result map for each loan bill
	resultOut := map[string]interface{}{
		"ClientFullName":   clnt.ClientLongName,
		"ClientSalutation": clnt.Salutation,
		"AddressLine1":     address.AddressLine1,
		"AddressLine2":     address.AddressLine2,
		"AddressLine3":     address.AddressLine3,
		"AddressLine4":     address.AddressLine4,
		"AddressState":     address.AddressState,
		"AddressCountry":   address.AddressCountry,
		"PaymentType":      pymt.TypeOfPayment,
		"PaymentId":        pymt.ID,
		"PaymentDate":      DateConvert(pymt.CurrentDate),
		"Amount":           glbal.ContractAmount,
		"policyid":         polenq.ID,
	}

	return resultOut
}

// # 217
// GetHealthRate - Get Annual Rate of the Coverage - No Model Discount/Staff Discount/SA/Prem Discount
//
// Inputs: Company, Coverage, Age (Attained Age), Gender(F/N/U), Term (2 Characters), Premium Method as
// PM001 - Term Based , PM002 Age Based, Mortality Clause "S" Smoker, "N" Non Smoker
//
// Outputs Annualized Premium as float (124.22)
//
// ©  FuturaInsTech
func GetHeathRate(iCompany, iPolicy uint, iCoverage, iPlan string, iDate string, iPremRateCode string, iPremAge uint) (float64, error) {

	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)

	ikey := iCoverage + iPlan

	var p0074data paramTypes.P0074Data
	var extradatap0074 paramTypes.Extradata = &p0074data
	err := GetItemD(int(iCompany), "P0074", ikey, iDate, &extradatap0074)
	if err != nil {
		return 0, err

	}

	var p0080data paramTypes.P0080Data
	var extradatap0080 paramTypes.Extradata = &p0080data
	err = GetItemD(int(iCompany), "P0080", ikey, iDate, &extradatap0080)
	if err != nil {
		return 0, err

	}

	//iAgeCalcMethod := q0006data.AgeCalcMethod
	//iPlanPremAge := p0074data.PlanPremAge

	// iPremMethod := q0006data.PremiumMethod
	// iDisType := q0006data.DiscType
	// iDisMethod := q0006data.DiscMethod
	// iFrqMethod := q0006data.FrqMethod
	// iWaiverMethod := q0006data.WaivMethod
	// iMRTA := q0006data.MrtaMethod
	// iPremCalcType := q0006data.PremCalcType
	// iHealthBenefitType := q0006data.HealthBenefitType
	// iPlanMaxLives := p0074data.PlanMaxLives

	// var premAge int
	// var PremRateCode string

	// var planlife models.PlanLife

	// if iPlanPremAge == "PL" {

	// 	result := initializers.DB.First(&planlife, "company_id=? and policy_id= ? and client_rel_code=?", iCompany, iPolicy, "I")

	// 	if result.Error != nil || result.RowsAffected == 0 {
	// 		return 0, errors.New("No Primary Life Assured In Plan Life")
	// 	}

	// } else if iPlanPremAge == "YL" {

	// 	result := initializers.DB.Order("p_age ASC").First(&planlife, "company_id = ? AND policy_id = ?", iCompany, iPolicy)

	// 	if result.Error != nil || result.RowsAffected == 0 {
	// 		return 0, errors.New("No Youngest Life Assured In Plan Life")
	// 	}

	// } else if iPlanPremAge == "EL" {
	// 	result := initializers.DB.Order("p_age DESC").First(&planlife, "company_id = ? AND policy_id = ?", iCompany, iPolicy)

	// 	if result.Error != nil || result.RowsAffected == 0 {
	// 		return 0, errors.New("No Eldest Life Assured In Plan Life")
	// 	}
	// }

	// premAge, _, _, _, _, _ = CalculateAge(iDate, planlife.PDOB, iAgeCalcMethod)
	// PremRateCode, err = GetPlanPremRateCode(iCompany, iPolicy, planlife.PSumAssured, ikey, planlife.PStartDate)

	// if err != nil {
	// 	return 0, err
	// }

	var q0010key string
	var prem float64

	// term := strconv.FormatUint(uint64(iTerm), 10)
	// premTerm := strconv.FormatUint(uint64(iPremTerm), 10)

	// if q0006data.PremCalcType == "A" || q0006data.PremCalcType == "U" {
	//  if q0006data.PremiumMethod == "PM002" {
	//      // END1 + Male
	//      q0010key = iCoverage + iGender
	//  }
	// } else if q0006data.PremCalcType == "P" {
	//  // END1 + Male + Term + Premium Term
	//  if q0006data.PremiumMethod == "PM001" || q0006data.PremiumMethod == "PM003" {
	//      q0010key = iCoverage + iGender + term + premTerm
	//  }

	// }

	if q0006data.PremCalcType == "H" {
		if q0006data.PremiumMethod == "PM005" {
			q0010key = iCoverage + iPlan + iPremRateCode
		}
	}

	var q0010data paramTypes.Q0010Data
	var extradataq0010 paramTypes.Extradata = &q0010data
	err = GetItemD(int(iCompany), "Q0010", q0010key, iDate, &extradataq0010)
	if err != nil {
		return 0, err

	}

	for i := 0; i < len(q0010data.Rates); i++ {
		if q0010data.Rates[i].Age == iPremAge {
			prem = q0010data.Rates[i].Rate
			break
		}
	}
	return prem, nil
}

// # 218
// GetPremRateCode - Get Prem Rate Code From P0080 Fro give parameter
//
// Input :- SumAssured , PlanLaCode list and PlanLACodeCount list
//
// Output :- PremRateCode from P0080
//
// ©  FuturaInsTech
func GetPlanPremRateCode(iCompany, iPolicy uint, iSumAssured uint, iKey, iDate string) (string, error) {

	var p0080data paramTypes.P0080Data
	var extradatap0080 paramTypes.Extradata = &p0080data
	err := GetItemD(int(iCompany), "P0080", iKey, iDate, &extradatap0080)
	if err != nil {
		return "", err

	}

	var filteredP0080Records []paramTypes.P0080
	if p0080data.PremRateCodes != nil {
		for _, record := range p0080data.PremRateCodes {
			if record.SumAssured == float64(iSumAssured) {
				filteredP0080Records = append(filteredP0080Records, record) // Append entire record
			}
		}
	}

	var planlifes []models.PlanLife

	result := initializers.DB.Find(&planlifes, "company_id = ? AND policy_id = ?", iCompany, iPolicy)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", errors.New("No Life Assured In Plan Life")
	}

	premiumLACodeCount := make(map[string]int)

	for _, planLife := range planlifes {
		premiumLACodeCount[planLife.PremuimLACode]++
	}

	// Separate keys and values into arrays
	var LACodes []string
	var LACounts []int

	for code, count := range premiumLACodeCount {
		LACodes = append(LACodes, code)
		LACounts = append(LACounts, count)
	}

	// Filter records
	var filteredRecords []paramTypes.P0080
	for _, record := range filteredP0080Records {

		val := reflect.ValueOf(record)

		// Check all lACode fields
		for i, expectedCode := range LACodes {
			codeField := fmt.Sprintf("LACode%d", i+1)
			countField := fmt.Sprintf("LACount%d", i+1)

			// Get field values dynamically
			codeValue := val.FieldByName(codeField).String()
			countValue := int(val.FieldByName(countField).Int())

			// If the code matches, ensure the count also matches
			if expectedCode != "" && codeValue == expectedCode {
				if countValue == LACounts[i] {
					filteredRecords = append(filteredRecords, record)
				}
			}
		}

	}

	if len(filteredRecords) > 1 {
		return "", errors.New("Multiple Premium Life Assured Code")
	}

	return filteredRecords[0].PremRateCode, nil
}

// #219
// GetParamPlanBenefit - Getting Plan BenefitCode ,BenefitUnit ,BenefitBasis
// BenefitPlanCover, PlanBenefitGroup, MaxBenefitAmount, MaxBenefitUnit, MaxBenefitBasis
// From Param P0075(Plan Benefits) and P0077(Plan Max Benefits)
//
// Input :- iCompany, iPolicy, iBenefit, iBCoverage, iBenefitPlan, current date
//
// Output :- PremRateCode from P0080
//
// ©  FuturaInsTech
func GetParamPlanBenefit(iCompany uint, iBCoverage, iBenefitPlan, iDate string) (error, []interface{}) {
	resp := make([]interface{}, 0)

	iKey := iBCoverage + iBenefitPlan

	var p0075data paramTypes.P0075Data
	var extradatap0075 paramTypes.Extradata = &p0075data
	err := GetItemD(int(iCompany), "P0075", iKey, iDate, &extradatap0075)
	if err != nil {
		return err, nil

	}

	var p0077data paramTypes.P0077Data
	var extradatap0077 paramTypes.Extradata = &p0077data
	err = GetItemD(int(iCompany), "P0077", iKey, iDate, &extradatap0077)
	if err != nil {
		return err, nil

	}

	for _, planBenefit := range p0075data.PlanBenefits {

		for _, planMaxBenefit := range p0077data.PlanMaxBenefits {

			if err != nil {
				return err, nil
			}

			benefitCodeDesc := GetP0050ItemCodeDesc(iCompany, "BenefitCode", 1, planBenefit.BenefitCode)
			benefitBasisDesc := GetP0050ItemCodeDesc(iCompany, "BenefitBasis", 1, planBenefit.BenefitBasis)
			benefitPlanCover := GetP0050ItemCodeDesc(iCompany, "BenefitPlanCover", 1, planBenefit.BenefitPlanCover)
			planBenefitGroup := GetP0050ItemCodeDesc(iCompany, "PlanBenefitGroup", 1, planBenefit.PlanBenefitGroup)
			maxBenefitBasis := GetP0050ItemCodeDesc(iCompany, "MaxBenefitBasis", 1, planMaxBenefit.MaxBenefitBasis)

			if planBenefit.BenefitCode == planMaxBenefit.BenefitCode {

				paramOut := map[string]interface{}{
					"BenefitCode":          planBenefit.BenefitCode,
					"BenefitCodeDesc":      benefitCodeDesc,
					"BenefitUnit":          planBenefit.BenefitUnit,
					"BenefitBasis":         planBenefit.BenefitBasis,
					"BenefitBasisDesc":     benefitBasisDesc,
					"BenefitPlanCover":     planBenefit.BenefitPlanCover,
					"BenefitPlanCoverDesc": benefitPlanCover,
					"PlanBenefitGroup":     planBenefit.PlanBenefitGroup,
					"PlanBenefitGroupDesc": planBenefitGroup,
					"MaxBenefitAmount":     planMaxBenefit.MaxBenefitAmount,
					"MaxBenefitUnit":       planMaxBenefit.MaxBenefitUnit,
					"MaxBenefitBasis":      planMaxBenefit.MaxBenefitBasis,
					"MaxBenefitBasisDesc":  maxBenefitBasis,
				}
				resp = append(resp, paramOut)
			}

		}
	}

	return nil, resp
}

// # 220
// CalcHealthDiscounts - Calculate Discounted Amount based on SA or Annualised Prem
//
// Inputs: Company, Discount Type (S/P) , Discount Method (As per Product), Annualised Prem
// SA Amount
//
// # Outputs Discounted Amount as float
//
// ©  FuturaInsTech
func CalcHealthDiscounts(iCompany uint, iDiscType string, iDiscMethod string, iAnnPrem float64, iSA uint, iDate string) float64 {
	// SA Discount

	if iDiscType == "S" {
		var q0017data paramTypes.Q0017Data
		var extradataq0017 paramTypes.Extradata = &q0017data
		err := GetItemD(int(iCompany), "Q0017", iDiscMethod, iDate, &extradataq0017)

		if err != nil {
			return 0

		}

		for i := 0; i < len(q0017data.SaBand); i++ {
			if int(iSA) <= int(q0017data.SaBand[i].Sa) {
				oDiscount := uint(q0017data.SaBand[i].Discount) * uint(iAnnPrem) / 100
				return float64(oDiscount)
			}
		}
	}
	// Premium Discount
	if iDiscType == "P" {
		var q0018data paramTypes.Q0018Data
		var extradataq0018 paramTypes.Extradata = &q0018data

		err := GetItemD(int(iCompany), "Q0018", iDiscMethod, iDate, &extradataq0018)

		if err != nil {
			return 0

		}

		for i := 0; i < len(q0018data.PremBand); i++ {
			if int(iAnnPrem) <= int(q0018data.PremBand[i].AnnPrem) {
				oDiscount := uint(q0018data.PremBand[i].Discount) * uint(iAnnPrem) / 100
				return float64(oDiscount)
			}
		}
	}
	return 0
}
