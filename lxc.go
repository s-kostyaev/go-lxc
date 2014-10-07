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

func GetPids(container string) ([]int, error) {
	result := []int{}
	pids, err := cgroup.GetParam("cpuacct/lxc/"+container,
		"cgroup.procs")
	if err != nil {
		pids, err = cgroup.GetParam("cpuacct/lxc/"+container,
			"cgroup.procs")
		if err != nil {
			return nil, err
		}
	}
	for _, str := range strings.Split(pids, "\n") {
		pid, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result = append(result, pid)
	}
	return result, nil
}
