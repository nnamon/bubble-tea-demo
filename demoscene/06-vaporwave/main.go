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

// Floating shape for visual interest
type floatingShape struct {
	x, y     float64
	vx, vy   float64
	size     float64
	rotation float64
	rotSpeed float64
	shape    string
	color    lipgloss.Color
	age      float64
}

// Particle for atmospheric effects
type particle struct {
	x, y   float64
	vx, vy float64
	life   float64
	char   string
	color  lipgloss.Color
}

// Color mode configuration
type colorMode struct {
	name     string
	skyGrad  []string
	sunColor []string
	gridGrad []string
	fogColor string
}

type model struct {
	// Display properties
	width  int
	height int
	grid   [][]string  // Grid-based rendering for performance
	
	// Animation state
	time   float64
	speed  float64
	paused bool
	frame  int
	
	// Scene elements
	shapes    []floatingShape
	particles []particle
	
	// Configuration
	mode         int
	modes        []colorMode
	showShapes   bool
	showFog      bool
	gridIntensity float64
	sunPulse     bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	m := model{
		width:         80,
		height:        24,
		speed:         1.0,
		mode:          0,
		showShapes:    true,
		showFog:       true,
		gridIntensity: 1.2,
		sunPulse:      true,
		modes: []colorMode{
			{
				name:     "Classic Vaporwave",
				skyGrad:  []string{"#FF1493", "#FF69B4", "#DA70D6", "#9370DB", "#8A2BE2", "#4B0082"},
				sunColor: []string{"#FFD700", "#FFA500", "#FF8C00", "#FF4500"},
				gridGrad: []string{"#FF1493", "#DA70D6", "#9370DB", "#663399", "#4B0082"},
				fogColor: "#FF69B4",
			},
			{
				name:     "Miami Vice",
				skyGrad:  []string{"#FF6EC7", "#FF8A80", "#FFB74D", "#4FC3F7", "#29B6F6", "#0277BD"},
				sunColor: []string{"#FFD54F", "#FF8A65", "#FF7043", "#E91E63"},
				gridGrad: []string{"#FF6EC7", "#AB47BC", "#7E57C2", "#5E35B1"},
				fogColor: "#FF6EC7",
			},
			{
				name:     "Outrun",
				skyGrad:  []string{"#FF073A", "#FF6B35", "#F7931E", "#FFD23F", "#06FFA5", "#4ECDC4"},
				sunColor: []string{"#FFD23F", "#F7931E", "#FF6B35", "#FF073A"},
				gridGrad: []string{"#06FFA5", "#4ECDC4", "#45B7D1", "#96CEB4"},
				fogColor: "#06FFA5",
			},
			{
				name:     "Synthwave",
				skyGrad:  []string{"#FF0099", "#FF6600", "#FFFF00", "#00FFFF", "#9900FF", "#000033"},
				sunColor: []string{"#FFFF00", "#FF6600", "#FF0099", "#9900FF"},
				gridGrad: []string{"#00FFFF", "#00CCFF", "#0099FF", "#0066FF"},
				fogColor: "#FF0099",
			},
		},
	}
	m.initGrid()
	m.generateShapes()
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

// Generate floating shapes with aesthetic vaporwave elements
func (m *model) generateShapes() {
	// More vaporwave aesthetic shapes
	shapes := []string{"◆", "◇", "▲", "△", "●", "○", "■", "□", "★", "☆", "♦", "◈", "▼", "▽", "◉", "◎"}
	m.shapes = make([]floatingShape, 12) // More shapes for richer visuals
	
	for i := range m.shapes {
		// Create depth layers - some shapes in foreground, some in background
		layer := rand.Float64()
		var yRange, speed, sizeRange float64
		
		if layer < 0.3 { // Background layer
			yRange = float64(m.height) * 0.7
			speed = 0.2
			sizeRange = 0.5
		} else if layer < 0.7 { // Mid layer  
			yRange = float64(m.height) * 0.5
			speed = 0.4
			sizeRange = 1.0
		} else { // Foreground layer
			yRange = float64(m.height) * 0.3
			speed = 0.6
			sizeRange = 1.5
		}
		
		m.shapes[i] = floatingShape{
			x:        rand.Float64() * float64(m.width),
			y:        rand.Float64() * yRange,
			vx:       (rand.Float64() - 0.5) * speed,
			vy:       (rand.Float64() - 0.5) * speed * 0.3,
			size:     0.5 + rand.Float64()*sizeRange,
			rotation: rand.Float64() * math.Pi * 2,
			rotSpeed: (rand.Float64() - 0.5) * 0.05,
			shape:    shapes[rand.Intn(len(shapes))],
			color:    lipgloss.Color(m.modes[m.mode].skyGrad[rand.Intn(len(m.modes[m.mode].skyGrad))]),
		}
	}
}

// Emit atmospheric particles with multiple types
func (m *model) emitParticles() {
	if len(m.particles) < 30 && rand.Float64() < 0.4 {
		particleType := rand.Float64()
		
		if particleType < 0.6 { // Floating sparkles
			chars := []string{"·", "•", "◦", "∘", "˙", "⋅", "∙"}
			m.particles = append(m.particles, particle{
				x:     rand.Float64() * float64(m.width),
				y:     rand.Float64() * float64(m.height/2), // Throughout sky
				vx:    (rand.Float64() - 0.5) * 0.3,
				vy:    (rand.Float64() - 0.5) * 0.2,
				life:  1.0,
				char:  chars[rand.Intn(len(chars))],
				color: lipgloss.Color(m.modes[m.mode].fogColor),
			})
		} else if particleType < 0.8 { // Rising stars  
			chars := []string{"✦", "✧", "⋆", "✶", "✷", "✸"}
			m.particles = append(m.particles, particle{
				x:     rand.Float64() * float64(m.width),
				y:     float64(m.height),
				vx:    (rand.Float64() - 0.5) * 0.1,
				vy:    -rand.Float64()*0.3 - 0.1, // Rising upward
				life:  1.5,
				char:  chars[rand.Intn(len(chars))],
				color: lipgloss.Color(m.modes[m.mode].skyGrad[rand.Intn(len(m.modes[m.mode].skyGrad))]),
			})
		} else { // Drifting glows
			chars := []string{"◉", "◎", "○", "●", "◯"}
			m.particles = append(m.particles, particle{
				x:     rand.Float64() * float64(m.width),
				y:     rand.Float64() * float64(m.height*2/3),
				vx:    (rand.Float64() - 0.5) * 0.15,
				vy:    (rand.Float64() - 0.5) * 0.1,
				life:  2.0, // Longer lived
				char:  chars[rand.Intn(len(chars))],
				color: lipgloss.Color(m.modes[m.mode].sunColor[rand.Intn(len(m.modes[m.mode].sunColor))]),
			})
		}
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
		m.generateShapes()
		return m, nil

	case tickMsg:
		if !m.paused {
			m.frame++
			m.time += 0.05 * m.speed
			m.updateScene()
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
			m.generateShapes()
			m.particles = []particle{}
		case "1", "2", "3", "4":
			oldMode := m.mode
			m.mode = int(msg.String()[0] - '1')
			if m.mode != oldMode {
				m.generateShapes() // Regenerate with new colors
			}
		case "s":
			m.showShapes = !m.showShapes
		case "f":
			m.showFog = !m.showFog
		case "p":
			m.sunPulse = !m.sunPulse
		case "up":
			m.speed = common.Clamp(m.speed+0.2, 0.1, 3.0)
		case "down":
			m.speed = common.Clamp(m.speed-0.2, 0.1, 3.0)
		case "left":
			m.gridIntensity = common.Clamp(m.gridIntensity-0.2, 0.2, 2.0)
		case "right":
			m.gridIntensity = common.Clamp(m.gridIntensity+0.2, 0.2, 2.0)
		}
	}

	return m, nil
}

// Update all scene elements using proper physics
func (m *model) updateScene() {
	// Update floating shapes with physics
	for i := range m.shapes {
		s := &m.shapes[i]
		s.x += s.vx * m.speed
		s.y += s.vy * m.speed
		s.rotation += s.rotSpeed * m.speed
		s.age += 0.01
		
		// Gentle floating motion
		s.y += math.Sin(s.age*2 + float64(i)) * 0.1
		
		// Wrap around screen boundaries
		if s.x < 0 {
			s.x = float64(m.width)
		} else if s.x > float64(m.width) {
			s.x = 0
		}
		if s.y < 0 {
			s.y = float64(m.height / 3)
		} else if s.y > float64(m.height/3) {
			s.y = 0
		}
	}
	
	// Emit and update particles
	if m.showFog {
		m.emitParticles()
		
		alive := []particle{}
		for i := range m.particles {
			p := &m.particles[i]
			p.x += p.vx * m.speed
			p.y += p.vy * m.speed
			p.life -= 0.02
			
			// Keep alive particles within bounds
			if p.life > 0 && p.x >= 0 && p.x < float64(m.width) && 
			   p.y >= 0 && p.y < float64(m.height) {
				alive = append(alive, *p)
			}
		}
		m.particles = alive
	}
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(m.modes[m.mode].skyGrad[0])).
		Padding(0, 1)

	title := titleStyle.Render("🌆 " + m.modes[m.mode].name)

	// Status with more information
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.modes[m.mode].fogColor))
	status := statusStyle.Render(fmt.Sprintf(
		"Speed: %.1f | Grid: %.1f | Shapes: %s | Fog: %s | Pulse: %s | %s",
		m.speed, m.gridIntensity,
		map[bool]string{true: "ON", false: "OFF"}[m.showShapes],
		map[bool]string{true: "ON", false: "OFF"}[m.showFog],
		map[bool]string{true: "ON", false: "OFF"}[m.sunPulse],
		map[bool]string{true: "⏸ PAUSED", false: "▶ FLOWING"}[m.paused],
	))

	// Check minimum size requirements
	minWidth, minHeight := 60, 16
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

	// Render the complete scene using grid-based approach
	scene := m.renderCompleteScene()

	// Enhanced help
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render(
		"[1-4] modes • [↑↓] speed • [←→] grid • [s]hapes • [f]og • [p]ulse • [space] pause • [r]eset • [q]uit",
	)

	return lipgloss.JoinVertical(lipgloss.Left, title, status, "", scene, help)
}

