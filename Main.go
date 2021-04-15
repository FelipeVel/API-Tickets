package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func obtenerBaseDeDatos() (db *sql.DB, e error) {
	usuario := "root"
	pass := "root"
	host := "tcp(127.0.0.1:3306)"
	nombreBaseDeDatos := "tickets"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		return nil, err
	}
	return db, nil
}

type ticket struct {
	idticket, usuario         string
	fCreacion, fActualizacion time.Time
	estatus                   bool
}

func crear(t ticket) (e error) {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("INSERT INTO tickets (idtickets, usuario, fCreacion, fActualizacion, estatus) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()

	_, err = sentenciaPreparada.Exec(t.idticket, t.usuario, t.fCreacion, time.Now(), t.estatus)
	if err != nil {
		return err
	}
	return nil
}

func recuperarTodosTickets() ([]ticket, error) {
	tickets := []ticket{}
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT idtickets, usuario, fCreacion, fActualizacion, estatus FROM tickets")

	if err != nil {
		return nil, err
	}

	defer filas.Close()

	var t ticket

	for filas.Next() {
		err = filas.Scan(&t.idticket, &t.usuario, &t.fCreacion, &t.fActualizacion, &t.estatus)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func recuperarUnTicket(id string) (ticket, error) {
	var t ticket
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return t, err
	}
	defer db.Close()
	err = db.QueryRow(fmt.Sprintf("SELECT idtickets, usuario, fCreacion, fActualizacion, estatus FROM tickets WHERE idtickets = '%s'", id)).Scan(&t.idticket, &t.usuario, &t.fCreacion, &t.fActualizacion, &t.estatus)

	if err != nil {
		return t, err
	}

	return t, err
}

func actualizar(t ticket) error {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("UPDATE tickets SET usuario = ?, fActualizacion = ?, estatus = ? WHERE idtickets = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	_, err = sentenciaPreparada.Exec(t.usuario, time.Now(), t.estatus, t.idticket)
	return err
}

func eliminar(id string) error {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = recuperarUnTicket(id)
	if err != nil {
		return err
	}

	sentenciaPreparada, err := db.Prepare("DELETE FROM tickets WHERE idtickets = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()

	_, err = sentenciaPreparada.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println(`Bienvenido`)
	menu_inicial := `

	¿Qué quieres hacer?
	
	[1] Crear ticket
	[2] Eliminar ticket
	[3] Editar ticket
	[4] Recuperar todos los tickets
	[5] Recuperar un ticket
	[0] Salir
	
	Digite la opción: `

	fmt.Printf(menu_inicial)

	var eleccion int
	fmt.Scan(&eleccion)
	for eleccion != 0 {
		switch eleccion {
		case 1:
			var t ticket
			fmt.Printf(`Ingrese el id del ticket: `)
			fmt.Scan(&t.idticket)
			fmt.Printf(`Ingrese el usuario del ticket: `)
			fmt.Scan(&t.usuario)
			fmt.Printf(`Ingrese el estatus del ticket (1:Abierto | 0:Cerrado): `)
			fmt.Scan(&t.estatus)
			t.fCreacion = time.Now()
			t.fActualizacion = time.Now()
			err := crear(t)
			if err != nil {
				fmt.Printf("Error creando ticket: %v\n", err)
			}
			fmt.Printf("Ticket creado\n")
			break
		case 2:
			var id string
			fmt.Printf(`Inserte el id del ticket que quiere eliminar:`)
			fmt.Scan(&id)
			err := eliminar(id)
			if err != nil {
				fmt.Printf("Error borrando ticket: %v\n", err)
				break
			}
			fmt.Printf("Ticket borrado\n")
			break
		case 3:
			var id string
			fmt.Println(`Ingrese el id del ticket que quiere actualizar: `)
			fmt.Scan(&id)
			t, err := recuperarUnTicket(id)
			if err != nil {
				fmt.Printf("Error recuperando ticket: %v", err)
				break
			}
			var seleccion1 int
			menuActualizacion := `
			Seleccione el dato que quiere actualizar:

			[1] Id
			[2] Usuario
			[3] Estatus

			Ingrese su seleccion: `
			fmt.Printf(menuActualizacion)
			fmt.Scan(&seleccion1)
			switch seleccion1 {
			case 1:
				var id string
				fmt.Printf(`Ingrese el nuevo id:`)
				fmt.Scan(&id)
				t.idticket = id
				break
			case 2:
				var usuario string
				fmt.Printf(`Ingrese el nuevo usuario:`)
				fmt.Scan(&usuario)
				t.usuario = usuario
				break
			case 3:
				var estatus bool
				fmt.Printf(`Ingrese el nuevo estatus:`)
				fmt.Scan(&estatus)
				t.estatus = estatus
				break
			}
			err = actualizar(t)
			if err != nil {
				fmt.Printf("Error actualizando datos: %v\n", err)
				break
			}
			fmt.Println("Datos actualizados correctamente")
			break
		case 4:
			tickets, err := recuperarTodosTickets()
			if err != nil {
				fmt.Printf("Error recuperando tickets: %v\n", err)
				break
			}
			for _, t := range tickets {
				fmt.Println("--------------------------------------------------")
				fmt.Printf("Id:%v\n", t.idticket)
				fmt.Printf("Usuario:%v\n", t.usuario)
				fmt.Printf("Fecha de creacion:%v\n", t.fCreacion)
				fmt.Printf("Fecha de actualizacion:%v\n", t.fActualizacion)
				fmt.Printf("Estatus:%v\n", t.estatus)
				fmt.Println("--------------------------------------------------")
			}
			break
		case 5:
			var id string
			fmt.Printf(`Ingrese el id del ticket a recuperar: `)
			fmt.Scan(&id)
			t, err := recuperarUnTicket(id)
			if err != nil {
				fmt.Printf("Error recuperando ticket: %v\n", err)
				break
			}
			fmt.Printf("Id: %v\n", t.idticket)
			fmt.Printf("Usuario: %v\n", t.usuario)
			fmt.Printf("Fecha de creacion: %v\n", t.fCreacion)
			fmt.Printf("Fecha de actualizacion: %v\n", t.fActualizacion)
			fmt.Printf("Estatus: %v\n", t.estatus)
			break
		}
		fmt.Print(menu_inicial)
		fmt.Scan(&eleccion)
	}

}
