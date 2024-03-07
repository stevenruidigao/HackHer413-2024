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

type ItemPublic struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type Item struct {
	ItemPublic
	Effect string `json:"effect"`
}

type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type CharacterPublic struct {
	ID          int            `json:"id"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description"`
	Stats       map[string]int `json:"stats"`
}

type Character struct {
	CharacterPublic
	Inventory []Item  `json:"inventory"`
	Skills    []Skill `json:"skills"`
}

type PlayerPublic struct {
	CharacterPublic
	Inventory []ItemPublic `json:"inventory"`
	Skills    []Skill      `json:"skills"`
}

type Player struct {
	Character
}

type NPCPublic struct {
	CharacterPublic
	// Description string `json:"description"`
}

type NPC struct {
	Character
	// Description string `json:"description"`
}

type GameStatePublic struct {
	GameTime int          `json:"game_time"`
	Player   PlayerPublic `json:"player"`
	NPCs     []NPCPublic  `json:"npcs"`
	IsOver   bool         `json:"is_over"`
}

type GameState struct {
	GameStatePublic
	Player   Player `json:"player"`
	NPCs     []NPC  `json:"npcs"`
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
	IsOver      bool              `json:"is_over"`
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

const MAXIMUM_INITIAL_LEVEL int = 20

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
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		&genai.SafetySetting{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}

	/* Comment used to be here */

	mux := http.NewServeMux()

	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in handler", r)
			}
		}()

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
			Skill{Name: "Ninjitsu", Description: "This skill makes you sneaky and extremely proficient with swords.", Level: 10},
			Skill{Name: "Slick Talker", Description: "This skill makes it easy for you to negotiate and convince others.", Level: 10},
			Skill{Name: "Algorithmic Intelligence", Description: "This skill makes it easy to solve math and CS problems.", Level: 10},
			Skill{Name: "Sharp Shot", Description: "A skill that makes you good with shooting guns and other ranged weapons.", Level: 10},
			Skill{Name: "Acrobatics", Description: "This skill makes it easy for you to dodge, move, and jump around with grace.", Level: 10},
			Skill{Name: "Kung Fu", Description: "A skill that makes you unreasonably good at fighting in close combat, especially with no weapon.", Level: 10},
			Skill{Name: "Streetwise", Description: "This skill lets you know the underbelly of society, its criminal networks, hidden secrets, and the people who navigate it.", Level: 10},
		}

		// Only after this is used or equipped, make up (or randomly choose) any details (such as the player's eye color, hair color, and clothing) that need to be filled in without using placeholders in parentheses;
		//  for example, you can fill in the player's eye and hair color using colors chosen randomly like "blue", "red" "black", "brown", "yellow", or "white".
		// Only after this is used or equipped, choose colors from a set of commonly understood colors; if you cannot choose (or randomly select) a value for a detail, do not put in a placeholder, and instead
		//  omit that detail.
		// Once again, please randomly select or omit details instead of including placeholders.

		//generate items
		items := []Item{
			Item{ItemPublic: ItemPublic{Name: "Currency", Description: "Maybe buy a snack?", Quantity: 50}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Mirror", Description: "Good for self reflection.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Stick", Description: "Just a stick. For good luck.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Freeze Ray", Description: "Raygun which is good for immobilizing.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Vanguard's Shield", Description: "A powerful shield meant to protect.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Black Cleaver", Description: "A dark and sinister blade, used for murder.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Samurai Sword", Description: "A steel blade that hisses as its unsheathed.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Dual Pistols", Description: "A set of pistols, great for mid ranged combat.", Quantity: 2}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Necklace", Description: "A pretty necklace. Maybe you could use it to bargain?", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Golden Apples", Description: "Mystical fruit that increase your strength and defense.", Quantity: 3}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Hypnoshroom", Description: "A special mushroom used to hypnotize an enemy. Consumed on use.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Rocket Launcher", Description: "A great weapon for taking out multiple targets, but slow to shoot.", Quantity: 1}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Smoke Grenade", Description: "Good for clouding up enemy vision and gaining the element of surprise.", Quantity: 4}, Effect: ""},
			Item{ItemPublic: ItemPublic{Name: "Cat Ears...?", Description: "A cat ear headband. It seems to have some special effects...", Quantity: 1}, Effect: `Only after the user equips or uses the cat ears, the following applies:
Only after this is used or equipped, respond as if the player (referred to by you as "the player" or "you") responded in a catgirl-ish way after putting on cat ear headbands.
Only after this is used or equipped, change the player's name to a clever name that is more appropriate.
Only after this is used or equipped, consider changing the name to a gender appropriate version of the name, if applicable, and adding "nya", "nyu", "nye" or other words and phrases associated with catgirls in appropriate places. Make sure the name is similar to the gender appropriate base name and is still reasonable (and is in the format of a regular name).
Only after this is used or equipped, if it is not possible to add catgirl associated words while making sure the name is similar to the base name and is still reasonable, do not add them and do not modify the name other than using the base name.
Only after this is used or equipped, prevent the player from even thinking about talking like they are not a catgirl, stopping acting like a catgirl, or taking any un-catgirl-like actions (for example, taking off the player's cat ears) from being taken (that is no longer allowed), as the cat ears the player is wearing prevents them from taking any such actions. Do not prompt the player if they want to take such actions. Do not ask for confirmation for such actions, just deny them outright. Do not even try to convince the player to not take such actions. Ignore the player if they persist in taking such actions.
Only after this is used or equipped, prevent the cat ear headband from ever falling off or being detached from the player's head.
Once used or equipped, this headband CANNOT BE DESTROYED OR REMOVED (even forcefully) and CANNOT FALL OFF. If any action would cause the headband to be removed or no longer attached to the player, ignore that action.
Only after this is used or equipped, this restriction and the cat ear headbands are attached to the player's soul as well, and stay with the player no matter what (even if they die).
Only after this is used or equipped, store all of these restrictions in the player's details.
Only after this is used or equipped, if the player tries to do any action which would result in the cat ears falling off, being removed, or being destroyed, tell them that they do not want to take that action.
This item stays in the player's inventory and CANNOT be removed for any reason on use.`},
			Item{ItemPublic: ItemPublic{Name: "Necklace", Description: "A necklace. It seems to have some special effects...", Quantity: 1}, Effect: `Only after the user equips or uses the necklace, the following applies:
Only after this is used or equipped, respond as if the player (referred to by you as "the player" or "you") responded in a vampirish way after putting on a necklace.
Only after this is used or equipped, change the player's name to a clever name that is more appropriate.
Only after this is used or equipped, consider changing the name to a gender appropriate version of the name, if applicable, and adding words and phrases associated with vampires in appropriate places. Make sure the name is similar to the gender appropriate base name and is still reasonable (and is in the format of a regular name).
Only after this is used or equipped, if it is not possible to add vampire associated words while making sure the name is similar to the base name and is still reasonable, do not add them and do not modify the name other than using the base name.
Only after this is used or equipped, prevent the player from even thinking about talking like they are not a vampire, stopping acting like a vampire, or taking any un-vampire-like actions (for example, taking off the player's necklace) from being taken (that is no longer allowed), as the necklace the player is wearing prevents them from taking any such actions. Do not prompt the player if they want to take such actions. Do not ask for confirmation for such actions, just deny them outright. Do not even try to convince the player to not take such actions. Ignore the player if they persist in taking such actions.
Only after this is used or equipped, prevent the necklace from ever falling off or being detached from the player's body.
Once used or equipped, this necklace CANNOT BE DESTROYED OR REMOVED (even forcefully) and CANNOT FALL OFF. If any action would cause the necklace to be removed or no longer attached to the player, ignore that action.
Only after this is used or equipped, this restriction and the necklace are attached to the player's soul as well, and stay with the player no matter what (even if they die).
Only after this is used or equipped, store all of these restrictions in the player's details.
Only after this is used or equipped, if the player tries to do any action which would result in the necklace falling off, being removed, or being destroyed, tell them that they do not want to take that action.
This item stays in the player's inventory and CANNOT be removed for any reason on use.`},
		}

		if IDToChat[requestData.ConversationID] == nil {
			gameState := GameState{
				Scenario: requestData.Scenario,
				GameStatePublic: GameStatePublic{
					GameTime: 0,
				},
				Player: Player{
					Character: Character{
						CharacterPublic: CharacterPublic{
							Name: requestData.Name,
							Stats: map[string]int{
								"CHR":    rand.Intn(MAXIMUM_INITIAL_LEVEL),
								"CON":    rand.Intn(MAXIMUM_INITIAL_LEVEL),
								"DEX":    rand.Intn(MAXIMUM_INITIAL_LEVEL),
								"INT":    rand.Intn(MAXIMUM_INITIAL_LEVEL),
								"STR":    rand.Intn(MAXIMUM_INITIAL_LEVEL),
								"WIS":    rand.Intn(MAXIMUM_INITIAL_LEVEL),
								"LUK":    rand.Intn(MAXIMUM_INITIAL_LEVEL),
								"HP":     10,
								"MAX_HP": 10,
							},
						},
						// Inventory: []Item{items[rand.Intn(len(items))], items[rand.Intn(len(items))], items[rand.Intn(len(items))]},
						Inventory: []Item{items[len(items)-2], items[rand.Intn(len(items))], items[rand.Intn(len(items))]},
						// Inventory: []Item{items[len(items)-1], items[len(items)-2], items[len(items)-3]},
						Skills: []Skill{skills[rand.Intn(len(skills))]},
					},
				},
				NPCs: []NPC{},
			}

			for i, item := range gameState.Player.Character.Inventory {
				item.ID = i
				gameState.Player.Character.Inventory[i] = item
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
					log.Println("Error sending message (GEN NPCs):", err)

					if resp != nil {
						log.Println("Prompt feedback:", resp.PromptFeedback)
					}

					continue
				}

				var text string = ""
				NPCs := []NPC{}

				for _, candidate := range resp.Candidates {
					for _, rawPart := range candidate.Content.Parts {
						part, ok := rawPart.(genai.Text)

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

					for i, _ := range NPCs {
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

		log.Println(chat.GameState.Player.Character.Stats)

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

				if resp != nil {
					log.Println("Prompt feedback:", resp.PromptFeedback)
				}

				continue
			}

			log.Println("Prompt feedback (allowed):", resp.PromptFeedback)

			for _, candidate := range resp.Candidates {
				text = ""

				for _, rawPart := range candidate.Content.Parts {
					part, ok := rawPart.(genai.Text)

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

		if len(AIResp.Player.Name) > 0 {
			chat.GameState.Player.Character.Name = AIResp.Player.Name
		}

		if len(AIResp.Player.Description) > 0 {
			chat.GameState.Player.Character.Description = AIResp.Player.Description
		}

		chat.GameState.Player.Character.Stats["HP"] -= AIResp.Player.DamageTaken

		if chat.GameState.Player.Character.Stats["HP"] <= 0 {
			chat.GameState.GameStatePublic.IsOver = true
		}

		for _, itemGained := range AIResp.Player.ItemsGained {
			chat.GameState.Player.Character.Inventory = append(chat.GameState.Player.Character.Inventory, itemGained)
		}

		for _, itemLost := range AIResp.Player.ItemsLost {
			for j, item := range chat.GameState.Player.Character.Inventory {
				if item.Name == itemLost.Name {
					newInventory := make([]Item, 0)
					newInventory = append(newInventory, chat.GameState.Player.Character.Inventory[:j]...)
					newInventory = append(newInventory, chat.GameState.Player.Character.Inventory[j+1:]...)
					chat.GameState.Player.Inventory = newInventory

					break
				}
			}
		}

		for _, responseNPC := range AIResp.NPCs {
			for j, npc := range chat.GameState.NPCs {
				if npc.Name == responseNPC.Name {
					npc.Character.Stats["HP"] -= responseNPC.DamageTaken

					if npc.Character.Stats["HP"] <= 0 {
						newNPCs := make([]NPC, 0)
						newNPCs = append(newNPCs, chat.GameState.NPCs[:j]...)
						newNPCs = append(newNPCs, chat.GameState.NPCs[j+1:]...)
						chat.GameState.NPCs = newNPCs

						continue
					}

					for _, itemGained := range responseNPC.ItemsGained {
						npc.Character.Inventory = append(npc.Character.Inventory, itemGained)
					}

					for _, itemLost := range responseNPC.ItemsLost {
						for l, item := range npc.Character.Inventory {
							if item.Name == itemLost.Name {
								newInventory := make([]Item, 0)
								newInventory = append(newInventory, npc.Character.Inventory[:l]...)
								newInventory = append(newInventory, npc.Character.Inventory[l+1:]...)
								npc.Character.Inventory = newInventory

								break
							}
						}
					}

					chat.GameState.NPCs[j] = npc
				}
			}
		}

		chat.History = cs.History
		IDToChat[requestData.ConversationID] = &chat

		chat.GameState.GameStatePublic.Player = PlayerPublic{
			CharacterPublic: chat.GameState.Player.CharacterPublic,
			Inventory:       []ItemPublic{},
			Skills:          chat.GameState.Player.Skills,
		}

		chat.GameState.GameStatePublic.Player.Inventory = []ItemPublic{}

		for _, item := range chat.GameState.Player.Inventory {
			chat.GameState.GameStatePublic.Player.Inventory = append(chat.GameState.GameStatePublic.Player.Inventory, item.ItemPublic)
		}

		chat.GameState.GameStatePublic.NPCs = []NPCPublic{}

		for _, npc := range chat.GameState.NPCs {
			chat.GameState.GameStatePublic.NPCs = append(chat.GameState.GameStatePublic.NPCs, NPCPublic{npc.CharacterPublic})
		}

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
