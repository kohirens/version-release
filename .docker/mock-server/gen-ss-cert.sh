#!/bin/sh

set -e

#
# This is a modified version of, see:
# https://devopscube.com/create-self-signed-certificates-openssl/
#
# See later: https://security.stackexchange.com/questions/74345/provide-subjectaltname-to-openssl-directly-on-the-command-line
# To add SANS, see: https://security.stackexchange.com/a/159537

city='Detroit'
company='Acme'
ec_type='ecdsa'
out_dir='/etc/ssl'
package="$(basename "${0}")"
sans=''
state='Michigan'
verbose=''

# See https://stackoverflow.com/a/7069755/419097
while test $# -gt 0; do
    case "${1}" in
        -h|--help)
            echo "$package - attempt to capture frames"
            echo
            echo "$package [options] application [arguments]"
            echo
            echo "options:"
            echo "-h, --help            show brief help"
            echo "--city                specify a city locale (default: Detroit)"
            echo "--company             specify a city locale (default: Acme)"
            echo "--ec-type             set the encryption type [rsa, ecdsa]"
            echo "--sans                specify Subnet Alternative Name(s) in the form: \"DNS:domain.com, DNS:domain2.com\""
            echo "--state=Michigan      specify a state locale (default: Detroit)"
            exit 0
            ;;
        --city*)
            city=`echo "${1}" | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --company*)
            company=`echo "${1}" | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --ec-type*)
            ec_type=`echo "${1}" | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --sans*)
            sans=`echo "${1}" | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --state*)
            state=`echo "${1}" | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --out-dir*)
            out_dir=`echo "${1}" | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        -v)
            verbose='true'
            shift
            ;;
        *)
            break
            ;;
    esac
done

common_name="${1}"

if [ -z "${common_name}" ]; then
    echo "please enter the Common Name for the certificate as the first argument"
    exit 1
fi

ROOT_CA_KEY="${out_dir}/private/${company}-Root-CA.key"
ROOT_CA_CRT="${out_dir}/certs/${company}-Root-CA.pem"
SSC_KEY="${out_dir}/private/${company}-server.key"
SSC_CSR="${out_dir}/certs/ca-cert-${company}-server.csr"
CSR_CONF="${out_dir}/certs/${company}-cert.conf"
SSC_CRT="${out_dir}/certs/ca-cert-${company}-CA.pem"

OS_CA_FILE="/etc/ssl/certs/ca-certificates.crt"

SUBJ="/C=US/ST=${state}/L=${city}/O=${company}, Inc/OU=Team Ultra/CN=${common_name}"
CA_SUBJ="/C=US/ST=${state}/L=${city}/O=${company}, LLC/CN=${company} Root CA"
SANS="subjectAltName = ${sans}"
DAYS=365
EC_LEVEL=2048

echo  "set up the output directories ${out_dir}"

# Step 1: Make directories to store the certs (in case we are not installing
# where the system keeps the certs.
mkdir -p "${out_dir}"/private "${out_dir}/certs"

MakeKey() {
    ecType="${1}"
    ecLevel="${2}"
    keyFile="${3}"

    if [ "${ecType}" = "rsa" ]; then
        openssl genrsa -out "${keyFile}" "${ecLevel}"
    else
        # add "ECDSA" (X25519 || ≥ secp384r1); see https://safecurves.cr.yp.to/
        # or list ECDSA the supported curves (openssl ecparam -list_curves)
        openssl ecparam -genkey -name secp384r1 -out "${keyFile}"
    fi
}

# Generate a Root CA key and the Root CA pem|crt file.
MakeRootCA() {
    days="${1}"
    ecType="${2}"
    ecLevel="${3}"
    subject="${4}"
    keyFile="${5}"
    crtFile="${6}"

    MakeKey "${ecType}" "${ecLevel}" "${keyFile}"

    openssl req -new -x509 -sha256 -nodes \
        -days "${days}" \
        -subj "${subject}" \
        -key "${keyFile}" \
        -out "${crtFile}"
}

MakeCSR() {
    days="${1}"
    ecType="${2}"
    ecLevel="${3}"
    subject="${4}"
    keyFile="${5}"
    crtFile="${6}"
    sans="${7}"

    MakeKey "${ecType}" "${ecLevel}" "${keyFile}"

    # Removed -x509 as that makes it a certificate instead of a request.
    openssl req -new -sha256 \
        -key "${keyFile}" \
        -subj "${subject}" \
        -addext "${sans}" \
        -out "${crtFile}"
}

# Step 1: Generate a Root CA private key and pem file.
MakeRootCA "${DAYS}" "${ec_type}" "${EC_LEVEL}" "${SUBJ}" "${ROOT_CA_KEY}" "${ROOT_CA_CRT}"

# Step 2: Generation of self-signed(x509) Root CA (PEM-encodings .pem|.crt) based on the Root CA private (.key)
MakeCSR "${DAYS}" "${ec_type}" "${EC_LEVEL}" "${SUBJ}" "${SSC_KEY}" "${SSC_CSR}" "${SANS}"

# Debug: Show whats in the request
if [ "${verbose}" = "true" ]; then
    openssl req -in ${SSC_CSR} -text
fi

echo "signing the CSR"

# Method 2: using the -copy_extensions to copy them from the CSR
openssl x509 -req \
    -days "${DAYS}" \
    -CA "${ROOT_CA_CRT}" \
    -CAkey "${ROOT_CA_KEY}" \
    -CAcreateserial \
    -in "${SSC_CSR}" \
    -out "${SSC_CRT}" \
    -copy_extensions "copyall"

# Debug: Show whats in the cert
if [ "${verbose}" = "true" ]; then
    openssl x509 -in "${SSC_CRT}" -text -noout
fi

# Add the cert to the OS chain, which prevent curl SSL errors inside the
# container. Not necessary, but cool to use -v and see cURL succeed.
echo "add new self-signed certificate to the OS chain of certificates"
cat "${SSC_CRT}" >> "${OS_CA_FILE}"
echo "" >> "${OS_CA_FILE}"

echo
echo "Generated self-signed certificate with private key"
echo "   SSC_CRT: ${SSC_CRT}"
echo "   SSC_KEY: ${SSC_KEY}"
