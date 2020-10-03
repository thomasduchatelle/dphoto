package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"os"
	"reflect"
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
}

func (sig *UDiskSignal) String() string {
	return fmt.Sprintf("UDiskSignal[key=(%s,%s), drive=%s, filesystem=%s, block=%s]", sig.DrivePath, sig.BlockPath, sig.Drive, sig.FileSystem, sig.Block)
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

		} else {
			fmt.Println(logPrefix)
			fmt.Printf("\t- type=%s\n", reflect.TypeOf(body))
			fmt.Printf("\t- String()=%s\n", signal)
		}
	}

	return &udisk, udisk.Drive != nil || udisk.FileSystem != nil || udisk.Block != nil
}

func parsePropsToUDiskObject(udisk *UDiskSignal, dbusType string, props map[string]dbus.Variant, pathPrefix string) {
	switch dbusType {
	case "org.freedesktop.UDisks2.Drive":
		udisk.Drive = &UDiskDrive{getString(props, "Id", ""), getBool(props, "Ejectable", false)}

	case "org.freedesktop.UDisks2.Block":
		udisk.Block = &UDiskBlock{getString(props, "IdUUID", ""), getObjectPathAsString(props, "Drive", "")}

	case "org.freedesktop.UDisks2.Filesystem":
		udisk.FileSystem = &UDiskFilesystem{getStringArray(props, "MountPoints", nil)}

	default:
		fmt.Printf("?> %s -> %s\n", pathPrefix, dbusType)
		for key, value := range props {
			fmt.Printf("\t- %s=%s [signature=%s ; value type=%s]\n", key, value, value.Signature(), reflect.TypeOf(value.Value()))
		}
	}
}

func getString(props map[string]dbus.Variant, key string, defaultValue string) string {
	if val, present := props[key]; present {
		return val.String()
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

func main() {
	fmt.Println("Starting dphoto - delegates")

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err = conn.AddMatchSignal(
		dbus.WithMatchPathNamespace("/org/freedesktop/UDisks2"),
	); err != nil {
		panic(err)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		if signal, ok := parseDbusObject(v); ok {
			fmt.Print("\n\nReceived: ", signal, "\n\n\n")
		}
	}
}
