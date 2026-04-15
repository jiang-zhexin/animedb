package editseries

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/jiang-zhexin/animedb/internal/bangumi"
	"github.com/jiang-zhexin/animedb/internal/tui/common"
)

type SeriesModel struct {
	ctx        common.Ctx
	seriesItem list.Model
}

func newSeriesModel(ctx common.Ctx) SeriesModel {
	items := []list.Item{}
	for seriesName, subjectID := range ctx.Model.SeriesNameToSubjectID {
		subject, _ := ctx.Model.GetSubjectById(subjectID)
		items = append(items, seriesItem{
			seriesName: seriesName,
			subject:    subject,
		})
	}

	l := list.New(items, list.NewDefaultDelegate(), ctx.Width, ctx.Height)
	l.Title = "edit series name"

	return SeriesModel{
		ctx:        ctx,
		seriesItem: l,
	}
}

func (m SeriesModel) Init() tea.Cmd {
	return nil
}

func (m SeriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateSeriesMsg:
		items := []list.Item{}
		for seriesName, subjectID := range m.ctx.Model.SeriesNameToSubjectID {
			subject, _ := m.ctx.Model.GetSubjectById(subjectID)
			items = append(items, seriesItem{
				seriesName: seriesName,
				subject:    subject,
			})
		}
		return m, m.seriesItem.SetItems(items)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Enter):
			si, ok := m.seriesItem.SelectedItem().(seriesItem)
			if ok {
				return m, common.CmdHandler(updateSearchMsg{seriesName: si.seriesName})
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.WindowSizeMsg = msg
		m.seriesItem.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.seriesItem, cmd = m.seriesItem.Update(msg)
	return m, cmd
}

func (m SeriesModel) View() tea.View {
	return tea.NewView(m.seriesItem.View())
}

type seriesItem struct {
	seriesName string
	subject    *bangumi.Subject
}

func (si seriesItem) Title() string {
	return fmt.Sprintf("series name: %s", si.seriesName)
}

func (si seriesItem) Description() string {
	return fmt.Sprintf("bgm id: %d, name: %s", si.subject.Id, si.subject.NameCn)
}

func (si seriesItem) FilterValue() string {
	return si.seriesName
}

type updateSeriesMsg struct{}

type updateSearchMsg struct {
	seriesName string
}
