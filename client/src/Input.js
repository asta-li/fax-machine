import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';

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
    console.log('Constructing FileSelector');
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
      <React.Fragment> 
        <Button
          variant="contained"
          component="label"
        >
          Select PDF
          <input
            type="file"
            accept='.pdf,application/pdf'
            style={{ display: "none" }}
            onChange={(event) => this.handleFileSelection(event)} 
          />
        </Button>
        {this.state.selectedFileStatus}
      </React.Fragment>
    );
  }
}

// Controls fax number input.
class FaxNumberInput extends React.Component {
  constructor(props) {
    super(props);
    console.log('Constructing FileNumberInput');
    this.state = {
      faxNumberStatus: "",
    };
  }
  
  // Update and validate the input fax number.
  handleInput(event) {
    // TODO: Read fax number from event and validate.
    const faxNumber = "12345";
    const faxNumberStatus = "Okay!"
    this.setState({
      faxNumberStatus: faxNumberStatus,
    });
    this.props.setFaxNumber(faxNumber);
  }

  // Render the element that controls file seletion. 
  render() {
    return (
      <React.Fragment>
        <TextField
          variant="outlined"
          margin="normal"
          required
          fullWidth
          id="fax"
          label="Fax Number"
          name="fax"
          autoComplete="fax"
          autoFocus
          onClick={(event) => this.handleInput(event)}
        />
        {this.state.faxNumberStatus}
      </React.Fragment>
    );
  }
}

export {FileSelector, FaxNumberInput};
