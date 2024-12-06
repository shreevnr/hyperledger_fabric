function swalBasic(data) {
    swal.fire({
        // toast: true,
        icon: `${data.icon}`,
        title: `${data.title}`,
        animation: true,
        position: 'center',
        showConfirmButton: true,
        footer: `${data.footer}`,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', swal.stopTimer)
            toast.addEventListener('mouseleave', swal.resumeTimer)
        }
    });
}

function swalBasicRefresh(data) {
    swal.fire({
        // toast: true,
        icon: `${data.icon}`,
        title: `${data.title}`,
        animation: true,
        position: 'center',
        showConfirmButton: true,
        footer: `${data.footer}`,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', swal.stopTimer)
            toast.addEventListener('mouseleave', swal.resumeTimer)
        }
    }).then(() => {
        location.reload();
    });
}

function reloadWindow() {
    window.location.reload();
}
const createRc = async (event) => {
    event.preventDefault();
    const rcId = document.getElementById('rcId').value;
    const assetType = document.getElementById('assetType').value;
    const make = document.getElementById('make').value;
    const model = document.getElementById('model').value;
    const color = document.getElementById('colour').value;
    const ownerName = document.getElementById('ownerName').value;
    const ownerAadhar = document.getElementById('ownerAadhar').value;
    const engineNumber = document.getElementById('engineNumber').value;
    const insuranceCert = document.getElementById('insuranceCert').value;
    const pollutionCert = document.getElementById('pollutionCert').value;
    const registeredState = document.getElementById('registeredState').value;
    console.log(rcId + make + model + color+registeredState);

    const RegistrationCertificate = {
        rcId: rcId,
        assetType:assetType,
        make: make,
        model: model,
        color: color,
        ownerName:ownerName,
        ownerAadhar:ownerAadhar,
        engineNumber:engineNumber,
        insuranceCert:insuranceCert,
        pollutionCert:pollutionCert,
        registeredState:registeredState,
    };
    if (
        rcId.length == 0 ||
        assetType==0||
        make.length == 0 ||
        model.length == 0 ||
        color.length == 0 ||
        ownerName.length == 0 ||
        ownerAadhar.length == 0||
        engineNumber.length==0||
        insuranceCert.length==0||
        pollutionCert.length==0||
        registeredState.length==0
    ) {
        const data = {
            title: "You might have missed something",
            footer: "Enter all mandatory fields to add a new car",
            icon: "warning"
        }
        swalBasic(data);
    } else {
        try {
            const response = await fetch("/api/rc/create", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(RegistrationCertificate),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            const rcStatus = {
                title: "Success",
                footer: "Added a new RC",
                icon: "success"
            }
            swalBasicRefresh(rcStatus);

        } catch (err) {
            // alert("Error");
            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }

}

const readTNRC = async (event) => {
    event.preventDefault();
    const rcId = document.getElementById("rcID").value;
    if (rcId.length == 0) {
        const data = {
            title: "Enter a valid RC Id",
            footer: "This is a mandatory field",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`/api/tnrto/rc/${rcId}`);
            let responseData = await response.json();
            console.log("response", responseData);
            // alert(JSON.stringify(responseData));
            const dataBuf = JSON.stringify(responseData)
            swal.fire({
                // toast: true,
                icon: `success`,
                title: `Ledger Status of RC with rcId ${rcId} :`,
                animation: false,
                position: 'center',
                html: `<h3>${dataBuf}</h3>`,
                showConfirmButton: true,
                timer: 3000,
                timerProgressBar: true,
                didOpen: (toast) => {
                    toast.addEventListener('mouseenter', swal.stopTimer)
                    toast.addEventListener('mouseleave', swal.resumeTimer)
                }
            })
        } catch (err) {

            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
};

const readKLRC = async (event) => {
    event.preventDefault();
    const rcId = document.getElementById("rcID").value;
    if (rcId.length == 0) {
        const data = {
            title: "Enter a valid RC Id",
            footer: "This is a mandatory field",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`/api/klrto/rc/${rcId}`);
            let responseData = await response.json();
            console.log("response", responseData);
            // alert(JSON.stringify(responseData));
            const dataBuf = JSON.stringify(responseData)
            swal.fire({
                // toast: true,
                icon: `success`,
                title: `Private data of RC with rcId ${rcId} :`,
                animation: false,
                position: 'center',
                html: `<h3>${dataBuf}</h3>`,
                showConfirmButton: true,
                timer: 3000,
                timerProgressBar: true,
                didOpen: (toast) => {
                    toast.addEventListener('mouseenter', swal.stopTimer)
                    toast.addEventListener('mouseleave', swal.resumeTimer)
                }
            })
        } catch (err) {

            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
};


const initiateRCTransfer = async (event) => {
    alert("triggering initiate transfer")
    event.preventDefault();
    const rcId = document.getElementById("rcid").value;
    const fromState = document.getElementById("fromState").value;
    const toState = document.getElementById("toState").value;

    const initiateTrasfer = {
        RCId:rcId,
        FromState:fromState,
        ToState:toState,
    };
    
    if (rcId.length == 0||
        fromState.length==0||
        toState.length==0

    ) {
        const data = {
            title: "Enter a valid RC Id",
            footer: "This is a mandatory field",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`api/rc/initiateTransfer`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(initiateTrasfer),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            const rcStatus = {
                title: "Success",
                footer: "Successfully Initiated RC",
                icon: "success"
            }
            swalBasicRefresh(rcStatus);
        } catch (err) {

            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
};
const approveRCTransfer = async (event) => {
    alert("triggering approve transfer")
    event.preventDefault();
    const rcId = document.getElementById("rcid").value;
    const fromState = document.getElementById("fromState").value;
    const toState = document.getElementById("toState").value;

    const approveTrasfer = {
        RCId:rcId,
        FromState:fromState,
        ToState:toState,
    };
    
    if (rcId.length == 0||
        fromState.length==0||
        toState.length==0

    ) {
        const data = {
            title: "Enter a valid RC Id",
            footer: "This is a mandatory field",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`api/rc/approveTransfer`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(approveTrasfer),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            const rcStatus = {
                title: "Success",
                footer: "Successfully Approved RC",
                icon: "success"
            }
            swalBasicRefresh(rcStatus);
        } catch (err) {

            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
};

const deleteRCTransfer = async (event) => {
    
    event.preventDefault();
    const rcId = document.getElementById("rcid").value;
    const fromState = document.getElementById("fromState").value;
    const toState = document.getElementById("toState").value;

    const deleteTransfer = {
        RCId:rcId,
        FromState:fromState,
        ToState:toState,
    };
    
    if (rcId.length == 0||
        fromState.length==0||
        toState.length==0

    ) {
        const data = {
            title: "Enter a valid RC Id",
            footer: "This is a mandatory field",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`api/rc/deleterc`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(deleteTransfer),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            const rcStatus = {
                title: "Success",
                footer: "Successfully Deleteted transfered RC",
                icon: "success"
            }
            swalBasicRefresh(rcStatus);
        } catch (err) {

            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
};

async function addTransferRC(RCId,assetType,make,model,color,
    ownerName,ownerAadhar,engineNumber,pollutionCert,insuranceCert,registeredState) {
    const TransferrcCertificate = {
        rcId: RCId,
        assetType:assetType,
        make: make,
        model: model,
        color: color,
        ownerName:ownerName,
        ownerAadhar:ownerAadhar,
        engineNumber:engineNumber,
        insuranceCert:insuranceCert,
        pollutionCert:pollutionCert,
        registeredState:registeredState,
    };
    if (
        RCId.length == 0 ||
        assetType==0||
        make.length == 0 ||
        model.length == 0 ||
        color.length == 0 ||
        ownerName.length == 0 ||
        ownerAadhar.length == 0||
        engineNumber.length==0||
        insuranceCert.length==0||
        pollutionCert.length==0||
        registeredState.length==0
    ) {
        const data = {
            title: "You might have missed something",
            footer: "Enter all mandatory fields to add a new RC",
            icon: "warning"
        }
        swalBasic(data);
    } else {
        try {
            const response = await fetch("/api/rc/addrc", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(TransferrcCertificate),
            });
            console.log("RESPONSE: ", response)
            console.log(JSON.stringify(TransferrcCertificate))
            const data = await response.json()
            console.log("DATA: ", data);
            const rcStatus = {
                title: "Success",
                footer: "Added the transferred RC",
                icon: "success"
            }
            swalBasicRefresh(rcStatus);
        } catch (err) {
            // alert("Error");
            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
}
//Method to get the history of an item
function getRCHistory(rcId) {
    console.log("carId====", rcId)
    window.location.href = '/api/rc/history?rcId=' + rcId;
}




function tnrtoCertificates() {
    window.location.href = '/api/tnrto/rc/all';
}

function klrtoCertificates() {
    window.location.href = '/api/klrto/rc/all';
}

function transferredCertificatestn() {
    window.location.href = '/api/tnrto/transferred_rc/all';
}
function transferredCertificateskl() {
    window.location.href = '/api/klrto/transferred_rc/all';
}



async function getEvent() {
    try {
        const response = await fetch("/api/event", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        });
        // console.log("RESPONSE: ", response)
        const data = await response.json()
        // console.log("DATA: ", data);

        const eventsData = data["rcEvent"]
        swal.fire({
            toast: true,
            // icon: `${data.icon}`,
            title: `Event : `,
            animation: false,
            position: 'top-right',
            html: `<h5>${eventsData}</h5>`,
            showConfirmButton: false,
            timer: 5000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', swal.stopTimer)
                toast.addEventListener('mouseleave', swal.resumeTimer)
            }
        })
    } catch (err) {
        swal.fire({
            toast: true,
            icon: `error`,
            title: `Error`,
            animation: false,
            position: 'top-right',
            showConfirmButton: true,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', swal.stopTimer)
                toast.addEventListener('mouseleave', swal.resumeTimer)
            }
        })
        console.log(err);
    }
}




