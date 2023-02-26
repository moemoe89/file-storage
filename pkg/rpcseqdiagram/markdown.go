package rpcseqdiagram

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func (r *rpcSeqDiagram) writeMD(mapRPC map[string]struct{}) error {
	// iterates handlers data and generate the sequence diagram for each handler from parameter.
	for i, handler := range r.handlers {
		if _, ok := mapRPC[handler.rpc]; !ok {
			continue
		}

		participantCodeDuplicate := make(map[string]struct{}, 0)

		md := "### " + handler.rpc + " RPC - Sequence Diagram\n\n"
		md += "```mermaid\n"
		md += "sequenceDiagram\n"
		md += "\tautonumber\n"
		md += "\tparticipant RPC as " + handler.rpc + " RPC\n"

		for _, usecase := range handler.usecases {
			num := ""
			if len(handler.usecases) > 1 {
				num = strconv.Itoa(i + 1)
			}

			md += "\tparticipant UC" + num + " as " + usecase + " UC\n"

			md = r.iteratesChildParticipant(md, usecase, participantCodeDuplicate)
		}

		md += "\n"

		for _, usecase := range handler.usecases {
			num := ""
			if len(handler.usecases) > 1 {
				num = strconv.Itoa(i + 1)
			}

			md += "\tRPC->>+UC" + num + ": Call\n"

			md += r.iteratesChildUsecase(num, r.mapUsecaseSequence[usecase], usecase)

			md += "\tUC" + num + "-->>-RPC: return\n"
		}

		md += "```\n"

		// remove loop end, alt end, else end without object.
		md = fixIncorrectMdFormat(md)

		err := writeSequenceDiagramMd(md, r.targetPath+"/"+handler.filename)
		if err != nil {
			return fmt.Errorf("failed to write sequence diagram md: %w", err)
		}

		log.Println(">>> " + handler.rpc + " RPC sequence diagram generated!!")
	}

	return nil
}

// after mermaid js content generated, sometimes the format is incorrect
// need to remove it to make the output can be generated
// e.g some loop doesn't have any usecase
// loop xx
// end
// or if else condition doesn't have any usecase
// alt
// else
// end.
func fixIncorrectMdFormat(md string) string {
	// scan the whole string
	scanner := bufio.NewScanner(bytes.NewReader([]byte(md)))

	mdArr := make([]string, 0)

	// put each of line in array string
	for scanner.Scan() {
		mdArr = append(mdArr, scanner.Text())
	}

	i := 0

	// iterates the array.
	for ; i < len(mdArr); i++ {
		// skip if the sequence code not contain end.
		if space.ReplaceAllString(mdArr[i], "") != "end" {
			continue
		}

		// checks if the previous code contains loop / alt
		// or the previous code contains else / default.
		if strings.Contains(mdArr[i-1], "loop") || strings.Contains(mdArr[i-1], "alt") {
			// if contains loop / alt then delete include the end
			// turn back the index in 3 position.
			mdArr = append(mdArr[:i-1], mdArr[i+1:]...)
			i -= 3
		} else if strings.Contains(mdArr[i-1], "else") || strings.Contains(mdArr[i-1], "default") {
			// if contains else / default then delete include the end
			// turn back the index in 2 position (still leaving the end code).
			mdArr = append(mdArr[:i-1], mdArr[i:]...)
			i -= 2
		}
	}

	// add new line.
	mdArr = append(mdArr, "\n")
	// return in string with new line in each array.
	return strings.Join(mdArr, "\n")
}

func createMDFilename(rpc string) []byte {
	filename := make([]byte, 0)

	for i := range rpc {
		if i > 0 && isUpper(rpc[i]) {
			filename = append(filename, '-')
		}

		filename = append(filename, toLower(rpc[i]))
	}

	filename = append(filename, []byte(mdExtension)...)

	return filename
}

// writeSequenceDiagramMd writes the sequence diagram md file.
func writeSequenceDiagramMd(md, path string) error {
	newReadmeFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer newReadmeFile.Close()

	_, err = newReadmeFile.WriteString(md)

	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
