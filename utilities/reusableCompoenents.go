package utilities

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"

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
func ValidateItem(iUserId uint64, iName string, iItem any, iFieldName string, iErros string) error {
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
	term := strconv.FormatUint(uint64(iTerm), 10)
	premTerm := strconv.FormatUint(uint64(iPremTerm), 10)
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
	iHistoryCD := transaction.Method
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

	iHistoryCD := transaction.Method
	oHistory = iHistoryCD
	var p0029data paramTypes.P0029Data
	var extradata paramTypes.Extradata = &p0029data
	fmt.Println("Transaction Foound !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", iHistoryCode)
	err := GetItemD(int(iCompany), "P0029", iHistoryCD, iDate, &extradata)

	fmt.Println("Newstatus", iStatus, p0029data.Statuses[1].CurrentStatus, p0029data.Statuses[1].ToBeStatus)
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
func GetDeathAmount(iCompany uint, iPolicy uint, iProduct string, iCoverage string, iEffectiveDate string, iCause string) (oAmount float64) {
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

	err = GetItemD(int(iCompany), "Q0005", iCoverage, iDate, &extradataq0005)
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

	otrancode = transaction.Method
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
	iRCD := policyenq.PRCD
	var benefitenq1 []models.Benefit

	results := initializers.DB.Find(&benefitenq1, "company_id =? and policy_id = ? ", iCompany, iPolicy)
	if results.Error != nil {
		return 0
	}

	for a := 0; a < len(benefitenq1); a++ {
		FromDate := iFromDate
		ToDate := iToDate

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
				iMonths := NewNoOfInstalments(iRCD, FromDate)
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
		oAmount = 0
		return oAmount
	case imatMethod == "MAT099": // No Maturity Value
		oAmount = 0
		break
	default:
		oAmount = 0
		return
	}

	return oAmount
}

// #86
// GetCompany Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Company Details as an Interface
//
// ©  FuturaInsTech
func GetCompanyData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iDate string) []interface{} {
	companyarray := make([]interface{}, 0)
	var company models.Company
	initializers.DB.Find(&company, "id = ?", iCompany)

	resultOut := map[string]interface{}{
		"ID":                       IDtoPrint(company.ID),
		"CompanyName":              company.CompanyName,
		"CompanyAddress1":          company.CompanyAddress1,
		"CompanyAddress2":          company.CompanyAddress2,
		"CompanyAddress3":          company.CompanyAddress3,
		"CompanyAddress4":          company.CompanyAddress4,
		"CompanyAddress5":          company.CompanyAddress5,
		"CompanyPostalCode":        company.CompanyPostalCode,
		"CompanyCountry":           company.CompanyCountry,
		"CompanyUid":               company.CompanyUid,
		"CompanyGst":               company.CompanyGst,
		"CompanyPan":               company.CompanyPan,
		"CompanyTan":               company.CompanyTan,
		"CompanyIncorporationDate": DateConvert(company.CompanyIncorporationDate),
		"CompanyTerminationDate":   DateConvert(company.CompanyTerminationDate),
		"LetterDate":               DateConvert(iDate),
	}
	companyarray = append(companyarray, resultOut)
	return companyarray
}

// #87
// GetClient Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Client Details as an Interface
//
// ©  FuturaInsTech
func GetClientData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	clientarray := make([]interface{}, 0)
	var client models.Client

	initializers.DB.Find(&client, "company_id = ? and id = ?", iCompany, iClient)
	resultOut := map[string]interface{}{
		"ID":              IDtoPrint(client.ID),
		"ClientShortName": client.ClientShortName,
		"ClientLongName":  client.ClientLongName,
		"ClientSurName":   client.ClientSurName,
		"Gender":          client.Gender,
		"Salutation":      client.Salutation,
		"Language":        client.Language,
		"ClientDob":       DateConvert(client.ClientDob),
		// "ClientDod":DateConvert(client.ClientDod),
		"ClientEmail":  client.ClientEmail,
		"ClientMobile": client.ClientMobile,
		"ClientStatus": client.ClientStatus,
	}
	clientarray = append(clientarray, resultOut)
	return clientarray
}

// #88
// GetAddressData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Address Details as an Interface
//
// ©  FuturaInsTech
func GetAddressData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	addressarray := make([]interface{}, 0)
	var address models.Address

	initializers.DB.Find(&address, "company_id = ? and id = ?", iCompany, iAddress)
	resultOut := map[string]interface{}{
		"ID":              IDtoPrint(address.ID),
		"AddressType":     address.AddressType,
		"AddressLine1":    address.AddressLine1,
		"AddressLine2":    address.AddressLine2,
		"AddressLine3":    address.AddressLine3,
		"AddressLine4":    address.AddressLine4,
		"AddressLine5":    address.AddressLine5,
		"AddressPostCode": address.AddressPostCode,
		"AddressState":    address.AddressState,
		"AddressCountry":  address.AddressCountry,
	}
	addressarray = append(addressarray, resultOut)
	return addressarray
}

// #89
// GetPolicyData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Policy Details as an Interface
//
// ©  FuturaInsTech
func GetPolicyData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	policyarray := make([]interface{}, 0)
	var policy models.Policy
	result := initializers.DB.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return nil
	}
	_, oStatus, _ := GetParamDesc(policy.CompanyID, "P0024", policy.PolStatus, 1)
	_, oFreq, _ := GetParamDesc(policy.CompanyID, "Q0009", policy.PFreq, 1)
	_, oProduct, _ := GetParamDesc(policy.CompanyID, "Q0005", policy.PProduct, 1)
	_, oBillCurr, _ := GetParamDesc(policy.CompanyID, "P0023", policy.PBillCurr, 1)
	_, oContCurr, _ := GetParamDesc(policy.CompanyID, "P0023", policy.PContractCurr, 1)

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)
	gracedate := AddLeadDays(policy.PaidToDate, q0005data.LapsedDays)
	premduedates := GetPremDueDates(policy.PRCD, policy.PFreq)
	iAnnivDate := Date2String(GetNextDue(policy.AnnivDate, "Y", "R"))

	var benefitenq []models.Benefit

	initializers.DB.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	oRiskCessDate := ""
	oPremCessDate := ""
	for i := 0; i < len(benefitenq); i++ {
		if oRiskCessDate < benefitenq[i].BRiskCessDate {
			oRiskCessDate = benefitenq[i].BRiskCessDate
		}
		if oPremCessDate < benefitenq[i].BPremCessDate {
			oPremCessDate = benefitenq[i].BPremCessDate
		}
	}

	oAnnivDate := String2Date(iAnnivDate)
	oPRCD := String2Date(policy.PRCD)
	ocompletedyears, _, _, _, _, _ := DateDiff(oAnnivDate, oPRCD, "")

	sRiskCessDate := String2Date(oRiskCessDate)
	sPRCD := String2Date(policy.PRCD)
	oRiskTerm, _, _, _, _, _ := DateDiff(sRiskCessDate, sPRCD, "")

	sPremCessDate := String2Date(oPremCessDate)
	sPRCD = String2Date(policy.PRCD)
	oPremTerm, _, _, _, _, _ := DateDiff(sPremCessDate, sPRCD, "")

	resultOut := map[string]interface{}{
		"ID":                 IDtoPrint(policy.ID),
		"CompanyID":          IDtoPrint(policy.CompanyID),
		"PRCD":               DateConvert(policy.PRCD),
		"PProduct":           oProduct,
		"PFreq":              oFreq,
		"PContractCurr":      oContCurr,
		"PBillCurr":          oBillCurr,
		"POffice":            policy.POffice,
		"PolStatus":          oStatus,
		"PReceivedDate":      DateConvert(policy.PReceivedDate),
		"ClientID":           IDtoPrint(policy.ClientID),
		"BTDate":             DateConvert(policy.BTDate),
		"PaidToDate":         DateConvert(policy.PaidToDate),
		"NxtBTDate":          DateConvert(policy.NxtBTDate),
		"AnnivDate":          DateConvert(policy.AnnivDate),
		"AgencyID":           IDtoPrint(policy.AgencyID),
		"InstalmentPrem":     NumbertoPrint(policy.InstalmentPrem),
		"GracePeriodEndDate": DateConvert(gracedate),
		"RiskCessDate":       DateConvert(oRiskCessDate),
		"PremCessDate":       DateConvert(oPremCessDate),
		"CompletedYears":     ocompletedyears,
		"PolicyRiskTerm":     oRiskTerm,
		"PolicyPremTerm":     oPremTerm,
		"GraceDays":          q0005data.LapsedDays,
		"PremiumDueDates":    premduedates,
		"PrevAnnivDate":      DateConvert(iAnnivDate),
		// "PUWDate":DateConvert(policy.PUWDate),
	}
	policyarray = append(policyarray, resultOut)

	fmt.Print(policyarray)
	return policyarray
}

