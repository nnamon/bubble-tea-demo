package main

import (
	"fmt"
	"math"
	"math/rand"
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
	fireField [][]float64
	intensity float64
	windForce float64
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
		intensity: 1.0,
		windForce: 0.0,
		paused:    false,
	}
}

func (m *model) initFireField() {
	m.fireField = make([][]float64, m.height)
	for i := range m.fireField {
		m.fireField[i] = make([]float64, m.width)
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
		m.initFireField()
		return m, nil

	case tickMsg:
		if !m.paused {
			m.updateFire()
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "r":
			m.initFireField()
		case "up":
			m.intensity = math.Min(m.intensity+0.1, 2.0)
		case "down":
			m.intensity = math.Max(m.intensity-0.1, 0.1)
		case "left":
			m.windForce = math.Max(m.windForce-0.1, -1.0)
		case "right":
			m.windForce = math.Min(m.windForce+0.1, 1.0)
		case "0":
			m.windForce = 0.0
		}
	}

	return m, nil
}

func (m *model) updateFire() {
	if len(m.fireField) == 0 {
		return
	}

	// Create new fire field
	newField := make([][]float64, m.height)
	for i := range newField {
		newField[i] = make([]float64, m.width)
	}

	// Add heat sources at the bottom
	bottomRow := m.height - 1
	for x := 0; x < m.width; x++ {
		// Create hot spots with some randomness
		if rand.Float64() < 0.7 {
			heat := (0.8 + rand.Float64()*0.2) * m.intensity
			// Create some variation in the base fire
			if x%3 == 0 || x%7 == 0 {
				heat *= 1.2
			}
			m.fireField[bottomRow][x] = heat
		}
	}

	// Propagate fire upward with cooling and wind
	for y := 1; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			// Gather heat from below and surrounding cells
			heat := 0.0
			samples := 0

			// Sample from below (main heat source)
			if y < m.height-1 {
				heat += m.fireField[y+1][x] * 0.4
				samples++

				// Sample diagonally below for spread
				if x > 0 {
					heat += m.fireField[y+1][x-1] * 0.2
					samples++
				}
				if x < m.width-1 {
					heat += m.fireField[y+1][x+1] * 0.2
					samples++
				}
			}

			// Apply wind effect
			windOffset := int(m.windForce * 2)
			windX := x - windOffset
			if windX >= 0 && windX < m.width && y < m.height-1 {
				heat += m.fireField[y+1][windX] * 0.2
				samples++
			}

			// Add some randomness and turbulence
			heat += (rand.Float64() - 0.5) * 0.1

			// Cool down as it rises
			coolingFactor := 0.95 - (float64(m.height-y)/float64(m.height))*0.3
			heat *= coolingFactor

			// Add horizontal spreading
			if x > 0 {
				heat += m.fireField[y][x-1] * 0.1
			}
			if x < m.width-1 {
				heat += m.fireField[y][x+1] * 0.1
			}

			newField[y][x] = math.Max(0, heat)
		}
	}

	m.fireField = newField
}

func (m model) View() string {
	if len(m.fireField) == 0 {
		return "Initializing fire..."
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF4500")).
		Padding(0, 1)

	title := titleStyle.Render("🔥 Fire Effect")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Yellow)
	status := statusStyle.Render(fmt.Sprintf(
		"Intensity: %.1f | Wind: %.1f | %s",
		m.intensity, m.windForce,
		map[bool]string{true: "⏸ Paused", false: "🔥 Burning"}[m.paused],
	))

	// Render fire
	lines := make([]string, m.height)
	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			heat := m.fireField[y][x]
			char, color := m.getFireChar(heat)
			style := lipgloss.NewStyle().Foreground(color)
			line.WriteString(style.Render(char))
		}
		lines[y] = line.String()
	}

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[↑↓] intensity • [←→] wind • [0] calm wind • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s  %s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) getFireChar(heat float64) (string, lipgloss.Color) {
	if heat < 0.1 {
		return " ", lipgloss.Color("#000000")
	} else if heat < 0.2 {
		chars := []string{".", "·", "∘"}
		return chars[rand.Intn(len(chars))], lipgloss.Color("#330000")
	} else if heat < 0.35 {
		chars := []string{"∘", "•", "◦"}
		return chars[rand.Intn(len(chars))], lipgloss.Color("#660000")
	} else if heat < 0.5 {
		chars := []string{"▁", "▂", "▃"}
		return chars[rand.Intn(len(chars))], lipgloss.Color("#990000")
	} else if heat < 0.65 {
		chars := []string{"▄", "▅", "▆"}
		return chars[rand.Intn(len(chars))], lipgloss.Color("#CC3300")
	} else if heat < 0.8 {
		chars := []string{"▇", "█", "▉"}
		return chars[rand.Intn(len(chars))], lipgloss.Color("#FF4500")
	} else if heat < 0.95 {
		chars := []string{"▓", "▒", "░"}
		return chars[rand.Intn(len(chars))], lipgloss.Color("#FF6600")
	} else {
		chars := []string{"▓", "▒", "░", "▔"}
		return chars[rand.Intn(len(chars))], lipgloss.Color("#FFAA00")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}