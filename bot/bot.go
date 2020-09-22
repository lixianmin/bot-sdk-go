package bot

import (
	"github.com/lixianmin/bot-sdk-go/bot/logger"
	"github.com/lixianmin/bot-sdk-go/bot/model"
	"reflect"
)

// 技能基础类
type Bot struct {
	intentHandler map[string]func(bot *Bot, request *model.IntentRequest) // 针对intent requset不同intent的处理函数
	//eventHandler               map[string]func(bot *Bot, request *model.EventRequest)  // 针对事件的处理函数
	eventHandler               map[string]func(bot *Bot, request interface{}) // 针对事件的处理函数
	defaultEventHandler        func(bot *Bot, request interface{})
	launchRequestHandler       func(bot *Bot, request *model.LaunchRequest)       // 针对技能打开的处理函数
	sessionEndedRequestHandler func(bot *Bot, request *model.SessionEndedRequest) // 针对技能关闭的处理函数
	Request                    interface{}                                        // 对当前request的封装，需要在使用时断言，判断当前的类型
	Session                    *model.Session                                     // 对session的封装
	Response                   *model.Response                                    // 对技能返回的封装
}

// 创建常驻bot类，可维持在内存状态中, addhandler 和 addEventer事件可以缩减为一次
func NewBot() *Bot {
	return &Bot{
		intentHandler: make(map[string]func(bot *Bot, request *model.IntentRequest)),
		eventHandler:  make(map[string]func(bot *Bot, request interface{})),
	}
}

// 根据每个请求分别处理
func (my *Bot) Handler(request string) string {
	//logger.Debug(request)
	
	my.Request = model.NewRequest(request)
	my.Session = model.NewSession(model.GetSessionData(request))
	my.Response = model.NewResponse(my.Session, my.Request)

	my.dispatch()

	return my.Response.Build()
}

// 添加对intent的处理函数
func (my *Bot) AddIntentHandler(intentName string, fn func(bot *Bot, request *model.IntentRequest)) {
	if intentName != "" {
		my.intentHandler[intentName] = fn
	}
}

// 添加对事件的处理函数
func (my *Bot) AddEventListener(eventName string, fn func(bot *Bot, request interface{})) {
	if eventName != "" {
		my.eventHandler[eventName] = fn
	}
}

// 添加事件默认处理函数
// 比如，在播放视频时，技能会收到各种事件的上报，如果不想一一处理可以使用这个来添加处理
func (my *Bot) AddDefaultEventListener(fn func(bot *Bot, request interface{})) {
	my.defaultEventHandler = fn
}

// 打开技能时的处理
func (my *Bot) OnLaunchRequest(fn func(bot *Bot, request *model.LaunchRequest)) {
	my.launchRequestHandler = fn
}

// 技能关闭的处理，比如可以做一些清理的工作
// TIP: 根据协议，技能关闭返回的结果，DuerOS不会返回给用户。
func (my *Bot) OnSessionEndedRequest(fn func(bot *Bot, request *model.SessionEndedRequest)) {
	my.sessionEndedRequestHandler = fn
}

func (my *Bot) dispatch() {
	switch request := my.Request.(type) {
	case model.IntentRequest:
		my.processIntentHandler(request)
		return
	case model.LaunchRequest:
		my.processLaunchHandler(request)
		return
	case model.SessionEndedRequest:
		my.processSessionEndedHandler(request)
		return
	}
	my.processEventHandler(my.Request)
}

func (my *Bot) processLaunchHandler(request model.LaunchRequest) {
	logger.Info("uid=%q, sn=%q", request.GetUserId(), request.GetOriginalDeviceId())

	if my.launchRequestHandler != nil {
		my.launchRequestHandler(my, &request)
	}
}

func (my *Bot) processSessionEndedHandler(request model.SessionEndedRequest) {
	logger.Info("uid=%q, sn=%q", request.GetUserId(), request.GetOriginalDeviceId())

	if my.sessionEndedRequestHandler != nil {
		my.sessionEndedRequestHandler(my, &request)
	}
}

func (my *Bot) processIntentHandler(request model.IntentRequest) {
	intentName, _ := request.GetIntentName()
	fn, ok := my.intentHandler[intentName]
	logger.Info("uid=%q, sn=%q, intentName=%q, hasHandler=%v", request.GetUserId(), request.GetOriginalDeviceId(), intentName, ok)

	if ok {
		fn(my, &request)
		return
	}
}

func (my *Bot) processEventHandler(request interface{}) {
	rVal := reflect.ValueOf(request)
	eventType := rVal.FieldByName("Type").Interface().(string)

	fn, ok := my.eventHandler[eventType]
	logger.Info("eventType=%q, hasHandler=%v", eventType, ok)

	if ok {
		fn(my, request)
		return
	}

	if my.defaultEventHandler != nil {
		my.defaultEventHandler(my, request)
	}
}
