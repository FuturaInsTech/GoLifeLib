package paramTypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
)

type Extradata interface {
	// Methods
	ParseData(map[string]interface{})
	GetFormattedData(datamap map[string]string) map[string]interface{}
}

// Q0005
type Q0005Data struct {
	FreeLookDays           int
	MaxLives               int
	MinLives               int
	MinSurrMonths          int
	ProductFamily          string
	ReinstatementMonth     int
	Renewable              string
	Single                 string
	Frequencies            [4]string // M Q H Y
	ContractCurr           []string  // INR USD SGD HKD LKR PKR DON IDR
	BillingCurr            []string
	ComponentAddAtAnyTime  string // Policy Anniversary or Any Time N / Y
	FuturePremAdj          string //Y or N
	FuturePremAdjYrs       int    // eg., 3 Yrs
	LapsedDays             int
	BillingLeadDays        int
	LapseInterest          float64
	AgencyChannel          []string //P0050
	BackDateAllowed        string   // P0050  YESNO
	NoLapseGuarantee       string   //P0050 YESNO
	NoLapseGuaranteeMonths int
	SpecialRevivalMonths   int
	AplLoanMethod          string
	NfoMethod              []string //P0050
	UnderwitingReq         string   // P0050 YESNO
	BankReq                string   //P0050 YESNO
}

func (m *Q0005Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)
	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

// func (m *Q0005Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

// 	return nil

// }

func (m *Q0005Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	coy, _ := strconv.Atoi(datamap["company_id"])
	langid, _ := strconv.Atoi(datamap["LanguageId"])

	if datamap["function"] == "BillingCurr" {
		resp := make(map[string]interface{})
		// allowedbilling := make([]string, 0)
		resultarray := make([]interface{}, 0)
		for i := 0; i < len(m.BillingCurr); i++ {
			if m.BillingCurr[i] == "" {
				break
			}
			short, long, _ := GetParamDesc(uint(coy), "P0023", m.BillingCurr[i], uint(langid))

			resultOut := map[string]interface{}{
				"Item":      m.BillingCurr[i],
				"ShortDesc": short,
				"LongDesc":  long,
			}

			resultarray = append(resultarray, resultOut)

		}
		resp["AllowedBillingCurriencies"] = resultarray
		return resp
	} else if datamap["function"] == "ContractCurr" {
		resp := make(map[string]interface{})
		//	contractcurr := make([]string, 0)
		resultarray := make([]interface{}, 0)
		for i := 0; i < len(m.ContractCurr); i++ {
			if m.ContractCurr[i] == "" {
				break
			}
			short, long, _ := GetParamDesc(uint(coy), "P0023", m.ContractCurr[i], uint(langid))

			resultOut := map[string]interface{}{
				"Item":      m.ContractCurr[i],
				"ShortDesc": short,
				"LongDesc":  long,
			}

			resultarray = append(resultarray, resultOut)
		}
		resp["AllowedContractCurriencies"] = resultarray
		return resp
	} else if datamap["function"] == "Frequencies" {
		resp := make(map[string]interface{})
		//allowedfreq := make([]string, 0)
		resultarray := make([]interface{}, 0)

		for i := 0; i < len(m.Frequencies); i++ {
			if m.Frequencies[i] == "" {
				break
			}

			short, long, _ := GetParamDesc(uint(coy), "Q0009", m.Frequencies[i], uint(langid))

			resultOut := map[string]interface{}{
				"Item":      m.Frequencies[i],
				"ShortDesc": short,
				"LongDesc":  long,
			}

			resultarray = append(resultarray, resultOut)
		}
		resp["AllowedFrequencies"] = resultarray
		return resp
	} else if datamap["function"] == "NfoMethod" {
		resp := make(map[string]interface{})
		//allowedfreq := make([]string, 0)
		resultarray := make([]interface{}, 0)

		for i := 0; i < len(m.NfoMethod); i++ {
			if m.NfoMethod[i] == "" {
				break
			}

			long := GetP0050ItemCodeDesc(uint(coy), "NfoMethod", 1, m.NfoMethod[i])

			resultOut := map[string]interface{}{
				"Item":     m.NfoMethod[i],
				"LongDesc": long,
			}

			resultarray = append(resultarray, resultOut)
		}
		resp["NfoMethod"] = resultarray
		return resp
	} else {
		return nil
	}

}

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

