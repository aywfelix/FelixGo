package eventsystem

import (
	"fmt"
)

type IEventObserver interface {
	Update()
	SetSubject(subject IEventSubject)
}

type EventObserver struct {
	subject IEventSubject
}

func NewEventObserver() IEventObserver {
	return &EventObserver{}
}

func (e *EventObserver) SetSubject(subject IEventSubject) {
	e.subject = subject
}

func (e *EventObserver) Update() {

}

type EnemyKilledUI struct {
	EventObserver
}

func NewEnemyKilledUI() IEventObserver {
	return &EnemyKilledUI{}
}

func (e *EnemyKilledUI) Update() {
	fmt.Println("enemy killed count=", e.subject.GetCount())
}

// 某一类任务
type EnemyKilledTask struct {
	EventObserver
	// taskSystem TaskSystem
}

func NewEnemyKilledTask() IEventObserver {
	return &EnemyKilledTask{}
}

func (e *EnemyKilledTask) Update() {
	params := e.subject.GetParams()
	if params == nil && len(params) == 0 {
		return
	}
	if params[0].(string) == "小怪" {
		// e.taskCount++
		// if e.taskCount >= 5 {
		// 	// 完成任务
		// }
	}
}
