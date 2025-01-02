package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BeauRussell/TournamentAssistant/components"
	"github.com/BeauRussell/TournamentAssistant/startgg"
)

func pullEventStandings(startClient *startgg.Start, event components.Option) {
	standings := startClient.GetEventStandings(event.ID)

	// Ensure the base directory exists
	baseDir := filepath.Join("event", "standings")
	err := os.MkdirAll(baseDir, os.ModePerm) // os.ModePerm is cross-platform default
	if err != nil {
		log.Printf("Failed to create directories: %v\n", err)
		return
	}

	// Write the event name to a file
	eventFile := filepath.Join("event", "name.txt")
	err = overwriteFile(eventFile, standings.Name)
	if err != nil {
		log.Printf("Failed to write event name: %v\n", err)
		return
	}

	// Write standings to individual files
	for _, standing := range standings.Standings.Nodes {
		filePath := filepath.Join(baseDir, fmt.Sprintf("%d.txt", standing.Placement))
		err = overwriteFile(filePath, standing.Entrant.Name)
		if err != nil {
			log.Printf("Failed to write standing file '%s': %v\n", filePath, err)
		}
	}
}

func overwriteFile(filePath, content string) error {
	// Create or truncate the file for writing
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", filePath, err)
	}

	return nil
}
