package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		inStage := make(Bi)

		go func(in Bi, out Out) {
			defer close(in)

			for {
				select {
				case <-done:
					return
				case v, ok := <-out:
					if !ok {
						return
					}

					in <- v
				}
			}
		}(inStage, out)

		out = stage(inStage)
	}

	return out
}
