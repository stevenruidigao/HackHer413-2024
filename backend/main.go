package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"

	"google.golang.org/api/option"
)

type PotentialActionResult struct {
	Text        string  `json:"text"`
	Probability float64 `json:"probability"`
}

type PotentialAction struct {
	Description      string                  `json:"description"`
	TimeCost         int                     `json:"time_cost"`
	PotentialResults []PotentialActionResult `json:"potential_results"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type AIResponse struct {
	Outcome     string            `json:"outcome"`
	Scenario    string            `json:"-"`
	Inventory   []Item            `json:"inventory"`
	GameTime    int               `json:"game_time"`
	Stats       map[string]int    `json:"stats"`
	Skills      []Skill           `json:"skills"`
	NextActions []PotentialAction `json:"next_actions"`
}

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	geminiAPIKey := viper.GetString("GEMINI_API_KEY")
	log.Println(geminiAPIKey)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	// Initialize the chat
	cs := model.StartChat()
	cs.History = []*genai.Content{
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(`{
					"action": "Begin the game.",
					"scenario": "A dragon has abducted the prince.",
					"inventory": [],
					"game_time": 0,
					"stats": {
						"INT": 1,
						"LUK": 1,
						"STR": 1,
					},
					"skills": []
				}
				
				Respond only in JSON. Do not include anything else in the response.`),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(`You are a storytelling game master. The user will tell you what they do (in JSON), and you will respond with the result (in JSON) and possible next actions.

				Example of user input:
				{
				 	"action": "Inspect the situation",
					"scenario": "Fighting a dragon",
					"inventory": [
						{
							"name": "potato",
							"description": "A potato.",
						},
						{
							"name": "wand",
							"description": "A magic wand.",
						}, 
						{
							"name": "computer",
							"description": "A Dell laptop.",
						}
					],
					"game_time": 10,
					"stats": {
						"INT": 100,
						"LUK": 30,
						"STR": 4,
					},
					"skills": [{
						"name": "Programming",
						"description": "Can code to defeat computer viruses.",
						"level": 10,
					}]
				}
				
				Example of a response you can give (in JSON):
				{
				  "outcome": "Outcome of action",
				  "scenario": "Updated scenario",
				  "inventory": [
					  {
						  "name": "potato",
						  "description": "A potato.",
					  },
					  {
						  "name": "wand",
						  "description": "A magic wand.",
					  }, 
					  {
						  "name": "computer",
						  "description": "A Dell laptop.",
					  }
				  ],
				  "game_time": 11,
				  "stats": {
					  "INT": 100,
					  "LUK": 30,
					  "STR": 4,
				  },
				  "skills": [{
					  "name": "Programming",
					  "description": "Can code to defeat computer viruses.",
					  "level": 10,
				  }],
				  "next_actions": [
					{
					  "description": "Go for the jugular (high risk, high reward)",
					  "time_cost": 1,
					  "potential_results": [
						{
						  "text": "You land a critical blow, dealing devastating damage! But beware the dragon's fiery breath!",
						  "probability": 0.25
						},
						{
						  "text": "The dragon deflects your attack and retaliates with a powerful swipe!",
						  "probability": 0.5
						},
						{
						  "text": "Your aim falters, missing the vulnerable spot entirely.",
						  "probability": 0.25
						}
					  ]
					},
					{
					  "description": "Weaken its defenses (moderate risk, moderate reward)",
					  "time_cost": 2,
					  "potential_results": [
						{
						  "text": "You manage to cripple a wing, hindering the dragon's flight and maneuverability!",
						  "probability": 0.35
						},
						{
						  "text": "Your attacks chip away at its scales, slowly wearing it down and exposing weak points.",
						  "probability": 0.5
						},
						{
						  "text": "The dragon shrugs off your blows, its thick hide proving resilient. Be wary of its tail swing!",
						  "probability": 0.15
						}
					  ]
					},
					{
					  "description": "Distract and escape (low risk, low reward)",
					  "time_cost": 1,
					  "potential_results": [
						{
						  "text": "You successfully divert the dragon's attention with a well-placed object, creating an opening to flee!",
						  "probability": 0.4
						},
						{
						  "text": "The distraction fails, angering the dragon further! It unleashes a fiery breath in your direction!",
						  "probability": 0.4
						},
						{
						  "text": "You stumble during your escape attempt, leaving yourself vulnerable to the dragon's sharp claws.",
						  "probability": 0.2
						}
					  ]
					}
				  ]
				}
				
				Respond only in JSON. Do not include anything else in the response. Do not allow the user to significantly modify the state of the game without good reason.`),
			},
			Role: "model",
		},
	}

	resp, err := cs.SendMessage(ctx, genai.Text(`{
		"action": "Assess the situation.",
		"scenario": "A dragon has abducted the prince.",
		"inventory": [],
		"game_time": 0,
		"stats": {
			"INT": 1,
			"LUK": 1,
			"STR": 1,
		},
		"skills": []
	}
				
	Respond only in JSON. Do not include anything else in the response.`))

	if err != nil {
		log.Fatal(err)
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)

	if !ok {
		log.Fatal("wrong")
	}

	AIResponse := AIResponse{}

	err = json.Unmarshal([]byte(text), &AIResponse)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(AIResponse)

	mux := http.NewServeMux()

	mux.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
