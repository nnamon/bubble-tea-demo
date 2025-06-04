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

type model struct {
	width    int
	height   int
	time     float64
	rotation float64
	zoom     float64
	offsetX  float64
	offsetY  float64
	pattern  int
	paused   bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width:   80,
		height:  24,
		zoom:    1.0,
		pattern: 0,
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
		if !m.paused {
			m.time += 0.1
			m.rotation += 0.02
			m.zoom = 1.0 + math.Sin(m.time*0.3)*0.8
			m.offsetX = math.Sin(m.time*0.15) * 20
			m.offsetY = math.Cos(m.time*0.2) * 15
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "r":
			m.time = 0
			m.rotation = 0
			m.zoom = 1.0
			m.offsetX = 0
			m.offsetY = 0
		case "1":
			m.pattern = 0 // Checkerboard
		case "2":
			m.pattern = 1 // Stripes
		case "3":
			m.pattern = 2 // Dots
		case "4":
			m.pattern = 3 // Mandala
		case "5":
			m.pattern = 4 // Circuit
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF8000")).
		Padding(0, 1)

	title := titleStyle.Render("🌀 Rotozoom Effect")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Orange)
	patterns := []string{"Checkerboard", "Stripes", "Dots", "Mandala", "Circuit"}
	status := statusStyle.Render(fmt.Sprintf(
		"Pattern: %s | Rotation: %.1f° | Zoom: %.2fx | %s",
		patterns[m.pattern], m.rotation*180/math.Pi, m.zoom,
		map[bool]string{true: "⏸ Paused", false: "🌀 Rotating"}[m.paused],
	))

	// Render rotozoom
	lines := m.renderRotozoom()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1-5] patterns • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) renderRotozoom() []string {
	lines := make([]string, m.height)
	centerX := float64(m.width) / 2
	centerY := float64(m.height) / 2

	// Precompute rotation matrix
	cosTheta := math.Cos(m.rotation)
	sinTheta := math.Sin(m.rotation)

	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			// Transform screen coordinates to texture coordinates
			screenX := float64(x) - centerX
			screenY := (float64(y) - centerY) * 2 // Adjust for character aspect ratio

			// Apply inverse rotation and zoom
			texX := (screenX*cosTheta + screenY*sinTheta) / m.zoom
			texY := (-screenX*sinTheta + screenY*cosTheta) / m.zoom

			// Add scrolling offset
			texX += m.offsetX
			texY += m.offsetY

			// Sample the pattern
			char, color := m.samplePattern(texX, texY)
			style := lipgloss.NewStyle().Foreground(color)
			line.WriteString(style.Render(char))
		}
		lines[y] = line.String()
	}

	return lines
}

func (m model) samplePattern(x, y float64) (string, lipgloss.Color) {
	switch m.pattern {
	case 0:
		return m.checkerboardPattern(x, y)
	case 1:
		return m.stripesPattern(x, y)
	case 2:
		return m.dotsPattern(x, y)
	case 3:
		return m.mandalaPattern(x, y)
	case 4:
		return m.circuitPattern(x, y)
	default:
		return m.checkerboardPattern(x, y)
	}
}

func (m model) checkerboardPattern(x, y float64) (string, lipgloss.Color) {
	tileSize := 4.0
	tileX := int(math.Floor(x / tileSize))
	tileY := int(math.Floor(y / tileSize))

	if (tileX+tileY)%2 == 0 {
		return "█", lipgloss.Color("#FFFFFF")
	} else {
		return "█", lipgloss.Color("#000000")
	}
}

func (m model) stripesPattern(x, y float64) (string, lipgloss.Color) {
	stripeWidth := 3.0
	stripeIndex := int(math.Floor(x / stripeWidth))

	colors := []lipgloss.Color{
		lipgloss.Color("#FF0000"),
		lipgloss.Color("#00FF00"),
		lipgloss.Color("#0000FF"),
		lipgloss.Color("#FFFF00"),
		lipgloss.Color("#FF00FF"),
		lipgloss.Color("#00FFFF"),
	}

	colorIndex := stripeIndex % len(colors)
	if colorIndex < 0 {
		colorIndex += len(colors)
	}

	return "█", colors[colorIndex]
}

func (m model) dotsPattern(x, y float64) (string, lipgloss.Color) {
	gridSize := 6.0
	dotRadius := 2.0

	// Find grid position
	gridX := math.Mod(x, gridSize)
	gridY := math.Mod(y, gridSize)

	// Distance from grid center
	centerX := gridSize / 2
	centerY := gridSize / 2
	distance := math.Sqrt((gridX-centerX)*(gridX-centerX) + (gridY-centerY)*(gridY-centerY))

	if distance < dotRadius {
		// Color based on position
		colorValue := math.Sin(x*0.1) * math.Cos(y*0.1)
		if colorValue > 0.3 {
			return "●", lipgloss.Color("#FF4080")
		} else if colorValue > -0.3 {
			return "●", lipgloss.Color("#4080FF")
		} else {
			return "●", lipgloss.Color("#80FF40")
		}
	} else {
		return " ", lipgloss.Color("#000000")
	}
}

func (m model) mandalaPattern(x, y float64) (string, lipgloss.Color) {
	// Distance from origin
	distance := math.Sqrt(x*x + y*y)
	// Angle from origin
	angle := math.Atan2(y, x)

	// Create mandala pattern
	rings := math.Sin(distance * 0.3)
	spokes := math.Sin(angle * 8)
	pattern := rings * spokes

	// Add time-based rotation
	timePattern := math.Sin(distance*0.2 - m.time*2) * math.Cos(angle*6 + m.time)

	combinedPattern := (pattern + timePattern) / 2

	var char string
	var color lipgloss.Color

	if combinedPattern > 0.6 {
		char = "◆"
		color = lipgloss.Color("#FFD700")
	} else if combinedPattern > 0.2 {
		char = "◇"
		color = lipgloss.Color("#FF8000")
	} else if combinedPattern > -0.2 {
		char = "○"
		color = lipgloss.Color("#FF4000")
	} else if combinedPattern > -0.6 {
		char = "∘"
		color = lipgloss.Color("#800040")
	} else {
		char = " "
		color = lipgloss.Color("#000000")
	}

	return char, color
}

func (m model) circuitPattern(x, y float64) (string, lipgloss.Color) {
	gridSize := 8.0
	lineWidth := 1.0

	// Grid coordinates
	gridX := math.Mod(x, gridSize)
	gridY := math.Mod(y, gridSize)

	// Circuit board traces
	isHorizontalTrace := math.Abs(gridY-gridSize/2) < lineWidth
	isVerticalTrace := math.Abs(gridX-gridSize/2) < lineWidth

	// Circuit pads at intersections
	isNearCenter := math.Abs(gridX-gridSize/2) < lineWidth*2 && math.Abs(gridY-gridSize/2) < lineWidth*2

	// Add some randomness based on position
	hash := math.Sin(math.Floor(x/gridSize)*12.345 + math.Floor(y/gridSize)*67.890)

	if isNearCenter && hash > 0.3 {
		return "●", lipgloss.Color("#00FF80")
	} else if isHorizontalTrace || isVerticalTrace {
		if hash > 0 {
			return "─", lipgloss.Color("#80FF80")
		} else {
			return "│", lipgloss.Color("#80FF80")
		}
	} else {
		// Background with occasional components
		if hash > 0.8 {
			return "▪", lipgloss.Color("#404040")
		} else {
			return " ", lipgloss.Color("#000000")
		}
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}