// Grid-based rendering for optimal performance
func (m model) renderCompleteScene() string {
	// Clear the grid
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			m.grid[y][x] = " "
		}
	}
	
	// Render layers in order: sky -> sun -> grid -> shapes -> particles
	m.renderSky()
	m.renderSun()
	m.renderPerspectiveGrid()
	
	if m.showShapes {
		m.renderFloatingShapes()
	}
	
	if m.showFog {
		m.renderParticles()
	}
	
	// Convert grid to string with styling
	return m.gridToString()
}

// Render sky gradient with enhanced atmospheric effects
func (m *model) renderSky() {
	skyHeight := m.height / 3
	if skyHeight < 1 {
		skyHeight = 1
	}
	
	for y := 0; y < skyHeight; y++ {
		intensity := float64(y) / float64(skyHeight)
		
		for x := 0; x < m.width; x++ {
			// Enhanced noise patterns for atmospheric texture
			noise1 := math.Sin(float64(x)*0.08 + float64(y)*0.12 + m.time*0.6) * 0.15
			noise2 := math.Sin(float64(x)*0.15 + float64(y)*0.08 + m.time*0.4) * 0.08
			noise3 := math.Sin(float64(x)*0.05 + float64(y)*0.2 + m.time*1.2) * 0.05
			
			totalNoise := noise1 + noise2 + noise3
			adjustedIntensity := common.Clamp(intensity + totalNoise, 0, 1)
			
			// Create atmospheric layers
			char := m.getEnhancedGradientChar(adjustedIntensity, x, y)
			color := m.getSkyColor(adjustedIntensity)
			
			// Add subtle atmospheric effects
			if adjustedIntensity > 0.7 && math.Sin(float64(x)*0.2 + m.time*2) > 0.8 {
				// Wispy cloud effects
				char = "░"
			} else if y < skyHeight/4 && math.Sin(float64(x)*0.3 + m.time*0.8) > 0.9 {
				// High altitude shimmer
				char = "·"
			}
			
			m.grid[y][x] = m.styleChar(char, color)
		}
	}
}

