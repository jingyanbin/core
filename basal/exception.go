package basal

// CallerShort 调用信息短文件名
func CallerShort(skip int) (file string, line int)

// Caller 调用信息长文件名
func Caller(skip int) (file string, line int)

// CallerInFunc 调用信息包括函数名
func CallerInFunc(skip int) (name string, file string, line int)

// CallerLineStack 调用信息一行堆栈信息
func CallerLineStack(stack string) (name string, file string)

// Exception 捕获异常
func Exception(catchs ...func(stack string, e error))

// Try 捕获异常
func Try(f func(), catch func(stack string, e error))
