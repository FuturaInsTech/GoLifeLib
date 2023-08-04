package utilities

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/types"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

// *********************************************************************************************
//
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
	//curr_date := time.Date(2024, 2, 29, 0, 0, 0, 0, time.Now().Location())

	//oneDayLater := curr_date.AddDate(0, 0, 1)
	// if iFrequency == "Y" {
	// 	a := iDate.AddDate(1, 0, 0)
	// 	return a
	// }

	// if iFrequency == "H" {
	// 	a := iDate.AddDate(0, 6, 0)
	// 	return a
	// }

	// if iFrequency == "Q" {
	// 	a := iDate.AddDate(0, 3, 0)
	// 	return a
	// }
	// if iFrequency == "M" {
	// 	a := iDate.AddDate(0, 1, 0)
	// 	return a
	// }
	// if iFrequency == "S" {
	// 	a := iDate.AddDate(0, 0, 0)
	// 	return a
	// }
	// return

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

// Add Month should ended with month
func AddMonth(t time.Time, m int) time.Time {
	x := t.AddDate(0, m, 0)
	if d := x.Day(); d != t.Day() {
		return x.AddDate(0, 0, -d)
	}
	return x
}

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

// func RoundAmt(iAmount float64, iMethod string) (oAmount float64) {

// 	if iMethod == "D" {
// 		oAmount = math.Floor(iAmount*100) / 100
// 	}
// 	if iMethod == "U" {
// 		oAmount = math.Ceil(iAmount*100) / 100
// 	}
// 	return
// }

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

func SimpleInterest(iPrincipal, iInterest, iDays float64) (oInterest float64) {
	oInterest = iPrincipal * (iInterest / 100) * (iDays / 365)
	return oInterest
}

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

func GetItemD(iCompany int, iTable string, iItem string, iFrom string, data *types.Extradata) error {

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

// func GetItemD1(iCompany int, iTable string, iItem string, iFrom string, data *types.Extradata1) error {

// 	//var sourceMap map[string]interface{}
// 	var itemparam models.Param

// 	results := initializers.DB.Find(&itemparam, "company_id =? and name= ? and item = ? and rec_type = ? and ? between start_date  and  end_date", iCompany, iTable, iItem, "IT", iFrom)

// 	if results.Error == nil {
// 		(*data).ParseData(itemparam.Data)
// 		return nil
// 	} else {
// 		return errors.New(results.Error.Error())
// 	}
// }

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

// GetMrtaBen
//
// # Inputs Original Term, SA, Frequency default Y (reset), Year elapsed, Interest Rate
//
// # Output Decreased SA
//
// ©  FuturaInsTech
// func GetMrtaBen(iTerm float64, iSA float64, iFrequency string, iYr float64, iInterest float64) (oSA float64) {

// 	i := 1 + (iInterest / 100)
// 	p := math.Pow(i, iTerm)
// 	comp := iSA * p
// 	emi := RoundFloat(comp/iTerm, 2)
// 	decreasesa := emi * (iTerm - iYr)
// 	decreasesa = RoundFloat(decreasesa, 2)

// 	return decreasesa
// }

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

// GetAnnualRate - Get Annual Rate of the Coverage - No Model Discount/Staff Discount/SA/Prem Discount
//
// Inputs: Company, Coverage, Age (Attained Age), Gender(F/N/U), Term (2 Characters), Premium Method as
// PM001 - Term Based , PM002 Age Based, Mortality Clause "S" Smoker, "N" Non Smoker
//
// Outputs Annualized Premium as float (124.22)
//
// ©  FuturaInsTech

func GetAnnualRate(iCompany uint, iCoverage string, iAge uint, iGender string, iTerm uint, iPremTerm uint, iPremMethod string, iDate string, iMortality string) (float64, error) {

	var q0006data types.Q0006Data
	var extradata types.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)

	var q0010data types.Q0010Data
	var extradataq0010 types.Extradata = &q0010data
	var q0010key string
	var prem float64
	term := strconv.FormatUint(uint64(iTerm), 10)
	premTerm := strconv.FormatUint(uint64(iPremTerm), 10)
	//fmt.Println("****************", iCompany, iCoverage, iAge, iGender, iTerm, iPremMethod, iDate, iMortality)
	if q0006data.PremCalcType == "A" {
		if q0006data.PremiumMethod == "PM002" {
			q0010key = iCoverage + iGender
		}
	} else if q0006data.PremCalcType == "P" {
		if q0006data.PremiumMethod == "PM001" {
			q0010key = iCoverage + iGender + term + premTerm
		}
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
			prem = q0010data.Rates[i].Rate
		}
	}
	fmt.Println("************", iCompany, iAge, q0010key, iDate, prem)
	return prem, nil
}

