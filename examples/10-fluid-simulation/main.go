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

type droplet struct {
	x, y     float64
	vx, vy   float64
	life     float64
	size     float64
	ripples  []ripple
}

type ripple struct {
	x, y     float64
	radius   float64
	strength float64
	age      float64
}

type model struct {
	width     int
	height    int
	droplets  []droplet
	surface   [][]float64
	time      float64
	gravity   float64
	viscosity float64
	paused    bool
	mode      string
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
		droplets:  []droplet{},
		gravity:   0.3,
		viscosity: 0.98,
		mode:      "rain",
	}
}

func (m *model) initSurface() {
	m.surface = make([][]float64, m.height)
	for i := range m.surface {
		m.surface[i] = make([]float64, m.width)
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
		m.initSurface()
		return m, nil

	case tickMsg:
		if !m.paused {
			m.time += 0.1
			m.updateSimulation()
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "r":
			m.droplets = []droplet{}
			m.initSurface()
			m.time = 0
		case "1":
			m.mode = "rain"
		case "2":
			m.mode = "drops"
		case "3":
			m.mode = "fountain"
		case "up":
			m.gravity = math.Min(m.gravity+0.1, 1.0)
		case "down":
			m.gravity = math.Max(m.gravity-0.1, 0.1)
		case "left":
			m.viscosity = math.Max(m.viscosity-0.01, 0.90)
		case "right":
			m.viscosity = math.Min(m.viscosity+0.01, 0.99)
		case "c":
			// Add manual droplet at center
			m.addDroplet(float64(m.width)/2, 5, 0, 0, 1.0)
		}
	}

	return m, nil
}

func (m *model) addDroplet(x, y, vx, vy, size float64) {
	if len(m.droplets) < 150 {
		d := droplet{
			x: x, y: y, vx: vx, vy: vy,
			life: 1.0, size: size,
			ripples: []ripple{},
		}
		m.droplets = append(m.droplets, d)
	}
}

func (m *model) updateSimulation() {
	if len(m.surface) == 0 {
		return
	}

	// Generate new droplets based on mode
	switch m.mode {
	case "rain":
		if rand.Float64() < 0.3 {
			x := rand.Float64() * float64(m.width)
			size := 0.5 + rand.Float64()*0.5
			m.addDroplet(x, 0, (rand.Float64()-0.5)*0.5, 0, size)
		}
	case "drops":
		if rand.Float64() < 0.1 {
			x := rand.Float64() * float64(m.width)
			y := rand.Float64() * float64(m.height/2)
			vx := (rand.Float64() - 0.5) * 2
			vy := rand.Float64() * 2
			size := 0.3 + rand.Float64()*0.4
			m.addDroplet(x, y, vx, vy, size)
		}
	case "fountain":
		if rand.Float64() < 0.4 {
			centerX := float64(m.width) / 2
			x := centerX + (rand.Float64()-0.5)*10
			y := float64(m.height) - 5
			vx := (rand.Float64() - 0.5) * 3
			vy := -3 - rand.Float64()*2
			size := 0.4 + rand.Float64()*0.3
			m.addDroplet(x, y, vx, vy, size)
		}
	}

	// Update droplets
	alive := []droplet{}
	for i := range m.droplets {
		d := &m.droplets[i]

		// Apply physics
		d.vy += m.gravity
		d.x += d.vx
		d.y += d.vy
		d.life -= 0.01

		// Update ripples
		newRipples := []ripple{}
		for j := range d.ripples {
			r := &d.ripples[j]
			r.radius += 0.5
			r.strength *= 0.95
			r.age += 0.1
			if r.strength > 0.01 && r.radius < 20 {
				newRipples = append(newRipples, *r)
			}
		}
		d.ripples = newRipples

		// Check for surface collision
		if d.y >= float64(m.height)-10 && d.vy > 0 {
			// Create ripple on impact
			if len(d.ripples) < 5 {
				impact := math.Min(math.Abs(d.vy)*d.size, 2.0)
				d.ripples = append(d.ripples, ripple{
					x: d.x, y: d.y,
					radius: 0, strength: impact, age: 0,
				})
			}
			// Bounce with energy loss
			d.vy = -d.vy * 0.3
			d.vx *= 0.7
			d.life -= 0.2
		}

		// Check bounds
		if d.x < 0 || d.x >= float64(m.width) {
			d.vx = -d.vx * 0.8
			d.x = math.Max(0, math.Min(float64(m.width-1), d.x))
		}

		// Keep alive droplets
		if d.life > 0 && d.y < float64(m.height) {
			alive = append(alive, *d)
		}
	}
	m.droplets = alive

	// Update surface waves
	m.updateSurface()
}

