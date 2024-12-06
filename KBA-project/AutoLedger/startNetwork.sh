#!/bin/bash

echo "------------Register the ca admin for each organization—----------------"

docker-compose -f Docker/docker-compose-ca.yaml up -d
sleep 3

sudo chmod -R 777 organizations/

echo "------------Register and enroll the users for each organization—-----------"

chmod +x registerEnroll.sh

./registerEnroll.sh
sleep 3

echo "—-------------Build the infrastructure—-----------------"

docker-compose -f Docker/docker-compose-org.yaml up -d
sleep 3

echo "-------------Generate the genesis block—-------------------------------"

export FABRIC_CFG_PATH=./config

export CHANNEL_NAME=vaahanchannel

configtxgen -profile ThreeRtoChannel -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
sleep 2

echo "------ Create the application channel------"

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/msp/tlscacerts/tlsca.vaahan.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/server.key

osnadmin channel join --channelID $CHANNEL_NAME --config-block ${PWD}/channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2


export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_LOCALMSPID=tnrtoMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/users/Admin@tnrto.vaahan.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/msp/tlscacerts/tlsca.vaahan.com-cert.pem
export TNRTO_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/ca.crt
export KLRTO_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/ca.crt
export KNRTO_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/ca.crt
sleep 2

echo "—---------------Join tnrto peer to the channel—-------------"

echo ${FABRIC_CFG_PATH}
sleep 2
peer channel join -b ${PWD}/channel-artifacts/${CHANNEL_NAME}.block
sleep 3

echo "—-------------tnrto anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.vaahan.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.tnrtoMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.tnrto.vaahan.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb

configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json

echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.vaahan.com --tls --cafile $ORDERER_CA

echo "—---------------package chaincode—-------------"

peer lifecycle chaincode package autoledger.tar.gz --path ${PWD}/../AutoLedger_Chaincode/ --lang golang --label autoledger_1.0
sleep 1

echo "—---------------install chaincode in tnrto peer—-------------"

peer lifecycle chaincode install autoledger.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled
sleep 1

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid autoledger.tar.gz)

echo "—---------------Approve chaincode in tnrto peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.vaahan.com --channelID $CHANNEL_NAME --name autoledger --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --collections-config ../AutoLedger_Chaincode/collection.json --tls --cafile $ORDERER_CA --waitForEvent
sleep 2
echo "—---------------KLRTO—-------------"
export CORE_PEER_LOCALMSPID=klrtoMSP 
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/klrto.vaahan.com/users/Admin@klrto.vaahan.com/msp

echo "—---------------Join klrto peer to the channel—-------------"

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
sleep 1
peer channel list

echo "—-------------KLRTO anchor peer update—-----------"
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.vaahan.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.klrtoMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.klrto.vaahan.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb

configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json

echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.vaahan.com --tls --cafile $ORDERER_CA
sleep 1

echo "—---------------install chaincode in KLRTO peer—-------------"

peer lifecycle chaincode install autoledger.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled

echo "—---------------Approve chaincode in KLRTO peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.vaahan.com --channelID $CHANNEL_NAME --name autoledger --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --collections-config ../AutoLedger_Chaincode/collection.json --tls --cafile $ORDERER_CA --waitForEvent
sleep 1
echo "----------------KNRTO----------------------------------------"

export CORE_PEER_LOCALMSPID=knrtoMSP 
export CORE_PEER_ADDRESS=localhost:11051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/knrto.vaahan.com/users/Admin@knrto.vaahan.com/msp

echo "—---------------Join KNRTO peer to the channel—-------------"

peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
sleep 1
peer channel list

echo "—-------------KNRTO anchor peer update—-----------"

peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.vaahan.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1
cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.knrtoMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.knrto.vaahan.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb

configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json

echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.vaahan.com --tls --cafile $ORDERER_CA

sleep 1

peer channel getinfo -c $CHANNEL_NAME

echo "—---------------install chaincode in KNRTO peer—-------------"
peer lifecycle chaincode install autoledger.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled
echo "—---------------Approve chaincode in KNRTO peer—-------------"
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.vaahan.com --channelID $CHANNEL_NAME --name autoledger --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --collections-config ../AutoLedger_Chaincode/collection.json --tls --cafile $ORDERER_CA --waitForEvent

echo "—---------------Commit chaincode in KNRTO peer—-------------"
peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name autoledger --version 1.0 --sequence 1 --collections-config ../AutoLedger_Chaincode/collection.json --tls --cafile $ORDERER_CA --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.vaahan.com --channelID $CHANNEL_NAME --name autoledger --version 1.0 --sequence 1 --collections-config ../AutoLedger_Chaincode/collection.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $TNRTO_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $KLRTO_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $KNRTO_PEER_TLSROOTCERT

sleep 1
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name autoledger --cafile $ORDERER_CA	


