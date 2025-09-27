package schedule

type Calendar []Event

func (c *Calendar) append(event Event) Calendar {
	events := append(c.events(), event)
	return Calendar(events)
}

func (c *Calendar) events() []Event {
	return *c
}
