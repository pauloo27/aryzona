pipeline {
  agent any

  tools {
      go 'go-1.16'
  }

  environment {
    WEBHOOK_URL = credentials('DISCORD_CI_WEBHOOK')
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
  post {
    always {
      echo env.WEBHOOK_URL
      discordSend description: "Jenkins Pipeline Build", footer: "Aryzona: " + currentBuild.currentResult, link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: env.WEBHOOK_URL
    }
  }
}
