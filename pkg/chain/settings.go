package chain

type Settings struct {
	RunParallelConditionEvaluations   bool
	ContinueFailedConditionEvaluation bool
	CustomScriptTimeout               int
}

func DefaultSettings() *Settings {
	return &Settings{
		RunParallelConditionEvaluations:   false,
		ContinueFailedConditionEvaluation: false,
		CustomScriptTimeout:               30,
	}
}

func NewSettings(override func(*Settings)) *Settings {
	s := DefaultSettings()

	if override != nil {
		override(s)
	}
	return s
}
