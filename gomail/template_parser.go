package gomail

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
)

func parseBodyTemplate(templateFileName string, data Data) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("config/%s", templateFileName))
	if err != nil {
		return "", errors.New("Invalid template name")
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	body := buf.String()
	return body, nil
}
