import './Story.css';
import './Menu/Menu.css';

function Story(props) {
    return (
        <div className="story">
            <button class="popup-btn btn" id="open-popup" onClick={()=>{document.getElementById('popup').className = 'overlay'}}>Show Current Story</button>
            <div id="popup" class="overlay hidden">
                <div class="popup-content">
                    <div id="turns">
                        {props.history.map((turn) => {
                            return (
                                <div className="turn">
                                    <i>{turn.action}</i>
                                    <p>{turn.outcome}</p>
                                </div>
                            );
                        })}
                    </div>
                    { /* <p>{props.story}</p> */ }
                    <button class="popup-btn btn" onClick = {()=>{
                        document.getElementById('popup').className = 'overlay hidden';

                        if (props.isOver) {
                            console.log('GG');
                            document.getElementById('ending').className = 'overlay';
                        }
                    }}>Back</button>
                </div>
            </div>
        </div>
    );
}

export default Story;
