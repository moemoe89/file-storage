package rpcseqdiagram

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// findContentUsecasesStructUsecaseMapParticipant get the content of usecases,
// struct usecase name and mermaid JS participant name + code.
func (r *rpcSeqDiagram) findContentUsecasesStructUsecaseMapParticipant() error {
	var structUsecase string

	var findStructUsecase bool

	contentUsecases := make([]string, 0)
	mapParticipantCode := make(map[string]struct{}, 0)
	mapParticipant := make(map[string]participantModel, 0)

	// read all files in usecase path.
	err := filepath.Walk(r.usecasePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to read each file: %w", err)
			}

			// only read go files exclude test and mock.
			if !strings.Contains(path, goExtension) ||
				strings.Contains(path, goTestExtension) ||
				strings.Contains(path, goMockExtension) {
				return nil
			}

			file, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			// put the content of usecase in array.
			contentUsecases = append(contentUsecases, string(file))

			return r.findStructUsecaseMapParticipant(
				mapParticipantCode, mapParticipant, file, &findStructUsecase, &structUsecase,
			)
		})
	if err != nil {
		return fmt.Errorf("failed to read all files: %w", err)
	}

	if structUsecase == "" {
		return errEmptyUsecaseStruct
	}

	r.contentUsecases = contentUsecases
	r.structUsecase = structUsecase
	r.mapParticipant = mapParticipant

	return nil
}

func (r *rpcSeqDiagram) findStructUsecaseMapParticipant(
	mapParticipantCode map[string]struct{}, mapParticipant map[string]participantModel,
	file []byte, findStructUsecase *bool, structUsecase *string,
) error {
	scanner := bufio.NewScanner(bytes.NewReader(file))

	// read line by line.
	for scanner.Scan() {
		text := scanner.Text()

		// if struct usecase found, get the struct usecase name,
		// if structUsecase is not empty, skip assign the value.
		if *findStructUsecase && *structUsecase == "" {
			*structUsecase = strings.ReplaceAll(text, typePrefix, "")
			*structUsecase = strings.ReplaceAll(*structUsecase, structPrefix, "")
		}

		// if the text contains the struct usecase.
		r.findStructUsecase(text, findStructUsecase)

		findLib := r.findLibOnParticipant(text)

		// skip if the struct usecase not found or the text not contains participant we want to catch.
		if !*findStructUsecase || !findLib {
			continue
		}

		r.pushToMapParticipant(mapParticipantCode, mapParticipant, text)
	}

	return nil
}

func (r *rpcSeqDiagram) findStructUsecase(text string, findStructUsecase *bool) {
	// if the text contains the struct usecase.
	if strings.Contains(text, r.isAStructUsecase) {
		*findStructUsecase = true
	}

	// if the text is close curly bracket } for the struct without leading space.
	if text == closeCurlyBracket {
		*findStructUsecase = false
	}
}

func (r *rpcSeqDiagram) findLibOnParticipant(text string) bool {
	for k := range r.mapParticipantLibs {
		if strings.Contains(text, k) {
			return true
		}
	}

	return false
}

func (r *rpcSeqDiagram) pushToMapParticipant(
	mapParticipantCode map[string]struct{}, mapParticipant map[string]participantModel,
	text string,
) {
	split := strings.Split(text, dot)
	if len(split) <= 1 || (len(split) > 1 && split[1] == "") {
		return
	}

	fieldLib, participantName := split[0], split[1]

	// checks if participantName exist in mapParticipantLibs and ahs true value
	// because some case has factory with the same multiple package
	// then use the field instead of the type name
	// e.g xFactory x.Factory, yFactory y.Factory etc.
	val, ok := r.mapParticipantLibs[participantName]
	if val && ok {
		splitSpace := strings.Split(split[0], " ")
		factoryType := splitSpace[len(splitSpace)-1]
		r := []rune(factoryType)
		r[0] = unicode.ToUpper(r[0])

		participantName = string(r) + split[1]
	}

	code := make([]byte, 0)

	// turn every single capitalize to code character
	// e.g. xyzClient client.XyzServiceClient as XSC.
	for i := range participantName {
		if isUpper(participantName[i]) {
			code = append(code, participantName[i])
		}
	}

	// if code already exists, then adding 'A' to avoid duplicate.
	if _, ok := mapParticipantCode[string(code)]; ok {
		code = append(code, 'A')
	}

	mapParticipantCode[string(code)] = struct{}{}

	fields := strings.Split(fieldLib, " ")

	// get the participant name e.g. xyzClient client.XyzServiceClient
	// as XyzServiceClient.
	participant := space.ReplaceAllString(fields[0], "")

	mapParticipant[participant] = participantModel{
		Code: string(code),
		Name: participantName,
	}
}

