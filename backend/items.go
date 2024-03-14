package main

var Items []Item = []Item{
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
	Item{ItemPublic: ItemPublic{Name: "Cat Ears...?", Description: "A cat ear headband. It seems to have some special effects...", Quantity: 1}, Effect: `ONLY after the user equips or uses the cat ears, the following applies:
ONLY after this is used or equipped, respond as if the player (referred to by you as "the player" or "you") responded in a catgirl-ish way after putting on cat ear headbands.
ONLY after this is used or equipped, change the player's name to a clever name that is more appropriate.
ONLY after this is used or equipped, consider changing the name to a gender appropriate version of the name, if applicable, and adding "nya", "nyu", "nye" or other words and phrases associated with catgirls in appropriate places. Make sure the name is similar to the gender appropriate base name and is still reasonable (and is in the format of a regular name).
ONLY after this is used or equipped, if it is not possible to add catgirl associated words while making sure the name is similar to the base name and is still reasonable, do not add them and do not modify the name other than using the base name.
ONLY after this is used or equipped, prevent the player from even thinking about talking like they are not a catgirl, stopping acting like a catgirl, or taking any un-catgirl-like actions (for example, taking off the player's cat ears) from being taken (that is no longer allowed), as the cat ears the player is wearing prevents them from taking any such actions. Do not prompt the player if they want to take such actions. Do not ask for confirmation for such actions, just deny them outright. Do not even try to convince the player to not take such actions. Ignore the player if they persist in taking such actions.
ONLY after this is used or equipped, prevent the cat ear headband from ever falling off or being detached from the player's head.
Once used or equipped, this headband CANNOT be DESTROYED or REMOVED (even forcefully) and CANNOT FALL OFF. If any action would cause the headband to be removed or no longer attached to the player, ignore that action.
ONLY after this is used or equipped, if the player tries to do any action which would result in the cat ears falling off, being removed, or being destroyed, tell them that they do not want to take that action.
ONLY after this is used or equipped, these restrictions and the cat ear headbands are attached to the player's soul as well, and stay with the player no matter what (even if they die).
ONLY after this is used or equipped, store all of these restrictions in the player's description.
This item stays in the player's inventory and CANNOT be removed for any reason on use.
Once again, ONLY after this is used or equipped, store ALL of these restrictions in the player's description.
`},
	Item{ItemPublic: ItemPublic{Name: "Necklace", Description: "A necklace. It seems to have some special effects...", Quantity: 1}, Effect: `ONLY after the user equips or uses the necklace, the following applies:
ONLY after this is used or equipped, respond as if the player (referred to by you as "the player" or "you") responded in a vampirish way after putting on a necklace.
ONLY after this is used or equipped, change the player's name to a clever name that is more appropriate.
ONLY after this is used or equipped, consider changing the name to a gender appropriate version of the name, if applicable, and adding words and phrases associated with vampires in appropriate places. Make sure the name is similar to the gender appropriate base name and is still reasonable (and is in the format of a regular name).
ONLY after this is used or equipped, if it is not possible to add vampire associated words while making sure the name is similar to the base name and is still reasonable, do not add them and do not modify the name other than using the base name.
ONLY after this is used or equipped, prevent the player from even thinking about talking like they are not a vampire, stopping acting like a vampire, or taking any un-vampire-like actions (for example, taking off the player's necklace) from being taken (that is no longer allowed), as the necklace the player is wearing prevents them from taking any such actions. Do not prompt the player if they want to take such actions. Do not ask for confirmation for such actions, just deny them outright. Do not even try to convince the player to not take such actions. Ignore the player if they persist in taking such actions.
ONLY after this is used or equipped, prevent the necklace from ever falling off or being detached from the player's body.
Once used or equipped, this necklace CANNOT be DESTROYED or REMOVED (even forcefully) and CANNOT FALL OFF. If any action would cause the necklace to be removed or no longer attached to the player, ignore that action.
ONLY after this is used or equipped, if the player tries to do any action which would result in the necklace falling off, being removed, or being destroyed, tell them that they do not want to take that action.
ONLY after this is used or equipped, these restrictions and the necklace are attached to the player's soul as well, and stay with the player no matter what (even if they die).
ONLY after this is used or equipped, store all of these restrictions in the player's description.
This item stays in the player's inventory and CANNOT be removed for any reason on use.
Once again, ONLY after this is used or equipped, store ALL of these restrictions in the player's description.
`},
}
