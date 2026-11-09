// virtusb: GUI to attach/detach host USB devices to libvirt VMs via virsh.
package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var uri = envOr("LIBVIRT_DEFAULT_URI", "qemu:///system")

type usbDev struct{ Vendor, Product, Name string }

func (d usbDev) id() string { return d.Vendor + ":" + d.Product }

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func virsh(args ...string) (string, error) {
	out, err := exec.Command("virsh", append([]string{"-c", uri}, args...)...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("virsh %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out), nil
}

func listVMs() []string {
	out, _ := virsh("list", "--all", "--name")
	return strings.Fields(out)
}

// hostUSB reads USB devices from sysfs, skipping root hubs.
func hostUSB() []usbDev {
	read := func(dir, f string) string {
		b, _ := os.ReadFile(filepath.Join(dir, f))
		return strings.TrimSpace(string(b))
	}
	dirs, _ := filepath.Glob("/sys/bus/usb/devices/*")
	var devs []usbDev
	for _, d := range dirs {
		v := read(d, "idVendor")
		if v == "" || v == "1d6b" {
			continue
		}
		name := strings.TrimSpace(read(d, "manufacturer") + " " + read(d, "product"))
		devs = append(devs, usbDev{v, read(d, "idProduct"), name})
	}
	return devs
}

// attachedUSB returns vendor:product ids of USB hostdevs in the VM's live (or inactive) XML.
func attachedUSB(vm string) map[string]bool {
	out, err := virsh("dumpxml", vm)
	set := map[string]bool{}
	if err != nil {
		return set
	}
	var dom struct {
		Hostdevs []struct {
			Type   string `xml:"type,attr"`
			Vendor struct {
				ID string `xml:"id,attr"`
			} `xml:"source>vendor"`
			Product struct {
				ID string `xml:"id,attr"`
			} `xml:"source>product"`
		} `xml:"devices>hostdev"`
	}
	xml.Unmarshal([]byte(out), &dom)
	for _, h := range dom.Hostdevs {
		if h.Type == "usb" {
			set[strings.TrimPrefix(h.Vendor.ID, "0x")+":"+strings.TrimPrefix(h.Product.ID, "0x")] = true
		}
	}
	return set
}

// toggle runs attach-device or detach-device with a hostdev XML for d.
// --current = live if the VM is running, persistent config otherwise.
func toggle(cmd, vm string, d usbDev) error {
	x := fmt.Sprintf(`<hostdev mode='subsystem' type='usb' managed='yes'><source><vendor id='0x%s'/><product id='0x%s'/></source></hostdev>`, d.Vendor, d.Product)
	f, err := os.CreateTemp("", "virtusb-*.xml")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	f.WriteString(x)
	f.Close()
	_, err = virsh(cmd, vm, f.Name(), "--current")
	return err
}

func main() {
	a := app.New()
	w := a.NewWindow("virtusb")

	var devs []usbDev
	var attached map[string]bool
	vm := ""

	var list *widget.List
	list = widget.NewList(
		func() int { return len(devs) },
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("Detach", nil), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			d := devs[i]
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s  %s", d.id(), d.Name))
			btn := c.Objects[1].(*widget.Button)
			cmd, label := "attach-device", "Attach"
			if attached[d.id()] {
				cmd, label = "detach-device", "Detach"
			}
			btn.SetText(label)
			btn.OnTapped = func() {
				if vm == "" {
					return
				}
				btn.Disable()
				go func() {
					err := toggle(cmd, vm, d)
					fyne.Do(func() {
						btn.Enable()
						if err != nil {
							dialog.ShowError(err, w)
						}
						attached = attachedUSB(vm)
						list.Refresh()
					})
				}()
			}
		},
	)

	refresh := func() {
		devs = hostUSB()
		attached = map[string]bool{}
		if vm != "" {
			attached = attachedUSB(vm)
		}
		list.Refresh()
	}

	vmSel := widget.NewSelect(listVMs(), func(s string) { vm = s; refresh() })
	top := container.NewBorder(nil, nil, widget.NewLabel("VM:"),
		widget.NewButton("Refresh", func() { vmSel.Options = listVMs(); refresh() }), vmSel)

	w.SetContent(container.NewBorder(top, nil, nil, nil, list))
	w.Resize(fyne.NewSize(600, 400))
	refresh()
	w.ShowAndRun()
}
