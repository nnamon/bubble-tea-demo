package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type model struct {
	inputs    []textinput.Model
	focused   int
	submitted bool
	values    []string
}

func initialModel() model {
	inputs := make([]textinput.Model, 5)

	// Name input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Enter your name"
	inputs[0].Focus()
	inputs[0].CharLimit = 50
	inputs[0].Width = 40

	// Email input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "email@example.com"
	inputs[1].CharLimit = 100
	inputs[1].Width = 40

	// Password input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Password"
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].EchoCharacter = '•'
	inputs[2].CharLimit = 50
	inputs[2].Width = 40

	// Number input
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "Age (numbers only)"
	inputs[3].CharLimit = 3
	inputs[3].Width = 40
	inputs[3].Validate = func(s string) error {
		for _, char := range s {
			if char < '0' || char > '9' {
				return fmt.Errorf("only numbers allowed")
			}
		}
		return nil
	}

	// Custom styled input
	inputs[4] = textinput.New()
	inputs[4].Placeholder = "Custom styled input"
	inputs[4].CharLimit = 100
	inputs[4].Width = 40
	inputs[4].PromptStyle = lipgloss.NewStyle().Foreground(common.Purple)
	inputs[4].TextStyle = lipgloss.NewStyle().Foreground(common.Cyan)
	inputs[4].PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	return model{
		inputs:  inputs,
		focused: 0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			s := msg.String()

			// Submit on Enter when all inputs are filled
			if s == "enter" && m.allInputsFilled() {
				m.submitted = true
				m.values = make([]string, len(m.inputs))
				for i, input := range m.inputs {
					m.values[i] = input.Value()
				}
				return m, nil
			}

			// Navigate between inputs
			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > len(m.inputs)-1 {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, tea.Batch(cmds...)

		case tea.KeyRunes:
			// Handle character input for number validation
			if m.focused == 3 { // Number input
				for _, r := range msg.Runes {
					if r < '0' || r > '9' {
						return m, nil // Ignore non-numeric input
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		for i := range m.inputs {
			m.inputs[i].Width = msg.Width - 20
		}
		return m, nil
	}

	// Update the focused input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) allInputsFilled() bool {
	for _, input := range m.inputs {
		if input.Value() == "" {
			return false
		}
	}
	return true
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Blue).
		Padding(0, 1).
		MarginBottom(1)

	title := titleStyle.Render("📝 Text Input Components")

	if m.submitted {
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(common.Green).
			MarginTop(2)

		result := successStyle.Render("✅ Form Submitted Successfully!\n\n")

		valueStyle := lipgloss.NewStyle().
			Foreground(common.Cyan).
			MarginLeft(2)

		labels := []string{"Name:", "Email:", "Password:", "Age:", "Custom:"}
		for i, value := range m.values {
			displayValue := value
			if i == 2 { // Password field
				displayValue = strings.Repeat("•", len(value))
			}
			result += valueStyle.Render(fmt.Sprintf("%s %s\n", labels[i], displayValue))
		}

		helpStyle := lipgloss.NewStyle().
			Faint(true).
			MarginTop(2)

		result += helpStyle.Render("\nPress [Esc] to quit")

		return title + "\n\n" + result
	}

	content := title + "\n\n"

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(common.Yellow).
		Width(20)

	focusedStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(common.Purple).
		Padding(0, 1)

	blurredStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	labels := []string{
		"Name:",
		"Email:",
		"Password:",
		"Age (numbers only):",
		"Custom Styled:",
	}

	for i, input := range m.inputs {
		label := labelStyle.Render(labels[i])

		var inputView string
		if i == m.focused {
			inputView = focusedStyle.Render(input.View())
		} else {
			inputView = blurredStyle.Render(input.View())
		}

		// Show validation error for number input
		errorMsg := ""
		if i == 3 && input.Err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(common.Red).Faint(true)
			errorMsg = "\n" + errorStyle.Render("⚠ " + input.Err.Error())
		}

		content += fmt.Sprintf("%s\n%s%s\n\n", label, inputView, errorMsg)
	}

	// Progress indicator
	filled := 0
	for _, input := range m.inputs {
		if input.Value() != "" {
			filled++
		}
	}

	progressStyle := lipgloss.NewStyle().Foreground(common.Green)
	progress := progressStyle.Render(fmt.Sprintf("Progress: %d/%d fields completed", filled, len(m.inputs)))

	// Help text
	helpStyle := lipgloss.NewStyle().
		Faint(true).
		MarginTop(1)

	var help string
	if m.allInputsFilled() {
		help = helpStyle.Render("[Tab] navigate • [Enter] submit • [Esc] quit")
	} else {
		help = helpStyle.Render("[Tab] navigate • [Esc] quit")
	}

	return content + progress + "\n" + help
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}