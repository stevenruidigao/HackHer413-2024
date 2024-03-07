import Item from './Item.js';

export default function Inventory({items}) {
    return (
        <div className="inventory right">
            <h2>ITEMS</h2>
            <div className="inventory-container" >
                {items.map((item) => {
                    return (<Item name={item.name} title={item.description}/>);
                })}
            </div>
        </div>
    );
}
