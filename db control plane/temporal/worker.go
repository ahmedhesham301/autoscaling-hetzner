package temporal

import (
	"log"

	"go.temporal.io/sdk/worker"
)

func StartWorker() {
	// The client and worker are heavyweight objects that should be created once per process.
	w := worker.New(TemporalClient, "task-queue", worker.Options{})

	// w.RegisterWorkflow(helloworld.Workflow)
	w.RegisterWorkflow(CreateServiceWorkflow)
	w.RegisterActivity(checkImageExist)
	w.RegisterActivity(buildImage)
	w.RegisterActivity(deployDB)
	w.RegisterActivity(GetOrCreateAllowAllFirewall)
	err := w.Run(nil)
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}
