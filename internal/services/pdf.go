package services

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

type IPDFService interface {
	GeneratePDFFromHTML(ctx context.Context, html string) ([]byte, error)
}

type PDFService struct {
	workerSem chan struct{}
}

func NewPDFService(maxConcurrent uint) *PDFService {
	return &PDFService{
		workerSem: make(chan struct{}, maxConcurrent),
	}
}

func (s *PDFService) GeneratePDFFromHTML(ctx context.Context, html string) ([]byte, error) {
	s.workerSem <- struct{}{}
	defer func() { <-s.workerSem }()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "weasyprint", "-", "-")

	cmd.Stdin = bytes.NewBufferString(html)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), err
}
