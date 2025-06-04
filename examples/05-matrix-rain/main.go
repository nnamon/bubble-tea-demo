package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type column struct {
	chars    []rune
	position int
	speed    int
	length   int
}

type model struct {
	width   int
	height  int
	columns []column
	tick    int
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width:   80,
		height:  24,
		columns: []column{},
	}
}

func (m *model) initColumns() {
	m.columns = make([]column, m.width)
	chars := []rune("ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ0123456789")
	
	for i := range m.columns {
		length := rand.Intn(m.height/2) + 5
		col := column{
			chars:    make([]rune, m.height),
			position: -rand.Intn(m.height),
			speed:    rand.Intn(3) + 1,
			length:   length,
		}
		
		for j := range col.chars {
			col.chars[j] = chars[rand.Intn(len(chars))]
		}
		
		m.columns[i] = col
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initColumns()
		return m, nil

	case tickMsg:
		m.tick++
		chars := []rune("ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ0123456789")
		
		for i := range m.columns {
			if m.tick%m.columns[i].speed == 0 {
				m.columns[i].position++
				
				if m.columns[i].position-m.columns[i].length > m.height {
					m.columns[i].position = -rand.Intn(m.height)
					m.columns[i].speed = rand.Intn(3) + 1
					m.columns[i].length = rand.Intn(m.height/2) + 5
				}
				
				if rand.Float64() < 0.1 {
					changePos := rand.Intn(m.height)
					m.columns[i].chars[changePos] = chars[rand.Intn(len(chars))]
				}
			}
		}
		
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.initColumns()
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	
	grid := make([][]string, m.height)
	for i := range grid {
		grid[i] = make([]string, m.width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}
	
	greenShades := []string{"#00FF00", "#00CC00", "#009900", "#006600", "#003300"}
	
	for col, column := range m.columns {
		for row := 0; row < m.height; row++ {
			if row >= column.position-column.length && row < column.position {
				distance := column.position - row
				colorIndex := distance * len(greenShades) / column.length
				if colorIndex >= len(greenShades) {
					colorIndex = len(greenShades) - 1
				}
				
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(greenShades[colorIndex]))
				
				if distance == 1 {
					style = lipgloss.NewStyle().
						Foreground(lipgloss.Color("#FFFFFF")).
						Bold(true)
				}
				
				if row >= 0 && row < m.height && col < m.width {
					grid[row][col] = style.Render(string(column.chars[row]))
				}
			}
		}
	}
	
	lines := make([]string, len(grid))
	for i, row := range grid {
		lines[i] = strings.Join(row, "")
	}
	
	return strings.Join(lines, "\n")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}