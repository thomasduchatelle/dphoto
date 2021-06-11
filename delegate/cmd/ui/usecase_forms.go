package ui

import "github.com/logrusorgru/aurora/v3"

func CreateAlbumForm(operations InteractiveActionsPort, record Record) error {
	creation := RecordCreation{}
	var ok bool

	creation.Name, ok = ReadString("Name of the album", record.Name)
	if !ok {
		return nil
	}

	creation.FolderName, _ = ReadString("Folder name (leave blank for automatically generated)", "")

	creation.Start, ok = ReadDate("Start date", record.Start)
	if !ok {
		return nil
	}

	creation.End, ok = ReadDate("End date", record.End)
	if !ok {
		return nil
	}

	return operations.Create(creation)
}

func DeleteAlbum(operations InteractiveActionsPort, record Record) error {
	const pattern = "02/01/2006"
	proceed, ok := ReadBool(aurora.Sprintf("Are you sure you want to delete %s (%s) [%s -> %s] with %d medias in it?", aurora.Cyan(record.Name), record.FolderName, record.Start.Format(pattern), record.End.Format(pattern), record.Count), "y/N")
	if ok && proceed {
		return operations.DeleteAlbum(record.FolderName)
	}

	return nil
}

func EditAlbumDates(operations InteractiveActionsPort, record Record) error {
	start, ok := ReadDate("Start date", record.Start)
	if !ok {
		return nil
	}

	end, ok := ReadDate("End date", record.End)
	if !ok {
		return nil
	}

	return operations.UpdateAlbum(record.FolderName, start, end)
}

func EditAlbumName(operations InteractiveActionsPort, record Record) error {
	newName, ok := ReadString("Name of the album", record.Name)
	if !ok {
		return nil
	}

	if newName != record.Name {
		proceed, ok := ReadBool(aurora.Sprintf("Re-generate folder name /%s ?", aurora.Cyan(record.FolderName)), "[Y/n]")
		return operations.RenameAlbum(record.FolderName, newName, !ok || proceed)
	}

	return nil
}
