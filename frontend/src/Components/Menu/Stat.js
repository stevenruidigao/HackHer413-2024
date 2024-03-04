import './Menu.css'
import Item from "./Item.js";

export default function Stat(props) {
  let stats = props.stats;
  let skills = props.skills;

  return (
    <div className="stats">
      <div className="stat-title-container"><h2 className = "stat-title">STAT</h2></div>
      <div className = "stat-container">
        {Object.keys(stats).map((stat) => {
          return (<Item name = {`${stat}: Lv. ${stats[stat]}`} title = {stats[stat]}/>);
        })}
      </div>
      <div className="skill-title-container"><h2 className = "skill-title">SKILL</h2></div>
      <div className = "skill-container">
        {skills.map((skill) => {
          return (<Item name = {skill.name} title = {skill.description}/>);
        })}

      </div>
    </div>
  );
}
