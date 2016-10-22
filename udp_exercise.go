package main
 
import (
    "fmt"
    "net"
    "os"
    "crypto/sha256"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    
)


type MessageByte struct {
  gorm.Model
  Code string
  Price uint
}

type MessageChunk struct {
  gorm.Model
  Code string
  Price uint
}

type Message struct {
  gorm.Model
  Code string
  Price uint
}


/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}
 
func main() {
    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    /* Lets prepare a address at any address at port 10001*/   
    ServerAddr,err := net.ResolveUDPAddr("udp",":6789")
    CheckError(err)
 
    /* Now listen at selected port */
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
 
    buf := make([]byte, 1024)
 
    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)

        flag := uint16(buf[1]) | (uint16(buf[0]) << 8)
        size := uint16(buf[3]) | (uint16(buf[2]) << 8)
        offset := uint32(buf[7]) | (uint32(buf[6]) << 8) | (uint32(buf[5]) << 16) | (uint32(buf[5]) << 24)
        trans_id := uint32(buf[11]) | (uint32(buf[10]) << 8) | (uint32(buf[9]) << 16) | (uint32(buf[8]) << 24)
        fmt.Println("FLAG ", flag)
        fmt.Println("Size ", size)
        fmt.Println("OFFS ", offset)
        fmt.Println("TrID ", trans_id)
        fmt.Println("Received ", n, /*string(buf[0:n]),*/ " from ",addr)
        fmt.Println(sha256.Sum256(buf[12:size+12]))
 
        if err != nil {
            fmt.Println("Error: ",err)
        } 
    }
}