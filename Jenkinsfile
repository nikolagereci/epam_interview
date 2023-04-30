pipeline {
    agent any

    stages {
        stage('Clone') {
            steps {
                checkout scm
            }
        }
        stage('Lint') {
            steps {
                sh 'golangci-lint run -c golangci.yaml'
            }
        }
        stage('Test') {
            steps {
                sh 'go test ./... -count=1'
            }
        }
        stage('Build') {
            steps {
                sh 'CGO_ENABLED=0 go build -C cmd -o ../bin/app'
            }
        }
        stage('Docker Build') {
            steps {
                sh 'sudo docker build -t companies:$(cat VERSION) .'
            }
        }
        stage('Deploy') {
            steps {
                sh 'sudo docker push companies:$(cat VERSION)'
            }
        }
    }
}
