#!/bin/bash

function createTnrto() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/tnrto.vaahan.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname tnrto --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-tnrto.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-tnrto.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-tnrto.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-tnrto.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy tnrto's CA cert to tnrto's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/msp/tlscacerts/ca.crt"

  # Copy tnrto's CA cert to tnrto's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/tlsca/tlsca.tnrto.vaahan.com-cert.pem"

  # Copy tnrto's CA cert to tnrto's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/ca"
  cp "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/ca/ca.tnrto.vaahan.com-cert.pem"

  echo "Registering peer0"
  set -x
  fabric-ca-client register --caname tnrto --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering user"
  set -x
  fabric-ca-client register --caname tnrto --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the org admin"
  set -x
  fabric-ca-client register --caname tnrto --id.name tnrtoadmin --id.secret tnrtoadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname tnrto -M "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/msp/config.yaml"

  echo "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname tnrto -M "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls" --enrollment.profile tls --csr.hosts peer0.tnrto.vaahan.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/server.key"

  echo "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname tnrto -M "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/users/User1@tnrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/users/User1@tnrto.vaahan.com/msp/config.yaml"

  echo "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://tnrtoadmin:tnrtoadminpw@localhost:7054 --caname tnrto -M "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/users/Admin@tnrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/tnrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/tnrto.vaahan.com/users/Admin@tnrto.vaahan.com/msp/config.yaml"
}

function createKlrto() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/klrto.vaahan.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/klrto.vaahan.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname klrto --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-klrto.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-klrto.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-klrto.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-klrto.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy klrto's CA cert to klrto's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/msp/tlscacerts/ca.crt"

  # Copy klrto's CA cert to klrto's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/tlsca/tlsca.klrto.vaahan.com-cert.pem"

  # Copy klrto's CA cert to klrto's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/ca"
  cp "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/ca/ca.klrto.vaahan.com-cert.pem"

  echo "Registering peer0"
  set -x
  fabric-ca-client register --caname klrto --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering user"
  set -x
  fabric-ca-client register --caname klrto --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the org admin"
  set -x
  fabric-ca-client register --caname klrto --id.name klrtoadmin --id.secret klrtoadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname klrto -M "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/msp/config.yaml"

  echo "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname klrto -M "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls" --enrollment.profile tls --csr.hosts peer0.klrto.vaahan.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/server.key"

  echo "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname klrto -M "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/users/User1@klrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/users/User1@klrto.vaahan.com/msp/config.yaml"

  echo "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://klrtoadmin:klrtoadminpw@localhost:8054 --caname klrto -M "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/users/Admin@klrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/klrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/klrto.vaahan.com/users/Admin@klrto.vaahan.com/msp/config.yaml"
}

function createKnrto() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/knrto.vaahan.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/knrto.vaahan.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:11054 --caname knrto --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-11054-knrto.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-11054-knrto.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-11054-knrto.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-11054-knrto.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy knrto's CA cert to knrto's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/msp/tlscacerts/ca.crt"

  # Copy knrto's CA cert to knrto's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/tlsca/tlsca.knrto.vaahan.com-cert.pem"

  # Copy knrto's CA cert to knrto's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/ca"
  cp "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem" "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/ca/ca.knrto.vaahan.com-cert.pem"

  echo "Registering peer0"
  set -x
  fabric-ca-client register --caname knrto --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering user"
  set -x
  fabric-ca-client register --caname knrto --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the org admin"
  set -x
  fabric-ca-client register --caname knrto --id.name knrtoadmin --id.secret knrtoadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname knrto -M "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/msp/config.yaml"

  echo "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname knrto -M "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls" --enrollment.profile tls --csr.hosts peer0.knrto.vaahan.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/server.key"

  echo "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:11054 --caname knrto -M "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/users/User1@knrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/users/User1@knrto.vaahan.com/msp/config.yaml"

  echo "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://knrtoadmin:knrtoadminpw@localhost:11054 --caname knrto -M "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/users/Admin@knrto.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/knrto/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/knrto.vaahan.com/users/Admin@knrto.vaahan.com/msp/config.yaml"
}

function createOrderer() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/vaahan.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/vaahan.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/vaahan.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy orderer org's CA cert to orderer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/vaahan.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/vaahan.com/msp/tlscacerts/tlsca.vaahan.com-cert.pem"

  # Copy orderer org's CA cert to orderer org's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/ordererOrganizations/vaahan.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/vaahan.com/tlsca/tlsca.vaahan.com-cert.pem"

  echo "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/vaahan.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/msp/config.yaml"

  echo "Generating the orderer-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls" --enrollment.profile tls --csr.hosts orderer.vaahan.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/server.key"

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/vaahan.com/orderers/orderer.vaahan.com/msp/tlscacerts/tlsca.vaahan.com-cert.pem"

  echo "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/vaahan.com/users/Admin@vaahan.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/vaahan.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/vaahan.com/users/Admin@vaahan.com/msp/config.yaml"
}

createTnrto
createKlrto
createKnrto
createOrderer
