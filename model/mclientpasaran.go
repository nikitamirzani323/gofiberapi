package model

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/gofiberapi/db"
	"github.com/nleeper/goment"
)

type Mclientpasaran struct {
	IdCompPasaran    int    `json:"IdCompPasaran"`
	IdComp           string `json:"IdCompany"`
	NmPasaran        string `json:"NamaPasaran"`
	PasaranJamTutup  string `json:"Pasaran_Tutup"`
	PasaranJamJadwal string `json:"Pasaran_Jadwal"`
	PasaranJamOpen   string `json:"Pasaran_Open"`
	StatusPasaran    string `json:"StatusPasaran"`
}

func FetchAll_MclientPasaran() (Response, error) {
	var obj Mclientpasaran
	var arraobj []Mclientpasaran
	var res Response
	con := db.CreateCon()

	sql := `SELECT 
		idcomppasaran, idcompany, 
		nmpasarantogel, jamtutup, jamjadwal, jamopen 
		FROM client_view_pasaran 
		WHERE statuspasaranactive = 'Y' 
	`
	rows, err := con.Query(sql)
	defer rows.Close()

	if err != nil {
		return res, err
	}
	for rows.Next() {
		var idcomppasaran int
		var idcompany, nmpasarantogel, jamtutup, jamjadwal, jamopen string

		err = rows.Scan(
			&idcomppasaran,
			&idcompany,
			&nmpasarantogel,
			&jamtutup,
			&jamjadwal,
			&jamopen)
		if err != nil {
			return res, err
		}
		obj.IdCompPasaran = idcomppasaran
		obj.IdComp = idcompany
		obj.NmPasaran = nmpasarantogel
		obj.PasaranJamTutup = jamtutup
		obj.PasaranJamJadwal = jamjadwal
		obj.PasaranJamOpen = jamopen
		obj.StatusPasaran = "ONLINE"
		arraobj = append(arraobj, obj)
	}
	tglnow, _ := goment.New()
	log.Println(arraobj)
	log.Println(tglnow.Format("YYYY-MM-DD HH:mm:ss"))

	for i := 0; i < len(arraobj); i++ {
		log.Println(arraobj[i].IdCompPasaran)
	}

	res.Status = fiber.StatusOK
	res.Message = "Success"
	res.Totalrecord = len(arraobj)
	res.Record = arraobj

	return res, nil
}
