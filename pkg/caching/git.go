package caching

import (
	"fmt"
)

func Git(action string) error {
	fmt.Printf(".")
	// do something
	switch action {
	case "startup":
		err := gitCloner()
		if err != nil {
			return err
		}
	case "refresh":
		err := gitFetcher()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown error encountered, possible buffer overflow")
	}
	fmt.Printf(".")
	return nil
}
