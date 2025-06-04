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

type cell struct {
	alive bool
	age   int
}

type model struct {
	width      int
	height     int
	grid       [][]cell
	generation int
	speed      time.Duration
	paused     bool
	pattern    string
}

type tickMsg time.Time

func tick(speed time.Duration) tea.Cmd {
	return tea.Tick(speed, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width:   80,
		height:  24,
		speed:   time.Millisecond * 200,
		pattern: "random",
	}
}

func (m *model) initGrid() {
	m.grid = make([][]cell, m.height)
	for i := range m.grid {
		m.grid[i] = make([]cell, m.width)
	}
	m.generation = 0
}

func (m model) Init() tea.Cmd {
	return tick(m.speed)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 4
		m.initGrid()
		m.seedPattern()
		return m, nil

	case tickMsg:
		if !m.paused {
			m.nextGeneration()
		}
		return m, tick(m.speed)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "r":
			m.initGrid()
			m.seedPattern()
		case "1":
			m.pattern = "random"
			m.initGrid()
			m.seedPattern()
		case "2":
			m.pattern = "glider"
			m.initGrid()
			m.seedPattern()
		case "3":
			m.pattern = "oscillator"
			m.initGrid()
			m.seedPattern()
		case "4":
			m.pattern = "spaceship"
			m.initGrid()
			m.seedPattern()
		case "5":
			m.pattern = "gosper"
			m.initGrid()
			m.seedPattern()
		case "up":
			m.speed = time.Duration(float64(m.speed) * 0.8)
			if m.speed < time.Millisecond*50 {
				m.speed = time.Millisecond * 50
			}
		case "down":
			m.speed = time.Duration(float64(m.speed) * 1.2)
			if m.speed > time.Second {
				m.speed = time.Second
			}
		}
	}

	return m, nil
}

func (m *model) seedPattern() {
	switch m.pattern {
	case "random":
		m.seedRandom()
	case "glider":
		m.seedGlider()
	case "oscillator":
		m.seedOscillator()
	case "spaceship":
		m.seedSpaceship()
	case "gosper":
		m.seedGosperGun()
	}
}

func (m *model) seedRandom() {
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if rand.Float64() < 0.3 {
				m.grid[y][x].alive = true
			}
		}
	}
}

func (m *model) seedGlider() {
	// Place multiple gliders
	patterns := []struct{ x, y int }{
		{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, // Glider 1
	}

	for i := 0; i < 3; i++ {
		offsetX := i * 20 + 5
		offsetY := i * 8 + 5
		
		for _, p := range patterns {
			x, y := offsetX+p.x, offsetY+p.y
			if x < m.width && y < m.height {
				m.grid[y][x].alive = true
			}
		}
	}
}

func (m *model) seedOscillator() {
	centerX, centerY := m.width/2, m.height/2
	
	// Blinker (period 2)
	for i := -1; i <= 1; i++ {
		if centerX+i >= 0 && centerX+i < m.width {
			m.grid[centerY][centerX+i].alive = true
		}
	}
	
	// Toad (period 2)
	offsetY := centerY - 5
	for i := 0; i < 3; i++ {
		if centerX+i >= 0 && centerX+i < m.width && offsetY >= 0 {
			m.grid[offsetY][centerX+i].alive = true
		}
		if centerX+i-1 >= 0 && centerX+i-1 < m.width && offsetY+1 < m.height {
			m.grid[offsetY+1][centerX+i-1].alive = true
		}
	}
	
	// Beacon (period 2)
	offsetY = centerY + 5
	beaconPattern := []struct{ x, y int }{
		{0, 0}, {1, 0}, {0, 1}, {3, 2}, {2, 3}, {3, 3},
	}
	for _, p := range beaconPattern {
		x, y := centerX+p.x-2, offsetY+p.y
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			m.grid[y][x].alive = true
		}
	}
}

