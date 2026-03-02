package pomodoro

const dateFormat = "02/01/2006"
const timeFormat = "15:04:05"

type exportArgs struct {
	Year            int
	Month           int
	IncludeTags     []string
	ExcludeTags     []string
	IncludeProjects []string
	ExcludeProjects []string
	Output          string
}
