package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Repository struct {
	DB *sql.DB
}

var repo Repository

func Init() {
	if db, err := sql.Open("sqlite3", "tracks.db"); err == nil {
		repo = Repository{DB: db}
	} else {
		log.Fatal("Database initialisation")
	}
}

func Create() int {
	const sql = "CREATE TABLE IF NOT EXISTS Tracks" +
		    "(Id TEXT PRIMARY KEY, Audio TEXT)"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Clear() int {
	const sql = "DELETE FROM Tracks"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Update(c Track) int64 {
	const sql = "UPDATE Tracks SET Audio = ? WHERE id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(c.Audio, c.Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n 
			}
		}
	}
	return -1
}

func Insert(c Track) int64 {
	const sql = "INSERT INTO Tracks(Id, Audio) VALUES (?, ?)"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(c.Id, c.Audio); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n 
			}
		}
	}
	return -1
}

func Read(id string) (Track, int64) {
	const sql = "SELECT * FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		var c Track
		row := stmt.QueryRow(id)
		if err := row.Scan(&c.Id, &c.Audio); err == nil {
			return c, 1
		} else {
			return Track{}, 0
		}
	}
	return Track{}, -1
}

func List() ([]string, int64) {
    const sql = "SELECT Id FROM Tracks"
    rows, err := repo.DB.Query(sql)
    if err != nil {
        return nil, -1
    }
    defer rows.Close()

    var ids []string
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, -1
        }
        ids = append(ids, id)
    }
    if err := rows.Err(); err != nil {
        return nil, -1
    }

    return ids, int64(len(ids))
}

func Delete(id string) int64 {
	const sql = "DELETE FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n 
			}
		}
	}
	return -1
}