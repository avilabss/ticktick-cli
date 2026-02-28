package main

const DateFormat = "02/01/2006"
const TimeFormat = "15:04:05"

type Args struct {
	Year            int
	Month           int
	IncludeTags     []string
	ExcludeTags     []string
	IncludeProjects []string
	ExcludeProjects []string
	Output          string
}
