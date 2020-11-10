import logo from './logo.svg';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          I am a fax machine.
        </p>
        <a
          className="App-link"
          href="https://github.com/asta-li/fax-machine"
          target="_blank"
          rel="noopener noreferrer"
        >
          Code
        </a>
      </header>
    </div>
  );
}

export default App;
