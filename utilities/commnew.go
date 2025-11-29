package utilities

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"
	"gorm.io/gorm"
)

func GetCompanyDataNew(iCompany uint, iDate string, txn *gorm.DB) ([]interface{}, models.TxnError) {
	companyarray := make([]interface{}, 0)
	var company models.Company
	result := txn.Find(&company, "id = ?", iCompany)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}

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
	return companyarray, models.TxnError{}
}

func GetClientDataNew(iCompany uint, iClient uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	clientarray := make([]interface{}, 0)
	var client models.Client

	result := txn.Find(&client, "company_id = ? and id = ?", iCompany, iClient)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
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
	return clientarray, models.TxnError{}
}

func GetAddressDataNew(iCompany uint, iAddress uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	addressarray := make([]interface{}, 0)
	var address models.Address

	result := txn.Find(&address, "company_id = ? and id = ?", iCompany, iAddress)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
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
	return addressarray, models.TxnError{}
}

func GetPolicyDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	policyarray := make([]interface{}, 0)
	var policy models.Policy
	result := txn.Find(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL003",
			DbError:   result.Error,
		}
	}
	_, oStatus, _ := GetParamDesc(policy.CompanyID, "P0024", policy.PolStatus, 1)
	_, oFreq, _ := GetParamDesc(policy.CompanyID, "Q0009", policy.PFreq, 1)
	_, oProduct, _ := GetParamDesc(policy.CompanyID, "Q0005", policy.PProduct, 1)
	_, oBillCurr, _ := GetParamDesc(policy.CompanyID, "P0023", policy.PBillCurr, 1)
	_, oContCurr, _ := GetParamDesc(policy.CompanyID, "P0023", policy.PContractCurr, 1)
	_, oBillingType, _ := GetParamDesc(policy.CompanyID, "P0055", policy.BillingType, 1)

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	errparam := "Q0005"
	err := GetItemD(int(iCompany), errparam, policy.PProduct, policy.PRCD, &extradataq0005)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: policy.PProduct}
	}
	gracedate := AddLeadDays(policy.PaidToDate, q0005data.LapsedDays)
	premduedates := GetPremDueDates(policy.PRCD, policy.PFreq)
	iAnnivDate := Date2String(GetNextDue(policy.AnnivDate, "Y", "R"))

	var benefitenq []models.Benefit

	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
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

	return policyarray, models.TxnError{}
}
func GetBenefitDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var policyenq models.Policy
	var benefit []models.Benefit
	var clientenq models.Client
	var addressenq models.Address
	result := txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL003",
			DbError:   result.Error,
		}
	}
	paidToDate := policyenq.PaidToDate
	nextDueDate := policyenq.NxtBTDate
	result = txn.Find(&benefit, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	benefitarray := make([]interface{}, 0)

	for k := 0; k < len(benefit); k++ {
		iCompany := benefit[k].CompanyID
		_, oGender, _ := GetParamDesc(iCompany, "P0001", benefit[k].BGender, 1)
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", benefit[k].BCoverage, 1)
		_, oStatus, _ := GetParamDesc(iCompany, "P0024", benefit[k].BStatus, 1)

		clientname := GetName(iCompany, benefit[k].ClientID)
		result = txn.Find(&clientenq, "company_id = ? and id = ?", iCompany, benefit[k].ClientID)
		if result.RowsAffected == 0 {
			return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
		}
		result = txn.Find(&addressenq, "company_id = ? and client_id = ?", iCompany, clientenq.ID)
		if result.RowsAffected == 0 {
			return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
		}
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
	return benefitarray, models.TxnError{}
}

func GetSurBDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var survb []models.SurvB
	result := txn.Find(&survb, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL941", DbError: result.Error}
	}
	var benefitenq models.Benefit
	result = txn.Find(&benefitenq, "company_id = ? and policy_id =? and id = ?", iCompany, iPolicy, survb[0].BenefitID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	basis := ""
	var q0006data paramTypes.Q0006Data
	var extradataq0006 paramTypes.Extradata = &q0006data
	errparam := "Q0006"
	err := GetItemD(int(iCompany), errparam, benefitenq.BCoverage, benefitenq.BStartDate, &extradataq0006)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: benefitenq.BCoverage}
	}
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
	return survbarray, models.TxnError{}
}

func GetMrtaDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var mrtaenq []models.Mrta
	result := txn.Find(&mrtaenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL848",
			DbError:   result.Error,
		}
	}
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
	return mrtaarray, models.TxnError{}
}

func GetReceiptDataNew(iCompany uint, iReceipt uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var receiptenq models.Receipt
	result := txn.Find(&receiptenq, "company_id = ? and id = ?", iCompany, iReceipt)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL014",
			DbError:   result.Error,
		}
	}
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

	return receiptarray, models.TxnError{}
}

func GetSaChangeDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var sachangeenq []models.SaChange
	result := txn.Find(&sachangeenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL426",
			DbError:   result.Error,
		}
	}
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
	return sachangearray, models.TxnError{}
}
func GetCompAddDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var addcomp []models.Addcomponent
	result := txn.Find(&addcomp, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL393",
			DbError:   result.Error,
		}
	}
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
	return addcomparray, models.TxnError{}
}
func GetSurrHDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) (interface{}, models.TxnError) {
	var surrhenq models.SurrH

	result := txn.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL942",
			DbError:   result.Error,
		}
	}
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

	result = txn.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL943",
			DbError:   result.Error,
		}
	}
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

	return surrcombinedvalue, models.TxnError{}

}

func GetDeathDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var surrhenq models.SurrH
	var surrdenq []models.SurrD
	result := txn.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL942",
			DbError:   result.Error,
		}
	}
	result = txn.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL943",
			DbError:   result.Error,
		}
	}
	surrarray := make([]interface{}, 0)

	return surrarray, models.TxnError{}
}
func GetMatHDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) (interface{}, models.TxnError) {
	var mathenq models.MaturityH

	result := txn.Find(&mathenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL822",
			DbError:   result.Error,
		}
	}
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

	result = txn.Find(&matdenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL823",
			DbError:   result.Error,
		}
	}
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
	return matcombineddata, models.TxnError{}

}

func GetSurvBPayNew(iCompany uint, iPolicy uint, iTranno uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var survbenq models.SurvB
	result := txn.Find(&survbenq, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL941", DbError: result.Error}
	}
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
	return survbparray, models.TxnError{}
}

