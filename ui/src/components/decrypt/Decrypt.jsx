import React from 'react';

// Styles
import './styles.scss';

class Decrypt extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      keysFileName: "Select keys file",
      encryptedFileName: "Select encrypted file",
      keys: [],
      showParsedKeys: false,
      encryptedB64Content: "",
      error: "",
      showResult: false,
      showAllDecrypted: false,
      plainContent: new Blob(),
      derivedKey: "",
      decryptedFileName: "",
    };

    this.handleDecryptRequest = this.handleDecryptRequest.bind(this);

    this.handleKeysFileSelection = this.handleKeysFileSelection.bind(this);
    this.handleEncryptedFileSelection = this.handleEncryptedFileSelection.bind(this);
    this.keysFileInput = React.createRef();
    this.encryptedFileInput = React.createRef();
  }

  handleKeysFileSelection(e) {
    if (this.keysFileInput.current.files.length !== 1) {
      this.setState({error: "you must select exactly 1 file to encrypt!"});
      return;
    }
    const fileName = this.keysFileInput.current.files[0].name;
    this.setState({
      keysFileName: fileName
    })

    const keysReader = new FileReader();
    keysReader.onload = function(e) {
      this.setState({
        keys: e.target.result.split('\n').filter((k) => k.length > 0),
        showParsedKeys: true,
      })
    }
    keysReader.onerror = function(e) {
      this.setState({
        error: "Failed reading keys file </3",
      })
    }
    keysReader.onload = keysReader.onload.bind(this);
    keysReader.onerror = keysReader.onerror.bind(this);
    keysReader.readAsText(this.keysFileInput.current.files[0], "UTF-8")
  }

  handleEncryptedFileSelection(e) {
    if (this.encryptedFileInput.current.files.length !== 1) {
      this.setState({error: "you must select exactly 1 file to encrypt!"});
      return;
    }
    const fileName = this.encryptedFileInput.current.files[0].name;
    this.setState({
      encryptedFileName: fileName
    })

    const reader = new FileReader();
    reader.onload = function(e) {
      this.setState({
        encryptedB64Content: e.target.result,
      })
    }
    reader.onerror = function(e) {
      this.setState({
        error: "Failed reading encrypted file </3",
      })
    }
    reader.onload = reader.onload.bind(this);
    reader.onerror = reader.onerror.bind(this);
    reader.readAsText(this.encryptedFileInput.current.files[0], "UTF-8")
  }

  handleDecryptRequest(e) {
    if (this.state.keys.length < 2) {
      this.setState({
        error: "obtained less than 2 keys from keys input",
      })
      return
    }
    if (this.state.encryptedB64Content === "") {
      this.setState({
        error: "please select an encrypted file"
      })
      return
    }

    const key = global.GoGetKeyFromKeyShares(this.state.keys);
    if (typeof(key) === 'string' && key.startsWith('ERROR')) {
      this.setState({
        error: key,
      });
      return
    }

    const content = global.GoDecrypt(key, this.state.encryptedB64Content);
    if (typeof(key) === 'string' && key.startsWith('ERROR')) {
      this.setState({
        error: key,
      });
      return
    }

    const byteContent = atob(content);
    const contentData = new Uint8Array(byteContent.length);
    for (var i = 0; i < byteContent.length; i++) {
      contentData[i] = byteContent.charCodeAt(i);
    }

    if (this.state.encryptedFileName.endsWith(".aes")) {
      this.setState({
        decryptedFileName: this.state.encryptedFileName.substr(0, this.state.encryptedFileName.length-4)
      })
    } else {
      this.setState({
        decryptedFileName: this.state.encryptedFileName,
      })
    }

    this.setState({
      error: "",
      derivedKey: key,
      showResult: true,
      plainContent: new Blob(contentData),
    })
  }

  render () {
    return this.props.show ? (
      <div id="encryption-section">
        <div className="container" id="decryption-input">
          <div className="input-group">
            <p className="accent">Keys</p>
            <label>
              A text file containing the decryption keys, each one on a new line.
            </label>
            <div className="file-input">
              <input
                type="file" id="file-input"
                ref={this.keysFileInput}
                onChange={this.handleKeysFileSelection}
              />
              <label htmlFor="file-input">
                <i className="fas fa-file"></i>
                <span>{ this.state.keysFileName} </span>
              </label>
            </div>
            <p className="encoded" style={{display: !this.state.showParsedKeys ? 'none' : ''}}>
              {this.state.keys.map((k) => 
                <span key={k}>{k}</span>
              )}
            </p>
          </div>
          <div className="input-group">
            <p className="accent">Encrypted file</p>
            <label>
              Encrytped file
            </label>
            <div className="file-input">
              <input
                type="file" id="efile-input"
                ref={this.encryptedFileInput}
                onChange={this.handleEncryptedFileSelection}
              />
              <label htmlFor="efile-input">
                <i className="fas fa-file"></i>
                <span>{ this.state.encryptedFileName} </span>
              </label>
            </div>
          </div>
          <div className="button">
            <button onClick={this.handleDecryptRequest}>Decrypt <i className="fas fa-play"></i></button>
          </div>
          <p style={{display: this.state.error !== "" ? 'block' : 'none'}} className='error-message'>
            <span>Error</span>: {this.state.error}
          </p>
        </div>
        <section className="container" id="encrypted-file">
          <p className="subsection-title">Decrypted file <i className="fas fa-copy"></i></p>
          <p className="small" style={{display: !this.state.showResult ? 'none' : ''}}>
            Derived master key:
            <span>{ this.state.derivedKey }</span>
          </p>
          <p style={{display: !this.state.showResult ? 'none' : ''}}>
            <a
              href={`data:application/octet-stream,${this.state.plainContent}`}
              download={this.state.decryptedFileName}
            >
              Save file <i className="fas fa-file-download"></i>
            </a>
          </p>
        </section>
      </div>
    ): null;
  }
}

export default Decrypt;