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
	width      int
	height     int
	time       float64
	speed      float64
	tunnelMode int
	paused     bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width:      80,
		height:     24,
		speed:      1.0,
		tunnelMode: 0,
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
			m.tunnelMode = 0 // Classic tunnel
		case "2":
			m.tunnelMode = 1 // Checkerboard tunnel
		case "3":
			m.tunnelMode = 2 // Spiral tunnel
		case "4":
			m.tunnelMode = 3 // Ripple tunnel
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
		Background(lipgloss.Color("#8800FF")).
		Padding(0, 1)

	title := titleStyle.Render("🕳️ Tunnel Effect")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Purple)
	modes := []string{"Classic", "Checkerboard", "Spiral", "Ripple"}
	status := statusStyle.Render(fmt.Sprintf(
		"Mode: %s | Speed: %.1f | %s",
		modes[m.tunnelMode], m.speed,
		map[bool]string{true: "⏸ Paused", false: "🕳️ Tunneling"}[m.paused],
	))

	// Render tunnel
	lines := m.renderTunnel()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1-4] tunnel modes • [↑↓] speed • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) renderTunnel() []string {
	lines := make([]string, m.height)
	centerX := float64(m.width) / 2
	centerY := float64(m.height) / 2

	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			// Calculate distance from center
			dx := float64(x) - centerX
			dy := (float64(y) - centerY) * 2 // Adjust for character aspect ratio
			distance := math.Sqrt(dx*dx + dy*dy)
			
			// Calculate angle
			angle := math.Atan2(dy, dx)
			
			// Apply tunnel effect based on mode
			var intensity float64
			var char string
			var color lipgloss.Color
			
			switch m.tunnelMode {
			case 0: // Classic tunnel
				intensity, char, color = m.classicTunnel(distance, angle)
			case 1: // Checkerboard tunnel
				intensity, char, color = m.checkerboardTunnel(distance, angle)
			case 2: // Spiral tunnel
				intensity, char, color = m.spiralTunnel(distance, angle)
			case 3: // Ripple tunnel
				intensity, char, color = m.rippleTunnel(distance, angle)
			}
			
			style := lipgloss.NewStyle().Foreground(color)
			if intensity < 0.1 {
				style = style.Faint(true)
			} else if intensity > 0.8 {
				style = style.Bold(true)
			}
			
			line.WriteString(style.Render(char))
		}
		lines[y] = line.String()
	}

	return lines
}

func (m model) classicTunnel(distance, angle float64) (float64, string, lipgloss.Color) {
	if distance < 1 {
		distance = 1
	}
	
	// Create tunnel depth effect
	depth := 50.0/distance + m.time*2
	ringPos := math.Mod(depth, 2.0)
	
	var intensity float64
	var char string
	
	if ringPos < 1.0 {
		intensity = ringPos
		char = "▓"
	} else {
		intensity = 2.0 - ringPos
		char = "▒"
	}
	
	// Color based on depth
	colorValue := math.Mod(depth*0.2, 1.0)
	color := m.getDepthColor(colorValue)
	
	return intensity, char, color
}

func (m model) checkerboardTunnel(distance, angle float64) (float64, string, lipgloss.Color) {
	if distance < 1 {
		distance = 1
	}
	
	depth := 30.0/distance + m.time*3
	angleSegments := int((angle + math.Pi) / (math.Pi / 8))
	depthSegments := int(depth)
	
	var intensity float64
	var char string
	
	if (angleSegments+depthSegments)%2 == 0 {
		intensity = 0.8
		char = "█"
	} else {
		intensity = 0.2
		char = "░"
	}
	
	colorValue := math.Mod(depth*0.1, 1.0)
	color := m.getDepthColor(colorValue)
	
	return intensity, char, color
}

func (m model) spiralTunnel(distance, angle float64) (float64, string, lipgloss.Color) {
	if distance < 1 {
		distance = 1
	}
	
	depth := 40.0/distance + m.time*2
	spiralAngle := angle + depth*0.5
	spiralValue := math.Sin(spiralAngle * 4)
	
	var intensity float64
	var char string
	
	if spiralValue > 0 {
		intensity = spiralValue
		char = "◤"
	} else {
		intensity = -spiralValue
		char = "◥"
	}
	
	colorValue := math.Mod(depth*0.15, 1.0)
	color := m.getSpiralColor(colorValue)
	
	return intensity, char, color
}

func (m model) rippleTunnel(distance, angle float64) (float64, string, lipgloss.Color) {
	if distance < 1 {
		distance = 1
	}
	
	depth := 35.0/distance + m.time*2.5
	ripple := math.Sin(distance*0.3 - m.time*4)
	wave := math.Sin(depth*2 + ripple*2)
	
	intensity := (wave + 1) / 2
	
	var char string
	if intensity > 0.7 {
		char = "●"
	} else if intensity > 0.4 {
		char = "◦"
	} else {
		char = "·"
	}
	
	colorValue := math.Mod(depth*0.25 + ripple*0.1, 1.0)
	color := m.getRippleColor(colorValue)
	
	return intensity, char, color
}

func (m model) getDepthColor(value float64) lipgloss.Color {
	// Blue to red gradient for depth
	if value < 0.33 {
		return lipgloss.Color("#0000FF")
	} else if value < 0.66 {
		return lipgloss.Color("#8800FF")
	} else {
		return lipgloss.Color("#FF00FF")
	}
}

func (m model) getSpiralColor(value float64) lipgloss.Color {
	// Green to yellow gradient for spiral
	if value < 0.33 {
		return lipgloss.Color("#00FF00")
	} else if value < 0.66 {
		return lipgloss.Color("#88FF00")
	} else {
		return lipgloss.Color("#FFFF00")
	}
}

func (m model) getRippleColor(value float64) lipgloss.Color {
	// Cyan to white gradient for ripples
	if value < 0.25 {
		return lipgloss.Color("#00FFFF")
	} else if value < 0.5 {
		return lipgloss.Color("#44FFFF")
	} else if value < 0.75 {
		return lipgloss.Color("#88FFFF")
	} else {
		return lipgloss.Color("#CCFFFF")
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}