package testcase

import (
	"go_dev/dns_project/monitor/plugins"
	"testing"
)

func TestMytest(t *testing.T) {

	wlog, f := plugins.RecordError("abc_1123")
	wlog.WriteString("this is my test aa a a a a")
	wlog.Flush()
	f.Close()

	t.Log("over")
}
