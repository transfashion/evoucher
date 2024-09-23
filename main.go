package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fgtago/fgweb"
	"github.com/go-chi/chi/v5"
	"github.com/transfashion/evoucher/libs"
	"github.com/transfashion/evoucher/models"
	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("Starting Server eVoucher")

	// baca parameter dari CLI (untuk keperluan debug)
	var cfgFileName string
	flag.StringVar(&cfgFileName, "conf", "config.yml", "file konfigurasi yang akan di load")
	flag.Parse()

	// ambil informasi root direktori program dijalankan
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// path ke file konfigurasi
	cfgpath := filepath.Join(rootDir, cfgFileName)

	// inisiasi webservernya
	ws, err := fgweb.New(rootDir, cfgpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// extended config
	cfgcontent := *fgweb.GetCfgContent()
	appcfg := &models.ApplicationConfig{}
	err = yaml.Unmarshal(cfgcontent, appcfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ws.ApplicationConfig = appcfg

	// load libraries
	libs.Load(ws)

	// set output log
	if appcfg.Logging.Enabled {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
		if appcfg.Logging.Output != "" {
			logfilepath := filepath.Join(ws.RootDir, "data", "logs", appcfg.Logging.Output)
			fmt.Println("Logging to", logfilepath)
			f, err := os.OpenFile(logfilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
				os.Exit(1)
			}
			defer f.Close()
			log.SetOutput(f)
		} else {
			fmt.Println("log to screen")
		}
	} else {
		fmt.Println("Logging is disabled")
		log.SetOutput(io.Discard)
	}

	// logfilepath := filepath.Join(ws.RootDir, "data", "logs", "application.log")
	// f, err := os.OpenFile(logfilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	// jalankan service webserver
	port := ws.Configuration.Port
	fmt.Println("Running on port: ", port)
	err = fgweb.StartService(port, func(mux *chi.Mux) error {
		return Router(mux)
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
