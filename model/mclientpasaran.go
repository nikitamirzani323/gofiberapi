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
type MclientpasaranResult struct {
	No      uint16 `json:"no"`
	Date    string `json:"date"`
	Periode string `json:"periode"`
	Result  string `json:"result"`
}
type MclientpasaranCheckPasaran struct {
	PasaranIdtansaction string `json:"pasaran_idtransaction"`
	PasaranName         string `json:"pasaran_name"`
	PasaranPeriode      string `json:"pasaran_periode"`
	PasaranIdcomp       string `json:"pasaran_idcomp"`
	PasaranStatus       string `json:"pasaran_status"`
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

	if len(arraobj) > 0 {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Totalrecord = len(arraobj)
		res.Record = arraobj
		res.Time = tglnow.Format("YYYY-MM-DD HH:mm:ss")
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = "Not Found"
		res.Totalrecord = 0
		res.Record = nil
		res.Time = tglnow.Format("YYYY-MM-DD HH:mm:ss")
	}

	return res, nil
}

func FetchAll_MclientPasaranResult(client_company, pasaran_code string) (Response, error) {
	var obj MclientpasaranResult
	var arraobj []MclientpasaranResult
	var res Response
	msg := "Error"
	con := db.CreateCon()
	tglnow, _ := goment.New()

	sqlresult := `SELECT 
		A.keluaranperiode, A.datekeluaran, 
		A.keluarantogel, B.idpasarantogel 
		FROM tbl_trx_keluarantogel as A 
		JOIN tbl_mst_company_game_pasaran as B ON B.idcomppasaran = A.idcomppasaran
		WHERE B.idcompany = ? 
		AND B.idpasarantogel = ?
	`
	rowresult, err := con.Query(sqlresult, client_company, pasaran_code)
	defer rowresult.Close()

	if err != nil {
		return res, err
	}
	var norecord uint16 = 1
	for rowresult.Next() {
		var (
			keluaranperiode                             string
			datekeluaran, keluarantogel, idpasarantogel string
		)

		err = rowresult.Scan(
			&keluaranperiode,
			&datekeluaran,
			&keluarantogel,
			&idpasarantogel)
		if err != nil {
			return res, err
		}

		obj.No = norecord
		obj.Date = datekeluaran
		obj.Periode = idpasarantogel + "-" + keluaranperiode
		obj.Result = keluarantogel
		arraobj = append(arraobj, obj)
		msg = "Success"
		norecord = norecord + 1
	}
	if len(arraobj) > 0 {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Totalrecord = len(arraobj)
		res.Record = arraobj
		res.Time = tglnow.Format("YYYY-MM-DD HH:mm:ss")
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = "Not Found"
		res.Totalrecord = 0
		res.Record = nil
		res.Time = tglnow.Format("YYYY-MM-DD HH:mm:ss")
	}

	return res, nil
}

func CheckPasaran(client_company, pasaran_code string) (Response, error) {
	var obj MclientpasaranCheckPasaran
	var arraobj []MclientpasaranCheckPasaran
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
		idcomppasaran, nmpasarantogel, 
		jamtutup, jamopen  
		FROM client_view_pasaran 
		WHERE idcompany = ? 
		AND idpasarantogel = ?
	`
	rowpasaran, err := con.Query(sqlpasaran, client_company, pasaran_code)
	defer rowpasaran.Close()

	if err != nil {
		return res, err
	}
	for rowpasaran.Next() {
		var (
			idcomppasaran, nmpasarantogel, jamtutup, jamopen string
			idtrxkeluaran, keluaranperiode, haripasaran      string
		)

		err = rowpasaran.Scan(
			&idcomppasaran,
			&nmpasarantogel,
			&jamtutup,
			&jamopen)
		if err != nil {
			return res, err
		}

		sqlkeluaran := `
			SELECT 
			idtrxkeluaran, keluaranperiode
			FROM 
				tbl_trx_keluarantogel 
			WHERE idcompany = ?
			AND idcomppasaran = ?
			AND keluarantogel = ''
			LIMIT 1
		`
		err := con.QueryRow(sqlkeluaran, client_company, idcomppasaran).Scan(&idtrxkeluaran, &keluaranperiode)
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

		obj.PasaranIdtansaction = idtrxkeluaran
		obj.PasaranName = nmpasarantogel
		obj.PasaranPeriode = keluaranperiode
		obj.PasaranIdcomp = idcomppasaran
		obj.PasaranStatus = statuspasaran
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	if len(arraobj) > 0 {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Totalrecord = len(arraobj)
		res.Record = arraobj
		res.Time = tglnow.Format("YYYY-MM-DD HH:mm:ss")
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = "Not Found"
		res.Totalrecord = 0
		res.Record = nil
		res.Time = tglnow.Format("YYYY-MM-DD HH:mm:ss")
	}

	return res, nil
}
