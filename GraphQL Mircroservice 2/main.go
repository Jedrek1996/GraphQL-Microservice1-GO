package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Tutorial struct {
	Id       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name     string
	Tutorial []int
}

type Comment struct {
	Body string
}

// Populate tutorial slice
func populate() []Tutorial {
	author := &Author{Name: "Handsome guy", Tutorial: []int{1}}
	tutorial := Tutorial{
		Id:       1,
		Title:    "GraphQL Tutorial",
		Author:   *author,
		Comments: []Comment{Comment{Body: "First Comment"}},
	}

	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)

	return tutorials
}

func main() {
	fmt.Println("Graphql tut")

	tutorials := populate()
	//A GraphQL object type has a name and fields*

	//Comment Type Object
	var commentType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"body": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.Int), //New list for array of int
				},
			},
		})

	var tutorialType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorials",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	//Define what fields on objects we want returned to us, so we have to define these fields within our Schema.
	fields := graphql.Fields{

		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get tutorial by id", //we want to be able to specify the ID of the tutorial we want to retrieve
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			//Resolver function that is triggered whenever this particular field is requested.
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int) //Takes in the Id Arguement
				if ok {
					for _, tutorial := range tutorials { //Check tutorial array for similar id
						if int(tutorial.Id) == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get full tutorial list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	//Define object config
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}

	//Define schema config (Query acts as the entry point for each graphQl entering the graphQl app)
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	//Creates schema
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("Failed to create graphQL schema, err %v", err)
	}

	query := `
    {
        list {
            id
            title
            comments {
                body
            }
            author {
                Name
                Tutorials
            }
        }
    }
`
	//Pass in the schema and q
	params := graphql.Params{Schema: schema, RequestString: query}

	r := graphql.Do(params)

	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphQL operation, errors %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s\n", rJSON)
}
