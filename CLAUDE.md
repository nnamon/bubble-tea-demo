# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running Examples
```bash
# Run the main showcase (interactive menu)
make run
# OR
go run showcase/main.go

# Run a specific example directly
go run examples/01-wave-animation/main.go
go run demoscene/02-tunnel/main.go
go run bubbles/03-table/main.go

# Build all examples to bin/ directory
make build

# Clean built binaries
make clean
```

### Requirements
- Go 1.21+ (currently using 1.23.0)
- Terminal with Unicode support and 256-color capability for optimal experience

## Architecture Overview

This is a **Bubble Tea showcase repository** demonstrating TUI animations, visual effects, and interactive components. The codebase follows The Elm Architecture (TEA) pattern with strict separation of concerns.

### Directory Structure

- **`examples/`** - Basic animations and visual effects (13 demos)
- **`demoscene/`** - Advanced demoscene-style effects (6 demos) 
- **`bubbles/`** - Interactive UI components using the Bubbles library (5 demos)
- **`showcase/`** - Main interactive launcher that runs other demos
- **`common/`** - Shared utilities for animations and styling

### Core Design Patterns

**TEA Implementation**
- Each demo implements the standard TEA interface: `Init()`, `Update()`, `View()`
- Animation state managed through `tickMsg` with 30fps target (`time.Second/30`)
- All side effects handled via commands, never in Update/View functions
- Immutable state updates - models are copied, never mutated

**Animation Architecture**
```go
type tickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

// In Update():
case tickMsg:
    m.animationTime += 0.05
    return m, tick() // Continue animation loop
```

**Shared Utilities (`common/` package)**
- `colors.go` - Predefined color palette and gradients (GradientBlue, GradientFire)
- `utils.go` - Mathematical helpers for animations:
  - `Lerp()`, `Clamp()`, `Map()` for value interpolation
  - `GetWaveChar()` for Unicode wave visualization
  - `GenerateGradient()` for color transitions

### Demo Categories

**Examples (Basic Effects)**
Progressive complexity from simple waves to advanced fractals. Each demonstrates specific animation techniques.

**Demoscene (Advanced Visual Effects)**  
Classic computer graphics effects with multiple rendering modes and complex mathematical calculations.

**Bubbles (UI Components)**
Interactive form elements and data display components using the Bubbles component library.

### Key Implementation Details

**Window Management**
All demos handle `tea.WindowSizeMsg` for responsive layouts. Most use `tea.WithAltScreen()` for full-screen rendering.

**Performance Considerations**  
- Grid-based rendering for pixel-like effects
- Pre-calculated mathematical values where possible  
- Efficient string building for complex visuals
- Frame skipping logic in computationally heavy demos

**Interactive Controls**
Standard keybindings across demos:
- `q` or `ctrl+c` - Quit
- `space` - Toggle pause/effects
- Arrow keys - Parameter adjustment
- `r` - Reset animation
- `h` - Toggle help display

### Module Import Pattern
Examples import the common utilities:
```go
import "github.com/yourusername/bubbletea-showcase/common"
```

Note: The module name in go.mod uses a placeholder GitHub URL and should be updated for actual deployment.

### Showcase Launcher Pattern
The main showcase (`showcase/main.go`) uses a Bubbles list component to present organized categories of demos. It executes selected demos using `exec.Command("go", "run", selectedPath)` with proper terminal handoff.

## Bubble Tea Framework Deep Knowledge

### The Elm Architecture (TEA) in Practice

Bubble Tea implements a strict unidirectional data flow pattern borrowed from Elm:

```
User Input → Message → Update(msg, model) → New Model → View(model) → Rendered UI
```

**Core Interface Every Application Must Implement:**
```go
type Model interface {
    Init() tea.Cmd                      // Initial setup and commands
    Update(tea.Msg) (tea.Model, tea.Cmd) // Handle events, return new state
    View() string                       // Pure render function
}
```

**Critical Principles:**
- **State Immutability**: Never modify model fields directly. Always return new model instances
- **Pure Functions**: `Update()` and `View()` must be pure - no side effects, network calls, or I/O
- **Commands for Side Effects**: All I/O, timers, API calls handled through the command system
- **Single Source of Truth**: All application state lives in the model

### Message System Architecture

**Built-in System Messages:**
- `tea.KeyMsg` - Keyboard input with `.String()` method for key identification
- `tea.MouseMsg` - Mouse events (click, wheel, movement) 
- `tea.WindowSizeMsg` - Terminal resize events with `.Width` and `.Height`
- `tea.QuitMsg` - Application quit signal

