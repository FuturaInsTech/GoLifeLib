package types

import (
	"encoding/json"
	"fmt"
)

type Extradata interface {
	// Methods
	ParseData(map[string]interface{})
	GetFormattedData(datamap map[string]string) map[string]interface{}
}

// Q0005 Structure T5688
type Q0005Data struct {
	FreeLookDays          int
	MaxLives              int
	MinLives              int
	MinSurrMonths         int
	ProductFamily         string
	ReinstatementMonth    int
	Renewable             string
	Single                string
	Frequencies           [4]string // M Q H Y
	ContractCurr          []string  // INR USD SGD HKD LKR PKR DON IDR
	BillingCurr           []string
	ComponentAddAtAnyTime string // Policy Anniversary or Any Time N / Y
	FuturePremAdj         string //Y or N
	FuturePremAdjYrs      int    // eg., 3 Yrs
	LapsedDays            int
	BillingLeadDays       int
	LapseInterest         float64
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

	if datamap["function"] == "BillingCurr" {
		resp := make(map[string]interface{})
		allowedbilling := make([]string, 0)
		for i := 0; i < len(m.BillingCurr); i++ {
			if m.BillingCurr[i] == "" {
				break
			}

			allowedbilling = append(allowedbilling, m.BillingCurr[i])

		}
		resp["AllowedBillingCurriencies"] = allowedbilling
		return resp
	} else if datamap["function"] == "ContractCurr" {
		resp := make(map[string]interface{})
		contractcurr := make([]string, 0)
		for i := 0; i < len(m.ContractCurr); i++ {
			if m.ContractCurr[i] == "" {
				break
			}

			contractcurr = append(contractcurr, m.ContractCurr[i])

		}
		resp["AllowedContractCurriencies"] = contractcurr
		return resp
	} else if datamap["function"] == "Frequencies" {
		resp := make(map[string]interface{})
		allowedfreq := make([]string, 0)
		for i := 0; i < len(m.Frequencies); i++ {
			if m.Frequencies[i] == "" {
				break
			}

			allowedfreq = append(allowedfreq, m.Frequencies[i])

		}
		resp["AllowedFrequencies"] = allowedfreq
		return resp
	} else {
		return nil
	}

}

// /
type Q0006Data struct {
	AgeCalcMethod      string // T or P
	AnnMethod          string // Anniversary Method
	AnnuityMethod      string //Annuity Method
	CommMethod         string // Commission
	DeathType          string // 1 SA 2 FV 3 GT(SA/FV) 4. SA + FV
	DeathMethod        string // DC****
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
	NFOMethod          string // NFO Method
	PartSurrMethod     string // Part Surrender Method
	PremInc            string // Premium Increase Allowed
	PremIncYrs         uint   // No of Yrs in Number
	PremiumMethod      string // Premium Method
	RevBonus           string // Reversionary Bonus Method
	SBType             string // Either Age Based or Term Based
	SBMethod           string // Q00012 OR Q0013
	SurrMethod         string
	TBonus             string
	ULDeductFrequency  string
	GSVMethod          string
	SSVMethod          string
	BSVMethod          string
	DivMethod          string
	DivIMethod         string
	Mortalities        []Q0006M // Smoker/Non Smoker/Combined - Array
	PremCalcType       string   // Either Age Based / PPT Based
	DiscType           string   // S for SA P for Pemium
	DiscMethod         string   // DM001 or DM002
	FrqMethod          string   // Frequency Factor
	WaivMethod         string   // Waiver Method Q0020
	// Unit Linked Components
	ULALMethod         string    // UL Prem Allocation Method  Q0021
	ULMortFreq         string    // UL Mortality Deduction Frequency
	ULMortCalcType     string    // 1 - SAR 2 - SA, 3 - Fund + SA  Q0022
	ULMortDeductMethod string    // UL Mortality Deudction Method Q0022 Attained Age
	ULFeeFreq          string    // UL Charges Deduction Frequency
	ULFeeType          string    // 1 SA Based 2 Annualised Premium 3 Fund Value
	ULFeeMethod        string    // UL Fee Method
	ULFundRules        string    // UL Fund Rules
	MrtaMethod         string    // MRTA Method
	MrtaInterest       []float64 // MRTA Interest - Array
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
	SBRates []Q0012
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
	SBRates []Q0013
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
	SABand []Q0017
}
type Q0017 struct {
	SA       uint
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

type Q0021Data struct {
	ALBand []Q0021
}

type Q0021 struct {
	Months     uint
	Percentage float64
}

func (m *Q0021Data) ParseData(datamap map[string]interface{}) {
	jsonStr, err := json.Marshal(datamap)

	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct

	if err := json.Unmarshal(jsonStr, &m); err != nil {
		fmt.Println(err)
	}

}

func (m *Q0021Data) GetFormattedData(datamap map[string]string) map[string]interface{} {

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
	GST []Q0023
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
	BIRates []Q0024
}
type Q0024 struct {
	BIType string
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
	PPT  uint
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

// Collection Bank Account
type P0030Data struct {
	BankAccount string
	GLAccount   string
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
	UWRules []P0032
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
	TemplateName         string
	SMSAllowed           string
	EmailAllowed         string
	WhatsAppAllowed      string
	AgentSMSAllowed      string
	AgentEmailAllowed    string
	AgentWhatsAppAllowed string
	CompanyEmail         string
	CompanyPhone         string
	DepartmentName       string
	DepartmentHead       string
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
	SABands []P0041
}

type P0041 struct {
	Ages []P0041Age
}
type P0041Age struct {
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
	BIYrInterval uint
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
	BankRequired string // Y or N
	BankCode     string // Bank Account No AND IFSC
	BankAccount  string
}

type P0055 struct {
	ExtractionDate string
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
	ExtractionDates      []P0055 // 5, 10,15,25,28, 30,31
	NetCollection        string  // Net Colleciton/Gross Collection (N/G)
	CollectionFee        float64 // Flat Fee for Each Collection
	CollectionPercentage float64 // Collection Fee
	AccountEntry         string  // Pass Accounting Entries (Y/N)
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
