<p align="center"><img src="icon.svg" width="128" alt="virtusb logo"></p>

# virtusb

A small desktop GUI for attaching and detaching host USB devices to libvirt VMs. It is a front end for `virsh attach-device` / `detach-device`, so you don't have to write hostdev XML by hand.

## Features

- Lists the host's USB devices from sysfs (root hubs are hidden).
- Shows which devices are already attached to the selected VM. **Attach** buttons are blue and **Detach** buttons are red.
- If the VM is running, changes apply to it live. If it is shut off, they go into its persistent config (`--current`).
- Selects the VM for you when exactly one is running.

## Requirements

- Linux with libvirt and `virsh` on `PATH`
- Permission to manage VMs on the libvirt URI, usually by being in the `libvirt` group

## Install

Download a binary for `amd64` or `arm64` from [Releases](https://github.com/Toshimichi0915/virtusb/releases), or build from source:

```sh
go build -o virtusb .
```

Building needs Go, a C compiler, and the OpenGL/X11/Wayland development headers that [Fyne](https://docs.fyne.io/started/) uses. `.github/workflows/release.yml` lists the Debian/Ubuntu packages.

## Usage

```sh
virtusb
```

Pick a VM, then click **Attach** or **Detach** next to a device. Click **Refresh** after you plug or unplug something.

virtusb connects to `qemu:///system` unless `LIBVIRT_DEFAULT_URI` is set:

```sh
LIBVIRT_DEFAULT_URI=qemu:///session virtusb
```

## Notes

Devices are matched by `vendor:product` ID. If two identical devices are plugged in, libvirt can't tell them apart.

## License

[MIT](LICENSE)
