package common

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	Port         = flag.Int("port", 7088, "the listening port")
	PrintVersion = flag.Bool("version", false, "print version and exit")
	PrintHelp    = flag.Bool("help", false, "print help and exit")
	LogDir       = flag.String("log-dir", "", "specify the log directory")
)

// UploadPath Maybe override by ENV_VAR
var UploadPath = "upload"

func printHelp() {
	fmt.Println("GPT Content Audit" + Version + "")
	fmt.Println("Copyright (C) 2024 Dean. All rights reserved.")
	//fmt.Println("GitHub: https://github.com/deanxv/gpt-content-audit ")
	fmt.Println("Usage: gpt-content-audit [--port <port>] [--log-dir <log directory>] [--version] [--help]")
}

func init() {
	flag.Parse()

	if *PrintVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *PrintHelp {
		printHelp()
		os.Exit(0)
	}

	if os.Getenv("UPLOAD_PATH") != "" {
		UploadPath = os.Getenv("UPLOAD_PATH")
	}
	if *LogDir != "" {
		var err error
		*LogDir, err = filepath.Abs(*LogDir)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stat(*LogDir); os.IsNotExist(err) {
			err = os.Mkdir(*LogDir, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
