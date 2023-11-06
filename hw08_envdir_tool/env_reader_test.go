package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	cases := [...]struct {
		name, dir string
		err       error
		result    Environment
	}{
		{
			name: "dir doesn't exist",
			dir:  "./not_found_dir",
			err:  os.ErrNotExist,
		},
		{
			name: "not dir, but file",
			dir:  "./README.md",
			err:  ErrNotADirectory,
		},
		{
			name: "empty folder",
			dir:  "./ENVS",
			err:  nil,
		},
		{
			name: "successful case",
			dir:  "./testdata/env",
			err:  nil,
			result: Environment{
				"BAR":   EnvValue{Value: "bar"},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line"},
				"HELLO": EnvValue{Value: "\"hello\""},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
		},
	}

	_ = os.Mkdir("ENVS", 0o777)

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result, err := ReadDir(c.dir)

			if c.err != nil {
				assert.Error(t, c.err, err)
			}

			for key, act := range result {
				exp, ok := c.result[key]
				assert.True(t, ok, fmt.Sprintf("env with key %s was not expected", key))
				assert.Equal(t, exp.Value, act.Value, fmt.Sprintf("env with key %s value is wrong", key))
				assert.Equal(t, exp.NeedRemove, act.NeedRemove, fmt.Sprintf("env with key %s need to remove is wrong", key))
			}
		})
	}

	_ = os.Remove("ENVS")
}
