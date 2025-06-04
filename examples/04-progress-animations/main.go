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

type progressBar struct {
	name     string
	progress float64
	speed    float64
	style    string
}

type model struct {
	bars   []progressBar
	width  int
	paused bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width: 40,
		bars: []progressBar{
			{name: "Classic", progress: 0, speed: 0.01, style: "classic"},
			{name: "Smooth", progress: 0, speed: 0.015, style: "smooth"},
			{name: "Gradient", progress: 0, speed: 0.012, style: "gradient"},
			{name: "Pulse", progress: 0, speed: 0.018, style: "pulse"},
			{name: "Wave", progress: 0, speed: 0.02, style: "wave"},
			{name: "Blocks", progress: 0, speed: 0.008, style: "blocks"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !m.paused {
			for i := range m.bars {
				m.bars[i].progress += m.bars[i].speed
				if m.bars[i].progress > 1 {
					m.bars[i].progress = 0
				}
			}
		}
		return m, tick()

	case tea.WindowSizeMsg:
		m.width = min(msg.Width-20, 60)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "r":
			for i := range m.bars {
				m.bars[i].progress = 0
			}
		}
	}

	return m, nil
}

func (m model) renderBar(bar progressBar) string {
	filled := int(bar.progress * float64(m.width))
	
	switch bar.style {
	case "classic":
		return m.renderClassic(filled)
	case "smooth":
		return m.renderSmooth(filled, bar.progress)
	case "gradient":
		return m.renderGradient(filled)
	case "pulse":
		return m.renderPulse(filled, bar.progress)
	case "wave":
		return m.renderWave(filled, bar.progress)
	case "blocks":
		return m.renderBlocks(filled)
	default:
		return ""
	}
}

func (m model) renderClassic(filled int) string {
	bar := strings.Repeat("█", filled) + strings.Repeat("░", m.width-filled)
	return lipgloss.NewStyle().Foreground(common.Blue).Render(bar)
}

func (m model) renderSmooth(filled int, progress float64) string {
	chars := []string{"░", "▒", "▓", "█"}
	bar := strings.Builder{}
	
	for i := 0; i < m.width; i++ {
		if i < filled {
			bar.WriteString("█")
		} else if i == filled {
			subProgress := (progress * float64(m.width)) - float64(filled)
			charIndex := int(subProgress * float64(len(chars)-1))
			bar.WriteString(chars[charIndex])
		} else {
			bar.WriteString("░")
		}
	}
	
	return lipgloss.NewStyle().Foreground(common.Green).Render(bar.String())
}

func (m model) renderGradient(filled int) string {
	bar := strings.Builder{}
	gradient := common.GradientFire
	
	for i := 0; i < m.width; i++ {
		colorIndex := int(float64(i) / float64(m.width) * float64(len(gradient)-1))
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(gradient[colorIndex]))
		
		if i < filled {
			bar.WriteString(style.Render("█"))
		} else {
			bar.WriteString(style.Faint(true).Render("░"))
		}
	}
	
	return bar.String()
}

func (m model) renderPulse(filled int, progress float64) string {
	bar := strings.Builder{}
	pulseIntensity := (math.Sin(progress*math.Pi*2) + 1) / 2
	
	for i := 0; i < m.width; i++ {
		if i < filled {
			alpha := 0.5 + pulseIntensity*0.5
			if alpha > 0.7 {
				bar.WriteString(lipgloss.NewStyle().Foreground(common.Purple).Render("█"))
			} else {
				bar.WriteString(lipgloss.NewStyle().Foreground(common.Purple).Faint(true).Render("█"))
			}
		} else {
			bar.WriteString("░")
		}
	}
	
	return bar.String()
}

func (m model) renderWave(filled int, progress float64) string {
	bar := strings.Builder{}
	waveChars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	
	for i := 0; i < m.width; i++ {
		if i < filled {
			waveHeight := (math.Sin(float64(i)*0.3 + progress*10) + 1) / 2
			charIndex := int(waveHeight * float64(len(waveChars)-1))
			bar.WriteString(lipgloss.NewStyle().Foreground(common.Cyan).Render(waveChars[charIndex]))
		} else {
			bar.WriteString(" ")
		}
	}
	
	return bar.String()
}

func (m model) renderBlocks(filled int) string {
	bar := strings.Builder{}
	blockChars := []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
	
	fullBlocks := filled / len(blockChars)
	remainder := filled % len(blockChars)
	
	bar.WriteString(strings.Repeat("█", fullBlocks))
	
	if remainder > 0 && fullBlocks < m.width {
		bar.WriteString(blockChars[remainder-1])
	}
	
	empty := m.width - fullBlocks
	if remainder > 0 {
		empty--
	}
	bar.WriteString(strings.Repeat(" ", max(0, empty)))
	
	return lipgloss.NewStyle().Foreground(common.Orange).Render(bar.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Purple).
		Padding(0, 1)
	
	content := titleStyle.Render("📊 Progress Bar Animations") + "\n\n"
	
	nameStyle := lipgloss.NewStyle().
		Width(10).
		Foreground(common.Yellow)
	
	percentStyle := lipgloss.NewStyle().
		Width(5).
		Align(lipgloss.Right).
		Foreground(common.Green)
	
	for _, bar := range m.bars {
		name := nameStyle.Render(bar.name)
		percent := percentStyle.Render(fmt.Sprintf("%3.0f%%", bar.progress*100))
		barRender := m.renderBar(bar)
		
		content += fmt.Sprintf("%s %s %s\n\n", name, barRender, percent)
	}
	
	statusStyle := lipgloss.NewStyle().Foreground(common.Cyan)
	status := "▶ Playing"
	if m.paused {
		status = "⏸ Paused"
	}
	content += statusStyle.Render(status) + "\n"
	
	helpStyle := lipgloss.NewStyle().Faint(true)
	content += helpStyle.Render("[space] pause/play • [r]eset • [q]uit")
	
	return content
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}