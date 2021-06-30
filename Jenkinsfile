def gv 
def FAILED_STAGE

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
        script {
          FAILED_STAGE=env.STAGE_NAME
        }
        sh 'make build'
      }
    }

    stage("test") {
      steps {
        script {
          FAILED_STAGE=env.STAGE_NAME
        }
        sh 'make test'
      }
    }

    stage("deploy") {
      when {
        branch 'master'
      }
      steps {
        script {
          FAILED_STAGE=env.STAGE_NAME
          gv = load "./.jenkins/script.groovy"
          gv.deployApp()
        }
      }
    }
  }
  post {
    failure {
      script {
        gv = load "./.jenkins/script.groovy"
        gv.notifyDiscord(FAILED_STAGE)
      }
    }
    fixed {
      script {
        gv = load "./.jenkins/script.groovy"
        gv.notifyDiscord(FAILED_STAGE)
      }
    }
  }
}
