package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type particle struct {
	x, y   float64
	vx, vy float64
	life   float64
	char   string
	color  lipgloss.Color
}

type model struct {
	width     int
	height    int
	particles []particle
	emitting  bool
	gravity   float64
	wind      float64
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
		particles: []particle{},
		emitting:  true,
		gravity:   0.1,
		wind:      0.0,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m *model) emitParticle() {
	chars := []string{"✦", "✧", "⋆", "◦", "•", "∘", "○", "◌"}
	colors := []lipgloss.Color{common.Yellow, common.Orange, common.Red, common.Pink}
	
	p := particle{
		x:     float64(m.width) / 2,
		y:     float64(m.height) - 5,
		vx:    (rand.Float64() - 0.5) * 3,
		vy:    -rand.Float64() * 2 - 1,
		life:  1.0,
		char:  chars[rand.Intn(len(chars))],
		color: colors[rand.Intn(len(colors))],
	}
	m.particles = append(m.particles, p)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.emitting && len(m.particles) < 100 {
			for i := 0; i < 3; i++ {
				m.emitParticle()
			}
		}
		
		alive := []particle{}
		for i := range m.particles {
			p := &m.particles[i]
			
			p.vy += m.gravity
			p.vx += m.wind
			p.x += p.vx
			p.y += p.vy
			p.life -= 0.02
			
			if p.life > 0 && p.y < float64(m.height) && p.x >= 0 && p.x < float64(m.width) {
				alive = append(alive, *p)
			}
		}
		m.particles = alive
		
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.emitting = !m.emitting
		case "g":
			m.gravity = -m.gravity
		case "left":
			m.wind -= 0.05
		case "right":
			m.wind += 0.05
		case "r":
			m.particles = []particle{}
			m.gravity = 0.1
			m.wind = 0
		}
	}

	return m, nil
}

func (m model) View() string {
	grid := make([][]string, m.height-3)
	for i := range grid {
		grid[i] = make([]string, m.width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}
	
	for _, p := range m.particles {
		x := int(p.x)
		y := int(p.y)
		if y >= 0 && y < len(grid) && x >= 0 && x < m.width {
			alpha := p.life
			style := lipgloss.NewStyle().Foreground(p.color)
			if alpha < 0.5 {
				style = style.Faint(true)
			}
			grid[y][x] = style.Render(p.char)
		}
	}
	
	lines := make([]string, len(grid))
	for i, row := range grid {
		lines[i] = strings.Join(row, "")
	}
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Orange).
		Padding(0, 1)
	
	title := titleStyle.Render("✨ Particle System")
	
	statusStyle := lipgloss.NewStyle().Foreground(common.Yellow)
	status := statusStyle.Render(fmt.Sprintf("Particles: %d | Gravity: %.1f | Wind: %.1f | %s",
		len(m.particles), m.gravity, m.wind,
		map[bool]string{true: "Emitting", false: "Paused"}[m.emitting]))
	
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render("[space] toggle • [g]ravity flip • [←→] wind • [r]eset • [q]uit")
	
	return fmt.Sprintf("%s\n%s\n\n%s\n%s", title, status, strings.Join(lines, "\n"), help)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}