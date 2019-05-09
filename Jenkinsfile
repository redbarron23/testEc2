node {
    def root = tool name: '1.12.4', type: 'go'

    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/grugrut/golang-ci-jenkins-pipeline") {
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
            
            stage 'Checkout'
        
            git url: 'https://github.com/redbarron23/testEc2.git'
        
            stage 'preTest'
            sh 'go version'
            sh "/usr/local/bin/dep init"
            sh "/usr/local/bin/dep ensure --add github.com/aws/aws-sdk-go"
            sh "/usr/local/bin/dep ensure -add github.com/gruntwork-io/terratest/modules/aws"
            
            stage 'Test'
            //sh 'go vet'
            //sh 'go test -cover'
            
            stage 'Build'
            sh 'go build .'
            
            stage 'Deploy'
            // Do nothing.
        }
    }
}
