package calendar

type Schedules []Event

func (c Schedules) append(event Event) Schedules {
    return append(c, event)
}
