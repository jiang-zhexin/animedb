package tui

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type MainMenuModel struct {
	ctx  *Ctx
	menu list.Model
}

const (
	editSeriesName = menuItem("edit series name to subject id")
	todo           = menuItem("TODO")
)

func NewMainMenuModel(ctx *Ctx) tea.Model {
	items := []list.Item{editSeriesName, todo}

	l := list.New(items, itemDelegate{styles: newStyles()}, ctx.Width, ctx.Height)
	l.Title = "animedb mange mode"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return MainMenuModel{
		menu: l,
		ctx:  ctx,
	}
}

func (m MainMenuModel) Init() tea.Cmd { return nil }

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			i, ok := m.menu.SelectedItem().(menuItem)
			if ok {
				switch i {
				case editSeriesName:
					return m, CmdHandler(NewEditSeriesModel)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.WindowSizeMsg = msg
		h, v := listStyle.GetFrameSize()
		m.menu.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m MainMenuModel) View() tea.View {
	return tea.NewView(listStyle.Render(m.menu.View()))
}

type menuItem string

func (i menuItem) FilterValue() string { return "" }

type itemDelegate struct {
	styles *styles
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type styles struct {
	item         lipgloss.Style
	selectedItem lipgloss.Style
}

func newStyles() *styles {
	var s styles
	s.item = lipgloss.NewStyle().PaddingLeft(4)
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	return &s
}
