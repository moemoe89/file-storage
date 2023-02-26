package rpcseqdiagram

import "strings"

// check the character is upper case or not.
func isUpper(b byte) bool {
	if b >= 'A' && b <= 'Z' {
		return true
	}

	return false
}

// lower the character.
func toLower(b byte) byte {
	if isUpper(b) {
		b += 'a' - 'A'
	}

	return b
}

// countLeadingSpaces until meet character in mapLeadingSpaceEnd.
func (r *rpcSeqDiagram) countLeadingSpaces(line string) int {
	count := 0

	for _, v := range line {
		if !r.mapEndLeadingSpace[v] {
			count++
		} else {
			break
		}
	}

	return count
}

// trimCondition for trim unnecessary character
// e.g in if { , else }, case xxx:
// for ; in Mermaid JS will decide as last character,
// any character after ; will not count, then replace ; with ,
// also replacing the text by specific mermaid js comment in the code
// e.g. if errors.Is(err, apperr.ErrNotFound) { // mermaid js replace: if data not found.
func (r *rpcSeqDiagram) trimCondition(text string) string {
	if strings.Contains(text, r.mermaidJSReplace) {
		split := strings.Split(text, r.mermaidJSReplace)
		text = split[1]
	}

	for k, v := range r.mapTrimConditions {
		text = strings.ReplaceAll(text, k, v)
	}

	text = strings.TrimSpace(text)

	return text
}
