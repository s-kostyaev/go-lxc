package lxc

import (
	"github.com/s-kostyaev/go-cgroup"
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
