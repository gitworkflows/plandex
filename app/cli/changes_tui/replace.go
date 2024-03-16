package changes_tui

import (
	"plandex/term"
	"strings"

	"github.com/fatih/color"
	"github.com/muesli/reflow/wrap"
)

const replacementPrependLines = 20
const replacementAppendLines = 20

type oldReplacementRes struct {
	old               string
	oldDisplay        string
	prependContent    string
	numLinesPrepended int
	appendContent     string
	numLinesAppended  int
}

func (m changesUIModel) getReplacementOldDisplay() oldReplacementRes {
	oldContent := m.selectionInfo.currentRep.Old
	originalFile := m.selectionInfo.currentFilesBeforeReplacement.Files[m.selectionInfo.currentPath]

	oldContent = strings.ReplaceAll(oldContent, "\\`\\`\\`", "```")
	originalFile = strings.ReplaceAll(originalFile, "\\`\\`\\`", "```")

	// log.Printf("oldContent: %v", oldContent)
	// log.Printf("originalFile: %v", originalFile)

	fileIdx := strings.Index(originalFile, oldContent)
	if fileIdx == -1 {
		panic("old content not found in full file") // should never happen
	}

	toPrepend := ""
	numLinesPrepended := 0
	for i := fileIdx - 1; i >= 0; i-- {
		s := string(originalFile[i])
		toPrepend = s + toPrepend
		if originalFile[i] == '\n' {
			numLinesPrepended++
			if numLinesPrepended == replacementPrependLines {
				break
			}
		}
	}
	prependedToStart := strings.Index(originalFile, toPrepend) == 0

	toPrepend = strings.TrimLeft(toPrepend, "\n")
	if !prependedToStart {
		toPrepend = "…\n" + toPrepend
	}

	toAppend := ""
	numLinesAppended := 0
	for i := fileIdx + len(oldContent); i < len(originalFile); i++ {
		s := string(originalFile[i])
		if s == "\t" {
			s = "  "
		}
		toAppend += s
		if originalFile[i] == '\n' {
			numLinesAppended++
			if numLinesAppended == replacementAppendLines {
				break
			}
		}
	}
	appendedToEnd := strings.Index(originalFile, toAppend) == len(originalFile)-len(toAppend)

	toAppend = strings.TrimRight(toAppend, "\n")

	if !appendedToEnd {
		toAppend += "\n…"
	}

	wrapWidth := m.changeOldViewport.Width - 6
	toPrepend = wrap.String(toPrepend, wrapWidth)
	oldContent = wrap.String(oldContent, wrapWidth)
	toAppend = wrap.String(toAppend, wrapWidth)

	toPrependLines := strings.Split(toPrepend, "\n")
	for i, line := range toPrependLines {
		toPrependLines[i] = color.New(color.FgWhite).Sprint(line)
	}
	toPrepend = strings.Join(toPrependLines, "\n")

	oldContentLines := strings.Split(oldContent, "\n")
	for i, line := range oldContentLines {
		oldContentLines[i] = color.New(term.ColorHiRed).Sprint(line)
	}
	oldContent = strings.Join(oldContentLines, "\n")

	toAppendLines := strings.Split(toAppend, "\n")
	for i, line := range toAppendLines {
		toAppendLines[i] = color.New(color.FgWhite).Sprint(line)
	}
	toAppend = strings.Join(toAppendLines, "\n")

	oldDisplayContent := toPrepend + oldContent + toAppend

	numLinesPrepended = len(toPrependLines)
	numLinesAppended = len(toAppendLines)

	return oldReplacementRes{
		oldContent,
		oldDisplayContent,
		toPrepend,
		numLinesPrepended,
		toAppend,
		numLinesAppended,
	}
}

func (m changesUIModel) getReplacementNewDisplay(prependContent, appendContent string) (string, string) {
	newContent := m.selectionInfo.currentRep.New

	newContent = strings.ReplaceAll(newContent, "\\`\\`\\`", "```")

	newContent = wrap.String(newContent, m.changeNewViewport.Width-6)

	newContentLines := strings.Split(newContent, "\n")
	for i, line := range newContentLines {
		newContentLines[i] = color.New(term.ColorHiGreen).Sprint(line)
	}
	newContent = strings.Join(newContentLines, "\n")

	return newContent, prependContent + newContent + appendContent
}
