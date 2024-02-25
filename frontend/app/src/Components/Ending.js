import './Story.css'
import './Menu/Menu.css'

function Ending(props) {
  
  return (

<div id="ending" class="overlay hidden">
  <div class="popup-content">
    <h1>Adventure Over!</h1>
    <p>Thank you for playing this AIdventure!</p>
    <button class="popup-btn btn" onClick={() => {document.getElementById("ending").className = "overlay hidden"}}>Continue?</button>
  </div>
</div>
  );
}

export default Ending;