func (m *model) seedSpaceship() {
	// Lightweight spaceship (LWSS)
	centerX, centerY := 10, m.height/2
	lwssPattern := []struct{ x, y int }{
		{1, 0}, {4, 0}, {0, 1}, {0, 2}, {4, 2}, {0, 3}, {1, 3}, {2, 3}, {3, 3},
	}
	
	for _, p := range lwssPattern {
		x, y := centerX+p.x, centerY+p.y
		if x < m.width && y < m.height {
			m.grid[y][x].alive = true
		}
	}
}

func (m *model) seedGosperGun() {
	// Gosper Glider Gun (simplified version)
	if m.width < 40 || m.height < 15 {
		m.seedRandom()
		return
	}
	
	gun := []struct{ x, y int }{
		// Left block
		{1, 5}, {1, 6}, {2, 5}, {2, 6},
		// Left part
		{11, 5}, {11, 6}, {11, 7}, {12, 4}, {12, 8}, {13, 3}, {13, 9}, {14, 3}, {14, 9},
		{15, 6}, {16, 4}, {16, 8}, {17, 5}, {17, 6}, {17, 7}, {18, 6},
		// Right part
		{21, 3}, {21, 4}, {21, 5}, {22, 3}, {22, 4}, {22, 5}, {23, 2}, {23, 6},
		{25, 1}, {25, 2}, {25, 6}, {25, 7},
		// Right block
		{35, 3}, {35, 4}, {36, 3}, {36, 4},
	}
	
	for _, p := range gun {
		if p.x < m.width && p.y < m.height {
			m.grid[p.y][p.x].alive = true
		}
	}
}

func (m *model) nextGeneration() {
	if len(m.grid) == 0 {
		return
	}
	
	newGrid := make([][]cell, m.height)
	for i := range newGrid {
		newGrid[i] = make([]cell, m.width)
	}
	
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			neighbors := m.countNeighbors(x, y)
			currentCell := m.grid[y][x]
			
			if currentCell.alive {
				// Survival rules
				newGrid[y][x].alive = neighbors == 2 || neighbors == 3
				if newGrid[y][x].alive {
					newGrid[y][x].age = currentCell.age + 1
				}
			} else {
				// Birth rule
				newGrid[y][x].alive = neighbors == 3
				if newGrid[y][x].alive {
					newGrid[y][x].age = 0
				}
			}
		}
	}
	
	m.grid = newGrid
	m.generation++
}

func (m model) countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < m.width && ny >= 0 && ny < m.height {
				if m.grid[ny][nx].alive {
					count++
				}
			}
		}
	}
	return count
}

func (m model) View() string {
	if len(m.grid) == 0 {
		return "Initializing Conway's Game of Life..."
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Green).
		Padding(0, 1)

	title := titleStyle.Render("🧬 Conway's Game of Life")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Yellow)
	population := m.countPopulation()
	status := statusStyle.Render(fmt.Sprintf(
		"Generation: %d | Population: %d | Pattern: %s | Speed: %dms | %s",
		m.generation, population, strings.Title(m.pattern), 
		m.speed.Milliseconds(),
		map[bool]string{true: "⏸ Paused", false: "🧬 Evolving"}[m.paused],
	))

	// Render grid
	lines := make([]string, m.height)
	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		for x := 0; x < m.width; x++ {
			char, color := m.getCellChar(m.grid[y][x])
			style := lipgloss.NewStyle().Foreground(color)
			line.WriteString(style.Render(char))
		}
		lines[y] = line.String()
	}

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1]random [2]glider [3]oscillator [4]spaceship [5]gosper gun • [↑↓] speed • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) getCellChar(c cell) (string, lipgloss.Color) {
	if !c.alive {
		return " ", lipgloss.Color("#000000")
	}
	
	// Color cells based on age
	if c.age < 5 {
		return "●", common.Green
	} else if c.age < 15 {
		return "●", common.Yellow
	} else if c.age < 30 {
		return "●", common.Orange
	} else {
		return "●", common.Red
	}
}

func (m model) countPopulation() int {
	count := 0
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if m.grid[y][x].alive {
				count++
			}
		}
	}
	return count
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}