package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Turns SVG data into PDF using Inkscape CLI
// Makes temporary files into disk
func SvgTo(SVG *bytes.Buffer, Type string) ([]byte, error) {
	tmpSvg, err := os.CreateTemp("", "cert-*.svg")
	if err != nil {
		return nil, fmt.Errorf("failed creating temp SVG file: %v", err)
	}
	defer os.Remove(tmpSvg.Name())

	if _, err := tmpSvg.Write(SVG.Bytes()); err != nil {
		tmpSvg.Close()
		return nil, fmt.Errorf("failed writing to temp SVG file: %v", err)
	}
	tmpSvg.Close()

	tmpPdfPath := tmpSvg.Name() + "." + Type
	defer os.Remove(tmpPdfPath)

	// Execute Inkscape command
	cmd := exec.Command("inkscape", tmpSvg.Name(), "-o", tmpPdfPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("inkscape rendering failed: %v, output: %s", err, string(out))
	}

	return os.ReadFile(tmpPdfPath)
}
