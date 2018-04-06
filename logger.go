// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fibre

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger ...
type Logger struct {
	*logrus.Logger
}

// NewLogger returns a new Logger instance.
func NewLogger(f *Fibre) *Logger {

	return &Logger{&logrus.Logger{
		Out:       os.Stderr,
		Level:     logrus.ErrorLevel,
		Hooks:     make(logrus.LevelHooks),
		Formatter: new(logrus.TextFormatter),
	}}

}

// SetLogger sets the logrus instance.
func (l *Logger) SetLogger(i *logrus.Logger) {
	l.Logger = i
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(v string) {
	switch v {
	case "debug", "DEBUG":
		l.Level = logrus.DebugLevel
	case "info", "INFO":
		l.Level = logrus.InfoLevel
	case "warning", "WARNING":
		l.Level = logrus.WarnLevel
	case "error", "ERROR":
		l.Level = logrus.ErrorLevel
	case "fatal", "FATAL":
		l.Level = logrus.FatalLevel
	case "panic", "PANIC":
		l.Level = logrus.PanicLevel
	}
}

// SetFormat sets the logging format.
func (l *Logger) SetFormat(v string) {
	switch v {
	case "text":
		l.Formatter = &logrus.TextFormatter{}
	case "json":
		l.Formatter = &logrus.JSONFormatter{}
	}
}
