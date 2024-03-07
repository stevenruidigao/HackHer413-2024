export default function Item(props) {
    return (
        <div className="item" title={props.title}>
            <h3>{props.name}</h3>
        </div>
    );
}
