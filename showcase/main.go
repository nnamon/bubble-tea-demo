package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title       string
	description string
	command     string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	choice string
}

func initialModel() model {
	items := []list.Item{
		item{
			title:       "🌊 Wave Animation",
			description: "Smooth sine wave animations with multiple layers",
			command:     "examples/01-wave-animation/main.go",
		},
		item{
			title:       "✨ Particle System",
			description: "Dynamic particle effects with physics simulation",
			command:     "examples/02-particle-system/main.go",
		},
		item{
			title:       "🔄 Loading Spinners",
			description: "Collection of various animated loading indicators",
			command:     "examples/03-loading-spinners/main.go",
		},
		item{
			title:       "📊 Progress Animations",
			description: "Different styles of animated progress bars",
			command:     "examples/04-progress-animations/main.go",
		},
		item{
			title:       "💻 Matrix Rain",
			description: "The classic Matrix digital rain effect",
			command:     "examples/05-matrix-rain/main.go",
		},
		item{
			title:       "🏀 Bouncing Ball",
			description: "Physics-based ball animation with trails",
			command:     "examples/06-bouncing-ball/main.go",
		},
		item{
			title:       "⭐ Starfield",
			description: "3D starfield simulation with depth perception",
			command:     "examples/07-starfield/main.go",
		},
		item{
			title:       "🎵 Audio Visualizer",
			description: "Simulated audio spectrum visualization",
			command:     "examples/08-audio-visualizer/main.go",
		},
	}

	// Add impressive animations section
	items = append(items,
		item{
			title:       "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
			description: "Advanced Animations - Impressive Visual Effects",
			command:     "",
		},
		item{
			title:       "🔥 Fire Effect",
			description: "Realistic fire simulation with heat propagation",
			command:     "examples/09-fire-effect/main.go",
		},
		item{
			title:       "💧 Fluid Simulation",
			description: "Water droplets with ripples and physics",
			command:     "examples/10-fluid-simulation/main.go",
		},
		item{
			title:       "🎲 3D Rotating Cube",
			description: "Real-time 3D wireframe cube with perspective",
			command:     "examples/11-rotating-cube/main.go",
		},
		item{
			title:       "🧬 Game of Life",
			description: "Conway's cellular automata with famous patterns",
			command:     "examples/12-game-of-life/main.go",
		},
		item{
			title:       "🌀 Mandelbrot Zoom",
			description: "Interactive fractal explorer with infinite zoom",
			command:     "examples/13-mandelbrot-zoom/main.go",
		},
	)

	// Add separator and Demoscene effects section
	items = append(items, 
		item{
			title:       "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
			description: "Demoscene Effects - Classic Computer Graphics",
			command:     "",
		},
		item{
			title:       "🌈 Plasma Effect",
			description: "Classic demoscene plasma with multiple color palettes",
			command:     "demoscene/01-plasma/main.go",
		},
		item{
			title:       "🕳️ Tunnel Effect",
			description: "Hypnotic tunnel with 4 different rendering modes",
			command:     "demoscene/02-tunnel/main.go",
		},
		item{
			title:       "🫧 Metaballs",
			description: "Organic metaball simulation with field visualization",
			command:     "demoscene/03-metaballs/main.go",
		},
		item{
			title:       "🌀 Rotozoom",
			description: "Rotating and zooming patterns with 5 different styles",
			command:     "demoscene/04-rotozoom/main.go",
		},
		item{
			title:       "📜 Scroller",
			description: "Demoscene text scroller with bitmap fonts and effects",
			command:     "demoscene/05-scroller/main.go",
		},
		item{
			title:       "🌆 Vaporwave",
			description: "Retro synthwave landscape with neon grid and floating shapes",
			command:     "demoscene/06-vaporwave/main.go",
		},
	)

	// Add separator and Bubbles components section
	items = append(items, 
		item{
			title:       "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
			description: "Bubbles Components - Interactive UI Elements",
			command:     "",
		},
		item{
			title:       "📝 Text Input",
			description: "Form inputs with validation and custom styling",
			command:     "bubbles/01-textinput/main.go",
		},
		item{
			title:       "📄 Textarea",
			description: "Multi-line text editor with preview mode",
			command:     "bubbles/02-textarea/main.go",
		},
		item{
			title:       "📊 Table",
			description: "Interactive data table with sorting and selection",
			command:     "bubbles/03-table/main.go",
		},
		item{
			title:       "📜 Viewport",
			description: "Scrollable content container for large documents",
			command:     "bubbles/04-viewport/main.go",
		},
		item{
			title:       "📁 File Picker",
			description: "File browser with filtering and navigation",
			command:     "bubbles/05-filepicker/main.go",
		},
	)

	l := list.New(items, list.NewDefaultDelegate(), 80, 20)
	l.Title = "🫧 Bubble Tea Showcase"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1)

	return model{list: l}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4) // Account for title and help text
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok && i.command != "" {
				m.choice = i.command
				return m, tea.Quit
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return ""
	}
	
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\n[↑↓] Navigate • [enter] Select • [q] Quit")
	
	return m.list.View() + help
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok && m.choice != "" {
		fmt.Printf("\033[2J\033[H")
		cmd := exec.Command("go", "run", m.choice)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running example: %v\n", err)
			os.Exit(1)
		}
	}
}