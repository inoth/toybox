package schedule

import (
	"github.com/robfig/cron/v3"
)

type ExecuteFunc func(string)

type ScheduleMission struct {
	// 任务id
	MissionId string
	// 执行id
	entryID cron.EntryID
	// 定时
	Spec string

	ExecuteBefore ExecuteFunc
	// 执行函数
	Execute ExecuteFunc

	ExecuteAfter ExecuteFunc
}

func NewMission(missionId, spec string, execute, executeBefore, executeAfter ExecuteFunc) *ScheduleMission {
	return &ScheduleMission{
		MissionId:     missionId,
		Spec:          spec,
		Execute:       execute,
		ExecuteBefore: executeBefore,
		ExecuteAfter:  executeAfter,
	}
}

func (sm *ScheduleMission) getEntryID() cron.EntryID {
	return sm.entryID
}

func (sm *ScheduleMission) exec() {
	go func(missionId string) {
		if sm.ExecuteBefore != nil {
			sm.ExecuteBefore(missionId)
		}
		sm.Execute(missionId)
		if sm.ExecuteAfter != nil {
			sm.ExecuteAfter(missionId)
		}
	}(sm.MissionId)
}
