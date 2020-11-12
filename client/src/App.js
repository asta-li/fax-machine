import axios from 'axios'; 
import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Link from '@material-ui/core/Link';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Container from '@material-ui/core/Container';
import { withStyles } from '@material-ui/styles';
import { makeStyles } from '@material-ui/core/styles';

// Custom styles.
const styles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}));

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
      formData.append('file', this.props.selectedFile); 
      formData.append('faxNumber', this.props.faxNumber); 

      const config = {     
        headers: { 'content-type': 'multipart/form-data' }
      }

      // Sends the file to the backend for payment processing, upload, and faxing.
      axios.post('/api/fax', formData, config)
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
    console.log(this.props.classes);
    const { classes } = this.props;
    return (
      <React.Fragment>
        <Button
          type="submit"
          fullWidth
          variant="contained"
          color="primary"
          className={classes.submit}
          onClick={this.handleFileFax}
        >
          Fax me!
        </Button>
        {/* TODO(asta): Display status and errors persistently. Use Redux? */}
        {this.state.faxFileStatus}
      </React.Fragment>
    );
  }
}

const StyledFileFaxer = withStyles(styles)(FileFaxer);

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://github.com/asta-li/fax-machine">
        Fax Machine Dev
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

class FaxMachineApp extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = { 
      selectedFile: null,
      faxNumber: '',
    };
    
    this.setSelectedFile = this.setSelectedFile.bind(this);
    this.setFaxNumber = this.setFaxNumber.bind(this);
  }

  // Sets the selected file.
  // We pass this callback to FileSelector in order maintain state at the top level.
  setSelectedFile(selectedFile) {
    this.setState({
      selectedFile: selectedFile,
    }); 
  } 
  
  // Sets the fax number.
  // We pass this callback to FaxNumberInput in order maintain state at the top level.
  setFaxNumber(faxNumber) {
    this.setState({
      faxNumber: faxNumber,
    }); 
  } 

  render() {
    const { classes } = this.props;
    return (
      <Container component="main" maxWidth="xs">
        <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap" />
        <CssBaseline />
        <div className={classes.paper}>
          <Typography component="h1" variant="h5">
            Sign in
          </Typography>
          <form className={classes.form} noValidate>
            {/* Controls fax number input. */}
            <FaxNumberInput
              setFaxNumber={this.setFaxNumber}
            />
            {/* Controls file selection and validation. This component allows a user to select a file,
                validates the file, and updates the file information in the app state. */}
            <FileSelector
              setSelectedFile={this.setSelectedFile}
            />
            {/* Controls file upload and faxing. */}
            <StyledFileFaxer
              selectedFile={this.state.selectedFile}
              faxNumber={this.state.faxNumber}
            />
          </form>
        </div>
        <Box mt={8}>
          <Copyright />
        </Box>
      </Container>
    );
  }
}

FaxMachineApp.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(FaxMachineApp);
