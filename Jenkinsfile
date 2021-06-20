def gv 

pipeline {
  agent any

  tools {
      go 'go-1.16'
  }

  environment {
    WEBHOOK_URL = credentials('DISCORD_CI_WEBHOOK')
    DEPLOY_HOST = credentials('DEPLOY_HOST')
    DEPLOY_USER = credentials('DEPLOY_USER')
    DEPLOY_IDENTITY_FILE = credentials('DEPLOY_IDENTITY_FILE')
    DEPLOY_FOLDER = credentials('DEPLOY_FOLDER')
  }

  stages {
    stage("build") {
      steps {
        sh 'make build'
      }
    }

    stage("deploy") {
      //when {
       // branch 'master'
      //}
      steps {
        script {
          gv = load "deploy.groovy"
          gv.deployApp()
        }
      }
    }
  }
  post {
    failure {
      discordSend description: "Jenkins Pipeline Build", footer: "Aryzona: " + currentBuild.currentResult, link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: env.WEBHOOK_URL
    }
    fixed {
      discordSend description: "Jenkins Pipeline Build", footer: "Aryzona: " + currentBuild.currentResult, link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: env.WEBHOOK_URL
    }
  }
}
