package rpcseqdiagram

import "errors"

var (
	errEmptyRPC                        = errors.New("please specify the RPC names! e.g. RPC=GetData,GetList")
	errEmptyHandlerStruct              = errors.New("please define handler struct annotation")
	errEmptyUsecaseStruct              = errors.New("please define usecase struct annotation")
	errEmptyTargetPath                 = errors.New("failed to set targetPath: value can't be empty")
	errEmptyReadmePath                 = errors.New("failed to set readmePath: value can't be empty")
	errEmptyGrpchandlerPath            = errors.New("failed to set grpchandlerPath: value can't be empty")
	errEmptyUsecasePath                = errors.New("failed to set usecasePath: value can't be empty")
	errEmptyIsAStructHandler           = errors.New("failed to set isAStructHandler: value can't be empty")
	errEmptyIsAStructUsecase           = errors.New("failed to set isAStructUsecase: value can't be empty")
	errEmptyStartRPCSequenceDiagramDoc = errors.New("failed to set startRPCSequenceDiagramDoc: value can't be empty")
	errEmptyEndRPCSequenceDiagramDoc   = errors.New("failed to set endRPCSequenceDiagramDoc: value can't be empty")
	errEmptyCommentForMermaidJS        = errors.New("failed to set commentForMermaidJS: value can't be empty")
	errEmptyMermaidJSReplace           = errors.New("failed to set mermaidJSReplace: value can't be empty")
	errEmptyEndLeadingSpace            = errors.New("failed to set mapLeadingSpaceEnd: value can't be empty")
	errEmptyTrimConditions             = errors.New("failed to set mapTrimConditions: value can't be empty")
	errEmptyParticipantList            = errors.New("failed to set mapParticipantLibs: value can't be empty")
)
