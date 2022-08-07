package basal

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const WINDOWS = `windows`
const LINUX = `linux`

const ERR_PID_NOT_FOUND = -1 //没有找到进程
const ERR_PID_ITOA = -2      //查找错误
const ERR_PID_OS = -3        //系统错误
const ERR_PID_NETWORK = -4   //网络协议错误

type systemCommand struct {
}

//系统命令
var SystemCmder = &systemCommand{}

//基本命令
func (self *systemCommand) Command(cmd string) (output []byte, err error) {
	if runtime.GOOS == WINDOWS {
		return exec.Command("cmd", "/c", cmd).CombinedOutput()
	} else if runtime.GOOS == LINUX {
		return exec.Command("sh", "-c", cmd).CombinedOutput()
	} else {
		return nil, NewError("os error: %v", runtime.GOOS)
	}
}

//根据协议地址端口, 返回使用该端口的进程ID(返回小于0为错误,大于等于0为端口)
func (self *systemCommand) GetPidByAddr(network string, ip string, port int) int {
	var pid = ERR_PID_NOT_FOUND //没有找到
	nwk := strings.ToLower(network)
	if nwk == "http" {
		nwk = "tcp"
	}
	if nwk == "tcp" || nwk == "udp" || nwk == "" {
		addr := fmt.Sprintf("%s:%d", ip, port)
		if runtime.GOOS == WINDOWS {
			var cmdStr string
			if nwk == "" {
				cmdStr = fmt.Sprintf("netstat -ano %s | findstr %s", "", addr)
			} else {
				cmdStr = fmt.Sprintf("netstat -ano %s | findstr %s", "-p "+nwk, addr)
			}
			output, _ := self.Command(cmdStr)
			reg := regexp.MustCompile(`[ ]+(\w+)[ ]+(\S+)[ ]+(\S+)[ ]+(\S*)[ ]+(\d+)`)
			results := reg.FindAllSubmatch(output, -1)
			for _, res := range results {
				//第0个元素是字符串,后面的元素才是分割的数据
				localAddr := string(res[2])
				state := string(res[4])
				pidStr := string(res[5])
				if strings.HasSuffix(localAddr, addr) && state == "LISTENING" {
					id, err := strconv.Atoi(strings.TrimSpace(pidStr))
					if err != nil {
						pid = ERR_PID_ITOA //转换错误
					} else {
						pid = id
					}
					break
				}
				//records := make([]string, 0, len(res)-1)
				//for _, r := range res[1:] {
				//	records = append(records, string(r))
				//}
				//LogInfo("%v, %v, %v", i, records, len(records))
			}
		} else if runtime.GOOS == LINUX {
			var cmdStr string
			if nwk == "" {
				cmdStr = fmt.Sprintf("netstat -tupln | grep '%s'", addr)
			} else {
				cmdStr = fmt.Sprintf("netstat -%cpln | grep '%s'", nwk[0], addr)
			}
			output, _ := self.Command(cmdStr)
			reg := regexp.MustCompile(`(\w+)[ ]+(\d+)[ ]+(\d+)[ ]+(\S+)[ ]+(\S+)[ ]+([A-Z]+)[ ]+([0-9]+)/(\S*)[ ]*\n`)
			results := reg.FindAllSubmatch(output, -1)
			for _, res := range results {
				localAddr := string(res[4])
				state := string(res[6])
				pidStr := string(res[7])
				if strings.HasSuffix(localAddr, addr) && state == "LISTEN" {
					id, err := strconv.Atoi(strings.TrimSpace(pidStr))
					if err != nil {
						pid = ERR_PID_ITOA //转换错误
					} else {
						pid = id
					}
					break
				}
			}
		} else {
			pid = ERR_PID_OS //平台错误
		}
	} else {
		pid = ERR_PID_NETWORK
	}
	return pid
}

//端口是否被占用
//network: “http”,“tcp“,“udp“
func (c *systemCommand) IsUsedPortByAddr(network string, ip string, port int) bool {
	pid := c.GetPidByAddr(network, ip, port)
	if pid < 0 {
		return false
	}
	return true
}

//端口是否被占用(port)
func (c *systemCommand) IsUsedPort(port int) bool {
	return c.IsUsedPortByAddr("", "", port)
}
