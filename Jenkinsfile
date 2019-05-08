pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                script {
                    /**
                     * To be able to access this Jenkins `tool` the https://wiki.jenkins.io/display/JENKINS/Go+Plugin plugin is needed.
                     * With more recent versions of Jenkins the documentation for adding a `go` installation is out of date. To properly
                     * configure a go installation go to the Jenkins tools configuration (Manage Jenkins -> Global Tool Configuration)
                     * find the "Go" and "Go installations" section and click "Add Go". The `name` specified below should
                     * line up with the "Go installation" to be used.
                     */
                    def root = tool name: '1.12.4', type: 'go'

                    /**
                     * Add in GOPATH, GOROOT, GOBIN to the environment and add go to the path for jenkins.
                     * The environment variables are needed for glide to properly work and adding go to the path is required to
                     */
                    withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
                        sh "mkdir -p ${env.WORKSPACE}/go/src"


                        echo "Configuring git to use ssh rather than https for downloading private repositories"
                        // This configures git settings to allow for private repositories to be downloaded with glide.
                        //sh "git config --local url.ssh://git@github.com/.insteadOf https://github.com/"

                        echo "Installing dependencies"
                        //sh "install"
                        sh "/usr/local/bin/dep init"
                        sh "/usr/local/bin/dep ensure --add github.com/aws/aws-sdk-go"
                        sh "/usr/local/bin/dep ensure -add github.com/gruntwork-io/terratest/modules/aws"

                        echo "Building Go Code"
                        //sh "go build ..."

                    }
                }
            }
        }
        stage('Stage') {
            steps {
                script {
                    def root = tool name: '1.12.4', type: 'go'
                    withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {

                        echo "Testing Go Code"
                        /**
                         * Since glide is installed, glide novendor or nv for short can be taken advantage of to list all
                         * files to test sans vendored dependencies.
                         */
                        //sh 'go test -v ./test/*.go'
                    }
                }
            }
        }
    }
    /**
     * This post step will always execute regardless of a build failing or passing to clean up the setting that allows glide
     * to install private dependencies from Github. When using `checkout scm` or the default Jenkins clone step for a git
     * multibranch pipeline this undo change is needed. If the url change is not undone it will fail subsequent builds because the
     * Jenkins Git plugin will fail to clone the repository correctly.
     */
    post {
        always {
            script {
                echo "Undoing config for git to use ssh rather than https for downloading private repositories"
                //sh "git config --local --unset url.ssh://git@github.com/.insteadOf https://github.com/"
            }
        }
    }
}
