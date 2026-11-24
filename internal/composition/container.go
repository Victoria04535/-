package composition

import (
	"io"

	"pos/internal/adapter/console"
	"pos/internal/application"
	"pos/internal/domain"
	csvexporter "pos/internal/infrastructure/csv"
	"pos/internal/infrastructure/memory"
)

type Container struct {
	Runner *console.Console
}

func NewContainer(input io.Reader, output io.Writer) *Container {

	studentRepository := memory.NewStudentRepository()
	criterionRepository := memory.NewCriterionRepository()
	evaluationRepository := memory.NewEvaluationRepository()

	exporter := csvexporter.NewExporter()

	strategy := domain.WeightedAverageStrategy{}
	calculator := domain.NewTotalScoreService(strategy)

	facade := application.NewPOSFacade(
		studentRepository,
		criterionRepository,
		evaluationRepository,
		calculator,
		exporter,
	)

	consoleAdapter := console.NewConsole(
		facade,
		input,
		output,
	)

	return &Container{
		Runner: consoleAdapter,
	}
}
