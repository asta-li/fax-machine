import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
} from 'react-router-dom';

import FaxMachineApp from './App';

ReactDOM.render(
    <React.StrictMode>
        <Router>
            <FaxMachineApp />
        </Router>
    </React.StrictMode>,
    document.getElementById('root')
);
