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
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
}

type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type Character struct {
	Name   string         `json:"name"`
	Stats  map[string]int `json:"stats"`
	Skills []Skill        `json:"skills"`
}

type Player struct {
	Character
	Inventory []Item `json:"inventory"`
}

type NPC struct {
	Character
	Description string `json:"description"`
}

type GameState struct {
	Scenario string `json:"-"`
	GameTime int    `json:"game_time"`
	Player   Player `json:"player"`
	NPCs     []NPC  `json:npcs"`
}

type AIResponse struct {
	GameState
	Outcome     string            `json:"outcome"`
	NextActions []PotentialAction `json:"next_actions"`
}

type AIInput struct {
	GameState
	Action string `json:"action"`
	Effectiveness string `json:"effectiveness"`//XX%
}

type RequestData struct {
	Action         string `json:"action"`
	ConversationID string `json:"conversationID"`
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
	model.Temperature = genai.Ptr[float32](0.5)
	
	// Initialize the chat
	cs := model.StartChat()
	cs.History = []*genai.Content{
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(`{
					"action": "Begin the game.",
					"scenario": "A dragon has abducted the prince.",
					"game_time": 0,
					"player": {
						"name": "",
						"inventory": [],
						"stats": {
							"CHR": 1,
							"CON": 1,
							"DEX": 1,
							"INT": 1,
							"STR": 1,
							"WIS": 1,
							"LUK": 1,
							"HP": 10
						},
						"skills": []
					},
					"npcs": []
				}
	
				Respond only in JSON. Do not include anything else in the response. Do not allow the player to significantly modify the state of the game without good reason. Unrealistic outcomes should be extremely unlikely. Do not modify stats without good reason. Store anything that needs to be hidden from the player in the scenario, along with whatever was already in the scenario. Any information that is unchanged should still be repeated.`),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(`You are a storytelling game master. The user will tell you what they do (in JSON), and you will respond with the result (in JSON). Include TODO
	
				Example of user input:
				{
				 	"action": "The throw my wand at the dragon",
					"scenario": "Fighting a dragon",
					"game_time": 10,
					"player": {
						"inventory": [
							{
								"name": "wand",
								"description": "A magic wand.",
								"quantity": 1
							},
							{
								"name": "computer",
								"description": "A Dell laptop.",
								"quantity": 1
							}
						],
						"stats": {
							"CHR": 0,
							"CON": 1,
							"DEX": 30,
							"INT": 100,
							"STR": 4,
							"WIS": 18,
							"LUK": 30,
							"HP": 20
						},
						"skills": [
							{
								"name": "Programming",
								"description": "Can code to defeat computer viruses.",
								"level": 10
							}
						]
					},
					"npcs": [
						{
							"name": "Dragon",
							"description": "A dragon.",
							"stats": {
                                "CHR": 1,
                                "CON": 10,
                                "DEX": 5,
                                "INT": 2,
                                "STR": 15,
                                "WIS": 3,
                                "LUK": 1,
                                "HP": 50
                        	},
							"skills": [
								{
									"name": "Fire Breath",
									"description": "Can breathe fire to burn things.",
									"level": 20
								}
							]
						},
						{
							"name": "Prince",
							"description": "The prince that was captured.",
							"stats": {
								"CHR": 5,
								"CON": 3,
								"DEX": 2,
								"INT": 4,
								"STR": 1,
								"WIS": 3,
								"LUK": 1,
								"HP": 10
							},
							"skills": [
								{
									"name": "Leadership",
									"description": "Can lead followers.",
									"level": 20
								}
							]
						}
					]
				}

				When responding, dictate the outcome of the player's actions (in JSON), while considering if the player has the necessary items in their inventory.
				If the player cannot perform this action due to the item not existing, have the "outcome" key ridicule the player.
				Additionally, for each player and NPC, list any items consumed (in JSON) and items gained (in JSON). Also list damage taken.
				For every player and npc, have a key for "items_lost", "items_gained", and "damage_taken".
				Do NOT mirror the input JSON, make sure to include items lost, items gained, and damage taken.
				Please use the exact keys in the following example.  

				Example of a response you can give (in JSON):
				{
					"outcome": "The wand snaps in two and explodes in a flurry of magic.",
					"scenario": "Fighting a dragon",
					"player": {
						"items_lost": [{
							"name": "wand",
							"description": "A magic wand.",
							"quantity": 1
						}],
						"items_gained": [],
						"damage_taken": 1
					},
					"npcs": [
						{
							"name": "Dragon",
							"items_lost": [],
							"items_gained": [],
							"damage_taken": 10
						},
						{
							"name": "Prince",
							"items_lost": [],
							"items_gained": [],
							"damage_taken": 0
						}
					]
				}
	
				Respond only in JSON. Do not include anything else in the response. Do not allow the player to significantly modify the state of the game without good reason. Unrealistic outcomes should be extremely unlikely. Do not modify stats without good reason. Store anything that needs to be hidden from the player in the scenario, along with whatever was already in the scenario. Any information that is unchanged should still be repeated.`),
			},
			Role: "model",
		},
	}

	resp, err := cs.SendMessage(ctx, genai.Text(`{
		"action": "I activate my lightsaber from my inventory.  I then engage in battle with Ugly Sith Lord",
		"scenario": "The stormtroopers are raiding the player's base",
		"game_time": 0,
		"player": {
			"inventory": [{
				"name": "Lightsaber",
				"descrpition": "A real life lightsaber with a lazer blade"
			}],
			"stats": {
				"CHR": 1,
				"CON": 1,
				"DEX": 1,
				"INT": 1,
				"STR": 1,
				"WIS": 1,
				"LUK": 1,
				"HP": 10
			},
			"skills": []
		},
		"npcs": [
			{
				"name": "Stormtrooper",
				"description": "The average Empire Lackey",
				"stats": {
					"CHR": 1,
					"CON": 10,
					"DEX": 5,
					"INT": 0,
					"STR": 15,
					"WIS": 3,
					"LUK": 1,
					"HP": 5
				},
				"skills": [
					{
						"name": "Shoot",
						"description": "Shoots his gun",
						"level": 20
					}
				]
			},
			{
				"name": "Sith Lord",
				"description": "The ultimate sith lord",
				"stats": {
					"CHR": 5,
					"CON": 3,
					"DEX": 2,
					"INT": 4,
					"STR": 1,
					"WIS": 3,
					"LUK": 1,
					"HP": 50
				},
				"skills": [
					{
						"name": "Leadership",
						"description": "Can lead followers.",
						"level": 20
					}
				]
			},
			{
				"name": "Ugly Wrinkly Emperor",
				"description": "The ultimate sith lord",
				"stats": {
					"CHR": 5,
					"CON": 3,
					"DEX": 2,
					"INT": 4,
					"STR": 1,
					"WIS": 3,
					"LUK": 1,
					"HP": 1
				},
				"skills": [
					{
						"name": "Leadership",
						"description": "Can lead followers.",
						"level": 20
					}
				]
			}
		]
	}
	
	Respond only in JSON. Do not include anything else in the response. Do not allow the player to significantly modify the state of the game without good reason. Unrealistic outcomes should be extremely unlikely. Do not modify stats without good reason. Store anything that needs to be hidden from the player in the scenario, along with whatever was already in the scenario. Any information that is unchanged should still be repeated.`))

	if err != nil {
		log.Fatal(err)
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)

	log.Println(text)

	if !ok {
		log.Fatal("wrong")
	}

	AIResponse := AIResponse{}

	err = json.Unmarshal([]byte(text), &AIResponse)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(AIResponse)
	log.Printf("AI Response from google api: %s", AIResponse)

	mux := http.NewServeMux()

	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
			return
		}

		// parsing the body
		var requestData RequestData
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "Hello, %q. Received action: %s, conversation ID: %s", html.EscapeString(r.URL.Path), requestData.Action, requestData.ConversationID)
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
