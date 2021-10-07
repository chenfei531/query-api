package graphql

import (
	"encoding/json"

	"github.com/graphql-go/graphql"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/model"
)

var UserSchema graphql.Schema

func Init(dm data.DataManager) {
	agentType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Agent",
		Description: "agent",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "id of agent",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if agent, ok := p.Source.(model.Agent); ok {
						return agent.ID, nil
					}
					return nil, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"agents": &graphql.Field{
				Type:        graphql.NewList(agentType),
				Description: "agents",
				Args: graphql.FieldConfigArgument{
					"offset": &graphql.ArgumentConfig{
						Description: "offset",
						Type:        graphql.Int,
					},
					"limit": &graphql.ArgumentConfig{
						Description: "limit",
						Type:        graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					offset, ok := (p.Args["offset"]).(int)
					if ok == false {
						return nil, nil
					}
					limit, ok := (p.Args["limit"]).(int)
					if ok == false {
						return nil, nil
					}
					return dm.GetAgents(offset, limit), nil
				},
			},
		},
	})
	UserSchema, _ = graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
}

func Execute(query string) string {
	r := graphql.Do(graphql.Params{
		Schema:        UserSchema,
		RequestString: query,
	})
	rJSON, _ := json.MarshalIndent(r, "", "    ")
	return string(rJSON)
}
