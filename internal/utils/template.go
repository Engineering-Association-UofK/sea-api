package utils

import (
	"bytes"
	"os"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"text/template"
)

type Templates string

const (
	EmailCertificateAr Templates = "email-certificate-ar"
	EmailCertificateEn Templates = "email-certificate-en"
	EmailTechnicalAr   Templates = "email-technical-ar"
	EmailTechnicalEn   Templates = "email-technical-en"
	EventCertificateAr Templates = "event-certificate-ar"
	EventCertificateEn Templates = "event-certificate-en"
)

func GetEmailTechnicalTemplate(data any) (string, error) {
	return ReadFile("email-technical-"+data.(models.TechnicalEmailTemplate).Lang, data)
}

func GetArabicCertificateTemplate(data any) (string, error) {
	return ReadFile("cert-ar", data)
}

func GetEnglishCertificateTemplate(data any) (string, error) {
	return ReadFile("cert-en", data)
}

func ReadFile(name string, data any) (string, error) {
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
