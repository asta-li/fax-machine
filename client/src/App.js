import axios from 'axios'; 
import React from 'react';

import logo from './logo.svg';
import './App.css';

     
class FaxMachineApp extends React.Component {
  constructor(props) {
    super(props);
    this.state = { 
      selectedFile: null,
      uploadedFile: 'Please select a file',
    };
    
    this.handleFileUpload = this.handleFileUpload.bind(this);
  }

  // Handles file selection.
  handleFileSelection(event) {
    this.setState({
      selectedFile: event.target.files[0],
    }); 
  } 
     
  // Handles file upload.
  handleFileUpload() {
    console.log('Uploading', this.state.selectedFile);

    if (this.state.selectedFile) {
      const formData = new FormData(); 
      formData.append( 
        'fileToFax', 
        this.state.selectedFile, 
        this.state.selectedFile.name 
      ); 
       
      // Sends the file to the backend.
      // TODO(asta): Handle this request.
      axios.post('/upload', formData);
  
      this.setState({
        uploadedFile: 'Uploaded!',
      });
    }
  } 
     
  render() {
    return (
      <div className='App'>
        <header className='App-header'>
          <img src={logo} className='App-logo' alt='logo' />
          <p>I am a fax machine.</p>
        </header>
        <div className='App-body'> 
          <label className='Select'>
            <input type='file' onChange={(event) => this.handleFileSelection} /> 
            Select PDF
          </label>
          <button className='Select' onClick={this.handleFileUpload}>
            Upload now
          </button>
          {this.state.uploadedFile}
        </div> 
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
      </div>
    );
  }
}

export default FaxMachineApp;
