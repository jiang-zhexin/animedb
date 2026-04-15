package common

import (
	tea "charm.land/bubbletea/v2"
)

type NextMsg = func(ctx Ctx) tea.Model

type ExitMSg struct {
	DeferFunc tea.Cmd
}

func NewExitMsg(msg tea.Msg) ExitMSg {
	return ExitMSg{
		DeferFunc: CmdHandler(msg),
	}
}
