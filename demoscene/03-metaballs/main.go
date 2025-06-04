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

type metaball struct {
	x, y       float64
	vx, vy     float64
	radius     float64
	strength   float64
	colorPhase float64
}

type model struct {
	width     int
	height    int
	metaballs []metaball
	time      float64
	threshold float64
	paused    bool
	colorMode int
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	// Create initial metaballs
	balls := []metaball{
		{x: 20, y: 10, vx: 0.8, vy: 0.3, radius: 8, strength: 1.0, colorPhase: 0},
		{x: 40, y: 15, vx: -0.5, vy: 0.7, radius: 6, strength: 0.8, colorPhase: math.Pi / 3},
		{x: 60, y: 8, vx: 0.6, vy: -0.4, radius: 7, strength: 0.9, colorPhase: 2 * math.Pi / 3},
		{x: 30, y: 20, vx: -0.7, vy: -0.6, radius: 5, strength: 0.7, colorPhase: math.Pi},
	}

	return model{
		width:     80,
		height:    24,
		metaballs: balls,
		threshold: 1.0,
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
			m.time += 0.1
			m.updateMetaballs()
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
			m = initialModel()
			m.width = msg.Width
			m.height = msg.Height - 4
		case "1":
			m.colorMode = 0 // Classic
		case "2":
			m.colorMode = 1 // Rainbow
		case "3":
			m.colorMode = 2 // Heat
		case "4":
			m.colorMode = 3 // Electric
		case "up":
			m.threshold = math.Min(m.threshold+0.1, 3.0)
		case "down":
			m.threshold = math.Max(m.threshold-0.1, 0.3)
		case "a":
			// Add new metaball
			if len(m.metaballs) < 8 {
				newBall := metaball{
					x:          float64(m.width) / 2,
					y:          float64(m.height) / 2,
					vx:         (math.Sin(m.time) * 0.8),
					vy:         (math.Cos(m.time) * 0.8),
					radius:     4 + math.Sin(m.time*2)*2,
					strength:   0.6 + math.Sin(m.time*3)*0.3,
					colorPhase: m.time,
				}
				m.metaballs = append(m.metaballs, newBall)
			}
		case "d":
			// Remove last metaball
			if len(m.metaballs) > 1 {
				m.metaballs = m.metaballs[:len(m.metaballs)-1]
			}
		}
	}

	return m, nil
}

func (m *model) updateMetaballs() {
	for i := range m.metaballs {
		ball := &m.metaballs[i]

		// Update position
		ball.x += ball.vx
		ball.y += ball.vy

		// Bounce off walls
		if ball.x <= ball.radius || ball.x >= float64(m.width)-ball.radius {
			ball.vx = -ball.vx
			ball.x = math.Max(ball.radius, math.Min(float64(m.width)-ball.radius, ball.x))
		}
		if ball.y <= ball.radius || ball.y >= float64(m.height)-ball.radius {
			ball.vy = -ball.vy
			ball.y = math.Max(ball.radius, math.Min(float64(m.height)-ball.radius, ball.y))
		}

		// Add some organic movement
		ball.vx += math.Sin(m.time*0.7+ball.colorPhase) * 0.05
		ball.vy += math.Cos(m.time*0.8+ball.colorPhase) * 0.05

		// Limit velocity
		maxVel := 1.5
		vel := math.Sqrt(ball.vx*ball.vx + ball.vy*ball.vy)
		if vel > maxVel {
			ball.vx = ball.vx / vel * maxVel
			ball.vy = ball.vy / vel * maxVel
		}

		// Animate radius and strength
		ball.radius = 4 + math.Sin(m.time*1.2+ball.colorPhase)*2
		ball.strength = 0.7 + math.Sin(m.time*0.9+ball.colorPhase)*0.3
	}
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF4080")).
		Padding(0, 1)

	title := titleStyle.Render("🫧 Metaballs")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Pink)
	colorModes := []string{"Classic", "Rainbow", "Heat", "Electric"}
	status := statusStyle.Render(fmt.Sprintf(
		"Balls: %d | Threshold: %.1f | Mode: %s | %s",
		len(m.metaballs), m.threshold, colorModes[m.colorMode],
		map[bool]string{true: "⏸ Paused", false: "🫧 Flowing"}[m.paused],
	))

	// Render metaballs
	lines := m.renderMetaballs()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[a]dd ball • [d]elete ball • [1-4] color modes • [↑↓] threshold • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) renderMetaballs() []string {
	lines := make([]string, m.height)

	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			// Calculate metaball field strength at this position
			totalStrength := 0.0
			colorInfluence := 0.0

			for _, ball := range m.metaballs {
				// Distance from this pixel to the metaball center
				dx := float64(x) - ball.x
				dy := (float64(y) - ball.y) * 2 // Adjust for character aspect ratio
				distance := math.Sqrt(dx*dx + dy*dy)

				if distance > 0 {
					// Metaball field strength (inverse square law)
					strength := ball.strength * (ball.radius * ball.radius) / (distance * distance)
					totalStrength += strength

					// Weight color influence by strength
					colorInfluence += strength * ball.colorPhase
				}
			}

			// Determine if we're inside the metaball surface
			if totalStrength >= m.threshold {
				char, color := m.getMetaballChar(totalStrength, colorInfluence)
				style := lipgloss.NewStyle().Foreground(color)
				if totalStrength > m.threshold*2 {
					style = style.Bold(true)
				}
				line.WriteString(style.Render(char))
			} else {
				// Outside metaballs - show field lines occasionally
				if totalStrength > m.threshold*0.3 {
					fieldChar := "·"
					if totalStrength > m.threshold*0.6 {
						fieldChar = "∘"
					}
					style := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Faint(true)
					line.WriteString(style.Render(fieldChar))
				} else {
					line.WriteString(" ")
				}
			}
		}
		lines[y] = line.String()
	}

	return lines
}

