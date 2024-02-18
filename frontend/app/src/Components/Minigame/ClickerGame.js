import './Minigame.css'

import target from './../../assets/target.png';

function ClickerGame(props) {

    return (
        <div>
            <h1>Click the target to deal damage!!</h1>
            <div className="target-outer">
                <div className="target-container">
                <img src={target} className="target" onClick = {props.handleClick} alt=""></img>
                </div>
            </div>

        </div>
    )


}

export default ClickerGame;
