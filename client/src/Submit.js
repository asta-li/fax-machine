import axios from 'axios'; 
import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import Link from '@material-ui/core/Link';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Container from '@material-ui/core/Container';


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
    return (
      <React.Fragment>
        {/* TODO(asta): Make this a type="submit" button and correctly route the form */}
        {/*<Button
          type="submit"
          fullWidth
          variant="contained"
          color="primary"
          className={classes.submit}
        >
        </Button>*/}
        <Button
          fullWidth
          variant="contained"
          color="primary"
          onClick={this.handleFileFax}
        >
          Fax me!
        </Button>
  {this.state.faxFileStatus}
      </React.Fragment>
    );
  }
}

export {FileFaxer};
