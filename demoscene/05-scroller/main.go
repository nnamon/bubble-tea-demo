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

// Character bitmap definition
type charBitmap []string

// Color mode configuration  
type colorMode struct {
	name   string
	colors []string
}

type model struct {
	// Display properties
	width  int
	height int
	grid   [][]string // Grid-based rendering for performance
	
	// Animation state
	time       float64
	frame      int
	scrollPos  float64
	waveHeight float64
	speed      float64
	paused     bool
	
	// Content and configuration
	message    string
	font       int
	colorMode  int
	modes      []colorMode
	bitmaps    map[rune]charBitmap
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	m := model{
		width:      80,
		height:     24,
		waveHeight: 3.0,
		speed:      1.0,
		message:    "DEMOSCENE GREETINGS! * BUBBLE TEA SHOWCASE * TERMINAL GRAPHICS RULE * ",
		font:       0,
		colorMode:  0,
		modes: []colorMode{
			{name: "Rainbow Wave", colors: []string{"#FF0000", "#FF8000", "#FFFF00", "#00FF00", "#0080FF", "#8000FF"}},
			{name: "Fire", colors: []string{"#FF0000", "#FF4000", "#FF8000", "#FFFF00"}},
			{name: "Matrix", colors: []string{"#004000", "#008000", "#00C000", "#00FF00"}},
			{name: "Plasma", colors: []string{"#FF0080", "#8000FF", "#0080FF", "#00FF80", "#80FF00"}},
		},
		bitmaps: initBitmaps(),
	}
	m.initGrid()
	return m
}

// Initialize grid for efficient rendering
func (m *model) initGrid() {
	m.grid = make([][]string, m.height)
	for i := range m.grid {
		m.grid[i] = make([]string, m.width)
		for j := range m.grid[i] {
			m.grid[i][j] = " "
		}
	}
}