// /
type Q0006Data struct {
	AgeCalcMethod string //N/L/X  P0050
	AnnMethod     string //Bonus Method  P/D P0050 (later)
	AnnuityMethod string //Annuity Method
	CommMethod    string //Commission
	//DeathType          string // 1 SA 2 FV 3 GT(SA/FV) 4. SA + FV
	DeathMethod        string // DC****   P0050
	GBonus             string
	IBonus             string
	LoanMethod         string
	LoyaltyBonus       string
	MatMethod          string
	AgeRange           []uint // Age Range - Array
	PptRange           []uint // Premium Paying Term Range - Array
	MaxSA              uint
	TermRange          []uint // Term Range - Array
	MinSA              uint
	MinRiskCessAge     uint
	MaxRiskCessAge     uint
	MinPremCessAge     uint
	MaxPremCessAge     uint
	MaxTermBeyongCover uint   // ETI
	NofMethod          string // NFO Method
	PartSurrMethod     string // Part Surrender Method
	PremInc            string // Premium Increase Allowed
	PremIncYrs         uint   // No of Yrs in Number
	PremiumMethod      string // Premium Method
	RevBonus           string // Reversionary Bonus Method
	SbType             string // Either Age Based or Term Based (later)
	SbMethod           string // Q00012 OR Q0013
	SurrMethod         string
	TBonus             string
	GsvMethod          string
	SsvMethod          string
	BsvMethod          string
	DivMethod          string
	DivIMethod         string
	Mortalities        []Q0006M // Smoker/Non Smoker/Combined - Array  //P0050
	PremCalcType       string   // Either Age Based / PPT Based
	DiscType           string   // S for SA P for Pemium  ?? Do we need?  (later)
	DiscMethod         string   // DM001 or DM002
	FrqMethod          string   // Frequency Factor
	WaivMethod         string   // Waiver Method Q0020
	// Unit Linked Components
	//UlDeductFrequency    string    // UL Fee Deduction Frequency
	UlAlMethod string // UL Prem Allocation Method  //P0050   P0060
	UlMortFreq string // UL Mortality Deduction Frequency //P0050
	// UlMortCalcType       string    // 1 - SAR 2 - SA, 3 - Fund + SA  Q0022  //P0050
	UlMorttMethod string // UL Mortality Deduction Method Q0022 Attained Age // premium rates
	UlFeeFreq     string // UL Charges Deduction Frequency //P0050
	// UlFeeType            string    // 1 SA Based 2 Annualised Premium 3 Fund Value  //P0050
	UlFeeMethod          string    // UL Fee Method  //P0050
	UlFundMethod         string    // UL Fund Rules  //P0050
	FUNDCODE             []string  //P0050
	UlTopUpMethod        string    //P0050
	UlWithdrawMethod     string    //P0050
	MrtaMethod           string    // MRTA Method
	MrtaInterest         []float64 // MRTA Interest - Array
	BenefitType          string    // Health,CI,Waiver,Pension etc., P0050
	CommissionOnExtraInd string    //P0050 Yes/No
	UlSwitchMethod       string    //P0050
	CovrFamily           string
	//Health Ins Components
	isHealthBenefitPlan string //P0050 Yes/No
	HealthBenefitType   string //P0050 Individual/Non-Floater/Family-Floater
	ColaMethod          string //P0050
	UinNo               string
}

func (m *Q0006Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)
	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0006Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	coy, _ := strconv.Atoi(datamap["company_id"])
	langid, _ := strconv.Atoi(datamap["LanguageId"])

	if datamap["function"] == "MrtaInterest" {
		resp := make(map[string]interface{})
		allowedinterest := make([]float64, 0)
		for i := 0; i < len(m.MrtaInterest); i++ {
			if m.MrtaInterest[i] == 0 {
				break
			}

			allowedinterest = append(allowedinterest, m.MrtaInterest[i])

		}
		resp["AllowedInterestRates"] = allowedinterest
		return resp

	} else if datamap["function"] == "UlAlMethod" {
		resp := make(map[string]interface{})
		// UlAlMethod := make([]string, 0)
		UlAlMethod := m.UlAlMethod
		// for i := 0; i < len(m.UlAlMethod); i++ {
		// 	if m.UlAlMethod[i] == "" {
		// 		break
		// 	}

		// 	UlAlMethod = append(UlAlMethod, m.UlAlMethod[i])

		// }
		resp["AllowedUlAlMethod"] = UlAlMethod
		return resp
	} else if datamap["function"] == "AgeRange" {
		resp := make(map[string]interface{})
		allowedagerange := make([]uint, 0)
		for i := 0; i < len(m.AgeRange); i++ {
			if m.AgeRange[i] == 0 {
				break
			}

			allowedagerange = append(allowedagerange, m.AgeRange[i])

		}
		resp["AllowedAgeRange"] = allowedagerange
		return resp

	} else if datamap["function"] == "FundCode" {
		resp := make(map[string]interface{})
		//funds := make([]string, 0)
		resultarray := make([]interface{}, 0)
		for i := 0; i < len(m.FUNDCODE); i++ {
			if m.FUNDCODE[i] == "" {
				break
			}
			short, long, _ := GetParamDesc(uint(coy), "P0061", m.FUNDCODE[i], uint(langid))

			resultOut := map[string]interface{}{
				"Item":      m.FUNDCODE[i],
				"ShortDesc": short,
				"LongDesc":  long,
			}

			resultarray = append(resultarray, resultOut)
		}
		resp["Funds"] = resultarray
		return resp
	} else if datamap["function"] == "PptRange" {
		resp := make(map[string]interface{})
		allowedpptrange := make([]uint, 0)
		for i := 0; i < len(m.PptRange); i++ {
			if m.PptRange[i] == 0 {
				break
			}

			allowedpptrange = append(allowedpptrange, m.PptRange[i])

		}
		resp["AllowedPptRange"] = allowedpptrange
		return resp
	} else if datamap["function"] == "TermRange" {
		resp := make(map[string]interface{})
		allowedtermrange := make([]uint, 0)
		for i := 0; i < len(m.TermRange); i++ {
			if m.TermRange[i] == 0 {
				break
			}

			allowedtermrange = append(allowedtermrange, m.TermRange[i])

		}
		resp["AllowedTermRange"] = allowedtermrange
		return resp
	} else if datamap["function"] == "PremCalcType" {
		resp := make(map[string]interface{})
		PremCalcType := m.PremCalcType
		resp["AllowedPremCalcType"] = PremCalcType
		return resp
	} else {
		return nil
	}
}

