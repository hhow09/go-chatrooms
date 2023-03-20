package repository

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/hhow09/go-chatrooms/chatroom/model"
	"github.com/hhow09/go-chatrooms/chatroom/util"
)

func (repo *RoomRepository) notUsed() bool {
	return repo.Db == nil
}

type RoomRepository struct {
	Db *sql.DB
}

func (repo *RoomRepository) AddRoom(room model.Room) {
	if repo.notUsed() {
		return
	}
	stmt, err := repo.Db.Prepare("INSERT OR IGNORE INTO room(name, private) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(room.GetName(), room.GetPrivate())
	checkErr(err)
}

func (repo *RoomRepository) FindRoomByName(name string, redisClient *redis.Client) model.Room {
	if repo.notUsed() {
		return nil
	}
	row := repo.Db.QueryRow("SELECT name, private FROM room where name = ? LIMIT 1", name)

	name, private := "", false
	if err := row.Scan(&name, &private); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return model.NewRoom(name, private, redisClient)

}

func (repo *RoomRepository) GetAllRooms() ([]string, error) {
	if repo.notUsed() {
		return nil, nil
	}
	rows, err := repo.Db.Query("SELECT name FROM room")
	if err != nil {
		util.Log("GetAllRooms error", err.Error())
		return nil, err
	}
	defer rows.Close()
	res := []string{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			util.Log("GetAllRooms error", err.Error())
			return nil, err
		}
		res = append(res, name)
	}

	return res, nil
}

func checkErr(err error) {
	if err != nil {
		util.Log(err.Error())
	}
}
