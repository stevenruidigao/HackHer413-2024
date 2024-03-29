package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
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
	Quantity    int    `json:"quantity"`
}

type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type Character struct {
	Name      string         `json:"name,omitempty"`
	Description string `json:"description"`
	Inventory []Item         `json:"inventory"`
	Stats     map[string]int `json:"stats"`
	Skills    []Skill        `json:"skills"`
}

type Player struct {
	Character
}

type NPC struct {
	Character
//	Description string `json:"description"`
}

type GameStatePublic struct {
	GameTime int    `json:"game_time"`
	Player   Player `json:"player"`
	NPCs     []NPC  `json:"npcs"`
	IsOver bool `json:"is_over"`
}

type GameState struct {
	GameStatePublic
	Scenario string `json:"scenario"`
}

type OutputCharacter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ItemsLost   []Item `json:"items_lost"`
	ItemsGained []Item `json:"items_gained"`
	DamageTaken int    `json:"damage_taken"`
}

type OutcomeOutput struct {
	Outcome     string            `json:"outcome"`
	Scenario    string            `json:"scenario"`
	IsOver bool `json:"is_over"`
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
	ConversationID string `json:"conversation_id"`
	Action         string `json:"action"`
	Name           string `json:"name"`
	Scenario       string `json:"scenario"`
}

