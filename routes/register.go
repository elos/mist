package routes

import (
	"log"

	"github.com/elos/ehttp/serve"
	"github.com/elos/ehttp/templates"
	"github.com/elos/mist/services"
	"github.com/elos/mist/views"
	"github.com/elos/models"
)

func RegisterGET(c *serve.Conn, engine services.Views) {
	templates.CatchError(c, engine.Execute(c, views.Register, nil))
}

func RegisterPOST(c *serve.Conn, db services.DB) {
	u := models.NewUser()
	u.SetID(db.NewID())
	err := db.Save(u)
	if err != nil {
		log.Printf("%s", err)
		c.Write([]byte(err.Error()))
		return
	}

	p := models.NewProfile()
	p.SetID(db.NewID())
	p.SetOwner(u)

	name := c.ParamVal("name")
	phone := c.ParamVal("phone")

	p.Name = name
	p.Phone = phone

	err = db.Save(p)
	if err != nil {
		log.Printf("%s", err)
		c.Write([]byte(err.Error()))
		return
	}

	c.Write([]byte("Heck Yeah"))
}
