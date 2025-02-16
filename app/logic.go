package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BeauRussell/TournamentAssistant/components"
	"github.com/BeauRussell/TournamentAssistant/startgg"
)

func pullEventStandings(startClient *startgg.Start, event components.Option) {
	standings := startClient.GetEventStandings(event.ID)

	baseDir := filepath.Join("./event", "standings")
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create directories: %v\n", err)
		return
	}

	eventFile := filepath.Join("./event", "name.txt")
	err = overwriteFile(eventFile, standings.Name)
	if err != nil {
		log.Printf("Failed to write event name: %v\n", err)
		return
	}

	for _, standing := range standings.Standings.Nodes {
		filePath := filepath.Join(baseDir, fmt.Sprintf("%d.txt", standing.Placement))
		err = overwriteFile(filePath, standing.Entrant.Name)
		if err != nil {
			log.Printf("Failed to write standing file '%s': %v\n", filePath, err)
		}
	}
}

func pullBracketData(startClient *startgg.Start, event components.Option) {
	phaseMatches := startClient.GetPhaseMatches(event.ID)

	for _, phase := range phaseMatches {
		if phase.Name == "Qualifiers" {
			continue
		}
		baseDir := filepath.Join("./event", "matches")
		phaseDir := filepath.Join(baseDir, phase.Name)
		err := os.MkdirAll(phaseDir, os.ModePerm)
		if err != nil {
			log.Printf("Failed to create directories: %v\n", err)
			return
		}

		replacePhaseMatchups(phaseDir, phase.Sets.Nodes)
	}
}

func replacePhaseMatchups(phaseDir string, matches []startgg.MatchNode) {
	for _, match := range matches {
		for index, slot := range match.Slots {
			fileName := match.MatchId + "_" + strconv.Itoa(index) + ".txt"
			overwriteFile(filepath.Join(phaseDir, fileName), slot.Entrant.Name)
		}
	}
}

func overwriteFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", path, err)
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", path, err)
	}

	return nil
}
