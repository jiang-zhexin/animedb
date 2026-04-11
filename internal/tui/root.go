package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type RootModel struct {
	ctx          *Ctx
	currentModel tea.Model
	models       []tea.Model
}

func NewRootModel(ctx *Ctx, maker Nexter) RootModel {
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
	switch msg := msg.(type) {
	case Nexter:
		newModel := msg(m.ctx)
		m.models = append(m.models, m.currentModel)
		m.currentModel = newModel
		return m, m.currentModel.Init()

	case Exiter:
		if len(m.models) == 0 {
			return m, tea.Quit
		}
		m.currentModel, m.models = pop(m.models)
		return m, tea.Batch(CmdHandler(m.ctx.WindowSizeMsg), msg.deferFunc)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			if len(m.models) == 0 {
				return m, tea.Quit
			}
			m.currentModel, m.models = pop(m.models)
			return m, CmdHandler(m.ctx.WindowSizeMsg)
		}
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m RootModel) View() tea.View {
	v := m.currentModel.View()
	v.AltScreen = true
	return v
}
