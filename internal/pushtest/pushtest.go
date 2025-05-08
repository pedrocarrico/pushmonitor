package pushtest

import (
	"context"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/pedrocarrico/pushmonitor/internal/logger"
)

type PushTest struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"`
	Retries  int    `yaml:"retries"`
	Command  string `yaml:"command"`
}

func (p *PushTest) shouldRun() bool {
	if p.Command == "" {
		return true
	}

	cmd := exec.Command("sh", "-c", p.Command)

	var outputBuf strings.Builder
	cmd.Stdout = &outputBuf
	cmd.Stderr = &outputBuf

	err := cmd.Run()
	output := outputBuf.String()

	if err != nil {
		logger.Warn("Command for test %s failed: %v\nOutput:\n%s", p.Name, err, output)
		return false
	}

	logger.Debug("Command for test %s succeeded. Output:\n%s", p.Name, output)
	return true
}

func (p *PushTest) executeRequest(client *http.Client) bool {
	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		logger.Error("Error creating request for test %s: %v", p.Name, err)
		return false
	}

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	logger.Debug("Test %s: Request URL: %s", p.Name, req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error executing push test %s: %v", p.Name, err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response for test %s: %v", p.Name, err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		logger.Info("Successfully executed push test %s", p.Name)
		return true
	}

	logger.Warn("Push test %s failed with status %d: %s", p.Name, resp.StatusCode, string(body))
	return false
}

func (p *PushTest) Run(client *http.Client, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	logger.Info("Starting push test: %s (interval: %d seconds, retries: %d)", p.Name, p.Interval, p.Retries)
	ticker := time.NewTicker(time.Duration(p.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Debug("Test %s received shutdown signal", p.Name)
			return
		case <-ticker.C:
			if !p.shouldRun() {
				logger.Debug("Skipping push test %s due to failed command", p.Name)
				continue
			}

			logger.Debug("Executing push test: %s", p.Name)
			for i := 0; i < p.Retries; i++ {
				logger.Debug("Test %s: Attempt %d/%d", p.Name, i+1, p.Retries)
				if p.executeRequest(client) {
					break
				}
			}
			logger.Debug("Completed push test cycle for: %s", p.Name)
		}
	}
}
