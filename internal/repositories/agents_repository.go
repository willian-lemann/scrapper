package repositories

import (
	"scrapper/internal/structs"

	supa "github.com/nedpals/supabase-go"
)

type AgentsRepository struct {
	client *supa.Client
}

type IAgentsRepository interface{}

func NewAgentsRepository(client *supa.Client) *AgentsRepository {
	return &AgentsRepository{
		client: client,
	}
}

func (repo AgentsRepository) SaveOne(agentItem structs.Agent) (bool, error) {
	defer repo.client.DB.CloseIdleConnections()

	var results []structs.Agent
	err := repo.client.DB.From("agent_list").Upsert(agentItem).Execute(&results)
	if err != nil {
		return false, err
	}
	return true, nil
}
