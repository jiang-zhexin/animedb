package model

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jiang-zhexin/animedb/internal/bangumi"
)

type Model struct {
	SeriesNameToSubjectID SeriesNameToSubjectID    `json:"series_name_to_subject_id"`
	ProcessedFiles        map[string]ProcessedFile `json:"processed_files"`
	BangumiCache          BangumiCache             `json:"bangumi_cache"`
	SourceDirPath         string                   `json:"-"`
	TargetDirPath         string                   `json:"-"`
	DryRun                bool                     `json:"-"`
}

func LoadFromFile(path string) (*Model, error) {
	m := Model{
		SeriesNameToSubjectID: make(SeriesNameToSubjectID),
		BangumiCache:          make(BangumiCache),
		ProcessedFiles:        make(map[string]ProcessedFile),
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &m, nil
		}
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}
	if len(data) == 0 {
		return &m, nil
	}

	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	return &m, nil
}

func (m *Model) SaveToFile(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}
	return nil
}

func (m *Model) WalkFile(clear bool) (sourcePath []string, err error) {
	err = filepath.WalkDir(m.SourceDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := filepath.Ext(path)
			if ext != ".mp4" && ext != ".mkv" {
				return nil
			}
			rel, err := filepath.Rel(m.SourceDirPath, path)
			if err != nil {
				return err
			}
			sourcePath = append(sourcePath, rel)
		}
		return nil
	})
	if clear {
		m.ClearProcessedFiles(sourcePath)
	}
	return
}

func (m *Model) IsProcessed(sourcePath string) bool {
	_, ok := m.ProcessedFiles[sourcePath]
	return ok
}

func (m *Model) ClearProcessedFiles(sourcePath []string) {
	allowed := make(map[string]struct{}, len(sourcePath))
	for _, p := range sourcePath {
		allowed[p] = struct{}{}
	}

	for k := range m.ProcessedFiles {
		if _, ok := allowed[k]; !ok {
			delete(m.ProcessedFiles, k)
		}
	}
}

func (m *Model) UpdateProcessedFiles(processedFiles *ProcessedFile) (err error) {
	pf, ok := m.ProcessedFiles[processedFiles.SourcePath]

	if pf.TargetPath == processedFiles.TargetPath {
		slog.Info("exit")
		return
	}
	slog.Info(
		"create link",
		slog.String("src", processedFiles.SourcePath),
		slog.String("dst", processedFiles.TargetPath),
	)
	if m.DryRun {
		return
	}

	sourcePath := filepath.Join(m.SourceDirPath, processedFiles.SourcePath)
	targetPath := filepath.Join(m.TargetDirPath, processedFiles.TargetPath)
	if err = os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return
	}
	if _, err = os.Stat(targetPath); !os.IsNotExist(err) {
		err = os.Remove(targetPath)
		if err != nil {
			slog.Error("can not removing file", slog.String("err", err.Error()))
			return
		}
	}
	if err = os.Link(sourcePath, targetPath); err != nil {
		return
	}

	if ok {
		slog.Warn("remove file", slog.String("path", pf.TargetPath))
		if m.DryRun {
			return
		}

		oldfile := filepath.Join(m.TargetDirPath, pf.TargetPath)
		if err = os.Remove(oldfile); err != nil && !os.IsNotExist(err) {
			return
		}

		nfo := strings.TrimSuffix(oldfile, filepath.Ext(oldfile)) + ".nfo"
		if err = os.Remove(nfo); err != nil && !os.IsNotExist(err) {
			return
		}

		dir := filepath.Dir(oldfile)
		if err = removeIfEmpty(dir); err != nil && !os.IsNotExist(err) {
			return
		}
	}

	m.ProcessedFiles[processedFiles.SourcePath] = *processedFiles
	return
}

func (m *Model) GetSubjectById(id bangumi.SubjectID) (*bangumi.Subject, error) {
	s, ok := m.BangumiCache[id]
	if ok {
		return s, nil
	}

	s, err := bangumi.GetSubject(id)
	if err != nil {
		return nil, err
	}

	m.BangumiCache[id] = s
	return s, nil
}

func (m *Model) GetSubjectBySeriesName(seriesName string) (*bangumi.Subject, error) {
	id, ok := m.SeriesNameToSubjectID[seriesName]
	if ok {
		return m.GetSubjectById(id)
	}

	results, err := m.SearchSubject(seriesName)
	if err != nil {
		return nil, err
	}
	s := &results[0]
	m.SeriesNameToSubjectID[seriesName] = s.Id
	return s, nil
}

func (m *Model) SearchSubject(seriesName string) ([]bangumi.Subject, error) {
	results, err := bangumi.SearchSubject(seriesName)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		m.BangumiCache[r.Id] = &r
	}
	return results, nil
}

func (m *Model) FindProcessedFilesBySeriesName(seriesName string) (result []ProcessedFile) {
	for _, pf := range m.ProcessedFiles {
		if pf.SeriesName == seriesName {
			result = append(result, pf)
		}
	}
	return
}

func (m *Model) UpdateSeriesNameToSubjectID(seriesName string, id bangumi.SubjectID) error {
	m.SeriesNameToSubjectID[seriesName] = id
	s, err := m.GetSubjectById(id)
	if err != nil {
		return err
	}

	ProcesseFiles := m.FindProcessedFilesBySeriesName(seriesName)
	if len(ProcesseFiles) == 0 {
		return nil
	}

	for _, pf := range ProcesseFiles {
		pf.UpdateTargetPath(s)
		if err = m.UpdateProcessedFiles(&pf); err != nil {
			return err
		}
	}
	return nil
}

func removeIfEmpty(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(entries) == 0 {
		return os.Remove(dir)
	}

	return nil
}
