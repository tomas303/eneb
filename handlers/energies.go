package handlers

import (
	"database/sql"
	"eneb/data"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_energies(r *gin.Engine, db *sql.DB) {
	getprevparam := MakeGetQueryParAsInt("prev", 0, abortWith(http.StatusBadRequest))
	getnextparam := MakeGetQueryParAsInt("next", 10, abortWith(http.StatusBadRequest))
	getpinparam := MakeGetQueryParAsInt64("pin", math.MaxInt64, abortWith(http.StatusBadRequest))

	cmdSelectBefore, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where created < ?
		order by created desc limit ?`,
		false,
		func(row data.RowScanner) (*data.Energy, error) {
			en := data.NewEnergy()
			err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created)
			if err != nil {
				return nil, err
			}
			return &en, nil
		})
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany[*data.Energy]([]ParamGetterFunc{getpinparam, getprevparam}, cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where created > ?
		order by created limit ?`,
		false,
		func(row data.RowScanner) (*data.Energy, error) {
			en := data.NewEnergy()
			err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created)
			if err != nil {
				return nil, err
			}
			return &en, nil
		})
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany[*data.Energy]([]ParamGetterFunc{getpinparam, getnextparam}, cmdSelectAfter)

	r.GET("/energies",
		func(c *gin.Context) {
			prev := getprevparam(c).(int)
			next := getnextparam(c).(int)
			pin := getpinparam(c)
			if prev != 0 && next != 0 {
				c.AbortWithError(http.StatusBadRequest, paramErr{message: "cannot specify both prev and next parameter"})
			}
			if (prev != 0 || next != 0) && pin == 0 {
				c.AbortWithError(http.StatusBadRequest, paramErr{message: "for prev or next parameter the pin parameter is mandatory"})
			}
			if prev > 0 {
				beforeHandler(c)
			} else if next > 0 {
				afterHandler(c)
			} else {
				c.AbortWithError(400, paramErr{message: "nor prev nor next parameter specified"})
			}
		})

	cmdSave, err := data.MakeDataCmdSaveOne[*data.Energy](db,
		"insert or replace into energies(id, kind, amount, info, created) VALUES(?,?,?,?,?)",
		func(en *data.Energy) []any {
			return []any{en.ID, en.Kind, en.Amount, en.Info, en.Created.Val}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne[*data.Energy](cmdSave)

	r.POST("/energies", postHandler)
}
