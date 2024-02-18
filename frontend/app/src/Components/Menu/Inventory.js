import Item from "./Item.js";

export default function Inventory({items}) {
  return (
    <div className="inventory">
      <h1>ITEMS</h1>
      <div className = "inventory-container">
      {items.map((item) => {
        return (<Item name = {item}/>);
      })}
      </div>
    </div>
  );
}