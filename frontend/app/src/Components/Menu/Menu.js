import React, { useState } from 'react';

import './Menu.css'

import Act from './Act.js'

function Menu() {
  const [tab, setTab] = useState(0);
  // let tab, setTab = (0, () => {});

  return (
    <div className = "menu">
      <div className = "sidebar">
      <button className={tab == 0 ? "tab selected" : "btn tab "} onClick={() => {setTab(0)}}>ACT</button>
      <button className={tab == 1 ? "tab selected" : "btn tab "} onClick={() => {setTab(1)}}>ITEM</button>
      <button className={tab == 2 ? "tab selected" : "btn tab "} onClick={() => {setTab(2)}}>STAT</button>
      </div>
        {
          tab == 0 ? (<Act />) :
            tab == 1 ? (<h2>poop</h2>) :
              (<h3>uwu</h3>)
        }

    </div>

  );
}

export default Menu;
