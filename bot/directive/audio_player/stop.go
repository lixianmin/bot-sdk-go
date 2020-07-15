package audio_player

import (
	"github.com/lixianmin/bot-sdk-go/bot/directive"
)

type StopDirective struct {
	directive.BaseDirective
}

func NewStopDirective() *StopDirective {
	stop := &StopDirective{}
	stop.Type = "AudioPlayer.Stop"
	return stop
}
