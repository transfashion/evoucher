package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fgtago/fgweb"
	"github.com/go-chi/chi/v5"
	"github.com/transfashion/evoucher/custdb"
	"github.com/transfashion/evoucher/models"
	"github.com/transfashion/evoucher/qiscus"
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

	// setup application data
	ws.ApplicationData = &models.ApplicationData{
		Qiscus: &qiscus.Qiscus{
			BaseUrl: appcfg.QiscusConfig.BaseUrl,
			AppCode: appcfg.QiscusConfig.AppCode,
			Secret:  appcfg.QiscusConfig.Secret,
			Sender:  appcfg.QiscusConfig.Sender,
		},
		CustomerDb: custdb.NewCustomerDB(),
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
