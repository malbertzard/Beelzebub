package session

import (
	"database/sql"
	"fmt"
	"os/exec"

	"github.com/malbertzard/beelzebub/pkg/db"
)

type Session struct {
	Name    string
	Command *exec.Cmd
}

func NewSession(name string, command string, args ...string) *Session {
	return &Session{
		Name:    name,
		Command: exec.Command(command, args...),
	}
}

func (s *Session) Start() error {
	return s.Command.Start()
}

func (s *Session) Stop() error {
	return s.Command.Process.Kill()
}

func (s *Session) IsRunning() bool {
	return s.Command.Process != nil && s.Command.ProcessState == nil
}

func (s *Session) Save() error {
	conn := db.GetDB()
	_, err := conn.Exec("INSERT OR REPLACE INTO sessions (name, command, args, status) VALUES (?, ?, ?, ?)",
		s.Name, s.Command.Path, s.Command.Args, "running")
	return err
}

func LoadSession(name string) (*Session, error) {
	conn := db.GetDB()
	row := conn.QueryRow("SELECT command, args FROM sessions WHERE name = ?", name)

	var command string
	var args string
	if err := row.Scan(&command, &args); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session %s not found", name)
		}
		return nil, err
	}

	cmd := exec.Command(command, args)
	return &Session{Name: name, Command: cmd}, nil
}
