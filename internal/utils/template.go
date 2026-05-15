package utils

import (
	"bytes"
	"fmt"
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

func GetCertTemplate(certType models.CertType, certVersion string, lang models.Language, data any) (string, error) {
	fileName := fmt.Sprintf("%s_%s", certType, certVersion)
	path := fmt.Sprintf(
		"%s/static-assets/certificates/%s/%s.html",
		config.App.ResourcesDir,
		lang,
		fileName,
	)

	return getTemplate(path, fileName, data)
}

func GetEmailTemplate(emailType models.EmailType, lang models.Language, data any) (string, error) {
	path := fmt.Sprintf(
		"%s/static-assets/emails/%s/%s.html",
		config.App.ResourcesDir,
		lang,
		emailType,
	)

	return getTemplate(path, string(emailType), data)
}

func getTemplate(path string, name string, data any) (string, error) {
	content, err := os.ReadFile(path)
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
