package session

import (
	"os/exec"
	"sync"

	"github.com/malbertzard/beelzebub/pkg/db"
)

type Manager struct {
	Sessions map[string]*Session
	Lock     sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Sessions: make(map[string]*Session),
	}
}

func (m *Manager) AddSession(name string, session *Session) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Sessions[name] = session
	session.Save()
}

func (m *Manager) RemoveSession(name string) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	delete(m.Sessions, name)
	m.RemoveSessionFromDB(name)
}

func (m *Manager) GetSession(name string) (*Session, bool) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	session, exists := m.Sessions[name]
	return session, exists
}

func (m *Manager) ListSessions() []*Session {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	var sessions []*Session
	for _, session := range m.Sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

func (m *Manager) LoadAllSessions() error {
	conn := db.GetDB()
	rows, err := conn.Query("SELECT name, command, args FROM sessions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, command, args string
		if err := rows.Scan(&name, &command, &args); err != nil {
			return err
		}
		cmd := exec.Command(command, args)
		m.Sessions[name] = &Session{Name: name, Command: cmd}
	}
	return nil
}

func (m *Manager) RemoveSessionFromDB(name string) error {
	conn := db.GetDB()
	_, err := conn.Exec("DELETE FROM sessions WHERE name = ?", name)
	return err
}
