package tui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/jiang-zhexin/animedb/internal/bangumi"
)

type EditSeriesModel struct {
	ctx        *Ctx
	seriesItem list.Model
}

func NewEditSeriesModel(ctx *Ctx) tea.Model {
	items := []list.Item{}
	for seriesName, subjectID := range ctx.Model.SeriesNameToSubjectID {
		subject, _ := ctx.Model.GetSubjectById(subjectID)
		items = append(items, seriesItem{
			seriesName: seriesName,
			subject:    subject,
		})
	}
	em := EditSeriesModel{
		ctx:        ctx,
		seriesItem: list.New(items, list.NewDefaultDelegate(), ctx.Width, ctx.Height),
	}
	em.seriesItem.Title = "animedb mange mode > edit series name to subject id"
	return em
}

func (m EditSeriesModel) Init() tea.Cmd { return nil }

func (m EditSeriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateList:
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
		case key.Matches(msg, keys.Enter):
			si, ok := m.seriesItem.SelectedItem().(seriesItem)
			if ok {
				return m, CmdHandler(func(ctx *Ctx) tea.Model {
					return newSearchModel(ctx, si.seriesName)
				})
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.WindowSizeMsg = msg
		h, v := listStyle.GetFrameSize()
		m.seriesItem.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.seriesItem, cmd = m.seriesItem.Update(msg)
	return m, cmd
}

type updateList struct{}

func (m EditSeriesModel) View() tea.View {
	return tea.NewView(listStyle.Render(m.seriesItem.View()))
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
