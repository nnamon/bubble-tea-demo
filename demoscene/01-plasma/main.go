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
	paused    bool
	palette   int
	intensity float64
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
		palette:   0,
		intensity: 1.0,
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
		case "1", "2", "3", "4":
			switch msg.String() {
			case "1":
				m.palette = 0 // Classic fire
			case "2":
				m.palette = 1 // Ocean
			case "3":
				m.palette = 2 // Psychedelic
			case "4":
				m.palette = 3 // Monochrome
			}
		case "up":
			m.speed = math.Min(m.speed+0.2, 3.0)
		case "down":
			m.speed = math.Max(m.speed-0.2, 0.1)
		case "left":
			m.intensity = math.Max(m.intensity-0.1, 0.3)
		case "right":
			m.intensity = math.Min(m.intensity+0.1, 2.0)
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF0080")).
		Padding(0, 1)

	title := titleStyle.Render("🌈 Plasma Effect")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Cyan)
	palettes := []string{"Fire", "Ocean", "Psychedelic", "Monochrome"}
	status := statusStyle.Render(fmt.Sprintf(
		"Palette: %s | Speed: %.1f | Intensity: %.1f | %s",
		palettes[m.palette], m.speed, m.intensity,
		map[bool]string{true: "⏸ Paused", false: "🌈 Flowing"}[m.paused],
	))

	// Render plasma
	lines := m.renderPlasma()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1-4] palettes • [↑↓] speed • [←→] intensity • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) renderPlasma() []string {
	lines := make([]string, m.height)

	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			// Calculate plasma value using multiple sine waves
			fx := float64(x) / float64(m.width) * 16
			fy := float64(y) / float64(m.height) * 16

			// Classic plasma formula with multiple frequency components
			value := math.Sin(fx*0.5+m.time) +
				math.Sin(fy*0.3+m.time*1.2) +
				math.Sin((fx+fy)*0.25+m.time*0.8) +
				math.Sin(math.Sqrt(fx*fx+fy*fy)*0.4+m.time*1.5) +
				math.Sin(fx*0.1+fy*0.2+m.time*0.6)

			// Normalize and apply intensity
			value = (value + 5) / 10 * m.intensity
			value = math.Max(0, math.Min(1, value))

			// Convert to character and color
			char, color := m.getPlasmaChar(value)
			style := lipgloss.NewStyle().Foreground(color)
			line.WriteString(style.Render(char))
		}
		lines[y] = line.String()
	}

	return lines
}

func (m model) getPlasmaChar(value float64) (string, lipgloss.Color) {
	// Choose character based on intensity
	chars := []string{" ", "·", "∘", "•", "◦", "○", "●", "▫", "▪", "▒", "▓", "█"}
	charIndex := int(value * float64(len(chars)-1))
	if charIndex >= len(chars) {
		charIndex = len(chars) - 1
	}
	char := chars[charIndex]

	// Choose color based on palette
	var color lipgloss.Color
	switch m.palette {
	case 0: // Fire palette
		color = m.getFireColor(value)
	case 1: // Ocean palette
		color = m.getOceanColor(value)
	case 2: // Psychedelic palette
		color = m.getPsychedelicColor(value)
	case 3: // Monochrome palette
		color = m.getMonochromeColor(value)
	default:
		color = m.getFireColor(value)
	}

	return char, color
}

func (m model) getFireColor(value float64) lipgloss.Color {
	if value < 0.2 {
		return lipgloss.Color("#330000")
	} else if value < 0.4 {
		return lipgloss.Color("#660000")
	} else if value < 0.6 {
		return lipgloss.Color("#990000")
	} else if value < 0.7 {
		return lipgloss.Color("#CC3300")
	} else if value < 0.8 {
		return lipgloss.Color("#FF4400")
	} else if value < 0.9 {
		return lipgloss.Color("#FF8800")
	} else {
		return lipgloss.Color("#FFCC00")
	}
}

func (m model) getOceanColor(value float64) lipgloss.Color {
	if value < 0.2 {
		return lipgloss.Color("#000033")
	} else if value < 0.4 {
		return lipgloss.Color("#000066")
	} else if value < 0.6 {
		return lipgloss.Color("#003399")
	} else if value < 0.7 {
		return lipgloss.Color("#0066CC")
	} else if value < 0.8 {
		return lipgloss.Color("#0099FF")
	} else if value < 0.9 {
		return lipgloss.Color("#33CCFF")
	} else {
		return lipgloss.Color("#66FFFF")
	}
}

func (m model) getPsychedelicColor(value float64) lipgloss.Color {
	// Cycle through rainbow colors
	hue := value * 360
	if hue < 60 {
		return lipgloss.Color("#FF0080")
	} else if hue < 120 {
		return lipgloss.Color("#8000FF")
	} else if hue < 180 {
		return lipgloss.Color("#0080FF")
	} else if hue < 240 {
		return lipgloss.Color("#00FF80")
	} else if hue < 300 {
		return lipgloss.Color("#80FF00")
	} else {
		return lipgloss.Color("#FF8000")
	}
}

func (m model) getMonochromeColor(value float64) lipgloss.Color {
	gray := int(value * 255)
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", gray, gray, gray))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}