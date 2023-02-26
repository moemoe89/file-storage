package rpcseqdiagram

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (r *rpcSeqDiagram) buildReadmeSeqDiagram() (string, error) {
	// readmeSequenceDiagram is a readme content for sequence diagram.
	readmeSequenceDiagram := "\n"

	i := 1

	// read all files in rpc sequence diagram.
	err := filepath.Walk(r.targetPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to read each file: %w", err)
			}

			if !strings.Contains(path, "md") {
				return nil
			}

			file, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			rpcName := ""

			scanner := bufio.NewScanner(bytes.NewReader(file))

			// reads content line by line.
			for scanner.Scan() {
				text := scanner.Text()

				if strings.Contains(text, "###") {
					rpcName = strings.ReplaceAll(text, "### ", "")
				}
			}

			if strings.Contains(path, "md") {
				// creates the list of sequence diagram on the readme.
				str := "%d. [%s](%s)\n"
				readmeSequenceDiagram += fmt.Sprintf(str, i, rpcName, path)

				i++
			}

			return nil
		})
	if err != nil {
		return "", fmt.Errorf("failed to read sequence diagram directory: %w", err)
	}

	return readmeSequenceDiagram, nil
}

// generateNewReadme generates the content of the new readme (including the old readme + sequence diagram readme).
func (r *rpcSeqDiagram) generateNewReadme(readmeSequenceDiagram string) (string, error) {
	readmeFile, err := os.Open(r.readmePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer readmeFile.Close()

	scanner := bufio.NewScanner(readmeFile)

	skipLine := false
	newReadme := ""

	for scanner.Scan() {
		if skipLine {
			// continue reads the content if end sequence diagram label found.
			if scanner.Text() == r.endRPCSequenceDiagramDoc {
				skipLine = false

				newReadme += scanner.Text() + "\n"
			}

			continue
		}

		newReadme += scanner.Text()

		// start skips the content if start sequence diagram label found
		// and put the readme sequence diagram on the content.
		if scanner.Text() == r.startRPCSequenceDiagramDoc {
			skipLine = true

			newReadme += readmeSequenceDiagram
		}

		newReadme += "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read scanner: %w", err)
	}

	return newReadme, nil
}

// writeNewReadme writes the new readme file.
func writeNewReadme(newReadme, path string) error {
	newReadmeFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer newReadmeFile.Close()

	_, err = newReadmeFile.WriteString(newReadme)

	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
