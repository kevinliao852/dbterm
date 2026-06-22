package pages

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevinliao852/dbterm/pkg/views"
	log "github.com/sirupsen/logrus"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = views.HelpStyle.PaddingLeft(2).PaddingBottom(1)
)

const (
	NotSaveSession int = iota
	SaveSession
)

type ConfirmPage struct {
	list        list.Model
	DB          *sql.DB
	driverType  string
	saveSession *int
	width       int
	height      int
}

type item string

func (i item) FilterValue() string {
	return ""
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprint(w, views.ActiveItemStyle.Render("● "+string(i)))
		return
	}

	fmt.Fprint(w, views.InactiveItemStyle.Render("○ "+string(i)))
}

func NewConfirmPage() ConfirmPage {
	items := []list.Item{
		item("Yes"),
		item("No"),
	}
	const defaultWidth = 40
	const defaultHeight = 10

	l := list.New(items, itemDelegate{}, defaultWidth, defaultHeight)

	l.Title = "Save this connection?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.Styles.Title = views.PageTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	confirmPage := ConfirmPage{
		list:        l,
		DB:          nil,
		saveSession: nil,
	}
	return confirmPage
}

var _ tea.Model = (*ConfirmPage)(nil)

func (c ConfirmPage) Init() tea.Cmd {
	return nil
}

func intPtr(v int) *int {
	return &v
}

func (c ConfirmPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		c.width = sizeMsg.Width
		c.height = sizeMsg.Height
		c.resize()
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i, ok := c.list.SelectedItem().(item)
			if ok {
				if string(i) == "Yes" {
					c.saveSession = intPtr(SaveSession)
					log.Println("Save to config")
				}
				if string(i) == "No" {
					c.saveSession = intPtr(NotSaveSession)
					log.Println("Not save to config")
				}
			}

			if c.saveSession != nil {
				options := &map[string]interface{}{}
				log.Println(c.DB)
				(*options)["db"] = c.DB
				(*options)["driverType"] = c.driverType
				n := Navigator{
					To:      QueryPageType,
					Options: options,
				}
				return c, n.NavigateTo
			}
		}
	}
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c *ConfirmPage) resize() {
	if c.width <= 0 || c.height <= 0 {
		return
	}
	c.list.SetSize(max(14, c.contentWidth()-6), max(5, c.height-12))
}

func (c ConfirmPage) View() string {
	if c.saveSession == nil {
		content := views.MutedStyle.Render("Keep these connection details for later.") + "\n\n" +
			c.list.View()
		return views.CardStyle(c.contentWidth()).Render(content)
	}
	if c.saveSession == intPtr(SaveSession) {
		return views.SuccessStyle.Render("✓ Connection selected for saving.")
	}
	return views.MutedStyle.Render("Connection will not be saved.")
}

func (c ConfirmPage) contentWidth() int {
	if c.width <= 0 {
		return 48
	}
	return max(20, min(64, c.width-8))
}
