import './App.css';
import Inventory from './Components/Inventory';

function App() {
  return (
    <div className="App">
      <div>
        <h1>testing</h1>
        <Inventory items={[1,2,3,4]} />
      </div>
    </div>
  );
}

export default App;
