package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func StirngInSlice(s string, l []string) bool {
	for _, b := range l {
		if b == s {
			return true
		}
	}
	return false
}

func DomainBelongsToProvider(d string, l []string) bool {
	for _, pd := range l {
		fmt.Println(pd, d)
		if strings.Contains(d, pd) {
			return true
		}
	}
	return false
}

func Dump(cls interface{}) {
	data, err := json.MarshalIndent(cls, "", "    ")
	if err != nil {
		fmt.Println("[ERROR] Oh no! There was an error on Dump command: ", err)
		return
	}
	fmt.Println(string(data))
}

// WaitFor polls the given function 'f', once every 'interval', up to 'timeout'.
func WaitFor(timeout, interval time.Duration, f func() (bool, error)) error {
	var lastErr string
	timeup := time.After(timeout)
	for {
		select {
		case <-timeup:
			return fmt.Errorf("Time limit exceeded. Last error: %s", lastErr)
		default:
		}

		stop, err := f()
		if stop {
			return nil
		}
		if err != nil {
			lastErr = err.Error()
		}

		time.Sleep(interval)
	}
}
