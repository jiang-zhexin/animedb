package model

import "github.com/jiang-zhexin/animedb/internal/bangumi"

type BangumiCache = map[bangumi.SubjectID]*bangumi.Subject

type SeriesNameToSubjectID = map[string]bangumi.SubjectID