// #90
// GetBenefitData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Benefit Details as an Interface
//
// ©  FuturaInsTech
func GetBenefitData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var policyenq models.Policy
	var benefit []models.Benefit
	var clientenq models.Client
	var addressenq models.Address
	initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	paidToDate := policyenq.PaidToDate
	nextDueDate := policyenq.NxtBTDate
	initializers.DB.Find(&benefit, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	benefitarray := make([]interface{}, 0)

	for k := 0; k < len(benefit); k++ {
		iCompany := benefit[k].CompanyID
		_, oGender, _ := GetParamDesc(iCompany, "P0001", benefit[k].BGender, 1)
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", benefit[k].BCoverage, 1)
		_, oStatus, _ := GetParamDesc(iCompany, "P0024", benefit[k].BStatus, 1)

		clientname := GetName(iCompany, benefit[k].ClientID)
		initializers.DB.Find(&clientenq, "company_id = ? and id = ?", iCompany, benefit[k].ClientID)
		initializers.DB.Find(&addressenq, "company_id = ? and client_id = ?", iCompany, clientenq.ID)
		clientdob := clientenq.ClientDob
		address := addressenq.AddressLine1 + "" + addressenq.AddressLine2 + "" + addressenq.AddressLine3 + "" + addressenq.AddressLine4 + "" + addressenq.AddressLine5 + "" + addressenq.AddressPostCode + "" + addressenq.AddressState

		resultOut := map[string]interface{}{

			"ID":                 IDtoPrint(benefit[k].ID),
			"CompanyID":          IDtoPrint(benefit[k].CompanyID),
			"ClientID":           IDtoPrint(benefit[k].ClientID),
			"PolicyID":           IDtoPrint(benefit[k].PolicyID),
			"BStartDate":         DateConvert(benefit[k].BStartDate),
			"BRiskCessDate":      DateConvert(benefit[k].BRiskCessDate),
			"BPremCessDate":      DateConvert(benefit[k].BPremCessDate),
			"BTerm":              benefit[k].BTerm,
			"BPTerm":             benefit[k].BPTerm,
			"BRiskCessAge":       benefit[k].BRiskCessAge,
			"BPremCessAge":       benefit[k].BPremCessAge,
			"BBasAnnualPrem":     NumbertoPrint(benefit[k].BBasAnnualPrem),
			"BLoadPrem":          NumbertoPrint(benefit[k].BLoadPrem),
			"BCoverage":          oCoverage,
			"BSumAssured":        NumbertoPrint(float64(benefit[k].BSumAssured)),
			"BPrem":              NumbertoPrint(benefit[k].BPrem),
			"BGender":            oGender,
			"BDOB":               benefit[k].BDOB,
			"BMortality":         benefit[k].BMortality,
			"BStatus":            oStatus,
			"BAge":               benefit[k].BAge,
			"BRerate":            benefit[k].BRerate,
			"LifeAssuredName":    clientname,
			"LifeAssuredAddress": address,
			"LifeAssuredDOB":     DateConvert(clientdob),
			"paidToDate":         DateConvert(paidToDate),
			"nextDueDate":        DateConvert(nextDueDate),
		}
		benefitarray = append(benefitarray, resultOut)
	}
	return benefitarray
}

// #91
// GetSurvivalBenefit - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Survival Details as an Interface
//
// ©  FuturaInsTech
func GetSurBData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var survb []models.SurvB
	initializers.DB.Find(&survb, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	var benefitenq models.Benefit
	initializers.DB.Find(&benefitenq, "company_id = ? and policy_id =? and id = ?", iCompany, iPolicy, survb[0].BenefitID)
	basis := ""
	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data

	GetItemD(int(iCompany), "Q0006", benefitenq.BCoverage, benefitenq.BStartDate, &extradataq0006)
	if q0006data.SbType == "A" {
		basis = "Age Based Survival Benefit"
	} else {
		basis = "Term Based Survival Benefit"
	}

	survbarray := make([]interface{}, 0)
	for k := 0; k < len(survb); k++ {
		resultOut := map[string]interface{}{
			"ID":            IDtoPrint(survb[k].ID),
			"CompanyID":     IDtoPrint(survb[k].CompanyID),
			"BenefitID":     IDtoPrint(survb[k].BenefitID),
			"PolicyID":      IDtoPrint(survb[k].PolicyID),
			"EffectiveDate": DateConvert(survb[k].EffectiveDate),
			"Basis":         basis,
			"SeqNo":         survb[k].Sequence,
			"SBPercentage":  survb[k].SBPercentage,
			"SBAmount":      survb[k].Amount,
			//		"PaidDate ":     DateConvert(survb[k].PaidDate),
		}
		survbarray = append(survbarray, resultOut)
	}
	return survbarray
}

// #92
// GetMRTAData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  MRTA Details as an Interface
//
// ©  FuturaInsTech
func GetMrtaData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var mrtaenq []models.Mrta
	initializers.DB.Find(&mrtaenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	mrtaarray := make([]interface{}, 0)
	for k := 0; k < len(mrtaenq); k++ {
		resultOut := map[string]interface{}{
			"ID":            IDtoPrint(mrtaenq[k].ID),
			"CompanyID":     IDtoPrint(mrtaenq[k].CompanyID),
			"Term":          mrtaenq[k].BTerm,
			"Ppt":           mrtaenq[k].PremPayingTerm,
			"ClientID":      IDtoPrint(mrtaenq[k].ClientID),
			"BenefitID":     IDtoPrint(mrtaenq[k].BenefitID),
			"PolicyID":      IDtoPrint(mrtaenq[k].PolicyID),
			"Coverage":      mrtaenq[k].BCoverage,
			"Product":       mrtaenq[k].Pproduct,
			"Interest":      mrtaenq[k].Interest,
			"DecreaseSA":    mrtaenq[k].BSumAssured,
			"InterimPeriod": mrtaenq[k].InterimPeriod,
			"StartDate":     DateConvert(mrtaenq[k].BStartDate),
		}
		mrtaarray = append(mrtaarray, resultOut)
	}
	return mrtaarray
}

// #93
// GetReceiptData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Receipt Details as an Interface
//
// ©  FuturaInsTech
func GetReceiptData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var receiptenq models.Receipt
	initializers.DB.Find(&receiptenq, "company_id = ? and id = ?", iCompany, iReceipt)
	amtinwords, csymbol := AmountinWords(receiptenq.CompanyID, receiptenq.AccAmount, receiptenq.AccCurry)
	receiptarray := make([]interface{}, 0)
	resultOut := map[string]interface{}{
		"ID":                IDtoPrint(receiptenq.ID),
		"CompanyID":         IDtoPrint(receiptenq.CompanyID),
		"Branch":            receiptenq.Branch,
		"AccCurry":          receiptenq.AccCurry,
		"AccAmount":         NumbertoPrint(receiptenq.AccAmount),
		"PolicyID":          IDtoPrint(receiptenq.PolicyID),
		"ClientID":          IDtoPrint(receiptenq.ClientID),
		"DateOfCollection":  DateConvert(receiptenq.DateOfCollection),
		"BankAccountNo":     receiptenq.BankAccountNo,
		"BankReferenceNo":   receiptenq.BankReferenceNo,
		"TypeOfReceipt":     receiptenq.TypeOfReceipt,
		"InstalmentPremium": receiptenq.InstalmentPremium,
		"AddressID":         IDtoPrint(receiptenq.AddressID),
		"AmountInWords":     amtinwords,
		"CurrSymbol":        csymbol,
		//		"PaidToDate":        DateConvert(receiptenq.PaidToDate),
		//		"ReconciledDate":    DateConvert(receiptenq.ReconciledDate),
		//		"CurrentDate":       DateConvert(receiptenq.CurrentDate),
	}
	receiptarray = append(receiptarray, resultOut)

	return receiptarray
}

// #94
// GetSAChangeData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  SA Change Details as an Interface
//
// ©  FuturaInsTech
func GetSaChangeData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var sachangeenq []models.SaChange
	initializers.DB.Find(&sachangeenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	sachangearray := make([]interface{}, 0)
	for k := 0; k < len(sachangeenq); k++ {
		_, oGender, _ := GetParamDesc(iCompany, "P0001", sachangeenq[k].BGender, 1)
		_, oFreq, _ := GetParamDesc(iCompany, "P0024", sachangeenq[k].Frequency, 1)
		resultOut := map[string]interface{}{
			"PolicyID    ": IDtoPrint(sachangeenq[k].PolicyID),
			"BenefitID":    IDtoPrint(sachangeenq[k].BenefitID),
			"BCoverage":    sachangeenq[k].BCoverage,
			"BStartDate":   DateConvert(sachangeenq[k].BStartDate),
			"BSumAssured":  sachangeenq[k].BSumAssured,
			"BTerm":        IDtoPrint(sachangeenq[k].BTerm),
			"BPTerm":       sachangeenq[k].BPTerm,
			"BPrem":        sachangeenq[k].BPrem,
			"BGender":      oGender,
			"BDOB":         DateConvert(sachangeenq[k].BDOB),
			"NSumAssured":  sachangeenq[k].NSumAssured,
			"NTerm":        sachangeenq[k].NTerm,
			"NPTerm":       sachangeenq[k].NPTerm,
			"NPrem":        sachangeenq[k].NPrem,
			"NAnnualPrem":  sachangeenq[k].NAnnualPrem,
			"Method":       sachangeenq[k].Method,
			"Frequency":    oFreq,
		}
		sachangearray = append(sachangearray, resultOut)
	}
	return sachangearray
}

