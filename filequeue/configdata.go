package filequeue

import (
	internal "github.com/jingyanbin/core/internal"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var log = internal.GetStdoutLogger()

func SetLogger(logger internal.ILogger) {
	log = logger
}

const msgEOF byte = 27      //文件结束符
const byteMB1 = 1024 * 1024 //1MB

func MBToByteCount(mb int64) int64 {
	return mb * byteMB1
}

//配置数据基础结构
type configDataBase struct {
	filename string   //配置文件名
	f        *os.File //配置文件
	fsync    int32
}

func (m *configDataBase) SetFsync(fsync bool) {
	if fsync {
		atomic.StoreInt32(&m.fsync, 1)
	} else {
		atomic.StoreInt32(&m.fsync, 0)
	}
}

//同步
func (m *configDataBase) Sync() error {
	if m.f == nil {
		return nil
	}
	return m.f.Sync()
}

//关闭
func (m *configDataBase) Close() error {
	if m.f == nil {
		return nil
	}
	return m.f.Close()
}

//重新打开
func (m *configDataBase) reopen(force bool) error {
	if m.f == nil || force {
		if m.f != nil {
			m.f.Sync()
			m.f.Close()
		}
		f, err := internal.OpenFileB(m.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			m.f = nil
			return err
		}
		m.f = f
	}
	return nil
}

//加载配置参数
func (m *configDataBase) Load(numCount int) ([]int64, error) {
	data, err := ioutil.ReadFile(m.filename)
	if err != nil {
		if os.IsNotExist(err) {
			//log.ErrorF("不存在配置: %v", err)
			return nil, nil
		} else {
			return nil, err
		}
	}
	if nums := internal.NewNextNumber(string(data)).Numbers(); len(nums) >= numCount {
		result := make([]int64, 0, len(nums))
		for _, n := range nums {
			result = append(result, int64(n))
		}
		return result, nil
	} else {
		return nil, internal.NewError("load config file num len error: %v", string(data))
	}
}

//保存配置数据
func (m *configDataBase) Save(clear bool, nums ...int64) (err error) {
	if err = m.reopen(false); err != nil {
		return err
	}
	data := make([]string, 0, len(nums))
	for _, n := range nums {
		data = append(data, strconv.FormatInt(n, 10))
	}
	dataStr := strings.Join(data, ",")
	//log.ErrorF("data: %v, %v", dataStr, nums)
	if clear {
		if err = m.f.Truncate(0); err != nil {
			return err
		}
	}
	_, err = m.f.WriteAt([]byte(dataStr), 0)
	if err != nil {
		if err = m.reopen(true); err != nil {
			return err
		}
		if clear {
			if err = m.f.Truncate(0); err != nil {
				return err
			}
		}
		_, err = m.f.WriteAt([]byte(dataStr), 0)
	}
	if err == nil && atomic.LoadInt32(&m.fsync) == 1 {
		m.f.Sync()
	}
	return err
}

//pusher 配置
type configDataPusher struct {
	options *Options
	configDataBase
	//数据
	index int64 //push 消息文件的index
}

func (m *configDataPusher) Next() error {
	index := m.index + 1
	err := m.configDataBase.Save(false, index)
	if err == nil {
		m.index = index
	}
	return err
}

func (m *configDataPusher) Load() error {
	nums, err := m.configDataBase.Load(1)
	if err != nil {
		return err
	}
	if len(nums) >= 1 {
		m.index = nums[0]
		return nil
	} else {
		m.Save()
	}
	return nil
}

func (m *configDataPusher) Save() error {
	return m.configDataBase.Save(false, m.index)
}

//popper配置
type configDataPopper struct {
	options *Options
	configDataBase
	//数据
	index  int64 //pop 消息文件的index
	offset int64 //pop 消息文件内容偏移量
}

func (m *configDataPopper) Load() error {
	nums, err := m.configDataBase.Load(2)
	if err != nil {
		return err
	}
	if len(nums) >= 2 {
		m.index = nums[0]
		m.offset = nums[1]
		return nil
	} else {
		m.Save()
	}
	return nil
}

func (m *configDataPopper) Save() error {
	if err := m.configDataBase.Save(false, m.index, m.offset); err != nil {
		return err
	}
	return nil
}

func (m *configDataPopper) SaveEx(index, offset int64) error {
	if err := m.configDataBase.Save(true, index, offset); err != nil {
		return err
	}
	m.index = index
	m.offset = offset
	return nil
}