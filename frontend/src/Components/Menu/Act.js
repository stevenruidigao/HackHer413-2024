import './Menu.css';

function Act(props) {
    function submit() {
        props.send(document.getElementsByClassName('text-input')[0].value);
        document.getElementsByClassName('text-input')[0].value = '';
    }

    return (
        <div className="act right">
            <h2>ACT - What will you do?</h2>
            <label for="action-input" hidden>Write your action here</label>
            <textarea id="action-input" className="text-input" placeholder="What do you do?" onKeyDown={(e) => {
                    if (e.key === 'Enter') submit();
                }
            }></textarea>
            <div align="right">
                <label for="submit-action" hidden>Submit</label>
                <button id="submit-action" className="btn submit" onClick={submit}>{'>'}</button>
            </div>
        </div>
    );
}

export default Act;
