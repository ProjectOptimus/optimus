package osc

import (
	"regexp"
	"strings"
)

/*
Dedent takes a string and returns the same string with all leading whitespace
characters removed, with a maximum of the smallest count of those characters. It
is very similar to Python's [textwrap.dedent()]. Dedent also removes any newline-only lines

Note that this implementation currently depends on the line separator being LF
(\n).

[textwrap.dedent()]: https://docs.python.org/3/library/textwrap.html#textwrap.dedent
*/
func Dedent(s string) string {
	lines := strings.Split(s, "\n")

	leadingWSRegex := regexp.MustCompile(`^\s+`)
	onlyWSRegex := regexp.MustCompile(`^\s+$`)

	// leastReps holds the value by which to dedent the input string lines
	leastReps := 1_000_000
	for _, line := range lines {
		// TODO: Don't bother processing the rest of the full string as soon as
		// you see a non-indented line (that isn't just an empty string), since
		// that would mean that there's no overall dedenting to do. Doing that
		// check here since this is the first loop we run in this func
		if !leadingWSRegex.MatchString(line) && line != "" {
			return s
		}

		// We also don't want to process whitespace-only lines, because e.g. the
		// final input line may be arbitrarily indented in the raw input, which
		// might not be the dedentation we actually want
		if onlyWSRegex.MatchString(line) {
			continue
		}

		reps := countLeadingWhitespace(line)

		if reps != 0 && reps <= leastReps {
			leastReps = reps
		}
	}

	var dedentedLines []string
	for i, line := range lines {
		// Skip empty lines in source post-Split(), and also any trailing all-WS lines
		if line == "" || (i == len(lines) && onlyWSRegex.MatchString(line)) {
			continue
		}

		if leadingWSRegex.MatchString(line) && !onlyWSRegex.MatchString(line) {
			dedentedLine := line[leastReps:]
			dedentedLines = append(dedentedLines, dedentedLine)
		}
	}

	dedented := strings.Join(dedentedLines, "\n")
	return dedented
}

func countLeadingWhitespace(s string) int {
	if !regexp.MustCompile(`^\s`).MatchString(s) {
		return 0
	}

	var reps int

	re := regexp.MustCompile(`\s`)

	for _, c := range s {
		if !re.Match([]byte(string(c))) {
			// since we already checked if the first rune was a whitespace
			// character, this will return early as soon as we hit any
			// non-whitespace runes
			return reps
		}
		reps++
	}

	// Don't know if we still need to return here, but
	return reps
}
