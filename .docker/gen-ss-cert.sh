#!/bin/sh

set -e

print_help() {
cat << HDOC
${package} - attempt to capture frames

${package} [options] application [arguments]

options:
-h, --help          show brief help
--ca-file           specify a CA file which to append the certificate CA (default: /etc/ssl/certs/ca-certificates.crt)
--city              specify a city locale (default: Detroit)
--company           specify a company (default: Acme)
--days              number days until the certificate expires(default: 365)
--ec-type           set the encryption type [rsa, ecdsa]
--ec-level          set the encryption level (default: 2048)
--out-dir           where to save the certificate files, will make sub directories (default: /etc/ssl)
--prefix            add an optional prefix to the filenames (ex: myprefix-server.pem), the hyphen will be added
--sans              specify Subnet Alternative Name(s) in the form: "DNS:domain.com, DNS:domain2.com"
--skip-append-to-ca do not append the certificate to a CA file
--state=Michigan    specify a state locale (default: Detroit)
-v                  turn on verbose output
HDOC
}

#
# This is a modified version of, see:
# https://devopscube.com/create-self-signed-certificates-openssl/
#
# See later: https://security.stackexchange.com/questions/74345/provide-subjectaltname-to-openssl-directly-on-the-command-line
# To add SANS, see: https://security.stackexchange.com/a/159537

# For details see:
# https://stackoverflow.com/questions/192249/how-do-i-parse-command-line-arguments-in-bash
getopt --test > /dev/null && true
if [ $? -ne 4 ]; then
    echo 'sorry, getopts --test` failed in this environment'
    exit 1
fi

# Options with a colon must have a value that follows, those without are just booleans.
LONG_OPTS=ca-file:,city:,company:,days,ec-level:,ec-type:,help,out-dir:,prefix:,sans:,skip-append-to-ca,state:
OPTIONS=h,v

PARSED=$(getopt --options=${OPTIONS} --longoptions=${LONG_OPTS} --name "$0" -- "${@}") || exit 1
eval set -- "${PARSED}"


ca_file='/etc/ssl/certs/ca-certificates.crt'
city='Detroit'
company='Acme'
days=365
ec_level=2048
ec_type='ecdsa'
out_dir='/etc/ssl'
prefix=''
sans=''
skip_append_to_ca='no'
state='Michigan'
verbose=''

package="$(basename "${0}")"

# Knowledge improved by reading https://stackoverflow.com/a/7069755/419097
while test $# -gt 0; do
    case "${1}" in
        -h|--help)
            print_help
            exit 0
            ;;
        --ca-file)
            ca_file="${2}"
            shift 2
            ;;
        --city)
            city="${2}"
            shift 2
            ;;
        --company)
            company="${2}"
            shift 2
            ;;
        --days)
            days="${2}"
            shift 2
            ;;
        --ec-type)
            ec_type="${2}"
            shift 2
            ;;
        --ec-level)
            ec_level="${2}"
            shift 2
            ;;
        --sans)
            sans="${2}"
            shift 2
            ;;
        --skip-append-to-ca)
            skip_append_to_ca="yes"
            shift
            ;;
        --state)
            state="${2}"
            shift 2
            ;;
        --out-dir)
            out_dir="${2}"
            shift 2
            ;;
        --prefix)
            prefix="${2}"
            shift 2
            ;;
        -v)
            verbose='true'
            shift
            ;;
        --) shift; break;;
        *) echo "unknown option '${1}'"; exit 1;;
    esac
done

if [ "${verbose}" = "true" ]; then
    echo "configuration:"
    echo "\tca_file=${ca_file}"
    echo "\tcity=${city}"
    echo "\tcompany=${company}"
    echo "\tdays=${days}"
    echo "\tec-type=${ec_type}"
    echo "\tec-level=${ec_level}"
    echo "\tout-dir=${out_dir}"
    echo "\tpackage=${package}"
    echo "\tprefix=${prefix}"
    echo "\tsans=${sans}"
    echo "\tskip-append-to-ca=${skip_append_to_ca}"
    echo "\tstate=${state}"
    echo "\tverbose=${verbose}"
fi

common_name="${1}"
if [ -z "${common_name}" ]; then
    echo "please enter the Common Name for the certificate as the first and only argument"
    exit 1
fi

# Add an optional prefix to the filenames.
if [ -n "${prefix}" ]; then
  prefix="${prefix}-"
fi

