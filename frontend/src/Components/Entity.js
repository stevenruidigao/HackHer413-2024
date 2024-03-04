
import './Entities.css'

export default function Entity(props) {
  return (
    <div className="entity">
      <h1>{props.name}</h1>
      <div className = "hp-bar" >
      <progress className="health" value={props.hp} max={props.maxHp}></progress>
      </div>
    </div>
  );
}