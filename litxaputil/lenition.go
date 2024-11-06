package litxaputil

import "strings"

func ApplyLenition(current string) (lenition string, next string) {
	switch {
	case strings.HasPrefix(current, "ts"):
		lenition = "ts→s"
		next = current[1:]
	case strings.HasPrefix(current, "tx"):
		lenition = "tx→t"
		next = "t" + current[2:]
	case strings.HasPrefix(current, "kx"):
		lenition = "kx→k"
		next = "k" + current[2:]
	case strings.HasPrefix(current, "px"):
		lenition = "px→p"
		next = "p" + current[2:]
	case strings.HasPrefix(current, "t"):
		lenition = "t→s"
		next = "s" + current[1:]
	case strings.HasPrefix(current, "k"):
		lenition = "k→h"
		next = "h" + current[1:]
	case strings.HasPrefix(current, "p"):
		lenition = "p→f"
		next = "f" + current[1:]
	case strings.HasPrefix(current, "'l"), strings.HasPrefix(current, "'r"):
		next = current
	case strings.HasPrefix(current, "'"):
		lenition = current[0:2] + "→" + current[1:2]
		next = current[1:]
	default:
		next = current
	}

	return
}
