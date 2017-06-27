package conditional

// ConditionStateObserver represents a type that listens on condition state
// changes.
type ConditionStateObserver interface {
	OnChange(bool)
}

type channelObserver struct {
	ch chan<- bool
}

// NewChannelObserver creates a new condition state observer that writes the
// state changes to the specified channel.
func NewChannelObserver(ch chan<- bool) ConditionStateObserver {
	return channelObserver{ch: ch}
}

func (o channelObserver) OnChange(state bool) {
	o.ch <- state
}
