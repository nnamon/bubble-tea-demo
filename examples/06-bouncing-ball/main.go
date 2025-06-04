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

type ball struct {
	x, y   float64
	vx, vy float64
	char   string
	color  lipgloss.Color
	trail  []position
}

type position struct {
	x, y  float64
	age   int
	color lipgloss.Color
}

type model struct {
	width    int
	height   int
	balls    []ball
	gravity  float64
	friction float64
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
		width:    80,
		height:   24,
		gravity:  0.5,
		friction: 0.98,
		balls: []ball{
			{
				x: 40, y: 10, vx: 2, vy: 0,
				char: "●", color: common.Red,
				trail: []position{},
			},
		},
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
			for i := range m.balls {
				ball := &m.balls[i]
				
				// Add current position to trail
				ball.trail = append(ball.trail, position{
					x: ball.x, y: ball.y, age: 0,
					color: ball.color,
				})
				
				// Age trail positions and remove old ones
				newTrail := []position{}
				for _, pos := range ball.trail {
					if pos.age < 10 {
						pos.age++
						newTrail = append(newTrail, pos)
					}
				}
				ball.trail = newTrail
				
				// Apply gravity
				ball.vy += m.gravity
				
				// Update position
				ball.x += ball.vx
				ball.y += ball.vy
				
				// Bounce off walls
				if ball.x <= 0 || ball.x >= float64(m.width-1) {
					ball.vx = -ball.vx * m.friction
					ball.x = math.Max(0, math.Min(float64(m.width-1), ball.x))
				}
				
				if ball.y <= 0 {
					ball.vy = -ball.vy * m.friction
					ball.y = 0
				}
				
				// Bounce off floor with some energy loss
				if ball.y >= float64(m.height-1) {
					ball.vy = -ball.vy * m.friction
					ball.vx *= m.friction
					ball.y = float64(m.height - 1)
					
					// Add some randomness to prevent settling
					if math.Abs(ball.vy) < 0.5 {
						ball.vy = -2
					}
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
			return initialModel(), nil
		case "g":
			m.gravity = -m.gravity
		case "up":
			if len(m.balls) > 0 {
				m.balls[0].vy -= 3
			}
		case "left":
			if len(m.balls) > 0 {
				m.balls[0].vx -= 1
			}
		case "right":
			if len(m.balls) > 0 {
				m.balls[0].vx += 1
			}
		case "a":
			// Add new ball
			if len(m.balls) < 5 {
				colors := []lipgloss.Color{common.Red, common.Blue, common.Green, common.Yellow, common.Purple}
				chars := []string{"●", "○", "◉", "⬤", "🔴"}
				newBall := ball{
					x: float64(m.width) / 2, y: 5,
					vx: (float64(len(m.balls))-2.5) * 0.8, vy: 0,
					char:  chars[len(m.balls)%len(chars)],
					color: colors[len(m.balls)%len(colors)],
					trail: []position{},
				}
				m.balls = append(m.balls, newBall)
			}
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
	
	// Draw trails
	for _, ball := range m.balls {
		for _, pos := range ball.trail {
			x, y := int(pos.x), int(pos.y)
			if y >= 0 && y < m.height && x >= 0 && x < m.width {
				alpha := float64(10-pos.age) / 10.0
				char := "·"
				if alpha > 0.7 {
					char = "•"
				} else if alpha > 0.4 {
					char = "∘"
				}
				
				style := lipgloss.NewStyle().Foreground(pos.color)
				if alpha < 0.5 {
					style = style.Faint(true)
				}
				grid[y][x] = style.Render(char)
			}
		}
	}
	
	// Draw balls
	for _, ball := range m.balls {
		x, y := int(ball.x), int(ball.y)
		if y >= 0 && y < m.height && x >= 0 && x < m.width {
			style := lipgloss.NewStyle().Foreground(ball.color).Bold(true)
			grid[y][x] = style.Render(ball.char)
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
		Background(common.Red).
		Padding(0, 1)
	
	title := titleStyle.Render("🏀 Bouncing Ball Physics")
	
	statusStyle := lipgloss.NewStyle().Foreground(common.Cyan)
	status := fmt.Sprintf("Balls: %d | Gravity: %.1f | %s",
		len(m.balls), m.gravity,
		map[bool]string{true: "⏸ Paused", false: "▶ Playing"}[m.paused])
	
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := "[space] pause • [↑←→] control • [a]dd ball • [g]ravity flip • [r]eset • [q]uit"
	
	return fmt.Sprintf("%s  %s\n\n%s\n%s", title, statusStyle.Render(status), 
		strings.Join(lines, "\n"), helpStyle.Render(help))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}