package naming_client

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type PushReceiver struct {
	port        int
	host        string
	hostReactor *HostReactor
}

type PushData struct {
	PushType    string `json:"type"`
	Data        string `json:"data"`
	LastRefTime int64  `json:"lastRefTime"`
}

func NewPushRecevier(hostReactor *HostReactor) *PushReceiver {
	pr := PushReceiver{
		hostReactor: hostReactor,
	}
	go pr.startServer()
	return &pr
}

func (us *PushReceiver) tryListen() (*net.UDPConn, bool) {
	addr, err := net.ResolveUDPAddr("udp", us.host+":"+strconv.Itoa(us.port))
	if err != nil {
		log.Printf("[ERROR]: Can't resolve address,err: %s \n", err.Error())
		return nil, false
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("Error listening %s:%d,err:%s \n", us.host, us.port, err.Error())
		return nil, false
	}

	return conn, true
}

func (us *PushReceiver) startServer() {
	var conn *net.UDPConn

	for i := 0; i < 3; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		port := r.Intn(1000) + 54951
		us.port = port
		conn1, ok := us.tryListen()

		if ok {
			conn = conn1
			log.Println("[INFO] udp server start, port: " + strconv.Itoa(port))
			break
		}

		if !ok && i == 2 {
			log.Panicf("failed to start udp server after trying 3 times.")
			os.Exit(1)
		}
	}

	defer conn.Close()
	for {
		us.handleClient(conn)
	}
}

func (us *PushReceiver) handleClient(conn *net.UDPConn) {
	data := make([]byte, 4024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Printf("[ERROR]:failed to read UDP msg because of %s \n", err.Error())
		return
	}

	s := utils.TryDecompressData(data[:n])
	log.Println("[INFO] receive push: "+s+" from: ", remoteAddr)

	var pushData PushData
	err1 := json.Unmarshal([]byte(s), &pushData)
	if err1 != nil {
		log.Printf("[ERROR] failed to process push data.err:%s \n", err1.Error())
		return
	}
	ack := make(map[string]string)

	if pushData.PushType == "dom" || pushData.PushType == "service" {
		us.hostReactor.ProcessServiceJson(pushData.Data)

		ack["type"] = "push-ack"
		ack["lastRefTime"] = strconv.FormatInt(pushData.LastRefTime, 10)
		ack["data"] = ""

	} else if pushData.PushType == "dump" {
		ack["type"] = "dump-ack"
		ack["lastRefTime"] = strconv.FormatInt(pushData.LastRefTime, 10)
		ack["data"] = utils.ToJsonString(us.hostReactor.serviceInfoMap)
	} else {
		ack["type"] = "unknow-ack"
		ack["lastRefTime"] = strconv.FormatInt(pushData.LastRefTime, 10)
		ack["data"] = ""
	}

	bs, _ := json.Marshal(ack)
	conn.WriteToUDP(bs, remoteAddr)
}
