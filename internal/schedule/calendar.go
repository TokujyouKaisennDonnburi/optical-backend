package schedule

type Calendar []Event

func (c Calendar) append(event Event) Calendar {
    return append(c, event)
}