**Custom Message Pattern:**
```go
type tickMsg time.Time
type dataLoadedMsg struct{ data []string }
type errorMsg struct{ err error }

// In Update():
switch msg := msg.(type) {
case tickMsg:
    // Handle animation frame
case dataLoadedMsg:
    m.data = msg.data
case errorMsg:
    m.error = msg.err
}
```

**Message Type Assertion Pattern:**
Always use type switches, not type assertions, for message handling to avoid panics.

### Command System Mastery

**Command Types:**
```go
// Simple command returning immediate message
func doSomething() tea.Msg {
    // Perform work
    return resultMsg{data: "done"}
}

// Command function (deferred execution)
func fetchData(url string) tea.Cmd {
    return func() tea.Msg {
        // HTTP request, file I/O, etc.
        result, err := http.Get(url)
        if err != nil {
            return errorMsg{err}
        }
        return dataMsg{result}
    }
}

// Timer/animation command
func tick() tea.Cmd {
    return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
```

**Command Composition:**
```go
// Run multiple commands in parallel
return m, tea.Batch(cmd1, cmd2, cmd3)

// Sequential command (rare, usually in Init)
return m, tea.Sequence(setupCmd, startCmd)

// No command needed
return m, nil
```

### Animation Patterns and Performance

**Frame Rate Management:**
- Standard: `time.Second/30` (30 FPS) for smooth animation
- High performance: `time.Second/60` (60 FPS) for demanding visuals
- Low performance: `time.Second/10` (10 FPS) for simple updates
- Variable rate: Adjust based on animation complexity

**Animation State Pattern:**
```go
type model struct {
    // Animation timing
    frame     int
    time      float64
    paused    bool
    
    // Visual state
    particles []particle
    effects   []effect
    
    // UI state
    width     int
    height    int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        if !m.paused {
            m.frame++
            m.time += 0.033 // ~30fps time increment
            m.updateAnimations()
        }
        return m, tick()
    }
    return m, nil
}
```

**Performance Optimization Strategies:**
- **Pre-calculate Constants**: Move expensive math outside animation loops
- **Grid-based Rendering**: Use 2D arrays for pixel-like effects instead of string concatenation
- **Efficient String Building**: Use `strings.Builder` for complex text assembly
- **State Diffing**: Only update changed portions of complex UIs
- **Frame Skipping**: Skip visual updates under high load while maintaining timing

### Lipgloss Styling Mastery

**Color System:**
```go
// True Color (16.7M colors)
lipgloss.Color("#FF5733")

// ANSI 256 colors
lipgloss.Color("196")

// ANSI 16 colors  
lipgloss.Color("9")

// Adaptive colors (light/dark theme aware)
lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}

// Complete color specification (no auto-degradation)
lipgloss.CompleteColor{
    TrueColor: "#FF5733", 
    ANSI256: "196", 
    ANSI: "9"
}
```

**Layout and Positioning:**
```go
style := lipgloss.NewStyle().
    Width(40).
    Height(10).
    Align(lipgloss.Center).          // Horizontal alignment
    
    Padding(1, 2, 1, 2).            // Top, Right, Bottom, Left
    Margin(0, 1).                   // Vertical, Horizontal
    
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#FF5733"))

// Advanced layout
content := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
row := lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
```

**Advanced Styling Techniques:**
```go
// Conditional styling
style := baseStyle
if isActive {
    style = style.Bold(true).Foreground(lipgloss.Color("#00FF00"))
}

// Style inheritance (only unset properties inherited)
childStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#FF0000")).
    Inherit(parentStyle)

// Style copying and modification
newStyle := existingStyle.Background(lipgloss.Color("#000000"))

// Measurement and sizing
width := lipgloss.Width(styledText)
height := lipgloss.Height(styledText)
w, h := lipgloss.Size(styledText)
```

### Bubbles Component Library Integration

**List Component:**
```go
import "github.com/charmbracelet/bubbles/list"

// Item interface implementation
type item struct {
    title, desc string
}
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// List setup
items := []list.Item{item{title: "Example", desc: "Description"}}
l := list.New(items, list.NewDefaultDelegate(), 50, 10)
l.Title = "My List"

// In Update:
var cmd tea.Cmd
m.list, cmd = m.list.Update(msg)
return m, cmd
```

