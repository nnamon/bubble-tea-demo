package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type complex128 struct {
	real, imag float64
}

func (c complex128) add(other complex128) complex128 {
	return complex128{c.real + other.real, c.imag + other.imag}
}

func (c complex128) mul(other complex128) complex128 {
	return complex128{
		c.real*other.real - c.imag*other.imag,
		c.real*other.imag + c.imag*other.real,
	}
}

func (c complex128) abs() float64 {
	return math.Sqrt(c.real*c.real + c.imag*c.imag)
}

type model struct {
	width      int
	height     int
	centerX    float64
	centerY    float64
	zoom       float64
	maxIter    int
	autoZoom   bool
	paused     bool
	zoomTarget complex128
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/15, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width:      80,
		height:     24,
		centerX:    -0.75,
		centerY:    0.1,
		zoom:       1.0,
		maxIter:    80,
		autoZoom:   true,
		zoomTarget: complex128{-0.7463, 0.1102}, // Interesting zoom point on boundary
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 4
		return m, nil

	case tickMsg:
		if !m.paused && m.autoZoom {
			// Gradually zoom into the target point
			m.zoom *= 1.03
			// Gradually move toward the zoom target
			factor := 0.01
			m.centerX += (m.zoomTarget.real - m.centerX) * factor
			m.centerY += (m.zoomTarget.imag - m.centerY) * factor
			
			// Increase iterations as we zoom deeper for more detail
			if m.zoom > 100 && m.maxIter < 150 {
				m.maxIter++
			}
			
			// Reset if zoom gets too high
			if m.zoom > 1e15 {
				m.zoom = 1.0
				m.centerX = -0.75
				m.centerY = 0.1
				m.maxIter = 80
			}
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "a":
			m.autoZoom = !m.autoZoom
		case "r":
			m.centerX = -0.75
			m.centerY = 0.1
			m.zoom = 1.0
			m.maxIter = 80
		case "up":
			if !m.autoZoom {
				m.centerY -= 0.1 / m.zoom
			}
		case "down":
			if !m.autoZoom {
				m.centerY += 0.1 / m.zoom
			}
		case "left":
			if !m.autoZoom {
				m.centerX -= 0.1 / m.zoom
			}
		case "right":
			if !m.autoZoom {
				m.centerX += 0.1 / m.zoom
			}
		case "+", "=":
			if !m.autoZoom {
				m.zoom *= 1.2
			}
		case "-":
			if !m.autoZoom {
				m.zoom /= 1.2
				if m.zoom < 0.1 {
					m.zoom = 0.1
				}
			}
		case "1":
			// Interesting boundary area with spirals
			m.zoomTarget = complex128{-0.7463, 0.1102}
			m.centerX = -0.75
			m.centerY = 0.1
			m.zoom = 1.0
			m.maxIter = 80
		case "2":
			// Edge of the main bulb
			m.zoomTarget = complex128{-0.16, 1.0405}
			m.centerX = -0.2
			m.centerY = 1.0
			m.zoom = 1.0
			m.maxIter = 80
		case "3":
			// Seahorse valley
			m.zoomTarget = complex128{-0.74529, 0.11307}
			m.centerX = -0.75
			m.centerY = 0.11
			m.zoom = 1.0
			m.maxIter = 80
		case "4":
			// Feather location
			m.zoomTarget = complex128{-0.235125, 0.827215}
			m.centerX = -0.24
			m.centerY = 0.83
			m.zoom = 1.0
			m.maxIter = 80
		case "i":
			m.maxIter = min(m.maxIter+10, 200)
		case "d":
			m.maxIter = max(m.maxIter-10, 20)
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#663399")).
		Padding(0, 1)

	title := titleStyle.Render("🌀 Mandelbrot Fractal Zoom")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Purple)
	status := statusStyle.Render(fmt.Sprintf(
		"Center: (%.6f, %.6f) | Zoom: %.2e | Iterations: %d | %s | %s",
		m.centerX, m.centerY, m.zoom, m.maxIter,
		map[bool]string{true: "Auto-zooming", false: "Manual control"}[m.autoZoom],
		map[bool]string{true: "⏸ Paused", false: "🌀 Exploring"}[m.paused],
	))

	// Render fractal
	lines := m.renderMandelbrot()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	var help string
	if m.autoZoom {
		help = "[a] manual • [1-4] targets • [i/d] iterations • [space] pause • [r]eset • [q]uit"
	} else {
		help = "[a] auto-zoom • [↑↓←→] move • [+/-] zoom • [1-4] targets • [i/d] iterations • [r]eset • [q]uit"
	}

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), helpStyle.Render(help))
}

func (m model) renderMandelbrot() []string {
	lines := make([]string, m.height)
	
	// Calculate the complex plane bounds
	aspect := float64(m.width) / float64(m.height) * 2.0 // Adjust for character aspect ratio
	scale := 3.0 / m.zoom
	
	minX := m.centerX - scale*aspect/2
	maxX := m.centerX + scale*aspect/2
	minY := m.centerY - scale/2
	maxY := m.centerY + scale/2
	
	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			// Map pixel to complex plane
			cx := minX + float64(x)*(maxX-minX)/float64(m.width)
			cy := maxY - float64(y)*(maxY-minY)/float64(m.height) // Flip Y axis
			
			// Calculate iterations for this point
			iterations := m.mandelbrotIterations(complex128{cx, cy})
			
			// Convert to character and color
			char, color := m.getPixelChar(iterations)
			style := lipgloss.NewStyle().Foreground(color)
			line.WriteString(style.Render(char))
		}
		lines[y] = line.String()
	}
	
	return lines
}

func (m model) mandelbrotIterations(c complex128) int {
	z := complex128{0, 0}
	
	for i := 0; i < m.maxIter; i++ {
		if z.abs() > 2.0 {
			return i
		}
		z = z.mul(z).add(c)
	}
	
	return m.maxIter
}

func (m model) getPixelChar(iterations int) (string, lipgloss.Color) {
	if iterations == m.maxIter {
		// Point is in the Mandelbrot set - use black
		return "█", lipgloss.Color("#000000")
	}
	
	// Use a logarithmic scale for better detail at boundaries
	logRatio := math.Log(float64(iterations+1)) / math.Log(float64(m.maxIter+1))
	
	if logRatio < 0.15 {
		chars := []string{"█", "▓", "▒"}
		return chars[iterations%len(chars)], lipgloss.Color("#FF0000") // Bright red
	} else if logRatio < 0.3 {
		chars := []string{"▒", "░", "▫"}
		return chars[iterations%len(chars)], lipgloss.Color("#FF4400") // Red-orange
	} else if logRatio < 0.45 {
		chars := []string{"▫", "•", "◦"}
		return chars[iterations%len(chars)], lipgloss.Color("#FF8800") // Orange
	} else if logRatio < 0.6 {
		chars := []string{"◦", "∘", "·"}
		return chars[iterations%len(chars)], lipgloss.Color("#FFCC00") // Yellow
	} else if logRatio < 0.75 {
		chars := []string{"·", ".", " "}
		return chars[iterations%len(chars)], lipgloss.Color("#88FF00") // Yellow-green
	} else if logRatio < 0.85 {
		return " ", lipgloss.Color("#00FF88") // Green
	} else if logRatio < 0.95 {
		return " ", lipgloss.Color("#0088FF") // Blue
	} else {
		return " ", lipgloss.Color("#8800FF") // Purple
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}