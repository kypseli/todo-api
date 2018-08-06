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
                sh 'mkdir -p $GOPATH/src/github.com/kypseli/todo-api'
                sh 'ln -s $WORKSPACE $GOPATH/src/github.com/kypseli/todo-api'
                sh 'ls $GOPATH/src/github.com/kypseli/todo-api'
                sh 'cd $GOPATH/src/github.com/kypseli/todo-api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o app'
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