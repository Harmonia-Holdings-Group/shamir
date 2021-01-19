import React from 'react';

// Components
import Encrypt from './components/encrypt/Encrypt';
import Decrypt from './components/decrypt/Decrypt';


// Styles
import './css/nav.scss';


class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      usage: 'encrypt',
    };
  }

  render () {
    return (
      <div className='root-container'>
        <nav>
          <div>
            <p
              onClick={() => {this.setState({usage: 'encrypt'})}}
              className={this.state.usage === 'encrypt'? 'active' : ''}>
                Encrypt
            </p>
          </div>
          <div>
            <p
              onClick={() => {this.setState({usage: 'decrypt'})}}
              className={this.state.usage === 'decrypt'? 'active' : ''}>
                Decrypt
            </p>
          </div>
        </nav>

        <Encrypt show={this.state.usage === 'encrypt'}/>
        <Decrypt show={this.state.usage === 'decrypt'}/>
      </div>
    )
  }
}
export default App;
