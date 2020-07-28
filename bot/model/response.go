package model

import (
	"encoding/json"

	"github.com/lixianmin/bot-sdk-go/bot/data"
	"github.com/lixianmin/bot-sdk-go/bot/util"
)

type Response struct {
	session *Session
	request interface{}
	data    map[string]interface{}
}

func NewResponse(session *Session, request interface{}) *Response {
	d := make(map[string]interface{})
	return &Response{
		data:    d,
		session: session,
		request: request,
	}
}

/**
 * 询问用户时，返回的speech.
 * 此时设备的麦克风会进入收音状态，比如设备灯光亮起
 * TIP: 一般技能要完成一项任务，还缺少一些信息，主动发起对用户的询问的时候使用
 */
func (my *Response) Ask(speech string) *Response {
	my.Tell(speech)
	my.HoldOn()
	return my
}

func (my *Response) AskSlot(speech string, slot string) *Response {
	my.Ask(speech)

	request, ok := my.request.(IntentRequest)
	if ok {
		request.Dialog.ElicitSlot(slot)
	}
	return my
}

/**
 * 回复用户，返回的speech
 */
func (my *Response) Tell(speech string) *Response {
	my.data["outputSpeech"] = util.FormatSpeech(speech)
	return my
}

/**
 * 回复用户，返回的speech
 */
func (my *Response) Reprompt(speech string) *Response {
	my.data["reprompt"] = map[string]interface{}{
		"outputSpeech": util.FormatSpeech(speech),
	}
	return my
}

/**
 * 返回卡片.
 * 针对有屏幕的设备，比如: 电视、show，可以呈现更多丰富的信息给用户
 * 卡片协议参考：TODO
 */
func (my *Response) DisplayCard(card interface{}) *Response {
	my.data["card"] = card

	return my
}

/**
 * 返回指令. 比如，返回音频播放指令，使设备开始播放音频
 * TIP: 可以同时返回多个指令，设备按返回顺序执行这些指令，指令协议参考TODO
 */
func (my *Response) Command(directive interface{}) *Response {
	_, ok := my.data["directives"]
	if !ok {
		my.data["directives"] = make([]interface{}, 0)
	}

	directives, ok := my.data["directives"].([]interface{})
	directives = append(directives, directive)

	my.data["directives"] = directives

	return my
}

/**
 * 保持会话.
 * 此时设备的麦克风会自动开启监听用户说话
 */
func (my *Response) HoldOn() *Response {
	my.data["shouldEndSession"] = false
	return my
}

/**
 * 保持会话.
 * 关闭麦克风
 */
func (my *Response) CloseMicrophone() *Response {
	my.data["expectSpeech"] = true
	return my
}

func (my *Response) Build() string {
	//session
	attributes := my.session.GetData().Attributes

	ret := map[string]interface{}{
		"version":  "2.0",
		"session":  data.SessionResponse{Attributes: attributes},
		"response": my.data,
	}

	//intent request
	request, ok := my.request.(IntentRequest)
	if ok {
		ret["context"] = data.ContextResponse{Intent: request.Dialog.Intents[0].GetData()}

		directive := request.Dialog.GetDirective()
		if directive != nil {
			my.Command(directive)
		}
	}

	response, _ := json.Marshal(ret)

	return string(response)
}

func (my *Response) GetData() map[string]interface{} {
	return my.data
}
