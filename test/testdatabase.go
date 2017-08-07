package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
const (
	PAGEACCESSTOKEN = "EAAEo7uCJDOEBAIM0NNK8LwfJ25YHOPkKqZCiVsCowsjMUdLUB2l0ABXTZCkZBIM5rFzzTmPvqsdxKmOBMl2P4ZAwxa5qe2Fgk1w6XF34SMOmvbzllwaN9HUIsObcxxBkikkp4ApNo0ceHOIgvhE25B3DiqBZBipgQskeDkBvOsZBWQTFuN8h6y"
	VERIFYTOKEN     = "1234"
	PORT = 8080

	DB_NAME="recordchatbot"
	DB_USER="root"
	DB_PASS="1"
)
type Msg struct {
	fbid string
	userid int
}
var dbchan =make(chan *sql.DB,1)

func main()  {

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASS+"@/"+DB_NAME )//"user:password@/dbname"
	fmt.Println("Opening connection")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	fmt.Println("checked opening connnection")
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	fmt.Println("Ping database")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("checked ping database")
	dbchan<-db
	fmt.Println("added database channel")
	//1.test ham getlastsample
	//var fbid string
	//var lastsample int
	//var timestamp string
	//var msg = Msg{
	//	fbid:"123",
	//	userid:123,
	//}
	//query := "Select * from UserState where FbId=?"
	//rows, err := db.Query(query,msg.fbid)
	//if err != nil {
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//}
	//for rows.Next() {
	//	// Scan the value to []byte
	//	err = rows.Scan(&fbid, &lastsample,&timestamp)
	//
	//	if err != nil {
	//		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	//	}
	//}
	//fmt.Println(fbid,lastsample,timestamp)
	var msg = Msg{
	fbid:"123",
	userid:123,
	}
	var Id int
	var sample string
	stateofuser := GetStateOfUser(db,msg.fbid)
	Id= stateofuser+1
	fmt.Println("id",Id)
	fmt.Println("state of user",stateofuser)
	query := "SELECT * FROM InputText WHERE Id= ?"
	fmt.Println(query)
	rows, err := db.Query(query,Id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next(){
		rows.Scan(&Id,&sample)
	}
	fmt.Println(sample)


}
//func GetSampleOfUser(db *sql.DB, FbId string) string {
//	var samplestring string
//	var Id int
//	var sample int
//	stateofuser := GetStateOfUser(db,FbId)
//	sample = stateofuser + 1
//	query := "Select * from InputText where FbId="+string(sample)
//	rows, err := db.Query(query)
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//	}
//	for rows.Next(){
//		rows.Scan(&Id,&samplestring)
//	}
//	return samplestring
//
//}


func GetStateOfUser(db *sql.DB ,FbId string) int {
	var fbid string
	var lastsample int
	var timestamp string
	query := "Select * from UserState where FbId="+FbId
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next() {
		// Scan the value to []byte
		err = rows.Scan(&fbid, &lastsample,&timestamp)

		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
	}
	return lastsample
}