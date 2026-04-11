package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/jiang-zhexin/animedb/internal/model"
)

type Nexter = func(ctx *Ctx) tea.Model

type Exiter struct {
	deferFunc tea.Cmd
}

func NewExiter(msg tea.Msg) Exiter {
	return Exiter{
		deferFunc: CmdHandler(msg),
	}
}

type Ctx struct {
	*model.Model
	tea.WindowSizeMsg
}