**Text Input Component:**
```go
import "github.com/charmbracelet/bubbles/textinput"

// Setup
ti := textinput.New()
ti.Placeholder = "Enter text..."
ti.Focus()

// In Update:
m.textInput, cmd = m.textInput.Update(msg)
```

**Common Integration Pattern:**
```go
type model struct {
    list      list.Model
    textInput textinput.Model
    table     table.Model
    focused   int // Which component has focus
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd
    
    // Update focused component
    switch m.focused {
    case 0:
        m.list, cmd = m.list.Update(msg)
    case 1:
        m.textInput, cmd = m.textInput.Update(msg)
    }
    cmds = append(cmds, cmd)
    
    return m, tea.Batch(cmds...)
}
```

### Advanced Patterns and Best Practices

**Error Handling:**
```go
type errorMsg struct{ err error }

func (e errorMsg) Error() string { return e.err.Error() }

// In Update:
case errorMsg:
    m.error = msg.err
    m.loading = false
    return m, nil
```

**State Management for Complex Apps:**
```go
// Nested models for complex applications
type model struct {
    currentView int
    homeModel   home.Model
    editorModel editor.Model
    helpModel   help.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.currentView {
    case 0:
        var cmd tea.Cmd
        m.homeModel, cmd = m.homeModel.Update(msg)
        return m, cmd
    // Handle other views...
    }
}
```

**Responsive Design:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        
        // Update child components
        m.list.SetSize(msg.Width, msg.Height-4)
        
        // Recalculate layouts
        m.recalculateLayout()
        return m, nil
    }
}
```

### Debugging and Development

**Debug Logging (when stdout is occupied by TUI):**
```go
if len(os.Getenv("DEBUG")) > 0 {
    f, err := tea.LogToFile("debug.log", "debug")
    if err != nil {
        fmt.Println("fatal:", err)
        os.Exit(1)
    }
    defer f.Close()
}
```

**Delve Debugging:**
```bash
# Terminal 1: Start headless debugger
dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 .

# Terminal 2: Connect to debugger
dlv connect 127.0.0.1:43000
```

### Common Pitfalls and Solutions

**❌ Avoid These Mistakes:**
- Using goroutines directly (use commands instead)
- Performing I/O in Update() or View() methods
- Modifying model state outside Update()
- Ignoring WindowSizeMsg for responsive layouts
- Using `interface{}` instead of specific message types
- Forgetting to handle all message types in switches

**✅ Best Practices:**
- Always handle WindowSizeMsg for responsive design
- Use type switches for message handling
- Batch commands when possible for efficiency
- Implement proper focus management for multi-component UIs
- Use Lipgloss consistently for styling
- Test with different terminal sizes and capabilities

## Advanced Simulation and Animation Techniques

### Grid-Based Rendering for Complex Visuals

For pixel-level control and complex visual effects, use 2D grids instead of string concatenation:

```go
type model struct {
    grid   [][]string  // Visual grid
    width  int
    height int
}

func (m *model) initGrid() {
    m.grid = make([][]string, m.height)
    for i := range m.grid {
        m.grid[i] = make([]string, m.width)
        for j := range m.grid[i] {
            m.grid[i][j] = " " // Initialize with spaces
        }
    }
}

func (m model) View() string {
    lines := make([]string, len(m.grid))
    for i, row := range m.grid {
        lines[i] = strings.Join(row, "")
    }
    return strings.Join(lines, "\n")
}
```

**Grid Benefits:**
- Pixel-precise control for complex animations
- Efficient for particle systems and cellular automata  
- Easy collision detection and spatial queries
- Natural coordinate system for physics simulations

### Physics Simulation Patterns

**Particle System Architecture:**
```go
type particle struct {
    // Position and velocity
    x, y   float64
    vx, vy float64
    
    // Lifecycle management
    life   float64  // 0.0 to 1.0
    age    int      // Frames alive
    
    // Visual properties
    char   string
    color  lipgloss.Color
    size   float64
}

func (p *particle) update(gravity, wind float64) {
    // Apply forces
    p.vy += gravity
    p.vx += wind
    
    // Update position
    p.x += p.vx
    p.y += p.vy
    
    // Age particle
    p.life -= 0.02
    p.age++
}

