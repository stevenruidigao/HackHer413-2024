import './Story.css'
import './Menu/Menu.css'

function Story(props) {
  
  return (
    <div className="story">
      <button class="popup-btn btn" id="open-popup" onClick = {()=>{document.getElementById("popup").className = "overlay"}}>Show Current Story</button>

<div id="popup" class="overlay hidden">
  <div class="popup-content">
    <p>{props.story}</p>
    <button class="popup-btn btn" onClick = {()=>{document.getElementById("popup").className = "overlay hidden"}}>Back</button>
  </div>
</div>


    </div>

  );
}

export default Story;
