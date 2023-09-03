package envexec

// Status 定义运行任务状态返回状态
type Status int

var statusToString = []string{
	"Invalid",
	"Accepted",
	"Wrong Answer",
	"Partially Correct",
	"Memory Limit Exceeded",
	"Time Limit Exceeded",
	"Output Limit Exceeded",
	"File Error",
	"Nonzero Exit Status",
	"Signalled",
	"Dangerous Syscall",
	"Judgement Failed",
	"Invalid Interaction",
	"Internal Error",
	"CGroup Error",
	"Container Error",
}

func (s Status) String() string {
	si := int(s)
	if si < 0 || si >= len(statusToString) {
		return statusToString[0]
	}
	return statusToString[si]
}
