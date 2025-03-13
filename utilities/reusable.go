package utilities

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"
	"gorm.io/gorm"
)

func FormatIndianNumber(amount float64) string {
	// Convert float to string with 2 decimal places
	formatted := fmt.Sprintf("%.2f", amount)

	// Split into whole and decimal parts
	parts := strings.Split(formatted, ".")
	whole := parts[0]
	decimal := parts[1]

	// Apply Indian grouping to the whole part
	n := len(whole)
	if n <= 3 {
		return whole + "." + decimal // No formatting needed for small numbers
	}

	// Last three digits remain unchanged, process the rest
	indianFormatted := whole[n-3:] // Take last 3 digits
	whole = whole[:n-3]            // Remaining digits

	// Insert commas every two digits from the right
	for len(whole) > 2 {
		indianFormatted = whole[len(whole)-2:] + "," + indianFormatted
		whole = whole[:len(whole)-2]
	}

	// Append the remaining part
	if len(whole) > 0 {
		indianFormatted = whole + "," + indianFormatted
	}

	return indianFormatted + "." + decimal
}

func ConvertYYYYMMDD(inputDate string) (string, error) {
	iDate, err := time.Parse("20060102", inputDate)
	if err != nil {
		return "", err
	}
	outputDateStr := iDate.Format("02/01/2006")

	return outputDateStr, nil
}

func CreateFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatNumber": formatNumber,
	}
}
func formatInteger(famt string, ctype byte) string {
	// Handle the Indian and Western numbering formats
	n := len(famt)
	if n <= 3 {
		// No formatting needed for numbers with 3 or fewer digits
		return famt
	}

	// Format the first three digits (thousands)
	result := famt[n-3:]

	if ctype == 'a' || ctype == 'A' {
		// Format the remaining digits in groups of two (lakhs, crores, etc.)
		for i := n - 3; i > 0; i -= 2 {
			start := i - 2
			if start < 0 {
				start = 0
			}
			result = famt[start:i] + "," + result
		}
	}
	if ctype == 'c' || ctype == 'C' {
		// Format the remaining digits in groups of three (millions, billions, etc.)
		for i := n - 3; i > 0; i -= 3 {
			start := i - 3
			if start < 0 {
				start = 0
			}
			result = famt[start:i] + "," + result
		}
	}

	// Return the formatted integer
	return result
}

func formatDecimal(famt string, ctype byte) string {
	decimalpart := ""
	decimallen, _ := strconv.Atoi(string(ctype))
	//damt, _ := strconv.Atoi(string(famt))
	var decimalFloat float64 = 0
	if decimallen == 0 {
		decimalFloat, _ = strconv.ParseFloat("0."+"0", 64)
	} else {
		decimalFloat, _ = strconv.ParseFloat("0."+string(famt), 64)
	}
	if decimallen == 0 {
		return decimalpart
	} else {
		formatString := fmt.Sprintf("%%.%df", decimallen) // Generate format string
		decimalpart = fmt.Sprintf(formatString, decimalFloat)
		decimalpart = decimalpart[1:] // Remove the leading "0"
		return decimalpart
	}
}

func formatNumber(value interface{}, fds string) string {
	integralpart := ""
	decimalpart := ""
	famt := ""
	famount := ""
	decimallen, _ := strconv.Atoi(string(fds[1]))
	if _, ok := value.(float64); ok {
		famt = strconv.FormatFloat(value.(float64), 'f', decimallen, 64) // Format float64 with d decimal places
	}
	if _, ok := value.(int); ok {
		famt = strconv.FormatFloat(float64(value.(int)), 'f', decimallen, 64) // Format float64 with d decimal places
	}
	if _, ok := value.(string); ok {
		famt = value.(string)
	}
	parts := strings.Split(famt, ".")
	if fds[0] == 'a' || fds[0] == 'A' {
		integralpart = formatInteger(parts[0], fds[0])
	} else if fds[0] == 'c' || fds[0] == 'C' {
		integralpart = formatInteger(parts[0], fds[0])
	} else if fds[0] == 'd' || fds[0] == 'D' {
		integralpart = parts[0]
	} else {
		integralpart = parts[0]
	}

	if len(parts) > 1 {
		decimalpart = formatDecimal(parts[1], fds[1])
		famount = integralpart + decimalpart
	} else {
		famount = integralpart
	}
	return fmt.Sprintf("%s", famount) // Format float64 with 2 decimal places and comma separators
}

