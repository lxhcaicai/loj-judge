package envexec

import "fmt"

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

// stringToStatus map string to corresponding Status
var stringToStatus = make(map[string]Status)

const (
	// 未初始化状态
	StatusInvalid = iota

	// 正常退出
	StatusAccepted
	StatusWrongAnswer
	StatusPartiallyCorrect

	// 错误退出
	StatusMemoryLimitExceeded
	StatusTimeLimitExceeded
	StatusOutputLimitExceeded
	StatusFileError
	StatusNonzeroExitStatus
	StatusSignalled
	StatusDangerousSyscall

	//StatusRuntimeError // RE

	// SPJ / interactor error
	StatusJudgementFailed
	StatusInvalidInteraction

	// 内部错误包括:cgroup 初始化失败，容器失败等
	StatusInternalError
)

func (s Status) String() string {
	si := int(s)
	if si < 0 || si >= len(statusToString) {
		return statusToString[0]
	}
	return statusToString[si]
}

func StringToStatus(s string) (Status, error) {
	v, ok := stringToStatus[s]
	if !ok {
		return 0, fmt.Errorf("invalid string converting %s", s)
	}
	return v, nil
}

func init() {
	for i, v := range statusToString {
		stringToStatus["\""+v+"\""] = Status(i)
	}
}
