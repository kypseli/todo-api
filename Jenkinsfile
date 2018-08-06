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
          yamlFile 'build-pod.yml'
        }
      }
      steps {
        container('golang') {
          sh 'mkdir -p $GOPATH/src/github.com/kypseli/todo-api'
          sh 'cp -r $WORKSPACE/. $GOPATH/src/github.com/kypseli/todo-api'
          sh 'cd $GOPATH/src/github.com/kypseli/todo-api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o $WORKSPACE/app .'
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
        container('golang') {
          sh 'mkdir -p $GOPATH/src/github.com/kypseli/todo-api'
          sh 'cp -r $WORKSPACE/. $GOPATH/src/github.com/kypseli/todo-api'
          sh 'cd $GOPATH/src/github.com/kypseli/todo-api && go test'
        }
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
        publishEvent simpleEvent('todo-api')
        dockerBuildPush('beedemo/todo-api', "${BUILD_NUMBER}",'./') {
            unstash 'app'
        }
      }
    }
  }
}