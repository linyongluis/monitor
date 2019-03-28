package ftpclient

import (
	"fmt"
	"go_dev/dns_project/monitor/plugins"
	"os"
	"os/exec"
	"path"
	"time"
)

func (f Clientinfo) Ftpclient() {

	var err error

	ftpserver := f.Protocol + "://" + f.User + ":" + f.Password + "@" + f.Ip + ":" + f.Port + "/" + f.Destdir + "/"

	cmd := exec.Command("curl", "-T", f.Tmpfilename, ftpserver+path.Base(f.Tmpfilename), "-Q", "-RNFR "+path.Base(f.Tmpfilename), "-Q", "-RNTO "+path.Base(f.Destfilename))

	wlog, f1 := plugins.RecordError(f.Ftplog)

	if _, err = cmd.Output(); err != nil {

		wlog.WriteString(fmt.Sprintf("%s put %s  ftp server %s failed , errors: %s.\n", time.Now().Format("2006-01-02 15:04:05"), f.Destfilename, f.Ip, err.Error()))
		wlog.Flush()
		f1.Close()

		os.Exit(1)
	}

	wlog.WriteString(fmt.Sprintf("%s put %s to ftp server %s successed.\n", time.Now().Format("2006-01-02 15:04:05"), f.Destfilename, f.Ip))
	wlog.Flush()
	f1.Close()

}
