package tui

import (
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jiang-zhexin/animedb/internal/bangumi"
)

type SearchModel struct {
	ctx         *Ctx
	seriesName  string
	subjectItem list.Model
	spinner     spinner.Model
	state       State
	err         error
}

type State uint

const (
	SearchLoadingState State = iota
	SearchResultsState
	SearchErrorState
)

func newSearchModel(ctx *Ctx, seriesName string) tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return SearchModel{
		ctx:        ctx,
		seriesName: seriesName,
		state:      SearchLoadingState,
		spinner:    s,
	}
}

func (m SearchModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			results, err := m.ctx.Model.SearchSubject(m.seriesName)
			slog.Debug("search results", slog.String("len", fmt.Sprint(len(results))))
			if err != nil {
				return searchErrorMsg{err: err}
			}
			return searchResultsMsg{results: results}
		},
	)

}

type searchErrorMsg struct {
	err error
}

type searchResultsMsg struct {
	results []bangumi.Subject
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case searchErrorMsg:
		m.err = msg.err
		m.state = SearchErrorState
		return m, nil

	case searchResultsMsg:
		items := []list.Item{}
		for _, subject := range msg.results {
			items = append(items, subjectItem{
				Subject: subject,
			})
		}
		m.subjectItem = list.New(items, list.NewDefaultDelegate(), m.ctx.Width, m.ctx.Height)
		m.subjectItem.Title = fmt.Sprintf("animedb mange mode > edit series name to subject id > choose subject for %s", m.seriesName)
		m.state = SearchResultsState
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			if m.state == SearchResultsState {
				si, ok := m.subjectItem.SelectedItem().(subjectItem)
				if ok {
					slog.Info("update SeriesNameToSubjectID", slog.String("seriesName", m.seriesName), slog.String("id", fmt.Sprint(si.Id)))
					m.err = m.ctx.Model.UpdateSeriesNameToSubjectID(m.seriesName, si.Id)
					if m.err != nil {
						m.state = SearchErrorState
						return m, nil
					}
					return m, CmdHandler(NewExiter(updateList{}))
				}
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.WindowSizeMsg = msg
		if m.state == SearchResultsState {
			h, v := listStyle.GetFrameSize()
			m.subjectItem.SetSize(msg.Width-h, msg.Height-v)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.state == SearchResultsState {
		var cmd tea.Cmd
		m.subjectItem, cmd = m.subjectItem.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SearchModel) View() tea.View {
	switch m.state {
	case SearchResultsState:
		return tea.NewView(listStyle.Render(m.subjectItem.View()))
	case SearchErrorState:
		return tea.NewView(m.err.Error())
	case SearchLoadingState:
		return tea.NewView(fmt.Sprintf("\n\n   %s Search bgm...\n\n", m.spinner.View()))
	default:
		panic(fmt.Sprintf("Unkown state: %d", m.state))
	}
}

type subjectItem struct {
	bangumi.Subject
}

func (si subjectItem) Title() string {
	if si.NameCn != "" {
		return si.NameCn
	} else {
		return si.Name
	}
}

func (si subjectItem) Description() string {
	return fmt.Sprintf("bgm id: %d, name: %s", si.Id, si.Name)
}

func (si subjectItem) FilterValue() string {
	return si.NameCn
}
