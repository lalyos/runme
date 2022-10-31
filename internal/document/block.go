package document

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
)

const (
	KindRaw = iota + 1
	KindCode
)

type Block2 struct {
	Kind int         `json:"kind"`
	Data interface{} `json:"data"`
}

type CodeBlock struct {
	inner        *ast.FencedCodeBlock
	nameResolver *nameResolver
	source       []byte
	cache        struct {
		attributes map[string]string
		content    string
		intro      string
		lines      []string
	}
}

func (b *CodeBlock) rawAttributes() []byte {
	if b.inner.Info == nil {
		return nil
	}

	segment := b.inner.Info.Segment
	info := segment.Value(b.source)
	start, stop := -1, -1

	for i := 0; i < len(info); i++ {
		if start == -1 && info[i] == '{' && i+1 < len(info) && info[i+1] != '}' {
			start = i + 1
		}
		if stop == -1 && info[i] == '}' {
			stop = i
			break
		}
	}

	if start >= 0 && stop >= 0 {
		return bytes.TrimSpace(info[start:stop])
	}

	return nil
}

func (b *CodeBlock) parseAttributes(raw []byte) map[string]string {
	items := bytes.Split(raw, []byte{' '})
	if len(items) == 0 {
		return nil
	}

	result := make(map[string]string)

	for _, item := range items {
		if !bytes.Contains(item, []byte{'='}) {
			continue
		}
		kv := bytes.Split(item, []byte{'='})
		result[string(kv[0])] = string(kv[1])
	}

	return result
}

// Attributes returns code block attributes detected in the first line.
// They are of a form: "sh { attr=value }".
func (b *CodeBlock) Attributes() map[string]string {
	if b.cache.attributes == nil {
		b.cache.attributes = b.parseAttributes(b.rawAttributes())
	}
	return b.cache.attributes
}

// Content returns unaltered snippet as a single blob of text.
func (b *CodeBlock) Content() string {
	if b.cache.content == "" {
		var buf strings.Builder
		for i := 0; i < b.inner.Lines().Len(); i++ {
			line := b.inner.Lines().At(i)
			_, _ = buf.Write(line.Value(b.source))
		}
		b.cache.content = buf.String()
	}
	return b.cache.content
}

// Executable returns an identifier of a program to execute the block.
func (b *CodeBlock) Executable() string {
	if lang := string(b.inner.Language(b.source)); lang != "" {
		return lang
	}
	return ""
}

var replaceEndingRe = regexp.MustCompile(`([^a-z0-9!?\.])$`)

func normalizeIntro(s string) string {
	return replaceEndingRe.ReplaceAllString(s, ".")
}

// Intro returns a normalized description of the code block
// based on the preceding paragraph.
func (b *CodeBlock) Intro() string {
	if b.cache.intro == "" {
		if prevNode := b.inner.PreviousSibling(); prevNode != nil {
			b.cache.intro = normalizeIntro(string(prevNode.Text(b.source)))
		}
	}
	return b.cache.intro
}

// Line returns a normalized code block line at index.
func (b *CodeBlock) Line(idx int) string {
	lines := b.lines()
	if idx >= len(lines) {
		return ""
	}
	return lines[idx]
}

// LineCount returns the number of code block lines.
func (b *CodeBlock) LineCount() int {
	return len(b.lines())
}

func normalizeLine(s string) string {
	return strings.TrimSpace(strings.TrimLeft(s, "$"))
}

func (b *CodeBlock) lines() []string {
	if b.cache.lines == nil {
		var result []string
		for i := 0; i < b.inner.Lines().Len(); i++ {
			line := b.inner.Lines().At(i)
			result = append(result, normalizeLine(string(line.Value(b.source))))
		}
		b.cache.lines = result
	}
	return b.cache.lines
}

// Lines returns all code block lines, normalized.
func (b *CodeBlock) Lines() (result []string) {
	return b.lines()
}

