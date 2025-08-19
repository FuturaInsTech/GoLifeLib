package utilities

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/paramTypes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/valyala/fasthttp"
	"gopkg.in/gomail.v2"
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
		"contains":     strings.Contains, // <C>.<Field> STR.text.</C> Tag
		"eq": func(a, b interface{}) bool { // <E>.<Field> STR.text..</E> Tag
			return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
		},
		"ne": func(a, b interface{}) bool { // <N>.<Field> STR.text..</N> Tag
			return fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b)
		},
		"in": func(val interface{}, options ...interface{}) bool { // <I>.<Field> [S T R].text..</I> Tag
			valStr := fmt.Sprintf("%v", val)
			for _, opt := range options {
				if valStr == fmt.Sprintf("%v", opt) {
					return true
				}
			}
			return false
		},
		"out": func(val string, options ...interface{}) bool { // <O>.<Field> [S T R].text...</O> Tag
			valStr := fmt.Sprintf("%v", val)
			for _, opt := range options {
				if valStr == fmt.Sprintf("%v", opt) {
					return false
				}
			}
			return true
		},
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
func EmailTrigger(icommuncation models.Communication, itempName string, pdfData []byte, txn *gorm.DB) error {

	var client models.Client
	result := txn.First(&client, "id = ?", icommuncation.ClientID)
	if result.Error != nil {
		return fmt.Errorf("failed to read Client")
	}

	if icommuncation.EmailAllowed == "Y" {
		sender := icommuncation.CompanyEmail
		var p0033data paramTypes.P0033Data
		var extradatap0033 paramTypes.Extradata = &p0033data
		err := GetItemD(int(icommuncation.CompanyID), "P0033", itempName, icommuncation.EffectiveDate, &extradatap0033)
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

		m.Attach(icommuncation.TemplateName+"_"+fmt.Sprint(icommuncation.ClientID)+"_"+fmt.Sprint(icommuncation.PolicyID)+"_"+iDateTime+".pdf", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(pdfData)
			return err
		}))

		// Following Lines are commented.  We have made it as as async
		// d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
		// d.SSL = true

		// if err := d.DialAndSend(m); err != nil {
		// 	log.Printf("Failed to send email: %v", err)
		// 	return err
		// }

		// Configure SMTP dialer
		d := gomail.NewDialer(smtpServer, smtpPort, sender, password)
		d.SSL = true      // Enables SSL
		d.TLSConfig = nil // Use default TLS settings

		// Send email asynchronously with proper logging
		go func() {
			sendStart := time.Now()
			if err := d.DialAndSend(m); err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("Email sent successfully to %s (CC: %s, BCC: %s) in %v",
					receiver, "", "", time.Since(sendStart))
			}
		}()

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
			return fmt.Errorf("failed to read Agency")
		}
		var agclient models.Client
		result = txn.First(&agclient, "id = ?", agntenq.ClientID)
		if result.Error != nil {
			return fmt.Errorf("failed to read Client")
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
	return nil
}

func CeateFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatNumber": formatNumber,
	}
}

