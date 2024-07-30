package dnsrcv

import app "github.com/ciphermarco/BOAST"

type ExportDNSHandler struct {
	dnsHandler
}

func NewExportDNSHandler(domain string, publicIP string, txt []string, strg app.Storage) *ExportDNSHandler {
	return &ExportDNSHandler{
		dnsHandler{
			domain:   domain,
			publicIP: publicIP,
			txt:      txt,
			storage:  strg,
		},
	}
}