ROOT_CA_KEY="${out_dir}/private/${prefix}Root-CA.key"
ROOT_CA_CRT="${out_dir}/certs/${prefix}Root-CA.pem"

SRV_KEY="${out_dir}/private/${prefix}server.key"
SRV_CSR="${out_dir}/certs/${prefix}server.csr"
SRV_CERT="${out_dir}/certs/${prefix}server.pem"

CA_FILE="${ca_file}"

SUBJ="/C=US/ST=${state}/L=${city}/O=${company}, Inc/OU=Team Ultra/CN=${common_name}"
CA_SUBJ="/C=US/ST=${state}/L=${city}/O=${company}, LLC/CN=${company} Root CA"
SANS="subjectAltName = ${sans}"
DAYS=${days}
EC_LEVEL="${ec_level}"

echo  "set up the output directories ${out_dir}"

# Step 1: Make directories to store the certs (in case we are not installing
# where the system keeps the certs.
mkdir -p "${out_dir}"/private "${out_dir}/certs"

make_key() {
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
make_root_ca_key_and_cert() {
    days="${1}"
    ecType="${2}"
    ecLevel="${3}"
    subject="${4}"
    keyFile="${5}"
    crtFile="${6}"

    make_key "${ecType}" "${ecLevel}" "${keyFile}"

    openssl req -new -x509 -sha256 -nodes \
        -days "${days}" \
        -subj "${subject}" \
        -key "${keyFile}" \
        -out "${crtFile}"

  chmod 600 "${keyFile}"
  chmod 600 "${crtFile}"
}

# Generate both a key and CSR for a server
make_key_and_csr() {
    days="${1}"
    ecType="${2}"
    ecLevel="${3}"
    subject="${4}"
    keyFile="${5}"
    crtFile="${6}"
    sans="${7}"

    make_key "${ecType}" "${ecLevel}" "${keyFile}"

    # Removed -x509 as that makes it a certificate instead of a request.
    openssl req -new -sha256 \
        -key "${keyFile}" \
        -subj "${subject}" \
        -addext "${sans}" \
        -out "${crtFile}"

  chmod 644 "${keyFile}"
  chmod 600 "${crtFile}"
}

sign_csr() {
    days="${1}"
    root_ca_crt="${2}"
    root_ca_key="${3}"
    csr="${4}"
    srv_cert="${5}"

    openssl x509 -req \
        -days "${days}" \
        -CA "${root_ca_crt}" \
        -CAkey "${root_ca_key}" \
        -CAcreateserial \
        -in "${csr}" \
        -out "${srv_cert}" \
        -copy_extensions "copyall"

    chmod 644 "${srv_cert}"
}

# Step 1: Generate a Root CA private key and pem file.
make_root_ca_key_and_cert "${DAYS}" "${ec_type}" "${EC_LEVEL}" "${SUBJ}" "${ROOT_CA_KEY}" "${ROOT_CA_CRT}"

# Step 2: Generation of self-signed(x509) Root CA (PEM-encodings .pem|.crt) based on the Root CA private (.key)
make_key_and_csr "${DAYS}" "${ec_type}" "${EC_LEVEL}" "${SUBJ}" "${SRV_KEY}" "${SRV_CSR}" "${SANS}"

# Debug: Show whats in the request
if [ "${verbose}" = "true" ]; then
    openssl req -in ${SRV_CSR} -text
fi

echo "signing the CSR"

# Step 3: Using the -copy_extensions to copy them from the CSR and produce a new signed certificate.
sign_csr "${DAYS}" "${ROOT_CA_CRT}" "${ROOT_CA_KEY}" "${SRV_CSR}" "${SRV_CERT}"

# Debug: Show whats in the cert
if [ "${verbose}" = "true" ]; then
    openssl x509 -in "${SRV_CERT}" -text -noout
fi

if [ "${skip_append_to_ca}" = "no" ]; then
    # Append the certificate to a chain
    echo "add new self-signed certificate to the ${ca_file} chain of certificates"
    cat "${SRV_CERT}" >> "${CA_FILE}"
    echo "" >> "${CA_FILE}"
fi

if [ "${verbose}" = "true" ]; then
    ls -la "${SRV_CERT}"
    ls -la "${SRV_KEY}"
fi

echo
echo "Generated self-signed certificate with private key"
echo "   server cert: ${SRV_CERT}"
echo "   server key: ${SRV_KEY}"
