import React, { useState } from 'react';

import './Menu.css'

import Act from './Act.js'
import Inventory from './Inventory.js'

function Menu(props) {
  const [tab, setTab] = useState(0);
  // let tab, setTab = (0, () => {});

  return (
    <div className="menu">
      <div className="sidebar">
        <label for="act" hidden>Action</label>
        <button id="act" className={tab === 0 ? "tab selected" : "btn tab "} onClick={() => {setTab(0)}}>ACT</button>
        <label for="item" hidden>Item</label>
        <button id="item" className={tab === 1 ? "tab selected" : "btn tab "} onClick={() => {setTab(1)}}>ITEM</button>
        <label for="stat" hidden>Item</label>
        <button id="stat" className={tab === 2 ? "tab selected" : "btn tab "} onClick={() => {setTab(2)}}>STAT</button>
      </div>
      <div className="right">
        {
          tab === 0 ? (<Act send={props.send}/>) :
            tab === 1 ? (<Inventory items={[1,"u",1,1,1,1,1]}/>) :
              (<h3>uwu</h3>)
        }
      </div>

    </div>

  );
}

export default Menu;
