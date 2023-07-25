#!/usr/bin/env bash

clean_terraform() {
  echo "Removing Terraform files from `pwd`"
  # Remove all Terraform-related files
  rm -rf .terraform
  rm -f .terraform.lock.hcl
  rm -f *.tfstate
  rm -f *.tfstate.backup
  rm -f *.tfstate.backup.*
  rm -f *.tfstate.tmp
  rm -f *.tfstate.tmp.*
  rm -f *.tfvars
  rm -f *.tfvars.json
  rm -f *.tfvars.auto.json
}

run_test() {
  echo "Running test in `pwd`"
  pushd "$1" >/dev/null
  if [[ $? -ne 0 ]]; then
    return 1
  fi

  # Clean up and remove any and all Terraform data for the test
  clean_terraform  
  terraform init
  TF_LOG=DEBUG terraform apply -auto-approve
  TF_LOG=DEBUG terraform show
  TF_LOG=DEBUG terraform destroy -auto-approve
  status_code=$?

  popd >/dev/null
  return $status_code
}

echo "Running tests"

# For each subdirectory in the current directory
for DIRECTORY in */; do
  DIRECTORY="${DIRECTORY%/}"  # Remove trailing slash

  # Run the Python server in the background
  python3 server.py "$DIRECTORY" &
  server_pid=$!

  # Wait for the server to be up and running
  MAX_WAIT_TIME=20  # Maximum wait time in seconds
  WAIT_INTERVAL=1   # Wait interval in seconds
  ELAPSED_TIME=0

  sleep "$WAIT_INTERVAL"

  until curl -sSf "http://localhost:8080" >/dev/null; do
    echo "Booting server"
    sleep "$WAIT_INTERVAL"
    ELAPSED_TIME=$((ELAPSED_TIME + WAIT_INTERVAL))

    if [[ $ELAPSED_TIME -ge $MAX_WAIT_TIME ]]; then
      echo "Server did not start within $MAX_WAIT_TIME seconds. Exiting."
      exit 1
    fi
  done

  echo "Booted server"

  # Run the test for the current directory
  run_test "$DIRECTORY"
  echo "Tests returned exit code $?"

  # Kill the Python server process after the test
  kill "$server_pid"
  wait "$server_pid"
done

echo "Done testing"
