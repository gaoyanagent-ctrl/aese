package application

import "fmt"

type RunStatus string

const (
	RunStatusPlanned              RunStatus = "planned"
	RunStatusInitializing         RunStatus = "initializing"
	RunStatusReady                RunStatus = "ready"
	RunStatusRunning              RunStatus = "running"
	RunStatusAwaitingAnalysis     RunStatus = "awaiting_analysis"
	RunStatusAnalyzing            RunStatus = "analyzing"
	RunStatusAwaitingVerification RunStatus = "awaiting_verification"
	RunStatusCompleted            RunStatus = "completed"
	RunStatusFailed               RunStatus = "failed"
	RunStatusResetting            RunStatus = "resetting"
	RunStatusReset                RunStatus = "reset"
)

type RunAction string

const (
	RunActionPreflight    RunAction = "preflight"
	RunActionInitialize   RunAction = "initialize"
	RunActionAdvance      RunAction = "advance"
	RunActionRunToEnd     RunAction = "run-to-end"
	RunActionAnalyze      RunAction = "analyze"
	RunActionVerify       RunAction = "verify"
	RunActionResetPlan    RunAction = "reset-plan"
	RunActionReset        RunAction = "reset"
)

type RunTransitionContext struct {
	CurrentAct int
	TotalActs  int
}

func (c RunTransitionContext) actsCompleted() bool {
	return c.CurrentAct >= c.TotalActs
}

// AllowedActions is a stateless decision helper for precondition-free actions.
func AllowedActions(status RunStatus) []RunAction {
	switch status {
	case RunStatusPlanned:
		return []RunAction{RunActionPreflight}
	case RunStatusInitializing:
		return []RunAction{RunActionInitialize, RunActionResetPlan}
	case RunStatusReady, RunStatusRunning:
		return []RunAction{RunActionAdvance, RunActionRunToEnd, RunActionResetPlan}
	case RunStatusAwaitingAnalysis:
		return []RunAction{RunActionAnalyze, RunActionResetPlan}
	case RunStatusAnalyzing:
		return []RunAction{RunActionVerify, RunActionResetPlan}
	case RunStatusAwaitingVerification:
		return []RunAction{RunActionVerify, RunActionResetPlan}
	case RunStatusCompleted:
		return []RunAction{RunActionPreflight, RunActionResetPlan}
	case RunStatusFailed:
		return []RunAction{RunActionResetPlan}
	case RunStatusResetting:
		return []RunAction{RunActionReset}
	case RunStatusReset:
		return []RunAction{RunActionPreflight}
	default:
		return []RunAction{}
	}
}

// NextStatus returns the next status after an action if the transition is legal.
//
// For advance/run-to-end transitions, `context` controls whether this action
// reaches the final act and therefore moves to awaiting_analysis.
func NextStatus(current RunStatus, action RunAction, context RunTransitionContext) (RunStatus, error) {
	switch current {
	case RunStatusPlanned:
		if action == RunActionPreflight {
			return RunStatusInitializing, nil
		}
	case RunStatusInitializing:
		if action == RunActionInitialize {
			return RunStatusReady, nil
		}
    case RunStatusReady, RunStatusRunning:
		if action == RunActionAdvance || action == RunActionRunToEnd {
			if context.actsCompleted() {
				return RunStatusAwaitingAnalysis, nil
			}
			return RunStatusRunning, nil
		}
	case RunStatusAwaitingAnalysis:
		if action == RunActionAnalyze {
			return RunStatusAnalyzing, nil
		}
	case RunStatusAnalyzing:
		if action == RunActionVerify {
			return RunStatusAwaitingVerification, nil
		}
	case RunStatusAwaitingVerification:
		if action == RunActionVerify {
			return RunStatusCompleted, nil
		}
	case RunStatusCompleted, RunStatusFailed:
		if action == RunActionResetPlan {
			return RunStatusResetting, nil
		}
	case RunStatusResetting:
		if action == RunActionReset {
			return RunStatusReset, nil
		}
	case RunStatusReset:
		if action == RunActionPreflight {
			return RunStatusInitializing, nil
		}
	}

	if action == RunActionResetPlan {
		switch current {
		case RunStatusReady, RunStatusRunning, RunStatusAwaitingAnalysis, RunStatusAnalyzing, RunStatusAwaitingVerification:
			return RunStatusResetting, nil
		}
	}

	return "", fmt.Errorf("transition %s -> %s is not allowed", current, action)
}
