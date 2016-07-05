package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	//	. "archiver/utils"
)

//
type Cmd struct {
	Command string
	Argv    string
	LogPath string
	Stdout  io.Reader
	Stderr  io.Reader
	Cmd     *exec.Cmd
}

// start system cli
func (c *Cmd) Start() error {
	cli, err := exec.LookPath(c.Command)
	if err != nil {
		ErrorLog(fmt.Sprintf("cli dont existed: %s\n", err))
		return err
	}
	cli = fmt.Sprintf("%s %s", cli, c.Argv)
	ErrorLog(fmt.Sprintf("%s\n", cli))
	c.Cmd = exec.Command("/bin/sh", "-c", cli)

	c.Stdout, err = c.Cmd.StdoutPipe()
	c.Stderr, err = c.Cmd.StderrPipe()

	err = c.Cmd.Start()
	if err != nil {
		ErrorLog(fmt.Sprintf("Cmd Run err:\n\t%s\n", err))
		return err
	}
	ErrorLog(fmt.Sprintf("-->Start command... Pid: %d\n", c.Cmd.Process.Pid))

	return nil
}

// wait sys command running
func (c *Cmd) Wait() (*os.ProcessState, error) {
	var p *os.ProcessState

	err := c.Cmd.Wait()
	if err != nil {
		return p, err
	}
	ErrorLog(fmt.Sprintf("-->Waiting for command to finish... Pid: %d\n", c.Cmd.Process.Pid))

	p = c.Cmd.ProcessState

	return p, nil
}

// get cli log
func (c *Cmd) WriteLog() ([]byte, error) {
	stdout, err := ioutil.ReadAll(c.Stdout)
	if err != nil {
		ErrorLog(err)
		return stdout, err
	}
	stderr, err := ioutil.ReadAll(c.Stderr)
	if err != nil {
		ErrorLog(err)
		return stderr, err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.log", c.LogPath), stdout, 0664)
	if err != nil {
		ErrorLog(err)
		return stdout, err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s.err", c.LogPath), stderr, 0664)
	if err != nil {
		ErrorLog(err)
		return stderr, err
	}

	return stderr, nil
}

// get cli log
func (c *Cmd) RemoveLog() error {
	err := os.Remove(fmt.Sprintf("%s.log", c.LogPath))
	if err != nil {
		ErrorLog(err)
		return err
	}
	err = os.Remove(fmt.Sprintf("%s.err", c.LogPath))
	if err != nil {
		ErrorLog(err)
		return err
	}

	return nil
}

// get cli log
func (c *Cmd) WriteRunLog() error {
	stdout, err := ioutil.ReadAll(c.Stdout)
	if err != nil {
		ErrorLog(err)
		return err
	}
	stderr, err := ioutil.ReadAll(c.Stderr)
	if err != nil {
		ErrorLog(err)
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.log", c.LogPath), stdout, 0664)
	if err != nil {
		ErrorLog(err)
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s.err", c.LogPath), stderr, 0664)
	if err != nil {
		ErrorLog(err)
		return err
	}

	return nil
}

// get cli log
func (c *Cmd) ReadRunLog() ([]byte, error) {
	var stdout, stderr []byte

	outfile, err := os.Open(fmt.Sprintf("%s.log", c.LogPath))
	if err != nil {
		ErrorLog(err)
		return stdout, err
	}
	defer outfile.Close()
	stdout, err = ioutil.ReadAll(outfile)
	if err != nil {
		ErrorLog(err)
		return stdout, err
	}

	errfile, err := os.Open(fmt.Sprintf("%s.err", c.LogPath))
	if err != nil {
		ErrorLog(err)
		return stderr, err
	}
	defer errfile.Close()
	stderr, err = ioutil.ReadAll(errfile)
	if err != nil {
		ErrorLog(err)
		return stderr, err
	}

	return stderr, nil
}
