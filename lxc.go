package lxc

import (
	"bytes"
	"github.com/s-kostyaev/go-cgroup"
	"os/exec"
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

func IsTmpTmpfs(container string) (bool, error) {
	cmd := exec.Command(
		"lxc-attach", "-e", "-n", container, "--", "/bin/mount",
	)
	cmd.Stdout = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return false, err
	}
	mounts := strings.Split(strings.Trim(cmd.Stdout.(*bytes.Buffer).String(),
		"\n"), "\n")
	for _, mount := range mounts {
		if strs := strings.Fields(mount); strs[2] == "/tmp" &&
			strs[0] == "tmpfs" {
			return true, nil
		}
	}
	return false, nil
}

func GetTmpUsageMb(container string) (int, error) {
	cmd := exec.Command(
		"lxc-attach", "-e", "-n", container, "--", "/usr/bin/du", "-ms", "/tmp",
	)
	cmd.Stdout = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return 0, err
	}
	usage := strings.Fields(cmd.Stdout.(*bytes.Buffer).String())[0]
	result, err := strconv.Atoi(usage)
	return result, err
}

func ClearTmp(container string) error {
	cmd := exec.Command(
		"lxc-attach", "-e", "-n", container, "--",
		"/bin/sh", "-c", "rm /tmp/* /tmp/.* -rf",
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
