package rpcseqdiagram

// handlerModel is a data model for handler.
type handlerModel struct {
	rpc      string
	receiver string
	filename string
	usecases []string
}

// usecaseModel is a data model for usecase.
type usecaseModel struct {
	callingName string
	callingFunc string
	callingType string
}

// participantModel is a data model for Mermaid JS participant.
type participantModel struct {
	Code string
	Name string
}

// operatorModel is a data model for operator like if, switch, for.
type operatorModel struct {
	operator     string
	leadingSpace int
	isSKip       bool
}

// rpcSeqDiagram is a struct for RPC Sequence Diagram generator.
type rpcSeqDiagram struct {
	// targetPath is the directory for the generated diagram.
	targetPath string
	// readmePath is the readme file for including the diagram content.
	readmePath string
	// grpchandlerPath is the grpc handler path for the project.
	grpchandlerPath string
	// usecasePath is the usecase path for the project.
	usecasePath string
	// isAStructHandler is a comment text for pointing to struct handler.
	isAStructHandler string
	// isAStructUsecase is a comment text for pointing to struct usecase.
	isAStructUsecase string
	// startRPCSequenceDiagramDoc is the start of rpc sequence diagram content.
	startRPCSequenceDiagramDoc string
	// endRPCSequenceDiagramDoc is the end of rpc sequence diagram content.
	endRPCSequenceDiagramDoc string
	// commentForMermaidJS is the comment indicates for mermaid JS syntax.
	commentForMermaidJS string
	// mermaidJSReplace is the code for replacing mermaid JS syntax.
	mermaidJSReplace string
	// contentHandlers is a content for each handler.
	contentHandlers []string
	// structHandler is a struct name for handler.
	structHandler string
	// ucVariable is a name for ucasecase (uc) variable.
	ucVariable string
	// handlers is an array of handlerModel.
	handlers []handlerModel
	// contentUsecases is a content for each usecase.
	contentUsecases []string
	// structUsecase is a struct name for usecase.
	structUsecase string
	// mapParticipant is map data with participantModel.
	mapParticipant map[string]participantModel
	// mapUsecaseMethod is map string key data with array usecaseModel.
	mapUsecaseMethod map[string][]usecaseModel
	// mapUsecaseFunction is map string key data with array usecaseModel.
	mapUsecaseFunction map[string][]usecaseModel
	// mapUsecaseSequence is map string key data with string.
	mapUsecaseSequence map[string]string
	// mapLeadingSpaceEnd for counting the leading space.
	// before meet some character listed in this map.
	mapEndLeadingSpace map[rune]bool
	// mapTrimConditions is map string key data with string.
	mapTrimConditions map[string]string
	// mapParticipantLibs is map string key data with bool.
	mapParticipantLibs map[string]bool
}
