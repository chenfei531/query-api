package graphql

import (
	"encoding/json"
	"github.com/graphql-go/graphql"

	"github.com/chenfei531/query-api/model"
	"github.com/chenfei531/query-api/query"
)

var UserSchema graphql.Schema

func Init(dm query.DataReader) {
	agentObject := graphql.NewObject(graphql.ObjectConfig{
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
			"createAt": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "create time of agent",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if agent, ok := p.Source.(model.Agent); ok {
						return agent.CreateAt, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "name of agent",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if agent, ok := p.Source.(model.Agent); ok {
						return agent.Name, nil
					}
					return nil, nil
				},
			},
		},
	})

	agentsField := &graphql.Field{
		Type:        graphql.NewList(agentObject),
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
	}

	userObject := graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "id of user",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(model.User); ok {
						return user.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "name of user",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(model.User); ok {
						return user.Name, nil
					}
					return nil, nil
				},
			},
			"agents": agentsField,
		},
	})

	userField := &graphql.Field{
		Type:        userObject,
		Description: "user",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Description: "id of user",
				Type:        graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := (p.Args["id"]).(int)
			if ok == false {
				return nil, nil
			}
			return dm.GetUserById(id), nil
		},
	}

	queryObject := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user":   userField,
			"agents": agentsField,
		},
	})
	//Schema->Object->Field
	//Field: describe struct of the Schema
	//Object: describe non primitive type fields
	//Interface: describe common attributes for objects implement it
	UserSchema, _ = graphql.NewSchema(graphql.SchemaConfig{Query: queryObject})
}

func Execute(query string) string {
	r := graphql.Do(graphql.Params{
		Schema:        UserSchema,
		RequestString: query,
	})
	rJSON, _ := json.MarshalIndent(r, "", "    ")
	return string(rJSON)
}
