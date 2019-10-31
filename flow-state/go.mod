module github.com/project-flogo/services/flow-state

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/project-flogo/core v0.9.4-rc.1
	github.com/project-flogo/flow v0.9.4-rc.2
	github.com/stretchr/testify v1.4.0
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.9.1
)

go 1.13

replace github.com/project-flogo/flow => ../../flow
