import axios from "axios";
import CircularProgress from "@material-ui/core/CircularProgress";
import React from "react";


const SendingFax =  (props) => {
    const FAILED = 0;
    const IN_PROGRESS = 1;
    const SUCCESS = 2;

    const [isFaxSuccess, setFaxSuccess] = React.useState(false);
    const [faxStatus, setFaxStatus] = React.useState(null);
    const [faxId, setFaxId] = React.useState(null);

    const config = {
        headers: { 'content-type': 'multipart/form-data' }
    };

    React.useEffect(() => {
        console.log("component updated, making an API call");

        const formData = new FormData();
        formData.append('transactionId', props.transactionId);

        setFaxStatus(IN_PROGRESS);
        // try to send fax. on the backend will validate that the person has indeed paid.
        axios.post('/api/fax', formData, config)
            .then((response) => {

                // TODO: get fax id out and poll
                console.log('Received successful fax response', response);

                setFaxId(response.data.FaxId);
                setFaxSuccess(true);
                setFaxStatus(SUCCESS);
            })
            .catch((error) => {
                console.log(error);
                setFaxStatus(FAILED)

            });

    }, []);

    let maindiv;

    if (faxStatus === SUCCESS) {
        maindiv = <div>"yay sucesssfully faxed,  you're all done here! and your fax id is" <h3>{faxId}</h3></div>
    } else if (faxStatus === IN_PROGRESS) {
        maindiv = <div>
            faxing... you can navigate away from this page, you will receive an email notification upon successful fax
            <CircularProgress />
        </div>
    } else if (faxStatus === FAILED) {
        maindiv = <div> "sorry somthing went wrong</div>
    }

    return (
        <div>
            {maindiv}

        </div>
    )

}

export default SendingFax;