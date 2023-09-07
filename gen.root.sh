#!/bin/bash

# Define default values
DEFAULT_ROOT_CA_SUBJECT="/C=CN/ST=Jiangsu/L=Wuxi/O=zzz/OU=zzz/CN=zzz Root CA"
DEFAULT_VALID_DAYS=73000
DEFAULT_SERIAL_NUMBER=1000

# Define variables with default values
ROOT_CA_SUBJECT="$DEFAULT_ROOT_CA_SUBJECT"
VALID_DAYS="$DEFAULT_VALID_DAYS"
SERIAL_NUMBER="$DEFAULT_SERIAL_NUMBER"

# Define other variables
SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
OUT_DIR="$SCRIPT_DIR/out"
ROOT_KEY="$OUT_DIR/root.key.pem"
ROOT_CERT="$OUT_DIR/root.crt"
CA_CONFIG="ca.cnf"
INDEX_FILE="$OUT_DIR/index.txt"
ATTR_FILE="$OUT_DIR/index.txt.attr"
SERIAL_FILE="$OUT_DIR/serial"

# Process command line arguments
while getopts "s:d:sn:" opt; do
  case "$opt" in
    -s|--subject) ROOT_CA_SUBJECT="$OPTARG";;
    -d|--valid-days) VALID_DAYS="$OPTARG";;
    -sn|--serial-number) SERIAL_NUMBER="$OPTARG";;
    \?) echo "Usage: $0 [-s root ca subject, defalut: $DEFAULT_ROOT_CA_SUBJECT ] 
         [-d validity days, default: $DEFAULT_VALID_DAYS ] 
         [-sn serial number, default: $DEFAULT_SERIAL_NUMBER ]" >&2
        exit 1;;
  esac
done

# Function to init 'out' directory if it doesn't exist
init_out_directory() {
    if [ ! -d "$OUT_DIR" ]; then
        echo "Creating output structure"
        mkdir -p "$OUT_DIR/newcerts"
        touch "$INDEX_FILE"
        echo "unique_subject = no" > "$ATTR_FILE"
        echo "$SERIAL_NUMBER" > "$SERIAL_FILE"
        echo "Done"
    fi
}

# Function to check if a file exists
file_exists() {
    [ -f "$1" ]
}

# Function to generate the root certificate and key
generate_root_cert() {
    # Check if root certificate already exists
    if file_exists "$ROOT_CERT"; then
        echo "######### Root certificate already exists, skip generating root certificate #########"
        return
    fi

    # Create 'out' directory if it doesn't exist
    [ ! -d "$OUT_DIR" ] && mkdir -p "$OUT_DIR"

    # Generate root cert along with root key
    openssl req -config "$CA_CONFIG" \
        -newkey rsa:4096 -nodes -keyout "$ROOT_KEY" \
        -new -x509 -days "$VALID_DAYS" -out "$ROOT_CERT" \
        -subj "$ROOT_CA_SUBJECT"
    
    echo "Root certificate generated."
}

# Main script execution
cd "$SCRIPT_DIR"
init_out_directory
generate_root_cert
