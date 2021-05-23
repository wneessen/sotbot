package handler

import (
	"fmt"
	"runtime"
)

// Let's the bot provide some memory indicators
func TellMemUsage() (string, bool) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memResponse := fmt.Sprintf(
		"Memory allocated: %v MiB\nTotal allocated: %v MiB\nSys Memory allocated: %v MiB\n"+
			"Number of GCs: %v",
		memStats.Alloc/1024/1024, memStats.TotalAlloc/1024/1024, memStats.Sys/1024/1024,
		memStats.NumGC)

	return memResponse, true
}