// #95
// GetCompAddData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Component Add Details as an Interface
//
// ©  FuturaInsTech
func GetCompAddData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var addcomp []models.Addcomponent
	initializers.DB.Find(&addcomp, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	addcomparray := make([]interface{}, 0)

	for k := 0; k < len(addcomp); k++ {
		_, oFreq, _ := GetParamDesc(iCompany, "P0024", addcomp[k].Frequency, 1)
		_, oGender, _ := GetParamDesc(iCompany, "P0001", addcomp[k].BGender, 1)
		resultOut := map[string]interface{}{

			"ID":          IDtoPrint(addcomp[k].ID),
			"Select":      addcomp[k].Select,
			"PolicyID":    IDtoPrint(addcomp[k].PolicyID),
			"ClientID":    IDtoPrint(addcomp[k].ClientID),
			"BCoverage":   addcomp[k].BCoverage,
			"BStartDate":  DateConvert(addcomp[k].BStartDate),
			"BSumAssured": NumbertoPrint(float64(addcomp[k].BSumAssured)),
			"BTerm":       addcomp[k].BTerm,
			"BPTerm":      addcomp[k].BPTerm,
			"BPrem":       NumbertoPrint(addcomp[k].BPrem),
			"BAnnualPrem": NumbertoPrint(addcomp[k].BAnnualPrem),
			"BGender":     oGender,
			"BDOB":        addcomp[k].BDOB,
			"Method":      addcomp[k].Method,
			"Frequency":   oFreq,
			"BAge":        addcomp[k].BAge,
		}
		addcomparray = append(addcomparray, resultOut)
	}
	return addcomparray
}

