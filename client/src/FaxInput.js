import { forwardRef } from 'react'
import TextField from '@material-ui/core/TextField'
import { makeStyles } from '@material-ui/core/styles'

import 'react-phone-number-input/style.css'
import PhoneInput from 'react-phone-number-input'

const MuiFaxInput = forwardRef((props, ref) => {
  return (
    <TextField
      {...props}
      inputRef={ref}
      fullWidth
      size="small"
      label="Fax Number"
      variant="outlined"
      name="fax-number"
      margin="normal"
      required
      id="fax-number"
      autoComplete="fax-number"
      autoFocus
      error={false}
      helperText={""}
    />
  )
});

/* TODO(asta): Bring back fax number error handling.
Forward refs on higher-level components: https://reactjs.org/docs/forwarding-refs.html
onChange={(event) => this.handleInput(event)}
error={this.state.faxNumberError !== ""}
helperText={this.state.faxNumberError}
*/

export default function FaxInput(props) {
  return (
    <div>
      <PhoneInput
         country={"US"}
         countries={["US"]}
         addInternationalOption={false}
         inputComponent={MuiFaxInput}
         onChange={(faxNumber) => props.onChange(faxNumber)}
      />
    </div>
  );
}
