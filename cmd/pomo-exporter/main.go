package main

import (
	"fmt"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func main() {
	fmt.Println("Hello pomo exporter")

	hello := ticktick.HelloTicktick()
	fmt.Println(hello)
}
