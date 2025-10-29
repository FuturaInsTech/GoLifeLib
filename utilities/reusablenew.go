package utilities

import (
	"bytes"
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

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
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
		return models.TxnError{ErrorCode: "GL003", DbError: results.Error}
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

	err1 := GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
		if err1 != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: "P0034", ParamItem: iKey}
		}
	}

	seqno := 0
	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			err := GetItemD(int(iCompany), "P0033", iKey, iDate, &extradatap0033)
			if err != nil {
				return models.TxnError{ErrorCode: "PARME", ParamName: "P0033", ParamItem: iKey}

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
					oData := GetPriorPolicyData(iCompany, iPolicy, iPageSize, iOrientation, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "43":
					oData := GetTermAndConditionData(iCompany, iPolicy, iPageSize, iOrientation, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "44":
					oData := GetpremiumCertificateData(iCompany, iPolicy, iPageSize, iOrientation, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "47":
					oData := GetPOLSCDEndowmentData(iCompany, iPolicy, iPageSize, iOrientation, p0033data, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "51":
					oData := PrtReceiptData(iCompany, iReceipt, iPolicy, iPa, p0033data, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "52":
					oData := PrtPolicyBillData(iCompany, iPolicy, iDate, p0033data, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "53":
					oData := PrtPolicyLapseData(iCompany, iPolicy, iDate, p0033data, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "54":
					oData := PrtCollectionData(iCompany, iPolicy, iDate, p0033data, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "55":
					oData := PrtAnniData(iCompany, iPolicy, iDate, p0033data, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "56":
					oData := PrtAnniILPData(iCompany, iPolicy, iDate, p0033data, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "58":
					oData := PrtExpiData(iCompany, iPolicy, iDate, p0033data, iTranno, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "61":
					oData := PrtPremstData(iCompany, iPolicy, iBenefit, iDate, p0033data, iTranno, iAgency, iFromDate, iToDate, iHistoryCode, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "65":
					oData := PrtFreqChangeData(iCompany, iPolicy, iDate, p0033data, iAgency, iHistoryCode, iTranno, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "66":
					oData := PrtSachangeData(iCompany, iPolicy, iDate, p0033data, iAgency, iTranno, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "67":
					oData := PrtCompaddData(iCompany, iPolicy, iDate, p0033data, iAgency, iTranno, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "68":
					oData := PrtSurrData(iCompany, iPolicy, iDate, p0033data, iAgency, iTranno, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "69":
					oData := PrtMatyData(iCompany, iPolicy, iDate, p0033data, iAgency, txn)
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

func EmailTriggerMNew(icommunication models.Communication, pdfData []byte, txn *gorm.DB) models.TxnError {
	var client models.Client
	result := txn.First(&client, "id = ?", icommunication.ClientID)
	if result.Error != nil {
		return models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	if client.ClientEmail == "" {
		return models.TxnError{ErrorCode: "GL475"}
	}

	iTemplate := icommunication.TemplateName
	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	err := GetItemD(int(icommunication.CompanyID), "P0033", iTemplate, icommunication.EffectiveDate, &extradatap0033)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: "P0033", ParamItem: iTemplate}
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
	go func() {
		m := gomail.NewMessage()
		m.SetHeader("From", sender)
		m.SetHeader("To", receiver)
		m.SetHeader("Subject", p0033data.Subject)
		m.SetBody("text/plain", emailBody)

		// Attach PDF file
		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(pdfData)
			return err
		}))

		d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
		d.SSL = true

		sendStart := time.Now()
		if err := d.DialAndSend(m); err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent successfully to %s in %v", receiver, time.Since(sendStart))
		}
	}()

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
		return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
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
		return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
	}
}

// 2025-10-16 Divya Changes
func ValidateBankNew(bankval models.Bank, userco uint, userlan uint, iKey string) models.TxnError {
	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	// Fetch validation rules
	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: "P0065", ParamItem: iKey}
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

////////////////////////////////////////////////////////////////

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
			return models.TxnError{ErrorCode: "GL706"}
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
	err := GetItemD(int(iCompany), "P0033", itempName, iEffDate, &extradatap0033)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: "P0033", ParamItem: itempName}
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
			return models.TxnError{ErrorCode: "GL037", DbError: result.Error}
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

	err1 := GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
	if err1 != nil {
		iKey = iTransaction
		err1 = GetItemD(int(iCompany), "P0034", iKey, iDate, &extradatap0034)
		if err1 != nil {
			return models.TxnError{ErrorCode: "PARME", ParamName: "P0034", ParamItem: iKey}
		}
	}

	for i := 0; i < len(p0034data.Letters); i++ {
		if p0034data.Letters[i].Templates != "" {
			iKey = p0034data.Letters[i].Templates
			err := GetItemD(int(iCompany), "P0033", iKey, iDate, &extradatap0033)
			if err != nil {
				return models.TxnError{ErrorCode: "PARME", ParamName: "P0033", ParamItem: iKey}
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
				case oLetType == "45":
					oData := ColaCancelData(iCompany, iPolicy, iHistoryCode, txn)
					for key, value := range oData {
						resultMap[key] = value
					}
				case oLetType == "46":
					oData := AplCancelData(iCompany, iPolicy, iHistoryCode, txn)
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

	result = initializers.DB.First(&permission, "company_id = ? and method = ?", iCompany, iMethod)
	if result.Error != nil {
		return iMethod, 0, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	iHistoryCode := permission.TransactionID
	var transaction models.Transaction
	result = initializers.DB.Find(&transaction, "ID = ?", iHistoryCode)
	if result.RowsAffected == 0 {
		return iMethod, 0, models.TxnError{ErrorCode: "GL058", DbError: result.Error}
	}
	iHistoryCD := transaction.TranCode
	var phistory models.PHistory
	var maxtranno float64 = 0

	fmt.Println(iCompany, iPolicy, iHistoryCD, iEffDate)

	result1 := initializers.DB.Table("p_histories").Where("company_id = ? and policy_id= ?", iCompany, iPolicy).Select("max(tranno)")

	if result1.Error != nil {
		fmt.Println(models.TxnError{ErrorCode: "GL058", DbError: result1.Error})

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

// 2025-10-21 Lakshmi Changes
func ValidateClientWorkNNew(clientwork models.ClientWork, userco uint, userlan uint, iDate string, iKey string, txn *gorm.DB) models.TxnError {

	var p0065data paramTypes.P0065Data
	var extradatap0065 paramTypes.Extradata = &p0065data

	err := GetItemD(int(userco), "P0065", iKey, "0", &extradatap0065)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: "P0065", ParamItem: iKey}
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
	result1 := initializers.DB.Find(&client, "company_id = ? and id = ?", userco, clientid)
	if result1.RowsAffected == 0 {
		return models.TxnError{
			ErrorCode: "GL050",
			DbError:   result1.Error,
		}
	}

	if client.ClientStatus != "AC" {
		return models.TxnError{ErrorCode: "GL221", DbError: result1.Error}
	}
	var employer models.Client
	employerid := clientwork.EmployerID
	result2 := initializers.DB.Find(&employer, "company_id = ? and id = ?", userco, employerid)
	if result2.RowsAffected == 0 {
		return models.TxnError{
			ErrorCode: "GL050",
			DbError:   result2.Error,
		}
	}

	if employer.ClientStatus != "AC" {
		return models.TxnError{ErrorCode: "GL221", DbError: result2.Error}
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
		txnerr = models.TxnError{ErrorCode: "GL058", DbError: result.Error}
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

	err := GetItemD(int(iCompany), "P0055", iTypeofPayment, iDate, &extradatap0055)
	if err != nil {
		txnerr = models.TxnError{ErrorCode: "PARME", ParamName: "P0055", ParamItem: iTypeofPayment}
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
	err := GetItemD(int(iCompany), "P0033", itempName, iEffDate, &extradatap0033)
	if err != nil {
		return models.TxnError{ErrorCode: "PARME", ParamName: "P0033", ParamItem: itempName}
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

	m.Attach(itempName+iClientnumstr+iTime+".pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfData)
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
	err := GetItemD(int(icommunication.CompanyID), "P0033", iTemplate, icommunication.EffectiveDate, &extradatap0033)
	errparam := "P0033"
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
	go func() {
		m := gomail.NewMessage()
		m.SetHeader("From", sender)
		m.SetHeader("To", receiver)
		m.SetHeader("Subject", p0033data.Subject)
		m.SetBody("text/plain", emailBody)

		// Attach PDF file
		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(pdfData)
			return err
		}))

		d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
		d.SSL = true

		sendStart := time.Now()
		if err := d.DialAndSend(m); err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent successfully to %s in %v", receiver, time.Since(sendStart))
		}
	}()

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
	txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_status = ? ", iCompany, iPolicy, "AC")
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
			txn.Delete(&tdfpolicy)
			return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: results.Error}
		}
	}
	return "", models.TxnError{}
}

func TDFLoanCapNew(iCompany uint, iPolicy uint, iFunction string, iTranno uint, txn *gorm.DB) (string, models.TxnError) {
	var loanenq []models.Loan
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule
	txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_status = ? ", iCompany, iPolicy, "AC")
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
			txn.Delete(&tdfpolicy)
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
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
		return "", models.TxnError{ErrorCode: "GL037", DbError: result.Error}
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

	result = initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)

	if result.Error != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	results := initializers.DB.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
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

			result = initializers.DB.Create(&tdfpolicy)
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
			return "", models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}

	} else {
		result = initializers.DB.Delete(&tdfpolicy)
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
		result = initializers.DB.Create(&tdfpolicy)
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
	initializers.DB.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	result := initializers.DB.Find(&extraenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return "", models.TxnError{ErrorCode: "GL058", DbError: result.Error}
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// End of Changes
