import React, { useState } from 'react';

import './App.css';

import Menu from './Components/Menu/Menu.js'
import ClickerGame from './Components/Minigame/ClickerGame.js'

function App() {
  const [minigame, setMinigame] = useState(true);
  const [playerTurn, setPlayerTurn] = useState(true);
  
  const [attackPower, setAttackpower] = useState(0);
  function incrPower(x) {
    setAttackpower(attackPower + x);
  }

  //send data/receive data
  function send(convId, userAction, name, scenario) {
    console.log("SENDING THE USER INPUT UWUW: "+userAction)
    let url = "http://localhost:8080/submit";
    const response = fetch(url, {
      method: "POST", // *GET, POST, PUT, DELETE, etc.
      headers: {
        "Content-Type": "application/json",
        // 'Content-Type': 'application/x-www-form-urlencoded',
      },
      
      body: JSON.stringify({
        conversation_id: convId,
        action: userAction,
        name: name,
        scenario: scenario
      }),
    }).then(data => data.json()).then(json => {
      //update state
      setAttackpower(99);
      console.log("SERVER DATA RECEIVED")
    });
  }


  return (
    <div className="App">
      <div>
        <Menu tab={1} send={(userAction)=>{send("", userAction, "arky", "The player is fighting against Kafka, a mind controlling fugitive who knows the secrets of the universe.")}}/>
        <h1>AIdventure</h1>

        {minigame ? (<ClickerGame handleClick = {()=>{incrPower(0.5)}}/>) : (<></>)}
  
      </div>
    </div>
  );
}


export default App;
