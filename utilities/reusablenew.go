package utilities

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"
	"github.com/FuturaInsTech/GoLifeLib/types"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// 2025-10-15 Divya Changes

func TDFCollDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iDate string, txn *gorm.DB) (string, models.TxnError) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		return "", models.TxnError{}
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {

			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		return "", models.TxnError{}
	}
}

func TdfhUpdateNNew(iCompany uint, iPolicy uint, txn *gorm.DB) models.TxnError {
	var tdfhupd models.Tdfh
	var tdfpolicyenq []models.TDFPolicy

	iDate := "29991231"

	results := txn.Find(&tdfpolicyenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if results.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL392", DbError: results.Error}
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
				return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
		} else {
			result = txn.Delete(&tdfhupd)
			if result.Error != nil {
				return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			var tdfhupd models.Tdfh
			tdfhupd.CompanyID = iCompany
			tdfhupd.PolicyID = iPolicy
			tdfhupd.EffectiveDate = iDate
			tdfhupd.ID = 0
			result = txn.Create(&tdfhupd)
			if result.Error != nil {
				return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
		}

	}
	return models.TxnError{}
}

func CreateCommunicationsMNew(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB, iBenefit uint, iPa uint, iClientWork uint) models.TxnError {

	var communication models.Communication
	var iKey string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	iTransaction := iHistoryCode
	iReceiptTranCode := "H0034"
	iReceiptFor := ""

	if iReceipt != 0 {
		var receipt models.Receipt
		result := txn.Find(&receipt, "company_id = ? and id = ?", iCompany, iReceipt)
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL014", DbError: result.Error}
		}
		iReceiptFor = receipt.ReceiptFor

		receiptMaxTRanNo, err := GetReceiptMaxTranNo(iCompany, iPolicy, iReceiptFor)
		if err != nil {
			return models.TxnError{ErrorCode: "GL722"}
		}
		communication.Tranno = receiptMaxTRanNo
	}

	if iPolicy != 0 {
		var policy models.Policy
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL003", DbError: result.Error}
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
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL671", DbError: result.Error}
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
	errparam := "P0034"
	err1 := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0034)
	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0034)
		if err1 != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
		}
	}

	seqno := 0
	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			errparam = "P0033"
			err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0033)
			if err != nil {
				return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}

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
					oData, funcErr := GetCompanyDataNew(iCompany, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData, funcErr := GetClientDataNew(iCompany, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData, funcErr := GetAddressDataNew(iCompany, iAddress, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData, funcErr := GetPolicyDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData, funcErr := GetBenefitDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData, funcErr := GetSurBDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData, funcErr := GetMrtaDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData, funcErr := GetReceiptDataNew(iCompany, iReceipt, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData, funcErr := GetSaChangeDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData, funcErr := GetCompAddDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData, funcErr := GetSurrHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData, funcErr := GetDeathDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData, funcErr := GetMatHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData, funcErr := GetSurvBPayNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData, funcErr := GetExpiNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData, funcErr := GetBonusValsNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData, funcErr := GetAgencyNew(iCompany, iAgency, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData, funcErr := GetNomiDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData, funcErr := GetGLDataNew(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData, funcErr := GetIlpSummaryDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData, funcErr := GetIlpAnnsummaryDataNew(iCompany, iPolicy, iHistoryCode, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData, funcErr := GetIlpTranctionDataNew(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData, funcErr := GetPremTaxGLDataNew(iCompany, iPolicy, iFromDate, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData, funcErr := GetIlpFundSwitchDataNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData, funcErr := GetPHistoryDataNew(iCompany, iPolicy, iHistoryCode, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData, funcErr := GetIlpFundDataNew(iCompany, iPolicy, iBenefit, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData, funcErr := GetPPolicyDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData, funcErr := GetPBenefitDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					fmt.Println(oData) // Dummy to avoid compilation error
				case oLetType == "29":
					oData, funcErr := GetPayingAuthorityDataNew(iCompany, iPa, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevBenefit"] = oData
				case oLetType == "30":
					oData, funcErr := GetClientWorkDataNew(iCompany, iClientWork, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ClientWork"] = oData
				case oLetType == "36":
					oData, funcErr := GetReqDataNew(iCompany, iPolicy, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for _, item := range oData {
						for key, value := range item.(map[string]interface{}) {
							resultMap[key] = value
						}
					}
				case oLetType == "37":
					oData, funcErr := PolicyDepDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "38":
					oData, funcErr := PolAgntChDataNew(iCompany, iPolicy, iAgency, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "39":
					oData, funcErr := GetBankDataNew(iCompany, iClientWork, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BankData"] = oData
				case oLetType == "40":
					iKey := iReceipt
					oData, funcErr := GetPaymentDataNew(iCompany, iPolicy, iKey, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "41":
					oData, funcErr := GetHIPPOLSCDDataNew(iCompany, iPolicy, iPageSize, iOrientation, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "42":
					oData, funcErr := GetPriorPolicyDataNew(iCompany, iPolicy, iPageSize, iOrientation, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "43":
					oData, funcErr := GetTermAndConditionDataNew(iCompany, iPolicy, iPageSize, iOrientation, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "44":
					oData, funcErr := GetpremiumCertificateDataNew(iCompany, iPolicy, iPageSize, iOrientation, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "47":
					oData, funcErr := GetPOLSCDEndowmentDataNew(iCompany, iPolicy, iPageSize, iOrientation, p0033data, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "51":
					oData, funcErr := PrtReceiptDataNew(iCompany, iReceipt, iPolicy, iPa, p0033data, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "52":
					oData, funcErr := PrtPolicyBillDataNew(iCompany, iPolicy, iDate, p0033data, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "53":
					oData, funcErr := PrtPolicyLapseDataNew(iCompany, iPolicy, iDate, p0033data, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "54":
					oData, funcErr := PrtCollectionDataNew(iCompany, iPolicy, iDate, p0033data, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "55":
					oData, funcErr := PrtAnniDataNew(iCompany, iPolicy, iDate, p0033data, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "56":
					oData, funcErr := PrtAnniILPDataNew(iCompany, iPolicy, iDate, p0033data, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "58":
					oData, funcErr := PrtExpiDataNew(iCompany, iPolicy, iDate, p0033data, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "61":
					oData, funcErr := PrtPremstDataNew(iCompany, iPolicy, iBenefit, iDate, p0033data, iTranno, iAgency, iFromDate, iToDate, iHistoryCode, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "65":
					oData, funcErr := PrtFreqChangeDataNew(iCompany, iPolicy, iDate, p0033data, iAgency, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "66":
					oData, funcErr := PrtSachangeDataNew(iCompany, iPolicy, iDate, p0033data, iAgency, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "67":
					oData, funcErr := PrtCompaddDataNew(iCompany, iPolicy, iDate, p0033data, iAgency, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "68":
					oData, funcErr := PrtSurrDataNew(iCompany, iPolicy, iDate, p0033data, iAgency, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "69":
					oData, funcErr := PrtMatyDataNew(iCompany, iPolicy, iDate, p0033data, iAgency, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
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
					//err := GetReportforOnline(communication, p0033data.TemplateName, txn)
					funcErr := GetReportforOnlineV3New(communication, p0033data.TemplateName, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
				}

				if p0033data.SMSAllowed == "Y" {
					funcErr := SendSMSTwilioNew(communication.CompanyID, communication.ClientID, p0033data.TemplateName, communication.EffectiveDate, p0033data.SMSBody, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
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
					return models.TxnError{ErrorCode: "DBERR", DbError: results.Error}

				}

				seqno++
			}
		}
	}
	return models.TxnError{}
}

func GetReportforOnlineV3New(icommunication models.Communication, itempName string, txn *gorm.DB) models.TxnError {
	defaultpath := os.Getenv("PDF_SAVE_PATH")
	parts := strings.Split(icommunication.TemplatePath, "/")
	templateFile := parts[len(parts)-1]
	imgFolder := strings.TrimSuffix(templateFile, "."+strings.Split(templateFile, ".")[1])
	remainingPath := strings.Join(parts[:len(parts)-1], "/")
	absolutePath, err := filepath.Abs(remainingPath)
	if err != nil {
		return models.TxnError{ErrorCode: "GL724"}
	}

	staticPath := filepath.Join(absolutePath, "static", imgFolder)
	staticPath = filepath.ToSlash(staticPath)
	//staticPath = "file:///" + strings.ReplaceAll(staticPath, " ", "%20")
	staticPath = toFileURL(staticPath)

	if icommunication.ExtractedData == nil {
		icommunication.ExtractedData = make(map[string]interface{})
	}
	icommunication.ExtractedData["Img"] = staticPath

	basePath := strings.TrimSuffix(templateFile, ".gohtml")
	templateFileWithPath := filepath.Join(remainingPath, templateFile)
	templateFileWithHeaderPath := filepath.Join(remainingPath, basePath+"-h.gohtml")
	templateFileWithFooterPath := filepath.Join(remainingPath, basePath+"-f.gohtml")

	iFile := filepath.Base(strings.TrimSuffix(templateFileWithPath, ".gohtml"))
	hFile := filepath.Base(strings.TrimSuffix(templateFileWithHeaderPath, ".gohtml"))
	fFile := filepath.Base(strings.TrimSuffix(templateFileWithFooterPath, ".gohtml"))

	cwdPath, _ := os.Getwd()
	iPath := filepath.Join(cwdPath, "reportTemplates", "static")
	imgPath := filepath.Join(iPath, iFile)
	imgHeaderPath := filepath.Join(iPath, hFile)
	imgFooterPath := filepath.Join(iPath, fFile)

	ifileContent, err := os.ReadFile(templateFileWithPath)
	if err != nil {
		return models.TxnError{ErrorCode: "GL725"}
	}
	hfileContent, err := os.ReadFile(templateFileWithHeaderPath)
	if err != nil {
		return models.TxnError{ErrorCode: "GL726"}
	}
	ffileContent, err := os.ReadFile(templateFileWithFooterPath)
	if err != nil {
		return models.TxnError{ErrorCode: "GL727"}
	}

	bodyTpl := strings.ReplaceAll(string(ifileContent), "{{.Img}}", imgPath)
	headTpl := strings.ReplaceAll(string(hfileContent), "{{.Img}}", imgHeaderPath)
	footTpl := strings.ReplaceAll(string(ffileContent), "{{.Img}}", imgFooterPath)

	ioutFile := filepath.Join(defaultpath, iFile+"-outfile.html")
	houtFile := filepath.Join(defaultpath, hFile+"-outfile.html")
	foutFile := filepath.Join(defaultpath, fFile+"-outfile.html")

	tempHTMLFiles := []string{ioutFile, houtFile, foutFile}

	if err := createhtml(bodyTpl, icommunication.ExtractedData, ioutFile, ioutFile); err != nil {
		return models.TxnError{ErrorCode: "GL728"}
	}
	if err := createhtml(headTpl, icommunication.ExtractedData, houtFile, houtFile); err != nil {
		return models.TxnError{ErrorCode: "GL729"}
	}
	if err := createhtml(footTpl, icommunication.ExtractedData, foutFile, foutFile); err != nil {
		return models.TxnError{ErrorCode: "GL730"}
	}

	houtFile = toFileURL(houtFile)
	foutFile = toFileURL(foutFile)

	var pdfBuf bytes.Buffer

	finalBody, _ := os.ReadFile(ioutFile)
	r := NewRequestPdfV3(string(finalBody), houtFile, foutFile)

	pdffileName := fmt.Sprintf("%s_%d_%d_%s.pdf", icommunication.TemplateName, icommunication.ClientID, icommunication.PolicyID, time.Now().Format("20060102150405"))

	_, funcErr := r.GeneratePDFPV3New(&pdfBuf, icommunication.CompanyID, icommunication.ClientID, txn)
	if funcErr.ErrorCode != "" {
		return funcErr

	}

	pdfFilePath := filepath.Join(defaultpath, pdffileName)
	if icommunication.PDFPath != "" {
		pdfFilePath = filepath.Join(icommunication.PDFPath, pdffileName)
	}

	pdfFilePath = filepath.ToSlash(filepath.Clean(pdfFilePath))
	if err := os.WriteFile(pdfFilePath, pdfBuf.Bytes(), 0644); err != nil {
		return models.TxnError{ErrorCode: "GL732"}
	}

	if icommunication.EmailAllowed == "Y" {
		funcErr := EmailTriggerMNew(icommunication, pdfBuf.Bytes(), txn)
		if funcErr.ErrorCode != "" {
			return funcErr
		}
	}

	for _, file := range tempHTMLFiles {
		_ = os.Remove(file)
	}

	return models.TxnError{}
}

func (r *RequestPdfV3) GeneratePDFPV3New(output io.Writer, iUserco, iClientid uint, txn *gorm.DB) (bool, models.TxnError) {
	opassword := "FuturaInsTech"
	var clntenq models.Client
	ipassword := ""

	result := txn.First(&clntenq, "company_id = ? and id = ?", iUserco, iClientid)
	if result.RowsAffected == 0 {
		ipassword = opassword
	} else {
		ipassword = strconv.Itoa(int(iClientid)) + clntenq.ClientMobile
	}

	tempHTML := "temp.html"
	if err := os.WriteFile(tempHTML, []byte(r.body), 0644); err != nil {
		return false, models.TxnError{ErrorCode: "GL734"}
	}

	wkhtmlDir := os.Getenv("WKHTMLTOPDF_PATH")
	if wkhtmlDir == "" {
		return false, models.TxnError{ErrorCode: "GL735"}
	}

	wkhtmlPath := filepath.Join(wkhtmlDir, "wkhtmltopdf.exe")
	wkhtmlPath = filepath.ToSlash(wkhtmlPath)

	tempPDF := "temp.pdf"

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(tempPDF), os.ModePerm); err != nil {
		return false, models.TxnError{ErrorCode: "GL736"}
	}

	cmd := exec.Command(
		wkhtmlPath,
		"--enable-local-file-access",
		"--header-html", r.HeaderFile,
		"--footer-html", r.FooterFile,
		"--margin-top", "40mm",
		"--margin-bottom", "50mm",
		"--margin-left", "20mm",
		"--margin-right", "15mm",
		tempHTML,
		filepath.ToSlash(tempPDF),
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return false, models.TxnError{ErrorCode: "GL737"}
	}

	// Password protect PDF
	protectedFile := "protected.pdf"
	if err := EncryptPDF(tempPDF, protectedFile, ipassword, opassword); err != nil {
		return false, models.TxnError{ErrorCode: "GL738"}
	}

	protectedData, err := os.ReadFile(protectedFile)
	if err != nil {
		return false, models.TxnError{ErrorCode: "GL739"}
	}
	if _, err := output.Write(protectedData); err != nil {
		return false, models.TxnError{ErrorCode: "GL740"}
	}

	os.Remove(tempHTML)
	os.Remove(tempPDF)
	os.Remove(protectedFile)

	return true, models.TxnError{}
}

// func EmailTriggerMNew(icommunication models.Communication, pdfData []byte, txn *gorm.DB) models.TxnError {
// 	var client models.Client
// 	result := txn.First(&client, "id = ?", icommunication.ClientID)
// 	if result.Error != nil {
// 		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
// 	}
// 	if client.ClientEmail == "" {
// 		return models.TxnError{ErrorCode: "GL475"}
// 	}

// 	iTemplate := icommunication.TemplateName
// 	var p0033data paramTypes.P0033Data
// 	var extradatap0033 paramTypes.Extradata = &p0033data
// 	errparam := "P0033"
// 	err := GetItemD(int(icommunication.CompanyID), errparam, iTemplate, icommunication.EffectiveDate, &extradatap0033)
// 	if err != nil {
// 		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iTemplate}
// 	}

// 	sender := icommunication.CompanyEmail
// 	receiver := client.ClientEmail
// 	password := p0033data.SenderPassword
// 	smtpServer := p0033data.SMTPServer
// 	smtpPort := p0033data.SMTPPort
// 	emailBody := p0033data.Body
// 	iDateTime := time.Now().Format("20060102150405")
// 	fileName := fmt.Sprintf("%s_%s.pdf", icommunication.TemplateName, iDateTime)

// 	// Send email asynchronously
// 	go func() {
// 		m := gomail.NewMessage()
// 		m.SetHeader("From", sender)
// 		m.SetHeader("To", receiver)
// 		m.SetHeader("Subject", p0033data.Subject)
// 		m.SetBody("text/plain", emailBody)

// 		// Attach PDF file
// 		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
// 			_, err := w.Write(pdfData)
// 			return models.TxnError{ErrorCode: "GL939", DbError: err}
// 		}))

// 		d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
// 		d.SSL = true

// 		sendStart := time.Now()
// 		if err := d.DialAndSend(m); err != nil {
// 			return models.TxnError{ErrorCode: "GL947", DbError: err}
// 		} else {
// 			log.Printf("Email sent successfully to %s in %v", receiver, time.Since(sendStart))
// 		}
// 	}()

// 	// Send Agent Email asynchronously if allowed
// 	if icommunication.AgentEmailAllowed == "Y" {
// 		var agntenq models.Agency
// 		result := txn.First(&agntenq, "id = ?", icommunication.AgencyID)
// 		if result.Error != nil {
// 			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
// 		}
// 		var agclient models.Client
// 		result = txn.First(&agclient, "id = ?", agntenq.ClientID)
// 		if result.Error != nil {
// 			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
// 		}

// 		if agclient.ClientEmail != "" {
// 			go func() {
// 				agentReceiver := agclient.ClientEmail
// 				iName := GetName(client.CompanyID, client.ID)
// 				agentEmailBody := fmt.Sprintf(
// 					"Hi Sir/Madam,\n\nFollowing Email was sent to your Customer %d %s\n\n"+
// 						"I am from Futura Instech..\n\nThank you!",
// 					client.ID, iName,
// 				)

// 				m := gomail.NewMessage()
// 				m.SetHeader("From", sender)
// 				m.SetHeader("To", agentReceiver)
// 				m.SetHeader("Subject", "Mail Sent to Your Customer")
// 				m.SetBody("text/plain", agentEmailBody)

// 				d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
// 				d.SSL = true

// 				sendStart := time.Now()
// 				if err := d.DialAndSend(m); err != nil {
// 					log.Printf("Failed to send email to Agent: %v", err)
// 				} else {
// 					log.Printf("Email sent successfully to agent %s in %v", agentReceiver, time.Since(sendStart))
// 				}
// 			}()
// 		}
// 	}

//		log.Println("Email sent successfully with attachment via office SMTP")
//		return models.TxnError{}
//	}
func EmailTriggerMNew(icommunication models.Communication, pdfData []byte, txn *gorm.DB) models.TxnError {

	var client models.Client
	result := txn.First(&client, "id = ?", icommunication.ClientID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	if client.ClientEmail == "" {
		return models.TxnError{ErrorCode: "GL475"}
	}

	// Get P0033 parameters
	iTemplate := icommunication.TemplateName
	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	errparam := "P0033"

	err := GetItemD(int(icommunication.CompanyID), errparam, iTemplate, icommunication.EffectiveDate, &extradatap0033)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iTemplate}
	}

	sender := icommunication.CompanyEmail
	receiver := client.ClientEmail
	password := p0033data.SenderPassword
	smtpServer := p0033data.SMTPServer
	smtpPort := p0033data.SMTPPort
	emailBody := p0033data.Body
	iDateTime := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.pdf", icommunication.TemplateName, iDateTime)

	// -----------------------------------------------
	//  CUSTOMER EMAIL WITH CHANNEL RETURN (sync)
	// -----------------------------------------------

	errChan := make(chan models.TxnError, 1)

	go func() {
		m := gomail.NewMessage()
		m.SetHeader("From", sender)
		m.SetHeader("To", receiver)
		m.SetHeader("Subject", p0033data.Subject)
		m.SetBody("text/plain", emailBody)

		// Attach PDF
		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(pdfData)
			errChan <- models.TxnError{
				ErrorCode: "GL939",
				DbError:   err,
			}
			return err

		}))

		d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
		d.SSL = true

		if err := d.DialAndSend(m); err != nil {
			txnErr := models.TxnError{
				ErrorCode: "GL947",
				DbError:   err,
			}
			errChan <- txnErr
			return
		}
		return
	}()

	// Wait for email result
	emailResult := <-errChan
	if emailResult.DbError != nil {
		return emailResult
	}

	log.Printf("Email sent successfully to %s", receiver)

	// -----------------------------------------------
	//  AGENT EMAIL (ASYNC – FIRE AND FORGET)
	// -----------------------------------------------

	if icommunication.AgentEmailAllowed == "Y" {

		var agntenq models.Agency
		result := txn.First(&agntenq, "id = ?", icommunication.AgencyID)
		if result.Error != nil {
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		var agclient models.Client
		result = txn.First(&agclient, "id = ?", agntenq.ClientID)
		if result.Error != nil {
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		if agclient.ClientEmail != "" {

			go func() {
				agentReceiver := agclient.ClientEmail
				iName := GetName(client.CompanyID, client.ID)

				agentEmailBody := fmt.Sprintf(
					"Hi Sir/Madam,\n\nFollowing Email was sent to your Customer %d %s\n\n"+
						"I am from Futura Instech..\n\nThank you!",
					client.ID, iName,
				)

				m := gomail.NewMessage()
				m.SetHeader("From", sender)
				m.SetHeader("To", agentReceiver)
				m.SetHeader("Subject", "Mail Sent to Your Customer")
				m.SetBody("text/plain", agentEmailBody)

				d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
				d.SSL = true

				if err := d.DialAndSend(m); err != nil {
					log.Printf("Failed to send email to Agent: %v", err)
				} else {
					log.Printf("Email sent successfully to agent %s", agentReceiver)
				}
			}()
		}
	}

	log.Println("Email sent successfully with attachment via office SMTP")
	return models.TxnError{}
}

func PostGlMoveNNew(iCompany uint, iContractCurry string, iEffectiveDate string,
	iTranno int, iGlAmount float64, iAccAmount float64, iAccountCodeID uint, iGlRdocno uint,
	iGlRldgAcct string, iSeqnno uint64, iGlSign string, iAccountCode string, iHistoryCode string,
	iRevInd string, iCoverage string, txn *gorm.DB) models.TxnError {

	iAccAmount = RoundFloat(iAccAmount, 2)

	var glmove models.GlMove
	var company models.Company
	glmove.ContractCurry = iContractCurry
	glmove.ContractAmount = iAccAmount
	result := txn.Find(&company, "id = ?", iCompany)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL002", DbError: result.Error}
	}
	var currency models.Currency
	// fmt.Println("Currency Code is .... ", company.CurrencyID)
	result = txn.Find(&currency, "id = ?", company.CurrencyID)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL672", DbError: result.Error}
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
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	//tx := initializers.DB.Save(&glmove)
	//tx.Commit()

	UpdateGlBalNNew(iCompany, iGlRldgAcct, iAccountCode, iContractCurry, iAccAmount, iGlSign, GlRdocno, txn)
	return models.TxnError{}
}

func UpdateGlBalNNew(iCompany uint, iGlRldgAcct string, iGlAccountCode string, iContCurry string, iAmount float64, iGLSign string, iGlRdocno string, txn *gorm.DB) (models.TxnError, float64) {
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
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}, 0
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
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}, glbal.ContractAmount
		}
		return models.TxnError{}, glbal.ContractAmount
	} else {
		iAmount := glbal.ContractAmount + temp
		// fmt.Println("I am inside update.....2", iAmount, glbal.ContractAmount)
		//initializers.DB.Model(&glbal).Where("company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno).Update("contract_amount", iAmount)
		result = txn.Model(&glbal).Where("company_id = ? and gl_accountno = ? and gl_rldg_acct = ? and contract_curry = ? and gl_rdocno = ?", iCompany, iGlAccountCode, iGlRldgAcct, iContCurry, iGlRdocno).Update("contract_amount", iAmount)
		if result.Error != nil {
			return models.TxnError{ErrorCode: "GL721"}, glbal.ContractAmount
		}

		return models.TxnError{}, glbal.ContractAmount
	}
	//results.Commit()

}

func TDFLoanDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iDate string, txn *gorm.DB) (string, models.TxnError) {
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = iDate
		tdfpolicy.Tranno = iTranno
		tdfpolicy.Seqno = tdfrule.Seqno
		txn.Create(&tdfpolicy)
		return "", models.TxnError{}
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
		return "", models.TxnError{}
	}
}

// 2025-10-15 Lakshmi Changes
func GetReportforOnlineNew(icommuncation models.Communication, itempName string, txn *gorm.DB) models.TxnError {
	defaultpath := os.Getenv("REPORTPDF_SAVE_PATH")
	parts := strings.Split(icommuncation.TemplatePath, "/")
	templateFile := parts[len(parts)-1] // Extract gohtml file name

	imgFolder := strings.TrimSuffix(templateFile, "."+strings.Split(templateFile, ".")[1])

	remainingPath := strings.Join(parts[:len(parts)-1], "/")
	absolutePath, err := filepath.Abs(remainingPath)
	if err != nil {
		return models.TxnError{ErrorCode: "GL701"}
	}

	iPath := filepath.Join(absolutePath, "static")
	imPath := filepath.Join(iPath, imgFolder)

	imagePath := strings.ReplaceAll(imPath, "\\", "/")

	// Ensure ExtractedData is initialized
	if icommuncation.ExtractedData == nil {
		icommuncation.ExtractedData = make(map[string]interface{})
	}
	icommuncation.ExtractedData["Img"] = imagePath

	// Parse and execute template
	funcMap := CreateFuncMap()
	tmpl, err := template.New(templateFile).Funcs(funcMap).ParseFiles(icommuncation.TemplatePath)
	if err != nil {
		return models.TxnError{ErrorCode: "GL702"}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, icommuncation.ExtractedData)
	if err != nil {
		return models.TxnError{ErrorCode: "GL703"}
	}

	// Create PDF from the template execution output
	r := NewRequestPdf(buf.String())
	pdffileName := fmt.Sprintf("%s_%d_%d_%s.pdf", icommuncation.TemplateName, icommuncation.ClientID, icommuncation.PolicyID, time.Now().Format("20060102150405"))

	var pdfBuf bytes.Buffer
	success, funcErr := r.GeneratePDFPN(&pdfBuf, icommuncation.CompanyID, icommuncation.ClientID, txn)
	if funcErr.ErrorCode != "" || !success {
		return funcErr
	}

	// Save the PDF to the file system if needed
	comFileName := filepath.Join(defaultpath, pdffileName)
	if icommuncation.PDFPath != "" {
		comFileName = filepath.Join(icommuncation.PDFPath, pdffileName)
	}
	comFileName = filepath.ToSlash(filepath.Clean(comFileName))

	err = os.WriteFile(comFileName, pdfBuf.Bytes(), 0644)
	if err != nil {
		return models.TxnError{ErrorCode: "GL705"}
	}

	// Send email if allowed
	if icommuncation.EmailAllowed == "Y" {
		err = EmailTrigger(icommuncation, itempName, pdfBuf.Bytes(), txn)
		if err != nil {
			return models.TxnError{ErrorCode: "GL706", DbError: err}
		}
	}

	// Return the generated PDF buffer
	return models.TxnError{}
}

func SendSMSTwilioNew(iCompany, iclientID uint, itempName, iEffDate string, message string, txn *gorm.DB) models.TxnError {
	// Fetch client details
	var client models.Client
	result := txn.First(&client, "id = ?", iclientID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	errparam := "P0033"
	err := GetItemD(int(iCompany), errparam, itempName, iEffDate, &extradatap0033)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: itempName}
	}

	toNumber := client.ClientMobCode + client.ClientMobile
	accountSID := p0033data.SMSSID
	authToken := p0033data.SMSAuthToken
	fromNumber := p0033data.SMSAuthPhoneNo
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"

	// Prepare message payload
	msgData := url.Values{}
	msgData.Set("To", toNumber)
	msgData.Set("From", fromNumber)
	msgData.Set("Body", message)
	msgDataReader := strings.NewReader(msgData.Encode())

	// Send SMS asynchronously
	go func() {
		startTime := time.Now()
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(urlStr)
		req.Header.SetMethod("POST")
		req.Header.SetContentType("application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Basic "+basicAuth(accountSID, authToken))
		req.SetBodyStream(msgDataReader, msgDataReader.Len())

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		client := fasthttp.Client{}
		err := client.Do(req, resp)
		if err != nil {
			log.Printf("Failed to send SMS to %s: %v", toNumber, err)
			return
		}

		if resp.StatusCode() == 201 {
			log.Printf("SMS sent successfully to %s in %v", toNumber, time.Since(startTime))
		} else {
			log.Printf("Failed to send SMS to %s, response: %v", toNumber, resp.StatusCode())
		}
	}()

	log.Println("SMS sending initiated asynchronously")
	return models.TxnError{}
}

func CreateCommunicationsNNew(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB, iBenefit uint, iPa uint, iClientWork uint) models.TxnError {

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
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL014", DbError: result.Error}
		}
		iReceiptFor = receipt.ReceiptFor
	}

	if iPolicy != 0 {
		var policy models.Policy
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL175", DbError: result.Error}
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
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL671", DbError: result.Error}
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
	errparam := "P0034"
	err1 := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0034)
	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0034)
		if err1 != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
		}
	}

	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			errparam := "P0033"
			err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0033)
			if err != nil {
				return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
					oData, funcErr := GetCompanyDataNew(iCompany, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData, funcErr := GetClientDataNew(iCompany, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData, funcErr := GetAddressDataNew(iCompany, iAddress, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData, funcErr := GetPolicyDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData, funcErr := GetBenefitDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData, funcErr := GetSurBDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData, funcErr := GetMrtaDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData, funcErr := GetReceiptDataNew(iCompany, iReceipt, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData, funcErr := GetSaChangeDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData, funcErr := GetCompAddDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData, funcErr := GetSurrHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData, funcErr := GetDeathDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData, funcErr := GetMatHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData, funcErr := GetSurvBPayNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData, funcErr := GetExpiNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData, funcErr := GetBonusValsNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData, funcErr := GetAgencyNew(iCompany, iAgency, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData, funcErr := GetNomiDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData, funcErr := GetGLDataNew(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData, funcErr := GetIlpSummaryDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData, funcErr := GetIlpAnnsummaryDataNew(iCompany, iPolicy, iHistoryCode, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData, funcErr := GetIlpTranctionDataNew(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData, funcErr := GetPremTaxGLDataNew(iCompany, iPolicy, iFromDate, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData, funcErr := GetIlpFundSwitchDataNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData, funcErr := GetPHistoryDataNew(iCompany, iPolicy, iHistoryCode, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData, funcErr := GetIlpFundDataNew(iCompany, iPolicy, iBenefit, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData, funcErr := GetPPolicyDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData, funcErr := GetPBenefitDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					fmt.Println(oData) // Dummy to avoid compilation error
				case oLetType == "29":
					oData, funcErr := GetPayingAuthorityDataNew(iCompany, iPa, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevBenefit"] = oData
				case oLetType == "30":
					oData, funcErr := GetClientWorkDataNew(iCompany, iClientWork, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ClientWork"] = oData
				case oLetType == "36":
					oData, funcErr := GetReqDataNew(iCompany, iPolicy, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for _, item := range oData {
						for key, value := range item.(map[string]interface{}) {
							resultMap[key] = value
						}
					}
				case oLetType == "37":
					oData, funcErr := PolicyDepDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "38":
					oData, funcErr := PolAgntChDataNew(iCompany, iPolicy, iAgency, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "39":
					oData, funcErr := GetBankDataNew(iCompany, iClientWork, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BankData"] = oData
				case oLetType == "40":
					iKey := iReceipt
					oData, funcErr := GetPaymentDataNew(iCompany, iPolicy, iKey, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "45":
					oData, funcErr := ColaCancelDataNew(iCompany, iPolicy, iHistoryCode, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "46":
					oData, funcErr := AplCancelDataNew(iCompany, iPolicy, iHistoryCode, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					for key, value := range oData {
						resultMap[key] = value
					}
				// case oLetType == "47":
				// 	oData := GetPOLSCDEndowmentData(iCompany, iPolicy, iPageSize, iOrientation, p0033data, txn)
				// 	for key, value := range oData {
				// 		resultMap[key] = value
				// 	}
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
				funcErr := GetReportforOnlineNew(communication, p0033data.TemplateName, txn)
				if funcErr.ErrorCode != "" {
					return funcErr
				}
			}
			if p0033data.SMSAllowed == "Y" {
				funcErr := SendSMSTwilioNew(communication.CompanyID, communication.ClientID, p0033data.TemplateName, communication.EffectiveDate, p0033data.SMSBody, txn)
				if funcErr.ErrorCode != "" {
					return funcErr
				}
			}
			communication.Print = "Y"
			communication.PrintDate = iDate
			communication.UpdatedID = 1
			communication.ID = 0
			// New Changes Ended
			results := txn.Create(&communication)

			if results.Error != nil {
				return models.TxnError{
					ErrorCode: "DBERR",
					DbError:   results.Error,
				}
			}

		}
	}
	return models.TxnError{}
}

func GetMaxTrannoNNew(iCompany uint, iPolicy uint, iMethod string, iEffDate string, iuser uint64, historyMap map[string]interface{}, txn *gorm.DB) (string, uint, models.TxnError) {
	var permission models.Permission
	var result *gorm.DB

	result = txn.First(&permission, "company_id = ? and method = ?", iCompany, iMethod)
	if result.Error != nil {
		return iMethod, 0, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	iHistoryCode := permission.TransactionID
	var transaction models.Transaction
	result = txn.Find(&transaction, "ID = ?", iHistoryCode)
	if result.RowsAffected == 0 {
		return iMethod, 0, models.TxnError{ErrorCode: "GL042", DbError: result.Error}
	}
	iHistoryCD := transaction.TranCode
	var phistory models.PHistory
	var maxtranno float64 = 0

	fmt.Println(iCompany, iPolicy, iHistoryCD, iEffDate)

	result1 := txn.Table("p_histories").Where("company_id = ? and policy_id= ?", iCompany, iPolicy).Select("max(tranno)")

	if result1.Error != nil {
		fmt.Println(models.TxnError{ErrorCode: "GL919", DbError: result1.Error})

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
		return phistory.HistoryCode, phistory.Tranno, models.TxnError{ErrorCode: "DBERR",
			DbError: result.Error}
	}

	return phistory.HistoryCode, phistory.Tranno, models.TxnError{}
}

// 2025-10-16 Divya Changes
func ValidateBankNew(bankval models.Bank, userco uint, userlan uint, iKey string) models.TxnError {
	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	// Fetch validation rules
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
	}

	// Loop through validation fields
	for i := 0; i < len(p0065data.FieldList); i++ {
		var fv interface{}
		r := reflect.ValueOf(bankval)
		f := reflect.Indirect(r).FieldByName(p0065data.FieldList[i].Field)

		if f.IsValid() {
			fv = f.Interface()
		} else {
			continue
		}

		if isFieldZero(fv) {
			errcode := p0065data.FieldList[i].ErrorCode
			return models.TxnError{
				ErrorCode: errcode,
			}
		}
	}

	// Special date check
	if bankval.StartDate > bankval.EndDate {
		return models.TxnError{
			ErrorCode: "GL563",
		}
	}

	return models.TxnError{} // no error
}

// 2025-10-21 Lakshmi Changes
func ValidateClientWorkNNew(clientwork models.ClientWork, userco uint, userlan uint, iDate string, iKey string, txn *gorm.DB) models.TxnError {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
			return models.TxnError{ErrorCode: "GL756"}
		}

	}

	var client models.Client
	clientid := clientwork.ClientID
	result1 := txn.Find(&client, "company_id = ? and id = ?", userco, clientid)
	if result1.RowsAffected == 0 {
		return models.TxnError{
			ErrorCode: "GL050",
			DbError:   result1.Error,
		}
	}

	if client.ClientStatus != "AC" {
		return models.TxnError{ErrorCode: "GL006", DbError: result1.Error}
	}
	var employer models.Client
	employerid := clientwork.EmployerID
	result2 := txn.Find(&employer, "company_id = ? and id = ?", userco, employerid)
	if result2.RowsAffected == 0 {
		return models.TxnError{
			ErrorCode: "GL050",
			DbError:   result2.Error,
		}
	}

	if employer.ClientStatus != "AC" {
		return models.TxnError{ErrorCode: "GL006", DbError: result2.Error}
	}

	if clientwork.StartDate > iDate {
		return models.TxnError{ErrorCode: "GL656"}
	}

	if clientwork.EndDate < iDate {
		return models.TxnError{ErrorCode: "GL657"}
	}
	return models.TxnError{}
}

func AutoPayCreateNew(iCompany uint, iPolicy uint, iClient uint, iAddress uint, iBank uint, iAccCurr string, iAmount float64, iDate string, iDrAcc string, iCrAcc string, iTypeofPayment string, iUserID uint, iReason string, iHistoryCode string, iTranno uint, iPayStatus string, iCoverage string, txn *gorm.DB) (oPayno uint, txnerr models.TxnError) {
	if iPayStatus == "PN" {
		var payosbal models.PayOsBal
		result := txn.Find(&payosbal, "company_id = ? and gl_accountno = ? and gl_rldg_acct =? and contract_curry = ?", iCompany, iDrAcc, iPolicy, iAccCurr)
		//	iErr := "Payment Already Processed"
		if result.RowsAffected > 0 {
			txnerr = models.TxnError{ErrorCode: "GL709", DbError: result.Error}
			return 0, txnerr
		}
	}

	oPayno = 0
	var bankenq models.Bank
	result := txn.Find(&bankenq, "id = ?", iBank)
	if result.RowsAffected == 0 {
		txnerr = models.TxnError{ErrorCode: "GL262", DbError: result.Error}
		return oPayno, txnerr
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
	errparam := "P0055"
	err := GetItemD(int(iCompany), errparam, iTypeofPayment, iDate, &extradatap0055)
	if err != nil {
		txnerr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iTypeofPayment}
		return oPayno, txnerr
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
		txnerr = models.TxnError{
			ErrorCode: "DBERR",
			DbError:   result.Error,
		}
		return oPayno, txnerr
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
			txnerr = models.TxnError{
				ErrorCode: "DBERR",
				DbError:   result.Error,
			}
			return 0, txnerr
		}

	}
	if iPayStatus == "AP" {
		// Debit
		glcode := iDrAcc
		var acccode models.AccountCode
		result = txn.First(&acccode, "company_id = ? and account_code = ? ", iCompany, glcode)
		if result.Error != nil {
			txnerr = models.TxnError{
				ErrorCode: "DBERR",
				DbError:   result.Error,
			}
			return oPayno, txnerr
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

		funcErr := PostGlMoveNNew(iCompany, iAccCurry, iEffectiveDate, int(iTranno), iGlAmount,
			iAccAmount, iAccountCodeID, uint(iGlRdocno), string(iGlRldgAcct), iSequenceno, iGlSign, iAccountCode, iHistoryCode, "", "", txn)

		if funcErr.ErrorCode != "" {
			txnerr = funcErr
			return oPayno, txnerr
		}
		// Credit

		glcode = iCrAcc
		var acccode1 models.AccountCode
		result = txn.First(&acccode1, "company_id = ? and account_code = ? ", iCompany, glcode)
		if result.Error != nil {
			txnerr = models.TxnError{
				ErrorCode: "DBERR",
				DbError:   result.Error,
			}
			return oPayno, txnerr
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

		funcErr = PostGlMoveNNew(iCompany, iAccCurry, iEffectiveDate, int(iTranno), iGlAmount,
			iAccAmount, iAccountCodeID, uint(iGlRdocno), string(iGlRldgAcct), iSequenceno, iGlSign, iAccountCode, iHistoryCode, "", "", txn)

		if funcErr.ErrorCode != "" {
			txnerr = funcErr
			return oPayno, txnerr
		}
	}

	return oPayno, txnerr
}

func EmailTriggerforReportNew(iCompany uint, iReference uint, iClient uint, iEmail string, iEffDate string, itempName string, pdfData []byte, txn *gorm.DB) models.TxnError {

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	errparam := "P0033"
	err := GetItemD(int(iCompany), errparam, itempName, iEffDate, &extradatap0033)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: itempName}
	}

	sender := p0033data.CompanyEmail
	receiver := iEmail
	password := p0033data.SenderPassword
	smtpServer := p0033data.SMTPServer
	smtpPort := p0033data.SMTPPort

	emailBody := p0033data.Body
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", p0033data.Subject)
	m.SetBody("text/plain", emailBody)
	iTime := time.Now().Format("20060102150405")
	iClientnumstr := strconv.Itoa(int(iClient))
	errChan := make(chan models.TxnError, 1)
	m.Attach(itempName+iClientnumstr+iTime+".pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfData)
		errChan <- models.TxnError{
			ErrorCode: "GL939",
			DbError:   err,
		}
		return err

	}))

	// Configure SMTP dialer
	d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
	d.SSL = true      // Enables SSL
	d.TLSConfig = nil // Use default TLS settings

	// Send email asynchronously with proper logging
	sendStart := time.Now()
	go func() {
		if err := d.DialAndSend(m); err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent successfully to %s (CC: %s, BCC: %s) in %v",
				receiver, "", "", time.Since(sendStart))
		}
	}()
	log.Printf("EmailTrigger function executed in %v", time.Since(sendStart))
	return models.TxnError{}
}

// 2025-10-22 Lakshmi Changes
func (r *RequestPdf) GeneratePDFPN(inputFile io.Writer, iUserco uint, iClientid uint, txn *gorm.DB) (bool, models.TxnError) {

	opassword := "FuturaInsTech"
	var clntenq models.Client
	ipassword := ""

	result := txn.First(&clntenq, "company_id = ? and id = ?", iUserco, iClientid)
	// In case no record found, use owner password as user password
	if result.RowsAffected == 0 {
		ipassword = opassword
	} else {
		ipassword = strconv.Itoa(int(iClientid)) + clntenq.ClientMobile
	}
	// Step 1: Generate the PDF
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return false, models.TxnError{ErrorCode: "GL754", DbError: err}
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(r.body))
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	//pdfg.Orientation.Set(wkhtmltopdf.)
	pdfg.Dpi.Set(300)

	// Save to temporary file
	tempFile := "temp.pdf"
	outFile, err := os.Create(tempFile)
	if err != nil {
		return false, models.TxnError{ErrorCode: "GL755", DbError: err}
	}
	defer outFile.Close()

	pdfg.SetOutput(outFile)
	err = pdfg.Create()
	if err != nil {
		return false, models.TxnError{ErrorCode: "GL712", DbError: err}
	}

	// Step 2: Protect the PDF using Python script
	protectedFile := "protected.pdf"
	err = EncryptPDF(tempFile, protectedFile, ipassword, opassword)
	if err != nil {
		return false, models.TxnError{ErrorCode: "GL756", DbError: err}
	}

	// Step 3: Write the password-protected PDF to the writer
	protectedData, err := os.ReadFile(protectedFile)
	if err != nil {
		return false, models.TxnError{ErrorCode: "GL715", DbError: err}
	}
	_, err = inputFile.Write(protectedData)
	if err != nil {
		return false, models.TxnError{ErrorCode: "GL757", DbError: err}
	}

	// Cleanup temporary files
	os.Remove(tempFile)
	os.Remove(protectedFile)

	return true, models.TxnError{}
}

// 2025-10-29 Lakshmi Changes
func EmailTriggerNNew(icommunication models.Communication, pdfData []byte, txn *gorm.DB) models.TxnError {
	var client models.Client
	result := txn.First(&client, "id = ?", icommunication.ClientID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	if client.ClientEmail == "" {
		return models.TxnError{ErrorCode: "GL770"}
	}

	iTemplate := icommunication.TemplateName
	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	errparam := "P0033"
	err := GetItemD(int(icommunication.CompanyID), errparam, iTemplate, icommunication.EffectiveDate, &extradatap0033)

	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iTemplate}
	}

	sender := icommunication.CompanyEmail
	receiver := client.ClientEmail
	password := p0033data.SenderPassword
	smtpServer := p0033data.SMTPServer
	smtpPort := p0033data.SMTPPort
	emailBody := p0033data.Body
	iDateTime := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.pdf", icommunication.TemplateName, iDateTime)

	// Send email asynchronously
	errChan := make(chan models.TxnError, 1)

	go func() {
		m := gomail.NewMessage()
		m.SetHeader("From", sender)
		m.SetHeader("To", receiver)
		m.SetHeader("Subject", p0033data.Subject)
		m.SetBody("text/plain", emailBody)

		// Attach PDF
		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(pdfData)
			errChan <- models.TxnError{
				ErrorCode: "GL939",
				DbError:   err,
			}
			return err

		}))

		d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
		d.SSL = true

		if err := d.DialAndSend(m); err != nil {
			txnErr := models.TxnError{
				ErrorCode: "GL947",
				DbError:   err,
			}
			errChan <- txnErr
			return
		}
		return
	}()

	// Wait for email result
	emailResult := <-errChan
	if emailResult.DbError != nil {
		return emailResult
	}

	log.Printf("Email sent successfully to %s", receiver)

	// Send Agent Email asynchronously if allowed
	if icommunication.AgentEmailAllowed == "Y" {
		var agntenq models.Agency
		result := txn.First(&agntenq, "id = ?", icommunication.AgencyID)
		if result.Error != nil {
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		var agclient models.Client
		result = txn.First(&agclient, "id = ?", agntenq.ClientID)
		if result.Error != nil {
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		if agclient.ClientEmail != "" {
			go func() {
				agentReceiver := agclient.ClientEmail
				iName := GetName(client.CompanyID, client.ID)
				agentEmailBody := fmt.Sprintf(
					"Hi Sir/Madam,\n\nFollowing Email was sent to your Customer %d %s\n\n"+
						"I am from Futura Instech..\n\nThank you!",
					client.ID, iName,
				)

				m := gomail.NewMessage()
				m.SetHeader("From", sender)
				m.SetHeader("To", agentReceiver)
				m.SetHeader("Subject", "Mail Sent to Your Customer")
				m.SetBody("text/plain", agentEmailBody)

				d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
				d.SSL = true

				sendStart := time.Now()
				if err := d.DialAndSend(m); err != nil {
					log.Printf("Failed to send email to Agent: %v", err)
				} else {
					log.Printf("Email sent successfully to agent %s in %v", agentReceiver, time.Since(sendStart))
				}
			}()
		}
	}

	log.Println("Email sent successfully with attachment via office SMTP")
	return models.TxnError{}
}

func TDFReraDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
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
	return "", models.TxnError{}
}

func TDFAnniPNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var annuity models.Annuity
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}
	result = txn.First(&annuity, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	if result.Error != nil {
		//	txn.Rollback()
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		return "", models.TxnError{}
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		return "", models.TxnError{}
	}
}

func TDFLoanIntNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var loanenq []models.Loan
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_status = ? ", iCompany, iPolicy, "AC")
	loandelete := "N"
	if result.Error != nil || result.RowsAffected == 0 {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
		} else {
			loandelete = "Y"
		}
	}

	if loandelete == "Y" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error == nil {
			results = txn.Delete(&tdfpolicy)
			if results.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
			}
			return "", models.TxnError{}
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
			results = txn.Create(&tdfpolicy)
			if results.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
			}
			return "", models.TxnError{}
		} else {
			results = txn.Delete(&tdfpolicy)
			if results.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
			}
			var tdfpolicy models.TDFPolicy
			tdfpolicy.CompanyID = iCompany
			tdfpolicy.PolicyID = iPolicy
			tdfpolicy.Seqno = tdfrule.Seqno
			tdfpolicy.TDFType = iFunction
			tdfpolicy.ID = 0
			tdfpolicy.EffectiveDate = oDate
			tdfpolicy.Tranno = iTranno

			results = txn.Create(&tdfpolicy)
			if results.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFLoanCapNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var loanenq []models.Loan
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_status = ? ", iCompany, iPolicy, "AC")
	loandelete := "N"
	if result.Error != nil || result.RowsAffected == 0 {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		} else {
			loandelete = "Y"
		}
	}

	if loandelete == "Y" {
		results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
		if results.Error == nil {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
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
			result = txn.Create(&tdfpolicy)
			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFLapsDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	errparam := "Q0005"
	err := GetItemD(int(iCompany), errparam, policy.PProduct, policy.PRCD, &extradataq0005)

	if err != nil {
		return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: policy.PProduct}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		return "", models.TxnError{}
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		return "", models.TxnError{}
	}

}

func TDFSurvbDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var survb models.SurvB
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}
	result = txn.First(&survb, "company_id = ? and policy_id = ? and paid_date = ?", iCompany, iPolicy, "")

	if result.Error != nil {
		//	txn.Rollback()
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		return "", models.TxnError{}
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		return "", models.TxnError{}
	}
}

func TDFAnniDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		return "", models.TxnError{}
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		return "", models.TxnError{}
	}
}

func TDFFundPNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string, txn *gorm.DB) (string, models.TxnError) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var ilptransenq []models.IlpTransaction
	odate := "00000000"

	result := txn.Where("company_id = ? and policy_id = ? and ul_process_flag = ?", iCompany, iPolicy, "P").Order("fund_eff_date").Find(&ilptransenq)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL137", DbError: result.Error}
	}

	for i := 0; i < len(ilptransenq); i++ {
		if ilptransenq[i].FundEffDate > odate {
			odate = ilptransenq[i].FundEffDate
		}
	}

	result = txn.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL017", DbError: result.Error}
	}

	result = txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		} else {
			txn.Delete(&tdfpolicy)
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFFundMNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {

	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"

	result := txn.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL017", DbError: result.Error}
	}

	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL018", DbError: result.Error}
	}
	for i := 0; i < len(benefitenq); i++ {

		iCoverage := benefitenq[i].BCoverage
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data
		errparam := "Q0006"
		err := GetItemD(int(iCompany), errparam, iCoverage, benefitenq[i].BStartDate, &extradataq0006)
		if err != nil {
			return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}
		}
		if q0006data.PremCalcType == "U" {

			if benefitenq[i].IlpMortalityDate > odate {
				odate = benefitenq[i].IlpMortalityDate
			}
		}
	}

	result = txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}

		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFFundFNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {

	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"

	result := txn.Find(&policy, "company_id = ? and id  = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL017", DbError: result.Error}
	}

	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL018", DbError: result.Error}
	}
	for i := 0; i < len(benefitenq); i++ {

		iCoverage := benefitenq[i].BCoverage
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data
		errparam := "Q0006"
		err := GetItemD(int(iCompany), errparam, iCoverage, benefitenq[i].BStartDate, &extradataq0006)
		if err != nil {
			return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}
		}
		if q0006data.PremCalcType == "U" {
			if benefitenq[i].IlpFeeDate > odate {
				odate = benefitenq[i].IlpFeeDate
			}
		}
	}

	result = txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFIBDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var incomeb models.IBenefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		// txn.Rollback()
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.First(&incomeb, "company_id = ? and policy_id = ? and paid_date = ?", iCompany, iPolicy, "")
	if result.Error != nil {
		// txn.Rollback()
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	result = txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ? ", iCompany, iPolicy, iFunction)

	if result.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = incomeb.NextPayDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = incomeb.NextPayDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		return "", models.TxnError{}
	}
	return "", models.TxnError{}
}

func TDFExtrDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var extraenq []models.Extra
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var policyenq models.Policy
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.Find(&extraenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return "", models.TxnError{ErrorCode: "GL236", DbError: result.Error}
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
	result = txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL017", DbError: result.Error}
	}
	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	errparam := "Q0005"
	err := GetItemD(int(iCompany), errparam, policyenq.PProduct, policyenq.PRCD, &extradataq0005)
	if err != nil {
		return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: policyenq.PProduct}
	}
	if oDate != "" {
		oDate = AddLeadDays(oDate, (-1 * q0005data.BillingLeadDays))
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFExpiDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL018", DbError: result.Error}
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		if benefits[i].BStatus != "EX" {
			iCoverage := benefits[i].BCoverage
			iDate := benefits[i].BStartDate
			var q0006data paramTypes.Q0006Data
			var extradataq0006 paramTypes.Extradata = &q0006data
			errparam := "Q0006"
			err := GetItemD(int(iCompany), errparam, iCoverage, iDate, &extradataq0006)
			if err != nil {
				return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}

			return "", models.TxnError{}
		} else {
			result = txn.Delete(&tdfpolicy)
			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}

			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFMatDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "IF")
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL018", DbError: result.Error}
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		iCoverage := benefits[i].BCoverage
		iDate := benefits[i].BStartDate
		var q0006data paramTypes.Q0006Data
		var extradataq0006 paramTypes.Extradata = &q0006data
		errparam := "Q0006"
		err := GetItemD(int(iCompany), errparam, iCoverage, iDate, &extradataq0006)
		if err != nil {
			return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}
		}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		} else {
			result = txn.Delete(&tdfpolicy)

			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

func TDFBillDNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string, txn *gorm.DB) (string, models.TxnError) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	var benefitenq []models.Benefit
	odate := "00000000"
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return "", models.TxnError{ErrorCode: "GL018", DbError: result.Error}
	}
	for i := 0; i < len(benefitenq); i++ {
		if benefitenq[i].BPremCessDate > odate {
			odate = benefitenq[i].BPremCessDate
		}
	}

	result = txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	if iRevFlag == "R" {
		var q0005data paramTypes.Q0005Data
		var extradataq0005 paramTypes.Extradata = &q0005data
		errparam := "Q0005"
		err := GetItemD(int(iCompany), errparam, policy.PProduct, policy.PRCD, &extradataq0005)
		if err != nil {
			return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: policy.PProduct}
		}

		nxtBtdate := AddLeadDays(policy.PaidToDate, (-1 * q0005data.BillingLeadDays))
		policy.NxtBTDate = nxtBtdate
	}

	if policy.PaidToDate >= odate {
		// return "Date Exceeded", errors.New("Premium Cessation Date is Exceeded")
		var tdfpolicyupd models.TDFPolicy
		result = txn.Find(&tdfpolicyupd, "company_id = ? AND policy_id = ? and tdf_type= ?", iCompany, iPolicy, "BILLD")
		if result.RowsAffected == 0 {
			return "", models.TxnError{ErrorCode: "GL392", DbError: result.Error}
		}
		result = txn.Delete(&tdfpolicyupd)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

		}
		return "", models.TxnError{}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		return "", models.TxnError{}
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
		return "", models.TxnError{}
	}
}

