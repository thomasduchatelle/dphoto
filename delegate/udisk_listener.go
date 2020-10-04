package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

type UDiskDrive struct {
	Id        string
	Ejectable bool
}

type UDiskBlock struct {
	IdUuid string
	Drive  string
}

type UDiskFilesystem struct {
	MountPoint []string
}

type UDiskSignal struct {
	DrivePath  string
	BlockPath  string
	Drive      *UDiskDrive
	Block      *UDiskBlock
	FileSystem *UDiskFilesystem

	removed   bool
	unmounted bool
}

func (sig *UDiskSignal) String() string {
	return fmt.Sprintf("UDiskSignal[key=(%s,%s), drive=%s, filesystem=%s, block=%s]", sig.DrivePath, sig.BlockPath, sig.Drive, sig.FileSystem, sig.Block)
}

func (sig *UDiskSignal) merge(update *UDiskSignal) {
	if sig.DrivePath == "" {
		sig.DrivePath = update.DrivePath
	}
	if sig.BlockPath == "" {
		sig.BlockPath = update.BlockPath
	}

	if sig.Block == nil {
		sig.Block = update.Block
	}
	if update.unmounted {
		sig.Drive = nil
		sig.FileSystem = nil
	} else {
		if sig.Drive == nil {
			sig.Drive = update.Drive
		}
		if sig.FileSystem == nil {
			sig.FileSystem = update.FileSystem
		} else if update.FileSystem != nil {
			sig.FileSystem.MountPoint = update.FileSystem.MountPoint
		}
	}
}

func (drive *UDiskDrive) String() string {
	return fmt.Sprintf("UDiskDrive[Id=%s, Ejectable=%t]", drive.Id, drive.Ejectable)
}

func (block *UDiskBlock) String() string {
	return fmt.Sprintf("UDiskBlock[IdUuid=%s, Drive=%s]", block.IdUuid, block.Drive)
}

func (fs *UDiskFilesystem) String() string {
	return fmt.Sprintf("UDiskFilesystem[mountpoints=%s]", fs.MountPoint)
}

func parseDbusObject(signal *dbus.Signal) (*UDiskSignal, bool) {
	udisk := UDiskSignal{}
	var header string

	log.WithField("signal", signal).Traceln("Received UDisk message")

	for idx, body := range signal.Body {
		logPrefix := fmt.Sprintf("?> %s -> %s#%d", signal.Path, header, idx)

		if objectPath, ok := body.(dbus.ObjectPath); idx == 0 && ok {
			header = string(objectPath)

		} else if objectPath, ok := body.(string); idx == 0 && ok && signal.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			header = objectPath

		} else if obj, ok := body.(map[string]map[string]dbus.Variant); ok {
			for dbusType, props := range obj {
				parsePropsToUDiskObject(&udisk, dbusType, props, logPrefix)
			}
			if udisk.Block != nil {
				udisk.DrivePath = udisk.Block.Drive
				udisk.BlockPath = header
			} else if udisk.Drive != nil {
				udisk.DrivePath = header
			}

		} else if props, ok := body.(map[string]dbus.Variant); ok && signal.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			parsePropsToUDiskObject(&udisk, header, props, logPrefix)
			udisk.BlockPath = string(signal.Path)

		} else if removedInterfaces, ok := body.([]string); ok && signal.Name == "org.freedesktop.DBus.ObjectManager.InterfacesRemoved" {
			for _, removed := range removedInterfaces {
				if removed == "org.freedesktop.UDisks2.Block" {
					udisk.BlockPath = header
					udisk.unmounted = true
				} else if removed == "org.freedesktop.UDisks2.Drive" {
					udisk.DrivePath = header
					udisk.removed = true
				}
			}
		}
	}

	return &udisk, udisk.Drive != nil || udisk.FileSystem != nil || udisk.Block != nil || udisk.removed || udisk.unmounted
}

func parsePropsToUDiskObject(udisk *UDiskSignal, dbusType string, props map[string]dbus.Variant, pathPrefix string) {
	switch dbusType {
	case "org.freedesktop.UDisks2.Drive":
		udisk.Drive = &UDiskDrive{getString(props, "Id", ""), getBool(props, "Ejectable", false)}

	case "org.freedesktop.UDisks2.Block":
		udisk.Block = &UDiskBlock{getString(props, "IdUUID", ""), getObjectPathAsString(props, "Drive", "")}

	case "org.freedesktop.UDisks2.Filesystem":
		udisk.FileSystem = &UDiskFilesystem{getStringArray(props, "MountPoints", nil)}
	}
}

func getString(props map[string]dbus.Variant, key string, defaultValue string) string {
	if val, present := props[key]; present {
		if str, ok := val.Value().(string); ok {
			return str
		}
	}

	return defaultValue
}

func getBool(props map[string]dbus.Variant, key string, defaultValue bool) bool {
	if val, present := props[key]; present {
		if b, ok := val.Value().(bool); ok {
			return b
		}
	}

	return defaultValue
}

func getStringArray(props map[string]dbus.Variant, key string, defaultValue []string) []string {
	if val, present := props[key]; present {
		if array, ok := val.Value().([][]uint8); ok {
			var result []string
			for _, p := range array {
				result = append(result, string(p))
			}

			return result
		}
	}

	return defaultValue
}

func getObjectPathAsString(props map[string]dbus.Variant, key string, defaultValue string) string {
	if val, present := props[key]; present {
		if objectPath, ok := val.Value().(dbus.ObjectPath); ok {
			return string(objectPath)
		}
	}

	return defaultValue
}

func startUDiskListener() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatalln("Failed to connect to session bus:", err)
	}
	defer conn.Close()

	if err = conn.AddMatchSignal(
		dbus.WithMatchPathNamespace("/org/freedesktop/UDisks2"),
	); err != nil {
		panic(err)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	var removableDrives []*UDiskSignal

	for v := range c {
		if signal, ok := parseDbusObject(v); ok {
			var found *UDiskSignal
			wasPlugged := false

			for i := 0; i < len(removableDrives) && found == nil; i++ {
				if equalsIfNotNil(removableDrives[i].DrivePath, signal.DrivePath) || equalsIfNotNil(removableDrives[i].BlockPath, signal.BlockPath) {
					found = removableDrives[i]
					wasPlugged = isMounted(found)
					if signal.removed {
						removableDrives = append(removableDrives[:i], removableDrives[i+1:]...)
						found.removed = true
					} else {
						removableDrives[i].merge(signal)
					}
				}
			}

			if found == nil && !signal.removed {
				found = signal
				removableDrives = append(removableDrives, signal)
			}

			withContext := log.WithFields(log.Fields{})
			if found != nil && found.Block != nil {
				withContext = withContext.WithField("uuid", found.Block.IdUuid)
			}
			if found != nil && found.FileSystem != nil {
				withContext = withContext.WithField("mountPoints", found.FileSystem.MountPoint)
			}

			if found != nil && !wasPlugged && isMounted(found) {
				withContext.Infoln("Disk plugged")
			} else if found != nil && found.Block != nil && wasPlugged && !isMounted(found) {
				withContext.Infoln("Disk unplugged")
			} else if found != nil && found.removed {
				withContext.WithField("found", found).Debugln("Interface removed")
			} else {
				withContext.WithField("found", found).Debugln("Disk updated")
			}
		}
	}
}

func isMounted(signal *UDiskSignal) bool {
	return !signal.removed && signal.FileSystem != nil && len(signal.FileSystem.MountPoint) > 0
}

func equalsIfNotNil(value1 string, value2 string) bool {
	return value1 != "" && value2 != "" && value1 == value2
}
