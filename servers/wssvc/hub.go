package wssvc

import (
	"context"

	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/parsers"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/process"
)

/*
	websocket 长连接服务
	* 是否启用缓存, 本地 OR Redis
	* 开放消息管道插件, 插件化录入
*/

var ChatHub *Hub

type Hub struct {
	done chan struct{}

	clients         map[string]*Client
	broadcastInput  chan []byte                 // 房间消息输入
	broadcastOutPut chan models.SendMessageBody // 房间消息输出
	register        chan *Client                // 加入通道
	unregister      chan *Client                // 退出通道

	parser  parsers.MsgParsers   // 消息体解析
	process []process.MsgProcess // 消息体处理
}

func (h *Hub) Init() error {
	var hubCfg models.HubConfig
	err := config.Cfg.UnmarshalKey("ChatHub", &hubCfg)
	if err != nil {
		return err
	}

	h = &Hub{
		clients:         make(map[string]*Client),
		broadcastInput:  make(chan []byte, hubCfg.MaxMsgChan),
		broadcastOutPut: make(chan models.SendMessageBody, hubCfg.MaxMsgChan),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
	}

	// 装配消息解析管道
	if hubCfg.Parser == "" {
		logger.Log.Warn("load default parser json")
		h.parser, _ = parsers.GetParsers("json")
	} else {
		h.parser, err = parsers.GetParsers(hubCfg.Parser)
		if err != nil {
			return err
		}
	}
	// 装配消息处理管道
	if len(hubCfg.Process) <= 0 {
		logger.Log.Warn("load default process default")
		tmp, _ := process.GetParsers("default")
		h.process = append(h.process, tmp)
	} else {
		for _, pro := range hubCfg.Process {
			tmp, err := process.GetParsers(pro)
			if err != nil {
				return err
			}
			h.process = append(h.process, tmp)
		}
	}

	ChatHub = h
	return nil
}

func (rh *Hub) Start(p_ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(p_ctx)
	defer cancel()

	// 消息处理通道
	logger.Log.Info("运行消息处理管道...")
	// go rh.msgProcessPipeline(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-rh.done:
			cancel()
			return
		case reg := <-rh.register:
			if _, ok := rh.clients[reg.user.Id]; !ok {
				rh.clients[reg.user.Id] = reg
				// logger.Log.Info(fmt.Sprintf("[%v:%v] join room [%v]", reg.Id, reg.Name, rh.Rid))
				// rh.broadcastOutPut <- rh.roomSystemMsg(fmt.Sprintf("[%v:%v] join room [%v]", reg.Id, reg.Name, rh.Rid))
			}
		case unreg := <-rh.unregister:
			if _, ok := rh.clients[unreg.user.Id]; ok {
				// logger.Log.Info(fmt.Sprintf("[%v:%v] exit room [%v]", unreg.Id, unreg.Name, rh.Rid))
				delete(rh.clients, unreg.user.Id)
				unreg.Close()
				// rh.broadcastOutPut <- rh.roomSystemMsg(fmt.Sprintf("[%v:%v] exit room [%v]", unreg.Id, unreg.Name, rh.Rid))
			}
		case broadcast := <-rh.broadcastOutPut:
			for _, id := range broadcast.Ids {
				if client, ok := rh.clients[id]; ok {
					client.send <- broadcast.Body.GetJsonBody()
				}
			}
		}
	}
}