// #96
// GetSurrenderH Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Surrender Header  and Details as an Interface
//
// ©  FuturaInsTech
func GetSurrHData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) interface{} {
	var surrhenq models.SurrH

	initializers.DB.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	surrharray := make([]interface{}, 0)
	_, oProduct, _ := GetParamDesc(iCompany, "Q0005", surrhenq.Product, 1)
	resultOut := map[string]interface{}{

		"ID":                IDtoPrint(surrhenq.ID),
		"PolicyID":          IDtoPrint(surrhenq.PolicyID),
		"ClientID":          IDtoPrint(surrhenq.ClientID),
		"EffectiveDate":     DateConvert(surrhenq.EffectiveDate),
		"SurrDate":          DateConvert(surrhenq.SurrDate),
		"Cause":             surrhenq.Cause,
		"Status":            surrhenq.Status,
		"BillDate":          DateConvert(surrhenq.BillDate),
		"PaidToDate":        DateConvert(surrhenq.PaidToDate),
		"Product":           oProduct,
		"AplAmount":         float64(surrhenq.AplAmount),
		"LoanAmount":        float64(surrhenq.LoanAmount),
		"PolicyDepost":      float64(surrhenq.PolicyDepost),
		"CashDeposit":       float64(surrhenq.CashDeposit),
		"RefundPrem":        float64(surrhenq.RefundPrem),
		"PremTolerance":     float64(surrhenq.PremTolerance),
		"TotalSurrPayable":  float64(surrhenq.TotalSurrPayable),
		"AdjustedAmount":    float64(surrhenq.AdjustedAmount),
		"ReasonDescription": surrhenq.ReasonDescription,
	}
	surrharray = append(surrharray, resultOut)

	var surrdenq []models.SurrD

	initializers.DB.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	surrdarray := make([]interface{}, 0)

	for k := 0; k < len(surrdenq); k++ {
		resultOut := map[string]interface{}{
			"ID":              IDtoPrint(surrdenq[k].ID),
			"PolicyID":        IDtoPrint(surrdenq[k].PolicyID),
			"ClientID":        IDtoPrint(surrdenq[k].ClientID),
			"BenefitID":       IDtoPrint(surrdenq[k].ID),
			"BCoverage":       surrdenq[k].BCoverage,
			"BSumAssured":     surrdenq[k].BSumAssured,
			"SurrAmount":      float64(surrdenq[k].SurrAmount),
			"RevBonus":        float64(surrdenq[k].RevBonus),
			"AddlBonus":       float64(surrdenq[k].AddlBonus),
			"TerminalBonus":   float64(surrdenq[k].TerminalBonus),
			"InterimBonus":    float64(surrdenq[k].InterimBonus),
			"LoyaltyBonus":    float64(surrdenq[k].LoyaltyBonus),
			"OtherAmount":     float64(surrdenq[k].OtherAmount),
			"AccumDividend":   float64(surrdenq[k].AccumDividend),
			"AccumDivInt":     float64(surrdenq[k].AccumDivInt),
			"TotalFundValue":  float64(surrdenq[k].TotalFundValue),
			"TotalSurrAmount": float64(surrdenq[k].TotalSurrAmount),
		}
		surrdarray = append(surrdarray, resultOut)
	}

	surrcombinedvalue := map[string]interface{}{
		"surrhdata":  surrharray,
		"surrdarray": surrdarray,
	}

	return surrcombinedvalue

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

// #97
// GetNominee Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Nominee Details as an Interface
//
// ©  FuturaInsTech
func GetNomiData(iCompany uint, iPolicy uint) []interface{} {

	var nomenq []models.Nominee

	initializers.DB.Find(&nomenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	nomarray := make([]interface{}, 0)
	var clientenq models.Client
	var policyenq models.Policy
	initializers.DB.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	for k := 0; k < len(nomenq); k++ {
		initializers.DB.Find(&clientenq, "company_id = ? and id = ?", iCompany, nomenq[k].ClientID)
		rcd := policyenq.PRCD
		oAge, _, _, _, _, _ := DateDiff(String2Date(clientenq.ClientDob), String2Date(rcd), "")
		resultOut := map[string]interface{}{
			"ID":                  IDtoPrint(nomenq[k].ID),
			"PolicyID":            IDtoPrint(nomenq[k].PolicyID),
			"ClientID":            IDtoPrint(nomenq[k].ClientID),
			"NomineeRelationship": nomenq[k].NomineeRelationship,
			"LongName":            nomenq[k].NomineeLongName,
			"Percentage":          nomenq[k].NomineePercentage,
			"Age":                 oAge,
			"DateofBirth":         DateConvert(clientenq.ClientDob),
		}
		nomarray = append(nomarray, resultOut)
	}

	return nomarray

}

// # 98 (Redundant) Not in Use
// GetDeath Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Death Details as an Interface
//
// ©  FuturaInsTech
// Not Required
func GetDeathData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var surrhenq models.SurrH
	var surrdenq []models.SurrD
	initializers.DB.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	initializers.DB.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	surrarray := make([]interface{}, 0)

	return surrarray
}

// #98
// GetMatH Data - Printing Purpose Only (both header and detail)
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Maturity Header and Details Interface
// ©  FuturaInsTech
func GetMatHData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) interface{} {
	var mathenq models.MaturityH

	initializers.DB.Find(&mathenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	matharray := make([]interface{}, 0)
	_, oProduct, _ := GetParamDesc(iCompany, "Q0005", mathenq.Product, 1)
	resultOut := map[string]interface{}{

		"ID":                   IDtoPrint(mathenq.ID),
		"PolicyID":             IDtoPrint(mathenq.PolicyID),
		"ClientID":             IDtoPrint(mathenq.ClientID),
		"EffectiveDate":        DateConvert(mathenq.EffectiveDate),
		"MaturityDate":         DateConvert(mathenq.MaturityDate),
		"Status":               mathenq.Status,
		"BillDate":             DateConvert(mathenq.BillDate),
		"PaidToDate":           DateConvert(mathenq.PaidToDate),
		"Product":              oProduct,
		"AplAmount":            mathenq.AplAmount,
		"LoanAmount":           mathenq.LoanAmount,
		"PolicyDepost":         mathenq.PolicyDepost,
		"CashDeposit":          mathenq.CashDeposit,
		"RefundPrem":           mathenq.RefundPrem,
		"PremTolerance":        mathenq.PremTolerance,
		"TotalMaturityPayable": mathenq.TotalMaturityPayable,
		"AdjustedAmount":       mathenq.AdjustedAmount,
	}
	matharray = append(matharray, resultOut)

	var matdenq []models.MaturityD

	initializers.DB.Find(&matdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	matdarray := make([]interface{}, 0)

	for k := 0; k < len(matdenq); k++ {
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0005", matdenq[k].BCoverage, 1)
		resultOut := map[string]interface{}{
			"MaturityHID":         IDtoPrint(matdenq[k].MaturityHID),
			"PolicyID":            matdenq[k].PolicyID,
			"ClientID":            matdenq[k].ClientID,
			"BenifitID":           matdenq[k].BenefitID,
			"BCoverage":           oCoverage,
			"BSumAssured":         matdenq[k].BSumAssured,
			"MaturityAmount":      matdenq[k].MaturityAmount,
			"RevBonus":            matdenq[k].RevBonus,
			"AddlBonus":           matdenq[k].AddlBonus,
			"TerminalBonus":       matdenq[k].TerminalBonus,
			"InterimBonus":        matdenq[k].InterimBonus,
			"LoyaltyBonus":        matdenq[k].LoyaltyBonus,
			"OtherAmount":         matdenq[k].OtherAmount,
			"AccumDividend":       matdenq[k].AccumDividend,
			"AccumDivInt":         matdenq[k].AccumDivInt,
			"TotalFundValue":      matdenq[k].TotalFundValue,
			"TotalMaturityAmount": matdenq[k].TotalMaturityAmount,
		}
		matdarray = append(matdarray, resultOut)
	}

	matcombineddata := map[string]interface{}{
		"matharray": matdarray,
		"matdarray": matdarray,
	}
	return matcombineddata

}

// #99
// GetSurvBPay Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Survival Benefit Payment Interface
// ©  FuturaInsTech
func GetSurvBPay(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iTranno uint) []interface{} {
	var survbenq models.SurvB
	initializers.DB.Find(&survbenq, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno)
	survbparray := make([]interface{}, 0)
	resultOut := map[string]interface{}{
		"ID":           IDtoPrint(survbenq.ID),
		"Sequence":     IDtoPrint(uint(survbenq.Sequence)),
		"PolicyID":     IDtoPrint(survbenq.PolicyID),
		"BenefitID":    IDtoPrint(survbenq.BenefitID),
		"DueDate":      DateConvert(survbenq.EffectiveDate),
		"PaidDate":     DateConvert(survbenq.PaidDate),
		"SBPercentage": survbenq.SBPercentage,
		"SBAmount":     survbenq.Amount,
	}
	survbparray = append(survbparray, resultOut)
	return survbparray
}

// #100
// GetBonsusValues Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Bonus Values
// ©  FuturaInsTech
func GetBonusVals(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iTranno uint) []interface{} {

	bonusarray := make([]interface{}, 0)

	oPolicyDeposit := GetGlBal(iCompany, uint(iPolicy), "PolicyDeposit")
	oRevBonus := GetGlBal(iCompany, uint(iPolicy), "ReversionaryBonus")
	oTermBonus := GetGlBal(iCompany, uint(iPolicy), "TerminalBonus")
	oIntBonus := GetGlBal(iCompany, uint(iPolicy), "InterimBonus")
	oAccumDiv := GetGlBal(iCompany, uint(iPolicy), "AccumDividend")
	oAccumDivInt := GetGlBal(iCompany, uint(iPolicy), "AccumDivInt")
	oAddBonus := GetGlBal(iCompany, uint(iPolicy), "AdditionalBonus")
	oLoyalBonus := GetGlBal(iCompany, uint(iPolicy), "LoyaltyBonus")
	oAplAmt := GetGlBal(iCompany, uint(iPolicy), "AplAmount")
	oPolLoan := GetGlBal(iCompany, uint(iPolicy), "PolicyLoan")
	oCashDep := GetGlBal(iCompany, uint(iPolicy), "CashDeposit")

	resultOut := map[string]interface{}{
		"ID":            IDtoPrint(iPolicy),
		"PolicyDeposit": NumbertoPrint(oPolicyDeposit),
		"RevBonus":      NumbertoPrint(oRevBonus),
		"TermBonus":     NumbertoPrint(oTermBonus),
		"IntBonus":      NumbertoPrint(oIntBonus),
		"AccDividend":   NumbertoPrint(oAccumDiv),
		"AccDivInt":     NumbertoPrint(oAccumDivInt),
		"AddBonus":      NumbertoPrint(oAddBonus),
		"LoyalBonus":    NumbertoPrint(oLoyalBonus),
		"AplAmount":     NumbertoPrint(oAplAmt),
		"PolLoan":       NumbertoPrint(oPolLoan),
		"CashDeposit":   NumbertoPrint(oCashDep),
	}
	bonusarray = append(bonusarray, resultOut)

	return bonusarray
}

// #101
// GetAgency Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt, Agency and Date
//
// # Outputs  Agency
// ©  FuturaInsTech
func GetAgency(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iTranno uint, iAgency uint) []interface{} {

	agencyarray := make([]interface{}, 0)
	var agencyenq models.Agency
	var clientenq models.Client
	initializers.DB.Find(&agencyenq, "company_id  = ? and id = ?", iCompany, iAgency)

	initializers.DB.Find(&clientenq, "company_id = ? and id = ?", iCompany, agencyenq.ClientID)
	oAgentName := clientenq.ClientLongName + " " + clientenq.ClientShortName + " " + clientenq.ClientSurName

	var addressenq models.Address
	initializers.DB.Find(&addressenq, "company_id = ? and client_id = ?", iCompany, clientenq.ID)
	oAddress := addressenq.AddressLine1 + "" + addressenq.AddressLine2 + "" + addressenq.AddressLine3 + "" + addressenq.AddressLine4 + "" + addressenq.AddressLine5 + "" + addressenq.AddressPostCode + "" + addressenq.AddressState

	resultOut := map[string]interface{}{
		"ID":              IDtoPrint(iAgency),
		"AgyChannelSt":    agencyenq.AgencyChannel,
		"AgyStatus":       agencyenq.AgencySt,
		"AgyBankId":       agencyenq.BankID,
		"AgyClientNo":     agencyenq.ClientID,
		"AgyLicNo":        agencyenq.LicenseNo,
		"AgyLicEndDate":   agencyenq.LicenseEndDate,
		"AgyLicStartDate": agencyenq.LicenseStartDate,
		"AgyOffice":       agencyenq.Office,
		"AgyTermReason":   agencyenq.TerminationReason,
		"AgyEndDate":      agencyenq.EndDate,
		"AgentName":       oAgentName,
		"oAddress":        oAddress,
		"PhoneNo":         clientenq.ClientMobile,
		"Email":           clientenq.ClientEmail,
	}
	agencyarray = append(agencyarray, resultOut)

	return agencyarray
}

// #102
// GetExpi Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Expiry Interface Information
// ©  FuturaInsTech
func GetExpi(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iTranno uint) []interface{} {
	var benefit []models.Benefit
	initializers.DB.Find(&benefit, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno)
	expiryarray := make([]interface{}, 0)

	for k := 0; k < len(benefit); k++ {
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", benefit[k].BCoverage, 1)
		resultOut := map[string]interface{}{
			"ID":             IDtoPrint(benefit[k].ID),
			"CompanyID":      IDtoPrint(benefit[k].CompanyID),
			"ClientID":       IDtoPrint(benefit[k].ClientID),
			"PolicyID":       IDtoPrint(benefit[k].PolicyID),
			"BStartDate":     DateConvert(benefit[k].BStartDate),
			"BRiskCessDate":  DateConvert(benefit[k].BRiskCessDate),
			"BPremCessDate":  DateConvert(benefit[k].BPremCessDate),
			"BTerm":          benefit[k].BTerm,
			"BPTerm":         benefit[k].BPTerm,
			"BRiskCessAge":   benefit[k].BRiskCessAge,
			"BPremCessAge":   benefit[k].BPremCessAge,
			"BBasAnnualPrem": NumbertoPrint(benefit[k].BBasAnnualPrem),
			"BLoadPrem":      NumbertoPrint(benefit[k].BLoadPrem),
			"BCoverage":      oCoverage,
			"BSumAssured":    NumbertoPrint(float64(benefit[k].BLoadPrem)),
			"BPrem":          NumbertoPrint(benefit[k].BPrem),
			"BGender":        benefit[k].BGender,
			"BDOB":           benefit[k].BDOB,
			"BMortality":     benefit[k].BMortality,
			"BStatus":        benefit[k].BStatus,
			"BAge":           benefit[k].BAge,
			"BRerate":        benefit[k].BRerate,
		}
		expiryarray = append(expiryarray, resultOut)
	}
	return expiryarray
}

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

// #104
// Create Communication
//
// # This function, Create Communication Records by getting input values as Company ID, History Code, Tranno, Date of Transaction, Policy Id, Client Id, Address Id, Receipt ID . Quotation ID, Agency ID
// 10 Input Variables
// # It returns success or failure.  Successful records written in Communciaiton Table
//
// ©  FuturaInsTech
func CreateCommunications(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string) error {

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	//utilities.LetterCreate(int(iCompany), uint(iPolicy), iHistoryCode, createreceipt.CurrentDate, idata)
	iTransaction := iHistoryCode
	var policy models.Policy

	result := initializers.DB.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.Error != nil {
		return result.Error
	}

	iKey := iTransaction + policy.PProduct
	err1 := GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)

	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
		if err1 != nil {
			return err1
		}
	}

	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			err := GetItemD(int(iCompany), "P0033", iKey, iDate, &extradatap0033)
			if err != nil {
				return err
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
			oLetType := ""

			signData := make([]interface{}, 0)
			resultOut := map[string]interface{}{
				"Department":     p0033data.DepartmentName,
				"DepartmentHead": p0033data.DepartmentHead,
				"CoEmail":        p0033data.CompanyEmail,
				"CoPhone":        p0033data.CompanyPhone,
			}

			signData = append(signData, resultOut)

			batchData := make([]interface{}, 0)
			resultOut = map[string]interface{}{
				"Date":     DateConvert(iDate),
				"FromDate": DateConvert(iFromDate),
				"ToDate":   DateConvert(iToDate),
			}

			batchData = append(batchData, resultOut)

			resultMap := make(map[string]interface{})

			//	iCompany uint, iPolicy uint, iAddress uint, iClient uint, iLanguage uint, iBankcode uint, iReceipt uint, iCommunciation uint, iQuotation uint
			for n := 0; n < len(p0034data.Letters[i].LetType); n++ {
				oLetType = p0034data.Letters[i].LetType[n]
				switch {
				case oLetType == "1":
					oData := GetCompanyData(iCompany, iPolicy, iClient, iAddress, iReceipt, iDate)
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData := GetClientData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData := GetAddressData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData := GetPolicyData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData := GetBenefitData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData := GetSurBData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData := GetMrtaData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData := GetReceiptData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData := GetSaChangeData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData := GetCompAddData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData := GetSurrHData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData := GetDeathData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData := GetMatHData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData := GetSurvBPay(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno)
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData := GetExpi(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno)
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData := GetBonusVals(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno)
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData := GetAgency(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno, iAgency)
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData := GetNomiData(iCompany, iPolicy)
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData := GetGLData(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign)
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData := GetIlpSummaryData(iCompany, iPolicy)
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData := GetIlpAnnsummaryData(iCompany, iPolicy, iHistoryCode)
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData := GetIlpTranctionData(iCompany, iPolicy, iHistoryCode, iToDate)
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData := GetPremTaxGLData(iCompany, iPolicy, iFromDate, iToDate)
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData := GetIlpFundSwitchData(iCompany, iPolicy, iTranno)
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData := GetPHistoryData(iCompany, iPolicy, iHistoryCode, iDate)
					resultMap["PolicyHistoryData"] = oData

				case oLetType == "98":
					resultMap["BatchData"] = batchData

				case oLetType == "99":
					resultMap["SignData"] = signData
				default:

				}
			}

			communication.ExtractedData = resultMap
			communication.PDFPath = p0034data.Letters[i].PdfLocation
			communication.TemplatePath = p0034data.Letters[i].ReportTemplateLocation

			results := initializers.DB.Create(&communication)

			if results.Error != nil {
				return results.Error
			}

		}
	}
	return nil
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

	if result != nil {
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

	num := int(camt)
	dec := int((camt - float64(num)) * 100)

	// Process the number in  trillion, billion, million, thousand, and units

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
	if trillionWords := HundredsInWords(trillions, ones, tens); len(trillionWords) > 0 {
		words += trillionWords + " Trillion "
	}

	// Convert billions to words
	if billionWords := HundredsInWords(billions, ones, tens); len(billionWords) > 0 {
		words += billionWords + " Billion "
	}

	// Convert millions to words
	if millionWords := HundredsInWords(millions, ones, tens); len(millionWords) > 0 {
		words += millionWords + " Million "
	}

	// Convert thousands to words
	if thousandWords := HundredsInWords(thousands, ones, tens); len(thousandWords) > 0 {
		words += thousandWords + " Thousand "
	}

	// Convert units to words
	if unitWords := HundredsInWords(units, ones, tens); len(unitWords) > 0 {
		words += unitWords
	}

	words += " " + cname

	// Convert dec to decwords

	if decWords := HundredsInWords(dec, ones, tens); len(decWords) > 0 {
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
	receiptupd.PolicyID = iPolicy
	receiptupd.InstalmentPremium = policyenq.InstalmentPrem
	receiptupd.PaidToDate = policyenq.PaidToDate
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

	err = CreateCommunications(iCompany, iMethod, uint(iTranno), iBusinssdate, iPolicy, receiptupd.ClientID, receiptupd.AddressID, receiptupd.ID, 0, iAgency, "", "", "", "", "")
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
// TDFExpidD - Time Driven Function - Expiry Date Updation
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
func PostUlpDeduction(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint) error {

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
		ilptrancrt.FundAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
		ilptrancrt.FundCurr = p0061data.FundCurr
		ilptrancrt.FundUnits = 0
		ilptrancrt.FundPrice = 0
		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
		ilptrancrt.OriginalAmount = RoundFloat(((iAmount * iFundValue) / iTotalFundValue), 2)
		ilptrancrt.ContractCurry = policyenq.PContractCurr
		ilptrancrt.HistoryCode = iHistoryCode
		ilptrancrt.InvNonInvFlag = "AC"
		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
		ilptrancrt.InvNonInvPercentage = RoundFloat((iFundValue / iTotalFundValue), 5)
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

		if isFieldZero(fv) == true {
			shortCode := p0065data.FieldList[i].ErrorCode
			longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return errors.New(shortCode + " : " + longDesc)
		}

	}

	// if clientval.ClientEmail == "" || !strings.Contains(clientval.ClientEmail, "@") || !strings.Contains(clientval.ClientEmail, ".") {
	// 	shortCode := "GL477"
	// 	longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
	// 	return errors.New(shortCode + " : " + longDesc)
	// }

	validemail := isValidEmail(clientval.ClientEmail)
	if !validemail {
		shortCode := "GL477"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)
	}

	_, err = strconv.Atoi(clientval.ClientMobile)
	if err != nil {
		shortCode := "GL478"
		longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return errors.New(shortCode + " : " + longDesc)
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

// # 139
// Get ILP Summary Data Printing Purpose Only
// Inputs: Company and Policy
//
// Outputs : Summary Data as interface
//
// ©  FuturaInsTech
func GetIlpSummaryData(iCompany uint, iPolicy uint) interface{} {
	var ilpsummary []models.IlpSummary
	initializers.DB.Find(&ilpsummary, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	ilpsummaryarray := make([]interface{}, 0)
	ilpsumtotfundvalue := make([]interface{}, 0)
	bpfv, opfv, _ := GetAllFundValueByPol(iCompany, iPolicy, "")

	resultOut := map[string]interface{}{
		"BpFundValue": RoundFloat(bpfv, 5),
		"OpFundValue": RoundFloat(opfv, 5),
	}
	ilpsumtotfundvalue = append(ilpsumtotfundvalue, resultOut)

	for k := 0; k < len(ilpsummary); k++ {
		resultOut := map[string]interface{}{
			"ID":        IDtoPrint(ilpsummary[k].ID),
			"CompanyID": IDtoPrint(ilpsummary[k].CompanyID),
			"BenefitID": IDtoPrint(ilpsummary[k].BenefitID),
			"PolicyID":  IDtoPrint(ilpsummary[k].PolicyID),
			"FundCode":  ilpsummary[k].FundCode,
			"FundType":  ilpsummary[k].FundType,
			"FundUnits": ilpsummary[k].FundUnits,
			"FundCurr":  ilpsummary[k].FundCurr,
		}

		ilpsummaryarray = append(ilpsummaryarray, resultOut)
	}
	ilpsumcombinedvalue := map[string]interface{}{
		"TotalFundValues": ilpsumtotfundvalue,
		"IlpSummaryFunds": ilpsummaryarray,
	}

	return ilpsumcombinedvalue
}

// # 140
// Get GL Data Printing Purpose Only
// Inputs: Company  Policy, From and To Date, History Code, GL Code,GL Sign
//
// Outputs : GL Data as interface
//
// ©  FuturaInsTech
func GetGLData(iCompany uint, iPolicy uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string) interface{} {
	var benefitenq []models.Benefit

	var covrcodes []string
	var covrnames []string

	initializers.DB.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	for i := 0; i < len(benefitenq); i++ {
		covrcode := benefitenq[i].BCoverage
		_, covrname, err := GetParamDesc(iCompany, "Q0006", covrcode, 1)
		if err != nil {
			continue
		}
		covrcodes = append(covrcodes, covrcode)
		covrnames = append(covrnames, covrname)
	}

	var glmoves []models.GlMove
	if iGlHistoryCode == "" && iGlAccountCode == "" && iGlSign == "" {
		initializers.DB.Where("company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date <=?", iCompany, iPolicy, iFromDate, iToDate).Order("effective_date , tranno").Find(&glmoves)
	} else if iGlHistoryCode != "" && iGlAccountCode == "" && iGlSign == "" {
		initializers.DB.Where("company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode).Order("history_code, effective_date , tranno").Find(&glmoves)
	} else if iGlHistoryCode != "" && iGlAccountCode != "" && iGlSign == "" {
		initializers.DB.Where("company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ? and account_code like ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, "%"+iGlAccountCode+"%").Order("history_code, account_code, effective_date , tranno").Find(&glmoves)
	} else if iGlHistoryCode != "" && iGlAccountCode != "" && iGlSign != "" {
		initializers.DB.Where("company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ? and account_code like ? and gl_sign = ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, "%"+iGlAccountCode+"%", iGlSign).Order("history_code, account_code, gl_sign, effective_date , tranno").Find(&glmoves)
	} else if iGlHistoryCode == "" && iGlAccountCode != "" && iGlSign != "" {
		initializers.DB.Where("company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and account_code like ? and gl_sign = ?", iCompany, iPolicy, iFromDate, iToDate, "%"+iGlAccountCode+"%", iGlSign).Order("account_code, gl_sign, effective_date , tranno").Find(&glmoves)
	} else if iGlHistoryCode == "" && iGlAccountCode != "" && iGlSign == "" {
		initializers.DB.Where("company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and account_code like ?", iCompany, iPolicy, iFromDate, iToDate, "%"+iGlAccountCode+"%").Order("account_code, effective_date , tranno").Find(&glmoves)
	}

	glsumtotarray := make([]interface{}, 0)
	glaccountcode := ""
	glaccounttotal := 0.0
	glcoveragecode := ""
	glcoveragename := ""
	for k := 0; k < len(glmoves); k++ {
		if k == 0 {
			glcoveragecode = glmoves[k].BCoverage
			glaccountcode = glmoves[k].AccountCode
			glaccounttotal = glaccounttotal + glmoves[k].ContractAmount
			continue
		}

		if glmoves[k].BCoverage == glcoveragecode && glmoves[k].AccountCode == glaccountcode {
			glaccounttotal = glaccounttotal + glmoves[k].ContractAmount
			continue
		} else {
			for i := 0; i < len(covrcodes); i++ {
				if glcoveragecode == covrcodes[i] {
					glcoveragename = covrnames[i]
					break
				}
			}
			resultOut := map[string]interface{}{
				"Glcoveragecode": glcoveragecode,
				"Glcoveragename": glcoveragename,
				"AccountCode":    glaccountcode,
				"GlAccountTotal": NumbertoPrint(glaccounttotal),
			}
			glsumtotarray = append(glsumtotarray, resultOut)
			// process the first record of next account code
			glaccounttotal = 0.0
			glaccountcode = glmoves[k].AccountCode
			glcoveragecode = glmoves[k].BCoverage
			glaccounttotal = glaccounttotal + glmoves[k].ContractAmount
			continue
		}
	}
	for i := 0; i < len(covrcodes); i++ {
		if glcoveragecode == covrcodes[i] {
			glcoveragename = covrnames[i]
			break
		}
	}
	resultOut := map[string]interface{}{
		"GlCoverageCode": glcoveragecode,
		"GlCoverageName": glcoveragename,
		"AccountCode":    glaccountcode,
		"GlAccountTotal": NumbertoPrint(glaccounttotal),
	}
	glsumtotarray = append(glsumtotarray, resultOut)

	glarray := make([]interface{}, 0)
	for k := 0; k < len(glmoves); k++ {
		glcoveragecode = glmoves[k].BCoverage
		for i := 0; i < len(covrcodes); i++ {
			if glcoveragecode == covrcodes[i] {
				glcoveragename = covrnames[i]
				break
			}
		}
		resultOut := map[string]interface{}{

			"GlRdocno":       glmoves[k].GlRdocno,
			"GlRldgAcct":     glmoves[k].GlRldgAcct,
			"GlCoverageCode": glcoveragecode,
			"GlCoverageName": glcoveragename,
			"GlCurry":        glmoves[k].GlCurry,
			"GlAmount":       NumbertoPrint(glmoves[k].GlAmount),
			"ContractCurry":  glmoves[k].ContractCurry,
			"ContractAmount": NumbertoPrint(glmoves[k].ContractAmount),
			"AccountCodeID":  glmoves[k].AccountCodeID,
			"AccountCode":    glmoves[k].AccountCode,
			"GlSign":         glmoves[k].GlSign,
			"SequenceNo":     glmoves[k].SequenceNo,
			"CurrencyRate":   NumbertoPrint(glmoves[k].CurrencyRate),
			"CurrentDate":    DateConvert(glmoves[k].CurrentDate),
			"EffectiveDate":  DateConvert(glmoves[k].EffectiveDate),
			"ReconciledDate": DateConvert(glmoves[k].ReconciledDate),
			"ExtractedDate":  DateConvert(glmoves[k].ExtractedDate),
			"HistoryCode":    glmoves[k].HistoryCode,
			"ReversalInd":    glmoves[k].ReversalIndicator,
		}
		glarray = append(glarray, resultOut)
	}

	glcombineddata := map[string]interface{}{
		"glsumtotarray": glsumtotarray,
		"glarray":       glarray,
	}
	return glcombineddata

}

// # 141
// GetPremTaxGLData
// Extract PremTax Data  (Printing Purpose Only)
// Input:  Company, Policy, From and To Date
// Output: Interface
//
// ©  FuturaInsTech
func GetPremTaxGLData(iCompany uint, iPolicy uint, iFromDate string, iToDate string) interface{} {
	var benefitenq []models.Benefit
	var codesql string = ""
	var covrcodes []string
	var covrnames []string

	var acodearray []string

	var p0067data paramTypes.P0067Data
	var extradatap0067 paramTypes.Extradata = &p0067data

	initializers.DB.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	for i := 0; i < len(benefitenq); i++ {
		covrcode := benefitenq[i].BCoverage
		_, covrname, err := GetParamDesc(iCompany, "Q0006", covrcode, 1)
		if err != nil {
			continue
		}
		covrcodes = append(covrcodes, covrcode)
		covrnames = append(covrnames, covrname)

		err = GetItemD(int(iCompany), "P0067", benefitenq[i].BCoverage, iFromDate, &extradatap0067)
		if err != nil {
			return nil
		}

		notFound := true
		for j := 0; j < len(p0067data.GlTax); j++ {
			for _, str := range acodearray {
				if str == p0067data.GlTax[j].AccountCode {
					notFound = false
					break
				}

			}
			if notFound {
				acodearray = append(acodearray, p0067data.GlTax[j].AccountCode)
			}

		}
	}

	for k := 0; k < len(acodearray); k++ {
		if k == 0 {
			codesql = " account_code like '%" + acodearray[k] + "%' "
		} else {
			codesql = codesql + " or account_code like '%" + acodearray[k] + "%' "
		}
	}
	var glmoves []models.GlMove
	initializers.DB.Where("("+codesql+") and company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date <=? ", iCompany, iPolicy, iFromDate, iToDate).Order("account_code, gl_sign, effective_date , tranno").Find(&glmoves)

	glsumtotarray := make([]interface{}, 0)
	glaccounttotal := 0.0
	glcoveragecode := ""
	glcoveragename := ""
	glaccountcode := ""
	gltaxsection := ""
	for k := 0; k < len(glmoves); k++ {
		if k == 0 {
			glcoveragecode = glmoves[k].BCoverage
			glaccountcode = glmoves[k].AccountCode
			glaccounttotal = glaccounttotal + glmoves[k].ContractAmount
			continue
		}
		if glmoves[k].BCoverage == glcoveragecode && glmoves[k].AccountCode == glaccountcode {
			glaccounttotal = glaccounttotal + glmoves[k].ContractAmount
			continue
		} else {
			for i := 0; i < len(covrcodes); i++ {
				if glcoveragecode == covrcodes[i] {
					glcoveragename = covrnames[i]
					break
				}
			}
			//gltaxsection
			err := GetItemD(int(iCompany), "P0067", glcoveragecode, iFromDate, &extradatap0067)
			if err != nil {
				return nil
			}
			for j := 0; j < len(p0067data.GlTax); j++ {
				glcodewithcover := p0067data.GlTax[j].AccountCode + glcoveragecode
				if glaccountcode == p0067data.GlTax[j].AccountCode ||
					glaccountcode == glcodewithcover {
					gltaxsection = p0067data.GlTax[j].TaxSection
					break
				}
			}
		}
		resultOut := map[string]interface{}{
			"GlCoverageCode": glcoveragecode,
			"GlCoverageName": glcoveragename,
			"GlAccountCode":  glaccountcode,
			"GlAccountTotal": NumbertoPrint(glaccounttotal),
			"GlTaxSection":   gltaxsection,
		}
		glsumtotarray = append(glsumtotarray, resultOut)
		// process the first record of next account code
		glaccounttotal = 0.0
		glaccountcode = glmoves[k].AccountCode
		glcoveragecode = glmoves[k].BCoverage
		glaccounttotal = glaccounttotal + glmoves[k].ContractAmount
		//gltaxsection
		err := GetItemD(int(iCompany), "P0067", glcoveragecode, iFromDate, &extradatap0067)
		if err != nil {
			return nil
		}
		for j := 0; j < len(p0067data.GlTax); j++ {
			glcodewithcover := p0067data.GlTax[j].AccountCode + glcoveragecode
			if glaccountcode == p0067data.GlTax[j].AccountCode ||
				glaccountcode == glcodewithcover {
				gltaxsection = p0067data.GlTax[j].TaxSection
				break
			}
		}
	}

	for i := 0; i < len(covrcodes); i++ {
		if glcoveragecode == covrcodes[i] {
			glcoveragename = covrnames[i]
			break
		}
	}
	//gltaxsection
	err := GetItemD(int(iCompany), "P0067", glcoveragecode, iFromDate, &extradatap0067)
	if err != nil {
		return nil
	}
	for j := 0; j < len(p0067data.GlTax); j++ {
		glcodewithcover := p0067data.GlTax[j].AccountCode + glcoveragecode
		if glaccountcode == p0067data.GlTax[j].AccountCode ||
			glaccountcode == glcodewithcover {
			gltaxsection = p0067data.GlTax[j].TaxSection
			break
		}
	}
	resultOut := map[string]interface{}{
		"GlCoverageCode": glcoveragecode,
		"GlCoverageName": glcoveragename,
		"GlAccountCode":  glaccountcode,
		"GlAccountTotal": NumbertoPrint(glaccounttotal),
		"GlTaxSection":   gltaxsection,
	}
	glsumtotarray = append(glsumtotarray, resultOut)

	glarray := make([]interface{}, 0)
	for k := 0; k < len(glmoves); k++ {
		glcoveragecode = glmoves[k].BCoverage
		for i := 0; i < len(covrcodes); i++ {
			if glcoveragecode == covrcodes[i] {
				glcoveragename = covrnames[i]
				break
			}
		}
		resultOut := map[string]interface{}{
			"GlRdocno":       glmoves[k].GlRdocno,
			"GlRldgAcct":     glmoves[k].GlRldgAcct,
			"GlCoverageCode": glmoves[k].BCoverage,
			"GlCoverageName": glcoveragename,
			"GlCurry":        glmoves[k].GlCurry,
			"GlAmount":       NumbertoPrint(glmoves[k].GlAmount),
			"ContractCurry":  glmoves[k].ContractCurry,
			"ContractAmount": NumbertoPrint(glmoves[k].ContractAmount),
			"AccountCodeID":  glmoves[k].AccountCodeID,
			"AccountCode":    glmoves[k].AccountCode,
			"GlSign":         glmoves[k].GlSign,
			"SequenceNo":     glmoves[k].SequenceNo,
			"CurrencyRate":   NumbertoPrint(glmoves[k].CurrencyRate),
			"CurrentDate":    DateConvert(glmoves[k].CurrentDate),
			"EffectiveDate":  DateConvert(glmoves[k].EffectiveDate),
			"ReconciledDate": DateConvert(glmoves[k].ReconciledDate),
			"ExtractedDate":  DateConvert(glmoves[k].ExtractedDate),
			"HistoryCode":    glmoves[k].HistoryCode,
			"ReversalInd":    glmoves[k].ReversalIndicator,
		}
		glarray = append(glarray, resultOut)
	}

	glcombineddata := map[string]interface{}{
		"glsumtotarray": glsumtotarray,
		"glarray":       glarray,
	}
	return glcombineddata

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

// # 146
//
// # GetIlpTranctionData - ILP transaction Data extraction for Communications
//
// ©  FuturaInsTech
func GetIlpTranctionData(iCompany uint, iPolicy uint, iHistoryCode string, iDate string) []interface{} {
	var policyenq models.Policy
	initializers.DB.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	iAnnivDate := Date2String(GetNextDue(policyenq.AnnivDate, "Y", "R"))
	iPrevAnnivDate := Date2String(GetNextDue(iAnnivDate, "Y", "R"))
	var ilptranction []models.IlpTransaction
	if iHistoryCode == "B0103" {
		initializers.DB.Where("company_id = ? and policy_id = ? and ul_process_flag = ? and inv_non_inv_flag != ? and transaction_date >= ? and transaction_date < ?", iCompany, iPolicy, "C", "NI", iPrevAnnivDate, iAnnivDate).Order("fund_code, transaction_date , tranno").Find(&ilptranction)
	} else if iHistoryCode == "B0115" {
		initializers.DB.Where("company_id = ? and policy_id = ? and ul_process_flag = ? and inv_non_inv_flag != ? and transaction_date >= ? and transaction_date <= ?", iCompany, iPolicy, "C", "NI", iAnnivDate, iDate).Order("fund_code, transaction_date , tranno").Find(&ilptranction)
	}

	ilptranctionarray := make([]interface{}, 0)

	for k := 0; k < len(ilptranction); k++ {
		resultOut := map[string]interface{}{
			"PolicyID": ilptranction[k].PolicyID,
			//"BenefitID":           ilptranction[k].BenefitID,
			"FundCode":        ilptranction[k].FundCode,
			"FundType":        ilptranction[k].FundType,
			"TransactionDate": DateConvert(ilptranction[k].TransactionDate),
			"FundEffDate":     DateConvert(ilptranction[k].FundEffDate),
			"FundAmount":      ilptranction[k].FundAmount,
			//"FundCurr":            ilptranction[k].FundCurr,
			"FundUnits": ilptranction[k].FundUnits,
			"FundPrice": NumbertoPrint(float64(ilptranction[k].FundPrice)),
			//"CurrentOrFuture":     ilptranction[k].CurrentOrFuture,
			"OriginalAmount": ilptranction[k].OriginalAmount,
			"ContractCurry":  ilptranction[k].ContractCurry,
			//"HistoryCode":         ilptranction[k].HistoryCode,
			//"InvNonInvFlag":       ilptranction[k].InvNonInvFlag,
			"InvNonInvPercentage": ilptranction[k].InvNonInvPercentage,
			"AccountCode":         ilptranction[k].AccountCode,
			//"CurrencyRate":        ilptranction[k].CurrencyRate,
			//"MortalityIndicator":  ilptranction[k].MortalityIndicator,
			//"SurrenderPercentage": ilptranction[k].SurrenderPercentage,
			//"Seqno":               ilptranction[k].Seqno,
			//"UlProcessFlag":       ilptranction[k].UlProcessFlag,
			"UlpPriceDate":       DateConvert(ilptranction[k].UlpPriceDate),
			"AllocationCategory": ilptranction[k].AllocationCategory,
			//"AdjustedDate":        DateConvert(ilptranction[k].AdjustedDate),
			"ID": IDtoPrint(ilptranction[k].ID),
		}

		ilptranctionarray = append(ilptranctionarray, resultOut)
	}
	return ilptranctionarray
}

// # 147
//
// # GetIlpAnnsummaryData - ILP Anniversary Summary Data extraction for Communications
//
// ©  FuturaInsTech
func GetIlpAnnsummaryData(iCompany uint, iPolicy uint, iHistoryCode string) interface{} {
	ilpannsumprevarray := make([]interface{}, 0)
	ilpannsumcurrarray := make([]interface{}, 0)
	var policyenq models.Policy
	initializers.DB.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	iAnnivDate := Date2String(GetNextDue(policyenq.AnnivDate, "Y", "R"))
	iPrevAnnivDate := Date2String(GetNextDue(iAnnivDate, "Y", "R"))

	var ilpannsumprev []models.IlpAnnSummary
	initializers.DB.Find(&ilpannsumprev, "company_id = ? and policy_id = ? and effective_date = ?", iCompany, iPolicy, iPrevAnnivDate)

	for k := 0; k < len(ilpannsumprev); k++ {
		resultOut := map[string]interface{}{
			"ID":            IDtoPrint(ilpannsumprev[k].ID),
			"PolicyID":      IDtoPrint(ilpannsumprev[k].PolicyID),
			"BenefitID":     IDtoPrint(ilpannsumprev[k].BenefitID),
			"FundCode":      ilpannsumprev[k].FundCode,
			"FundType":      ilpannsumprev[k].FundType,
			"FundUnits":     ilpannsumprev[k].FundUnits,
			"FundCurr":      ilpannsumprev[k].FundCurr,
			"EffectiveDate": DateConvert(ilpannsumprev[k].EffectiveDate),
		}

		ilpannsumprevarray = append(ilpannsumprevarray, resultOut)
	}

	var ilpannsumcurr []models.IlpAnnSummary
	initializers.DB.Find(&ilpannsumcurr, "company_id = ? and policy_id = ? and effective_date = ?", iCompany, iPolicy, iAnnivDate)

	for k := 0; k < len(ilpannsumcurr); k++ {
		resultOut := map[string]interface{}{
			"ID":            IDtoPrint(ilpannsumcurr[k].ID),
			"PolicyID":      IDtoPrint(ilpannsumcurr[k].PolicyID),
			"BenefitID":     IDtoPrint(ilpannsumcurr[k].BenefitID),
			"FundCode":      ilpannsumcurr[k].FundCode,
			"FundType":      ilpannsumcurr[k].FundType,
			"FundUnits":     ilpannsumcurr[k].FundUnits,
			"FundCurr":      ilpannsumcurr[k].FundCurr,
			"EffectiveDate": DateConvert(ilpannsumcurr[k].EffectiveDate),
		}

		ilpannsumcurrarray = append(ilpannsumcurrarray, resultOut)
	}
	ilpannsumdata := map[string]interface{}{
		"IlpAnnSumPrevData": ilpannsumprevarray,
		"IlpAnnSumCurrData": ilpannsumcurrarray,
	}
	return ilpannsumdata
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

// #157
// GetPHistoryData  (Printing Purpose)
// Input: Company, Policy No, Transaction Code and Effective Date
// Output: An Interface Record (History Information)
//
// ©  FuturaInsTech
func GetPHistoryData(iCompany uint, iPolicy uint, iHistoryCode string, iDate string) []interface{} {
	var policyhistory []models.PHistory
	initializers.DB.Find(&policyhistory, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	policyhistoryarray := make([]interface{}, 0)

	for k := 0; k < len(policyhistory); k++ {
		resultOut := map[string]interface{}{
			"PolicyID":      policyhistory[k].PolicyID,
			"CompanyID":     policyhistory[k].CompanyID,
			"HistoryCode":   policyhistory[k].HistoryCode,
			"EffectiveDate": DateConvert(policyhistory[k].EffectiveDate),
			"CurrentDate":   DateConvert(policyhistory[k].CurrentDate),
			"PrevData":      policyhistory[k].PrevData,
			"ID":            IDtoPrint(policyhistory[k].ID),
		}

		policyhistoryarray = append(policyhistoryarray, resultOut)
	}
	return policyhistoryarray
}

// #158
// ILP Products Only.  GetIlpFundSwitchData (Printing Purpose)
// Input: Company, Policy No, Tranno
// Output: An Interface Record
//
// ©  FuturaInsTech
func GetIlpFundSwitchData(iCompany uint, iPolicy uint, iTranno uint) interface{} {
	ilpswitchfundarray := make([]interface{}, 0)
	ilpfundarray := make([]interface{}, 0)
	var policyenq models.Policy
	initializers.DB.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	var ilpswitchheader []models.IlpSwitchHeader

	initializers.DB.Where("company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno).Order("tranno").Find(&ilpswitchheader)

	for k := 0; k < len(ilpswitchheader); k++ {
		resultOut := map[string]interface{}{
			"PolicyID":      ilpswitchheader[k].PolicyID,
			"BenefitID":     ilpswitchheader[k].BenefitID,
			"CompanyID":     ilpswitchheader[k].CompanyID,
			"EffectiveDate": DateConvert(ilpswitchheader[k].EffectiveDate),
			"ID":            IDtoPrint(ilpswitchheader[k].ID),
		}

		ilpswitchfundarray = append(ilpswitchfundarray, resultOut)
	}

	var ilpswitchfund []models.IlpSwitchFund
	initializers.DB.Where(" policy_id = ? and tranno = ?", iPolicy, iTranno).Order("fund_code").Find(&ilpswitchfund)

	for j := 0; j < len(ilpswitchfund); j++ {
		resultOut := map[string]interface{}{
			"PolicyID":           ilpswitchfund[j].PolicyID,
			"BenefitID":          ilpswitchfund[j].BenefitID,
			"CompanyID":          ilpswitchfund[j].CompanyID,
			"EffectiveDate":      DateConvert(ilpswitchfund[j].EffectiveDate),
			"ID":                 IDtoPrint(ilpswitchfund[j].ID),
			"FundSwitchHeaderID": ilpswitchfund[j].IlpSwitchHeaderID,
			"SwitchDirection":    ilpswitchfund[j].SwitchDirection,
			"SequenceNo":         ilpswitchfund[j].SequenceNo,
			"FundCode":           ilpswitchfund[j].FundCode,
			"FundPercentage":     ilpswitchfund[j].FundPercentage,
			"FundUnits":          ilpswitchfund[j].FundUnits,
			"FundAmount":         ilpswitchfund[j].FundAmount,
			"FundType":           ilpswitchfund[j].FundType,
			"FundCurruncy":       ilpswitchfund[j].FundCurr,
			"FundPrice":          ilpswitchfund[j].FundPrice,
		}

		ilpfundarray = append(ilpfundarray, resultOut)
	}
	IlpFundSwitchData := map[string]interface{}{
		"SwitchHeader": ilpswitchfundarray,
		"SwitchFund":   ilpfundarray,
	}
	return IlpFundSwitchData

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
		if p0070data.FreeSwitches < uint(switchcount) {
			return nil, 0, 0
		} else {
			return nil, 0, p0070data.FeePercentage
		}
	}
	// Fixed Amount
	if p0070data.SwitchFeeBasis == "F" {
		if p0070data.FreeSwitches < uint(switchcount) {
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
	txn.Save(&glmove)
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
	initializers.DB.First(&company, "id = ?", iCompany)

	results := initializers.DB.First(&glbal, "company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno)
	if results.Error != nil {
		return errors.New("Account Code Not Found"), glbal.ContractAmount
	}
	if results.RowsAffected == 0 {
		glbal.ContractAmount = temp
		glbal.CompanyID = iCompany
		glbal.GlAccountno = iGlAccountCode
		glbal.GlRldgAcct = iGlRldgAcct
		glbal.ContractCurry = iContCurry
		glbal.GlRdocno = iGlRdocno
		//initializers.DB.Save(&glbal)
		txn.Save(&glbal)
		return nil, glbal.ContractAmount
	} else {
		iAmount := glbal.ContractAmount + temp
		// fmt.Println("I am inside update.....2", iAmount, glbal.ContractAmount)
		//initializers.DB.Model(&glbal).Where("company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno).Update("contract_amount", iAmount)
		txn.Model(&glbal).Where("company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno).Update("contract_amount", iAmount)
		return nil, glbal.ContractAmount
	}
	//results.Commit()

}

// #104
// Create Communication (New Version with Rollback)
//
// # This function, Create Communication Records by getting input values as Company ID, History Code, Tranno, Date of Transaction, Policy Id, Client Id, Address Id, Receipt ID . Quotation ID, Agency ID
// 10 Input Variables
// # It returns success or failure.  Successful records written in Communciaiton Table
//
// ©  FuturaInsTech
func CreateCommunicationsN(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB) error {

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	//utilities.LetterCreate(int(iCompany), uint(iPolicy), iHistoryCode, createreceipt.CurrentDate, idata)
	iTransaction := iHistoryCode
	var policy models.Policy

	result := initializers.DB.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.Error != nil {
		return result.Error
	}

	iKey := iTransaction + policy.PProduct
	err1 := GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)

	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
		if err1 != nil {
			return err1
		}
	}

	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			err := GetItemD(int(iCompany), "P0033", iKey, iDate, &extradatap0033)
			if err != nil {
				return err
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
			oLetType := ""

			signData := make([]interface{}, 0)
			resultOut := map[string]interface{}{
				"Department":     p0033data.DepartmentName,
				"DepartmentHead": p0033data.DepartmentHead,
				"CoEmail":        p0033data.CompanyEmail,
				"CoPhone":        p0033data.CompanyPhone,
			}

			signData = append(signData, resultOut)

			batchData := make([]interface{}, 0)
			resultOut = map[string]interface{}{
				"Date":     DateConvert(iDate),
				"FromDate": DateConvert(iFromDate),
				"ToDate":   DateConvert(iToDate),
			}

			batchData = append(batchData, resultOut)

			resultMap := make(map[string]interface{})

			//	iCompany uint, iPolicy uint, iAddress uint, iClient uint, iLanguage uint, iBankcode uint, iReceipt uint, iCommunciation uint, iQuotation uint
			for n := 0; n < len(p0034data.Letters[i].LetType); n++ {
				oLetType = p0034data.Letters[i].LetType[n]
				switch {
				case oLetType == "1":
					oData := GetCompanyData(iCompany, iPolicy, iClient, iAddress, iReceipt, iDate)
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData := GetClientData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData := GetAddressData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData := GetPolicyData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData := GetBenefitData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData := GetSurBData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData := GetMrtaData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData := GetReceiptData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData := GetSaChangeData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData := GetCompAddData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData := GetSurrHData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData := GetDeathData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData := GetMatHData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData := GetSurvBPay(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno)
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData := GetExpi(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno)
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData := GetBonusVals(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno)
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData := GetAgency(iCompany, iPolicy, iClient, iAddress, iReceipt, iTranno, iAgency)
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData := GetNomiData(iCompany, iPolicy)
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData := GetGLData(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign)
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData := GetIlpSummaryData(iCompany, iPolicy)
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData := GetIlpAnnsummaryData(iCompany, iPolicy, iHistoryCode)
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData := GetIlpTranctionData(iCompany, iPolicy, iHistoryCode, iToDate)
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData := GetPremTaxGLData(iCompany, iPolicy, iFromDate, iToDate)
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData := GetIlpFundSwitchData(iCompany, iPolicy, iTranno)
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData := GetPHistoryData(iCompany, iPolicy, iHistoryCode, iDate)
					resultMap["PolicyHistoryData"] = oData

				case oLetType == "98":
					resultMap["BatchData"] = batchData

				case oLetType == "99":
					resultMap["SignData"] = signData
				default:

				}
			}

			communication.ExtractedData = resultMap
			communication.PDFPath = p0034data.Letters[i].PdfLocation
			communication.TemplatePath = p0034data.Letters[i].ReportTemplateLocation

			results := initializers.DB.Create(&communication)

			if results.Error != nil {
				return results.Error
			}

		}
	}
	return nil
}

// #82
// TdfhUpdate - Time Driven Function - Update TDF Header File
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
			result = txn.Create(&tdfhupd)
			if result.Error != nil {
				return errors.New("Error")
			}
		} else {
			result = initializers.DB.Delete(&tdfhupd)
			var tdfhupd models.Tdfh
			tdfhupd.CompanyID = iCompany
			tdfhupd.PolicyID = iPolicy
			tdfhupd.EffectiveDate = iDate
			tdfhupd.ID = 0
			result = txn.Create(&tdfhupd)
			if result.Error != nil {
				return errors.New("Error")
			}
		}

	}
	return nil
}
