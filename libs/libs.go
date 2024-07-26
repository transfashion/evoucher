package libs

import (
	"github.com/fgtago/fgweb/appsmodel"
	"github.com/transfashion/evoucher/libs/custdb"
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

	Qiscus = qiscus.NewQiscus(&qiscus.QiscusConfig{
		BaseUrl: appcfg.QiscusConfig.BaseUrl,
		AppCode: appcfg.QiscusConfig.AppCode,
		Secret:  appcfg.QiscusConfig.Secret,
		Sender:  appcfg.QiscusConfig.Sender,
	})

	CustomerDb = custdb.NewCustomerDB()
	VoucherDb = voucher.NewVoucherDB()

}
