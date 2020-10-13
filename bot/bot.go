package bot

import (
	"github.com/lixianmin/bot-sdk-go/bot/logger"
	"github.com/lixianmin/bot-sdk-go/bot/model"
	"reflect"
)

type (
	IntentFunc              func(bot *Bot, request *model.IntentRequest, response *model.Response)
	EventFunc               func(bot *Bot, request interface{}, response *model.Response)
	LaunchRequestFunc       func(bot *Bot, request *model.LaunchRequest, response *model.Response)
	SessionEndedRequestFunc func(bot *Bot, request *model.SessionEndedRequest, response *model.Response)

	// 技能基础类
	Bot struct {
		intentHandler              map[string]IntentFunc // 针对intent request不同intent的处理函数
		eventHandler               map[string]EventFunc  // 针对事件的处理函数
		defaultEventHandler        EventFunc
		launchRequestHandler       LaunchRequestFunc       // 针对技能打开的处理函数
		sessionEndedRequestHandler SessionEndedRequestFunc // 针对技能关闭的处理函数
	}
)

// 创建常驻bot类，可维持在内存状态中, addhandler 和 addEventer事件可以缩减为一次
func NewBot() *Bot {
	return &Bot{
		intentHandler: make(map[string]IntentFunc),
		eventHandler:  make(map[string]EventFunc),
	}
}

// 根据每个请求分别处理
func (my *Bot) Handler(rawData []byte) string {
	//logger.Debug(rawData)

	var request = model.NewRequest(rawData)                       // 对当前request的封装，需要在使用时断言，判断当前的类型
	var session = model.NewSession(model.GetSessionData(rawData)) // 对session的封装
	var response = model.NewResponse(session, request)            // 对技能返回的封装

	my.dispatch(request, response)

	return response.Build()
}

// 添加对intent的处理函数
func (my *Bot) AddIntentHandler(intentName string, fn IntentFunc) {
	if intentName != "" {
		my.intentHandler[intentName] = fn
	}
}

// 添加对事件的处理函数
func (my *Bot) AddEventListener(eventName string, fn EventFunc) {
	if eventName != "" {
		my.eventHandler[eventName] = fn
	}
}

// 添加事件默认处理函数
// 比如，在播放视频时，技能会收到各种事件的上报，如果不想一一处理可以使用这个来添加处理
func (my *Bot) AddDefaultEventListener(fn EventFunc) {
	my.defaultEventHandler = fn
}

// 打开技能时的处理
func (my *Bot) OnLaunchRequest(fn LaunchRequestFunc) {
	my.launchRequestHandler = fn
}

// 技能关闭的处理，比如可以做一些清理的工作
// TIP: 根据协议，技能关闭返回的结果，DuerOS不会返回给用户。
func (my *Bot) OnSessionEndedRequest(fn SessionEndedRequestFunc) {
	my.sessionEndedRequestHandler = fn
}

func (my *Bot) dispatch(request interface{}, response *model.Response) {
	logger.Debug(request)
	switch request := request.(type) {
	case model.IntentRequest:
		my.processIntentHandler(request, response)
		return
	case model.LaunchRequest:
		my.processLaunchHandler(request, response)
		return
	case model.SessionEndedRequest:
		my.processSessionEndedHandler(request, response)
		return
	}

	my.processEventHandler(request, response)
}

func (my *Bot) processLaunchHandler(request model.LaunchRequest, response *model.Response) {
	logger.Info("uid=%q, sn=%q", request.GetUserId(), request.GetOriginalDeviceId())

	if my.launchRequestHandler != nil {
		my.launchRequestHandler(my, &request, response)
	}
}

func (my *Bot) processSessionEndedHandler(request model.SessionEndedRequest, response *model.Response) {
	logger.Info("uid=%q, sn=%q", request.GetUserId(), request.GetOriginalDeviceId())

	if my.sessionEndedRequestHandler != nil {
		my.sessionEndedRequestHandler(my, &request, response)
	}
}

func (my *Bot) processIntentHandler(request model.IntentRequest, response *model.Response) {
	intentName, _ := request.GetIntentName()
	fn, ok := my.intentHandler[intentName]
	logger.Info("uid=%q, sn=%q, intentName=%q, hasHandler=%v", request.GetUserId(), request.GetOriginalDeviceId(), intentName, ok)

	if ok {
		fn(my, &request, response)
		return
	}
}

func (my *Bot) processEventHandler(request interface{}, response *model.Response) {
	rVal := reflect.ValueOf(request)
	eventType := rVal.FieldByName("Type").Interface().(string)

	fn, ok := my.eventHandler[eventType]
	logger.Info("eventType=%q, hasHandler=%v", eventType, ok)

	if ok {
		fn(my, request, response)
		return
	}

	if my.defaultEventHandler != nil {
		my.defaultEventHandler(my, request, response)
	}
}
