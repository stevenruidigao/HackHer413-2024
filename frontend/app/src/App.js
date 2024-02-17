import './App.css';
import Inventory from './Components/Inventory';
import Menu from './Components/Menu/Menu.js'

function App() {
  return (
    <div className="App">
      <div>
        <Menu tab={1}/>
        <h1>testing</h1>
        <Inventory items={[1,2,3,4]} />
      </div>
    </div>
  );
}

export default App;
