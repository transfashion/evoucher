package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fgtago/fgweb"
	"github.com/go-chi/chi/v5"
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
