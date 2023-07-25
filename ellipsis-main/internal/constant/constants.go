package constant

const AppName = "UnifyLogic"

const (
	Admin      string = "admin"
	User       string = "user"
	CounterKey string = "counter"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

var Roles = map[string]int{
	Admin: 1,
	User:  2,
}
