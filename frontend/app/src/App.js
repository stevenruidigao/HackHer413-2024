import React, { useState } from 'react';

import './App.css';

import Menu from './Components/Menu/Menu.js'
import ClickerGame from './Components/Minigame/ClickerGame.js'
import Story from './Components/Story.js'

function App() {
  const [minigame, setMinigame] = useState(true);
  const [playerTurn, setPlayerTurn] = useState(true);
  const [attackPower, setAttackpower] = useState(0);
  const [result, setResult] = useState({});
  const [convId, setConvId] = useState("");
  function incrPower(x) {
    setAttackpower(attackPower + x);
  }

  //send data/receive data
  function send(conversation_id, userAction, name, scenario) {
    console.log("SENDING THE USER INPUT UWUW: "+userAction)
    console.log("into the conversation: "+conversation_id)
    let url = "http://localhost:8080/submit";
    const response = fetch(url, {
      method: "POST", // *GET, POST, PUT, DELETE, etc.
      headers: {
        "Content-Type": "application/json",
        // 'Content-Type': 'application/x-www-form-urlencoded',
      },
      
      body: JSON.stringify({
        conversation_id: conversation_id,
        action: userAction,
        name: name,
        scenario: scenario
      }),
    }).then(data => data.json()).then(json => {
      //update state
      setResult(json) //save all data
      setConvId(json.conversation_id) //save the conversation id
      console.log(json);

      //open the story panel
      document.getElementById("popup").className = "overlay"
    });
  }


  return (
    <div className="App">
      <div>
        <Menu tab={1} send={(userAction)=>{
          send(convId, userAction, "Arky", "The player is fighting against Kafka, a mind controlling fugitive who knows the secrets of the universe.")
          }} result = {result}/>
        <h1>AIdventure</h1>

        {minigame ? (<ClickerGame handleClick = {()=>{incrPower(0.5)}}/>) : (<></>)}
  
      </div>
      <Story story = {result.outcome}/>
    </div>
  );
}


export default App;