// Mortality
type Q0006M struct {
	Mortality string
}

func (m *Q0006M) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

// func (m *Q0006M) GetFormattedData(datamap map[string]string) map[string]interface{} {

// 	return nil

// }

// /
type Q0010Data struct {
	Rates []Q0010
}

type Q0010 struct {
	Age  uint
	Rate float64
}

func (m *Q0010Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)
	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0010Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type Q0011Data struct {
	Coverages []Q0011
}

type Q0011 struct {
	CoverageName  string
	Mandatory     string
	BasicorRider  string
	TermCanExceed string
	PptCanExceed  string
	SaCanExceed   string
}

func (m *Q0011Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)
	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

// func (m *Q0011Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

// 	resp := make(map[string]interface{})
// 	coverages := make([]string, 0)
// 	for i := 0; i < len(m.Coverages); i++ {
// 		if m.Coverages[i].CoverageName == "" {
// 			break

// 		}
// 		coverages = append(coverages, m.Coverages[i].CoverageName)

// 	}
// 	resp["coverages"] = coverages
// 	return resp

// }

func (m *Q0011Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	resp := make(map[string]interface{})
	allowedcoverages := make([]interface{}, 0)
	for i := 0; i < len(m.Coverages); i++ {
		if m.Coverages[i].CoverageName == "" {
			break
		}

		coverage := map[string]interface{}{
			"Coverage":           m.Coverages[i].CoverageName,
			"Mandatory":          m.Coverages[i].Mandatory,
			"Basic":              m.Coverages[i].BasicorRider,
			"TermCanExceedBasic": m.Coverages[i].TermCanExceed,
			"PPTCanExceedBasic":  m.Coverages[i].PptCanExceed,
			"SACanExceedBasic":   m.Coverages[i].SaCanExceed,
		}
		allowedcoverages = append(allowedcoverages, coverage)

	}
	resp["AllowedCoverages"] = allowedcoverages
	return resp

}

// /
type Q0012Data struct {
	SbRates []Q0012
}

type Q0012 struct {
	Term       uint
	Percentage float64
}

func (m *Q0012Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {

		fmt.Println(err)
	}

}

