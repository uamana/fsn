package main

import (
	"flag"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	// Version - program version
	Version = "0.2.4"
)

// Job contains data describing the task for worker
type Job struct {
	timestamp time.Time
	params    *PathParams
}

func main() {
	var (
		configFile string
	)
	flag.StringVar(&configFile, "config", "fsn.cfg", "Path to config file")
	flag.Parse()

	log.Print("Copyright (c) 2020 Sergey Juriev (sjuriev@gmail.com)\n")
	log.Printf("FS watcher, ver. %s\n", Version)
	log.Printf("OS: %s\n", runtime.GOOS)

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error reading config '%s': %s", configFile, err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Watcher error: %s", err)
	}

	for path := range config.Watch {
		err := watcher.Add(path)
		if err != nil {
			log.Fatalf("Add path %q error: %s", path, err)
		}
		log.Printf("Add path '%s' to watcher", path)
	}
	var (
		wg      sync.WaitGroup
		jobs    = make(chan Job)
		opcache = make(map[string]time.Time)
	)

	for wid := 0; wid < config.Workers; wid++ {
		wg.Add(1)
		go worker(wid, &wg, jobs)
	}
	log.Printf("Total workers started: %d", config.Workers)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case err := <-watcher.Errors:
				log.Printf("Watch error: %s", err)
			case ev := <-watcher.Events:
				var params *PathParams
				log.Printf("Event: %s", ev.String())

				cleanPath := filepath.Clean(ev.Name)
				if p, ok := config.Watch[cleanPath]; ok {
					params = p
					log.Printf("Found full path params for '%s'", cleanPath)
				} else {
					dir := filepath.Dir(cleanPath)
					if p, ok := config.Watch[dir]; ok {
						params = p
						log.Printf("Found dir params for '%s'", dir)
					}
				}
				if params == nil || len(params.Cmd) == 0 || (params.Op != 0 && params.Op&ev.Op == 0) {
					log.Printf("Skip event: no params or empty command or wrong mode for event")
					continue
				}
				if ts, ok := opcache[cleanPath]; ok {
					elapsed := time.Since(ts)
					if elapsed <= params.Throttle {
						log.Printf("Skip event '%s' (throttling), elapsed: %d ms", ev.String(), elapsed.Milliseconds())
						continue
					}
				}
				opcache[cleanPath] = time.Now()
				jobs <- Job{timestamp: time.Now(), params: params}
			}
		}
	}()
	log.Println("Event dispatcher started")
	log.Println("Wait for events...")

	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup, jobs <-chan Job) {
	defer wg.Done()

	for job := range jobs {
		params := job.params
		if params.Pause > 0 {
			log.Printf("[%d] Sleep for: %d ms", id, params.Pause.Milliseconds())
			time.Sleep(params.Pause)
		}
		log.Printf("[%d] Try to run: '%s'", id, params.Cmd)
		cmd := exec.Command(params.Cmd[0], params.Cmd[1:]...)
		if params.WorkDir != "" {
			cmd.Dir = params.WorkDir
		}
		if params.Env != nil {
			cmd.Env = params.Env
		}
		if params.LogOutput {
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("[%d] Command run error: %s", id, err)
			} else {
				log.Printf("[%d] Command output: %s", id, output)
			}
		} else {
			err := cmd.Run()
			if err != nil {
				log.Printf("[%d] Command run error: %s", id, err)
			} else {
				log.Printf("[%d] Command '%s' OK", id, params.Cmd)
			}
		}
	}
}
