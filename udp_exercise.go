package main
 
import (
    "fmt"
    "net"
    "os"
    // "crypto/sha256"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "time"
    
)


type MessageByte struct {
  Data byte
  Id uint32
  MessageId uint32
  gorm.Model
}

type MessageChunk struct {
  Data []byte
  Id uint32
  Size uint16
  MessageId uint32
  Offset uint32
  gorm.Model
}


type Message struct {
  Id uint64
  Size uint32
  Bytes []byte
  ShaBytes []byte
  Sha string
  MessageChunks []MessageChunk `gorm:"ForeignKey:MessageId"`
  MessageBytes []MessageByte `gorm:"ForeignKey:MessageId"`
  gorm.Model
}


/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func CheckBytes(){
    for {
        for ix:=0; ix<10; ix++{
            allBytes := make([]byte, 0, 1024*1024)
            for _, mByte := range MessageBytes{
                if mByte.MessageId == uint32(ix) {
                    allBytes[mByte.Id] = mByte.Data
                }
            }
            fmt.Println(len(allBytes), cap(allBytes))
            fmt.Println(len(MessageBytes), cap(MessageBytes))
        }
        time.Sleep(time.Second)
    }
}

var MessageBytes = make([]MessageByte, 0, 1024*32)

func main() {
    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()
    // messages := []Message
    // var messageChunks []MessageChunk

    /* Lets prepare a address at any address at port 10001*/   
    ServerAddr,err := net.ResolveUDPAddr("udp",":6789")
    CheckError(err)
 
    /* Now listen at selected port */
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
 
    buf := make([]byte, 1024)

    go CheckBytes()
 
    for {
        // n,addr,err := ServerConn.ReadFromUDP(buf)
        _,_,err := ServerConn.ReadFromUDP(buf)

        flag := uint16(buf[1]) | (uint16(buf[0]) << 8)
        size := uint16(buf[3]) | (uint16(buf[2]) << 8)
        offset := uint32(buf[7]) | (uint32(buf[6]) << 8) | (uint32(buf[5]) << 16) | (uint32(buf[5]) << 24)
        trans_id := uint32(buf[11]) | (uint32(buf[10]) << 8) | (uint32(buf[9]) << 16) | (uint32(buf[8]) << 24)
        if flag>0 {
            fmt.Println(trans_id, "FLAG")
        }
        // fmt.Println("FLAG ", flag)
        // fmt.Println("Size ", size)
        // fmt.Println("OFFS ", offset)
        // fmt.Println("TrID ", trans_id)
        // fmt.Println("Received ", n, /*string(buf[0:n]),*/ " from ",addr)
        // fmt.Println(sha256.Sum256(buf[12:size+12]))
        for ix, b := range buf[12:size+12] {
            mByte := MessageByte{}
            mByte.Data = b
            mByte.Id = offset + uint32(ix)
            mByte.MessageId = trans_id
            append(MessageBytes, mByte)
        }
        // chunk := MessageChunk{}
        // chunk.Data := make([]byte, size)
        // copy(chunk.Data, buf[12:size+12])
        // chunk.Size = size
        // chunk.MessageId = trans_id
        // chunk.Offset = offset

        // fmt.Println(chunk)
        // append(messageChunks, chunk)
 
        if err != nil {
            fmt.Println("Error: ",err)
        } 
    }
}