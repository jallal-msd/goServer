package main

import (
    "net/smtp"
	"database/sql"
	"fmt"
	"log"
	"net/http"
    "sql3go/hashing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type comments struct {

    User string
    Comment string
    Rating int8 `json"rating, string"`
}

func storeComment(comment comments) (bool, error){
    tx, err := DB.Begin()

    stmnt, err := tx.Prepare("insert into comments values(?,?,?)")

    checkErr(err)

    defer stmnt.Close()

    _, err = stmnt.Exec(&comment.User, &comment.Comment, &comment.Rating)

    checkErr(err)

    tx.Commit()

    return true,nil
}

func postComments(c *gin.Context){

    var comment comments

    err := c.ShouldBindJSON(&comment)
    checkErr(err)

    succes, err := storeComment(comment)
    checkErr(err)

    if succes {
        
        c.JSON(http.StatusOK, gin.H{"message": "success"})
    }else {
        c.JSON(http.StatusBadRequest, gin.H{"error":"usuck"})
    }
}


func savingComments() ([]comments, error){
    
    rows, err := DB.Query("select * from comments")
    checkErr(err)
    
    commentt := make([]comments, 0)
    for rows.Next(){
        comment := comments{}
        err = rows.Scan(&comment.User, &comment.Comment,
                &comment.Rating)

        if err != nil {
            return nil,err

        commentt = append(commentt, comment)
    }

    err = rows.Err()

    if err != nil {
        return nil, err
    }
}
  return commentt, nil

}

func getComments(c *gin.Context){

    rows, err := savingComments()
    checkErr(err)
    
    if rows == nil {
        
        c.JSON(http.StatusBadRequest, gin.H{"error": "no rows"})
    }else {
        
        c.JSON(http.StatusOK, gin.H{"data":rows})
    }

}

func sendMail() {
    
    from := "gerard52@ethereal.email"
    passowrd := "sPMU8YSF2wPJ8GKJCn"
    
    to := []string {"jallalblack@gmail.com" ,}


    stmpHost :="smtp.ethereal.email"
    stmpPort :="587"

    message := []byte("this is a test message")

    auth := smtp.PlainAuth("", from, passowrd, stmpHost)

    err := smtp.SendMail(stmpHost+":"+stmpPort, auth, from, to, message)

    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Email sent")
}




var DB *sql.DB

type infos struct {
    Name string
    Lname string
}

func addInfo(newInfo infos) (bool,error) {
    
    tx,err :=DB.Begin()

    stmt, err := tx.Prepare("insert into info values(?,?)")
    if err != nil {
        return false, err
    }

    defer stmt.Close()

    _, err = stmt.Exec(&newInfo.Name, &newInfo.Lname)

    if err != nil {
        return false, err
    }
    tx.Commit()

    return true, err
}

func sendData(c *gin.Context) {

    var json infos
    
    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
    }

    success, err := addInfo(json);

    if err != nil {
        log.Fatal(err)
    }
    
    if success {
        
        c.JSON(http.StatusOK, gin.H{"message": "Success"})

    }else {
        c.JSON(http.StatusBadRequest,gin.H{"message": "ur mom is a hoe"} )
    }

}

func getInfo() ([]infos, error){
    
    rows, err := DB.Query("select * from info")

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    info := make([]infos, 0)

    for rows.Next() {
        singleInfo := infos{}
        err = rows.Scan(&singleInfo.Name, &singleInfo.Lname)

        if err != nil {
            return nil, err
        }
        
        info = append(info, singleInfo)
    }
    err = rows.Err()

    if err != nil {
        return nil, err
    }
    return info, nil
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func getInfos(c *gin.Context){
    
    info, err := getInfo()
    if err != nil {
        checkErr(err)
    }

    if info == nil {
        c.JSON(404, gin.H{"error": "No records found"})
        return
    }else {
        c.JSON(200, gin.H{"message": info})
    }


    // stmt, err := DB.Prepare("insert into info values(?,?)")
    //
    // if err != nil {
    //     log.Fatal(err)
    // }
    // stmt.Exec("samuel", "eto")

    

}


func connect() {
    
    db, err := sql.Open("sqlite3", "./infos.db")
    if err != nil {
        log.Fatal(err)
    }
    DB = db
    fmt.Println("connected")
    
}
func main() {
    sendMail()
    hashing.Check()
    connect()

   router :=gin.Default()
   router.Use(cors.Default())
   

   router.GET("/albums", getComments)
   router.POST("/albums", postComments) 
   router.Run("localhost:8080")

}

