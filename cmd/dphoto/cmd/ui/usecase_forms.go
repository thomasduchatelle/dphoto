package ui

import "github.com/logrusorgru/aurora/v3"

type AlbumFormSession struct {
	actionPort CreateAlbumPort
	form       *FormUseCase
}
type DeleteSession struct {
	actionPort DeleteAlbumPort
	form       *FormUseCase
}
type EditDateSession struct {
	actionPort UpdateAlbumPort
	form       *FormUseCase
}
type RenameSession struct {
	actionPort RenameAlbumPort
	form       *FormUseCase
}

func NewCreateAlbumForm(actionsPort CreateAlbumPort, terminalPort PrintReadTerminalPort) *AlbumFormSession {
	return &AlbumFormSession{
		actionPort: actionsPort,
		form: &FormUseCase{
			TerminalPort: terminalPort,
		},
	}
}

func NewDeleteAlbumForm(actionPort DeleteAlbumPort, terminalPort PrintReadTerminalPort) *DeleteSession {
	return &DeleteSession{
		actionPort: actionPort,
		form: &FormUseCase{
			TerminalPort: terminalPort,
		},
	}
}

func NewEditAlbumDateForm(actionPort UpdateAlbumPort, terminalPort PrintReadTerminalPort) *EditDateSession {
	return &EditDateSession{
		actionPort: actionPort,
		form: &FormUseCase{
			TerminalPort: terminalPort,
		},
	}
}

func NewRenameAlbumForm(actionPort RenameAlbumPort, terminalPort PrintReadTerminalPort) *RenameSession {
	return &RenameSession{
		actionPort: actionPort,
		form: &FormUseCase{
			TerminalPort: terminalPort,
		},
	}
}

func (a *AlbumFormSession) AlbumForm(owner string, record Record) (bool, error) {
	creation := RecordCreation{
		Owner: owner,
	}
	var ok bool

	creation.Name, ok = a.form.ReadString("Name of the album", record.Name)
	if !ok {
		return false, nil
	}

	creation.FolderName, _ = a.form.ReadString("Folder name (leave blank for automatically generated)", "")

	creation.Start, ok = a.form.ReadDate("Start date", record.Start)
	if !ok {
		return false, nil
	}

	creation.End, ok = a.form.ReadDate("End date", record.End)
	if !ok {
		return false, nil
	}

	return true, a.actionPort.Create(creation)
}

func (s *DeleteSession) DeleteAlbum(record Record) (bool, error) {
	const pattern = "02/01/2006"
	proceed, ok := s.form.ReadBool(aurora.Sprintf("Are you sure you want to delete %s (%s) [%s -> %s] with %d medias in it?", aurora.Cyan(record.Name), record.FolderName, record.Start.Format(pattern), record.End.Format(pattern), record.Count), "y/N")
	if ok && proceed {
		return true, s.actionPort.DeleteAlbum(record.FolderName)
	}

	return false, nil
}

func (s *EditDateSession) EditAlbumDates(record Record) error {
	start, ok := s.form.ReadDate("Start date", record.Start)
	if !ok {
		return nil
	}

	end, ok := s.form.ReadDate("End date", record.End)
	if !ok {
		return nil
	}

	return s.actionPort.UpdateAlbum(record.FolderName, start, end)
}

func (s *RenameSession) EditAlbumName(record Record) error {
	newName, ok := s.form.ReadString("Name of the album", record.Name)
	if !ok {
		return nil
	}

	if newName != record.Name {
		proceed, ok := s.form.ReadBool(aurora.Sprintf("Re-generate folder name /%s", aurora.Cyan(record.FolderName)), "Y/n")
		return s.actionPort.RenameAlbum(record.FolderName, newName, !ok || proceed)
	}

	return nil
}