// Render sun with enhanced dramatic effects
func (m *model) renderSun() {
	sunCenterX := m.width / 2
	sunCenterY := m.height / 4
	baseRadius := 5.0
	
	// Enhanced pulsing effect
	pulseIntensity := 1.0
	if m.sunPulse {
		pulseIntensity = 1.0 + math.Sin(m.time*2.5)*0.4 + math.Sin(m.time*4)*0.15
	}
	sunRadius := baseRadius * pulseIntensity
	
	for y := 0; y < m.height/2; y++ {
		for x := 0; x < m.width; x++ {
			dx := float64(x - sunCenterX)
			dy := float64(y - sunCenterY) * 1.6 // Character aspect ratio adjustment
			distance := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx)
			
			if distance < sunRadius-2 {
				// Sun core with varying intensity
				coreIntensity := 1.0 - (distance/(sunRadius-2))*0.3
				coreChar := "●"
				if distance < sunRadius-3 {
					coreChar = "◉"
				}
				m.grid[y][x] = m.styleChar(coreChar, m.getSunColor(coreIntensity))
			} else if distance < sunRadius {
				// Sun edge with animated glow
				edgeIntensity := 0.6 + math.Sin(m.time*3 + distance)*0.3
				glowChar := "◎"
				if math.Sin(m.time*2 + distance) > 0.5 {
					glowChar = "○"
				}
				m.grid[y][x] = m.styleChar(glowChar, m.getSunColor(edgeIntensity))
			} else if distance < sunRadius+3 {
				// Enhanced ray system
				rayIntensity := (sunRadius + 3 - distance) / 3
				rayPattern := int((angle + math.Pi) * 16 / (2 * math.Pi)) // More rays
				timeOffset := m.time*3 + float64(rayPattern)*0.5
				
				if math.Sin(timeOffset) > 0.3 {
					rayChar := "─"
					if rayPattern%2 == 0 {
						rayChar = "━"
					}
					if rayPattern%4 == 0 {
						rayChar = "═"
					}
					intensity := rayIntensity * (0.5 + math.Sin(timeOffset)*0.5)
					m.grid[y][x] = m.styleChar(rayChar, m.getSunColor(intensity))
				}
			} else if distance < sunRadius+6 {
				// Extended glow with scan lines for retro effect
				glowIntensity := (sunRadius + 6 - distance) / 6 * 0.3
				if y%2 == int(m.time*10)%2 { // Moving scan lines
					m.grid[y][x] = m.styleChar("▒", m.getSunColor(glowIntensity))
				}
			}
		}
	}
}

