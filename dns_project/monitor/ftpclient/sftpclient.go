package ftpclient

import (
	"bufio"
	"fmt"
	"github.com/pkg/sftp"
	"go_dev/dns_project/monitor/plugins"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path"
	"time"
)

func (f Clientinfo) Sftpclient() {

	wlog, f1 := plugins.RecordError(f.Ftplog)

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	config := &ssh.ClientConfig{
		User: f.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(f.Password),
		},
		HostKeyCallback: hostKeyCallbk,
	}
	addr := f.Ip + ":" + f.Port
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {

		wlog.WriteString(fmt.Sprintf("%s Failed to dial %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), f.Ip, err.Error()))
		wlog.Flush()
		f1.Close()
		return
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	sftp, err := sftp.NewClient(client)
	if err != nil {

		wlog.WriteString(fmt.Sprintf("%s Failed to create sftp client: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		wlog.Flush()
		f1.Close()
		return
	}
	defer sftp.Close()

	srcfile, err := os.Open(f.Tmpfilename)
	defer srcfile.Close()
	if err != nil {
		wlog.WriteString(fmt.Sprintf("%s Src file %s read failed: %s\n", time.Now().Format("2006-01-02 15:04:05"), f.Tmpfilename, err.Error()))
		wlog.Flush()
		f1.Close()
		return

	}

	sftptmpfilename := "/" + f.Destdir + "/" + path.Base(f.Tmpfilename)
	sftpdestfilename := "/" + f.Destdir + "/" + path.Base(f.Destfilename)

	//fmt.Println(sftptmpfilename,sftpdestfilename)

	desfile, err := sftp.Create(sftptmpfilename)
	defer desfile.Close()
	if err != nil {
		//fmt.Println(sftptmpfilename)
		wlog.WriteString(fmt.Sprintf("%s Sftp create file %s error: %s\n", time.Now().Format("2006-01-02 15:04:05"), sftptmpfilename, err.Error()))
		wlog.Flush()
		f1.Close()
		return

	}

	rbuf := bufio.NewReader(srcfile)

	for {
		wdata, err := rbuf.ReadBytes('\n')

		if err == io.EOF {
			break
		}
		desfile.Write(wdata)
	}

	//buf := make([]byte,1024)
	//for {
	//	n,_ := srcfile.Read(buf)
	//	if n == 0 {
	//		break
	//	}
	//	desfile.Write(buf)
	//}

	_, err = sftp.Lstat(sftpdestfilename)
	if err == nil {

		err = sftp.Remove(sftpdestfilename)
		if err != nil {

			wlog.WriteString(fmt.Sprintf("%s Sftp  remove file %s  failed: %s\n", time.Now().Format("2006-01-02 15:04:05"), sftptmpfilename, err.Error()))
			wlog.Flush()
			f1.Close()
			return
		}
	}

	err = sftp.Rename(sftptmpfilename, sftpdestfilename)
	if err != nil {

		wlog.WriteString(fmt.Sprintf("%s Sftp  rename file %s to %s failed: %s\n", time.Now().Format("2006-01-02 15:04:05"), sftptmpfilename, sftpdestfilename, err.Error()))
		wlog.Flush()
		f1.Close()
		return

	}
	//sftpserver := f.Protocol + "://" + f.User + ":" + f.Password + "@" + f.Ip + ":"+ f.Port + "/" + f.Destdir + "/"
	//
	//
	//cmd := exec.Command("curl",sftpserver,"-Q", "-rename "+ sftptmpfilename + " " + sftpdestfilename)
	//err = cmd.Run()
	//if err != nil {
	//	wlog.Fatalf("%s Sftp  rename file error: %s" ,time.Now().Format("2006-01-02 15:04:05"), err.Error())
	//	f1.Close()
	//	return
	//}

	wlog.WriteString(fmt.Sprintf("%s Sftp file %s put successed.\n", time.Now().Format("2006-01-02 15:04:05"), f.Destfilename))
	wlog.Flush()
	f1.Close()

}
