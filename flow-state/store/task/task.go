package task

import (
	"github.com/project-flogo/flow/model"
	"github.com/project-flogo/flow/state"
	flowEvent "github.com/project-flogo/flow/support/event"

	"sort"
	"strings"

	"github.com/project-flogo/flow/state/change"
)

type Task struct {
	Id         string                 `json:"id"`
	StepId     int64                  `json:"stepId"`
	SubflowId  int                    `json:"subflowId"`
	Input      map[string]interface{} `json:"input,omitempty"`
	Output     map[string]interface{} `json:"output,omitempty"`
	Links      []*Link                `json:"links,omitempty"`
	Status     flowEvent.Status       `json:"status"`
	FlowStatus flowEvent.Status       `json:"flow_status"`
	StartTask  bool                   `json:"startTask"`
	//For subflow
	NewSubflow bool   `json:"newSubflow, omitempty"`
	Flowname   string `json:"flowname"`
}

type Link struct {
	From           string           `json:"from"`
	To             string           `json:"to"`
	ToTaskStutus   flowEvent.Status `json:"toTaskStutus"`
	FromTaskStutus flowEvent.Status `json:"fromTaskStutus"`
	Status         string           `json:"status"`
}

func StepToTask(step *state.Step) ([]*Task, error) {
	return StepToTaskReady(step, false)
}

func StepToTaskWithReadyTaskInput(step *state.Step) ([]*Task, error) {
	return StepToTaskReady(step, true)
}

func StepToTaskReady(step *state.Step, incluedReadyInput bool) ([]*Task, error) {
	if len(step.FlowChanges) > 1 {
		var tasks []*Task
		//1. Subflow started
		//2. Subflow end
		//Subflow activity + new flow start task
		var subflowIds []int
		for k, _ := range step.FlowChanges {
			subflowIds = append(subflowIds, k)
		}

		sort.Ints(subflowIds)

		parentFlowId := subflowIds[0]
		parentTask, err := flowChangeToTask(step.Id, parentFlowId, step.FlowChanges[parentFlowId], incluedReadyInput)
		if err != nil {
			return nil, err
		}

		subflowID := subflowIds[1]
		subflowTask, err := flowChangeToTask(step.Id, subflowID, step.FlowChanges[subflowID], incluedReadyInput)
		if err != nil {
			return nil, err
		}
		subflowTask.NewSubflow = step.FlowChanges[subflowID].NewFlow
		if subflowTask.NewSubflow {
			subflowTask.Id = ""
			tasks = append(tasks, parentTask, subflowTask)

		} else {
			tasks = append(tasks, subflowTask, parentTask)

		}

		return tasks, nil
	} else {
		var tasks []*Task
		for subflowId, flowChange := range step.FlowChanges {
			tak, err := flowChangeToTask(step.Id, subflowId, flowChange, incluedReadyInput)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, tak)
		}
		return tasks, nil
	}
	return nil, nil
}

