package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
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

type Character struct {
	Name      string         `json:"name,omitempty"`
	Inventory []Item         `json:"inventory"`
	Stats     map[string]int `json:"stats"`
	Skills    []Skill        `json:"skills"`
}

type Player struct {
	Character
}

type NPC struct {
	Character
	Description string `json:"description"`
}

type GameStatePublic struct {
	GameTime int    `json:"game_time"`
	Player   Player `json:"player"`
	NPCs     []NPC  `json:"npcs"`
}

type GameState struct {
	GameStatePublic
	Scenario string `json:"scenario"`
}

type OutputCharacter struct {
	Name        string `json:"name"`
	ItemsLost   []Item `json:"items_lost"`
	ItemsGained []Item `json:"items_gained"`
	DamageTaken int    `json:"damage_taken"`
}

type OutcomeOutput struct {
	Outcome     string            `json:"outcome"`
	Scenario    string            `json:"scenario"`
	Player      OutputCharacter   `json:"player"`
	NPCs        []OutputCharacter `json:"npcs"`
	NextActions []PotentialAction `json:"next_actions"`
}

type ActionInput struct {
	GameState
	Action string `json:"action"`
}

type Chat struct {
	ConversationID string           `json:"conversation_id"`
	GameState      GameState        `json:"game_state"`
	History        []*genai.Content `json:"-"`
}

type ChatResponse struct {
	ConversationID string          `json:"conversation_id"`
	GameState      GameStatePublic `json:"game_state"`
	Outcome        string          `json:"outcome"`
}

type RequestData struct {
	ConversationID string `json:"conversationID"`
	Action         string `json:"action"`
	Name           string `json:"name"`
	Scenario       string `json:"scenario"`
}

