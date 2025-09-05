package entity

type Scenario string

const (
	ScenarioDefault   Scenario = "default"
	ScenarioEvaluator Scenario = "evaluator"
)

func ScenarioValue(scenario *Scenario) Scenario {
	if scenario == nil {
		return ScenarioDefault
	}
	return *scenario
}
