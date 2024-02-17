import './Menu.css'

function Act(props) {
  
  return (
    <div className="act">
      <h1>ACT - What will you do?</h1>
      <label for="text-input" hidden>Write your action here</label>
      <textarea className="text-input" placeholder="What do you do?"></textarea>
      <div align="right">
        <label for="submit-action" hidden>Submit</label>
        <button id="submit-action" className="btn submit">{">"}</button>
      </div>
    </div>

  );
}

export default Act;
