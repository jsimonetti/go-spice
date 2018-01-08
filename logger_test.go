package spice

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggerAdapt(t *testing.T) {
	adapter := cleanLogger(t)

	msg := "myInfoMessageString"
	adapter.Info(msg)

	if !strings.Contains(adapter.String(), "msg="+msg) || !strings.Contains(adapter.String(), "level=info") {
		t.Error("did not find string")
	}
}

func TestLoggerWithError(t *testing.T) {
	adapter := cleanLogger(t)

	msg := "myErrorMessageString"
	error := errors.New(msg)
	adapter.WithError(error).Info()

	if !strings.Contains(adapter.String(), "error="+msg) {
		t.Error("did not find short error message")
	}

	adapter = cleanLogger(t)
	msg = "myErrorMessageString longer"
	error = errors.New(msg)
	adapter.WithError(error).Info()

	if !strings.Contains(adapter.String(), "error=\""+msg+"\"") {
		t.Error("did not find long error message")
	}
}

func TestLoggerWithField(t *testing.T) {
	adapter := cleanLogger(t)

	field := "myField"
	msg := "myFieldMessageString"
	adapter.WithFields(field, msg).Info()

	if !strings.Contains(adapter.String(), field+"="+msg) {
		t.Error("did not find short field message")
	}

	adapter = cleanLogger(t)
	field = "myField"
	msg = "myFieldMessageString longer"
	adapter.WithFields(field, msg).Info()

	if !strings.Contains(adapter.String(), field+"=\""+msg+"\"") {
		t.Error("did not find long field message")
	}
}

func cleanLogger(t *testing.T) *clean {
	t.Helper()

	var buf bytes.Buffer
	logger := logrus.New()
	logger.Out = &buf

	adapter := &clean{
		Adapt(logger.WithField("logger", "clean")),
		&buf,
	}
	return adapter
}

type clean struct {
	Logger
	buf *bytes.Buffer
}

func (c *clean) String() string {
	return c.buf.String()
}