func (m *Q0012Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Survival Benefits
type Q0013Data struct {
	SbRates []Q0013
}

type Q0013 struct {
	Age        uint
	Percentage float64
}

func (m *Q0013Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {

		fmt.Println(err)
	}

}
func (m *Q0013Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Bonus
type Q0014Data struct {
	BRates []Q0014
}

type Q0014 struct {
	Term       uint
	Percentage float64
}

func (m *Q0014Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {

		fmt.Println(err)
	}

}

// func (m *Q0014Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

// 	return nil

// }

//Term

func (m *Q0014Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	function := datamap["function"]
	resp := make(map[string]interface{})
	bonuses := make([]interface{}, 0)
	for i := 0; i < len(m.BRates); i++ {
		if m.BRates[i].Term == 0 {
			break
		}
		if function == "getTerm" {
			bonus := map[string]interface{}{
				"Term": m.BRates[i].Term,
			}
			bonuses = append(bonuses, bonus)
		}
		if function == "getPercentage" {
			bonus := map[string]interface{}{
				"Percentage": m.BRates[i].Percentage,
			}
			bonuses = append(bonuses, bonus)
		}
		if function == "" {
			bonus := map[string]interface{}{
				"Term":       m.BRates[i].Term,
				"Percentage": m.BRates[i].Percentage,
			}
			bonuses = append(bonuses, bonus)
		}
	}
	resp["Bonuses"] = bonuses
	return resp

}

type Q0015Data struct {
	Terms []Q0015
}

type Q0015 struct {
	Term uint
}

func (m *Q0015Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {

		fmt.Println(err)
	}

}

// func (m *Q0015Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

// 	return nil

// }
func (m *Q0015Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	resp := make(map[string]interface{})
	ppt := make([]int, 0)
	for i := 0; i < len(m.Terms); i++ {
		// var x uint
		// x = m.PTerms[i]
		if m.Terms[i].Term == 0 {
			// if uint(x) == 0 {
			break

		}
		//ppt = append(ppt, m.PTerms[i])
		ppt = append(ppt, int(m.Terms[i].Term))

	}
	resp["ppt"] = ppt
	return resp

}

// PPT
type Q0016Data struct {
	PTerms []Q0016
}

type Q0016 struct {
	PTerm uint
}

func (m *Q0016Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {

		fmt.Println(err)
	}

}

// func (m *Q0016Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

// 	return nil

// }

func (m *Q0016Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	resp := make(map[string]interface{})
	ppt := make([]int, 0)
	for i := 0; i < len(m.PTerms); i++ {
		// var x uint
		// x = m.PTerms[i]
		if m.PTerms[i].PTerm == 0 {
			// if uint(x) == 0 {
			break

		}
		//ppt = append(ppt, m.PTerms[i])
		ppt = append(ppt, int(m.PTerms[i].PTerm))

	}
	resp["ppt"] = ppt
	return resp

}

// SA Discount
type Q0017Data struct {
	SaBand []Q0017
}
type Q0017 struct {
	Sa       uint
	Discount float64
}

func (m *Q0017Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0017Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Premium  Discount
type Q0018Data struct {
	PremBand []Q0018
}
type Q0018 struct {
	AnnPrem  float64
	Discount float64
}

func (m *Q0018Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0018Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Frequency Factor
type Q0019Data struct {
	FreqFactor []Q0019
}
type Q0019 struct {
	Frequency string
	Factor    float64
}

func (m *Q0019Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0019Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Waiver of Premium
type Q0020Data struct {
	WaiverCoverages []Q0020
}
type Q0020 struct {
	Coverage string
}

func (m *Q0020Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0020Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Allocation Method
// Key AL001 + Transaction
// AL01H0007
// AL01B0102
// 5000 REGULAR PREMIUM 10000 REGULAR TOP UP FOR 5 YEARS

type P0060Data struct {
	AlBand []P0060
}

type P0060 struct {
	Months     uint
	Percentage float64
}

func (m *P0060Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0060Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type Q0022Data struct {
	Rates []Q0022
}

type Q0022 struct {
	Age  uint
	Rate float64
}

func (m *Q0022Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0022Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type Q0023Data struct {
	Gst []Q0023
}
type Q0023 struct {
	Month uint
	Rate  float64
}

func (m *Q0023Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0023Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type Q0024Data struct {
	BiRates []Q0024
}
type Q0024 struct {
	BiType string
	Rate   float64
}

func (m *Q0024Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0024Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// GL Movements

type P0027Data struct {
	GlMovements []P0027
}
type P0027 struct {
	AccountCode string
	AccountAmt  float64
	SeqNo       uint
	GlSign      string
}

func (m *P0027Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0027Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// commission Rates

type P0028Data struct {
	Commissions []P0028
}
type P0028 struct {
	Ppt  uint
	Rate float64
}

func (m *P0028Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0028Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Allowed Status

type P0029Data struct {
	Statuses []P0029
}
type P0029 struct {
	CurrentStatus string
	ToBeStatus    string
}

func (m *P0029Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0029Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Collection Bank Account  Redundant
type P0030Data struct {
	BankAccount string
	GlAccount   string
}

func (m *P0030Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0030Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Currency Rates
type P0031Data struct {
	CurrencyRates []P0031
}
type P0031 struct {
	Action string
	Rate   float64
}

func (m *P0031Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0031Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// UW Rules

type P0032Data struct {
	UwRules []P0032
}

type P0032 struct {
	NoOfMonths uint
	Factor     float64
}

func (m *P0032Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0032Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Templates

type P0033Data struct {
	TemplateName          string
	SMSAllowed            string
	EmailAllowed          string
	WhatsAppAllowed       string
	AgentSMSAllowed       string
	AgentEmailAllowed     string
	AgentWhatsAppAllowed  string
	CompanyEmail          string
	CompanyPhone          string
	DepartmentName        string
	DepartmentHead        string
	SenderPassword        string
	SMTPServer            string
	SMTPPort              int
	Body                  string
	Subject               string
	Online                string // (Y/N)
	CarbonCopy            string // (Y/N)
	BlindCarbonCopy       string // (Y/N)
	SMSSID                string // SMS ID
	SMSAuthToken          string // SMS Authorized Token
	SMSAuthPhoneNo        string // SMS No which has been registered with the service provider
	SMSBody               string
	WhatsappPhoneNumberID string
	WhatsappAuthToken     string
	WhatsappBody          string
}

func (m *P0033Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0033Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Letter Groups

type P0034Data struct {
	Letters []P0034
}

type P0034 struct {
	Templates              string
	ReportTemplateLocation string
	PdfLocation            string
	LetType                []string //1,2,3,4,5,6,7,...
	PageSize               string
	Orientation            string
}

func (m *P0034Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0034Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// FreeLook Rules
type P0035Data struct {
	PremProp           string
	CommRecovPercetage float64
	MedicalFeeRecovery string
	GstRecovery        string
	StampDuty          string
}

func (m *P0035Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0035Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Stamp Duty
type P0036Data struct {
	StampDuties []P0036
}
type P0036 struct {
	Noofinstalments int
	Sa              float64
	Rate            float64
}

func (m *P0036Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0036Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// P0040 - Medical Fee Company Level

type P0040Data struct {
	MedicalFee  float64
	MedicalCurr string
}

func (m *P0040Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0040Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// P0041 SA Band for Medical Requirements

type P0041Data struct {
	SumAssured []P0041
}
type P0041 struct {
	Sa    float64
	Age   uint
	Codes []string
}

func (m *P0041Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0041Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// P0043 - Tolerance

type P0043Data struct {
	Frequencies []P0043
}

type P0043 struct {
	Frequency string
	Amount    float64
}

func (m *P0043Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0043Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Menu Switching Parameter
type P0044Data struct {
	Actions []P0044
}

type P0044 struct {
	Action      string
	Description string
	Url         string
	Trancode    string
}

func (m *P0044Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)
	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0044Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	resp := make(map[string]interface{})
	allowedmenus := make([]interface{}, 0)
	for i := 0; i < len(m.Actions); i++ {
		if m.Actions[i].Action == "" {
			break
		}

		action := map[string]interface{}{
			"Action":      m.Actions[i].Action,
			"Description": m.Actions[i].Description,
			"URL":         m.Actions[i].Url,
			"Trancode":    m.Actions[i].Trancode,
		}
		allowedmenus = append(allowedmenus, action)

	}
	resp["AllowedMenus"] = allowedmenus
	return resp

}

// p0045data

type P0045Data struct {
	Gender    string
	RelatedTo string
}

func (m *P0045Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0045Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

// Cause of Death Rule
// P0043 - Tolerance

// P0049 Death Claim Calculation by Cause of Death
type P0049Data struct {
	Months []P0049
}

type P0049 struct {
	Month       uint
	Percentage  float64
	DeathMethod string
}

func (m *P0049Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0049Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Benefit Illsutration Interval Years by Coverage

type Q0025Data struct {
	BiYrInterval uint
}

func (m *Q0025Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0025Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Surrender Parameters

type P0053Data struct {
	Rates []P0053
}

type P0053 struct {
	Month      uint
	Percentage float64
}

func (m *P0053Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0053Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Transaction Table
type P0054Data struct {
	Trancodes []P0054
}

type P0054 struct {
	Transaction string
}

func (m *P0054Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0054Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Billing Type
type P0055Data struct {
	BankRequired    string // Y or N  Client should have bank account Y/N
	BankCode        string // IFSC No of Insurance Company
	BankAccount     string // Bank Account No of Insurance Company
	GlAccount       string // GL Code for Posting for the Billing Type
	PayingAuthority string // Y or N through P0050
	Vpa             string // Virtual Payment Authority Code say 8825761193@upi etc.,
}

func (m *P0055Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0055Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Bank Code Rules
type P0056Data struct {
	NoOfDishours         int     // No of Dishnours Allowed
	ProcessFlag          string  // Process Flag
	ExtractionDates      []P0056 // 5, 10,15,25,28, 30,31
	NetCollection        string  // Net Colleciton/Gross Collection (N/G)
	CollectionFee        float64 // Flat Fee for Each Collection
	CollectionPercentage float64 // Collection Fee
	AccountEntry         string  // Pass Accounting Entries (Y/N)
}
type P0056 struct {
	ExtractionDate string
}

func (m *P0056Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0056Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Branch Code

type P0018Data struct {
	BankIfsc    string
	BankAccount string
	ClientID    uint
}

func (m *P0018Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0018Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Critical Illness Rules
type P0057Data struct {
	Rules []P0057
}

type P0057 struct {
	Months     uint
	Percentage float64
}

func (m *P0057Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0057Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Income Benefit Rules

type P0058Data struct {
	Percentage                     float64 // 100.. IT CAN BE 200 300
	NoOfYears                      int
	AdjustPayTerm                  string // Y OR N
	LiabilityPosting               string
	CertificateOfExistanceRequired string // Y OR N
	CertficiateOfExistanceLeadDays int    // DAYS IN ADVANCE
	FollowBenefitRCD               string `gorm:"type:varchar(1)"` // Y means Benefit RCD N Menas Incident RCD
}

func (m *P0058Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0058Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// ILP Rules
// Transaction Code + Coverage Code
type P0059Data struct {
	CurrentOrFuture    string `gorm:"type:varchar(1)"` // P0050
	SeqNo              int
	AllocationCategory string // P0050 2 Character
	AccountCode        string `gorm:"type:varchar(30)"`
	// NegativeUnits            string // P0050 YES/NO Y
	// NegativeUnitsPeriod      uint   // In Months    36
	// NegativeAmounts          string //P0050 YES/NO  N
	// NegativeAmountsPeriod    uint   // In Months    26
	// NegUnitsOrAmtRecovPeriod uint   //In Months     60 MONTHS

	NegativeAccum         string  // P0050 Y/N  (ILP  Only)
	NegativeAccumMonths   float64 // No of Months of Negative Accum
	NegativeUnitsOrAmt    string  // unit or Amounts P0050 U/A
	RecoverFromTopUpFirst string  // YES/NO

}

func (m *P0059Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0059Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Fund Information
type P0061Data struct {
	FundCode         string // P0050  FUNDCODE
	FundType         string // P0050  FUNDTYPE
	FundCategory     string // P0050  FUNDCATEGORY
	FundCurr         string // P0050  FUNDCURR
	FundMinUnits     string
	FundMaxUnits     string
	FundChargeMethod string //P0050   FUNDCHARGEMETHOD
}

func (m *P0061Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0061Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	//coy, _ := strconv.Atoi(datamap["company_id"])
	//langid, _ := strconv.Atoi(datamap["LanguageId"])

	if datamap["function"] == "P0061" {
		resp := make(map[string]interface{})
		resultarray := make([]interface{}, 0)

		UlFundCatory := m.FundCategory
		UlFundCurr := m.FundCurr
		UlFundType := m.FundType

		resultOut := map[string]interface{}{
			"FundCategory": UlFundCatory,
			"FundCurr":     UlFundCurr,
			"FundType":     UlFundType,
		}

		resultarray = append(resultarray, resultOut)

		resp["P0061"] = resultarray
		return resp
	} else {
		return nil
	}
}

// Minimum and Maximum Premium Limit
// Coverage Level + Currency

type P0062Data struct {
	MinMaxRule []P0062
}
type P0062 struct {
	Frequency string // FREQ  P0050
	Premium   float64
}

func (m *P0062Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0062Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

/*func GetParamDesc(iCompany uint, iParam string, iItem string, iLanguage uint) (string, string, error) {
	type Descs struct {
		Longdesc  string
		Shortdesc string
	}

	var descs Descs

	results := initializers.DB.Table("param_descs").Select("longdesc", "shortdesc").Where("company_id = ? AND name = ? and item = ? and language_id = ?", iCompany, iParam, iItem, iLanguage).Scan(&descs)
	if results.Error != nil || results.RowsAffected == 0 {

		return "", "", errors.New(" -" + strconv.FormatUint(uint64(iCompany), 10) + "-" + iParam + "-" + "-" + iItem + "-" + strconv.FormatUint(uint64(iLanguage), 10) + "-" + " is missing")
		//return errors.New(results.Error.Error())
	}

	return descs.Shortdesc, descs.Longdesc, nil
}
*/
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

// P0063 - Flat Amount or Flat % on Fund Value  or Simple or Compount

type P0063Data struct {
	FlatAmount        float64 // Flat Amount Per Year
	FundValPercentage float64 // Flat Percentage on Fund Value Per Year
	Percentage        float64 // Increase % Per Year
	SimpleOrCompound  string  // Simple or Compunt Indicator P0050
	CapAmount         float64 // Max Amount
}

func (m *P0063Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0063Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// P0064 - ILP Surrender Penalty

type P0064Data struct {
	SurrenderPenalty []P0064
}
type P0064 struct {
	NoOfMonths        int
	PenaltyPercentage float64
}

func (m *P0064Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0064Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// P0023 - Currency Printing Details

type P0023Data struct {
	CurSymbol   string // Currency Symbol $, £, €, etc
	CurBill     string // Dollars, Pounds, Euros, etc
	CurCoin     string // Cents, Pence, Cents, etc
	CurWordType string // M - Millions, Billions, Trilions and  L - Lakhs and Crores

}

func (m *P0023Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0023Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type P0065Data struct {
	FieldList []P0065
}
type P0065 struct {
	Field     string // Field Name of the Table
	Mandatory string // P0050 Yes/No
	ErrorCode string // Error Code Table
}

func (m *P0065Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0065Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Country, Phone No, Country Code and Flag
type P0066Data struct {
	Name     string
	DialCode string
	Code     string
	Flag     string
}

func (m *P0066Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0066Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Tax Rules
type P0067Data struct {
	GlTax []P0067
}
type P0067 struct {
	AccountCode string
	TaxSection  string
}

func (m *P0067Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0067Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// ILP SA Rules

type P0068Data struct {
	RangeArray []P0068Range
}
type P0068Range struct {
	P0068Basis string // R - Range M - Multiplier  P0050
	Age        uint
	FromSA     float64
	ToSA       float64
	Factor     float64
}

func (m *P0068Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0068Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

// Extended Lapse Rule

type P0069Data struct {
	P0069 []P0069Lapse
}
type P0069Lapse struct {
	Months            uint
	ToBeStatus        string // P0024
	SaProportion      string // P0050 Y/N (Trad Only)
	LiquidatedIlpFund string // P0050 Y/N (ILP  Only)
	RecoverFromFund   string // P0050 Y/N (ILP  Only)
	LiquidFundCode    string // P0050 FUNDCODE
}

func (m *P0069Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0069Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

// Fund Switch Rules

type P0070Data struct {
	SwitchFeeBasis string // P0050 - Fixed/Percentage
	FreeSwitches   uint
	FeeAmount      float64
	FeePercentage  float64
}

func (m *P0070Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0070Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type P0071Data struct {
	P0071Array []P0071
}
type P0071 struct {
	BenDataType string // P0050
	ManOrOpt    string // P0050
}

func (m *P0071Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0071Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type P0072Data struct {
	MinLoanAmount           float64
	MaxLoanPercentage       int
	LoanInterestType        string // P0050
	IntPayableFreq          string // P0050
	RateOfInterest          float64
	StampDutyRate           float64
	LoanCapitalization      string // LOANCAP // P0050
	CapitalizationFrequency string // P0050
	ToleranceAmount         uint
	PrevLoanToBeClosed      string
	AllowCapDuringLoanYear  string // P0050 Yes/No
}

func (m *P0072Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0072Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type P0073Data struct {
	P0073Array []P0073
}
type P0073 struct {
	NoOfYears  float64
	LBFunction string // P0050
}

func (m *P0073Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0073Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

// Plan Params

type P0074Data struct {
	BenefitPlanType   string // P0050 STD, SPL, PRM [single select]
	PlanMaxLives      int64  //2/5/10/15/…  [single select]
	PlanLARelations   string // P0050 self, spouse/live-in partner, Parents, Kids & Others [multi select]
	BenefitPlanSA     string // P0050 3L,5L,10L  [multi select]
	RestorePlanSA     string //P0050(YESNO) Yes/No
	PlanPremAge       string
	PlanDiscountTypes string //P0050 Preferred Lives, Loyalty, Digital, Floater [multi select]
	CoPayRate         float64
	CoPayMinAmount    float64
	WPinMonths        int64 // Waiting Period
	AbroadWPinMonths  int64
}

func (m *P0074Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0074Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type P0075Data struct {
	PlanBenefits []P0075
}
type P0075 struct {
	BenefitCode      string //P0050
	BenefitUnit      float64
	BenefitBasis     string //P0050 % of SAPD, % of SA, FlatAmount,
	BenefitPlanCover string //P0050
	PlanBenefitGroup string //P0050 PlanBenefitGroup
}

func (m *P0075Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0075Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type P0076Data struct {
	ClaimsAllowedLimits []P0076
}
type P0076 struct {
	LimitCode string //P0050
	NoOfTimes uint
}

func (m *P0076Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0076Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type P0077Data struct {
	PlanMaxBenefits []P0077
}
type P0077 struct {
	BenefitCode      string //P0050
	MaxBenefitAmount float64
	MaxBenefitUnit   float64
	MaxBenefitBasis  string
}

func (m *P0077Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0077Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type P0078Data struct {
	PlanLARelationRules []P0078
}
type P0078 struct {
	PlanLARelations  string // P0050 self, spouse/live-in partner, Parents, Kids & Others [single select]
	NoOfLivesAllowed uint
	Underwriting     string // P0050 Yes/No
	PremLACode       string //P0050

}

func (m *P0078Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0078Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type P0079Data struct {
	PlanBenefitDiscounts []P0079
}
type P0079 struct {
	DiscountType  string // P0050(PlanDiscountTypes) Preferred Health, Loyalty, Digital, Floater
	DiscountRate  float64
	DiscountBasis string //P0050 % of BasePremium, FlatAmount, % of SA

}

func (m *P0079Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0079Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type P0080Data struct {
	PremRateCodes []P0080
}
type P0080 struct {
	SumAssured   float64
	LACode1      string // P0050 PremLACode
	LACount1     uint
	LACode2      string // P0050 PremLACode
	LACount2     uint
	LACode3      string // P0050 PremLACode
	LACount3     uint
	PremRateCode string //

}

func (m *P0080Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0080Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type P0081Data struct {
	AgeFrom uint
	AgeTo   uint
}

func (m *P0081Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0081Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

// Policy Addl Data  [same like P0071 Benefit Addl Data]

type P0082Data struct {
	P0082Array []P0082
}
type P0082 struct {
	PolDataType string // P0050
	ManOrOpt    string // P0050
}

func (m *P0082Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0082Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

// Cola Parameter

type P0083Data struct {
	NewOrExistId     string //P0050 (Existing/New)
	SimpleOrCompound string // P0050
}

func (m *P0083Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0083Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Cola Rates

type P0084Data struct {
	Rates []P0084
}
type P0084 struct {
	Yrs             uint
	PercentIncrease float64
}

func (m *P0084Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *P0084Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Workflow Param
type W0001Data struct {
	RequestType     string  // P0050
	Department      string  // W0006
	Team            string  // W0008
	SLADuration     float64 // W0009
	SLADurationType string  // W0009
}

func (m *W0001Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0001Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type W0002Data struct {
	TaskType        string  // P0050
	Department      string  // W0006
	Team            string  // W0008
	SLADuration     float64 // W0009
	SLADurationType string  // W0009
	ReassignedBy    string  // assignedto, manager, head
}

func (m *W0002Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0002Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type W0003Data struct {
	ActionType      string  // P0050
	Department      string  // W0006
	Team            string  // W0008
	SLADuration     float64 // W0009
	SLADurationType string  // W0009
	TranCode        string  // table:Transaction
	ReassignedBy    string  // assignedto, manager, head
}

func (m *W0003Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0003Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type W0004Data struct {
	FieldArray []W0004
}
type W0004 struct {
	Action          string // W0003
	SeqNo           float64
	MandInd         string  //P0050
	SLADuration     float64 // W0009
	SLADurationType string  // W0009
	Red             float64
	Amber           float64
	Green           float64
}

func (m *W0004Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0004Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type W0005Data struct {
	FieldArray []W0005
}
type W0005 struct {
	Task            string // W0002
	SeqNo           float64
	MandInd         string  //P0050
	SLADuration     float64 // W0009
	SLADurationType string  // W0009
	Red             float64
	Amber           float64
	Green           float64
}

func (m *W0005Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0005Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type W0006Data struct {
	DepartmentHead uint // table:users
	MaxNoOfTeams   uint
}

func (m *W0006Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0006Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type W0007Data struct {
	TeamLead     uint // table:users
	MaxNoOfUsers uint
}

func (m *W0007Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0007Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type W0008Data struct {
	FieldArray []W0008
}
type W0008 struct {
	Team string // W0007
}

func (m *W0008Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0008Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type W0009Data struct {
	SLAMethod string // sla calc program name
}

func (m *W0009Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0009Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

type W0010Data struct {
	FieldArray []W0010
}
type W0010 struct {
	NextStatus string // W0009
}

func (m *W0010Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0010Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

type W0011Data struct {
	ReminderMethod string // reminder method trigger program name
}

func (m *W0011Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0011Data) GetFormattedData(datamap map[string]string) map[string]interface{} {
	return nil

}

// Task Status
type W0012Data struct {
	FieldArray []W0012
}
type W0012 struct {
	TaskStatus string // p0050
}

func (m *W0012Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0012Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Action Status
type W0013Data struct {
	FieldArray []W0013
}
type W0013 struct {
	ActionStatus string // p0050
}

func (m *W0013Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0013Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Request Status
type W0014Data struct {
	FieldArray []W0014
}
type W0014 struct {
	RequestStatus string // p0050
}

func (m *W0014Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0014Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Task Completed Status
type W0015Data struct {
	FieldArray []W0015
}
type W0015 struct {
	TaskCompletedStatus string // p0050
}

func (m *W0015Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0015Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}

// Request Completed Status
type W0016Data struct {
	FieldArray []W0016
}
type W0016 struct {
	RequestCompletedStatus string // p0050
}

func (m *W0016Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *W0016Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

	return nil

}
