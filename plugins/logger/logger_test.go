package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestLogger(t *testing.T) {
	logger := Logger()
	if logger == nil {
		t.Errorf("logger is nil, expected a non-nil logger")
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name  string
		logFn func(string, ...interface{})
		level hclog.Level
	}{
		{"Trace", Trace, hclog.Trace},
		{"Debug", Debug, hclog.Debug},
		{"Info", Info, hclog.Info},
		{"Warn", Warn, hclog.Warn},
		{"Error", Error, hclog.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				os.Setenv("YATAS_LOG", "")
			}()

			os.Setenv("YATAS_LOG", strings.ToLower(tt.name))

			var buf bytes.Buffer

			logger := hclog.New(&hclog.LoggerOptions{
				Level:                    hclog.LevelFromString(strings.ToLower(tt.name)),
				Output:                   &buf,
				TimeFormat:               "15:04:05",
				IncludeLocation:          true,
				AdditionalLocationOffset: 1,
				Color:                    hclog.AutoColor,
				ColorHeaderOnly:          true,
			})
			SetLogger(logger)

			message := "test message"
			tt.logFn(message)

			logOutput := buf.String()
			if !strings.Contains(logOutput, message) {
				t.Errorf("Expected log message '%s' in output: %s", message, logOutput)
			}

			// tt.name but uppercased
			if !strings.Contains(logOutput, strings.ToUpper(tt.name)) {
				t.Errorf("Expected log level '%s' in output: %s", tt.name, logOutput)
			}
		})
	}
}
