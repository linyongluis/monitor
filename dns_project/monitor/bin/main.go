package main

import (
	"bufio"
	"flag"
	"fmt"
	"go_dev/day11/logagent/config"
	"go_dev/dns_project/monitor/ftpclient"
	"go_dev/dns_project/monitor/plugins"
	"os"
	"strings"
	"time"
)

type Worker struct {
	jobchan    chan string
	resultchan chan string
	user       string
	errorfile  string
}

var timeout time.Duration

func (w Worker) CreatWorker() {
	for {
		j := <-w.jobchan
		go w.GetMonitorInfo(j)

	}

}

func (w Worker) GetMonitorInfo(ip string) {

	host := plugins.GetHostname(ip, w.user, w.errorfile, timeout)
	cpudata := plugins.GetCpuInfo(ip, w.user, w.errorfile, timeout)
	memalldata := plugins.GetMem(ip, w.user, w.errorfile, timeout)

	timedata := time.Now().Format("2006/01/02 15:04")

	result := strings.TrimRight(host, "\n") + "(" + ip + ")" + "#" + timedata + "#" + strings.TrimRight(cpudata, "\n") + "#" + strings.TrimRight(memalldata, "\n")

	w.resultchan <- result

}

func main() {

	var filename string

	flag.StringVar(&filename, "f", "", "config file")
	flag.Parse()

	if len(filename) == 0 {
		fmt.Println("please input your config file name")
		return
	}

	conf, err := config.Initconfig("ini", filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := Worker{
		make(chan string, 100),
		make(chan string, 100),
		conf.String("User::user"),
		conf.String("Errorlog::errorlog"),
	}

	tmptimeout, err := conf.Int("MonitorMachineList::timeout")

	if err != nil {
		//wlog,f := plugins.RecordError(w.errorfile)
		fmt.Printf("%s timeout config error, used timeout 10 second.\n", time.Now().Format("2006-01-02 15:04:05"))
		//f.Close()

		timeout = 10

	} else {
		timeout = time.Duration(tmptimeout)
	}

	//fmt.Println("start")
	outlog := conf.String("OutPutLog::OutLog")
	destoutlog := outlog + "." + time.Now().Format("200601021504")

	wirter, err := os.OpenFile(destoutlog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer wirter.Close()

	if err != nil {
		fmt.Println(err)
		wlog, f := plugins.RecordError(w.errorfile)
		wlog.WriteString(fmt.Sprintf("%s open wirter file is error,erros: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		wlog.Flush()
		f.Close()
		return
	}
	buf := bufio.NewWriter(wirter)

	for i := 0; i < 10; i++ {
		go w.CreatWorker()
	}

	iplist := conf.String("MonitorMachineList::iplist")
	montiorlist := strings.Split(iplist, ",")

	alljob := len(montiorlist)
	count := 0

	for _, v := range montiorlist {
		w.jobchan <- v
	}

	for count < alljob {

		select {
		case result := <-w.resultchan:

			_, err := buf.WriteString(result + "\n")
			count++
			if err != nil {
				wlog, f := plugins.RecordError(w.errorfile)
				wlog.WriteString(fmt.Sprintf("%s writing data is error,erros:%s.\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
				wlog.Flush()
				f.Close()
				continue
			}

		}
	}

	_ = buf.Flush()

	transferprotocol := conf.String("TargetFtp::Protocol")

	f := ftpclient.Clientinfo{
		Protocol:     transferprotocol,
		Ip:           conf.String("TargetFtp::Ip"),
		Port:         conf.String("TargetFtp::Port"),
		User:         conf.String("TargetFtp::User"),
		Password:     conf.String("TargetFtp::Password"),
		Tmpfilename:  destoutlog,
		Destfilename: outlog,
		Ftplog:       conf.String("ftplog::ftplog"),
		Destdir:      conf.String("TargetFtp::DestDir"),
	}

	switch {
	case transferprotocol == "ftp":
		f.Ftpclient()
	case transferprotocol == "sftp":
		f.Sftpclient()
	default:
		fmt.Printf("Error protocol is %s,Please input protocol is ftp or sftp.", transferprotocol)

	}

}
