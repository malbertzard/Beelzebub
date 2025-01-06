package cli

import (
	"fmt"
	"github.com/malbertzard/beelzebub/pkg/logging"
	"github.com/malbertzard/beelzebub/pkg/session"
	"github.com/spf13/cobra"
)

var manager = session.NewManager()

func Execute() error {
	// Load sessions from the database
	if err := manager.LoadAllSessions(); err != nil {
		logging.Errorf("Failed to load sessions: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "daemon-wrapper",
		Short: "A CLI for managing long-running processes",
	}

	rootCmd.AddCommand(
		startCommand(),
		stopCommand(),
		listCommand(),
		reattachCommand(),
	)

	return rootCmd.Execute()
}

func startCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start [name] [command]",
		Short: "Start a new process",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			command := args[1]
			session := session.NewSession(name, command, args[2:]...)
			if err := session.Start(); err != nil {
				fmt.Printf("Failed to start session %s: %v\n", name, err)
				return
			}
			manager.AddSession(name, session)
			fmt.Printf("Session %s started successfully\n", name)
		},
	}
}

func stopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [name]",
		Short: "Stop a running process",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			session, exists := manager.GetSession(name)
			if !exists {
				fmt.Printf("Session %s not found\n", name)
				return
			}
			if err := session.Stop(); err != nil {
				fmt.Printf("Failed to stop session %s: %v\n", name, err)
				return
			}
			manager.RemoveSession(name)
			fmt.Printf("Session %s stopped successfully\n", name)
		},
	}
}

func listCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all processes",
		Run: func(cmd *cobra.Command, args []string) {
			manager.LoadAllSessions()
			sessions := manager.ListSessions()
			if len(sessions) == 0 {
				fmt.Println("No active sessions")
				return
			}
			fmt.Println("Active sessions:")
			for _, session := range sessions {
				status := "stopped"
				if session.IsRunning() {
					status = "running"
				}
				fmt.Printf("- %s: %s\n", session.Name, status)
			}
		},
	}
}

func reattachCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reattach [name]",
		Short: "Reattach to a running process",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			session, exists := manager.GetSession(name)
			if !exists {
				fmt.Printf("Session %s not found\n", name)
				return
			}
			if !session.IsRunning() {
				fmt.Printf("Session %s is not running\n", name)
				return
			}
			fmt.Printf("Reattaching to session %s...\n", name)
			err := session.Command.Wait()
			if err != nil {
				fmt.Printf("Error while reattaching to session %s: %v\n", name, err)
			}
		},
	}
}
