package chiweb

import (
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

func ParseBody(r *http.Request, outptr any) error {
	// 读取Body，使用sonic进行反序列化
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("http read body. %w", err)
	}
	if err := sonic.Unmarshal(data, outptr); err != nil {
		return fmt.Errorf("http unmarshal body. %w", err)
	}
	return nil
}
