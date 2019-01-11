library 'kypseli'
pipeline {
  options { 
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    disableConcurrentBuilds()
    skipDefaultCheckout() 
  }
  agent none
  stages {
    stage('Build') {
      agent {
        kubernetes {
          label 'golang-build'
          yamlFile 'build-pod.yml'
        }
      }
      steps {
        checkout scm
        container('golang') {
          sh """
            mkdir -p /go/src/github.com/kypseli
            ln -s `pwd` /go/src/github.com/kypseli/todo-api
            cd /go/src/github.com/kypseli/todo-api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o ./app .
          """
        }
        stash name: 'app', includes: 'app, Dockerfile'
      }
    }
    stage('Test') {
      agent {
        kubernetes {
          label 'golang-test'
          yamlFile 'test-pod.yml'
        }
      }
      when { 
        beforeAgent true
        not { branch 'pr*' } 
      }
      steps {
        checkout scm
        container('golang') {
          sh 'mkdir -p $GOPATH/src/github.com/kypseli'
          sh 'ln -s `pwd` $GOPATH/src/github.com/kypseli/todo-api'
          sh 'cd $GOPATH/src/github.com/kypseli/todo-api && go test'
        }
        //stash modified deploy yml
        stash name: 'deploy', includes: 'todo-api-deploy.yml'
      }
    }
    stage('Docker Build and Push') {
      //don't need an agent as one is provided in shared pipeline library -> kypseli
      agent none
      when {
        beforeAgent true
        branch 'master'
      }
      steps {
        checkpoint('Post Tests')
        dockerBuildPush('todo-api', "${BUILD_NUMBER}",'./') {
            unstash 'app'
        }
        publishEvent simpleEvent('todo-api')
      }
    }
    stage('Deploy') {
      //don't need an agent as one is provided in shared pipeline library -> kypseli
      agent none
      when {
        beforeAgent true
        branch 'master'
      }
      steps {
        kubeDeploy('todo-api', 'kypseli', "${BUILD_NUMBER}", './todo-api-deploy.yml')
      }
    }
  }
}
