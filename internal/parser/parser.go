package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

type Parser struct {
	patterns []*regexp.Regexp
}

type ParseResult struct {
	SeriesName string `json:"series_name"`
	Episode    int    `json:"episode"`
	Group      string `json:"group"`
	RawName    string `json:"raw_name"`
}

func NewParser() *Parser {
	return &Parser{
		patterns: reg,
	}
}

func (p *Parser) Parse(filename string) (*ParseResult, error) {
	result := &ParseResult{
		RawName: filename,
	}

	for _, re := range p.patterns {
		if matches := re.FindStringSubmatch(filename); matches != nil {
			result.Group = matches[1]
			result.SeriesName = matches[2]
			e, _ := strconv.Atoi(matches[3])
			result.Episode = e
			return result, nil
		}
	}

	return nil, fmt.Errorf("all match fail: \"%s\"", filename)
}

var patterns = []string{
	`\[(.+?)]\s*(.+) - (\d+)(?:v\d+)? \[(.+)]`,
	`\[(.+?)]\s*(.+) - s\d+e(\d+)? -? \[(.+)]`,
	`\[(.+?)]\s*(.+) S\d+ \[(\d+)(?:v\d+)?] \[(.+)]`,
	`\[(.+?)]\s*(.+) \[(\d+)(?:v\d+)?]\[(.+)]`,
	`\[(.+?)]\s*(.+) \[(\d+)(?:v\d+)?] \[(.+)]`,
	`\[(.+?)]\s*\[(.+)]\[.+]\[(\d+)(?:v\d+)?](?:\[(.+)])?`,
	`\[(.+?)]\s*\[(.+)]\[(\d+)(?:v\d+)?]\[(.+)]`,
	`\[(.+?)]\s*\[(.+)]\[(\d+)-(?:v\d+)?]\[(.+)]`,
	`\[(.+?)]\s*\[(.+)]\[(\d+) (?:v\d+)?]\[(.+)]`,
	`\[(.+?)]\s*\[(.+) - (\d+)(?:v\d+)?]\[(.+)]`,
}

var reg = make([]*regexp.Regexp, len(patterns))

func init() {
	for i, r := range patterns {
		reg[i] = regexp.MustCompile(r)
	}
}