// Pre-calculate all character bitmaps for performance
func initBitmaps() map[rune]charBitmap {
	return map[rune]charBitmap{
		'A': {"01110", "10001", "11111", "10001", "10001"},
		'B': {"11110", "10001", "11110", "10001", "11110"},
		'C': {"01111", "10000", "10000", "10000", "01111"},
		'D': {"11110", "10001", "10001", "10001", "11110"},
		'E': {"11111", "10000", "11110", "10000", "11111"},
		'F': {"11111", "10000", "11110", "10000", "10000"},
		'G': {"01111", "10000", "10011", "10001", "01111"},
		'H': {"10001", "10001", "11111", "10001", "10001"},
		'I': {"11111", "00100", "00100", "00100", "11111"},
		'J': {"11111", "00010", "00010", "10010", "01100"},
		'K': {"10010", "10100", "11000", "10100", "10010"},
		'L': {"10000", "10000", "10000", "10000", "11111"},
		'M': {"10001", "11011", "10101", "10001", "10001"},
		'N': {"10001", "11001", "10101", "10011", "10001"},
		'O': {"01110", "10001", "10001", "10001", "01110"},
		'P': {"11110", "10001", "11110", "10000", "10000"},
		'Q': {"01110", "10001", "10101", "10010", "01101"},
		'R': {"11110", "10001", "11110", "10010", "10001"},
		'S': {"01111", "10000", "01110", "00001", "11110"},
		'T': {"11111", "00100", "00100", "00100", "00100"},
		'U': {"10001", "10001", "10001", "10001", "01110"},
		'V': {"10001", "10001", "10001", "01010", "00100"},
		'W': {"10001", "10001", "10101", "11011", "10001"},
		'X': {"10001", "01010", "00100", "01010", "10001"},
		'Y': {"10001", "10001", "01010", "00100", "00100"},
		'Z': {"11111", "00010", "00100", "01000", "11111"},
		' ': {"00000", "00000", "00000", "00000", "00000"},
		'*': {"00100", "10101", "01110", "10101", "00100"},
		'!': {"00100", "00100", "00100", "00000", "00100"},
		'.': {"00000", "00000", "00000", "00000", "00100"},
		',': {"00000", "00000", "00000", "00100", "01000"},
		'?': {"01110", "10001", "00110", "00000", "00100"},
		'-': {"00000", "00000", "11111", "00000", "00000"},
		'+': {"00000", "00100", "01110", "00100", "00000"},
		'0': {"01110", "10001", "10001", "10001", "01110"},
		'1': {"00100", "01100", "00100", "00100", "01110"},
		'2': {"01110", "10001", "00110", "01000", "11111"},
		'3': {"01110", "10001", "00110", "10001", "01110"},
		'4': {"10001", "10001", "11111", "00001", "00001"},
		'5': {"11111", "10000", "11110", "00001", "11110"},
		'6': {"01110", "10000", "11110", "10001", "01110"},
		'7': {"11111", "00001", "00010", "00100", "01000"},
		'8': {"01110", "10001", "01110", "10001", "01110"},
		'9': {"01110", "10001", "01111", "00001", "01110"},
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
		m.initGrid()
		return m, nil

	case tickMsg:
		if !m.paused {
			m.frame++
			m.time += 0.05 * m.speed
			
			// Update scroll position with smooth movement
			m.scrollPos += 0.8 * m.speed
			
			// Reset when message completely scrolls off screen
			messageWidth := float64(len(m.message) * 6) // 5 chars + 1 space per character
			if m.scrollPos > messageWidth + float64(m.width) {
				m.scrollPos = -float64(m.width)
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
			m.time = 0
			m.frame = 0
			m.scrollPos = -float64(m.width)
		case "1", "2", "3":
			newFont := int(msg.String()[0] - '1')
			if newFont >= 0 && newFont < 3 {
				m.font = newFont
			}
		case "4", "5", "6", "7":
			newMode := int(msg.String()[0] - '4')
			if newMode < len(m.modes) {
				m.colorMode = newMode
			}
		case "up":
			m.speed = common.Clamp(m.speed+0.2, 0.1, 4.0)
		case "down":
			m.speed = common.Clamp(m.speed-0.2, 0.1, 4.0)
		case "left":
			m.waveHeight = common.Clamp(m.waveHeight-0.5, 0.0, 8.0)
		case "right":
			m.waveHeight = common.Clamp(m.waveHeight+0.5, 0.0, 8.0)
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#00FF80")).
		Padding(0, 1)

	title := titleStyle.Render("📜 Demoscene Scroller")

	// Status with enhanced information
	statusStyle := lipgloss.NewStyle().Foreground(common.Green)
	fonts := []string{"Block", "Outline", "Dotted"}
	status := statusStyle.Render(fmt.Sprintf(
		"Font: %s | Color: %s | Speed: %.1f | Wave: %.1f | %s",
		fonts[m.font], m.modes[m.colorMode].name, m.speed, m.waveHeight,
		map[bool]string{true: "⏸ PAUSED", false: "📜 SCROLLING"}[m.paused],
	))

	// Check minimum size requirements
	minWidth, minHeight := 60, 12
	if m.width < minWidth || m.height < minHeight {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
		
		sizeError := errorStyle.Render(fmt.Sprintf(
			"Terminal too small!\nMinimum size: %dx%d\nCurrent size: %dx%d\n\nPlease resize your terminal window.",
			minWidth, minHeight+4, m.width, m.height+4,
		))
		
		helpStyle := lipgloss.NewStyle().Faint(true)
		help := helpStyle.Render("[q]uit")

		return lipgloss.JoinVertical(lipgloss.Left, title, status, "", sizeError, help)
	}

	// Render the complete scroller using grid-based approach
	scene := m.renderCompleteScroller()

	// Enhanced help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1-3] fonts • [4-7] colors • [↑↓] speed • [←→] wave • [space] pause • [r]eset • [q]uit",
	)

	return lipgloss.JoinVertical(lipgloss.Left, title, status, "", scene, help)
}

// Grid-based rendering for optimal performance
func (m model) renderCompleteScroller() string {
	// Clear the grid
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			m.grid[y][x] = " "
		}
	}
	
	// Render scrolling text to grid
	m.renderScrollingText()
	
	// Convert grid to string with styling
	return m.gridToString()
}

