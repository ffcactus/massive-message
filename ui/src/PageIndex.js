import React from 'react';
import './App.css';

function PageIndex(props) {
  if (props.index === props.value) {
    return <span style={{ backgroundColor: 'lightblue' }} className="PageIndex" onClick={props.onClick} >{props.value}</span >
  } else {
    return <span className="PageIndex" onClick={props.onClick}>{props.value}</span>
  }
}

export default PageIndex;