package certservice

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"text/template"
	"time"
	"unicode"

	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
	"sea-api/internal/utils"
)

// Mock function to represent your byte-fetching logic for overlay images
func fetchAssetBytes(filepath string) []byte {
	b, _ := os.ReadFile(filepath)
	return b
}

func (s *CertService) GenerateTestImage(
	req *certmodels.IssueRequest,
	ctx context.Context,
) (string, error) {
	certTemplate, err := s.repo.GetTemplateByID(req.TemplateID)
	if err != nil {
		return "", fmt.Errorf("failed getting template raw: %v", err)
	}

	// Create unique hash cor the certificate
	qr, err := utils.GenerateGearQR("This is a test QR", 512, 512)
	if err != nil {
		return "", fmt.Errorf("failed generating QR: %v", err)
	}

	// Create certificate data
	tmplData, err := fillTemplate(req, certTemplate, qr)
	if err != nil {
		return "", fmt.Errorf("failed filling template struct: %v", err)
	}

	// Parse and Execute SVG Template
	filledSVG, err := parseTemplate(tmplData, certTemplate.Language, certTemplate.Version)
	if err != nil {
		return "", fmt.Errorf("failed parsing template: %v", err)
	}

	// Generate PNG file
	png, err := utils.SvgTo(filledSVG, "png")
	if err != nil {
		return "", fmt.Errorf("failed reading Inkscape PNG output: %v", err)
	}

	return s.s3.UploadTemp(ctx, png, "image/png", ",png")
}

func (s *CertService) GeneratePdf(
	req *certmodels.IssueRequest,
) ([]byte, error) {
	certTemplate, err := s.repo.GetTemplateByID(req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed getting template raw: %v", err)
	}

	// Create unique hash cor the certificate
	stringToHash := fmt.Sprintf(`%s|%d|%s|%s`, req.RecipientName, req.TemplateID, time.Now().Format("02-01-2006"), config.App.SecretSalt)
	hash := sha256.Sum256([]byte(stringToHash))
	hashString := hex.EncodeToString(hash[:])

	url := config.Links.CertVerify + "/" + hashString
	qr, err := utils.GenerateGearQR(url, 512, 512)
	if err != nil {
		return nil, fmt.Errorf("failed generating QR: %v", err)
	}

	// Create certificate data
	tmplData, err := fillTemplate(req, certTemplate, qr)
	if err != nil {
		return nil, fmt.Errorf("failed filling template struct: %v", err)
	}

	// Parse and Execute SVG Template
	filledSVG, err := parseTemplate(tmplData, certTemplate.Language, certTemplate.Version)
	if err != nil {
		return nil, fmt.Errorf("failed parsing template: %v", err)
	}

	// Generate PDF file
	PdfFileBytes, err := utils.SvgTo(filledSVG, "pdf")
	if err != nil {
		return nil, fmt.Errorf("failed reading Inkscape PNG output: %v", err)
	}

	return PdfFileBytes, nil
}

func parseTemplate(tmplData *certmodels.V0_1, lang models.Language, version string) (*bytes.Buffer, error) {
	path := fmt.Sprintf("%s/static-assets/certificates/sea-certificate.%s.%s.svg", config.App.ResourcesDir, lang, version)
	slog.Debug("Generating certificate path", "path", path)
	svgTemplateBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading SVG: %v", err)
	}
	slog.Debug("Got SVG File", "file", string(svgTemplateBytes[0:100]))

	tmpl, err := template.New("cert").Parse(string(svgTemplateBytes))
	if err != nil {
		return nil, fmt.Errorf("failed Creating template object: %v", err)
	}

	var filledSVG bytes.Buffer
	if err := tmpl.Execute(&filledSVG, tmplData); err != nil {
		return nil, fmt.Errorf("failed executing template: %v", err)
	}

	return &filledSVG, nil
}

func fillTemplate(req *certmodels.IssueRequest, template *certmodels.CertificateTemplate, qr []byte) (*certmodels.V0_1, error) {
	var layout certmodels.Layout
	err := json.Unmarshal(template.LayoutConfig, &layout)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling layout config: %v", err)
	}

	sign1, err := req.SignerSignatureOne.Open()
	if err != nil {
		return nil, fmt.Errorf("failed opening first signature: %v", err)
	}
	defer sign1.Close()

	signature1, err := io.ReadAll(sign1)
	if err != nil {
		return nil, fmt.Errorf("failed reading first signature: %v", err)
	}

	sign2, err := req.SignerSignatureTwo.Open()
	if err != nil {
		return nil, fmt.Errorf("failed opening second signature: %v", err)
	}
	defer sign2.Close()

	signature2, err := io.ReadAll(sign2)
	if err != nil {
		return nil, fmt.Errorf("failed reading second signature: %v", err)
	}

	stamp, err := os.ReadFile(fmt.Sprintf(`%s/secrets/stamp.png`, config.App.ResourcesDir))
	if err != nil {
		return nil, fmt.Errorf("failed reading stamp file: %v", err)
	}

	return &certmodels.V0_1{
		Name:      req.RecipientName,
		Title:     layout.Title,
		Subtitle:  layout.Subtitle,
		Statement: layout.Statement,

		CollabNameOne: req.SignerNameOne,
		CollabRoleOne: req.SignerRoleOne,
		CollabNameTwo: req.SignerNameTwo,
		CollabRoleTwo: req.SignerRoleTwo,

		SignOneBase64: base64.StdEncoding.EncodeToString(signature1),
		SignTwoBase64: base64.StdEncoding.EncodeToString(signature2),

		QRBase64:    base64.RawStdEncoding.EncodeToString(qr),
		StampBase64: base64.RawStdEncoding.EncodeToString(stamp),
	}, nil
}

// isArabic checks if the string contains characters from the Arabic Unicode block.
func isArabic(text string) bool {
	for _, r := range text {
		if unicode.In(r, unicode.Arabic) {
			return true
		}
	}
	return false
}