func main() {
	IDToChat := make(map[string]Chat)

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

	/* Comment used to be here */

	mux := http.NewServeMux()

	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// parsing the body
		var requestData RequestData

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		log.Println("Conversation ID:", requestData.ConversationID, "Length:", len(requestData.ConversationID))

		if len(requestData.ConversationID) == 0 {
			gameState := GameState{
				Scenario: requestData.Scenario,
				GameStatePublic: GameStatePublic{
					GameTime: 0,
					Player: Player{
						Character: Character{
							Name:      requestData.Name,
							Inventory: []Item{},
							Stats: map[string]int{
								"CHR": 1,
								"CON": 1,
								"DEX": 1,
								"INT": 1,
								"STR": 1,
								"WIS": 1,
								"LUK": 1,
								"HP":  10,
							},
							Skills: []Skill{},
						},
					},
					NPCs: []NPC{},
				},
			}

			action := ActionInput{
				Action:    "Begin the game.",
				GameState: gameState,
			}

			jsonAction, err := json.Marshal(action)

			log.Println(genai.Text(strings.Join([]string{string(jsonAction), PROMPT_POSTFIX}, "\n\n")))

			if err != nil {
				log.Fatal("Error marshalling action to JSON:", err)
			}

			chat := Chat{
				ConversationID: uuid.NewString(),
				GameState:      gameState,
				History: []*genai.Content{
					&genai.Content{
						Parts: []genai.Part{
							genai.Text(strings.Join([]string{string(jsonAction), PROMPT_POSTFIX}, "\n\n")),
						},
						Role: "user",
					},
					&genai.Content{
						Parts: []genai.Part{
							genai.Text(strings.Join([]string{SYSTEM_PROMPT, PROMPT_POSTFIX}, "\n\n")),
						},
						Role: "model",
					},
				},
			}

			IDToChat[chat.ConversationID] = chat
			requestData.ConversationID = chat.ConversationID
		}

		// construct the input for AI api
		chat := IDToChat[requestData.ConversationID]

		log.Println(chat.GameState.GameStatePublic.Player.Character.Stats)

		action := ActionInput{
			Action:    requestData.Action,
			GameState: chat.GameState,
		}

		inputJSON, err := json.Marshal(action)

		if err != nil {
			log.Fatal("Error marshalling input to JSON:", err)
		}

		log.Println("Input JSON:", string(inputJSON))

		cs := model.StartChat()

		var text string
		var AIResp OutcomeOutput

		for i := 0; i < 3; i++ {
			var err error = nil

			cs.History = chat.History
			resp, err := cs.SendMessage(ctx, genai.Text(inputJSON))

			if err != nil {
				log.Println("Error sending message:", err)
				continue
			}

			for j := 0; j < len(resp.Candidates); j++ {
				text = ""

				for i := 0; i < len(resp.Candidates[0].Content.Parts); i++ {
					part, ok := resp.Candidates[0].Content.Parts[i].(genai.Text)

					if !ok {
						log.Println("The response was not a text response.")
						err = errors.New("The response was not a text response.")
						continue
					}

					text = strings.Join([]string{text, string(part)}, "")
				}

				log.Println(text)
				err = json.Unmarshal([]byte(text), &AIResp)

				if err != nil {
					log.Println(err)
					continue
				}

				break
			}

			if err != nil {
				continue
			}

			break
		}

		chat.GameState.Scenario = AIResp.Scenario
		chat.GameState.GameStatePublic.Player.Character.Stats["HP"] -= AIResp.Player.DamageTaken

		for i := 0; i < len(AIResp.Player.ItemsGained); i++ {
			chat.GameState.GameStatePublic.Player.Character.Inventory = append(chat.GameState.GameStatePublic.Player.Character.Inventory, AIResp.Player.ItemsGained[i])
		}

		for i := 0; i < len(AIResp.Player.ItemsLost); i++ {
			for j := 0; j < len(chat.GameState.GameStatePublic.Player.Character.Inventory); j++ {
				if chat.GameState.GameStatePublic.Player.Character.Inventory[j].Name == AIResp.Player.ItemsLost[i].Name {
					newInventory := make([]Item, 0)
					newInventory = append(newInventory, chat.GameState.GameStatePublic.Player.Character.Inventory[:j]...)
					newInventory = append(newInventory, chat.GameState.GameStatePublic.Player.Character.Inventory[j+1:]...)
					chat.GameState.GameStatePublic.Player.Inventory = newInventory

					break
				}
			}
		}

		for i := 0; i < len(AIResp.NPCs); i++ {
			for j := 0; j < len(chat.GameState.GameStatePublic.NPCs); j++ {
				if chat.GameState.GameStatePublic.NPCs[j].Name == AIResp.NPCs[i].Name {
					chat.GameState.GameStatePublic.NPCs[j].Stats["HP"] -= AIResp.NPCs[i].DamageTaken

					for i := 0; i < len(AIResp.NPCs[i].ItemsGained); i++ {
						chat.GameState.GameStatePublic.NPCs[j].Character.Inventory = append(chat.GameState.GameStatePublic.NPCs[j].Character.Inventory, AIResp.NPCs[i].ItemsGained[i])
					}

					for i := 0; i < len(AIResp.NPCs[i].ItemsLost); i++ {
						for j := 0; j < len(chat.GameState.GameStatePublic.NPCs[j].Character.Inventory); j++ {
							if chat.GameState.GameStatePublic.NPCs[j].Character.Inventory[j].Name == AIResp.NPCs[i].ItemsLost[i].Name {
								newInventory := make([]Item, 0)
								newInventory = append(newInventory, chat.GameState.GameStatePublic.NPCs[j].Character.Inventory[:j]...)
								newInventory = append(newInventory, chat.GameState.GameStatePublic.NPCs[j].Character.Inventory[j+1:]...)
								chat.GameState.GameStatePublic.NPCs[j].Character.Inventory = newInventory

								break
							}
						}
					}
				}
			}
		}

		chat.History = cs.History
		IDToChat[requestData.ConversationID] = chat

		chatResponse := ChatResponse{
			ConversationID: chat.ConversationID,
			GameState:      chat.GameState.GameStatePublic,
			Outcome:        AIResp.Outcome,
		}

		json.NewEncoder(w).Encode(chatResponse)
	})

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", http.StripPrefix("/", fs))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
