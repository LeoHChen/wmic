// +build windows

package wmic

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	disk, err := GetDiskDriveInfo()

	if err != nil {
		t.Errorf("failed: %v", err)
	} else {
		fmt.Println("disk =>", disk)
	}

	port, err := GetPortInfo()

	if err != nil {
		t.Errorf("failed: %v", err)
	} else {
		fmt.Println("port =>", port)
	}

	vol, err := GetVolumeInfo()

	if err != nil {
		t.Errorf("failed: %v", err)
	} else {
		fmt.Println("volume =>", vol)
	}
}
