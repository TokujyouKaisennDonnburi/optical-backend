package option

const (
	OPTION_PR_REVIEW_PENDING_COUNT = 1
	OPTION_REVIEW_LOAD_STATUS = 2
)

type Option struct {
	Id   int32
	Name string
	Deprecated bool
}

