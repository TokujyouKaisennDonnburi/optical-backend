package tool

import "fmt"

func statusChunk(status string) []byte {
	statusJson := fmt.Sprintf("{\"status\":\"%s\"}", status)
	return []byte(statusJson)
}
