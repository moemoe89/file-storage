package rpcseqdiagram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTargetPath(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		targetPath string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set target path value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "target",
				},
				want:    "target",
				wantErr: nil,
			}
		},
		"Failed set target path value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					targetPath: "initial",
				},
				want:    "initial",
				wantErr: errEmptyTargetPath,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				targetPath: tt.fields.targetPath,
			}

			err := WithTargetPath(tt.args.value)(r)

			assert.Equal(t, tt.want, r.targetPath)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithReadmePath(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		readmePath string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set readme path value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "readme",
				},
				want:    "readme",
				wantErr: nil,
			}
		},
		"Failed set target readme value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					readmePath: "initial",
				},
				want:    "initial",
				wantErr: errEmptyReadmePath,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				readmePath: tt.fields.readmePath,
			}

			err := WithReadmePath(tt.args.value)(r)

			assert.Equal(t, tt.want, r.readmePath)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithGrpcHandlerPath(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		grpchandlerPath string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set grpc handler path value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "grpc/handler",
				},
				want:    "grpc/handler",
				wantErr: nil,
			}
		},
		"Failed set grpc handler path value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					grpchandlerPath: "initial",
				},
				want:    "initial",
				wantErr: errEmptyGrpchandlerPath,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				grpchandlerPath: tt.fields.grpchandlerPath,
			}

			err := WithGrpcHandlerPath(tt.args.value)(r)

			assert.Equal(t, tt.want, r.grpchandlerPath)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithUsecasePath(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		usecasePath string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set usecase path value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "usecase",
				},
				want:    "usecase",
				wantErr: nil,
			}
		},
		"Failed set usecase path value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					usecasePath: "initial",
				},
				want:    "initial",
				wantErr: errEmptyUsecasePath,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				usecasePath: tt.fields.usecasePath,
			}

			err := WithUsecasePath(tt.args.value)(r)

			assert.Equal(t, tt.want, r.usecasePath)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithIsAStructHandler(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		isAStructHandler string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set is a struct handler value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "isastructhandler",
				},
				want:    "isastructhandler",
				wantErr: nil,
			}
		},
		"Failed set is a struct handler value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					isAStructHandler: "initial",
				},
				want:    "initial",
				wantErr: errEmptyIsAStructHandler,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				isAStructHandler: tt.fields.isAStructHandler,
			}

			err := WithIsAStructHandler(tt.args.value)(r)

			assert.Equal(t, tt.want, r.isAStructHandler)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithIsAStructUsecase(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		isAStructUsecase string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set is a struct usecase value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "isastructusecase",
				},
				want:    "isastructusecase",
				wantErr: nil,
			}
		},
		"Failed set is a struct usecase value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					isAStructUsecase: "initial",
				},
				want:    "initial",
				wantErr: errEmptyIsAStructUsecase,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				isAStructUsecase: tt.fields.isAStructUsecase,
			}

			err := WithIsAStructUsecase(tt.args.value)(r)

			assert.Equal(t, tt.want, r.isAStructUsecase)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithStartRPCSequenceDiagramDoc(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		startRPCSequenceDiagramDoc string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set start rpc sequence diagram doc value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "startrpcsequencediagramdoc",
				},
				want:    "startrpcsequencediagramdoc",
				wantErr: nil,
			}
		},
		"Failed set start rpc sequence diagram doc value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					startRPCSequenceDiagramDoc: "initial",
				},
				want:    "initial",
				wantErr: errEmptyStartRPCSequenceDiagramDoc,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				startRPCSequenceDiagramDoc: tt.fields.startRPCSequenceDiagramDoc,
			}

			err := WithStartRPCSequenceDiagramDoc(tt.args.value)(r)

			assert.Equal(t, tt.want, r.startRPCSequenceDiagramDoc)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithEndRPCSequenceDiagramDoc(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		endRPCSequenceDiagramDoc string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set end rpc sequence diagram doc value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "endrpcsequencediagramdoc",
				},
				want:    "endrpcsequencediagramdoc",
				wantErr: nil,
			}
		},
		"Failed set end rpc sequence diagram doc value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					endRPCSequenceDiagramDoc: "initial",
				},
				want:    "initial",
				wantErr: errEmptyEndRPCSequenceDiagramDoc,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				endRPCSequenceDiagramDoc: tt.fields.endRPCSequenceDiagramDoc,
			}

			err := WithEndRPCSequenceDiagramDoc(tt.args.value)(r)

			assert.Equal(t, tt.want, r.endRPCSequenceDiagramDoc)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithCommentForMermaidJS(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		commentForMermaidJS string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set comment for mermaid js value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "commentformermaidjs",
				},
				want:    "commentformermaidjs",
				wantErr: nil,
			}
		},
		"Failed set comment for mermaid js value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					commentForMermaidJS: "initial",
				},
				want:    "initial",
				wantErr: errEmptyCommentForMermaidJS,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				commentForMermaidJS: tt.fields.commentForMermaidJS,
			}

			err := WithCommentForMermaidJS(tt.args.value)(r)

			assert.Equal(t, tt.want, r.commentForMermaidJS)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithMermaidJSReplace(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		mermaidJSReplace string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set mermaid js replace value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "mermaidjsreplace",
				},
				want:    "mermaidjsreplace",
				wantErr: nil,
			}
		},
		"Failed set mermaid js replace value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					mermaidJSReplace: "initial",
				},
				want:    "initial",
				wantErr: errEmptyMermaidJSReplace,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				mermaidJSReplace: tt.fields.mermaidJSReplace,
			}

			err := WithMermaidJSReplace(tt.args.value)(r)

			assert.Equal(t, tt.want, r.mermaidJSReplace)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithEndLeadingSpace(t *testing.T) {
	type args struct {
		value map[rune]bool
	}

	type fields struct {
		mapEndLeadingSpace map[rune]bool
	}

	type test struct {
		args    args
		fields  fields
		want    map[rune]bool
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set map end leading space value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: map[rune]bool{
						'f': true,
					},
				},
				want: map[rune]bool{
					'f': true,
				},
				wantErr: nil,
			}
		},
		"Failed set map end leading space value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: nil,
				},
				fields: fields{
					mapEndLeadingSpace: map[rune]bool{
						'f': true,
					},
				},
				want: map[rune]bool{
					'f': true,
				},
				wantErr: errEmptyEndLeadingSpace,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				mapEndLeadingSpace: tt.fields.mapEndLeadingSpace,
			}

			err := WithEndLeadingSpace(tt.args.value)(r)

			assert.Equal(t, tt.want, r.mapEndLeadingSpace)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithTrimConditions(t *testing.T) {
	type args struct {
		value map[string]string
	}

	type fields struct {
		mapTrimConditions map[string]string
	}

	type test struct {
		args    args
		fields  fields
		want    map[string]string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set trim conditions value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: map[string]string{
						"k": "v",
					},
				},
				want: map[string]string{
					"k": "v",
				},
				wantErr: nil,
			}
		},
		"Failed set trim conditions value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: nil,
				},
				fields: fields{
					mapTrimConditions: map[string]string{
						"k": "v",
					},
				},
				want: map[string]string{
					"k": "v",
				},
				wantErr: errEmptyTrimConditions,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				mapTrimConditions: tt.fields.mapTrimConditions,
			}

			err := WithTrimConditions(tt.args.value)(r)

			assert.Equal(t, tt.want, r.mapTrimConditions)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithParticipantLibs(t *testing.T) {
	type args struct {
		value map[string]bool
	}

	type fields struct {
		mapParticipantLibs map[string]bool
	}

	type test struct {
		args    args
		fields  fields
		want    map[string]bool
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set participant libs value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: map[string]bool{
						"participant": true,
					},
				},
				want: map[string]bool{
					"participant": true,
				},
				wantErr: nil,
			}
		},
		"Failed set participant libs value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: nil,
				},
				fields: fields{
					mapParticipantLibs: map[string]bool{
						"participant": true,
					},
				},
				want: map[string]bool{
					"participant": true,
				},
				wantErr: errEmptyParticipantList,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			r := &rpcSeqDiagram{
				mapParticipantLibs: tt.fields.mapParticipantLibs,
			}

			err := WithParticipantLibs(tt.args.value)(r)

			assert.Equal(t, tt.want, r.mapParticipantLibs)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
