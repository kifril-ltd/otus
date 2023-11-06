package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	cases := []struct {
		name               string
		command            []string
		env                Environment
		expectedReturnCode int
	}{
		{
			name:               "simple exec",
			command:            []string{"ls", "-l"},
			env:                make(Environment),
			expectedReturnCode: 0,
		},
		{
			name:               "exec with fail exit code",
			command:            []string{"/bin/bash", "xxx"},
			env:                make(Environment),
			expectedReturnCode: 127,
		},
		{
			name:    "exec with set env var",
			command: []string{"ls"},
			env: Environment{
				"CUSTOM_VAR": EnvValue{Value: "VALUE"},
			},
			expectedReturnCode: 0,
		},
		{
			name:    "exec with unset env var",
			command: []string{"ls"},
			env: Environment{
				"CUSTOM_VAR": EnvValue{Value: "", NeedRemove: true},
			},
			expectedReturnCode: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			returnCode := RunCmd(c.command, c.env)
			assert.Equal(t, c.expectedReturnCode, returnCode)

			for key, val := range c.env {
				if val.NeedRemove {
					_, ok := os.LookupEnv(key)
					assert.False(t, ok)
				} else {
					act, ok := os.LookupEnv(key)
					assert.True(t, ok)
					assert.Equal(t, val.Value, act)
				}
			}
		})
	}
}