// findMapUsecaseMethodFunction is for collects all method and function in usecases
// then put it on the map.
func (r *rpcSeqDiagram) findMapUsecaseMethodFunction() {
	mapUsecaseMethod := make(map[string][]usecaseModel, 0)
	mapUsecaseFunction := make(map[string][]usecaseModel, 0)

	// read the content of usecases.
	for _, c := range r.contentUsecases {
		scanner := bufio.NewScanner(bytes.NewReader([]byte(c)))

		// read line by line.
		for scanner.Scan() {
			text := scanner.Text()

			// if found the child method e.g. u.getChildUsecase().
			if strings.Contains(text, funcOpenParenthesesPrefix) &&
				strings.Contains(text, r.structUsecase+closeParenthesesSpace) {
				split := strings.Split(text, r.structUsecase+closeParenthesesSpace)

				// find the usecase method name.
				split = strings.Split(split[1], openParentheses)
				usecase := split[0]

				mapUsecaseMethod[usecase] = []usecaseModel{}
			}

			// if found the child function e.g. getChildUsecase (without receiver variable usecase).
			if strings.Contains(text, funcPrefix) && !strings.Contains(text, r.structUsecase+closeParenthesesSpace) {
				text = strings.ReplaceAll(text, funcPrefix, "")

				split := strings.Split(text, openParentheses)

				// find the usecase function name.
				function := split[0]

				if space.ReplaceAllString(function, "") != "" {
					mapUsecaseFunction[function] = []usecaseModel{}
				}
			}
		}
	}

	r.mapUsecaseMethod = mapUsecaseMethod
	r.mapUsecaseFunction = mapUsecaseFunction
}

// findMapUsecaseSequence is for build the mermaid js sequence diagram code based on usecases.
func (r *rpcSeqDiagram) findMapUsecaseSequence() {
	var lastUsecaseType, lastMethodFunction, lastReceiver, lastSequence string

	mapUsecaseSequence := make(map[string]string, 0)
	operators := make([]operatorModel, 0)
	tabs := make([]string, 0)

	var findSwitch, findFirstCase, findMermaidCode bool

	// generate sequence diagram code for every usecases.
	for _, c := range r.contentUsecases {
		scanner := bufio.NewScanner(bytes.NewReader([]byte(c)))

		// read line by line.
		for scanner.Scan() {
			text := scanner.Text()

			// if text contains beginning of method, reset the sequence diagram generator state.
			if strings.Contains(text, funcOpenParenthesesPrefix) &&
				strings.Contains(text, r.structUsecase+closeParenthesesSpace) {
				split := strings.Split(text, r.structUsecase+closeParenthesesSpace)

				lastReceiver = strings.ReplaceAll(split[0], funcOpenParenthesesPrefix, "")
				lastReceiver = strings.ReplaceAll(lastReceiver, spacePointer, "")

				usecase := strings.Split(split[1], openParentheses)[0]

				lastUsecaseType = lastTypeMethod
				lastMethodFunction = usecase

				resetValue(&lastSequence, &operators, &tabs, &findSwitch, &findFirstCase)
			}

			// if text contains beginning of function, reset the sequence diagram generator state.
			if strings.Contains(text, funcPrefix) &&
				!strings.Contains(text, r.structUsecase+closeParenthesesSpace) {
				text = strings.ReplaceAll(text, funcPrefix, "")

				function := strings.Split(text, openParentheses)[0]

				if space.ReplaceAllString(function, "") != "" {
					lastUsecaseType = lastTypeFunction
					lastMethodFunction = function

					resetValue(&lastSequence, &operators, &tabs, &findSwitch, &findFirstCase)
				}
			}

			r.findLastUsecaseIsMethod(
				text, &lastUsecaseType, &lastMethodFunction, &lastReceiver, &lastSequence,
				&findSwitch, &findFirstCase, &findMermaidCode, &tabs, &operators, mapUsecaseSequence,
			)
		}
	}

	r.mapUsecaseSequence = mapUsecaseSequence
}