// ValidateCoverageQ0011 - Rider is Allowed for Product or Not Validation
//
// Inputs: Company, Product, Coverage and Date String in YYYYMMDD
//
// Outputs Product Found or Not  "Y" Means Found "N" Means Not Found
//
// ©  FuturaInsTech
func ValidateCoverageQ0011(iCompany uint, iProduct, iCoverage, iDate string) string {

	fmt.Println("Coverages Q0011", iCompany, iProduct, iCoverage, iDate)
	var q0011data types.Q0011Data
	var extradataq0011 types.Extradata = &q0011data
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

// Levels
func CustomizedPreload(d *gorm.DB) *gorm.DB {
	return d.Preload("Levels", CustomizedPreload)
}

// ValidateQ0012 - Survival Benefit (Term Based)
//
// Inputs: Company, Coverage and Date String in YYYYMMDD
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func ValidateQ0012(iCompany uint, iCoverage string, iDate string) error {
	var q0012data types.Q0012Data
	var extradataq0012 types.Extradata = &q0012data

	err := GetItemD(int(iCompany), "Q0012", "AED1", iDate, &extradataq0012)

	if err != nil {
		return err

	}

	for i := 0; i < len(q0012data.SBRates); i++ {
		fmt.Println("Survival Benefits .......")
		fmt.Println(q0012data.SBRates[i].Term)
		fmt.Println(q0012data.SBRates[i].Percentage)
	}
	return nil
}

// ValidateQ0013 - Survival Benefit (Age Based)
//
// Inputs: Company, Coverage and Date String in YYYYMMDD
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func ValidateQ0013(iCompany uint, iCoverage string, iDate string) error {
	var q0013data types.Q0013Data
	var extradataq0013 types.Extradata = &q0013data

	err := GetItemD(int(iCompany), "Q0013", "AEDR", iDate, &extradataq0013)

	if err != nil {
		return err

	}
	fmt.Println(q0013data.SBRates[0].Percentage)
	for i := 0; i < len(q0013data.SBRates); i++ {
		fmt.Println("Survival Benefits .......")
		fmt.Println(q0013data.SBRates[i].Age)
		fmt.Println(q0013data.SBRates[i].Percentage)
	}
	return nil
}

// ValidateQ0013 - Survival Benefit (Age Based)
//
// Inputs: Company, Coverage and Date String in YYYYMMDD
//
// # Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech
func GetSBByYear(iCompany uint, iCoverage string, iDate string, iSA float64, iType string, iMethod string, iYear int, iAge int) float64 {

	if iType == "T" {
		var q0012data types.Q0012Data
		var extradataq0012 types.Extradata = &q0012data
		// fmt.Println("SB Parameters", iCompany, iType, iMethod, iYear, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0012", iMethod, iDate, &extradataq0012)

		if err != nil {
			return 0

		}
		// fmt.Println(q0012data.SBRates[0].Percentage)
		for i := 0; i < len(q0012data.SBRates); i++ {
			if iYear == int(q0012data.SBRates[i].Term) {
				oSB := q0012data.SBRates[i].Percentage * iSA
				return oSB
			}
		}
	}
	if iType == "A" {
		var q0013data types.Q0013Data
		var extradataq0013 types.Extradata = &q0013data
		fmt.Println("SB Parameters", iCompany, iType, iMethod, iAge, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0013", iMethod, iDate, &extradataq0013)
		fmt.Println("SB Parameters", iCompany, iCoverage, iDate)

		if err != nil {
			return 0

		}
		fmt.Println(q0013data.SBRates[0].Percentage)
		for i := 0; i < len(q0013data.SBRates); i++ {
			if iAge == int(q0013data.SBRates[i].Age) {
				oSB := q0013data.SBRates[i].Percentage * iSA
				return oSB
			}
		}
	}
	return 0
}

// GetBonus - Get Bonus for a Given Duration
//
// Inputs: Company, Bonus Method, Status, Coverage Start Date, Year of Policy, Policy Status, SA
//
// Date in YYYYMMDD as a string
//
// Outputs SB Term and SB Percentage
//
// ©  FuturaInsTech

func GetBonus(iCompany uint, iCoverage string, iStartDate string, iEndDate string, iStatus string, iTerm uint, iSA uint) uint {

	var q0006data types.Q0006Data
	var extradata types.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iStartDate, &extradata)

	iRBMethod := q0006data.RevBonus
	// iIBMethod := q0006data.IBonus
	// iTBMethod := q0006data.TBonus
	// iLBMethod := q0006data.LoyaltyBonus
	// iSSVMethod := q0006data.SSVMethod
	// iGSVMethod := q0006data.GSVMethod
	// iBSVMethod := q0006data.BSVMethod
	var q0014data types.Q0014Data
	var extradata1 types.Extradata = &q0014data

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
	var q0006data types.Q0006Data
	var extradata types.Extradata = &q0006data
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
		key1 = q0006data.SSVMethod
	} else if iBonusMethod == "GSV" {
		key1 = q0006data.GSVMethod
	} else if iBonusMethod == "BSV" {
		key1 = q0006data.BSVMethod
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
	var q0014data types.Q0014Data
	var extradata1 types.Extradata = &q0014data
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

func GetTerm(iCompany uint, iCoverage string, iDate string) {
	var q0015data types.Q0015Data
	var extradata types.Extradata = &q0015data
	iKey := iCoverage
	fmt.Println(iKey)

	GetItemD(int(iCompany), "Q0015", iKey, iDate, &extradata)
	for i := 0; i < len(q0015data.Terms); i++ {
		term := q0015data.Terms[i].Term
		fmt.Println(term)

	}
	return
}
func GetPTerm(iCompany uint, iCoverage string, iDate string) {
	var q0016data types.Q0016Data
	var extradata types.Extradata = &q0016data
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
		var q0017data types.Q0017Data
		var extradataq0017 types.Extradata = &q0017data
		err := GetItemD(int(iCompany), "Q0017", iDiscMethod, iDate, &extradataq0017)

		if err != nil {
			return 0

		}

		for i := 0; i < len(q0017data.SABand); i++ {
			if int(iSA) <= int(q0017data.SABand[i].SA) {
				oDiscount := uint(q0017data.SABand[i].Discount) * uint(iAnnPrem) / 100
				return float64(oDiscount)
			}
		}
	}
	// Premium Discount
	if iDiscType == "P" {
		var q0018data types.Q0018Data
		var extradataq0018 types.Extradata = &q0018data

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

// CalcFrequencyPrem - Calculate Frequency Premium as per Model Factor Provided
//
// Inputs: Company, Frequency Factor Method as mentioned in Q0006, Current Frequency, Annualized Premium of the Coverage
//
// Output Model Premium =  Model Factor * Annualized Premium.
//
// ©  FuturaInsTech
func CalcFrequencyPrem(iCompany uint, iDate, iFreqMethod string, iFreq string, iAnnPrem float64) float64 {
	var q0019data types.Q0019Data
	var extradataq0019 types.Extradata = &q0019data
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

// GetWaiverSA - Calculate Waiver SA of a Policy
//
// Inputs: Company, All Coverages under the policy , Waiver Method as per Q0006, Waiver Coverage Start Date, Premium of the Current Coverage
//
// First Check whether given coverage is available in Q0020 for the Waiver Method.
// Foud, Add it in output and return SA
//
// ©  FuturaInsTech
func GetWaiverSA(iCompany uint, iCoverage string, iMethod string, iDate string, iPrem float64) float64 {

	var q0020data types.Q0020Data
	var extradataq0020 types.Extradata = &q0020data
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

	var q0021data types.Q0021Data
	var extradataq0021 types.Extradata = &q0021data

	err := GetItemD(int(iCompany), "Q0021", iAllMethod, iDate, &extradataq0021)
	if err != nil {
		return 0

	}

	noofdues := GetNoIstalments(iFromDate, iToDate, "M")
	fmt.Println("Inside Allocation", iCompany, iDate, iAllMethod, iFrequency, iFromDate, iToDate)

	for i := 0; i < len(q0021data.ALBand); i++ {
		if uint(noofdues) <= uint(q0021data.ALBand[i].Months) {
			iRate := q0021data.ALBand[i].Percentage
			return iRate
		}
	}
	return 0
}

// GetULMortPrem - Get Unit Linked Mortality Prem for a given duration
//
// Inputs: Company,  Coverage and Date String in YYYYMMDD, SA, Fund Value, Attained Age, Gender
//
// # Outputs Premium
//
// ©  FuturaInsTech
func GetULMortPrem(iCompany uint, iCoverage string, iDate string, iSA uint64, iFund uint64, iAge uint, iGender string) float64 {

	var q0006data types.Q0006Data
	var extradataq0006 types.Extradata = &q0006data
	// Get Coverage Rules
	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradataq0006)
	if err != nil {
		return 0

	}
	// Check Basis  1 = SAR  2 = SA  3 = SA + Fund
	var oSA uint64
	if q0006data.ULMortCalcType == "1" {
		oSA = iSA - iFund
	} else if q0006data.ULMortCalcType == "2" {
		oSA = iSA
	} else if q0006data.ULMortCalcType == "1" {
		oSA = iSA + iFund
	}

	var q0022data types.Q0022Data
	var extradataq0022 types.Extradata = &q0022data
	key := q0006data.ULMortDeductMethod + iGender
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
	if q0006data.ULMortFreq == "M" {
		aPrem = aPrem * 0.0833
	} else if q0006data.ULMortFreq == "Q" {
		aPrem = aPrem * 0.25
	} else if q0006data.ULMortFreq == "H" {
		aPrem = aPrem * 0.5
	}

	return aPrem

}

// GetGSTPercentage - Get GST Percemtage for a given months
//
// Inputs: Company,  Coverage and Date String in YYYYMMDD (Current Date), Key is Coverage Code, No of Months, Amount to be charged
//
// # Outputs GST Amount
//
// ©  FuturaInsTech
func GetGSTAmount(iCompany uint, iDate string, iKey string, iMonths uint64, iAmount float64) float64 {

	var q0023data types.Q0023Data
	var extradataq0023 types.Extradata = &q0023data

	// Get Premium Rate
	err := GetItemD(int(iCompany), "Q0023", iKey, iDate, &extradataq0023)
	if err != nil {
		return 0
	}

	for i := 0; i < len(q0023data.GST); i++ {
		if uint(iMonths) <= q0023data.GST[i].Month {
			oAmount := iAmount * q0023data.GST[i].Rate
			oAmount = RoundFloat(oAmount, 2)
			return oAmount
		}
	}
	return 0
}

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

func PostGlMove(iCompany uint, iContractCurry string, iEffectiveDate string,
	iTranno int, iGlAmount float64, iAccAmount float64, iAccountCodeID uint, iGlRdocno uint,
	iGlRldgAcct string, iSeqnno uint64, iGlSign string, iAccountCode string, iHistoryCode string) error {

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
		var p0031data types.P0031Data
		var extradata types.Extradata = &p0031data
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
	tx := initializers.DB.Save(&glmove)
	tx.Commit()

	UpdateGlBal(iCompany, iGlRldgAcct, iAccountCode, iContractCurry, iAccAmount, iGlSign, GlRdocno)
	return nil
}

// GetCommissionRates - Get Commission Rates
//
// Inputs: Company,  Coverage, Nof Instalments Collected (so far) and Date String in YYYYMMDD
//
// # Outputs Commission Rate
//
// ©  FuturaInsTech
func GetCommissionRates(iCompany uint, iCoverage string, iNofInstalemnts uint, iDate string) float64 {

	var p0028data types.P0028Data
	var extradatap0028 types.Extradata = &p0028data
	iKey := iCoverage
	fmt.Println("commission Rates **********", iCompany, iCoverage, iDate, iNofInstalemnts, iKey)
	// Get Premium Rate
	err := GetItemD(int(iCompany), "P0028", iKey, iDate, &extradatap0028)
	if err != nil {
		return 0
	}

	for i := 0; i < len(p0028data.Commissions); i++ {
		if uint(iNofInstalemnts) <= p0028data.Commissions[i].PPT {
			fmt.Println("Iam inside the array", p0028data.Commissions[i].PPT)
			oRate := p0028data.Commissions[i].Rate
			fmt.Println("i am getting in ", p0028data.Commissions[i].Rate)
			return oRate
		}
	}
	return 0
}

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
	var p0029data types.P0029Data
	var extradata types.Extradata = &p0029data
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

// func CreateTMFPolicy(iCompany uint, iPolicy uint, iFunction string) error {
// 	var tdfrule models.TDFRule
// 	var policy models.Policy
// 	var tdfpolicy models.TDFPolicy
// 	//var benefits []models.Benefit

// 	if iFunction == "ALL" {
// 		tdfrulereads := initializers.DB.Find(&tdfrule, "company_id = ? ", iCompany)
// 		if tdfrulereads != nil {
// 			if tdfrulereads.Error != nil {
// 				return errors.New(tdfrulereads.Error.Error())
// 			}
// 		}
// 	} else {
// 		tdfruleread := initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
// 		if tdfruleread.Error != nil {
// 			return errors.New(tdfruleread.Error.Error())
// 		}
// 	}

// 	switch {
// 	case iFunction == "ALL":

// 	case iFunction == "BILLD":
// 		tdfpolicy := initializers.DB.Find(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type= iFunction", iCompany, iPolicy, iFunction)
// 		if tdfpolicy.Error != nil {
// 			// Create Record
// 			tdfpolicy.CompanyID = iCompany
// 			tdfpolicy.PolicyID = iPolicy
// 			tdfpolicy.TDFType = tdfrule.TDFType
// 			tdfpolicy.TDFID = uint(tdfrule.Seqno)
// 			tdfpolicy.EffectiveDate = policy.BTDate

// 			result := initializers.DB.Create(&tdfpolicy)
// 			if result.Error != nil {
// 				return errors.New(tdfpolicy(tdfpolicy.Error.Error()))
// 			}

// 		} else {
// 			// Delete Record
// 			result := initializers.DB.Delete(&tdfpolicy)
// 			if result != nil {
// 				return errors.New(tdfpolicy.Error.Error())
// 			}
// 			// Create Record
// 			tdfpolicy.CompanyID = iCompany
// 			tdfpolicy.PolicyID = iPolicy
// 			tdfpolicy.TDFType = tdfrule.TDFType
// 			tdfpolicy.TDFID = uint(tdfrule.Seqno)
// 			tdfpolicy.EffectiveDate = policy.BTDate

// 			result := initializers.DB.Create(&tdfpolicy)
// 			if result.Error != nil {
// 				return errors.New(tdfpolicy(tdfpolicy.Error.Error()))
// 			}

// 		}
// 	case iFunction == "BILLD":
// 	case iFunction == "BILLD":

// 	default:
// 		break
// 	}
// }

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
		var q0005data types.Q0005Data
		var extradataq0005 types.Extradata = &q0005data
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

// TDFReraD - Time Driven Function - ReRate Date Updation
//
// Inputs: Company, Policy, Function RERAD, Transaction No.
//
// # Outputs  Old Record is Soft Deleted and New Record is Created
//
// ©  FuturaInsTech
func TDFReraD(iCompany uint, iPolicy uint, iFunction string, iTranno uint) (string, error) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var extraenq []models.Extra

	oDate := ""

	results := initializers.DB.Find(&extraenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if results.Error == nil {
		if results.RowsAffected > 1 {
			for i := 0; i < len(extraenq); i++ {
				if oDate == "" {
					oDate = extraenq[i].ToDate
				}
			}
		}

	}
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	if result.Error != nil {
		return "", result.Error
	}

	for i := 0; i < len(benefits); i++ {
		if benefits[i].BPremCessDate > benefits[i].BRerate {
			if oDate == "" {
				oDate = benefits[i].BRerate
			}

			if benefits[i].BRerate < oDate {
				oDate = benefits[i].BRerate
			}
		}

	}
	if oDate != "" {
		results = initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
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
			var q0006data types.Q0006Data
			var extradataq0006 types.Extradata = &q0006data
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
			var q0006data types.Q0006Data
			var extradataq0006 types.Extradata = &q0006data
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
		var q0006data types.Q0006Data
		var extradataq0006 types.Extradata = &q0006data
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
		var q0012data types.Q0012Data
		var extradataq0012 types.Extradata = &q0012data
		// fmt.Println("SB Parameters", iCompany, iType, iMethod, iYear, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0012", iMethod, iDate, &extradataq0012)
		fmt.Println("I am inside Term Based ")
		if err != nil {
			return err

		}
		// fmt.Println(q0012data.SBRates[0].Percentage)
		for x1 := 0; x1 <= iYear; x1++ {
			fmt.Println("X1Values are ", x1)
			for i := 0; i < len(q0012data.SBRates); i++ {
				fmt.Println("i Values are ", x1, i)
				if x1 == int(q0012data.SBRates[i].Term) {
					oSB := q0012data.SBRates[i].Percentage * iSA / 100
					// Write it in SB Table
					fmt.Println("Values of X and I", x1, i, iYear)
					survb.CompanyID = iCompany
					survb.PolicyID = iPolicy
					survb.PaidDate = ""
					survb.EffectiveDate = AddYears2Date(iDate, x1, 0, 0)
					survb.SBPercentage = q0012data.SBRates[i].Percentage
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
		var q0013data types.Q0013Data
		var extradataq0013 types.Extradata = &q0013data
		fmt.Println("SB Parameters", iCompany, iType, iMethod, iAge, iCoverage, iDate)
		err := GetItemD(int(iCompany), "Q0013", iMethod, iDate, &extradataq0013)
		fmt.Println("SB Parameters", iCompany, iCoverage, iDate)

		if err != nil {
			return err

		}
		fmt.Println(q0013data.SBRates[0].Percentage)
		for x := 0; x <= iAge; x++ {
			for i := 0; i < len(q0013data.SBRates); i++ {
				if x == int(q0013data.SBRates[i].Age) {
					oSB := q0013data.SBRates[i].Percentage * iSA / 100
					// Write it in SB Table
					survb.CompanyID = iCompany
					survb.PolicyID = iPolicy
					survb.PaidDate = ""
					survb.EffectiveDate = AddYears2Date(iDate, x, 0, 0)
					survb.SBPercentage = q0013data.SBRates[i].Percentage
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
	var p0034data types.P0034Data
	var extradatap0034 types.Extradata = &p0034data

	var p0033data types.P0033Data
	var extradatap0033 types.Extradata = &p0033data

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
	var p0036data types.P0036Data
	var extradata types.Extradata = &p0036data
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

// GetTolerance - To Get Tolerance for a Given Freqquency
//
// Inputs: Company,  Transaciton Code, Currency, Product, Date
//
// # Output Tolerance Amount
//
// ©  FuturaInsTech
func GetTolerance(iCompany uint, iTransaction string, iCurrency string, iProduct string, iFrequency string, iDate string) float64 {
	var p0043data types.P0043Data
	var extradata types.Extradata = &p0043data
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

// GetDeathAmount - Give Death Amount based on coverage and reason of death
//
// Inputs: Company Code, Policy, Coverage, Effective Date and cause of Death
//
// # Death Amount
//
// ©  FuturaInsTech
func GetDeathAmount(iCompany uint, iPolicy uint, iCoverage string, iEffectiveDate string, iCause string) (oAmount float64) {
	var benefit models.Benefit
	result := initializers.DB.Find(&benefit, "company_id = ? and policy_id = ? and b_coverage = ?", iCompany, iPolicy, iCoverage)

	if result.Error != nil {
		oAmount = 0
		return
	}

	iFund := float64(70000.00)
	iSA := float64(benefit.BSumAssured)
	iStartDate := benefit.BStartDate
	var q0006data types.Q0006Data
	var extradata types.Extradata = &q0006data
	iDate := benefit.BStartDate

	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)
	if err != nil {
		oAmount = 0
		return
	}

	ideathMethod := q0006data.DeathMethod //DC001
	oAmount = 0
	var p0049data types.P0049Data
	var extradata1 types.Extradata = &p0049data
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

	switch {
	case ideathMethod == "DC001": // Return of SA
		oAmount = iSA
		break
	case ideathMethod == "DC002": // Return of FV
		oAmount = iFund
		break
	case ideathMethod == "DC003": // Return of SA or Fund Value whichever is Highter
		if iSA >= iFund {
			oAmount = iSA
		} else {
			oAmount = iFund
		}
		break
	case ideathMethod == "DC004": // Return of SA + Fund Value
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
func WrapInArray(obj interface{}) interface{} {
	sliceType := reflect.SliceOf(reflect.TypeOf(obj))
	slice := reflect.MakeSlice(sliceType, 1, 1)
	slice.Index(0).Set(reflect.ValueOf(obj))
	return slice.Interface()
}

func NumberFunc(iAmount float64) (oAmount string) {

	p := message.NewPrinter(language.English)
	oAmount = p.Sprintf("%15.2f", iAmount)
	return
}
func NumbertoPrint(iAmount float64) (oAmount string) {

	p := message.NewPrinter(language.English)
	oAmount = p.Sprintf("%15.2f", iAmount)
	return
}

func IDtoPrint(iID uint) (oID string) {
	oID = strconv.FormatUint(uint64(iID), 10)
	fmt.Println(oID, reflect.TypeOf(oID))
	return
}

//reusable function

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
		var q0023data types.Q0023Data
		var extradataq0023 types.Extradata = &q0023data
		iKey := benefitenq1[a].BCoverage
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
			for i := 0; i < len(q0023data.GST); i++ {
				if uint(iMonths) <= q0023data.GST[i].Month {
					oAmount = float64(iAmount)*q0023data.GST[i].Rate + oAmount
					oAmount = RoundFloat(oAmount, 2)
					break

				}
			}

		}

	}
	return oAmount
}

func GetMRTABen(iSA float64, iInterest float64, iPolYear float64, iInterimPeriod float64, iTerm float64) float64 {
	a := math.Pow((1 + ((iInterest / 100) / 12)), ((iPolYear - iInterimPeriod) * 12))
	b := math.Pow((1 + ((iInterest / 100) / 12)), (iTerm * 12))
	c := (1 - (a-1)/(b-1))
	oSA := RoundFloat(iSA*c, 2)
	return oSA

}

func GetMrtaPrem(iCompany uint, iPolicy uint, iCoverage string, iAge uint, iGender string, iTerm uint, iPremTerm uint, iPremMethod string, iDate string, iMortality string) (float64, error) {

	var q0006data types.Q0006Data
	var extradata types.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)

	var q0010data types.Q0010Data
	var extradataq0010 types.Extradata = &q0010data
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
		//glmoveupd.UpdatedID = userid
		err := PostGlMove(uint(iCompany), iContractCurry, iEffectiveDate, int(iTranno), iGlAmount, iAccAmount, iAccountCodeID, uint(iGlRdocno), iGlRldgAcct, iSeqnno, iGlSign, iAccountCode, iHistoryCode)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSurrenderAmount(iCompany uint, iPolicy uint, iCoverage string, iEffectiveDate string, iTerm uint, iPremTerm uint, iStatus string, iSumAssured float64, iPaidTerm int, iStartDate string, iSurrMethod string, iInstallments int) (oAmount float64) {

	oAmount = 0
	var p0053data types.P0053Data
	var extradatap0053 types.Extradata = &p0053data

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
func CalculateStampDutyByPolicy(iCompanyId uint, iPolicyId uint) float64 {

	tStampDuty := 0.0
	var policyenq models.Policy
	result := initializers.DB.First(&policyenq, "company_id =? and id = ?", iCompanyId, iPolicyId)
	iDate := policyenq.PRCD

	if result.Error != nil {
		return 0.0

	}

	var benefitsenq []models.Benefit
	results := initializers.DB.Find(&benefitsenq, "company_id = ? and policy_id = ? ", iCompanyId, iPolicyId)

	if results.Error != nil {
		return 0.0
	}

	for i := 0; i < len(benefitsenq); i++ {
		iCoverage := benefitsenq[i].BCoverage
		iSA := benefitsenq[i].BSumAssured
		iInstalmentPaid := GetNoIstalments(benefitsenq[i].BStartDate, policyenq.PaidToDate, policyenq.PFreq)

		iStampDuty := CalculateStampDuty(iCompanyId, iCoverage, iInstalmentPaid, iDate, float64(iSA))

		tStampDuty = tStampDuty + iStampDuty

	}

	return tStampDuty

}

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

func GetBusinessDate(iCompany uint, iUser uint, iDepartment string) (oDate string) {
	var businessdate models.BusinessDate
	// Get with User
	result := initializers.DB.Find(&businessdate, "company_id = ? and user_id = ? and department = ?", iCompany, iUser, iDepartment)
	if result.Error == nil {
		// If User Not Found, get with Department
		result = initializers.DB.Find(&businessdate, "company_id = ? and department = ?", iCompany, iDepartment)
		if result.Error == nil {
			// If Department Not Found, get with comapny
			result = initializers.DB.Find(&businessdate, "company_id = ?", iCompany)
			if result.Error == nil {
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
	return Date2String(time.Now())
}

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

	var q0005data types.Q0005Data
	var extradataq0005 types.Extradata = &q0005data
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

// TdfhUpdate - Time Driven Function - Update TDF Header File
//
// Inputs: Company, Policy
//
// It has to loop through TDFPOLICIES and update earliest due in Tdfh
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

func GetFutureDue(iFromDate string, iToDate string, iFreq string) (oDate string) {
	// iFrom is Paid To Date
	// iTo is Current Date
	// Frequency is Policy Freqency
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
	var q0006data types.Q0006Data
	var extradata types.Extradata = &q0006data
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

// Function # 1
func GetCompanyData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
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
	}
	companyarray = append(companyarray, resultOut)
	return companyarray
}

// Function # 2
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

// Function # 3
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

// Function # 4
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

	var q0005data types.Q0005Data
	var extradataq0005 types.Extradata = &q0005data
	GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)
	gracedate := AddLeadDays(policy.PaidToDate, q0005data.LapsedDays)

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

	oAnnivDate := String2Date(policy.AnnivDate)
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
		// "PUWDate":DateConvert(policy.PUWDate),
	}
	policyarray = append(policyarray, resultOut)

	fmt.Print(policyarray)
	return policyarray
}

// Function # 5
func GetBenefitData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var benefit []models.Benefit

	initializers.DB.Find(&benefit, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	benefitarray := make([]interface{}, 0)

	for k := 0; k < len(benefit); k++ {
		iCompany := benefit[k].CompanyID
		_, oGender, _ := GetParamDesc(iCompany, "P0001", benefit[k].BGender, 1)
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", benefit[k].BCoverage, 1)
		_, oStatus, _ := GetParamDesc(iCompany, "P0024", benefit[k].BStatus, 1)

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
			"BSumAssured":    NumbertoPrint(float64(benefit[k].BSumAssured)),
			"BPrem":          NumbertoPrint(benefit[k].BPrem),
			"BGender":        oGender,
			"BDOB":           benefit[k].BDOB,
			"BMortality":     benefit[k].BMortality,
			"BStatus":        oStatus,
			"BAge":           benefit[k].BAge,
			"BRerate":        benefit[k].BRerate,
		}
		benefitarray = append(benefitarray, resultOut)
	}
	return benefitarray
}

// Function # 6
func GetSurBData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var survb []models.SurvB
	initializers.DB.Find(&survb, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	var benefitenq models.Benefit
	initializers.DB.Find(&benefitenq, "company_id = ? and policy_id =? and id = ?", iCompany, iPolicy, survb[0].BenefitID)
	basis := ""
	var q0006data types.Q0006Data
	var extradataq0006 types.Extradata = &q0006data

	GetItemD(int(iCompany), "Q0006", benefitenq.BCoverage, benefitenq.BStartDate, &extradataq0006)
	if q0006data.SBType == "A" {
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

// Function # 7
func GetMrtaData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var mrtaenq []models.Mrta
	initializers.DB.Find(&mrtaenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	mrtaarray := make([]interface{}, 0)
	for k := 0; k < len(mrtaenq); k++ {
		resultOut := map[string]interface{}{
			"ID":         IDtoPrint(mrtaenq[k].ID),
			"CompanyID":  IDtoPrint(mrtaenq[k].CompanyID),
			"Term":       mrtaenq[k].BTerm,
			"Ppt":        mrtaenq[k].PremPayingTerm,
			"ClientID":   IDtoPrint(mrtaenq[k].ClientID),
			"BenefitID":  IDtoPrint(mrtaenq[k].BenefitID),
			"PolicyID":   IDtoPrint(mrtaenq[k].PolicyID),
			"Coverage":   mrtaenq[k].BCoverage,
			"Product":    mrtaenq[k].Pproduct,
			"Interest":   mrtaenq[k].Interest,
			"DecreaseSA": mrtaenq[k].BSumAssured,
			"StartDate":  DateConvert(mrtaenq[k].BStartDate),
		}
		mrtaarray = append(mrtaarray, resultOut)
	}
	return mrtaarray
}

func GetReceiptData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var receiptenq models.Receipt
	initializers.DB.Find(&receiptenq, "company_id = ? and id = ?", iCompany, iReceipt)

	receiptarray := make([]interface{}, 0)
	resultOut := map[string]interface{}{
		"ID":                IDtoPrint(receiptenq.ID),
		"CompanyID":         IDtoPrint(receiptenq.CompanyID),
		"Branch":            receiptenq.Branch,
		"AccCurry":          receiptenq.AccCurry,
		"AccAmount":         receiptenq.AccAmount,
		"PolicyID":          IDtoPrint(receiptenq.PolicyID),
		"ClientID":          IDtoPrint(receiptenq.ClientID),
		"DateOfCollection":  DateConvert(receiptenq.DateOfCollection),
		"BankAccountNo":     receiptenq.BankAccountNo,
		"BankReferenceNo":   receiptenq.BankReferenceNo,
		"TypeOfReceipt":     receiptenq.TypeOfReceipt,
		"InstalmentPremium": receiptenq.InstalmentPremium,
		"AddressID":         IDtoPrint(receiptenq.AddressID),
		//		"PaidToDate":        DateConvert(receiptenq.PaidToDate),
		//		"ReconciledDate":    DateConvert(receiptenq.ReconciledDate),
		//		"CurrentDate":       DateConvert(receiptenq.CurrentDate),
	}
	receiptarray = append(receiptarray, resultOut)

	return receiptarray
}

func GetSaChangeData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var sachangeenq []models.SaChange
	initializers.DB.Find(&sachangeenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	sachangearray := make([]interface{}, 0)
	for k := 0; k < len(sachangeenq); k++ {
		resultOut := map[string]interface{}{
			"PolicyID    ": IDtoPrint(sachangeenq[k].PolicyID),
			"BenefitID":    IDtoPrint(sachangeenq[k].BenefitID),
			"BCoverage":    sachangeenq[k].BCoverage,
			"BStartDate":   DateConvert(sachangeenq[k].BStartDate),
			"BSumAssured":  sachangeenq[k].BSumAssured,
			"BTerm":        IDtoPrint(sachangeenq[k].BTerm),
			"BPTerm":       sachangeenq[k].BPTerm,
			"BPrem":        sachangeenq[k].BPrem,
			"BGender":      sachangeenq[k].BGender,
			"BDOB":         DateConvert(sachangeenq[k].BDOB),
			"NSumAssured":  sachangeenq[k].NSumAssured,
			"NTerm":        sachangeenq[k].NTerm,
			"NPTerm":       sachangeenq[k].NPTerm,
			"NPrem":        sachangeenq[k].NPrem,
			"NAnnualPrem":  sachangeenq[k].NAnnualPrem,
			"Method":       sachangeenq[k].Method,
			"Frequency":    sachangeenq[k].Frequency,
		}
		sachangearray = append(sachangearray, resultOut)
	}
	return sachangearray
}

func GetCompAddData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var addcomp []models.Addcomponent
	initializers.DB.Find(&addcomp, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	addcomparray := make([]interface{}, 0)

	for k := 0; k < len(addcomp); k++ {
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
			"BGender":     addcomp[k].BGender,
			"BDOB":        addcomp[k].BDOB,
			"Method":      addcomp[k].Method,
			"Frequency":   addcomp[k].Frequency,
			"BAge":        addcomp[k].BAge,
		}
		addcomparray = append(addcomparray, resultOut)
	}
	return addcomparray
}
func GetSurrHData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var surrhenq models.SurrH

	initializers.DB.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	surrharray := make([]interface{}, 0)

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
		"Product":           surrhenq.Product,
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

	return surrharray

}
func GetSurrDData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {

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

	return surrdarray

}
func GetNomiData(iCompany uint, iPolicy uint) []interface{} {

	var nomenq []models.Nominee

	initializers.DB.Find(&nomenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	nomarray := make([]interface{}, 0)

	for k := 0; k < len(nomenq); k++ {
		resultOut := map[string]interface{}{
			"ID":                  IDtoPrint(nomenq[k].ID),
			"PolicyID":            IDtoPrint(nomenq[k].PolicyID),
			"ClientID":            IDtoPrint(nomenq[k].ClientID),
			"NomineeRelationship": nomenq[k].NomineeRelationship,
			"LongName":            nomenq[k].NomineeLongName,
			"Percentage":          nomenq[k].NomineePercentage,
		}
		nomarray = append(nomarray, resultOut)
	}

	return nomarray

}

// Not Required
func GetDeathData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var surrhenq models.SurrH
	var surrdenq []models.SurrD
	initializers.DB.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	initializers.DB.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	surrarray := make([]interface{}, 0)

	return surrarray
}

func GetMatHData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {
	var mathenq models.MaturityH

	initializers.DB.Find(&mathenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	matharray := make([]interface{}, 0)

	resultOut := map[string]interface{}{

		"ID":                   IDtoPrint(mathenq.ID),
		"PolicyID":             IDtoPrint(mathenq.PolicyID),
		"ClientID":             IDtoPrint(mathenq.ClientID),
		"EffectiveDate":        DateConvert(mathenq.EffectiveDate),
		"MaturityDate":         DateConvert(mathenq.MaturityDate),
		"Status":               mathenq.Status,
		"BillDate":             DateConvert(mathenq.BillDate),
		"PaidToDate":           DateConvert(mathenq.PaidToDate),
		"Product":              mathenq.Product,
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

	return matharray

}
func GetMatDData(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint) []interface{} {

	var matdenq []models.MaturityD

	initializers.DB.Find(&matdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	matdarray := make([]interface{}, 0)

	for k := 0; k < len(matdenq); k++ {
		resultOut := map[string]interface{}{
			"MaturityHID":         IDtoPrint(matdenq[k].MaturityHID),
			"PolicyID":            matdenq[k].PolicyID,
			"ClientID":            matdenq[k].ClientID,
			"BenifitID":           matdenq[k].BenefitID,
			"BCoverage":           matdenq[k].BCoverage,
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

	return matdarray

}

// SANDHYA
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
		"PolicyDeposit": oPolicyDeposit,
		"RevBonus":      oRevBonus,
		"TermBonus":     oTermBonus,
		"IntBonus":      oIntBonus,
		"AccDividend":   oAccumDiv,
		"AccDivInt":     oAccumDivInt,
		"AddBonus":      oAddBonus,
		"LoyalBonus":    oLoyalBonus,
		"AplAmount":     oAplAmt,
		"PolLoan":       oPolLoan,
		"CashDeposit":   oCashDep,
	}
	bonusarray = append(bonusarray, resultOut)

	return bonusarray
}

func GetAgency(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iTranno uint, iAgency uint) []interface{} {

	agencyarray := make([]interface{}, 0)
	var agencyenq models.Agency
	var clientenq models.Client
	initializers.DB.Find(&agencyenq, "company_id  = ? and id = ?", iCompany, iAgency)

	initializers.DB.Find(&clientenq, "company_id = ? and id = ?", iCompany, agencyenq.ClientID)
	oAgentName := clientenq.ClientLongName + " " + clientenq.ClientShortName + " " + clientenq.ClientSurName

	resultOut := map[string]interface{}{
		"ID":              IDtoPrint(iAgency),
		"AgyChannelSt":    agencyenq.AgencyChannelSt,
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
	}
	agencyarray = append(agencyarray, resultOut)

	return agencyarray
}

func GetExpi(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iTranno uint) []interface{} {
	var benefit []models.Benefit
	initializers.DB.Find(&benefit, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno)
	expiryarray := make([]interface{}, 0)

	for k := 0; k < len(benefit); k++ {
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
			"BCoverage":      benefit[k].BCoverage,
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

// Check Status
//
// # This function, take company code, history code, date and status as inputs
//
// # It returns status which is boolean and also output status which is string
//
// ©  FuturaInsTech
func CheckStatus(iCompany uint, iHistoryCD string, iDate string, iStatus string) (status bool, oStatus string) {
	var p0029data types.P0029Data
	var extradata types.Extradata = &p0029data

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

// Create Communication
//
// # This function, Create Communication Records by getting input values as Company ID, History Code, Tranno, Date of Transaction, Policy Id, Client Id, Address Id, Receipt ID . Quotation ID, Agency ID
// 10 Input Variables
// # It returns success or failure.  Successful records written in Communciaiton Table
//
// ©  FuturaInsTech
func CreateCommunications(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint) error {

	var p0034data types.P0034Data
	var extradatap0034 types.Extradata = &p0034data

	var p0033data types.P0033Data
	var extradatap0033 types.Extradata = &p0033data
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

			resultMap := make(map[string]interface{})

			//	iCompany uint, iPolicy uint, iAddress uint, iClient uint, iLanguage uint, iBankcode uint, iReceipt uint, iCommunciation uint, iQuotation uint
			for n := 0; n < len(p0034data.Letters[i].LetType); n++ {
				oLetType = p0034data.Letters[i].LetType[n]
				switch {
				case oLetType == "1":
					oData := GetCompanyData(iCompany, iPolicy, iClient, iAddress, iReceipt)
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
					resultMap["SurrHData"] = oData
					oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData := GetDeathData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData := GetMatHData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["MatHData"] = oData
					oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					resultMap["MatDData"] = oData
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
