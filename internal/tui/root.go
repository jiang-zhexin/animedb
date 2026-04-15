package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jiang-zhexin/animedb/internal/tui/common"
)

type RootModel struct {
	ctx          common.Ctx
	currentModel tea.Model
	models       []tea.Model
}

func NewRootModel(ctx common.Ctx, maker common.NextMsg) RootModel {
	return RootModel{
		ctx:          ctx,
		currentModel: maker(ctx),
		models:       make([]tea.Model, 0),
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case common.NextMsg:
		newModel := msg(m.ctx)
		m.models = append(m.models, m.currentModel)
		m.currentModel = newModel
		return m, tea.Sequence(m.currentModel.Init(), common.CmdHandler(m.ctx.WindowSizeMsg))

	case common.ExitMSg:
		if len(m.models) == 0 {
			return m, tea.Quit
		}
		m.currentModel, m.models = common.Pop(m.models)
		return m, tea.Batch(common.CmdHandler(m.ctx.WindowSizeMsg), msg.DeferFunc)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			if len(m.models) == 0 {
				return m, tea.Quit
			}
			m.currentModel, m.models = common.Pop(m.models)
			return m, common.CmdHandler(m.ctx.WindowSizeMsg)
		}

	case tea.WindowSizeMsg:
		m.ctx.WindowSizeMsg = msg

		h, v := RootStyle.GetFrameSize()
		msg.Height = msg.Height - v
		msg.Width = msg.Width - h

		m.currentModel, cmd = m.currentModel.Update(msg)
		return m, cmd
	}

	m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m RootModel) View() tea.View {
	v := m.currentModel.View()
	v.AltScreen = true
	v.Content = RootStyle.Render(v.Content)
	return v
}

var RootStyle = lipgloss.NewStyle().Margin(1, 2)
