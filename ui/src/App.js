import React, { useState, useEffect } from 'react';
import './App.css';
import PageIndex from './PageIndex'
import axios from 'axios'

function Server(props) {
  return (
    <div className="Server">
      <span className="ServerName">{props.name}</span>
      <span className="ServerCriticals">{props.criticals}</span>
      <span className="ServerWarnings">{props.warnings}</span>
    </div>
  );
}

function App() {
  const [orderby] = useState("Name");
  const [index, setIndex] = useState(1);
  const [servers, setServers] = useState([]);

  useEffect(() => {
    const interval = setInterval(() => {
      axios.get('http://10.93.81.79/api/v1/servers', {
        params: {
          start: (index - 1) * 1000,
          count: 1000,
          orderby: orderby,
        }
      }).then(response => {
        setServers(response.data.Member)
      }).catch(err => console.log(err));
    }, 1000);
    return () => clearInterval(interval);
  }, [index, orderby]);

  return (
    <div className="App">
      <div className="ServerList">
        {servers.map(server => <Server key={server.Name} name={server.Name} criticals={server.Criticals} warnings={server.Warnings} />)}
      </div>
      <div className="Pages">
        {[...Array(50).keys()].map(i => <PageIndex className="PageIndex" key={i} index={index} value={i + 1} onClick={() => {
          setIndex(i + 1);
          axios.get('http://10.93.81.79/api/v1/servers', {
            params: {
              start: i * 1000,
              count: 1000,
              orderby: orderby,
            }
          }).then(response => {
            setServers(response.data.Member)
          }).catch(err => console.log(err));
        }} />)}
      </div>
    </div>
  );
}

export default App;
