package main

import (
	"os"
	"testing"
)

// Live check against a real VM: VIRTUSB_VM=win11 VIRTUSB_DEV=46f8:0855 go test
// Detaches (if attached), attaches, and verifies the device shows in the live XML. Leaves it attached.
func TestAttachDetach(t *testing.T) {
	vm, id := os.Getenv("VIRTUSB_VM"), os.Getenv("VIRTUSB_DEV")
	if vm == "" || id == "" {
		t.Skip("set VIRTUSB_VM and VIRTUSB_DEV")
	}
	d := usbDev{Vendor: id[:4], Product: id[5:]}
	if attachedUSB(vm)[id] {
		if err := toggle("detach-device", vm, d); err != nil {
			t.Fatal(err)
		}
	}
	if attachedUSB(vm)[id] {
		t.Fatal("still attached after detach")
	}
	if err := toggle("attach-device", vm, d); err != nil {
		t.Fatal(err)
	}
	if !attachedUSB(vm)[id] {
		t.Fatal("not attached")
	}
}
