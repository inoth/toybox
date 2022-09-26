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
	// 执行函数
	Execute       ExecuteFunc
	ExecuteBefore ExecuteFunc
	ExecuteAfter  ExecuteFunc
}

func NewMission(missionId, spec string, once bool, runner, before, after ExecuteFunc) *ScheduleMission {
	s := &ScheduleMission{
		MissionId:     missionId,
		Spec:          spec,
		once:          once,
		Execute:       runner,
		ExecuteBefore: before,
		ExecuteAfter:  after,
	}
	if once {
		s.ExecuteAfter = func(missionId string) {
			if after != nil {
				after(missionId)
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
