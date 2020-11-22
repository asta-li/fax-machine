import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import FormHelperText from '@material-ui/core/FormHelperText';

// Controls selection of a local file.
const FileInput = props => {
  // Update and validate the selected file.
  const handleFileSelection = event => {
    if (event.target.files.length === 0) {
      return;
    }
    const selectedFile = event.target.files[0];
    const { error : selectedFileError }  = validateFile(selectedFile);
    if (selectedFileError) {
      props.setSelectedFileError(selectedFileError);
      props.setSelectedFile(null);
    } else {
      // Update the state with the selected file only after it passes validation.
      props.setSelectedFile(selectedFile);
    }
  }
    
  let status = 'Please select a file.';
  let error = false;
  if (props.selectedFileError) {
    error = true;
    status = props.selectedFileError;
  } else if (props.selectedFile) {
    status = props.selectedFile.name; 
  }
  
  // Render the element that controls file seletion.
  return (
    <React.Fragment>
      <Button variant="contained" component="label" required fullWidth>
        Select PDF
        <input
          type="file"
          accept=".pdf,application/pdf"
          hidden
          onChange={handleFileSelection}
        />
      </Button>
      <FormHelperText error={error}>{status}</FormHelperText>
    </React.Fragment>
  );
}

export { FileInput };


// File validation helper functions.
// ===================================================================

// TODO(asta): Perform additional client-side validation,
// such as checking for JavaScript in the file.
//
// Perform basic file validation. Return true if validation is successful.
function validateFile(file) {
  let error = '';

  if (!file) {
    error = 'Please select a PDF file';
    return { error };
  }

  if (file.type !== 'application/pdf') {
    error = 'Selected file must be a PDF';
    return { error };
  }

  const fileSizeMB = file.size / 1024 / 1024;
  const MAX_SIZE_MB = 5;
  if (fileSizeMB > MAX_SIZE_MB) {
    error = 'Selected file is ' + fileSizeMB + 'MB but max is ' + MAX_SIZE_MB + 'MB';
    return { error };
  }

  return { error };
}

