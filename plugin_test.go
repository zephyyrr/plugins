package plugins

type TestPlugin struct {
	e        Event
	args     Args
	recv_err error
	send_err error
}

func (t TestPlugin) Name() string {
	return "TestPlugin"
}

func (t TestPlugin) Provides() []Event {
	return []Event{"test.data", "test.alternative"}
}

func (t TestPlugin) Subscribes() []Event {
	return []Event{"test.answer", "test.alternative"}
}

func (t TestPlugin) Send(Event, Args) error {
	return t.send_err
}

func (t TestPlugin) Recieve() (Event, Args, error) {
	return t.e, t.args, t.recv_err
}
