package rpcseqdiagram

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (r *rpcSeqDiagram) getMapRPC(rpc string) (map[string]struct{}, error) {
	if space.ReplaceAllString(rpc, "") == "" {
		return nil, errEmptyRPC
	}

	rpcs := strings.Split(rpc, ",")

	mapRPC := make(map[string]struct{})

	if rpc == "*" {
		for _, handler := range r.handlers {
			mapRPC[handler.rpc] = struct{}{}
		}
	} else {
		for i := range rpcs {
			foundRPC := false

			for _, handler := range r.handlers {
				if handler.rpc == rpcs[i] {
					foundRPC = true
					mapRPC[rpcs[i]] = struct{}{}

					break
				}
			}

			if !foundRPC {
				return nil, fmt.Errorf("something error, rpc: %s not found", rpcs[i])
			}
		}
	}

	return mapRPC, nil
}

// findContentHandlersStructHandlerUCVariable will get the every content handler as array string
// and name of handler struct also the name of usecase variable.
func (r *rpcSeqDiagram) findContentHandlersStructHandlerUCVariable() error {
	var structHandler, ucVariable string

	var findStructHandler, findUCVariable bool

	contentHandlers := make([]string, 0)

	// read all files in handler path.
	err := filepath.Walk(r.grpchandlerPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to read each file: %w", err)
			}

			// skip if the file is not Go files
			if !strings.Contains(path, goExtension) {
				return nil
			}

			file, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			// put every content handlers in array.
			contentHandlers = append(contentHandlers, string(file))

			r.findStructHandlerUCVar(file, &findStructHandler, &findUCVariable, &structHandler, &ucVariable)

			return nil
		})
	if err != nil {
		return fmt.Errorf("failed to read all files: %w", err)
	}

	if structHandler == "" {
		return errEmptyHandlerStruct
	}

	r.contentHandlers = contentHandlers
	r.structHandler = structHandler
	r.ucVariable = ucVariable

	return nil
}

func (r *rpcSeqDiagram) findStructHandlerUCVar(
	file []byte,
	findStructHandler, findUCVariable *bool,
	structHandler, ucVariable *string,
) {
	scanner := bufio.NewScanner(bytes.NewReader(file))

	// reads content line by line.
	for scanner.Scan() {
		text := scanner.Text()

		// if struct for handler found, find the name of struct.
		if *findStructHandler {
			*structHandler = strings.ReplaceAll(text, typePrefix, "")
			*structHandler = strings.ReplaceAll(*structHandler, structPrefix, "")
			*findStructHandler = false
		}

		// if line contains annotation for handler struct, set some flags to true.
		if strings.Contains(text, r.isAStructHandler) {
			*findStructHandler = true
			*findUCVariable = true
		}

		// if line contains usecase prefix, find the uc variable name.
		if *findUCVariable && strings.Contains(text, usecasePrefix) {
			split := strings.Split(text, usecasePrefix)
			*ucVariable = space.ReplaceAllString(split[0], "")
			*findUCVariable = false
		}
	}
}

// findHandlers get the data of handler includes rpc name, receiver name, usecases in rpc, filename for md file.
func (r *rpcSeqDiagram) findHandlers() {
	handlers := make([]handlerModel, 0)

	findHandler := false
	lastIndex := -1

	// iterates all content handlers array.
	for _, c := range r.contentHandlers {
		scanner := bufio.NewScanner(bytes.NewReader([]byte(c)))

		// read line by line.
		for scanner.Scan() {
			text := scanner.Text()

			// if line contains methods func(x *X) X() X {.
			if strings.Contains(text, funcOpenParenthesesPrefix) && strings.Contains(text, r.structHandler+closeParenthesesSpace) {
				r.addToHandler(&handlers, text, &findHandler, &lastIndex)
			}

			// if the handler found, find the close curly bracket } of the handler without leading space.
			if findHandler && text == closeCurlyBracket {
				findHandler = false
			}

			// skip if handler not found
			if !findHandler {
				continue
			}

			// if strings contains access to usecases,
			// put the usecase inside the handler usecases array.
			if strings.Contains(text, handlers[lastIndex].receiver+dot+r.ucVariable+dot) {
				split := strings.Split(text, handlers[lastIndex].receiver+dot+r.ucVariable+dot)
				split = strings.Split(split[1], openParentheses)
				usecase := split[0]

				handlers[lastIndex].usecases = append(handlers[lastIndex].usecases, usecase)
			}
		}
	}

	r.handlers = handlers
}

func (r *rpcSeqDiagram) addToHandler(
	handlers *[]handlerModel,
	text string,
	findHandler *bool,
	lastIndex *int,
) {
	split := strings.Split(text, r.structHandler+closeParenthesesSpace)

	// get the receiver name func(x *X) -> x.
	receiver := strings.ReplaceAll(split[0], funcOpenParenthesesPrefix, "")
	receiver = strings.ReplaceAll(receiver, spacePointer, "")

	// get the RPC name func(x *X) X() -> X.
	split = strings.Split(split[1], openParentheses)
	rpc := split[0]

	// create a md filename based on RPC name e.g. GetXxx -> get-xxx.md.
	filename := createMDFilename(rpc)

	// RPC name always start with capital, if it does, put on handlers array.
	if isUpper(rpc[0]) {
		*findHandler = true
		*lastIndex++

		*handlers = append(*handlers, handlerModel{
			rpc:      rpc,
			receiver: receiver,
			filename: string(filename),
		})
	}
}