func resetValue(
	lastSequence *string, operators *[]operatorModel, tabs *[]string,
	findSwitch, findFirstCase *bool,
) {
	*lastSequence = ""
	*operators = make([]operatorModel, 0)
	*tabs = make([]string, 0)
	*findSwitch = false
	*findFirstCase = false
}

func (r *rpcSeqDiagram) findLastUsecaseIsMethod(
	text string,
	lastUsecaseType, lastMethodFunction, lastReceiver, lastSequence *string,
	findSwitch, findFirstCase, findMermaidCode *bool,
	tabs *[]string,
	operators *[]operatorModel,
	mapUsecaseSequence map[string]string,
) {
	if *lastUsecaseType != lastTypeMethod {
		return
	}

	if *findMermaidCode {
		if strings.Contains(text, "*/") {
			*findMermaidCode = false
		} else {
			*lastSequence += strings.Join(*tabs, "") + strings.ReplaceAll(text, "//", "") + "\n"
			mapUsecaseSequence[*lastMethodFunction] = *lastSequence
		}
	}

	// this is fall into mermaid js comment in the code.
	if strings.Contains(text, r.commentForMermaidJS) {
		*findMermaidCode = true
	}

	// this is fall into `for` condition.
	r.checksForCondition(text, tabs, lastSequence, operators)

	// this is fall into `if` condition.
	r.checksIfCondition(text, tabs, lastReceiver, lastSequence, operators)

	// this is fall into `switch` condition.
	r.checksSwitchCondition(text, tabs, findSwitch, findFirstCase, operators)

	// this is fall into `case` condition.
	r.checksCaseCondition(text, tabs, findSwitch, findFirstCase, lastSequence)

	// this is fall into `default` condition (from switch case).
	r.checksDefaultCondition(text, tabs, findSwitch, lastSequence)

	// this is fall into close curly `}` for all operator condition (for, if, switch).
	r.checksCloseCurly(text, tabs, lastMethodFunction, lastSequence, operators, mapUsecaseSequence)

	// if the text call another usecase method / another participant (RPC, repo).
	if !strings.Contains(text, *lastReceiver+dot) {
		return
	}

	splitText := strings.Split(text, *lastReceiver+dot)

	split := strings.Split(splitText[1], openParentheses)

	// if calling another participant (RPC, repo).
	skip := r.checksCallAnotherParticipant(split, tabs, lastMethodFunction, lastSequence, mapUsecaseSequence)
	if skip {
		return
	}

	// if calling another usecase method.
	r.checksCallAnotherMethod(split, tabs, lastMethodFunction, lastSequence, mapUsecaseSequence)
}

func (r *rpcSeqDiagram) checksForCondition(
	text string, tabs *[]string, lastSequence *string, operators *[]operatorModel,
) {
	if strings.Contains(text, rangeLoop) {
		split := strings.Split(text, rangeLoop)
		object := strings.ReplaceAll(split[1], spaceOpenCurlyBracket, "")

		*tabs = append(*tabs, "\t")

		*lastSequence += strings.Join(*tabs, "") + "loop Iterates " + object + "\n"

		*operators = append(*operators, operatorModel{
			operator:     text,
			leadingSpace: r.countLeadingSpaces(text),
			isSKip:       false,
		})
	}
}

func (r *rpcSeqDiagram) checksIfCondition(
	text string, tabs *[]string, lastReceiver, lastSequence *string, operators *[]operatorModel,
) {
	if strings.Contains(text, ifCondition) && !strings.Contains(text, elseIfCondition) &&
		strings.Contains(text, spaceOpenCurlyBracket) &&
		(!strings.Contains(text, comment) || strings.Contains(text, mermaidJSReplace)) {
		isSkip := false
		if strings.Contains(text, *lastReceiver+dot) {
			isSkip = true
		} else {
			*tabs = append(*tabs, "\t")
			*lastSequence += strings.Join(*tabs, "") + "alt " + r.trimCondition(text) + "\n"
		}

		*operators = append(*operators, operatorModel{
			operator:     text,
			leadingSpace: r.countLeadingSpaces(text),
			isSKip:       isSkip,
		})
	}
}

