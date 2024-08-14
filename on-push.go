package main

import "github.com/shykes/gha/internal/dagger"

// Add a trigger to execute a Dagger pipeline on a git push
func (m *Gha) OnPush(
	// The Dagger command to execute
	// Example 'build --source=.'
	command string,
	// The Dagger module to load
	// +optional
	// +default="."
	module string,
	// Github secrets to inject into the pipeline environment.
	// For each secret, an env variable with the same name is created.
	// Example: ["PROD_DEPLOY_TOKEN", "PRIVATE_SSH_KEY"]
	// +optional
	secrets []string,
	// Run only on push to specific branches
	// +optional
	branches []string,
	// Run only on push to specific branches
	// +optional
	tags []string,
	// Run only on push to specific paths
	// +optional
	paths []string,
	// Dispatch jobs to the given runner
	// +optional
	runner string,
) *Gha {
	if err := validateSecretNames(secrets); err != nil {
		panic(err) // FIXME
	}
	m.PushTriggers = append(m.PushTriggers, PushTrigger{
		Event: PushEvent{
			Branches: branches,
			Tags:     tags,
			Paths:    paths,
		},
		Pipeline: m.pipeline(command, module, runner, secrets),
	})
	return m
}

type PushTrigger struct {
	Event    PushEvent
	Pipeline Pipeline
}

func (t PushTrigger) asWorkflow() Workflow {
	var workflow = t.Pipeline.asWorkflow()
	workflow.On = WorkflowTriggers{Push: &(t.Event)}
	return workflow
}

func (t PushTrigger) Config(filename string, asJson bool) *dagger.Directory {
	return t.asWorkflow().Config(filename, asJson)
}