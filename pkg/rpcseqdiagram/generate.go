package rpcseqdiagram

import (
	"fmt"
	"os"
	"regexp"
)

// regex for eliminating \n, \t and " ".
var space = regexp.MustCompile(`\s+`)

// RPCSeqDiagram is an interface for rpcSeqDiagram.
type RPCSeqDiagram interface {
	Generate(rpc string) error
}

// New returns RPCSeqDiagram implementation and error.
func New(opts ...Option) (RPCSeqDiagram, error) {
	r := new(rpcSeqDiagram)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(r); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return r, nil
}

// Generate will generate the RPC sequence diagram.
func (r *rpcSeqDiagram) Generate(rpc string) error {
	// create rpc sequence diagram directory if not exist.
	if _, err := os.Stat(r.targetPath); os.IsNotExist(err) {
		// creates a diagram directory
		if err := os.MkdirAll(r.targetPath, dirPermission); err != nil {
			return fmt.Errorf("failed to create a dir: %w", err)
		}
	}

	// find content handlers, struct handler name and uc variable name.
	err := r.findContentHandlersStructHandlerUCVariable()
	if err != nil {
		return fmt.Errorf("failed to find content handlers, struct handler and uc variable: %w", err)
	}

	// find handlers RPC.
	r.findHandlers()

	mapRPC, err := r.getMapRPC(rpc)
	if err != nil {
		return fmt.Errorf("failed to get map rpc: %w", err)
	}

	// find content usecases, struct usecase, map participant.
	err = r.findContentUsecasesStructUsecaseMapParticipant()
	if err != nil {
		return fmt.Errorf("failed to find content usecase, struct usecase and map participant: %w", err)
	}

	// find mapUsecaseMethod.
	r.findMapUsecaseMethodFunction()

	// find mapUsecaseSequence.
	r.findMapUsecaseSequence()

	err = r.writeMD(mapRPC)
	if err != nil {
		return fmt.Errorf("failed to write md file: %w", err)
	}

	readmeSequenceDiagram, err := r.buildReadmeSeqDiagram()
	if err != nil {
		return fmt.Errorf("failed to build readme: %w", err)
	}

	newReadme, err := r.generateNewReadme(readmeSequenceDiagram)
	if err != nil {
		return fmt.Errorf("failed to generate new readme: %w", err)
	}

	err = writeNewReadme(newReadme, readmePath)
	if err != nil {
		return fmt.Errorf("failed to write new readme: %w", err)
	}

	return nil
}
