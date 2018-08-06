library 'kypseli'
pipeline {
  options { 
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    disableConcurrentBuilds()
  }
  agent none
  stages {
        stage('Build') {
            agent {
              kubernetes {
                label 'golang-build'
                yamlFile 'jenkins-agent-pod.yml'
              }
            }
            steps {
                container('golang') {
                   sh 'CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o app .'
                }
                stash name: 'app', includes: 'app'
            }
        }
        stage('Docker Build and Push') {
            //don't need an agent as one is provided in shared pipeline library -> kypseli
            agent none
            steps {
                dockerBuildPush('beedemo/todo-api', "${BUILD_NUMBER}",'./') {
                    unstash 'app'
                }
            }
        }
  }
}