func flowChangeToTask(stepId int64, subflowId int, flowChange *change.Flow, incluedReadyInput bool) (*Task, error) {
	task := &Task{StepId: stepId}
	task.SubflowId = subflowId
	task.Id = flowChange.TaskId
	task.Flowname = flowChange.FlowURI
	if flowChange.NewFlow {
		task.Output = attrToOutput(flowChange.Attrs)
		task.StartTask = true
		task.Status = flowEvent.STARTED
		task.FlowStatus = convertFlowStatus(flowChange.Status)
		for taskId, t := range flowChange.Tasks {
			if incluedReadyInput {
				if t.Status == 40 || t.Status == 30 || t.Status == 100 || t.Status == 20 {
					task.Id = taskId
					task.Input = t.Input
					task.Status = convertTaskStatus(t.Status)
				}
			} else {
				if t.Status == 40 || t.Status == 30 || t.Status == 100 {
					task.Id = taskId
					task.Input = t.Input
					task.Status = convertTaskStatus(t.Status)
				}
			}

		}
	} else {
		task.FlowStatus = convertFlowStatus(flowChange.Status)
		task.Output = attrToOutput(flowChange.Attrs)

		if flowChange.ReturnData != nil && len(flowChange.ReturnData) > 0 {
			task.Input = flowChange.ReturnData
		}

		for taskId, t := range flowChange.Tasks {
			if incluedReadyInput && task.Id == taskId {
				if t.Status == 40 || t.Status == 30 || t.Status == 100 || t.Status == 20 {
					task.Id = taskId
					if len(t.Input) > 0 {
						task.Input = t.Input
					}
					task.Status = convertTaskStatus(t.Status)
				}
			} else {
				if t.Status == 40 || t.Status == 30 || t.Status == 100 {
					task.Id = taskId
					if len(t.Input) > 0 {
						task.Input = t.Input
					}
					task.Status = convertTaskStatus(t.Status)
				}
			}

		}

		taskId := task.Id
		for _, lik := range flowChange.Links {
			if lik.Status == 0 && (lik.From == "" || lik.To == "") {
				continue
			}
			if taskId == lik.From {
				fromTask := getTaskByTaskId(flowChange.Tasks, lik.From)
				toTask := getTaskByTaskId(flowChange.Tasks, lik.To)
				var fromTaskStatus flowEvent.Status
				if fromTask != nil {
					fromTaskStatus = convertTaskStatus(fromTask.Status)
				}

				var toTaskStatus flowEvent.Status
				if toTask != nil {
					toTaskStatus = convertTaskStatus(toTask.Status)
				}
				link := &Link{
					From:           lik.From,
					FromTaskStutus: fromTaskStatus,
					ToTaskStutus:   toTaskStatus,
					To:             lik.To,
					Status:         LinkStatusToString(lik.Status),
				}
				task.Links = append(task.Links, link)
			}
		}
	}
	return task, nil
}

func getTaskByTaskId(tasks map[string]*change.Task, taskId string) *change.Task {
	for k, v := range tasks {
		if k == taskId {
			return v
		}
	}
	return nil
}

func attrToOutput(attrs map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range attrs {
		if strings.HasPrefix(k, "_A.") {
			fieldName := k[3:]
			index := strings.Index(fieldName, ".")
			if index > 0 {
				fieldName = fieldName[index+1:]
			}
			output[fieldName] = v
		} else {
			output[k] = v
		}
	}
	return output
}

func LinkStatusToString(flowStatus int) string {
	switch flowStatus {
	case 1:
		return "false"
	case 2:
		return "ture"
	case 3:
		return "skipped"
	}
	return "unknow"
}

func convertTaskStatus(code int) flowEvent.Status {
	switch model.TaskStatus(code) {
	case model.TaskStatusNotStarted:
		return flowEvent.CREATED
	case model.TaskStatusEntered:
		return flowEvent.SCHEDULED
	case model.TaskStatusSkipped:
		return flowEvent.SKIPPED
	case model.TaskStatusReady:
		return flowEvent.STARTED
	case model.TaskStatusFailed:
		return flowEvent.FAILED
	case model.TaskStatusDone:
		return flowEvent.COMPLETED
	case model.TaskStatusWaiting:
		return flowEvent.WAITING
	}
	return flowEvent.UNKNOWN
}

func convertFlowStatus(code int) flowEvent.Status {
	switch model.FlowStatus(code) {
	case model.FlowStatusNotStarted:
		return flowEvent.CREATED
	case model.FlowStatusActive:
		return flowEvent.STARTED
	case model.FlowStatusCancelled:
		return flowEvent.CANCELLED
	case model.FlowStatusCompleted:
		return flowEvent.COMPLETED
	case model.FlowStatusFailed:
		return flowEvent.FAILED
	}
	return flowEvent.UNKNOWN
}
