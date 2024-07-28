package common

import (
	"github.com/pkoukk/tiktoken-go"
	logger "gpt-content-audit/common/loggger"
)

var (
	Tke *tiktoken.Tiktoken
)

func init() {
	// gpt-4-turbo encoding
	tke, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		logger.FatalLog(err.Error())
	}
	Tke = tke

}

func CountTokens(text string) int {
	return len(Tke.Encode(text, nil, nil))
}
