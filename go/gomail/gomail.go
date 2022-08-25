package gomail

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Data struct {
	ComponentName string
	Subject       string
	Body          string
	EmailTo       string
}

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  "",
		RefreshToken: "",
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail service: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		log.Println("Email service is initialized")
	}
}

func SendEmailOAuth2(data Data) (bool, error) {
	var message gmail.Message

	emailTo := "To: " + data.EmailTo + "\r\n"
	subject := "Subject: " + data.Subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + data.Body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateBody(data Data) (string, error) {
	body, err := parseBodyTemplate("template.jinja", data)
	if err != nil {
		return "", errors.New("unable to parse email template")
	}

	return body, nil
}

func SendEmailOAuth2WithAttachment(data Data, fileNames []string) (bool, error) {
	var message gmail.Message

	boundary := randStr(32, "alphanum")

	var messageBody = []byte(
		"Content-Type: multipart/mixed; boundary=" + boundary + "\n" +
			"MIME-Version: 1.0\n" +
			"to: " + data.EmailTo + "\n" +
			"subject: " + data.Subject + "\n\n" +

			"--" + boundary + "\n" +
			"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
			"MIME-Version: 1.0\n" +
			"Content-Transfer-Encoding: 7bit\n\n" +
			data.Body + "\n\n")

	attachments := AddAttachments(boundary, fileNames)
	messageBody = append(messageBody, attachments...)

	endMessage := []byte("--" + boundary + "--")
	messageBody = append(messageBody, endMessage...)

	message.Raw = base64.URLEncoding.EncodeToString(messageBody)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		log.Println("Message sent!")
	}

	return true, nil
}

func AddAttachments(boundary string, files []string) []byte {
	mimeBody := []byte("")

	for _, file := range files {
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fileMIMEType := http.DetectContentType(fileBytes)

		fileData := base64.StdEncoding.EncodeToString(fileBytes)

		filename := filepath.Base(file)

		var attachment = []byte("--" + boundary + "\n" +

			"Content-Type: " + fileMIMEType + "; name=" + string('"') + filename + string('"') + " \n" +
			"MIME-Version: 1.0\n" +
			"Content-Transfer-Encoding: base64\n" +
			"Content-Disposition: attachment; filename=" + string('"') + filename + string('"') + " \n\n" +

			chunkSplit(fileData, 76, "\n"))

		mimeBody = append(mimeBody, attachment...)
	}

	return mimeBody
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for k, v := range strBytes {
		strBytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(strBytes)
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}
