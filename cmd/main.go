// Copyright (c) 2020 Alec Randazzo

package main

import (
	"archive/zip"
	"fmt"
	"github.com/AlecRandazzo/Packrat/internal/collector"
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"time"
)

var CLI struct {
	Collect struct {
		Throttle bool   `short:"t" optional:"" help:"Throttle the process to a single thread."`
		Output   string `short:"o" optional:"" help:"Output file. If not specified, the file name defaults to the host name and a timestamp."`
		Debug    bool   `short:"d" optional:"" help:"Debug mode"`
	} `cmd help:"Collect forensic data."`
	Parse struct {
	} `cmd help:"Parse forensic data."`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "collect":
		if CLI.Collect.Throttle {
			runtime.GOMAXPROCS(1)
		}
		if CLI.Collect.Debug {
			debugLog, _ := os.Create("debug.log")
			log.SetOutput(debugLog)
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetOutput(os.Stdout)
			log.SetLevel(log.ErrorLevel)
		}

		systemDrive := os.Getenv("SYSTEMDRIVE")
		exportList := collector.FileExportList{
			{
				FullPath:      fmt.Sprintf(`%s\$MFT`, systemDrive),
				FullPathRegex: false,
				FileName:      `$MFT`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\Windows\\System32\\winevt\\Logs\\.*\.evtx$`, systemDrive),
				FullPathRegex: true,
				FileName:      `.*\.evtx$`,
				FileNameRegex: true,
			},
			{
				FullPath:      fmt.Sprintf(`%s\Windows\System32\config\SYSTEM`, systemDrive),
				FullPathRegex: false,
				FileName:      `SYSTEM`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\Windows\System32\config\SOFTWARE`, systemDrive),
				FullPathRegex: false,
				FileName:      `SOFTWARE`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\users\\([^\\]+)\\ntuser.dat`, systemDrive),
				FullPathRegex: true,
				FileName:      `ntuser.dat`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\usrclass.dat`, systemDrive),
				FullPathRegex: true,
				FileName:      `usrclass.dat`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\WebCache\\WebCacheV01.dat`, systemDrive),
				FullPathRegex: true,
				FileName:      `WebCacheV01.dat`,
				FileNameRegex: false,
			},
		}

		var zipName string
		if CLI.Collect.Output != "" {
			zipName = CLI.Collect.Output
		} else {
			hostName, _ := os.Hostname()
			zipName = fmt.Sprintf("%s_%s.zip", hostName, time.Now().Format("2006-01-02T15.04.05Z"))
		}
		fileHandle, err := os.Create(zipName)
		if err != nil {
			err = fmt.Errorf("failed to create zip file %s", zipName)
		}
		defer fileHandle.Close()

		zipWriter := zip.NewWriter(fileHandle)
		//resultWriter := collector.ZipResultWriter{
		//	ZipWriter:  zipWriter,
		//	FileHandle: fileHandle,
		//}
		defer zipWriter.Close()
		volumeHandler := collector.NewVolumeHandler(strings.Trim(os.Getenv("SYSTEMDRIVE"), ":"))
		err = collector.Collect(volumeHandler, exportList, zipWriter)
		if err != nil {
			log.Panic(err)
		}
	default:
		ctx.Command()
	}
}
