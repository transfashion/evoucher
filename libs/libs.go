package libs

import (
	"fmt"

	"github.com/fgtago/fgweb/appsmodel"
	"github.com/transfashion/evoucher/libs/custdb"
	"github.com/transfashion/evoucher/libs/kalistadb"
	"github.com/transfashion/evoucher/libs/qiscus"
	"github.com/transfashion/evoucher/libs/voucher"
	"github.com/transfashion/evoucher/models"
)

var ws *appsmodel.Webservice

var Qiscus *qiscus.Qiscus
var CustomerDb *custdb.CustomerDB
var VoucherDb *voucher.VoucherDB

func Load(w *appsmodel.Webservice) {
	ws = w
	appcfg := ws.ApplicationConfig.(*models.ApplicationConfig)

	fmt.Println("Connecting to Kalista Database...")
	err := kalistadb.ConnectDatabase(
		appcfg.Kalista.Database.Server,
		appcfg.Kalista.Database.Name,
		appcfg.Kalista.Database.Port,
		appcfg.Kalista.Database.Username,
		appcfg.Kalista.Database.Password)
	if err != nil {
		panic(err)
	}
	fmt.Println("Kalista Database Connected.")

	Qiscus = qiscus.NewQiscus(&qiscus.QiscusConfig{
		BaseUrl: appcfg.QiscusConfig.BaseUrl,
		AppCode: appcfg.QiscusConfig.AppCode,
		Secret:  appcfg.QiscusConfig.Secret,
		Sender:  appcfg.QiscusConfig.Sender,
	})

	CustomerDb = custdb.NewCustomerDB(kalistadb.GetConnection())
	VoucherDb = voucher.NewVoucherDB(kalistadb.GetConnection())

}
