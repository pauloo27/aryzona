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

return this
