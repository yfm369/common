/*!
 * <性能分析>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package profile

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	pub "public_yfm"
	"runtime/debug"
	"runtime/pprof"
)

const (
	PROFILE_LOG = "profile"

	PROFILE_CPU    = "cpu.prof"
	PROFILE_MEMORY = "memory.prof"
	PROFILE_STACK  = "stack.prof"
)

var cpufile *os.File

func StartCpuProfile() {
	if cpufile != nil {
		pub.PrintFileLog(PROFILE_LOG, "StartCpuProfile error cpufile != nil")
		return
	}

	var errc error
	cpufile, errc = os.Create(PROFILE_CPU)
	if errc != nil {
		pub.PrintFileLog(PROFILE_LOG, "StartCpuProfile error :", errc.Error())
		return
	}

	errc = pprof.StartCPUProfile(cpufile)
	if errc != nil {
		pub.PrintFileLog(PROFILE_LOG, "StartCpuProfile start error :", errc.Error())
		cpufile.Close()
		return
	}
}

func StopCpuProfile() {
	if cpufile == nil {
		return
	}

	pprof.StopCPUProfile()

	cpufile.Close()
	cpufile = nil
}

func StartMemoryProfile() {
	file, err := os.Create(PROFILE_MEMORY)
	if err != nil {
		pub.PrintFileLog(PROFILE_LOG, "StartMemoryProfile error :", err.Error())
		return
	}
	defer file.Close()

	err = pprof.WriteHeapProfile(file)
	if err != nil {
		pub.PrintFileLog(PROFILE_LOG, "StartMemoryProfileb writeheap error :", err.Error())
		return
	}
}

func PrintStack() {
	data := debug.Stack()
	file, err := os.Create(PROFILE_STACK)
	if err != nil {
		pub.PrintFileLog(PROFILE_LOG, "PrintStack error :", err.Error())
		return
	}
	defer file.Close()

	file.Write(data)
}

func HttpProfile(port string) {
	http.ListenAndServe(":"+port, nil)
}