func FormatNumber(value interface{}, fds string) string {
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

func FormatInteger(famt string, ctype byte) string {
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

func FormatDecimal(famt string, ctype byte) string {
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

func MoveFile(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer destFile.Close()

	// Copy file content
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	// Remove the original file after copying
	if err := os.Remove(sourcePath); err != nil {
		return fmt.Errorf("error deleting source file: %w", err)
	}

	return nil
}

func GetReportforOnline(icommuncation models.Communication, itempName string, txn *gorm.DB) error {
	defaultpath := os.Getenv("REPORTPDF_SAVE_PATH")
	parts := strings.Split(icommuncation.TemplatePath, "/")
	templateFile := parts[len(parts)-1] // Extract gohtml file name

	imgFolder := strings.TrimSuffix(templateFile, "."+strings.Split(templateFile, ".")[1])

	remainingPath := strings.Join(parts[:len(parts)-1], "/")
	absolutePath, err := filepath.Abs(remainingPath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %w", err)
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
		return fmt.Errorf("error loading template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, icommuncation.ExtractedData)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Create PDF from the template execution output
	r := NewRequestPdf(buf.String())
	pdffileName := fmt.Sprintf("%s_%d_%d_%s.pdf", icommuncation.TemplateName, icommuncation.ClientID, icommuncation.PolicyID, time.Now().Format("20060102150405"))

	var pdfBuf bytes.Buffer
	success, err := r.GeneratePDFP(&pdfBuf, icommuncation.CompanyID, icommuncation.ClientID, txn)
	if err != nil || !success {
		return fmt.Errorf("error generating PDF: %w", err)
	}

	// Save the PDF to the file system if needed
	comFileName := filepath.Join(defaultpath, pdffileName)
	if icommuncation.PDFPath != "" {
		comFileName = filepath.Join(icommuncation.PDFPath, pdffileName)
	}
	comFileName = filepath.ToSlash(filepath.Clean(comFileName))

	err = os.WriteFile(comFileName, pdfBuf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error saving PDF: %w", err)
	}

	// Send email if allowed
	if icommuncation.EmailAllowed == "Y" {
		err = EmailTrigger(icommuncation, itempName, pdfBuf.Bytes(), txn)
		if err != nil {
			return fmt.Errorf("error sending email: %w", err)
		}
	}

	// Return the generated PDF buffer
	return nil
}

func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

type RequestPdf struct {
	body string
}

// parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()

	return nil
}

func (r *RequestPdf) GeneratePDFP(inputFile io.Writer, iUserco uint, iClientid uint, txn *gorm.DB) (bool, error) {

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
		return false, fmt.Errorf("failed to create PDF generator: %w", err)
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(r.body))
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Dpi.Set(300)

	// Save to temporary file
	tempFile := "temp.pdf"
	outFile, err := os.Create(tempFile)
	if err != nil {
		return false, fmt.Errorf("failed to create temp PDF file: %w", err)
	}
	defer outFile.Close()

	pdfg.SetOutput(outFile)
	err = pdfg.Create()
	if err != nil {
		return false, fmt.Errorf("PDF generation failed: %w", err)
	}

	// Step 2: Protect the PDF using Python script
	protectedFile := "protected.pdf"
	err = EncryptPDF(tempFile, protectedFile, ipassword, opassword)
	if err != nil {
		return false, fmt.Errorf("failed to protect PDF: %w", err)
	}

	// Step 3: Write the password-protected PDF to the writer
	protectedData, err := os.ReadFile(protectedFile)
	if err != nil {
		return false, fmt.Errorf("failed to read protected PDF: %w", err)
	}
	_, err = inputFile.Write(protectedData)
	if err != nil {
		return false, fmt.Errorf("failed to write protected PDF to output: %w", err)
	}

	// Cleanup temporary files
	os.Remove(tempFile)
	os.Remove(protectedFile)

	return true, nil
}

func EncryptPDF(inputFile, outputFile, userPassword, ownerPassword string) error {
	cmd := exec.Command("pdfcpu", "encrypt",
		"-upw", userPassword, // User password
		"-opw", ownerPassword, // Owner password
		"-mode", "rc4", // Encryption mode (RC4)
		"-key", "128", // Key length (128-bit)
		"-perm", "all", // Full permissions
		inputFile, outputFile, // Input and Output files
	)
	//pdfcpu encrypt -upw -opw -mode -key -perm

	// Capture errors if any
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pdfcpu encrypt error: %s", stderr.String())
	}

	fmt.Println("PDF encrypted successfully:", outputFile)
	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return encodeBase64(auth)
}

func encodeBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
func SendSMSTwilio(iCompany, iclientID uint, itempName, iEffDate string, message string, txn *gorm.DB) error {
	// Fetch client details
	var client models.Client
	result := txn.First(&client, "id = ?", iclientID)
	if result.Error != nil {
		return fmt.Errorf("failed to read Client: %v", result.Error)
	}

	var p0033data paramTypes.P0033Data
	var extradatap0033 paramTypes.Extradata = &p0033data
	err := GetItemD(int(iCompany), "P0033", itempName, iEffDate, &extradatap0033)
	if err != nil {
		return err
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
	return nil
}

func CalLoanOS(iCompany uint, iPolicy uint, iBenID uint, iLoanSeq uint, iEffectiveDate string, iLoantype string, txn *gorm.DB) (oCapAmount float64, oBilledAmt float64, oUnBilledAmt float64, oError error) {
	var polenq models.Policy
	result := txn.Find(&polenq, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return 0, 0, 0, result.Error
	}

	var benenq models.Benefit
	result = txn.Find(&benenq, "company_id = ? and id = ?", iCompany, iBenID)
	if result.Error != nil {
		return 0, 0, 0, result.Error
	}

	var loanenq []models.Loan
	if iLoantype == "" {
		result = txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_seq_number = ? and loan_status = ?", iCompany, iPolicy, iLoanSeq, "AC")
		if result.Error != nil {
			return 0, 0, 0, result.Error
		}
	} else {
		result = txn.Find(&loanenq, "company_id = ? and policy_id = ? and loan_seq_number = ? and loan_status = ? and loan_type = ?", iCompany, iPolicy, iLoanSeq, "AC", iLoantype)
		if result.Error != nil {
			return 0, 0, 0, result.Error
		}

	}

	var q0006data paramTypes.Q0006Data
	var extradata paramTypes.Extradata = &q0006data
	GetItemD(int(iCompany), "Q0006", benenq.BCoverage, benenq.BStartDate, &extradata)

	var p0072data paramTypes.P0072Data
	var extradata1 paramTypes.Extradata = &p0072data
	GetItemD(int(iCompany), "P0072", q0006data.LoanMethod, iEffectiveDate, &extradata1)

	for i := 0; i < len(loanenq); i++ {
		iIntBilled := 0.0
		iIntUnBilled := 0.0
		var loanbill []models.LoanBill
		if iLoantype == "" {
			result = txn.Find(&loanbill, "company_id = ? and policy_id = ? and loan_id = ? and billing_status = ?", iCompany, iPolicy, loanenq[i].ID, "OP")
			if result.Error != nil {
				return 0, 0, 0, result.Error
			}
		} else {
			result = txn.Find(&loanbill, "company_id = ? and policy_id = ? and loan_id = ? and billing_status = ? and loan_type = ?", iCompany, iPolicy, loanenq[i].ID, "OP", iLoantype)
			if result.Error != nil {
				return 0, 0, 0, result.Error
			}
		}
		for j := 0; j < len(loanbill); j++ {
			iIntBilled += loanbill[j].LoanIntAmount
		}

		oLoanOS := loanenq[i].LastCapAmount
		oLoanInt := loanenq[i].LoanIntRate
		_, _, _, iNoOfDays, _, _, _, _ := NoOfDays(iEffectiveDate, loanenq[i].LastIntBillDate)

		if p0072data.LoanInterestType == "C" {
			iIntUnBilled = CompoundInterestForLoan(oLoanOS, oLoanInt, loanenq[i].NextIntBillDate, loanenq[i].LastIntBillDate, loanenq[i].LastCapDate)
		} else if p0072data.LoanInterestType == "S" {
			iIntUnBilled = SimpleInterest(oLoanOS, oLoanInt, float64(iNoOfDays))
		}

		oCapAmount += iIntBilled + iIntUnBilled + oLoanOS
		oBilledAmt += iIntBilled
		oUnBilledAmt += iIntUnBilled

	}

	return RoundFloat(oCapAmount, 2), RoundFloat(oBilledAmt, 2), RoundFloat(oUnBilledAmt, 2), nil

}

// Converts 0 → A, 1 → B, ..., 25 → Z, 26 → AA, 27 → AB, ...
func ColumnIndexToName(index int) string {
	name := ""
	for index >= 0 {
		name = string(rune(index%26+'A')) + name
		index = index/26 - 1
	}
	return name
}

func CompoundInterestForLoan(iPrincipal float64, iInterest float64, iDate1 string, iDate2 string, iDate3 string) (ointrest float64) {

	iDate := GetBusinessDate(1, 0, 0)
	var iNoOfDays1 int64
	var iNoOfDays2 int64
	if iDate < iDate1 {
		_, _, _, iNoOfDays1, _, _, _, _ = NoOfDays(iDate, iDate3)
		_, _, _, iNoOfDays2, _, _, _, _ = NoOfDays(iDate2, iDate3)
	} else {
		_, _, _, iNoOfDays1, _, _, _, _ = NoOfDays(iDate1, iDate3)
		_, _, _, iNoOfDays2, _, _, _, _ = NoOfDays(iDate2, iDate3)
	}

	// Convert days to years
	oDays1 := float64(iNoOfDays1) / 365
	oDays2 := float64(iNoOfDays2) / 365

	annualRate := 1 + (iInterest / 100)

	// Compound amounts for both periods
	amount1 := iPrincipal * math.Pow(annualRate, oDays1)
	amount2 := iPrincipal * math.Pow(annualRate, oDays2)

	interest1 := amount1 - iPrincipal
	interest2 := amount2 - iPrincipal

	ointrest = interest1 - interest2
	return RoundFloat(ointrest, 2)
}

func CalculateCommision(iCompany uint, ipolicy uint, txn *gorm.DB) float64 {
	var polenq models.Policy
	txn.Find(&polenq, "id=?", ipolicy)

	var benefitenq []models.Benefit
	txn.Find(&benefitenq, "policy_id = ? and b_status =  ?", ipolicy, "IF")

	var commissionamt float64

	//Noofinstalment := GetNoIstalments(polenq.PaidToDate, polenq.BTDate, polenq.PFreq)

	Totalinstfromrcd := GetNoIstalments(polenq.PRCD, polenq.PaidToDate, polenq.PFreq)

	nomonths, _ := ConvertInstallmentsToMonths(Totalinstfromrcd, polenq.PFreq)

	for i := 0; i < len(benefitenq); i++ {

		commrate := GetCommissionRates(iCompany, benefitenq[i].BCoverage, uint(nomonths), polenq.PaidToDate)
		commission := benefitenq[i].BPrem * commrate

		commissionamt += commission
	}

	return commissionamt
}

func ConvertInstallmentsToMonths(installments int, frequency string) (int, int) {
	var monthsPerInstallment int

	switch frequency {
	case "M":
		monthsPerInstallment = 1

	case "Q":
		monthsPerInstallment = 3
	case "H":
		monthsPerInstallment = 6
	case "Y":
		monthsPerInstallment = 12

	}

	return installments * monthsPerInstallment, monthsPerInstallment
}

// calculate commission for lapse policy

func CalculateCommisionLA(iCompany uint, ipolicy uint, txn *gorm.DB) float64 {
	var polenq models.Policy
	txn.Find(&polenq, "id = ?", ipolicy)

	var benefitenq []models.Benefit
	txn.Find(&benefitenq, "policy_id = ? AND b_status = ?", ipolicy, "LA")

	// Get number of months between PRCD and PaidToDate
	currentDate := Date2String(time.Now())
	totalInstFromRCD := GetNoIstalmentsLA(polenq.PaidToDate, currentDate, polenq.PFreq)
	var commissionamt float64

	for j := 0; j < totalInstFromRCD; j++ {

		nomonths1, freqinno := ConvertInstallmentsToMonths(totalInstFromRCD, polenq.PFreq)
		noofmonths := nomonths1 + freqinno*j
		for i := 0; i < len(benefitenq); i++ {

			if polenq.PaidToDate <= benefitenq[i].BPremCessDate {
				commrate := GetCommissionRates(iCompany, benefitenq[i].BCoverage, uint(noofmonths), polenq.PaidToDate)
				commission := benefitenq[i].BPrem * commrate
				commissionamt += commission
			} else {

				return 0
			}
		}
	}
	return commissionamt
}

func GetNoIstalmentsLA(iFromDate, iToDate, iFrequency string) (oInstalments int) {

	fromDate := String2Date(iFromDate)
	toDate := String2Date(iToDate)
	method := "M"
	var noinstalments float64

	year1, month1, _, _, _, _ := DateDiff(fromDate, toDate, method)
	fmt.Println("Shubham", iFromDate, iToDate, iFrequency, year1, month1)
	var dd1 string
	var dd2 string
	dd1 = iFromDate[6:8] // DD paid to date
	dd2 = iToDate[6:8]   //DD current date
	switch {

	case iFrequency == "Y":
		// 10 and 0  10
		noinstalments = (float64(year1) / 1)
		noinst1 := float32((year1 * 1) + (month1))
		noinst2 := int(noinst1)
		if noinst1 == float32(noinst2) {

			if dd2 <= dd1 {
				noinstalments = noinstalments - 1
			}
		}
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments
	case iFrequency == "H":
		// 10 and 6  10*2 + 6/6 =  21
		noinstalments := float32((year1 * 2) + (month1 / 6))
		noinst1 := float32((year1 * 2) + (month1))
		noinst2 := int(noinst1)
		if noinst1 == float32(noinst2) {

			if dd2 <= dd1 {
				noinstalments = noinstalments - 1
			}
		}
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments

	case iFrequency == "Q":
		// 5 9   = 5 * 4  + 9/3    20 + 3
		noinstalments := float32((year1 * 4) + (month1 / 3))
		noinst1 := float32((year1 * 4) + (month1))

		noinst2 := int(noinst1)
		if noinst1 == float32(noinst2) {

			if dd2 <= dd1 {
				noinstalments = noinstalments - 1
			}
		}
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments
	case iFrequency == "M":
		noinstalments := float32((year1 * 12) + (month1))
		noinst1 := float32((year1 * 12) + (month1))
		noinst2 := int(noinst1)
		if noinst1 == float32(noinst2) {

			if dd2 <= dd1 {
				noinstalments = noinstalments - 1
			}
		}
		noinstalments = noinstalments + 1
		oInstalments = int(noinstalments)
		return oInstalments
	case iFrequency == "S":
		oInstalments := 1
		return oInstalments

	}
	return
}

func TDFCOLAD(iCompany uint, iPolicy uint, iFunction string, iTranno uint, iRevFlag string, txn *gorm.DB) (string, error) {
	var policy models.Policy
	var tdfpolicy models.TDFPolicy
	var tdfrule models.TDFRule

	var q0005data paramTypes.Q0005Data
	var extradataq0005 paramTypes.Extradata = &q0005data

	result := txn.First(&tdfrule, "company_id = ? and tdf_type = ?", iCompany, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
	}
	result = txn.First(&policy, "company_id = ? and id = ?", iCompany, iPolicy)
	if result.Error != nil {
		return "", result.Error
	}

	err := GetItemD(int(iCompany), "Q0005", policy.PProduct, policy.PRCD, &extradataq0005)
	if err != nil {
		txn.Rollback()
		return "", err
	}

	iNewDate := AddLeadDays(policy.AnnivDate, (-1 * q0005data.BillingLeadDays))

	results := txn.First(&tdfpolicy, "company_id = ? and policy_id = ? and tdf_type = ?", iCompany, iPolicy, iFunction)
	if result.Error != nil {
		txn.Rollback()
		return "", result.Error
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
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	} else {
		result = txn.Delete(&tdfpolicy)
		if result.Error != nil {
			txn.Rollback()
			return "", result.Error
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
			txn.Rollback()
			return "", result.Error
		}

		return "", nil
	}
}

func RoundUpToNext1000(n float64) float64 {
	return math.Ceil(n/1000.0) * 1000
}

func RoundToNearest100(val float64) float64 {
	return math.Round(val/100.0) * 100.0
}

func DiscountCalculationForCola(iCompany uint, iNewSA float64, iNewAnnPrem float64, iDiscType string, iDiscMethod string, iDate string, prevAnnPrem *int, prevBandIndex *int) (float64, error) {

	var oValue float64

	switch iDiscType {
	case "S":
		// Fetch Sum Assured bands
		var q0017data paramTypes.Q0017Data
		var extradata paramTypes.Extradata = &q0017data
		if err := GetItemD(int(iCompany), "Q0017", iDiscMethod, iDate, &extradata); err != nil {
			return 0.0, err
		}

		currentSA := int(iNewSA)
		currentPrem := int(iNewAnnPrem)
		var oDiscount uint
		currentBandIndex := -1

		// Find applicable SA band
		for i := len(q0017data.SaBand) - 1; i >= 0; i-- {
			if currentSA >= int(q0017data.SaBand[i].Sa) {
				currentBandIndex = i
				break
			}
		}

		if currentBandIndex != -1 {
			if *prevBandIndex != currentBandIndex {
				// First time in this band → full discount
				oDiscount = uint(q0017data.SaBand[currentBandIndex].Discount) * uint(currentPrem) / 100
			} else {
				// Same band → incremental discount
				diff := currentPrem - *prevAnnPrem
				if diff > 0 {
					oDiscount = uint(q0017data.SaBand[currentBandIndex].Discount) * uint(diff) / 100
				}
			}
			*prevAnnPrem = currentPrem
			*prevBandIndex = currentBandIndex
		}

		return float64(oDiscount), nil

	case "P":
		// Fetch Premium bands
		var q0018data paramTypes.Q0018Data
		var extradata paramTypes.Extradata = &q0018data
		if err := GetItemD(int(iCompany), "Q0018", iDiscMethod, iDate, &extradata); err != nil {
			return 0.0, err
		}

		currentPrem := int(iNewAnnPrem)
		var oDiscount uint
		currentBandIndex := -1

		// Find applicable Premium band
		for i := len(q0018data.PremBand) - 1; i >= 0; i-- {
			if currentPrem >= int(q0018data.PremBand[i].AnnPrem) {
				currentBandIndex = i
				break
			}
		}

		if currentBandIndex != -1 {
			if *prevBandIndex != currentBandIndex {
				oDiscount = uint(q0018data.PremBand[currentBandIndex].Discount) * uint(currentPrem) / 100
			} else {
				diff := currentPrem - *prevAnnPrem
				if diff > 0 {
					oDiscount = uint(q0018data.PremBand[currentBandIndex].Discount) * uint(diff) / 100
				}
			}
			*prevAnnPrem = currentPrem
			*prevBandIndex = currentBandIndex
		}

		return float64(oDiscount), nil
	}

	return oValue, nil
}

func CalculateFutureCommission(iCompany uint, ipolicy uint, ifuturedate string, txn *gorm.DB) float64 {
	var polenq models.Policy
	txn.Find(&polenq, "id = ?", ipolicy)

	var benefitenq []models.Benefit
	txn.Find(&benefitenq, "policy_id = ? AND b_status = ?", ipolicy, "IF")

	// Get number of months between PRCD and PaidToDate
	totalInstFromPTD := GetNoIstalmentsLA(polenq.PaidToDate, ifuturedate, polenq.PFreq)
	var commissionamt float64

	for j := 0; j < totalInstFromPTD; j++ {

		nomonths1, freqinno := ConvertInstallmentsToMonths(totalInstFromPTD, polenq.PFreq)
		noofmonths := nomonths1 + freqinno*j
		for i := 0; i < len(benefitenq); i++ {

			if polenq.PaidToDate <= benefitenq[i].BPremCessDate {
				commrate := GetCommissionRates(iCompany, benefitenq[i].BCoverage, uint(noofmonths), polenq.PaidToDate)
				commission := benefitenq[i].BPrem * commrate
				commissionamt += commission
			} else {

				return 0
			}
		}
	}
	return commissionamt
}
