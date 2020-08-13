pipeline {
  agent none
  environment {
    APP_NAME = 'go-gen'
  }
  stages {
    stage('notify-start') {
      agent { docker 'makeblock/bsdo' }
      when {
        anyOf {
          branch 'master'; branch 'release/*'; tag "release*"; tag "preview*"
        }
      }
      steps {
        script {
          sh "dingtalk link -t ${COMMON_DING_TOKEN} -i \"${APP_NAME} CI start on ${BRANCH_NAME}\" -e \"${GIT_COMMIT}\" -u \"${RUN_DISPLAY_URL}\" -p \"${COMMON_DING_PICTURE}\""
        }
      }
    }
    stage('test') {
      agent {
        docker {
          image 'makeblock/bsdo'
          args '-v /devops/ssh/id_rsa:/root/.ssh/id_rsa -u root'
        }
      }
      when {
        anyOf {
          branch 'master'; branch 'release/*'
        }
      }
      steps {
        sh "make utest"
        stash 'coverage.data'
      }
    }
    stage('sonarqube') {
      agent { label 'docker' }
      when {
        anyOf {
          branch 'master'; branch 'release/*'
        }
      }
      environment {
        scannerHome = tool 'SonarQubeScanner'
        APP_VERSION = sh(script: "cat package.json | grep version | head -1 | awk -F: '{ print \$2 }' | sed 's/[\",]//g' | tr -d '[[:space:]]'", returnStdout: true).trim()
      }
      steps {
        withSonarQubeEnv('SonarQube Server') {
          unstash 'coverage.data'
          script {
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${APP_NAME}:key -Dsonar.projectName=${APP_NAME} -Dsonar.projectVersion=${APP_VERSION} -Dsonar.sources=. -Dsonar.exclusions=**/proto/** -Dsonar.language=go -Dsonar.tests=. -Dsonar.test.inclusions=**/*_test.go -Dsonar.test.exclusions=**/vendor/**,**/proto/** -Dsonar.go.coverage.reportPaths=coverage.data -Dsonar.coverage.dtdVerification=false"
          }
        }
      }
    }
    stage('build') {
      agent {
        docker {
          image 'makeblock/bsdo'
          args '-v /devops/ssh/id_rsa:/root/.ssh/id_rsa -v /var/run/docker.sock:/var/run/docker.sock -v /devops/docker/:/root/.docker -u root'
        }
      }
      when {
        anyOf {
          branch 'master'; branch 'release/*'
        }
      }
      steps {
        script {
          if (env.BRANCH_NAME == 'master') {
            sh "make build-master"
          } else {
            sh "make build-release"
          }
        }
      }
    }
    stage('deploy') {
      agent {
        docker {
          image 'makeblock/bsdo'
          args '-v /devops/kubectl:/root/.kube -u root'
        }
      }
      when {
        anyOf {
          branch 'master'; branch 'release/*'; tag "release*"; tag "preview*"
        }
      }
      steps {
        script {
          if (env.TAG_NAME != null) {
            if (env.TAG_NAME.matches("release(.*)")) {
              sh "make deploy-prod-preview"
              sh "dingtalk link -t ${COMMON_DING_TOKEN} -i \"${APP_NAME} ${BRANCH_NAME} 确认部署\" -e \"${GIT_COMMIT}\" -u \"${RUN_DISPLAY_URL}\" -p \"${COMMON_DING_PICTURE}\""
              input "确认要部署线上环境吗？"
              sh "make deploy-prod"
              checkout scm
              sh "make deploy-we-prod-preview"
              sh "dingtalk link -t ${COMMON_DING_TOKEN} -i \"${APP_NAME} ${BRANCH_NAME} 欧服确认部署\" -e \"${GIT_COMMIT}\" -u \"${RUN_DISPLAY_URL}\" -p \"${COMMON_DING_PICTURE}\""
              input "确认要部署欧服线上环境吗？"
              sh "make deploy-we-prod"              
            } else {
              sh "make deploy-pre"
            }
          } else {
            if (env.BRANCH_NAME == 'master') {
              sh "make deploy-dev"
            } else {
              sh "make deploy-test"
              checkout scm
              sh 'make deploy-we-test'              
            }
          }
        }
      }
    }
    stage('notify-success') {
      agent { docker 'makeblock/bsdo' }
      when {
        anyOf {
          branch 'master'; branch 'release/*'; tag "release*"; tag "preview*"
        }
      }
      steps {
        sh "dingtalk link -t ${COMMON_DING_TOKEN} -i \"${APP_NAME} CI success on ${BRANCH_NAME}\" -e \"${GIT_COMMIT}\" -u \"${RUN_DISPLAY_URL}\" -p \"${COMMON_DING_PICTURE}\""
      }
    }
  }
  post {
    failure {
      node('docker'){
        sh "docker run makeblock/dingtalk dingtalk text -t ${COMMON_DING_TOKEN} -c \"${APP_NAME} CI failed on ${BRANCH_NAME}\" -a"
      }
    }
  }
}