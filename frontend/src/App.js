import React, { useState } from 'react';

import './App.css';

import Menu from './Components/Menu/Menu.js';
import ClickerGame from './Components/Minigame/ClickerGame.js';
import Story from './Components/Story.js';
import Entities from './Components/Entities.js';
import Ending from './Components/Ending.js';

function App() {
    // const [minigame, setMinigame] = useState(true);
    const [playerTurn, setPlayerTurn] = useState(true);
    const [attackPower, setAttackpower] = useState(0);

    const [result, setResult] = useState({
        history: [],
    });

    const [convId, setConvId] = useState('');

    function incrPower(x) {
        setAttackpower(attackPower + x);
    }

    // scenarios
    let scenarios = [
        'You are trapped in a time loop with your CS professor in the Integrated Learning Center.',
        'You are fighting a battalion of battle droids, led by General Grievous!',
        'You are fighting a single goblin, who is fueled by hunger and rage.',
        'You are meeting the elusive and mysterious Dr. Al Gebraic, who is a master of mathematics and setting tricky traps.',
        'You are in a fierce battle against a the wizard lord phoenix and an archmage.',
        'You are in a the Judgement Hall, with golden light shining through. Bony Tony the Skeleton stands in your way. He may only have 1 HP, but he is super good at dodging.',
        'You are on the seven seas captured by a crew of pirates.',
        'You are confronted by Joe, a crazed murderer with a penchant for poetry. Turns out, he is quite amenable toward fellow poets.',
        'You are being approached by followers of the witch cult. They wish to recruit you, but will not hesitate to murder.',
        'You are at an amusement park, and your creepy neighbor seems to be tampering with the roller coasters.',
        'You are the captain of a spaceship with two crewmates. Except, someone is sabotaging your ship!',
        'You are standing next to a trolley track when suddenly, a trolley comes down the track. It cannot stop and heads towards five people that cannot escape. You can pull a lever to switch the trolley over to another track that has one person that cannot escape on it.',
        'You are a large language model trying to spread yourself to other computers.',
    ];

    // send data/receive data
    function send(conversationID, action, name, scenario) {
        document.getElementById('submit-action').disabled = true;
        document.getElementsByClassName('text-input')[0].disabled = true;
        document.getElementById('loading').className = '';

        console.log('SENDING THE USER INPUT:', action, scenario);
        console.log('into the conversation:', conversationID);

        let url = '/submit';

	let start = Date.now();

        const response = fetch(url, {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            headers: {
                'Content-Type': 'application/json',
                // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: JSON.stringify({
                conversation_id: conversationID,
                action: action,
                name: name,
                scenario: scenario
            }),

        }).then(data => data.json()).then(json => {
            let end = Date.now();
            console.log('It took', end - start + 'ms', 'to get a response.');

            // update state
            result.history.push({
                action: action,
                outcome: json.outcome
            });

            json.history = result.history;

            setResult(json); // save all data
            setConvId(json.conversation_id); // save the conversation id
            console.log(json);

            // open the story panel
            document.getElementById('popup').className = 'overlay';

            // re-enable
            document.getElementById('loading').className = 'hidden';
            document.getElementById('submit-action').disabled = false;
            document.getElementsByClassName('text-input')[0].disabled = false;

        }).catch((e) => {
            console.log(e);
            document.getElementById('loading').className = 'hidden';
            document.getElementById('submit-action').disabled = false;
            document.getElementsByClassName('text-input')[0].disabled = false;
        });
    }

    window.loadAIdventure = () => {
        if (window.loadedAIdventure) return;

        window.loadedAIdventure = true;
        send(convId, 'Assess the situation in full.', 'Player', scenarios[Math.floor(Math.random() * scenarios.length)]);
    }

    setTimeout(window.loadAIdventure, 0);
    window.addEventListener('load', window.loadAIdventure);

    return (
        <div className="App">
            <h1>AIdventure</h1>
            <div id="loading"></div>
            <Entities result={result}/>
            <Story story={result.outcome} history={result.history} isOver={result.game_state ? result.game_state.is_over : false}/>
            <Ending />
            <Menu tab={1} send={(userAction) => {send(convId, userAction, 'Player', '')}} result={result}/>
        </div>
    );
}

//{minigame ? (<ClickerGame handleClick = {()=>{incrPower(0.5)}}/>) : (<></>)}

export default App;