func (m *model) updateParticles() {
    alive := []particle{}
    for i := range m.particles {
        p := &m.particles[i]
        p.update(m.gravity, m.wind)
        
        // Keep alive particles within bounds
        if p.life > 0 && p.x >= 0 && p.x < float64(m.width) && 
           p.y >= 0 && p.y < float64(m.height) {
            alive = append(alive, *p)
        }
    }
    m.particles = alive
}
```

**Cellular Automata (Game of Life) Pattern:**
```go
type cell struct {
    alive bool
    age   int  // For visual aging effects
}

func (m *model) evolveGeneration() {
    newGrid := make([][]cell, m.height)
    for i := range newGrid {
        newGrid[i] = make([]cell, m.width)
    }
    
    for y := 0; y < m.height; y++ {
        for x := 0; x < m.width; x++ {
            neighbors := m.countNeighbors(x, y)
            current := m.grid[y][x]
            
            // Conway's rules
            if current.alive {
                newGrid[y][x].alive = neighbors == 2 || neighbors == 3
                if newGrid[y][x].alive {
                    newGrid[y][x].age = current.age + 1
                }
            } else {
                newGrid[y][x].alive = neighbors == 3
                newGrid[y][x].age = 0
            }
        }
    }
    
    m.grid = newGrid
    m.generation++
}
```

### Mathematical and 3D Projection Techniques

**3D Starfield with Perspective:**
```go
type star struct {
    x, y, z    float64  // 3D coordinates
    prevX, prevY float64  // For trail effects
}

func (s *star) project(centerX, centerY, focalLength float64) (int, int, bool) {
    if s.z <= 0 {
        return 0, 0, false // Behind camera
    }
    
    // Perspective projection
    screenX := s.x * focalLength / s.z + centerX
    screenY := s.y * focalLength / s.z + centerY
    
    return int(screenX), int(screenY), true
}

func (m *model) updateStarfield() {
    for i := range m.stars {
        s := &m.stars[i]
        
        // Store previous position for trails
        s.prevX, s.prevY, _ = s.project(m.centerX, m.centerY, 50.0)
        
        // Move star towards camera
        s.z -= m.speed
        
        // Reset if star passes camera
        if s.z <= 0 {
            s.x = (rand.Float64() - 0.5) * 2
            s.y = (rand.Float64() - 0.5) * 2
            s.z = 1.0
        }
    }
}
```

**Wave and Ripple Mathematics:**
```go
// Sine wave with multiple parameters
func calculateWave(x, time, amplitude, frequency, phase, speed float64) float64 {
    return amplitude * math.Sin(2*math.Pi*(frequency*x + speed*time) + phase)
}

// Ripple effect from point source
func calculateRipple(x, y, centerX, centerY, time, amplitude, wavelength float64) float64 {
    distance := math.Sqrt((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY))
    return amplitude * math.Sin(2*math.Pi*(distance/wavelength - time)) / (1 + distance*0.1)
}

// Combine multiple wave sources
func (m model) calculateWaveHeight(x, y float64) float64 {
    height := 0.0
    for _, wave := range m.waves {
        height += calculateWave(x, m.time, wave.amplitude, wave.frequency, wave.phase, wave.speed)
    }
    for _, ripple := range m.ripples {
        height += calculateRipple(x, y, ripple.x, ripple.y, m.time, ripple.amplitude, ripple.wavelength)
    }
    return height
}
```

### Unicode Art and Visual Effects

**Character Progression for Visual Intensity:**
```go
// Height-based characters (waves, bars)
var heightChars = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// Density-based characters (particles, dots)  
var densityChars = []string{" ", "░", "▒", "▓", "█"}

// Particle characters by type
var sparkleChars = []string{"✦", "✧", "⋆", "◦", "•", "∘", "○", "◌"}
var fireChars = []string{"▄", "▀", "█", "▓", "▒", "░"}

// Matrix-style characters (authentic Japanese Katakana)
var matrixChars = []rune("ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ0123456789")

func getCharacterByIntensity(intensity float64, chars []string) string {
    index := int(common.Clamp(intensity * float64(len(chars)), 0, float64(len(chars)-1)))
    return chars[index]
}
```

**Color Gradients for Depth and Effects:**
```go
// Distance-based color fading (starfield)
func getStarColor(distance float64) lipgloss.Color {
    colors := []string{"#FFFFFF", "#CCCCCC", "#999999", "#666666", "#333333"}
    index := int(distance * float64(len(colors)-1))
    index = int(common.Clamp(float64(index), 0, float64(len(colors)-1)))
    return lipgloss.Color(colors[index])
}

