package session

import (
	"os"
	"os/exec"
	"github.com/creack/pty"
)

type PTYSession struct {
	Session
	PTY *os.File
}

func NewPTYSession(name string, command string, args ...string) (*PTYSession, error) {
	cmd := exec.Command(command, args...)
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return &PTYSession{
		Session: Session{
			Name:    name,
			Command: cmd,
		},
		PTY: ptyFile,
	}, nil
}

func (p *PTYSession) Read(output []byte) (int, error) {
	return p.PTY.Read(output)
}

func (p *PTYSession) Write(input []byte) (int, error) {
	return p.PTY.Write(input)
}

func (p *PTYSession) Close() error {
	err := p.PTY.Close()
	if err != nil {
		return err
	}
	return p.Session.Stop()
}
