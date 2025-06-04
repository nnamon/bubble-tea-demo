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

type bar struct {
	height   float64
	target   float64
	peak     float64
	peakTime int
}

type model struct {
	width     int
	height    int
	bars      []bar
	time      float64
	paused    bool
	beatTime  int
	intensity float64
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
		bars:      make([]bar, 64),
		time:      0,
		paused:    false,
		intensity: 1.0,
		mode:      "music",
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
		// Adjust number of bars to fit width
		numBars := m.width / 2
		if numBars > 128 {
			numBars = 128
		}
		if numBars < 16 {
			numBars = 16
		}
		
		// Preserve existing bar data if possible
		oldBars := m.bars
		m.bars = make([]bar, numBars)
		for i := range m.bars {
			if i < len(oldBars) {
				m.bars[i] = oldBars[i]
			}
		}
		return m, nil

	case tickMsg:
		if !m.paused {
			m.time += 0.1
			
			// Simulate different audio patterns
			for i := range m.bars {
				freq := float64(i) / float64(len(m.bars))
				
				var newTarget float64
				switch m.mode {
				case "music":
					// Simulate music with bass, mids, and treble
					bass := math.Sin(m.time*0.5) * math.Exp(-freq*2)
					mids := math.Sin(m.time*1.2+freq*math.Pi) * math.Exp(-(freq-0.3)*(freq-0.3)*10)
					treble := math.Sin(m.time*2.5+freq*math.Pi*2) * math.Exp(-(freq-0.8)*(freq-0.8)*15)
					newTarget = (bass + mids + treble) * m.intensity
					
				case "bass":
					// Heavy bass emphasis
					newTarget = math.Sin(m.time*0.8) * math.Exp(-freq*4) * m.intensity * 1.5
					
				case "electronic":
					// Sharp electronic beats
					beat := math.Sin(m.time * 4)
					if beat > 0.7 {
						newTarget = (1 - freq) * m.intensity
					} else {
						newTarget = math.Sin(m.time*3+freq*math.Pi*4) * (1-freq) * m.intensity * 0.3
					}
				}
				
				// Add some randomness
				newTarget += (rand.Float64() - 0.5) * 0.2 * m.intensity
				newTarget = math.Max(0, newTarget)
				
				// Smooth movement towards target
				m.bars[i].target = newTarget
				diff := m.bars[i].target - m.bars[i].height
				m.bars[i].height += diff * 0.3
				
				// Peak detection and decay
				if m.bars[i].height > m.bars[i].peak {
					m.bars[i].peak = m.bars[i].height
					m.bars[i].peakTime = 0
				} else {
					m.bars[i].peakTime++
					if m.bars[i].peakTime > 10 {
						m.bars[i].peak *= 0.95
					}
				}
			}
			
			// Beat detection for intensity changes
			m.beatTime++
			if m.beatTime%30 == 0 {
				m.intensity = 0.5 + rand.Float64()*0.8
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
			for i := range m.bars {
				m.bars[i] = bar{}
			}
			m.time = 0
		case "1":
			m.mode = "music"
		case "2":
			m.mode = "bass"
		case "3":
			m.mode = "electronic"
		case "up":
			m.intensity = math.Min(m.intensity+0.2, 2.0)
		case "down":
			m.intensity = math.Max(m.intensity-0.2, 0.1)
		}
	}

	return m, nil
}

func (m model) View() string {
	if len(m.bars) == 0 {
		return "Initializing..."
	}
	
	// Create visualization
	lines := make([]string, m.height)
	barWidth := math.Max(1, float64(m.width)/float64(len(m.bars)))
	
	for y := 0; y < m.height; y++ {
		line := strings.Builder{}
		normalizedY := 1.0 - float64(y)/float64(m.height-1)
		
		for i, bar := range m.bars {
			x := int(float64(i) * barWidth)
			
			// Skip if we've moved past this x position
			if x >= line.Len() {
				// Fill gaps
				for line.Len() < x {
					line.WriteString(" ")
				}
				
				normalizedHeight := bar.height * 0.8 // Scale to fit nicely
				normalizedPeak := bar.peak * 0.8
				
				var char string
				var style lipgloss.Style
				
				if normalizedY <= normalizedPeak && normalizedY > normalizedPeak-0.05 {
					// Peak indicator
					char = "▄"
					style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
				} else if normalizedY <= normalizedHeight {
					// Main bar
					intensity := normalizedHeight
					if intensity > 0.8 {
						char = "█"
						style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
					} else if intensity > 0.6 {
						char = "▆"
						style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6600"))
					} else if intensity > 0.4 {
						char = "▄"
						style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
					} else if intensity > 0.2 {
						char = "▂"
						style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
					} else {
						char = "▁"
						style = lipgloss.NewStyle().Foreground(lipgloss.Color("#0088FF"))
					}
				} else {
					char = " "
					style = lipgloss.NewStyle()
				}
				
				// Fill bar width
				for w := 0; w < int(barWidth) && line.Len() < m.width; w++ {
					if w == 0 || char != " " {
						line.WriteString(style.Render(char))
					} else {
						line.WriteString(" ")
					}
				}
			}
		}
		
		// Fill remaining width
		for line.Len() < m.width {
			line.WriteString(" ")
		}
		
		lines[y] = line.String()
	}
	
	// Title and UI
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#8B008B")).
		Padding(0, 1)
	
	title := titleStyle.Render("🎵 Audio Spectrum Visualizer")
	
	statusStyle := lipgloss.NewStyle().Foreground(common.Yellow)
	status := fmt.Sprintf("Mode: %s | Intensity: %.1f | Bars: %d | %s",
		strings.Title(m.mode), m.intensity, len(m.bars),
		map[bool]string{true: "⏸ Paused", false: "🎶 Playing"}[m.paused])
	
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := "[space] pause • [1]music [2]bass [3]electronic • [↑↓] intensity • [r]eset • [q]uit"
	
	return fmt.Sprintf("%s\n%s\n\n%s\n%s", title, statusStyle.Render(status),
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