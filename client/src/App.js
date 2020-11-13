import axios from 'axios'; 
import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import Link from '@material-ui/core/Link';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Container from '@material-ui/core/Container';
import { withStyles } from '@material-ui/styles';
import { makeStyles } from '@material-ui/core/styles';

import { FileSelector, FaxNumberInput } from './Input.js';

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
    } else if (!this.props.faxNumber) {
      this.setState({
        faxFileStatus: 'Enter a valid fax number',
      });
    } else {
      console.log('Faxing', this.props.selectedFile);
      console.log('Destination', this.props.faxNumber);
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
            faxFileStatus: 'Successfully faxed ' + response.data.FaxId + ' for $' + response.data.Price + '!',
          });
        })
        .catch((error) => {
          console.log(error);
          this.setState({
            faxFileStatus: 'Unable to fax!',
          });
        });
    }
  }

  // Render the element that controls faxing the selected file. 
  render() {
    // const { classes } = this.props;
    return (
      <React.Fragment>
        {/* TODO(asta): Make this a type="submit" button and correctly route the form */}
        <Button
          fullWidth
          variant="contained"
          color="primary"
          /*className={classes.submit}*/
          onClick={this.handleFileFax}
        >
          Fax me!
        </Button>
        {this.state.faxFileStatus}
      </React.Fragment>
    );
  }
}

// TODO(asta): Debug StyledFileFaxer.
//
// FileFaxer.propTypes = {
//   classes: PropTypes.object.isRequired,
// };
// 
// const StyledFileFaxer = withStyles(styles)(FileFaxer);

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
      faxNumber: '+16504344807',
    };
    
    this.setSelectedFile = this.setSelectedFile.bind(this);
    this.setFaxNumber = this.setFaxNumber.bind(this);
  }

  // Sets the selected file.
  // We pass this callback to FileSelector in order maintain state at the top level.
  setSelectedFile(selectedFile) {
    if (selectedFile) {
      this.setState({
        selectedFile: selectedFile,
      }); 
    }
  } 
  
  // Sets the fax number.
  // We pass this callback to FaxNumberInput in order maintain state at the top level.
  setFaxNumber(faxNumber) {
    if (faxNumber) {
      this.setState({
        faxNumber: faxNumber,
      }); 
    }
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
            <form className={classes.form}>
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
            {/*<StyledFileFaxer*/}
            <FileFaxer
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
