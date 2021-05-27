package probe

import (
	"bytes"
	"cube/log"
	"cube/model"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	UNIQUE_NAMES = map[string]string{
		"\x00": "Workstation Service",
		"\x03": "Messenger Service",
		"\x06": "RAS Server Service",
		"\x1F": "NetDDE Service",
		"\x20": "Server Service",
		"\x21": "RAS Client Service",
		"\xBE": "Network Monitor Agent",
		"\xBF": "Network Monitor Application",
		"\x1D": "Master Browser",
		"\x1B": "Domain Master Browser",
	}

	GROUP_NAMES = map[string]string{
		"\x00": "Domain Name",
		"\x1C": "Domain Controllers",
		"\x1E": "Browser Service Elections",
	}

	NetBIOS_ITEM_TYPE = map[string]string{
		"\x01\x00": "NetBIOS computer name",
		"\x02\x00": "NetBIOS domain name",
		"\x03\x00": "DNS computer name",
		"\x04\x00": "DNS domain name",
		"\x05\x00": "DNS tree name",
		"\x07\x00": "Time stamp",
	}
)

var SmbV2D1 = []byte{
	0x00, 0x00, 0x00, 0x45, 0xFF, 0x53, 0x4D, 0x42, 0x72, 0x00,
	0x00, 0x00, 0x00, 0x18, 0x01, 0x48, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF,
	0xAC, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22, 0x00, 0x02,
	0x4E, 0x54, 0x20, 0x4C, 0x4D, 0x20, 0x30, 0x2E, 0x31, 0x32,
	0x00, 0x02, 0x53, 0x4D, 0x42, 0x20, 0x32, 0x2E, 0x30, 0x30,
	0x32, 0x00, 0x02, 0x53, 0x4D, 0x42, 0x20, 0x32, 0x2E, 0x3F,
	0x3F, 0x3F, 0x00,
}
var SmbV2D2 = []byte{
	0x00, 0x00, 0x00, 0x68, 0xFE, 0x53, 0x4D, 0x42, 0x40, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00,
	0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x02, 0x02, 0x10, 0x02,
}
var NtlmV2 = []byte{
	0x00, 0x00, 0x00, 0x9A, 0xFE, 0x53, 0x4D, 0x42, 0x40, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19, 0x00,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x58, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x60, 0x40, 0x06, 0x06, 0x2B, 0x06, 0x01, 0x05,
	0x05, 0x02, 0xA0, 0x36, 0x30, 0x34, 0xA0, 0x0E, 0x30, 0x0C,
	0x06, 0x0A, 0x2B, 0x06, 0x01, 0x04, 0x01, 0x82, 0x37, 0x02,
	0x02, 0x0A, 0xA2, 0x22, 0x04, 0x20, 0x4E, 0x54, 0x4C, 0x4D,
	0x53, 0x53, 0x50, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x05, 0x80,
	0x08, 0xa0,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var SmbBufferV2 = []byte{
	0x00, 0x00, 0x00, 0xfc, 0xff, 0x53, 0x4d, 0x42,
	0x73, 0x00, 0x00, 0x00, 0x00, 0x18, 0x07, 0xc8,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xfe,
	0x00, 0x00, 0x40, 0x00, 0x0c, 0xff, 0x00, 0x0a,
	0x01, 0x04, 0x41, 0x32, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x4c, 0x00, 0x00, 0x00, 0x00,
	0x00, 0xd4, 0x00, 0x00, 0xa0, 0xc1, 0x00, 0x60,
	0x48, 0x06, 0x06, 0x2b, 0x06, 0x01, 0x05, 0x05,
	0x02, 0xa0, 0x3e, 0x30, 0x3c, 0xa0, 0x0e, 0x30,
	0x0c, 0x06, 0x0a, 0x2b, 0x06, 0x01, 0x04, 0x01,
	0x82, 0x37, 0x02, 0x02, 0x0a, 0xa2, 0x2a, 0x04,
	0x28, 0x4e, 0x54, 0x4c, 0x4d, 0x53, 0x53, 0x50,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x07, 0x82, 0x08,
	0xa2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x05, 0x02, 0xce, 0x0e, 0x00, 0x00, 0x00,
	0x0f, 0x00, 0x00, 0x00, 0x57, 0x00, 0x69, 0x00,
	0x6e, 0x00, 0x64, 0x00, 0x6f, 0x00, 0x77, 0x00,
	0x73, 0x00, 0x20, 0x00, 0x53, 0x00, 0x65, 0x00,
	0x72, 0x00, 0x76, 0x00, 0x65, 0x00, 0x72, 0x00,
	0x20, 0x00, 0x32, 0x00, 0x30, 0x00, 0x30, 0x00,
	0x33, 0x00, 0x20, 0x00, 0x33, 0x00, 0x37, 0x00,
	0x39, 0x00, 0x30, 0x00, 0x20, 0x00, 0x53, 0x00,
	0x65, 0x00, 0x72, 0x00, 0x76, 0x00, 0x69, 0x00,
	0x63, 0x00, 0x65, 0x00, 0x20, 0x00, 0x50, 0x00,
	0x61, 0x00, 0x63, 0x00, 0x6b, 0x00, 0x20, 0x00,
	0x32, 0x00, 0x00, 0x00, 0x57, 0x00, 0x69, 0x00,
	0x6e, 0x00, 0x64, 0x00, 0x6f, 0x00, 0x77, 0x00,
	0x73, 0x00, 0x20, 0x00, 0x32, 0x00, 0x30, 0x00,
	0x30, 0x00, 0x33, 0x00, 0x20, 0x00, 0x35, 0x00,
	0x2e, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00,
}
var SmbBufferV1 = []byte{
	0x00, 0x00, 0x00, 0x85, 0xff, 0x53, 0x4d, 0x42,
	0x72, 0x00, 0x00, 0x00, 0x00, 0x18, 0x53, 0xc8,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xfe,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x62, 0x00, 0x02,
	0x50, 0x43, 0x20, 0x4e, 0x45, 0x54, 0x57, 0x4f,
	0x52, 0x4b, 0x20, 0x50, 0x52, 0x4f, 0x47, 0x52,
	0x41, 0x4d, 0x20, 0x31, 0x2e, 0x30, 0x00, 0x02,
	0x4c, 0x41, 0x4e, 0x4d, 0x41, 0x4e, 0x31, 0x2e,
	0x30, 0x00, 0x02, 0x57, 0x69, 0x6e, 0x64, 0x6f,
	0x77, 0x73, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x57,
	0x6f, 0x72, 0x6b, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x73, 0x20, 0x33, 0x2e, 0x31, 0x61, 0x00, 0x02,
	0x4c, 0x4d, 0x31, 0x2e, 0x32, 0x58, 0x30, 0x30,
	0x32, 0x00, 0x02, 0x4c, 0x41, 0x4e, 0x4d, 0x41,
	0x4e, 0x31, 0x2e, 0x31, 0x00, 0x02, 0x4e, 0x54,
	0x20, 0x4c, 0x4d, 0x20, 0x30, 0x2e, 0x31, 0x32,
	0x00,
}

var d1 = []byte{
	0x00, 0x00, 0x00, 0x85, 0xFF, 0x53, 0x4D, 0x42, 0x72, 0x00, 0x00, 0x00, 0x00, 0x18, 0x53, 0xC8,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFE,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x62, 0x00, 0x02, 0x50, 0x43, 0x20, 0x4E, 0x45, 0x54, 0x57, 0x4F,
	0x52, 0x4B, 0x20, 0x50, 0x52, 0x4F, 0x47, 0x52, 0x41, 0x4D, 0x20, 0x31, 0x2E, 0x30, 0x00, 0x02,
	0x4C, 0x41, 0x4E, 0x4D, 0x41, 0x4E, 0x31, 0x2E, 0x30, 0x00, 0x02, 0x57, 0x69, 0x6E, 0x64, 0x6F,
	0x77, 0x73, 0x20, 0x66, 0x6F, 0x72, 0x20, 0x57, 0x6F, 0x72, 0x6B, 0x67, 0x72, 0x6F, 0x75, 0x70,
	0x73, 0x20, 0x33, 0x2E, 0x31, 0x61, 0x00, 0x02, 0x4C, 0x4D, 0x31, 0x2E, 0x32, 0x58, 0x30, 0x30,
	0x32, 0x00, 0x02, 0x4C, 0x41, 0x4E, 0x4D, 0x41, 0x4E, 0x32, 0x2E, 0x31, 0x00, 0x02, 0x4E, 0x54,
	0x20, 0x4C, 0x4D, 0x20, 0x30, 0x2E, 0x31, 0x32, 0x00,
}
var d2 = []byte{
	0x00, 0x00, 0x01, 0x0A, 0xFF, 0x53, 0x4D, 0x42, 0x73, 0x00, 0x00, 0x00, 0x00, 0x18, 0x07, 0xC8,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFE,
	0x00, 0x00, 0x40, 0x00, 0x0C, 0xFF, 0x00, 0x0A, 0x01, 0x04, 0x41, 0x32, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x4A, 0x00, 0x00, 0x00, 0x00, 0x00, 0xD4, 0x00, 0x00, 0xA0, 0xCF, 0x00, 0x60,
	0x48, 0x06, 0x06, 0x2B, 0x06, 0x01, 0x05, 0x05, 0x02, 0xA0, 0x3E, 0x30, 0x3C, 0xA0, 0x0E, 0x30,
	0x0C, 0x06, 0x0A, 0x2B, 0x06, 0x01, 0x04, 0x01, 0x82, 0x37, 0x02, 0x02, 0x0A, 0xA2, 0x2A, 0x04,
	0x28, 0x4E, 0x54, 0x4C, 0x4D, 0x53, 0x53, 0x50, 0x00, 0x01, 0x00, 0x00, 0x00, 0x07, 0x82, 0x08,
	0xA2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x05, 0x02, 0xCE, 0x0E, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x57, 0x00, 0x69, 0x00, 0x6E, 0x00,
	0x64, 0x00, 0x6F, 0x00, 0x77, 0x00, 0x73, 0x00, 0x20, 0x00, 0x53, 0x00, 0x65, 0x00, 0x72, 0x00,
	0x76, 0x00, 0x65, 0x00, 0x72, 0x00, 0x20, 0x00, 0x32, 0x00, 0x30, 0x00, 0x30, 0x00, 0x33, 0x00,
	0x20, 0x00, 0x33, 0x00, 0x37, 0x00, 0x39, 0x00, 0x30, 0x00, 0x20, 0x00, 0x53, 0x00, 0x65, 0x00,
	0x72, 0x00, 0x76, 0x00, 0x69, 0x00, 0x63, 0x00, 0x65, 0x00, 0x20, 0x00, 0x50, 0x00, 0x61, 0x00,
	0x63, 0x00, 0x6B, 0x00, 0x20, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x57, 0x00, 0x69, 0x00,
	0x6E, 0x00, 0x64, 0x00, 0x6F, 0x00, 0x77, 0x00, 0x73, 0x00, 0x20, 0x00, 0x53, 0x00, 0x65, 0x00,
	0x72, 0x00, 0x76, 0x00, 0x65, 0x00, 0x72, 0x00, 0x20, 0x00, 0x32, 0x00, 0x30, 0x00, 0x30, 0x00,
	0x33, 0x00, 0x20, 0x00, 0x35, 0x00, 0x2E, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00,
}

type NbnsName struct {
	unique    string
	group     string
	msg       string
	osversion string
}

func smbScan2(task model.ProbeTask) (result model.ProbeTaskResult) {
	result = model.ProbeTaskResult{ProbeTask: task, Result: "", Err: nil}
	realhost := fmt.Sprintf("%s:%v", task.Ip, task.Port)
	conn, err := net.DialTimeout("tcp", realhost, 3*time.Second)
	if err != nil {
		log.Debug(err)
		return
	}
	_, err = conn.Write(d1)
	if err != nil {
		return
	}
	buf := make([]byte, 4096)
	conn.Read(buf)

	_, err = conn.Write(d2)
	//_, err = conn.Write(SmbV2D2)
	if err != nil {
		return
	}
	//buf := make([]byte, 4096)
	//conn.Read(buf)
	//fmt.Println(buf)

	ret, err := readbytes(conn)
	start1 := bytes.Index(ret, []byte("NTLMSSP"))
	fmt.Println(start1)
	if err != nil || len(ret) < 45 {
		return
	}

	//_, err = conn.Write(NtlmV2)

	num1, err := bytetoint(ret[43:44][0])
	if err != nil {
		return
	}
	num2, err := bytetoint(ret[44:45][0])
	if err != nil {
		return
	}
	fmt.Printf("ret: %x\n", ret)
	length := num1 + num2*256
	osVersion := ret[47+length:]
	osVersion = bytes.ReplaceAll(osVersion, []byte{0x00}, []byte{})
	fmt.Printf("Version: %x\n", osVersion)
	result.Result = string(osVersion[:])

	R := bytes.ReplaceAll(ret, []byte{0x00}, []byte{})
	rs := []rune(string(R)) // 将字符串转为字节rune切片
	fmt.Println(rs)         // 输出rune切片
	fmt.Println(string(rs)) // 将rune切片转为字符串

	return result
}

func SmbProbe(task model.ProbeTask) (result model.ProbeTaskResult) {
	result = model.ProbeTaskResult{ProbeTask: task, Result: "", Err: nil}
	nbname, _ := NetBIOS1(task)

	var msg, isdc string
	//result.Result = nbname

	if strings.Contains(nbname.msg, "Domain Controllers") {
		isdc = "[+]DC"
	}
	msg += fmt.Sprintf("[*] %-15s%-5s %s\\%s   %s", task.Ip, isdc, nbname.group, nbname.unique, nbname.osversion)
	fmt.Printf("Group: %s\nMsg: %s\nOsVersion: %s\nUniqe: %s", nbname.group, nbname.msg, nbname.osversion, nbname.unique)
	//if info.Scantype == "netbios" {
	//	msg += "\n-------------------------------------------\n" + nbname.msg
	//}
	//if len(nbname.group) > 0 || len(nbname.unique) > 0 {
	//	common.LogSuccess(msg)
	//}
	return result
}

func NetBIOS1(task model.ProbeTask) (nbname NbnsName, err error) {
	nbname, err = GetNbnsname(task)
	var payload0 []byte
	if err == nil {
		name := netbiosEncode(nbname.unique)
		payload0 = append(payload0, []byte("\x81\x00\x00D ")...)
		payload0 = append(payload0, name...)
		payload0 = append(payload0, []byte("\x00 EOENEBFACACACACACACACACACACACACA\x00")...)
	}
	realhost := fmt.Sprintf("%s:%v", task.Ip, task.Port)
	conn, err := net.DialTimeout("tcp", realhost, 3*time.Second)
	if err != nil {
		log.Debug(err)
		return
	}
	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return
	}
	defer conn.Close()

	if task.Port == 139 && len(payload0) > 0 {
		_, err1 := conn.Write(payload0)
		if err1 != nil {
			return
		}
		_, err1 = readbytes(conn)
		if err1 != nil {
			return
		}
	}

	//payload1 := []byte("\x00\x00\x00\x85\xff\x53\x4d\x42\x72\x00\x00\x00\x00\x18\x53\xc8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xfe\x00\x00\x00\x00\x00\x62\x00\x02\x50\x43\x20\x4e\x45\x54\x57\x4f\x52\x4b\x20\x50\x52\x4f\x47\x52\x41\x4d\x20\x31\x2e\x30\x00\x02\x4c\x41\x4e\x4d\x41\x4e\x31\x2e\x30\x00\x02\x57\x69\x6e\x64\x6f\x77\x73\x20\x66\x6f\x72\x20\x57\x6f\x72\x6b\x67\x72\x6f\x75\x70\x73\x20\x33\x2e\x31\x61\x00\x02\x4c\x4d\x31\x2e\x32\x58\x30\x30\x32\x00\x02\x4c\x41\x4e\x4d\x41\x4e\x32\x2e\x31\x00\x02\x4e\x54\x20\x4c\x4d\x20\x30\x2e\x31\x32\x00")
	//payload2 := []byte("\x00\x00\x01\x0a\xff\x53\x4d\x42\x73\x00\x00\x00\x00\x18\x07\xc8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xfe\x00\x00\x40\x00\x0c\xff\x00\x0a\x01\x04\x41\x32\x00\x00\x00\x00\x00\x00\x00\x4a\x00\x00\x00\x00\x00\xd4\x00\x00\xa0\xcf\x00\x60\x48\x06\x06\x2b\x06\x01\x05\x05\x02\xa0\x3e\x30\x3c\xa0\x0e\x30\x0c\x06\x0a\x2b\x06\x01\x04\x01\x82\x37\x02\x02\x0a\xa2\x2a\x04\x28\x4e\x54\x4c\x4d\x53\x53\x50\x00\x01\x00\x00\x00\x07\x82\x08\xa2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x02\xce\x0e\x00\x00\x00\x0f\x00\x57\x00\x69\x00\x6e\x00\x64\x00\x6f\x00\x77\x00\x73\x00\x20\x00\x53\x00\x65\x00\x72\x00\x76\x00\x65\x00\x72\x00\x20\x00\x32\x00\x30\x00\x30\x00\x33\x00\x20\x00\x33\x00\x37\x00\x39\x00\x30\x00\x20\x00\x53\x00\x65\x00\x72\x00\x76\x00\x69\x00\x63\x00\x65\x00\x20\x00\x50\x00\x61\x00\x63\x00\x6b\x00\x20\x00\x32\x00\x00\x00\x00\x00\x57\x00\x69\x00\x6e\x00\x64\x00\x6f\x00\x77\x00\x73\x00\x20\x00\x53\x00\x65\x00\x72\x00\x76\x00\x65\x00\x72\x00\x20\x00\x32\x00\x30\x00\x30\x00\x33\x00\x20\x00\x35\x00\x2e\x00\x32\x00\x00\x00\x00\x00")
	//_, err = conn.Write(payload1)
	//if err != nil {
	//	return
	//}
	//_, err = readbytes(conn)
	//if err != nil {
	//	return
	//}
	_, err = conn.Write(d1)
	if err != nil {
		return
	}
	buf := make([]byte, 4096)
	conn.Read(buf)
	//fmt.Println(buf)
	_, err = conn.Write(d2)
	//_, err = conn.Write(SmbV2D2)
	if err != nil {
		return
	}
	//buf := make([]byte, 4096)
	//conn.Read(buf)
	//fmt.Println(buf)

	ret, err := readbytes(conn)
	start1 := bytes.Index(ret, []byte("NTLMSSP"))
	fmt.Println(start1)
	if err != nil || len(ret) < 45 {
		return
	}

	//_, err = conn.Write(NtlmV2)

	num1, err := bytetoint(ret[43:44][0])
	if err != nil {
		return
	}
	num2, err := bytetoint(ret[44:45][0])
	if err != nil {
		return
	}
	length := num1 + num2*256
	os_version := ret[47+length:]
	fmt.Printf("Version: %s", os_version)
	tmp1 := bytes.ReplaceAll(os_version, []byte{0x00, 0x00}, []byte{124})
	tmp1 = bytes.ReplaceAll(tmp1, []byte{0x00}, []byte{})
	msg1 := string(tmp1[:len(tmp1)-1])

	nbname.osversion = msg1

	index1 := strings.Index(msg1, "|")
	if index1 > 0 {
		nbname.osversion = nbname.osversion[:index1]
	}
	nbname.msg += "-------------------------------------------\n"
	nbname.msg += msg1 + "\n"
	start := bytes.Index(ret, []byte("NTLMSSP"))
	num1, err = bytetoint(ret[start+40 : start+41][0])
	if err != nil {
		return
	}
	num2, err = bytetoint(ret[start+41 : start+42][0])
	if err != nil {
		return
	}
	length = num1 + num2*256
	num1, err = bytetoint(ret[start+44 : start+45][0])
	if err != nil {
		return
	}
	offset, err := bytetoint(ret[start+44 : start+45][0])
	if err != nil {
		return
	}
	index := start + offset
	for index < start+offset+length {
		item_type := ret[index : index+2]
		num1, err = bytetoint(ret[index+2 : index+3][0])
		if err != nil {
			return
		}
		num2, err = bytetoint(ret[index+3 : index+4][0])
		if err != nil {
			return
		}
		item_length := num1 + num2*256
		item_content := bytes.ReplaceAll(ret[index+4:index+4+item_length], []byte{0x00}, []byte{})
		index += 4 + item_length
		if string(item_type) == "\x07\x00" {
			//Time stamp, 暂时不想处理
		} else if NetBIOS_ITEM_TYPE[string(item_type)] != "" {
			nbname.msg += fmt.Sprintf("%-22s: %s\n", NetBIOS_ITEM_TYPE[string(item_type)], string(item_content))
		} else if string(item_type) == "\x00\x00" {
			break
		} else {
			nbname.msg += fmt.Sprintf("Unknown: %s\n", string(item_content))
		}
	}
	return nbname, err
}

