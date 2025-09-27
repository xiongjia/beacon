package util

import (
	"net/http"
	"strings"
)

type (
	GroupMux struct {
		mainMux *http.ServeMux
	}
)

func NewMainMuxGroup() *GroupMux {
	return NewMuxGroup(http.NewServeMux())
}

func NewMuxGroup(mainMux *http.ServeMux) *GroupMux {
	return &GroupMux{mainMux: mainMux}
}

func (g *GroupMux) Group(prefix string, groupMux *http.ServeMux) *http.ServeMux {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	g.mainMux.Handle(prefix+"/", http.StripPrefix(prefix, groupMux))
	return g.mainMux
}

func (g *GroupMux) MainMux() *http.ServeMux {
	return g.mainMux
}
