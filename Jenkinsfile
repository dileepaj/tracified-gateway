node {

  try {
    currentBuild.result = "SUCCESS"
    env.AWS_ECR_LOGIN = true

    checkout scm

    docker.image('golang:latest').inside('-u root') {
        stage('Setup') {
                   echo 'Setting up environment'
                   echo scm.branches
                   echo env.BRANCH_NAME
                   if (env.BRANCH_NAME == 'staging') {
                        configFileProvider([configFile(fileId: '4e86e233-697c-4371-aad3-dae58c04a62a', targetLocation: './')]) {
                        load './staging.properties'
                        }
                   }
        }

    }
    stage('Deploy to Staging') {
              if (env.BRANCH_NAME == 'staging') {
                echo 'Building and pushing image'
                docker.withRegistry('https://453230908534.dkr.ecr.ap-south-1.amazonaws.com/tracified/gateway-staging', 'ecr:ap-south-1:aws-ecr-credentials') {
                  echo 'Building image'
                  echo ${env.BUILD_ID}
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
                  ansiblePlaybook inventory: 'deploy/hosts', playbook: 'deploy/staging.yml', extras: '-u ubuntu'
                }
              }
    }



    }
    catch (exc) {
        currentBuild.result = "FAILURE"
        echo 'Something went wrong'
    }
    finally {
        echo 'All done. Cleaning up docker'
        sh 'docker system prune -af'
        echo "Completed pipeline: ${currentBuild.fullDisplayName} with status of ${currentBuild.result}"
    }
}