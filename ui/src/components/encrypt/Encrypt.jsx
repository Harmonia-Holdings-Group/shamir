import React from 'react';

// Styles
import './styles.scss';

class Encrypt extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      keyShares: 3,
      keyThreshold: 2,
      encryptionPassword: "",
      error: "",
      fileInputMessage: "Select file to encrypt",
      genKey: "",
      showResult: false,
      cipherContent: "",
      showAllCipher: false,
    };

    this.handleInputChange = this.handleInputChange.bind(this);
    this.handleEncryptRequest = this.handleEncryptRequest.bind(this);
    this.handleFileSelection = this.handleFileSelection.bind(this);

    this.fileInput = React.createRef();
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.value;
    const name = target.name;
    this.setState({
      [name]: value,
    });
  }

  handleEncryptRequest(event) {
    if (this.fileInput.current.files.length !== 1) {
      this.setState({error: "you must select exactly 1 file to encrypt!"});
      return;
    }
    if (this.state.keyShares < 3) {
      this.setState({error: "key shares value must be at least 3"});
      return;
    }
    if (this.state.keyThreshold < 2 || this.state.keyThreshold > this.state.keyShares) {
      this.setState({error: "keys threshold value must be between 2 and the number of key shares."});
      return;
    }
    if (this.state.encryptionPassword === "") {
      this.setState({error: "encryption password must not be empty!"});
      return;
    }
    this.setState({
      error: ""
    })

    const reader = new FileReader();
    reader.onload = function(e) {
      const fileContent = e.target.result;
      const content8 = new Uint8Array(fileContent);

      const wasmOut = global.GoEncrypt(this.state.encryptionPassword, content8);
      if (wasmOut.length !== 2) {
        this.setState({error: `(unexpected) ${wasmOut}`});
        return;
      }

      this.setState({
        genKey: wasmOut[0],
        cipherContent: wasmOut[1],
        showResult: true,
      });
    }
    reader.onload = reader.onload.bind(this);
    reader.readAsArrayBuffer(this.fileInput.current.files[0]);
  }

  handleFileSelection(event) {
    if (this.fileInput.current.files.length !== 1) {
      this.setState({error: "you must select exactly 1 file to encrypt!"});
      return;
    }
    const fileName = this.fileInput.current.files[0].name;
    this.setState({
      fileInputMessage: fileName
    })
  }

  render () {
    return this.props.show ? (
      <div id="encryption-section">
        <section className="container" id="encryption-input">
          <div className="file-input">
            <input
              type="file" id="file-input"
              ref={this.fileInput}
              onChange={this.handleFileSelection}
            />
            <label htmlFor="file-input"><i className="fas fa-file"></i><span>{this.state.fileInputMessage}</span></label>
          </div>
          <div className="input-group">
            <p className="accent">Key shares</p>
            <input
              type="number"
              min="3"
              placeholder="Key shares"
              id="key-shares-input"
              name="keyShares"
              value={this.state.keyShares}
              onChange={this.handleInputChange}
            />
            <label htmlFor="key-shares-input">Number of keys to generate after encryption. Min: 3</label>
          </div>
          <div className="input-group">
            <p className="accent">Keys threshold</p>
            <input
              type="number"
              placeholder="Key threshold"
              id="key-threshold-input"
              name="keyThreshold"
              min="2"
              value={this.state.keyThreshold}
              onChange={this.handleInputChange}
            />
            <label htmlFor="key-threshold-input">
              Minimum number of keys required to decrypt the file, must bet least 2, and at most
              the number of key shares.
            </label>
          </div>
          <div className="input-group">
            <p className="accent">Encryption password</p>
            <input
              type="password"
              placeholder="Enter a master password"
              id="encryption-password"
              name="encryptionPassword"
              value={this.state.encryptionPassword}
              onChange={this.handleInputChange}
            />
            <label htmlFor="encryption-password">Master password to encrypt the file.</label>
          </div>
          <div className="button">
            <button onClick={this.handleEncryptRequest} id="encrypt-button">Encrypt <i className="fas fa-play"></i></button>
          </div>
          <p style={{display: this.state.error !== "" ? 'block' : 'none'}} className='error-message'>
            <span>Error</span>: {this.state.error}
          </p>
        </section>

        <section className="container" id="encryption-keys">
          <p className="subsection-title">Keys <i className="fas fa-copy"></i></p>
          <p className="small" style={{display: !this.state.showResult ? 'none' : ''}}>
            Derived master key:
            <span>{ this.state.genKey }</span>
          </p>
          <div className="keys">
            <p className="encoded" style={{display: this.state.showResult ? 'none' : ''}} >
              Waiting for input...
            </p>
            <ol style={{display: !this.state.showResult ? 'none' : ''}} id="gen-keys-result">
            </ol>
          </div>
        </section>
      
        <section className="container" id="encrypted-file">
          <p className="subsection-title">Encrypted file <i className="fas fa-copy"></i></p>
          <p style={{display: !this.state.showResult ? 'none' : ''}}>
            <a href="data:application/octet-stream,charset=utf-16le;base64,hola">
              Save file <i className="fas fa-file-download"></i>
            </a>
          </p>
          <p className="small">
            Base64 encoded:
          </p>
          <p className="encoded" style={{display: this.state.showResult ? 'none' : ''}} >
              Waiting for input...
          </p>
          <p className="encoded" style={{display: !this.state.showResult ? 'none' : ''}}>
            <p style={{display: this.state.cipherContent.length >= 500 && this.state.showAllCipher ? 'block': 'none'}}>
              <button onClick={() => { this.setState({showAllCipher: false})}} >... hide</button>
            </p>
            { this.state.showAllCipher ? this.state.cipherContent : this.state.cipherContent.substr(0, 500) }
            <p style={{display: this.state.cipherContent.length >= 500 && !this.state.showAllCipher ? 'block': 'none'}}>
              { this.state.cipherContent.substr(0, 500) }
              <button onClick={() => { this.setState({showAllCipher: true})}} >... show all</button>
            </p>
          </p>
        </section>
      </div>
    ): null;
  }
}

export default Encrypt;