func (r *rpcSeqDiagram) checksSwitchCondition(
	text string, tabs *[]string, findSwitch, findFirstCase *bool, operators *[]operatorModel,
) {
	if strings.Contains(text, switchCondition) && strings.Contains(text, spaceOpenCurlyBracket) {
		*findSwitch = true
		*findFirstCase = false

		*tabs = append(*tabs, "\t")

		*operators = append(*operators, operatorModel{
			operator:     text,
			leadingSpace: r.countLeadingSpaces(text),
			isSKip:       false,
		})
	}
}

func (r *rpcSeqDiagram) checksCaseCondition(
	text string, tabs *[]string, findSwitch, findFirstCase *bool, lastSequence *string,
) {
	if *findSwitch && strings.Contains(text, caseCondition) && strings.Contains(text, colon) {
		altElse := "else"
		if !*findFirstCase {
			altElse = "alt"
			*findFirstCase = true
		}

		*lastSequence += strings.Join(*tabs, "") + altElse + " " + r.trimCondition(text) + "\n"
	}
}

func (r *rpcSeqDiagram) checksDefaultCondition(
	text string, tabs *[]string, findSwitch *bool, lastSequence *string,
) {
	if *findSwitch && strings.Contains(text, defaultCondition) {
		altElse := "else"

		*lastSequence += strings.Join(*tabs, "") + altElse + " " + r.trimCondition(text) + "\n"
	}
}

func (r *rpcSeqDiagram) checksCloseCurly(
	text string, tabs *[]string,
	lastMethodFunction, lastSequence *string,
	operators *[]operatorModel, mapUsecaseSequence map[string]string,
) {
	if strings.Contains(text, closeCurlyBracket) && len(*operators) > 0 {
		unpointerOperators := *operators
		// if the text is contains `else` condition
		if strings.Contains(text, elseCondition) && !unpointerOperators[len(unpointerOperators)-1].isSKip {
			*lastSequence += strings.Join(*tabs, "") + "else " + r.trimCondition(text) + "\n"
			return
		}

		leadingSpace := r.countLeadingSpaces(text)

		if unpointerOperators[len(unpointerOperators)-1].leadingSpace == leadingSpace {
			if !unpointerOperators[len(unpointerOperators)-1].isSKip {
				*lastSequence += strings.Join(*tabs, "") + "end\n"
				unPointerTabs := *tabs
				*tabs = unPointerTabs[:len(unPointerTabs)-1]
			}

			*operators = unpointerOperators[:len(unpointerOperators)-1]

			mapUsecaseSequence[*lastMethodFunction] = *lastSequence
		}
	}
}

func (r *rpcSeqDiagram) checksCallAnotherParticipant(
	split []string,
	tabs *[]string,
	lastMethodFunction, lastSequence *string,
	mapUsecaseSequence map[string]string,
) bool {
	if strings.Contains(split[0], dot) {
		split = strings.Split(split[0], dot)

		callingName := split[0]

		split = strings.Split(split[1], openParentheses)
		callingFunc := split[0]

		if _, ok := r.mapParticipant[callingName]; ok {
			*lastSequence += strings.Join(*tabs, "") + callingTypeParticipant + "|" + callingName + "|" + callingFunc + "|up\n"
			*lastSequence += strings.Join(*tabs, "") + callingTypeParticipant + "|" + callingName + "|" + callingFunc + "|down\n"

			mapUsecaseSequence[*lastMethodFunction] = *lastSequence

			r.mapUsecaseMethod[*lastMethodFunction] = append(r.mapUsecaseMethod[*lastMethodFunction], usecaseModel{
				callingName: callingName,
				callingFunc: callingFunc,
				callingType: callingTypeParticipant,
			})
		}

		return true
	}

	return false
}

func (r *rpcSeqDiagram) checksCallAnotherMethod(
	split []string,
	tabs *[]string,
	lastMethodFunction, lastSequence *string,
	mapUsecaseSequence map[string]string,
) {
	callingName := split[0]
	if _, ok := r.mapUsecaseMethod[callingName]; ok {
		*lastSequence += strings.Join(*tabs, "") + callingTypeUCMethod + "|" + callingName + "\n"

		mapUsecaseSequence[*lastMethodFunction] = *lastSequence

		r.mapUsecaseMethod[*lastMethodFunction] = append(r.mapUsecaseMethod[*lastMethodFunction], usecaseModel{
			callingName: callingName,
			callingType: callingTypeUCMethod,
		})
	}
}
