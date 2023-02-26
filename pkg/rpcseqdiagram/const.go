package rpcseqdiagram

// const used in this script.
const (
	sequenceDiagramPath        = "docs/sequence-diagrams/rpc"
	readmePath                 = "README.md"
	grpchandlerPath            = "internal/adapters/grpchandler"
	usecasePath                = "internal/usecases"
	goExtension                = ".go"
	goTestExtension            = "_test.go"
	goMockExtension            = "_mock.go"
	mdExtension                = ".md"
	typePrefix                 = "type "
	structPrefix               = " struct {"
	usecasePrefix              = " usecases."
	funcOpenParenthesesPrefix  = "func ("
	funcPrefix                 = "func "
	spacePointer               = " *"
	openParentheses            = "("
	closeParenthesesSpace      = ") "
	closeCurlyBracket          = "}"
	spaceOpenCurlyBracket      = " {"
	dot                        = "."
	callingTypeParticipant     = "participant"
	callingTypeUCMethod        = "usecase-method"
	lastTypeMethod             = "method"
	lastTypeFunction           = "function"
	serviceClientLib           = "ServiceClient"
	serviceLib                 = "Service"
	repoLib                    = "Repo"
	factoryLib                 = "Factory"
	isAStructHandler           = "is a struct for handler"
	isAStructUsecase           = "is a struct for usecase"
	rangeLoop                  = ":= range "
	ifCondition                = "if "
	elseIfCondition            = " else if "
	elseCondition              = "} else"
	switchCondition            = "switch "
	caseCondition              = "case "
	defaultCondition           = "default:"
	colon                      = ":"
	comment                    = "//"
	startRPCSequenceDiagramDoc = "<!-- start rpc sequence diagram doc -->"
	endRPCSequenceDiagramDoc   = "<!-- end rpc sequence diagram doc -->"
	commentForMermaidJS        = "NOTE: this comments for generating mermaid js code"
	mermaidJSReplace           = "// mermaid js replace: "
	dirPermission              = 0755
)