// Render perspective grid with enhanced dramatic effects
func (m *model) renderPerspectiveGrid() {
	gridStart := m.height / 3
	
	for y := gridStart; y < m.height; y++ {
		depth := float64(y - gridStart + 1)
		if depth <= 0 {
			depth = 1
		}
		
		// Enhanced perspective with dramatic scaling
		scale := 25.0 / (depth * 1.2)
		offset := m.time * m.speed * scale * 1.5
		
		// Add horizontal scan line effect
		scanLineIntensity := math.Sin(float64(y)*0.5 + m.time*8) * 0.1
		
		for x := 0; x < m.width; x++ {
			gridX := (float64(x) - float64(m.width)/2) / scale
			gridZ := depth + offset
			
			// Enhanced multi-layered wave effects
			wave1 := math.Sin(gridX*0.4 + gridZ*0.25 + m.time*1.2) * 3.0 / depth
			wave2 := math.Sin(gridX*0.2 + gridZ*0.15 + m.time*0.8) * 1.5 / depth
			wave3 := math.Sin(gridX*0.1 + gridZ*0.05 + m.time*2.0) * 0.8 / depth
			waveOffset := wave1 + wave2 + wave3
			
			// Enhanced grid line detection
			gridSpacing := 1.8
			lineThickness := (0.12/scale) * m.gridIntensity
			
			isGridLineX := math.Abs(math.Mod(gridX+0.5, gridSpacing)-gridSpacing/2) < lineThickness
			isGridLineZ := math.Abs(math.Mod(gridZ+0.5, gridSpacing)-gridSpacing/2) < lineThickness
			
			if isGridLineX || isGridLineZ {
				// Distance-based intensity with enhanced falloff
				intensity := (1.0 / (depth*0.08 + 1)) * m.gridIntensity
				
				// Major grid line emphasis (every 4th line)
				majorLineX := math.Abs(math.Mod(gridX+0.5, gridSpacing*4)-gridSpacing*2) < lineThickness*2
				majorLineZ := math.Abs(math.Mod(gridZ+0.5, gridSpacing*4)-gridSpacing*2) < lineThickness*2
				
				if majorLineX || majorLineZ {
					intensity *= 2.0
				}
				
				// Grid intersection with extra emphasis
				if isGridLineX && isGridLineZ {
					intensity *= 1.6
					if majorLineX && majorLineZ {
						intensity *= 1.4 // Super bright intersections
					}
				}
				
				// Enhanced glow and wave effects
				glowIntensity := intensity + math.Sin(waveOffset+m.time*2.5)*0.5 + scanLineIntensity
				glowIntensity = common.Clamp(glowIntensity, 0, 1.5)
				
				// More varied characters based on intensity and position
				char := m.getEnhancedGridChar(isGridLineX, isGridLineZ, majorLineX, majorLineZ, glowIntensity)
				m.grid[y][x] = m.styleChar(char, m.getGridColor(glowIntensity))
			} else if math.Sin(float64(y)*0.3 + m.time*5) > 0.95 {
				// Occasional scan line artifacts for retro CRT effect
				m.grid[y][x] = m.styleChar("▁", m.getGridColor(0.2))
			}
		}
	}
}

