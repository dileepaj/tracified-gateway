
node {
    def root = tool name: 'Go 1.11.2', type: 'go'
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
            
            sh 'mkdir bin'
            sh 'mkdir src'
            sh 'mkdir src/github.com'
            sh 'mkdir src/github.com/tracified-gateway'
            sh 'ls'
            ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/tracified-gateway") {
              stage('Checkout'){
                  echo 'Checking out SCM'
                  // sh 'cd src'
                  checkout scm
              }  
          
              stage 'preTest'
              sh 'go version'
              sh 'go env'

              stage 'Test'

              
              stage 'Build'
              sh 'pwd'
              sh 'ls -la'
              sh 'go get -u github.com/golang/dep/cmd/dep'
              sh 'dep ensure'
              sh 'ls ./../'
              // sh 'ls ./../github.com@tmp'
              sh 'go build'
              sh 'ls -l'
              // sh "usermod -a -G jenkins jenkins"
              // sh "cd ${GOPATH}src/main/"
              // sh 'ls -l'
              // sh "go get ${GOPATH}src/main/"
              // sh "go build ${GOPATH}src/main/"
              
              stage 'Deploy'
              // Do nothing.
            }
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