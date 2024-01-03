// SPDX-FileCopyrightText: Stefan Tatschner
//
// SPDX-License-Identifier: MIT

package penrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/klauspost/compress/zstd"
)

type TestRun struct {
	*Config

	artifactsDirNormalized string
	createLatestLink       bool
}

func NewTestRun(config *Config) *TestRun {
	return &TestRun{Config: config}
}

func (tr *TestRun) createEnvironment() []string {
	// TODO: sudo cleans env variables; reuse this filer maybe?
	env := os.Environ()
	env = append(env, fmt.Sprintf("PENRUN_CLI_STRING=%s", strings.Join(tr.CLI, " ")))
	env = append(env, fmt.Sprintf("PENRUN_ARTIFACTS_DIR=%s", tr.artifactsDirNormalized))
	env = append(env, fmt.Sprintf("PENRUN_INVOCATION=%s", strings.Join(os.Args, " ")))
	return env
}

func (tr *TestRun) prepareRun() error {
	artifactsDir := ""
	if tr.ArtifactsDir != "" {
		artifactsDir = tr.ArtifactsDir
		tr.createLatestLink = false
	} else {
		if tr.ArtifactsBase != "" {
			artifactsDir = tr.ArtifactsBase
		} else {
			artifactsDir = "penrun-artifacts"
		}

		artifactsDir = filepath.Join(artifactsDir, tr.CLI[0], fmt.Sprintf("run-%s", time.Now().Format(time.RFC3339Nano)))
		tr.createLatestLink = true
	}

	var err error
	artifactsDir, err = filepath.Abs(artifactsDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(artifactsDir, 0755); err != nil {
		return err
	}
	tr.artifactsDirNormalized = artifactsDir

	return nil
}

func (tr *TestRun) runHook(cli string) error {
	cmd := exec.Command("/bin/sh", "-c", cli)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (tr *TestRun) runPreHook() error {
	if tr.PreHook != "" {
		return tr.runHook(tr.PreHook)
	}
	return nil
}

func (tr *TestRun) runPostHook() error {
	if tr.PostHook != "" {
		return tr.runHook(tr.PostHook)
	}
	return nil
}

func (tr *TestRun) run() (*exec.Cmd, error) {
	cmd := exec.Command(tr.CLI[0], tr.CLI[1:]...)

	outfile, err := os.Create(filepath.Join(tr.artifactsDirNormalized, "OUTPUT.zst"))
	if err != nil {
		return nil, err
	}

	zstdEncoder, err := zstd.NewWriter(outfile)
	if err != nil {
		return nil, err
	}

	cmd.Env = tr.createEnvironment()

	rp, wp := io.Pipe()
	cmd.Stderr = wp
	cmd.Stdout = wp
	cmd.Stdin = os.Stdin

	teeReader := io.TeeReader(rp, zstdEncoder)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Copy stdout/stderr into the teeReader, aka. the terminal and OUTPUT.zst.
	errCh := make(chan error)
	go func() {
		_, err := io.Copy(os.Stderr, teeReader)
		errCh <- err
	}()

	// Signal Handling
	go func() {
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		recvSig := <-sigCh

		// Forward the catched signal to the child process.
		cmd.Process.Signal(recvSig)
	}()

	// Ignore errors here, since failed runs need to be
	// tracked as well.
	// TODO: are there errors that need to be handled?
	cmd.Wait()

	if err := wp.Close(); err != nil {
		return nil, err
	}

	err = <-errCh
	if err != nil {
		return nil, err
	}

	if err := zstdEncoder.Close(); err != nil {
		return nil, err
	}
	if err := outfile.Close(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (tr *TestRun) getFinishedRun(cmd *exec.Cmd, startTime time.Time) *FinishedRun {
	buildInfoStr := "unknown"
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		buildInfoStr = buildInfo.String()
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}

	endTime := time.Now()

	exitCode := cmd.ProcessState.ExitCode()

	// If the process was terminated by a signal, then extract the real exit code.
	// Typically it is 128 + SIGNAL_NUMBER.
	if exitCode == -1 {
		waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = int(128 + waitStatus.Signal())
	}

	return &FinishedRun{
		CLI:           tr.CLI,
		CLIString:     strings.Join(tr.CLI, " "),
		ArtifactsDir:  tr.artifactsDirNormalized,
		ExitCode:      exitCode,
		StartTime:     startTime,
		EndTime:       endTime,
		Duration:      endTime.Sub(startTime),
		Environ:       tr.createEnvironment(),
		Hostname:      hostname,
		CWD:           cwd,
		PenrunVersion: buildInfoStr,
	}
}

func (tr *TestRun) Run() (*FinishedRun, error) {
	if err := tr.prepareRun(); err != nil {
		return nil, err
	}

	if tr.PreHook != "" {
		if err := tr.runPreHook(); err != nil {
			slog.Warn(fmt.Sprintf("pre hook failed: %s", err))
		}
	}

	startTime := time.Now()

	cmd, err := tr.run()
	if err != nil {
		return nil, err
	}

	signal.Reset()

	if tr.PostHook != "" {
		if err := tr.runPostHook(); err != nil {
			slog.Warn(fmt.Sprintf("post hook failed: %s", err))
		}
	}

	return tr.getFinishedRun(cmd, startTime), nil
}

type FinishedRun struct {
	ArtifactsDir  string        `json:"artifacts_dir"`
	CLI           []string      `json:"cli"`
	CLIString     string        `json:"cli_string"`
	CWD           string        `json:"cwd"`
	ExitCode      int           `json:"exit_code"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	Duration      time.Duration `json:"duration"`
	Environ       []string      `json:"environ"`
	Hostname      string        `json:"hostname"`
	PenrunVersion string        `json:"penrun_version"`
}

func RunCommand(config *Config) (*FinishedRun, error) {
	tr := NewTestRun(config)
	fr, err := tr.Run()
	if err != nil {
		return nil, err
	}

	if err := fr.createArtifacts(); err != nil {
		return nil, err
	}

	return fr, nil
}

func (fr *FinishedRun) createMetaFile() error {
	data, err := json.Marshal(fr)
	if err != nil {
		return err
	}

	var (
		metaPath   = filepath.Join(fr.ArtifactsDir, "META.json.zst")
		encoder, _ = zstd.NewWriter(nil)
		outData    = make([]byte, 0, len(data))
	)
	outData = encoder.EncodeAll(append(data, '\n'), outData)

	if err := os.WriteFile(metaPath, outData, 0644); err != nil {
		return nil
	}
	return nil
}

func (fr *FinishedRun) createLatestSymlink() error {
	latestLinkPath := filepath.Join(filepath.Dir(fr.ArtifactsDir), "LATEST")
	if _, err := os.Stat(latestLinkPath); !errors.Is(err, os.ErrNotExist) {
		if err := os.Remove(latestLinkPath); err != nil {
			return err
		}
	}

	// The LATEST symlink must have a relative target, such that the
	// directory hierarchy can be moved around without breaking the link.
	if err := os.Symlink(filepath.Base(fr.ArtifactsDir), latestLinkPath); err != nil {
		return err
	}
	return nil
}

func (fr *FinishedRun) createArtifacts() error {
	if err := fr.createMetaFile(); err != nil {
		return err
	}
	if err := fr.createLatestSymlink(); err != nil {
		return err
	}
	return nil
}

func (fr *FinishedRun) Success() bool {
	return fr.ExitCode == 0
}
