import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
// import MuiPhoneNumber from "material-ui-phone-number";

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
    const {isValid: fileIsValid, status: selectedFileStatus} = validateFile(selectedFile);
      
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
    this.state = {
      faxNumberStatus: "",
    };
  }
  
  // Update and validate the input fax number.
  handleInput(event) {
    const faxNumber = event.target.value;
    const {isValid: faxNumberIsValid, status: faxNumberStatus} = validateFaxNumber(faxNumber);

    this.setState({
      faxNumberStatus: faxNumberStatus,
    });
    if (faxNumberIsValid) {
      this.props.setFaxNumber("+1" + faxNumber);
    }
  }

  // Render the element that controls file seletion. 
  render() {
    return (
      <div>
        {/* <MuiPhoneNumber onlyCountries={'us'} onChange={this.handleInput}/> */}
        <TextField
          variant="outlined"
          margin="normal"
          required
          /*fullWidth*/
          id="fax"
          label="Fax Number"
          name="fax"
          autoComplete="fax"
          autoFocus
          onClick={(event) => this.handleInput(event)}
        />
        {this.state.faxNumberStatus}
      </div>
    );
  }
}

export {FileSelector, FaxNumberInput};


// Validation helper functions
// ===================================================================

// TODO(asta): Perform additional client-side validation,
// such as checking for JavaScript in the file.
//
// Perform basic file validation. Returns a pair {isValid, status}.
// If validation is successful then isValid is true and the status is the file name.
// Otherwise, isValid is false and the status contains an error message.
function validateFile(file) {
  let isValid = false;
  let status = 'Error';

  if (!file) {
    isValid = false;
    status = 'Error: Please select a PDF file';
    return {isValid, status};
  }

  if (file.type !== 'application/pdf') {
    isValid = false;
    status = 'Error: Selected file must be a PDF';
    return {isValid, status};
  }
 
  const fileSizeMB = file.size / 1024 / 1024;
  const MAX_SIZE_MB = 5;
  if (fileSizeMB > MAX_SIZE_MB) {
    isValid = false;
    status = 'Error: Selected file is ' + fileSizeMB + 'MB but max is ' + MAX_SIZE_MB + 'MB';
    return {isValid, status};
  }
  
  isValid = true;
  status = file.name;
  return {isValid, status};
}

// Validate the given fax number string.
function validateFaxNumber(faxNumber) {
  let isValid = false;
  let status = 'Error';

  if (!faxNumber) {
    isValid = false;
    status = 'Error: Please enter a fax number';
    return {isValid, status};
  }

  if (faxNumber.length !== 10) {
    isValid = false;
    status = 'Error: Fax number must be 10 digits long';
    return {isValid, status};
  }
 
  isValid = true;
  status = "";
  return {isValid, status};
}
