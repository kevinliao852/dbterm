package pages

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type ConfirmPage struct {
	list        list.Model
	DB          *sql.DB
	saveSession *int
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

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewConfirmPage() ConfirmPage {
	items := []list.Item{
		item("Yes"),
		item("No"),
	}
	const defaultWidth = 40
	const defaultHeight = 10

	l := list.New(items, itemDelegate{}, defaultWidth, defaultHeight)

	l.Title = "Do you want to save this session?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i, ok := c.list.SelectedItem().(item)
			if ok {
				if string(i) == "Yes" {
					c.saveSession = intPtr(1)
					log.Println("Save to config")
				}
				if string(i) == "No" {
					c.saveSession = intPtr(0)
					log.Println("Not save to config")
				}
			}

			if c.saveSession != nil {
				options := &map[string]interface{}{}
				log.Println(c.DB)
				(*options)["db"] = c.DB
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

func (c ConfirmPage) View() string {
	if c.saveSession == nil {
		return "\n" + c.list.View()
	}
	if c.saveSession == intPtr(1) {
		return fmt.Sprintln("Save this session")
	}
	return fmt.Sprintln("Dont save this session")
}
