package entities

import "time"

type DiagnosticRange struct {
	Filename string `json:"filename"`
	Start    struct {
		Line   int `json:"line"`
		Column int `json:"column"`
		Byte   int `json:"byte"`
	} `json:"start"`
	End struct {
		Line   int `json:"line"`
		Column int `json:"column"`
		Byte   int `json:"byte"`
	} `json:"end"`
}

type Snippet struct {
	Context              string   `json:"context"`
	Code                 string   `json:"code"`
	StartLine            int      `json:"start_line"`
	HighlightStartOffset int      `json:"highlight_start_offset"`
	HighlightEndOffset   int      `json:"highlight_end_offset"`
	Values               []string `json:"values"`
}

type Diagnostic struct {
	Severity string          `json:"severity"`
	Summary  string          `json:"summary"`
	Detail   string          `json:"detail"`
	Address  string          `json:"address"`
	Range    DiagnosticRange `json:"range"`
	Snippet  Snippet         `json:"snippet"`
}

type IacTerraformDiagnosticJson struct {
	Level      string     `json:"@level"`
	Message    string     `json:"@message"`
	Module     string     `json:"@module"`
	Timestamp  time.Time  `json:"@timestamp"`
	Diagnostic Diagnostic `json:"diagnostic"`
	Type       string     `json:"type"`
}