// 2025-10-30 Divya Changes
// Commented as Lakshmi also give this code...
// func PostUlpDeductionByUnitsNNew(iCompany uint, iPolicy uint, iBenefit uint, iSurrPercentage float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string, txn *gorm.DB) (txnErr models.TxnError) {

// 	var policyenq models.Policy

// 	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
// 	if result.RowsAffected == 0 {
// 		txnErr = models.TxnError{ErrorCode: "GL003", DbError: result.Error}
// 		return
// 	}

// 	var p0061data paramTypes.P0061Data
// 	var extradatap0061 paramTypes.Extradata = &p0061data

// 	var p0059data paramTypes.P0059Data
// 	var extradatap0059 paramTypes.Extradata = &p0059data

// 	iKey := iHistoryCode + iBenefitCode + iallocType
// 	errparam := "P0059"
// 	err := GetItemD(int(iCompany), errparam, iKey, iStartDate, &extradatap0059)
// 	if err != nil {
// 		txnErr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
// 		return
// 	}

// 	var ilpfundenq []models.IlpFund

// 	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
// 	if result.RowsAffected == 0 {
// 		txnErr = models.TxnError{ErrorCode: "GL746", DbError: result.Error}
// 		return
// 	}

// 	var ilpsumenq []models.IlpSummary

