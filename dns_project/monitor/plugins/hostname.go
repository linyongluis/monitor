package plugins

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func GetHostname(ip string, use string, errorfile string, timeout time.Duration) (result string) {

	var out, err1 bytes.Buffer
	var done chan error
	done = make(chan error)

	cmd := exec.Command("ssh", ip, "hostname")
	User, err := user.Lookup(use)
	if err != nil {
		wlog, f := RecordError(errorfile)
		wlog.WriteString(fmt.Sprintf("%s %s hostname plugins get user  is error ,erros: %s\n", time.Now().Format("2006-01-02 15:04:05"), ip, err.Error()))
		wlog.Flush()
		f.Close()
		return "Null"

	}

	uid, _ := strconv.Atoi(User.Uid)
	gid, _ := strconv.Atoi(User.Gid)

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}

	cmd.Stdout = &out
	cmd.Stderr = &err1
	cmd.Start()

	go func() { done <- cmd.Wait() }()

	select {
	case <-time.Tick(timeout * time.Second):
		wlog, f := RecordError(errorfile)
		wlog.WriteString(fmt.Sprintf("%s %s get hostname timeout , the timeout set %d second.\n", time.Now().Format("2006-01-02 15:04:05"), ip, timeout))
		wlog.Flush()
		f.Close()

		return "NULL"
	case <-done:

	}
	stdout := TrimOutput(out)
	stderr := TrimOutput(err1)

	if len(stdout) == 0 && len(stderr) != 0 {
		wlog, f := RecordError(errorfile)
		wlog.WriteString(fmt.Sprintf("%s %s get hostname failed ,erros: %s\n", time.Now().Format("2006-01-02 15:04:05"), ip, stderr))
		wlog.Flush()
		f.Close()

		return "NUll"

	}

	return string(stdout)

}
