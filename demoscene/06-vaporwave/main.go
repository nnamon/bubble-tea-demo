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
	width     int
	height    int
	time      float64
	speed     float64
	colorMode int
	paused    bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width:     80,
		height:    24,
		speed:     1.0,
		colorMode: 0,
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
			m.time += 0.1 * m.speed
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
		case "1":
			m.colorMode = 0 // Classic vaporwave
		case "2":
			m.colorMode = 1 // Cyberpunk
		case "3":
			m.colorMode = 2 // Synthwave
		case "4":
			m.colorMode = 3 // Retrowave
		case "up":
			m.speed = math.Min(m.speed+0.2, 3.0)
		case "down":
			m.speed = math.Max(m.speed-0.2, 0.1)
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF00FF")).
		Padding(0, 1)

	title := titleStyle.Render("🌆 Vaporwave Landscape")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Pink)
	colorModes := []string{"Classic", "Cyberpunk", "Synthwave", "Retrowave"}
	status := statusStyle.Render(fmt.Sprintf(
		"Mode: %s | Speed: %.1f | %s",
		colorModes[m.colorMode], m.speed,
		map[bool]string{true: "⏸ Paused", false: "🌆 Flowing"}[m.paused],
	))

	// Check minimum size requirements
	minWidth, minHeight := 60, 16
	if m.width < minWidth || m.height < minHeight {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
		
		sizeError := errorStyle.Render(fmt.Sprintf(
			"Terminal too small!\nMinimum size: %dx%d\nCurrent size: %dx%d\n\nPlease resize your terminal window.",
			minWidth, minHeight+4, m.width, m.height+4,
		))
		
		// Help
		helpStyle := lipgloss.NewStyle().Faint(true)
		help := helpStyle.Render(
			"[q]uit",
		)

		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			status,
			"",
			sizeError,
			help,
		)
	}

	// Render vaporwave scene
	scene := m.renderScene()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1-4] color modes • [↑↓] speed • [space] pause • [r]eset • [q]uit",
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		status,
		"",
		scene,
		help,
	)
}

