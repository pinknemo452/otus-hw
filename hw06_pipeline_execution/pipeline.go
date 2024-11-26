package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	resultCh := make(Bi)
	outChannels := createPipelines(in, stages...)
	go handleDone(done, outChannels, resultCh)
	go collectResultsOrdered(done, outChannels, resultCh)

	return resultCh
}

func createPipelines(in In, stages ...Stage) []Out {
	outChannels := []Out{}
	for data := range in {
		stageInCh := make(Bi)

		outChannels = append(outChannels, chainStages(stageInCh, stages...))
		stageInCh <- data
		close(stageInCh)
	}
	return outChannels
}

func chainStages(dataIn In, stages ...Stage) Out {
	var pipelineOut Out
	prev := dataIn
	for i := range stages {
		prev = stages[i](prev)
	}
	pipelineOut = prev

	return pipelineOut
}

func handleDone(done In, outChannels []Out, resCh Bi) {
	<-done
	close(resCh)
	for _, ch := range outChannels {
		go func(in In) { <-in }(ch)
	}
}

func collectResultsOrdered(done In, outChannels []Out, resCh Bi) {
	for _, ch := range outChannels {
		select {
		case <-done:
			return
		case result := <-ch:
			resCh <- result
		}
	}
	close(resCh)
}
