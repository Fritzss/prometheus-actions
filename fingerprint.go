package main

import (
	"github.com/shirou/gopsutil/host"
)

type Fingerprint struct {
	OSName        string
	OSFamily      string
	OSVersion     string
	KernelName    string
	KernelVersion string
	Hostname      string
}

func BuildFingerprint() (*Fingerprint, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	f := &Fingerprint{
		OSName:        hostInfo.Platform,
		OSFamily:      hostInfo.PlatformFamily,
		OSVersion:     hostInfo.PlatformVersion,
		KernelName:    hostInfo.OS,
		KernelVersion: hostInfo.KernelVersion,
		Hostname:      hostInfo.Hostname,
	}
	return f, nil
}
