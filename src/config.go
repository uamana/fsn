package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/anmitsu/go-shlex"
	"github.com/fsnotify/fsnotify"
)

// WatchItem is a single item for fsevent watch
type pathParams struct {
	Cmd       string   `json:"cmd"`
	Pause     int      `json:"pause"`
	Throttle  int      `json:"throttle"`
	Modes     []string `json:"modes"`
	LogOutput bool     `json:"log_output"`
	WorkDir   string   `json:"work_dir"`
	Env       []string `json:"env"`
}

// Config is a program config
type jsonConfig struct {
	Workers int                   `json:"workers"`
	Watch   map[string]pathParams `json:"watch"`
}

// PathParams contains parameters describing the watched path
type PathParams struct {
	Pause     time.Duration
	Throttle  time.Duration
	Cmd       []string
	Op        fsnotify.Op
	LogOutput bool
	WorkDir   string
	Env       []string
}

type Config struct {
	Workers int
	Watch   map[string]*PathParams
}

// LoadConfig load program config in JSON format from file specified by fname
func LoadConfig(fname string) (*Config, error) {
	var (
		tmpConfig jsonConfig
		config    Config
	)
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &tmpConfig)
	if err != nil {
		return nil, err
	}
	config.Workers = tmpConfig.Workers
	config.Watch = make(map[string]*PathParams)

	var op fsnotify.Op
	for path, p := range tmpConfig.Watch {
		cmd, err := shlex.Split(p.Cmd, runtime.GOOS != "windows")
		if err != nil {
			return nil, err
		}
		if len(p.Modes) == 0 {
			op = 0
		} else {
			if find(p.Modes, "create") >= 0 {
				op |= fsnotify.Create
			}
			if find(p.Modes, "write") >= 0 {
				op |= fsnotify.Write
			}
			if find(p.Modes, "remove") >= 0 {
				op |= fsnotify.Remove
			}
			if find(p.Modes, "rename") >= 0 {
				op |= fsnotify.Rename
			}
			if find(p.Modes, "chmod") >= 0 {
				op |= fsnotify.Chmod
			}

		}
		var env []string
		if len(p.Env) > 0 {
			env = append(os.Environ(), p.Env...)
		}
		config.Watch[filepath.Clean(path)] = &PathParams{
			Pause:     time.Duration(p.Pause) * time.Millisecond,
			Throttle:  time.Duration(p.Throttle) * time.Millisecond,
			Cmd:       cmd,
			Op:        op,
			LogOutput: p.LogOutput,
			WorkDir:   filepath.Clean(p.WorkDir),
			Env:       env,
		}
	}
	return &config, nil
}

func find(s []string, v string) int {
	for i, item := range s {
		if item == v {
			return i
		}
	}
	return -1
}
