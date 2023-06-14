node {
  withCredentials([
    usernamePassword(credentialsId: 'container-registry', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')
  ]) {
    try {
      currentBuild.result = "SUCCESS"
      env.AWS_ECR_LOGIN = true

      checkout scm

      stage('Deploy to qa') {
        if (env.BRANCH_NAME == 'master') {
          echo './qa.properties going to load.'
          configFileProvider([configFile(fileId: 'qa-env-file', targetLocation: './')]) {
            load './qa.properties'
          }

          echo 'Building and pushing image'
          docker.withRegistry('https://registry.tracified.com/tracified/gateway-qa', 'container-registry') {
            echo 'Building image'
            echo "${env.BUILD_ID}"
            def releaseImage = docker.build("tracified/gateway-qa:${env.BUILD_ID}")
            releaseImage.push()
            releaseImage.push('latest')
          }

          echo 'Deploying image in server'
          withCredentials([[
            $class: 'AmazonWebServicesCredentialsBinding',
            accessKeyVariable: 'AWS_ACCESS_KEY_ID',
            credentialsId: 'aws-ecr-credentials',
            secretKeyVariable: 'AWS_SECRET_ACCESS_KEY'
          ]]) {
            ansiblePlaybook(
              inventory: 'deploy/hosts',
              playbook: 'deploy/qa.yml',
              extras: '-u ubuntu -e GATEWAY_PORT=$GATEWAY_PORT -e DOCKER_USER=${DOCKER_USER} -e DOCKER_PASS=${DOCKER_PASS}'
            )
          }
        }
      }

      stage('Deploy to Staging') {
        if (env.BRANCH_NAME == 'master') {
          echo './staging.properties going to load.'
          configFileProvider([configFile(fileId: 'staging-env-file', targetLocation: './')]) {
            load './staging.properties'
          }

          echo 'Building and pushing image'
          docker.withRegistry('https://registry.tracified.com/tracified/gateway-staging', 'container-registry') {
            echo 'Building image'
            echo "${env.BUILD_ID}"
            def releaseImage = docker.build("tracified/gateway-staging:${env.BUILD_ID}")
            releaseImage.push()
            releaseImage.push('latest')
          }

          echo 'Deploying image in server'
          withCredentials([[
            $class: 'AmazonWebServicesCredentialsBinding',
            accessKeyVariable: 'AWS_ACCESS_KEY_ID',
            credentialsId: 'aws-ecr-credentials',
            secretKeyVariable: 'AWS_SECRET_ACCESS_KEY'
          ]]) {
            ansiblePlaybook(
              inventory: 'deploy/hosts',
              playbook: 'deploy/staging.yml',
              extras: '-u ubuntu -e GATEWAY_PORT=$GATEWAY_PORT -e DOCKER_USER=${DOCKER_USER} -e DOCKER_PASS=${DOCKER_PASS}'
            )
          }
        }
      }

    }
    catch (exc) {
      currentBuild.result = "FAILURE"
      echo 'Something went wrong'
      echo exc.toString()
    }
    finally {
      echo 'All done. Cleaning up docker'
      discordSend(
        description: "Tracified Gateway - ${currentBuild.currentResult}",
        footer: "#${env.BUILD_ID} ${currentBuild.getBuildCauses()[0].shortDescription}",
        link: env.BUILD_URL,
        result: currentBuild.currentResult,
        title: JOB_NAME,
        webhookURL: env.DISCORD_BUILD
      )
      sh 'docker system prune -af'
      echo "Completed pipeline: ${currentBuild.fullDisplayName} with status of ${currentBuild.result}"
    }
  }
}
