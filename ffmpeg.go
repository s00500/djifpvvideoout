package main

import (
	"io"
	"os/exec"
	"syscall"

	log "github.com/s00500/env_logger"
)

// TODO: implement ways of ovveriding args
// TODO: Implement gst as well
func StartFFMPEGInstance() (io.WriteCloser, func()) {
	additionalArgs := []string{"-fast", "-flags2", "fast", "-fflags", "nobuffer", "-flags", "low_delay", "-strict", "experimental", "-vf", "setpts=N/60/TB", "-framedrop", "-sync", "ext", "-probesize", "32", "-analyzeduration", "0"}
	args := []string{"-i", "-"}
	args = append(args, additionalArgs...)
	cmd := exec.Command("ffplay", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("Could not get ffmpeg stdin")
	}
	log.MustFatal(cmd.Start())
	return stdin, func() {
		cmd.Process.Signal(syscall.SIGKILL) // Not ellegant... could try sigterm and wait before...
	}
}
