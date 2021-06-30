def deployApp() {
  def remote = [:]
  remote.allowAnyHosts = true
  remote.name = "target"
  remote.host = env.DEPLOY_HOST
  remote.user = env.DEPLOY_USER
  remote.identityFile = env.DEPLOY_IDENTITY_FILE

  sshRemove remote: remote, path: env.DEPLOY_FOLDER + "/aryzona", failOnError: false
  sshPut remote: remote, from: './aryzona', into: env.DEPLOY_FOLDER
  sshCommand remote: remote, command: "ary-restart"
}

def notifyDiscord(FAILED_STAGE) {
  discordSend description: "Took " + currentBuild.durationString, footer: "Aryzona: " + (currentBuild.currentResult == "FAILURE" ? FAILED_STAGE + " " : "" ) + currentBuild.currentResult, link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME + " " + currentBuild.displayName, webhookURL: env.WEBHOOK_URL
}

return this
