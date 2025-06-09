package player

import (
	db "cocontador/internal/database"
	"time"
)

type Score struct {
	id        uint32
	id_player uint32
	date      time.Time
}

func (s *Score) Insert() (id int32, err error) {
	conn, err := db.OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	sql := `INSERT INTO players (player_id, createdAt) VALUES ($1, $2) RETURNING id`

	err = conn.QueryRow(sql, s.id_player, s.date).Scan(&id)

	return id, err
}