// 	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
// 	if result.RowsAffected == 0 {
// 		txnErr = models.TxnError{ErrorCode: "GL747", DbError: result.Error}
// 		return
// 	}

// 	// Get Total Fund Value
// 	iTotalFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, "", iEffDate)

// 	for j := 0; j < len(ilpsumenq); j++ {
// 		iBusinessDate := GetBusinessDate(iCompany, 0, 0)
// 		if p0059data.CurrentOrFuture == "F" {
// 			iBusinessDate = AddLeadDays(iBusinessDate, 1)
// 		} else if p0059data.CurrentOrFuture == "E" {
// 			iBusinessDate = iEffDate
// 		}
// 		iFundCode := ilpsumenq[j].FundCode
// 		iFundValue, _, _ := GetAllFundValueByBenefit(iCompany, iPolicy, iBenefit, iFundCode, iEffDate)
// 		var ilptrancrt models.IlpTransaction
// 		iKey := ilpsumenq[j].FundCode
// 		errparam = "P0061"
// 		err := GetItemD(int(iCompany), "P0061", iKey, iStartDate, &extradatap0061)
// 		if err != nil {
// 			txnErr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
// 			return
// 		}

// 		ilptrancrt.CompanyID = iCompany
// 		ilptrancrt.PolicyID = iPolicy
// 		ilptrancrt.BenefitID = iBenefit
// 		ilptrancrt.FundCode = ilpsumenq[j].FundCode
// 		ilptrancrt.FundType = ilpsumenq[j].FundType
// 		ilptrancrt.TransactionDate = iEffDate
// 		ibidprice, _, ipriceuseddate := GetFundCPrice(iCompany, ilpsumenq[j].FundCode, iBusinessDate)
// 		ilptrancrt.FundPrice = ibidprice
// 		ilptrancrt.FundEffDate = ipriceuseddate
// 		iUnits, _ := GetIlpFundUnits(iCompany, iPolicy, iBenefit, iFundCode)
// 		// Full Withdrawl is -100% and Part Withdrawl is -20% or -30% etc
// 		iSurrUnits := iUnits * iSurrPercentage / 100
// 		ilptrancrt.FundUnits = RoundFloat(iSurrUnits, 5)
// 		//utilities.RoundFloat(ilptrancrt.FundAmount/ibidprice, 5)
// 		ilptrancrt.FundAmount = RoundFloat((iSurrUnits * ibidprice), 2)
// 		ilptrancrt.FundCurr = p0061data.FundCurr
// 		ilptrancrt.CurrentOrFuture = p0059data.CurrentOrFuture
// 		ilptrancrt.OriginalAmount = RoundFloat((iSurrUnits * ibidprice), 2)
// 		ilptrancrt.ContractCurry = policyenq.PContractCurr
// 		ilptrancrt.SurrenderPercentage = RoundFloat(((ilptrancrt.FundAmount / iFundValue) * 100), 2)
// 		ilptrancrt.HistoryCode = iHistoryCode
// 		ilptrancrt.InvNonInvFlag = "AC"
// 		ilptrancrt.AllocationCategory = p0059data.AllocationCategory
// 		ilptrancrt.InvNonInvPercentage = RoundFloat(((ilptrancrt.FundAmount / iTotalFundValue) * 100), 2)
// 		ilptrancrt.AccountCode = p0059data.AccountCode

