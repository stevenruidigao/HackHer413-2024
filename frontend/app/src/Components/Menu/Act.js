import './Menu.css'

function Act(props) {
  
  return (
    <div className = "act">
      <h1>ACT - What will you do?</h1>
      <textarea className = "text-input"></textarea>
      <div align="right">
        <button className = "btn submit">{">"}</button>
      </div>
    </div>

  );
}

export default Act;
