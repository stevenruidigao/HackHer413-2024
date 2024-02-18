

export default function Entity(props) {
  return (
    <div className="entity">
      <h1>df</h1>
      <div className = "inventory-container" >
      {items.map((item) => {
        return (<Item name = {item.name} title = {item.description}/>);
      })}
      </div>
    </div>
  );
}