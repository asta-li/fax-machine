import React from 'react';
import { forwardRef } from 'react';

import TextField from '@material-ui/core/TextField';
import { makeStyles } from '@material-ui/core/styles';

import 'react-phone-number-input/style.css';
import PhoneInput from 'react-phone-number-input';

const MuiFaxInput = forwardRef((props, ref) => {
  // TODO(asta): Make the fax number error handling nicer, then enable.
  // const [faxNumberError, setFaxNumberError] = React.useState('');
  let faxNumberError = validateFaxNumber(props.value);
  faxNumberError = ''
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
      error={faxNumberError !== ''}
      helperText={faxNumberError}
    />
  );
});

export default function FaxInput(props) {
  return (
    <div>
      <PhoneInput
        country={'US'}
        countries={['US']}
        addInternationalOption={false}
        inputComponent={MuiFaxInput}
        onChange={(faxNumber) => props.onChange(faxNumber)}
      />
    </div>
  );
}

// Returns empty string if the fax number is valid.
function validateFaxNumber(faxNumber) {
  let status = '';

  if (!faxNumber) {
    status = 'Error: Please enter a fax number';
    return status;
  }
  
  let numDigits = 0;
  for (let i = 0; i < faxNumber.length; ++i) {
    if (i === 0 && faxNumber[i] === '1') {
      // Skip the US country code +1.
      continue;
    }
    if (faxNumber[i] === ' ') {
      // Skip spaces and parentheses.
      continue;
    }
    if (isNaN(faxNumber[i])) {
      // Skip non-numbers.
      continue;
    }
    numDigits = numDigits + 1;
  }
  
  // Check the fax number length, which is 10 digits plus the US country code (+1)
  if (numDigits !== 10) {
    status = 'Error: Fax number must be 10 digits long';
    return status ;
  }

  return status;
}
