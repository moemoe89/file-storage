package rpcseqdiagram

// Option configures "rpcSeqDiagram" structure.
type Option func(*rpcSeqDiagram) error

// defaultMapEndLeadingSpace for counting the leading space
// before meet some character listed in this map.
var defaultMapEndLeadingSpace = map[rune]bool{
	'f': true, // for
	'}': true, // } end for, switch, if
	'i': true, // if
	's': true, // switch
}

// defaultMapTrimConditions for replacing some characters.
var defaultMapTrimConditions = map[string]string{
	" {": "",
	"} ": "",
	";":  ",",
	":":  "",
}

// defaultMapParticipantLibs will find the lib on the code based on key
// for the value, will use as using the field name instead the type
// this is for some lib that have the same type e.g. Factory.
var defaultMapParticipantLibs = map[string]bool{
	serviceClientLib: false,
	repoLib:          false,
	serviceLib:       false,
	factoryLib:       true,
}

// defaultOptions is default option value for the data.
var defaultOptions = []Option{
	WithTargetPath(sequenceDiagramPath),
	WithReadmePath(readmePath),
	WithGrpcHandlerPath(grpchandlerPath),
	WithUsecasePath(usecasePath),
	WithIsAStructHandler(isAStructHandler),
	WithIsAStructUsecase(isAStructUsecase),
	WithStartRPCSequenceDiagramDoc(startRPCSequenceDiagramDoc),
	WithEndRPCSequenceDiagramDoc(endRPCSequenceDiagramDoc),
	WithCommentForMermaidJS(commentForMermaidJS),
	WithMermaidJSReplace(mermaidJSReplace),
	WithEndLeadingSpace(defaultMapEndLeadingSpace),
	WithTrimConditions(defaultMapTrimConditions),
	WithParticipantLibs(defaultMapParticipantLibs),
}

// WithTargetPath returns Option that sets the value to the rpcSeqDiagram.targetPath
// by default the target path will be in `docs/sequence-diagrams/rpc`.
func WithTargetPath(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyTargetPath
		}

		r.targetPath = val

		return nil
	}
}

// WithReadmePath returns Option that sets the value to the rpcSeqDiagram.readmePath
// by default the readme path will be in `README.md` in the root directory.
func WithReadmePath(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyReadmePath
		}

		r.readmePath = val

		return nil
	}
}

// WithGrpcHandlerPath returns Option that sets the value to the rpcSeqDiagram.grpchandlerPath
// by default the GRPC handler path will be in `internal/adapters/grpchandler`.
func WithGrpcHandlerPath(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyGrpchandlerPath
		}

		r.grpchandlerPath = val

		return nil
	}
}

// WithUsecasePath returns Option that sets the value to the rpcSeqDiagram.usecasePath
// by default the usecase path will be in `internal/usecases`.
func WithUsecasePath(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyUsecasePath
		}

		r.usecasePath = val

		return nil
	}
}

// WithIsAStructHandler returns Option that sets the value to the rpcSeqDiagram.isAStructHandler
// by default is a struct handler will be in `is a struct for handler`.
func WithIsAStructHandler(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyIsAStructHandler
		}

		r.isAStructHandler = val

		return nil
	}
}

// WithIsAStructUsecase returns Option that sets the value to the rpcSeqDiagram.isAStructUsecase
// by default is a struct usecase will be in `is a struct for usecase`.
func WithIsAStructUsecase(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyIsAStructUsecase
		}

		r.isAStructUsecase = val

		return nil
	}
}

// WithStartRPCSequenceDiagramDoc returns Option that sets the value to the rpcSeqDiagram.startRPCSequenceDiagramDoc
// by default the start RPC sequence diagram doc will be in `<!-- start rpc sequence diagram doc -->`.
func WithStartRPCSequenceDiagramDoc(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyStartRPCSequenceDiagramDoc
		}

		r.startRPCSequenceDiagramDoc = val

		return nil
	}
}

// WithEndRPCSequenceDiagramDoc returns Option that sets the value to the rpcSeqDiagram.endRPCSequenceDiagramDoc
// by default the end RPC sequence diagram doc will be in `<!-- end rpc sequence diagram doc -->`.
func WithEndRPCSequenceDiagramDoc(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyEndRPCSequenceDiagramDoc
		}

		r.endRPCSequenceDiagramDoc = val

		return nil
	}
}

// WithCommentForMermaidJS returns Option that sets the value to the rpcSeqDiagram.commentForMermaidJS
// by default the comment for mermaid JS will be in `NOTE: this comments for generating mermaid js code`.
func WithCommentForMermaidJS(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyCommentForMermaidJS
		}

		r.commentForMermaidJS = val

		return nil
	}
}

// WithMermaidJSReplace returns Option that sets the value to the rpcSeqDiagram.mermaidJSReplace
// by default the mermaid JS replace will be in `// mermaid js replace: `.
func WithMermaidJSReplace(val string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == "" {
			return errEmptyMermaidJSReplace
		}

		r.mermaidJSReplace = val

		return nil
	}
}

// WithEndLeadingSpace returns Option that sets the value to the rpcSeqDiagram.mapLeadingSpaceEnd
// by default the map value will be use `defaultMapLeadingSpaceEnd`.
func WithEndLeadingSpace(val map[rune]bool) Option {
	return func(r *rpcSeqDiagram) error {
		if val == nil {
			return errEmptyEndLeadingSpace
		}

		r.mapEndLeadingSpace = val

		return nil
	}
}

// WithTrimConditions returns Option that sets the value to the rpcSeqDiagram.mapTrimConditions
// by default the map value will be use `defaultMapTrimConditions`.
func WithTrimConditions(val map[string]string) Option {
	return func(r *rpcSeqDiagram) error {
		if val == nil {
			return errEmptyTrimConditions
		}

		r.mapTrimConditions = val

		return nil
	}
}

// WithParticipantLibs returns Option that sets the value to the rpcSeqDiagram.mapParticipantLibs
// by default the map value will be use `defaultMapParticipantLibs`.
func WithParticipantLibs(val map[string]bool) Option {
	return func(r *rpcSeqDiagram) error {
		if val == nil {
			return errEmptyParticipantList
		}

		r.mapParticipantLibs = val

		return nil
	}
}
