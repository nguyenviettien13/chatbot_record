package main

import (
	"github.com/michlabs/fbbot"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
)

const (
	PAGEACCESSTOKEN = "EAAEo7uCJDOEBAIM0NNK8LwfJ25YHOPkKqZCiVsCowsjMUdLUB2l0ABXTZCkZBIM5rFzzTmPvqsdxKmOBMl2P4ZAwxa5qe2Fgk1w6XF34SMOmvbzllwaN9HUIsObcxxBkikkp4ApNo0ceHOIgvhE25B3DiqBZBipgQskeDkBvOsZBWQTFuN8h6y"
	VERIFYTOKEN     = "1234"
	PORT = 8080

	DB_NAME="recordchatbot"
	DB_USER="root"
	DB_PASS="1"
)

type Record struct {

}

var dbchan =make(chan *sql.DB,1)

func main()  {
	//processing database
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

	var record Record
	bot := fbbot.New(PORT,VERIFYTOKEN,PAGEACCESSTOKEN)
	bot.AddMessageHandler(record)
	bot.Run()
}

func ( record Record ) HandleMessage( bot *fbbot.Bot , msg *fbbot.Message ) {
	//check whether is newUser
	db := <- dbchan
	//if arrived messeger is text and user is new
	if IsNewUser( msg , db ) && !isAudioMessage( msg ){
		greeting := "xin chao!"
		m := fbbot.NewTextMessage(greeting)
		bot.Send(msg.Sender,m)

		tutorialmesseger :="huong dan"
		m1 := fbbot.NewTextMessage(tutorialmesseger)
		bot.Send(msg.Sender,m1)

		stmtInsNewUser, err := db.Prepare("INSERT INTO UserState (Fbid,LastSample)VALUES( ?, ? )") // ? = placeholder
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtInsNewUser.Close() // Close the statement when we leave main() / the program terminates
		_, err = stmtInsNewUser.Exec(string(msg.Sender.ID),0)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		sample1 := GetSampleOfUser(db,msg.Sender.ID)
		m2 := fbbot.NewTextMessage(sample1)
		bot.Send(msg.Sender,m2)

		//add new user and state 0

	}else if isAudioMessage(msg) {
		//luu database
		stmtInsAudio , err := db.Prepare("INSERT INTO Outputs (FbId,Gender,SampleId,State,NumberTime, UrlRecord, RecordTime)VALUES( ?, ? , ? , ?, ?, ?, ? )")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		SampleId := 1 + GetStateOfUser(db,msg.Sender.ID)
		stmtInsAudio.Query(msg.Sender.ID,msg.Sender.Gender(),SampleId,false,1,msg.Audios[0].URL,time.Now())
		//send confirm messager
		confirmmessager := "ban co muon gui audio : y/n"
		m := fbbot.NewTextMessage(confirmmessager)
		bot.Send(msg.Sender,m)
	} else if msg.Text =="y" || msg.Text =="Y"  {
		//cap nhat output
		fmt.Println("confirm Y")
		stmtInsAudio , err := db.Prepare("UPDATE Outputs  SET State = ?  WHERE FbId = ? AND SampleId = ?")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		SampleId := 1 + GetStateOfUser(db,msg.Sender.ID)
		stmtInsAudio.Query(true,msg.Sender.ID,SampleId)

		//cap nhat trang thai cho user
		stmtUpdateUserState , err := db.Prepare("UPDATE UserState SET LastSample = ? WHERE FbId = ? ")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		stmtUpdateUserState.Query(SampleId,msg.Sender.ID)
		//??? co nen dong statement o day khong?
		stmtUpdateUserState.Close()

		nextsample := GetSampleOfUser( db , msg.Sender.ID )
		m1 :=fbbot.NewTextMessage(nextsample)
		bot.Send(msg.Sender,m1)

	} else if msg.Text == "n" || msg.Text == "N"  {
		fmt.Println("confirm N")
		//gui lai text yeu cau doc lai
		sample := GetSampleOfUser(db,msg.Sender.ID)
		m := fbbot.NewTextMessage(sample)
		bot.Send(msg.Sender,m)
	} else {
		next := "cau tiep theo 2: "
		m := fbbot.NewTextMessage(next)
		bot.Send(msg.Sender,m)

		nextsample := GetSampleOfUser( db , msg.Sender.ID )
		m1 :=fbbot.NewTextMessage(nextsample)
		bot.Send(msg.Sender,m1)
	}
	dbchan <- db
}

//in my opinion, a newUser is who nerver been sent Audio Messager
func IsNewUser(msg *fbbot.Message , db *sql.DB ) bool {
	query := "SELECT * FROM UserState WHERE FbId="+msg.Sender.ID
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	if rows.Next() {
		return false
	}
	return true
}

func isAudioMessage(msg *fbbot.Message) bool  {
	if len(msg.Audios) == 0 {
		return false
	}
	return true
}

func GetSampleOfUser(db *sql.DB, FbId string) string {
	var sample string
	var id int
	stateofuser := GetStateOfUser(db,FbId)
	id = stateofuser + 1
	query := "Select * from InputText where Id= ?"
	rows, err := db.Query(query,id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next(){
		rows.Scan(&id,&sample)
	}
	return sample

}

func GetStateOfUser(db *sql.DB ,FbId string) int {
	var fbid string
	var lastsample int
	var lasttime string
	fmt.Println("running function GetStateOfUser ....")
	query := "Select * from UserState where FbId = ? "
	rows, err := db.Query(query,FbId)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next() {
		// Scan the value to []byte
		err = rows.Scan(&fbid, &lastsample,&lasttime)

		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic

		}
	}
	return lastsample
}

//tai sao khi cho khoi tao bot o duoi phan database thi no lai khong chay?
//tai sao cau lenh dbchannel<- db lai block ham main