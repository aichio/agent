package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	net "net"
	"os"
	"strconv"
	"strings"
)

var fileSize int64 = 20480000

func GetFileMD5(path string) (string, error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	// log.Println(fileinfo.Size())
	if fileinfo.Size() >= fileSize  {
		return "", errors.New("big file")
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	md5Ctx := md5.New()
	if _, err = io.Copy(md5Ctx, file); err != nil {
		return "", err
	}
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}

func BigEndianPut(v uint32) uint32 {
	var testBytes []byte = make([]byte, 4)
	binary.BigEndian.PutUint32(testBytes, uint32(v))
	convInt := binary.LittleEndian.Uint32(testBytes)
	return convInt
}

func LittleEndianPut(v uint32) uint32 {
	var testBytes []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(testBytes, uint32(v))
	convInt := binary.BigEndian.Uint32(testBytes)
	return convInt
}

type IP uint32
// 根据IP和mask换算内网IP范围
func Table(ipNet *net.IPNet) []IP {
	ip := ipNet.IP.To4()
	fmt.Printf("本机ip: %s", ip)
	var min, max IP
	var data []IP
	for i := 0; i < 4; i++ {
		b := IP(ip[i] & ipNet.Mask[i])
		min += b << ((3 - uint(i)) * 8)
	}
	one, _ := ipNet.Mask.Size()
	max = min | IP(math.Pow(2, float64(32 - one)) - 1)
	fmt.Println("内网IP范围:",min," --- ", max)
	// max 是广播地址，忽略
	// i & 0x000000ff  == 0 是尾段为0的IP，根据RFC的规定，忽略
	for i := min; i < max; i++ {
		if i & 0x000000ff == 0 {
			continue
		}
		data = append(data, i)
	}
	return data
}


func(ip IP)BigEndianPut() IP{
	var testBytes []byte = make([]byte, 4)
	binary.BigEndian.PutUint32(testBytes, uint32(ip))
	convInt := binary.LittleEndian.Uint32(testBytes)
	return IP(convInt)
}

func(ip IP)LittleEndianPut() IP{
	var testBytes []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(testBytes, uint32(ip))
	convInt := binary.BigEndian.Uint32(testBytes)
	return IP(convInt)
}

// []byte --> IP
func ParseIP(b []byte) IP {
	return IP(IP(b[0]) << 24 + IP(b[1]) << 16 + IP(b[2]) << 8 + IP(b[3]))
}

// string --> IP
func ParseIPString(s string) IP{
	var b []byte
	for _, i := range strings.Split(s, ".") {
		v, _ := strconv.Atoi(i)
		b = append(b, uint8(v))
	}
	return ParseIP(b)
}

func (ip IP) String() string {
	var bf bytes.Buffer
	for i := 1; i <= 4; i++ {
		bf.WriteString(strconv.Itoa(int((ip >> ((4 - uint(i)) * 8)) & 0xff)))
		if i != 4 {
			bf.WriteByte('.')
		}
	}
	return bf.String()
}