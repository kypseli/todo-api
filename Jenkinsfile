library 'kypseli'
pipeline {
  options { 
    skipDefaultCheckout()
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    disableConcurrentBuilds()
  }
  agent none
  stages {
          stage('Docker Build and Push') {
          //Dockerfile has a multi-stage build to build app and then build docker image
          //don't need an agent as one is provided in shared pipeline library -> kypseli
          agent none
          steps {
            dockerBuildPush('beedemo/todo-api', "${BUILD_NUMBER}",'./') {
              checkout scm
            }
          }
        }
  }
}