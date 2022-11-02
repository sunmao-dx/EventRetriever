pipeline {
    agent any  	

    stages {
        stage('Hello') {
            steps {
                echo 'Hello World'
            }
        }
        // 这里的hello2 是我加的，就是说明，这是stages下的第二个任务 ,就是在pipeline中加单行注释 用 // 就行
        stage('Hello2') {
            steps {
                echo 'Hello World，i 应该是 可以了 ！！！'
            }
        }
        
        stage('上传文件到docker服务器-并运行'){
            steps {
                sshPublisher(publishers: [sshPublisherDesc(configName: 'docker', transfers: [sshTransfer(cleanRemote: false, excludes: '', execCommand: '''cd /home/jianshan/jk_docker/event-retriever

                echo "环境变量设置"
                export gitee_token=9f0b0f5cb7fe819a87013210024ec84d
                export api_url=http://119.8.116.2:9278/api/dataCache/pushGiteeIssue

                echo "opensource" | sudo -S docker build -t event-retriever:v1 .
                echo "opensource" | sudo -S docker stop event-retriver
                echo "opensource" | sudo -S docker container rm event-retriver
                echo "opensource" | sudo -S docker run -dit -p:9277:8001 --name event-retriver event-retriever:v1
                echo "opensource" | sudo -S docker image prune -f''', execTimeout: 120000, flatten: false, makeEmptyDirs: false, noDefaultExcludes: false, patternSeparator: '[, ]+', remoteDirectory: 'event-retriever', remoteDirectorySDF: false, removePrefix: '', sourceFiles: '**/*')], usePromotionTimestamp: false, useWorkspaceInPromotion: false, verbose: false)])
            }
        }
    }
}
