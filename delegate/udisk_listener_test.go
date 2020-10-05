package main

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetectPlugAndUnplugDisk(t *testing.T) {
	assert := assert.New(t)
	log.SetLevel(log.TraceLevel)

	// when
	var removableDrives []*UDiskSignal
	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2", Name: "org.freedesktop.DBus.ObjectManager.InterfacesAdded", Body: []interface{}{
		dbus.ObjectPath("/org/freedesktop/UDisks2/drives/F8_T5_11100000917"),
		map[string]map[string]dbus.Variant{
			"org.freedesktop.UDisks2.Drive": {
				"Id":        dbus.MakeVariant("12345"),
				"Ejectable": dbus.MakeVariant(true),
			},
		},
	}}, removableDrives)

	// then
	if len(removableDrives) == 0 {
		t.Error("new drives should have been added")
		t.FailNow()
	}
	assert.Equal("/org/freedesktop/UDisks2/drives/F8_T5_11100000917", removableDrives[0].DrivePath, "DrivePath")
	assert.Equal(true, removableDrives[0].Drive.Ejectable, "Ejectable")

	// when
	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2", Name: "org.freedesktop.DBus.ObjectManager.InterfacesAdded", Body: []interface{}{
		dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sdb"),
		map[string]map[string]dbus.Variant{
			"org.freedesktop.UDisks2.Block": {
				"Drive":  dbus.MakeVariant(dbus.ObjectPath("/org/freedesktop/UDisks2/drives/F8_T5_11100000917")),
				"IdUUID": dbus.MakeVariant("001B-9622"),
			},
			"org.freedesktop.UDisks2.Filesystem": {
				"MountPoints": dbus.MakeVariant([][]uint8{}),
			},
		},
	}}, removableDrives)

	// then
	if len(removableDrives) != 1 {
		t.Error("update should have updated the list, not create a new record.", removableDrives)
		t.FailNow()
	}

	assert.Equal("/org/freedesktop/UDisks2/block_devices/sdb", removableDrives[0].BlockPath)
	assert.NotNil(removableDrives[0].Block)
	assert.Equal("001B-9622", removableDrives[0].Block.IdUuid)
	assert.NotNil(removableDrives[0].FileSystem)
	assert.Equal(0, len(removableDrives[0].FileSystem.MountPoint))

	// when
	const mountPoint = "run/media/stark/001B-9622"
	bytes := make([]uint8, len(mountPoint))
	for i, c := range mountPoint {
		bytes[i] = uint8(c)
	}
	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2/block_devices/sdb", Name: "org.freedesktop.DBus.Properties.PropertiesChanged", Body: []interface{}{
		"org.freedesktop.UDisks2.Filesystem",
		map[string]dbus.Variant{
			"MountPoints": dbus.MakeVariant([][]uint8{bytes}),
		},
	}}, removableDrives)

	// then
	if len(removableDrives) != 1 {
		t.Error("update should have updated the list, not create a new record.", removableDrives)
		t.FailNow()
	}

	assert.Equal([]string{mountPoint}, removableDrives[0].FileSystem.MountPoint)
}

func TestUmount(t *testing.T) {
	assert := assert.New(t)
	// given
	removableDrives := []*UDiskSignal{
		{
			DrivePath: "/org/freedesktop/UDisks2/drives/F8_T5_11100000917",
			BlockPath: "/org/freedesktop/UDisks2/block_devices/sdb",
			Drive: &UDiskDrive{
				Id:        "12345",
				Ejectable: true,
			},
			Block: &UDiskBlock{
				IdUuid: "001B-9622",
				Drive:  "/org/freedesktop/UDisks2/drives/F8_T5_11100000917",
			},
			FileSystem: &UDiskFilesystem{
				MountPoint: []string{"run/media/stark/001B-9622"},
			},
		},
	}

	// when
	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2", Name: "org.freedesktop.DBus.ObjectManager.InterfacesRemoved", Body: []interface{}{
		dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sdb"),
		[]string{
			"org.freedesktop.UDisks2.Filesystem",
			"org.freedesktop.UDisks2.Block",
		},
	}}, removableDrives)

	// then - note: Block is kept
	assert.Equal(1, len(removableDrives))
	assert.Nil(removableDrives[0].Drive)
}

func TestRemoved(t *testing.T) {
	assert := assert.New(t)
	// given
	removableDrives := []*UDiskSignal{
		{
			DrivePath: "/org/freedesktop/UDisks2/drives/F8_T5_11100000917",
			BlockPath: "/org/freedesktop/UDisks2/block_devices/sdb",
			Drive: &UDiskDrive{
				Id:        "12345",
				Ejectable: true,
			},
			Block: &UDiskBlock{
				IdUuid: "001B-9622",
				Drive:  "/org/freedesktop/UDisks2/drives/F8_T5_11100000917",
			},
			FileSystem: &UDiskFilesystem{
				MountPoint: []string{"run/media/stark/001B-9622"},
			},
		},
	}

	// when
	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2", Name: "org.freedesktop.DBus.ObjectManager.InterfacesRemoved", Body: []interface{}{
		dbus.ObjectPath("/org/freedesktop/UDisks2/drives/F8_T5_11100000917"),
		[]string{
			"org.freedesktop.UDisks2.Drive",
		},
	}}, removableDrives)

	// then
	assert.Equal(0, len(removableDrives))
}
