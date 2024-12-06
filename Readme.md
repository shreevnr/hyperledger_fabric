
git clone "https://github.com/shreevnr/hyperledger_fabric.git"



cd KBA-project_V3/
cd AutoLedger
============================

**To setup the network

run ./startNetwork.sh

-------------------------------------
RUN THE CLIENT

cd ../AutoLedger_app
go run .

launch on browser http://localhost:8080/

//The index page AutoLedger App provides options to Login as any of the state (TamilnaduRTO,KeralaRTO,KarnatakaRTO)

Create new Registration Certicate
=================================

Login as TamilnaduRTO

		Click on [Create New RC]
		Enter Details ['TN01CX2525','RegistrationCertificate','Tata','Tiago','Red','Srimathi','AD111111','EN11111','IC11111','PC111111','tnrto']
		Click Submit

		Success message will be shown after RC get created


		Get the confidential/private details of all RegistrationCertificate belong TamilNaduRTO (tnrtoPDC)
		==================================================================================================
		click on [Show All Certificates]
		
		You will be able to see the list of RCs that belong to TamilnaduRTO
		
		Get the on-chain details of  RegistrationCertificate 
		====================================================
		Click [RC Details]
		
		Enter RC ID [TN01CX2525]
		
		You will be able to see the on-chain details of the given RC.
		
		

Initate Registration Certicate Transfer
========================================

Login as TamilnaduRTO


		click on [Transfer]
		Enter ['TN01CX2525','tnrto','klrto']
		click [Initiate]
		
		you will get the message as Successfully Initiated Transfer.
		
		Get the on-chain details of  RegistrationCertificate 
		====================================================
		Click [RC Details]
		
		Enter RC ID [TN01CX2525]
		
		You will be able to see the on-chain details of the given RC with status as 'Transfer Initiated'
		
		
Approve Registration Certicate Transfer
========================================

Login as KeralaRTO


		click on [Transfer]
		Enter ['TN01CX2525','tnrto','klrto']
		click [Approve]
		
		you will get the message as Successfully Approved Transfer.
		
		Get the on-chain details of  RegistrationCertificate 
		====================================================
		Click [RC Details]
		
		Enter RC ID [TN01CX2525]
		
		You will be able to see the on-chain details of the given RC with status as 'Transfer Approved'
		
		
Delete Registration Certicate
========================================

Login as TamilnaduRTO


		click on [Transfer]
		Enter ['TN01CX2525','tnrto','klrto']
		click [Delete]
		
		you will get the message as Successfully Deleted RC.
		
		Get the on-chain details of  RegistrationCertificate 
		====================================================
		Click [RC Details]
		
		Enter RC ID [TN01CX2525]
		
		You will be able to see the on-chain details of the given RC with status as 'Dis-Owned by tnrto'
		
		Get the confidential/private details of all RegistrationCertificate belong TamilNaduRTO (tnrtoPDC)
		==================================================================================================
		click on [Show All Certificates]
		
		'TN01CX2525' will be not available now in this page
		
		Get the transferred certificates [this page will list out all the RCs that are diswoned by from  state ]
		========================================================================================================
		click on [Transferred Certificates]
		you will be able to see the entry of 'TN01CX2525' 
		

Add Transferred Registration Certicate
========================================
Login as KeralaRTO

		Get the transferred certificates [this page will list out all the RCs that are diswoned by from  state ]
		========================================================================================================
		click on [Transferred Certificates]
		you will be able to see the entry of 'TN01CX2525' 


		Click on [Add] 
		
		You will get success message upon RC get added successfully
		
		Get the on-chain details of  added RegistrationCertificate 
		====================================================
		Click [RC Details]
		
		Enter RC ID [KL04CX2525]
		
		You will be able to see the on-chain details of the given RC with status as 'Active'
		
		Get the confidential/private details of all RegistrationCertificate belong TamilNaduRTO (tnrtoPDC)
		==================================================================================================
		click on [Show All Certificates]
		
		You will be able to see the RC details of KL04CX2525


Get History of Registration Certificate
========================================
go to index page

Click on any entry of the shown RC List.

Click [Get History]

You will be able to see all the logs of the particular RC


to stop the Network

cd ../AutoLedger

./stopNetwork.sh




		
		
		
		
		
		
		
		
		
		

