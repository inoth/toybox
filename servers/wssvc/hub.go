package wssvc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/parsers"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/process"
)

/*
	websocket 长连接服务
	* 是否启用缓存, 本地 OR Redis
	* 开放消息管道插件, 插件化录入
*/

var ChatHub *HubServer

// ChatHub:
//   MaxMsgChan: 10
//   Parser: "json"
//   Process:
//     - "source"
//     - "normal"
type HubServer struct {
	clients         map[string]*Client
	broadcastInput  chan []byte             // 房间消息输入
	broadcastOutPut chan models.MessageBody // 房间消息输出
	register        chan *Client            // 加入通道
	unregister      chan *Client            // 退出通道

	parser  parsers.MsgParsers   // 消息体解析
	process []process.MsgProcess // 消息体处理
}

func (h *HubServer) Init() error {
	var hubCfg models.HubConfig
	err := config.Cfg.UnmarshalKey("ChatHub", &hubCfg)
	if err != nil {
		return err
	}

	h.clients = make(map[string]*Client)
	h.broadcastInput = make(chan []byte, hubCfg.MaxMsgChan)
	h.broadcastOutPut = make(chan models.MessageBody, hubCfg.MaxMsgChan)
	h.register = make(chan *Client, 5)
	h.unregister = make(chan *Client, 5)

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
		tmp, _ := process.GetParsers("source")
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

func (h *HubServer) Start() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 消息处理通道
	logger.Log.Info("run msg process pipeline...")
	go h.msgProcessPipeline(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case reg := <-h.register:
			if _, ok := h.clients[reg.user.Id]; !ok {
				h.clients[reg.user.Id] = reg
				msg := fmt.Sprintf("[%v:%v] Connect.", reg.user.Id, reg.user.Name)
				logger.Log.Info(msg)
				h.broadcastInput <- h.sysMessage(msg)
			}
		case unreg := <-h.unregister:
			if _, ok := h.clients[unreg.user.Id]; ok {
				msg := fmt.Sprintf("[%v:%v] Disconnect.", unreg.user.Id, unreg.user.Name)
				logger.Log.Info(msg)
				delete(h.clients, unreg.user.Id)
				h.broadcastInput <- h.sysMessage(msg)
				unreg.Close()
			}
		case broadcast := <-h.broadcastOutPut:
			for _, id := range broadcast.GetTargets() {
				if client, ok := h.clients[id]; ok {
					client.send <- broadcast.GetJsonBody()
				}
			}
		}
	}
}

func (h *HubServer) SendMessage(msg models.MessageBody) error {
	h.broadcastOutPut <- msg
	return nil
}

// 消息处理管道
//  ______     ┌────────┐    ┌─────────┐     ______
// ()_____)──▶ │ parser │──▶ │ process │──▶ ()_____)
//             └────────┘    └─────────┘
// broadcastOutPut -> parser -> process -> broadcastOutPut
func (h *HubServer) msgProcessPipeline(ctx context.Context) {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				fmt.Printf("%v\n", err)
			}
		}
	}()

	// 输出
	next, out := h.startOuput()
	// 处理
	next, duts := h.startProcess(next)
	// 解析
	bufNext, put := h.startParser(next)
	// 输入
	iut := h.startInput(bufNext)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.runOutput(out)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.runProcess(duts)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.runParser(put)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.runInput(iut)
	}()

	wg.Wait()
	// 消息处理管道结束
	logger.Log.Info("room msg process pipeline done.")
}

type inputUnit struct {
	dst   chan<- []byte
	input <-chan []byte
}

func (h *HubServer) startInput(dst chan<- []byte) *inputUnit {
	return &inputUnit{
		dst:   dst,
		input: h.broadcastInput,
	}
}

func (h *HubServer) runInput(iut *inputUnit) {
	for val := range iut.input {
		iut.dst <- val
	}
}

type parserUnit struct {
	src    <-chan []byte
	dst    chan<- models.MessageBody
	parser parsers.MsgParsers
}

func (h *HubServer) startParser(dst chan<- models.MessageBody) (chan<- []byte, *parserUnit) {
	src := make(chan []byte, 20)
	put := &parserUnit{
		src:    src,
		dst:    dst,
		parser: h.parser,
	}
	return src, put
}

func (h *HubServer) runParser(put *parserUnit) {
	for val := range put.src {
		acc := accumulator.NewAccumulator(put.dst)
		put.parser.Parser(val, acc)
		if acc.Error() != nil {
			logger.Log.Error(acc.Error().Error())
		}
	}
}

type processUnit struct {
	src <-chan models.MessageBody
	dst chan<- models.MessageBody

	process process.MsgProcess
}

func (h *HubServer) startProcess(dst chan<- models.MessageBody) (chan<- models.MessageBody, []*processUnit) {
	var duts []*processUnit

	var src chan models.MessageBody
	for _, process := range h.process {
		src = make(chan models.MessageBody, 20)
		duts = append(duts, &processUnit{
			src:     src,
			dst:     dst,
			process: process,
		})
		dst = src
	}
	return src, duts
}

func (h *HubServer) runProcess(duts []*processUnit) {
	var wg sync.WaitGroup
	for _, dut := range duts {
		wg.Add(1)
		go func(dut *processUnit) {
			defer wg.Done()
			acc := accumulator.NewAccumulator(dut.dst)
			for val := range dut.src {
				dut.process.Process(val, acc)
				if acc.Error() != nil {
					logger.Log.Errorf("[chatsvc] Processor channel; %v", acc.Error().Error())
				}
			}
		}(dut)
	}
	wg.Wait()
}

type outputUnit struct {
	src    <-chan models.MessageBody
	output chan<- models.MessageBody
}

func (h *HubServer) startOuput() (chan<- models.MessageBody, *outputUnit) {
	src := make(chan models.MessageBody, 20)
	out := &outputUnit{
		src:    src,
		output: h.broadcastOutPut,
	}
	return src, out
}

func (h *HubServer) runOutput(out *outputUnit) {
	for val := range out.src {
		out.output <- val
	}
}

type SourceMessage struct {
	Token string `json:"token"`
	// msg | auth | call
	Event     string    `json:"event"`
	EventBody EventBody `json:"eventBody"`
}
type EventBody struct {
	Targets   []string `json:"targets"`
	Msg       string   `json:"msg"`
	Timestamp int64    `json:"timestamp"`
}

func (sm SourceMessage) GetJsonBody() []byte {
	buf, _ := json.Marshal(sm)
	return buf
}

func (h *HubServer) sysMessage(msg string) []byte {
	return (&SourceMessage{
		Event: "system",
		EventBody: EventBody{
			Msg:       msg,
			Timestamp: time.Now().Unix(),
		},
	}).GetJsonBody()
}
