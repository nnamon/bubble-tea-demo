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
	waves      []wave
	showHelp   bool
}

type wave struct {
	amplitude  float64
	frequency  float64
	phase      float64
	speed      float64
	color      lipgloss.Color
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
		time:     0,
		showHelp: true,
		waves: []wave{
			{amplitude: 0.3, frequency: 0.05, phase: 0, speed: 0.05, color: common.Blue},
			{amplitude: 0.2, frequency: 0.08, phase: math.Pi/3, speed: 0.08, color: common.Cyan},
			{amplitude: 0.25, frequency: 0.03, phase: math.Pi/2, speed: 0.03, color: common.Purple},
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
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.time += 0.05
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h":
			m.showHelp = !m.showHelp
		case "r":
			m.time = 0
		case "space":
			if len(m.waves) < 5 {
				m.waves = append(m.waves, wave{
					amplitude: 0.1 + math.Mod(m.time, 0.3),
					frequency: 0.02 + math.Mod(m.time, 0.08),
					phase:     m.time,
					speed:     0.02 + math.Mod(m.time, 0.06),
					color:     lipgloss.Color(common.GradientBlue[int(m.time)%len(common.GradientBlue)]),
				})
			}
		case "backspace":
			if len(m.waves) > 1 {
				m.waves = m.waves[:len(m.waves)-1]
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	lines := make([]string, m.height-4)
	
	for y := range lines {
		line := strings.Builder{}
		normalizedY := float64(y) / float64(len(lines)-1)
		
		for x := 0; x < m.width; x++ {
			normalizedX := float64(x) / float64(m.width-1)
			
			height := 0.5
			for _, w := range m.waves {
				waveHeight := w.amplitude * math.Sin(2*math.Pi*(w.frequency*normalizedX+w.speed*m.time)+w.phase)
				height += waveHeight
			}
			
			if math.Abs(normalizedY-(0.5-height/2)) < 0.05 {
				colorIndex := int((height + 1) * float64(len(common.GradientBlue)-1) / 2)
				colorIndex = int(common.Clamp(float64(colorIndex), 0, float64(len(common.GradientBlue)-1)))
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(common.GradientBlue[colorIndex]))
				line.WriteString(style.Render("█"))
			} else if normalizedY > (0.5 - height/2) {
				waterChar := "░"
				if math.Mod(float64(x)+m.time*10, 3) < 1 {
					waterChar = "▒"
				}
				style := lipgloss.NewStyle().Foreground(common.Blue).Faint(true)
				line.WriteString(style.Render(waterChar))
			} else {
				line.WriteString(" ")
			}
		}
		lines[y] = line.String()
	}
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Blue).
		Padding(0, 1)
	
	title := titleStyle.Render("🌊 Wave Animation")
	
	help := ""
	if m.showHelp {
		helpStyle := lipgloss.NewStyle().Faint(true)
		help = helpStyle.Render("\n[h]ide help • [space] add wave • [backspace] remove • [r]eset • [q]uit")
	} else {
		help = "\n[h] show help"
	}
	
	countStyle := lipgloss.NewStyle().Foreground(common.Cyan)
	count := countStyle.Render(fmt.Sprintf("Waves: %d", len(m.waves)))
	
	return fmt.Sprintf("%s  %s\n\n%s%s", title, count, strings.Join(lines, "\n"), help)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}