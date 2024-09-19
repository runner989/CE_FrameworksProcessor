import './style.css';
import './app.css';

import logo from './assets/images/CFLogo.png';
// import {Greet} from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
`;
document.getElementById('logo').src = logo;