// Render floating shapes with physics
func (m *model) renderFloatingShapes() {
	for _, shape := range m.shapes {
		x, y := int(shape.x), int(shape.y)
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			// Apply rotation effect through character variation
			rotatedShape := shape.shape
			if int(shape.rotation*4) % 2 == 1 {
				// Simple rotation simulation
				switch shape.shape {
				case "◆": rotatedShape = "◇"
				case "▲": rotatedShape = "▼"
				case "■": rotatedShape = "▪"
				}
			}
			
			// Age-based transparency effect
			alpha := 0.7 + math.Sin(shape.age*2)*0.3
			color := shape.color
			if alpha < 0.5 {
				color = lipgloss.Color("#666666") // Fade effect
			}
			
			m.grid[y][x] = m.styleChar(rotatedShape, color)
		}
	}
}

// Render atmospheric particles
func (m *model) renderParticles() {
	for _, particle := range m.particles {
		x, y := int(particle.x), int(particle.y)
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			// Life-based alpha blending
			if particle.life > 0.5 || int(m.frame*3) % 2 == 0 {
				m.grid[y][x] = m.styleChar(particle.char, particle.color)
			}
		}
	}
}

// Convert grid to styled string
func (m model) gridToString() string {
	lines := make([]string, m.height)
	for y := 0; y < m.height; y++ {
		lines[y] = strings.Join(m.grid[y], "")
	}
	return strings.Join(lines, "\n")
}

// Helper functions for color and character selection
func (m model) getSkyColor(intensity float64) lipgloss.Color {
	colors := m.modes[m.mode].skyGrad
	index := common.Clamp(intensity * float64(len(colors)-1), 0, float64(len(colors)-1))
	return lipgloss.Color(colors[int(index)])
}

func (m model) getSunColor(intensity float64) lipgloss.Color {
	colors := m.modes[m.mode].sunColor
	index := common.Clamp(intensity * float64(len(colors)-1), 0, float64(len(colors)-1))
	return lipgloss.Color(colors[int(index)])
}

func (m model) getGridColor(intensity float64) lipgloss.Color {
	colors := m.modes[m.mode].gridGrad
	index := common.Clamp(intensity * float64(len(colors)-1), 0, float64(len(colors)-1))
	return lipgloss.Color(colors[int(index)])
}

func (m model) getGradientChar(intensity float64) string {
	chars := []string{"▓", "▒", "░", " "}
	index := int(intensity * float64(len(chars)-1))
	if index >= len(chars) {
		index = len(chars) - 1
	}
	return chars[index]
}

func (m model) getEnhancedGradientChar(intensity float64, x, y int) string {
	// More varied atmospheric characters
	if intensity < 0.1 {
		return "█"
	} else if intensity < 0.25 {
		return "▓"
	} else if intensity < 0.45 {
		return "▒"
	} else if intensity < 0.65 {
		return "░"
	} else if intensity < 0.85 {
		// Add some variation for atmospheric effect
		if (x+y)%3 == 0 {
			return "·"
		}
		return " "
	} else {
		return " "
	}
}

func (m model) getGridChar(isLineX, isLineZ bool, intensity float64) string {
	if isLineX && isLineZ {
		return "+"
	} else if isLineX {
		return "|"
	} else if isLineZ {
		return "-"
	}
	
	// Intensity-based characters
	chars := []string{"░", "▒", "▓", "█"}
	index := int(intensity * float64(len(chars)-1))
	if index >= len(chars) {
		index = len(chars) - 1
	}
	return chars[index]
}

func (m model) getEnhancedGridChar(isLineX, isLineZ, majorLineX, majorLineZ bool, intensity float64) string {
	if isLineX && isLineZ {
		if majorLineX && majorLineZ {
			return "╬" // Major intersection
		}
		if majorLineX || majorLineZ {
			return "┼" // Semi-major intersection
		}
		return "+" // Regular intersection
	} else if isLineX {
		if majorLineX {
			return "┃" // Major vertical
		}
		return "|" // Regular vertical
	} else if isLineZ {
		if majorLineZ {
			return "━" // Major horizontal  
		}
		return "─" // Regular horizontal
	}
	
	// Enhanced intensity-based characters
	if intensity > 1.2 {
		return "█"
	} else if intensity > 0.9 {
		return "▓"
	} else if intensity > 0.6 {
		return "▒"
	} else if intensity > 0.3 {
		return "░"
	} else {
		return "·"
	}
}

func (m model) styleChar(char string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Render(char)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}