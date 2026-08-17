package main

import (
	"os"
	"strings"
	"testing"
)

func TestFilterPorts(t *testing.T) {
	in := []PortEntry{{Port: 80}, {Port: 8080}, {Port: 9000}}
	out := FilterPorts(in, 8000, 8999)
	if len(out) != 1 || out[0].Port != 8080 {
		t.Fatalf("unexpected result: %#v", out)
	}
}

func TestWriteCSV(t *testing.T) {
	path := t.TempDir() + "/out.csv"
	if err := WriteCSV(path, []PortEntry{{Protocol: "TCP", Port: 80, PID: 1, Process: "x", User: "u"}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "protocol,port,pid,process,user") {
		t.Fatalf("missing header: %s", data)
	}
}