func (m model) renderScene() string {
	var lines []string

	// Ensure we have enough height
	if m.height < 6 {
		return "Scene too small to render"
	}

	// Render sky (top third)
	skyHeight := max(1, m.height/3)
	for y := 0; y < skyHeight; y++ {
		lines = append(lines, m.renderSkyLine(y, skyHeight))
	}

	// Render sun area - overlay on existing sky line
	sunY := max(1, m.height/4)
	if sunY >= 0 && sunY < len(lines) {
		lines[sunY] = m.renderSunLine(sunY)
	}

	// Render grid (remaining height)
	gridStart := skyHeight
	for y := gridStart; y < m.height; y++ {
		lines = append(lines, m.renderGridLine(y, gridStart))
	}

	// Ensure we don't exceed expected height
	if len(lines) > m.height {
		lines = lines[:m.height]
	}

	return strings.Join(lines, "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) renderSkyLine(y, skyHeight int) string {
	intensity := float64(y) / float64(skyHeight)
	char, color := m.getSkyColor(intensity)
	
	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(strings.Repeat(char, m.width))
}

func (m model) renderSunLine(y int) string {
	line := strings.Builder{}
	sunCenterX := m.width / 2
	sunCenterY := m.height / 4
	sunRadius := 4.0

	for x := 0; x < m.width; x++ {
		dx := float64(x - sunCenterX)
		dy := float64(y - sunCenterY) * 2 // Adjust for character aspect ratio
		distance := math.Sqrt(dx*dx + dy*dy)

		var char string
		var color lipgloss.Color

		if distance < sunRadius-1 {
			// Sun core
			char = "●"
			color = m.getSunColor(0.9)
		} else if distance < sunRadius {
			// Sun edge
			char = "○"
			color = m.getSunColor(0.7)
		} else if distance < sunRadius+1 && int(distance+m.time*4)%3 == 0 {
			// Sun rays
			char = "▬"
			color = m.getSunColor(0.5)
		} else {
			// Sky background
			intensity := float64(y) / float64(m.height/3)
			char, color = m.getSkyColor(intensity)
		}

		style := lipgloss.NewStyle().Foreground(color)
		line.WriteString(style.Render(char))
	}

	return line.String()
}

func (m model) renderGridLine(y, gridStart int) string {
	line := strings.Builder{}
	depth := float64(y - gridStart + 1)
	if depth <= 0 {
		depth = 1
	}

	// Perspective scaling
	scale := 20.0 / depth
	offset := m.time * m.speed * scale

	for x := 0; x < m.width; x++ {
		// Grid coordinates with perspective
		gridX := (float64(x) - float64(m.width)/2) / scale
		gridZ := depth + offset

		// Add wave effect
		waveOffset := math.Sin(gridX*0.3+gridZ*0.2+m.time) * 2.0 / depth

		// Determine if we're on a grid line
		isGridLineX := math.Abs(math.Mod(gridX+0.5, 2.0)-1.0) < 0.2/scale
		isGridLineZ := math.Abs(math.Mod(gridZ+0.5, 2.0)-1.0) < 0.2/scale

		var char string
		var color lipgloss.Color

		if isGridLineX || isGridLineZ {
			intensity := 1.0 / (depth*0.1 + 1)
			if isGridLineX && isGridLineZ {
				char = "+"
				intensity *= 1.5
			} else if isGridLineX {
				char = "|"
			} else {
				char = "-"
			}

			// Add glow effect
			glowIntensity := intensity + math.Sin(waveOffset+m.time*2)*0.3
			glowIntensity = math.Max(0, math.Min(1, glowIntensity))

			char, color = m.getGridColor(glowIntensity)
		} else {
			char = " "
			color = lipgloss.Color("#000000")
		}

		style := lipgloss.NewStyle().Foreground(color)
		line.WriteString(style.Render(char))
	}

	return line.String()
}

func (m model) getSkyColor(intensity float64) (string, lipgloss.Color) {
	switch m.colorMode {
	case 0: // Classic vaporwave
		if intensity < 0.3 {
			return "▓", lipgloss.Color("#FF00FF")
		} else if intensity < 0.6 {
			return "▒", lipgloss.Color("#FF0080")
		} else {
			return "░", lipgloss.Color("#8000FF")
		}
	case 1: // Cyberpunk
		if intensity < 0.3 {
			return "▓", lipgloss.Color("#00FFFF")
		} else if intensity < 0.6 {
			return "▒", lipgloss.Color("#0080FF")
		} else {
			return "░", lipgloss.Color("#000080")
		}
	case 2: // Synthwave
		if intensity < 0.3 {
			return "▓", lipgloss.Color("#FF4080")
		} else if intensity < 0.6 {
			return "▒", lipgloss.Color("#FF8040")
		} else {
			return "░", lipgloss.Color("#FFFF00")
		}
	case 3: // Retrowave
		if intensity < 0.3 {
			return "▓", lipgloss.Color("#FF0040")
		} else if intensity < 0.6 {
			return "▒", lipgloss.Color("#8000FF")
		} else {
			return "░", lipgloss.Color("#4000FF")
		}
	default:
		return "░", lipgloss.Color("#FF00FF")
	}
}

func (m model) getSunColor(intensity float64) lipgloss.Color {
	switch m.colorMode {
	case 0: // Classic vaporwave
		if intensity > 0.8 {
			return lipgloss.Color("#FFFF00")
		} else if intensity > 0.5 {
			return lipgloss.Color("#FFCC00")
		} else {
			return lipgloss.Color("#FF8800")
		}
	case 1: // Cyberpunk
		if intensity > 0.8 {
			return lipgloss.Color("#00FFFF")
		} else if intensity > 0.5 {
			return lipgloss.Color("#0088FF")
		} else {
			return lipgloss.Color("#0044FF")
		}
	case 2: // Synthwave
		if intensity > 0.8 {
			return lipgloss.Color("#FF8080")
		} else if intensity > 0.5 {
			return lipgloss.Color("#FF4040")
		} else {
			return lipgloss.Color("#FF0000")
		}
	case 3: // Retrowave
		if intensity > 0.8 {
			return lipgloss.Color("#FF80FF")
		} else if intensity > 0.5 {
			return lipgloss.Color("#FF40FF")
		} else {
			return lipgloss.Color("#FF00FF")
		}
	default:
		return lipgloss.Color("#FFFF00")
	}
}

func (m model) getGridColor(intensity float64) (string, lipgloss.Color) {
	var char string
	if intensity > 0.7 {
		char = "█"
	} else if intensity > 0.4 {
		char = "▓"
	} else if intensity > 0.2 {
		char = "▒"
	} else {
		char = "░"
	}

	switch m.colorMode {
	case 0: // Classic vaporwave
		if intensity > 0.6 {
			return char, lipgloss.Color("#FF00FF")
		} else if intensity > 0.3 {
			return char, lipgloss.Color("#8000FF")
		} else {
			return char, lipgloss.Color("#400080")
		}
	case 1: // Cyberpunk
		if intensity > 0.6 {
			return char, lipgloss.Color("#00FFFF")
		} else if intensity > 0.3 {
			return char, lipgloss.Color("#008080")
		} else {
			return char, lipgloss.Color("#004040")
		}
	case 2: // Synthwave
		if intensity > 0.6 {
			return char, lipgloss.Color("#FF4080")
		} else if intensity > 0.3 {
			return char, lipgloss.Color("#FF0040")
		} else {
			return char, lipgloss.Color("#800020")
		}
	case 3: // Retrowave
		if intensity > 0.6 {
			return char, lipgloss.Color("#FF8000")
		} else if intensity > 0.3 {
			return char, lipgloss.Color("#FF4000")
		} else {
			return char, lipgloss.Color("#802000")
		}
	default:
		return char, lipgloss.Color("#FF00FF")
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}