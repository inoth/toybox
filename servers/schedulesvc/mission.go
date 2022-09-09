package schedulesvc

import (
	"github.com/robfig/cron/v3"
)

type MissionRunner interface {
	MissionRunBefore(string)
	MissionRun(string)
	MissionRunAfter(string)
}

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
	Execute MissionRunner
}

func NewMission(missionId, spec string, once bool, runner MissionRunner) *ScheduleMission {
	s := &ScheduleMission{
		MissionId: missionId,
		Execute:   runner,
		Spec:      spec,
		once:      once,
	}
	return s
}

func (sm *ScheduleMission) getEntryID() cron.EntryID {
	return sm.entryID
}

func (sm *ScheduleMission) exec() {
	go func(missionId string) {
		sm.Execute.MissionRunBefore(missionId)
		sm.Execute.MissionRun(missionId)
		sm.Execute.MissionRunAfter(missionId)
		if sm.once {
			Schedule.RemoveMission(missionId)
		}
	}(sm.MissionId)
}
