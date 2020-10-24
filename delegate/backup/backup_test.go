package backup

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestStartBackupRunner(t *testing.T) {
	a := assert.New(t)

	var backup []LocalMedia
	var signatures []SimpleMediaSignature

	// given
	fileHandlerAdapter, imageDetailsReaderAdapter, onlineStorageAdapter, volumeRepositoryAdapter := mockAdapters()
	LocalMediaPath = path.Join(os.TempDir(), "dphoto_"+time.Now().Format("20060102_1504050700"))

	volume := RemovableVolume{
		UniqueId:   "disk-uuid-1",
		MountPaths: []string{"/mnt/disk1", "/dev/automount/disk-uuid-1"},
	}

	// and
	volumeRepositoryAdapter.On("RestoreLastSnapshot", volume.UniqueId).Maybe().Return([]SimpleMediaSignature{
		{"photos/image1.png", 42},
		{"photos/image2.png", 48},
	}, nil)

	// and
	fileHandlerAdapter.On("FindMediaRecursively", "/mnt/disk1", mock.Anything).Once().Run(func(args mock.Arguments) {
		channel := args.Get(1).(chan FoundMedia)
		channel <- NewFoundMedia(IMAGE, "/mnt/disk1", "photos/image1.png", 42)
		channel <- NewFoundMedia(IMAGE, "/mnt/disk1", "photos/image2.png", 42)
		channel <- NewFoundMedia(VIDEO, "/mnt/disk1", "video/video001.avi", 42)
	}).Return(nil)

	fileHandlerAdapter.On("CopyToLocal", "/mnt/disk1/photos/image2.png", mockHasSuffix("photos/image2.png")).Return("image2-sha256-1", nil)
	fileHandlerAdapter.On("CopyToLocal", "/mnt/disk1/video/video001.avi", mockHasSuffix("video/video001.avi")).Return("video1-sha256-2", nil)

	// and
	imageDetailsReaderAdapter.On("ReadImageDetails", mockHasSuffix("photos/image2.png")).Return(&MediaDetails{
		Width:    1024,
		Height:   768,
		DateTime: time.Date(2020, 10, 8, 18, 48, 0, 0, time.UTC),
		Make:     "Barry",
		Model:    "Allen",
	}, nil)

	onlineStorageAdapter.On("BackupOnline", mock.Anything).Run(func(args mock.Arguments) {
		origin := args.Get(0).(chan LocalMedia)
		for m := range origin {
			backup = append(backup, m)
		}
	}).Return(nil)

	volumeRepositoryAdapter.On("StoreSnapshot", volume.UniqueId, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		signatures = args.Get(2).([]SimpleMediaSignature)
	}).Return(nil)

	// when
	runner, err := StartBackupRunner(volume)
	if a.NoError(err) {
		report := <-runner.CompletionChannel

		a.Empty(report.Errors)

		a.Equal([3]uint32{2, 1, 1}, report.Counter.found)
		a.Equal(2, len(backup), backup)
		a.Equal(3, len(signatures), signatures)
	}
}

func TestStartBackupRunnerWithErrors(t *testing.T) {
	a := assert.New(t)

	var backup []LocalMedia
	var signatures []SimpleMediaSignature

	// given
	fileHandlerAdapter, imageDetailsReaderAdapter, onlineStorageAdapter, volumeRepositoryAdapter := mockAdapters()
	LocalMediaPath = path.Join(os.TempDir(), "dphoto_"+time.Now().Format("20060102_1504050700"))

	volume := RemovableVolume{
		UniqueId:   "disk-uuid-1",
		MountPaths: []string{"/mnt/disk1", "/dev/automount/disk-uuid-1"},
	}

	// and
	volumeRepositoryAdapter.On("RestoreLastSnapshot", volume.UniqueId).Maybe().Return([]SimpleMediaSignature{}, nil)

	// and
	fileHandlerAdapter.On("FindMediaRecursively", "/mnt/disk1", mock.Anything).Once().Run(func(args mock.Arguments) {
		channel := args.Get(1).(chan FoundMedia)
		channel <- NewFoundMedia(IMAGE, "/mnt/disk1", "photos/image1.png", 42)
		channel <- NewFoundMedia(IMAGE, "/mnt/disk1", "photos/image2.png", 42)
	}).Return(errors.Errorf("[test] FindMediaRecursively reports an error"))

	fileHandlerAdapter.On("CopyToLocal", "/mnt/disk1/photos/image1.png", mockHasSuffix("photos/image1.png")).Return("", errors.Errorf("[test] CopyToLocal reports an error"))
	fileHandlerAdapter.On("CopyToLocal", "/mnt/disk1/photos/image2.png", mockHasSuffix("photos/image2.png")).Return("image2-sha256-1", nil)

	// and
	imageDetailsReaderAdapter.On("ReadImageDetails", mockHasSuffix("photos/image2.png")).Return(nil, errors.Errorf("[test] ReadImageDetails reports an error"))

	onlineStorageAdapter.On("BackupOnline", mock.Anything).Run(func(args mock.Arguments) {
		origin := args.Get(0).(chan LocalMedia)
		for m := range origin {
			backup = append(backup, m)
		}
	}).Return(errors.Errorf("[test] BackupOnline reports an error"))

	volumeRepositoryAdapter.On("StoreSnapshot", volume.UniqueId, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		signatures = args.Get(2).([]SimpleMediaSignature)
	}).Return(errors.Errorf("[test] StoreSnapshot reports an error"))

	// when
	runner, err := StartBackupRunner(volume)
	if a.NoError(err) {
		report := <-runner.CompletionChannel

		a.Len(report.Errors, 3)
		sort.Slice(report.Errors, func(i, j int) bool {
			return report.Errors[i].Error() < report.Errors[j].Error()
		})
		a.Contains(report.Errors[0].Error(), "[test] BackupOnline reports an error")
		a.Contains(report.Errors[1].Error(), "[test] CopyToLocal reports an error")
		a.Contains(report.Errors[2].Error(), "[test] FindMediaRecursively reports an error")

		a.Equal([3]uint32{2, 2, 0}, report.Counter.found)
		a.Equal(1, len(backup), backup)
		a.Equal(0, len(signatures), signatures)
	}
}

func mockHasSuffix(suffix string) interface{} {
	return mock.MatchedBy(func(dest string) bool {
		return strings.HasSuffix(dest, suffix)
	})
}

func mockAdapters() (*MockFileHandlerAdapter, *MockImageDetailsReaderAdapter, *MockOnlineStorageAdapter, *MockVolumeRepositoryAdapter) {
	fileHandlerMock := new(MockFileHandlerAdapter)
	FileHandler = fileHandlerMock

	imageDetailsReaderMock := new(MockImageDetailsReaderAdapter)
	ImageDetailsReader = imageDetailsReaderMock

	onlineStorageMock := new(MockOnlineStorageAdapter)
	OnlineStorage = onlineStorageMock

	volumeRepositoryMock := new(MockVolumeRepositoryAdapter)
	VolumeRepository = volumeRepositoryMock

	return fileHandlerMock, imageDetailsReaderMock, onlineStorageMock, volumeRepositoryMock
}
