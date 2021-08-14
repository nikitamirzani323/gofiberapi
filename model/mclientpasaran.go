package model

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/gofiberapi/db"
	"github.com/nleeper/goment"
)

type Mclientpasaran struct {
	PasaranId          string `json:"pasaran_id"`
	PasaranTogel       string `json:"pasaran_togel"`
	PasaranPeriode     string `json:"pasaran_periode"`
	PasaranTglKeluaran string `json:"pasaran_tglkeluaran"`
	PasaranJamTutup    string `json:"pasaran_marketclose"`
	PasaranJamJadwal   string `json:"pasaran_marketschedule"`
	PasaranJamOpen     string `json:"pasaran_marketopen"`
	PasaranStatus      string `json:"pasaran_status"`
}

func FetchAll_MclientPasaran(client_company string) (Response, error) {
	var obj Mclientpasaran
	var arraobj []Mclientpasaran
	var res Response
	var myDays = []string{"minggu", "senin", "selasa", "rabu", "kamis", "jumat", "sabtu"}
	statuspasaran := "ONLINE"
	msg := "Error"
	con := db.CreateCon()

	tglnow, _ := goment.New()
	daynow := tglnow.Format("d")
	intVar, _ := strconv.ParseInt(daynow, 0, 8)
	daynowhari := myDays[intVar]

	sqlpasaran := `SELECT 
		idcomppasaran, idpasarantogel, nmpasarantogel, jamtutup, jamjadwal, jamopen 
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
		var (
			idcomppasaran                                                int
			idpasarantogel, nmpasarantogel, jamtutup, jamjadwal, jamopen string
			tglkeluaran, periodekerluaran, haripasaran                   string
		)

		err = rowspasaran.Scan(
			&idcomppasaran,
			&idpasarantogel,
			&nmpasarantogel,
			&jamtutup,
			&jamjadwal,
			&jamopen)
		if err != nil {
			return res, err
		}

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
			return res, err
		}

		sqlpasaranonline := `
			SELECT
				haripasaran
			FROM
				tbl_mst_company_game_pasaran_offline
			WHERE idcomppasaran = ?
			AND idcompany = ? 
			AND haripasaran = ? 
		`

		errpasaranonline := con.QueryRow(sqlpasaranonline, idcomppasaran, client_company, daynowhari).Scan(&haripasaran)

		if errpasaranonline != sql.ErrNoRows {
			jamtutup := tglnow.Format("YYYY-MM-DD") + " " + jamtutup
			jamopen := tglnow.Format("YYYY-MM-DD") + " " + jamopen
			tutup, _ := goment.New(jamtutup)
			open, _ := goment.New(jamopen)
			nowconvert := tglnow.Format("x")
			tutupconvert := tutup.Format("x")
			openconvert := open.Format("x")

			intNow, _ := strconv.Atoi(nowconvert)
			intTutup, _ := strconv.Atoi(tutupconvert)
			intOpen, _ := strconv.Atoi(openconvert)

			if intNow >= intTutup && intNow <= intOpen {
				statuspasaran = "OFFLINE"
			}

		}

		obj.PasaranId = idpasarantogel
		obj.PasaranTogel = nmpasarantogel
		obj.PasaranPeriode = "#" + periodekerluaran + "-" + idpasarantogel
		obj.PasaranTglKeluaran = tglkeluaran
		obj.PasaranJamTutup = tglkeluaran + " " + jamtutup
		obj.PasaranJamJadwal = tglkeluaran + " " + jamjadwal
		obj.PasaranJamOpen = tglkeluaran + " " + jamopen
		obj.PasaranStatus = statuspasaran
		arraobj = append(arraobj, obj)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Totalrecord = len(arraobj)
	res.Record = arraobj
	res.Time = tglnow.Format("YYYY-MM-DD HH:mm:ss")

	return res, nil
}
