package rpcseqdiagram

import "strings"

func (r *rpcSeqDiagram) iteratesChildParticipant(
	md string,
	usecase string,
	participantCodeDuplicate map[string]struct{},
) string {
	if _, ok := r.mapUsecaseMethod[usecase]; !ok {
		return md
	}

	for _, v := range r.mapUsecaseMethod[usecase] {
		if v.callingType == callingTypeParticipant {
			if _, ok := r.mapParticipant[v.callingName]; !ok {
				continue
			}

			participant := r.mapParticipant[v.callingName]

			if _, ok := participantCodeDuplicate[participant.Code]; !ok {
				md += "\tparticipant " + participant.Code + " as " + participant.Name + "\n"
				participantCodeDuplicate[participant.Code] = struct{}{}
			}
		} else if v.callingType == callingTypeUCMethod {
			md = r.iteratesChildParticipant(md, v.callingName, participantCodeDuplicate)
		}
	}

	return md
}

func (r *rpcSeqDiagram) iteratesChildUsecase(num, usecaseFunctionLoopSequence string, usecase string) string {
	if _, ok := r.mapUsecaseMethod[usecase]; !ok {
		return usecaseFunctionLoopSequence
	}

	for _, v := range r.mapUsecaseMethod[usecase] {
		if v.callingType == callingTypeParticipant {
			var upReplace, downReplace string

			if _, ok := r.mapParticipant[v.callingName]; ok {
				participant := r.mapParticipant[v.callingName]

				upReplace += "\tUC" + num + "->>+" + participant.Code + ": Call `" + v.callingFunc + "`\n"
				downReplace += "\t" + participant.Code + "-->>-UC" + num + ": return\n"
			}

			oldReplace := callingTypeParticipant + "|" + v.callingName + "|" + v.callingFunc + "|"

			usecaseFunctionLoopSequence = strings.ReplaceAll(usecaseFunctionLoopSequence, oldReplace+"up\n", upReplace)
			usecaseFunctionLoopSequence = strings.ReplaceAll(usecaseFunctionLoopSequence, oldReplace+"down\n", downReplace)
		} else if v.callingType == callingTypeUCMethod {
			usecaseFunctionLoopSequence = strings.ReplaceAll(
				usecaseFunctionLoopSequence,
				callingTypeUCMethod+"|"+v.callingName+"\n",
				r.mapUsecaseSequence[v.callingName],
			)

			usecaseFunctionLoopSequence = r.iteratesChildUsecase(num, usecaseFunctionLoopSequence, v.callingName)
		}
	}

	return usecaseFunctionLoopSequence
}
