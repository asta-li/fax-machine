import axios from 'axios'; 
import React from 'react';

import logo from './logo.svg';
import './App.css';


class FileSelector extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedFile: null,
      selectedFileStatus: 'Select a file',
      uploadedFileStatus: '',
    };
    this.handleFileUpload = this.handleFileUpload.bind(this);
  }

  handleFileSelection(event) {  
    const selectedFile = event.target.files[0];
    this.setState({
      selectedFile: selectedFile,
      selectedFileStatus: selectedFile.name,
    });
  }

  renderSelectFile() {
    return (
      <div>
        <label className='Select'>
          <input type='file' onChange={(event) => this.handleFileSelection(event)} /> 
          Select PDF
        </label>
        {this.state.selectedFileStatus}
      </div>
    );
  }

  handleFileUpload() {
    console.log('Uploading', this.state.selectedFile);

    if (!this.state.selectedFile) {
      this.setState({
        uploadedFileStatus: 'Please select a file to upload',
      });
    } else {
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
        uploadedFileStatus: 'Uploaded!',
      });
      this.props.setUploadedFile(this.state.selectedFile);
    }
  }

  renderUploadFile() {
    return (
      <div>
        <button className='Upload' onClick={this.handleFileUpload}>
          Upload
        </button>
        {this.state.uploadedFileStatus}
      </div>
    );
  }

  render() {
    return (
      <div className='App-body'>
        {this.renderSelectFile()}
        {this.renderUploadFile()}
      </div>
    );
  }
}

class FileFaxer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      faxFileStatus: '',
    };
    this.handleFileFax = this.handleFileFax.bind(this);
  }

  handleFileFax() {
    console.log('Faxing', this.props.uploadedFile);

//    if (!this.props.uploadedFile) {
//      this.setState({
//        faxFileStatus: 'Upload a file to fax',
//      });
//    } else {
//      const formData = new FormData(); 
//      formData.append( 
//        'fileToFax', 
//        this.state.selectedFile, 
//        this.state.selectedFile.name 
//      ); 
//       
//      // Sends the file to the backend.
//      // TODO(asta): Handle this request.
//      axios.post('/upload', formData);
  
    this.setState({
      faxFileStatus: 'Successfully faxed!',
    });
  }

  renderFaxFile() {
    return (
      <div>
        <button className='Upload' onClick={this.handleFileFax}>
          Fax me!
        </button>
        {this.state.faxFileStatus}
      </div>
    );
  }

  render() {
    return (
      <div className='App-body'>
        {this.renderFaxFile()}
      </div>
    );
  }
}

class FaxMachineApp extends React.Component {
  constructor(props) {
    super(props);
    this.state = { 
      uploadedFile: null,
    };
    
    this.setUploadedFile = this.setUploadedFile.bind(this);
  }

  // Handles file selection.
  setUploadedFile(uploadedFile) {
    this.setState({
      uploadedFile: uploadedFile,
    }); 
  } 

  render() {
    return (
      <div className='App'>
        <header className='App-header'>
          <img src={logo} className='App-logo' alt='logo' />
          <p>I am a fax machine.</p>
        </header>
        <FileSelector
          setUploadedFile={this.setUploadedFile}
        />
        <FileFaxer
          uploadedFile={this.uploadedFile}
        />
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
