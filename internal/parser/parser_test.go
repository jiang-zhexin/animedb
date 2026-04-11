package parser

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		wantSeriesName string
		wantEpisode    int
		wantGroup      string
	}{
		{
			name:           "ANi Re:Zero S3",
			filename:       "[ANi] Re：從零開始的異世界生活 第三季 - 01v2 [1080P][Baha][WEB-DL][AAC AVC][CHT]",
			wantSeriesName: "Re：從零開始的異世界生活 第三季",
			wantEpisode:    1,
			wantGroup:      "ANi",
		},
		{
			name:           "ANi Dandadan",
			filename:       "[ANi] 膽大黨 - 01 [1080P][Baha][WEB-DL][AAC AVC][CHT]",
			wantSeriesName: "膽大黨",
			wantEpisode:    1,
			wantGroup:      "ANi",
		},
		{
			name:           "ANi Maou 2099",
			filename:       "[ANi] 魔王 2099 - 02 [1080P][Baha][WEB-DL][AAC AVC][CHT]",
			wantSeriesName: "魔王 2099",
			wantEpisode:    2,
			wantGroup:      "ANi",
		},
		{
			name:           "ANi Puniru",
			filename:       "[ANi] 噗妮露是可愛史萊姆 - 01 [1080P][Baha][WEB-DL][AAC AVC][CHT]",
			wantSeriesName: "噗妮露是可愛史萊姆",
			wantEpisode:    1,
			wantGroup:      "ANi",
		},
		{
			name:           "LoliHouse Maou 2099",
			filename:       "[LoliHouse] Maou 2099 - 01 [WebRip 1080p HEVC-10bit AAC SRTx2]",
			wantSeriesName: "Maou 2099",
			wantEpisode:    1,
			wantGroup:      "LoliHouse",
		},
		{
			name:           "Nekomoe Ao no Hako",
			filename:       "[Nekomoe kissaten&LoliHouse] Ao no Hako - 04v2 [WebRip 1080p HEVC-10bit AAC ASSx2]",
			wantSeriesName: "Ao no Hako",
			wantEpisode:    4,
			wantGroup:      "Nekomoe kissaten&LoliHouse",
		},
		{
			name:           "Tsukigakirei Ao no Hako",
			filename:       "[Tsukigakirei][Ao no Hako][07][WEBrip][1080P][CHS&JPN]",
			wantSeriesName: "Ao no Hako",
			wantEpisode:    7,
			wantGroup:      "Tsukigakirei",
		},
		{
			name:           "KitaujiSub Chi. Chikyuu no Undou ni Tsuite",
			filename:       "[KitaujiSub&LoliHouse] Chi. Chikyuu no Undou ni Tsuite - 01 [WebRip 1080p HEVC-10bit AAC ASSx2]",
			wantSeriesName: "Chi. Chikyuu no Undou ni Tsuite",
			wantEpisode:    1,
			wantGroup:      "KitaujiSub&LoliHouse",
		},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.filename)

			if err != nil {
				t.Error(err.Error())
			}
			if got.SeriesName != tt.wantSeriesName {
				t.Errorf("SeriesName = %q, want %q", got.SeriesName, tt.wantSeriesName)
			}
			if got.Episode != tt.wantEpisode {
				t.Errorf("Episode = %d, want %d", got.Episode, tt.wantEpisode)
			}
			if got.Group != tt.wantGroup {
				t.Errorf("Group = %q, want %q", got.Group, tt.wantGroup)
			}
		})
	}
}