// 		ilptrancrt.CurrencyRate = 1.00 // ranga
// 		ilptrancrt.MortalityIndicator = ""
// 		//ilptrancrt.SurrenderPercentage = 0
// 		ilptrancrt.Tranno = iTranno
// 		ilptrancrt.Seqno = uint(p0059data.SeqNo)
// 		ilptrancrt.UlProcessFlag = "C"
// 		result = txn.Create(&ilptrancrt)
// 		if result.Error != nil {
// 			txnErr = models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
// 			return
// 		}
// 		//update ilpsummary
// 		var ilpsummupd models.IlpSummary
// 		result = txn.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

// 		if result.RowsAffected != 0 {
// 			ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
// 			txn.Save(&ilpsummupd)
// 		} else if result.Error != nil {
// 			txnErr = models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
// 			return
// 		}
// 	}
// 	return models.TxnError{}
// }

func CreateCommunicationsNew(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, iBenefit uint, iPa uint, iClientWork uint, txn *gorm.DB) models.TxnError {

	var communication models.Communication
	var iP0033Key string
	var iP0034Key string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data
	// txn = initializers.DB.Begin()

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data

	var policy models.Policy
	if iPolicy != 0 {
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL175", DbError: result.Error}
		}
	}
	var payingauth models.PayingAuthority
	if iPa != 0 {
		result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL671", DbError: result.Error}
		}
	}

	iReceiptTranCode := "H0034"
	iReceiptFor := ""
	if iHistoryCode == iReceiptTranCode {
		var receipt models.Receipt
		result := txn.Find(&receipt, "company_id = ? and id = ?", iCompany, iReceipt)
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL014", DbError: result.Error}
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
	errparam := "P0034"
	err1 := GetItemD(int(iCompany), errparam, iP0034Key, iDate, &extradatap0034)
	if err1 != nil {
		iP0034Key = iHistoryCode
		err1 = GetItemD(int(iCompany), errparam, iP0034Key, iDate, &extradatap0034)
		if err1 != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iP0034Key}
		}
	}

	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iP0033Key = p0034data.Letters[i].Templates
			errparam := "P0033"
			err := GetItemD(int(iCompany), errparam, iP0033Key, iDate, &extradatap0033)
			if err != nil {
				return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iP0033Key}
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
					oData, funcErr := GetCompanyDataNew(iCompany, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData, funcErr := GetClientDataNew(iCompany, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData, funcErr := GetAddressDataNew(iCompany, iAddress, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData, funcErr := GetPolicyDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData, funcErr := GetBenefitDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData, funcErr := GetSurBDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData, funcErr := GetMrtaDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData, funcErr := GetReceiptDataNew(iCompany, iReceipt, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData, funcErr := GetSaChangeDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData, funcErr := GetCompAddDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData, funcErr := GetSurrHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData, funcErr := GetDeathDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData, funcErr := GetMatHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData, funcErr := GetSurvBPayNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData, funcErr := GetExpiNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData, funcErr := GetBonusValsNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData, funcErr := GetAgencyNew(iCompany, iAgency, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData, funcErr := GetNomiDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData, funcErr := GetGLDataNew(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData, funcErr := GetIlpSummaryDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData, funcErr := GetIlpAnnsummaryDataNew(iCompany, iPolicy, iHistoryCode, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData, funcErr := GetIlpTranctionDataNew(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData, funcErr := GetPremTaxGLDataNew(iCompany, iPolicy, iFromDate, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData
				case oLetType == "24":
					oData, funcErr := GetIlpFundSwitchDataNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SwitchData"] = oData
				case oLetType == "25":
					oData, funcErr := GetPHistoryDataNew(iCompany, iPolicy, iHistoryCode, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData, funcErr := GetIlpFundDataNew(iCompany, iPolicy, iBenefit, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData, funcErr := GetPPolicyDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData, funcErr := GetPBenefitDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevBenefit"] = oData
				case oLetType == "29":
					oData, funcErr := GetPayingAuthorityDataNew(iCompany, iPa, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PaData"] = oData
				case oLetType == "30":
					oData, funcErr := GetClientWorkDataNew(iCompany, iClientWork, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
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

			results := txn.Create(&communication)

			if results.Error != nil {
				return models.TxnError{
					ErrorCode: "DBERR",
					DbError:   results.Error,
				}
			}

		}
	}
	return models.TxnError{}
}

func ValidateBillTypeNNew(policyenq models.Policy, userco uint, userlan uint, iDate string, iBillType string, iPayingAuthority uint, txn *gorm.DB) (txnErr models.TxnError) {

	var p0055data paramTypes.P0055Data
	var extradatap0055 paramTypes.Extradata = &p0055data

	errparam := "P0055"

	err := GetItemD(int(userco), errparam, iBillType, iDate, &extradatap0055)
	if err != nil {
		txnErr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iBillType}
		return txnErr
	}
	// Validate SSI Bill Type

	if p0055data.PayingAuthority == "N" &&
		iBillType == policyenq.BillingType {
		txnErr = models.TxnError{ErrorCode: "GL637"}
		return txnErr
	}

	if p0055data.PayingAuthority == "Y" &&
		iBillType == policyenq.BillingType &&
		iPayingAuthority == policyenq.PayingAuthority {
		txnErr = models.TxnError{ErrorCode: "GL638"}
		return txnErr
	}

	if p0055data.PayingAuthority == "N" {
		if iPayingAuthority != 0 {
			txnErr = models.TxnError{ErrorCode: "GL700"}
			return txnErr

		}
	}

	if p0055data.PayingAuthority == "Y" {
		if iPayingAuthority == 0 {
			txnErr = models.TxnError{ErrorCode: "GL701"}
			return txnErr

		}
	}

	// validate Paying authority
	funcErr := ValidatePayingAuthorityNNew(userco, userlan, iDate, iPayingAuthority, txn)
	if funcErr.ErrorCode != "" {
		txnErr = funcErr
		return txnErr
	}

	// P0055 Bank Extration Types like cBank,DBank,NEFT,UPI validation are to be added

	return txnErr
}

func ValidatePayingAuthorityNNew(userco uint, userlan uint, iDate string, iPayingAuthority uint, txn *gorm.DB) (txnErr models.TxnError) {

	var payingauth models.PayingAuthority
	result := txn.First(&payingauth, "company_id = ? and id = ?", userco, iPayingAuthority)
	if result.Error != nil {
		txnErr = models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		return txnErr
	}

	if payingauth.PaStatus != "AC" {
		txnErr = models.TxnError{ErrorCode: "GL640"}
		return txnErr
	}

	if payingauth.StartDate > iDate {
		txnErr = models.TxnError{ErrorCode: "GL641"}
		return txnErr
	}

	if payingauth.EndDate < iDate {
		txnErr = models.TxnError{ErrorCode: "GL642"}
		return txnErr
	}

	return txnErr
}

// 2025-10-30 Lakshmi Changes
// Commented out as duplicate given by Lakshmi
// func GetAnnualRateNNew(iCompany uint, iCoverage string, iAge uint, iGender string, iTerm uint, iPremTerm uint, iPremMethod string, iDate string, iMortality string, txn *gorm.DB) (float64, models.TxnError) {

// 	var q0006data paramTypes.Q0006Data
// 	var extradata paramTypes.Extradata = &q0006data
// 	errparam := "Q0006"
// 	err := GetItemD(int(iCompany), "Q0006", iCoverage, iDate, &extradata)
// 	if err != nil {
// 		return 0, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}

// 	}

// 	var q0010data paramTypes.Q0010Data
// 	var extradataq0010 paramTypes.Extradata = &q0010data
// 	var q0010key string
// 	var prem float64
// 	//term := strconv.FormatUint(uint64(iTerm), 10)
// 	//premTerm := strconv.FormatUint(uint64(iPremTerm), 10)

// 	term := fmt.Sprintf("%02d", iTerm)
// 	premTerm := fmt.Sprintf("%02d", iPremTerm)

// 	//fmt.Println("****************", iCompany, iCoverage, iAge, iGender, iTerm, iPremMethod, iDate, iMortality)
// 	if q0006data.PremCalcType == "A" || q0006data.PremCalcType == "U" {
// 		if q0006data.PremiumMethod == "PM002" {
// 			// END1 + Male
// 			q0010key = iCoverage + iGender
// 		}
// 	} else if q0006data.PremCalcType == "P" {
// 		// END1 + Male + Term + Premium Term
// 		if q0006data.PremiumMethod == "PM001" || q0006data.PremiumMethod == "PM003" {
// 			q0010key = iCoverage + iGender + term + premTerm

// 		}

// 	} else if q0006data.PremCalcType == "H" {
// 		// HIP1 + Male
// 		if q0006data.PremiumMethod == "PM005" {
// 			q0010key = iCoverage + iGender
// 		}
// 	}
// 	errparam = "Q0010"
// 	fmt.Println("Premium Key ******", iCoverage, iGender, term, premTerm, q0006data.PremCalcType, q0010key)
// 	err = GetItemD(int(iCompany), errparam, q0010key, iDate, &extradataq0010)
// 	if err != nil {
// 		return 0, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: q0010key}

// 	}
// 	fmt.Println("************", iCompany, iAge, q0010key, iDate)

// 	for i := 0; i < len(q0010data.Rates); i++ {
// 		if q0010data.Rates[i].Age == uint(iAge) {
// 			prem = q0010data.Rates[i].Rate
// 			break
// 		}
// 	}
// 	fmt.Println("************", iCompany, iAge, q0010key, iDate, prem)
// 	return prem, models.TxnError{}
// }

// 2025-10-31 Divya Changes
func SBCreateNNew(iCompany uint, iPolicy uint, iBenefitID uint, iCoverage string, iDate string, iSA float64, iType string, iMethod string, iYear int, iAge int, iTranno uint, txn *gorm.DB) models.TxnError {

	var survb models.SurvB
	fmt.Println("Values", iCompany, iPolicy, iBenefitID, iCoverage, iDate, iSA, iType, iMethod, iYear, iAge, iTranno)
	if iType == "T" {
		var q0012data paramTypes.Q0012Data
		var extradataq0012 paramTypes.Extradata = &q0012data
		// fmt.Println("SB Parameters", iCompany, iType, iMethod, iYear, iCoverage, iDate)
		errparam := "Q0012"
		err := GetItemD(int(iCompany), errparam, iMethod, iDate, &extradataq0012)
		if err != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iMethod}

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
					result := txn.Create(&survb)
					if result.Error != nil {
						return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

					}

				}

			}

		}
	}
	if iType == "A" {
		var q0013data paramTypes.Q0013Data
		var extradataq0013 paramTypes.Extradata = &q0013data
		errparam := "Q0013"
		fmt.Println("SB Parameters", iCompany, iType, iMethod, iAge, iCoverage, iDate)
		err := GetItemD(int(iCompany), errparam, iMethod, iDate, &extradataq0013)
		fmt.Println("SB Parameters", iCompany, iCoverage, iDate)
		if err != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iMethod}

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
					result := txn.Create(&survb)
					if result.Error != nil {
						return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

					}
					continue
				}

			}

		}

	}
	return models.TxnError{}
}

// 2025-10-31 Lakshmi Changes
func CheckPendingILPNew(iCompany uint, iPolicy uint, iLanguage uint, txn *gorm.DB) models.TxnError {

	var ilptransenq models.IlpTransaction

	result := txn.Find(&ilptransenq, "company_id = ? and policy_id = ? and ul_process_flag = ?", iCompany, iPolicy, "P")
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL137", DbError: result.Error}
	}
	return models.TxnError{}
}

