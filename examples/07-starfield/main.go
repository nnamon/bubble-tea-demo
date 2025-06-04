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

type star struct {
	x, y, z float64
	prevX, prevY float64
}

type model struct {
	width     int
	height    int
	stars     []star
	speed     float64
	centerX   float64
	centerY   float64
	paused    bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	m := model{
		width:   80,
		height:  24,
		speed:   0.05,
		paused:  false,
	}
	m.centerX = float64(m.width) / 2
	m.centerY = float64(m.height) / 2
	m.initStars()
	return m
}

func (m *model) initStars() {
	m.stars = make([]star, 200)
	for i := range m.stars {
		m.stars[i] = star{
			x: (rand.Float64() - 0.5) * 2,
			y: (rand.Float64() - 0.5) * 2,
			z: rand.Float64(),
		}
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
		m.centerX = float64(m.width) / 2
		m.centerY = float64(m.height) / 2
		return m, nil

	case tickMsg:
		if !m.paused {
			for i := range m.stars {
				star := &m.stars[i]
				
				// Store previous position for trail effect
				star.prevX = star.x / star.z * m.centerX + m.centerX
				star.prevY = star.y / star.z * m.centerY + m.centerY
				
				// Move star towards viewer
				star.z -= m.speed
				
				// Reset star if it's too close
				if star.z <= 0 {
					star.x = (rand.Float64() - 0.5) * 2
					star.y = (rand.Float64() - 0.5) * 2
					star.z = 1.0
				}
			}
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "r":
			m.initStars()
		case "up":
			m.speed = math.Min(m.speed+0.01, 0.2)
		case "down":
			m.speed = math.Max(m.speed-0.01, 0.01)
		case "+", "=":
			m.speed = math.Min(m.speed+0.02, 0.3)
		case "-":
			m.speed = math.Max(m.speed-0.02, 0.005)
		}
	}

	return m, nil
}

func (m model) View() string {
	// Create grid
	grid := make([][]string, m.height)
	for i := range grid {
		grid[i] = make([]string, m.width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}
	
	// Draw stars
	for _, star := range m.stars {
		// Calculate screen position
		screenX := star.x / star.z * m.centerX + m.centerX
		screenY := star.y / star.z * m.centerY + m.centerY
		
		x, y := int(screenX), int(screenY)
		
		// Only draw if on screen
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			// Choose character and color based on distance
			brightness := 1.0 - star.z
			
			var char string
			var color lipgloss.Color
			
			if brightness > 0.9 {
				char = "✦"
				color = lipgloss.Color("#FFFFFF")
			} else if brightness > 0.8 {
				char = "★"
				color = lipgloss.Color("#FFFF99")
			} else if brightness > 0.6 {
				char = "✧"
				color = lipgloss.Color("#CCCCCC")
			} else if brightness > 0.4 {
				char = "•"
				color = lipgloss.Color("#999999")
			} else if brightness > 0.2 {
				char = "∘"
				color = lipgloss.Color("#666666")
			} else {
				char = "·"
				color = lipgloss.Color("#444444")
			}
			
			style := lipgloss.NewStyle().Foreground(color)
			if brightness > 0.8 {
				style = style.Bold(true)
			} else if brightness < 0.3 {
				style = style.Faint(true)
			}
			
			grid[y][x] = style.Render(char)
			
			// Draw trail for fast-moving stars
			if m.speed > 0.08 && brightness > 0.5 {
				prevX, prevY := int(star.prevX), int(star.prevY)
				if prevX >= 0 && prevX < m.width && prevY >= 0 && prevY < m.height &&
					(prevX != x || prevY != y) {
					trailStyle := lipgloss.NewStyle().Foreground(color).Faint(true)
					if grid[prevY][prevX] == " " {
						grid[prevY][prevX] = trailStyle.Render("·")
					}
				}
			}
		}
	}
	
	// Render grid
	lines := make([]string, len(grid))
	for i, row := range grid {
		lines[i] = strings.Join(row, "")
	}
	
	// Title and UI
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#000080")).
		Padding(0, 1)
	
	title := titleStyle.Render("⭐ 3D Starfield")
	
	statusStyle := lipgloss.NewStyle().Foreground(common.Cyan)
	status := fmt.Sprintf("Speed: %.3f | Stars: %d | %s",
		m.speed, len(m.stars),
		map[bool]string{true: "⏸ Paused", false: "🚀 Warping"}[m.paused])
	
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := "[space] pause • [↑↓] speed • [+/-] turbo • [r]eset • [q]uit"
	
	return fmt.Sprintf("%s  %s\n\n%s\n%s", title, statusStyle.Render(status),
		strings.Join(lines, "\n"), helpStyle.Render(help))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}