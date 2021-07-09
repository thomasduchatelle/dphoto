package daemon

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/mocks"
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var volumeManagerPort *mocks.VolumeManagerPort

func init() {
	log.SetLevel(log.TraceLevel)
}

func TestDetectPlugAndUnplugDisk(t *testing.T) {
	a := assert.New(t)

	volumeManagerPort = new(mocks.VolumeManagerPort)
	VolumeManager = volumeManagerPort

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
	if !a.Len(removableDrives, 1, "new drives should have been added") {
		t.FailNow()
	}
	a.Equal("/org/freedesktop/UDisks2/drives/F8_T5_11100000917", removableDrives[0].DrivePath, "DrivePath")
	a.Equal(true, removableDrives[0].Drive.Ejectable, "Ejectable")

	// and
	volumeManagerPort.AssertNotCalled(t, "OnMountedVolume", mock.Anything)

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
	if !a.Len(removableDrives, 1, "update should have updated the list, not create a new record.") {
		t.FailNow()
	}

	a.Equal("/org/freedesktop/UDisks2/block_devices/sdb", removableDrives[0].BlockPath)
	a.NotNil(removableDrives[0].Block)
	a.Equal("001B-9622", removableDrives[0].Block.IdUuid)
	a.NotNil(removableDrives[0].FileSystem)
	a.Equal(0, len(removableDrives[0].FileSystem.MountPoint))

	// and
	volumeManagerPort.AssertNotCalled(t, "OnMountedVolume", mock.Anything)

	// when
	const mountPoint = "run/media/stark/001B-9622"
	bytes := make([]uint8, len(mountPoint))
	for i, c := range mountPoint {
		bytes[i] = uint8(c)
	}

	volumeManagerPort.On("OnMountedVolume", backupmodel.VolumeToBackup{
		UniqueId: "001B-9622",
		Path:     mountPoint,
	}).Return()

	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2/block_devices/sdb", Name: "org.freedesktop.DBus.Properties.PropertiesChanged", Body: []interface{}{
		"org.freedesktop.UDisks2.Filesystem",
		map[string]dbus.Variant{
			"MountPoints": dbus.MakeVariant([][]uint8{bytes}),
		},
	}}, removableDrives)

	// then
	if !a.Len(removableDrives, 1, "update should have updated the list, not create a new record.") {
		t.FailNow()
	}

	a.Equal([]string{mountPoint}, removableDrives[0].FileSystem.MountPoint)

	// and
	volumeManagerPort.AssertExpectations(t)
}

func TestUmount(t *testing.T) {
	a := assert.New(t)

	volumeManagerPort = new(mocks.VolumeManagerPort)
	VolumeManager = volumeManagerPort

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

	volumeManagerPort.On("OnUnMountedVolume", "001B-9622").Return()

	// when
	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2", Name: "org.freedesktop.DBus.ObjectManager.InterfacesRemoved", Body: []interface{}{
		dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sdb"),
		[]string{
			"org.freedesktop.UDisks2.Filesystem",
			"org.freedesktop.UDisks2.Block",
		},
	}}, removableDrives)

	// then - note: Block is kept
	a.Equal(1, len(removableDrives))
	a.Nil(removableDrives[0].Drive)

	volumeManagerPort.AssertExpectations(t)
}

func TestRemoved(t *testing.T) {
	a := assert.New(t)

	volumeManagerPort = new(mocks.VolumeManagerPort)
	VolumeManager = volumeManagerPort

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

	volumeManagerPort.On("OnUnMountedVolume", "001B-9622").Return()

	// when
	removableDrives = HandleDBusSignal(&dbus.Signal{Sender: "unit", Path: "/org/freedesktop/UDisks2", Name: "org.freedesktop.DBus.ObjectManager.InterfacesRemoved", Body: []interface{}{
		dbus.ObjectPath("/org/freedesktop/UDisks2/drives/F8_T5_11100000917"),
		[]string{
			"org.freedesktop.UDisks2.Drive",
		},
	}}, removableDrives)

	// then
	a.Empty(removableDrives)
	volumeManagerPort.AssertExpectations(t)
}
