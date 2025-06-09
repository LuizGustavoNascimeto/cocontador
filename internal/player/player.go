package player

import (
	db "cocontador/internal/database"
)

type Player struct {
	Id          uint32
	Name        string
	Total_score int32
}

func (p *Player) Insert() (id int32, err error) {
	conn, err := db.OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	sql := `INSERT INTO players (name, total_score) VALUES ($1, $2) RETURNING id`

	err = conn.QueryRow(sql, p.Name, p.Total_score).Scan(&id)

	return id, err
}