func PostUlpDeductionByUnitsNNew(iCompany uint, iPolicy uint, iBenefit uint, iSurrPercentage float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string, txn *gorm.DB) models.TxnError {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL175", DbError: result.Error}
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	errparam := "P0059"
	err := GetItemD(int(iCompany), errparam, iKey, iStartDate, &extradatap0059)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
	}

	var ilpfundenq []models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL784", DbError: result.Error}
	}

	var ilpsumenq []models.IlpSummary

	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL135", DbError: result.Error}
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
		errparam = "P0061"
		err := GetItemD(int(iCompany), errparam, iKey, iStartDate, &extradatap0061)
		if err != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

		//update ilpsummary
		var ilpsummupd models.IlpSummary
		result = txn.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

		if result.RowsAffected != 0 {
			ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
			txn.Save(&ilpsummupd)
		} else {
			return models.TxnError{ErrorCode: "GL135", DbError: result.Error}
		}
	}
	return models.TxnError{}
}

func GetAnnualRateNNew(iCompany uint, iCoverage string, iAge uint, iGender string, iTerm uint, iPremTerm uint, iPremMethod string, iDate string, iMortality string, txn *gorm.DB) (float64, models.TxnError) {

	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	errparam := "Q0006"
	err := GetItemD(int(iCompany), errparam, iCoverage, iDate, &extradata)
	if err != nil {
		return 0, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}

	}

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
	errparam = "Q0010"
	fmt.Println("Premium Key ******", iCoverage, iGender, term, premTerm, q0006data.PremCalcType, q0010key)
	err = GetItemD(int(iCompany), errparam, q0010key, iDate, &extradataq0010)
	if err != nil {
		return 0, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: q0010key}

	}
	fmt.Println("************", iCompany, iAge, q0010key, iDate)

	for i := 0; i < len(q0010data.Rates); i++ {
		if q0010data.Rates[i].Age == uint(iAge) {
			prem = q0010data.Rates[i].Rate
			break
		}
	}
	fmt.Println("************", iCompany, iAge, q0010key, iDate, prem)
	return prem, models.TxnError{}
}

// 2025-11-03 Lakshmi Changes
func GetVersionIdN(iCompany uint, lockedType types.LockedType, lockedTypeKey string, txn *gorm.DB) (string, error) {
	var tranLock models.TransactionLock
	result := txn.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if !recordNotFound && result.Error != nil {
		return "", result.Error
	}

	if recordNotFound {
		fmt.Println("creating the entity as it does not exist:" + lockedTypeKey + ":" + lockedTypeKey)
		versionid, err := CreateTheEntityN(iCompany, lockedType, lockedTypeKey, txn)
		if err != nil {
			return "", errors.New("entity did not exist,error while trying to create :" + err.Error())
		} else {
			return versionid, nil
		}
	}

	if !tranLock.IsValid {
		return "", errors.New("entity is not valid")
	}

	/*if tranLock.IsLocked {
		return "", errors.New("entity is locked")

	} */
	return tranLock.VersionId, nil

}

func GetVersionIdNNew(iCompany uint, lockedType types.LockedType, lockedTypeKey string, txn *gorm.DB) (string, models.TxnError) {
	var tranLock models.TransactionLock
	result := txn.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if !recordNotFound && result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	if recordNotFound {
		fmt.Println("creating the entity as it does not exist:" + lockedTypeKey + ":" + lockedTypeKey)
		versionid, funcErr := CreateTheEntityNNew(iCompany, lockedType, lockedTypeKey, txn)
		if funcErr.ErrorCode != "" {
			return "", funcErr
		} else {
			return versionid, models.TxnError{}
		}
	}

	if !tranLock.IsValid {
		return "", models.TxnError{ErrorCode: "GL783"}
	}

	/*if tranLock.IsLocked {
		return "", errors.New("entity is locked")

	} */
	return tranLock.VersionId, models.TxnError{}

}

func CreateTheEntityN(iCompany uint, lockedType types.LockedType, lockedTypeKey string, txn *gorm.DB) (string, error) {

	var tranLock models.TransactionLock
	result := txn.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if !recordNotFound && result.Error != nil {
		return "", result.Error
	}

	if !recordNotFound {
		return "", errors.New("entity already exists")
	}

	tranLock.CompanyID = iCompany
	tranLock.LockedTypeKey = lockedTypeKey
	tranLock.LockedType = lockedType
	tranLock.IsLocked = false
	tranLock.IsValid = true
	tranLock.CreatedAt = time.Now()
	tranLock.VersionId = uuid.New().String()

	result = txn.Create(&tranLock)

	if result.Error != nil {
		return "", result.Error
	}

	return tranLock.VersionId, nil

}

func CreateTheEntityNNew(iCompany uint, lockedType types.LockedType, lockedTypeKey string, txn *gorm.DB) (string, models.TxnError) {

	var tranLock models.TransactionLock
	result := txn.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if !recordNotFound && result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	if !recordNotFound {
		return "", models.TxnError{ErrorCode: "GL782"}
	}

	tranLock.CompanyID = iCompany
	tranLock.LockedTypeKey = lockedTypeKey
	tranLock.LockedType = lockedType
	tranLock.IsLocked = false
	tranLock.IsValid = true
	tranLock.CreatedAt = time.Now()
	tranLock.VersionId = uuid.New().String()

	result = txn.Create(&tranLock)

	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	return tranLock.VersionId, models.TxnError{}

}

// 2025-11-04 Lakshmi Changes
func CalcSwitchFeeN(iCompany uint, iPolicy uint, iFeeMethod string, iEffectiveDate string, txn *gorm.DB) (oError error, oAmount float64, oPercentage float64) {

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
	results := txn.Find(&policyhistory, "company_id = ? and policy_id = ? and effective_date >=? and effective_date <=? and is_reversed = ? and history_code = ?", iCompany, iPolicy, iPolicyAnniversary, iBusinessDate, 0, "H0139")
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

func CalcSwitchFeeNNew(iCompany uint, iPolicy uint, iFeeMethod string, iEffectiveDate string, txn *gorm.DB) (txnerr models.TxnError, oAmount float64, oPercentage float64) {

	var p0070data paramTypes.P0070Data
	var extradatap0070 paramTypes.Extradata = &p0070data
	iKey := iFeeMethod
	errparam := "P0070"
	err := GetItemD(int(iCompany), errparam, iKey, iEffectiveDate, &extradatap0070)

	if err != nil {
		txnerr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
		return txnerr, 0, 0

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
	results := txn.Find(&policyhistory, "company_id = ? and policy_id = ? and effective_date >=? and effective_date <=? and is_reversed = ? and history_code = ?", iCompany, iPolicy, iPolicyAnniversary, iBusinessDate, 0, "H0139")
	switchcount := 0
	if results.Error != nil {
		switchcount = 0
	} else {
		switchcount = len(policyhistory)

	}
	// Percentage
	if p0070data.SwitchFeeBasis == "P" {
		if uint(switchcount) <= p0070data.FreeSwitches {
			return txnerr, 0, 0
		} else {
			return txnerr, 0, p0070data.FeePercentage
		}
	}
	// Fixed Amount
	if p0070data.SwitchFeeBasis == "F" {
		if uint(switchcount) <= p0070data.FreeSwitches {
			return txnerr, 0, 0
		} else {
			return txnerr, p0070data.FeeAmount, 0
		}
	}

	return txnerr, 0, 0
}

func ValidateBenefitFieldsN(benefitval models.Benefit, userco uint, userlan uint, iHistoryCode string, iCoverage string, txn *gorm.DB) (string error) {

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

func ValidateBenefitFieldsNNew(benefitval models.Benefit, userco uint, userlan uint, iHistoryCode string, iCoverage string, txn *gorm.DB) (string models.TxnError) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	iKey := iHistoryCode + iCoverage
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		iKey = iHistoryCode
		err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
		if err != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
			// longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			// return errors.New(shortCode + " : " + longDesc)
			return models.TxnError{ErrorCode: shortCode}
		}

	}

	return
}

func ValidateMrtaFieldsN(mrtaval models.Mrta, userco uint, userlan uint, iHistoryCode string, iCoverage string, txn *gorm.DB) (string error) {

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

func ValidateMrtaFieldsNNew(mrtaval models.Mrta, userco uint, userlan uint, iHistoryCode string, iCoverage string, txn *gorm.DB) (string models.TxnError) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	iKey := iHistoryCode + iCoverage
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		iKey = iHistoryCode
		err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
		if err != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
			// longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			// return errors.New(shortCode + " : " + longDesc)
			return models.TxnError{ErrorCode: shortCode}
		}

	}

	return
}

func CalcILPSAN(iCompany uint, iPolicy uint, iCoverage string, iHistoryCD string, iDate string, iAge uint, iPrem float64, iFreq string, iSA float64, txn *gorm.DB) (oErr error, oSA float64) {
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

func CalcILPSANNew(iCompany uint, iPolicy uint, iCoverage string, iHistoryCD string, iDate string, iAge uint, iPrem float64, iFreq string, iSA float64, txn *gorm.DB) (txnerr models.TxnError, oSA float64) {
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
	errparam := "P0068"
	err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0068)
	if err != nil {
		txnerr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
		return txnerr, 0

	}
	// err = errors.New("In Valid SA")
	// Multiplier Logic
	for i := 0; i < len(p0068data.RangeArray); i++ {
		if p0068data.RangeArray[i].P0068Basis == "M" {
			if iAge <= p0068data.RangeArray[i].Age {
				oSA = p0068data.RangeArray[i].Factor * prem
				return txnerr, oSA
			}
		}
		if p0068data.RangeArray[i].P0068Basis == "R" {
			if iAge <= p0068data.RangeArray[i].Age {
				if iSA < p0068data.RangeArray[i].FromSA {
					oSA = p0068data.RangeArray[i].FromSA
					return txnerr, oSA
				}
				if iSA > p0068data.RangeArray[i].ToSA {
					oSA = p0068data.RangeArray[i].ToSA
					return txnerr, oSA
				}
				return txnerr, iSA
			}
		}
	}

	return txnerr, oSA
}

