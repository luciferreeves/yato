package lib

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/draw"
)

type ImageRenderer struct {
	method string
}

func NewImageRenderer() *ImageRenderer {
	method := determineRenderMethod()
	return &ImageRenderer{method: method}
}

func determineRenderMethod() string {
	if os.Getenv("TERM") == "xterm-kitty" {
		return "kitty"
	} else if os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		return "iterm2"
	} else if os.Getenv("TERM") == "xterm-256color" && os.Getenv("VTE_VERSION") != "" {
		return "sixel"
	}
	return "none"
}

func (r *ImageRenderer) RenderImage(img image.Image, width, height int) string {
	switch r.method {
	case "kitty":
		return r.renderKitty(img, width, height)
	case "iterm2":
		return r.renderITerm2(img, width, height)
	case "sixel":
		return r.renderSixel(img, width, height)
	case "ascii":
		return r.renderASCII(img, width, height)
	default:
		return ""
	}
}

func (r *ImageRenderer) renderKitty(img image.Image, width, height int) string {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	png.Encode(&buf, resized)
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Split the encoded data into chunks
	const chunkSize = 4096
	chunks := make([]string, 0, (len(encoded)+chunkSize-1)/chunkSize)
	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunks = append(chunks, encoded[i:end])
	}

	// Build the Kitty graphics protocol command
	var result strings.Builder
	for i, chunk := range chunks {
		if i == 0 {
			result.WriteString(fmt.Sprintf("\033_Ga=T,f=100,s=%d,v=%d,m=1;", width, height))
		} else {
			result.WriteString("\033_Gm=1;")
		}
		result.WriteString(chunk)
		result.WriteString("\033\\")
	}

	// Final chunk
	result.WriteString("\033_Gm=0;\033\\")

	return result.String()
}

func (r *ImageRenderer) renderITerm2(img image.Image, width, height int) string {
	// Implement iTerm2 inline image protocol
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("\033]1337;File=inline=1;width=%dpx;height=%dpx:%s\a", width, height, encoded)
}

func (r *ImageRenderer) renderSixel(img image.Image, width, height int) string {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)

	// Convert to Sixel
	var sb strings.Builder
	sb.WriteString("\033Pq") // Start Sixel sequence
	sb.WriteString("\"1;1;") // Set color mode and aspect ratio
	sb.WriteString(fmt.Sprintf("%d;%d", width, height))
	sb.WriteString("\n")

	// Simple color quantization (this can be improved)
	palette := make(map[color.Color]int)
	colorIndex := 0

	for y := 0; y < height; y++ {
		sixelRow := make([]int, width)
		for x := 0; x < width; x++ {
			c := resized.At(x, y)
			if _, exists := palette[c]; !exists {
				palette[c] = colorIndex
				colorIndex++
				r, g, b, _ := c.RGBA()
				sb.WriteString(fmt.Sprintf("#%d;2;%d;%d;%d", palette[c], r>>8, g>>8, b>>8))
			}
			sixelRow[x] = palette[c]
		}

		// Encode sixel data
		for i := 0; i < 6; i++ {
			for _, colorIdx := range sixelRow {
				sb.WriteByte(byte('?' + ((colorIdx >> i) & 1)))
			}
			sb.WriteByte('-')
		}
		sb.WriteByte('\n')
	}

	sb.WriteString("\033\\") // End Sixel sequence
	return sb.String()
}

func (r *ImageRenderer) renderASCII(img image.Image, width, height int) string {
	// Implement a simple ASCII art renderer
	// This is a very basic implementation and can be improved
	bounds := img.Bounds()
	ascii := ""
	for y := bounds.Min.Y; y < bounds.Max.Y; y += height / 10 {
		for x := bounds.Min.X; x < bounds.Max.X; x += width / 20 {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			avg := (r + g + b) / 3
			if avg > 32768 {
				ascii += " "
			} else {
				ascii += "#"
			}
		}
		ascii += "\n"
	}
	return ascii
}
