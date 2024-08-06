package caching

import (
	"fmt"
)

func Git(action string) error {
	switch action {
	case "startup":
		err := gitCloner()
		if err != nil {
			return err
		}

	case "refresh":
		err := gitPull()
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown error encountered, possible buffer overflow")
	}

	return nil
}
