package data

type User struct {
	ID         int
	Name       string
	Lastname   string
	Student_id int //legajo
	Email      string
	Phone      int
	Role       string
	Dni        int
	//Creator_id int   --> SE AGREGA DESPUES, MUCHO LABURO PARA SOLO ESTAR PROBANDO COSAS
	School string
	Hash   []byte
	Salt   string
}
