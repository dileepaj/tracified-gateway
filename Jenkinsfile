node {
    try{
        notifyBuild('STARTED')
        echo 'buildState INPROGRESS'
        
        ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
            withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
                env.PATH="${GOPATH}/bin:$PATH"
                
                stage('Checkout'){
                    echo 'Checking out SCM'
                    checkout scm
                }        
            
                stage('Build'){
                    echo 'Building Executable'
                
                    //Produced binary is $GOPATH/src/cmd/project/project
                    sh 'cd $GOPATH/src/main/ && env GOOS=linux GOARCH=arm go build'
                    sh 'chmod u+x main'
                }
                
                stage('Deploy'){

                }
            }
        }
    }catch (e) {
        // If there was an exception thrown, the build failed
        currentBuild.result = "FAILED"
        
        echo 'buildState FAILED'

    } finally {
        // Success or failure, always send notifications
        notifyBuild(currentBuild.result)
        
        def bs = currentBuild.result ?: 'SUCCESSFUL'
        if(bs == 'SUCCESSFUL'){
            echo 'buildState SUCCESSFUL' 
        }
    }
}