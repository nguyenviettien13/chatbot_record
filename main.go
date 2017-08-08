package main

import (
	"github.com/michlabs/fbbot"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
	"log"
)

const (
	PAGEACCESSTOKEN = "EAAXu4QafGd8BANrhhNs9KOuhdBMYpv8wdi0GyCaBt62fcvtZBYvgLdfsgVlk4m2p9P5h3M2Cxf80IU41P3YRUao0f80D88i8q8SoWaD6RxZA7C2fQlFOF6ZBn3ZAvZBw2vtcuBc2XhkUM5WZCrQW4zQtKa1cPoF4r6yPbTzVCqCMCcK1jEk0dJ"
	VERIFYTOKEN     = "neitteiv1234"
	PORT = 2102

	DB_NAME="record_chatbot"
	DB_USER="root"
	DB_PASS="1"
	MaxSample=3
)

var db *sql.DB

type Record struct {}


func (r Record) HandlePostback(bot *fbbot.Bot, pbk *fbbot.Postback)  {
	switch pbk.Payload {
	case "Yes":
		stmtInsAudio , err := db.Prepare("UPDATE Outputs  SET State = ?  WHERE FbId = ? AND SampleId = ?")
		if err != nil {
			log.Println("error when create stmtInsAudio") // proper error handling instead of panic in your app
		}

		SampleId := 1 + GetStateOfUser(db,pbk.Sender.ID)
		stmtInsAudio.Query(true,pbk.Sender.ID,SampleId)

		//cap nhat trang thai cho user
		stmtUpdateUserState , err := db.Prepare("UPDATE UserState SET LastSample = ? WHERE FbId = ? ")
		if err != nil {
			log.Println("error when create stmtUpdateUserState")
		}
		stmtUpdateUserState.Query(SampleId,pbk.Sender.ID)
		//??? co nen dong statement o day khong?
		stmtUpdateUserState.Close()

		if SampleId == MaxSample {
			goobye := "Bạn đã hoàn thành, xin chào bạn"
			m := fbbot.NewTextMessage(goobye)
			bot.Send(pbk.Sender,m)

		} else {
			nextsample := GetSampleOfUser( db , pbk.Sender.ID )

			m1 :=fbbot.NewTextMessage(nextsample)
			bot.Send(pbk.Sender,m1)
			m4:=fbbot.NewTextMessage("....")
			bot.Send(pbk.Sender,m4)
			bot.Send(pbk.Sender,m4)
			bot.Send(pbk.Sender,m4)
		}
	case "No":
		sample := GetSampleOfUser(db,pbk.Sender.ID)
		m := fbbot.NewTextMessage(sample)
		bot.Send(pbk.Sender,m)
		m4:=fbbot.NewTextMessage("....")
		bot.Send(pbk.Sender,m4)
		bot.Send(pbk.Sender,m4)
		bot.Send(pbk.Sender,m4)

	default:
		log.Println("Switch case does not exist")
	}
}

func main()  {
	//processing database
	var err error
	db, err = sql.Open("mysql", DB_USER+":"+DB_PASS+"@/"+DB_NAME )//"user:password@/dbname"
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
	fmt.Println("added database channel")


	var r Record
	bot := fbbot.New(PORT,VERIFYTOKEN,PAGEACCESSTOKEN)
	bot.AddMessageHandler(r)
	bot.AddPostbackHandler(r)

	bot.Run()
}

