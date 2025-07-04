package blueprint

import (
	"fmt"
	"os"
	"strings"
)

// debug is a local wrapper around tool package debug functionality
func debug(format string, args ...interface{}) {
	// We'll use fmt.Fprintf to stderr for now since we can't easily access tool.debug
	// This is a simple implementation that will work for our debugging needs
	if true { // Always log for now during debugging
		fmt.Fprintf(os.Stderr, "[Studio MCP Blueprint] "+format+"\n", args...)
	}
}

// FromArgs creates a new Blueprint from command arguments using tokenization
func FromArgs(args []string) (*Blueprint, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("cannot create blueprint: no command provided")
	}

	if strings.TrimSpace(args[0]) == "" {
		return nil, fmt.Errorf("cannot create blueprint: empty command provided")
	}

	// Debug logging: log the exact arguments received
	debug("FromArgs called with %d arguments:", len(args))
	for i, arg := range args {
		debug("  args[%d]: %q", i, arg)
	}

	bp := &Blueprint{
		BaseCommand: args[0],
		ShellWords:  make([][]Token, len(args)),
	}

	// Tokenize each shell word
	for i, arg := range args {
		tokens := tokenizeShellWord(arg)
		bp.ShellWords[i] = tokens
		debug("  shellword[%d] %q -> %d tokens", i, arg, len(tokens))
		for j, token := range tokens {
			debug("    token[%d]: %T %q", j, token, token.String())
		}
	}

	return bp, nil
}

// tokenizeShellWord tokenizes a single shell word into tokens
func tokenizeShellWord(word string) []Token {
	// Parse mixed content
	tokens := []Token{}
	pos := 0

	for pos < len(word) {
		templateStart := findNextTemplate(word, pos)
		if templateStart == nil {
			// No more templates, add remaining text
			if pos < len(word) {
				tokens = append(tokens, TextToken{Value: word[pos:]})
			}
			break
		}

		// Add text before template
		if templateStart.Start > pos {
			tokens = append(tokens, TextToken{Value: word[pos:templateStart.Start]})
		}

		// Parse template
		templateText := word[templateStart.Start:templateStart.End]
		if token := parseField(templateText); token != nil {
			tokens = append(tokens, token)
		} else {
			tokens = append(tokens, TextToken{Value: templateText})
		}

		pos = templateStart.End
	}

	// Ensure we always return at least one token
	if len(tokens) == 0 {
		tokens = append(tokens, TextToken{Value: word})
	}

	return tokens
}

// templateMatch represents a found template in the text
type templateMatch struct {
	Start int
	End   int
	Type  string
}

// findNextTemplate finds the next template starting from the given position
func findNextTemplate(word string, startPos int) *templateMatch {
	remaining := word[startPos:]

	requiredStart := strings.Index(remaining, "{{")
	optionalStart := strings.Index(remaining, "[")

	// Find the closest template start
	var nextStart int
	var templateType string

	if requiredStart != -1 && (optionalStart == -1 || requiredStart < optionalStart) {
		nextStart = requiredStart
		templateType = "required"
	} else if optionalStart != -1 {
		nextStart = optionalStart
		templateType = "optional"
	} else {
		return nil // No templates found
	}

	absoluteStart := startPos + nextStart

	// Find template end
	var endMarker string
	var markerLength int
	if templateType == "required" {
		endMarker = "}}"
		markerLength = 2
	} else {
		endMarker = "]"
		markerLength = 1
	}

	endIndex := strings.Index(remaining[nextStart:], endMarker)
	if endIndex == -1 {
		// Malformed template - treat rest as text by returning no match
		return nil
	}

	absoluteEnd := absoluteStart + endIndex + markerLength

	return &templateMatch{
		Start: absoluteStart,
		End:   absoluteEnd,
		Type:  templateType,
	}
}

// parseField parses a field enclosed in {{ }} or [ ]
func parseField(field string) Token {
	var content string
	var required bool

	// Determine field type and extract content
	if strings.HasPrefix(field, "{{") && strings.HasSuffix(field, "}}") {
		content = field[2 : len(field)-2] // Remove {{ }}
		required = true
	} else if strings.HasPrefix(field, "[") && strings.HasSuffix(field, "]") {
		content = field[1 : len(field)-1] // Remove [ ]
		required = false
	} else {
		return nil // Not a valid field
	}

	// Parse content for name, description, and modifiers
	var name, description string
	isArray := false
	var originalFlag string

	// Check for description (split on #)
	parts := strings.SplitN(content, "#", 2)
	name = strings.TrimSpace(parts[0])

	// If name is empty, this is not a valid field (e.g., {{}})
	if name == "" {
		return nil
	}

	if len(parts) > 1 {
		description = strings.TrimSpace(parts[1])
	}

	// Check for array notation (...)
	if strings.HasSuffix(name, "...") {
		isArray = true
		name = strings.TrimSuffix(name, "...")
		name = strings.TrimSpace(name)
	}

	// Check for boolean flag (starts with - or --)
	if !required && (strings.HasPrefix(name, "-") || strings.HasPrefix(name, "--")) {
		originalFlag = name
		name = strings.TrimLeft(name, "-")
		if description == "" {
			description = fmt.Sprintf("Enable %s flag", originalFlag)
		}
	}

	return FieldToken{
		Name:         name,
		Description:  description,
		Required:     required,
		IsArray:      isArray,
		OriginalFlag: originalFlag,
	}
}
