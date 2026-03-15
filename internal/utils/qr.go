package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sea-api/internal/config"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
)

func GenerateGearQR(data string, width, height int) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR matrix: %w", err)
	}

	matrix := qr.Bitmap()
	gridSize := len(matrix)

	dc := gg.NewContext(width, height)

	// Draw outer frame
	frameThickness := 5.0
	dc.SetColor(color.Black)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), 0) // corner radius 0
	dc.Fill()

	// Draw inner area
	usableWidth := float64(width) - 2*frameThickness
	usableHeight := float64(height) - 2*frameThickness
	dc.SetColor(color.White)
	dc.DrawRoundedRectangle(frameThickness, frameThickness, usableWidth, usableHeight, 0)
	dc.Fill()

	// 3. Setup drawing variables
	dc.SetHexColor("#05AFDA")

	moduleSize := usableWidth / float64(gridSize)
	offsetX := (float64(width) - float64(gridSize)*moduleSize) / 2.0
	offsetY := (float64(height) - float64(gridSize)*moduleSize) / 2.0

	// Loop through addresses and draw Gears
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if matrix[y][x] {
				xPos := offsetX + (float64(x) * moduleSize)
				yPos := offsetY + (float64(y) * moduleSize)

				if isInsideEye(x, y, gridSize) {
					dc.SetHexColor("#05AFDA")
					dc.DrawRectangle(xPos, yPos, moduleSize, moduleSize)
					dc.Fill()
				} else {
					dc.SetColor(color.Black)
					drawGear(dc, xPos, yPos, moduleSize)
				}
			}
		}
	}

	// Put logo in the middle
	logoFile, err := os.Open(config.App.ResourcesDir + "/logo.png")
	if err == nil {
		defer logoFile.Close()
		logoImg, _, err := image.Decode(logoFile)
		if err == nil {
			logoWidth := width / 5
			logoHeight := height / 5

			scaledLogo := image.NewRGBA(image.Rect(0, 0, logoWidth, logoHeight))
			draw.CatmullRom.Scale(scaledLogo, scaledLogo.Bounds(), logoImg, logoImg.Bounds(), draw.Over, nil)

			dc.DrawImageAnchored(scaledLogo, width/2, height/2, 0.5, 0.5)
		}
	}

	// Output to []byte
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("failed to encode image to PNG bytes: %w", err)
	}

	return buf.Bytes(), nil
}

func drawGear(dc *gg.Context, x, y, size float64) {
	dc.Push()      // Save context state
	defer dc.Pop() // Restore context state

	centerX := x + size/2.0
	centerY := y + size/2.0
	toothCount := 8
	toothLength := size / 6.0
	coreRadius := (size / 2.0) - (toothLength / 2.0)

	dc.Translate(centerX, centerY)

	// Draw Teeth
	angleStep := (2.0 * math.Pi) / float64(toothCount)
	for i := 0; i < toothCount; i++ {
		dc.Push()
		dc.Rotate(float64(i) * angleStep)
		dc.DrawRectangle(-(size / 8.0), -(size / 2.0), size/4.0, size)
		dc.Fill()
		dc.Pop()
	}

	// Draw Body
	dc.DrawCircle(0, 0, coreRadius)
	dc.Fill()
}

// This function is used to skip the three corner squares
func isInsideEye(x, y, gridSize int) bool {
	if x < 12 && y < 12 {
		return true
	} // Top Left
	if x > gridSize-13 && y < 12 {
		return true
	} // Top Right
	if x < 12 && y > gridSize-13 {
		return true
	} // Bottom Left
	return false
}