func GetExpiNew(iCompany uint, iPolicy uint, iTranno uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var benefit []models.Benefit
	result := txn.Find(&benefit, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
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
	return expiryarray, models.TxnError{}
}

func GetBonusValsNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {

	bonusarray := make([]interface{}, 0)

	oPolicyDeposit, funcErr1 := GetGlBalNew(iCompany, uint(iPolicy), "PolicyDeposit", txn)
	if funcErr1.ErrorCode != "" {
		return nil, funcErr1
	}
	oRevBonus, funcErr2 := GetGlBalNew(iCompany, uint(iPolicy), "ReversionaryBonus", txn)
	if funcErr2.ErrorCode != "" {
		return nil, funcErr2
	}
	oTermBonus, funcErr3 := GetGlBalNew(iCompany, uint(iPolicy), "TerminalBonus", txn)
	if funcErr3.ErrorCode != "" {
		return nil, funcErr3
	}
	oIntBonus, funcErr4 := GetGlBalNew(iCompany, uint(iPolicy), "InterimBonus", txn)
	if funcErr4.ErrorCode != "" {
		return nil, funcErr4
	}
	oAccumDiv, funcErr5 := GetGlBalNew(iCompany, uint(iPolicy), "AccumDividend", txn)
	if funcErr5.ErrorCode != "" {
		return nil, funcErr5
	}
	oAccumDivInt, funcErr6 := GetGlBalNew(iCompany, uint(iPolicy), "AccumDivInt", txn)
	if funcErr6.ErrorCode != "" {
		return nil, funcErr6
	}
	oAddBonus, funcErr7 := GetGlBalNew(iCompany, uint(iPolicy), "AdditionalBonus", txn)
	if funcErr7.ErrorCode != "" {
		return nil, funcErr7
	}
	oLoyalBonus, funcErr8 := GetGlBalNew(iCompany, uint(iPolicy), "LoyaltyBonus", txn)
	if funcErr8.ErrorCode != "" {
		return nil, funcErr8
	}
	oAplAmt, funcErr9 := GetGlBalNew(iCompany, uint(iPolicy), "AplAmount", txn)
	if funcErr9.ErrorCode != "" {
		return nil, funcErr9
	}
	oPolLoan, funcErr10 := GetGlBalNew(iCompany, uint(iPolicy), "PolicyLoan", txn)
	if funcErr10.ErrorCode != "" {
		return nil, funcErr10
	}
	oCashDep, funcErr11 := GetGlBalNew(iCompany, uint(iPolicy), "CashDeposit", txn)
	if funcErr11.ErrorCode != "" {
		return nil, funcErr11
	}

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

	return bonusarray, models.TxnError{}
}

func GetAgencyNew(iCompany uint, iAgency uint, txn *gorm.DB) ([]interface{}, models.TxnError) {

	agencyarray := make([]interface{}, 0)
	var agencyenq models.Agency
	var clientenq models.Client
	result := txn.Find(&agencyenq, "company_id  = ? and id = ?", iCompany, iAgency)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	result = txn.Find(&clientenq, "company_id = ? and id = ?", iCompany, agencyenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	oAgentName := clientenq.ClientLongName + " " + clientenq.ClientShortName + " " + clientenq.ClientSurName

	var addressenq models.Address
	result = txn.Find(&addressenq, "company_id = ? and client_id = ?", iCompany, clientenq.ID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
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

	return agencyarray, models.TxnError{}
}

func GetNomiDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) ([]interface{}, models.TxnError) {

	var nomenq []models.Nominee

	result := txn.Find(&nomenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL256", DbError: result.Error}
	}
	nomarray := make([]interface{}, 0)
	var clientenq models.Client
	var policyenq models.Policy
	result = txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL003",
			DbError:   result.Error,
		}
	}
	for k := 0; k < len(nomenq); k++ {
		result = txn.Find(&clientenq, "company_id = ? and id = ?", iCompany, nomenq[k].ClientID)
		if result.RowsAffected == 0 {
			return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
		}
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

	return nomarray, models.TxnError{}

}

func GetGLDataNew(iCompany uint, iPolicy uint, iFromDate string, iToDate string, iGlHistoryCode string, iGlAccountCode string, iGlSign string, txn *gorm.DB) (interface{}, models.TxnError) {
	var benefitenq []models.Benefit

	var covrcodes []string
	var covrnames []string

	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
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
		result = txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date <=?", iCompany, iPolicy, iFromDate, iToDate).Order("effective_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL830",
				DbError:   result.Error,
			}
		}
	} else if iGlHistoryCode != "" && iGlAccountCode == "" && iGlSign == "" {
		result = txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode).Order("history_code, effective_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL830",
				DbError:   result.Error,
			}
		}
	} else if iGlHistoryCode != "" && iGlAccountCode != "" && iGlSign == "" {
		result = txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ? and account_code like ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, "%"+iGlAccountCode+"%").Order("history_code, account_code, effective_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL830",
				DbError:   result.Error,
			}
		}
	} else if iGlHistoryCode != "" && iGlAccountCode != "" && iGlSign != "" {
		result = txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and history_code = ? and account_code like ? and gl_sign = ?", iCompany, iPolicy, iFromDate, iToDate, iGlHistoryCode, "%"+iGlAccountCode+"%", iGlSign).Order("history_code, account_code, gl_sign, effective_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL830",
				DbError:   result.Error,
			}
		}
	} else if iGlHistoryCode == "" && iGlAccountCode != "" && iGlSign != "" {
		result = txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and account_code like ? and gl_sign = ?", iCompany, iPolicy, iFromDate, iToDate, "%"+iGlAccountCode+"%", iGlSign).Order("account_code, gl_sign, effective_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL830",
				DbError:   result.Error,
			}
		}
	} else if iGlHistoryCode == "" && iGlAccountCode != "" && iGlSign == "" {
		result = txn.Find(&glmoves, "company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date<=? and account_code like ?", iCompany, iPolicy, iFromDate, iToDate, "%"+iGlAccountCode+"%").Order("account_code, effective_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL830",
				DbError:   result.Error,
			}
		}
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
	return glcombineddata, models.TxnError{}

}
func GetIlpSummaryDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) (interface{}, models.TxnError) {
	var ilpsummary []models.IlpSummary
	result := txn.Find(&ilpsummary, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL135",
			DbError:   result.Error,
		}
	}
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

	return ilpsumcombinedvalue, models.TxnError{}
}
func GetIlpAnnsummaryDataNew(iCompany uint, iPolicy uint, iHistoryCode string, txn *gorm.DB) (interface{}, models.TxnError) {
	ilpannsumprevarray := make([]interface{}, 0)
	ilpannsumcurrarray := make([]interface{}, 0)
	var policyenq models.Policy
	result := txn.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	iAnnivDate := Date2String(GetNextDue(policyenq.AnnivDate, "Y", "R"))
	iPrevAnnivDate := Date2String(GetNextDue(iAnnivDate, "Y", "R"))

	var ilpannsumprev []models.IlpAnnSummary
	result = txn.Find(&ilpannsumprev, "company_id = ? and policy_id = ? and effective_date = ?", iCompany, iPolicy, iPrevAnnivDate)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL944",
			DbError:   result.Error,
		}
	}
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
	result = txn.Find(&ilpannsumcurr, "company_id = ? and policy_id = ? and effective_date = ?", iCompany, iPolicy, iAnnivDate)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL944",
			DbError:   result.Error,
		}
	}
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
	return ilpannsumdata, models.TxnError{}
}
func GetIlpTranctionDataNew(iCompany uint, iPolicy uint, iHistoryCode string, iDate string, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var policyenq models.Policy
	result := txn.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	iAnnivDate := Date2String(GetNextDue(policyenq.AnnivDate, "Y", "R"))
	iPrevAnnivDate := Date2String(GetNextDue(iAnnivDate, "Y", "R"))
	var ilptranction []models.IlpTransaction
	if iHistoryCode == "B0103" {
		result = txn.Find(&ilptranction, "company_id = ? and policy_id = ? and ul_process_flag = ? and inv_non_inv_flag != ? and transaction_date >= ? and transaction_date < ?", iCompany, iPolicy, "C", "NI", iPrevAnnivDate, iAnnivDate).Order("fund_code, transaction_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL137",
				DbError:   result.Error,
			}
		}
	} else if iHistoryCode == "B0115" {
		result = txn.Find(&ilptranction, "company_id = ? and policy_id = ? and ul_process_flag = ? and inv_non_inv_flag != ? and transaction_date >= ? and transaction_date <= ?", iCompany, iPolicy, "C", "NI", iAnnivDate, iDate).Order("fund_code, transaction_date , tranno")
		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL137",
				DbError:   result.Error,
			}
		}
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
	return ilptranctionarray, models.TxnError{}
}
func GetPremTaxGLDataNew(iCompany uint, iPolicy uint, iFromDate string, iToDate string, txn *gorm.DB) (interface{}, models.TxnError) {
	var benefitenq []models.Benefit
	var codesql string = ""
	var covrcodes []string
	var covrnames []string

	var acodearray []string

	var p0067data paramTypes.P0067Data
	var extradatap0067 paramTypes.Extradata = &p0067data

	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	for i := 0; i < len(benefitenq); i++ {
		covrcode := benefitenq[i].BCoverage
		_, covrname, err := GetParamDesc(iCompany, "Q0006", covrcode, 1)
		if err != nil {
			continue
		}
		covrcodes = append(covrcodes, covrcode)
		covrnames = append(covrnames, covrname)
		errparam := "P0067"
		err = GetItemD(int(iCompany), errparam, benefitenq[i].BCoverage, iFromDate, &extradatap0067)
		if err != nil {
			return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: benefitenq[i].BCoverage}
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
	result = txn.Find(&glmoves, "("+codesql+") and company_id = ? and gl_rdocno = ? and effective_date >=? and effective_date <=? ", iCompany, iPolicy, iFromDate, iToDate).Order("account_code, gl_sign, effective_date , tranno")
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL830",
			DbError:   result.Error,
		}
	}
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
			errparam := "P0067"
			err := GetItemD(int(iCompany), errparam, glcoveragecode, iFromDate, &extradatap0067)
			if err != nil {
				return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: glcoveragecode}
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
		errparam := "P0067"
		err := GetItemD(int(iCompany), errparam, glcoveragecode, iFromDate, &extradatap0067)
		if err != nil {
			return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: glcoveragecode}
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
	errparam := "P0067"
	err := GetItemD(int(iCompany), errparam, glcoveragecode, iFromDate, &extradatap0067)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: glcoveragecode}
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
	return glcombineddata, models.TxnError{}

}
func GetIlpFundSwitchDataNew(iCompany uint, iPolicy uint, iTranno uint, txn *gorm.DB) (interface{}, models.TxnError) {
	ilpswitchfundarray := make([]interface{}, 0)
	ilpfundarray := make([]interface{}, 0)
	var policyenq models.Policy
	result := txn.First(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	var ilpswitchheader []models.IlpSwitchHeader

	result = txn.Where("company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicy, iTranno).Order("tranno").Find(&ilpswitchheader)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "GL945", DbError: result.Error}
	}
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
	result = txn.Where(" policy_id = ? and tranno = ?", iPolicy, iTranno).Order("fund_code").Find(&ilpswitchfund)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL946", DbError: result.Error}
	}
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
	return IlpFundSwitchData, models.TxnError{}

}

