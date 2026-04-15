package common

import (
	tea "charm.land/bubbletea/v2"
	"github.com/jiang-zhexin/animedb/internal/model"
)

type Ctx struct {
	*model.Model
	tea.WindowSizeMsg
}
