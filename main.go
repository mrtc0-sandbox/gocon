package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func catch(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"level": "error",
		}).Info(err)
	}
	os.Exit(1)
}

func init() {
	// 初期化処理
	if os.Args[0] == "gocon-init" {
		root, _ := os.Getwd()
		// Mount
		if err := syscall.Mount("proc", filepath.Join(root, "/proc"), "proc", 0, ""); err != nil {
			catch(err)
		}
		// chroot
		if err := syscall.Chroot(root); err != nil {
			catch(err)
		}
		// chdir
		if err := syscall.Chdir("/"); err != nil {
			catch(err)
		}
		if err := syscall.Sethostname([]byte("box")); err != nil {
			catch(err)
		}
		if err := syscall.Exec(os.Args[1], os.Args[1:], os.Environ()); err != nil {
			catch(err)
		}
	}
}

func main() {
	// このプロセスの引数にARGSを与えて実行して、隔離する
	cmd := exec.Command("/proc/self/exe", os.Args[1:]...)
	cmd.Args[0] = "gocon-init"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
	}

	cmd.Run()
}
