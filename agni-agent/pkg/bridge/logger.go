package bridge

import "fmt"

func FormatError(message string, err error) error {
	if err == nil {
		return fmt.Errorf("[Agni Agent] %s ", message)
	}
	return fmt.Errorf("[Agni Agent] %s -- %v", message, err)
}
