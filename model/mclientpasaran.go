package model

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/gofiberapi/db"
	"github.com/nleeper/goment"
)

type Mclientpasaran struct {
	IdCompPasaran    int    `json:"IdCompPasaran"`
	IdComp           string `json:"IdCompany"`
	NmPasaran        string `json:"NamaPasaran"`
	Periode          string `json:"Periode"`
	PasaranJamTutup  string `json:"Pasaran_Tutup"`
	PasaranJamJadwal string `json:"Pasaran_Jadwal"`
	PasaranJamOpen   string `json:"Pasaran_Open"`
	StatusPasaran    string `json:"StatusPasaran"`
}

func FetchAll_MclientPasaran(client_company string) (Response, error) {
	var obj Mclientpasaran
	var arraobj []Mclientpasaran
	var res Response
	msg := "Error"
	con := db.CreateCon()

	sqlpasaran := `SELECT 
		idcomppasaran, idcompany, idpasarantogel, 
		nmpasarantogel, jamtutup, jamjadwal, jamopen 
		FROM client_view_pasaran 
		WHERE statuspasaranactive = 'Y' 
		AND idcompany = ?
	`
	rowspasaran, err := con.Query(sqlpasaran, client_company)
	defer rowspasaran.Close()

	if err != nil {
		return res, err
	}
	for rowspasaran.Next() {
		var idcomppasaran int
		var idpasarantogel, idcompany, nmpasarantogel, jamtutup, jamjadwal, jamopen string

		err = rowspasaran.Scan(
			&idcomppasaran,
			&idcompany,
			&idpasarantogel,
			&nmpasarantogel,
			&jamtutup,
			&jamjadwal,
			&jamopen)
		if err != nil {
			return res, err
		}

		var tglkeluaran, periodekerluaran string

		sqlkeluaran := `
		SELECT 
		datekeluaran, keluaranperiode
		FROM 
			tbl_trx_keluarantogel 
		WHERE idcomppasaran = ?
		ORDER BY datekeluaran DESC
		LIMIT 1
		`
		err := con.QueryRow(sqlkeluaran, idcomppasaran).Scan(&tglkeluaran, &periodekerluaran)

		if err != nil {
			return res, errors.New("Not Found")
		}

		obj.IdCompPasaran = idcomppasaran
		obj.IdComp = idcompany
		obj.NmPasaran = nmpasarantogel
		obj.Periode = "#" + periodekerluaran + "-" + idpasarantogel
		obj.PasaranJamTutup = tglkeluaran + " " + jamtutup
		obj.PasaranJamJadwal = tglkeluaran + " " + jamjadwal
		obj.PasaranJamOpen = tglkeluaran + " " + jamopen
		obj.StatusPasaran = "ONLINE"
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	tglnow, _ := goment.New()
	// log.Println(arraobj)
	log.Println(tglnow.Format("YYYY-MM-DD HH:mm:ss"))

	// for i := 0; i < len(arraobj); i++ {
	// 	log.Println(arraobj[i].IdCompPasaran)
	// }

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Totalrecord = len(arraobj)
	res.Record = arraobj

	return res, nil
}