// Render scrolling text using efficient grid-based approach
func (m *model) renderScrollingText() {
	centerY := m.height / 2
	textStartX := int(-m.scrollPos)
	
	// Render each character of the message
	for charIndex, char := range m.message {
		charX := textStartX + charIndex*6 // 5 char width + 1 space
		
		// Only render if character is potentially visible
		if charX > -6 && charX < m.width+6 {
			m.renderCharacterToGrid(char, charX, centerY, charIndex)
		}
	}
}

// Render a single character to the grid using bitmap font
func (m *model) renderCharacterToGrid(char rune, startX, centerY, charIndex int) {
	// Get bitmap, fallback to default if not found
	bitmap, exists := m.bitmaps[char]
	if !exists {
		// Fallback to a simple block pattern
		bitmap = charBitmap{"11111", "10001", "10001", "10001", "11111"}
	}
	
	bitmapHeight := len(bitmap)
	startY := centerY - bitmapHeight/2
	
	for y := 0; y < bitmapHeight; y++ {
		for x := 0; x < len(bitmap[y]); x++ {
			if bitmap[y][x] == '1' {
				screenX := startX + x
				screenY := startY + y
				
				// Apply sine wave effect
				waveOffset := math.Sin(float64(screenX)*0.08 + m.time*2.5) * m.waveHeight
				finalY := screenY + int(waveOffset)
				
				// Check bounds and render
				if screenX >= 0 && screenX < m.width && finalY >= 0 && finalY < m.height {
					char, color := m.getStyledCharacter(screenX, finalY, charIndex)
					m.grid[finalY][screenX] = m.styleChar(string(char), color)
				}
			}
		}
	}
}

// Get styled character and color based on current configuration
func (m model) getStyledCharacter(x, y, charIndex int) (rune, lipgloss.Color) {
	// Character selection based on font
	var char rune
	switch m.font {
	case 0: char = '█' // Block
	case 1: char = '▓' // Outline
	case 2: char = '●' // Dotted
	default: char = '█'
	}
	
	// Color calculation based on mode
	var colorIntensity float64
	switch m.colorMode {
	case 0: // Rainbow Wave
		colorIntensity = math.Mod(float64(x+charIndex*20)*0.05 + m.time, 1.0)
	case 1: // Fire
		colorIntensity = (math.Sin(float64(x)*0.1 + m.time*2) + 1) / 2
	case 2: // Matrix
		colorIntensity = (math.Sin(float64(y)*0.2 + m.time*3) + 1) / 2
	case 3: // Plasma
		plasma := math.Sin(float64(x)*0.1) + math.Sin(float64(y)*0.15) + math.Sin(m.time*2)
		colorIntensity = (plasma + 3) / 6
	default:
		colorIntensity = 1.0
	}
	
	color := m.getColorFromIntensity(colorIntensity)
	return char, color
}

// Get color from intensity using current color mode
func (m model) getColorFromIntensity(intensity float64) lipgloss.Color {
	colors := m.modes[m.colorMode].colors
	index := common.Clamp(intensity * float64(len(colors)-1), 0, float64(len(colors)-1))
	return lipgloss.Color(colors[int(index)])
}

// Convert grid to styled string
func (m model) gridToString() string {
	lines := make([]string, m.height)
	for y := 0; y < m.height; y++ {
		lines[y] = strings.Join(m.grid[y], "")
	}
	return strings.Join(lines, "\n")
}

// Helper function to style characters
func (m model) styleChar(char string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Render(char)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}