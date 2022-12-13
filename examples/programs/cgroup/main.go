package main

import (
	"bytes"
	_ "embed"

	manager "github.com/DataDog/ebpf-manager"

	"github.com/sirupsen/logrus"
)

//go:embed ebpf/bin/probe.o
var Probe []byte

var m = &manager.Manager{
	Probes: []*manager.Probe{
		{
			ProbeIdentificationPair: manager.ProbeIdentificationPair{
				EBPFFuncName: "egress",
			},
			CGroupPath: "/sys/fs/cgroup/unified",
		},
	},
}

func main() {
	// Initialize the manager
	if err := m.Init(bytes.NewReader(Probe)); err != nil {
		logrus.Fatal(err)
	}

	// Start the manager
	if err := m.Start(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("successfully started, head over to /sys/kernel/debug/tracing/trace_pipe")

	// Generate some network traffic to trigger the probe
	trigger()

	// Close the manager
	if err := m.Stop(manager.CleanAll); err != nil {
		logrus.Fatal(err)
	}
}
