#!/bin/bash
EXPECTED_OUTPUT=truffle-expected-output.json
RESULT_OUTPUT=truffle-result-output.json
PRUNED_OUTPUT=truffle-pruned-output.json

doGithubWorkflowProcessing () {
  if [ "" != "$GITHUB_ACTION" ] ; then
    echo "Running within github actions... processing output file"
    # running in a github action, output results for next workflow action
    make -f ../Makefile github-action-openzeppelin || exit 1

    if [ ! -f $EXPECTED_OUTPUT ] ; then
      echo "Expected output not found -" $EXPECTED_OUTPUT
      doGithubWorkflowProcessingResult=1
      return
    fi

    if [ -e $RESULT_OUTPUT ] ; then
      echo Successfully copied output results from docker container
      make -f ../Makefile truffle-parser-docker
    else
      echo "Failed to find output results in docker container"
      doGithubWorkflowProcessingResult=-1
      return
    fi
    
    doGithubWorkflowProcessingResult=$?
  else
    echo "Not running within github actions, skipping processing of output results"
  fi
}
cleanupDocker () {
  echo Shutting down docker-compose containers
  docker-compose -f docker-compose-openzeppelin.yml -p ci kill
  docker-compose -f docker-compose-openzeppelin.yml -p ci rm -f
}
trap 'cleanupDocker ; echo "Tests Failed For Unexpected Reasons"' HUP INT QUIT PIPE TERM
docker-compose -p ci -f docker-compose-openzeppelin.yml build && docker-compose -p ci -f docker-compose-openzeppelin.yml up -d
if [ $? -ne 0 ] ; then
  echo "Docker Compose Failed"
  exit 1
fi
# docker logs qtum_seeded_testchain -f&
# docker logs ci_charon_1 -f&
docker logs ci_openzeppelin_1 -f&
EXIT_CODE=`docker wait ci_openzeppelin_1`
echo "Processing openzeppelin test results with exit code of:" $EXIT_CODE
doGithubWorkflowProcessingResult=$EXIT_CODE

if [ -e $RESULT_OUTPUT ] ; then
  echo "Deleting existing output results"
  rm $RESULT_OUTPUT
fi

echo "Copying output results from docker container to local filesystem"
CONTAINER=ci_openzeppelin_1 INPUT=/openzeppelin-contracts/output.json make -f ../Makefile truffle-parser-extract-result-docker || exit 1
echo "Successfully copied output results from docker container to local filesystem"

doGithubWorkflowProcessing
EXIT_CODE=$doGithubWorkflowProcessingResult
if [ -z ${EXIT_CODE+z} ] || [ -z ${EXIT_CODE} ] || ([ "0" != "$EXIT_CODE" ] && [ "" != "$EXIT_CODE" ]) ; then
  # TODO: is it even worth outputting the logs? they will overflow the actual results
  # these logs are so large we can't print them out into github actions
  # docker logs qtum_seeded_testchain
  # docker logs ci_charon_1
  # docker logs ci_openzeppelin_1
  echo "Tests Failed - Exit Code: $EXIT_CODE (truffle exit code indicates how many tests failed)"
else
  echo "Tests Passed"
fi
cleanupDocker
exit $EXIT_CODE
