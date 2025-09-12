package schedule

type Calender []Event

func (c *Calender) append(event Event) Calender {
	events := append(c.events(), event)
	return Calender(events)
}

func (c *Calender) events() []Event {
	return *c
}
