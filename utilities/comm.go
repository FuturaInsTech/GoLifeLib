package utilities

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"
	"gorm.io/gorm"
)

// #104
// Create Communication (New Version with Rollback)
//
// # This function, Create Communication Records by getting input values as Company ID, History Code, Tranno, Date of Transaction, Policy Id, Client Id, Address Id, Receipt ID . Quotation ID, Agency ID
// 10 Input Variables
// # It returns success or failure.  Successful records written in Communciaiton Table
//
// ©  FuturaInsTech
func CreateCommunicationsN(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB, iBenefit uint, iPa uint, iClientWork uint) error {

	var communication models.Communication
	var iKey string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	//utilities.LetterCreate(int(iCompany), uint(iPolicy), iHistoryCode, createreceipt.CurrentDate, idata)
	iTransaction := iHistoryCode
	iReceiptTranCode := "H0034"
	iReceiptFor := ""

	if iReceipt != 0 {
		var receipt models.Receipt
		result := txn.Find(&receipt, "company_id = ? and id = ?", iCompany, iReceipt)
		if result.Error != nil {
			return result.Error
		}
		iReceiptFor = receipt.ReceiptFor
	}

	if iPolicy != 0 {
		var policy models.Policy
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.Error != nil {
			return result.Error
		}
		communication.CompanyID = uint(iCompany)
		communication.AgencyID = policy.AgencyID
		communication.ClientID = policy.ClientID
		communication.PolicyID = policy.ID
		communication.Tranno = policy.Tranno
		communication.EffectiveDate = policy.PRCD
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPolicy
		if iTransaction == iReceiptTranCode {
			iKey = iTransaction + iReceiptFor
		} else {
			iKey = iTransaction + policy.PProduct
		}
	}

	if iPolicy == 0 && iTransaction == iReceiptTranCode && iPa != 0 {
		var payingauth models.PayingAuthority
		result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)
		if result.Error != nil {
			return result.Error
		}

		communication.CompanyID = uint(iCompany)
		communication.AgencyID = 0
		communication.ClientID = payingauth.ClientID
		communication.PolicyID = 0
		communication.Tranno = 0
		communication.EffectiveDate = iDate
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPa
		iKey = iTransaction + iReceiptFor
	}

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

			communication.AgentEmailAllowed = p0033data.AgentEmailAllowed
			communication.AgentSMSAllowed = p0033data.AgentSMSAllowed
			communication.AgentWhatsAppAllowed = p0033data.AgentWhatsAppAllowed
			communication.EmailAllowed = p0033data.EmailAllowed
			communication.SMSAllowed = p0033data.SMSAllowed
			communication.WhatsAppAllowed = p0033data.WhatsAppAllowed
			communication.DepartmentHead = p0033data.DepartmentHead
			communication.DepartmentName = p0033data.DepartmentName
			communication.CompanyPhone = p0033data.CompanyPhone
			communication.CompanyEmail = p0033data.CompanyEmail

			communication.TemplateName = iKey
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
					oData := GetCompanyData(iCompany, iDate, txn)
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData := GetClientData(iCompany, iClient, txn)
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData := GetAddressData(iCompany, iAddress, txn)
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData := GetPolicyData(iCompany, iPolicy, txn)
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData := GetBenefitData(iCompany, iPolicy, txn)
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData := GetSurBData(iCompany, iPolicy, txn)
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData := GetMrtaData(iCompany, iPolicy, txn)
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData := GetReceiptData(iCompany, iReceipt, txn)
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData := GetSaChangeData(iCompany, iPolicy, txn)
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData := GetCompAddData(iCompany, iPolicy, txn)
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData := GetSurrHData(iCompany, iPolicy, txn)
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData := GetDeathData(iCompany, iPolicy, txn)
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData := GetMatHData(iCompany, iPolicy, txn)
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData := GetSurvBPay(iCompany, iPolicy, iTranno, txn)
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData := GetExpi(iCompany, iPolicy, iTranno, txn)
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData := GetBonusVals(iCompany, iPolicy, txn)
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData := GetAgency(iCompany, iAgency, txn)
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData := GetNomiData(iCompany, iPolicy, txn)
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData := GetGLData(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData := GetIlpSummaryData(iCompany, iPolicy, txn)
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData := GetIlpAnnsummaryData(iCompany, iPolicy, iHistoryCode, txn)
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData := GetIlpTranctionData(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData := GetPremTaxGLData(iCompany, iPolicy, iFromDate, iToDate, txn)
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData := GetIlpFundSwitchData(iCompany, iPolicy, iTranno, txn)
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData := GetPHistoryData(iCompany, iPolicy, iHistoryCode, iDate, txn)
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData := GetIlpFundData(iCompany, iPolicy, iBenefit, iDate, txn)
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData := GetPPolicyData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData := GetPBenefitData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					fmt.Println(oData) // Dummy to avoid compilation error
				case oLetType == "29":
					oData := GetPayingAuthorityData(iCompany, iPa, txn)
					resultMap["PrevBenefit"] = oData
				case oLetType == "30":
					oData := GetClientWorkData(iCompany, iClientWork, txn)
					resultMap["ClientWork"] = oData
				case oLetType == "36":
					oData := GetReqData(iCompany, iPolicy, iClient, txn)
					for _, item := range oData {
						for key, value := range item.(map[string]interface{}) {
							resultMap[key] = value
						}
					}
				case oLetType == "37":
					oData := PolicyDepData(iCompany, iPolicy, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "38":
					oData := PolAgntChData(iCompany, iPolicy, iAgency, iClient, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "39":
					oData := GetBankData(iCompany, iClientWork, txn)
					resultMap["BankData"] = oData
				case oLetType == "40":
					iKey := iReceipt
					oData := GetPaymentData(iCompany, iPolicy, iKey, txn)
					for key, value := range oData {
						resultMap[key] = value
					}

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
			// New Changes for Online Print and Email Trigger
			if p0033data.Online == "Y" {
				err := GetReportforOnline(communication, p0033data.TemplateName, txn)
				if err != nil {
					log.Fatalf("Failed to generate report: %v", err)
				}
			}
			if p0033data.SMSAllowed == "Y" {
				err := SendSMSTwilio(communication.CompanyID, communication.ClientID, p0033data.TemplateName, communication.EffectiveDate, p0033data.SMSBody, txn)
				if err != nil {
					log.Fatalf("Failed to send SMS: %v", err)
				}
			}
			communication.Print = "Y"
			communication.PrintDate = iDate
			communication.UpdatedID = 1
			communication.ID = 0
			// New Changes Ended
			results := txn.Create(&communication)

			if results.Error != nil {
				return results.Error
			}

		}
	}
	return nil
}

// #104
// Create Communication
//
// # This function, Create Communication Records by getting input values as Company ID, History Code, Tranno, Date of Transaction, Policy Id, Client Id, Address Id, Receipt ID . Quotation ID, Agency ID
// 10 Input Variables
// # It returns success or failure.  Successful records written in Communciaiton Table
//
// ©  FuturaInsTech
func CreateCommunications(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, iBenefit uint, iPa uint, iClientWork uint) error {

	var communication models.Communication
	var iP0033Key string
	var iP0034Key string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data
	txn := initializers.DB.Begin()

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data

	var policy models.Policy
	if iPolicy != 0 {
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.Error != nil {
			return result.Error
		}
	}
	var payingauth models.PayingAuthority
	if iPa != 0 {
		result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)
		if result.Error != nil {
			return result.Error
		}
	}

	iReceiptTranCode := "H0034"
	iReceiptFor := ""
	if iHistoryCode == iReceiptTranCode {
		var receipt models.Receipt
		result := txn.Find(&receipt, "company_id = ? and id = ?", iCompany, iReceipt)
		if result.Error != nil {
			return result.Error
		}
		iReceiptFor = receipt.ReceiptFor
		iP0034Key = iHistoryCode + iReceiptFor
	}
	if iReceiptFor == "" {
		communication.CompanyID = uint(iCompany)
		communication.AgencyID = policy.AgencyID
		communication.ClientID = policy.ClientID
		communication.PolicyID = policy.ID
		communication.Tranno = policy.Tranno
		communication.EffectiveDate = policy.PRCD
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPolicy
		iP0034Key = iHistoryCode + policy.PProduct
	}

	if iReceiptFor == "01" {
		communication.CompanyID = uint(iCompany)
		communication.AgencyID = policy.AgencyID
		communication.ClientID = policy.ClientID
		communication.PolicyID = policy.ID
		communication.Tranno = policy.Tranno
		communication.EffectiveDate = policy.PRCD
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPolicy
	}

	if iReceiptFor == "02" {
		communication.CompanyID = uint(iCompany)
		communication.AgencyID = 0
		communication.ClientID = payingauth.ClientID
		communication.PolicyID = 0
		communication.Tranno = 0
		communication.EffectiveDate = iDate
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPa
	}

	if iReceiptFor == "03" {
		communication.CompanyID = uint(iCompany)
		communication.AgencyID = 0
		communication.ClientID = iClient
		communication.PolicyID = 0
		communication.Tranno = 0
		communication.EffectiveDate = iDate
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iClient
	}

	err1 := GetItemD(int(iCompany), "P0034", iP0034Key, iDate, &extradatap0034)
	if err1 != nil {
		iP0034Key = iHistoryCode
		err1 = GetItemD(int(iCompany), "P0034", iP0034Key, iDate, &extradatap0034)
		if err1 != nil {
			return err1
		}
	}

	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iP0033Key = p0034data.Letters[i].Templates
			err := GetItemD(int(iCompany), "P0033", iP0033Key, iDate, &extradatap0033)
			if err != nil {
				return err
			}

			communication.AgentEmailAllowed = p0033data.AgentEmailAllowed
			communication.AgentSMSAllowed = p0033data.AgentSMSAllowed
			communication.AgentWhatsAppAllowed = p0033data.AgentWhatsAppAllowed
			communication.EmailAllowed = p0033data.EmailAllowed
			communication.SMSAllowed = p0033data.SMSAllowed
			communication.WhatsAppAllowed = p0033data.WhatsAppAllowed
			communication.DepartmentHead = p0033data.DepartmentHead
			communication.DepartmentName = p0033data.DepartmentName
			communication.CompanyPhone = p0033data.CompanyPhone
			communication.CompanyEmail = p0033data.CompanyEmail

			communication.TemplateName = iP0033Key
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
					oData := GetCompanyData(iCompany, iDate, txn)
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData := GetClientData(iCompany, iClient, txn)
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData := GetAddressData(iCompany, iAddress, txn)
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData := GetPolicyData(iCompany, iPolicy, txn)
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData := GetBenefitData(iCompany, iPolicy, txn)
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData := GetSurBData(iCompany, iPolicy, txn)
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData := GetMrtaData(iCompany, iPolicy, txn)
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData := GetReceiptData(iCompany, iReceipt, txn)
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData := GetSaChangeData(iCompany, iPolicy, txn)
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData := GetCompAddData(iCompany, iPolicy, txn)
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData := GetSurrHData(iCompany, iPolicy, txn)
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData := GetDeathData(iCompany, iPolicy, txn)
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData := GetMatHData(iCompany, iPolicy, txn)
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData := GetSurvBPay(iCompany, iPolicy, iTranno, txn)
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData := GetExpi(iCompany, iPolicy, iTranno, txn)
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData := GetBonusVals(iCompany, iPolicy, txn)
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData := GetAgency(iCompany, iAgency, txn)
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData := GetNomiData(iCompany, iPolicy, txn)
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData := GetGLData(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData := GetIlpSummaryData(iCompany, iPolicy, txn)
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData := GetIlpAnnsummaryData(iCompany, iPolicy, iHistoryCode, txn)
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData := GetIlpTranctionData(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData := GetPremTaxGLData(iCompany, iPolicy, iFromDate, iToDate, txn)
					resultMap["GLData"] = oData
				case oLetType == "24":
					oData := GetIlpFundSwitchData(iCompany, iPolicy, iTranno, txn)
					resultMap["SwitchData"] = oData
				case oLetType == "25":
					oData := GetPHistoryData(iCompany, iPolicy, iHistoryCode, iDate, txn)
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData := GetIlpFundData(iCompany, iPolicy, iBenefit, iDate, txn)
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData := GetPPolicyData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData := GetPBenefitData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					resultMap["PrevBenefit"] = oData
				case oLetType == "29":
					oData := GetPayingAuthorityData(iCompany, iPa, txn)
					resultMap["PaData"] = oData
				case oLetType == "30":
					oData := GetClientWorkData(iCompany, iClientWork, txn)
					resultMap["ClientWork"] = oData
				// case oLetType == "36":
				// 	oData := GetReqData(iCompany, iPolicy)
				// 	resultMap["ReqWork"] = oData
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

// 194
// CreateCommunicationL for loan
// Inputs: CompanyID, PolicyID, ...
//
// # Outputs: error
//
// ©  FuturaInsTech
func CreateCommunicationsL(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB, iBenefit uint, iPa uint, iClientWork uint, iAmount1 float64, iAmount2 float64, iNo1 uint, iNo2 uint) error {

	var communication models.Communication
	var iKey string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data

	//utilities.LetterCreate(int(iCompany), uint(iPolicy), iHistoryCode, createreceipt.CurrentDate, idata)
	iTransaction := iHistoryCode
	iReceiptTranCode := "H0034"
	iReceiptFor := ""

	if iReceipt != 0 {
		var receipt models.Receipt
		result := txn.Find(&receipt, "company_id = ? and id = ?", iCompany, iReceipt)
		if result.Error != nil {
			return result.Error
		}
		iReceiptFor = receipt.ReceiptFor
	}

	if iPolicy != 0 {
		var policy models.Policy
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.Error != nil {
			return result.Error
		}
		communication.CompanyID = uint(iCompany)
		communication.AgencyID = policy.AgencyID
		communication.ClientID = policy.ClientID
		communication.PolicyID = policy.ID
		communication.Tranno = policy.Tranno
		communication.EffectiveDate = policy.PRCD
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPolicy
		if iTransaction == iReceiptTranCode {
			iKey = iTransaction + iReceiptFor
		} else {
			iKey = iTransaction + policy.PProduct
		}
	}

	if iPolicy == 0 && iTransaction == iReceiptTranCode && iPa != 0 {
		var payingauth models.PayingAuthority
		result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)
		if result.Error != nil {
			return result.Error
		}

		communication.CompanyID = uint(iCompany)
		communication.AgencyID = 0
		communication.ClientID = payingauth.ClientID
		communication.PolicyID = 0
		communication.Tranno = 0
		communication.EffectiveDate = iDate
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPa
		iKey = iTransaction + iReceiptFor
	}

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

			communication.AgentEmailAllowed = p0033data.AgentEmailAllowed
			communication.AgentSMSAllowed = p0033data.AgentSMSAllowed
			communication.AgentWhatsAppAllowed = p0033data.AgentWhatsAppAllowed
			communication.EmailAllowed = p0033data.EmailAllowed
			communication.SMSAllowed = p0033data.SMSAllowed
			communication.WhatsAppAllowed = p0033data.WhatsAppAllowed
			communication.DepartmentHead = p0033data.DepartmentHead
			communication.DepartmentName = p0033data.DepartmentName
			communication.CompanyPhone = p0033data.CompanyPhone
			communication.CompanyEmail = p0033data.CompanyEmail

			communication.TemplateName = iKey
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
					oData := GetCompanyData(iCompany, iDate, txn)
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData := GetClientData(iCompany, iClient, txn)
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData := GetAddressData(iCompany, iAddress, txn)
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData := GetPolicyData(iCompany, iPolicy, txn)
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData := GetBenefitData(iCompany, iPolicy, txn)
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData := GetSurBData(iCompany, iPolicy, txn)
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData := GetMrtaData(iCompany, iPolicy, txn)
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData := GetReceiptData(iCompany, iReceipt, txn)
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData := GetSaChangeData(iCompany, iPolicy, txn)
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData := GetCompAddData(iCompany, iPolicy, txn)
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData := GetSurrHData(iCompany, iPolicy, txn)
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData := GetDeathData(iCompany, iPolicy, txn)
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData := GetMatHData(iCompany, iPolicy, txn)
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData := GetSurvBPay(iCompany, iPolicy, iTranno, txn)
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData := GetExpi(iCompany, iPolicy, iTranno, txn)
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData := GetBonusVals(iCompany, iPolicy, txn)
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData := GetAgency(iCompany, iAgency, txn)
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData := GetNomiData(iCompany, iPolicy, txn)
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData := GetGLData(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData := GetIlpSummaryData(iCompany, iPolicy, txn)
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData := GetIlpAnnsummaryData(iCompany, iPolicy, iHistoryCode, txn)
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData := GetIlpTranctionData(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData := GetPremTaxGLData(iCompany, iPolicy, iFromDate, iToDate, txn)
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData := GetIlpFundSwitchData(iCompany, iPolicy, iTranno, txn)
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData := GetPHistoryData(iCompany, iPolicy, iHistoryCode, iDate, txn)
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData := GetIlpFundData(iCompany, iPolicy, iBenefit, iDate, txn)
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData := GetPPolicyData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData := GetPBenefitData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					fmt.Println(oData) // Dummy to avoid compilation error
				case oLetType == "29":
					oData := GetPayingAuthorityData(iCompany, iPa, txn)
					resultMap["PrevBenefit"] = oData
				case oLetType == "30":
					oData := GetClientWorkData(iCompany, iClientWork, txn)
					resultMap["ClientWork"] = oData
				case oLetType == "31":
					oData := GetLoanData(iCompany, iPolicy, iDate, iAmount1, txn)
					resultMap["LoanData"] = oData
				case oLetType == "32":
					oData := GetAllLoanInterestData(iCompany, iPolicy, iDate, txn)
					resultMap["LoanInterestData"] = oData
				case oLetType == "33":
					oData := LoanCapData(iCompany, iPolicy, iDate, iFromDate, iToDate, iAmount1, iAmount2, iNo1, txn)
					resultMap["LoanCap"] = oData
				case oLetType == "34":
					oData := LoanBillData(iCompany, iPolicy, iDate, txn)
					resultMap["LoanBillData"] = oData
				case oLetType == "35":
					oData := LoanBillsInterestData(iCompany, iPolicy, iNo1, iAmount1, txn)
					resultMap["LoanBillsInterest"] = oData
				case oLetType == "98":
					resultMap["BatchData"] = batchData

				case oLetType == "99":
					resultMap["SignData"] = signData
				default:

				}
			}

			if p0033data.Online == "Y" {
				err := GetReportforOnline(communication, p0033data.TemplateName, txn)
				if err != nil {
					log.Fatalf("Failed to generate report: %v", err)
				}
			}
			if p0033data.SMSAllowed == "Y" {
				err := SendSMSTwilio(communication.CompanyID, communication.ClientID, p0033data.TemplateName, communication.EffectiveDate, p0033data.SMSBody, txn)
				if err != nil {
					log.Fatalf("Failed to send SMS: %v", err)
				}
			}
			communication.Print = "Y"
			communication.PrintDate = iDate
			communication.UpdatedID = 1

			communication.ExtractedData = resultMap
			communication.PDFPath = p0034data.Letters[i].PdfLocation
			communication.TemplatePath = p0034data.Letters[i].ReportTemplateLocation
			communication.ID = 0
			results := txn.Create(&communication)

			if results.Error != nil {
				return results.Error
			}

		}
	}
	return nil
}

