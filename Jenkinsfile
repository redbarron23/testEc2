node {
    def root = tool name: '1.12.4', type: 'go'

    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/redbarron23/testEc2") {
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
            parameters {
                string(name: 'ip',
                defaultValue: '172.31.22.136',
                description: 'target host')
                string(name: 'ami',
                defaultValue: 'ami-0fb176954360127fc',
                description: 'latest ami')
            }
         
            
            stage ('Checkout') {
                git url: 'https://github.com/redbarron23/testEc2.git'
            }
        
            stage ('Dependencies') {
                sh 'go version'
                sh "/usr/local/bin/dep init"
                sh "/usr/local/bin/dep ensure --add github.com/aws/aws-sdk-go"
                sh "/usr/local/bin/dep ensure -add github.com/gruntwork-io/terratest/modules/aws"
            }
            
            stage ('Test') {
                sh 'go vet'
                // sh "$HOME/go/bin/golint"
            }
            
            stage ('Build') {
                sh 'go build'
            }
            
            stage ('Deploy') {
                withCredentials([[
                    $class: 'AmazonWebServicesCredentialsBinding',
                    credentialsId: 'tenant-acct-1',
                    accessKeyVariable: 'AWS_ACCESS_KEY_ID',
                    secretKeyVariable: 'AWS_SECRET_ACCESS_KEY'
                ]]) {
                    sh "./testEc2 -ip ${params.ip} -ami ${params.ami}"
                }
            }
        }
    }
}
