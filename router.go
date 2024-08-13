package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fgtago/fgweb"
	"github.com/fgtago/fgweb/appsmodel"
	"github.com/fgtago/fgweb/defaulthandlers"
	"github.com/fgtago/fgweb/midware"
	"github.com/go-chi/chi/v5"
	"github.com/transfashion/evoucher/apis"
	"github.com/transfashion/evoucher/handlers"
)

func Router(mux *chi.Mux) error {
	ws := appsmodel.GetWebservice()

	mux.Use(PageSetup)

	// Default handler
	fgweb.Get(mux, "/favicon.ico", defaulthandlers.FaviconHandler)
	fgweb.Get(mux, "/asset/*", defaulthandlers.AssetHandler)
	fgweb.Get(mux, "/template/*", defaulthandlers.TemplateHandler)

	hnd := handlers.New(ws)
	fgweb.Get(mux, "/", hnd.Home)
	fgweb.Get(mux, "/{vouid}/voucherqr.svg", hnd.VoucherQrSVG)
	fgweb.Get(mux, "/{vouid}/voucherqr.png", hnd.VoucherQrPNG)
	fgweb.Get(mux, "/{vouid}/voucherqr.jpg", hnd.VoucherQrJPG)

	fgweb.Get(mux, "/form", hnd.Form)
	fgweb.Post(mux, "/form", hnd.Form)

	fgweb.Get(mux, "/result", hnd.Result)
	fgweb.Get(mux, "/sent", hnd.Sent)
	fgweb.Get(mux, "/view/{vouid}", hnd.View)

	api := apis.New(appsmodel.GetWebservice())
	fgweb.Post(mux, "/api/requestvoucher", api.RequestVoucher)

	stadir := filepath.Join(ws.RootDir, "data", "vouchers")
	fileServer := http.FileServer(http.Dir(stadir))
	mux.Handle("/vouchers/*", noDirListing(http.StripPrefix("/vouchers", fileServer)))

	return nil
}

func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func PageSetup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if midware.IsAsset(r.URL.Path) || midware.IsTemplate(r.URL.Path) {
			next.ServeHTTP(w, r)
		} else {
			ctx := r.Context()
			pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)

			var base_href string
			if r.Header.Get("Base_href") != "" {
				base_href = r.Header.Get("Base_href")
			} else {
				base_href = "/"
			}

			pv.Setup = &handlers.PageSetup{
				BaseUrl:        base_href,
				ShowHeader:     true,
				ShowFooter:     true,
				ShowFooterRow3: false,
			}
			next.ServeHTTP(w, r)
		}

	})
}
