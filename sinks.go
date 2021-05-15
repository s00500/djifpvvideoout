package main

import (
	"io"
	"os/exec"
	"syscall"

	log "github.com/s00500/env_logger"
)

type StreamSink interface {
	StartInstance() (io.WriteCloser, func())
}

type FFPlaySink struct {
}

func (sink FFPlaySink) StartInstance() (io.WriteCloser, func()) {
	additionalArgs := []string{"-fast", "-flags2", "fast", "-fflags", "nobuffer", "-flags", "low_delay", "-strict", "experimental", "-vf", "setpts=N/60/TB", "-framedrop", "-sync", "ext", "-probesize", "32", "-analyzeduration", "0"}
	args := []string{"-i", "-"}
	args = append(args, additionalArgs...)
	cmd := exec.Command("ffplay", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("Could not get ffplay stdin")
	}
	log.MustFatal(cmd.Start())
	return stdin, func() {
		cmd.Process.Signal(syscall.SIGKILL) // Not ellegant... could try sigterm and wait before...
	}
}

type GstSink struct {
}

func (sink GstSink) StartInstance() (io.WriteCloser, func()) {
	args := []string{"fdsrc", "fd=0", "!", "decodebin", "!", "videoconvert", "!", "autovideosink"}
	cmd := exec.Command("gst-launch-1.0", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("Could not get gstreamer stdin")
	}
	log.MustFatal(cmd.Start())
	return stdin, func() {
		cmd.Process.Signal(syscall.SIGKILL) // Not ellegant... could try sigterm and wait before...
	}
}

type UdpSink struct {
}

func (sink UdpSink) StartInstance() (io.WriteCloser, func()) {
	additionalArgs := []string{"-vcodec", "mpeg4", "-v", "0", "-f", "mpegts", "udp://127.0.0.1:23000"}
	args := []string{"-i", "-"}
	args = append(args, additionalArgs...)
	cmd := exec.Command("ffmpeg", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("Could not get ffmpeg stdin")
	}
	log.MustFatal(cmd.Start())
	return stdin, func() {
		cmd.Process.Signal(syscall.SIGKILL) // Not ellegant... could try sigterm and wait before...
	}
}