// #86
// GetCompany Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Company Details as an Interface
//
// ©  FuturaInsTech
func GetCompanyData(iCompany uint, iDate string, txn *gorm.DB) []interface{} {
	companyarray := make([]interface{}, 0)
	var company models.Company
	txn.Find(&company, "id = ?", iCompany)

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

// #88
// GetAddressData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Address Details as an Interface
//
// ©  FuturaInsTech
func GetAddressData(iCompany uint, iAddress uint, txn *gorm.DB) []interface{} {
	addressarray := make([]interface{}, 0)
	var address models.Address

	txn.Find(&address, "company_id = ? and id = ?", iCompany, iAddress)
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
func GetPolicyData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {
	policyarray := make([]interface{}, 0)
	var policy models.Policy
	result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return nil
	}
	_, oStatus, _ := GetParamDesc(policy.CompanyID, "P0024", policy.PolStatus, 1)
	_, oFreq, _ := GetParamDesc(policy.CompanyID, "Q0009", policy.PFreq, 1)
	_, oProduct, _ := GetParamDesc(policy.CompanyID, "Q0005", policy.PProduct, 1)
	_, oBillCurr, _ := GetParamDesc(policy.CompanyID, "P0023", policy.PBillCurr, 1)
	_, oContCurr, _ := GetParamDesc(policy.CompanyID, "P0023", policy.PContractCurr, 1)
	_, oBillingType, _ := GetParamDesc(policy.CompanyID, "P0055", policy.BillingType, 1)

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)
	gracedate := AddLeadDays(policy.PaidToDate, q0005data.LapsedDays)
	premduedates := GetPremDueDates(policy.PRCD, policy.PFreq)
	iAnnivDate := Date2String(GetNextDue(policy.AnnivDate, "Y", "R"))

	var benefitenq []models.Benefit

	txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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
		"BillingType":        oBillingType,
		"PaingAuthorityId":   IDtoPrint(policy.PayingAuthority),

		// "PUWDate":DateConvert(policy.PUWDate),
	}
	policyarray = append(policyarray, resultOut)

	return policyarray
}

