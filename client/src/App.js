import axios from 'axios'; 
import React from 'react';
import Avatar from '@material-ui/core/Avatar';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import Link from '@material-ui/core/Link';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Container from '@material-ui/core/Container';
import { withStyles } from '@material-ui/styles';
import { makeStyles } from '@material-ui/core/styles';

import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Paper from '@material-ui/core/Paper';

import { FileSelector, FaxNumberInput } from './Input.js';
import { FileFaxer } from './Submit.js';
import { ReactComponent as Logo } from './logo.svg';

const styles = theme => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  logo: {
    margin: theme.spacing(1),
    width: theme.spacing(9),
    height: theme.spacing(9),
  },
  form: {
    width: '100%',
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
});

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
      <CssBaseline />
      <div className={classes.paper}>
        <Logo className={classes.logo}/>
        <Typography component="h1" variant="h4" gutterBottom>
          I am a fax machine.
        </Typography>
        <form className={classes.form} noValidate>
          {/* Controls file selection and validation. This component allows a user to select a file,
              validates the file, and updates the file information in the app state. */}
          <FileSelector
            setSelectedFile={this.setSelectedFile}
          />
          {/* Controls fax number input. */}
          <FaxNumberInput
            setFaxNumber={this.setFaxNumber}
          />
          {/* Controls file upload and faxing. */}
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

export default withStyles(styles, {withTheme: true})(FaxMachineApp);