// This Function has to be used when an email is trigger through communication (Online)
func EmailTrigger(icommuncationId, itempName string, pdfData []byte, txn *gorm.DB) error {
	var communication models.Communication
	result := txn.First(&communication, "id = ? and template_name = ?", icommuncationId, itempName)
	if result.Error != nil {
		return fmt.Errorf("failed to read communication")
	}

	var client models.Client
	result = txn.First(&client, "id = ?", communication.ClientID)
	if result.Error != nil {
		return fmt.Errorf("failed to read Client")
	}

	if communication.EmailAllowed == "Y" {
		sender := communication.CompanyEmail
		var p0033data paramTypes.P0033Data
		var extradatap0033 paramTypes.Extradata = &p0033data
		err := GetItemD(int(communication.CompanyID), "P0033", itempName, communication.EffectiveDate, &extradatap0033)
		if err != nil {
			return err
		}
		receiver := client.ClientEmail
		password := p0033data.SenderPassword
		smtpServer := p0033data.SMTPServer
		smtpPort := p0033data.SMTPPort
		emailBody := p0033data.Body
		m := gomail.NewMessage()
		m.SetHeader("From", sender)
		m.SetHeader("To", receiver)
		m.SetHeader("Subject", p0033data.Subject)
		m.SetBody("text/plain", emailBody)
		iDateTime := time.Now().Format("20060102150405")

		m.Attach(communication.TemplateName+iDateTime+".pdf", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(pdfData)
			return err
		}))

		d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
		d.SSL = true

		if err := d.DialAndSend(m); err != nil {
			log.Printf("Failed to send email: %v", err)
			return err
		}

		log.Println("Email sent successfully with attachment via office SMTP")
		return nil
	}
	return nil
}

// This Function has to be used when  an email is triggered through communication.  Especially fpr Batch
func EmailTriggerN(icommunication models.Communication, pdfData []byte, txn *gorm.DB) error {

	var client models.Client
	result := txn.First(&client, "id = ?", icommunication.ClientID)
	if result.Error != nil {
		return fmt.Errorf("failed to read Client")
	}
	if client.ClientEmail == "" {
		return fmt.Errorf("Email is not Found")
	}
	iTemplate := icommunication.TemplateName
	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	err := GetItemD(int(icommunication.CompanyID), "P0033", iTemplate, icommunication.EffectiveDate, &extradatap0033)
	if err != nil {
		return err
	}
	sender := icommunication.CompanyEmail
	receiver := client.ClientEmail
	password := p0033data.SenderPassword
	smtpServer := p0033data.SMTPServer
	smtpPort := p0033data.SMTPPort
	emailBody := p0033data.Body
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", p0033data.Subject)
	m.SetBody("text/plain", emailBody)
	iDateTime := time.Now().Format("20060102150405")

	m.Attach(icommunication.TemplateName+iDateTime+".pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfData)
		return err
	}))

	d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	// Send Mail to Agent if it is configured in Communication (P0033)
	if icommunication.AgentEmailAllowed == "Y" {
		var agntenq models.Agency
		result := txn.First(&agntenq, "id = ?", icommunication.AgencyID)
		if result.Error != nil {
			return fmt.Errorf("failed to read Agency")
		}
		var agclient models.Client
		result = txn.First(&agclient, "id = ?", agntenq.ClientID)
		if result.Error != nil {
			return fmt.Errorf("failed to read Client")
		}

		if agclient.ClientEmail != "" {
			sender := icommunication.CompanyEmail

			receiver := agclient.ClientEmail
			iName := GetName(client.CompanyID, client.ID)
			emailBody := fmt.Sprintf(
				"Hi Sir/Madam,\n\n"+
					"Following Email was sent to your Customer %d %s\n\n"+
					"I am from Futura Instech..\n\n"+
					"Thank you!",
				client.ID, iName,
			)

			m := gomail.NewMessage()
			m.SetHeader("From", sender)
			m.SetHeader("To", receiver)
			m.SetHeader("Subject", "Mail Sent to Your Customer")
			m.SetBody("text/plain", emailBody)
			d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
			d.SSL = true

			if err := d.DialAndSend(m); err != nil {
				log.Printf("Failed to send email: %v", err)
				return err
			}

			log.Println("Email sent successfully with attachment via office SMTP")
			return nil
		}

	}
	return nil
}

// This function has to be used when an email is triggered with out Communication. Especially for Batch Reports
func EmailTriggerforReport(iCompany uint, iReference uint, iClient uint, iEmail string, iEffDate string, itempName string, pdfData []byte, txn *gorm.DB) error {

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	err := GetItemD(int(iCompany), "P0033", itempName, iEffDate, &extradatap0033)
	if err != nil {
		return err
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

	d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully with attachment via office SMTP")
	return nil
}
