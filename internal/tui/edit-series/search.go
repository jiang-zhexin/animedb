package editseries

import (
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jiang-zhexin/animedb/internal/bangumi"
	"github.com/jiang-zhexin/animedb/internal/tui/common"
)

type state uint

const (
	searchNothing state = iota
	searchLoading
	searchResults
	searchError
)

type SearchModel struct {
	ctx         common.Ctx
	seriesName  string
	subjectItem list.Model
	state       state
	spinner     spinner.Model
	err         error
}

func newSearchModel(ctx common.Ctx) SearchModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return SearchModel{
		ctx:     ctx,
		spinner: s,
	}
}

func (m SearchModel) Init() tea.Cmd {
	return nil
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateSearchMsg:
		m.seriesName = msg.seriesName
		m.state = searchLoading
		return m, tea.Batch(
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

	case searchErrorMsg:
		m.err = msg.err
		m.state = searchError
		return m, nil

	case searchResultsMsg:
		items := []list.Item{}
		for _, subject := range msg.results {
			items = append(items, subjectItem{
				Subject: subject,
			})
		}
		m.subjectItem = list.New(items, list.NewDefaultDelegate(), m.ctx.Width, m.ctx.Height)
		m.subjectItem.Title = fmt.Sprintf("choose subject for %s", m.seriesName)
		m.state = searchResults
		return m, nil

	case tea.KeyMsg:
		if m.state == searchResults {
			switch {
			case key.Matches(msg, common.Keys.Enter):
				si, ok := m.subjectItem.SelectedItem().(subjectItem)
				if ok {
					slog.Info("update SeriesNameToSubjectID", slog.String("seriesName", m.seriesName), slog.String("id", fmt.Sprint(si.Id)))
					m.err = m.ctx.Model.UpdateSeriesNameToSubjectID(m.seriesName, si.Id)
					if m.err != nil {
						m.state = searchError
						return m, nil
					}
					return m, common.CmdHandler(updateSeriesMsg{})
				}
			default:
				var cmd tea.Cmd
				m.subjectItem, cmd = m.subjectItem.Update(msg)
				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.WindowSizeMsg = msg
		if m.state == searchResults {
			m.subjectItem.SetSize(msg.Width, msg.Height)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SearchModel) View() tea.View {
	switch m.state {
	case searchResults:
		return tea.NewView(m.subjectItem.View())
	case searchError:
		return tea.NewView(m.err.Error())
	case searchLoading:
		return tea.NewView(fmt.Sprintf("\n\n   %s Search bgm...\n\n", m.spinner.View()))
	default:
		return tea.NewView("")
	}
}

type searchErrorMsg struct {
	err error
}

type searchResultsMsg struct {
	results []bangumi.Subject
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
