pipeline {
  agent any

  tools {
      go 'go-1.16'
  }

  stages {
    stage("build") {
      steps {
        sh 'make build'
      }
    }

    stage("deploy") {
      when {
        branch 'master'
      }
      steps {
        echo 'deploy? not yet...'
      }
    }
  }
}
