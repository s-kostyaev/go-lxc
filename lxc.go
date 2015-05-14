package lxc

import (
	"bytes"
	"errors"
	"github.com/s-kostyaev/go-cgroup"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

const (
	configPathPrefix = "/var/lib/lxc/"
)

type Container struct {
	Name   string   `json:"name"`
	Host   string   `json:"host"`
	Status string   `json:"status"`
	Image  []string `json:"image"`
	Ip     string   `json:"ip"`
	Key    string   `json:"key"`
}

func GetMemoryLimit(container string) (int, error) {
	limit, err := cgroup.GetParamInt("memory/lxc/"+container,
		cgroup.MemoryLimit)
	return limit, err
}

func GetMemoryUsage(container string) (int, error) {
	usage, err := cgroup.GetParamInt("memory/lxc/"+container,
		cgroup.MemoryUsage)
	return usage, err
}

func GetMemoryPids(container string) ([]int32, error) {
	result := []int32{}
	pids, err := cgroup.GetParam("memory/lxc/"+container,
		"tasks")
	if err != nil {
		return nil, err
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

func GetCpuTicks() (ticks int, err error) {
	stats, err := cgroup.GetParam("cpu/lxc/", "cpuacct.stat")
	if err != nil {
		return 0, err
	}

	userTime, err := strconv.Atoi(strings.Fields(string(stats))[1])
	if err != nil {
		return 0, err
	}

	systemTime, err := strconv.Atoi(strings.Fields(string(stats))[3])
	if err != nil {
		return 0, err
	}

	return userTime + systemTime, nil
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

func GetRootFS(container string) (string, error) {
	return GetConfigItem(container, "lxc.rootfs")
}

func GetContainers() ([]string, error) {
	cmd := exec.Command("/usr/bin/lxc-ls", "-1")
	list, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.Trim(string(list), "\n"), "\n"), nil
}

func GetState(container string) (string, error) {
	cmd := exec.Command("/usr/bin/lxc-info", "-sHn", container)
	info, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(info), "\n"), nil
}

func GetConfigItem(container, key string) (string, error) {
	prefix := key + " = "
	buf, err := ioutil.ReadFile(configPathPrefix + container + "/config")
	if err != nil {
		return "", err
	}
	for _, string := range strings.Split(string(buf), "\n") {
		if strings.HasPrefix(string, prefix) {
			return strings.TrimPrefix(string, prefix), nil
		}
	}
	return "", errors.New("Could not find item in config")
}
