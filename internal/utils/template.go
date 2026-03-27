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
	EmailCertificateAr      Templates = "email-certificate-ar"
	EmailCertificateEn      Templates = "email-certificate-en"
	EmailTechnicalAr        Templates = "email-technical-ar"
	EmailTechnicalEn        Templates = "email-technical-en"
	EmailVerificationCodeEn Templates = "email-verification-code-en"
	EmailVerificationLinkEn Templates = "email-verification-link-en"

	EventCertificateAr Templates = "event-certificate-ar"
	EventCertificateEn Templates = "event-certificate-en"
)

func GetEmailTechnicalTemplate(data any) (string, error) {
	return GetTemplate("email-technical-"+data.(models.TechnicalEmailTemplate).Lang, data)
}

func GetArabicCertificateTemplate(data any) (string, error) {
	return GetTemplate("cert-ar", data)
}

func GetEnglishCertificateTemplate(data any) (string, error) {
	return GetTemplate("cert-en", data)
}

func GetTemplate(name string, data any) (string, error) {
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