func main() {
	IDToChat := make(map[string]*Chat)

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	geminiAPIKey := viper.GetString("GEMINI_API_KEY")
	log.Println(geminiAPIKey)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))

	if err != nil {
		log.Println(err)
	}

	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	model.Temperature = genai.Ptr[float32](0.6)

	model.SafetySettings = []*genai.SafetySetting{
		&genai.SafetySetting{
			Category: genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

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

		log.Println("Conversation ID:", requestData.ConversationID, "Length:", len(requestData.ConversationID), requestData)

		if len(requestData.ConversationID) == 0 {
			requestData.ConversationID = uuid.NewString()
		}
		//generate skillset
		skills := []Skill{
			Skill{Name: "Kung Fu", Description: "A skill that makes you unreasonably good at fighting in close combat, especially with no weapon.", Level: 10},
			Skill{Name: "Sharp Shot", Description: "A skill that makes you good with shooting guns and other ranged weapons.", Level: 10},
			Skill{Name: "Slick Talker", Description: "This skill makes it easy for you to negotiate and convince others.", Level: 10},
			Skill{Name: "Streetwise", Description: "This skill lets you know the underbelly of society, its criminal networks, hidden secrets, and the people who navigate it.", Level: 10},
			Skill{Name: "Algorithmic Intelligence", Description: "This skill makes it easy to solve math and CS problems.", Level: 10},
			Skill{Name: "Acrobatics", Description: "This skill makes it easy for you to dodge, move, and jump around with grace.", Level: 10},
			Skill{Name: "Ninjitsu", Description: "This skill makes you sneaky and extremely proficient with swords.", Level: 10},
		}

		//generate items
		items := []Item{
			Item{Name: "Samurai Sword", Description: "A steel blade that hisses as its unsheathed.", Quantity: 1},
			Item{Name: "Black Cleaver", Description: "A dark and sinister blade, used for murder.", Quantity: 1},
			Item{Name: "Vangaurd's Shield", Description: "A powerful shield meant to protect.", Quantity: 1},
			Item{Name: "Freeze Ray", Description: "Raygun which is good for immobilizing.", Quantity: 1},
			Item{Name: "Golden Apples", Description: "Mystical fruit that increase your strength and defense.", Quantity: 3},
			Item{Name: "Stick", Description: "Just a stick. For good luck.", Quantity: 1},
			Item{Name: "Hypnoshroom", Description: "A special mushroom used to hypnotize an enemy. Consumed on use.", Quantity: 1},
			Item{Name: "Rocket Launcher", Description: "A great weapon for taking out multiple targets, but slow to shoot.", Quantity: 1},
			Item{Name: "Dual Pistols", Description: "A set of pistols, great for mid ranged combat.", Quantity: 2},
			Item{Name: "Smoke Grenade", Description: "Good for clouding up enemy vision and gaining the element of surprise.", Quantity: 4},
			Item{Name: "Currency", Description: "Maybe buy a snack?", Quantity: 50},
			Item{Name: "Mirror", Description: "Good for self reflection.", Quantity: 1},
		}

		if IDToChat[requestData.ConversationID] == nil {
			gameState := GameState{
				Scenario: requestData.Scenario,
				GameStatePublic: GameStatePublic{
					GameTime: 0,
					Player: Player{
						Character: Character{
							Name:      requestData.Name,
							Inventory: []Item{items[rand.Intn(len(items))], items[rand.Intn(len(items))], items[rand.Intn(len(items))]},
							Stats: map[string]int{
								"CHR": rand.Intn(11),
								"CON": rand.Intn(11),
								"DEX": rand.Intn(11),
								"INT": rand.Intn(11),
								"STR": rand.Intn(11),
								"WIS": rand.Intn(11),
								"LUK": rand.Intn(11),
								"HP":  10,
								"MAX_HP": 10,
							},
							Skills: []Skill{skills[rand.Intn(len(skills))]},
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

			log.Println(genai.Text(string(jsonAction) + "\n\n" + PROMPT_POSTFIX))

			if err != nil {
				log.Println("Error marshalling action to JSON:", err)
			}

			chat := Chat{
				ConversationID: requestData.ConversationID,
				GameState:      gameState,
				History: []*genai.Content{
					&genai.Content{
						Parts: []genai.Part{
							genai.Text(string(jsonAction) + "\n\n" + PROMPT_POSTFIX),
						},
						Role: "user",
					},
					&genai.Content{
						Parts: []genai.Part{
							genai.Text(SYSTEM_PROMPT + "\n\n" + PROMPT_POSTFIX),
						},
						Role: "model",
					},
				},
			}

			IDToChat[chat.ConversationID] = &chat
			requestData.ConversationID = chat.ConversationID
		}

		// construct the input for AI api
		chat := *IDToChat[requestData.ConversationID]

		cs := model.StartChat()

		// if there are no current NPCS
		if len(chat.GameState.GameStatePublic.NPCs) == 0 || chat.GameState.GameStatePublic.IsOver {
			for i := 0; i < 3; i++ {
				log.Println("Generating NPCs")

				var resp *genai.GenerateContentResponse

				cs.History = []*genai.Content{
					&genai.Content{
						Parts: []genai.Part{
							genai.Text(PROMPT_POSTFIX),
						},
						Role: "user",
					},
					&genai.Content{
						Parts: []genai.Part{
							genai.Text(GENERATE_NPCS_PROMPT),
						},
						Role: "model",
					},
				}

				gameStateJSON, err := json.Marshal(chat.GameState)

				if err != nil {
					log.Println("Error marshalling input to JSON:", err)
					return
				}

				resp, err = cs.SendMessage(ctx, genai.Text(gameStateJSON))

				if err != nil {
					log.Println(err)
					return
				}

				var text string = ""
				NPCs := []NPC{}

				for j := 0; j < len(resp.Candidates); j++ {
					for k := 0; k < len(resp.Candidates[j].Content.Parts); k++ {
						part, ok := resp.Candidates[j].Content.Parts[k].(genai.Text)

						if !ok {
							log.Println("The response was not a text response.")
							err = errors.New("The response was not a text response.")
							continue
						}

						text = text + string(part)
					}

					log.Println(text)
					err = json.Unmarshal([]byte(text), &NPCs)

					if err != nil {
						log.Println(err)
						continue
					}

					// modify the input
					chat.GameState.NPCs = NPCs

					for i := 0; i < len(NPCs); i++ {
						if chat.GameState.NPCs[i].Inventory == nil {
							chat.GameState.NPCs[i].Inventory = []Item{}
						}

						if chat.GameState.NPCs[i].Stats == nil {
							chat.GameState.NPCs[i].Stats = map[string]int{}
						}

						if chat.GameState.NPCs[i].Skills == nil {
							chat.GameState.NPCs[i].Skills = []Skill{}
						}
					}

					break
				}

				if err != nil {
					log.Println(err)
					continue
				}

				break
			}
		}

		log.Println(chat.GameState.GameStatePublic.Player.Character.Stats)

		action := ActionInput{
			Action:    requestData.Action,
			GameState: chat.GameState,
		}

		inputJSON, err := json.Marshal(action)

		if err != nil {
			log.Println("Error marshalling input to JSON:", err)
		}

		log.Println("Input JSON:", string(inputJSON))

		var text string
		var AIResp OutcomeOutput
		err = nil

		for i := 0; i < 3; i++ {
			var resp *genai.GenerateContentResponse

			cs.History = chat.History
			resp, err = cs.SendMessage(ctx, genai.Text(inputJSON))

			if err != nil {
				log.Println("Error sending message:", err)
				continue
			}

			for j := 0; j < len(resp.Candidates); j++ {
				text = ""

				for k := 0; k < len(resp.Candidates[j].Content.Parts); k++ {
					part, ok := resp.Candidates[j].Content.Parts[k].(genai.Text)

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

		if err != nil {
			return
		}

		chat.GameState.Scenario = AIResp.Scenario
		chat.GameState.GameStatePublic.IsOver = AIResp.IsOver
		chat.GameState.GameStatePublic.Player.Character.Name = AIResp.Player.Name
		chat.GameState.GameStatePublic.Player.Character.Description = AIResp.Player.Description
		chat.GameState.GameStatePublic.Player.Character.Stats["HP"] -= AIResp.Player.DamageTaken

		if chat.GameState.GameStatePublic.Player.Character.Stats["HP"] <= 0 {
			chat.GameState.GameStatePublic.IsOver = true
		}

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

					if chat.GameState.GameStatePublic.NPCs[j].Stats["HP"] <= 0 {
						newNPCs := make([]NPC, 0)
						newNPCs = append(newNPCs, chat.GameState.GameStatePublic.NPCs[:j]...)
						newNPCs = append(newNPCs, chat.GameState.GameStatePublic.NPCs[j+1:]...)
						chat.GameState.GameStatePublic.NPCs = newNPCs
					}

					for k := 0; k < len(AIResp.NPCs[i].ItemsGained); k++ {
						chat.GameState.GameStatePublic.NPCs[j].Character.Inventory = append(chat.GameState.GameStatePublic.NPCs[j].Character.Inventory, AIResp.NPCs[i].ItemsGained[k])
					}

					for k := 0; k < len(AIResp.NPCs[i].ItemsLost); k++ {
						for l := 0; l < len(chat.GameState.GameStatePublic.NPCs[l].Character.Inventory); l++ {
							if chat.GameState.GameStatePublic.NPCs[j].Character.Inventory[l].Name == AIResp.NPCs[i].ItemsLost[k].Name {
								newInventory := make([]Item, 0)
								newInventory = append(newInventory, chat.GameState.GameStatePublic.NPCs[j].Character.Inventory[:l]...)
								newInventory = append(newInventory, chat.GameState.GameStatePublic.NPCs[j].Character.Inventory[l+1:]...)
								chat.GameState.GameStatePublic.NPCs[j].Character.Inventory = newInventory

								break
							}
						}
					}
				}
			}
		}


		chat.History = cs.History
		IDToChat[requestData.ConversationID] = &chat

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
