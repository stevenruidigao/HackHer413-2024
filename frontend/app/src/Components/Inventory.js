import Item from "./Item";

export default function Inventory({items}) {
  return (
    <div className="inventory">
      {items.map((item) => {
        return Item(item);
      })}
    </div>
  );
}