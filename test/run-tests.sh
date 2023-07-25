#!/usr/bin/env bash

overall_status=0

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
  TF_LOG=DEBUG terraform apply -auto-approve || return 1
  TF_LOG=DEBUG terraform show || return 1
  TF_LOG=DEBUG terraform destroy -auto-approve || return 1
  status_code=$?

  popd >/dev/null
  return $status_code
}

echo "Running tests"

# For each subdirectory in the current directory
for DIRECTORY in */; do
  DIRECTORY="${DIRECTORY%/}"  # Remove trailing slash

  # Run the Python server in the background
  echo "Booting server"
  python3 server.py "$DIRECTORY" &
  server_pid=$!

  # Wait for the server to be up and running
  MAX_WAIT_TIME=5  # Maximum wait time in seconds
  WAIT_INTERVAL=1   # Wait interval in seconds
  ELAPSED_TIME=0

  sleep "$WAIT_INTERVAL"

  until curl -sSf "http://127.0.0.1:8080" ; do
    echo "Waiting for server..."
    sleep "$WAIT_INTERVAL"
    ELAPSED_TIME=$((ELAPSED_TIME + WAIT_INTERVAL))

    if [[ $ELAPSED_TIME -ge $MAX_WAIT_TIME ]]; then
      echo "Server did not start within $MAX_WAIT_TIME seconds. Exiting."
      exit 1
    fi
  done

  echo "Server ready"

  echo "Booted server"

  # Run the test for the current directory
  run_test "$DIRECTORY"
  test_status=$?
  echo "Tests returned exit code $test_status"

  # Update the overall status based on the current test's result
  if [[ $test_status -ne 0 ]]; then
    overall_status=1
  fi

  # Kill the Python server process after the test
  echo "Shutting down the server"
  kill "$server_pid"
  wait "$server_pid"
done

echo "Done testing"
echo "Exiting with code $test_status"

# Exit with the overall status
exit $overall_status