func (m model) getMetaballChar(strength, colorInfluence float64) (string, lipgloss.Color) {
	// Choose character based on field strength
	chars := []string{"▒", "▓", "█", "▉", "▊", "▋", "▌", "▍", "▎", "▏"}
	normalizedStrength := math.Min(1.0, (strength-m.threshold)/(m.threshold*2))
	charIndex := int(normalizedStrength * float64(len(chars)-1))
	if charIndex >= len(chars) {
		charIndex = len(chars) - 1
	}
	char := chars[charIndex]

	// Choose color based on mode
	var color lipgloss.Color
	switch m.colorMode {
	case 0: // Classic - blue to white
		color = m.getClassicColor(normalizedStrength)
	case 1: // Rainbow
		color = m.getRainbowColor(colorInfluence + m.time)
	case 2: // Heat - black to red to yellow to white
		color = m.getHeatColor(normalizedStrength)
	case 3: // Electric - electric blue variations
		color = m.getElectricColor(normalizedStrength, m.time)
	default:
		color = m.getClassicColor(normalizedStrength)
	}

	return char, color
}

func (m model) getClassicColor(strength float64) lipgloss.Color {
	if strength < 0.3 {
		return lipgloss.Color("#0044FF")
	} else if strength < 0.6 {
		return lipgloss.Color("#4488FF")
	} else if strength < 0.8 {
		return lipgloss.Color("#88CCFF")
	} else {
		return lipgloss.Color("#CCFFFF")
	}
}

func (m model) getRainbowColor(phase float64) lipgloss.Color {
	hue := math.Mod(phase*60, 360)
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

func (m model) getHeatColor(strength float64) lipgloss.Color {
	if strength < 0.25 {
		return lipgloss.Color("#440000")
	} else if strength < 0.5 {
		return lipgloss.Color("#880000")
	} else if strength < 0.75 {
		return lipgloss.Color("#FF4400")
	} else {
		return lipgloss.Color("#FFFF00")
	}
}

func (m model) getElectricColor(strength, time float64) lipgloss.Color {
	flicker := math.Sin(time*20) * 0.2
	intensity := strength + flicker
	
	if intensity < 0.3 {
		return lipgloss.Color("#001188")
	} else if intensity < 0.6 {
		return lipgloss.Color("#0044FF")
	} else if intensity < 0.8 {
		return lipgloss.Color("#00AAFF")
	} else {
		return lipgloss.Color("#88FFFF")
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}