package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type model struct {
	viewport viewport.Model
	content  string
	ready    bool
}

func generateLongContent() string {
	content := lipgloss.NewStyle().Bold(true).Foreground(common.Yellow).Render("📜 Welcome to the Viewport Component Demo\n\n")

	sections := []struct {
		title string
		text  string
	}{
		{
			"What is a Viewport?",
			"A viewport is a scrollable container that allows you to display content that's larger than the available screen space. It's perfect for documents, logs, file contents, or any lengthy text that needs to be navigable.",
		},
		{
			"Key Features",
			"• Smooth scrolling with keyboard navigation\n• Mouse wheel support\n• Customizable styling\n• Automatic content wrapping\n• Scroll position indicators\n• Line-by-line or page-by-page navigation",
		},
		{
			"Navigation Controls",
			"• ↑/↓ - Scroll line by line\n• Page Up/Page Down - Scroll page by page\n• Home/End - Jump to top/bottom\n• Mouse wheel - Smooth scrolling\n• g/G - Go to top/bottom (vim-style)",
		},
		{
			"Use Cases",
			"Viewports are commonly used for:\n• Documentation viewers\n• Log file displays\n• Code editors\n• File browsers\n• Chat message history\n• Terminal output\n• Configuration file editors",
		},
		{
			"Lorem Ipsum Content",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n\nDuis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n\nSed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
		},
		{
			"Sample Code",
			"Here's how you might use a viewport in your Bubble Tea application:\n\n```go\nfunc (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {\n    switch msg := msg.(type) {\n    case tea.KeyMsg:\n        switch msg.String() {\n        case \"q\":\n            return m, tea.Quit\n        }\n    }\n    \n    var cmd tea.Cmd\n    m.viewport, cmd = m.viewport.Update(msg)\n    return m, cmd\n}\n```",
		},
		{
			"Advanced Features",
			"The viewport component supports many advanced features:\n\n• Custom key bindings for navigation\n• Programmatic scrolling to specific lines\n• Dynamic content updates\n• Search functionality (when combined with other components)\n• Custom scroll indicators\n• Integration with other Bubble Tea components",
		},
		{
			"Performance Considerations",
			"Viewports are designed to handle large amounts of content efficiently:\n\n• Only visible content is rendered\n• Smooth scrolling animations\n• Memory-efficient content handling\n• Responsive to terminal size changes\n• Optimized for both small and large documents",
		},
		{
			"Styling Options",
			"You can customize the appearance of viewports:\n\n• Border styles and colors\n• Background colors\n• Text formatting\n• Scroll indicators\n• Focus states\n• Custom themes",
		},
		{
			"More Lorem Ipsum",
			"Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Vestibulum tortor quam, feugiat vitae, ultricies eget, tempor sit amet, ante.\n\nDonec eu libero sit amet quam egestas semper. Aenean ultricies mi vitae est. Mauris placerat eleifend leo. Quisque sit amet est et sapien ullamcorper pharetra.\n\nVestibulum erat wisi, condimentum sed, commodo vitae, ornare sit amet, wisi. Aenean fermentum, elit eget tincidunt condimentum, eros ipsum rutrum orci, sagittis tempus lacus enim ac dui.",
		},
	}

	for i, section := range sections {
		// Add section number
		sectionNum := lipgloss.NewStyle().
			Bold(true).
			Foreground(common.Cyan).
			Render(fmt.Sprintf("%d. ", i+1))

		// Style section title
		title := lipgloss.NewStyle().
			Bold(true).
			Foreground(common.Purple).
			Render(section.title)

		// Add section content
		content += sectionNum + title + "\n\n"
		content += section.text + "\n\n"

		// Add separator
		if i < len(sections)-1 {
			separator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render(strings.Repeat("─", 50))
			content += separator + "\n\n"
		}
	}

	// Add footer
	footer := lipgloss.NewStyle().
		Bold(true).
		Foreground(common.Green).
		Render("\n🎉 End of Content\n\nYou've reached the bottom! Press 'g' to go back to the top.")

	content += footer

	return content
}

func initialModel() model {
	return model{
		content: generateLongContent(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "g":
			m.viewport.GotoTop()
			return m, nil
		case "G":
			m.viewport.GotoBottom()
			return m, nil
		case "r":
			// Refresh content
			m.content = generateLongContent()
			m.viewport.SetContent(m.content)
			return m, nil
		}

	case tea.WindowSizeMsg:
		headerHeight := 4
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "Initializing viewport..."
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Blue).
		Padding(0, 1)

	title := titleStyle.Render("📄 Viewport Component")

	// Stats
	statsStyle := lipgloss.NewStyle().
		Foreground(common.Cyan)

	stats := statsStyle.Render(fmt.Sprintf(
		"Position: %d/%d (%.0f%%) | Content lines: %d",
		m.viewport.YOffset+1,
		len(strings.Split(m.content, "\n")),
		m.viewport.ScrollPercent()*100,
		len(strings.Split(m.content, "\n")),
	))

	// Viewport with border
	viewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(common.Purple).
		Padding(0, 1)

	viewportView := viewportStyle.Render(m.viewport.View())

	// Help text
	helpStyle := lipgloss.NewStyle().
		Faint(true)

	help := helpStyle.Render(
		"[↑↓] scroll • [PgUp/PgDn] page • [Home/End] top/bottom • [g/G] vim-style • [r]efresh • [q]uit",
	)

	// Scroll indicator
	scrollStyle := lipgloss.NewStyle().
		Foreground(common.Yellow).
		Bold(true)

	var scrollIndicator string
	if m.viewport.AtTop() {
		scrollIndicator = scrollStyle.Render("▲ TOP")
	} else if m.viewport.AtBottom() {
		scrollIndicator = scrollStyle.Render("▼ BOTTOM")
	} else {
		scrollIndicator = scrollStyle.Render("● SCROLLING")
	}

	// Combine all elements
	header := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", stats, "  ", scrollIndicator)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		viewportView,
		help,
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}