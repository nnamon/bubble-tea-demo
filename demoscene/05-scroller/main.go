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
	scrollPos  float64
	waveHeight float64
	speed      float64
	message    string
	font       int
	colorMode  int
	paused     bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		width:      80,
		height:     24,
		waveHeight: 3.0,
		speed:      1.0,
		message:    "DEMOSCENE GREETINGS! * BUBBLE TEA SHOWCASE * TERMINAL GRAPHICS RULE * ",
		font:       0,
		colorMode:  0,
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
			m.time += 0.1
			m.scrollPos += 0.5 * m.speed
			if m.scrollPos > float64(len(m.message)*8) {
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
			m.scrollPos = -float64(m.width)
		case "1":
			m.font = 0 // Block
		case "2":
			m.font = 1 // Outline
		case "3":
			m.font = 2 // Dotted
		case "4":
			m.colorMode = 0 // Rainbow wave
		case "5":
			m.colorMode = 1 // Fire
		case "6":
			m.colorMode = 2 // Matrix
		case "7":
			m.colorMode = 3 // Plasma
		case "up":
			m.speed = math.Min(m.speed+0.2, 3.0)
		case "down":
			m.speed = math.Max(m.speed-0.2, 0.1)
		case "left":
			m.waveHeight = math.Max(m.waveHeight-0.5, 0.5)
		case "right":
			m.waveHeight = math.Min(m.waveHeight+0.5, 8.0)
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

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.Green)
	fonts := []string{"Block", "Outline", "Dotted"}
	colors := []string{"Rainbow", "Fire", "Matrix", "Plasma"}
	status := statusStyle.Render(fmt.Sprintf(
		"Font: %s | Color: %s | Speed: %.1f | Wave: %.1f | %s",
		fonts[m.font], colors[m.colorMode], m.speed, m.waveHeight,
		map[bool]string{true: "⏸ Paused", false: "📜 Scrolling"}[m.paused],
	))

	// Render scroller
	lines := m.renderScroller()

	// Help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1-3] fonts • [4-7] colors • [↑↓] speed • [←→] wave • [space] pause • [r]eset • [q]uit",
	)

	return fmt.Sprintf("%s\n%s\n\n%s\n%s",
		title, status, strings.Join(lines, "\n"), help)
}

func (m model) renderScroller() []string {
	lines := make([]string, m.height)
	
	// Initialize all lines
	for i := range lines {
		lines[i] = strings.Repeat(" ", m.width)
	}

	// Calculate text position
	textY := m.height / 2
	textStartX := int(-m.scrollPos)

	// Render each character of the message
	for charIndex, char := range m.message {
		charX := textStartX + charIndex*8

		// Only render if character is potentially visible
		if charX > -8 && charX < m.width+8 {
			m.renderCharacter(lines, char, charX, textY, charIndex)
		}
	}

	return lines
}

func (m model) renderCharacter(lines []string, char rune, startX, centerY, charIndex int) {
	// Get the bitmap for this character
	bitmap := m.getCharacterBitmap(char)
	
	for y := 0; y < len(bitmap); y++ {
		for x := 0; x < len(bitmap[y]); x++ {
			if bitmap[y][x] == '1' {
				screenX := startX + x
				screenY := centerY - len(bitmap)/2 + y

				// Apply sine wave effect
				waveOffset := math.Sin(float64(screenX)*0.1 + m.time*2) * m.waveHeight
				screenY += int(waveOffset)

				// Check bounds
				if screenX >= 0 && screenX < m.width && screenY >= 0 && screenY < m.height {
					char, color := m.getStyledCharacter(screenX, screenY, charIndex)
					
					// Replace character in line
					line := []rune(lines[screenY])
					if screenX < len(line) {
						style := lipgloss.NewStyle().Foreground(color)
						styledChar := style.Render(string(char))
						
						// Replace the character at this position
						lineStr := string(line[:screenX]) + styledChar + string(line[screenX+1:])
						lines[screenY] = lineStr
					}
				}
			}
		}
	}
}

