pipeline {
  agent any
  environment {
    registry = "airren/metric-exporter-demo"
    registryCredential="dockerhub"
  }
  stages {
    stage('Building image') {
      steps{
        script {
          dockerImage = docker.build(registry)
        }
      }
    }
    stage('Docker Push') {
      steps {
        script{
          docker.withRegistry( '', registryCredential ) {
            dockerImage.push("v0.0.$BUILD_NUMBER")
            dockerImage.push("latest")
          }
        }
      }
    }
    stage('Remove Unused docker image') {
      steps{
        sh "docker rmi $registry"
        sh "docker rmi $registry:v0.0.$BUILD_NUMBER"
      }
    }
  }
}