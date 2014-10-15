package lxc

import (
	"github.com/s-kostyaev/go-cgroup"
	"strconv"
	"strings"
)

func GetMemoryLimit(container string) (int, error) {
	limit, err := cgroup.GetParamInt("memory/lxc/"+container,
		cgroup.MemoryLimit)
	if err != nil {
		return cgroup.GetParamInt("memory/lxc/"+container+"-1",
			cgroup.MemoryLimit)
	}
	return limit, err
}

func GetMemoryUsage(container string) (int, error) {
	usage, err := cgroup.GetParamInt("memory/lxc/"+container,
		cgroup.MemoryUsage)
	if err != nil {
		return cgroup.GetParamInt("memory/lxc/"+container+"-1",
			cgroup.MemoryUsage)
	}
	return usage, err
}

func GetMemoryPids(container string) ([]int32, error) {
	result := []int32{}
	pids, err := cgroup.GetParam("memory/lxc/"+container,
		"tasks")
	if err != nil {
		pids, err = cgroup.GetParam("memory/lxc/"+container,
			"tasks")
		if err != nil {
			return nil, err
		}
	}
	for _, str := range strings.Split(pids, "\n") {
		pid, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result = append(result, int32(pid))
	}
	return result, nil
}
