package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
    "sql3go/hashing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

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
    hashing.Check()
    connect()

   router :=gin.Default()
   router.Use(cors.Default())
   

   router.GET("/albums", getInfos)
   router.POST("/albums", sendData) 
   router.Run("localhost:8080")

}