// Age-based color transitions (Game of Life)
func getCellColor(age int) lipgloss.Color {
    switch {
    case age < 5:  return lipgloss.Color("#00FF00")  // Young - bright green
    case age < 15: return lipgloss.Color("#FFFF00")  // Mature - yellow  
    case age < 30: return lipgloss.Color("#FF8800")  // Old - orange
    default:       return lipgloss.Color("#FF0000")  // Ancient - red
    }
}
```

### Performance Optimization for Complex Simulations

**Memory-Efficient Particle Management:**
```go
// Pre-allocate particle slice and reuse
type model struct {
    particles     []particle
    activeCount   int      // Track active particles
    maxParticles  int      // Pool size
}

func (m *model) emitParticle() {
    if m.activeCount < m.maxParticles {
        // Reuse existing slot
        p := &m.particles[m.activeCount]
        p.reset() // Reset particle properties
        m.activeCount++
    }
}

func (m *model) updateParticles() {
    aliveCount := 0
    for i := 0; i < m.activeCount; i++ {
        p := &m.particles[i]
        p.update()
        
        if p.isAlive() {
            // Move alive particle to front of slice
            if aliveCount != i {
                m.particles[aliveCount] = *p
            }
            aliveCount++
        }
    }
    m.activeCount = aliveCount
}
```

**Spatial Optimization for Large Simulations:**
```go
// Spatial partitioning for efficient neighbor queries
type spatialGrid struct {
    cellSize int
    cells    map[string][]int // Key: "x,y", Value: particle indices
}

func (sg *spatialGrid) addParticle(index int, x, y float64) {
    cellX := int(x / float64(sg.cellSize))
    cellY := int(y / float64(sg.cellSize))
    key := fmt.Sprintf("%d,%d", cellX, cellY)
    sg.cells[key] = append(sg.cells[key], index)
}

func (sg *spatialGrid) getNeighbors(x, y float64) []int {
    neighbors := []int{}
    cellX := int(x / float64(sg.cellSize))
    cellY := int(y / float64(sg.cellSize))
    
    // Check surrounding cells
    for dx := -1; dx <= 1; dx++ {
        for dy := -1; dy <= 1; dy++ {
            key := fmt.Sprintf("%d,%d", cellX+dx, cellY+dy)
            neighbors = append(neighbors, sg.cells[key]...)
        }
    }
    return neighbors
}
```

### Interactive Parameter Control Patterns

**Real-time Parameter Adjustment:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        // Simulation speed
        case "+", "=":
            m.speed = common.Clamp(m.speed*1.1, 0.01, 2.0)
        case "-", "_":
            m.speed = common.Clamp(m.speed*0.9, 0.01, 2.0)
            
        // Physics parameters
        case "g":
            m.gravity = -m.gravity  // Flip gravity
        case "up":
            m.gravity -= 0.01
        case "down":
            m.gravity += 0.01
        case "left":
            m.wind -= 0.05
        case "right":
            m.wind += 0.05
            
        // Visual modes
        case "1", "2", "3", "4", "5":
            m.mode = msg.String()
            m.applyModeSettings()
            
        // Simulation control
        case "space":
            m.paused = !m.paused
        case "r":
            m.reset()
        case "c":
            m.clear()
        }
    }
    return m, nil
}
```

### Multi-Mode Application Architecture

**Mode-Based Behavior:**
```go
type model struct {
    mode     string
    modes    map[string]modeConfig
}

type modeConfig struct {
    name        string
    description string
    settings    map[string]interface{}
}

func (m *model) applyModeSettings() {
    config, exists := m.modes[m.mode]
    if !exists {
        return
    }
    
    // Apply mode-specific settings
    if gravity, ok := config.settings["gravity"].(float64); ok {
        m.gravity = gravity
    }
    if speed, ok := config.settings["speed"].(float64); ok {
        m.speed = speed
    }
    if particleCount, ok := config.settings["particles"].(int); ok {
        m.maxParticles = particleCount
    }
}

func (m model) View() string {
    // Show current mode in UI
    modeInfo := lipgloss.NewStyle().
        Foreground(common.Cyan).
        Render(fmt.Sprintf("Mode: %s", m.modes[m.mode].name))
    
    return m.renderSimulation() + "\n" + modeInfo + "\n" + m.renderControls()
}
```

These advanced techniques enable sophisticated simulations and visual effects while maintaining the clean TEA architecture and good performance in terminal environments.