node {
    def root = tool name: '1.12.4', type: 'go'

    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/redbarron23/testEc2") {
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
            env.AWS_DEFAULT_REGION = "eu-west-2"
            env.AWS_ACCESS_KEY_ID = credentials('jenkins-aws-secret-key-id')
            env.AWS_SECRET_ACCESS_KEY = credentials('jenkins-aws-secret-access-key')
            
            stage 'Checkout'
                git url: 'https://github.com/redbarron23/testEc2.git'
        
            stage 'Dependencies'
                sh 'go version'
                sh "/usr/local/bin/dep init"
                sh "/usr/local/bin/dep ensure --add github.com/aws/aws-sdk-go"
                sh "/usr/local/bin/dep ensure -add github.com/gruntwork-io/terratest/modules/aws"
            
            stage 'Test'
                //sh 'go vet'
                //sh 'go test -cover'
            
            stage 'Build'
                sh 'go build'
                //sh 'ls -l'
            
            stage 'Deploy'
                //withAWS(credentials:'tenant-acct-1', region:'eu-west-2') {
                //    awsIdentity()
                //    sh './testEc2 -ip 172.31.22.136 -ami ami-020ddcd8686c4bc95'
                //}
                //withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', credentialsId: 'AWS_ID']]) {
                //    sh 'aws s3api list-buckets --query "Buckets[].Name"'
                //}
            withCredentials([[
                $class: 'AmazonWebServicesCredentialsBinding',
                credentialsId: 'tenant-acct-1',
                accessKeyVariable: 'AWS_ACCESS_KEY_ID',
                secretKeyVariable: 'AWS_SECRET_ACCESS_KEY'
            ]]) {             
                {
                    sh 'env | sort -u'
                    sh 'aws s3api list-buckets --query "Buckets[].Name"'
                    sh './testEc2 -ip 172.31.22.136 -ami ami-020ddcd8686c4bc95'
                    sh 'aws ec2 describe-instances'
                }
            } // withCredentials
        }
    }
}