func (m *model) updateSurface() {
	// Clear surface
	for y := range m.surface {
		for x := range m.surface[y] {
			m.surface[y][x] = 0
		}
	}

	// Add ripple effects
	for _, d := range m.droplets {
		for _, r := range d.ripples {
			m.addRippleToSurface(r)
		}
	}

	// Add base wave motion
	waterLevel := float64(m.height) - 8
	for x := 0; x < m.width; x++ {
		wave := math.Sin(float64(x)*0.2+m.time*2) * 0.5
		wave += math.Sin(float64(x)*0.1+m.time*1.5) * 0.3
		y := int(waterLevel + wave)
		if y >= 0 && y < m.height {
			m.surface[y][x] = math.Max(m.surface[y][x], 0.3)
		}
	}
}

func (m *model) addRippleToSurface(r ripple) {
	centerX, centerY := int(r.x), int(r.y)
	radius := int(r.radius)

	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			x, y := centerX+dx, centerY+dy
			if x >= 0 && x < m.width && y >= 0 && y < m.height {
				dist := math.Sqrt(float64(dx*dx + dy*dy))
				if dist <= r.radius {
					// Calculate wave height based on distance
					waveHeight := r.strength * math.Cos(dist*math.Pi/(r.radius*2))
					if waveHeight > 0 {
						m.surface[y][x] = math.Max(m.surface[y][x], waveHeight)
					}
				}
			}
		}
	}
}

func (m model) View() string {
	if len(m.surface) == 0 {
		return "Initializing fluid simulation..."
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Blue).
		Padding(0, 1)

	title := titleStyle.Render("💧 Fluid Simulation")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Cyan)
	status := statusStyle.Render(fmt.Sprintf(
		"Mode: %s | Droplets: %d | Gravity: %.1f | Viscosity: %.2f | %s",
		strings.Title(m.mode), len(m.droplets), m.gravity, m.viscosity,
		map[bool]string{true: "⏸ Paused", false: "💧 Flowing"}[m.paused],
	))

	// Render simulation
	lines := make([]string, m.height)
	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			char, color := m.getFluidChar(x, y)
			style := lipgloss.NewStyle().Foreground(color)
			line.WriteString(style.Render(char))
		}
		lines[y] = line.String()
	}

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1]rain [2]drops [3]fountain • [↑↓] gravity • [←→] viscosity • [c] add drop • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) getFluidChar(x, y int) (string, lipgloss.Color) {
	// Check for droplets first
	for _, d := range m.droplets {
		if int(d.x) == x && int(d.y) == y {
			if d.size > 0.7 {
				return "●", common.Blue
			} else {
				return "•", common.Cyan
			}
		}
	}

	// Check surface waves
	if m.surface[y][x] > 0 {
		intensity := m.surface[y][x]
		if intensity > 0.6 {
			chars := []string{"█", "▓", "▒"}
			return chars[rand.Intn(len(chars))], lipgloss.Color("#0066CC")
		} else if intensity > 0.3 {
			chars := []string{"▒", "░", "▫"}
			return chars[rand.Intn(len(chars))], lipgloss.Color("#0088FF")
		} else {
			chars := []string{"░", "▫", "·"}
			return chars[rand.Intn(len(chars))], lipgloss.Color("#00AAFF")
		}
	}

	// Water level background
	waterLevel := float64(m.height) - 8
	if float64(y) >= waterLevel {
		return "░", lipgloss.Color("#003366")
	}

	return " ", lipgloss.Color("#000000")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}