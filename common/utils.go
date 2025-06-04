package common

import (
	"github.com/charmbracelet/lipgloss"
)

// Lerp performs linear interpolation
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Clamp constrains a value between min and max
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Map maps a value from one range to another
func Map(value, inMin, inMax, outMin, outMax float64) float64 {
	return (value-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// GetWaveChar returns a Unicode character for wave visualization
func GetWaveChar(height float64) string {
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	index := int(Clamp(height*float64(len(chars)), 0, float64(len(chars)-1)))
	return chars[index]
}

// GenerateGradient creates a gradient between two colors
func GenerateGradient(steps int) []lipgloss.Color {
	gradient := make([]lipgloss.Color, steps)
	for i := range gradient {
		gradient[i] = lipgloss.Color(GradientBlue[i%len(GradientBlue)])
	}
	return gradient
}