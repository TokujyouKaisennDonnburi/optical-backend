package tool

import "fmt"

func statusChunk(status string) []byte {
	statusJson := fmt.Sprintf("{\"progress\":\"%s\"}", status)
	return []byte(statusJson)
}
