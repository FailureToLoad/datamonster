package logger

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/failuretoload/datamonster/request"
	"github.com/failuretoload/datamonster/testenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextHandlerAddsCorrelationIDAndUserID(t *testing.T) {
	var buf bytes.Buffer
	base := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	log := slog.New(NewContextHandler(base))

	req := httptest.NewRequest(http.MethodGet, "http://example.test/", nil)
	rec := httptest.NewRecorder()
	settlementID := testenv.UUID()

	var handler http.Handler = http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = request.SetUserID(ctx, "user-123")
		ctx = request.SetSettlementID(ctx, settlementID)
		log.InfoContext(ctx, "hello")
	})
	handler = request.CorrelationIDMiddleware(handler)
	handler.ServeHTTP(rec, req)

	line := strings.TrimSpace(buf.String())
	require.NotEmpty(t, line, "expected log output")
	assert.Contains(t, line, "userID=user-123")
	assert.Contains(t, line, fmt.Sprintf("settlementID=%s", settlementID.String()))
	assert.Regexp(t, regexp.MustCompile(`\bcorrelationID=\S+`), line)
}