func ( record Record ) HandleMessage( bot *fbbot.Bot , msg *fbbot.Message ) {
	//check whether is newUser
	//if arrived messeger is text and user is new
	if IsNewUser( msg , db ) && !isAudioMessage( msg ){

		greeting := "Xin chào!"+msg.Sender.FirstName()+"Chúng tôi đang thực hiện một dự án thu thập dữ liệu ghi âm giọng nói và rất vui khi nhận được sự hợp tác của bạn"
		m := fbbot.NewTextMessage(greeting)
		bot.Send(msg.Sender,m)

		tutorialmesseger :="Hướng dẫn: Bây h tôi sẽ gửi cho bạn một đoạn text, bạn hãy đọc và ghi âm rồi gửi chúng lại cho t ôi"
		m1 := fbbot.NewTextMessage(tutorialmesseger)
		bot.Send(msg.Sender,m1)

		start:="Oki! Bây h chúng ta sẽ bắt đầu"
		m2 := fbbot.NewTextMessage(start)
		bot.Send(msg.Sender,m2)

		stmtInsNewUser, err := db.Prepare("INSERT INTO UserState ( Fbid,LastSample )VALUES( ?, ? )") // ? = placeholder
		if err != nil {
			log.Println("error when create stminsertNewUser")
		}
		_, err = stmtInsNewUser.Exec(string(msg.Sender.ID),0)
		if err != nil {
			log.Println("error when exec stminsertNewUser")
		}

		sample1 := GetSampleOfUser(db,msg.Sender.ID)
		m3 := fbbot.NewTextMessage(sample1)
		bot.Send(msg.Sender,m3)

		m4:=fbbot.NewTextMessage("....")
		bot.Send(msg.Sender,m4)
		bot.Send(msg.Sender,m4)
		bot.Send(msg.Sender,m4)



	}else if isAudioMessage(msg) {
		//luu database
		us := GetStateOfUser(db,msg.Sender.ID)
		if ! isExistOutput(msg.Sender.ID,us+1) {
			stmtInsAudio , err := db.Prepare("INSERT INTO Outputs (FbId,Gender,SampleId,State, UrlRecord, RecordTime)VALUES( ?, ? , ? , ?, ?, ? )")
			if err != nil {
				log.Println("error when create stminsertAudio")
			}
			_ , err = stmtInsAudio.Query(msg.Sender.ID,msg.Sender.Gender(),us+1,false,msg.Audios[0].URL,time.Now())
			if err !=nil {
				log.Println("error when exec stminsertOutput")
			}
		} else {
			stmtInsAudio , err := db.Prepare("UPDATE Outputs SET UrlRecord=? WHERE FbId= ? AND SampleId=?")
			if err != nil {
				log.Println("error when create stminsertAudio")
			}
			stmtInsAudio.Query(msg.Audios[0].URL, msg.Sender.ID,us+1)
			_ , err = stmtInsAudio.Query(msg.Sender.ID,msg.Sender.Gender(),us+1,false,msg.Audios[0].URL,time.Now())
			if err !=nil {
				log.Println("error when exec stminsertOutput")
			}
		}
		//send confirm messager
		//confirmmessager := "ban co muon gui audio : y/n"
		//m := fbbot.NewTextMessage(confirmmessager)
		//bot.Send(msg.Sender,m)

		pkb :=fbbot.NewButtonMessage()
		pkb.Text = "Bạn muốn ghi âm lại hay ghi âm câu tiếp theo"
		pkb.Noti ="REGULAR"
		pkb.AddPostbackButton("Ghi âm lại","No")
		pkb.AddPostbackButton("Câu tiếp theo","Yes")
		bot.Send(msg.Sender,pkb)
	}  else {
		us := GetStateOfUser(db,msg.Sender.ID)
		sample1 := GetSampleOfUser(db,msg.Sender.ID)
		if us == MaxSample {
			m3 := fbbot.NewTextMessage("Bạn đã hoàn thành quá trình ghi âm!")
			bot.Send(msg.Sender,m3)

			m4 := fbbot.NewTextMessage("Cảm ơn bạn")
			bot.Send(msg.Sender,m4)
		} else {
			m3 := fbbot.NewTextMessage(sample1)
			bot.Send(msg.Sender,m3)
			m4:=fbbot.NewTextMessage("....")
			bot.Send(msg.Sender,m4)
			bot.Send(msg.Sender,m4)
			bot.Send(msg.Sender,m4)
		}
	}
}

//in my opinion, a newUser is who nerver been sent Audio Messager
func IsNewUser(msg *fbbot.Message , db *sql.DB ) bool {
	query := "SELECT * FROM UserState WHERE FbId="+msg.Sender.ID
	rows, err := db.Query(query)
	if err != nil {
		log.Println("error when SELECT * FROM UserState WHERE FbId")
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
	if id <= MaxSample {
		query := "Select * from InputText where Id= ?"
		rows, err := db.Query(query,id)
		if err != nil {
			log.Println("err when exec func GetSampleOfUser")
		}
		for rows.Next(){
			rows.Scan(&id,&sample)
		}
		return sample
	} else {
		return ""
	}


}

func GetStateOfUser(db *sql.DB ,FbId string) int {
	var fbid string
	var lastsample int
	var lasttime string
	fmt.Println("running function GetStateOfUser ....")
	query := "Select * from UserState where FbId = ? "
	rows, err := db.Query(query,FbId)
	if err != nil {
		log.Println("err when exec func GetStateOfUser")
	}
	for rows.Next() {
		// Scan the value to []byte
		err = rows.Scan(&fbid, &lastsample,&lasttime)

		if err != nil {
			log.Println("err when scan to &lastsample")

		}
	}
	return lastsample
}
func isExistOutput(FbId string, SampleId int) bool {
	rows, err := db.Query("SELECT * FROM Outputs WHERE FbId=? AND  SampleId=?",FbId,SampleId)

	if err != nil {
		log.Println("err when exec func isExistOutput")
	}
	return  rows.Next()
}
//tai sao khi cho khoi tao bot o duoi phan database thi no lai khong chay?
//tai sao cau lenh dbchannel<- db lai block ham main
//gettime