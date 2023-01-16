//mail  service logic

package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer" // premailer converts css formats easily to email
	mail "github.com/xhit/go-simple-mail/v2"     //// go mailer package
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From         string //email address
	FromName     string //name associated with the email address
	To           string
	Subject      string
	Attcachments []string
	Data         any //interface{}
	DataMap      map[string]any
}

// send email
func (m *Mail) SendSMTPMessage(msg Message) error {
	//creating a valid from address and name setup
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}
	//calling a template for mail
	data := map[string]any{
		"message": msg.Data,
	}
	msg.Data = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}
	//palin text version of message
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}
	//create a server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	//establish client
	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}
	//create email message
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	//plain text msg body
	email.SetBody(mail.TextPlain, plainMessage)
	//text htmlbody
	email.AddAlternative(mail.TextHTML, formattedMessage)

	//add attachments if there are any
	if len(msg.Attcachments) > 0 {
		for _, x := range msg.Attcachments {
			email.AddAttachment(x)
		}
	}
	//sending the mail
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.html"
	//calling the html template
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	//store templates in temporary memory
	var tpl bytes.Buffer

	//describe the data associated with t that has the given name to the
	// specified data object and writes the output to wr(io.writer= &tpl)
	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)
	if err != nil {
		return "", err
	}
	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.html.html"
	//calling the html template
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	//store templates in temporary memory
	var tpl bytes.Buffer

	//describe the data associated with t that has the given name to the
	// specified data object and writes the output to wr(io.writer= &tpl)
	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)
	if err != nil {
		return "", err
	}
	formattedMessage := tpl.String()
	//inline the css
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}
	return formattedMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		// premailer converts css formats easily to email
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}
	return html, nil

}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls": //transport layer security to protect data sent over the internet by encryption
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSL //secure socket layer to establish an encryption link btwn server and client
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
