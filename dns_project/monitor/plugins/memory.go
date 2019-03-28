package plugins

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func GetMem(ip string, use string, errorfile string, timeout time.Duration) (result string) {

	var out, err1 bytes.Buffer
	var done chan error
	done = make(chan error)

	cmd := exec.Command("ssh", ip, "sar", "-r", "1", "5", "|", "grep", "Average")

	User, err := user.Lookup(use)
	if err != nil {
		wlog, f := RecordError(errorfile)
		wlog.WriteString(fmt.Sprintf("%s %s memory plugins get user  is error ,erros: %s\n", time.Now().Format("2006-01-02 15:04:05"), ip, err.Error()))
		wlog.Flush()
		f.Close()
		return "NULL"

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
		wlog.WriteString(fmt.Sprintf("%s %s get memory timeout , the timeout set %d second.\n", time.Now().Format("2006-01-02 15:04:05"), ip, timeout))
		wlog.Flush()
		f.Close()

		return "NULL"
	case <-done:

	}

	stdout := TrimOutput(out)
	stderr := TrimOutput(err1)

	if len(stdout) == 0 && len(stderr) != 0 {
		wlog, f := RecordError(errorfile)
		wlog.WriteString(fmt.Sprintf("%s %s get memory  failed ,erros: %s\n", time.Now().Format("2006-01-02 15:04:05"), ip, stderr))
		wlog.Flush()
		f.Close()

		return "NUll"

	}

	d := strings.Split(string(stdout), " ")
	var new_array []float64
	count := 0
	for _, v := range d {
		if len(v) != 0 {
			count++
			if count > 1 {
				data, _ := strconv.ParseFloat(v, 64)

				new_array = append(new_array, data)
			}

		}
	}

	kbmemfree := new_array[0]
	kbmemused := new_array[1]
	kbbuffers := new_array[3]
	kbcached := new_array[4]

	totalmem := kbmemfree + kbmemused
	freemem := kbmemfree + kbbuffers + kbcached

	free := freemem / totalmem

	usepercent := (1 - free) * 100

	return (fmt.Sprintf("%.2f%%", usepercent))

}
