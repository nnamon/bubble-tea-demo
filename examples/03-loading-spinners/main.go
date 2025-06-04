package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type spinner struct {
	name   string
	frames []string
	index  int
	color  lipgloss.Color
}

type model struct {
	spinners []spinner
	ticks    int
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		spinners: []spinner{
			{
				name:   "Dots",
				frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
				color:  common.Blue,
			},
			{
				name:   "Line",
				frames: []string{"-", "\\", "|", "/"},
				color:  common.Green,
			},
			{
				name:   "Circle",
				frames: []string{"◐", "◓", "◑", "◒"},
				color:  common.Yellow,
			},
			{
				name:   "Square",
				frames: []string{"◰", "◳", "◲", "◱"},
				color:  common.Red,
			},
			{
				name:   "Triangle",
				frames: []string{"◢", "◣", "◤", "◥"},
				color:  common.Purple,
			},
			{
				name:   "Box",
				frames: []string{"▖", "▘", "▝", "▗"},
				color:  common.Cyan,
			},
			{
				name:   "Arc",
				frames: []string{"◜", "◠", "◝", "◞", "◡", "◟"},
				color:  common.Orange,
			},
			{
				name:   "Bounce",
				frames: []string{"⠁", "⠂", "⠄", "⠂"},
				color:  common.Pink,
			},
			{
				name:   "Pulse",
				frames: []string{"▁", "▃", "▄", "▅", "▆", "▇", "▆", "▅", "▄", "▃"},
				color:  common.Blue,
			},
			{
				name:   "Points",
				frames: []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"},
				color:  common.Green,
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.ticks++
		for i := range m.spinners {
			if m.ticks%(i+1) == 0 {
				m.spinners[i].index = (m.spinners[i].index + 1) % len(m.spinners[i].frames)
			}
		}
		return m, tick()

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Green).
		Padding(0, 1).
		MarginBottom(1)
	
	spinnerStyle := lipgloss.NewStyle().
		Width(18).
		Height(5).
		Padding(1, 2).
		MarginRight(1).
		MarginBottom(1).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center)
	
	content := titleStyle.Render("🔄 Loading Spinners Gallery") + "\n\n"
	
	var rows []string
	var currentRow []string
	
	for i, s := range m.spinners {
		frame := s.frames[s.index]
		
		style := spinnerStyle.BorderForeground(s.color)
		spinnerContent := lipgloss.NewStyle().
			Foreground(s.color).
			Bold(true).
			Render(frame)
		
		name := lipgloss.NewStyle().
			Foreground(s.color).
			Faint(true).
			Render(s.name)
		
		box := style.Render(fmt.Sprintf("%s\n\n%s", spinnerContent, name))
		currentRow = append(currentRow, box)
		
		if (i+1)%4 == 0 || i == len(m.spinners)-1 {
			row := lipgloss.JoinHorizontal(lipgloss.Top, currentRow...)
			rows = append(rows, row)
			currentRow = []string{}
		}
	}
	
	content += lipgloss.JoinVertical(lipgloss.Left, rows...)
	
	helpStyle := lipgloss.NewStyle().Faint(true).MarginTop(2)
	content += "\n\n" + helpStyle.Render("Press [q] to quit")
	
	return content
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}