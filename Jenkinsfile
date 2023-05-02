node {

  try {
    currentBuild.result = "SUCCESS"
    env.AWS_ECR_LOGIN = true

    checkout scm

    docker.image('golang:1.15.6').inside('-u root') {
        stage('Setup') {
                   echo 'Setting up environment'
                   echo env.ENVIRONMENT
                   if (env.ENVIRONMENT == 'staging') {
                        echo './staging.properties going to load.'
                        configFileProvider([configFile(fileId: 'staging-env-file', targetLocation: './')]) {
                        load './staging.properties'                        
                        }
                     echo 'load properties done.'
                   }
                   if (env.ENVIRONMENT == 'qa') {
                        echo './qa.properties going to load.'
                        configFileProvider([configFile(fileId: 'qa-env-file', targetLocation: './')]) {
                        load './qa.properties'
                        }
                      echo 'load properties done.'
                   }
                   if (env.ENVIRONMENT == 'production') {
                        echo './production.properties going to load.'
                        configFileProvider([configFile(fileId: 'production-env-file', targetLocation: './')]) {
                        load './production.properties'
                        }
                      echo 'load properties done.'
                   }
        }

    }
    stage('Deploy to Staging') {
              echo env.ENVIRONMENT
              if (env.ENVIRONMENT == 'staging') {
                echo 'Building and pushing image'
                docker.withRegistry('https://453230908534.dkr.ecr.ap-south-1.amazonaws.com/tracified/gateway-staging', 'ecr:ap-south-1:aws-ecr-credentials') {
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
                  ansiblePlaybook inventory: 'deploy/hosts', playbook: 'deploy/staging.yml', extras: '-u ubuntu -e GATEWAY_PORT=$GATEWAY_PORT'
                }
              }
    }
        stage('Deploy to qa') {
              echo env.ENVIRONMENT
              if (env.ENVIRONMENT == 'qa') {
                echo 'Building and pushing image'
                docker.withRegistry('https://453230908534.dkr.ecr.ap-south-1.amazonaws.com/tracified/gateway-qa', 'ecr:ap-south-1:aws-ecr-credentials') {
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
                  ansiblePlaybook inventory: 'deploy/hosts', playbook: 'deploy/qa.yml', extras: '-u ubuntu -e GATEWAY_PORT=$GATEWAY_PORT'
                }
              }
    }
        stage('Deploy to production') {
              echo env.ENVIRONMENT
              if (env.ENVIRONMENT == 'production') {
                echo 'Building and pushing image'
                docker.withRegistry('https://453230908534.dkr.ecr.ap-south-1.amazonaws.com/tracified/gateway-prod', 'ecr:ap-south-1:aws-ecr-credentials') {
                  echo 'Building image'
                  echo "${env.BUILD_ID}"                  
                  def releaseImage = docker.build("tracified/gateway-prod:${env.BUILD_ID}")
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
                  // ansiblePlaybook inventory: 'deploy/hosts', playbook: 'deploy/production.yml', extras: '-u ubuntu -e GATEWAY_PORT=$GATEWAY_PORT'
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
        sh 'docker system prune -af'
        echo "Completed pipeline: ${currentBuild.fullDisplayName} with status of ${currentBuild.result}"
    }
}
