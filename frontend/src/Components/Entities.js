import React, { useState } from 'react';

import Entity from './Entity.js'
import './Entities.css'


function Menu(props) {
  let result = props.result;
  if (result.game_state) {
    let npcs = result.game_state.npcs;
    let player = result.game_state.player;
    return (
      <div className="entities">
          <div className="entity-left"><Entity name={player.name} hp={player.stats.HP} maxHp={player.stats.MAX_HP} /></div>
          <div className="entity-right">
            {npcs.map((npc) => {
              return (<Entity name={npc.name} hp={npc.stats.HP} maxHp={npc.stats.MAX_HP} />)
            })}
            </div>
      </div>
  
    );
  } else {
    return (
    <div className="entities">
    <div className="entity-left"><Entity name="Player" hp={100} maxHp={100} /></div>
    <div className="entity-right">
      </div>
</div>)
  }
  }
  

export default Menu;

