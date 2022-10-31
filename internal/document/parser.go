package document

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type ParsedSource struct {
	data []byte
	root ast.Node
}

func (s *ParsedSource) Root() ast.Node {
	return s.root
}

func (s *ParsedSource) Source() []byte {
	return s.data
}

func (s *ParsedSource) hasChildOfKind(node ast.Node, kind ast.NodeKind) (ast.Node, bool) {
	if node.Type() != ast.TypeBlock {
		return nil, false
	}

	if node.Kind() == kind {
		return node, true
	}

	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		if node, ok := s.hasChildOfKind(c, kind); ok {
			return node, ok
		}
	}
	return nil, false
}

func (s *ParsedSource) findBlocks(nameRes *nameResolver, docNode ast.Node) (result Blocks) {
	for c := docNode.FirstChild(); c != nil; c = c.NextSibling() {
		switch c.Kind() {
		case ast.KindFencedCodeBlock:
			result = append(result, &CodeBlock{
				source:       s.data,
				inner:        c.(*ast.FencedCodeBlock),
				nameResolver: nameRes,
			})
		default:
			if innerCodeBlock, ok := s.hasChildOfKind(c, ast.KindFencedCodeBlock); ok {
				fmt.Printf("found inner fenced code block\n")

				switch c.Kind() {
				case ast.KindList:
					listItem := innerCodeBlock.Parent()

					// move the code block into the root node
					listItem.RemoveChild(listItem, innerCodeBlock)
					docNode.InsertAfter(docNode, c, innerCodeBlock)

					// split the list if there are any list items
					// after listItem
					if listItem.NextSibling() != nil {
						newList := ast.NewList(c.(*ast.List).Marker)
						for item := listItem.NextSibling(); item != nil; item = item.NextSibling() {
							c.RemoveChild(c, item)
							newList.AppendChild(newList, item)
						}
						docNode.InsertAfter(docNode, innerCodeBlock, newList)
					}
				case ast.KindBlockquote:
					nextParagraph := innerCodeBlock.NextSibling()

					// move the code block into the root node
					c.RemoveChild(c, innerCodeBlock)
					docNode.InsertAfter(docNode, c, innerCodeBlock)

					// move all paragraphs after the code block
					// into the new block quote
					if nextParagraph != nil {
						newBlockQuote := ast.NewBlockquote()
						for item := nextParagraph; item != nil; item = item.NextSibling() {
							c.RemoveChild(c, item)
							newBlockQuote.AppendChild(newBlockQuote, item)
						}
						docNode.InsertAfter(docNode, innerCodeBlock, newBlockQuote)
					}
				}
			}

			result = append(result, &MarkdownBlock{
				source: s.data,
				inner:  c,
			})
		}
	}
	return
}

func (s *ParsedSource) Blocks() Blocks {
	nameRes := &nameResolver{
		namesCounter: map[string]int{},
		cache:        map[interface{}]string{},
	}
	return s.findBlocks(nameRes, s.root)
}

func (s *ParsedSource) CodeBlocks() CodeBlocks {
	var result CodeBlocks

	nameRes := &nameResolver{
		namesCounter: map[string]int{},
		cache:        map[interface{}]string{},
	}

	// TODO(adamb): check the case when a paragraph is immediately
	// followed by a code block without a new line separating them.
	// Currently, such a code block is not detected at all.
	for c := s.root.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindFencedCodeBlock {
			result = append(result, &CodeBlock{
				source:       s.data,
				inner:        c.(*ast.FencedCodeBlock),
				nameResolver: nameRes,
			})
		}
	}

	return result
}

func getRange(source []byte, n ast.Node, start int, stop int) (string, error) {
	var content strings.Builder
	switch n.Kind() {
	case ast.KindHeading:
		heading := n.(*ast.Heading)
		offset := 1 + heading.Level
		// shield from inital ======= vs ### heading
		if start-offset < 0 {
			offset = 0
		}
		_, _ = content.Write(source[start-offset : stop])
	default:
		_, _ = content.Write(source[start:stop])
	}
	return content.String(), nil
}

func getPrevStart(n ast.Node) int {
	curr := n
	prev := n.PreviousSibling()
	if prev != nil && prev.Lines().Len() > 0 {
		curr = prev
		l := curr.Lines().Len()
		return curr.Lines().At(l - 1).Stop
	}
	return curr.Lines().At(0).Stop
}

func getNextStop(n ast.Node) int {
	stop := 0
	curr := n

	if n.Kind() != ast.KindDocument {
		next := curr.NextSibling()

		if next != nil && next.Lines().Len() > 0 {
			curr = next
			stop = curr.Lines().At(0).Start
		} else {
			l := curr.Lines().Len()
			stop = curr.Lines().At(l - 1).Start
		}
	}

	// add back markdown heading levels
	if stop > 0 && curr.Kind() == ast.KindHeading {
		heading := curr.(*ast.Heading)
		// simple math to add back ## incl trailing space
		stop = stop - 1 - heading.Level
	}

	return stop
}

func (s *ParsedSource) SquashedBlocks() (blocks Blocks, err error) {
	nameRes := &nameResolver{
		namesCounter: map[string]int{},
		cache:        map[interface{}]string{},
	}

	lastCodeBlock := s.root
	remainingNode := s.root

	err = ast.Walk(s.root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() != ast.KindFencedCodeBlock || !entering {
			return ast.WalkContinue, nil
		}

		start := getPrevStart(n)

		if lastCodeBlock != nil {
			stop := getNextStop(lastCodeBlock)
			// check for existence of markdown in between code blocks
			if start > stop {
				markdown, _ := getRange(s.data, n, stop, start)
				blocks = append(blocks, &MarkdownBlock{content: markdown})
			}
			lastCodeBlock = n
		}

		blocks = append(blocks, &CodeBlock{
			source:       s.data,
			inner:        n.(*ast.FencedCodeBlock),
			nameResolver: nameRes,
		})

		remainingNode = n.NextSibling()

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Never encounter a code block, stuck on document node
	if remainingNode == s.root {
		remainingNode = remainingNode.FirstChild()
	}

	// Skip remainingNodes unless it's got lines
	for remainingNode != nil && remainingNode.Lines().Len() == 0 {
		remainingNode = remainingNode.NextSibling()
	}

	if remainingNode != nil {
		start := remainingNode.Lines().At(0).Start
		stop := len(s.data) - 1
		markdown, _ := getRange(s.data, remainingNode, start, stop)
		blocks = append(blocks, &MarkdownBlock{content: markdown})
	}

	// Handle a single code block
	if len(blocks) == 2 {
		if b, ok := blocks[0].(*MarkdownBlock); ok && strings.HasPrefix(b.content, "```") {
			blocks = blocks[1:]
		}
	}

	return
}

type defaultParser struct {
	parser parser.Parser
}

func newDefaultParser() *defaultParser {
	return &defaultParser{parser: goldmark.DefaultParser()}
}

func (p *defaultParser) Parse(data []byte) *ParsedSource {
	root := p.parser.Parse(text.NewReader(data))
	return &ParsedSource{data: data, root: root}
}

type nameResolver struct {
	namesCounter map[string]int
	cache        map[interface{}]string
}

func (r *nameResolver) Get(obj interface{}, name string) string {
	if v, ok := r.cache[obj]; ok {
		return v
	}

	var result string

	r.namesCounter[name]++

	if r.namesCounter[name] == 1 {
		result = name
	} else {
		result = fmt.Sprintf("%s-%d", name, r.namesCounter[name])
	}

	r.cache[obj] = result

	return result
}
