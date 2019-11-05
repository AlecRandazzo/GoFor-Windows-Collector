package windowscollector

import (
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/windows"
	"io"
	"os"
)

type DataRunsReader struct {
	VolumeHandler          *VolumeHandler
	DataRuns               mft.DataRuns
	fileName               string
	dataRunTracker         int
	bytesLeftToReadTracker int64
	initialized            bool
}

func (dataRunReader *DataRunsReader) Read(byteSliceToPopulate []byte) (numberOfBytesRead int, err error) {
	numberOfBytesToRead := len(byteSliceToPopulate)

	// Sanity checking
	if len(dataRunReader.DataRuns) == 0 {
		err = io.ErrUnexpectedEOF
		log.Warnf("failed to read %s, received: %v", dataRunReader.fileName, err)
		return
	}

	// Check if this reader has been initialized, if not, do so.
	if dataRunReader.initialized != true {
		dataRunReader.dataRunTracker = 0
		dataRunReader.bytesLeftToReadTracker = dataRunReader.DataRuns[dataRunReader.dataRunTracker].Length
		dataRunReader.VolumeHandler.lastReadVolumeOffset, _ = syscall.Seek(dataRunReader.VolumeHandler.Handle, dataRunReader.DataRuns[dataRunReader.dataRunTracker].AbsoluteOffset, 0)
		dataRunReader.VolumeHandler.lastReadVolumeOffset -= int64(numberOfBytesToRead)
		dataRunReader.initialized = true

		// These are for debug purposes
		if log.GetLevel() == log.DebugLevel {
			totalSize := int64(0)
			for _, dataRun := range dataRunReader.DataRuns {
				totalSize += dataRun.Length
			}
			log.Debugf("Reading data run number 1 of %d for file '%s' which has a length of %d bytes at absolute offset %d",
				len(dataRunReader.DataRuns),
				dataRunReader.fileName,
				totalSize,
				dataRunReader.DataRuns[0].AbsoluteOffset,
			)
		}

	}

	// Figure out how many bytes are left to read
	if dataRunReader.bytesLeftToReadTracker-int64(numberOfBytesToRead) == 0 {
		dataRunReader.bytesLeftToReadTracker -= int64(numberOfBytesToRead)
	} else if dataRunReader.bytesLeftToReadTracker-int64(numberOfBytesToRead) < 0 {
		numberOfBytesToRead = int(dataRunReader.bytesLeftToReadTracker)
		dataRunReader.bytesLeftToReadTracker = 0
	} else {
		dataRunReader.bytesLeftToReadTracker -= int64(numberOfBytesToRead)
	}

	// Read from the data run
	buffer := make([]byte, numberOfBytesToRead)
	dataRunReader.VolumeHandler.lastReadVolumeOffset += int64(len(buffer))
	numberOfBytesRead, _ = syscall.Read(dataRunReader.VolumeHandler.Handle, buffer)
	copy(byteSliceToPopulate, buffer)

	// Check to see if there are any bytes left to read in the current data run
	if dataRunReader.bytesLeftToReadTracker == 0 {
		// Check to see if we have read all the data runs.
		if dataRunReader.dataRunTracker+1 == len(dataRunReader.DataRuns) {
			err = io.EOF
			return
		}

		// Increment our tracker
		dataRunReader.dataRunTracker++

		// Get the size of the next datarun
		dataRunReader.bytesLeftToReadTracker = dataRunReader.DataRuns[dataRunReader.dataRunTracker].Length

		// Seek to the offset of the next datarun
		dataRunReader.VolumeHandler.lastReadVolumeOffset, _ = syscall.Seek(dataRunReader.VolumeHandler.Handle, dataRunReader.DataRuns[dataRunReader.dataRunTracker].AbsoluteOffset, 0)
		dataRunReader.VolumeHandler.lastReadVolumeOffset -= int64(len(buffer))

		log.Debugf("Reading data run number %d of %d for file '%s' which has a length of %d bytes at absolute offset %d",
			dataRunReader.dataRunTracker+1,
			len(dataRunReader.DataRuns),
			dataRunReader.fileName,
			dataRunReader.DataRuns[dataRunReader.dataRunTracker].Length,
			dataRunReader.VolumeHandler.lastReadVolumeOffset+int64(len(buffer)),
		)
	}

	return
}

func apiFileReader(file foundFile) (reader io.Reader, err error) {
	reader, err = os.Open(file.fullPath)
	return
}

func rawFileReader(handler *VolumeHandler, file foundFile) (reader io.Reader) {
	reader = &DataRunsReader{
		VolumeHandler:          handler,
		DataRuns:               file.dataRuns,
		fileName:               file.fullPath,
		dataRunTracker:         0,
		bytesLeftToReadTracker: 0,
		initialized:            false,
	}
	return
}
