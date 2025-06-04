package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type model struct {
	textarea textarea.Model
	mode     string
	saved    bool
	content  string
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Start typing your message here..."
	ta.Focus()
	ta.CharLimit = 1000
	ta.SetWidth(60)
	ta.SetHeight(10)
	ta.ShowLineNumbers = true
	ta.KeyMap.InsertNewline.SetEnabled(true)

	return model{
		textarea: ta,
		mode:     "edit",
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.mode == "preview" {
				m.mode = "edit"
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyCtrlS:
			// Save content
			m.saved = true
			m.content = m.textarea.Value()
			return m, nil

		case tea.KeyCtrlP:
			// Toggle preview mode
			if m.mode == "edit" {
				m.mode = "preview"
				m.content = m.textarea.Value()
			} else {
				m.mode = "edit"
			}
			return m, nil

		case tea.KeyCtrlR:
			// Reset/clear
			m.textarea.Reset()
			m.saved = false
			m.content = ""
			return m, nil

		case tea.KeyF1:
			// Toggle line numbers
			m.textarea.ShowLineNumbers = !m.textarea.ShowLineNumbers
			return m, nil

		case tea.KeyF2:
			// Toggle word wrap
			m.textarea.KeyMap.InsertNewline.SetEnabled(!m.textarea.KeyMap.InsertNewline.Enabled())
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 10)
		m.textarea.SetHeight(msg.Height - 10)
		return m, nil
	}

	// Only update textarea in edit mode
	if m.mode == "edit" {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Green).
		Padding(0, 1).
		MarginBottom(1)

	title := titleStyle.Render("📄 Textarea Component")

	// Mode indicator
	modeStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		MarginLeft(2)

	var modeIndicator string
	if m.mode == "edit" {
		modeIndicator = modeStyle.
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(common.Blue).
			Render("✏️ EDIT MODE")
	} else {
		modeIndicator = modeStyle.
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(common.Purple).
			Render("👁️ PREVIEW MODE")
	}

	header := lipgloss.JoinHorizontal(lipgloss.Center, title, modeIndicator)

	// Stats
	statsStyle := lipgloss.NewStyle().
		Foreground(common.Cyan).
		MarginTop(1).
		MarginBottom(1)

	lines := len(strings.Split(m.textarea.Value(), "\n"))
	chars := len([]rune(m.textarea.Value()))
	words := len(strings.Fields(m.textarea.Value()))

	stats := statsStyle.Render(fmt.Sprintf(
		"Lines: %d | Words: %d | Characters: %d/%d",
		lines, words, chars, m.textarea.CharLimit,
	))

	// Save indicator
	var saveIndicator string
	if m.saved {
		saveStyle := lipgloss.NewStyle().
			Foreground(common.Green).
			Bold(true)
		saveIndicator = saveStyle.Render(" ✅ Saved")
	}

	// Main content area
	var content string
	if m.mode == "edit" {
		// Edit mode - show textarea
		textareaStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(common.Blue).
			Padding(1)

		content = textareaStyle.Render(m.textarea.View())
	} else {
		// Preview mode - show formatted content
		previewStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(common.Purple).
			Padding(1).
			Width(m.textarea.Width() + 2).
			Height(m.textarea.Height() + 2)

		previewContent := m.content
		if previewContent == "" {
			previewContent = "Nothing to preview yet..."
		}

		// Simple markdown-like formatting
		lines := strings.Split(previewContent, "\n")
		formattedLines := make([]string, len(lines))
		for i, line := range lines {
			if strings.HasPrefix(line, "# ") {
				// Header
				headerStyle := lipgloss.NewStyle().
					Bold(true).
					Foreground(common.Yellow).
					MarginBottom(1)
				formattedLines[i] = headerStyle.Render(strings.TrimPrefix(line, "# "))
			} else if strings.HasPrefix(line, "- ") {
				// Bullet point
				bulletStyle := lipgloss.NewStyle().
					Foreground(common.Green)
				formattedLines[i] = bulletStyle.Render("• " + strings.TrimPrefix(line, "- "))
			} else if strings.HasPrefix(line, "*") && strings.HasSuffix(line, "*") {
				// Italic
				italicStyle := lipgloss.NewStyle().
					Italic(true).
					Foreground(common.Cyan)
				formattedLines[i] = italicStyle.Render(strings.Trim(line, "*"))
			} else {
				formattedLines[i] = line
			}
		}

		content = previewStyle.Render(strings.Join(formattedLines, "\n"))
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Faint(true).
		MarginTop(1)

	var help string
	if m.mode == "edit" {
		help = helpStyle.Render(
			"[Ctrl+S] save • [Ctrl+P] preview • [Ctrl+R] reset • [F1] line numbers • [F2] word wrap • [Esc] quit",
		)
	} else {
		help = helpStyle.Render(
			"[Ctrl+P] back to edit • [Esc] back to edit • Preview supports: # headers, - bullets, *italic*",
		)
	}

	// Feature indicators
	featureStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Faint(true)

	features := []string{}
	if m.textarea.ShowLineNumbers {
		features = append(features, "Line Numbers: ON")
	} else {
		features = append(features, "Line Numbers: OFF")
	}

	if m.textarea.KeyMap.InsertNewline.Enabled() {
		features = append(features, "Word Wrap: ON")
	} else {
		features = append(features, "Word Wrap: OFF")
	}

	featureInfo := featureStyle.Render(strings.Join(features, " | "))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		stats+saveIndicator,
		featureInfo,
		"",
		content,
		help,
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}