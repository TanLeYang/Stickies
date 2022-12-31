package command

type CommandStage string

// Shared stages
const (
	Start CommandStage = "START"
	End   CommandStage = "END"
)

func nextStage(stages []CommandStage, currentStage CommandStage, repeats bool) CommandStage {
	currentStageIdx := -1
	for i, s := range stages {
		if s == currentStage {
			currentStageIdx = i
			break
		}
	}

	nextStage := currentStageIdx + 1
	if currentStageIdx == -1 || (nextStage == len(stages) && !repeats) {
		return End
	} else {
		return stages[nextStage%len(stages)]
	}
}