func (m model) getCharacterBitmap(char rune) []string {
	// Simple 7x5 bitmap font
	switch char {
	case 'A':
		return []string{
			"01110",
			"10001",
			"11111",
			"10001",
			"10001",
		}
	case 'B':
		return []string{
			"11110",
			"10001",
			"11110",
			"10001",
			"11110",
		}
	case 'C':
		return []string{
			"01111",
			"10000",
			"10000",
			"10000",
			"01111",
		}
	case 'D':
		return []string{
			"11110",
			"10001",
			"10001",
			"10001",
			"11110",
		}
	case 'E':
		return []string{
			"11111",
			"10000",
			"11110",
			"10000",
			"11111",
		}
	case 'F':
		return []string{
			"11111",
			"10000",
			"11110",
			"10000",
			"10000",
		}
	case 'G':
		return []string{
			"01111",
			"10000",
			"10011",
			"10001",
			"01111",
		}
	case 'H':
		return []string{
			"10001",
			"10001",
			"11111",
			"10001",
			"10001",
		}
	case 'I':
		return []string{
			"11111",
			"00100",
			"00100",
			"00100",
			"11111",
		}
	case 'L':
		return []string{
			"10000",
			"10000",
			"10000",
			"10000",
			"11111",
		}
	case 'M':
		return []string{
			"10001",
			"11011",
			"10101",
			"10001",
			"10001",
		}
	case 'N':
		return []string{
			"10001",
			"11001",
			"10101",
			"10011",
			"10001",
		}
	case 'O':
		return []string{
			"01110",
			"10001",
			"10001",
			"10001",
			"01110",
		}
	case 'P':
		return []string{
			"11110",
			"10001",
			"11110",
			"10000",
			"10000",
		}
	case 'R':
		return []string{
			"11110",
			"10001",
			"11110",
			"10010",
			"10001",
		}
	case 'S':
		return []string{
			"01111",
			"10000",
			"01110",
			"00001",
			"11110",
		}
	case 'T':
		return []string{
			"11111",
			"00100",
			"00100",
			"00100",
			"00100",
		}
	case 'U':
		return []string{
			"10001",
			"10001",
			"10001",
			"10001",
			"01110",
		}
	case 'W':
		return []string{
			"10001",
			"10001",
			"10101",
			"11011",
			"10001",
		}
	case 'Y':
		return []string{
			"10001",
			"10001",
			"01010",
			"00100",
			"00100",
		}
	case ' ':
		return []string{
			"00000",
			"00000",
			"00000",
			"00000",
			"00000",
		}
	case '*':
		return []string{
			"00100",
			"10101",
			"01110",
			"10101",
			"00100",
		}
	case '!':
		return []string{
			"00100",
			"00100",
			"00100",
			"00000",
			"00100",
		}
	default:
		// Default to a block for unknown characters
		return []string{
			"11111",
			"10001",
			"10001",
			"10001",
			"11111",
		}
	}
}

func (m model) getStyledCharacter(x, y, charIndex int) (rune, lipgloss.Color) {
	var char rune
	switch m.font {
	case 0: // Block
		char = '█'
	case 1: // Outline
		char = '▓'
	case 2: // Dotted
		char = '●'
	default:
		char = '█'
	}

	var color lipgloss.Color
	switch m.colorMode {
	case 0: // Rainbow wave
		hue := float64(x+charIndex*20) * 0.05 + m.time
		color = m.getRainbowColor(hue)
	case 1: // Fire
		intensity := math.Sin(float64(x)*0.1 + m.time*2) * 0.5 + 0.5
		color = m.getFireColor(intensity)
	case 2: // Matrix
		intensity := math.Sin(float64(y)*0.2 + m.time*3) * 0.5 + 0.5
		color = m.getMatrixColor(intensity)
	case 3: // Plasma
		plasma := math.Sin(float64(x)*0.1) + math.Sin(float64(y)*0.15) + math.Sin(m.time*2)
		color = m.getPlasmaColor((plasma + 3) / 6)
	default:
		color = lipgloss.Color("#FFFFFF")
	}

	return char, color
}

func (m model) getRainbowColor(hue float64) lipgloss.Color {
	phase := math.Mod(hue, 6)
	if phase < 1 {
		return lipgloss.Color("#FF0000")
	} else if phase < 2 {
		return lipgloss.Color("#FF8000")
	} else if phase < 3 {
		return lipgloss.Color("#FFFF00")
	} else if phase < 4 {
		return lipgloss.Color("#00FF00")
	} else if phase < 5 {
		return lipgloss.Color("#0080FF")
	} else {
		return lipgloss.Color("#8000FF")
	}
}

func (m model) getFireColor(intensity float64) lipgloss.Color {
	if intensity < 0.3 {
		return lipgloss.Color("#FF0000")
	} else if intensity < 0.6 {
		return lipgloss.Color("#FF4000")
	} else if intensity < 0.8 {
		return lipgloss.Color("#FF8000")
	} else {
		return lipgloss.Color("#FFFF00")
	}
}

func (m model) getMatrixColor(intensity float64) lipgloss.Color {
	if intensity < 0.3 {
		return lipgloss.Color("#004000")
	} else if intensity < 0.6 {
		return lipgloss.Color("#008000")
	} else if intensity < 0.8 {
		return lipgloss.Color("#00C000")
	} else {
		return lipgloss.Color("#00FF00")
	}
}

func (m model) getPlasmaColor(value float64) lipgloss.Color {
	if value < 0.2 {
		return lipgloss.Color("#FF0080")
	} else if value < 0.4 {
		return lipgloss.Color("#8000FF")
	} else if value < 0.6 {
		return lipgloss.Color("#0080FF")
	} else if value < 0.8 {
		return lipgloss.Color("#00FF80")
	} else {
		return lipgloss.Color("#80FF00")
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}