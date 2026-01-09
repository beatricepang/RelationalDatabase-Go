package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	// Capture connection properties.
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "recordings3"

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// ping error to terminal
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

}

// QUERY FOR A SINGLE DATA FIELD
func albumByArtist(name string) ([]Album, error) {
	// create an array of album object for the album data
	var albums []Album
	// collect the data from the artist names in the db
	// by writting it this way it prevents injection risks.
	rows, err := db.Query("SELECT * FROM album WHERE artist =?", name)

	// error handling for if the
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v, name, err")
	}
	// any resoruce it holds is released when the function exists
	defer rows.Close()

	// loop the returned rows through row.scan for the album struct fields
	// scans for a list of pointers to the go value
	// collumn values are written
	// inside the loop check for an error scanning the column values within the struct fields
	// append the new alb to the album slice
	for rows.Next() {

		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q; %v", name, err)

		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist &q: %v", name, err)

	}
	return albums, nil
}

// QUERY FOR A SINGLE ROW

func albumByID(id int64) (Album, error) {
	var alb Album
	row := db.QueryRow("SELECT * FROM album WHERE id=?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}
