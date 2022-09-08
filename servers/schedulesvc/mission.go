package schedulesvc

import (
	"github.com/robfig/cron/v3"
)

type ExecuteFunc func(string)

type ScheduleMission struct {
	// 任务id
	MissionId string
	// 执行id
	entryID cron.EntryID
	// 仅执行一次
	once bool
	// 定时 cron 语句
	Spec string
	// 执行前函数
	ExecuteBefore ExecuteFunc
	// 执行函数
	Execute ExecuteFunc
	// 执行后函数
	ExecuteAfter ExecuteFunc
}

func NewMission(missionId, spec string, once bool, execute, executeBefore, executeAfter ExecuteFunc) *ScheduleMission {
	s := &ScheduleMission{
		MissionId:     missionId,
		Execute:       execute,
		Spec:          spec,
		once:          once,
		ExecuteBefore: executeBefore,
		ExecuteAfter:  executeAfter,
	}
	if once {
		s.ExecuteAfter = func(missionId string) {
			if executeAfter != nil {
				executeAfter(missionId)
			}
			Schedule.RemoveMission(missionId)
		}
	}
	return s
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
