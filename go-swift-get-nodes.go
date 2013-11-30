package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GetNodesConfig struct {
	RingFile       string
	GetNodesTarget string
	HashPrefix     string
	HashSuffix     string
	PartShift      uint
}

type Device struct {
	Id              uint    `json:"id"`
	Device          string  `json:"device"`
	Ip              string  `json:"ip"`
	Meta            string  `json:"meta"`
	Port            uint    `json:"port"`
	Region          uint    `json:"region"`
	ReplicationIp   string  `json:"replication_ip"`
	ReplicationPort uint    `json:"replication_port"`
	Weight          float64 `json:"weight"`
	Zone            uint    `json:"zone"`
}

type Ring struct {
	Devs                []Device `json:"devs"`
	ReplicaCount        uint     `json:"replica_count"`
	PartShift           uint     `json:"part_shift"`
	replica2part2dev_id [][]uint16
}

func (r Ring) GetPartitionNodes(partition uint) []Device {
	var response []Device
	for i := uint(0); i < r.ReplicaCount; i++ {
		response = append(response, r.Devs[r.replica2part2dev_id[i][partition]])
	}
	return response
}

func LoadRing(path string) Ring {
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	gz, err := gzip.NewReader(fp)
	if err != nil {
		panic(err)
	}
	magic_buf := make([]byte, 4)
	io.ReadFull(gz, magic_buf)
	// TODO: assert magic_buf == "R1NG"
	var ring_version uint16
	binary.Read(gz, binary.BigEndian, &ring_version)
	// TODO: assert ring_version == 1
	var json_len uint32
	binary.Read(gz, binary.BigEndian, &json_len)
	json_buf := make([]byte, json_len)
	io.ReadFull(gz, json_buf)
	var ring Ring
	json.Unmarshal(json_buf, &ring)
	partition_count := 1 << (32 - ring.PartShift)
	for i := uint(0); i < ring.ReplicaCount; i++ {
		part2dev := make([]uint16, partition_count)
		binary.Read(gz, binary.LittleEndian, &part2dev)
		ring.replica2part2dev_id = append(ring.replica2part2dev_id, part2dev)
	}
	return ring
}

func GetPart(path string, conf *GetNodesConfig) (uint, []byte) {
	hash := md5.New()
	target := conf.HashPrefix + path + conf.HashSuffix
	io.WriteString(hash, target)
	buf := bytes.NewBuffer([]byte(hash.Sum(nil)))
	var part uint32
	err := binary.Read(buf, binary.BigEndian, &part)
	if err != nil {
		panic(err)
	}
	return uint(part >> conf.PartShift), hash.Sum(nil)
}

func readConf(configFile string, gnc *GetNodesConfig) {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "swift_hash_path_prefix =") {
			gnc.HashPrefix = strings.TrimSpace(strings.TrimPrefix(line, "swift_hash_path_prefix ="))
			continue
		}
		if strings.HasPrefix(line, "swift_hash_path_suffix =") {
			gnc.HashSuffix = strings.TrimSpace(strings.TrimPrefix(line, "swift_hash_path_suffix ="))
			continue
		}
	}
}

func doHead(id int, url string) error {
	res, err := http.Head(url)
	if err != nil {
		fmt.Println(err)
		return errors.New("Host error")
	} else {
		res.Body.Close()
		fmt.Printf("\nNode %d: %s\n", id+1, res.Status)
		fmt.Println("Url:", url)
		for k, _ := range res.Header {
			fmt.Printf(" %v: %v\n", k, res.Header[k])
		}
		if res.StatusCode/200 != 1 {
			return errors.New("Non 2xx status")
		}
	}
	return nil
}

func main() {
	conf := GetNodesConfig{"/etc/swift/object.ring.gz", "", "", "", 14}
	if len(os.Args) != 3 {
		fmt.Println("Usage: ")
		fmt.Println("\t", os.Args[0], "/path/to/ring.gz /account/[container]/[object]")
		fmt.Println()
		os.Exit(1)
	}
	readConf("/etc/swift/swift.conf", &conf)
	conf.RingFile = os.Args[1]
	conf.GetNodesTarget = os.Args[2]
	ring := LoadRing(conf.RingFile)
	conf.PartShift = ring.PartShift

	fmt.Printf("Ring File:\t%s\n", conf.RingFile)
	fmt.Printf("Target:\t\t%s\n", conf.GetNodesTarget)

	part, hash := GetPart(os.Args[2], &conf)
	fmt.Printf("Partion:\t%v\n", part)
	fmt.Printf("Hash:\t\t%x\n", hash)

	ecount := 0
	for i := 0; i <= 2; i++ {
		err := doHead(i, fmt.Sprintf("http://%v:%v/%v/%v%v", ring.GetPartitionNodes(part)[i].Ip, ring.GetPartitionNodes(part)[i].Port, ring.GetPartitionNodes(part)[i].Device, part, conf.GetNodesTarget))
		if err != nil {
			ecount = 1
		}
	}
	os.Exit(ecount)
}
