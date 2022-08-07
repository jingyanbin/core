package basal

type KwArgs map[string]interface{}

func formatBytes(format []byte, signByte, startByte, endByte byte, kws KwArgs, ignore bool) []byte {
	sLen := len(format)
	buf := make([]byte, 0, sLen)
	for pos := 0; pos < sLen; pos++ {
		if format[pos] == signByte {
			posStartByte := pos + 1
			if posStartByte < sLen {
				if format[posStartByte] == startByte {
					var end int
					for end = posStartByte + 1; end < sLen; end++ {
						if format[end] == endByte {
							key := string(format[posStartByte+1 : end])
							if v, ok := kws[key]; ok {
								buf = append(buf, Sprintf("%v", v)...)
								pos = end
								break
							} else {
								panic(Sprintf("format bytes kws not found key: '%s', format: %s", key, string(format)))
							}
						}
					}
					if end >= sLen {
						if ignore {
							buf = append(buf, format[pos])
						} else {
							panic(Sprintf("format bytes not found end byte: '%s', format: %s", string(endByte), string(format)))
						}
					}
				} else {
					if ignore {
						buf = append(buf, format[pos])
					} else {
						panic(Sprintf("format bytes not found start byte: '%s', format: %s", string(startByte), string(format)))
					}
				}
			} else {
				if ignore {
					buf = append(buf, format[pos])
				} else {
					panic(Sprintf("format bytes not found start byte out of range: '%s', format: %s", string(startByte), string(format)))
				}
			}
		} else {
			buf = append(buf, format[pos])
		}
	}
	return buf
}

func Format(format string, sign, startByte, endByte byte, kws KwArgs, ignore bool) string {
	return string(formatBytes([]byte(format), sign, startByte, endByte, kws, ignore))
}

//格式器
type Formatter struct {
	signByte  byte
	startByte byte
	endByte   byte
	ignore    bool
}

func (m *Formatter) Format(format string, kws KwArgs) string {
	return string(formatBytes([]byte(format), m.signByte, m.startByte, m.endByte, kws, m.ignore))
}

//新建格式器
//signByte标志字符
//startByte 开始字符
//endByte 结束字符
//ignore: true 忽略非完整格式, false 不忽略出现非完整格式后panic
func NewFormatter(signByte, startByte, endByte byte, ignore bool) *Formatter {
	return &Formatter{signByte: signByte, startByte: startByte, endByte: endByte, ignore: ignore}
}
