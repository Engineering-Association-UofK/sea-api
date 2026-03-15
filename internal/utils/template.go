package utils

import (
	"bytes"
	"os"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"text/template"
)

func GetEmailTechnicalTemplate(data any) (string, error) {
	return readFile("email-technical-"+data.(models.TechnicalEmailTemplate).Lang, data)
}

func GetCertificateEmailTemplate(data any) (string, error) {
	return readFile("email-technical", data)
}

func GetArabicCertificateTemplate(data any) (string, error) {
	return readFile("cert-ar", data)
}

func GetEnglishCertificateTemplate(data any) (string, error) {
	return readFile("cert-en", data)
}

func readFile(name string, data any) (string, error) {
	content, err := os.ReadFile(config.App.ResourcesDir + "/templates/" + name + ".html")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
