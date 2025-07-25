package initializers

import "github.com/FuturaInsTech/GoLifeLib/models"

// Table Name should start with Capital Letter
func SyncDatabase() {
	// All Drop Down Come First
	//DB.AutoMigrate(&models.Language{})
	// DB.AutoMigrate(&models.Currency{})
	// DB.AutoMigrate(&models.Gender{})
	// DB.AutoMigrate(&models.Salutation{})
	// DB.AutoMigrate(&models.UserStatus{})
	// DB.AutoMigrate(&models.UserGroup{})
	// DB.AutoMigrate(&models.ContHeader{})
	// DB.AutoMigrate(&models.CompanyStatus{})
	// DB.AutoMigrate(&models.Company{})
	// DB.AutoMigrate(&models.Transaction{})
	// DB.AutoMigrate(&models.TDFRule{})
	// DB.AutoMigrate(&models.Permission{})
	// DB.AutoMigrate(&models.LeadChannel{})
	// DB.AutoMigrate(&models.User{})
	// DB.AutoMigrate(&models.Coverage{})
	// DB.AutoMigrate(&models.Load{})
	// DB.AutoMigrate(&models.Client{})
	// DB.AutoMigrate(&models.ParamDesc{})
	// DB.AutoMigrate(&models.Param{})
	// DB.AutoMigrate(&models.Address{})
	// DB.AutoMigrate(&models.Relationship{})
	// DB.AutoMigrate(&models.Nominee{})
	// DB.AutoMigrate(&models.UinMaster{})
	// DB.AutoMigrate(&models.Error{})
	// DB.AutoMigrate(&models.Bank{})
	// DB.AutoMigrate(&models.Agency{})
	// DB.AutoMigrate(models.LeadDetail{})
	// DB.AutoMigrate(models.LeadFollowup{})

	// DB.AutoMigrate(&models.LeadAllocation{})
	// DB.AutoMigrate(&models.Campaign{})
	// DB.AutoMigrate(&models.CampaignComp{})

	// DB.AutoMigrate(&models.FieldValidator{})
	// DB.AutoMigrate(&models.Level{})
	DB.AutoMigrate(&models.Policy{})
	// DB.AutoMigrate(&models.Extra{})
	// DB.AutoMigrate(&models.PHistory{})
	// DB.AutoMigrate(&models.GlType{})
	// DB.AutoMigrate(&models.AccountCode{})
	// DB.AutoMigrate(&models.GlMove{})
	// DB.AutoMigrate(&models.GlBal{})
	// DB.AutoMigrate(&models.Receipt{})
	// DB.AutoMigrate(&models.Tdfh{})
	// DB.AutoMigrate(&models.TDFPolicy{})
	//DB.AutoMigrate(&sudmodels.UIMaster{})
	// DB.AutoMigrate(&models.TDFParam{})
	// DB.AutoMigrate(&models.SurvB{})
	// DB.AutoMigrate(&models.Communication{})
	// DB.AutoMigrate(&models.Uwreason{})
	// DB.AutoMigrate(&models.MedProvider{})
	// DB.AutoMigrate(&models.MedReq{})
	// DB.AutoMigrate(&models.Nominee{})
	// DB.AutoMigrate(&models.DeathH{})
	// DB.AutoMigrate(&models.DeathD{})

	// Following Should be deleted
	// DB.AutoMigrate(&quotation.Qcommunication{})
	// DB.AutoMigrate(&models.Quotation{})
	// DB.AutoMigrate(&quotation.QHeader{})
	// DB.AutoMigrate(&quotation.QDetail{})
	// DB.AutoMigrate(&quotation.QBenIllValue{})

	// DB.AutoMigrate(&models.Payer{})
	// DB.AutoMigrate(&models.SaChange{})
	// DB.AutoMigrate(&models.Assignee{})
	// DB.AutoMigrate(&models.Addcomponent{})
	// DB.AutoMigrate(&models.Mrta{})
	// DB.AutoMigrate(&models.SurrH{})
	// DB.AutoMigrate(&models.SurrD{})
	// DB.AutoMigrate(&models.BusinessDate{})
	// DB.AutoMigrate(&models.MaturityH{})
	// DB.AutoMigrate(&models.MaturityD{})
	// DB.AutoMigrate(&models.PolBill{})
	// DB.AutoMigrate(&models.CriticalIllness{})
	// DB.AutoMigrate(&models.TransactionLock{})
	// DB.AutoMigrate(&models.IBenefit{})
	DB.AutoMigrate(&models.Benefit{})
	// DB.AutoMigrate(&models.Payment{})
	// DB.AutoMigrate(&models.IlpFund{})
	// DB.AutoMigrate(&models.IlpPrice{})
	// DB.AutoMigrate(&models.IlpSummary{})
	// DB.AutoMigrate(&models.IlpTransaction{})
	// DB.AutoMigrate(&models.IlpAnnSummary{})
	// DB.AutoMigrate(&models.IlpSwitchHeader{})
	// DB.AutoMigrate(&models.IlpSwitchFund{})
	// DB.AutoMigrate(&models.IlpStatementPrint{})
	// DB.AutoMigrate(&models.PayingAuthority{})
	// DB.AutoMigrate(&models.PaBillSummary{})
	// DB.AutoMigrate(&models.ClientWork{})
	// DB.AutoMigrate(&models.Loan{})
	// DB.AutoMigrate(&models.LoanBill{})

	// DB.AutoMigrate(&models.WorkflowComments{})
	// DB.AutoMigrate(&models.WorkflowRules{})
	// DB.AutoMigrate(&models.WorkflowPolicy{})

	DB.AutoMigrate(&models.WfRequest{})
	// DB.AutoMigrate(&models.WfTask{})
	// DB.AutoMigrate(&models.WfTaskAssignment{})
	// DB.AutoMigrate(&models.WfTaskExecutionLog{})

	// DB.AutoMigrate(&models.WfTaskAssignment{})
	// DB.AutoMigrate(&models.WfActionAssignment{})
	// DB.AutoMigrate(&models.WfTaskExecutionLog{})
	DB.AutoMigrate(&models.WfAction{})
	// DB.AutoMigrate(&models.WfTask{})
	// DB.AutoMigrate(&models.WfRequest{})
	// DB.AutoMigrate(&models.UserDepartment{})
	// DB.AutoMigrate(&models.WfUserReminder{})
	DB.AutoMigrate(&models.ReqCall{})
	DB.AutoMigrate(&models.VideoProof{})
	DB.AutoMigrate(&models.ReqProof{})
	// DB.AutoMigrate(&models.WfComment{})
	// DB.AutoMigrate(&models.UserLimit{})
	// DB.AutoMigrate(&models.Annuity{})
	// DB.AutoMigrate(&models.AgtExt{})
	DB.AutoMigrate(&models.UserExt{})
	// DB.AutoMigrate(&models.PayOsBal{})
	// DB.AutoMigrate(&models.PlanLife{})
	// DB.AutoMigrate(&models.PlanLifeBenefit{})
	DB.AutoMigrate(&models.PlanLifeDiscount{})
	DB.AutoMigrate(&models.PriorPolicy{})
	DB.AutoMigrate(&models.TranReversal{})
	DB.AutoMigrate(&models.CbUser{})
	DB.AutoMigrate(&models.CbLog{})
	DB.AutoMigrate(&models.Cola{})
	DB.AutoMigrate(&models.ApClient{})
	DB.AutoMigrate(&models.ApAddress{})

}