func GetMrtaPremN(iCompany uint, iCoverage string, iSA float64, iAge uint, iGender string, iTerm uint, iPremTerm uint, iInterest float64, iInterimPeriod uint, iDate string, txn *gorm.DB) (float64, error) {

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

func GetMrtaPremNNew(iCompany uint, iCoverage string, iSA float64, iAge uint, iGender string, iTerm uint, iPremTerm uint, iInterest float64, iInterimPeriod uint, iDate string, txn *gorm.DB) (float64, models.TxnError) {

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	errparam := "Q0006"
	err := GetItemD(int(iCompany), errparam, iCoverage, iDate, &extradataq0006)
	if err != nil {
		return 0, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}
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
	errparam = "Q0010"
	err = GetItemD(int(iCompany), errparam, q0010key, iDate, &extradataq0010)
	if err != nil {
		return 0, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: q0010key}
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

	return prem, models.TxnError{}

}

func ValidateBenefitDataN(benefitenq models.Benefit, langid uint, iHistoryCode string, txn *gorm.DB) (string error) {
	//businessdate := GetBusinessDate(benefitenq.CompanyID, 0, 0)
	var clientenq models.Client
	result := txn.First(&clientenq, "company_id  = ? and id = ?", benefitenq.CompanyID, benefitenq.ClientID)
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

func ValidateBenefitDataNNew(benefitenq models.Benefit, langid uint, iHistoryCode string, txn *gorm.DB) (string models.TxnError) {
	//businessdate := GetBusinessDate(benefitenq.CompanyID, 0, 0)
	var clientenq models.Client
	result := txn.First(&clientenq, "company_id  = ? and id = ?", benefitenq.CompanyID, benefitenq.ClientID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	errparam := "Q0006"
	err := GetItemD(int(benefitenq.CompanyID), errparam, benefitenq.BCoverage, benefitenq.BStartDate, &extradataq0006)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: benefitenq.BCoverage}
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
		return models.TxnError{ErrorCode: "GL548"} // Age Not Allowed
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
		return models.TxnError{ErrorCode: "GL549"} // Policy Term not Allowed
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
		return models.TxnError{ErrorCode: "GL550"} // Premium Paying Term not Allowed
	}

	//#004 Risk cess Age not Allowed
	benriskcessage := benefitenq.BAge + benefitenq.BTerm
	if benriskcessage < q0006data.MinRiskCessAge ||
		benriskcessage > q0006data.MaxRiskCessAge {
		return models.TxnError{ErrorCode: "GL551"} // Risk cess Age not Allowed
	}

	//#005 Premium cess Age not Allowed
	benpremcessage := benefitenq.BAge + benefitenq.BPTerm
	if benpremcessage < q0006data.MinPremCessAge ||
		benpremcessage > q0006data.MaxPremCessAge {
		return models.TxnError{ErrorCode: "GL552"} // Premium cess Age not Allowed
	}

	//#006 Min Sum Assured not met
	if uint(benefitenq.BSumAssured) < q0006data.MinSA {
		return models.TxnError{ErrorCode: "GL553"} // Min Sum Assured not met
	}

	return models.TxnError{}
}

func ValidatePolicyDataN(policyenq models.Policy, langid uint, iHistoryCode string, txn *gorm.DB) (string error) {
	businessdate := GetBusinessDate(policyenq.CompanyID, 0, 0)
	var clientenq models.Client
	result := txn.First(&clientenq, "company_id  = ? and id = ?", policyenq.CompanyID, policyenq.ClientID)
	if result.Error != nil {
		shortCode := "GL212" // Client Not Found
		longDesc, _ := GetErrorDesc(policyenq.CompanyID, langid, shortCode)
		return errors.New(shortCode + ":" + longDesc)

	}

	var agencyenq models.Agency
	result = txn.First(&agencyenq, "company_id  = ? and id = ?", policyenq.CompanyID, policyenq.AgencyID)
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

func ValidatePolicyDataNNew(policyenq models.Policy, langid uint, iHistoryCode string, txn *gorm.DB) (string models.TxnError) {
	businessdate := GetBusinessDate(policyenq.CompanyID, 0, 0)
	var clientenq models.Client
	result := txn.First(&clientenq, "company_id  = ? and id = ?", policyenq.CompanyID, policyenq.ClientID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}

	var agencyenq models.Agency
	result = txn.First(&agencyenq, "company_id  = ? and id = ?", policyenq.CompanyID, policyenq.AgencyID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	errparam := "Q0005"
	err := GetItemD(int(policyenq.CompanyID), errparam, policyenq.PProduct, policyenq.PRCD, &extradataq0005)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: policyenq.PProduct}
	}

	//#001 RCD is less than PropsalDate
	if q0005data.BackDateAllowed == "N" {
		if policyenq.PRCD < policyenq.ProposalDate {
			return models.TxnError{ErrorCode: "GL539"} // RCD is less than PropsalDate
		}
	}

	//#002 UW Date is less than PropsalDate
	if policyenq.PUWDate != "" {
		if policyenq.PUWDate < policyenq.ProposalDate {
			return models.TxnError{ErrorCode: "GL540"} // UW Date is less than PropsalDate
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
		return models.TxnError{ErrorCode: "GL541"} // Frequency is Inalid
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
		return models.TxnError{ErrorCode: "GL542"} // Contract Curr is Inalid
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
		return models.TxnError{ErrorCode: "GL543"} // Billing Curr is Inalid
	}

	//#006 Backdataing not Allowed
	if policyenq.PRCD < businessdate &&
		q0005data.BackDateAllowed == "N" {
		return models.TxnError{ErrorCode: "GL544"} // Backdataing not Allowed
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
		return models.TxnError{ErrorCode: "GL545"} // Agency Channel Not Allowed
	}

	//#008 Client is Invalid
	if clientenq.ClientStatus != "AC" {
		return models.TxnError{ErrorCode: "GL546"} // Invalid Client
	}

	//#009 Deceased Client
	if !isFieldZero(clientenq.ClientDod) {
		return models.TxnError{ErrorCode: "GL547"} // Deceased Client
	}

	if policyenq.PRCD > businessdate {
		return models.TxnError{ErrorCode: "GL568"} // RCD is greter than businessdate
	}

	return models.TxnError{}
}

// 2025-11-04 Divya Changes
func CreateCommunicationsLNew(iCompany uint, iHistoryCode string, iTranno uint, iDate string, iPolicy uint, iClient uint, iAddress uint, iReceipt uint, iQuotation uint, iAgency uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB, iBenefit uint, iPa uint, iClientWork uint, iAmount1 float64, iAmount2 float64, iNo1 uint, iNo2 uint) models.TxnError {

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
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL014", DbError: result.Error}
		}
		iReceiptFor = receipt.ReceiptFor
	}

	if iPolicy != 0 {
		var policy models.Policy
		result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL003", DbError: result.Error}
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
		if result.RowsAffected == 0 {
			return models.TxnError{ErrorCode: "GL671", DbError: result.Error}
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
	errparam := "P0034"

	err1 := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0034)
	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0034)
		if err1 != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
		}
	}

	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			errparam := "P0033"
			err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0033)
			if err != nil {
				return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
					oData, funcErr := GetCompanyDataNew(iCompany, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["CompanyData"] = oData
				case oLetType == "2":
					oData, funcErr := GetClientDataNew(iCompany, iClient, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ClientData"] = oData
				case oLetType == "3":
					oData, funcErr := GetAddressDataNew(iCompany, iAddress, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["AddressData"] = oData
				case oLetType == "4":
					oData, funcErr := GetPolicyDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyData"] = oData
				case oLetType == "5":
					oData, funcErr := GetBenefitDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BenefitData"] = oData
				case oLetType == "6":
					oData, funcErr := GetSurBDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurBData"] = oData
				case oLetType == "7":
					oData, funcErr := GetMrtaDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MRTAData"] = oData
				case oLetType == "8":
					oData, funcErr := GetReceiptDataNew(iCompany, iReceipt, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ReceiptData"] = oData
				case oLetType == "9":
					oData, funcErr := GetSaChangeDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SAChangeData"] = oData
				case oLetType == "10":
					oData, funcErr := GetCompAddDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ComponantAddData"] = oData
				case oLetType == "11":
					oData, funcErr := GetSurrHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurrData"] = oData
					// oData = GetSurrDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["SurrDData"] = oData
				case oLetType == "12":
					oData, funcErr := GetDeathDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["DeathData"] = oData
				case oLetType == "13":
					oData, funcErr := GetMatHDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["MaturityData"] = oData
					// oData = GetMatDData(iCompany, iPolicy, iClient, iAddress, iReceipt)
					// resultMap["MatDData"] = oData
				case oLetType == "14":
					oData, funcErr := GetSurvBPayNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SurvbPay"] = oData
				case oLetType == "15":
					oData, funcErr := GetExpiNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ExpiryData"] = oData
				case oLetType == "16":
					oData, funcErr := GetBonusValsNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["BonusData"] = oData
				case oLetType == "17":
					oData, funcErr := GetAgencyNew(iCompany, iAgency, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Agency"] = oData
				case oLetType == "18":
					oData, funcErr := GetNomiDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["Nominee"] = oData
				case oLetType == "19":
					oData, funcErr := GetGLDataNew(iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, iGlAccountCode, iGlSign, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData
				case oLetType == "20":
					oData, funcErr := GetIlpSummaryDataNew(iCompany, iPolicy, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlPSummaryData"] = oData
				case oLetType == "21":
					oData, funcErr := GetIlpAnnsummaryDataNew(iCompany, iPolicy, iHistoryCode, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPANNSummaryData"] = oData
				case oLetType == "22":
					oData, funcErr := GetIlpTranctionDataNew(iCompany, iPolicy, iHistoryCode, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ILPTransactionData"] = oData
				case oLetType == "23":
					oData, funcErr := GetPremTaxGLDataNew(iCompany, iPolicy, iFromDate, iToDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["GLData"] = oData

				case oLetType == "24":
					oData, funcErr := GetIlpFundSwitchDataNew(iCompany, iPolicy, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["SwitchData"] = oData

				case oLetType == "25":
					oData, funcErr := GetPHistoryDataNew(iCompany, iPolicy, iHistoryCode, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PolicyHistoryData"] = oData
				case oLetType == "26":
					oData, funcErr := GetIlpFundDataNew(iCompany, iPolicy, iBenefit, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["IlpFundData"] = oData
				case oLetType == "27":
					oData, funcErr := GetPPolicyDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevPolicy"] = oData
				case oLetType == "28":
					oData, funcErr := GetPBenefitDataNew(iCompany, iPolicy, iHistoryCode, iTranno, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					fmt.Println(oData) // Dummy to avoid compilation error
				case oLetType == "29":
					oData, funcErr := GetPayingAuthorityDataNew(iCompany, iPa, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["PrevBenefit"] = oData
				case oLetType == "30":
					oData, funcErr := GetClientWorkDataNew(iCompany, iClientWork, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["ClientWork"] = oData
				case oLetType == "31":
					oData, funcErr := GetLoanDataNew(iCompany, iPolicy, iDate, iAmount1, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["LoanData"] = oData
				case oLetType == "32":
					oData, funcErr := GetAllLoanInterestDataNew(iCompany, iPolicy, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["LoanInterestData"] = oData
				case oLetType == "33":
					oData, funcErr := LoanCapDataNew(iCompany, iPolicy, iDate, iFromDate, iToDate, iAmount1, iAmount2, iNo1, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["LoanCap"] = oData
				case oLetType == "34":
					oData, funcErr := LoanBillDataNew(iCompany, iPolicy, iDate, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["LoanBillData"] = oData
				case oLetType == "35":
					oData, funcErr := LoanBillsInterestDataNew(iCompany, iPolicy, iNo1, iAmount1, txn)
					if funcErr.ErrorCode != "" {
						return funcErr
					}
					resultMap["LoanBillsInterest"] = oData
				case oLetType == "98":
					resultMap["BatchData"] = batchData

				case oLetType == "99":
					resultMap["SignData"] = signData
				default:

				}
			}

			if p0033data.Online == "Y" {
				funcErr := GetReportforOnlineNew(communication, p0033data.TemplateName, txn)
				if funcErr.ErrorCode != "" {
					return funcErr
				}
			}
			if p0033data.SMSAllowed == "Y" {
				funcErr := SendSMSTwilioNew(communication.CompanyID, communication.ClientID, p0033data.TemplateName, communication.EffectiveDate, p0033data.SMSBody, txn)
				if funcErr.ErrorCode != "" {
					return funcErr
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
				return models.TxnError{ErrorCode: "DBERR", DbError: results.Error}

			}

		}
	}
	return models.TxnError{}
}

func GetMaxLoanSeqNoNNew(iCompany uint, iPolicy uint, txn *gorm.DB) (models.TxnError, uint) {
	var result *gorm.DB
	result1 := txn.Table("loans").Where("company_id = ? and policy_id= ?", iCompany, iPolicy).Select("max(loan_seq_number)")

	if result1.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL750", DbError: result.Error}, 0
	}

	var loanseqno uint = 0
	var newloanseqno uint = 0
	err := result1.Row().Scan(&loanseqno)
	fmt.Println("Error ", err)
	newloanseqno = uint(loanseqno) + 1
	return models.TxnError{}, newloanseqno
}

func EmailTriggerforReportNNew(iCompany uint, iReference uint, iClient uint, iEmail string, iEffDate string, itempName string, pdfData []byte, txn *gorm.DB) models.TxnError {

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	errparam := "P0033"
	err := GetItemD(int(iCompany), errparam, itempName, iEffDate, &extradatap0033)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: itempName}

	}

	sender := p0033data.CompanyEmail
	receiver := iEmail
	password := p0033data.SenderPassword
	smtpServer := p0033data.SMTPServer
	smtpPort := p0033data.SMTPPort

	emailBody := p0033data.Body
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", p0033data.Subject)
	m.SetBody("text/plain", emailBody)
	iTime := time.Now().Format("20060102150405")
	iClientnumstr := strconv.Itoa(int(iClient))
	errChan := make(chan models.TxnError, 1)
	m.Attach(itempName+iClientnumstr+iTime+".pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfData)
		errChan <- models.TxnError{
			ErrorCode: "GL939",
			DbError:   err,
		}
		return err

	}))

	// Configure SMTP dialer
	d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
	d.SSL = true      // Enables SSL
	d.TLSConfig = nil // Use default TLS settings

	// Send email asynchronously with proper logging
	sendStart := time.Now()
	go func() {
		if err := d.DialAndSend(m); err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent successfully to %s (CC: %s, BCC: %s) in %v",
				receiver, "", "", time.Since(sendStart))
		}
	}()
	log.Printf("EmailTrigger function executed in %v", time.Since(sendStart))
	return models.TxnError{}
}

func ValidateNomineeNNew(nomineeval models.Nominee, userco uint, userlan uint, iKey string, txn *gorm.DB) (string models.TxnError) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}

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
			return models.TxnError{ErrorCode: "GL894"}
		}
	}
	var clientenq models.Client
	result := txn.First(&clientenq, "company_id  = ? and id = ?", nomineeval.CompanyID, nomineeval.ClientID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}

	if clientenq.ClientStatus != "AC" ||
		clientenq.ClientDod != "" {
		return models.TxnError{ErrorCode: "GL546"}
	}

	var p0045data paramTypes.P0045Data
	var extradatap0045 paramTypes.Extradata = &p0045data
	errparam = "P0045"
	err = GetItemD(int(nomineeval.CompanyID), errparam, nomineeval.NomineeRelationship, "0", &extradatap0045)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: nomineeval.NomineeRelationship}

	}

	var iGender bool = false
	for i := 0; i < len(p0045data.Gender); i++ {
		if clientenq.Gender == p0045data.Gender {
			iGender = true
			break
		}
	}
	if !iGender {
		return models.TxnError{ErrorCode: "GL572"}
	}

	// Owner cannot be Nominee
	var policyenq models.Policy
	result = txn.First(&policyenq, "company_id  = ? and id = ?", nomineeval.CompanyID, nomineeval.PolicyID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}

	if nomineeval.ClientID == policyenq.ClientID {
		return models.TxnError{ErrorCode: "GL589"}
	}

	return
}

func GetAllowedFundsNNew(iCompany uint, iCoverage string, iDate string, txn *gorm.DB) ([]interface{}, models.TxnError) {

	fundlist := make([]interface{}, 0)

	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	errparam := "Q0006"
	err := GetItemD(int(iCompany), errparam, iCoverage, iDate, &extradataq0006)
	if err != nil {
		return fundlist, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}

	}
	if q0006data.FUNDCODE == nil {
		return nil, models.TxnError{ErrorCode: "GL452"}
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	for i := 0; i < len(q0006data.FUNDCODE); i++ {
		errparam = "P0061"
		err = GetItemD(int(iCompany), errparam, q0006data.FUNDCODE[i], iDate, &extradatap0061)
		if err != nil {
			return fundlist, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: q0006data.FUNDCODE[i]}

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
	return fundlist, models.TxnError{}
}

// 2025-11-05 Lakshmi Changes
func GetParamPlanBenefitN(iCompany uint, iBCoverage, iBenefitPlan, iDate string, txn *gorm.DB) (error, []interface{}) {
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

func GetParamPlanBenefitNNew(iCompany uint, iBCoverage, iBenefitPlan, iDate string, txn *gorm.DB) (models.TxnError, []interface{}) {
	resp := make([]interface{}, 0)

	iKey := iBCoverage + iBenefitPlan

	var p0075data paramTypes.P0075Data
	var extradatap0075 paramTypes.Extradata = &p0075data
	errparam := "P0075"
	err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0075)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}, nil

	}

	var p0077data paramTypes.P0077Data
	var extradatap0077 paramTypes.Extradata = &p0077data
	errparam = "P0077"
	err = GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0077)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}, nil

	}

	for _, planBenefit := range p0075data.PlanBenefits {

		for _, planMaxBenefit := range p0077data.PlanMaxBenefits {

			if err != nil {
				return models.TxnError{}, nil
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

	return models.TxnError{}, resp
}

// 2025-11-07 Divya Changes
func ValidatePolicyFieldsN(policyval models.Policy, userco uint, userlan uint, iKey string, txn *gorm.DB) (string error) {

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

func ValidatePolicyFieldsNNew(policyval models.Policy, userco uint, userlan uint, iKey string, txn *gorm.DB) (string models.TxnError) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
			//longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return models.TxnError{ErrorCode: shortCode}
		}

	}

	return
}

func LockTheEntityN(iCompany uint, lockedType types.LockedType, lockedTypeKey string, versionID string, iUserId uint64, txn *gorm.DB) error {

	var tranLock models.TransactionLock
	result := txn.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if recordNotFound {
		return errors.New("entity does not exist")
	}

	if result.Error != nil {
		return result.Error
	}

	if !tranLock.IsValid {
		return errors.New("entity does not exist")
	}

	if tranLock.IsLocked {
		return errors.New("entity is already locked")

	}

	if versionID != tranLock.VersionId {
		return errors.New("entity version mismatch")

	}

	tranLock.IsLocked = true
	tranLock.UpdatedID = iUserId
	tranLock.UpdatedAt = time.Now()

	//result = initializers.DB.Save(&tranLock)
	result = txn.Model(&tranLock).Updates(tranLock)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func LockTheEntityNNew(iCompany uint, lockedType types.LockedType, lockedTypeKey string, versionID string, iUserId uint64, txn *gorm.DB) models.TxnError {

	var tranLock models.TransactionLock
	result := txn.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if recordNotFound {
		return models.TxnError{ErrorCode: "GL683"}
	}

	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}

	if !tranLock.IsValid {
		return models.TxnError{ErrorCode: "GL683"}
	}

	if tranLock.IsLocked {
		return models.TxnError{ErrorCode: "GL684"}

	}

	if versionID != tranLock.VersionId {
		return models.TxnError{ErrorCode: "GL685"}

	}

	tranLock.IsLocked = true
	tranLock.UpdatedID = iUserId
	tranLock.UpdatedAt = time.Now()

	//result = initializers.DB.Save(&tranLock)
	result = txn.Model(&tranLock).Updates(tranLock)

	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}

	return models.TxnError{}

}

