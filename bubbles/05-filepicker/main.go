package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

func initialModel() model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".go", ".md", ".txt", ".json", ".yaml", ".yml", ".toml", ".csv"}
	fp.CurrentDirectory, _ = os.Getwd()
	fp.ShowHidden = false
	fp.DirAllowed = true
	fp.FileAllowed = true

	// Custom styles
	fp.Styles.Cursor = lipgloss.NewStyle().Foreground(common.Purple)
	fp.Styles.Symlink = lipgloss.NewStyle().Foreground(common.Cyan)
	fp.Styles.Directory = lipgloss.NewStyle().Foreground(common.Blue).Bold(true)
	fp.Styles.File = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	fp.Styles.Permission = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	fp.Styles.Selected = lipgloss.NewStyle().Foreground(common.Yellow).Bold(true)
	fp.Styles.DisabledCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	fp.Styles.DisabledFile = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return model{
		filepicker: fp,
	}
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "h":
			// Toggle hidden files
			m.filepicker.ShowHidden = !m.filepicker.ShowHidden
			return m, m.filepicker.Init()

		case "r":
			// Refresh directory
			return m, m.filepicker.Init()

		case "~":
			// Go to home directory
			home, err := os.UserHomeDir()
			if err == nil {
				m.filepicker.CurrentDirectory = home
				return m, m.filepicker.Init()
			}

		case "ctrl+h":
			// Go up one directory level
			parent := filepath.Dir(m.filepicker.CurrentDirectory)
			if parent != m.filepicker.CurrentDirectory {
				m.filepicker.CurrentDirectory = parent
				return m, m.filepicker.Init()
			}
		}

	case tea.WindowSizeMsg:
		m.filepicker.Height = msg.Height - 8
		return m, nil

	// Did the user select a file?
	case filepicker.FileSelectedMsg:
		m.selectedFile = msg.Path
		return m, nil

	// Did the user select a disabled file?
	case filepicker.FileSelectedDisabledMsg:
		m.err = fmt.Errorf("file type not allowed: %s", filepath.Ext(msg.Path))
		return m, nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = fmt.Errorf("file type not allowed: %s", filepath.Ext(path))
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Green).
		Padding(0, 1)

	title := titleStyle.Render("📁 File Picker Component")

	// Current directory info
	dirStyle := lipgloss.NewStyle().
		Foreground(common.Cyan).
		Bold(true)

	currentDir := dirStyle.Render(fmt.Sprintf("Current: %s", m.filepicker.CurrentDirectory))

	// File type filter info
	filterStyle := lipgloss.NewStyle().
		Foreground(common.Yellow)

	allowedTypes := strings.Join(m.filepicker.AllowedTypes, ", ")
	if allowedTypes == "" {
		allowedTypes = "All files"
	}
	filter := filterStyle.Render(fmt.Sprintf("Filter: %s", allowedTypes))

	// Hidden files status
	hiddenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))

	hiddenStatus := "Hidden files: "
	if m.filepicker.ShowHidden {
		hiddenStatus += hiddenStyle.Foreground(common.Green).Render("ON")
	} else {
		hiddenStatus += hiddenStyle.Foreground(common.Red).Render("OFF")
	}

	// Header info
	header := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		currentDir,
		filter,
		hiddenStyle.Render(hiddenStatus),
		"",
	)

	// File picker view
	fpView := m.filepicker.View()

	// Selected file or error display
	var footer string
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(common.Red).
			Bold(true)
		footer = errorStyle.Render(fmt.Sprintf("❌ Error: %s", m.err.Error()))
		m.err = nil // Clear error after displaying
	} else if m.selectedFile != "" {
		selectedStyle := lipgloss.NewStyle().
			Foreground(common.Green).
			Bold(true)

		// Get file info
		info, err := os.Stat(m.selectedFile)
		var fileInfo string
		if err == nil {
			if info.IsDir() {
				fileInfo = fmt.Sprintf("📁 Directory selected: %s", m.selectedFile)
			} else {
				size := info.Size()
				var sizeStr string
				if size < 1024 {
					sizeStr = fmt.Sprintf("%d B", size)
				} else if size < 1024*1024 {
					sizeStr = fmt.Sprintf("%.1f KB", float64(size)/1024)
				} else {
					sizeStr = fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
				}
				fileInfo = fmt.Sprintf("📄 File selected: %s (%s)", filepath.Base(m.selectedFile), sizeStr)
			}
		} else {
			fileInfo = fmt.Sprintf("📄 Selected: %s", m.selectedFile)
		}

		footer = selectedStyle.Render(fileInfo)
		
		// Show file path
		pathStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Faint(true)
		footer += "\n" + pathStyle.Render(fmt.Sprintf("Path: %s", m.selectedFile))
		
		// Clear selection after a moment
		m.selectedFile = ""
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Faint(true).
		MarginTop(1)

	help := helpStyle.Render(
		"[↑↓] navigate • [Enter] select • [h] toggle hidden • [r] refresh • [~] home • [Ctrl+H] parent • [q] quit",
	)

	// Combine all elements
	content := header + fpView

	if footer != "" {
		content += "\n\n" + footer
	}

	content += "\n" + help

	return content
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}