package dbg

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDebugServer(t *testing.T) {
	assert := require.New(t)
	dbgMgr := NewDebugManager(
		DebugManagerWithAddr("", 0),
		DebugManagerWithBaseUrl("/internal/debug"))
	assert.NotNil(dbgMgr, "new debug manager")

	err := dbgMgr.Start()
	assert.NoError(err, "start debug server")
	defer func() {
		err = dbgMgr.Stop(30 * time.Second)
		assert.NoError(err, "stop debug server")
	}()

	addr, err := dbgMgr.DebugServerAddr()
	assert.NoError(err, "debug server start error")
	assert.NotEmpty(addr)
	t.Logf("debug server addr %s", addr)

	dbgResp, err := http.Get(fmt.Sprintf("http://%s/internal/debug/pprof/", addr))
	assert.NoError(err, "debug server request error")
	assert.NotNil(dbgResp, "debug server return an invalid response")
	assert.Equal(http.StatusOK, dbgResp.StatusCode, "invalid response code")

	respContent, err := io.ReadAll(dbgResp.Body)
	assert.NoError(err, "debug server return an invalid response")
	assert.NotEmpty(respContent, "invalid response content")
	t.Logf("Response: %s", respContent)
}
