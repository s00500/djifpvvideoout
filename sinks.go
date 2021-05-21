package main

import (
	"bufio"
	"io"
	"os"
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
	Args []string
}

func (sink GstSink) StartInstance() (io.WriteCloser, func()) {
	cmd := exec.Command("gst-launch-1.0", sink.Args...)
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

type FifoSink struct {
	Path string
}

func (sink FifoSink) StartInstance() (io.WriteCloser, func()) {
	os.Remove(sink.Path)
	err := syscall.Mkfifo(sink.Path, 0666)
	if err != nil {
		log.Fatal("Make named pipe file error ", sink.Path, " ", err)
	}
	f, err := os.OpenFile(sink.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	log.MustFatal(err)

	return f, func() {
		f.Close()
		os.Remove(sink.Path)
	}
}

type HelloVideoSink struct {
}

func (sink HelloVideoSink) StartInstance() (io.WriteCloser, func()) {
	logger := log.GetLoggerForPrefix("hellosink")

	cmd := exec.Command("/etc/hello_video.bin")
	cmd.Env = append(cmd.Env, "LD_LIBRARY_PATH=/opt/vc/lib")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Fatal("Could not get hellovideo stdin")
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Fatal("Could not get hellovideo stderr")
	}

	go func() {
		errScanner := bufio.NewScanner(stderr)

		for errScanner.Scan() {
			logger.Error(errScanner.Text())
		}
	}()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Fatal("Could not get hellovideo stdout")
	}
	go func() {
		outScanner := bufio.NewScanner(stdout)
		for outScanner.Scan() {
			logger.Warn(outScanner.Text())
		}
	}()

	log.MustFatal(cmd.Start())

	return stdin, func() {
		cmd.Process.Signal(syscall.SIGKILL) // Not ellegant... could try sigterm and wait before...
	}
}
