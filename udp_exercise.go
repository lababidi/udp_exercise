package main

import (
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type messageByte struct {
	Data    byte
	ID      uint32
	transID uint32
	gorm.Model
}

type messageChunk struct {
	Data    []byte
	ID      uint32
	Size    uint16
	transID uint32
	Offset  uint32
	gorm.Model
}

type message struct {
	ID            uint64
	Size          uint32
	Bytes         []byte
	ShaBytes      []byte
	Sha           string
	messageChunks []messageChunk `gorm:"ForeignKey:transID"`
	messageBytes  []messageByte  `gorm:"ForeignKey:transID"`
	gorm.Model
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

var messageBytes = make([][]messageByte, 11)

func findMax(bytes []messageByte) uint32 {
	maxID := uint32(0)
	for _, mByte := range bytes {
		if maxID < mByte.ID {
			maxID = mByte.ID
		}
	}
	return maxID
}

func checkBytes() {
	for {
		allBytes := make([][]byte, 11)
		for ix := 1; ix < 11; ix++ {
			allBytes[ix] = make([]byte, findMax(messageBytes[ix])+1)
			for _, mByte := range messageBytes[ix] {
				allBytes[ix][mByte.ID] = mByte.Data
			}
			h := sha256.New()
			h.Write(allBytes[ix])
			hout := fmt.Sprintf("%x", h.Sum(nil))
			fmt.Println(ix, len(allBytes[ix]), hout)
		}
		time.Sleep(time.Second * 10)
	}
}

func produceByte(buf []byte) []messageByte {
	bufBytes := make([]messageByte, 0, 1024)
	flag := uint16(buf[1]) | (uint16(buf[0]) << 8)
	size := uint16(buf[3]) | (uint16(buf[2]) << 8)
	offset := uint32(buf[7]) | (uint32(buf[6]) << 8) | (uint32(buf[5]) << 16) | (uint32(buf[5]) << 24)
	transID := uint32(buf[11]) | (uint32(buf[10]) << 8) | (uint32(buf[9]) << 16) | (uint32(buf[8]) << 24)
	if flag > 0 {
		fmt.Println(transID, flag, "FLAG")
	}

	// for ix, b := range buf[12 : size+12] {
	for ix := uint16(0); ix < size; ix++ {
		mByte := messageByte{}
		mByte.Data = buf[12+ix]
		mByte.ID = offset + uint32(ix)
		mByte.transID = transID
		bufBytes = append(bufBytes, mByte)
	}
	return bufBytes
}

func main() {

	for ix := range messageBytes {
		messageBytes[ix] = make([]messageByte, 0, 1024*1024)
	}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	ServerAddr, err := net.ResolveUDPAddr("udp", ":6789")
	checkError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	checkError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	go checkBytes()

	for {
		_, _, err := ServerConn.ReadFromUDP(buf)
		checkError(err)
		for _, mb := range produceByte(buf) {
			messageBytes[mb.transID] = append(messageBytes[mb.transID], mb)
		}
	}
	// chunk := messageChunk{}
	// chunk.Data := make([]byte, size)
	// copy(chunk.Data, buf[12:size+12])
	// chunk.Size = size
	// chunk.transID = transID
	// chunk.Offset = offset

	// append(messageChunks, chunk)

}
