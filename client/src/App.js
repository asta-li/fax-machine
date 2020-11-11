import axios from 'axios'; 
import React from 'react';

import logo from './logo.svg';
import './App.css';

// TODO(asta): Perform additional client-side validation,
// such as checking for JavaScript in the file.
//
// Perform basic file validation. Returns a pair {fileIsValid, status}.
// If validation is successful then fileIsValid is true and the status is the file name.
// Otherwise, fileIsValid is false and the status contains an error message.
function validateFile(file) {
  let fileIsValid = false;
  let status = 'Error';

  if (!file) {
    fileIsValid = false;
    status = 'Error: Please select a PDF file';
    return {fileIsValid, status};
  }

  if (file.type !== 'application/pdf') {
    fileIsValid = false;
    status = 'Error: Selected file must be a PDF';
    return {fileIsValid, status};
  }
 
  const fileSizeMB = file.size / 1024 / 1024;
  const MAX_SIZE_MB = 5;
  if (fileSizeMB > MAX_SIZE_MB) {
    fileIsValid = false;
    status = 'Error: Selected file is ' + fileSizeMB + 'MB but max is ' + MAX_SIZE_MB + 'MB';
    return {fileIsValid, status};
  }
  
  fileIsValid = true;
  status = file.name;
  return {fileIsValid, status};
}

// Controls selection of a local file.
class FileSelector extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedFileStatus: 'Select a file',
    };
  }
  
  // Update and validate the selected file.
  handleFileSelection(event) {
    if (event.target.files.length === 0) {
        return;
    }
    const selectedFile = event.target.files[0];
    const {fileIsValid, status: selectedFileStatus} = validateFile(selectedFile);
      
    this.setState({
      selectedFileStatus: selectedFileStatus,
    });
    
    if (fileIsValid) {
      // Update the state with the selected file only after validation.
      this.props.setSelectedFile(selectedFile);
    }
  }

  // Render the element that controls file seletion. 
  render() {
    return (
      <div className='App-body'>
        <label className='select-file'>
          <input type='file' onChange={(event) => this.handleFileSelection(event)} /> 
          Select PDF
        </label>
        {this.state.selectedFileStatus}
      </div>
    );
  }
}

// Controls faxing a selected file.
class FileFaxer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      faxFileStatus: '',
    };
    this.handleFileFax = this.handleFileFax.bind(this);
  }

  handleFileFax() {
    if (!this.props.selectedFile) {
      this.setState({
        faxFileStatus: 'Select a file to fax',
      });
    } else {
      console.log('Faxing', this.props.selectedFile);
      this.setState({
        faxFileStatus: 'Faxing...',
      });
  
      // Create form containing the file data.
      const formData = new FormData(); 
      formData.append( 
        'file', 
        this.props.selectedFile,
      ); 

      const config = {     
        headers: { 'content-type': 'multipart/form-data' }
      }

      // Sends the file to the backend for payment processing, upload, and faxing.
      axios.post('/fax', formData, config)
        .then((response) => {
          console.log('Received successful fax response', response);
          
          this.setState({
            faxFileStatus: 'Successfully faxed for $' + response.data.Price + '!',
          });
        })
        .catch((error) => {
          console.log(error);
          this.setState({
            faxFileStatus: error,
          });
        });
    }
  }

  // Render the element that controls faxing the selected file. 
  render() {
    return (
      <div className='App-body'>
        <button className='fax-file' onClick={this.handleFileFax}>
          Fax me!
        </button>
        {this.state.faxFileStatus}
      </div>
    );
  }
}

class FaxMachineApp extends React.Component {
  constructor(props) {
    super(props);
    this.state = { 
      selectedFile: null,
    };
    
    this.setSelectedFile = this.setSelectedFile.bind(this);
  }

  // Sets the selected file.
  // We pass this callback to FileSelector in order maintain file state at the top level.
  setSelectedFile(selectedFile) {
    this.setState({
      selectedFile: selectedFile,
    }); 
  } 

  render() {
    return (
      <div className='App'>
        {/* App header. This content is static and does not change. */}
        <header className='App-header'>
          <img src={logo} className='App-logo' alt='logo' />
          <p>I am a fax machine.</p>
        </header>
        {/* Controls file selection and validation. This component allows a user to select a file,
            validates the file, and updates the file information in the app state. */}
        <FileSelector
          setSelectedFile={this.setSelectedFile}
        />
        {/* Controls file upload and faxing. */}
        <FileFaxer
          selectedFile={this.state.selectedFile}
        />
        {/* App footer. This content is static and does not change. */}
        <footer>
          <p>
          <a
            className='App-link'
            href='https://github.com/asta-li/fax-machine'
            target='_blank'
            rel='noopener noreferrer'
          >
            Code
          </a>
          </p>
        </footer>
      </div>
    );
  }
}

export default FaxMachineApp;
