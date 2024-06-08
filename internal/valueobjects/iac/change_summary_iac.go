package iac

var (
	EmptySummary = ChangeSummaryIac{}
)

type ChangeSummaryIac struct {
	Add    int
	Change int
	Remove int
}
