package apis

import "strconv"

func parsePathId(value string) int {
	id, _ := strconv.Atoi(value)
	return id
}