func (b *CodeBlock) MapLines(fn func(string) (string, error)) error {
	var result []string
	for _, line := range b.lines() {
		v, err := fn(line)
		if err != nil {
			return err
		}
		result = append(result, v)
	}
	b.cache.lines = result
	return nil
}

func sanitizeName(s string) string {
	// Handle cases when the first line is defining a variable.
	if idx := strings.Index(s, "="); idx > 0 {
		return sanitizeName(s[:idx])
	}

	limit := len(s)
	if limit > 32 {
		limit = 32
	}
	s = s[0:limit]

	fragments := strings.Split(s, " ")
	if len(fragments) > 1 {
		s = strings.Join(fragments[:2], " ")
	} else {
		s = fragments[0]
	}

	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r == ' ' && b.Len() > 0 {
			_, _ = b.WriteRune('-')
		} else if r >= '0' && r <= '9' || r >= 'a' && r <= 'z' {
			_, _ = b.WriteRune(r)
		}
	}
	return b.String()
}

// Name returns a code block name.
func (b *CodeBlock) Name() string {
	var name string
	if n, ok := b.Attributes()["name"]; ok && n != "" {
		name = n
	} else {
		name = sanitizeName(b.Line(0))
	}
	return b.nameResolver.Get(b, name)
}

func (b *CodeBlock) MarshalJSON() ([]byte, error) {
	type codeBlock struct {
		Attributes map[string]string `json:"attributes,omitempty"`
		Content    string            `json:"content,omitempty"`
		Name       string            `json:"name,omitempty"`
		Language   string            `json:"language,omitempty"`
		Lines      []string          `json:"lines,omitempty"`
	}

	block := codeBlock{
		Attributes: b.Attributes(),
		Content:    b.Content(),
		Name:       b.Name(),
		Language:   b.Executable(),
		Lines:      b.Lines(),
	}

	return json.Marshal(block)
}

type MarkdownBlock struct {
	inner  ast.Node
	source []byte

	content string
}

func findFirstBlockChild(node ast.Node) ast.Node {
	if node == nil {
		return nil
	}

	if !node.HasChildren() || node.FirstChild().Type() != ast.TypeBlock {
		return node
	}

	node = node.FirstChild()
	for node.FirstChild() != nil && node.FirstChild().Type() == ast.TypeBlock {
		node = node.FirstChild()
	}
	return node
}

func findLastBlockChild(node ast.Node) ast.Node {
	if node == nil {
		return nil
	}

	if !node.HasChildren() || node.LastChild().Type() != ast.TypeBlock {
		return node
	}

	node = node.LastChild()
	for node.LastChild() != nil && node.LastChild().Type() == ast.TypeBlock {
		node = node.LastChild()
	}
	return node
}

func (b *MarkdownBlock) Content() string {
	if b.content == "" {
		startNode := findFirstBlockChild(b.inner)
		lastNode := findLastBlockChild(b.inner)

		start := startNode.Lines().At(0).Start - startNode.Lines().At(0).Padding

		switch b.inner.Kind() {
		case ast.KindBlockquote:
			start -= 2
		}

		end := lastNode.Lines().At(lastNode.Lines().Len() - 1).Stop

		b.content = string(b.source[start:end])
	}
	return b.content
}

func (b *MarkdownBlock) MarshalJSON() ([]byte, error) {
	type markdownBlock struct {
		Markdown string `json:"markdown,omitempty"`
	}

	block := markdownBlock{
		Markdown: b.Content(),
	}

	return json.Marshal(block)
}

type Block interface {
	json.Marshaler
}

type Blocks []Block

func (b Blocks) CodeBlocks() (result CodeBlocks) {
	for _, block := range b {
		if v, ok := block.(*CodeBlock); ok {
			result = append(result, v)
		}
	}
	return
}

type CodeBlocks []*CodeBlock

func (b CodeBlocks) Lookup(name string) *CodeBlock {
	for _, block := range b {
		if block.Name() == name {
			return block
		}
	}
	return nil
}

func (b CodeBlocks) Names() (result []string) {
	for _, block := range b {
		result = append(result, block.Name())
	}
	return result
}
