import React from 'react';
import PropTypes from 'prop-types';
import CssBaseline from '@material-ui/core/CssBaseline';
import MLink from '@material-ui/core/Link';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Container from '@material-ui/core/Container';
import {withStyles} from '@material-ui/styles';

import {FaxNumberInput, FileSelector} from './Input.js';
import {FileFaxer} from './Submit.js';
import {ReactComponent as Logo} from './logo.svg';
import {Link, Route, Switch, useLocation} from 'react-router-dom';
import SendingFax from "./components/SendingFax";

const styles = (theme) => ({
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

// // TODO(asta): Debug StyledFileFaxer.
// //
// // FileFaxer.propTypes = {
// //   classes: PropTypes.object.isRequired,
// // };
// //
// // const StyledFileFaxer = withStyles(styles)(FileFaxer);
// >>>>>>> 612be14 (paypal_integration)

// A custom hook that builds on useLocation to parse
// the query string for you.
function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function Copyright() {
  return (
      <Typography variant="body2" color="textSecondary" align="center">
        {'Copyright © '}
          Fax Machine Dev
        {new Date().getFullYear()}
        {'.'}
      </Typography>
  );
}

const FaxMachineApp = () => {

  let query = useQuery();

  return <Container component="main" maxWidth="xs">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap" />
    <CssBaseline />
    <div>
      <nav>
        <ul>
          <li>
            <Link to="/">Home</Link>
          </li>
        </ul>
      </nav>

      {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
      <Switch>
        <Route exact path="/">
          <Home action={query.get("action")} transactionId={query.get("transactionId")} token={query.get("token")} payerId={query.get("PayerId")} />
        </Route>

      </Switch>
    </div>
    <Box mt={8}>
      <Copyright />
    </Box>
  </Container>
};


const Home = props => {
  const [selectedFile, setSelectedFile] = React.useState(null);
  const [faxNumber, setFaxNumber] = React.useState('+16504344807');
  const [isUploadSuccess, setUploadSuccess] = React.useState(false);

  const uploadSuccessHandler = (id) => {
    setUploadSuccess(true);
  };

  let showFileFaxer;

  if (!isUploadSuccess && !props.action) {
    showFileFaxer =
        <div className={props.paper}>
          <Logo className={props.logo} />
          <Typography component="h1" variant="h4" gutterBottom>
            I am a fax machine.
          </Typography>
          <form className={props.form} noValidate>
            {/* Controls fax number input. */}
            <FaxNumberInput
              setFaxNumber={setFaxNumber}
          />
            {/*/!* Controls file selection and validation. This component allows a user to select a file,*/}
            {/*    validates the file, and updates the file information in the app state. *!/*/}
            <FileSelector
                setSelectedFile={setSelectedFile}
            />
            {/* Controls file upload and faxing. */}
            {/*<StyledFileFaxer*/}
            <FileFaxer
                selectedFile={selectedFile}
                faxNumber={faxNumber}
                uploadSuccessHandler={uploadSuccessHandler}
            />
          </form>
        </div>;
  }
  else {
    showFileFaxer = <div></div>;
  }

  return <div>
    <div>
      {showFileFaxer}
      {
        isUploadSuccess ?
            <div>
              file succesfully uploaded!
              you're redirected to paypal for payment!
            </div> :
            <div></div>
      }
    </div>
    <div>
      {
        props.action === "process" ? (
            <div>
              <SendingFax transactionId={props.transactionId}/>
            </div>
        ) : (
            <h3></h3>
        )
      }
    </div>
  </div>
};

FaxMachineApp.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles, { withTheme: true })(FaxMachineApp);
