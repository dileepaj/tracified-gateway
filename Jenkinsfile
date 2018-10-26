node {
    try{
        currentBuild.result = "SUCCESS"
        echo 'buildState INPROGRESS'
        
        ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
            withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
                // env.PATH="${GOPATH}/bin:$PATH"

                // Install the desired Go version
                def root = tool name: 'Go 1.10', type: 'go'
                
                stage('Checkout'){
                    echo 'Checking out SCM'
                    checkout scm
                }        
            
                stage('Build'){
                    echo 'Building Executable'
                
                    //Produced binary is $GOPATH/src/cmd/project/project
                    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
                        sh 'go env'
                        sh "cd ./src/main/ && go get && env GOOS=linux GOARCH=arm go build"
                        sh 'chmod u+x main'
                    }
                    // sh "cd $GOPATH/src/main/ && env GOOS=linux GOARCH=arm go build"
                    // sh 'chmod u+x main'
                }
                
                stage('Deploy'){
                    echo 'Deployed'
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