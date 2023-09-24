package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func stageWrapper(done In, in In, stage Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}

				out <- val
			}
		}
	}()

	return stage(out)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	curStage := stageWrapper(done, in, func(in In) Out { return in })

	for _, stage := range stages {
		curStage = stageWrapper(done, curStage, stage)
	}

	return curStage
}