func GetPHistoryDataNew(iCompany uint, iPolicy uint, iHistoryCode string, iDate string, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var policyhistory []models.PHistory
	result := txn.Find(&policyhistory, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL919",
			DbError:   result.Error,
		}
	}
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
	return policyhistoryarray, models.TxnError{}
}
func GetIlpFundDataNew(iCompany uint, iPolicy uint, iBenefit uint, iDate string, txn *gorm.DB) (interface{}, models.TxnError) {
	var ilpfund []models.IlpFund
	result := txn.Find(&ilpfund, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, iBenefit)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL784",
			DbError:   result.Error,
		}
	}
	var ibenfit models.Benefit
	result = txn.Find(&ibenfit, "company_id = ? and policy_id = ? and id = ?", iCompany, iPolicy, iBenefit)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	ilpfundtarray := make([]interface{}, 0)

	for k := 0; k < len(ilpfund); k++ {
		var p0061data paramTypes.P0061Data
		var extradatap0061 paramTypes.Extradata = &p0061data
		errparam := "P0061"
		err := GetItemD(int(iCompany), errparam, ilpfund[k].FundCode, iDate, &extradatap0061)

		if err != nil {
			return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: ilpfund[k].FundCode}

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

	return ilpfundtarray, models.TxnError{}
}
func GetPPolicyDataNew(iCompany uint, iPolicy uint, iHistoryCode string, iTranno uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	ppolicyarray := make([]interface{}, 0)
	var phistory models.PHistory
	result := txn.Find(&phistory, "company_id = ? and policy_id = ? and history_code = ?  and tranno =  ?", iCompany, iPolicy, iHistoryCode, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL919",
			DbError:   result.Error,
		}
	}
	previousPolicy := phistory.PrevData["Policy"]
	ppolicyarray = append(ppolicyarray, previousPolicy)
	return ppolicyarray, models.TxnError{}

}
func GetPBenefitDataNew(iCompany uint, iPolicy uint, iHistoryCode string, iTranno uint, txn *gorm.DB) (interface{}, models.TxnError) {
	var phistory models.PHistory
	result := txn.Find(&phistory, "company_id = ? and policy_id = ? and history_code = ?  and tranno =  ?", iCompany, iPolicy, iHistoryCode, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL919",
			DbError:   result.Error,
		}
	}
	previousBenefit := phistory.PrevData["Benefits"]
	return previousBenefit, models.TxnError{}

}
func GetPayingAuthorityDataNew(iCompany uint, iPa uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	payingautharray := make([]interface{}, 0)
	var payingauth models.PayingAuthority
	result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL626",
			DbError:   result.Error,
		}
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

	return payingautharray, models.TxnError{}
}
func GetClientWorkDataNew(iCompany uint, iClientWork uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	clientworkarray := make([]interface{}, 0)
	var clientwork models.ClientWork

	result := txn.Find(&clientwork, "company_id = ? and id = ?", iCompany, iClientWork)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL647",
			DbError:   result.Error,
		}
	}
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
	return clientworkarray, models.TxnError{}
}
func GetReqDataNew(iCompany uint, iPolicy uint, iClient uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	reqArray := make([]interface{}, 0)
	resultMap, funcErr := GetReqCommNew(iCompany, iPolicy, iClient, txn)
	if funcErr.ErrorCode != "" {
		// fmt.Println("Error:", err)
		return nil, funcErr // Return nil if there was an error
	}

	// Append resultMap (from GetReqComm) as the last element
	reqArray = append(reqArray, resultMap)

	// Return the array of all results
	return reqArray, models.TxnError{}
}

