package kill_cord

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/flant/dapp/pkg/dapp"
)

func Init() error {
	if os.Getenv("GITLAB_CI") == "" {
		return nil
	}

	parentsPids, err := getParents()
	if err != nil {
		return err
	}

	go run(parentsPids)

	return nil
}

func run(parentsPids []int) {
	for {
		for _, pid := range parentsPids {
			if isProcessAlive(pid) {
				continue
			}

			ownPid := os.Getpid()

			err := writePidToFile(os.Getpid(), filepath.Join(dapp.GetHomeDir(), ".killed_pids"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Kill cord error: %s\n", err)
			}

			syscall.Kill(ownPid, syscall.SIGINT)

			return
		}

		time.Sleep(1 * time.Second)
	}
}

func writePidToFile(pid int, path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%d\n", pid))
	if err != nil {
		return err
	}

	return nil
}

func getParents() ([]int, error) {
	pid := os.Getpid()
	parents := []int{}

	for {
		ppid, err := getProcessParentPid(pid)
		if err != nil {
			return nil, err
		}
		if ppid == 0 {
			break
		}
		parents = append(parents, ppid)
		pid = ppid
	}

	return parents, nil
}

func getProcessParentPid(pid int) (int, error) {
	path := filepath.Join("/proc", fmt.Sprintf("%d", pid), "status")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineParts := strings.Split(strings.TrimSpace(scanner.Text()), "\t")
		if lineParts[0] == "PPid:" {
			i, err := strconv.ParseInt(lineParts[1], 10, 32)
			if err != nil {
				return 0, err
			}
			return int(i), nil
		}
	}

	return 0, nil
}

func isProcessAlive(pid int) bool {
	err := syscall.Kill(pid, syscall.Signal(0))
	if err == syscall.EPERM || err == nil {
		return true
	}
	return false
}
