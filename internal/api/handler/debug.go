package handler

import (
	"github.com/xiongjia/beacon/internal/core"
)

type (
	DebugHandler struct {
		dbgService core.DebugService
	}
)

func NewDebugHandler() *DebugHandler {
	return &DebugHandler{}
}
