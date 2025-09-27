package runtime

import (
	"context"
	"strings"

	"github.com/kiosk404/airi-go/backend/api/crossdomain/message"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
)

func (art *AgentRuntime) ChatflowRun(ctx context.Context, imagex imagex.ImageX) (err error) {

	// todo： 暂不支持 LangChain 中的对话流

	return err
}

func concatWfInput(rtDependence *AgentRuntime) string {
	var input string
	for _, content := range rtDependence.RunMeta.Content {
		if content.Type == message.InputTypeText {
			input = content.Text + "," + input
		} else {
			for _, file := range content.FileData {
				input += file.Url + ","
			}
		}
	}
	return strings.Trim(input, ",")
}
