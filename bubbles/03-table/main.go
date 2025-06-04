package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubbletea-showcase/common"
)

type model struct {
	table       table.Model
	selected    table.Row
	action      string
	showDetails bool
	width       int
	height      int
}

func generateSampleData() []table.Row {
	companies := []string{"Apple", "Google", "Microsoft", "Amazon", "Meta", "Tesla", "Netflix", "Adobe", "Salesforce", "Oracle"}
	departments := []string{"Engineering", "Marketing", "Sales", "HR", "Finance", "Support", "Product", "Design"}
	statuses := []string{"Active", "Inactive", "Pending", "Archived"}

	rows := make([]table.Row, 25)
	for i := range rows {
		company := companies[rand.Intn(len(companies))]
		dept := departments[rand.Intn(len(departments))]
		status := statuses[rand.Intn(len(statuses))]
		salary := 50000 + rand.Intn(150000)
		experience := 1 + rand.Intn(15)

		rows[i] = table.Row{
			strconv.Itoa(i + 1001),
			fmt.Sprintf("Employee %d", i+1),
			company,
			dept,
			fmt.Sprintf("$%,d", salary),
			fmt.Sprintf("%d years", experience),
			status,
		}
	}
	return rows
}

func initialModel() model {
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Name", Width: 15},
		{Title: "Company", Width: 12},
		{Title: "Department", Width: 12},
		{Title: "Salary", Width: 10},
		{Title: "Experience", Width: 12},
		{Title: "Status", Width: 10},
	}

	rows := generateSampleData()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	// Custom styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(common.Purple).
		BorderBottom(true).
		Bold(true).
		Foreground(common.Purple)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(common.Purple).
		Bold(true)

	s.Cell = s.Cell.
		Foreground(lipgloss.Color("252"))

	t.SetStyles(s)

	return model{
		table:  t,
		width:  80,
		height: 24,
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

		case "enter":
			// Toggle details panel
			if len(m.table.Rows()) > 0 {
				m.selected = m.table.SelectedRow()
				m.showDetails = !m.showDetails
				m.action = "selected"
				
				// Adjust table height immediately
				tableHeight := m.height - 8
				if m.showDetails {
					tableHeight = m.height - 16
				}
				m.table.SetHeight(tableHeight)
			}
			return m, nil

		case "d":
			// Delete row
			if len(m.table.Rows()) > 0 {
				rows := m.table.Rows()
				cursor := m.table.Cursor()
				if cursor < len(rows) {
					// Remove the selected row
					newRows := append(rows[:cursor], rows[cursor+1:]...)
					m.table.SetRows(newRows)
					m.action = "deleted"
				}
			}
			return m, nil

		case "a":
			// Add new row
			rows := m.table.Rows()
			newID := strconv.Itoa(2000 + len(rows))
			newRow := table.Row{
				newID,
				"New Employee",
				"TechCorp",
				"Engineering",
				"$75,000",
				"2 years",
				"Active",
			}
			rows = append(rows, newRow)
			m.table.SetRows(rows)
			m.action = "added"
			return m, nil

		case "r":
			// Refresh data
			rows := generateSampleData()
			m.table.SetRows(rows)
			m.action = "refreshed"
			m.showDetails = false
			return m, nil

		case "s":
			// Sort by different columns (simple demonstration)
			m.action = "sorted"
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Adjust table size based on whether details are shown
		tableHeight := m.height - 8
		if m.showDetails {
			tableHeight = m.height - 16 // Leave room for details panel at bottom
		}
		
		m.table.SetWidth(m.width - 4)
		m.table.SetHeight(tableHeight)
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(common.Purple).
		Padding(0, 1).
		MarginBottom(1)

	title := titleStyle.Render("📊 Table Component")

	// Action feedback
	actionStyle := lipgloss.NewStyle().
		Foreground(common.Green).
		Bold(true)

	var actionMsg string
	switch m.action {
	case "selected":
		if m.showDetails {
			actionMsg = actionStyle.Render("✓ Details panel opened")
		} else {
			actionMsg = actionStyle.Render("✓ Details panel closed")
		}
	case "deleted":
		actionMsg = actionStyle.Render("🗑️ Row deleted")
		m.showDetails = false // Close details when row is deleted
	case "added":
		actionMsg = actionStyle.Render("➕ Row added")
	case "refreshed":
		actionMsg = actionStyle.Render("🔄 Data refreshed")
	case "sorted":
		actionMsg = actionStyle.Render("↕️ Table sorted")
	}

	// Stats
	statsStyle := lipgloss.NewStyle().
		Foreground(common.Cyan).
		MarginBottom(1)

	stats := statsStyle.Render(fmt.Sprintf(
		"Total rows: %d | Selected: %d",
		len(m.table.Rows()),
		m.table.Cursor()+1,
	))

	// Header
	header := title
	if actionMsg != "" {
		header += "\n" + actionMsg
	}
	header += "\n" + stats

	// Main table
	tableView := m.table.View()

	// Create main content layout
	var mainContent string
	if m.showDetails && len(m.selected) > 0 {
		// Vertical layout: table on top, details on bottom
		detailStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(common.Blue).
			Padding(0, 1).
			Width(m.width - 6).
			MarginTop(1)

		detailContent := lipgloss.NewStyle().Foreground(common.Yellow).Bold(true).Render("Selected Employee Details") + "\n\n"
		
		// Format details in a horizontal layout to save vertical space
		col1 := fmt.Sprintf("ID: %s\nName: %s\nCompany: %s", 
			m.selected[0], m.selected[1], m.selected[2])
		col2 := fmt.Sprintf("Department: %s\nSalary: %s\nExperience: %s", 
			m.selected[3], m.selected[4], m.selected[5])
		col3 := fmt.Sprintf("Status: %s", m.selected[6])
		
		// Create columns for compact display
		col1Style := lipgloss.NewStyle().Width((m.width - 10) / 3)
		col2Style := lipgloss.NewStyle().Width((m.width - 10) / 3)
		col3Style := lipgloss.NewStyle().Width((m.width - 10) / 3)
		
		detailsRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			col1Style.Render(col1),
			col2Style.Render(col2),
			col3Style.Render(col3),
		)
		
		detailContent += detailsRow
		details := detailStyle.Render(detailContent)
		
		// Join table and details vertically
		mainContent = lipgloss.JoinVertical(lipgloss.Left, tableView, details)
	} else {
		// Just the table
		mainContent = tableView
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Faint(true).
		MarginTop(1)

	var helpText string
	if m.showDetails {
		helpText = "[↑↓] navigate • [Enter] hide details • [a]dd row • [d]elete row • [r]efresh • [q]uit"
	} else {
		helpText = "[↑↓] navigate • [Enter] show details • [a]dd row • [d]elete row • [r]efresh • [q]uit"
	}
	help := helpStyle.Render(helpText)

	// Combine all elements vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		mainContent,
		help,
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}