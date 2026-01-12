package openrouter

type Effort string

const (
	EffortXHigh   Effort = "xhigh"
	EffortHigh    Effort = "high"
	EffortMedium  Effort = "medium"
	EffortLow     Effort = "low"
	EffortMinimal Effort = "minimal"
	EffortNone    Effort = "none"
)

type Summary string

const (
	SummaryAuto     Summary = "auto"
	SummaryConcise  Summary = "concise"
	SummaryDetailed Summary = "detailed"
)

type Reasoning struct {
	Effort  Effort  `json:"effort,omitempty"`
	Summary Summary `json:"summary,omitempty"`
}
