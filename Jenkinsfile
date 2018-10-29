node {
    def root = tool name: 'Go 1.8', type: 'go'
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
            
            stage 'Checkout'
        
            stage 'preTest'
            sh 'go version'
            sh 'go env'

            stage 'Test'

            
            stage 'Build'
            sh 'ls -l'
            // sh 'cd src/main'
            // sh 'go get'
            // sh 'go build .'
            
            stage 'Deploy'
            // Do nothing.
        }
    }
}
// node {
//     try{
//         currentBuild.result = "SUCCESS"
//         echo 'buildState INPROGRESS'

        
        
//         ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
//             withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}", "PATH+GO=${root}/bin"]) {
//                 env.PATH="${GOPATH}/bin:$PATH"
 
//                 // Install the desired Go version
//                 //  def root = tool name: 'Go1.8', type: 'go'
//                 // sh "${root}"
//                 // sh 'go version'

//                 stage('Checkout'){
//                     echo 'Checking out SCM'
//                     checkout scm

//                 }        
            
           

//                 stage('Build'){
//                     echo 'Building Executable'
                
//                     // Produced binary is $GOPATH/src/cmd/project/project
//                     // withEnv(["GOROOT=${root}/bin", "PATH+GO=${root}/bin"]) {
//                         sh 'go env'
//                         // sh "cd $GOPATH/src/main/ && go get && go build"
//                         sh "go get && go build"
//                         sh 'chmod u+x main'
//                     // }
                   
//                 }
                
//                 stage('Deploy'){
//                     echo 'Deployed'
//                 }
//             }
//         }
//     }catch (e) {
//         // If there was an exception thrown, the build failed
//         currentBuild.result = "FAILED"
        
//         echo 'buildState FAILED'

//     } finally {
//         // Success or failure, always send notifications
//         notifyBuild(currentBuild.result)
        
//         def bs = currentBuild.result ?: 'SUCCESSFUL'
//         if(bs == 'SUCCESSFUL'){
//             echo 'buildState SUCCESSFUL' 
//         }
//     }
// }