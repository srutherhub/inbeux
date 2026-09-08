package agent

import (
	"embed"
	_ "embed"
	"fmt"
	platform "inbeux/internal/platform/agent"
)

//go:embed *.md
var agentInstructions embed.FS

//go:embed *.json
var agentSchemas embed.FS

type AgentService struct {
	MessageClassifier platform.Agent
}

func New() (*AgentService, error) {
	agent := &AgentService{}

	var AGENT_MESSAGE_CLASSIFIER string = "message_classifier"

	agent.MessageClassifier = *platform.NewAgent(AGENT_MESSAGE_CLASSIFIER)
	agent.MessageClassifier.SetJSONRequest(platform.LLMPaidGeminiFlashThreeSeven)

	err := agent.MessageClassifier.SetInstructions(agentInstructions)
	if err != nil {
		panic(err)
	}
	err = agent.MessageClassifier.SetSchema(agentSchemas)
	if err != nil {
		panic(err)
	}

	return agent, nil
}

func (as *AgentService) ClassifyMessage(content string) {
	output, err := as.MessageClassifier.Send(content)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(output)
	}

}
