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

type point3D struct {
	x, y, z float64
}

type edge struct {
	start, end int
}

type model struct {
	width       int
	height      int
	vertices    []point3D
	edges       []edge
	rotationX   float64
	rotationY   float64
	rotationZ   float64
	scale       float64
	autoRotate  bool
	perspective float64
	paused      bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	// Define cube vertices
	vertices := []point3D{
		{-1, -1, -1}, {1, -1, -1}, {1, 1, -1}, {-1, 1, -1}, // Back face
		{-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1},     // Front face
	}

	// Define cube edges
	edges := []edge{
		// Back face
		{0, 1}, {1, 2}, {2, 3}, {3, 0},
		// Front face
		{4, 5}, {5, 6}, {6, 7}, {7, 4},
		// Connecting edges
		{0, 4}, {1, 5}, {2, 6}, {3, 7},
	}

	return model{
		width:       80,
		height:      24,
		vertices:    vertices,
		edges:       edges,
		scale:       8,
		autoRotate:  true,
		perspective: 4,
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
		if !m.paused && m.autoRotate {
			m.rotationX += 0.02
			m.rotationY += 0.03
			m.rotationZ += 0.01
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "a":
			m.autoRotate = !m.autoRotate
		case "r":
			m.rotationX = 0
			m.rotationY = 0
			m.rotationZ = 0
		case "up":
			if !m.autoRotate {
				m.rotationX -= 0.1
			}
		case "down":
			if !m.autoRotate {
				m.rotationX += 0.1
			}
		case "left":
			if !m.autoRotate {
				m.rotationY -= 0.1
			}
		case "right":
			if !m.autoRotate {
				m.rotationY += 0.1
			}
		case "+", "=":
			m.scale = math.Min(m.scale+1, 20)
		case "-":
			m.scale = math.Max(m.scale-1, 2)
		case "z":
			if !m.autoRotate {
				m.rotationZ -= 0.1
			}
		case "x":
			if !m.autoRotate {
				m.rotationZ += 0.1
			}
		case "p":
			m.perspective = math.Max(m.perspective-0.5, 1)
		case "o":
			m.perspective = math.Min(m.perspective+0.5, 10)
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Purple).
		Padding(0, 1)

	title := titleStyle.Render("🎲 3D Rotating Cube")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Yellow)
	status := statusStyle.Render(fmt.Sprintf(
		"Scale: %.0f | Perspective: %.1f | %s | %s",
		m.scale, m.perspective,
		map[bool]string{true: "Auto-rotating", false: "Manual control"}[m.autoRotate],
		map[bool]string{true: "⏸ Paused", false: "🎲 Spinning"}[m.paused],
	))

	// Create 3D visualization
	lines := m.render3D()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	var help string
	if m.autoRotate {
		help = "[a] manual control • [space] pause • [+/-] scale • [p/o] perspective • [r]eset • [q]uit"
	} else {
		help = "[a] auto-rotate • [↑↓←→] rotate • [z/x] roll • [+/-] scale • [p/o] perspective • [r]eset • [q]uit"
	}

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), helpStyle.Render(help))
}

func (m model) render3D() []string {
	// Create screen buffer
	screen := make([][]string, m.height)
	for i := range screen {
		screen[i] = make([]string, m.width)
		for j := range screen[i] {
			screen[i][j] = " "
		}
	}

	// Transform vertices
	transformed := make([]point3D, len(m.vertices))
	for i, v := range m.vertices {
		// Apply rotations
		transformed[i] = m.rotatePoint(v)
	}

	// Project to 2D and draw edges
	projected := make([][2]int, len(transformed))
	for i, v := range transformed {
		projected[i] = m.project(v)
	}

	// Draw all edges
	for _, edge := range m.edges {
		start := projected[edge.start]
		end := projected[edge.end]
		m.drawLine(screen, start[0], start[1], end[0], end[1])
	}

	// Draw vertices as dots
	for i, p := range projected {
		x, y := p[0], p[1]
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			// Different colors for front and back vertices
			vertex := transformed[i]
			var style lipgloss.Style
			if vertex.z > 0 {
				style = lipgloss.NewStyle().Foreground(common.Red).Bold(true)
			} else {
				style = lipgloss.NewStyle().Foreground(common.Blue)
			}
			screen[y][x] = style.Render("●")
		}
	}

	// Convert screen buffer to strings
	lines := make([]string, len(screen))
	for i, row := range screen {
		lines[i] = strings.Join(row, "")
	}

	return lines
}

func (m model) rotatePoint(p point3D) point3D {
	// Rotate around X axis
	cosX, sinX := math.Cos(m.rotationX), math.Sin(m.rotationX)
	y1 := p.y*cosX - p.z*sinX
	z1 := p.y*sinX + p.z*cosX

	// Rotate around Y axis
	cosY, sinY := math.Cos(m.rotationY), math.Sin(m.rotationY)
	x2 := p.x*cosY + z1*sinY
	z2 := -p.x*sinY + z1*cosY

	// Rotate around Z axis
	cosZ, sinZ := math.Cos(m.rotationZ), math.Sin(m.rotationZ)
	x3 := x2*cosZ - y1*sinZ
	y3 := x2*sinZ + y1*cosZ

	return point3D{x3, y3, z2}
}

func (m model) project(p point3D) [2]int {
	// Perspective projection
	distance := m.perspective + p.z
	if distance <= 0.1 {
		distance = 0.1
	}

	// Project to screen coordinates
	screenX := (p.x * m.scale / distance) + float64(m.width)/2
	screenY := (-p.y * m.scale / distance) + float64(m.height)/2

	return [2]int{int(screenX), int(screenY)}
}

func (m model) drawLine(screen [][]string, x0, y0, x1, y1 int) {
	// Bresenham's line algorithm
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := sign(x1 - x0)
	sy := sign(y1 - y0)
	err := dx - dy

	x, y := x0, y0

	for {
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			// Choose character based on line direction
			char := m.getLineChar(x0, y0, x1, y1, x, y)
			style := lipgloss.NewStyle().Foreground(common.Green)
			screen[y][x] = style.Render(char)
		}

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func (m model) getLineChar(x0, y0, x1, y1, x, y int) string {
	// Determine line character based on direction
	dx := x1 - x0
	dy := y1 - y0

	if abs(dx) > abs(dy) {
		// More horizontal
		return "─"
	} else if abs(dy) > abs(dx) {
		// More vertical
		return "│"
	} else {
		// Diagonal
		if (dx > 0 && dy > 0) || (dx < 0 && dy < 0) {
			return "╲"
		} else {
			return "╱"
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}