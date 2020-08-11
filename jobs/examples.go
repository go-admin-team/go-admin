package jobs

import (
	"fmt"
	"strconv"
)

func (t *EXEC)ExamplesWithParam(i int , s string) string{
	fmt.Println("call method PrintInfo i", i, ",s :", s)
	return s +strconv.Itoa(i)
}

func (t *EXEC) ExamplesNoParam() string {
	fmt.Println("\nshow msg input 'call reflect'")
	return "ShowMsg"
}