func PolicyDepDataNew(iCompany uint, iPolicy uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var clnt models.Client
	result = txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var address models.Address
	result = txn.Find(&address, "company_id = ? AND id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	var pymt models.Payment
	result = txn.Find(&pymt, "company_id = ? AND policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL034", DbError: result.Error}
	}
	var glbal models.GlBal
	result = txn.Find(&glbal, "company_id = ? AND gl_rdocno = ? ", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL838", DbError: result.Error}
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
	return resultOut, models.TxnError{}
}

func PolAgntChDataNew(iCompany uint, iPolicy uint, iAgent uint, iClient uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var agntaddress models.Address
	result = txn.Find(&agntaddress, "company_id = ? AND client_id = ?", iCompany, iClient)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var poladdress models.Address
	result = txn.Find(&poladdress, "company_id = ? AND id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
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

	return resultOut, models.TxnError{}
}
func GetBankDataNew(iCompany uint, iBank uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	bankarray := make([]interface{}, 0)
	var bank models.Bank
	result := txn.Find(&bank, "id = ?", iBank)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL262", DbError: result.Error}
	}
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
	return bankarray, models.TxnError{}
}
func GetPaymentDataNew(iCompany uint, iPolicyID uint, iPayment uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var clt models.Client
	result = txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	result = txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	var paymentenq models.Payment
	result = txn.Find(&paymentenq, "company_id = ? and id = ?", iCompany, iPayment)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL034", DbError: result.Error}
	}
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
	return resultOut, models.TxnError{}

}
func ColaCancelDataNew(iCompany uint, iPolicy uint, iHistoryCode string, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	iDate := GetBusinessDate(iCompany, 0, 0)
	var iP0033Key string
	var iP0034Key string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data

	iP0034Key = iHistoryCode
	errparam := "P0034"
	err1 := GetItemD(int(iCompany), errparam, iP0034Key, iDate, &extradatap0034)
	if err1 != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iP0034Key}
	}

	iP0033Key = p0034data.Letters[0].Templates
	errparam = "P0033"
	err := GetItemD(int(iCompany), errparam, iP0033Key, iDate, &extradatap0033)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iP0033Key}
	}

	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var clnt models.Client
	result = txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var address models.Address
	result = txn.Find(&address, "company_id = ? AND id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var company models.Company
	result = txn.Find(&company, "id = ?", iCompany)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}
	compAddress := fmt.Sprintf("%s %s %s %s %s %s", company.CompanyAddress1, company.CompanyAddress2, company.CompanyAddress3, company.CompanyAddress4, company.CompanyAddress5, company.CompanyPostalCode)

	// Create a result map for each loan bill
	resultOut := map[string]interface{}{
		"ClientFullName": clnt.ClientLongName,
		"AddressLine1":   address.AddressLine1,
		"AddressLine2":   address.AddressLine2,
		"AddressLine3":   address.AddressLine3,
		"AddressLine4":   address.AddressLine4,
		"AddressLine5":   address.AddressState,
		"compname":       company.CompanyName,
		"compadd":        compAddress,
		"PolicyID":       strconv.Itoa(int(polenq.ID)),
		"Date":           DateConvert(iDate),
		"Department":     p0033data.DepartmentName,
		"SignedBy":       p0033data.DepartmentHead,
		"CompanyEmail":   p0033data.CompanyEmail,
		"CompanyNo":      p0033data.CompanyPhone,
	}

	return resultOut, models.TxnError{}
}
func AplCancelDataNew(iCompany uint, iPolicy uint, iHistoryCode string, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	iDate := GetBusinessDate(iCompany, 0, 0)
	var iP0033Key string
	var iP0034Key string

	var p0034data paramTypes.P0034Data
	var extradatap0034 paramTypes.Extradata = &p0034data

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data

	iP0034Key = iHistoryCode
	errparam := "P0034"
	err1 := GetItemD(int(iCompany), errparam, iP0034Key, iDate, &extradatap0034)
	if err1 != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iP0034Key}
	}

	iP0033Key = p0034data.Letters[0].Templates
	errparam = "P0033"
	err := GetItemD(int(iCompany), errparam, iP0033Key, iDate, &extradatap0033)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iP0033Key}
	}

	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var clnt models.Client
	result = txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var address models.Address
	result = txn.Find(&address, "company_id = ? AND id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var company models.Company
	result = txn.Find(&company, "id = ?", iCompany)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}
	compAddress := fmt.Sprintf("%s %s %s %s %s %s", company.CompanyAddress1, company.CompanyAddress2, company.CompanyAddress3, company.CompanyAddress4, company.CompanyAddress5, company.CompanyPostalCode)

	// Create a result map for each loan bill
	resultOut := map[string]interface{}{
		"ClientName":   clnt.ClientLongName,
		"AddressLine1": address.AddressLine1,
		"AddressLine2": address.AddressLine2,
		"AddressLine3": address.AddressLine3,
		"AddressLine4": address.AddressLine4,
		"AddressLine5": address.AddressState,
		"compname":     company.CompanyName,
		"compadd":      compAddress,
		"PolicyID":     strconv.Itoa(int(polenq.ID)),
		"Date":         DateConvert(iDate),
		"Department":   p0033data.DepartmentName,
		"SignedBy":     p0033data.DepartmentHead,
		"CompanyEmail": p0033data.CompanyEmail,
		"CompanyNo":    p0033data.CompanyPhone,
	}

	return resultOut, models.TxnError{}
}
func GetHIPPOLSCDDataNew(iCompany uint, iPolicyID uint, iPageSize, iOrientation string, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var priorpolenq models.PriorPolicy
	result = txn.Where("company_id = ? AND policy_id = ?", iCompany, iPolicyID).Order("created_at  DESC").First(&priorpolenq)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	var receiptenq models.Receipt
	result = txn.Where("company_id = ? AND receipt_ref_no = ?", iCompany, iPolicyID).Order("created_at  DESC").First(&receiptenq)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	var benefitenq []models.Benefit
	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	var cmp models.Company
	result = txn.Find(&cmp, "id = ?", polenq.CompanyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}
	var clt models.Client
	result = txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var agency models.Agency
	result = txn.Find(&agency, "company_id = ? and id = ?", iCompany, polenq.AgencyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	var agent models.Client
	result = txn.Find(&agent, "company_id = ? and id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	result = txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	//amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)

	// compadd := cmp.CompanyAddress1+", "+cmp.CompanyAddress2+", "+cmp.CompanyAddress3+", "+cmp.CompanyAddress4+", "+cmp.CompanyAddress5+", "+cmp.CompanyCountry+", "+cmp.CompanyPostalCode

	compadd := GetCompanyFullAddress(cmp)

	yr, _ := extractYear(polenq.PRCD)
	prcd, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.PRCD)
	riskcessdate := benefitenq[0].BRiskCessDate
	enddate, _ := ConvertYYYYMMDDtoDDMMYYYY(riskcessdate)
	clientfulladdress, funcErr := GetFullAddressNew(iCompany, polenq.AddressID, txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}

	benpalntypedesc := GetP0050ItemCodeDesc(iCompany, "HealthBenefitType", 1, benefitenq[0].BenefitType)

	slno, laname, ladob, lagender, laoccup, larel, pecdecl, ppolstdate, bprem, lanominee, lanomrel, funcErr := GetPlanLifeDataNew(iCompany, iPolicyID, txn)

	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}

	totalpreflifedisamt, totalloyaltydisamt, totalfloaterdisamt, totalonlinedisamt := GetTotalDiscountsByPolicy(iCompany, iPolicyID)
	premwogstamt, gstamt, stampdutyamt, totalpremiumamt := GetPremiumDetailsFromGL(iCompany, iPolicyID)

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
	return resultOut, models.TxnError{}

}
func GetFullAddressNew(iCompany, iAddressId uint, txn *gorm.DB) (string, models.TxnError) {
	var addr models.Address
	err := txn.First(&addr, "id = ? AND company_id = ?", iAddressId, iCompany).Error
	if err != nil {
		return "", models.TxnError{ErrorCode: "DBERR", DbError: err}
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

	return strings.Join(nonEmptyParts, ", "), models.TxnError{}
}
func GetPlanLifeDataNew(iCompany, iPolicy uint, txn *gorm.DB) ([]int, []string, []string, []string, []string, []string, []string, []string, []float64, []string, []string, models.TxnError) {

	var planLives []models.PlanLife
	result := txn.Find(&planLives, "company_id = ? and policy_id =?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, models.TxnError{ErrorCode: "GL839", DbError: result.Error}
	}
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
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
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

	return slno, laname, ladob, lagender, laoccup, larel, pecdecl, ppolstdate, bprem, lanominee, lanomrel, models.TxnError{}
}
func GetPriorPolicyDataNew(iCompany uint, iPolicyID uint, iPageSize, iOrientation string, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var cmp models.Company
	result = txn.Find(&cmp, "id = ?", polenq.CompanyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}
	var clt models.Client
	result = txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	result = txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	var priorpolenq []models.PriorPolicy
	result = txn.Find(&priorpolenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL909", DbError: result.Error}
	}
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
		"Layout": map[string]string{
			"PageSize":    iPageSize,
			"Orientation": iOrientation,
		},
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
	return resultOut, models.TxnError{}

}
func GetTermAndConditionDataNew(iCompany uint, iPolicyID uint, iPageSize, iOrientation string, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var benefitenq []models.Benefit
	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	var cmp models.Company
	result = txn.Find(&cmp, "id = ?", polenq.CompanyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}
	var clt models.Client
	result = txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	result = txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	var priorpolenq []models.PriorPolicy
	result = txn.Find(&priorpolenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL909", DbError: result.Error}
	}
	prcd, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.PRCD)
	riskcessdate := benefitenq[0].BRiskCessDate
	enddate, _ := ConvertYYYYMMDDtoDDMMYYYY(riskcessdate)

	premwogstamt, gstamt, stampdutyamt, totalpremiumamt := GetPremiumDetailsFromGL(iCompany, iPolicyID)
	totalpremiumamtwords := AmountInWords(totalpremiumamt)

	fmt.Println(premwogstamt, gstamt, stampdutyamt)

	//amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)

	// compadd := cmp.CompanyAddress1+", "+cmp.CompanyAddress2+", "+cmp.CompanyAddress3+", "+cmp.CompanyAddress4+", "+cmp.CompanyAddress5+", "+cmp.CompanyCountry+", "+cmp.CompanyPostalCode

	compadd := GetCompanyFullAddress(cmp)
	//yr := time.Now().Year()

	resultOut := map[string]interface{}{
		"Layout": map[string]string{
			"PageSize":    iPageSize,
			"Orientation": iOrientation,
		},
		"compname":         cmp.CompanyName,
		"compadd":          compadd,
		"uinno":            cmp.CompanyUid,
		"policyno":         polenq.ID,
		"ClientSalutation": clt.Salutation,
		"ClientFullName":   clt.ClientLongName,
		"totalpremiumamt":  totalpremiumamt,
		"amountinwords":    totalpremiumamtwords,
		"startdate":        prcd,
		"enddate":          enddate,
	}
	return resultOut, models.TxnError{}

}
func GetpremiumCertificateDataNew(iCompany uint, iPolicyID uint, iPageSize, iOrientation string, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var benefitenq []models.Benefit
	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	var cmp models.Company
	result = txn.Find(&cmp, "id = ?", polenq.CompanyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}
	var clt models.Client
	result = txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	result = txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	var priorpolenq []models.PriorPolicy
	result = txn.Find(&priorpolenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL909", DbError: result.Error}
	}
	prcd, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.PRCD)
	riskcessdate := benefitenq[0].BRiskCessDate
	enddate, _ := ConvertYYYYMMDDtoDDMMYYYY(riskcessdate)

	premwogstamt, gstamt, stampdutyamt, totalpremiumamt := GetPremiumDetailsFromGL(iCompany, iPolicyID)
	totalpremiumamtwords := AmountInWords(totalpremiumamt)

	fmt.Println(premwogstamt, gstamt, stampdutyamt)

	//amtinwords, csymbol := AmountinWords(paymentenq.CompanyID, paymentenq.AccAmount, paymentenq.AccCurry)

	// compadd := cmp.CompanyAddress1+", "+cmp.CompanyAddress2+", "+cmp.CompanyAddress3+", "+cmp.CompanyAddress4+", "+cmp.CompanyAddress5+", "+cmp.CompanyCountry+", "+cmp.CompanyPostalCode

	compadd := GetCompanyFullAddress(cmp)
	//yr := time.Now().Year()

	resultOut := map[string]interface{}{
		"Layout": map[string]string{
			"PageSize":    iPageSize,
			"Orientation": iOrientation,
		},
		"compname":         cmp.CompanyName,
		"compadd":          compadd,
		"uinno":            cmp.CompanyUid,
		"policyno":         polenq.ID,
		"ClientSalutation": clt.Salutation,
		"ClientFullName":   clt.ClientLongName,
		"totalpremiumamt":  totalpremiumamt,
		"amountinwords":    totalpremiumamtwords,
		"startdate":        prcd,
		"enddate":          enddate,
	}
	return resultOut, models.TxnError{}

}
func GetPOLSCDEndowmentDataNew(iCompany uint, iPolicyID uint, iPageSize, iOrientation string, p0033Data paramTypes.P0033Data, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var benefitenq []models.Benefit
	result = txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	var cmp models.Company
	result = txn.Find(&cmp, "id = ?", polenq.CompanyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL055", DbError: result.Error}
	}
	var clt models.Client
	result = txn.Find(&clt, "company_id = ? and id = ?", iCompany, polenq.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var cltAdd models.Address
	result = txn.Find(&cltAdd, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	var agency models.Agency
	result = txn.Find(&agency, "company_id = ? and id = ?", iCompany, polenq.AgencyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	var agent models.Client
	result = txn.Find(&agent, "company_id = ? and id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	result = txn.Find(&add, "company_id = ? and id = ?", iCompany, polenq.AddressID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}
	var laclient models.Client
	result = txn.Find(&laclient, "company_id = ? and id = ?", iCompany, benefitenq[0].ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	//prcd, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.PRCD)
	DistMktg := ""

	if agency.AgencyChannel == "DM" {
		DistMktg = "Yes"
	} else {
		DistMktg = "No"
	}

	riskcessdate := benefitenq[0].BRiskCessDate
	baseRCessDate, _ := ConvertYYYYMMDDtoDDMMYYYY(riskcessdate)

	BasePremDueDate, _ := GetPremDueDate(benefitenq[0].BPremCessDate, polenq.PFreq)
	clientfulladdress, funcErr := GetFullAddressNew(iCompany, polenq.AddressID, txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}
	prcd, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.PRCD)
	paidToDate, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.PaidToDate)
	polAnniDate, _ := ConvertYYYYMMDDtoDDMMYYYY(polenq.AnnivDate)
	ownerAge, _ := GetAgeFromDate(clt.ClientDob)
	ownerDOB, _ := ConvertYYYYMMDDtoDDMMYYYY(clt.ClientDob)

	agentfulladdress, funcErr := GetFullAddressNew(iCompany, agency.AddressID, txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}

	laAge, _ := GetAgeFromDate(laclient.ClientDob)
	laDOB, _ := ConvertYYYYMMDDtoDDMMYYYY(laclient.ClientDob)

	//benpalntypedesc := GetP0050ItemCodeDesc(iCompany, "HealthBenefitType", 1, benefitenq[0].BenefitType)

	AppointeeAge := []string{}
	AppointeeGender := []string{}
	AppointeeName := []string{}

	nomineeAge, nomineeLaRel, nomineeName, nomineeShare, nomineeGender, funcErr := GetNomineeDataNew(iCompany, iPolicyID, txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}
	riderCover, riderInstPrem, riderPremDueDate, riderRCessDate, riderSa := ExtractRiderDetails(benefitenq, polenq)

	resultOut := map[string]interface{}{
		"Layout": map[string]string{
			"PageSize":    iPageSize,
			"Orientation": iOrientation,
		},
		"AgentEmailId":     agent.ClientEmail,
		"AgentFullAddress": agentfulladdress,
		"AgentMobileNo":    agent.ClientMobile,
		"AgentName":        agent.ClientShortName,
		"AgentNo":          agent.ID,
		"AgentTelNo":       agent.ClientMobile,
		"AppointeeAge":     AppointeeAge,
		"AppointeeGender":  AppointeeGender,
		"AppointeeName":    AppointeeName,
		"AuthSignatory":    p0033Data.DepartmentHead, // Sign data values
		"BaseAnnPrem":      benefitenq[0].BBasAnnualPrem,
		"BaseCover":        benefitenq[0].BCoverage,
		"BaseInstPrem":     benefitenq[0].BPrem,
		"BasePremDueDate":  BasePremDueDate,
		"BaseRCessDate":    baseRCessDate,
		"BaseSa":           benefitenq[0].BSumAssured,
		"CSCAddress":       clientfulladdress,
		"ClientNo":         polenq.ClientID,
		"DistMktg":         DistMktg, // as field is not present for now keeping it as default Yes
		"Freq":             polenq.PFreq,
		"InstalPrem":       polenq.InstalmentPrem,
		"IssueDate":        "", // To be developed later

		// base cover client
		"LAName":        laclient.ClientShortName,
		"LaAge":         laAge,
		"LaAgeAdm":      "", // To be developed later
		"LaClientID":    laclient.ID,
		"LaDob":         laDOB,
		"LaFullAddress": "", // To be developed later
		"LaGender":      laclient.Gender,

		"MaturityAmt": benefitenq[0].BSumAssured, // sum assured

		"NomineeAge":   nomineeAge,
		"NomineeLaRel": nomineeLaRel,
		"NomineeName":  nomineeName,
		"NomineeShare": nomineeShare,
		"NomneeGender": nomineeGender,

		"OwnerAdd1":        cltAdd.AddressLine1,
		"OwnerAdd2":        cltAdd.AddressLine2,
		"OwnerAdd3":        cltAdd.AddressLine3,
		"OwnerAdd4":        cltAdd.AddressLine4,
		"OwnerAdd5":        cltAdd.AddressLine5,
		"OwnerAge":         ownerAge,
		"OwnerAgeAdm":      "", // To be developed later
		"OwnerClientID":    clt.ID,
		"OwnerDob":         ownerDOB,
		"OwnerFullAddress": clientfulladdress,
		"OwnerName":        clt.ClientShortName,
		"OwnerPostcode":    cltAdd.AddressPostCode,
		"OwnerTelNo":       clt.ClientMobile,
		"PRCD":             prcd,
		"PaidToDate":       paidToDate,
		"Place":            cltAdd.AddressLine4,
		"PolAnnDate":       polAnniDate,
		"PolicyNo":         polenq.ID,
		"Ppt":              benefitenq[0].BPTerm,
		"RiderCover":       riderCover,
		"RiderInstPrem":    riderInstPrem,
		"RiderPremDueDate": riderPremDueDate,
		"RiderRCessDate":   riderRCessDate,
		"RiderSa":          riderSa,

		"StaffFlag": "", // To be developed later
		"Term":      benefitenq[0].BTerm,
	}
	return resultOut, models.TxnError{}

}
func GetNomineeDataNew(iCompany, iPolicy uint, txn *gorm.DB) ([]int, []string, []string, []float64, []string, models.TxnError) {
	var nominees []models.Nominee
	err := txn.Find(&nominees, "company_id = ? AND policy_id = ?", iCompany, iPolicy).Error
	if err != nil {
		return nil, nil, nil, nil, nil, models.TxnError{ErrorCode: "GL256", DbError: err}
	}

	nomineeAge := make([]int, len(nominees))
	nomineeLaRel := make([]string, len(nominees))
	nomineeName := make([]string, len(nominees))
	nomineeShare := make([]float64, len(nominees)) // nominee percentage
	nomineeGender := make([]string, len(nominees))

	for i, nom := range nominees {
		// Fetch client details for nominee
		var client models.Client
		err := txn.First(&client, "company_id = ? AND id = ?", nom.CompanyID, nom.ClientID).Error
		if err != nil {
			return nil, nil, nil, nil, nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
		}

		// Calculate age
		age, ageErr := GetAgeFromDate(client.ClientDob) // yyyymmdd format
		if ageErr != nil {
			nomineeAge[i] = 0
		} else {
			nomineeAge[i] = age
		}

		// Relation description
		// relationDesc := GetP0050ItemCodeDesc(iCompany, "PlanLARelations", 1, nom.NomineeRelationship)
		// nomineeLaRel[i] = strings.TrimSpace(relationDesc)

		nomineeLaRel[i] = strings.TrimSpace(nom.NomineeRelationship)

		// Name
		nomineeName[i] = strings.TrimSpace(client.ClientLongName)

		// Share percentage (mapped directly)
		nomineeShare[i] = nom.NomineePercentage // <-- mapped here

		// Gender
		nomineeGender[i] = strings.TrimSpace(client.Gender)
	}

	return nomineeAge, nomineeLaRel, nomineeName, nomineeShare, nomineeGender, models.TxnError{}
}
func PrtReceiptDataNew(iCompany uint, iReceipt uint, iPolicy uint, iPa uint, p0033data paramTypes.P0033Data, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicy); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var receiptenq models.Receipt
	if result := txn.Find(&receiptenq, "company_id = ? AND id = ?", iCompany, iReceipt); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL014",
			DbError:   result.Error,
		}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var agency models.Agency
	if result := txn.Find(&agency, "company_id = ? AND id = ?", iCompany, polenq.AgencyID); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var payingauth models.PayingAuthority
	if iPa > 0 {
		result := txn.Find(&payingauth, "company_id = ? and id = ?", iCompany, iPa)

		if result.RowsAffected == 0 {
			return nil, models.TxnError{
				ErrorCode: "GL626",
				DbError:   result.Error,
			}
		}
	}
	amtinwords, csymbol := AmountinWords(receiptenq.CompanyID, receiptenq.AccAmount, receiptenq.AccCurry)

	var agent models.Client
	result := txn.Find(&agent, "company_id = ? and id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	resultOut := map[string]interface{}{
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"PayingauthID":       IDtoPrint(payingauth.ID),
		"PaName":             payingauth.PaName,
		"ReceiptFor":         receiptenq.ReceiptFor,
		"Branch":             receiptenq.Branch,
		"AccCurry":           receiptenq.AccCurry,
		"ReceiptID":          IDtoPrint(receiptenq.ID),
		"ReceiptRefNo":       IDtoPrint(receiptenq.ReceiptRefNo),
		"DateOfCollection":   DateConvert(receiptenq.DateOfCollection),
		"ClientName":         clnt.ClientLongName + " " + clnt.ClientShortName,
		"Salutation":         clnt.Salutation,
		"AgentName":          agent.ClientLongName,
		"AgentCode":          agent.ID,
		"AccAmount":          NumbertoPrint(receiptenq.AccAmount),
		"AmountInWords":      amtinwords,
		"CurrSymbol":         csymbol,
		"TypeOfReceipt":      receiptenq.TypeOfReceipt,
		"ClientMobile":       clnt.ClientMobile,
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
		"ClientID":           IDtoPrint(clnt.ID),
		"ReceiptAmount":      receiptenq.ReceiptAmount,
		"ReceiptDueDate":     DateConvert(receiptenq.ReceiptDueDate),
	}
	return resultOut, models.TxnError{}
}
func PrtPolicyBillDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data
	errparam := "Q0005"
	err := GetItemD(int(iCompany), errparam, polenq.PProduct, polenq.PRCD, &extradataq0005)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: polenq.PProduct}
	}

	gracedate := AddLeadDays(polenq.PaidToDate, q0005data.LapsedDays)

	resultOut := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"Salutation":         clnt.Salutation,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"PolicyID":           IDtoPrint(polenq.ID),
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"GracePeriodEndDate": DateConvert(gracedate),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
	}
	return resultOut, models.TxnError{}
}
func PrtPolicyLapseDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}
	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}
	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	resultOut := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"Salutation":         clnt.Salutation,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"PolicyID":           IDtoPrint(polenq.ID),
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
	}
	return resultOut, models.TxnError{}

}
func PrtCollectionDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	oRiskCessDate := ""
	for _, b := range benefitenq {
		if oRiskCessDate < b.BRiskCessDate {
			oRiskCessDate = b.BRiskCessDate
		}
	}

	resulout := map[string]interface{}{
		"CompanyName":       cmp.CompanyName,
		"CompanyAddress1":   cmp.CompanyAddress1,
		"CompanyAddress2":   cmp.CompanyAddress2,
		"CompanyAddress3":   cmp.CompanyAddress3,
		"CompanyPostalCode": cmp.CompanyPostalCode,
		"LetterDate":        DateConvert(iDate),
		"ClientShortName":   clnt.ClientShortName,
		"ClientLongName":    clnt.ClientLongName,
		"Salutation":        clnt.Salutation,
		"AddressLine1":      add.AddressLine1,
		"AddressLine2":      add.AddressLine2,
		"AddressLine3":      add.AddressLine3,
		"AddressLine4":      add.AddressLine4,
		"AddressLine5":      add.AddressLine5,
		"AddressPostCode":   add.AddressPostCode,
		"PolicyID":          IDtoPrint(polenq.ID),
		"PaidToDate":        DateConvert(polenq.PaidToDate),
		"ClientID":          IDtoPrint(clnt.ID),
		"PRCD":              DateConvert(polenq.PRCD),
		"PFreq":             polenq.PFreq,
		"InstalmentPrem":    NumbertoPrint(polenq.InstalmentPrem),
		"RiskCessDate":      DateConvert(oRiskCessDate),
		"Department":        p0033data.DepartmentName,
		"DepartmentHead":    p0033data.DepartmentHead,
		"CoEmail":           p0033data.CompanyEmail,
		"CoPhone":           p0033data.CompanyPhone,
	}
	return resulout, models.TxnError{}
}
func PrtAnniDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
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

	AnnivDate := String2Date(polenq.AnnivDate)
	oPRCD := String2Date(polenq.PRCD)
	completedyears, _, _, _, _, _ := DateDiff(AnnivDate, oPRCD, "")

	// RiskCessDate := String2Date(oRiskCessDate)
	// sPRCD := String2Date(polenq.PRCD)
	// RiskTerm, _, _, _, _, _ := DateDiff(sPRCD, RiskCessDate)

	// PremCessDate := String2Date(oPremCessDate)
	// sPRCD = String2Date(polenq.PRCD)
	// PremTerm, _, _, _, _, _ := DateDiff(PremCessDate, sPRCD, "")

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"Salutation":         clnt.Salutation,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"PolicyID":           IDtoPrint(polenq.ID),
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"RiskCessDate":       DateConvert(oRiskCessDate),
		"PremCessDate":       DateConvert(oPremCessDate),
		"AnnivDate":          DateConvert(polenq.AnnivDate),
		"CompletedYears":     completedyears,
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
	}
	return resultout, models.TxnError{}
}
func PrtAnniILPDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
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

	sAnnivDate := String2Date(polenq.AnnivDate)
	sPRCD := String2Date(polenq.PRCD)
	ocompletedyears, _, _, _, _, _ := DateDiff(sAnnivDate, sPRCD, "")

	sPremCessDate := String2Date(oPremCessDate)
	oPremTerm, _, _, _, _, _ := DateDiff(sPremCessDate, sPRCD, "")

	oRevBonus, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "ReversionaryBonus", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}
	oAccumDiv, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "AccumDividend", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}
	oAccumDivInt, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "AccumDivInt", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}
	oCashDep, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "CashDeposit", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}
	oPolicyDeposit, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "PolicyDeposit", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}
	oAplAmt, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "AplAmount", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"Salutation":         clnt.Salutation,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"PolicyID":           IDtoPrint(polenq.ID),
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"RiskCessDate":       DateConvert(oRiskCessDate),
		"PremCessDate":       DateConvert(oPremCessDate),
		"CompletedYears":     ocompletedyears,
		"PremiumTerm":        oPremTerm,
		"AnnivDate":          DateConvert(polenq.AnnivDate),
		"RevBonus":           NumbertoPrint(oRevBonus),
		"AccDividend":        NumbertoPrint(oAccumDiv),
		"AccDivInt":          NumbertoPrint(oAccumDivInt),
		"CashDeposit":        NumbertoPrint(oCashDep),
		"PolicyDeposit":      NumbertoPrint(oPolicyDeposit),
		"AplAmount":          NumbertoPrint(oAplAmt),
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
	}
	return resultout, models.TxnError{}
}
func PrtExpiDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, iTranno uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	if result := txn.Find(&benefitenq, "company_id = ? AND policy_id = ? ", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	var bRiskCessDate string
	BCoverage := []string{}
	bPremCessDate := []string{}
	bclientId := []uint{}

	for _, val := range benefitenq {
		BCoverage = append(BCoverage, val.BCoverage)
		bPremCessDate = append(bPremCessDate, DateConvert(val.BPremCessDate))
		bclientId = append(bclientId, val.ClientID)
	}
	if len(benefitenq) > 0 {
		bRiskCessDate = benefitenq[0].BRiskCessDate
	}

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"PolicyID":           IDtoPrint(polenq.ID),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"BRiskCessDate":      DateConvert(bRiskCessDate),
		"BCoverage":          BCoverage,
		"BPremCessDate":      bPremCessDate,
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
		"LifeAClntId":        bclientId,
		"PaidToDate":         DateConvert(polenq.PaidToDate),
	}
	return resultout, models.TxnError{}
}
func PrtPremstDataNew(iCompany uint, iPolicyID uint, iBenifitID uint, iDate string, p0033data paramTypes.P0033Data, iTranno uint, iAgency uint, iFromDate string, iToDate string, iHistoryCode string, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	if result := txn.Find(&benefitenq, "company_id = ? AND policy_id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	var agency models.Agency
	result := txn.Find(&agency, "company_id = ? AND id = ?", iCompany, iAgency)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	var agent models.Client
	result = txn.Find(&agent, "company_id = ? AND id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var ilpsummary []models.IlpSummary
	result = txn.Find(&ilpsummary, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL135",
			DbError:   result.Error,
		}
	}
	var ilpfund []models.IlpFund
	result = txn.Find(&ilpfund, "company_id = ? and policy_id = ? and benefit_id = ? and history_code = ?  and tranno =  ?", iCompany, iPolicyID, iBenifitID, iHistoryCode, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL784",
			DbError:   result.Error,
		}
	}
	var policyhistory models.PHistory
	result = txn.Find(&policyhistory, "company_id = ? and policy_id = ? and history_code = ?  and tranno =  ?", iCompany, iPolicyID, iHistoryCode, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL919",
			DbError:   result.Error,
		}
	}
	oldFundCode := []string{}
	oldFType := []string{}
	oldFCurr := []string{}
	oldFPercentage := []float64{}
	//oldFuints := []float64{}

	if funds, ok := policyhistory.PrevData["IlpFunds"].([]interface{}); ok {
		for _, f := range funds {

			if fund, ok := f.(map[string]interface{}); ok {
				if fname, ok := fund["FundCode"].(string); ok {
					oldFundCode = append(oldFundCode, fname)
				}
				if fper, ok := fund["FundPercentage"].(float64); ok {
					oldFPercentage = append(oldFPercentage, fper)
				}
				if fcur, ok := fund["FundCurr"].(string); ok {
					oldFCurr = append(oldFCurr, fcur)
				}
				if ft, ok := fund["FundType"].(string); ok {
					oldFType = append(oldFType, ft)
				}
			}
		}
	}
	newFundCode := []string{}
	newPercentage := []float64{}
	newFType := []string{}
	newFCurr := []string{}
	for _, log := range ilpfund {
		newFundCode = append(newFundCode, log.FundCode)
		newPercentage = append(newPercentage, log.FundPercentage)
		newFType = append(newFType, log.FundType)
		newFCurr = append(newFCurr, log.FundCurr)
	}

	summFundCode := []string{}
	summFundUnits := []float64{}
	summFType := []string{}
	summFCurr := []string{}
	summPolId := []uint{}
	summFBenId := []uint{}
	for _, log := range ilpsummary {
		summFundCode = append(summFundCode, log.FundCode)
		summFundUnits = append(summFundUnits, log.FundUnits)
		summFType = append(summFType, log.FundType)
		summFCurr = append(summFCurr, log.FundCurr)
		summPolId = append(summPolId, log.PolicyID)
		summFBenId = append(summFBenId, log.BenefitID)
	}
	// Prepare benefit details
	fundpolID := make([]uint, len(newFundCode))
	fundBenId := make([]uint, len(newFundCode))
	fundBCov := make([]string, len(newFundCode))

	for i := range newFundCode {
		fundpolID[i] = iPolicyID
		fundBenId[i] = iBenifitID
		fundBCov[i] = benefitenq[0].BCoverage
	}
	var (
		bRiskCessDate string
		BCoverage     string
		bPremCessDate string
		BSumAssured   uint64
		BInsPrem      float64
	)
	if len(benefitenq) > 0 {
		bRiskCessDate = benefitenq[0].BRiskCessDate
		BCoverage = benefitenq[0].BCoverage
		bPremCessDate = benefitenq[0].BPremCessDate
		BSumAssured = benefitenq[0].BSumAssured
		BInsPrem = benefitenq[0].BPrem
	}

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"Salutation":         clnt.Salutation,
		"PProduct":           polenq.PProduct,
		"AnnivDate":          DateConvert(polenq.AnnivDate),
		"PolicyID":           polenq.ID,
		"ClientID":           clnt.ID,
		"PolCurr":            polenq.PBillCurr,
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"InsPrem":            BInsPrem,
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"BSumAssured":        BSumAssured,
		"BRiskCessDate":      DateConvert(bRiskCessDate),
		"BCoverage":          BCoverage,
		"BPremCessDate":      DateConvert(bPremCessDate),
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
		"Agentname":          agent.ClientLongName,
		"FromDate":           DateConvert(iFromDate),
		"ToDate":             DateConvert(iToDate),
		"OldFCode":           oldFundCode,
		"OldFType":           oldFType,
		"OldFCurr":           oldFCurr,
		"OldPercentage":      oldFPercentage,
		"NewFCode":           newFundCode,
		"NewFCurr":           newFCurr,
		"NewFType":           newFType,
		"NewFPercentage":     newPercentage,
		"NewFPoId":           fundpolID,
		"NFBenId":            fundBenId,
		"SummFCode":          summFundCode,
		"SummFCurr":          summFCurr,
		"SummFType":          summFType,
		"SummFUnits":         summFundUnits,
		"SummFPolId":         summPolId,
		"SummBenId":          summFBenId,
	}
	return resultout, models.TxnError{}
}
func PrtFreqChangeDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, iAgency uint, iHistoryCode string, iTranno uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.First(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if result := txn.First(&cmp, "id = ?", iCompany); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	var clnt models.Client
	if result := txn.First(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	var add models.Address
	if result := txn.First(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}

	bCoverage := []string{}
	obPrem := []float64{}
	nbPrem := []float64{}
	bPremCessDate := []string{}

	var benefitenq []models.Benefit
	result := txn.Find(&benefitenq, "company_id = ? AND policy_id = ?", iCompany, iPolicyID)
	if result.Error != nil {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}
	var agency models.Agency
	result = txn.First(&agency, "company_id = ? AND id = ?", iCompany, iAgency)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	var agent models.Client
	result = txn.First(&agent, "company_id = ? AND id = ?", iCompany, agency.ClientID)
	if result.Error != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
	}
	var phistory models.PHistory
	result = txn.Find(&phistory, "company_id = ? and policy_id = ? and history_code = ?  and tranno =  ?", iCompany, iPolicyID, iHistoryCode, iTranno)
	if result.Error != nil {
		return nil, models.TxnError{
			ErrorCode: "GL919",
			DbError:   result.Error,
		}
	}

	for _, ben := range benefitenq {
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", ben.BCoverage, 1)
		bCoverage = append(bCoverage, oCoverage)
		nbPrem = append(nbPrem, ben.BPrem)
		bPremCessDate = append(bPremCessDate, DateConvert(ben.BPremCessDate))
	}
	if benefits, ok := phistory.PrevData["Benefits"].([]interface{}); ok {
		for _, b := range benefits {
			if ben, ok := b.(map[string]interface{}); ok {
				if prem, ok := ben["BPrem"].(float64); ok {
					obPrem = append(obPrem, prem)
				}
			}
		}
	}

	prevPFreq := ""
	prevInsta := 0.0
	previousPolicy := phistory.PrevData["Policy"]
	if policyMap, ok := previousPolicy.(map[string]interface{}); ok {
		if pf, ok := policyMap["PFreq"].(string); ok {
			prevPFreq = pf
		}
		if pf, ok := policyMap["InstalmentPrem"].(float64); ok {
			prevInsta = pf
		}
	}
	var bRiskCessDate string
	if len(benefitenq) > 0 {
		bRiskCessDate = benefitenq[0].BRiskCessDate
	}
	pdate := make([]string, len(bCoverage))

	for i := range bCoverage {
		pdate[i] = DateConvert(polenq.PaidToDate)
	}

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"Salutation":         clnt.Salutation,
		"PProduct":           polenq.PProduct,
		"PolicyID":           IDtoPrint(polenq.ID),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"NPFreq":             polenq.PFreq,
		"OPFreq":             prevPFreq,
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"NInstalmentPrem":    polenq.InstalmentPrem,
		"OInstalmentPrem":    prevInsta,
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
		"BRiskCessDate":      bPremCessDate,
		"BCoverage":          bCoverage,
		"NBPrem":             nbPrem,
		"OBPrem":             obPrem,
		"BPaidToDate":        pdate,
		"PolEndDate":         DateConvert(bRiskCessDate),
	}

	return resultout, models.TxnError{}
}
func PrtSachangeDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, iAgency uint, iTranno uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	if result := txn.Find(&benefitenq, "company_id = ? AND policy_id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	var agency models.Agency
	result := txn.Find(&agency, "company_id = ? AND id = ?", iCompany, iAgency)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	var agent models.Client
	result = txn.Find(&agent, "company_id = ? AND id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var sachangeenq []models.SaChange
	result = txn.Find(&sachangeenq, "company_id = ? and policy_id = ? and tranno = ?", iCompany, iPolicyID, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL426",
			DbError:   result.Error,
		}
	}
	osumAssured := []uint64{}
	oBterm := []uint{}
	oPbterm := []uint{}
	oPrem := []float64{}
	nsumAssured := []uint64{}
	nBterm := []uint{}
	nPbterm := []uint{}
	nPrem := []float64{}
	bCoverage := []string{}
	prevInstalmentPrem := 0.0
	newInstalmentPrem := 0.0
	for _, saChan := range sachangeenq {
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", saChan.BCoverage, 1)
		bCoverage = append(bCoverage, oCoverage)
		osumAssured = append(osumAssured, saChan.BSumAssured)
		oBterm = append(oBterm, saChan.BTerm)
		oPbterm = append(oPbterm, saChan.BPTerm)
		oPrem = append(oPrem, saChan.BPrem)
		nsumAssured = append(nsumAssured, saChan.NSumAssured)
		nBterm = append(nBterm, saChan.NTerm)
		nPbterm = append(nPbterm, saChan.NPTerm)
		nPrem = append(nPrem, saChan.NPrem)
		prevInstalmentPrem += saChan.BPrem
		newInstalmentPrem += saChan.NPrem
	}

	var (
		bRiskCessDate string
	)
	if len(benefitenq) > 0 {
		bRiskCessDate = benefitenq[0].BRiskCessDate

	}

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"Salutation":         clnt.Salutation,
		"PProduct":           polenq.PProduct,
		"PolicyID":           IDtoPrint(polenq.ID),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"EndDate":            DateConvert(bRiskCessDate),
		"PFreq":              polenq.PFreq,
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"InstalmentPrem":     prevInstalmentPrem,
		"NInstalmentPrem":    newInstalmentPrem,
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
		"BCoverage":          bCoverage,
		"OSumAssured":        osumAssured,
		"OBterm":             oBterm,
		"OBPterm":            oPbterm,
		"oPrem":              oPrem,
		"NSumAssured":        nsumAssured,
		"NBterm":             nBterm,
		"NBPterm":            nPbterm,
		"NPrem":              nPrem,
	}

	return resultout, models.TxnError{}
}
func PrtCompaddDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, iAgency uint, iTranno uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	if result := txn.Find(&benefitenq, "company_id = ? AND policy_id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	var agency models.Agency
	result := txn.Find(&agency, "company_id = ? AND id = ?", iCompany, iAgency)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	var agent models.Client
	result = txn.Find(&agent, "company_id = ? AND id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var sachangeenq []models.SaChange
	result = txn.Find(&sachangeenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL426",
			DbError:   result.Error,
		}
	}
	var (
		bRiskCessDate string
	)
	if len(benefitenq) > 0 {
		bRiskCessDate = benefitenq[0].BRiskCessDate

	}
	var addcomp []models.Addcomponent
	result = txn.Find(&addcomp, "company_id = ? and policy_id =? and tranno = ?", iCompany, iPolicyID, iTranno)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL393",
			DbError:   result.Error,
		}
	}
	bCoverage := []string{}
	bTerm := []uint{}
	bPrem := []float64{}
	bCaRiskCessDate := []string{}
	bStartDate := []string{}
	bClientId := []uint{}
	bSumAssured := []uint64{}
	bPTerm := []uint{}

	for _, adcomp := range addcomp {
		_, oCoverage, _ := GetParamDesc(iCompany, "Q0006", adcomp.BCoverage, 1)
		bCoverage = append(bCoverage, oCoverage)
		bTerm = append(bTerm, adcomp.BTerm)
		bPrem = append(bPrem, adcomp.BPrem)
		bCaRiskCessDate = append(bCaRiskCessDate, DateConvert(adcomp.BRiskCessDate))
		bStartDate = append(bStartDate, DateConvert(adcomp.BStartDate))
		bClientId = append(bClientId, adcomp.ClientID)
		bSumAssured = append(bSumAssured, adcomp.BSumAssured)
		bPTerm = append(bPTerm, adcomp.BTerm)
	}

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"Salutation":         clnt.Salutation,
		"PProduct":           polenq.PProduct,
		"PolicyID":           IDtoPrint(polenq.ID),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"PolEndDate":         DateConvert(bRiskCessDate),
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
		"BCoverage":          bCoverage,
		"BStartDate":         bStartDate,
		"BLAClientID":        bClientId,
		"BSumAssured":        bSumAssured,
		"BTerm":              bTerm,
		"BPTerm":             bPTerm,
		"BPrem":              bPrem,
		"BRiskCessDate":      bCaRiskCessDate,
	}

	return resultout, models.TxnError{}
}
func PrtSurrDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, iAgency uint, iTronno uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	if result := txn.Find(&benefitenq, "company_id = ? AND policy_id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	var agency models.Agency
	result := txn.Find(&agency, "company_id = ? AND id = ?", iCompany, iAgency)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	var agent models.Client
	result = txn.Find(&agent, "company_id = ? AND id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var sachangeenq []models.SaChange
	result = txn.Find(&sachangeenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL426",
			DbError:   result.Error,
		}
	}
	var (
		bRiskCessDate string
		bSumA         float64
	)
	if len(benefitenq) > 0 {
		bRiskCessDate = benefitenq[0].BRiskCessDate
		bSumA = float64(benefitenq[0].BSumAssured)

	}
	var addcomp []models.Addcomponent
	result = txn.Find(&addcomp, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL393",
			DbError:   result.Error,
		}
	}
	var surrhenq models.SurrH
	result = txn.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL942",
			DbError:   result.Error,
		}
	}
	var surrdenq models.SurrD
	result = txn.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL942",
			DbError:   result.Error,
		}
	}
	oCashDep, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "CashDeposit", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}

	resultout := map[string]interface{}{
		"CompanyName":        cmp.CompanyName,
		"CompanyFullAddress": cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":         DateConvert(iDate),
		"ClientShortName":    clnt.ClientShortName,
		"ClientLongName":     clnt.ClientLongName,
		"AddressLine1":       add.AddressLine1,
		"AddressLine2":       add.AddressLine2,
		"AddressLine3":       add.AddressLine3,
		"AddressLine4":       add.AddressLine4,
		"AddressLine5":       add.AddressLine5,
		"AddressPostCode":    add.AddressPostCode,
		"Salutation":         clnt.Salutation,
		"PProduct":           polenq.PProduct,
		"PolicyID":           IDtoPrint(polenq.ID),
		"ClientID":           IDtoPrint(clnt.ID),
		"PRCD":               DateConvert(polenq.PRCD),
		"PFreq":              polenq.PFreq,
		"PaidToDate":         DateConvert(polenq.PaidToDate),
		"InstalmentPrem":     NumbertoPrint(polenq.InstalmentPrem),
		"RiskCessDate":       DateConvert(bRiskCessDate),
		"Department":         p0033data.DepartmentName,
		"DepartmentHead":     p0033data.DepartmentHead,
		"CoEmail":            p0033data.CompanyEmail,
		"CoPhone":            p0033data.CompanyPhone,
		"EffectiveDate":      DateConvert(surrhenq.EffectiveDate),
		"SurrAmount":         surrdenq.SurrAmount,
		"RevBonus":           surrdenq.RevBonus,
		"InterimBonus":       surrdenq.InterimBonus,
		"TerminalBonus":      surrdenq.TerminalBonus,
		"AccumDividend":      surrdenq.AccumDividend,
		"AccumDivInt":        surrdenq.AccumDivInt,
		"CashDeposit":        NumbertoPrint(oCashDep),
		"PolicyDepost":       surrhenq.PolicyDepost,
		"AplAmount":          surrhenq.AplAmount,
		"LoanAmount":         surrhenq.LoanAmount,
		"TotalSurrPayable":   surrhenq.TotalSurrPayable,
		"SurrenderDate":      surrhenq.EffectiveDate,
		"BsumAssured":        bSumA,
	}

	return resultout, models.TxnError{}
}
func PrtMatyDataNew(iCompany uint, iPolicyID uint, iDate string, p0033data paramTypes.P0033Data, iAgency uint, txn *gorm.DB) (map[string]interface{}, models.TxnError) {

	var polenq models.Policy
	if result := txn.Find(&polenq, "company_id = ? AND id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL003", DbError: result.Error}
	}

	var cmp models.Company
	if err := txn.First(&cmp, polenq.CompanyID).Error; err != nil {
		return nil, models.TxnError{ErrorCode: "DBERR", DbError: err}
	}

	var clnt models.Client
	if result := txn.Find(&clnt, "company_id = ? AND id = ?", iCompany, polenq.ClientID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}

	var add models.Address
	if result := txn.Find(&add, "company_id = ? AND id = ?", iCompany, polenq.AddressID); result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL035", DbError: result.Error}
	}

	var benefitenq []models.Benefit
	if result := txn.Find(&benefitenq, "company_id = ? AND policy_id = ?", iCompany, iPolicyID); result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
	}

	var agency models.Agency
	result := txn.Find(&agency, "company_id = ? AND id = ?", iCompany, iAgency)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL312",
			DbError:   result.Error,
		}
	}
	var agent models.Client
	result = txn.Find(&agent, "company_id = ? AND id = ?", iCompany, agency.ClientID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{ErrorCode: "GL050", DbError: result.Error}
	}
	var sachangeenq []models.SaChange
	result = txn.Find(&sachangeenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL426",
			DbError:   result.Error,
		}
	}
	var (
		bRiskCessDate string
		BSumAssured   uint64
	)
	if len(benefitenq) > 0 {
		bRiskCessDate = benefitenq[0].BRiskCessDate
		BSumAssured = benefitenq[0].BSumAssured

	}
	var addcomp []models.Addcomponent
	result = txn.Find(&addcomp, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL393",
			DbError:   result.Error,
		}
	}
	var surrhenq models.SurrH
	result = txn.Find(&surrhenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL942",
			DbError:   result.Error,
		}
	}
	var surrdenq models.SurrD
	result = txn.Find(&surrdenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL942",
			DbError:   result.Error,
		}
	}
	var mathenq models.MaturityH
	result = txn.Find(&mathenq, "company_id = ? and policy_id = ?", iCompany, iPolicyID)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL822",
			DbError:   result.Error,
		}
	}
	oCashDep, funcErr := GetGlBalNew(iCompany, uint(iPolicyID), "CashDeposit", txn)
	if funcErr.ErrorCode != "" {
		return nil, funcErr
	}

	resultout := map[string]interface{}{
		"CompanyName":          cmp.CompanyName,
		"CompanyFullAddress":   cmp.CompanyAddress1 + " " + cmp.CompanyAddress2 + " " + cmp.CompanyAddress3 + " " + cmp.CompanyPostalCode,
		"LetterDate":           DateConvert(iDate),
		"ClientShortName":      clnt.ClientShortName,
		"ClientLongName":       clnt.ClientLongName,
		"AddressLine1":         add.AddressLine1,
		"AddressLine2":         add.AddressLine2,
		"AddressLine3":         add.AddressLine3,
		"AddressLine4":         add.AddressLine4,
		"AddressLine5":         add.AddressLine5,
		"AddressPostCode":      add.AddressPostCode,
		"Salutation":           clnt.Salutation,
		"PProduct":             polenq.PProduct,
		"PolicyID":             IDtoPrint(polenq.ID),
		"ClientID":             IDtoPrint(clnt.ID),
		"PRCD":                 DateConvert(polenq.PRCD),
		"PFreq":                polenq.PFreq,
		"PaidToDate":           DateConvert(polenq.PaidToDate),
		"InstalmentPrem":       NumbertoPrint(polenq.InstalmentPrem),
		"BRiskCessDate":        DateConvert(bRiskCessDate),
		"Department":           p0033data.DepartmentName,
		"DepartmentHead":       p0033data.DepartmentHead,
		"CoEmail":              p0033data.CompanyEmail,
		"CoPhone":              p0033data.CompanyPhone,
		"ComponantAddData":     addcomp,
		"EffectiveDate":        DateConvert(surrhenq.EffectiveDate),
		"SurrAmount":           surrdenq.SurrAmount,
		"RevBonus":             surrdenq.RevBonus,
		"InterimBonus":         surrdenq.InterimBonus,
		"TerminalBonus":        surrdenq.TerminalBonus,
		"AccumDividend":        surrdenq.AccumDividend,
		"AccumDivInt":          surrdenq.AccumDivInt,
		"CashDeposit":          NumbertoPrint(oCashDep),
		"PolicyDepost":         surrhenq.PolicyDepost,
		"AplAmount":            surrhenq.AplAmount,
		"LoanAmount":           surrhenq.LoanAmount,
		"TotalSurrPayable":     surrhenq.TotalSurrPayable,
		"BSumAssured":          BSumAssured,
		"TotalMaturityPayable": mathenq.TotalMaturityPayable,
	}

	return resultout, models.TxnError{}
}
func GetLoanDataNew(iCompany uint, iPolicy uint, iEffectiveDate string, iOsLoanInterest float64, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	combinedData := make(map[string]interface{})

	loanArray := make([]interface{}, 0)
	extraData := make([]map[string]interface{}, 0)
	overallData := make(map[string]interface{}) // Create a map for overall data

	var loanenq []models.Loan

	// Fetch loans for the specified company and policy
	result := txn.Find(&loanenq, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL658",
			DbError:   result.Error,
		}
	}
	var overallLoanAmount float64
	var overallstampduty float64
	var finalLoanAmountTotal float64

	var p0072data paramTypes.P0072Data
	var extradata4 paramTypes.Extradata = &p0072data
	errparam := "P0072"
	err := GetItemD(int(iCompany), errparam, "LN001", iEffectiveDate, &extradata4)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: "LN001"}
	}
	// Map to keep track of already printed LoanSeqNumber values
	printedLoanSeqNumbers := make(map[uint]bool)

	for i := 0; i < len(loanenq); i++ {
		var benefit models.Benefit
		result = txn.First(&benefit, "company_id = ? and policy_id = ? and benefit_id = ?", iCompany, iPolicy, loanenq[i].BenefitID)
		if result.Error != nil {
			return nil, models.TxnError{ErrorCode: "DBERR", DbError: result.Error}
		}
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

	return combinedData, models.TxnError{}
}
func GetAllLoanInterestDataNew(iCompany uint, iPolicy uint, iEffectiveDate string, txn *gorm.DB) ([]interface{}, models.TxnError) {
	var benefitenq []models.Benefit
	allLoanOs := make([]interface{}, 0)
	result := txn.Find(&benefitenq, "company_id = ? and policy_id = ? ", iCompany, iPolicy)
	if result.Error != nil {
		return nil, models.TxnError{
			ErrorCode: "GL018",
			DbError:   result.Error,
		}
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
		errparam := "Q0006"
		err := GetItemD(int(iCompany), errparam, iKey, iDate, &extradataq0006)
		if err != nil {
			return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: iKey}
		}
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

	return allLoanOs, models.TxnError{}
}
func LoanCapDataNew(iCompany uint, iPolicy uint, iEffectiveDate string, minLoanBillDueDate string, maxLoanBillDueDate string, itotalcapamount float64, itotalInterest float64, itotalOsDue uint, txn *gorm.DB) ([]interface{}, models.TxnError) {
	allLoanCap := make([]interface{}, 0)

	var loanenq []models.Loan
	var prevloancapamount float64
	var oLoanInt float64

	result := txn.Find(&loanenq, "company_id = ? and policy_id = ? and  loan_type = ? and loan_status = ? and next_cap_date <=?", iCompany, iPolicy, "P", "AC", iEffectiveDate)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL658",
			DbError:   result.Error,
		}
	}
	var policyenq models.Policy

	result = txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL003",
			DbError:   result.Error,
		}
	}
	// var minLoanDate string
	// var maxLoanDate string
	var p0072data paramTypes.P0072Data
	var extradata paramTypes.Extradata = &p0072data
	errparam := "P0072"
	err := GetItemD(int(iCompany), errparam, "LN001", iEffectiveDate, &extradata)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: "LN001"}
	}
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
	return allLoanCap, models.TxnError{}

}
func LoanBillDataNew(iCompany uint, iPolicy uint, iEffectiveDate string, txn *gorm.DB) ([]interface{}, models.TxnError) {

	var loanenq []models.Loan

	result := txn.Find(&loanenq, "company_id = ? and policy_id = ? and  loan_type = ? and loan_status = ? and next_int_bill_date<=?", iCompany, iPolicy, "P", "AC", iEffectiveDate)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL658",
			DbError:   result.Error,
		}
	}
	var policyenq models.Policy

	result = txn.Find(&policyenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL003",
			DbError:   result.Error,
		}
	}
	var loanbillupd models.LoanBill
	var oLoanOS float64
	var oLoanIntOS float64

	var p0072data paramTypes.P0072Data
	var extradata paramTypes.Extradata = &p0072data
	errparam := "P0072"
	err := GetItemD(int(iCompany), errparam, "LN001", iEffectiveDate, &extradata)
	if err != nil {
		return nil, models.TxnError{ErrorCode: "PARME", ParamName: errparam, ParamItem: "LN001"}
	}
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

	return loanbill, models.TxnError{}

}
func LoanBillsInterestDataNew(iCompany uint, iPolicy uint, iSeqNo uint, iCurrentIntDue float64, txn *gorm.DB) (map[string]interface{}, models.TxnError) {
	var loanbillupd1 []models.LoanBill

	result := txn.Order("CASE WHEN loan_seq_number = 1 THEN 0 WHEN loan_seq_number = 2 THEN 1 ELSE 2 END").
		Find(&loanbillupd1, "company_id = ? and policy_id = ?", iCompany, iPolicy)
	if result.RowsAffected == 0 {
		return nil, models.TxnError{
			ErrorCode: "GL821",
			DbError:   result.Error,
		}
	}
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
	return combinedData, models.TxnError{}
}
