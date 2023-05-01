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
                sh 'make lint'
            }
        }
        stage('Test') {
            steps {
                sh 'make unit-test'
            }
        }
        stage('Integration Test') {
            steps {
                sh 'make integration-test'
            }
        }
        stage('Build') {
            steps {
                sh 'make build'
            }
        }
        stage('Docker Build') {
            steps {
                sh 'make dockerbuild'
            }
        }
        stage('Deploy') {
            steps {
                sh 'sudo docker push companies:$(cat VERSION)'
            }
        }
    }
}