func UnLockTheEntityNew(iCompany uint, lockedType types.LockedType, lockedTypeKey string, iUserId uint64, changeVersion bool, txn *gorm.DB) models.TxnError {

	var tranLock models.TransactionLock
	result := txn.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)
	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if recordNotFound {
		return models.TxnError{ErrorCode: "GL878"}
	}

	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}

	if !tranLock.IsValid {
		return models.TxnError{ErrorCode: "GL878"}
	}

	if !tranLock.IsLocked {
		return models.TxnError{ErrorCode: "GL879"}

	}

	dataMap := make(map[string]interface{})

	dataMap["is_locked"] = false
	dataMap["updated_id"] = iUserId
	if changeVersion {
		dataMap["version_id"] = uuid.New().String()
	}

	result = txn.Model(&tranLock).Updates(dataMap)

	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	return models.TxnError{}

}

func ValidateAgencyN(agencyenq models.Agency, userco uint, userlan uint, iDate string, txn *gorm.DB) (string error) {

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

func ValidateAgencyNNew(agencyenq models.Agency, userco uint, userlan uint, iDate string, txn *gorm.DB) (string models.TxnError) {

	if agencyenq.AgencySt != "AC" {
		return models.TxnError{ErrorCode: "GL221"}
	}

	if agencyenq.LicenseStartDate > iDate {
		return models.TxnError{ErrorCode: "GL577"}

	}

	if agencyenq.LicenseEndDate < iDate {
		return models.TxnError{ErrorCode: "GL578"}

	}
	return models.TxnError{}
}

func PostAllocationNNew(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iFrequency string, iStartDate string, iEffDate string, iGender string, iAllocMethod string, iTranno uint, txn *gorm.DB) models.TxnError {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL017", DbError: result.Error}

	}

	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data
	iDate := iStartDate
	iKey := iAllocMethod + iGender
	errparam := "P0060"
	err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0060)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}

	}
	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey = iHistoryCode + iBenefitCode
	errparam = "P0059"
	err = GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0059)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}

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
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL784", DbError: result.Error}
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
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

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
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}
	return models.TxnError{}
}

func TDFCOLADNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string, txn *gorm.DB) (string, models.TxnError) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}
	errparam := "Q0005"
	err := GetItemD(int(iCompany), errparam, policy.PProduct, policy.PRCD, &extradataq0005)
	if err != nil {
		return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: policy.PProduct}

	}

	iNewDate := AddLeadDays(policy.AnnivDate, (-1 * q0005data.BillingLeadDays))

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

	}
	if results.Error != nil {
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.EffectiveDate = iNewDate
		tdfpolicy.Tranno = iTranno
		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

		}

		return "", models.TxnError{}
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

		}

		iNxtAnnDate := GetNextYr(policy.AnnivDate)
		oAnnivDate := AddLeadDays(iNxtAnnDate, (-1 * q0005data.BillingLeadDays))
		var tdfpolicy models.TDFPolicy
		tdfpolicy.CompanyID = iCompany
		tdfpolicy.PolicyID = iPolicy
		tdfpolicy.Seqno = tdfrule.Seqno
		tdfpolicy.TDFType = iFunction
		tdfpolicy.ID = 0
		tdfpolicy.EffectiveDate = oAnnivDate
		tdfpolicy.Tranno = iTranno

		result = txn.Create(&tdfpolicy)
		if result.Error != nil {
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}

		}

		return "", models.TxnError{}
	}
}

func ValidatePolicyBenefitsDataN(policyenq models.Policy, benefitenq []models.Benefit, langid uint, txn *gorm.DB) (string error) {

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

func ValidatePolicyBenefitsDataNNew(policyenq models.Policy, benefitenq []models.Benefit, langid uint, txn *gorm.DB) (string models.TxnError) {

	var q0011data paramTypes.Q0011Data
	var extradataq0011 paramTypes.Extradata = &q0011data
	errparam := "Q0011"
	err := GetItemD(int(policyenq.CompanyID), errparam, policyenq.PProduct, policyenq.PRCD, &extradataq0011)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: policyenq.PProduct}

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

		return models.TxnError{ErrorCode: "GL619"}
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

		return models.TxnError{ErrorCode: "GL621"}
	}

	// Benefits Validation
	for i := 0; i < len(benefitenq); i++ {
		//#001 Policy RCD > Benefit Start Date
		if policyenq.PRCD > benefitenq[i].BStartDate {
			return models.TxnError{ErrorCode: "GL622"}
		}
		if benefitenq[i].BCoverage != basicbenefit.BCoverage {
			for j := 0; j < len(q0011data.Coverages); j++ {
				if benefitenq[i].BCoverage == q0011data.Coverages[j].CoverageName {
					//#002 Benefit Term Exceed Basic Benefit Term
					if q0011data.Coverages[i].TermCanExceed == "N" &&
						benefitenq[i].BTerm > basicbenefit.BTerm {
						return models.TxnError{ErrorCode: "GL623"}
					}
					//#003 Benefit Prem Term Exceed Basic Benefit Prem Term
					if q0011data.Coverages[i].PptCanExceed == "N" &&
						benefitenq[i].BPTerm > basicbenefit.BPTerm {
						return models.TxnError{ErrorCode: "GL624"}
					}
					//#004 Benefit SA Exceed Basic Benefit SA
					if q0011data.Coverages[i].SaCanExceed == "N" &&
						benefitenq[i].BSumAssured > basicbenefit.BSumAssured {
						return models.TxnError{ErrorCode: "GL625"}
					}
				}
			}
		}
	}

	return models.TxnError{}
}

// 2025-11-10 Lakshmi Changes
func ValidateAddressN(addressval models.Address, userco uint, userlan uint, iKey string, txn *gorm.DB) (string error) {

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

func ValidateAddressNNew(addressval models.Address, userco uint, userlan uint, iKey string, txn *gorm.DB) (string models.TxnError) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
			// longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return models.TxnError{ErrorCode: shortCode}
		}

	}

	return
}

// 2025-11-11 Divya Changes
func ValidateClientN(clientval models.Client, userco uint, userlan uint, iKey string, txn *gorm.DB) (string error) {

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

func ValidateFieldsN(iFunction string, iFieldName string, iFieldVal string, iUserId uint64, iFieldType string, txn *gorm.DB) error {
	var fieldvalidators models.FieldValidator
	var getUser models.User
	results := txn.First(&getUser, "id = ?", iUserId)

	if results.Error != nil {
		return errors.New(results.Error.Error())
	} else {
		oLanguageId := getUser.LanguageID
		oCompanyId := getUser.CompanyID
		results := txn.First(&fieldvalidators, "function_name = ? and company_id = ? and language_id =? and field_name = ?", iFunction, oCompanyId, oLanguageId, iFieldName)

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

func ValidateFieldsNNew(iFunction string, iFieldName string, iFieldVal string, iUserId uint64, iFieldType string, txn *gorm.DB) models.TxnError {
	var fieldvalidators models.FieldValidator
	var getUser models.User
	results := txn.First(&getUser, "id = ?", iUserId)

	if results.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: results.Error}

	} else {
		oLanguageId := getUser.LanguageID
		oCompanyId := getUser.CompanyID
		results := txn.First(&fieldvalidators, "function_name = ? and company_id = ? and language_id =? and field_name = ?", iFunction, oCompanyId, oLanguageId, iFieldName)

		if results.Error != nil {
			return models.TxnError{ErrorCode: "DBERR", DbError: results.Error}

		} else {
			if fieldvalidators.ParamName != "" {
				fmt.Println("I am here Ranga")
				//	fmt.Println(iFieldName, iFieldVal, fieldvalidators.ParamName, fieldvalidators.ErrorDescription)
				// fmt.Println("I am inside ", fieldvalidators.ParamName, iFieldVal, iFieldName)
				funcErr := ValidateItemNNew(iUserId, fieldvalidators.ParamName, iFieldVal, iFieldName, fieldvalidators.ErrorDescription, txn)
				if funcErr.ErrorCode != "" {
					return funcErr

				}
			} else if iFieldType == "string" {
				if fieldvalidators.BlankAllowed == "N" && iFieldVal == "" {
					return models.TxnError{ErrorCode: "GL906"}
				}

			} else if iFieldType != "string" {
				fmt.Println("Field Type is ", iFieldType, iFieldVal)

				if fieldvalidators.ZeroAllowed == "N" && iFieldVal == "0" {
					return models.TxnError{ErrorCode: "GL906"}
				}
			}

		}
		return models.TxnError{}
	}

}

func ValidateItemN(iUserId uint64, iName string, iItem string, iFieldName string, iErros string, txn *gorm.DB) error {
	var getUser models.User
	results := txn.First(&getUser, "id = ?", iUserId)
	if results.Error != nil {
		fmt.Println(results.Error)
		return errors.New(results.Error.Error())
	}
	var valdiateparam models.ParamDesc
	oLanguageId := getUser.LanguageID
	oCompanyId := getUser.CompanyID
	results = txn.Where("company_id = ? AND name = ? and item = ? and language_id = ?", oCompanyId, iName, iItem, oLanguageId).Find(&valdiateparam)
	if results.Error != nil || results.RowsAffected == 0 {

		return errors.New(" -" + strconv.FormatUint(uint64(oCompanyId), 10) + "-" + iName + "-" + strconv.FormatUint(uint64(oLanguageId), 10) + "-" + "-" + iFieldName + iErros + " is missing")
		//return errors.New(results.Error.Error())
	}
	return nil
}

func ValidateItemNNew(iUserId uint64, iName string, iItem string, iFieldName string, iErros string, txn *gorm.DB) models.TxnError {
	var getUser models.User
	results := txn.First(&getUser, "id = ?", iUserId)
	if results.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: results.Error}

	}
	var valdiateparam models.ParamDesc
	oLanguageId := getUser.LanguageID
	oCompanyId := getUser.CompanyID
	results = txn.Where("company_id = ? AND name = ? and item = ? and language_id = ?", oCompanyId, iName, iItem, oLanguageId).Find(&valdiateparam)
	if results.Error != nil || results.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
		//return errors.New(" -" + strconv.FormatUint(uint64(oCompanyId), 10) + "-" + iName + "-" + strconv.FormatUint(uint64(oLanguageId), 10) + "-" + "-" + iFieldName + iErros + " is missing")
		//return errors.New(results.Error.Error())
	}
	return models.TxnError{}
}

// 2025-11-12 Lakshmi Changes
func PostTopAllocationNNew(iCompany uint, iPolicy uint, iBenefit uint, iAmount float64, iHistoryCode string, iBenefitCode string, iFrequency string, iStartDate string, iEffDate string, iGender string, iAllocMethod string, iTranno uint, txn *gorm.DB) models.TxnError {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL136", DbError: result.Error}
	}

	var p0060data paramTypes.P0060Data
	var extradatap0060 paramTypes.Extradata = &p0060data
	iDate := iStartDate
	iKey := iAllocMethod + iGender
	errparam := "P0060"
	err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0060)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
	}
	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey = iHistoryCode + iBenefitCode
	errparam = "P0059"
	err = GetItemD(int(iCompany), errparam, iKey, iDate, &extradatap0059)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL784", DbError: result.Error}
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
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	// Delete Newly Cleared Fund Rules which is created for Top-up
	result = txn.Delete(ilpfundenq)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	return models.TxnError{}

}

