import React from 'react';

// Styles
import './styles.scss';

class Decrypt extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      keysFileName: "Select keys file",
      encryptedFileName: "Select encrypted file",
      keysTextValue: "",
      encryptedB64Content: "",
      error: "",
      showKeysTextArea: false,
      showEncFileTextArea: false,
    };

    this.handleInputChange = this.handleInputChange.bind(this);

    this.handleKeysFileSelection = this.handleKeysFileSelection.bind(this);
    this.handleEncryptedFileSelection = this.handleEncryptedFileSelection.bind(this);
    this.keysFileInput = React.createRef();
    this.encryptedFileInput = React.createRef();
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.value;
    const name = target.name;
    this.setState({
      [name]: value,
    });
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

    var reader = new FileReader();
    reader.onload = function(e) {
      this.setState({
        keysTextValue: e.target.result,
      })
    }
    reader.onerror = function(e) {
      this.setState({
        error: "Failed reading keys file </3",
      })
    }
    reader.onload = reader.onload.bind(this);
    reader.onerror = reader.onerror.bind(this);
    reader.readAsText(this.keysFileInput.current.files[0], "UTF-8")
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
            <label>
              <b>OR</b> Paste the keys 
              <button
                onClick={(e) => {this.setState({showKeysTextArea: !this.state.showKeysTextArea})}}
              >
                here
              </button>.
            </label>
            <textarea
              style={{display: this.state.showKeysTextArea ? 'block' : 'none'}}
              name="keysTextValue"
              id=""
              value={this.state.keysTextValue}
              onChange={this.handleInputChange}
            />
          </div>
          <div className="input-group">
            <p className="accent">Encrypted file</p>
            <label>
              Encrytped file
            </label>
            <div className="file-input">
              <input
                type="file" id="file-input"
                ref={this.encryptedFileInput}
                onChange={this.handleEncryptedFileSelection}
              />
              <label htmlFor="file-input">
                <i className="fas fa-file"></i>
                <span>{ this.state.encryptedFileName} </span>
              </label>
            </div>
            <label>
              <b>OR</b> Paste the encrypted base 64 encoded content 
              <button
                onClick={(e) => {this.setState({showEncFileTextArea: !this.state.showEncFileTextArea})}}
              >
                here
              </button>.
            </label>
            <textarea
              style={{display: this.state.showEncFileTextArea ? 'block' : 'none'}}
              name="encryptedB64Content"
              value={this.state.keysTextValue}
              onChange={this.handleInputChange}
            />
          </div>
          <p style={{display: this.state.error !== "" ? 'block' : 'none'}} className='error-message'>
            <span>Error</span>: {this.state.error}
          </p>
        </div>
      </div>
    ): null;
  }
}

export default Decrypt;