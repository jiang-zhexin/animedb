package editseries

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jiang-zhexin/animedb/internal/tui/common"
)

type sessionState uint

const (
	seriesView sessionState = iota
	searchView
)

type RootModel struct {
	ctx              common.Ctx
	state            sessionState
	searchModel      tea.Model
	seriesModel      tea.Model
	searchModelStyle lipgloss.Style
	seriesModelStyle lipgloss.Style
}

func New(ctx common.Ctx) tea.Model {
	m := RootModel{
		ctx: ctx,
		searchModelStyle: lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")),

		seriesModelStyle: lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")),
	}
	m.updateStyle()

	h, v := m.searchModelStyle.GetFrameSize()
	searchCtx := common.Ctx{
		Model: ctx.Model,
		WindowSizeMsg: tea.WindowSizeMsg{
			Height: ctx.Height - v,
			Width:  ctx.Width*70/100 - h,
		},
	}

	h, v = m.searchModelStyle.GetFrameSize()
	seriesCtx := common.Ctx{
		Model: ctx.Model,
		WindowSizeMsg: tea.WindowSizeMsg{
			Height: ctx.Height - v,
			Width:  ctx.Width*30/100 - h,
		},
	}
	m.searchModel = newSearchModel(searchCtx)
	m.seriesModel = newSeriesModel(seriesCtx)

	return m
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
		m.searchModel.Init(),
		m.seriesModel.Init(),
	)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.WindowSizeMsg = msg
		m.updateStyle()

		h, v := m.searchModelStyle.GetFrameSize()
		m.searchModel, cmd = m.searchModel.Update(tea.WindowSizeMsg{
			Height: msg.Height - v,
			Width:  msg.Width*70/100 - h,
		})
		cmds = append(cmds, cmd)

		h, v = m.searchModelStyle.GetFrameSize()
		m.seriesModel, cmd = m.seriesModel.Update(tea.WindowSizeMsg{
			Height: msg.Height - v,
			Width:  msg.Width*30/100 - h,
		})
		cmds = append(cmds, cmd)

	case updateSearchMsg:
		m.state = searchView
		m.updateStyle()
		m.searchModel, cmd = m.searchModel.Update(msg)
		cmds = append(cmds, cmd)

	case updateSeriesMsg:
		m.state = seriesView
		m.updateStyle()
		m.seriesModel, cmd = m.seriesModel.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Tab):
			switch m.state {
			case searchView:
				m.state = seriesView
			case seriesView:
				m.state = searchView
			}

			m.updateStyle()
			return m, nil
		}

		switch m.state {
		case searchView:
			m.searchModel, cmd = m.searchModel.Update(msg)
			cmds = append(cmds, cmd)
		case seriesView:
			m.seriesModel, cmd = m.seriesModel.Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		m.searchModel, cmd = m.searchModel.Update(msg)
		cmds = append(cmds, cmd)
		m.seriesModel, cmd = m.seriesModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() tea.View {
	var s strings.Builder

	s.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.seriesModelStyle.Render(m.seriesModel.View().Content),
		m.searchModelStyle.Render(m.searchModel.View().Content),
	))

	return tea.NewView(s.String())
}

func (m *RootModel) updateStyle() {
	switch m.state {
	case searchView:
		m.searchModelStyle = m.searchModelStyle.BorderStyle(lipgloss.NormalBorder())
		m.seriesModelStyle = m.seriesModelStyle.BorderStyle(lipgloss.HiddenBorder())

	case seriesView:
		m.seriesModelStyle = m.seriesModelStyle.BorderStyle(lipgloss.NormalBorder())
		m.searchModelStyle = m.searchModelStyle.BorderStyle(lipgloss.HiddenBorder())
	}

	msg := m.ctx.WindowSizeMsg

	m.searchModelStyle = m.searchModelStyle.
		Height(msg.Height).
		Width(msg.Width * 70 / 100)

	m.seriesModelStyle = m.seriesModelStyle.
		Height(msg.Height).
		Width(msg.Width * 30 / 100)
}