// #90
// GetBenefitData - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Benefit Details as an Interface
//
// ©  FuturaInsTech
func GetBenefitData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {
	var policyenq models.Policy
	var benefit []models.Benefit
	var clientenq models.Client
	var addressenq models.Address
	txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	paidToDate := policyenq.PaidToDate
	nextDueDate := policyenq.NxtBTDate
	txn.Find(&benefit, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	benefitarray := make([]interface{}, 0)

	for k := 0; k < len(benefit); k++ {
		iCompany := benefit[k].CompanyID
		_, oGender, _ := GetParamDesc(iCompany, "P0001", benefit[k].BGender, 1)
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", benefit[k].BCoverage, 1)
		_, oStatus, _ := GetParamDesc(iCompany, "P0024", benefit[k].BStatus, 1)

		clientname := GetName(iCompany, benefit[k].ClientID)
		txn.Find(&clientenq, "company_id = ? and id = ?", iCompany, benefit[k].ClientID)
		txn.Find(&addressenq, "company_id = ? and client_id = ?", iCompany, clientenq.ID)
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
func GetSurBData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {
	var survb []models.SurvB
	txn.Find(&survb, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	var benefitenq models.Benefit
	txn.Find(&benefitenq, "company_id = ? and policy_id =? and id = ?", iCompany, iPolicy, survb[0].BenefitID)
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
func GetMrtaData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {
	var mrtaenq []models.Mrta
	txn.Find(&mrtaenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)

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
func GetReceiptData(iCompany uint, iReceipt uint, txn *gorm.DB) []interface{} {
	var receiptenq models.Receipt
	txn.Find(&receiptenq, "company_id = ? and id = ?", iCompany, iReceipt)
	amtinwords, csymbol := AmountinWords(receiptenq.CompanyID, receiptenq.AccAmount, receiptenq.AccCurry)
	receiptarray := make([]interface{}, 0)
	resultOut := map[string]interface{}{
		"ID":               IDtoPrint(receiptenq.ID),
		"CompanyID":        IDtoPrint(receiptenq.CompanyID),
		"Branch":           receiptenq.Branch,
		"AccCurry":         receiptenq.AccCurry,
		"AccAmount":        NumbertoPrint(receiptenq.AccAmount),
		"ReceiptFor":       receiptenq.ReceiptFor,
		"ReceiptRefNo":     IDtoPrint(receiptenq.ReceiptRefNo),
		"ReceiptDueDate":   DateConvert(receiptenq.ReceiptDueDate),
		"ClientID":         IDtoPrint(receiptenq.ClientID),
		"DateOfCollection": DateConvert(receiptenq.DateOfCollection),
		"BankAccountNo":    receiptenq.BankAccountNo,
		"BankReferenceNo":  receiptenq.BankReferenceNo,
		"TypeOfReceipt":    receiptenq.TypeOfReceipt,
		"ReceiptAmount":    receiptenq.ReceiptAmount,
		"AddressID":        IDtoPrint(receiptenq.AddressID),
		"AmountInWords":    amtinwords,
		"CurrSymbol":       csymbol,
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
func GetSaChangeData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {
	var sachangeenq []models.SaChange
	txn.Find(&sachangeenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)

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
func GetCompAddData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {
	var addcomp []models.Addcomponent
	txn.Find(&addcomp, "company_id = ? and policy_id = ?", iCompany, iPolicy)

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
func GetSurrHData(iCompany uint, iPolicy uint, txn *gorm.DB) interface{} {
	var surrhenq models.SurrH

	txn.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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

	txn.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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

// # 98 (Redundant) Not in Use
// GetDeath Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Death Details as an Interface
//
// ©  FuturaInsTech
// Not Required
func GetDeathData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {
	var surrhenq models.SurrH
	var surrdenq []models.SurrD
	txn.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	txn.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	surrarray := make([]interface{}, 0)

	return surrarray
}

// #98
// GetMatH Data - Printing Purpose Only (both header and detail)
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Maturity Header and Details Interface
// ©  FuturaInsTech
func GetMatHData(iCompany uint, iPolicy uint, txn *gorm.DB) interface{} {
	var mathenq models.MaturityH

	txn.Find(&mathenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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

	txn.Find(&matdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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
func GetSurvBPay(iCompany uint, iPolicy uint, iTranno uint, txn *gorm.DB) []interface{} {
	var survbenq models.SurvB
	txn.Find(&survbenq, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno)
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

// #102
// GetExpi Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Expiry Interface Information
// ©  FuturaInsTech
func GetExpi(iCompany uint, iPolicy uint, iTranno uint, txn *gorm.DB) []interface{} {
	var benefit []models.Benefit
	txn.Find(&benefit, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno)
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

// #100
// GetBonsusValues Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Bonus Values
// ©  FuturaInsTech
func GetBonusVals(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {

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
func GetAgency(iCompany uint, iAgency uint, txn *gorm.DB) []interface{} {

	agencyarray := make([]interface{}, 0)
	var agencyenq models.Agency
	var clientenq models.Client
	txn.Find(&agencyenq, "company_id  = ? and id = ?", iCompany, iAgency)

	txn.Find(&clientenq, "company_id = ? and id = ?", iCompany, agencyenq.ClientID)
	oAgentName := clientenq.ClientLongName + " " + clientenq.ClientShortName + " " + clientenq.ClientSurName

	var addressenq models.Address
	txn.Find(&addressenq, "company_id = ? and client_id = ?", iCompany, clientenq.ID)
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

// #183
// Get ClientWork Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Client Work Details as an Interface
//
// ©  FuturaInsTech
func GetClientWorkData(iCompany uint, iClientWork uint, txn *gorm.DB) []interface{} {
	clientworkarray := make([]interface{}, 0)
	var clientwork models.ClientWork

	txn.Find(&clientwork, "company_id = ? and id = ?", iCompany, iClientWork)
	resultOut := map[string]interface{}{
		"ID":            IDtoPrint(clientwork.ID),
		"ClientID":      IDtoPrint(clientwork.ClientID),
		"EmployerID":    IDtoPrint(clientwork.EmployerID),
		"PayRollNumber": clientwork.PayRollNumber,
		"Designation":   clientwork.Designation,
		"Department":    clientwork.Department,
		"Location":      clientwork.Location,
		"StartDate":     DateConvert(clientwork.StartDate),
		"EndDate":       clientwork.EndDate,
		"WorkType":      clientwork.WorkType,
	}
	clientworkarray = append(clientworkarray, resultOut)
	return clientworkarray
}

// #97
// GetNominee Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Nominee Details as an Interface
//
// ©  FuturaInsTech
func GetNomiData(iCompany uint, iPolicy uint, txn *gorm.DB) []interface{} {

	var nomenq []models.Nominee

	txn.Find(&nomenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	nomarray := make([]interface{}, 0)
	var clientenq models.Client
	var policyenq models.Policy
	txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	for k := 0; k < len(nomenq); k++ {
		txn.Find(&clientenq, "company_id = ? and id = ?", iCompany, nomenq[k].ClientID)
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

// # 140
// Get GL Data Printing Purpose Only
// Inputs: Company  Policy, From and To Date, History Code, GL Code,GL Sign
//
// Outputs : GL Data as interface
//
// ©  FuturaInsTech
func GetGLData(iCompany uint, iPolicy uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB) interface{} {
	var benefitenq []models.Benefit

	var covrcodes []string
	var covrnames []string

	txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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
		txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date <=?", iCompany, iPolicy, iFromDate, iToDate).Order("effective_date , tranno")
	} else if iGlHistoryCode != "" && iGlAccountCode == "" && iGlSign == "" {
		txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode).Order("history_code, effective_date , tranno")
	} else if iGlHistoryCode != "" && iGlAccountCode != "" && iGlSign == "" {
		txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ? and account_code like ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, "%"+iGlAccountCode+"%").Order("history_code, account_code, effective_date , tranno")
	} else if iGlHistoryCode != "" && iGlAccountCode != "" && iGlSign != "" {
		txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ? and account_code like ? and gl_sign = ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, "%"+iGlAccountCode+"%", iGlSign).Order("history_code, account_code, gl_sign, effective_date , tranno")
	} else if iGlHistoryCode == "" && iGlAccountCode != "" && iGlSign != "" {
		txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and account_code like ? and gl_sign = ?", iCompany, iPolicy, iFromDate, iToDate, "%"+iGlAccountCode+"%", iGlSign).Order("account_code, gl_sign, effective_date , tranno")
	} else if iGlHistoryCode == "" && iGlAccountCode != "" && iGlSign == "" {
		txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and account_code like ?", iCompany, iPolicy, iFromDate, iToDate, "%"+iGlAccountCode+"%").Order("account_code, effective_date , tranno")
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

// # 139
// Get ILP Summary Data Printing Purpose Only
// Inputs: Company and Policy
//
// Outputs : Summary Data as interface
//
// ©  FuturaInsTech
func GetIlpSummaryData(iCompany uint, iPolicy uint, txn *gorm.DB) interface{} {
	var ilpsummary []models.IlpSummary
	txn.Find(&ilpsummary, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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

// # 147
//
// # GetIlpAnnsummaryData - ILP Anniversary Summary Data extraction for Communications
//
// ©  FuturaInsTech
func GetIlpAnnsummaryData(iCompany uint, iPolicy uint, iHistoryCode string, txn *gorm.DB) interface{} {
	ilpannsumprevarray := make([]interface{}, 0)
	ilpannsumcurrarray := make([]interface{}, 0)
	var policyenq models.Policy
	txn.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	iAnnivDate := Date2String(GetNextDue(policyenq.AnnivDate, "Y", "R"))
	iPrevAnnivDate := Date2String(GetNextDue(iAnnivDate, "Y", "R"))

	var ilpannsumprev []models.IlpAnnSummary
	txn.Find(&ilpannsumprev, "company_id = ? and policy_id = ? and effective_date = ?", iCompany, iPolicy, iPrevAnnivDate)

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
	txn.Find(&ilpannsumcurr, "company_id = ? and policy_id = ? and effective_date = ?", iCompany, iPolicy, iAnnivDate)

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

// # 146
//
// # GetIlpTranctionData - ILP transaction Data extraction for Communications
//
// ©  FuturaInsTech
func GetIlpTranctionData(iCompany uint, iPolicy uint, iHistoryCode string, iDate string, txn *gorm.DB) []interface{} {
	var policyenq models.Policy
	txn.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	iAnnivDate := Date2String(GetNextDue(policyenq.AnnivDate, "Y", "R"))
	iPrevAnnivDate := Date2String(GetNextDue(iAnnivDate, "Y", "R"))
	var ilptranction []models.IlpTransaction
	if iHistoryCode == "B0103" {
		txn.Find(&ilptranction, "company_id = ? and policy_id = ? and ul_process_flag = ? and inv_non_inv_flag != ? and transaction_date >= ? and transaction_date < ?", iCompany, iPolicy, "C", "NI", iPrevAnnivDate, iAnnivDate).Order("fund_code, transaction_date , tranno")
	} else if iHistoryCode == "B0115" {
		txn.Find(&ilptranction, "company_id = ? and policy_id = ? and ul_process_flag = ? and inv_non_inv_flag != ? and transaction_date >= ? and transaction_date <= ?", iCompany, iPolicy, "C", "NI", iAnnivDate, iDate).Order("fund_code, transaction_date , tranno")
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

// # 141
// GetPremTaxGLData
// Extract PremTax Data  (Printing Purpose Only)
// Input:  Company, Policy, From and To Date
// Output: Interface
//
// ©  FuturaInsTech
func GetPremTaxGLData(iCompany uint, iPolicy uint, iFromDate string, iToDate string, txn *gorm.DB) interface{} {
	var benefitenq []models.Benefit
	var codesql string = ""
	var covrcodes []string
	var covrnames []string

	var acodearray []string

	var p0067data paramTypes.P0067Data
	var extradatap0067 paramTypes.Extradata = &p0067data

	txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
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
	txn.Find(&glmoves, "("+codesql+") and company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date <=? ", iCompany, iPolicy, iFromDate, iToDate).Order("account_code, gl_sign, effective_date , tranno")

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

// #158
// ILP Products Only.  GetIlpFundSwitchData (Printing Purpose)
// Input: Company, Policy No, Tranno
// Output: An Interface Record
//
// ©  FuturaInsTech
func GetIlpFundSwitchData(iCompany uint, iPolicy uint, iTranno uint, txn *gorm.DB) interface{} {
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

// #157
// GetPHistoryData  (Printing Purpose)
// Input: Company, Policy No, Transaction Code and Effective Date
// Output: An Interface Record (History Information)
//
// ©  FuturaInsTech
func GetPHistoryData(iCompany uint, iPolicy uint, iHistoryCode string, iDate string, txn *gorm.DB) []interface{} {
	var policyhistory []models.PHistory
	txn.Find(&policyhistory, "company_id = ? and policy_id = ?", iCompany, iPolicy)

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

// # 162
// GetIlpFundData - Get ILP Fund DAta  (Printing Purpose Only)
//
// Inputs: Company,
//
// # Outputs:
//
// ©  FuturaInsTech
func GetIlpFundData(iCompany uint, iPolicy uint, iBenefit uint, iDate string, txn *gorm.DB) interface{} {
	var ilpfund []models.IlpFund
	txn.Find(&ilpfund, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)

	var ibenfit models.Benefit
	txn.Find(&ibenfit, "company_id = ? and policy_id = ? and id = ?", iCompany, iPolicy, iBenefit)

	ilpfundtarray := make([]interface{}, 0)

	for k := 0; k < len(ilpfund); k++ {
		var p0061data paramTypes.P0061Data
		var extradatap0061 paramTypes.Extradata = &p0061data

		err := GetItemD(int(iCompany), "P0061", ilpfund[k].FundCode, iDate, &extradatap0061)

		if err != nil {
			shortCode := "GL442"
			longDesc, _ := GetErrorDesc(iCompany, 1, shortCode)
			return errors.New(shortCode + " : " + longDesc)

		}

		resultOut := map[string]interface{}{
			"ID":             IDtoPrint(ilpfund[k].ID),
			"CompanyID":      IDtoPrint(ilpfund[k].CompanyID),
			"BenefitID":      IDtoPrint(ilpfund[k].BenefitID),
			"PolicyID":       IDtoPrint(ilpfund[k].PolicyID),
			"FundCategory":   p0061data.FundCategory + " - " + GetP0050ItemCodeDesc(iCompany, "FUNDCATEGORY", 1, p0061data.FundCategory),
			"FundName":       ilpfund[k].FundCode + " - " + GetP0050ItemCodeDesc(iCompany, "FUNDCODE", 1, ilpfund[k].FundCode),
			"FundType":       ilpfund[k].FundType + " - " + GetP0050ItemCodeDesc(iCompany, "FUNDTYPE", 1, ilpfund[k].FundType),
			"FundCurr":       ilpfund[k].FundCurr + " - " + GetP0050ItemCodeDesc(iCompany, "FUNDCURR", 1, ilpfund[k].FundCurr),
			"BenefitName":    ibenfit.BCoverage + " - " + GetP0050ItemCodeDesc(iCompany, "COVR", 1, ibenfit.BCoverage),
			"FundPercentage": ilpfund[k].FundPercentage,
		}
		ilpfundtarray = append(ilpfundtarray, resultOut)

	}

	return ilpfundtarray
}

// # 169
// Get Previous Policy Data (New Version)
// Inputs: CompanyID, PolicyID, HistoryCode, Tranno
//
// # Outputs: JSON Policy Data
//
// ©  FuturaInsTech
func GetPPolicyData(iCompany uint, iPolicy uint, iHistoryCode string, iTranno uint, txn *gorm.DB) []interface{} {
	ppolicyarray := make([]interface{}, 0)
	var phistory models.PHistory
	result := txn.Find(&phistory, "company_id = ? and policy_id = ? and history_code = ?  and tranno =  ?", iCompany, iPolicy, iHistoryCode, iTranno)
	if result.Error != nil {
		return nil
	}
	previousPolicy := phistory.PrevData["Policy"]
	ppolicyarray = append(ppolicyarray, previousPolicy)
	return ppolicyarray

}

// # 170
// Get Previous Benefits Data (New Version)
// Inputs: CompanyID, PolicyID, HistoryCode, Tranno
//
// # Outputs: JSON Benefits Data
//
// ©  FuturaInsTech
func GetPBenefitData(iCompany uint, iPolicy uint, iHistoryCode string, iTranno uint, txn *gorm.DB) interface{} {
	var phistory models.PHistory
	result := txn.Find(&phistory, "company_id = ? and policy_id = ? and history_code = ?  and tranno =  ?", iCompany, iPolicy, iHistoryCode, iTranno)
	if result.Error != nil {
		return nil
	}
	previousBenefit := phistory.PrevData["Benefits"]
	return previousBenefit

}

// #174
// GetPayingAuthorityData - Printing Purpose Only
// Inputs: Company, Pa, Client, Address, Receipt and Date
//
// # Outputs  Paying Authority Details as an Interface
//
// ©  FuturaInsTech
func GetPayingAuthorityData(iCompany uint, iPa uint, txn *gorm.DB) []interface{} {
	payingautharray := make([]interface{}, 0)
	var payingauth models.PayingAuthority
	result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)
	if result.Error != nil {
		return nil
	}
	_, oStatus, _ := GetParamDesc(payingauth.CompanyID, "P0021", payingauth.PaStatus, 1)
	_, oCurr, _ := GetParamDesc(payingauth.CompanyID, "P0023", payingauth.PaCurrency, 1)

	resultOut := map[string]interface{}{
		"ID":                IDtoPrint(payingauth.ID),
		"CompanyID":         IDtoPrint(payingauth.CompanyID),
		"ClientID":          IDtoPrint(payingauth.ClientID),
		"AddressID":         IDtoPrint(payingauth.AddressID),
		"PaName":            payingauth.PaName,
		"PaType":            payingauth.PaType,
		"StartDate":         DateConvert(payingauth.StartDate),
		"EndDate":           DateConvert(payingauth.EndDate),
		"PaStatus":          oStatus,
		"ExtractionDay":     payingauth.ExtrationDay,
		"PayDay":            payingauth.PayDay,
		"PaToleranceAmount": NumbertoPrint(payingauth.PaToleranceAmt),
		"PaCurrency":        oCurr,
	}
	payingautharray = append(payingautharray, resultOut)

	return payingautharray
}

// #87
// GetClient Data - Printing Purpose Only
// Inputs: Company, Policy, Client, Address, Receipt and Date
//
// # Outputs  Client Details as an Interface
//
// ©  FuturaInsTech
func GetClientData(iCompany uint, iClient uint, txn *gorm.DB) []interface{} {
	clientarray := make([]interface{}, 0)
	var client models.Client

	txn.Find(&client, "company_id = ? and id = ?", iCompany, iClient)
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

func GetBankData(iCompany uint, iBank uint, txn *gorm.DB) []interface{} {
	bankarray := make([]interface{}, 0)
	var bank models.Bank
	txn.Find(&bank, "id = ?", iBank)
	_, oBanktype, _ := GetParamDesc(iCompany, "P0021", bank.BankType, 1)

	resultOut := map[string]interface{}{
		"ID":            IDtoPrint(bank.ID),
		"BankCode":      bank.BankCode,
		"BankAccountNo": bank.BankAccountNo,
		"StartDate":     DateConvert(bank.StartDate),
		"EndDate":       DateConvert(bank.EndDate),
		"BankType":      oBanktype,
	}
	bankarray = append(bankarray, resultOut)
	return bankarray
}

func PolAgntChData(iCompany uint, iPolicy uint, iAgent uint, iClient uint, txn *gorm.DB) map[string]interface{} {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Policy not found"}
	}
	var agntaddress models.Address
	result = txn.Find(&agntaddress, "company_id = ? AND client_id = ?", iCompany, iClient)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Address not found"}
	}

	var poladdress models.Address
	result = txn.Find(&poladdress, "company_id = ? AND id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Address not found"}
	}

	// Create a result map for each loan bill
	resultOut := map[string]interface{}{
		"addressline1": poladdress.AddressLine1,
		"postcode":     poladdress.AddressPostCode,
		"country":      poladdress.AddressCountry,
		"state":        poladdress.AddressState,
		"agntaddress":  agntaddress.AddressLine1,
		"agntpostcode": agntaddress.AddressPostCode,
		"agntstate":    agntaddress.AddressState,
		"agntcountry":  agntaddress.AddressCountry,
		"policyid":     polenq.ID,
		"agentid":      iAgent,
	}

	return resultOut
}

func PolicyDepData(iCompany uint, iPolicy uint, txn *gorm.DB) map[string]interface{} {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Policy not found"}
	}

	var clnt models.Client
	result = txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Client not found"}
	}

	var address models.Address
	result = txn.Find(&address, "company_id = ? AND id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Address not found"}
	}

	var pymt models.Payment
	result = txn.Find(&pymt, "company_id = ? AND policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return map[string]interface{}{"error": "Payment not found"}
	}

	var glbal models.GlBal
	result = txn.Find(&glbal, "company_id = ? AND gl_rdocno = ? ", iCompany, iPolicy)
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

// #210
// Get Requirement Data - Printing Purpose Only
// Inputs: Company, Policy
//
// # Outputs  Requirement Data for the Policyas an Interface
//
// ©  FuturaInsTech

func GetReqData(iCompany uint, iPolicy uint, iClient uint, txn *gorm.DB) []interface{} {
	reqArray := make([]interface{}, 0)
	resultMap, err := GetReqComm(iCompany, iPolicy, iClient, txn)
	if err != nil {
		fmt.Println("Error:", err)
		return nil // Return nil if there was an error
	}

	// Append resultMap (from GetReqComm) as the last element
	reqArray = append(reqArray, resultMap)

	// Return the array of all results
	return reqArray
}

// #186
// Get Loan Data
// Inputs: CompanyID, PolicyID and EffectiveDate
//
// # Outputs  All Loan Data with Total Loan Details
//
// ©  FuturaInsTech
func GetLoanData(iCompany uint, iPolicy uint, iEffectiveDate string, iOsLoanInterest float64, txn *gorm.DB) map[string]interface{} {
	combinedData := make(map[string]interface{})

	loanArray := make([]interface{}, 0)
	extraData := make([]map[string]interface{}, 0)
	overallData := make(map[string]interface{}) // Create a map for overall data

	var loanenq []models.Loan

	// Fetch loans for the specified company and policy
	txn.Find(&loanenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	var overallLoanAmount float64
	var overallstampduty float64
	var finalLoanAmountTotal float64

	var p0072data paramTypes.P0072Data
	var extradata4 paramTypes.Extradata = &p0072data
	GetItemD(int(iCompany), "P0072", "LN001", iEffectiveDate, &extradata4)

	// Map to keep track of already printed LoanSeqNumber values
	printedLoanSeqNumbers := make(map[uint]bool)

	for i := 0; i < len(loanenq); i++ {
		var benefit models.Benefit
		txn.First(&benefit, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, loanenq[i].BenefitID)

		// Calculate stamp duty based on the benefit
		stampDuty := CalculateStampDutyforLoan(iCompany, p0072data.StampDutyRate, iEffectiveDate, loanenq[i].LoanAmount, loanenq[i].PolicyID)
		overallstampduty += stampDuty

		// Calculate final loan amount by adding loan amount and stamp duty
		finalLoanAmount := loanenq[i].LoanAmount - stampDuty

		// Update finalLoanAmountTotal with the final loan amount
		finalLoanAmountTotal += finalLoanAmount

		// Construct resultOut map for the loan
		resultOut := map[string]interface{}{
			"ID":              IDtoPrint(loanenq[i].PolicyID),
			"BenefitID":       IDtoPrint(loanenq[i].BenefitID),
			"PProduct":        loanenq[i].PProduct,
			"BCoverage":       loanenq[i].BCoverage,
			"ClientID":        IDtoPrint(loanenq[i].ClientID),
			"LoanSeqNumber":   loanenq[i].LoanSeqNumber,
			"TranDate":        DateConvert(loanenq[i].TranDate),
			"TranNumber":      loanenq[i].TranNumber,
			"LoanEffDate":     DateConvert(loanenq[i].LoanEffDate),
			"LoanType":        loanenq[i].LoanType,
			"LoanStatus":      loanenq[i].LoanStatus,
			"LoanCurrency":    loanenq[i].LoanCurrency,
			"LoanAmount":      loanenq[i].LoanAmount,
			"LoanIntRate":     loanenq[i].LoanIntRate,
			"LoanIntType":     loanenq[i].LoanIntType,
			"LastCapAmount":   NumbertoPrint(loanenq[i].LastCapAmount),
			"LastCapDate":     DateConvert(loanenq[i].LastCapDate),
			"NextCapDate":     DateConvert(loanenq[i].NextCapDate),
			"LastIntBillDate": DateConvert(loanenq[i].LastIntBillDate),
			"NextIntBillDate": DateConvert(loanenq[i].NextIntBillDate),
			"StampDuty":       stampDuty,

			"FinalLoanAmount": finalLoanAmount,
		}

		// Append resultOut to loanArray
		loanArray = append(loanArray, resultOut)

		// Add LoanSeqNumber to printedLoanSeqNumbers if it's not already added
		if _, ok := printedLoanSeqNumbers[loanenq[i].LoanSeqNumber]; !ok {
			printedLoanSeqNumbers[loanenq[i].LoanSeqNumber] = true

			// Construct extraData map for the LoanSeqNumber
			extraData = append(extraData, map[string]interface{}{
				"LoanSeqNumber": loanenq[i].LoanSeqNumber,
				"LoanEffDate":   DateConvert(loanenq[i].LoanEffDate),
				"NextCapDate":   DateConvert(loanenq[i].NextCapDate),
			})
		}

		// Incrementally sum up the loan amounts to calculate overall loan amount
		overallLoanAmount += loanenq[i].LoanAmount
	}

	// Create a map to store overall loan data
	overallData["overallLoanAmount"] = overallLoanAmount
	overallData["overallstampduty"] = overallstampduty
	overallData["finalLoanAmountTotal"] = finalLoanAmountTotal
	overallData["OsLoaInterest"] = iOsLoanInterest

	// Assign arrays to their respective keys
	combinedData["a1"] = loanArray
	combinedData["a2"] = extraData
	combinedData["a3"] = []map[string]interface{}{overallData}

	return combinedData
}

// #188
// OS Loan and OS Loan Interest Calculation
// Inputs: CompanyID, PolicyID, Effective Date
//
// # Outputs  Outstanding Loan and Loan Interest
//
// ©  FuturaInsTech
func GetAllLoanInterestData(iCompany uint, iPolicy uint, iEffectiveDate string, txn *gorm.DB) []interface{} {
	var benefitenq []models.Benefit
	allLoanOs := make([]interface{}, 0)
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		return nil
	}

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data

	// Initialize variables
	var totalLoan float64
	var totalInt float64
	var totalOsAmount float64

	// Calculate total amount outside the loop
	for s := 0; s < len(benefitenq); s++ {
		iKey := benefitenq[s].BCoverage
		iDate := benefitenq[s].BStartDate
		iBenefit := benefitenq[s].ID

		GetItemD(int(iCompany), "Q0006", iKey, iDate, &extradataq0006)
		if q0006data.SurrMethod != "" && q0006data.LoanMethod != "" {
			oLoanOSP, oLoanIntOSP, _, _, _ := GetAllLoanOSByType(iCompany, iPolicy, iBenefit, iEffectiveDate)
			totalLoanP := oLoanOSP + oLoanIntOSP

			//note:= if we need we can revisite
			// oLoanOS_A, oLoanIntOS_A, , , _ := GetAllLoanOSByType(iCompany, iPolicy, iBenefit, iEffectiveDate, "A")
			// totalLoanA := oLoanOS_A + oLoanIntOS_A

			// Calculate the total amount
			totalLoan += totalLoanP //+ totalLoanA

			// Calculate the total interest
			totalInt += oLoanIntOSP // + oLoanIntOS_A
		}
		totalOsAmount += totalLoan + totalInt
	}

	allLoanOs = append(allLoanOs, map[string]interface{}{
		"TotLoanOS":     RoundFloat(totalLoan, 0),
		"TotIntOs":      RoundFloat(totalInt, 0),
		"TotalOsAmount": RoundFloat(totalOsAmount, 0),
	})

	return allLoanOs
}

// #193
// Get Loan Capitalized Amount
// Inputs: CompanyID, PolicyID, Effective Date, Minimum Loan Date, Maximum Loan Date
//
// # Outputs: Loan Capitalized Amount, OpenLoanBal Date, CloseLoanBal Date
//
// ©  FuturaInsTech
func LoanCapData(iCompany uint, iPolicy uint, iEffectiveDate string, minLoanBillDueDate string, maxLoanBillDueDate string, itotalcapamount float64, itotalInterest float64, itotalOsDue uint, txn *gorm.DB) []interface{} {
	allLoanCap := make([]interface{}, 0)

	var loanenq []models.Loan
	var prevloancapamount float64
	var oLoanInt float64

	txn.Find(&loanenq, "company_id = ? and policy_id = ? and  loan_type = ? and loan_status = ? and next_cap_date <=?", iCompany, iPolicy, "P", "AC", iEffectiveDate)

	var policyenq models.Policy

	txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	// var minLoanDate string
	// var maxLoanDate string
	var p0072data paramTypes.P0072Data
	var extradata paramTypes.Extradata = &p0072data
	GetItemD(int(iCompany), "P0072", "LN001", iEffectiveDate, &extradata)

	for i := 0; i < len(loanenq); i++ {

		prevloancapamount = loanenq[i].LastCapAmount
		oLoanInt = loanenq[i].LoanIntRate

	}

	resultOut := map[string]interface{}{
		"prevcapamount":  RoundFloat(prevloancapamount, 0),
		"minLoanDate":    DateConvert(minLoanBillDueDate),
		"maxLoanDate":    DateConvert(maxLoanBillDueDate),
		"InterestRate":   RoundFloat(oLoanInt, 0),
		"totalcapamount": itotalcapamount,
		"totalInterest":  itotalInterest,
		"totalOsDue":     itotalOsDue,
	}
	allLoanCap = append(allLoanCap, resultOut)
	return allLoanCap

}

// #192
// Get Loan Bill Data For Letters
// Inputs: CompanyID, PolicyID, Effective Date
//
// # Outputs: Loan Bill Data
//
// ©  FuturaInsTech
func LoanBillData(iCompany uint, iPolicy uint, iEffectiveDate string, txn *gorm.DB) []interface{} {

	var loanenq []models.Loan

	txn.Find(&loanenq, "company_id = ? and policy_id = ? and  loan_type = ? and loan_status = ? and next_int_bill_date<=?", iCompany, iPolicy, "P", "AC", iEffectiveDate)

	var policyenq models.Policy

	txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	var loanbillupd models.LoanBill
	var oLoanOS float64
	var oLoanIntOS float64

	var p0072data paramTypes.P0072Data
	var extradata paramTypes.Extradata = &p0072data
	GetItemD(int(iCompany), "P0072", "LN001", iEffectiveDate, &extradata)
	loanbill := make([]interface{}, 0)

	for i := 0; i < len(loanenq); i++ {
		var itemp float64

		// loanbillupd.PolicyID = newtranno
		loanbillupd.PolicyID = loanenq[i].PolicyID
		loanbillupd.BenefitID = loanenq[i].BenefitID
		loanbillupd.ClientID = loanenq[i].ClientID
		loanbillupd.LoanID = loanenq[i].ID
		loanbillupd.LoanBillCurr = loanenq[i].LoanCurrency
		loanbillupd.LoanType = loanenq[i].LoanType
		loanbillupd.LoanBillDueDate = loanenq[i].NextIntBillDate
		loanbillupd.LoanIntAmount = loanenq[i].LastCapAmount
		loanbillupd.CompanyID = iCompany

		oLoanOS = loanenq[i].LastCapAmount
		oLoanInt := loanenq[i].LoanIntRate
		_, _, _, iNoOfDays, _, _, _, _ := NoOfDays(loanenq[i].NextIntBillDate, loanenq[i].LastIntBillDate)

		if p0072data.LoanInterestType == "C" {
			itemp = CompoundInterest(oLoanOS, oLoanInt, float64(iNoOfDays))
		} else if p0072data.LoanInterestType == "S" {
			itemp = SimpleInterest(oLoanOS, oLoanInt, float64(iNoOfDays))
		}

		oLoanIntOS += itemp

		loanbillupd.LoanIntAmount = RoundFloat(itemp, 2)

		loanenq[i].LastIntBillDate = loanenq[i].NextIntBillDate
		a := GetNextDue(loanenq[i].NextIntBillDate, p0072data.IntPayableFreq, "")
		NextbillDate := Date2String(a)
		loanenq[i].NextIntBillDate = NextbillDate

		resultOut := map[string]interface{}{

			"PolicyID":  IDtoPrint(loanbillupd.PolicyID),
			"CompanyID": IDtoPrint(loanbillupd.CompanyID),
			"ClientID":  IDtoPrint(loanbillupd.ClientID),
			"BenefitID": IDtoPrint(loanbillupd.BenefitID),
			// "PaidToDate":       DateConvert(loanbillupd.PaidToDate),
			"LoanBillCurrency": loanbillupd.LoanBillCurr,
			"LoanType":         loanbillupd.LoanType,
			"LoanBillDueDate":  DateConvert(loanbillupd.LoanBillDueDate),
			"LoanIntAmount":    loanbillupd.LoanIntAmount,
		}
		loanbill = append(loanbill, resultOut)

	}

	return loanbill

}

// #196
// Get Loan Bills Interest For Letters
// Inputs: CompanyID, PolicyID, SeqNo,CurrentInterestDueAmount
//
// # Outputs: LoanBills,TotalInterest,totUnpaidInterest
//
// ©  FuturaInsTech

func LoanBillsInterestData(iCompany uint, iPolicy uint, iSeqNo uint, iCurrentIntDue float64, txn *gorm.DB) map[string]interface{} {
	var loanbillupd1 []models.LoanBill

	txn.Order("CASE WHEN loan_seq_number = 1 THEN 0 WHEN loan_seq_number = 2 THEN 1 ELSE 2 END").
		Find(&loanbillupd1, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	var totalInterest float64
	var totUnpaidIntetest float64

	var loanbills []interface{}
	overallInt := make(map[string]interface{})

	for i := 0; i < len(loanbillupd1); i++ {
		// Add LoanIntAmount to totalInterest
		totalInterest += loanbillupd1[i].LoanIntAmount

		// Create a result map for each loan bill
		resultOut := map[string]interface{}{
			"PolicyID":         IDtoPrint(loanbillupd1[i].PolicyID),
			"CompanyID":        IDtoPrint(loanbillupd1[i].CompanyID),
			"ClientID":         IDtoPrint(loanbillupd1[i].ClientID),
			"BenefitID":        IDtoPrint(loanbillupd1[i].BenefitID),
			"LoanBillCurrency": loanbillupd1[i].LoanBillCurr,
			"LoanType":         loanbillupd1[i].LoanType,
			"LoanBillDueDate":  DateConvert(loanbillupd1[i].LoanBillDueDate),
			"LoanIntAmount":    loanbillupd1[i].LoanIntAmount,
			"LoanSeqNo":        loanbillupd1[i].LoanSeqNumber,
		}
		// Append result map to loanbills slice
		loanbills = append(loanbills, resultOut)
	}
	totUnpaidIntetest += totalInterest + iCurrentIntDue

	overallInt["TotalInterest"] = totalInterest
	overallInt["totUnpaidIntetest"] = RoundFloat(totUnpaidIntetest, 2)

	// Combine loanbills and totalInterestMap into a single map
	combinedData := map[string]interface{}{
		"LoanBills": loanbills,
		"TotOsInt":  []map[string]interface{}{overallInt},
	}

	// Return combinedData map
	return combinedData
}

// #197
func GetPaymentData(iCompany uint, iPolicyID uint, iPayment uint, txn *gorm.DB) map[string]interface{} {
	var polenq models.Policy
	txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	var clt models.Client
	txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	var add models.Address
	txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	var paymentenq models.Payment
	txn.Find(&paymentenq, "company_id = ? and id = ?", iCompany, iPayment)
	amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)
	resultOut := map[string]interface{}{
		"ID":               IDtoPrint(paymentenq.ID),
		"ClientSalutation": clt.Salutation,
		"ClientFullName":   clt.ClientLongName,
		"AddressLine1":     add.AddressLine1,
		"AddressLine2":     add.AddressLine2,
		"AddressLine3":     add.AddressLine3,
		"AddressLine4":     add.AddressLine4,
		"AddressState":     add.AddressState,
		"AddressCountry":   add.AddressCountry,
		"AccCurry":         paymentenq.AccCurry,
		"AccAmount":        paymentenq.AccAmount,
		"Reason":           paymentenq.Reason,
		"PolicyId":         IDtoPrint(paymentenq.PolicyID),
		"DateOfPayment":    DateConvert(paymentenq.DateOfPayment),
		"BankReferenceNo":  paymentenq.BankReferenceNo,
		"TypeOfPayment":    paymentenq.TypeOfPayment,
		"AmountInWords":    amtinwords,
		"CurrSymbol":       csymbol,
	}
	return resultOut

}

// #198
// Create Communication multi (New Version with Rollback)
//
// # This function, Create Communication Records by getting input values as Company ID, History Code, Tranno, Date of Transaction, Policy Id, Client Id, Address Id, Receipt ID . Quotation ID, Agency ID
// 10 Input Variables
// # It returns success or failure.  Successful records written in Communciaiton Table
//
// ©  FuturaInsTech
func CreateCommunicationsM(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB, iBenefit uint, iPa uint, iClientWork uint) error {

	var communication models.Communication
	var iKey string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	//utilities.LetterCreate(int(iCompany), uint(iPolicy), iHistoryCode, createreceipt.CurrentDate, idata)
	iTransaction := iHistoryCode
	iReceiptTranCode := "H0034"
	iReceiptFor := ""

	if iReceipt != 0 {
		var receipt models.Receipt
		result := txn.Find(&receipt, "company_id = ? and id = ?", iCompany, iReceipt)
		if result.Error != nil {
			return result.Error
		}
		iReceiptFor = receipt.ReceiptFor
	}

	if iPolicy != 0 {
		var policy models.Policy
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.Error != nil {
			return result.Error
		}
		communication.CompanyID = uint(iCompany)
		communication.AgencyID = policy.AgencyID
		communication.ClientID = policy.ClientID
		communication.PolicyID = policy.ID
		communication.Tranno = policy.Tranno
		communication.EffectiveDate = policy.PRCD
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPolicy
		if iTransaction == iReceiptTranCode {
			iKey = iTransaction + iReceiptFor
		} else {
			iKey = iTransaction + policy.PProduct
		}
	}

	if iPolicy == 0 && iTransaction == iReceiptTranCode && iPa != 0 {
		var payingauth models.PayingAuthority
		result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)
		if result.Error != nil {
			return result.Error
		}

		communication.CompanyID = uint(iCompany)
		communication.AgencyID = 0
		communication.ClientID = payingauth.ClientID
		communication.PolicyID = 0
		communication.Tranno = 0
		communication.EffectiveDate = iDate
		communication.ReceiptFor = iReceiptFor
		communication.ReceiptRefNo = iPa
		iKey = iTransaction + iReceiptFor
	}

	err1 := GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
		if err1 != nil {
			return err1
		}
	}

	seqno := 0
	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			err := GetItemD(int(iCompany), "P0033", iKey, iDate, &extradatap0033)
			if err != nil {
				return err
			}

			iPageSize := p0034data.Letters[i].PageSize
			iOrientation := p0034data.Letters[i].Orientation

			communication.AgentEmailAllowed = p0033data.AgentEmailAllowed
			communication.AgentSMSAllowed = p0033data.AgentSMSAllowed
			communication.AgentWhatsAppAllowed = p0033data.AgentWhatsAppAllowed
			communication.EmailAllowed = p0033data.EmailAllowed
			communication.SMSAllowed = p0033data.SMSAllowed
			communication.WhatsAppAllowed = p0033data.WhatsAppAllowed
			communication.DepartmentHead = p0033data.DepartmentHead
			communication.DepartmentName = p0033data.DepartmentName
			communication.CompanyPhone = p0033data.CompanyPhone
			communication.CompanyEmail = p0033data.CompanyEmail

			communication.TemplateName = iKey
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
					oData := GetCompanyData(iCompany, iDate, txn)
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData := GetClientData(iCompany, iClient, txn)
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData := GetAddressData(iCompany, iAddress, txn)
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData := GetPolicyData(iCompany, iPolicy, txn)
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData := GetBenefitData(iCompany, iPolicy, txn)
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData := GetSurBData(iCompany, iPolicy, txn)
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData := GetMrtaData(iCompany, iPolicy, txn)
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData := GetReceiptData(iCompany, iReceipt, txn)
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData := GetSaChangeData(iCompany, iPolicy, txn)
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData := GetCompAddData(iCompany, iPolicy, txn)
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData := GetSurrHData(iCompany, iPolicy, txn)
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData := GetDeathData(iCompany, iPolicy, txn)
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData := GetMatHData(iCompany, iPolicy, txn)
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData := GetSurvBPay(iCompany, iPolicy, iTranno, txn)
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData := GetExpi(iCompany, iPolicy, iTranno, txn)
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData := GetBonusVals(iCompany, iPolicy, txn)
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData := GetAgency(iCompany, iAgency, txn)
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData := GetNomiData(iCompany, iPolicy, txn)
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData := GetGLData(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData := GetIlpSummaryData(iCompany, iPolicy, txn)
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData := GetIlpAnnsummaryData(iCompany, iPolicy, iHistoryCode, txn)
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData := GetIlpTranctionData(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData := GetPremTaxGLData(iCompany, iPolicy, iFromDate, iToDate, txn)
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData := GetIlpFundSwitchData(iCompany, iPolicy, iTranno, txn)
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData := GetPHistoryData(iCompany, iPolicy, iHistoryCode, iDate, txn)
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData := GetIlpFundData(iCompany, iPolicy, iBenefit, iDate, txn)
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData := GetPPolicyData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData := GetPBenefitData(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					fmt.Println(oData) // Dummy to avoid compilation error
				case oLetType == "29":
					oData := GetPayingAuthorityData(iCompany, iPa, txn)
					resultMap["PrevBenefit"] = oData
				case oLetType == "30":
					oData := GetClientWorkData(iCompany, iClientWork, txn)
					resultMap["ClientWork"] = oData
				case oLetType == "36":
					oData := GetReqData(iCompany, iPolicy, iClient, txn)
					for _, item := range oData {
						for key, value := range item.(map[string]interface{}) {
							resultMap[key] = value
						}
					}
				case oLetType == "37":
					oData := PolicyDepData(iCompany, iPolicy, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "38":
					oData := PolAgntChData(iCompany, iPolicy, iAgency, iClient, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "39":
					oData := GetBankData(iCompany, iClientWork, txn)
					resultMap["BankData"] = oData
				case oLetType == "40":
					iKey := iReceipt
					oData := GetPaymentData(iCompany, iPolicy, iKey, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "41":
					oData := GetHIPPOLSCDData(iCompany, iPolicy, iPageSize, iOrientation, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "42":
					oData := GetPriorPolicyData(iCompany, iPolicy, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "43":
					oData := GetTermAndConditionData(iCompany, iPolicy, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "44":
					oData := GetpremiumCertificateData(iCompany, iPolicy, txn)
					for key, value := range oData {
						resultMap[key] = value
					}

				case oLetType == "98":
					resultMap["BatchData"] = batchData

				case oLetType == "99":
					resultMap["SignData"] = signData
				default:

				}

				communication.ExtractedData = resultMap
				communication.PDFPath = p0034data.Letters[i].PdfLocation
				communication.TemplatePath = p0034data.Letters[i].ReportTemplateLocation
				// New Changes for Online Print and Email Trigger
				if p0033data.Online == "Y" {
					err := GetReportforOnline(communication, p0033data.TemplateName, txn)
					if err != nil {
						log.Fatalf("Failed to generate report: %v", err)
					}
				}
				if p0033data.SMSAllowed == "Y" {
					err := SendSMSTwilio(communication.CompanyID, communication.ClientID, p0033data.TemplateName, communication.EffectiveDate, p0033data.SMSBody, txn)
					if err != nil {
						log.Fatalf("Failed to send SMS: %v", err)
					}
				}
				communication.Print = "Y"
				communication.PrintDate = iDate
				communication.UpdatedID = 1
				communication.ID = 0
				communication.Seqno = uint16(seqno)
				// New Changes Ended
				results := txn.Create(&communication)

				if results.Error != nil {
					return results.Error
				}

				seqno++
			}
		}
	}
	return nil
}

// #199
func GetHIPPOLSCDData(iCompany uint, iPolicyID uint, iPageSize, iOrientation string, txn *gorm.DB) map[string]interface{} {
	var polenq models.Policy
	txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	var priorpolenq models.PriorPolicy
	txn.Where("company_id = ? AND policy_id = ?", iCompany, iPolicyID).Order("created_at  DESC").First(&priorpolenq)
	var receiptenq models.Receipt
	txn.Where("company_id = ? AND receipt_ref_no = ?", iCompany, iPolicyID).Order("created_at  DESC").First(&receiptenq)
	var benefitenq []models.Benefit
	txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)

	var cmp models.Company
	txn.Find(&cmp, "id = ?", polenq.CompanyID)
	var clt models.Client
	txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	var agency models.Agency
	txn.Find(&agency, "company_id = ? and id = ?", iCompany, polenq.AgencyID)
	var agent models.Client
	txn.Find(&agent, "company_id = ? and id = ?", iCompany, agency.ClientID)
	var add models.Address
	txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)

	//amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)

	// compadd := cmp.CompanyAddress1+", "+cmp.CompanyAddress2+", "+cmp.CompanyAddress3+", "+cmp.CompanyAddress4+", "+cmp.CompanyAddress5+", "+cmp.CompanyCountry+", "+cmp.CompanyPostalCode

	compadd := GetCompanyFullAddress(cmp)

	yr, _ := extractYear(polenq.PRCD)
	prcd, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.PRCD)
	riskcessdate := benefitenq[0].BRiskCessDate
	enddate, _ := ConvertYYYYMMDDtoDDMMYYYY(riskcessdate)
	clientfulladdress, _ := GetFullAddress(iCompany, polenq.AddressID)

	benpalntypedesc := GetP0050ItemCodeDesc(iCompany, "HealthBenefitType", 1, benefitenq[0].BenefitType)

	slno, laname, ladob, lagender, laoccup, larel, pecdecl, ppolstdate, bprem, lanominee, lanomrel, err := GetPlanLifeData(iCompany, iPolicyID)

	if err != nil {
		log.Println("Error fetching plan life data:", err)
	}

	totalpreflifedisamt, totalloyaltydisamt, totalfloaterdisamt, totalonlinedisamt := GetTotalDiscountsByPolicy(iCompany, iPolicyID)
	premwogstamt, gstamt, stampdutyamt, totalpremiumamt := GetPremiumDetailsFromGL(1, 1001)

	receiptdate, _ := ConvertYYYYMMDDtoDDMMYYYY(receiptenq.DateOfCollection)

	resultOut := map[string]interface{}{
		"Layout": map[string]string{
			"PageSize":    iPageSize,
			"Orientation": iOrientation,
		},
		"compname":          cmp.CompanyName,
		"compadd":           compadd,
		"uinno":             cmp.CompanyUid,
		"clientno":          clt.ID,
		"yr":                yr,
		"policyno":          polenq.ID,
		"startdate":         prcd,
		"enddate":           enddate,
		"ClientSalutation":  clt.Salutation,
		"ClientFullName":    clt.ClientLongName,
		"ClientFullAddress": clientfulladdress,
		"agentname":         agent.ClientLongName,
		"agentcode":         agent.ID,
		"agentphoneno":      agent.ClientMobile,
		"agentemailaddress": agent.ClientEmail,
		"datime":            time.Now(),
		// "prevyr":,
		"prevpol":             priorpolenq.PolicyID,
		"clientemailaddress":  clt.ClientEmail,
		"clientphoneno":       clt.ClientMobile,
		"benefitplantypedesc": benpalntypedesc,
		"bsumassured":         benefitenq[0].BSumAssured,
		"slno":                slno,
		"laname":              laname,
		"ladob":               ladob,
		"lagender":            lagender,
		"larel":               larel,
		"laoccup":             laoccup,
		"pecdecl":             pecdecl,
		"ppolstdate":          ppolstdate,
		"lanominee":           lanominee,
		"lanomrel":            lanomrel,
		"bprem":               bprem,
		"totalbasepremium":    benefitenq[0].BBasAnnualPrem,
		"totalloadingprem":    benefitenq[0].BLoadPrem,
		"totalpreflifedisamt": totalpreflifedisamt,
		"totalloyaltydisamt":  totalloyaltydisamt,
		"totalfloaterdisamt":  totalfloaterdisamt,
		"totalonlinedisamt":   totalonlinedisamt,
		"premwogstamt":        premwogstamt,
		"gstamt":              gstamt,
		"stampdutyamt":        stampdutyamt,
		"totalpremiumamt":     totalpremiumamt,
		"receiptno":           receiptenq.ID,
		"receiptdate":         receiptdate,
		"clientgstno":         "",
		"saccode":             "",
		"gstofficeno":         "",
		"gstinvoiceno":        "",
		"gstdate":             "",
		"day":                 "",
		"month":               "",
		"year":                "",
	}
	return resultOut

}

// #199
func GetPriorPolicyData(iCompany uint, iPolicyID uint, txn *gorm.DB) map[string]interface{} {
	var polenq models.Policy
	txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	var cmp models.Company
	txn.Find(&cmp, "id = ?", polenq.CompanyID)
	var clt models.Client
	txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	var add models.Address
	txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)

	var priorpolenq []models.PriorPolicy
	txn.Find(&priorpolenq, "company_id = ? and id = ?", iCompany, iPolicyID)

	//amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)

	// compadd := cmp.CompanyAddress1+", "+cmp.CompanyAddress2+", "+cmp.CompanyAddress3+", "+cmp.CompanyAddress4+", "+cmp.CompanyAddress5+", "+cmp.CompanyCountry+", "+cmp.CompanyPostalCode

	parts := []string{
		cmp.CompanyAddress1,
		cmp.CompanyAddress2,
		cmp.CompanyAddress3,
		cmp.CompanyAddress4,
		cmp.CompanyAddress5,
		cmp.CompanyCountry,
		cmp.CompanyPostalCode,
	}

	nonEmptyParts := []string{}
	for _, part := range parts {
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	compadd := strings.Join(nonEmptyParts, ", ")
	yr := time.Now().Year()

	slno := []int{}
	psumAssured := []float64{}

	priorPolicyNo := []string{}
	priorInsurerName := []string{}
	pstartDates := []string{}
	pendDates := []string{}

	for _, prior := range priorpolenq {
		slno = append(slno, int(prior.SeqNo))                   // int
		priorPolicyNo = append(priorPolicyNo, prior.PriorPolNo) // string
		priorInsurerName = append(priorInsurerName, prior.PriorInsurerName)

		pstart, _ := ConvertYYYYMMDDtoDDMMYYYY(prior.PStartDate)
		pend, _ := ConvertYYYYMMDDtoDDMMYYYY(prior.PEndDate)
		pstartDates = append(pstartDates, pstart)
		pendDates = append(pendDates, pend)
		psumAssured = append(psumAssured, float64(prior.PSumAssured)) // float64
	}

	resultOut := map[string]interface{}{
		"compname":        cmp.CompanyName,
		"compadd":         compadd,
		"uinno":           cmp.CompanyUid,
		"clientno":        clt.ID,
		"yr":              yr,
		"policyno":        polenq.ID,
		"slno":            slno,
		"previnsurername": priorInsurerName,
		"prevpolicyno":    priorPolicyNo,
		"prevstartdate":   pstartDates,
		"prevenddate":     pendDates,
		"prevsumassured":  psumAssured,
	}
	return resultOut

}

// #200
// term and condition
func GetTermAndConditionData(iCompany uint, iPolicyID uint, txn *gorm.DB) map[string]interface{} {
	var polenq models.Policy
	txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	var cmp models.Company
	txn.Find(&cmp, "id = ?", polenq.CompanyID)
	var clt models.Client
	txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	var add models.Address
	txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)

	var priorpolenq []models.PriorPolicy
	txn.Find(&priorpolenq, "company_id = ? and id = ?", iCompany, iPolicyID)

	//amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)

	// compadd := cmp.CompanyAddress1+", "+cmp.CompanyAddress2+", "+cmp.CompanyAddress3+", "+cmp.CompanyAddress4+", "+cmp.CompanyAddress5+", "+cmp.CompanyCountry+", "+cmp.CompanyPostalCode

	compadd := GetCompanyFullAddress(cmp)
	//yr := time.Now().Year()

	resultOut := map[string]interface{}{
		"compname":         cmp.CompanyName,
		"compadd":          compadd,
		"uinno":            cmp.CompanyUid,
		"policyno":         polenq.ID,
		"ClientSalutation": clt.Salutation,
		"ClientFullName":   clt.ClientLongName,
	}
	return resultOut

}

// #201
// premuim certificate
func GetpremiumCertificateData(iCompany uint, iPolicyID uint, txn *gorm.DB) map[string]interface{} {
	var polenq models.Policy
	txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	var cmp models.Company
	txn.Find(&cmp, "id = ?", polenq.CompanyID)
	var clt models.Client
	txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	var add models.Address
	txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)

	var priorpolenq []models.PriorPolicy
	txn.Find(&priorpolenq, "company_id = ? and id = ?", iCompany, iPolicyID)

	//amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)

	// compadd := cmp.CompanyAddress1+", "+cmp.CompanyAddress2+", "+cmp.CompanyAddress3+", "+cmp.CompanyAddress4+", "+cmp.CompanyAddress5+", "+cmp.CompanyCountry+", "+cmp.CompanyPostalCode

	compadd := GetCompanyFullAddress(cmp)
	//yr := time.Now().Year()

	resultOut := map[string]interface{}{
		"compname":         cmp.CompanyName,
		"compadd":          compadd,
		"uinno":            cmp.CompanyUid,
		"policyno":         polenq.ID,
		"ClientSalutation": clt.Salutation,
		"ClientFullName":   clt.ClientLongName,
	}
	return resultOut

}

// ConvertYYYYMMDDtoDDMMYYYY converts a date from "yyyymmdd" to "dd/mm/yyyy"
func ConvertYYYYMMDDtoDDMMYYYY(dateStr string) (string, error) {
	parsedDate, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("02/01/2006"), nil
}

func extractYear(date string) (string, error) {
	if len(date) != 8 {
		return "", fmt.Errorf("invalid date format, expected 8 characters")
	}
	year := date[6:] // Get the last 4 characters (year)
	if len(year) != 4 {
		return "", fmt.Errorf("year format incorrect")
	}
	return year[2:], nil // Return last two digits
}

func GetCompanyFullAddress(cmp models.Company) string {
	parts := []string{
		cmp.CompanyAddress1,
		cmp.CompanyAddress2,
		cmp.CompanyAddress3,
		cmp.CompanyAddress4,
		cmp.CompanyAddress5,
		cmp.CompanyCountry,
		cmp.CompanyPostalCode,
	}

	nonEmptyParts := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	return strings.Join(nonEmptyParts, ", ")
}

func GetFullAddress(iCompany, iAddressId uint) (string, error) {
	var addr models.Address
	err := initializers.DB.First(&addr, "id = ? AND company_id = ?", iAddressId, iCompany).Error
	if err != nil {
		return "", err
	}

	parts := []string{
		addr.AddressLine1,
		addr.AddressLine2,
		addr.AddressLine3,
		addr.AddressLine4,
		addr.AddressLine5,
		addr.AddressState,
		addr.AddressCountry,
		addr.AddressPostCode,
	}

	nonEmptyParts := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	return strings.Join(nonEmptyParts, ", "), nil
}

func GetPlanLifeData(iCompany, iPolicy uint) ([]int, []string, []string, []string, []string, []string, []string, []string, []float64, []string, []string, error) {

	var planLives []models.PlanLife
	initializers.DB.Find(&planLives, "company_id = ? and policy_id =?", iCompany, iPolicy)

	slno := make([]int, len(planLives))
	laname := make([]string, len(planLives))
	ladob := make([]string, len(planLives))
	lagender := make([]string, len(planLives))
	laoccup := make([]string, len(planLives))
	larel := make([]string, len(planLives))
	pecdecl := make([]string, len(planLives))
	ppolstdate := make([]string, len(planLives))
	bprem := make([]float64, len(planLives))

	lanominee := make([]string, len(planLives))
	lanomrel := make([]string, len(planLives))

	for i, plan := range planLives {
		slno[i] = i + 1

		var client models.Client
		err := initializers.DB.First(&client, "company_id = ? AND id = ?", plan.CompanyID, plan.ClientID).Error
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to fetch client for plan life %d: %w", i+1, err)
		}

		laname[i] = strings.TrimSpace(client.ClientLongName)

		formattedDOB, err := ConvertYYYYMMDDtoDDMMYYYY(client.ClientDob)
		if err != nil {
			ladob[i] = ""
		} else {
			ladob[i] = formattedDOB
		}

		lagender[i] = strings.TrimSpace(client.Gender)
		laoccup[i] = strings.TrimSpace(client.Occupation)

		relationDesc := GetP0050ItemCodeDesc(iCompany, "PlanLARelations", 1, plan.ClientRelCode)
		larel[i] = strings.TrimSpace(relationDesc)

		pecdecl[i] = "None"

		pstart, err := ConvertYYYYMMDDtoDDMMYYYY(plan.PStartDate)
		if err != nil {
			ppolstdate[i] = ""
		} else {
			ppolstdate[i] = pstart
		}

		bprem[i] = plan.PBasAnnualPrem

		lanominee[i] = ""
		lanomrel[i] = ""
	}

	return slno, laname, ladob, lagender, laoccup, larel, pecdecl, ppolstdate, bprem, lanominee, lanomrel, nil
}

func GetTotalDiscountsByPolicy(companyID, iPolicyID uint) (float64, float64, float64, float64) {
	var discounts []models.PlanLifeDiscount
	initializers.DB.Find(&discounts, "company_id = ? and policy_id = ?", companyID, iPolicyID)

	var totalPrefLifeDisAmt float64
	var totalLoyaltyDisAmt float64
	var totalFloaterDisAmt float64
	var totalOnlineDisAmt float64 = 0 // always zero

	for _, d := range discounts {
		if d.PDisType01 != "" {
			switch d.PDisType01 {
			case "P":
				totalPrefLifeDisAmt += d.PDisPrem01
			case "L":
				totalLoyaltyDisAmt += d.PDisPrem01
			case "F":
				totalFloaterDisAmt += d.PDisPrem01
			}
		}
		if d.PDisType02 != "" {
			switch d.PDisType02 {
			case "P":
				totalPrefLifeDisAmt += d.PDisPrem02
			case "L":
				totalLoyaltyDisAmt += d.PDisPrem02
			case "F":
				totalFloaterDisAmt += d.PDisPrem02
			}
		}
		if d.PDisType03 != "" {
			switch d.PDisType03 {
			case "P":
				totalPrefLifeDisAmt += d.PDisPrem03
			case "L":
				totalLoyaltyDisAmt += d.PDisPrem03
			case "F":
				totalFloaterDisAmt += d.PDisPrem03
			}
		}
		if d.PDisType04 != "" {
			switch d.PDisType04 {
			case "P":
				totalPrefLifeDisAmt += d.PDisPrem04
			case "L":
				totalLoyaltyDisAmt += d.PDisPrem04
			case "F":
				totalFloaterDisAmt += d.PDisPrem04
			}
		}
		if d.PDisType05 != "" {
			switch d.PDisType05 {
			case "P":
				totalPrefLifeDisAmt += d.PDisPrem05
			case "L":
				totalLoyaltyDisAmt += d.PDisPrem05
			case "F":
				totalFloaterDisAmt += d.PDisPrem05
			}
		}
	}

	return totalPrefLifeDisAmt, totalLoyaltyDisAmt, totalFloaterDisAmt, totalOnlineDisAmt
}

func GetPremiumDetailsFromGL(iCompany, iPolicy uint) (float64, float64, float64, float64) {
	var glMoves []models.GlMove

	// Fetch GL moves for given company and policy
	initializers.DB.Find(&glMoves, "company_id = ? AND gl_rdocno = ?", iCompany, iPolicy)

	var premwogstamt float64
	var gstamt float64
	var stampdutyamt float64

	for _, move := range glMoves {
		switch {
		case strings.HasPrefix(move.AccountCode, "PremiumAccount"):
			premwogstamt += move.ContractAmount
		case strings.HasPrefix(move.AccountCode, "GSTPayable"):
			gstamt += move.ContractAmount
		case strings.HasPrefix(move.AccountCode, "StampDuty"):
			stampdutyamt += move.ContractAmount
		}
	}

	totalpremiumamt := premwogstamt + gstamt + stampdutyamt

	return premwogstamt, gstamt, stampdutyamt, totalpremiumamt
}