func GetMaxTranno2New(iCompany uint, iPolicy uint, iMethod string, iEffDate string, iuser uint64, txn *gorm.DB) (models.TxnError, uint) {

	var maxtranno = 0
	var phistories models.PHistory

	result1 := txn.Order("tranno DESC").Find(&phistories, "company_id = ? and policy_id = ?", iCompany, iPolicy)

	if result1.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL919", DbError: result1.Error}, 0
	}
	maxtranno = int(phistories.Tranno)
	return models.TxnError{}, uint(maxtranno)
}

func CreatePHistoryNew(iCompany uint, iPolicy uint, iMethod string, iEffDate string, maxTranno uint, iuser uint64, historyMap map[string]interface{}, txn *gorm.DB) models.TxnError {

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
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	return models.TxnError{}
}

// 2025-11-13 Lakshmi Changes
func ValidatePayerN(payerval models.Payer, userco uint, userlan uint, iKey string, txn *gorm.DB) (string error) {
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
	result := txn.Find(&policy, "company_id = ? and id = ?", userco, iPolicy)
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

func ValidatePayerNNew(payerval models.Payer, userco uint, userlan uint, iKey string, txn *gorm.DB) (string models.TxnError) {
	businessdate := GetBusinessDate(payerval.CompanyID, 0, 0)
	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
			//longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
			return models.TxnError{ErrorCode: shortCode}
		}
	}

	iPolicy := payerval.PolicyID
	var policy models.Policy
	result := txn.Find(&policy, "company_id = ? and id = ?", userco, iPolicy)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "GL175", DbError: result.Error}

	}
	if payerval.FromDate > businessdate {
		return models.TxnError{ErrorCode: "GL616"}
	}
	if payerval.FromDate < policy.PRCD {
		return models.TxnError{ErrorCode: "GL617"}
	}

	if payerval.FromDate > payerval.ToDate {
		return models.TxnError{ErrorCode: "GL901"}
	}

	return
}

// 2025-11-17 Lakshmi Changes
func RevGlMoveN(tranno, userco, ipolicy float64, txn *gorm.DB) error {
	var glmoveenq []models.GlMove
	opol := strconv.Itoa(int(ipolicy))
	results := txn.Where("gl_rldg_acct LIKE ?", "%"+opol+"%").Find(&glmoveenq, "tranno = ? and company_id = ? ", tranno, userco)

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
		err := PostGlMoveN(uint(iCompany), iContractCurry, iEffectiveDate, int(iTranno), iGlAmount, iAccAmount, iAccountCodeID, uint(iGlRdocno), iGlRldgAcct, iSeqnno, iGlSign, iAccountCode, iHistoryCode, iRevInd, iCoverage, txn)
		if err != nil {
			return err
		}
	}

	return nil
}

func RevGlMoveNNew(tranno, userco, ipolicy float64, txn *gorm.DB) models.TxnError {
	var glmoveenq []models.GlMove
	opol := strconv.Itoa(int(ipolicy))
	results := txn.Where("gl_rldg_acct LIKE ?", "%"+opol+"%").Find(&glmoveenq, "tranno = ? and company_id = ? ", tranno, userco)

	if results.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL830", DbError: results.Error}
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
		funcErr := PostGlMoveNNew(uint(iCompany), iContractCurry, iEffectiveDate, int(iTranno), iGlAmount, iAccAmount, iAccountCodeID, uint(iGlRdocno), iGlRldgAcct, iSeqnno, iGlSign, iAccountCode, iHistoryCode, iRevInd, iCoverage, txn)
		if funcErr.ErrorCode != "" {
			return funcErr
		}
	}

	return models.TxnError{}
}

// 2025-11-18 Lakshmi Changes
func CreateReceiptBN(iCompany uint, iPolicy uint, iAmount float64, iCollDate string, iCollCurr string, iCollType string, iRef string, iMethod string, iIFSC string, iBankAc string, txn *gorm.DB) (oreceipt uint, oerror error) {
	iBusinssdate := GetBusinessDate(iCompany, 1, 2)

	var policyenq models.Policy
	var receiptupd models.Receipt
	var result *gorm.DB
	var clientenq models.Client

	result = txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.Error != nil {
		return 0, errors.New(result.Error.Error())
	}

	iClient := policyenq.ClientID

	result = txn.Find(&clientenq, "company_id = ? and Id = ?", iCompany, iClient)

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
	result = txn.Create(&receiptupd)

	// Debit Entry
	glcode := p0027data.GlMovements[0].AccountCode
	var acccode models.AccountCode
	result = txn.First(&acccode, "company_id = ? and account_code = ? ", iCompany, glcode)
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
	result = txn.First(&acccode1, "company_id = ? and account_code = ? ", iCompany, glcode)
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

func CreateReceiptBNNew(iCompany uint, iPolicy uint, iAmount float64, iCollDate string, iCollCurr string, iCollType string, iRef string, iMethod string, iIFSC string, iBankAc string, txn *gorm.DB) (oreceipt uint, txnerr models.TxnError) {
	iBusinssdate := GetBusinessDate(iCompany, 1, 2)

	var policyenq models.Policy
	var receiptupd models.Receipt
	var result *gorm.DB
	var clientenq models.Client

	result = txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)

	if result.RowsAffected == 0 {
		txnerr = models.TxnError{ErrorCode: "GL175", DbError: result.Error}
		return 0, txnerr
	}

	iClient := policyenq.ClientID

	result = txn.Find(&clientenq, "company_id = ? and Id = ?", iCompany, iClient)

	if result.RowsAffected == 0 {
		txnerr = models.TxnError{ErrorCode: "GL212", DbError: result.Error}
		return 0, txnerr
	}

	if clientenq.ClientStatus != "AC" {
		txnerr = models.TxnError{ErrorCode: "GL006"}
		return 0, txnerr
	}

	var p0055data paramTypes.P0055Data
	var extradatap0055 paramTypes.Extradata = &p0055data
	iKey := iCollType
	errparam := "P0055"
	err := GetItemD(int(iCompany), errparam, iKey, iBusinssdate, &extradatap0055)

	if err != nil {
		txnerr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
		return 0, txnerr
	}

	var p0027data paramTypes.P0027Data
	var extradata paramTypes.Extradata = &p0027data
	errparam = "P0027"
	err = GetItemD(int(iCompany), errparam, iMethod, iBusinssdate, &extradata)

	if err != nil {
		txnerr = models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iMethod}
		return 0, txnerr
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
	result = txn.Create(&receiptupd)
	if result.Error != nil {
		txnerr = models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		return 0, txnerr
	}
	// Debit Entry
	glcode := p0027data.GlMovements[0].AccountCode
	var acccode models.AccountCode
	result = txn.First(&acccode, "company_id = ? and account_code = ? ", iCompany, glcode)
	if result.Error != nil {
		txnerr = models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		return 0, txnerr
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

	funcErr := PostGlMoveNNew(iCompany, iCollCurr, iEffectiveDate, int(iTranno), iGlAmount,
		iAccAmount, iAccountCodeID, uint(iGlRdocno), string(iGlRldgAcct), iSequenceno, iGlSign, iAccountCode, iMethod, "", "", txn)

	if funcErr.ErrorCode != "" {
		txnerr = funcErr
		return 0, txnerr
	}

	// Credit Entry
	glcode = p0027data.GlMovements[1].AccountCode
	var acccode1 models.AccountCode
	result = txn.First(&acccode1, "company_id = ? and account_code = ? ", iCompany, glcode)
	if result.Error != nil {
		txnerr = models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		return 0, txnerr
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

	funcErr = PostGlMoveNNew(iCompany, iCollCurr, iEffectiveDate, int(iTranno), iGlAmount,
		iAccAmount, iAccountCodeID, uint(iGlRdocno), iGlRldgAcct, iSequenceno, iGlSign, iAccountCode, iMethod, "", "", txn)

	if funcErr.ErrorCode != "" {
		txnerr = funcErr
		return 0, txnerr

	}
	if policyenq.PolStatus == "IF" {
		iNextDueDate := Date2String(GetNextDue(policyenq.PaidToDate, policyenq.PFreq, ""))
		gstamountneeded := GetTotalGSTAmount(iCompany, iPolicy, policyenq.PaidToDate, iNextDueDate)
		iPolicyDeposit := GetGlBal(iCompany, iPolicy, "PolicyDeposit")
		iStampDuty := CalculateStampDutyByPolicy(iCompany, iPolicy)
		iPayable := policyenq.InstalmentPrem + gstamountneeded + iStampDuty + iPolicyDeposit
		if iPayable <= 0 {
			_, funcErr = TDFCollDNNew(iCompany, iPolicy, "COLLD", 0, policyenq.PaidToDate, txn)
			if funcErr.ErrorCode != "" {
				txnerr = funcErr
				return 0, txnerr
			}
			funcErr = TdfhUpdateNNew(iCompany, iPolicy, txn)
			if funcErr.ErrorCode != "" {
				txnerr = funcErr
				return 0, txnerr
			}
		}
	}

	iAgency := policyenq.AgencyID
	//////
	funcErr = CreateCommunicationsNew(iCompany, iMethod, uint(iTranno), iBusinssdate, iPolicy, receiptupd.ClientID, receiptupd.AddressID, receiptupd.ID, 0, iAgency, "", "", "", "", "", 0, 0, 0, txn)
	if funcErr.ErrorCode != "" {
		txnerr = funcErr
		return 0, txnerr
	}

	return receiptupd.ID, txnerr
}

func ValidateClientNNew(clientval models.Client, userco uint, userlan uint, iKey string, txn *gorm.DB) (string models.TxnError) {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	iClientType := clientval.ClientType
	if iClientType == "I" || iClientType == "C" {
		iKey = iKey + iClientType
	}
	errparam := "P0065"
	err := GetItemD(int(userco), errparam, iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}

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
			// longDesc, _ := GetErrorDesc(userco, userlan, shortCode)

			return models.TxnError{
				ErrorCode: shortCode,
			}
		}

	}

	validemail := isValidEmail(clientval.ClientEmail)
	if !validemail {
		return models.TxnError{ErrorCode: "GL477"}
	}

	_, err = strconv.Atoi(clientval.ClientMobile)
	if err != nil {
		return models.TxnError{ErrorCode: "GL478"}
	}

	ibusinessdate := GetBusinessDate(userco, 0, 0)
	if clientval.ClientDob > ibusinessdate {
		shortCode := ""
		if clientval.ClientType == "C" {
			shortCode = "GL586" // Incorrect Date of Incorporation
		} else {
			shortCode = "GL566" // Incorrect Date of Birth
		}
		// longDesc, _ := GetErrorDesc(userco, userlan, shortCode)
		return models.TxnError{
			ErrorCode: shortCode,
		}
	}

	if clientval.ClientDod != "" {
		if clientval.ClientDod <= clientval.ClientDob {
			shortCode := ""
			if clientval.ClientType == "C" {
				shortCode = "GL587" // Incorrect Date of Termination
			} else {
				shortCode = "GL567" // Date of Birth/Death Incorrect
			}
			// longDesc, _ := GetErrorDesc(userco, userlan, shortCode)

			return models.TxnError{
				ErrorCode: shortCode,
			}
		}
	}

	return
}

func GetReqCommNew(iCompany uint, iPolicy uint, iClient uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var reqcall []models.ReqCall
	var client models.Client
	var address models.Address

	medDetailsArray := make([]string, 0) // Array for medDetails
	reqCodeArray := make([]string, 0)    // Array for ReqCode
	reqIDArray := make([]uint, 0)        // Array for Req.ID
	remiderdateArray := make([]string, 0)

	// txn := initializers.DB.Begin()

	result := txn.Find(&reqcall, "company_id = ? and policy_id = ? and req_status = ?", iCompany, iPolicy, "P")
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL791",
			DbError:   result.Error,
		}
	}
	result = txn.Find(&client, "company_id = ? and id = ?", iCompany, iClient)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	result = txn.Find(&address, "company_id = ? and client_id = ?", iCompany, iClient)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

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

	return resultMap, models.TxnError{}
}

// this function does not required handle txn
func GetGlBalNew(iCompany uint, iPolicy uint, iGlaccount string, txn *gorm.DB) (float64, models.TxnError) {
	var glbal models.GlBal
	result := txn.Find(&glbal, "company_id = ? and gl_rdocno = ? and gl_accountno = ?", iCompany, iPolicy, iGlaccount)
	if result.Error != nil {
		return 0, models.TxnError{}
	}
	return glbal.ContractAmount, models.TxnError{}

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// End of Changes
////////////////////////////////////////////////////////

func PostUlpDeductionByFundAmountNNew(iCompany uint, iPolicy uint, iBenefit uint, iFundCode string, iAmount float64, iHistoryCode string, iBenefitCode string, iStartDate string, iEffDate string, iTranno uint, iallocType string, txn *gorm.DB) models.TxnError {

	var policyenq models.Policy

	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL175", DbError: result.Error}
	}

	var p0061data paramTypes.P0061Data
	var extradatap0061 paramTypes.Extradata = &p0061data

	var p0059data paramTypes.P0059Data
	var extradatap0059 paramTypes.Extradata = &p0059data

	iKey := iHistoryCode + iBenefitCode + iallocType
	errparam := "P0059"
	err := GetItemD(int(iCompany), errparam, iKey, iStartDate, &extradatap0059)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
	}

	var ilpfundenq models.IlpFund

	result = txn.Find(&ilpfundenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL784", DbError: result.Error}
	}

	var ilpsumenq models.IlpSummary

	result = txn.Find(&ilpsumenq, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, iBenefit, iFundCode)
	if result.RowsAffected == 0 {
		return models.TxnError{ErrorCode: "GL135", DbError: result.Error}
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
	errparam = "P0061"
	err = GetItemD(int(iCompany), errparam, iKey, iStartDate, &extradatap0061)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
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
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	//update ilpsummary
	var ilpsummupd models.IlpSummary
	result = txn.Find(&ilpsummupd, "company_id = ? and policy_id = ? and benefit_id = ? and fund_code = ?", iCompany, iPolicy, ilptrancrt.BenefitID, ilptrancrt.FundCode)

	if result.RowsAffected != 0 {
		ilpsummupd.FundUnits = RoundFloat(ilptrancrt.FundUnits+ilpsummupd.FundUnits, 5)
		result = txn.Save(&ilpsummupd)
		if result.Error != nil {
			return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
	} else {
		return models.TxnError{ErrorCode: "GL135", DbError: result.Error}
	}

	return models.TxnError{}
}

/////////////////////////////////////////////////////////

func TDFExpiDSNNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var benefits []models.Benefit
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	result = txn.Find(&benefits, "company_id = ? and policy_id = ? and b_status = ? ", iCompany, iPolicy, "SP")
	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "GL408", DbError: result.Error}
	}
	oDate := ""
	for i := 0; i < len(benefits); i++ {
		if benefits[i].BStatus != "EX" {
			iCoverage := benefits[i].BCoverage
			iDate := benefits[i].BStartDate
			var q0006data paramTypes.Q0006Data
			var extradataq0006 paramTypes.Extradata = &q0006data
			errparam := "Q0006"
			err := GetItemD(int(iCompany), errparam, iCoverage, iDate, &extradataq0006)
			if err != nil {
				return "", models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iCoverage}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		} else {
			result = txn.Delete(&tdfpolicy)

			if result.Error != nil {
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
				return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
			}
			return "", models.TxnError{}
		}
	}
	return "", models.TxnError{}
}

// ////////////////////////////////////////////////////
func GetUserNameN(iCompany uint, iUserId uint, txn *gorm.DB) (oName string, oErr error) {
	var usrenq models.User
	result := txn.Find(&usrenq, "company_id = ? and id = ?", iCompany, iUserId)
	if result.Error != nil {
		return "", result.Error
	}
	return usrenq.Name, nil
}

////////////////////////////////////////////////////

func GetUserNameNNew(iCompany uint, iUserId uint, txn *gorm.DB) (oName string, txnErr models.TxnError) {
	var usrenq models.User
	result := txn.Find(&usrenq, "company_id = ? and id = ?", iCompany, iUserId)
	if result.Error != nil {
		txnErr = models.TxnError{ErrorCode: "GL120", DbError: result.Error}
		return "", txnErr
	}
	return usrenq.Name, txnErr
}
