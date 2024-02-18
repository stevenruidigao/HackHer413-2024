import './App.css';

import Inventory from './Components/Menu/Inventory.js';
import Menu from './Components/Menu/Menu.js'

function App() {
  return (
    <div className="App">
      <div>
        <Menu tab={1}/>
        <h1>AIdventure</h1>
      </div>
    </div>
  );
}

export default App;