func wmi(task model.ProbeTask) {
	data1 := []byte{5, 0, 11, 3, 16, 0, 0, 0, 120, 0, 40, 0, 3, 0, 0, 0, 184, 16, 184, 16, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 160, 1, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 70, 0, 0, 0, 0, 4, 93, 136, 138, 235, 28, 201, 17, 159, 232, 8, 0, 43, 16, 72, 96, 2, 0, 0, 0, 10, 2, 0, 0, 0, 0, 0, 0, 78, 84, 76, 77, 83, 83, 80, 0, 1, 0, 0, 0, 7, 130, 8, 162, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 1, 177, 29, 0, 0, 0, 15}
	realhost := fmt.Sprintf("%s:%v", task.Ip, 135)
	conn, _ := net.DialTimeout("tcp", realhost, 3*time.Second)
	conn.Write(data1)
	ret, _ := readbytes(conn)
	start := bytes.Index(ret, []byte("NTLMSSP"))
	fmt.Println(start)
}

func GetNbnsname(task model.ProbeTask) (nbname NbnsName, err error) {
	senddata1 := []byte{102, 102, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 32, 67, 75, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 0, 0, 33, 0, 1}
	realhost := fmt.Sprintf("%s:%v", task.Ip, 137)
	conn, err := net.DialTimeout("udp", realhost, 3*time.Second)
	if err != nil {
		return
	}
	err = conn.SetDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Write(senddata1)
	if err != nil {
		return
	}
	text, err := readbytes(conn)
	if err != nil {
		return
	}
	if len(text) < 57 {
		return nbname, fmt.Errorf("no names available")
	}
	num, err := bytetoint(text[56:57][0])
	if err != nil {
		return
	}
	data := text[57:]
	var msg string
	for i := 0; i < num; i++ {
		name := string(data[18*i : 18*i+15])
		flag_bit := data[18*i+15 : 18*i+16]
		if GROUP_NAMES[string(flag_bit)] != "" && string(flag_bit) != "\x00" {
			msg += fmt.Sprintf("%s G %s\n", name, GROUP_NAMES[string(flag_bit)])
		} else if UNIQUE_NAMES[string(flag_bit)] != "" && string(flag_bit) != "\x00" {
			msg += fmt.Sprintf("%s U %s\n", name, UNIQUE_NAMES[string(flag_bit)])
		} else if string(flag_bit) == "\x00" {
			name_flags := data[18*i+16 : 18*i+18][0]
			if name_flags >= 128 {
				nbname.group = strings.Replace(name, " ", "", -1)
				msg += fmt.Sprintf("%s G %s\n", name, GROUP_NAMES[string(flag_bit)])
			} else {
				nbname.unique = strings.Replace(name, " ", "", -1)
				msg += fmt.Sprintf("%s U %s\n", name, UNIQUE_NAMES[string(flag_bit)])
			}
		} else {
			msg += fmt.Sprintf("%s \n", name)
		}
	}
	nbname.msg += msg
	return
}

func readbytes(conn net.Conn) (result []byte, err error) {
	buf := make([]byte, 4096)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			break
		}
		result = append(result, buf[0:count]...)
		if count < 4096 {
			break
		}
	}
	return result, err
}

func bytetoint(text byte) (int, error) {
	num1 := fmt.Sprintf("%v", text)
	num, err := strconv.Atoi(num1)
	return num, err
}

func netbiosEncode(name string) (output []byte) {
	var names []int
	src := fmt.Sprintf("%-16s", name)
	for _, a := range src {
		char_ord := int(a)
		high_4_bits := char_ord >> 4
		low_4_bits := char_ord & 0x0f
		names = append(names, high_4_bits, low_4_bits)
	}
	for _, one := range names {
		out := (one + 0x41)
		output = append(output, byte(out))
	}
	return
}
