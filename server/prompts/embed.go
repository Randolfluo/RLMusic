package prompts

import "embed"

// FS embeds all prompt .md files in this directory into the binary.
// The *.md pattern excludes this .go file and the README.md (documentation only, not used at runtime).
//
//go:embed prompt_*.md
var FS embed.FS
