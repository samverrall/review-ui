package syntax

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Highlighter handles syntax highlighting for code
type Highlighter struct {
	formatter chroma.Formatter
	style     *chroma.Style
}

// NewHighlighter creates a new syntax highlighter with Tailwind moon theme colors
func NewHighlighter() *Highlighter {
	// Create a custom Tailwind-inspired moon theme
	// Using Tailwind color palette for a cohesive look
	baseStyle := styles.Get("monokai") // Start with monokai as base

	builder := baseStyle.Builder()

	// Tailwind-inspired moon theme colors
	// Don't set background - let it inherit from the line style (cursor line, selection, etc.)
	builder.Add(chroma.Background, "") // No background - transparent

	// Text colors using Tailwind slate scale
	builder.Add(chroma.Text, "fg:248")                    // slate-100 (light text)
	builder.Add(chroma.Comment, "fg:100")                 // slate-500 (muted comments)
	builder.Add(chroma.CommentMultiline, "fg:100")
	builder.Add(chroma.CommentSingle, "fg:100")

	// Keywords: violet/purple family
	builder.Add(chroma.Keyword, "fg:147")                 // violet-400
	builder.Add(chroma.KeywordConstant, "fg:147")
	builder.Add(chroma.KeywordDeclaration, "fg:147")
	builder.Add(chroma.KeywordReserved, "fg:147")
	builder.Add(chroma.KeywordType, "fg:135")             // violet-500

	// Names and identifiers
	builder.Add(chroma.Name, "fg:248")                    // slate-100
	builder.Add(chroma.NameFunction, "fg:158")            // teal-300
	builder.Add(chroma.NameClass, "fg:167")               // orange-300
	builder.Add(chroma.NameBuiltin, "fg:135")             // violet-500
	builder.Add(chroma.NameAttribute, "fg:147")           // violet-400

	// Literals
	builder.Add(chroma.LiteralString, "fg:149")           // emerald-300 (greenish)
	builder.Add(chroma.LiteralStringDouble, "fg:149")
	builder.Add(chroma.LiteralStringSingle, "fg:149")
	builder.Add(chroma.LiteralNumber, "fg:173")           // amber-300
	builder.Add(chroma.LiteralNumberInteger, "fg:173")
	builder.Add(chroma.LiteralNumberFloat, "fg:173")

	// Operators and punctuation
	builder.Add(chroma.Operator, "fg:248")                // slate-100
	builder.Add(chroma.Punctuation, "fg:248")

	// Generic styles
	builder.Add(chroma.GenericEmph, "italic")
	builder.Add(chroma.GenericStrong, "bold")
	builder.Add(chroma.GenericUnderline, "underline")

	// Error handling
	builder.Add(chroma.Error, "fg:197")                   // red-500

	customStyle, err := builder.Build()
	if err != nil {
		// Fallback to monokai if custom style fails
		customStyle = baseStyle
	}

	formatter := formatters.TTY256

	return &Highlighter{
		formatter: formatter,
		style:     customStyle,
	}
}

// Highlight applies syntax highlighting to the given code for the specified file
func (h *Highlighter) Highlight(filename, code string) (string, error) {
	// Detect language from filename
	lexer := lexers.Match(filename)
	if lexer == nil {
		// Try to detect from extension if filename matching fails
		ext := filepath.Ext(filename)
		if ext != "" {
			lexer = lexers.Get(ext[1:]) // Remove the leading dot
		}
		if lexer == nil {
			// Fallback to plain text
			return code, nil
		}
	}

	// Ensure we have a lexer
	lexer = chroma.Coalesce(lexer)

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code, err
	}

	// Format with TTY colors
	var buf strings.Builder
	err = h.formatter.Format(&buf, h.style, iterator)
	if err != nil {
		return code, err
	}

	return buf.String(), nil
}

