package schedule

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var Schedule *ScheduleServer

type ScheduleServer struct {
	m          sync.RWMutex
	cron       *cron.Cron
	missionMap map[string]ScheduleMission
}

func (s *ScheduleServer) Init() error {
	s.m = sync.RWMutex{}
	s.cron = cron.New(cron.WithLocation(&time.Location{}))
	s.missionMap = make(map[string]ScheduleMission)
	return nil
}

func (s *ScheduleServer) Start() (err error) {
	for _, mission := range s.missionMap {
		mission.entryID, err = s.cron.AddFunc(mission.Spec, func() {
			mission.exec()
		})
	}
	s.cron.Run()
	return
}

// 服务启动前添加
func (s *ScheduleServer) AddMission(mission ScheduleMission) (err error) {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.missionMap[mission.MissionId]; !ok {
		s.missionMap[mission.MissionId] = mission
		return
	}
	return fmt.Errorf("repeat registration tasks [%v]", mission.MissionId)
}

// 服务启动后追加
func (s *ScheduleServer) AppendMission(mission ScheduleMission) (err error) {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.missionMap[mission.MissionId]; !ok {
		mission.entryID, err = s.cron.AddFunc(mission.Spec, func() {
			mission.exec()
		})
		s.missionMap[mission.MissionId] = mission
		return
	}
	return fmt.Errorf("repeat registration tasks [%v]", mission.MissionId)
}

func (s *ScheduleServer) RemoveMission(missionId string) error {
	s.m.Lock()
	defer s.m.Unlock()
	if mm, ok := s.missionMap[missionId]; ok {
		s.cron.Remove(mm.getEntryID())
		delete(s.missionMap, missionId)
		return nil
	}
	return fmt.Errorf("unknown registration tasks [%v]", missionId)
}
