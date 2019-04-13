// utils
package gpss

import (
	"fmt"
	"math/rand"
	"time"
)

func PrintlnVerbose(verbose bool, a ...interface{}) {
	if !verbose {
		return
	}
	var s string
	for _, v := range a {
		s = fmt.Sprint(s, v)
	}
	fmt.Println(s)
}

func GetRandom(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}
