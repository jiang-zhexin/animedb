package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jiang-zhexin/animedb/internal/bangumi"
	"github.com/jiang-zhexin/animedb/internal/parser"
)

type ProcessedFile struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
	*parser.ParseResult
}

func NewProcessedFiles(sourcePath string) (*ProcessedFile, error) {
	p := parser.NewParser()
	filename := filepath.Base(sourcePath)

	pr, err := p.Parse(strings.TrimSuffix(filename, filepath.Ext(filename)))
	if err != nil {
		return nil, err
	}
	return &ProcessedFile{
		SourcePath:  sourcePath,
		ParseResult: pr,
	}, nil
}

func (pf *ProcessedFile) UpdateTargetPath(s *bangumi.Subject) {
	name := s.NameCn
	if name == "" {
		name = s.Name
	}

	pf.TargetPath = filepath.Join(name, fmt.Sprintf("S01E%02d.mkv", pf.Episode))
}
