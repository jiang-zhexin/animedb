package common

import tea "charm.land/bubbletea/v2"

func CmdHandler(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func Pop[T any](a []T) (T, []T) {
	if l := len(a); l > 0 {
		return a[l-1], a[:l-1]
	} else {
		var x T
		return x, a
	}
}
