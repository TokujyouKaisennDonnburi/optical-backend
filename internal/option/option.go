package option

const (
	OPTION_PR_REVIEW_PENDING_COUNT = 1
	OPTION_REVIEW_LOAD_STATUS = 2
	OPTION_INVITE_MEMBERS = 3
	OPTION_MILESTONE_ITEMS = 4
	OPTION_MILESTONE_STATUS = 5
)

type Option struct {
	Id   int32
	Name string
	Deprecated bool
}

